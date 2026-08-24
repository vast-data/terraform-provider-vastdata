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

VERSION = "1.2.5"

# Reuse the same rename map as migration_script.py so import addresses match
# converted .tf configuration resource types.
try:
    from migration_script import resource_type_rename_map
except ImportError:  # pragma: no cover - running as installed module
    resource_type_rename_map = {}


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

def _is_data_source(resource: Dict) -> bool:
    """Return True when *resource* is a Terraform data source state entry."""
    return resource.get('mode') == 'data'


def _is_vastdata_data_source(resource: Dict) -> bool:
    """Return True when *resource* is a vastdata_* data source in state."""
    return (
        resource.get('type', '').startswith('vastdata_')
        and _is_data_source(resource)
    )


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


# Resource types that cannot be re-imported from the cluster in MigrateMode
# because they have no stable cluster-side identity to search by.
# These are preserved in state verbatim rather than stripped and re-adopted.
NON_IMPORTABLE_RESOURCE_TYPES = {
    "vastdata_user_key",
}

# Resources that require a composite import ID (pipe-separated field values)
# instead of a plain numeric ID.  The list of fields maps to the provider's
# ImportFields hint.  The import string is built as field1|field2|…
#
# If any field value is missing the script falls back to the numeric ``id``.
RESOURCE_IMPORT_FIELDS: Dict[str, List[str]] = {
    "vastdata_view":        ["path", "tenant_name"],
    "vastdata_view_policy": ["name", "tenant_name"],
}

# Attributes the v3 provider cannot repopulate from API GET (create-only, etc.).
# After terraform import these are copied from the original v1 state when null.
PRESERVE_ATTRIBUTES_FROM_V1: Dict[str, List[str]] = {
    "vastdata_view": ["create_dir"],
    "vastdata_quota": ["is_physical_quota"],
    "vastdata_quota_group": ["is_physical_quota"],
    "vastdata_qos_policy": ["attached_users"],
}

# JSON string attributes rewritten to indented form after import so migrated state
# matches typical heredoc formatting in converted .tf files.
PRETTIFY_JSON_ATTRIBUTES: Dict[str, List[str]] = {
    "vastdata_s3_policy": ["policy"],
}


def _map_resource_type(rtype: str) -> str:
    """Map v1/v2 state resource types to v3 configuration resource types."""
    return resource_type_rename_map.get(rtype, rtype)


def _build_terraform_address(
    module: str,
    rtype: str,
    rname: str,
    *,
    index_key=None,
    instance_index: int = None,
) -> str:
    """Build a full Terraform resource address including module path."""
    if index_key is not None:
        resource_addr = f"{rtype}.{rname}[{json.dumps(index_key)}]"
    elif instance_index is not None:
        resource_addr = f"{rtype}.{rname}[{instance_index}]"
    else:
        resource_addr = f"{rtype}.{rname}"

    if module:
        return f"{module}.{resource_addr}"
    return resource_addr


def _build_tenant_id_map(state: Dict) -> Dict[int, str]:
    """Return a {tenant_id: tenant_name} mapping built from any resource in the
    state that carries both ``tenant_id`` and ``tenant_name`` attributes.
    Falls back to "default" for tenant_id=1 if no explicit mapping is found.
    """
    mapping: Dict[int, str] = {1: "default"}  # safe default
    for resource in state.get('resources', []):
        for instance in resource.get('instances', []):
            attrs = instance.get('attributes', {})
            tid = attrs.get('tenant_id')
            tname = attrs.get('tenant_name')
            if tid is not None and tname:
                try:
                    mapping[int(tid)] = str(tname)
                except (ValueError, TypeError):
                    pass
    return mapping


