#!/usr/bin/env python3
# Copyright (c) HashiCorp, Inc.
"""
Terraform State Migration Tool for VastData Provider

Migrates an existing tfstate file to work with VastData provider v3.0.

Strategy:
  1. Read the customer's existing .tfstate file (which may contain both
     VastData and non-VastData resources.
  2. Remove all VastData resources (type starts with "vastdata_") from
     the state, leaving non-VastData resources untouched.
  3. Write the cleaned state as terraform.tfstate in the working directory.
  4. Run `terraform init` to initialise the v3.0 provider.
  5. Run `VASTDATA_MIGRATE_MODE=1 terraform apply --auto-approve` so the
     provider re-reads every VastData resource from the cluster and
     re-populates the state with the new schema.

The result is a terraform.tfstate that is fully compatible with the v3.0
provider and still contains all non-VastData resources from the original
state.
"""

import argparse
import json
import os
import shutil
import subprocess
import sys
from pathlib import Path
from typing import Dict, List, Tuple

VERSION = "1.0.0"


# ---------------------------------------------------------------------------
# Logging helpers
# ---------------------------------------------------------------------------

class Colors:
    """ANSI color codes for terminal output"""
    RED = '\033[0;31m'
    GREEN = '\033[0;32m'
    YELLOW = '\033[1;33m'
    BLUE = '\033[0;34m'
    NC = '\033[0m'  # No Color


def log_info(msg: str) -> None:
    print(f"{Colors.BLUE}[INFO]{Colors.NC} {msg}")


def log_success(msg: str) -> None:
    print(f"{Colors.GREEN}[SUCCESS]{Colors.NC} {msg}")


def log_warning(msg: str) -> None:
    print(f"{Colors.YELLOW}[WARNING]{Colors.NC} {msg}")


def log_error(msg: str) -> None:
    print(f"{Colors.RED}[ERROR]{Colors.NC} {msg}")


# ---------------------------------------------------------------------------
# State manipulation
# ---------------------------------------------------------------------------

def parse_tfstate(state_file: str) -> Dict:
    """Parse a Terraform state file and return the JSON dict."""
    log_info(f"Parsing state file: {state_file}")
    try:
        with open(state_file, 'r') as fh:
            state = json.load(fh)
        log_success(
            f"Parsed state file (version: {state.get('version', 'unknown')}, "
            f"serial: {state.get('serial', 'unknown')})"
        )
        return state
    except json.JSONDecodeError as exc:
        log_error(f"Invalid JSON in state file: {exc}")
        sys.exit(1)
    except FileNotFoundError:
        log_error(f"State file not found: {state_file}")
        sys.exit(1)


def separate_vast_resources(resources: List[Dict]) -> Tuple[List[Dict], List[Dict]]:
    """Separate VastData resources from other-provider resources.

    Returns:
        (vast_resources, non_vast_resources)
    """
    vast_resources: List[Dict] = []
    non_vast_resources: List[Dict] = []

    for resource in resources:
        resource_type = resource.get('type', '')
        if resource_type.startswith('vastdata_'):
            vast_resources.append(resource)
        else:
            non_vast_resources.append(resource)

    log_info(
        f"Separated resources: {len(vast_resources)} VastData, "
        f"{len(non_vast_resources)} non-VastData"
    )
    return vast_resources, non_vast_resources


def strip_vast_resources(state: Dict) -> Tuple[Dict, int, int]:
    """Remove all VastData resources from a state dict.

    Returns:
        (cleaned_state, vast_count, non_vast_count)
    """
    if 'resources' not in state:
        log_warning("State file has no 'resources' key — nothing to strip.")
        return state, 0, 0

    all_resources = state['resources']
    vast, non_vast = separate_vast_resources(all_resources)

    cleaned = dict(state)
    cleaned['resources'] = non_vast
    # Bump the serial so Terraform sees this as a newer state.
    cleaned['serial'] = cleaned.get('serial', 0) + 1

    return cleaned, len(vast), len(non_vast)


def write_state(state: Dict, path: str) -> None:
    """Write a state dict to a JSON file."""
    with open(path, 'w') as fh:
        json.dump(state, fh, indent=2)
    log_success(f"Wrote cleaned state to {path}")


# ---------------------------------------------------------------------------
# Terraform execution helpers
# ---------------------------------------------------------------------------

