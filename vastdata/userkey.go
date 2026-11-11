// // Copyright (c) HashiCorp, Inc.
package provider

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"slices"

	"github.com/ProtonMail/gopenpgp/v2/helper"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
)

// VMS rejects secrets that are not exactly 40 chars of [A-Za-z0-9/+].
var userKeySecretKeyPattern = regexp.MustCompile(`^[A-Za-z0-9/+]{40}$`)

var UserKeySchemaRef = is.NewSchemaReference(
	http.MethodPost,
	"users/{id}/access_keys",
	http.MethodGet,
	"users/{id}/access_keys",
)

type UserKey struct {
	tfstate is.ResourceState
}

func (m *UserKey) NewResourceManager(raw map[string]attr.Value, schema any) ResourceManager {
	return &UserKey{tfstate: is.NewResourceState(
		raw,
		schema,
		&is.TFStateHints{
			Importable:              &notImportable,
			SchemaRef:               UserKeySchemaRef,
			OptionalSchemaFields:    []string{"access_key", "secret_key"},
			ComputedSchemaFields:    []string{"access_key", "secret_key"},
			PreserveUserValueFields: []string{"access_key", "secret_key"},
			// Not WriteOnly: cluster-generated secrets must remain readable in state
			// (TERF-318). User-supplied secrets stay in state too — Terraform forbids
			// nulling a non-write-only sensitive plan value after apply.
			SensitiveFields: []string{"secret_key"},
			AdditionalSchemaAttributes: map[string]any{
				"user_id": rschema.Int64Attribute{
					Optional:    true,
					Computed:    true,
					Description: "The ID of the user to which this key belongs. If not provided, it will be derived from the username.",
				},
				"username": rschema.StringAttribute{
					Optional:    true,
					Computed:    true,
					Description: "The username of the user to which this key belongs.",
				},
				"pgp_public_key": rschema.StringAttribute{
					Optional:    true,
					Sensitive:   true,
					Description: "Optional PGP public key to encrypt the secret key.",
				},
				"encrypted_secret_key": rschema.StringAttribute{
					Computed:    true,
					Description: "The encrypted secret key, returned if pgp_public_key is used",
				},
				"enabled": rschema.BoolAttribute{
					Optional:    true,
					Description: "Whether the key is enabled.",
				},
			},
			SearchableFields: []string{"user_id", "username"},
		},
	)}
}

func (m *UserKey) TfState() is.ResourceState {
	return m.tfstate
}

// MigrateModePassThroughUpdate allows schema-reconciliation updates during
// VASTDATA_MIGRATE_MODE.  UserKey is non-importable and carried over from
// the old state verbatim, so it may need a schema update on first apply.
func (m *UserKey) MigrateModePassThroughUpdate() {}

func (m *UserKey) API(rest *VMSRest) VastResourceAPIWithContext {
	return rest.Users
}

func (m *UserKey) ReadResource(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	ts := m.tfstate
	if !ts.IsKnownAndNotNull("user_id") {
		userRecord, err := rest.Users.GetWithContext(ctx, params{"name": ts.String("username")})
		if err != nil {
			return nil, err
		}
		ts.Set("user_id", userRecord.RecordID())
	}
	return nil, nil
}

func validatePgpPublicKey(ts is.ResourceState) error {
	if ts.IsNull("pgp_public_key") {
		return nil
	}
	_, err := helper.EncryptMessageArmored(ts.String("pgp_public_key"), "######")
	return err
}

func validateSecretKeyFormat(secretKey string) error {
	if !userKeySecretKeyPattern.MatchString(secretKey) {
		return fmt.Errorf("secret_key must be 40 characters long and can include alphanumeric characters, '/', and '+'")
	}
	return nil
}

func validateCustomUserKeyPair(ts is.ResourceState) error {
	hasAccessKey := ts.IsKnownAndNotNull("access_key")
	hasSecretKey := ts.IsKnownAndNotNull("secret_key")
	if hasAccessKey != hasSecretKey {
		return fmt.Errorf("access_key and secret_key must both be specified or both omitted")
	}
	if hasSecretKey && ts.IsKnownAndNotNull("pgp_public_key") {
		return fmt.Errorf("pgp_public_key cannot be used when secret_key is specified")
	}
	if hasSecretKey {
		if err := validateSecretKeyFormat(ts.String("secret_key")); err != nil {
			return err
		}
	}
	return validatePgpPublicKey(ts)
}

// markUserSuppliedSecretKeyEphemeral asks ClearWriteOnlyFields to drop secret_key
// after CopyNonEmptyFieldsTo would otherwise copy a user-supplied secret back
// into state. Generated secrets never call this, so they remain in state.
func markUserSuppliedSecretKeyEphemeral(ts is.ResourceState) {
	hints := ts.StateHints()
	if hints == nil {
		return
	}
	if !slices.Contains(hints.WriteOnlyFields, "secret_key") {
		hints.WriteOnlyFields = append(hints.WriteOnlyFields, "secret_key")
	}
}