def _build_import_id(rtype: str, attrs: Dict, tenant_map: Dict[int, str]) -> str:
    """Return the correct import ID string for *rtype*.

    For resources listed in RESOURCE_IMPORT_FIELDS the ID is a pipe-separated
    composite (e.g. ``/my-bucket|default``).  If any required field is absent
    the function falls back to the plain numeric ``id``.

    For all other resource types the plain numeric ``id`` is returned.
    """
    import_fields = RESOURCE_IMPORT_FIELDS.get(rtype)
    if import_fields is None:
        return str(attrs.get('id', ''))

    values = []
    for field in import_fields:
        val = attrs.get(field)
        if val is None and field == 'tenant_name':
            # Derive tenant_name from tenant_id when not stored directly.
            tenant_id = attrs.get('tenant_id')
            if tenant_id is not None:
                val = tenant_map.get(int(tenant_id), 'default')
        if not val:
            # Missing field — fall back to plain numeric id.
            log_warning(
                f"  Cannot build composite import ID for {rtype}: "
                f"field '{field}' is missing. Falling back to numeric id."
            )
            return str(attrs.get('id', ''))
        values.append(str(val))

    return '|'.join(values)


def extract_resource_ids(state: Dict) -> Dict[str, str]:
    """Extract {terraform_address: import_id} for importable VastData resources.

    For most resources the import ID is the numeric ``id`` stored in the state,
    which lets ``terraform import`` perform a direct by-ID API call
    (``GET /<endpoint>/<id>/``) — O(1) and fast.

    For resources listed in RESOURCE_IMPORT_FIELDS a composite pipe-separated
    ID is built from the relevant field values (e.g. ``/bucket-path|default``
    for vastdata_view).  This matches the format expected by the provider's
    ImportState implementation for those resource types.

    Resources without an ``id`` attribute (or those in NON_IMPORTABLE_RESOURCE_TYPES)
    are skipped.

    Import addresses use the v3 resource type names from ``migration_script.resource_type_rename_map``
    so they match converted ``.tf`` configuration (e.g. ``vastdata_administators_roles``
    in state → ``vastdata_administrator_role`` for import).

    Module paths from the state ``module`` field are included so resources declared
    inside modules import correctly (e.g. ``module.views["env1"].vastdata_view.this``).
    """
    tenant_map = _build_tenant_id_map(state)
    result: Dict[str, str] = {}
    for resource in state.get('resources', []):
        rtype = resource.get('type', '')
        if _is_data_source(resource):
            # Data sources are refreshed on read — never terraform import-ed.
            continue
        mapped_rtype = _map_resource_type(rtype)
        rname = resource.get('name', '')
        module = resource.get('module', '')
        if not rtype.startswith('vastdata_'):
            continue
        if mapped_rtype in NON_IMPORTABLE_RESOURCE_TYPES:
            continue
        instances = resource.get('instances', [])
        if not instances:
            continue
        for instance in instances:
            attrs = instance.get('attributes', {})
            if attrs.get('id') is None:
                continue
            # Prefer the numeric ID from the existing state. It enables a direct
            # by-ID API call and avoids ambiguous composite keys (e.g. view path="/").
            import_id = str(attrs['id'])
            if not import_id:
                continue
            index_key = instance.get('index_key')
            if index_key is not None:
                addr = _build_terraform_address(
                    module, mapped_rtype, rname, index_key=index_key,
                )
            elif len(instances) > 1:
                addr = _build_terraform_address(
                    module,
                    mapped_rtype,
                    rname,
                    instance_index=instances.index(instance),
                )
            else:
                addr = _build_terraform_address(module, mapped_rtype, rname)
            result[addr] = import_id
    return result