def run_terraform_init(workdir: str) -> None:
    """Run `terraform init` in *workdir*.

    If provider development overrides are active (``dev_overrides`` in the
    CLI config), ``terraform init`` is expected to fail with a registry
    lookup error.  In that case we log a warning and continue — the
    override already points Terraform at the local provider binary, so
    ``init`` is not required.
    """
    log_info("Running: terraform init")
    result = subprocess.run(
        ['terraform', 'init'],
        cwd=workdir,
        capture_output=True,
        text=True,
    )
    combined_output = (result.stdout or '') + (result.stderr or '')

    if result.returncode != 0:
        if 'provider development overrides are in effect' in combined_output.lower() \
                or 'dev_overrides' in combined_output.lower():
            log_warning(
                "terraform init failed but provider dev_overrides are in effect. "
                "Skipping init — this is expected during development."
            )
            return
        # Print captured output so the user can see what went wrong
        if result.stdout:
            print(result.stdout)
        if result.stderr:
            print(result.stderr, file=sys.stderr)
        log_error("terraform init failed.")
        sys.exit(1)
    log_success("terraform init completed.")


def get_vastdata_provider_version(workdir: str) -> str:
    """Return the installed VastData provider version string, or '' on failure.

    Reads ``terraform version -json`` which parses the lock file and does
    not require a live cluster connection.

    Example ``provider_selections`` payload::

        {
          "registry.terraform.io/vast-data/vastdata": "3.1.1",
          ...
        }
    """
    result = subprocess.run(
        ['terraform', 'version', '-json'],
        cwd=workdir,
        capture_output=True,
        text=True,
    )
    if result.returncode != 0:
        return ''
    try:
        data = json.loads(result.stdout)
        selections = data.get('provider_selections', {})
        for key, ver in selections.items():
            if 'vastdata' in key.lower():
                return str(ver)
    except (json.JSONDecodeError, AttributeError, KeyError):
        pass
    return ''


def verify_migrate_mode_support(workdir: str) -> None:
    """Verify that the installed VastData provider supports VASTDATA_MIGRATE_MODE.

    Primary check: inspect the provider version via ``terraform version -json``.
    The feature was introduced in v3.0, so any v3.0+ provider is safe.

    Fallback check: run ``terraform plan`` with VASTDATA_MIGRATE_MODE=1 and
    TF_LOG=WARN and look for the banner the provider emits in Configure().
    This fallback is retained for dev-override / unusual install scenarios, but
    is NOT used as the primary gate because Configure() only runs after a
    successful cluster connection — an unreachable cluster would cause a false
    negative even with the correct provider version installed.
    """
    log_info("Verifying that the VastData provider supports VASTDATA_MIGRATE_MODE …")

    # --- Primary: version-based check (no cluster connection required) ---
    provider_version = get_vastdata_provider_version(workdir)
    if provider_version:
        log_info(f"Detected VastData provider version: {provider_version}")
        try:
            major = int(provider_version.split('.')[0])
        except (ValueError, IndexError):
            major = -1

        if major >= 3:
            log_success(
                f"Provider v{provider_version} supports VASTDATA_MIGRATE_MODE — safe to proceed."
            )
            return

        if major >= 0:
            # Version parsed successfully but it is too old.
            log_error(
                f"VastData provider v{provider_version} does NOT support VASTDATA_MIGRATE_MODE.\n"
                "Running 'terraform apply' with this provider could CREATE or DESTROY\n"
                "real resources — this is NOT safe for migration.\n\n"
                "Please upgrade the VastData Terraform provider to v3.0+ and try again."
            )
            sys.exit(1)

        log_warning(
            f"Could not parse provider version '{provider_version}' — "
            "falling back to banner detection."
        )
    else:
        log_warning(
            "Could not determine the installed provider version from "
            "'terraform version -json' — falling back to banner detection."
        )

    # --- Fallback: banner detection via terraform plan ---
    log_info("Running 'terraform plan' with VASTDATA_MIGRATE_MODE=1 to detect banner …")
    env = os.environ.copy()
    env['VASTDATA_MIGRATE_MODE'] = '1'
    env['TF_LOG'] = 'WARN'

    result = subprocess.run(
        ['terraform', 'plan', '-input=false'],
        cwd=workdir,
        capture_output=True,
        env=env,
        text=True,
    )

    combined = (result.stdout or '') + (result.stderr or '')

    if 'VASTDATA MIGRATE MODE ENABLED' not in combined:
        log_error(
            "The installed VastData provider does NOT support VASTDATA_MIGRATE_MODE.\n"
            "Running 'terraform apply' with an older provider could CREATE or DESTROY\n"
            "real resources — this is NOT safe for migration.\n\n"
            "Please upgrade the VastData Terraform provider to v3.0+ and try again."
        )
        sys.exit(1)

    log_success("Provider supports VASTDATA_MIGRATE_MODE — safe to proceed.")


def run_terraform_apply_migrate(workdir: str) -> None:
    """Run `VASTDATA_MIGRATE_MODE=1 terraform apply -auto-approve`."""
    log_info("Running: VASTDATA_MIGRATE_MODE=1 terraform apply -auto-approve")
    env = os.environ.copy()
    env['VASTDATA_MIGRATE_MODE'] = '1'
    result = subprocess.run(
        ['terraform', 'apply', '-auto-approve'],
        cwd=workdir,
        env=env,
        text=True,
    )
    if result.returncode != 0:
        log_error("terraform apply (migrate mode) failed.")
        sys.exit(1)
    log_success("terraform apply (migrate mode) completed.")


