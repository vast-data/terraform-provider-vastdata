// Copyright (c) HashiCorp, Inc.

package provider

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	planmodifiers "github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
)

type S3PolicyAttachment struct {
	tfstate *is.TFState
}

func (m *S3PolicyAttachment) NewResourceManager(raw map[string]attr.Value, schema any) ResourceManager {
	return &S3PolicyAttachment{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			Importable: &notImportable,
			TFStateHintsForCustom: &is.TFStateHintsForCustom{
				Description: "One-to-one association between an S3 policy and a non-local group or user. This resource attaches a single S3 policy to either a group (identified by 'gid') or a user (identified by 'uid').",
				SchemaAttributes: map[string]any{
					"gid": rschema.Int64Attribute{
						Optional:    true,
						Description: "The GID of the non-local group to attach the policy to.",
						PlanModifiers: []planmodifiers.Int64{
							int64planmodifier.RequiresReplace(),
						},
					},
					"uid": rschema.Int64Attribute{
						Optional:    true,
						Description: "The UID of the non-local user to attach the policy to.",
						PlanModifiers: []planmodifiers.Int64{
							int64planmodifier.RequiresReplace(),
						},
					},
					"s3_policy_id": rschema.Int64Attribute{
						Required:    true,
						Description: "The ID of the S3 policy to attach.",
					},
					"ignore_present": rschema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Description: "If set to true, the resource will not return an error if the specified S3 policy is already attached to the user or group. This is useful for gracefully handling pre-existing attachments.",
						Default:     booldefault.StaticBool(false),
					},
					"context": rschema.StringAttribute{
						Optional:    true,
						Description: "Specify the context for the user/group query.",
						PlanModifiers: []planmodifiers.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
				},
			},
		},
	)}
}

func (m *S3PolicyAttachment) TfState() *is.TFState {
	return m.tfstate
}

func (m *S3PolicyAttachment) API(_ *VMSRest) VastResourceAPIWithContext {
	return nil
}

func (m *S3PolicyAttachment) ValidateResourceConfig(context.Context) error {
	return validateOneOf(m.tfstate, "gid", "uid")
}

func (m *S3PolicyAttachment) ReadResource(_ context.Context, _ *VMSRest) (DisplayableRecord, error) {
	return nil, nil
}

func (m *S3PolicyAttachment) CreateResource(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	var (
		ts         = m.tfstate
		s3PolicyId = ts.Int64("s3_policy_id")
		key        string
		val        int64
		getFn      RestFn
		updateFn   RestFn
	)

	switch {
	case ts.IsKnownAndNotNull("gid"):
		key = "gid"
		val = ts.Int64("gid")
		getFn = rest.NonLocalGroups.GetWithContext
		updateFn = rest.NonLocalGroups.UpdateNonLocalGroupWithContext
		defer rest.NonLocalGroups.Lock(key, val)()

	case ts.IsKnownAndNotNull("uid"):
		key = "uid"
		val = ts.Int64("uid")
		getFn = rest.NonLocalUsers.GetWithContext
		updateFn = rest.NonLocalUsers.UpdateNonLocalUserWithContext
		defer rest.NonLocalUsers.Lock(key, val)()

	default:
		return nil, errors.New("either 'gid' or 'uid' must be specified")
	}

	searchParams := params{key: val}
	ts.SetToMapIfAvailable(searchParams, "context")
	record, err := getFn(ctx, searchParams)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch record by %s=%d: %w", key, val, err)
	}

	set := is.Must(is.NewSetFromAny[int64](record["s3_policies_ids"]))

	if set.Add(s3PolicyId) {
		searchParams["s3_policies_ids"] = set.ToSlice()
		return updateFn(ctx, searchParams)
	} else if ts.IsKnownAndNotNull("ignore_present") && !ts.Bool("ignore_present") {
		return nil, fmt.Errorf("s3 policy ID %d is already attached to %s=%d", s3PolicyId, key, val)
	}
	return nil, nil

}

