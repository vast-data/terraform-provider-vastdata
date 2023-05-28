package resources

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api_latest "github.com/vast-data/terraform-provider-vastdata.git/codegen/latest"
	metadata "github.com/vast-data/terraform-provider-vastdata.git/metadata"
	utils "github.com/vast-data/terraform-provider-vastdata.git/utils"
	vast_client "github.com/vast-data/terraform-provider-vastdata.git/vast-client"
	vast_versions "github.com/vast-data/terraform-provider-vastdata.git/vast_versions"
	"io"
	"net/url"
	"reflect"
	"strconv"
)

func ResourceViewPolicy() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceViewPolicyRead,
		DeleteContext: resourceViewPolicyDelete,
		CreateContext: resourceViewPolicyCreate,
		UpdateContext: resourceViewPolicyUpdate,
		Importer: &schema.ResourceImporter{
			StateContext: resourceViewPolicyImporter,
		},
		Schema: getResourceViewPolicySchema(),
	}
}

func getResourceViewPolicySchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{

		"guid": &schema.Schema{
			Type: schema.TypeString,

			Computed: true,
			Optional: true,
		},

		"name": &schema.Schema{
			Type: schema.TypeString,

			Required: true,
		},

		"gid_inheritance": &schema.Schema{
			Type: schema.TypeString,

			Computed: true,
			Optional: true,
		},

		"flavor": &schema.Schema{
			Type: schema.TypeString,

			Computed: true,
			Optional: true,
		},

		"access_flavor": &schema.Schema{
			Type: schema.TypeString,

			Computed: true,
			Optional: true,
		},

		"path_length": &schema.Schema{
			Type: schema.TypeString,

			Computed: true,
			Optional: true,
		},

		"allowed_characters": &schema.Schema{
			Type: schema.TypeString,

			Computed: true,
			Optional: true,
		},

		"use32bit_fileid": &schema.Schema{
			Type: schema.TypeBool,

			Computed: true,
			Optional: true,
		},

		"expose_id_in_fsid": &schema.Schema{
			Type: schema.TypeBool,

			Computed: true,
			Optional: true,
		},

		"use_auth_provider": &schema.Schema{
			Type: schema.TypeBool,

			Computed: true,
			Optional: true,
		},

		"auth_source": &schema.Schema{
			Type: schema.TypeString,

			Computed: true,
			Optional: true,
		},

		"read_write": &schema.Schema{
			Type: schema.TypeList,

			Computed: true,
			Optional: true,

			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},

		"read_only": &schema.Schema{
			Type: schema.TypeList,

			Computed: true,
			Optional: true,

			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},

		"nfs_read_write": &schema.Schema{
			Type: schema.TypeList,

			Computed: true,
			Optional: true,

			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},

		"nfs_read_only": &schema.Schema{
			Type: schema.TypeList,

			Computed: true,
			Optional: true,

			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},

		"smb_read_write": &schema.Schema{
			Type: schema.TypeList,

			Computed: true,
			Optional: true,

			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},

		"smb_read_only": &schema.Schema{
			Type: schema.TypeList,

			Computed: true,
			Optional: true,

			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},

		"s3_read_write": &schema.Schema{
			Type: schema.TypeList,

			Computed: true,
			Optional: true,

			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},

		"s3_read_only": &schema.Schema{
			Type: schema.TypeList,

			Computed: true,
			Optional: true,

			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},

		"trash_access": &schema.Schema{
			Type: schema.TypeList,

			Computed: true,
			Optional: true,

			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},

		"nfs_posix_acl": &schema.Schema{
			Type: schema.TypeBool,

			Computed: true,
			Optional: true,
		},

		"nfs_return_open_permissions": &schema.Schema{
			Type: schema.TypeBool,

			Computed: true,
			Optional: true,
		},

		"nfs_no_squash": &schema.Schema{
			Type: schema.TypeList,

			Computed: true,
			Optional: true,

			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},

		"nfs_root_squash": &schema.Schema{
			Type: schema.TypeList,

			Computed: true,
			Optional: true,

			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},

		"nfs_all_squash": &schema.Schema{
			Type: schema.TypeList,

			Computed: true,
			Optional: true,

			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},

		"s3_bucket_full_control": &schema.Schema{
			Type: schema.TypeString,

			Computed: true,
			Optional: true,
		},

		"s3_bucket_listing": &schema.Schema{
			Type: schema.TypeString,

			Computed: true,
			Optional: true,
		},

		"s3_bucket_read": &schema.Schema{
			Type: schema.TypeString,

			Computed: true,
			Optional: true,
		},

		"s3_bucket_read_acp": &schema.Schema{
			Type: schema.TypeString,

			Computed: true,
			Optional: true,
		},

		"s3_bucket_write": &schema.Schema{
			Type: schema.TypeString,

			Computed: true,
			Optional: true,
		},

		"s3_bucket_write_acp": &schema.Schema{
			Type: schema.TypeString,

			Computed: true,
			Optional: true,
		},

		"s3_object_full_control": &schema.Schema{
			Type: schema.TypeString,

			Computed: true,
			Optional: true,
		},

		"s3_object_read": &schema.Schema{
			Type: schema.TypeString,

			Computed: true,
			Optional: true,
		},

		"s3_object_read_acp": &schema.Schema{
			Type: schema.TypeString,

			Computed: true,
			Optional: true,
		},

		"s3_object_write": &schema.Schema{
			Type: schema.TypeString,

			Computed: true,
			Optional: true,
		},

		"s3_object_write_acp": &schema.Schema{
			Type: schema.TypeString,

			Computed: true,
			Optional: true,
		},

		"smb_file_mode": &schema.Schema{
			Type: schema.TypeInt,

			Computed: true,
			Optional: true,
		},

		"smb_directory_mode": &schema.Schema{
			Type: schema.TypeInt,

			Computed: true,
			Optional: true,
		},

		"smb_file_mode_padded": &schema.Schema{
			Type: schema.TypeString,

			Computed: true,
			Optional: true,
		},

		"smb_directory_mode_padded": &schema.Schema{
			Type: schema.TypeString,

			Computed: true,
			Optional: true,
		},

		"cluster": &schema.Schema{
			Type: schema.TypeString,

			Computed: true,
			Optional: true,
		},

		"cluster_id": &schema.Schema{
			Type: schema.TypeInt,

			Computed: true,
			Optional: true,
		},

		"tenant_id": &schema.Schema{
			Type: schema.TypeInt,

			Computed: true,
			Optional: true,
		},

		"tenant_name": &schema.Schema{
			Type: schema.TypeString,

			Computed: true,
			Optional: true,
		},

		"url": &schema.Schema{
			Type: schema.TypeString,

			Computed: true,
			Optional: true,
		},

		"atime_frequency": &schema.Schema{
			Type: schema.TypeString,

			Computed: true,
			Optional: true,
		},

		"sync": &schema.Schema{
			Type: schema.TypeString,

			Computed: true,
			Optional: true,
		},

		"vip_pools": &schema.Schema{
			Type: schema.TypeList,

			Computed: true,
			Optional: true,

			Elem: &schema.Schema{
				Type: schema.TypeInt,
			},
		},

		"sync_time": &schema.Schema{
			Type: schema.TypeString,

			Computed: true,
			Optional: true,
		},

		"nfs_minimal_protection_level": &schema.Schema{
			Type: schema.TypeString,

			Computed: true,
			Optional: true,
		},

		"s3_visibility": &schema.Schema{
			Type: schema.TypeList,

			Computed: true,
			Optional: true,

			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},

		"s3_visibility_groups": &schema.Schema{
			Type: schema.TypeList,

			Computed: true,
			Optional: true,

			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},

		"apple_sid": &schema.Schema{
			Type: schema.TypeBool,

			Computed: true,
			Optional: true,
		},

		"protocols": &schema.Schema{
			Type: schema.TypeList,

			Computed: true,
			Optional: true,

			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},

		"data_create_delete": &schema.Schema{
			Type: schema.TypeBool,

			Computed: true,
			Optional: true,
		},

		"data_modify": &schema.Schema{
			Type: schema.TypeBool,

			Computed: true,
			Optional: true,
		},

		"data_read": &schema.Schema{
			Type: schema.TypeBool,

			Computed: true,
			Optional: true,
		},

		"log_full_path": &schema.Schema{
			Type: schema.TypeBool,

			Computed: true,
			Optional: true,
		},

		"log_hostname": &schema.Schema{
			Type: schema.TypeBool,

			Computed: true,
			Optional: true,
		},

		"log_username": &schema.Schema{
			Type: schema.TypeBool,

			Computed: true,
			Optional: true,
		},

		"log_deleted": &schema.Schema{
			Type: schema.TypeBool,

			Computed: true,
			Optional: true,
		},

		"count_views": &schema.Schema{
			Type: schema.TypeInt,

			Computed: true,
			Optional: true,
		},

		"enable_snapshot_lookup": &schema.Schema{
			Type: schema.TypeBool,

			Computed: true,
			Optional: true,
		},

		"enable_listing_of_snapshot_dir": &schema.Schema{
			Type: schema.TypeBool,

			Computed: true,
			Optional: true,
		},
	}
}