def extract_v1_attributes(state: Dict) -> Dict[str, Dict]:
    """Return {terraform_address: attributes} for importable VastData resources.

    Used to preserve create-only and other non-round-tripping attributes after
    ``terraform import`` refreshes state from the cluster API.
    """
    result: Dict[str, Dict] = {}
    for resource in state.get('resources', []):
        rtype = resource.get('type', '')
        if _is_data_source(resource):
            continue
        mapped_rtype = _map_resource_type(rtype)
        rname = resource.get('name', '')
        module = resource.get('module', '')
        if not rtype.startswith('vastdata_'):
            continue
        if mapped_rtype in NON_IMPORTABLE_RESOURCE_TYPES:
            continue
        instances = resource.get('instances', [])
        if not instances:
            continue
        for instance in instances:
            attrs = instance.get('attributes', {})
            if attrs.get('id') is None:
                continue
            index_key = instance.get('index_key')
            if index_key is not None:
                addr = _build_terraform_address(
                    module, mapped_rtype, rname, index_key=index_key,
                )
            elif len(instances) > 1:
                addr = _build_terraform_address(
                    module,
                    mapped_rtype,
                    rname,
                    instance_index=instances.index(instance),
                )
            else:
                addr = _build_terraform_address(module, mapped_rtype, rname)
            result[addr] = attrs
    return result


def patch_imported_attributes(
    workdir: str,
    v1_attrs: Dict[str, Dict],
    imported_addresses: Dict[str, str],
) -> int:
    """Merge selected v1 attributes into state when import left them null.

    Returns the number of attribute values patched.
    """
    state_path = os.path.join(workdir, 'terraform.tfstate')
    state = parse_tfstate(state_path)
    patched = 0

    for resource in state.get('resources', []):
        if _is_data_source(resource):
            continue
        rtype = resource.get('type', '')
        mapped_rtype = _map_resource_type(rtype)
        preserve_fields = PRESERVE_ATTRIBUTES_FROM_V1.get(mapped_rtype)
        if not preserve_fields:
            continue

        rname = resource.get('name', '')
        module = resource.get('module', '')
        instances = resource.get('instances', [])
        for instance in instances:
            index_key = instance.get('index_key')
            if index_key is not None:
                addr = _build_terraform_address(
                    module, mapped_rtype, rname, index_key=index_key,
                )
            elif len(instances) > 1:
                addr = _build_terraform_address(
                    module,
                    mapped_rtype,
                    rname,
                    instance_index=instances.index(instance),
                )
            else:
                addr = _build_terraform_address(module, mapped_rtype, rname)

            if addr not in imported_addresses:
                continue

            v1 = v1_attrs.get(addr, {})
            attrs = instance.setdefault('attributes', {})
            for field in preserve_fields:
                if field not in v1 or v1[field] is None:
                    continue
                if attrs.get(field) is not None:
                    continue
                attrs[field] = v1[field]
                patched += 1
                log_info(f"  Preserved {addr}.{field} = {v1[field]!r} from v1 state")

    if patched:
        write_state(state, state_path)
        log_success(f"Preserved {patched} attribute value(s) from v1 state.")
    else:
        log_info("No create-only attributes needed patching.")

    return patched


def _json_strings_equivalent(a, b) -> bool:
    """Return True when two values contain equivalent JSON documents."""
    if a == b:
        return True
    if not isinstance(a, str) or not isinstance(b, str):
        return False
    try:
        return json.loads(a) == json.loads(b)
    except (json.JSONDecodeError, TypeError):
        return False


def _canonicalize_json_string(value: str) -> str:
    """Return compact canonical JSON matching provider CanonicalJSONString."""
    parsed = json.loads(value)
    return json.dumps(parsed, separators=(',', ':'), sort_keys=True)


