# Copyright (c) HashiCorp, Inc.

"""
Targeted regressions for write-only / credential-state bugs (TERF-318/319/321).

Unlike test_e2e_resources.py (apply/import/destroy smoke), these assert Terraform
state contents and destroy behaviour that the Crossplane suite cannot see.
"""

from __future__ import annotations

import json
import random
import string
from pathlib import Path

import pytest

VALID_SECRET = "abcdefghijklmnopqrstuvwxyz0123456789/+AB"
VALID_SECRET_2 = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789/+cd"
VALID_ACCESS = "ABCDEFGHIJ0123456789"
VALID_ACCESS_2 = "ZYXWVUTSRQ9876543210"


def _uniq() -> str:
    return "".join(random.choices(string.ascii_lowercase, k=4))


def _write_tf(workdir, name: str, body: str) -> None:
    (Path(workdir) / name).write_text(body)


def _attr(workdir, resource_type: str, name: str) -> dict:
    state = json.loads((Path(workdir) / "terraform.tfstate").read_text())
    for res in state.get("resources", []):
        if (
            res.get("mode") == "managed"
            and res.get("type") == resource_type
            and res.get("name") == name
        ):
            return res["instances"][0]["attributes"]
    raise AssertionError(f"resource {resource_type}.{name} not in state")


@pytest.fixture
def tf_ready(terraform_cmd):
    terraform_cmd["init", "-input=false"]()
    return terraform_cmd


@pytest.mark.e2e
def test_terf318_generated_and_assigned_user_key_keep_secret_in_state(tf_ready, terraform_workdir):
    """Both cluster-generated and user-supplied secrets must remain readable in state."""
    u = _uniq()
    _write_tf(
        terraform_workdir,
        "terf318.tf",
        f"""
data "vastdata_tenant" "t" {{
  name = "default"
}}

resource "vastdata_user" "u_gen" {{
  name              = "vastdb{u}_g"
  uid               = {random.randint(20000, 60000)}
  local_provider_id = 1
}}

resource "vastdata_user" "u_asg" {{
  name              = "vastdb{u}_a"
  uid               = {random.randint(20000, 60000)}
  local_provider_id = 1
}}

resource "vastdata_user_key" "generated" {{
  user_id   = vastdata_user.u_gen.id
  tenant_id = data.vastdata_tenant.t.id
  enabled   = true
}}

resource "vastdata_user_key" "assigned" {{
  user_id    = vastdata_user.u_asg.id
  tenant_id  = data.vastdata_tenant.t.id
  access_key = "{VALID_ACCESS}"
  secret_key = "{VALID_SECRET}"
  enabled    = true
}}
""",
    )
    try:
        tf_ready["apply", "-auto-approve", "-input=false", "-lock=false"]()

        gen = _attr(terraform_workdir, "vastdata_user_key", "generated")
        assert gen.get("access_key"), "generated access_key missing"
        assert gen.get("secret_key"), "generated secret_key must stay in state (TERF-318)"
        assert not gen.get("secret_key_md5"), "generated path must not invent secret_key_md5"

        assigned = _attr(terraform_workdir, "vastdata_user_key", "assigned")
        assert assigned.get("access_key") == VALID_ACCESS
        assert assigned.get("secret_key") == VALID_SECRET, "assigned secret_key must stay in state"
        assert not assigned.get("secret_key_md5"), "secret_key_md5 was dropped for user_key"

        rc, stdout, stderr = tf_ready[
            "plan", "-detailed-exitcode", "-input=false", "-lock=false"
        ].run(retcode=None)
        assert rc == 0, f"unexpected drift after apply:\n{stdout}\n{stderr}"
    finally:
        tf_ready["destroy", "-auto-approve", "-input=false", "-lock=false"].run(retcode=None)