func (m *S3PolicyAttachment) UpdateResource(ctx context.Context, plan UpdateResource, rest *VMSRest) (DisplayableRecord, error) {
	var (
		ts          = m.tfstate
		oldPolicyId = ts.Int64("s3_policy_id")
		planManager = plan.(*S3PolicyAttachment)
		planTs      = planManager.tfstate
		newPolicyId = planTs.Int64("s3_policy_id")
		key         string
		val         int64
		getFn       RestFn
		updateFn    RestFn
	)

	switch {
	case ts.IsKnownAndNotNull("gid"):
		key = "gid"
		val = ts.Int64("gid")
		getFn = rest.NonLocalGroups.GetWithContext
		updateFn = rest.NonLocalGroups.UpdateNonLocalGroupWithContext
		defer rest.NonLocalGroups.Lock(key, val)()

	case ts.IsKnownAndNotNull("uid"):
		key = "uid"
		val = ts.Int64("uid")
		getFn = rest.NonLocalUsers.GetWithContext
		updateFn = rest.NonLocalUsers.UpdateNonLocalUserWithContext
		defer rest.NonLocalUsers.Lock(key, val)()

	default:
		return nil, errors.New("either 'gid' or 'uid' must be specified")
	}

	// No-op if the policy ID hasn’t changed
	if oldPolicyId == newPolicyId {
		return nil, nil
	}

	searchParams := params{key: val}
	ts.SetToMapIfAvailable(searchParams, "context")
	record, err := getFn(ctx, searchParams)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch record by %s=%d: %w", key, val, err)
	}

	set := is.Must(is.NewSetFromAny[int64](record["s3_policies_ids"]))

	// Remove old policy if it existed
	removed := set.Remove(oldPolicyId)

	// Add new policy (if not already present)
	added := set.Add(newPolicyId)

	// policy already attached and unchanged
	if !added && ts.IsKnownAndNotNull("ignore_present") && !ts.Bool("ignore_present") {
		return nil, fmt.Errorf("s3 policy ID %d is already attached to %s=%d", newPolicyId, key, val)
	}

	if removed || added {
		// If we removed the old policy or added a new one, we need to update
		searchParams["s3_policies_ids"] = set.ToSlice()
		return updateFn(ctx, searchParams)
	}

	return nil, nil
}

func (m *S3PolicyAttachment) DeleteResource(ctx context.Context, rest *VMSRest) error {
	ts := m.tfstate
	s3PolicyId := ts.Int64("s3_policy_id")

	var (
		key      string
		val      int64
		getFn    RestFn
		updateFn RestFn
	)

	switch {
	case ts.IsKnownAndNotNull("gid"):
		key = "gid"
		val = ts.Int64("gid")
		getFn = rest.NonLocalGroups.GetWithContext
		updateFn = rest.NonLocalGroups.UpdateNonLocalGroupWithContext
		defer rest.NonLocalGroups.Lock(key, val)()

	case ts.IsKnownAndNotNull("uid"):
		key = "uid"
		val = ts.Int64("uid")
		getFn = rest.NonLocalUsers.GetWithContext
		updateFn = rest.NonLocalUsers.UpdateNonLocalUserWithContext
		defer rest.NonLocalUsers.Lock(key, val)()

	default:
		return fmt.Errorf("either 'gid' or 'uid' must be specified")
	}

	searchParams := params{key: val}
	ts.SetToMapIfAvailable(searchParams, "context")
	record, err := getFn(ctx, searchParams)
	if err != nil {
		return fmt.Errorf("failed to fetch record by %s=%d: %w", key, val, err)
	}

	set := is.Must(is.NewSetFromAny[int64](record["s3_policies_ids"]))

	if !set.Remove(s3PolicyId) {
		// Policy was not present — nothing to do
		return nil
	}

	searchParams["s3_policies_ids"] = set.ToSlice()
	_, err = updateFn(ctx, searchParams)
	return err
}