def prettify_json_fields(
    workdir: str,
    v1_attrs: Dict[str, Dict],
    imported_addresses: Dict[str, str],
) -> int:
    """Canonicalize configured JSON attributes in migrated state.

    The v3 provider stores JSON policy fields in compact canonical form so
    Terraform plan does not drift on whitespace.  This step rewrites imported
    state to that form when the semantic content is unchanged.
    """
    state_path = os.path.join(workdir, 'terraform.tfstate')
    state = parse_tfstate(state_path)
    prettified = 0

    for resource in state.get('resources', []):
        if _is_data_source(resource):
            continue
        rtype = resource.get('type', '')
        mapped_rtype = _map_resource_type(rtype)
        json_fields = PRETTIFY_JSON_ATTRIBUTES.get(mapped_rtype)
        if not json_fields:
            continue

        rname = resource.get('name', '')
        module = resource.get('module', '')
        instances = resource.get('instances', [])
        for instance in instances:
            index_key = instance.get('index_key')
            if index_key is not None:
                addr = _build_terraform_address(
                    module, mapped_rtype, rname, index_key=index_key,
                )
            elif len(instances) > 1:
                addr = _build_terraform_address(
                    module,
                    mapped_rtype,
                    rname,
                    instance_index=instances.index(instance),
                )
            else:
                addr = _build_terraform_address(module, mapped_rtype, rname)

            if addr not in imported_addresses:
                continue

            v1 = v1_attrs.get(addr, {})
            attrs = instance.setdefault('attributes', {})
            for field in json_fields:
                current = attrs.get(field)
                if not isinstance(current, str) or not current.strip():
                    continue

                try:
                    new_value = _canonicalize_json_string(current)
                except json.JSONDecodeError:
                    log_warning(
                        f"  Skipping {addr}.{field}: value is not valid JSON"
                    )
                    continue

                if new_value != current:
                    attrs[field] = new_value
                    prettified += 1
                    log_info(f"  Prettified {addr}.{field}")

    if prettified:
        write_state(state, state_path)
        log_success(f"Prettified {prettified} JSON attribute value(s) in state.")
    else:
        log_info("No JSON attributes needed prettifying.")

    return prettified


def separate_vast_resources(resources: List[Dict]) -> Tuple[List[Dict], List[Dict], List[Dict]]:
    """Separate resources into three buckets:

    - vast_importable:     vastdata_* resources that MigrateMode will re-adopt
                           from the cluster.
    - vast_preserved:      vastdata_* resources that cannot be re-imported
                           (e.g. vastdata_user_key) and are carried over verbatim.
    - non_vast:            resources belonging to other providers (AWS, GCP, …)
                           that are kept untouched.

    Returns:
        (vast_importable, vast_preserved, non_vast)
    """
    vast_importable: List[Dict] = []
    vast_preserved: List[Dict] = []
    non_vast: List[Dict] = []

    for resource in resources:
        resource_type = resource.get('type', '')
        if not resource_type.startswith('vastdata_'):
            non_vast.append(resource)
        elif _is_vastdata_data_source(resource):
            log_info(
                f"Preserving VastData data source "
                f"'{resource.get('name', '?')}' (type: {resource_type}) "
                f"— refreshed on read, not re-imported."
            )
            vast_preserved.append(resource)
        elif resource_type in NON_IMPORTABLE_RESOURCE_TYPES:
            log_info(
                f"Preserving non-importable VastData resource "
                f"'{resource.get('name', '?')}' (type: {resource_type}) "
                f"— cannot be re-adopted from cluster."
            )
            vast_preserved.append(resource)
        else:
            vast_importable.append(resource)

    log_info(
        f"Classified resources: "
        f"{len(vast_importable)} VastData to re-import, "
        f"{len(vast_preserved)} VastData preserved (non-importable), "
        f"{len(non_vast)} non-VastData kept as-is"
    )
    return vast_importable, vast_preserved, non_vast