func finalizeUserKeyRecord(record Record, ts is.ResourceState) (Record, error) {
	if ts.IsKnownAndNotNull("pgp_public_key") {
		pgp := ts.String("pgp_public_key")
		secretKey := record["secret_key"].(string)
		encrypted, err := helper.EncryptMessageArmored(pgp, secretKey)
		if err != nil {
			return nil, err
		}
		record["encrypted_secret_key"] = encrypted
		record["secret_key"] = types.StringNull()
		markUserSuppliedSecretKeyEphemeral(ts)
	} else if ts.IsKnownAndNotNull("secret_key") {
		// User-supplied: keep plaintext in state (same attribute as generated).
		record["encrypted_secret_key"] = types.StringNull()
		record["secret_key"] = ts.String("secret_key")
	} else {
		// Cluster-generated: keep the one-shot secret in state for S3 clients.
		record["encrypted_secret_key"] = types.StringNull()
	}
	if ts.IsKnownAndNotNull("access_key") {
		record["access_key"] = ts.String("access_key")
	}
	return record, nil
}

func userKeyCredentialsChanged(planTs, stateTs is.ResourceState) bool {
	changed := planTs.GetChangedParams(stateTs)
	_, accessKeyChanged := changed["access_key"]
	_, secretKeyChanged := changed["secret_key"]
	return accessKeyChanged || secretKeyChanged
}

func userKeyCreateParams(stateTs, planTs is.ResourceState) params {
	createParams := params{}
	if planTs.IsKnownAndNotNull("tenant_id") {
		createParams["tenant_id"] = planTs.Int64("tenant_id")
	} else if stateTs.IsKnownAndNotNull("tenant_id") {
		createParams["tenant_id"] = stateTs.Int64("tenant_id")
	}
	planTs.SetToMapIfAvailable(createParams, "access_key", "secret_key")
	return createParams
}

func restoreUserKeyCredentials(ctx context.Context, rest *VMSRest, userId int64, accessKey, secretKey string, tenantId any) error {
	if accessKey == "" || secretKey == "" {
		return errors.New("cannot restore previous access key: prior secret_key is not available in state")
	}
	body := params{"access_key": accessKey, "secret_key": secretKey}
	if tenantId != nil {
		body["tenant_id"] = tenantId
	}
	_, err := rest.Users.UserAccessKeysWithContext_POST(ctx, userId, body)
	return err
}

// replaceUserKeyCredentials rotates credentials without leaving the user keyless
// when the new key is rejected (TERF-321). Prefer create-before-delete when the
// access_key changes; when it stays the same, validate first, delete, create,
// and restore the prior key from state if create fails.
func replaceUserKeyCredentials(ctx context.Context, rest *VMSRest, userId int64, stateTs, planTs is.ResourceState) (Record, error) {
	oldAccessKey := stateTs.String("access_key")
	if oldAccessKey == "" {
		return nil, errors.New("cannot update credentials: existing access_key is unknown in state")
	}
	oldSecretKey := ""
	if stateTs.IsKnownAndNotNull("secret_key") {
		oldSecretKey = stateTs.String("secret_key")
	}
	if planTs.IsKnownAndNotNull("secret_key") {
		if err := validateSecretKeyFormat(planTs.String("secret_key")); err != nil {
			return nil, err
		}
	}

	createParams := userKeyCreateParams(stateTs, planTs)
	newAccessKey := ""
	if v, ok := createParams["access_key"].(string); ok {
		newAccessKey = v
	}

	// Different access_key: create the new key first so a rejected create leaves
	// the old key intact, then delete the old one.
	if newAccessKey != "" && newAccessKey != oldAccessKey {
		record, err := rest.Users.UserAccessKeysWithContext_POST(ctx, userId, createParams)
		if err != nil {
			return nil, err
		}
		if delErr := rest.Users.UserAccessKeysWithContext_DELETE(ctx, userId, params{"access_key": oldAccessKey}); delErr != nil {
			return nil, fmt.Errorf("created new access key %q but failed to delete previous access key %q: %w", newAccessKey, oldAccessKey, delErr)
		}
		return finalizeUserKeyRecord(record, planTs)
	}

	// Same access_key (or cluster-generated): must delete before recreate.
	if err := rest.Users.UserAccessKeysWithContext_DELETE(ctx, userId, params{"access_key": oldAccessKey}); err != nil {
		return nil, err
	}

	record, err := rest.Users.UserAccessKeysWithContext_POST(ctx, userId, createParams)
	if err != nil {
		var tenantId any
		if v, ok := createParams["tenant_id"]; ok {
			tenantId = v
		}
		if restoreErr := restoreUserKeyCredentials(ctx, rest, userId, oldAccessKey, oldSecretKey, tenantId); restoreErr != nil {
			return nil, fmt.Errorf("failed to create replacement access key after deleting %q: %v; additionally failed to restore previous key: %w", oldAccessKey, err, restoreErr)
		}
		return nil, fmt.Errorf("failed to create replacement access key after deleting %q (previous key restored): %w", oldAccessKey, err)
	}
	return finalizeUserKeyRecord(record, planTs)
}

