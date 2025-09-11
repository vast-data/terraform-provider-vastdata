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
						Optional:    true,
						Description: "The ID of the S3 policy to attach. Either 's3_policy_id' or 's3_policy_guid' must be provided.",
					},
					"s3_policy_guid": rschema.StringAttribute{
						Optional:    true,
						Computed:    true,
						Description: "The GUID of the S3 policy to attach. Either 's3_policy_id' or 's3_policy_guid' must be provided.",
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
					"tenant_id": rschema.Int64Attribute{
						Optional:    true,
						Description: "The ID of the tenant to which the user or group belongs.",
						PlanModifiers: []planmodifiers.Int64{
							int64planmodifier.RequiresReplace(),
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

// ensurePolicyIDAndGUID ensures s3_policy_guid is populated in tfstate.
// Supports both ways: user provides either s3_policy_id OR s3_policy_guid.
// Returns the current s3_policy_id to use for operations.
func (m *S3PolicyAttachment) ensurePolicyIDAndGUID(ctx context.Context, rest *VMSRest, ts *is.TFState) (int64, error) {

	if ts.IsKnownAndNotNull("s3_policy_guid") {
		// User provided GUID, lookup current ID (but don't store it in state)
		s3PolicyGuid := ts.String("s3_policy_guid")
		record, err := rest.S3Policies.GetWithContext(ctx, params{"guid": s3PolicyGuid})
		if err != nil {
			return 0, err
		}
		return record.RecordID(), nil
	}

	if ts.IsKnownAndNotNull("s3_policy_id") {
		// User provided ID, lookup GUID and store it
		s3PolicyId := ts.Int64("s3_policy_id")
		record, err := rest.S3Policies.GetByIdWithContext(ctx, s3PolicyId)
		if err != nil {
			return 0, err
		}
		ts.Set("s3_policy_guid", record.RecordGUID())

		return s3PolicyId, nil
	}

	return 0, fmt.Errorf("either s3_policy_id or s3_policy_guid must be provided")
}

// validateS3PolicyAttachmentConfig validates that exactly one of gid/uid is set and
// exactly one of s3_policy_id/s3_policy_guid is set. This validation is performed
// at runtime when resource references can be resolved.
func (m *S3PolicyAttachment) validateS3PolicyAttachmentConfig() error {
	if err := validateOneOf(m.tfstate, "gid", "uid"); err != nil {
		return err
	}
	if err := validateOneOf(m.tfstate, "s3_policy_id", "s3_policy_guid"); err != nil {
		return err
	}
	return nil
}

func (m *S3PolicyAttachment) ReadResource(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	var (
		ts               = m.tfstate
		key              string
		val              int64
		getFn            RestFn
		updateFn         RestFn
		s3PolicyId       int64
		actualS3PolicyId int64
		err              error
	)

	if ts.IsKnownAndNotNull("s3_policy_id") {
		s3PolicyId = ts.Int64("s3_policy_id")
	}

	// Ensure both s3_policy_id and s3_policy_guid are set
	actualS3PolicyId, err = m.ensurePolicyIDAndGUID(ctx, rest, ts)
	if err != nil {
		return nil, err
	}

	if s3PolicyId != actualS3PolicyId {
		// ID has changed out from under us. We need to detach the old ID and atach the new one.
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
		ts.SetToMapIfAvailable(searchParams, "context", "tenant_id")
		record, err := getFn(ctx, searchParams)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch record by %s=%d: %w", key, val, err)
		}

		set := is.Must(is.NewSetFromAny[int64](record["s3_policies_ids"]))
		set.Remove(s3PolicyId)
		if set.Add(actualS3PolicyId) {
			searchParams["s3_policies_ids"] = set.ToSlice()
			if _, err := updateFn(ctx, searchParams); err != nil {
				return nil, err
			}
		}
	}
	return nil, nil

}

func (m *S3PolicyAttachment) CreateResource(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	var (
		ts         = m.tfstate
		key        string
		val        int64
		getFn      RestFn
		updateFn   RestFn
		record     Record
		s3PolicyId int64
		err        error
	)

	// Validate configuration now that resource references are resolved
	if err := m.validateS3PolicyAttachmentConfig(); err != nil {
		return nil, err
	}

	// Ensure both s3_policy_id and s3_policy_guid are set
	s3PolicyId, err = m.ensurePolicyIDAndGUID(ctx, rest, ts)
	if err != nil {
		return nil, err
	}

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
	ts.SetToMapIfAvailable(searchParams, "context", "tenant_id")
	record, err = getFn(ctx, searchParams)
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
		ts                = m.tfstate
		actualOldPolicyId int64
		planManager       = plan.(*S3PolicyAttachment)
		planTs            = planManager.tfstate
		newPolicyId       int64
		key               string
		val               int64
		getFn             RestFn
		updateFn          RestFn
		record            Record
		err               error
	)

	// Fetch actual old policy ID from GUID, in case it changed out from under us
	actualOldPolicyId, err = m.ensurePolicyIDAndGUID(ctx, rest, ts)
	if err != nil && !isNotFoundErr(err) {
		return nil, err
	}

	// Handle new policy: lookup policy info from plan
	newPolicyId, err = m.ensurePolicyIDAndGUID(ctx, rest, planTs)
	if err != nil {
		return nil, err
	}

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

	// No-op if the policy ID hasn't changed
	if actualOldPolicyId == newPolicyId {
		return nil, nil
	}

	searchParams := params{key: val}
	ts.SetToMapIfAvailable(searchParams, "context", "tenant_id")
	record, err = getFn(ctx, searchParams)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch record by %s=%d: %w", key, val, err)
	}

	set := is.Must(is.NewSetFromAny[int64](record["s3_policies_ids"]))

	// Remove old policy if it existed
	var removed bool
	if ts.IsKnownAndNotNull("s3_policy_id") {
		oldPolicyId := ts.Int64("s3_policy_id")
		removed = set.Remove(oldPolicyId) || set.Remove(actualOldPolicyId)
	} else {
		// We only had GUID before, so just remove the actual old ID
		removed = set.Remove(actualOldPolicyId)
	}

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
	var (
		ts             = m.tfstate
		key            string
		val            int64
		getFn          RestFn
		updateFn       RestFn
		actualPolicyId int64
		err            error
	)

	// Fetch actual old policy ID from GUID, in case it changed out from under us
	actualPolicyId, err = m.ensurePolicyIDAndGUID(ctx, rest, ts)
	if err != nil && !isNotFoundErr(err) {
		return err
	}

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
	ts.SetToMapIfAvailable(searchParams, "context", "tenant_id")
	record, err := getFn(ctx, searchParams)
	if err != nil {
		return fmt.Errorf("failed to fetch record by %s=%d: %w", key, val, err)
	}

	set := is.Must(is.NewSetFromAny[int64](record["s3_policies_ids"]))

	var removed bool
	if ts.IsKnownAndNotNull("s3_policy_id") {
		s3PolicyId := ts.Int64("s3_policy_id")
		// Remove by both old ID and actual ID, in case it changed out from under us
		removed = set.Remove(s3PolicyId) || set.Remove(actualPolicyId)
	} else {
		// We only had GUID before, so just remove the actual old ID
		removed = set.Remove(actualPolicyId)
	}

	if !removed {
		// Policy was not present — nothing to do
		return nil
	}

	searchParams["s3_policies_ids"] = set.ToSlice()
	_, err = updateFn(ctx, searchParams)
	return err
}
