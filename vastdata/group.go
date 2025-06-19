package vastdata

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type Group struct {
	// A unique identifier of the group
	Id types.Int64 `tfsdk:"id" datasource:"optional"`
	// A unique GUID assigned to the group
	Guid types.String `tfsdk:"guid" datasource:"optional"`
	// A unique name given to the group
	Name types.String `tfsdk:"name" datasource:"optional" resource:"required,searchable"`
	// The group's Linux GID
	Gid types.Int64 `tfsdk:"gid" resource:"required"`
	// The group's SID
	Sid types.String `tfsdk:"sid"`
	// List of S3 policy IDs
	S3PoliciesIds types.Set `tfsdk:"s3_policies_ids" resource:"optional" element_type:"int64"`
}
