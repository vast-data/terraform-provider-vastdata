// // Copyright (c) HashiCorp, Inc.
package provider

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	vast_client "github.com/vast-data/go-vast-client"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
)

var UserCopySchemaRef = is.NewSchemaReference(
	http.MethodPost,
	"users/copy",
	"",
	"",
)

type UserCopy struct {
	tfstate *is.TFState
}

func (m *UserCopy) NewResourceManager(raw map[string]attr.Value, schema any) ResourceManager {
	return &UserCopy{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			Importable: &notImportable,
			SchemaRef:  UserCopySchemaRef,
		},
	)}
}

func (m *UserCopy) TfState() *is.TFState {
	return m.tfstate
}

func (m *UserCopy) API(rest *VMSRest) VastResourceAPIWithContext {
	return rest.Users
}

func (m *UserCopy) ReadResource(_ context.Context, _ *VMSRest) (DisplayableRecord, error) {
	return nil, nil
}

func (m *UserCopy) CreateResource(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	ts := m.tfstate

	// Validate that either tenant_id or user_ids is provided, but not both
	hasTenantID := ts.IsKnownAndNotNull("tenant_id")
	hasUserIDs := ts.IsKnownAndNotNull("user_ids")

	if !hasTenantID && !hasUserIDs {
		return nil, fmt.Errorf("either tenant_id or user_ids must be provided")
	}

	if hasTenantID && hasUserIDs {
		return nil, fmt.Errorf("cannot provide both tenant_id and user_ids")
	}

	// Prepare the copy parameters
	params := vast_client.UsersCopyParams{
		DestinationProviderID: ts.Int64("destination_provider_id"),
	}

	if hasTenantID {
		params.TenantID = ts.Int64("tenant_id")
	}

	if hasUserIDs {
		userIDs := ts.ToSlice("user_ids")
		params.UserIDs = make([]int64, len(userIDs))
		for i, userID := range userIDs {
			if id, ok := userID.(int64); ok {
				params.UserIDs[i] = id
			} else {
				return nil, fmt.Errorf("invalid user_id type: expected int64, got %T", userID)
			}
		}
	}

	// Execute the copy operation using the CopyWithContext method
	err := rest.Users.CopyWithContext(ctx, params)
	return nil, err
}

func (m *UserCopy) UpdateResource(ctx context.Context, plan UpdateResource, rest *VMSRest) (DisplayableRecord, error) {
	// With force_new modifiers, Terraform will handle replacements automatically
	// This method should not be called for updates since all fields have RequiresReplace()
	// But we'll keep it as a safety net in case it's called
	return nil, fmt.Errorf("user copy operations should be replaced, not updated")
}

func (m *UserCopy) DeleteResource(ctx context.Context, rest *VMSRest) error {
	// No-op: UserCopy cannot be deleted - it's a one-time operation
	return nil
}