# ---------------------------------------------------------------------------
# Main
# ---------------------------------------------------------------------------

def main() -> None:
    parser = argparse.ArgumentParser(
        description=(
            "Migrate an existing Terraform state file to VastData provider v3.0.\n\n"
            "This script strips all vastdata_* resources from the provided state file,\n"
            "then runs `VASTDATA_MIGRATE_MODE=1 terraform apply -auto-approve` to\n"
            "re-populate them from the VAST cluster using the v3.0 provider schema.\n"
            "Non-VastData resources (AWS, GCP, etc.) are preserved untouched."
        ),
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog="""
Examples:
  # Basic usage — run from the directory containing converted .tf files
  cd /path/to/output/configs
  python3 state_migration.py /path/to/existing/terraform.tfstate

  # Dry-run: only strip the state, do not run terraform
  python3 state_migration.py /path/to/terraform.tfstate --dry-run
""",
    )

    parser.add_argument('--version', action='version', version=f'%(prog)s {VERSION}')

    parser.add_argument(
        'state_file',
        help='Path to the existing terraform.tfstate file to migrate.',
    )
    parser.add_argument(
        '--dry-run',
        action='store_true',
        help=(
            'Only strip VastData resources and write the cleaned state. '
            'Do NOT run terraform init / apply.'
        ),
    )

    args = parser.parse_args()

    # Resolve paths
    state_file = os.path.abspath(args.state_file)
    workdir = os.getcwd()

    # Validate inputs
    if not os.path.isfile(state_file):
        log_error(f"State file does not exist: {state_file}")
        sys.exit(1)

    # Check that .tf files exist in workdir
    tf_files = list(Path(workdir).glob('*.tf'))
    if not tf_files:
        log_error(
            f"No .tf files found in working directory: {workdir}\n"
            "Please run this script from the directory containing your "
            "converted v3.0 .tf files."
        )
        sys.exit(1)

    # ---- Banner ----
    print()
    log_info("=" * 70)
    log_info("VastData Terraform State Migration Tool v" + VERSION)
    log_info("=" * 70)
    log_info(f"State file : {state_file}")
    log_info(f"Work dir   : {workdir}")
    log_info(f"Dry run    : {args.dry_run}")
    log_info("=" * 70)
    print()

    # Step 1: Parse the existing state
    state = parse_tfstate(state_file)

    # Step 2: Strip VastData resources
    cleaned_state, vast_count, non_vast_count = strip_vast_resources(state)

    if vast_count == 0:
        log_warning("No VastData resources found in the state file — nothing to migrate.")
        sys.exit(0)

    log_info(f"Removed {vast_count} VastData resource(s) from state.")
    log_info(f"Preserved {non_vast_count} non-VastData resource(s).")

    # Step 3: Write cleaned state into workdir
    dest_state = os.path.join(workdir, 'terraform.tfstate')

    # Safety: back up any existing state in the workdir
    if os.path.exists(dest_state):
        backup = dest_state + '.backup-pre-migration'
        log_warning(f"Backing up existing state to {backup}")
        shutil.copy2(dest_state, backup)

    write_state(cleaned_state, dest_state)

    if args.dry_run:
        log_success("Dry run complete. Cleaned state written; terraform was NOT executed.")
        print()
        log_info("To continue manually:")
        log_info(f"  cd {workdir}")
        log_info("  terraform init")
        log_info("  VASTDATA_MIGRATE_MODE=1 terraform apply -auto-approve")
        sys.exit(0)

    # Step 4: terraform init
    run_terraform_init(workdir)

    # Step 5: verify the provider supports VASTDATA_MIGRATE_MODE
    verify_migrate_mode_support(workdir)

    # Step 6: terraform apply in migrate mode
    run_terraform_apply_migrate(workdir)

    # ---- Summary ----
    print()
    log_success("=" * 70)
    log_success("State migration completed successfully!")
    log_success("=" * 70)
    log_info(f"  VastData resources re-imported : {vast_count}")
    log_info(f"  Non-VastData resources kept    : {non_vast_count}")
    log_info(f"  New state file                 : {dest_state}")
    print()
    log_info("Next steps:")
    log_info("  1. Verify: terraform plan   (should show no changes)")
    log_info("  2. If using a remote backend, upload the new state file.")
    log_info("  3. Continue using the v3.0 provider normally (no VASTDATA_MIGRATE_MODE).")
    print()


if __name__ == '__main__':
    try:
        main()
    except KeyboardInterrupt:
        print("\n\nInterrupted by user.")
        sys.exit(130)
    except Exception as exc:
        log_error(f"Unexpected error: {exc}")
        import traceback
        traceback.print_exc()
        sys.exit(1)
