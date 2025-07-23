# Adding New Resources or Data Sources

To add a new resource or data source, follow these steps:

## 1. Create a Manager

Create a new file (e.g., `user.go`). The file name is not important.

Example implementation for a `User` resource:

```go
package provider

import (
  "github.com/hashicorp/terraform-plugin-framework/attr"
  is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
  "net/http"
)

var UserSchemaRef = is.NewSchemaReference(
  http.MethodPost,
  "users",
  http.MethodGet,
  "users",
)

type User struct {
  tfstate *is.TFState
}

func (m *User) NewResourceManager(raw map[string]attr.Value, schema any) ResourceManager {
  return &User{tfstate: is.NewTFStateMust(
    raw,
    schema,
    &is.TFStateHints{
      SchemaRef:       UserSchemaRef,
    },
  )}
}

func (m *User) NewDatasourceManager(raw map[string]attr.Value, schema any) DataSourceManager {
  return &User{tfstate: is.NewTFStateMust(
    raw,
    schema,
    &is.TFStateHints{
      SchemaRef: UserSchemaRef,
    }),
  }
}

func (m *User) TfState() *is.TFState {
  return m.tfstate
}

func (m *User) API(rest *VMSRest) VastResourceAPIWithContext {
  return rest.Users
}

```

**Key Points:**

This tells the plugin:

- Use the POST /users operation for generating the create/update schema.
- Use the GET /users operation for generating the read schema (response fields).

```go
var UserSchemaRef = is.NewSchemaReference(
  http.MethodPost, "users",  // used for create/update
  http.MethodGet,  "users",  // used for read
)
```

Example:
```yaml
/users:
  post:
    operationId: "users_create"
    requestBody:
      content:
        application/json:
          schema:
            type: object
            properties:
              name:
                type: string
              role:
                type: string
  get:
    operationId: "users_list"
    responses:
      "200":
        content:
          application/json:
            schema:
              type: object
              properties:
               ...
```

