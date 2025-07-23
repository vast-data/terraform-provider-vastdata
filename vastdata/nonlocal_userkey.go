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

var NonlocalUserKeySchemaRef = is.NewSchemaReference(
	http.MethodPost,
	"users/non_local_keys",
	http.MethodGet,
	"users/non_local_keys",
)

type NonlocalUserKey struct {
	tfstate *is.TFState
}

func (m *NonlocalUserKey) NewResourceManager(raw map[string]attr.Value, schema any) ResourceManager {
	return &NonlocalUserKey{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef:            NonlocalUserKeySchemaRef,
			SensitiveFields:      []string{"secret_key"},
			ExcludedSchemaFields: []string{"login_name"},
			SearchableFields:     []string{"uid", "username"}, // User can be found by uid or username
			ComputedSchemaFields: []string{"uid", "username"},
			AdditionalSchemaAttributes: map[string]any{
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
		},
	)}
}

func (m *NonlocalUserKey) TfState() *is.TFState {
	return m.tfstate
}

func (m *NonlocalUserKey) API(rest *VMSRest) VastResourceAPIWithContext {
	return rest.NonLocalUserKeys
}

func (m *NonlocalUserKey) ReadResource(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	ts := m.tfstate
	if !ts.IsKnownAndNotNull("uid") {
		userRecord, err := rest.Users.GetWithContext(ctx, params{"name": ts.String("username")})
		if err != nil {
			return nil, err
		}
		ts.Set("uid", is.Must(toInt(userRecord["uid"])))
	}
	return nil, nil
}

func (m *NonlocalUserKey) PrepareCreateResource(_ context.Context, _ *VMSRest) error {
	ts := m.tfstate
	if ts.IsKnownAndNotNull("pgp_public_key") {
		if _, err := helper.EncryptMessageArmored(
			ts.String("pgp_public_key"), "######",
		); err != nil {
			return err
		}
	}
	return nil
}

func (m *NonlocalUserKey) CreateResource(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	ts := m.tfstate
	if _, err := m.ReadResource(ctx, rest); err != nil {
		return nil, err
	}
	uid := ts.Int64("uid")
	createParams := params{"uid": uid}
	ts.SetToMapIfAvailable(createParams, "tenant_id", "enabled")
	record, err := rest.NonLocalUserKeys.CreateWithContext(ctx, createParams)
	if err != nil {
		return nil, err
	}
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
	record["uid"] = uid
	return record, err
}

func (m *NonlocalUserKey) UpdateResource(ctx context.Context, plan UpdateResource, rest *VMSRest) (DisplayableRecord, error) {
	var (
		ts          = m.tfstate
		uid         = ts.Int64("uid")
		planManager = plan.(*NonlocalUserKey)
		planTs      = planManager.tfstate
	)

	// Handle enabled/disabled status toggle
	if planTs.IsKnownAndNotNull("enabled") {
		updateParams := params{
			"uid":        uid,
			"access_key": ts.String("access_key"),
			"enabled":    planTs.Bool("enabled"),
		}
		if _, err := rest.NonLocalUserKeys.UpdateNonIdWithContext(ctx, updateParams); err != nil {
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

func (m *NonlocalUserKey) DeleteResource(ctx context.Context, rest *VMSRest) error {
	ts := m.tfstate
	accessKey := ts.String("access_key")
	if accessKey == "" {
		return errors.New("access_key must be specified for deletion")
	}
	if _, err := m.ReadResource(ctx, rest); err != nil {
		return err
	}
	deleteParams := params{"access_key": accessKey, "uid": ts.Int64("uid")}
	_, err := rest.NonLocalUserKeys.DeleteNonIdWithContext(ctx, deleteParams)
	if ignoreStatusCodes(err, http.StatusNotFound) != nil {
		return err
	}
	return nil
}
