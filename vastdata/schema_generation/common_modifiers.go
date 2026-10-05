// Copyright (c) HashiCorp, Inc.

package schema_generation

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/float64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

const (
	ModifierForceNew             = "force_new"
	ModifierForceNewIfConfigured = "force_new_if_configured"
)

// Common plan modifiers for string attributes (e.g., "path", "name")
var commonStringModifiers = map[string][]planmodifier.String{
	ModifierForceNew:             {stringplanmodifier.RequiresReplace()},
	ModifierForceNewIfConfigured: {stringplanmodifier.RequiresReplaceIfConfigured()},
}

// Common plan modifiers for int64 attributes (e.g., "retry_count")
var commonIntModifiers = map[string][]planmodifier.Int64{
	ModifierForceNew:             {int64planmodifier.RequiresReplace()},
	ModifierForceNewIfConfigured: {int64planmodifier.RequiresReplaceIfConfigured()},
}

// Common plan modifiers for float64 attributes (e.g., "threshold")
var commonFloatModifiers = map[string][]planmodifier.Float64{
	ModifierForceNew:             {float64planmodifier.RequiresReplace()},
	ModifierForceNewIfConfigured: {float64planmodifier.RequiresReplaceIfConfigured()},
}
