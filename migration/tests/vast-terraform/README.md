# vast-terraform
Terraform modules to maintain the configuration of Vast clusters

The `vastdata` terraform provider must be setup in your environment. It is maintained in terraforms default repos so an environment

## Secure Backend Configuration

This project uses Vault to securely manage the Terraform backend connection string. The backend configuration has been moved out of the code for security reasons.

### Prerequisites

1. Vault CLI installed and configured
2. Authenticated with Vault:
   ```bash
   vault login -method=ldap username=<your-username>
   ```
   Note: The `vault login` command automatically stores your token locally, no need to set VAULT_TOKEN manually.
3. Access to the secret at `kv/teams/linfra/vast-terraform`

### Setup

Run the setup script to retrieve the connection string from Vault and create a backend configuration file:

```bash
./scripts/setup_backend.sh
terraform init -backend-config=backend.conf
```

## Terraform Commands

#### Basic usage

After making changes to the repo, changes can be viewed using the `plan` command below and applied using the `apply` command

```commandline
terraform plan  -out vast.plan
terraform apply vast.plan
```

### Repo Structure

When configuring a specific view, for an application, group or dataset, all required config should exist in `view_name.tf`.
This should include:

* `vastdata_view` (mandatory)
* `vastdata_view_policy` (optional) the view does not need a dedicated view_policy but the terraform resource must always reference one
* `vastdata_quota` (mandatory)

Other configs are found as follows:

* `main.ft` will hold the provider config and basic cluster admin configuration such as `ldap`, `groups` etc
* `views.tf` will hold shared views that we expect to serve multiple teams/services/functions
* `view_policies.tf` will hold any view policy associated with a shared view


### Views

Recall that within Vast a `View` is a set of access policies on top of a given filesystem path. A `View` can also map to a bucket,
the bucket name must be unique within the tenant (we currently only use the default tenant, so all bucket names have to be unique regardless of owner)
All `View`s require a `View Policy` to be attached to them which contains the rules for accessing data within that view.

#### Creating a new view

The minimum requirement for creating a new view in Vast via terraform looks like below
```commandline
resource "vastdata_view" "pcaps-raw" {
    path       = "/data/pcaps/raw"
    bucket     = "raw"
    create_dir = "true"
    policy_id  = vastdata_view_policy.data-pcaps.id
    protocols  = ["NFS","S3"]
    bucket_owner = "svcpcaps@mavensecurities.com"
}
```
The fields are:
* `path`: the filesystem path within Vast that this view applies to.
* `bucket` (optional): the name of the bucket associated with this view
* `create_dir`: Create the path if it does not already exist. This is necessary because we can create views on paths that already exist
* `policy_id`: The view policy to be applied, this must exist already _and_ be managed by Terraform. This resource can be imported, see below
* `protocols`: The protocols that are allowed access the data in the view
* `bucket_ownder`: Owner of the bucket if it is to be created as well

##### NFS Owners and Groups

The creation of a view _can_ set the owner of a path when a `bucket` and `bucket_owner` are both supplied but by default,
neither are required and so paths will be created with `root:root` ownership.

Currently within this repo there is no mechanism to automatically assign a paths owner and group. This has to be done manually
as a second step via `cdplatform01` where the vast root is mounted without root perms being squashed. Example below

```commandline
sudo chown svcptq:sgdatastrikepnl  /mnt/vast-root/data/strikepnl/
```


##### Creating S3 buckets

There is a strong preference to always create buckets via Vasts tooling rather than the native S3 protocol.
Unless a bucket owner only creates buckets within a single View, the path linked to the bucket cannot be guarenteed.

More info is available on https://wiki.mavensecurities.com/display/PE/VAST+Access%3A+S3


#### Import an existing resource
Sometimes resources have been defined within the cluster UI, like when Vast engineers do inital setup and troubleshoot.

For any resource managed within terraform, all dependencies need to exist within the terraform state. In order to import
a resource, if needs to be defined within a `.tf` file. This does not necessarily mean that it's full configuration needs
to be present, although that will always be better.

Example for the VIP pool 'Prod'
Add the following to a `.tf` file
```commandline
resource "vastdata_vip_pool" "prod" {
  name = "Prod"
}
```
Import the resource
```commandline
[18/02 14:38:15] michael.moyles@local vast-terraform $ terraform import  vastdata_vip_pool.prod Prod
vastdata_vip_pool.prod: Importing from ID "Prod"...
vastdata_vip_pool.prod: Import prepared!
  Prepared vastdata_vip_pool for import
vastdata_vip_pool.prod: Refreshing state... [id=2]

Import successful!

The resources that were imported are shown above. These resources are now in
your Terraform state and will henceforth be managed by Terraform.

```

#### Importing a view
Vast does not identify views by name, it uses a combination of path, numerical ID and tenant. To import a view use
```commandline
[19/02 16:14:14] michael.moyles@local $ terraform import vastdata_view.pcaps-raw '/data/pcaps/raw|default'
vastdata_view.pcaps-raw: Importing from ID "/data/pcaps/raw|default"...
vastdata_view.pcaps-raw: Import prepared!
  Prepared vastdata_view for import
vastdata_view.pcaps-raw: Refreshing state... [id=10]

Import successful!

The resources that were imported are shown above. These resources are now in
your Terraform state and will henceforth be managed by Terraform.


```

## Recreating resources
Sometimes there will be a need to delete a resource from the cluster and recreate it via Terraform.
This may produce the following error
```╷
│ Error: Error occured while obtaining data from the vastdata cluster
│
│   with vastdata_view.pcaps-raw,
│   on view_pcaps.tf line 11, in resource "vastdata_view" "pcaps-raw":
│   11: resource "vastdata_view" "pcaps-raw" {
│
│ Response Status code is 404 , which is not allowed
╵
```
In order for this to work, once the resource has been removed from the cluster (say via UI) then it needs
to be removed from `.tfstate`. Never edit this file, remove it via
````commandline
[19/02 16:44:44] michael.moyles@local vast-terraform $ terraform state rm vastdata_view.pcaps-raw
Removed vastdata_view.pcaps-raw
Successfully removed 1 resource instance(s).

````
