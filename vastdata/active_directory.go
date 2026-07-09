// Copyright (c) HashiCorp, Inc.

package provider

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
)

const (
	AD_UNKNOWN             = "unknown"
	AD_NOT_A_MEMBER        = "not_a_member"
	AD_JOINED              = "joined"
	AD_JOINING_IN_PROGRESS = "joining_in_progress"
	AD_LEAVING_IN_PROGRESS = "leaving_in_progress"
	AD_JOINED_FAILED       = "joined_failed"
	AD_LEAVE_FAILED        = "leave_failed"
)

var ActiveDirectorySchemaRef = is.NewSchemaReference(
	http.MethodPost,
	"activedirectory",
	http.MethodGet,
	"activedirectory",
)

type ActiveDirectory struct {
	tfstate *is.TFState
}

func (m *ActiveDirectory) NewResourceManager(raw map[string]attr.Value, schema any) ResourceManager {
	return &ActiveDirectory{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef:        ActiveDirectorySchemaRef,
			SearchableFields: []string{"ldap_id", "domain_name", "machine_account_name"},
			EditOnlyFields:   []string{"admin_username", "admin_passwd", "enabled"},
			AdditionalSchemaAttributes: map[string]any{
				"admin_username": rschema.StringAttribute{
					Optional:    true,
					Description: "An Active Directory admin user with permission to join the Active Directory server.",
				},
				"admin_passwd": rschema.StringAttribute{
					Sensitive:   true,
					Optional:    true,
					Description: "The password for the specified Active Directory admin user.",
				},
				"enabled": rschema.BoolAttribute{
					Optional:    true,
					Description: "Set to true to join Active Directory. Set to false to leave Active Directory.",
				},
			},
		},
	)}
}

func (m *ActiveDirectory) NewDatasourceManager(raw map[string]attr.Value, schema any) DataSourceManager {
	return &ActiveDirectory{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef: ActiveDirectorySchemaRef,
		},
	)}
}

func (m *ActiveDirectory) TfState() *is.TFState {
	return m.tfstate
}

func (m *ActiveDirectory) API(rest *VMSRest) VastResourceAPIWithContext {
	return rest.ActiveDirectories
}

func (m *ActiveDirectory) AfterCreateResource(ctx context.Context, rest *VMSRest, record Record) error {
	stateTs := m.TfState()
	return m.processADJoinLeave(ctx, stateTs, nil, rest, record, nil)
}

func (m *ActiveDirectory) AfterUpdateResource(ctx context.Context, plan AfterUpdateResource, rest *VMSRest, record Record) error {
	stateTs := m.tfstate
	planTs := plan.(*ActiveDirectory).TfState()
	return m.processADJoinLeave(ctx, stateTs, planTs, rest, record, nil)
}

func (m *ActiveDirectory) PrepareDeleteResource(ctx context.Context, rest *VMSRest) error {
	stateTs := m.TfState()

	searchParams := getSearchParams(ctx, stateTs, nil)
	record, err := rest.ActiveDirectories.GetWithContext(ctx, searchParams)
	if err != nil {
		if isNotFoundErr(err) {
			return nil
		}
		return err
	}

	// Force enabled=false to trigger leave operation before deletion
	forceLeave := false
	return m.processADJoinLeave(ctx, stateTs, nil, rest, record, &forceLeave)
}

// getADState returns the current state of the Active Directory in lowercase
func (m *ActiveDirectory) getADState(record Record) string {
	if state, ok := record["state"].(string); ok {
		return strings.ToLower(state)
	}
	return AD_UNKNOWN
}

// shouldJoinAD checks if we should attempt to join AD based on current state
func (m *ActiveDirectory) shouldJoinAD(currentState string, enabled bool) bool {
	if !enabled {
		return false
	}
	// Don't join if already joined or joining is in progress
	if currentState == AD_JOINED || currentState == AD_JOINING_IN_PROGRESS {
		return false
	}
	// Join if not a member, join failed, or left AD
	return currentState == AD_NOT_A_MEMBER ||
		currentState == AD_JOINED_FAILED ||
		currentState == AD_LEAVE_FAILED ||
		currentState == AD_UNKNOWN
}

// shouldLeaveAD checks if we should attempt to leave AD based on current state
func (m *ActiveDirectory) shouldLeaveAD(currentState string, enabled bool) bool {
	if enabled {
		return false
	}
	// Don't leave if not joined or leaving is already in progress
	if currentState == AD_NOT_A_MEMBER ||
		currentState == AD_LEAVING_IN_PROGRESS ||
		currentState == AD_UNKNOWN {
		return false
	}
	// Leave if joined or join failed
	return currentState == AD_JOINED || currentState == AD_JOINED_FAILED
}

// processADJoinLeave handles the join or leave operation for Active Directory
// Takes both state and plan TFStates (plan can be nil) and uses fallback pattern
// to get values from plan first, then state
// forceEnabled can be used to override the enabled value (e.g., for deletion to force leave)
func (m *ActiveDirectory) processADJoinLeave(ctx context.Context, stateTs, planTs *is.TFState, rest *VMSRest, record Record, forceEnabled *bool) error {
	// Get enabled from plan first, fallback to state, or use forced value
	var enabled bool
	var hasEnabled bool

	if forceEnabled != nil {
		enabled = *forceEnabled
		hasEnabled = true
	} else if planTs != nil {
		enabled, hasEnabled = planTs.BoolWithFallback(stateTs, "enabled")
	} else {
		if stateTs.IsKnownAndNotNull("enabled") {
			enabled = stateTs.Bool("enabled")
			hasEnabled = true
		}
	}

	if !hasEnabled {
		return nil
	}

	currentState := m.getADState(record)

	// Determine if we should join or leave
	shouldPerformOperation := m.shouldJoinAD(currentState, enabled) || m.shouldLeaveAD(currentState, enabled)

	if !shouldPerformOperation {
		// No operation needed, credentials don't matter
		return nil
	}

	// We need to perform join/leave operation, so credentials are required
	var username string
	var hasUsername bool
	var password string
	var hasPassword bool

	if planTs != nil {
		username, hasUsername = planTs.StringWithFallback(stateTs, "admin_username")
		password, hasPassword = planTs.StringWithFallback(stateTs, "admin_passwd")
	} else {
		if stateTs.IsKnownAndNotNull("admin_username") {
			username = stateTs.String("admin_username")
			hasUsername = true
		}
		if stateTs.IsKnownAndNotNull("admin_passwd") {
			password = stateTs.String("admin_passwd")
			hasPassword = true
		}
	}

	if !hasUsername || !hasPassword {
		return fmt.Errorf("admin credentials (admin_username and admin_passwd) are required to %s Active Directory, current AD state: %s",
			map[bool]string{true: "join", false: "leave"}[enabled], currentState)
	}

	// Perform the join/leave operation
	params := params{
		"enabled":        enabled,
		"admin_username": username,
		"admin_passwd":   password,
	}
	id := record.RecordID()
	record, err := rest.ActiveDirectories.UpdateWithContext(ctx, id, params)
	if err != nil {
		return err
	}
	if err := handleMaybeAsyncTask(ctx, rest, record, nil); err != nil {
		return err
	}
	return nil
}