@pytest.mark.e2e
def test_terf321_user_key_rotation_keeps_usable_credentials(tf_ready, terraform_workdir):
    """Rotation must succeed and leave the new secret in state."""
    u = _uniq()
    uid = random.randint(20000, 60000)
    base = f"""
data "vastdata_tenant" "t" {{
  name = "default"
}}

resource "vastdata_user" "u" {{
  name              = "vastdb{u}_r"
  uid               = {uid}
  local_provider_id = 1
}}
"""
    _write_tf(
        terraform_workdir,
        "terf321.tf",
        base
        + f"""
resource "vastdata_user_key" "k" {{
  user_id    = vastdata_user.u.id
  tenant_id  = data.vastdata_tenant.t.id
  access_key = "{VALID_ACCESS}"
  secret_key = "{VALID_SECRET}"
  enabled    = true
}}
""",
    )
    try:
        tf_ready["apply", "-auto-approve", "-input=false", "-lock=false"]()
        before = _attr(terraform_workdir, "vastdata_user_key", "k")
        assert before["secret_key"] == VALID_SECRET

        _write_tf(
            terraform_workdir,
            "terf321.tf",
            base
            + f"""
resource "vastdata_user_key" "k" {{
  user_id    = vastdata_user.u.id
  tenant_id  = data.vastdata_tenant.t.id
  access_key = "{VALID_ACCESS_2}"
  secret_key = "{VALID_SECRET_2}"
  enabled    = true
}}
""",
        )
        tf_ready["apply", "-auto-approve", "-input=false", "-lock=false"]()
        after = _attr(terraform_workdir, "vastdata_user_key", "k")
        assert after["access_key"] == VALID_ACCESS_2
        assert after["secret_key"] == VALID_SECRET_2
    finally:
        tf_ready["destroy", "-auto-approve", "-input=false", "-lock=false"].run(retcode=None)


@pytest.mark.e2e
def test_terf319_active_directory_destroy_with_admin_passwd_in_state(tf_ready, terraform_workdir):
    """admin_passwd must remain in state so destroy can leave the domain."""
    u = _uniq()
    # NetBIOS machine names are capped at 15 chars: vastdb + 4 + join = 14
    machine = f"vastdb{u}join"
    _write_tf(
        terraform_workdir,
        "terf319.tf",
        f"""
resource "vastdata_active_directory" "ad" {{
  machine_account_name            = "{machine}"
  group_login_name                = "uid"
  organizational_unit             = "OU=VASTs,OU=VastENG,DC=VastENG,DC=lab"
  use_auto_discovery              = true
  binddn                          = "CN=vastdata,OU=SmbOU,OU=VastENG,DC=VastENG,DC=lab"
  bindpw                          = "vastdata"
  use_ldaps                       = false
  domain_name                     = "VastENG.lab"
  method                          = "simple"
  query_groups_mode               = "COMPATIBLE"
  use_tls                         = false
  is_vms_auth_provider            = false
  smb_allowed                     = true
  ntlm_enabled                    = true
  port                            = 389
  posix_attributes_source         = "JOINED_DOMAIN"
  reverse_lookup                  = false
  gid_number                      = "gidNumber"
  uid                             = "sAMAccountName"
  uid_number                      = "thisAttributeDoesNotExist"
  use_multi_forest                = true
  match_user                      = "sAMAccountName"
  uid_member_value_property_name  = "sAMAccountName"
  uid_member                      = "member"
  posix_account                   = "user"
  posix_group                     = "group"
  username_property_name          = "name"
  user_login_name                 = "sAMAccountName"
  mail_property_name              = "mail"
  monitor_action                  = "PING"
  abac_read_only_value_name       = "ro"
  abac_read_write_value_name      = "rw"
  scheduled_ma_pwd_change_enabled = false
  admin_username                  = "vastpower"
  admin_passwd                    = "vastdata"
  enabled                         = true
}}
""",
    )
    try:
        tf_ready["apply", "-auto-approve", "-input=false", "-lock=false"]()
        ad = _attr(terraform_workdir, "vastdata_active_directory", "ad")
        assert ad.get("admin_passwd") == "vastdata", "admin_passwd must stay in state for leave"
        assert ad.get("enabled") is True
        # bindpw is write-only: cleared from state, fingerprint kept
        assert not ad.get("bindpw")
        assert ad.get("bindpw_md5")

        rc, stdout, stderr = tf_ready[
            "destroy", "-auto-approve", "-input=false", "-lock=false"
        ].run(retcode=None)
        assert rc == 0, (
            f"destroy/leave failed — admin_passwd missing from state?\n{stdout}\n{stderr}"
        )
    except Exception:
        tf_ready["destroy", "-auto-approve", "-input=false", "-lock=false"].run(retcode=None)
        raise