var ViewPolicy_names_mapping map[string][]string = map[string][]string{}

func ResourceViewPolicyReadStructIntoSchema(ctx context.Context, resource api_latest.ViewPolicy, d *schema.ResourceData) diag.Diagnostics {
	var diags diag.Diagnostics
	var err error

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "Guid", resource.Guid))

	err = d.Set("guid", resource.Guid)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"guid\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "Name", resource.Name))

	err = d.Set("name", resource.Name)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"name\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "GidInheritance", resource.GidInheritance))

	err = d.Set("gid_inheritance", resource.GidInheritance)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"gid_inheritance\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "Flavor", resource.Flavor))

	err = d.Set("flavor", resource.Flavor)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"flavor\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "AccessFlavor", resource.AccessFlavor))

	err = d.Set("access_flavor", resource.AccessFlavor)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"access_flavor\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "PathLength", resource.PathLength))

	err = d.Set("path_length", resource.PathLength)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"path_length\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "AllowedCharacters", resource.AllowedCharacters))

	err = d.Set("allowed_characters", resource.AllowedCharacters)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"allowed_characters\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "Use32bitFileid", resource.Use32bitFileid))

	err = d.Set("use32bit_fileid", resource.Use32bitFileid)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"use32bit_fileid\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "ExposeIdInFsid", resource.ExposeIdInFsid))

	err = d.Set("expose_id_in_fsid", resource.ExposeIdInFsid)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"expose_id_in_fsid\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "UseAuthProvider", resource.UseAuthProvider))

	err = d.Set("use_auth_provider", resource.UseAuthProvider)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"use_auth_provider\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "AuthSource", resource.AuthSource))

	err = d.Set("auth_source", resource.AuthSource)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"auth_source\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "ReadWrite", resource.ReadWrite))

	err = d.Set("read_write", utils.FlattenListOfPrimitives(&resource.ReadWrite))

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"read_write\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "ReadOnly", resource.ReadOnly))

	err = d.Set("read_only", utils.FlattenListOfPrimitives(&resource.ReadOnly))

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"read_only\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "NfsReadWrite", resource.NfsReadWrite))

	err = d.Set("nfs_read_write", utils.FlattenListOfPrimitives(&resource.NfsReadWrite))

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"nfs_read_write\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "NfsReadOnly", resource.NfsReadOnly))

	err = d.Set("nfs_read_only", utils.FlattenListOfPrimitives(&resource.NfsReadOnly))

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"nfs_read_only\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "SmbReadWrite", resource.SmbReadWrite))

	err = d.Set("smb_read_write", utils.FlattenListOfPrimitives(&resource.SmbReadWrite))

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"smb_read_write\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "SmbReadOnly", resource.SmbReadOnly))

	err = d.Set("smb_read_only", utils.FlattenListOfPrimitives(&resource.SmbReadOnly))

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"smb_read_only\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "S3ReadWrite", resource.S3ReadWrite))

	err = d.Set("s3_read_write", utils.FlattenListOfPrimitives(&resource.S3ReadWrite))

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"s3_read_write\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "S3ReadOnly", resource.S3ReadOnly))

	err = d.Set("s3_read_only", utils.FlattenListOfPrimitives(&resource.S3ReadOnly))

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"s3_read_only\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "TrashAccess", resource.TrashAccess))

	err = d.Set("trash_access", utils.FlattenListOfPrimitives(&resource.TrashAccess))

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"trash_access\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "NfsPosixAcl", resource.NfsPosixAcl))

	err = d.Set("nfs_posix_acl", resource.NfsPosixAcl)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"nfs_posix_acl\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "NfsReturnOpenPermissions", resource.NfsReturnOpenPermissions))

	err = d.Set("nfs_return_open_permissions", resource.NfsReturnOpenPermissions)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"nfs_return_open_permissions\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "NfsNoSquash", resource.NfsNoSquash))

	err = d.Set("nfs_no_squash", utils.FlattenListOfPrimitives(&resource.NfsNoSquash))

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"nfs_no_squash\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "NfsRootSquash", resource.NfsRootSquash))

	err = d.Set("nfs_root_squash", utils.FlattenListOfPrimitives(&resource.NfsRootSquash))

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"nfs_root_squash\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "NfsAllSquash", resource.NfsAllSquash))

	err = d.Set("nfs_all_squash", utils.FlattenListOfPrimitives(&resource.NfsAllSquash))

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"nfs_all_squash\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "S3BucketFullControl", resource.S3BucketFullControl))

	err = d.Set("s3_bucket_full_control", resource.S3BucketFullControl)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"s3_bucket_full_control\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "S3BucketListing", resource.S3BucketListing))

	err = d.Set("s3_bucket_listing", resource.S3BucketListing)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"s3_bucket_listing\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "S3BucketRead", resource.S3BucketRead))

	err = d.Set("s3_bucket_read", resource.S3BucketRead)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"s3_bucket_read\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "S3BucketReadAcp", resource.S3BucketReadAcp))

	err = d.Set("s3_bucket_read_acp", resource.S3BucketReadAcp)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"s3_bucket_read_acp\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "S3BucketWrite", resource.S3BucketWrite))

	err = d.Set("s3_bucket_write", resource.S3BucketWrite)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"s3_bucket_write\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "S3BucketWriteAcp", resource.S3BucketWriteAcp))

	err = d.Set("s3_bucket_write_acp", resource.S3BucketWriteAcp)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"s3_bucket_write_acp\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "S3ObjectFullControl", resource.S3ObjectFullControl))

	err = d.Set("s3_object_full_control", resource.S3ObjectFullControl)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"s3_object_full_control\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "S3ObjectRead", resource.S3ObjectRead))

	err = d.Set("s3_object_read", resource.S3ObjectRead)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"s3_object_read\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "S3ObjectReadAcp", resource.S3ObjectReadAcp))

	err = d.Set("s3_object_read_acp", resource.S3ObjectReadAcp)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"s3_object_read_acp\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "S3ObjectWrite", resource.S3ObjectWrite))

	err = d.Set("s3_object_write", resource.S3ObjectWrite)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"s3_object_write\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "S3ObjectWriteAcp", resource.S3ObjectWriteAcp))

	err = d.Set("s3_object_write_acp", resource.S3ObjectWriteAcp)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"s3_object_write_acp\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "SmbFileMode", resource.SmbFileMode))

	err = d.Set("smb_file_mode", resource.SmbFileMode)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"smb_file_mode\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "SmbDirectoryMode", resource.SmbDirectoryMode))

	err = d.Set("smb_directory_mode", resource.SmbDirectoryMode)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"smb_directory_mode\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "SmbFileModePadded", resource.SmbFileModePadded))

	err = d.Set("smb_file_mode_padded", resource.SmbFileModePadded)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"smb_file_mode_padded\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "SmbDirectoryModePadded", resource.SmbDirectoryModePadded))

	err = d.Set("smb_directory_mode_padded", resource.SmbDirectoryModePadded)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"smb_directory_mode_padded\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "Cluster", resource.Cluster))

	err = d.Set("cluster", resource.Cluster)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"cluster\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "ClusterId", resource.ClusterId))

	err = d.Set("cluster_id", resource.ClusterId)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"cluster_id\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "TenantId", resource.TenantId))

	err = d.Set("tenant_id", resource.TenantId)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"tenant_id\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "TenantName", resource.TenantName))

	err = d.Set("tenant_name", resource.TenantName)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"tenant_name\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "Url", resource.Url))

	err = d.Set("url", resource.Url)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"url\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "AtimeFrequency", resource.AtimeFrequency))

	err = d.Set("atime_frequency", resource.AtimeFrequency)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"atime_frequency\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "Sync", resource.Sync))

	err = d.Set("sync", resource.Sync)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"sync\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "VipPools", resource.VipPools))

	err = d.Set("vip_pools", utils.FlattenListOfPrimitives(&resource.VipPools))

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"vip_pools\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "SyncTime", resource.SyncTime))

	err = d.Set("sync_time", resource.SyncTime)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"sync_time\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "NfsMinimalProtectionLevel", resource.NfsMinimalProtectionLevel))

	err = d.Set("nfs_minimal_protection_level", resource.NfsMinimalProtectionLevel)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"nfs_minimal_protection_level\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "S3Visibility", resource.S3Visibility))

	err = d.Set("s3_visibility", utils.FlattenListOfPrimitives(&resource.S3Visibility))

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"s3_visibility\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "S3VisibilityGroups", resource.S3VisibilityGroups))

	err = d.Set("s3_visibility_groups", utils.FlattenListOfPrimitives(&resource.S3VisibilityGroups))

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"s3_visibility_groups\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "AppleSid", resource.AppleSid))

	err = d.Set("apple_sid", resource.AppleSid)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"apple_sid\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "Protocols", resource.Protocols))

	err = d.Set("protocols", utils.FlattenListOfPrimitives(&resource.Protocols))

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"protocols\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "DataCreateDelete", resource.DataCreateDelete))

	err = d.Set("data_create_delete", resource.DataCreateDelete)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"data_create_delete\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "DataModify", resource.DataModify))

	err = d.Set("data_modify", resource.DataModify)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"data_modify\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "DataRead", resource.DataRead))

	err = d.Set("data_read", resource.DataRead)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"data_read\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "LogFullPath", resource.LogFullPath))

	err = d.Set("log_full_path", resource.LogFullPath)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"log_full_path\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "LogHostname", resource.LogHostname))

	err = d.Set("log_hostname", resource.LogHostname)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"log_hostname\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "LogUsername", resource.LogUsername))

	err = d.Set("log_username", resource.LogUsername)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"log_username\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "LogDeleted", resource.LogDeleted))

	err = d.Set("log_deleted", resource.LogDeleted)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"log_deleted\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "CountViews", resource.CountViews))

	err = d.Set("count_views", resource.CountViews)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"count_views\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "EnableSnapshotLookup", resource.EnableSnapshotLookup))

	err = d.Set("enable_snapshot_lookup", resource.EnableSnapshotLookup)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"enable_snapshot_lookup\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "EnableListingOfSnapshotDir", resource.EnableListingOfSnapshotDir))

	err = d.Set("enable_listing_of_snapshot_dir", resource.EnableListingOfSnapshotDir)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"enable_listing_of_snapshot_dir\"",
			Detail:   err.Error(),
		})
	}

	return diags

}
func resourceViewPolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(vast_client.JwtSession)

	ViewPolicyId := d.Id()
	response, err := client.Get(ctx, fmt.Sprintf("/api/viewpolicies/%v", ViewPolicyId), "", map[string]string{})

	utils.VastVersionsWarn(ctx)

	tflog.Info(ctx, response.Request.URL.String())
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured while obtaining data from the vastdata cluster",
			Detail:   err.Error(),
		})
		return diags

	}
	resource := api_latest.ViewPolicy{}
	body, err := utils.DefaultProcessingFunc(ctx, response)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured reading data recived from VastData cluster",
			Detail:   err.Error(),
		})
		return diags

	}
	err = json.Unmarshal(body, &resource)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured while parsing data recived from VastData cluster",
			Detail:   err.Error(),
		})
		return diags

	}
	diags = ResourceViewPolicyReadStructIntoSchema(ctx, resource, d)
	return diags
}

func resourceViewPolicyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(vast_client.JwtSession)

	ViewPolicyId := d.Id()
	response, err := client.Delete(ctx, fmt.Sprintf("/api/viewpolicies/%v/", ViewPolicyId), "", map[string]string{})
	tflog.Info(ctx, fmt.Sprintf("Removing Resource"))
	tflog.Info(ctx, response.Request.URL.String())
	tflog.Info(ctx, utils.GetResponseBodyAsStr(response))

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured while deleting a resource from the vastdata cluster",
			Detail:   err.Error(),
		})

	}

	return diags

}

func resourceViewPolicyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	names_mapping := utils.ContextKey("names_mapping")
	new_ctx := context.WithValue(ctx, names_mapping, ViewPolicy_names_mapping)

	var diags diag.Diagnostics
	data := make(map[string]interface{})
	client := m.(vast_client.JwtSession)
	tflog.Info(ctx, fmt.Sprintf("Creating Resource ViewPolicy"))
	reflect_ViewPolicy := reflect.TypeOf((*api_latest.ViewPolicy)(nil))
	utils.PopulateResourceMap(new_ctx, reflect_ViewPolicy.Elem(), d, &data, "", false)

	version_compare := utils.VastVersionsWarn(ctx)

	if version_compare != metadata.CLUSTER_VERSION_EQUALS {
		cluster_version := metadata.ClusterVersionString()
		t, t_exists := vast_versions.GetVersionedType(cluster_version, "ViewPolicy")
		if t_exists {
			versions_error := utils.VersionMatch(t, data)
			if versions_error != nil {
				tflog.Warn(ctx, versions_error.Error())
				version_validation_mode, version_validation_mode_exists := metadata.GetClusterConfig("version_validation_mode")
				tflog.Warn(ctx, fmt.Sprintf("Version Validation Mode Detected %s", version_validation_mode))
				if version_validation_mode_exists && version_validation_mode == "strict" {
					diags = append(diags, diag.Diagnostic{
						Severity: diag.Error,
						Summary:  "Cluster Version & Build Version Are Too Differant",
						Detail:   versions_error.Error(),
					})
					return diags
				}
			}
		} else {
			tflog.Warn(ctx, fmt.Sprintf("Could have not found resource %s in version %s , things might not work properly", "ViewPolicy", cluster_version))
		}
	}
	tflog.Debug(ctx, fmt.Sprintf("Data %v", data))
	b, err := json.MarshalIndent(data, "", "   ")
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Could have not generate request json",
			Detail:   err.Error(),
		})
		return diags
	}
	tflog.Debug(ctx, fmt.Sprintf("Request json created %v", string(b)))

	response, create_err := client.Post(ctx, "/api/viewpolicies/", bytes.NewReader(b), map[string]string{})
	tflog.Info(ctx, fmt.Sprintf("Server Error for  ViewPolicy %v", create_err))

	if create_err != nil {
		error_message := create_err.Error() + " Server Response: " + utils.GetResponseBodyAsStr(response)
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Object Creation Failed",
			Detail:   error_message,
		})
		return diags
	}
	response_body, _ := io.ReadAll(response.Body)
	tflog.Debug(ctx, fmt.Sprintf("Object created , server response %v", string(response_body)))
	resource := api_latest.ViewPolicy{}
	err = json.Unmarshal(response_body, &resource)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to convert response body into ViewPolicy",
			Detail:   err.Error(),
		})
		return diags
	}

	d.SetId(strconv.FormatInt((int64)(resource.Id), 10))
	resourceViewPolicyRead(ctx, d, m)
	return diags
}

func resourceViewPolicyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	names_mapping := utils.ContextKey("names_mapping")
	new_ctx := context.WithValue(ctx, names_mapping, ViewPolicy_names_mapping)

	var diags diag.Diagnostics
	data := make(map[string]interface{})
	version_compare := utils.VastVersionsWarn(ctx)
	if version_compare != metadata.CLUSTER_VERSION_EQUALS {
		cluster_version := metadata.ClusterVersionString()
		t, t_exists := vast_versions.GetVersionedType(cluster_version, "ViewPolicy")
		if t_exists {
			versions_error := utils.VersionMatch(t, data)
			if versions_error != nil {
				tflog.Warn(ctx, versions_error.Error())
				version_validation_mode, version_validation_mode_exists := metadata.GetClusterConfig("version_validation_mode")
				tflog.Warn(ctx, fmt.Sprintf("Version Validation Mode Detected %s", version_validation_mode))
				if version_validation_mode_exists && version_validation_mode == "strict" {
					diags = append(diags, diag.Diagnostic{
						Severity: diag.Error,
						Summary:  "Cluster Version & Build Version Are Too Differant",
						Detail:   versions_error.Error(),
					})
					return diags
				}
			}
		} else {
			tflog.Warn(ctx, fmt.Sprintf("Could have not found resource %s in version %s , things might not work properly", "ViewPolicy", cluster_version))
		}
	}

	client := m.(vast_client.JwtSession)
	ViewPolicyId := d.Id()
	tflog.Info(ctx, fmt.Sprintf("Updating Resource ViewPolicy"))
	reflect_ViewPolicy := reflect.TypeOf((*api_latest.ViewPolicy)(nil))
	utils.PopulateResourceMap(new_ctx, reflect_ViewPolicy.Elem(), d, &data, "", false)
	tflog.Debug(ctx, fmt.Sprintf("Data %v", data))
	b, err := json.MarshalIndent(data, "", "   ")
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Could have not generate request json",
			Detail:   err.Error(),
		})
		return diags
	}
	tflog.Debug(ctx, fmt.Sprintf("Request json created %v", string(b)))
	response, patch_err := client.Patch(ctx, fmt.Sprintf("/api/viewpolicies//%v", ViewPolicyId), "application/json", bytes.NewReader(b), map[string]string{})
	tflog.Info(ctx, fmt.Sprintf("Server Error for  ViewPolicy %v", patch_err))
	if patch_err != nil {
		error_message := patch_err.Error() + " Server Response: " + utils.GetResponseBodyAsStr(response)
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Object Creation Failed",
			Detail:   error_message,
		})
		return diags
	}
	resourceViewPolicyRead(ctx, d, m)
	return diags

}