def strip_vast_resources(state: Dict) -> Tuple[Dict, int, int, int]:
    """Remove importable VastData resources from state; preserve the rest.

    Non-importable VastData resources and non-VastData resources are both
    written into the cleaned state so they survive the migration untouched.

    Returns:
        (cleaned_state, vast_importable_count, vast_preserved_count, non_vast_count)
    """
    if 'resources' not in state:
        log_warning("State file has no 'resources' key — nothing to strip.")
        return state, 0, 0, 0

    all_resources = state['resources']
    vast_importable, vast_preserved, non_vast = separate_vast_resources(all_resources)

    cleaned = dict(state)
    cleaned['resources'] = vast_preserved + non_vast
    # Bump the serial so Terraform sees this as a newer state.
    cleaned['serial'] = cleaned.get('serial', 0) + 1

    return cleaned, len(vast_importable), len(vast_preserved), len(non_vast)


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
    # The dev-override warning only appears in plan/apply/init output, not in
    # 'terraform version'. So we run plan once and check for EITHER the banner
    # (emitted by the v3 provider's Configure()) OR the dev_overrides notice
    # (which means the locally-built binary IS the v3 provider we just built).
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

    if 'VASTDATA MIGRATE MODE ENABLED' in combined:
        log_success("Provider supports VASTDATA_MIGRATE_MODE — safe to proceed.")
        return

    # When dev_overrides are active the lock file is never written, so the
    # version cannot be detected from 'terraform version -json'. The plan
    # output itself contains the dev_overrides notice — if it is present we
    # trust the locally built binary as the v3 provider.
    if 'provider development overrides are in effect' in combined.lower() \
            or 'dev_overrides' in combined.lower():
        log_warning(
            "Provider development overrides are in effect. "
            "Trusting the locally built binary as v3 — skipping banner check."
        )
        log_success("Provider supports VASTDATA_MIGRATE_MODE — safe to proceed.")
        return

    log_error(
        "The installed VastData provider does NOT support VASTDATA_MIGRATE_MODE.\n"
        "Running 'terraform apply' with an older provider could CREATE or DESTROY\n"
        "real resources — this is NOT safe for migration.\n\n"
        "Please upgrade the VastData Terraform provider to v3.0+ and try again."
    )
    sys.exit(1)


def release_state_lock(workdir: str) -> None:
    """Force-release a stale Terraform state lock if one exists.

    A stale lock (``OperationTypeInvalid``) is left behind when a previous
    Terraform command exits abnormally.  Without releasing it every subsequent
    ``terraform import`` call fails immediately.

    This function runs ``terraform force-unlock -force <lock-id>`` for any
    lock whose operation type indicates it is stale (i.e. not an active
    plan/apply/import).
    """
    result = subprocess.run(
        ['terraform', 'state', 'list'],
        cwd=workdir,
        capture_output=True,
        text=True,
    )
    # We intentionally ignore the exit code here — even a locked state returns
    # the lock info in stderr which we parse below.
    combined = (result.stdout or '') + (result.stderr or '')

    import re
    lock_id_match = re.search(r'ID:\s+([0-9a-f-]{36})', combined)
    if not lock_id_match:
        return  # no lock found

    lock_id = lock_id_match.group(1)
    # Only force-unlock stale locks (OperationTypeInvalid) to avoid accidentally
    # breaking an active in-progress operation.
    if 'OperationTypeInvalid' not in combined and 'resource temporarily unavailable' not in combined:
        return

    log_warning(f"Detected stale state lock (ID: {lock_id}). Force-unlocking …")
    unlock_result = subprocess.run(
        ['terraform', 'force-unlock', '-force', lock_id],
        cwd=workdir,
        capture_output=True,
        text=True,
    )
    if unlock_result.returncode == 0:
        log_success(f"State lock {lock_id} released.")
    else:
        log_warning(
            f"Could not auto-release lock {lock_id}: {unlock_result.stderr.strip()}\n"
            f"Run manually:  terraform force-unlock -force {lock_id}"
        )


