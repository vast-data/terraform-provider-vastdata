// // Copyright (c) HashiCorp, Inc.
package provider

import (
	"context"
	"errors"
	"fmt"
	"github.com/ProtonMail/gopenpgp/v2/helper"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
	"net/http"
)

var UserKeySchemaRef = is.NewSchemaReference(
	http.MethodPost,
	"users/{id}/access_keys",
	http.MethodGet,
	"users/{id}/access_keys",
)

type UserKey struct {
	tfstate *is.TFState
}

func (m *UserKey) NewResourceManager(raw map[string]attr.Value, schema any) ResourceManager {
	return &UserKey{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef: UserKeySchemaRef,
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
			SensitiveFields:  []string{"secret_key"},
		},
	)}
}

func (m *UserKey) TfState() *is.TFState {
	return m.tfstate
}

func (m *UserKey) API(rest *VMSRest) VastResourceAPIWithContext {
	return rest.UserKeys
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
	if ts.IsKnownAndNotNull("pgp_public_key") {
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
	if ts.IsKnownAndNotNull("enabled") && !ts.Bool("enabled") {
		if _, err = rest.UserKeys.DisableKeyWithContext(ctx, userId, record["access_key"].(string)); err != nil {
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
		err         error
	)

	// Handle enabled/disabled status toggle
	if planTs.IsKnownAndNotNull("enabled") {
		accessKey := ts.String("access_key")
		if planTs.Bool("enabled") {
			_, err = rest.UserKeys.EnableKeyWithContext(ctx, userId, accessKey)
		} else {
			_, err = rest.UserKeys.DisableKeyWithContext(ctx, userId, accessKey)
		}
		if err != nil {
			return nil, err
		}
	}
	// Conditionally encrypt secret_key if not yet encrypted
	if planTs.IsKnownAndNotNull("pgp_public_key") {
		if ts.IsNull("secret_key") {
			return nil, fmt.Errorf("secret key %q is already encrypted, cannot encrypt again", ts.String("access_key"))
		} else {
			secretKey := ts.String("secret_key")
			pgp := planTs.String("pgp_public_key")
			encrypted, err := helper.EncryptMessageArmored(pgp, secretKey)
			if err != nil {
				return nil, err
			}
			ts.Set("encrypted_secret_key", encrypted)
			ts.Set("secret_key", types.StringNull())
		}
	}
	// Nothing else to do, return nil to keep state unchanged
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
	_, err := rest.UserKeys.DeleteKeyWithContext(ctx, userId, accessKey)
	if ignoreStatusCodes(err, http.StatusNotFound) != nil {
		return err
	}
	return nil
}