func resourceViewPolicyImporter(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {

	result := []*schema.ResourceData{}
	client := m.(vast_client.JwtSession)
	guid := d.Id()
	values := url.Values{}
	values.Add("guid", fmt.Sprintf("%v", guid))

	response, err := client.Get(ctx, "/api/viewpolicies/", values.Encode(), map[string]string{})

	if err != nil {
		return result, err
	}

	resource_l := []api_latest.ViewPolicy{}

	body, err := utils.DefaultProcessingFunc(ctx, response)
	if err != nil {
		return result, err
	}
	err = json.Unmarshal(body, &resource_l)
	if err != nil {
		return result, err
	}

	if len(resource_l) == 0 {
		return result, errors.New("Cluster provided 0 elements matchin gthis guid")
	}

	resource := resource_l[0]

	Id := (int64)(resource.Id)
	d.SetId(strconv.FormatInt(Id, 10))
	diags := ResourceViewPolicyReadStructIntoSchema(ctx, resource, d)
	if diags.HasError() {
		all_errors := "Errors occured while importing:\n"
		for _, dig := range diags {
			all_errors += fmt.Sprintf("Summary:%s\nDetails:%s\n", dig.Summary, dig.Detail)
		}
		return result, errors.New(all_errors)
	}
	result = append(result, d)

	return result, err

}