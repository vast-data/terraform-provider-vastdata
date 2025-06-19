package vastdata

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type User struct {
	// A uniq id given to user
	Id types.Int64 `tfsdk:"id" datasource:"optional"`
	// A uniq guid given to the user
	Guid types.String `tfsdk:"guid" datasource:"optional"`
	// A uniq name given to the user
	Name types.String `tfsdk:"name" datasource:"optional" resource:"required,searchable"`
	// The user unix UID
	Uid types.Int64 `tfsdk:"uid" datasource:"optional" resource:"optional"`
	// The user leading unix GID
	LeadingGid types.Int64 `tfsdk:"leading_gid" resource:"optional"`
	// List of supplementary GID list
	Gids types.Set `tfsdk:"gids" resource:"optional" element_type:"int64"`
	// List of supplementary Group list
	Groups types.Set `tfsdk:"groups" resource:"optional" element_type:"string"`
	// Group Count
	GroupCount types.Int64 `tfsdk:"group_count"`
	// Leading Group Name
	LeadingGroupName types.String `tfsdk:"leading_group_name"`
	// Leading Group GID
	LeadingGroupGid types.Int64 `tfsdk:"leading_group_gid"`
	// The user SID
	Sid types.String `tfsdk:"sid"`
	// The user primary group SID
	PrimaryGroupSid types.String `tfsdk:"primary_group_sid" resource:"optional"`
	// supplementary SID list
	Sids types.Set `tfsdk:"sids" element_type:"string"`
	// Is this a local user
	Local types.Bool `tfsdk:"local"`
	// Allow to create bucket
	AllowCreateBucket types.Bool `tfsdk:"allow_create_bucket" resource:"optional"`
	// Allow delete bucket
	AllowDeleteBucket types.Bool `tfsdk:"allow_delete_bucket" resource:"optional"`
	// Is S3 superuser
	S3Superuser types.Bool `tfsdk:"s3_superuser" resource:"optional"`
	// List S3 policies IDs
	S3PoliciesIds types.Set `tfsdk:"s3_policies_ids" resource:"optional" element_type:"int64"`
}
