// // Copyright (c) HashiCorp, Inc.
package provider

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/ProtonMail/gopenpgp/v2/helper"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
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
			Importable:           &notImportable,
			SchemaRef:            NonlocalUserKeySchemaRef,
			SensitiveFields:      []string{"secret_key"},
			ExcludedSchemaFields: []string{"login_name"},
			SearchableFields:     []string{"uid", "sid", "username"}, // User can be found by uid, sid, or username
			ComputedSchemaFields: []string{"uid", "sid", "username"},
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
	return rest.Users
}

func (m *NonlocalUserKey) ReadResource(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	ts := m.tfstate
	if !ts.IsKnownAndNotNull("uid") && !ts.IsKnownAndNotNull("sid") {
		userRecord, err := rest.Users.GetWithContext(ctx, params{"name": ts.String("username")})
		if err != nil {
			return nil, err
		}
		if uid, ok := userRecord["uid"]; ok {
			if uid.(string) == "" {
				// Fallback to sid if "uid" is empty string.
				ts.Set("sid", userRecord["sid"])
			} else {
				ts.Set("uid", is.Must(toInt(userRecord["uid"])))
			}
		}
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

	// Use uid if available, otherwise use sid
	createParams := params{}
	if ts.IsKnownAndNotNull("uid") {
		createParams["uid"] = ts.Int64("uid")
	} else if ts.IsKnownAndNotNull("sid") {
		createParams["sid"] = ts.String("sid")
	} else {
		return nil, errors.New("either uid or sid must be set")
	}

	ts.SetToMapIfAvailable(createParams, "tenant_id", "enabled")
	record, err := rest.Users.UserNonLocalKeysWithContext_POST(ctx, createParams)
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

	// Preserve uid or sid in the record
	if ts.IsKnownAndNotNull("uid") {
		record["uid"] = ts.Int64("uid")
	}
	if ts.IsKnownAndNotNull("sid") {
		record["sid"] = ts.String("sid")
	}

	return record, err
}

func (m *NonlocalUserKey) UpdateResource(ctx context.Context, plan UpdateResource, rest *VMSRest) (DisplayableRecord, error) {
	var (
		ts          = m.tfstate
		planManager = plan.(*NonlocalUserKey)
		planTs      = planManager.tfstate
	)

	// Handle enabled/disabled status toggle
	if planTs.IsKnownAndNotNull("enabled") {
		updateParams := params{
			"access_key": ts.String("access_key"),
			"enabled":    planTs.Bool("enabled"),
		}

		// Use uid if available, otherwise use sid
		if ts.IsKnownAndNotNull("uid") {
			updateParams["uid"] = ts.Int64("uid")
		} else if ts.IsKnownAndNotNull("sid") {
			updateParams["sid"] = ts.String("sid")
		} else {
			return nil, errors.New("either uid or sid must be set")
		}

		if err := rest.Users.UserNonLocalKeysWithContext_PATCH(ctx, updateParams); err != nil {
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

	deleteParams := params{"access_key": accessKey}

	// Use uid if available, otherwise use sid
	if ts.IsKnownAndNotNull("uid") {
		deleteParams["uid"] = ts.Int64("uid")
	} else if ts.IsKnownAndNotNull("sid") {
		deleteParams["sid"] = ts.String("sid")
	} else {
		return errors.New("either uid or sid must be set for deletion")
	}

	err := rest.Users.UserNonLocalKeysWithContext_DELETE(ctx, deleteParams)
	if ignoreStatusCodes(err, http.StatusNotFound) != nil {
		return err
	}
	return nil
}