func (m *UserKey) PrepareCreateResource(_ context.Context, _ *VMSRest) error {
	return validateCustomUserKeyPair(m.tfstate)
}

func (m *UserKey) PrepareUpdateResource(_ context.Context, plan PrepareUpdateResource, _ *VMSRest) error {
	planTs := plan.(*UserKey).tfstate
	if userKeyCredentialsChanged(planTs, m.tfstate) {
		return validateCustomUserKeyPair(planTs)
	}
	return validatePgpPublicKey(planTs)
}

func (m *UserKey) CreateResource(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	ts := m.tfstate
	if _, err := m.ReadResource(ctx, rest); err != nil {
		return nil, err
	}
	userId := ts.Int64("user_id")
	createParams := params{}

	// Get tenant_id from tfstate if provided, otherwise try to get it from user record
	if ts.IsKnownAndNotNull("tenant_id") {
		createParams["tenant_id"] = ts.Int64("tenant_id")
	} else {
		userRecord, err := rest.Users.GetByIdWithContext(ctx, userId)
		if err != nil {
			return nil, fmt.Errorf("failed to get user details for tenant_id: %w", err)
		}
		if tenantId, ok := userRecord["tenant_id"].(int64); ok {
			createParams["tenant_id"] = tenantId
		}
	}

	ts.SetToMapIfAvailable(createParams, "access_key", "secret_key")

	record, err := rest.Users.UserAccessKeysWithContext_POST(ctx, userId, createParams)
	if err != nil {
		return nil, err
	}
	record["user_id"] = userId
	record["username"] = ts.String("username")
	record, err = finalizeUserKeyRecord(record, ts)
	if err != nil {
		return nil, err
	}
	if ts.IsKnownAndNotNull("enabled") && !ts.Bool("enabled") {
		if err = rest.Users.UserAccessKeysWithContext_PATCH(ctx, userId, params{"access_key": record["access_key"].(string), "enabled": false}); err != nil {
			return nil, err
		}
	}
	return record, err
}

func (m *UserKey) UpdateResource(ctx context.Context, plan UpdateResource, rest *VMSRest) (DisplayableRecord, error) {
	var (
		ts          = m.tfstate
		userId      = ts.Int64("user_id")
		planManager = plan.(*UserKey)
		planTs      = planManager.tfstate
		record      Record
		err         error
	)

	if userKeyCredentialsChanged(planTs, ts) {
		record, err = replaceUserKeyCredentials(ctx, rest, userId, ts, planTs)
		if err != nil {
			return nil, err
		}
	}

	if planTs.IsKnownAndNotNull("enabled") {
		accessKey := ts.String("access_key")
		if record != nil {
			if ak, ok := record["access_key"].(string); ok && ak != "" {
				accessKey = ak
			}
		}
		err = rest.Users.UserAccessKeysWithContext_PATCH(ctx, userId, params{"access_key": accessKey, "enabled": planTs.Bool("enabled")})
		if err != nil {
			return nil, err
		}
	}
	if planTs.IsKnownAndNotNull("pgp_public_key") {
		if ts.IsNull("secret_key") {
			return nil, fmt.Errorf("secret key %q is already encrypted, cannot encrypt again", ts.String("access_key"))
		}
		secretKey := ts.String("secret_key")
		pgp := planTs.String("pgp_public_key")
		encrypted, err := helper.EncryptMessageArmored(pgp, secretKey)
		if err != nil {
			return nil, err
		}
		ts.Set("encrypted_secret_key", encrypted)
		ts.Set("secret_key", types.StringNull())
		markUserSuppliedSecretKeyEphemeral(ts)
	}
	if record != nil {
		record["user_id"] = userId
		record["username"] = ts.String("username")
		return record, nil
	}
	return nil, nil
}

func (m *UserKey) DeleteResource(ctx context.Context, rest *VMSRest) error {
	ts := m.tfstate
	accessKey := ts.String("access_key")
	if accessKey == "" {
		return errors.New("access_key must be specified for deletion")
	}
	if _, err := m.ReadResource(ctx, rest); err != nil {
		return err
	}
	userId := ts.Int64("user_id")
	err := rest.Users.UserAccessKeysWithContext_DELETE(ctx, userId, params{"access_key": accessKey})
	if ignoreStatusCodes(err, http.StatusNotFound) != nil {
		return err
	}
	return nil
}