- `NewResourceManager` and `NewDatasourceManager`: Implement these to register as a resource, data source, or both.
- `TfState`: Returns the resource state.
- `API`: Returns the VAST resource API. It uses [vast-go-client](https://github.com/vast-data/go-vast-client/blob/main/vast_resource.go).
If resource is not in the client, consider contributing. See this [Developer.md](https://github.com/vast-data/go-vast-client/blob/main/docs/DEVELOPER.md)

## 2. Register the Manager

In `component.go`, add your manager to the list:

```go
var allTFComponents = []TFManager{
	&User{}, // Register User resource and data source
}
```

## 3. Customize Schema (Optional)

Use `is.TFStateHints` to adjust schema fields (e.g., mark fields as required, computed, or sensitive).  
This is only needed for advanced scenarios or if the OpenAPI schema is incomplete. Will be explained later.


Note:
- For most resources, no further changes are needed.
- For non-standard resources (e.g., with async tasks or custom fields), see advanced documentation.
- If a resource is missing in the [go-vast-client](https://github.com/vast-data/go-vast-client/blob/main/vast_resource.go), consider contributing. See [go-vast-client developer guide](https://github.com/vast-data/go-vast-client/blob/main/docs/DEVELOPER.md).

---

### Tuning resource behavior

Plugin defines 2 main structs for handling resources and data sources: 
- `Resource` - you can find in `resource.go`
- `DataSource` - you can find in `datasource.go`

`Datasource` has standard method `Read` which is used to read the resource state from VAST API.
`Resource` has standard methods `Create`, `Read`, `Update`, and `Delete` for managing the resource lifecycle.


We have bunch of method  interceptors to customize the behavior of the resource or data source.

Interceptors are declared in `components.go` implementing any of them on manager (eg `User`) will change the behavior of the resource or data source.

You can implement any of the following interfaces in your manager (e.g. `User`, `UserKeys`) to customize behavior at various points in the resource lifecycle:

##### Import

- `PrepareImportResourceState`  
  Called before import begins. Useful for loading custom state from external sources.

- `ImportResourceState`  
  Provides custom logic to import and populate internal state from VAST API.

- `AfterImportResourceState`  
  Called after import finishes. Useful for additional validation or enrichment.

---

##### Create

- `PrepareCreateResource`  
  Runs before the resource is created. Use this to validate input, derive values, or enrich the plan.

- `CreateResource`  
  Custom implementation of the create logic (instead of the default one).  
  If not implemented, the default logic uses OpenAPI schema and `CreateWithContext()`.

- `AfterCreateResource`  
  Runs after the resource is created and the internal state is filled.  
  Ideal for post-processing, validation, or triggering secondary actions.

---

##### Read

- `PrepareReadResource`  
  Runs before reading from the API. You can alter lookup logic here.

- `ReadResource`  
  Custom read logic.

- `AfterReadResource`  
  Runs after the resource is read and the internal state is filled.

---

##### Update

- `PrepareUpdateResource`  
  Invoked before update, allowing modification of the plan or setup state.

- `UpdateResource`  
  Custom update logic. If not implemented, the default update uses schema diffs and `UpdateWithContext()`.

- `AfterUpdateResource`  
  Runs after update and internal state is populated.  
  Can be used for logging, validation, etc.

---

##### Delete

- `PrepareDeleteResource`  
  Pre-delete hook, useful for validation or cleanup before deleting the resource.

- `DeleteResource`  
  Custom delete logic (e.g., cascading deletes, async deletion).  
  If not defined, the default behavior will delete by `id` or fallback to search-based deletion.

- `AfterDeleteResource`  
  Runs after the resource is deleted from VAST backend.

---

##### Practical examples:

See real example of interceptors usage here: `vastdata/userkey.go`

Example one:

```go
func (m *UserKey) ReadResource(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	ts := m.tfstate
	if !ts.IsKnown("user_id") && ts.IsNull("user_id") {
		userRecord, err := rest.Users.GetWithContext(ctx, params{"name": ts.String("username")})
		if err != nil {
			return nil, err
		}
		ts.Set("user_id", userRecord.RecordID())
	}
	return nil, nil
```

Here we are checking if `user_id` is known and not null, if it is not known we are fetching the user by `username` and setting the `user_id` in the state.
NOTE: We are pretty sure that `username` exists. To be more precise one of `user_id` or `username` should be known. It is handled by upstream `Resource` and `DataSource` structs validators.
So no need to check presence of `username` in the state.

Example above shows usage of `ReadResource` whre resource has no "required" fields. 
You can provide `user_id` or `username` to the resource, but if you provide `username` it will be used to fetch `user_id` from VAST API.

Example two:

```go
func (m *UserKey) PrepareCreateResource(_ context.Context, _ *VMSRest) error {
	ts := m.tfstate
	if !ts.IsNull("pgp_public_key") {
		if _, err := helper.EncryptMessageArmored(
			ts.String("pgp_public_key"), "######",
		); err != nil {
			return err
		}
	}
	return nil
}
```

If `pgp_public_key` is provided, we try to encrypt dummy string before real request to VAST API.
`PrepareCreateResource` - is good for implemeting some "prerequisites" before the resource is created.


Example three:

```go
func (m *UserKey) CreateResource(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	ts := m.tfstate
	if _, err := m.ReadResource(ctx, rest); err != nil {
		return nil, err
	}
	userId := ts.Int64("user_id")
	record, err := rest.UserKeys.CreateKeyWithContext(ctx, userId)
	if err != nil {
		return nil, err
	}
	record["user_id"] = userId
	record["username"] = ts.String("username")
	if !ts.IsNull("pgp_public_key") {
		pgp := ts.String("pgp_public_key")
		secretKey := record["secret_key"].(string)
		encrypted, err := helper.EncryptMessageArmored(pgp, secretKey)
		if err != nil {
			return nil, err
		}
		record["encrypted_secret_key"] = encrypted
		record["secret_key"] = types.StringNull()
	} else {
		record["encrypted_secret_key"] = types.StringNull()
	}
	if !ts.IsNull("enabled") && !ts.Bool("enabled") {
		if _, err = rest.UserKeys.DisableKeyWithContext(ctx, userId, record["access_key"].(string)); err != nil {
			return nil, err
		}
	}
	return record, err
}
```

here we are implementing `CreateResource` method which is used to create New Access key and Secret key pair for the user.
record here:
```go
record, err := rest.UserKeys.CreateKeyWithContext(ctx, userId)
```
is map of string to any, which is the record returned by VAST API.
By default succesfull POST request to userKey resource returns 2 map of 2 fields:
- secret_key - the secret key for the user
- access_key - the access key for the user

Eventually all key/value pairs from the record will be stored in the terraform state. It is upstream responsibility of the `Resource` or `DataSource` to copy all from record to the state.
So here our goal is to customize record that we're going to return upon end of the `CreateResource` method.

Note: when `pgp_public_key` is provided we reset secret key (set it to `types.StringNull()`) and set `encrypted_secret_key` to the encrypted value of the secret key.

Note: it is always a choice modify the record returned by VAST API or to modify the state directly in method.
You can also do:
```go
ts.Set("encrypted_secret_key", encrypted)
```

Instead of passing it via record:
```go
record["encrypted_secret_key"] = encrypted
```

For final result doesn't matter which one you choose but it is recommended to use `record` to keep all "set record-to-state" logic in one place.


Example four:

See: `vastdata/s3policy.go`

We have small tuning in `AfterCreateResource` method.
By default, S3 policy cannot be disabled upon creation (POST mehthod). To disable it we need to use another PATCH method.

Example:
```go
func (m *S3Policy) AfterCreateResource(ctx context.Context, rest *VMSRest, record Record) error {
	var (
		ts = m.tfstate
		id = record.RecordID()
	)

	if !ts.IsNull("enabled") {
		enabled := ts.Bool("enabled")
		if _, err := rest.S3Policies.UpdateWithContext(ctx, id, params{"enabled": enabled}); err != nil {
			return err
		}
		record["enabled"] = enabled
	}
	return nil

}
```

Note: `AfterCreateResource` method signature expects also `record` which is the record returned by VAST API after successful POST request. (Or by your custom `CreateResource` method).
You can modify the record additionally before it is converted to the state.

So here we are checking if `enabled` is set in the state and if it is not null, we are updating the S3 policy with the `enabled` value.
Also, we are setting the `enabled` value in the record to be returned to the state.


### Advanced Schema Customization

In "manager" (For intance `User`) you can use `is.TFStateHints` to customize the schema fields.

It is mostly list of fields that requires modification.
For more info see: `vastdata/internalstate/tf_state_hints.go`

`vastdata/userkey.go` is good example of usage of `is.TFStateHints` to customize the schema fields.
Here we declare manually `user_id`, `username` `pgp_public_key` etc. (Because OpenAPI schema contains only `access_key` and `secret_key`).


### Request/Response Transformation

Some resources require additional transformation logic for request or response payloads — for example, to inject derived fields,
adjust naming conventions, or normalize backend responses. The provider supports this via two transformation interfaces:

`TransformRequestBody`:
```go
type TransformRequestBody interface {
	TransformRequestBody(body params) params
}
```

`TransformResponseRecord`:
```go
type TransformResponseRecord interface {
	TransformResponseRecord(response Record) Record
}
```

## Generate OpenAPI Tarball

The updated OpenAPI schema must be saved as api.tar.gz and placed in the following path:

- `vastdata/client/api/5.3.0/api.tar.gz`

Assuming you have Orion cloned locally, run:

Execute:
```bash
make gen-openapi-tar [orion base path]/management/api/vast_doc.yaml
```


NOTE: you need to install `ruamel.yaml` using command `pip install ruamel.yaml`

After running this command, you should see the following two files in your current working directory:
- openapi.json
- openapi.tar.gz

We only need openapi.tar.gz.
Move it to: `vastdata/client/api/5.3.0/api.tar.gz`

##### Make sure newly generated shema can be parsed properly

After new schema generation execute two commands

To verify all resources can be parsed properly:
```bash
make show r
```

To verify all data-sources can be parsed properly:
```bash
make show d
```

#### Update schemas documentation

If fields were added, renamed, or removed in the new schema, regenerate the Terraform documentation:

```bash
make generate-docs
```

After that go to `docs/data-sources` and `docs/resources` to verify all uncommited changes in terraform schemas.
Pay close attention to any added or removed fields.
If you're expecting specific changes, verify that the updated documentation reflects them.