def run_terraform_import_all(workdir: str, resource_ids: Dict[str, str]) -> None:
    """Import each resource using ``terraform import``.

    For most resources the import token is the numeric ``id`` which triggers a
    direct by-ID API call (``GET /<endpoint>/<id>/``) — O(1) and fast.

    For resources whose provider ImportState expects a composite key
    (e.g. ``vastdata_view`` → ``/path|tenant_name``) the script builds the
    correct pipe-separated string from the original state attributes.

    ``-lock=false`` is passed to every import call so that a stale lock file
    from a previous failed run does not block re-execution.

    Resources that fail to import are reported; the script exits with a
    non-zero code if any import fails.
    """
    if not resource_ids:
        log_warning("No importable VastData resources found in the original state.")
        return

    # Try to release any stale lock before starting.
    release_state_lock(workdir)

    log_info(f"Importing {len(resource_ids)} resource(s) using IDs from original state …")
    failed: List[Tuple[str, str]] = []

    for address, import_id in resource_ids.items():
        log_info(f"  terraform import {address} {import_id!r}")
        result = subprocess.run(
            ['terraform', 'import', '-lock=false', address, import_id],
            cwd=workdir,
            text=True,
        )
        if result.returncode != 0:
            log_error(f"  ✗ Failed to import {address} (import_id={import_id!r})")
            failed.append((address, import_id))
        else:
            log_success(f"  ✓ {address}")

    if failed:
        log_error(f"{len(failed)} resource(s) could not be imported:")
        for addr, import_id in failed:
            log_error(f"    terraform import {addr} {import_id!r}")
        sys.exit(1)

    log_success(f"All {len(resource_ids)} resource(s) imported successfully.")


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

    # Step 1b: Extract resource IDs BEFORE stripping (needed for terraform import)
    resource_ids = extract_resource_ids(state)
    v1_attributes = extract_v1_attributes(state)
    if resource_ids:
        log_info(f"Extracted {len(resource_ids)} resource ID(s) for direct import:")
        for addr, import_id in resource_ids.items():
            log_info(f"  {addr} → {import_id!r}")
    else:
        log_warning("No importable resource IDs found in the state file.")

    # Step 2: Strip VastData resources
    cleaned_state, vast_importable_count, vast_preserved_count, non_vast_count = strip_vast_resources(state)

    if vast_importable_count == 0 and vast_preserved_count == 0:
        log_warning("No VastData resources found in the state file — nothing to migrate.")
        sys.exit(0)

    log_info(f"Removed {vast_importable_count} VastData resource(s) for re-import.")
    log_info(f"Preserved {vast_preserved_count} VastData resource(s) (non-importable, carried over verbatim).")
    log_info(f"Preserved {non_vast_count} non-VastData resource(s) (other providers).")

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
        for addr, import_id in resource_ids.items():
            log_info(f"  terraform import -lock=false {addr} {import_id!r}")
        sys.exit(0)

    # Step 4: terraform init
    run_terraform_init(workdir)

    # Step 5: Import each resource directly by its numeric ID.
    # terraform import <address> <id> calls GetByIdWithContext which is a single
    # direct API call — no slow name-based list+filter across the whole cluster.
    run_terraform_import_all(workdir, resource_ids)

    # Step 6: Restore create-only attributes that import cannot read from the API.
    patch_imported_attributes(workdir, v1_attributes, resource_ids)

    # Step 7: Canonicalize JSON policy fields so state matches provider storage.
    prettify_json_fields(workdir, v1_attributes, resource_ids)

    # ---- Summary ----
    print()
    log_success("=" * 70)
    log_success("State migration completed successfully!")
    log_success("=" * 70)
    log_info(f"  VastData resources imported    : {len(resource_ids)}")
    log_info(f"  VastData resources preserved   : {vast_preserved_count}")
    log_info(f"  Non-VastData resources kept    : {non_vast_count}")
    log_info(f"  New state file                 : {dest_state}")
    print()
    log_info("Next steps:")
    log_info("  1. Verify: terraform plan   (should show no changes)")
    log_info("  2. If using a remote backend, upload the new state file.")
    log_info("  3. Continue using the v3.0 provider normally.")
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
