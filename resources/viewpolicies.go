package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"reflect"

	//        "net/url"
	"errors"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api_latest "github.com/vast-data/terraform-provider-vastdata/codegen/latest"
	codegen_configs "github.com/vast-data/terraform-provider-vastdata/codegen_tools/configs"
	metadata "github.com/vast-data/terraform-provider-vastdata/metadata"
	utils "github.com/vast-data/terraform-provider-vastdata/utils"
	vast_client "github.com/vast-data/terraform-provider-vastdata/vast-client"
	vast_versions "github.com/vast-data/terraform-provider-vastdata/vast_versions"
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

		Description: ``,
		Schema:      getResourceViewPolicySchema(),
	}
}

func getResourceViewPolicySchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{

		"guid": &schema.Schema{
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("ViewPolicy").GetConflictingFields("guid"),

			Computed:    true,
			Optional:    false,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0) A uniqe guid given to the view policy`,
		},

		"name": &schema.Schema{
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("ViewPolicy").GetConflictingFields("name"),

			Required: true,
		},

		"gid_inheritance": &schema.Schema{
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("ViewPolicy").GetConflictingFields("gid_inheritance"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0) Determine the way a file inherits GID`,
		},

		"flavor": &schema.Schema{
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("ViewPolicy").GetConflictingFields("flavor"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0) Security flavor, which determines how file and directory permissions are applied in multiprotocol views.`,
		},

		"access_flavor": &schema.Schema{
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("ViewPolicy").GetConflictingFields("access_flavor"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0) Applicable with MIXED_LAST_WINS security flavor (Access can be set via NFSv3 regardless of this option)`,
		},

		"path_length": &schema.Schema{
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("ViewPolicy").GetConflictingFields("path_length"),

			Computed:  true,
			Optional:  true,
			Sensitive: false,

			ValidateDiagFunc: utils.OneOf([]string{"LCD", "NPL"}),
			Description:      `(Valid for versions: 5.0.0,5.1.0) How to determine the maximum allowed path length Allowed Values are [LCD NPL]`,
		},

		"allowed_characters": &schema.Schema{
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("ViewPolicy").GetConflictingFields("allowed_characters"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0) How to determine the allowed characters in a path`,
		},

		"use32bit_fileid": &schema.Schema{
			Type:          schema.TypeBool,
			ConflictsWith: codegen_configs.GetResourceByName("ViewPolicy").GetConflictingFields("use32bit_fileid"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: ``,
		},

		"expose_id_in_fsid": &schema.Schema{
			Type:          schema.TypeBool,
			ConflictsWith: codegen_configs.GetResourceByName("ViewPolicy").GetConflictingFields("expose_id_in_fsid"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0) `,
		},

		"use_auth_provider": &schema.Schema{
			Type:          schema.TypeBool,
			ConflictsWith: codegen_configs.GetResourceByName("ViewPolicy").GetConflictingFields("use_auth_provider"),

			Computed:    false,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0) Use configured Auth Provider(s) to enforce group permissions when set to true , if set to ture with out specifing auth_source , the auth_source set to "PROVIDERS". if set to false than auth_source set to RPC. Due to the nature or terrafrom simply changing use_auth_provider from false to true or the other way around will not change the value auth_source as terrafrom will keep hold on the previous value. therefor it is adviasable to always specify the value of auth_source`,

			Default: false,
		},

		"auth_source": &schema.Schema{
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("ViewPolicy").GetConflictingFields("auth_source"),

			Computed:  true,
			Optional:  true,
			Sensitive: false,

			ValidateDiagFunc: utils.OneOf([]string{"PROVIDERS", "RPC", "RPC_AND_PROVIDERS"}),
			Description:      `(Valid for versions: 5.0.0,5.1.0) The source of authentication Allowed Values are [PROVIDERS RPC RPC_AND_PROVIDERS]`,
		},

		"read_write": &schema.Schema{
			Type:          schema.TypeList,
			ConflictsWith: codegen_configs.GetResourceByName("ViewPolicy").GetConflictingFields("read_write"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0) Hosts with NFS read/write permissions`,

			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},

		"read_only": &schema.Schema{
			Type:          schema.TypeList,
			ConflictsWith: codegen_configs.GetResourceByName("ViewPolicy").GetConflictingFields("read_only"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0) Hosts with NFS read only permissions`,

			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},

		"nfs_read_write": &schema.Schema{
			Type:          schema.TypeList,
			ConflictsWith: codegen_configs.GetResourceByName("ViewPolicy").GetConflictingFields("nfs_read_write"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0) Hosts with NFS read/write permissions. when creating a new View Policy if the value is not set than an empty list is sent to the VastData cluster resulting in empty list of addresses However during update if nfs_all_squash is removed from the resource nothing is changed to preserve terraform default behaviour in such cases. If there is a need to change the value an empty list it must be secifed and set to [].`,

			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},

		"nfs_read_only": &schema.Schema{
			Type:          schema.TypeList,
			ConflictsWith: codegen_configs.GetResourceByName("ViewPolicy").GetConflictingFields("nfs_read_only"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0) Hosts with NFS read only permissions. when creating a new View Policy if the value is not set than an empty list is sent to the VastData cluster resulting in empty list of addresses However during update if nfs_all_squash is removed from the resource nothing is changed to preserve terraform default behaviour in such cases. If there is a need to change the value an empty list it must be secifed and set to [].`,

			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},

		"smb_read_write": &schema.Schema{
			Type:          schema.TypeList,
			ConflictsWith: codegen_configs.GetResourceByName("ViewPolicy").GetConflictingFields("smb_read_write"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0) Hosts with SMB read/write permissions. when creating a new View Policy if the value is not set than an empty list is sent to the VastData cluster resulting in empty list of addresses However during update if nfs_all_squash is removed from the resource nothing is changed to preserve terraform default behaviour in such cases. If there is a need to change the value an empty list it must be secifed and set to [].`,

			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},

		"smb_read_only": &schema.Schema{
			Type:          schema.TypeList,
			ConflictsWith: codegen_configs.GetResourceByName("ViewPolicy").GetConflictingFields("smb_read_only"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0) Hosts with SMB read only permissions. when creating a new View Policy if the value is not set than an empty list is sent to the VastData cluster resulting in empty list of addresses However during update if nfs_all_squash is removed from the resource nothing is changed to preserve terraform default behaviour in such cases. If there is a need to change the value an empty list it must be secifed and set to [].`,

			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},

		"s3_read_write": &schema.Schema{
			Type:          schema.TypeList,
			ConflictsWith: codegen_configs.GetResourceByName("ViewPolicy").GetConflictingFields("s3_read_write"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0) Hosts with S3 read/write permissions. when creating a new View Policy if the value is not set than an empty list is sent to the VastData cluster resulting in empty list of addresses However during update if nfs_all_squash is removed from the resource nothing is changed to preserve terraform default behaviour in such cases. If there is a need to change the value an empty list it must be secifed and set to [].`,

			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},

		"s3_read_only": &schema.Schema{
			Type:          schema.TypeList,
			ConflictsWith: codegen_configs.GetResourceByName("ViewPolicy").GetConflictingFields("s3_read_only"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0) Hosts with S3 read only permissions. when creating a new View Policy if the value is not set than an empty list is sent to the VastData cluster resulting in empty list of addresses However during update if nfs_all_squash is removed from the resource nothing is changed to preserve terraform default behaviour in such cases. If there is a need to change the value an empty list it must be secifed and set to [].`,

			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},

		"trash_access": &schema.Schema{
			Type:          schema.TypeList,
			ConflictsWith: codegen_configs.GetResourceByName("ViewPolicy").GetConflictingFields("trash_access"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0) Hosts with trash permissions`,

			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},

		"nfs_posix_acl": &schema.Schema{
			Type:          schema.TypeBool,
			ConflictsWith: codegen_configs.GetResourceByName("ViewPolicy").GetConflictingFields("nfs_posix_acl"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0) Enable POSIX ACL`,
		},

		"nfs_return_open_permissions": &schema.Schema{
			Type:          schema.TypeBool,
			ConflictsWith: codegen_configs.GetResourceByName("ViewPolicy").GetConflictingFields("nfs_return_open_permissions"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0) when using smb use open permissions for files`,
		},

		"nfs_no_squash": &schema.Schema{
			Type:          schema.TypeList,
			ConflictsWith: codegen_configs.GetResourceByName("ViewPolicy").GetConflictingFields("nfs_no_squash"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0) Hosts with no squash policy. when creating a new View Policy if the value is not set than an empty list is sent to the VastData cluster resulting in empty list of addresses However during update if nfs_all_squash is removed from the resource nothing is changed to preserve terraform default behaviour in such cases. If there is a need to change the value an empty list it must be secifed and set to [].`,

			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},

		"nfs_root_squash": &schema.Schema{
			Type:          schema.TypeList,
			ConflictsWith: codegen_configs.GetResourceByName("ViewPolicy").GetConflictingFields("nfs_root_squash"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0) Hosts with root squash policy. when creating a new View Policy if the value is not set than an empty list is sent to the VastData cluster resulting in empty list of addresses However during update if nfs_all_squash is removed from the resource nothing is changed to preserve terraform default behaviour in such cases. If there is a need to change the value an empty list it must be secifed and set to [].`,

			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},

		"nfs_all_squash": &schema.Schema{
			Type:          schema.TypeList,
			ConflictsWith: codegen_configs.GetResourceByName("ViewPolicy").GetConflictingFields("nfs_all_squash"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0) Hosts with all squash policy. when creating a new View Policy if the value is not set than an empty list is sent to the VastData cluster resulting in empty list of addresses However during update if nfs_all_squash is removed from the resource nothing is changed to preserve terraform default behaviour in such cases. If there is a need to change the value an empty list it must be secifed and set to []`,

			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},

		"s3_bucket_full_control": &schema.Schema{
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("ViewPolicy").GetConflictingFields("s3_bucket_full_control"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0) Hosts with full permissions`,
		},

		"s3_bucket_listing": &schema.Schema{
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("ViewPolicy").GetConflictingFields("s3_bucket_listing"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0) Hosts with full permissions`,
		},

		"s3_bucket_read": &schema.Schema{
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("ViewPolicy").GetConflictingFields("s3_bucket_read"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0) Hosts with full permissions`,
		},

		"s3_bucket_read_acp": &schema.Schema{
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("ViewPolicy").GetConflictingFields("s3_bucket_read_acp"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0) Hosts with full permissions`,
		},

		"s3_bucket_write": &schema.Schema{
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("ViewPolicy").GetConflictingFields("s3_bucket_write"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0) Hosts with full permissions`,
		},

		"s3_bucket_write_acp": &schema.Schema{
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("ViewPolicy").GetConflictingFields("s3_bucket_write_acp"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0) Hosts with full permissions`,
		},

		"s3_object_full_control": &schema.Schema{
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("ViewPolicy").GetConflictingFields("s3_object_full_control"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0) Hosts with full permissions`,
		},

		"s3_object_read": &schema.Schema{
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("ViewPolicy").GetConflictingFields("s3_object_read"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0) Hosts with full permissions`,
		},

		"s3_object_read_acp": &schema.Schema{
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("ViewPolicy").GetConflictingFields("s3_object_read_acp"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0) Hosts with full permissions`,
		},

		"s3_object_write": &schema.Schema{
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("ViewPolicy").GetConflictingFields("s3_object_write"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0) Hosts with full permissions`,
		},

		"s3_object_write_acp": &schema.Schema{
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("ViewPolicy").GetConflictingFields("s3_object_write_acp"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0) Hosts with full permissions`,
		},

		"smb_file_mode": &schema.Schema{
			Type:          schema.TypeInt,
			ConflictsWith: codegen_configs.GetResourceByName("ViewPolicy").GetConflictingFields("smb_file_mode"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0) Default unix type permissions on new file`,
		},

		"smb_directory_mode": &schema.Schema{
			Type:          schema.TypeInt,
			ConflictsWith: codegen_configs.GetResourceByName("ViewPolicy").GetConflictingFields("smb_directory_mode"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0) Default unix type permissions on new folder`,
		},

		"smb_file_mode_padded": &schema.Schema{
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("ViewPolicy").GetConflictingFields("smb_file_mode_padded"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0) Default unix type permissions on new file`,
		},

		"smb_directory_mode_padded": &schema.Schema{
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("ViewPolicy").GetConflictingFields("smb_directory_mode_padded"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0) Default unix type permissions on new folder`,
		},

		"cluster": &schema.Schema{
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("ViewPolicy").GetConflictingFields("cluster"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0) Parent Cluster`,
		},

		"cluster_id": &schema.Schema{
			Type:          schema.TypeInt,
			ConflictsWith: codegen_configs.GetResourceByName("ViewPolicy").GetConflictingFields("cluster_id"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0) Parent Cluster ID`,
		},

		"tenant_id": &schema.Schema{
			Type:          schema.TypeInt,
			ConflictsWith: codegen_configs.GetResourceByName("ViewPolicy").GetConflictingFields("tenant_id"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0) Tenant ID`,
		},

		"tenant_name": &schema.Schema{
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("ViewPolicy").GetConflictingFields("tenant_name"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0) Tenant Name`,
		},

		"url": &schema.Schema{
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("ViewPolicy").GetConflictingFields("url"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0) `,
		},

		"atime_frequency": &schema.Schema{
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("ViewPolicy").GetConflictingFields("atime_frequency"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0) Frequency for updating the atime attribute of NFS files. atime is updated on read operations if the difference between the current time and the file's atime value is greater than the atime frequency. Specify as time in seconds.`,
		},

		"vip_pools": &schema.Schema{
			Type:          schema.TypeList,
			ConflictsWith: codegen_configs.GetResourceByName("ViewPolicy").GetConflictingFields("vip_pools"),

			DiffSuppressOnRefresh: false,
			DiffSuppressFunc:      codegen_configs.GetResourceByName("ViewPolicy").GetAttributeDiffFunc("vip_pools"),
			Computed:              true,
			Optional:              true,
			Sensitive:             false,
			Description:           `(Valid for versions: 5.0.0,5.1.0) Comma separated vip pool ids, this attribute conflicts with vippool_permissions and can not be provided togather. Also due to the lack of ability to configure vippool permissions using this attibute , vippool permissions are always defined as read/write`,

			Elem: &schema.Schema{
				Type: schema.TypeInt,
			},
		},

		"nfs_minimal_protection_level": &schema.Schema{
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("ViewPolicy").GetConflictingFields("nfs_minimal_protection_level"),

			Computed:  true,
			Optional:  true,
			Sensitive: false,

			ValidateDiagFunc: utils.OneOf([]string{"NONE", "SYSTEM", "KRB_AUTH_ONLY", "KRB_INTEGRITY", "KRB_PRIVACY"}),
			Description:      `(Valid for versions: 5.0.0,5.1.0) NFS 4.1 minimal protection level Allowed Values are [NONE SYSTEM KRB_AUTH_ONLY KRB_INTEGRITY KRB_PRIVACY]`,
		},

		"s3_visibility": &schema.Schema{
			Type:          schema.TypeList,
			ConflictsWith: codegen_configs.GetResourceByName("ViewPolicy").GetConflictingFields("s3_visibility"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0) A list of usernames for bucket listing permissions`,

			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},

		"s3_visibility_groups": &schema.Schema{
			Type:          schema.TypeList,
			ConflictsWith: codegen_configs.GetResourceByName("ViewPolicy").GetConflictingFields("s3_visibility_groups"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0) A list of group names for bucket listing permissions`,

			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},

		"apple_sid": &schema.Schema{
			Type:          schema.TypeBool,
			ConflictsWith: codegen_configs.GetResourceByName("ViewPolicy").GetConflictingFields("apple_sid"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0) Apple sid`,
		},

		"protocols": &schema.Schema{
			Type:          schema.TypeList,
			ConflictsWith: codegen_configs.GetResourceByName("ViewPolicy").GetConflictingFields("protocols"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0) Protocols to audit`,

			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},

		"data_create_delete": &schema.Schema{
			Type:          schema.TypeBool,
			ConflictsWith: codegen_configs.GetResourceByName("ViewPolicy").GetConflictingFields("data_create_delete"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0) Create/Delete Files/Directories/Objects`,
		},

		"data_modify": &schema.Schema{
			Type:          schema.TypeBool,
			ConflictsWith: codegen_configs.GetResourceByName("ViewPolicy").GetConflictingFields("data_modify"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0) Modify data/MD`,
		},

		"data_read": &schema.Schema{
			Type:          schema.TypeBool,
			ConflictsWith: codegen_configs.GetResourceByName("ViewPolicy").GetConflictingFields("data_read"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0) Read data`,
		},

		"log_full_path": &schema.Schema{
			Type:          schema.TypeBool,
			ConflictsWith: codegen_configs.GetResourceByName("ViewPolicy").GetConflictingFields("log_full_path"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0) Log full path`,
		},

		"log_hostname": &schema.Schema{
			Type:          schema.TypeBool,
			ConflictsWith: codegen_configs.GetResourceByName("ViewPolicy").GetConflictingFields("log_hostname"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0) Log hostname`,
		},

		"log_username": &schema.Schema{
			Type:          schema.TypeBool,
			ConflictsWith: codegen_configs.GetResourceByName("ViewPolicy").GetConflictingFields("log_username"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0) Log username`,
		},

		"log_deleted": &schema.Schema{
			Type:          schema.TypeBool,
			ConflictsWith: codegen_configs.GetResourceByName("ViewPolicy").GetConflictingFields("log_deleted"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0) Log deleted files/dirs from trash dir`,
		},

		"count_views": &schema.Schema{
			Type:          schema.TypeInt,
			ConflictsWith: codegen_configs.GetResourceByName("ViewPolicy").GetConflictingFields("count_views"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0) Number of Policy related Views`,
		},

		"enable_snapshot_lookup": &schema.Schema{
			Type:          schema.TypeBool,
			ConflictsWith: codegen_configs.GetResourceByName("ViewPolicy").GetConflictingFields("enable_snapshot_lookup"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0) Specifies whether to make the .snapshot directory accessible in subdirectories of the View.`,
		},

		"enable_listing_of_snapshot_dir": &schema.Schema{
			Type:          schema.TypeBool,
			ConflictsWith: codegen_configs.GetResourceByName("ViewPolicy").GetConflictingFields("enable_listing_of_snapshot_dir"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0) Specifies whether to make the .snapshot directory visible in subdirectories of the View.`,
		},

		"s3_special_chars_support": &schema.Schema{
			Type:          schema.TypeBool,
			ConflictsWith: codegen_configs.GetResourceByName("ViewPolicy").GetConflictingFields("s3_special_chars_support"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.1.0) This will enable object names that contain “//“ or “/../“ and are incompatible with other protocols.`,
		},

		"smb_is_ca": &schema.Schema{
			Type:          schema.TypeBool,
			ConflictsWith: codegen_configs.GetResourceByName("ViewPolicy").GetConflictingFields("smb_is_ca"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.1.0) When enabled, the SMB share exposed by the view is set as continuously available, which allows SMB3 clients to request use of persistent file handles and keep their connections to this share in case of a failover event.`,
		},

		"nfs_case_insensitive": &schema.Schema{
			Type:          schema.TypeBool,
			ConflictsWith: codegen_configs.GetResourceByName("ViewPolicy").GetConflictingFields("nfs_case_insensitive"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.1.0) Force case insensitivity for NFSv3 and NFSv4`,
		},

		"enable_access_to_snapshot_dir_in_subdirs": &schema.Schema{
			Type:          schema.TypeBool,
			ConflictsWith: codegen_configs.GetResourceByName("ViewPolicy").GetConflictingFields("enable_access_to_snapshot_dir_in_subdirs"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.1.0) Specifies whether to make the .snapshot directory visible in subdirectories of the View.`,
		},

		"enable_visibility_of_snapshot_dir": &schema.Schema{
			Type:          schema.TypeBool,
			ConflictsWith: codegen_configs.GetResourceByName("ViewPolicy").GetConflictingFields("enable_visibility_of_snapshot_dir"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.1.0) Specifies whether to make the .snapshot directory visible in subdirectories of the View.`,
		},

		"nfs_enforce_tls": &schema.Schema{
			Type:          schema.TypeBool,
			ConflictsWith: codegen_configs.GetResourceByName("ViewPolicy").GetConflictingFields("nfs_enforce_tls"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.1.0) Accept NFSv3 and NFSv4.1 client mounts only if they are TLS-encrypted. Use only with Minimal Protection Level set to System or None.`,
		},

		"protocols_audit": &schema.Schema{
			Type:          schema.TypeList,
			ConflictsWith: codegen_configs.GetResourceByName("ViewPolicy").GetConflictingFields("protocols_audit"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0) `,

			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{

					"create_delete_files_dirs_objects": &schema.Schema{
						Type:          schema.TypeBool,
						ConflictsWith: codegen_configs.GetResourceByName("ProtocolsAudit").GetConflictingFields("create_delete_files_dirs_objects"),

						Computed:    false,
						Optional:    true,
						Sensitive:   false,
						Description: `(Valid for versions: 5.1.0) Audit operations that create or delete files, directories, or objects.`,

						Default: false,
					},

					"log_deleted_files_dirs": &schema.Schema{
						Type:          schema.TypeBool,
						ConflictsWith: codegen_configs.GetResourceByName("ProtocolsAudit").GetConflictingFields("log_deleted_files_dirs"),

						Computed:    false,
						Optional:    true,
						Sensitive:   false,
						Description: `(Valid for versions: 5.1.0) Log deleted files and directories.`,

						Default: false,
					},

					"log_full_path": &schema.Schema{
						Type:          schema.TypeBool,
						ConflictsWith: codegen_configs.GetResourceByName("ProtocolsAudit").GetConflictingFields("log_full_path"),

						Computed:    false,
						Optional:    true,
						Sensitive:   false,
						Description: `(Valid for versions: 5.1.0) Log full Element Store path to the requested resource. Enabled by default. May affect performance. When disabled, the view path is recorded.`,

						Default: true,
					},

					"log_username": &schema.Schema{
						Type:          schema.TypeBool,
						ConflictsWith: codegen_configs.GetResourceByName("ProtocolsAudit").GetConflictingFields("log_username"),

						Computed:    false,
						Optional:    true,
						Sensitive:   false,
						Description: `(Valid for versions: 5.1.0) Log username of requesting user. Disabled by default`,

						Default: false,
					},

					"log_hostname": &schema.Schema{
						Type:          schema.TypeBool,
						ConflictsWith: codegen_configs.GetResourceByName("ProtocolsAudit").GetConflictingFields("log_hostname"),

						Computed:    true,
						Optional:    true,
						Sensitive:   false,
						Description: `(Valid for versions: 5.1.0) Log the accessing Hostname`,
					},

					"modify_data_md": &schema.Schema{
						Type:          schema.TypeBool,
						ConflictsWith: codegen_configs.GetResourceByName("ProtocolsAudit").GetConflictingFields("modify_data_md"),

						Computed:    false,
						Optional:    true,
						Sensitive:   false,
						Description: `(Valid for versions: 5.1.0) Audit operations that modify data (including operations that change the file size) and metadata`,

						Default: false,
					},

					"read_data": &schema.Schema{
						Type:          schema.TypeBool,
						ConflictsWith: codegen_configs.GetResourceByName("ProtocolsAudit").GetConflictingFields("read_data"),

						Computed:    false,
						Optional:    true,
						Sensitive:   false,
						Description: `(Valid for versions: 5.1.0) Audit operations that read data and metadata`,

						Default: false,
					},

					"modify_data": &schema.Schema{
						Type:          schema.TypeBool,
						ConflictsWith: codegen_configs.GetResourceByName("ProtocolsAudit").GetConflictingFields("modify_data"),

						Computed:    false,
						Optional:    true,
						Sensitive:   false,
						Description: `(Valid for versions: 5.1.0) `,

						Default: false,
					},

					"read_data_md": &schema.Schema{
						Type:          schema.TypeBool,
						ConflictsWith: codegen_configs.GetResourceByName("ProtocolsAudit").GetConflictingFields("read_data_md"),

						Computed:    false,
						Optional:    true,
						Sensitive:   false,
						Description: `(Valid for versions: 5.1.0) `,

						Default: false,
					},
				},
			},
		},

		"vippool_permissions": &schema.Schema{
			Type:          schema.TypeList,
			ConflictsWith: codegen_configs.GetResourceByName("ViewPolicy").GetConflictingFields("vippool_permissions"),

			DiffSuppressOnRefresh: false,
			DiffSuppressFunc:      codegen_configs.GetResourceByName("ViewPolicy").GetAttributeDiffFunc("vippool_permissions"),
			Computed:              true,
			Optional:              true,
			Sensitive:             false,
			Description:           `(Valid for versions: 5.1.0) List of VIP pool permissions this attribute conflicts with vip_pools and can not be provided togather`,

			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{

					"vippool_id": &schema.Schema{
						Type:          schema.TypeString,
						ConflictsWith: codegen_configs.GetResourceByName("PermissionsPerVipPool").GetConflictingFields("vippool_id"),

						Computed:    true,
						Optional:    true,
						Sensitive:   false,
						Description: `(Valid for versions: 5.1.0) The Vippool ID`,
					},

					"vippool_permissions": &schema.Schema{
						Type:          schema.TypeString,
						ConflictsWith: codegen_configs.GetResourceByName("PermissionsPerVipPool").GetConflictingFields("vippool_permissions"),

						Computed:  false,
						Optional:  true,
						Sensitive: false,

						ValidateDiagFunc: utils.OneOf([]string{"RW", "RO"}),
						Description:      `(Valid for versions: 5.1.0) VIP pool permissions  Allowed Values are [RW RO]`,

						Default: "RW",
					},
				},
			},
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

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "VipPools", resource.VipPools))

	err = d.Set("vip_pools", utils.FlattenListOfPrimitives(&resource.VipPools))

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"vip_pools\"",
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

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "S3SpecialCharsSupport", resource.S3SpecialCharsSupport))

	err = d.Set("s3_special_chars_support", resource.S3SpecialCharsSupport)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"s3_special_chars_support\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "SmbIsCa", resource.SmbIsCa))

	err = d.Set("smb_is_ca", resource.SmbIsCa)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"smb_is_ca\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "NfsCaseInsensitive", resource.NfsCaseInsensitive))

	err = d.Set("nfs_case_insensitive", resource.NfsCaseInsensitive)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"nfs_case_insensitive\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "EnableAccessToSnapshotDirInSubdirs", resource.EnableAccessToSnapshotDirInSubdirs))

	err = d.Set("enable_access_to_snapshot_dir_in_subdirs", resource.EnableAccessToSnapshotDirInSubdirs)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"enable_access_to_snapshot_dir_in_subdirs\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "EnableVisibilityOfSnapshotDir", resource.EnableVisibilityOfSnapshotDir))

	err = d.Set("enable_visibility_of_snapshot_dir", resource.EnableVisibilityOfSnapshotDir)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"enable_visibility_of_snapshot_dir\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "NfsEnforceTls", resource.NfsEnforceTls))

	err = d.Set("nfs_enforce_tls", resource.NfsEnforceTls)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"nfs_enforce_tls\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "ProtocolsAudit", resource.ProtocolsAudit))

	tflog.Debug(ctx, fmt.Sprintf("Found a pointer object %v", resource.ProtocolsAudit))
	err = d.Set("protocols_audit", utils.FlattenModelAsList(ctx, resource.ProtocolsAudit))

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"protocols_audit\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "VippoolPermissions", resource.VippoolPermissions))

	err = d.Set("vippool_permissions", utils.FlattenListOfModelsToList(ctx, resource.VippoolPermissions))

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"vippool_permissions\"",
			Detail:   err.Error(),
		})
	}

	return diags

}
func resourceViewPolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(vast_client.JwtSession)
	resource_config := codegen_configs.GetResourceByName("ViewPolicy")
	attrs := map[string]interface{}{"path": utils.GenPath("viewpolicies"), "id": d.Id()}
	response, err := resource_config.GetFunc(ctx, client, attrs, d, map[string]string{})
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
	body, err := resource_config.ResponseProcessingFunc(ctx, response)

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
	resource_config := codegen_configs.GetResourceByName("ViewPolicy")
	attrs := map[string]interface{}{"path": utils.GenPath("viewpolicies"), "id": d.Id()}

	response, err := resource_config.DeleteFunc(ctx, client, attrs, nil, map[string]string{})

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
	resource_config := codegen_configs.GetResourceByName("ViewPolicy")
	tflog.Info(ctx, fmt.Sprintf("Creating Resource ViewPolicy"))
	reflect_ViewPolicy := reflect.TypeOf((*api_latest.ViewPolicy)(nil))
	utils.PopulateResourceMap(new_ctx, reflect_ViewPolicy.Elem(), d, &data, "", false)

	var before_post_error error
	data, before_post_error = resource_config.BeforePostFunc(data, client, ctx, d)
	if before_post_error != nil {
		return diag.FromErr(before_post_error)
	}

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
	attrs := map[string]interface{}{"path": utils.GenPath("viewpolicies")}
	response, create_err := resource_config.CreateFunc(ctx, client, attrs, data, map[string]string{})
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

	id_err := resource_config.IdFunc(ctx, client, resource.Id, d)
	if id_err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to set Id",
			Detail:   err.Error(),
		})
		return diags
	}
	ctx_with_resource := context.WithValue(ctx, utils.ContextKey("resource"), resource)
	resourceViewPolicyRead(ctx_with_resource, d, m)

	return diags
}

func resourceViewPolicyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	names_mapping := utils.ContextKey("names_mapping")
	new_ctx := context.WithValue(ctx, names_mapping, ViewPolicy_names_mapping)
	var diags diag.Diagnostics
	data := make(map[string]interface{})
	version_compare := utils.VastVersionsWarn(ctx)
	resource_config := codegen_configs.GetResourceByName("ViewPolicy")
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
	tflog.Info(ctx, fmt.Sprintf("Updating Resource ViewPolicy"))
	reflect_ViewPolicy := reflect.TypeOf((*api_latest.ViewPolicy)(nil))
	utils.PopulateResourceMap(new_ctx, reflect_ViewPolicy.Elem(), d, &data, "", false)

	var before_patch_error error
	data, before_patch_error = resource_config.BeforePatchFunc(data, client, ctx, d)
	if before_patch_error != nil {
		return diag.FromErr(before_patch_error)
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
	attrs := map[string]interface{}{"path": utils.GenPath("viewpolicies"), "id": d.Id()}
	response, patch_err := resource_config.UpdateFunc(ctx, client, attrs, data, d, map[string]string{})
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
	resource_config := codegen_configs.GetResourceByName("ViewPolicy")
	attrs := map[string]interface{}{"path": utils.GenPath("viewpolicies")}
	response, err := resource_config.ImportFunc(ctx, client, attrs, d, resource_config.Importer.GetFunc())

	if err != nil {
		return result, err
	}

	resource_l := []api_latest.ViewPolicy{}
	body, err := resource_config.ResponseProcessingFunc(ctx, response)

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
	id_err := resource_config.IdFunc(ctx, client, resource.Id, d)
	if id_err != nil {
		return result, id_err
	}

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
