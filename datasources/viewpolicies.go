package datasources

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api_latest "github.com/vast-data/terraform-provider-vastdata/codegen/latest"
	codegen_configs "github.com/vast-data/terraform-provider-vastdata/codegen_tools/configs"
	utils "github.com/vast-data/terraform-provider-vastdata/utils"
	vast_client "github.com/vast-data/terraform-provider-vastdata/vast-client"
	"net/url"
)

func DataSourceViewPolicy() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceViewPolicyRead,
		Description: ``,
		Schema: map[string]*schema.Schema{

			"id": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) The unique ID of the view policy.`,
			},

			"guid": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) The unique GUID of the view policy.`,
			},

			"name": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    false,
				Required:    true,
				Optional:    false,
				Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) The unique name of the view policy.`,
			},

			"gid_inheritance": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Determines the way a file inherits GID.`,
			},

			"flavor": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Security flavor, which determines how file and directory permissions are applied in multiprotocol views.`,
			},

			"access_flavor": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Controls which protocols can set file permissions when the security flavor is set to 'MIXED_LAST_WINS'. Note that with 'MIXED_LAST_WINS', NFSv3 access is unaffected by this option.`,
			},

			"path_length": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) How to determine the maximum allowed path length. Allowed Values are [LCD NPL]`,
			},

			"allowed_characters": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) How to determine the allowed characters in a path.`,
			},

			"use32bit_fileid": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: ``,
			},

			"expose_id_in_fsid": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) `,
			},

			"use_auth_provider": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    false,
				Required:    false,
				Optional:    true,
				Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Enables or disables use of the configured authentication provider(s) to enforce group permissions. If set to 'true' without specifing 'auth_source', the 'auth_source' is set to 'PROVIDERS'. If set to 'false', the 'auth_source' is set to RPC. Due to the nature of Terraform, simply changing 'use_auth_provider' from 'false' to 'true' or the other way around will not change the value of the 'auth_source' parameter because Terrafrom will keep holding the previous value. Therefore it is advisable to always specify the value of 'auth_source'.`,
			},

			"auth_source": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) The source of authentication. Allowed Values are [PROVIDERS RPC RPC_AND_PROVIDERS]`,
			},

			"read_write": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Hosts with NFS read/write permissions.`,

				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"read_only": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Hosts with NFS read-only permissions.`,

				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"nfs_read_write": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Hosts with NFS read/write permissions. When creating a new view policy, if this value is not set, then an empty list is sent to the VAST cluster, resulting in an empty list of addresses. However, during an update, if 'nfs_read_write' is removed from the resource, nothing is changed to preserve Terraform default behavior in such cases. If you need to change the value to an empty list, specify this parameter and set it to [].`,

				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"nfs_read_only": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Hosts with NFS read-only permissions. When creating a new view policy, if this value is not set, then an empty list is sent to the VAST cluster, resulting in an empty list of addresses. However, during an update, if 'nfs_read_only' is removed from the resource, nothing is changed to preserve Terraform default behavior in such cases. If you need to change the value to an empty list, specify this parameter and set it to [].`,

				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"smb_read_write": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Hosts with SMB read/write permissions. When creating a new view policy, if this value is not set, then an empty list is sent to the VAST cluster, resulting in an empty list of addresses. However, during an update, if 'smb_read_write' is removed from the resource, nothing is changed to preserve Terraform default behavior in such cases. If you need to change the value to an empty list, specify this parameter and set it to [].`,

				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"smb_read_only": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Hosts with SMB read-only permissions. When creating a new view policy, if this value is not set, then an empty list is sent to the VAST cluster, resulting in an empty list of addresses. However, during an update, if 'smb_read_only' is removed from the resource, nothing is changed to preserve Terraform default behavior in such cases. If you need to change the value to an empty list, specify this parameter and set it to [].`,

				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"s3_read_write": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Hosts with S3 read/write permissions. When creating a new view policy, if this value is not set, then an empty list is sent to the VAST cluster, resulting in an empty list of addresses. However, during an update, if 's3_read_write' is removed from the resource, nothing is changed to preserve Terraform default behavior in such cases. If you need to change the value to an empty list, specify this parameter and set it to [].`,

				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"s3_read_only": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Hosts with S3 read-only permissions. When creating a new view policy, if this value is not set, then an empty list is sent to the VAST cluster, resulting in an empty list of addresses. However, during an update, if 's3_read_only' is removed from the resource, nothing is changed to preserve Terraform default behavior in such cases. If you need to change the value to an empty list, specify this parameter and set it to [].`,

				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"trash_access": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Hosts with trash folder permissions.`,

				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"nfs_posix_acl": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Enables or disables POSIX ACL.`,
			},

			"nfs_return_open_permissions": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) When using SMB, enables or disables use of open permissions for files.`,
			},

			"nfs_no_squash": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Hosts with the no-squash policy. When creating a new view policy, if this value is not set, then an empty list is sent to the VAST cluster, resulting in an empty list of addresses. However, during an update, if 'nfs_no_squash' is removed from the resource, nothing is changed to preserve Terraform default behavior in such cases. If you need to change the value to an empty list, specify this parameter and set it to [].`,

				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"nfs_root_squash": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Hosts with the root squash policy. When creating a new view policy, if this value is not set, then an empty list is sent to the VAST cluster, resulting in an empty list of addresses. However, during an update, if 'nfs_root_squash' is removed from the resource, nothing is changed to preserve Terraform default behavior in such cases. If you need to change the value to an empty list, specify this parameter and set it to [].`,

				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"nfs_all_squash": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Hosts with the all-squash policy. When creating a new view policy, if this value is not set, then an empty list is sent to the VAST cluster, resulting in an empty list of addresses. However, during an update, if 'nfs_all_squash' is removed from the resource, nothing is changed to preserve Terraform default behavior in such cases. If you need to change the value to an empty list, specify this parameter and set it to [].`,

				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"s3_bucket_full_control": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Hosts with full permissions`,
			},

			"s3_bucket_listing": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Hosts with full permissions`,
			},

			"s3_bucket_read": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Hosts with full permissions`,
			},

			"s3_bucket_read_acp": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Hosts with full permissions`,
			},

			"s3_bucket_write": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Hosts with full permissions`,
			},

			"s3_bucket_write_acp": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Hosts with full permissions`,
			},

			"s3_object_full_control": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Hosts with full permissions`,
			},

			"s3_object_read": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Hosts with full permissions`,
			},

			"s3_object_read_acp": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Hosts with full permissions`,
			},

			"s3_object_write": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Hosts with full permissions`,
			},

			"s3_object_write_acp": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Hosts with full permissions`,
			},

			"smb_file_mode": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Default UNIX-type permissions on a new file.`,
			},

			"smb_directory_mode": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Default UNIX-type permissions on a new folder.`,
			},

			"smb_file_mode_padded": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Default UNIX-type permissions on a new file.`,
			},

			"smb_directory_mode_padded": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Default UNIX-type permissions on a new folder.`,
			},

			"cluster": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Parent cluster.`,
			},

			"cluster_id": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Parent cluster ID.`,
			},

			"tenant_id": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Tenant ID.`,
			},

			"tenant_name": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Tenant name.`,
			},

			"url": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) `,
			},

			"atime_frequency": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Frequency for updating the atime attribute of NFS files. atime is updated on read operations if the difference between the current time and the file's atime value is greater than the atime frequency. Specify as time in seconds.`,
			},

			"vip_pools": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) A comma-separated list of virtual IP pool IDs. This attribute conflicts with 'vippool_permissions' and cannot be supplied together with it. Also, due to the lack of ability to configure virtual IP pool permissions using this attibute, the pool permissions are always defined as read/write.`,

				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},

			"nfs_minimal_protection_level": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) NFSv4.1 minimal protection level. Allowed Values are [NONE SYSTEM KRB_AUTH_ONLY KRB_INTEGRITY KRB_PRIVACY]`,
			},

			"s3_visibility": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) A list of usernames for bucket listing permissions.`,

				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"s3_visibility_groups": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) A list of group names for bucket listing permissions.`,

				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"apple_sid": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Enables or disables returning Security IDs (SIDs) in Apple compatible representation.`,
			},

			"protocols": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Protocols to audit.`,

				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"data_create_delete": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) If 'true', logs operations that create or delete files, directories, or objects.`,
			},

			"data_modify": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) If 'true', logs operations that modify data or metadata.`,
			},

			"data_read": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) If 'true', logs operations that read data or metadata.`,
			},

			"log_full_path": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Enables or disables logging of full Element Store path to the requested resource.`,
			},

			"log_hostname": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Enables or disables logging of the accessing host name.`,
			},

			"log_username": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Enables or disables logging of username of the requesting user. Disabled by default.`,
			},

			"log_deleted": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Enables or disables logging of deleted files and directories from the trash folder.`,
			},

			"count_views": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) The total number of views associated with this view policy.`,
			},

			"enable_snapshot_lookup": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Specifies whether to make the .snapshot directory accessible in subdirectories of the view.`,
			},

			"enable_listing_of_snapshot_dir": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Specifies whether to make the .snapshot directory visible in subdirectories of the view.`,
			},

			"s3_special_chars_support": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `(Valid for versions: 5.1.0,5.2.0) Allows or prohibits object names that contain “//“ or “/../“ and are incompatible with other protocols.`,
			},

			"smb_is_ca": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `(Valid for versions: 5.1.0,5.2.0) If 'true', the SMB share exposed by the view is set as continuously available, which allows SMBv3 clients to request use of persistent file handles and keep their connections to this share in case of a failover event.`,
			},

			"nfs_case_insensitive": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `(Valid for versions: 5.1.0,5.2.0) If 'true', forces case insensitivity for NFSv3 and NFSv4.`,
			},

			"enable_access_to_snapshot_dir_in_subdirs": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `(Valid for versions: 5.1.0,5.2.0) Specifies whether to make the .snapshot directory visible in subdirectories of the view.`,
			},

			"enable_visibility_of_snapshot_dir": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `(Valid for versions: 5.1.0,5.2.0) Specifies whether to make the .snapshot directory visible in subdirectories of the view.`,
			},

			"nfs_enforce_tls": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `(Valid for versions: 5.1.0,5.2.0) If 'true', NFSv3 and NFSv4.1 client mounts are accepted only if they are TLS-encrypted. Use only with 'nfs_minimal_protection_level' set to 'System' or 'None'.`,
			},

			"vippool_permissions": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `(Valid for versions: 5.1.0,5.2.0) A list of virtual IP pool permissions. This attribute conflicts with 'vip_pools' and cannot be supplied together with it.`,

				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						"vippool_id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Required:    false,
							Optional:    false,
							Description: `(Valid for versions: 5.1.0,5.2.0) The ID of the virtual IP pool.`,
						},

						"vippool_permissions": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    false,
							Required:    false,
							Optional:    true,
							Description: `(Valid for versions: 5.1.0,5.2.0) Permissions for the virtual IP pool. Allowed Values are [RW RO]`,
						},
					},
				},
			},
		},
	}
}

func dataSourceViewPolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(*vast_client.VMSSession)
	values := url.Values{}
	datasource_config := codegen_configs.GetDataSourceByName("ViewPolicy")

	name := d.Get("name")
	values.Add("name", fmt.Sprintf("%v", name))

	_path := fmt.Sprintf(
		"viewpolicies",
	)
	response, err := client.Get(ctx, utils.GenPath(_path), values.Encode(), map[string]string{})
	tflog.Info(ctx, response.Request.URL.String())
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred while obtaining data from the vastdata cluster",
			Detail:   err.Error(),
		})
		return diags

	}
	resource_l := []api_latest.ViewPolicy{}
	body, err := datasource_config.ResponseProcessingFunc(ctx, response, d)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred reading data received from VastData cluster",
			Detail:   err.Error(),
		})
		return diags

	}
	err = json.Unmarshal(body, &resource_l)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred while parsing data received from VastData cluster",
			Detail:   err.Error(),
		})
		return diags

	}
	if len(resource_l) == 0 {
		d.SetId("")
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Could not find a resource that matches those attributes",
			Detail:   "Could not find a resource that matches those attributes",
		})
		return diags
	}
	if len(resource_l) > 1 {
		d.SetId("")
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Multiple results returned, you might want to add more attributes to get a specific resource",
			Detail:   "Multiple results returned, you might want to add more attributes to get a specific resource",
		})
		return diags
	}

	resource := resource_l[0]

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "Id", resource.Id))

	err = d.Set("id", resource.Id)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"id\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "Guid", resource.Guid))

	err = d.Set("guid", resource.Guid)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"guid\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "Name", resource.Name))

	err = d.Set("name", resource.Name)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"name\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "GidInheritance", resource.GidInheritance))

	err = d.Set("gid_inheritance", resource.GidInheritance)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"gid_inheritance\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "Flavor", resource.Flavor))

	err = d.Set("flavor", resource.Flavor)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"flavor\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "AccessFlavor", resource.AccessFlavor))

	err = d.Set("access_flavor", resource.AccessFlavor)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"access_flavor\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "PathLength", resource.PathLength))

	err = d.Set("path_length", resource.PathLength)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"path_length\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "AllowedCharacters", resource.AllowedCharacters))

	err = d.Set("allowed_characters", resource.AllowedCharacters)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"allowed_characters\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "Use32bitFileid", resource.Use32bitFileid))

	err = d.Set("use32bit_fileid", resource.Use32bitFileid)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"use32bit_fileid\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "ExposeIdInFsid", resource.ExposeIdInFsid))

	err = d.Set("expose_id_in_fsid", resource.ExposeIdInFsid)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"expose_id_in_fsid\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "UseAuthProvider", resource.UseAuthProvider))

	err = d.Set("use_auth_provider", resource.UseAuthProvider)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"use_auth_provider\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "AuthSource", resource.AuthSource))

	err = d.Set("auth_source", resource.AuthSource)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"auth_source\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "ReadWrite", resource.ReadWrite))

	err = d.Set("read_write", utils.FlattenListOfPrimitives(&resource.ReadWrite))

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"read_write\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "ReadOnly", resource.ReadOnly))

	err = d.Set("read_only", utils.FlattenListOfPrimitives(&resource.ReadOnly))

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"read_only\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "NfsReadWrite", resource.NfsReadWrite))

	err = d.Set("nfs_read_write", utils.FlattenListOfPrimitives(&resource.NfsReadWrite))

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"nfs_read_write\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "NfsReadOnly", resource.NfsReadOnly))

	err = d.Set("nfs_read_only", utils.FlattenListOfPrimitives(&resource.NfsReadOnly))

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"nfs_read_only\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "SmbReadWrite", resource.SmbReadWrite))

	err = d.Set("smb_read_write", utils.FlattenListOfPrimitives(&resource.SmbReadWrite))

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"smb_read_write\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "SmbReadOnly", resource.SmbReadOnly))

	err = d.Set("smb_read_only", utils.FlattenListOfPrimitives(&resource.SmbReadOnly))

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"smb_read_only\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "S3ReadWrite", resource.S3ReadWrite))

	err = d.Set("s3_read_write", utils.FlattenListOfPrimitives(&resource.S3ReadWrite))

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"s3_read_write\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "S3ReadOnly", resource.S3ReadOnly))

	err = d.Set("s3_read_only", utils.FlattenListOfPrimitives(&resource.S3ReadOnly))

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"s3_read_only\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "TrashAccess", resource.TrashAccess))

	err = d.Set("trash_access", utils.FlattenListOfPrimitives(&resource.TrashAccess))

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"trash_access\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "NfsPosixAcl", resource.NfsPosixAcl))

	err = d.Set("nfs_posix_acl", resource.NfsPosixAcl)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"nfs_posix_acl\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "NfsReturnOpenPermissions", resource.NfsReturnOpenPermissions))

	err = d.Set("nfs_return_open_permissions", resource.NfsReturnOpenPermissions)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"nfs_return_open_permissions\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "NfsNoSquash", resource.NfsNoSquash))

	err = d.Set("nfs_no_squash", utils.FlattenListOfPrimitives(&resource.NfsNoSquash))

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"nfs_no_squash\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "NfsRootSquash", resource.NfsRootSquash))

	err = d.Set("nfs_root_squash", utils.FlattenListOfPrimitives(&resource.NfsRootSquash))

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"nfs_root_squash\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "NfsAllSquash", resource.NfsAllSquash))

	err = d.Set("nfs_all_squash", utils.FlattenListOfPrimitives(&resource.NfsAllSquash))

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"nfs_all_squash\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "S3BucketFullControl", resource.S3BucketFullControl))

	err = d.Set("s3_bucket_full_control", resource.S3BucketFullControl)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"s3_bucket_full_control\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "S3BucketListing", resource.S3BucketListing))

	err = d.Set("s3_bucket_listing", resource.S3BucketListing)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"s3_bucket_listing\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "S3BucketRead", resource.S3BucketRead))

	err = d.Set("s3_bucket_read", resource.S3BucketRead)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"s3_bucket_read\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "S3BucketReadAcp", resource.S3BucketReadAcp))

	err = d.Set("s3_bucket_read_acp", resource.S3BucketReadAcp)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"s3_bucket_read_acp\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "S3BucketWrite", resource.S3BucketWrite))

	err = d.Set("s3_bucket_write", resource.S3BucketWrite)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"s3_bucket_write\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "S3BucketWriteAcp", resource.S3BucketWriteAcp))

	err = d.Set("s3_bucket_write_acp", resource.S3BucketWriteAcp)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"s3_bucket_write_acp\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "S3ObjectFullControl", resource.S3ObjectFullControl))

	err = d.Set("s3_object_full_control", resource.S3ObjectFullControl)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"s3_object_full_control\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "S3ObjectRead", resource.S3ObjectRead))

	err = d.Set("s3_object_read", resource.S3ObjectRead)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"s3_object_read\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "S3ObjectReadAcp", resource.S3ObjectReadAcp))

	err = d.Set("s3_object_read_acp", resource.S3ObjectReadAcp)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"s3_object_read_acp\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "S3ObjectWrite", resource.S3ObjectWrite))

	err = d.Set("s3_object_write", resource.S3ObjectWrite)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"s3_object_write\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "S3ObjectWriteAcp", resource.S3ObjectWriteAcp))

	err = d.Set("s3_object_write_acp", resource.S3ObjectWriteAcp)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"s3_object_write_acp\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "SmbFileMode", resource.SmbFileMode))

	err = d.Set("smb_file_mode", resource.SmbFileMode)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"smb_file_mode\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "SmbDirectoryMode", resource.SmbDirectoryMode))

	err = d.Set("smb_directory_mode", resource.SmbDirectoryMode)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"smb_directory_mode\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "SmbFileModePadded", resource.SmbFileModePadded))

	err = d.Set("smb_file_mode_padded", resource.SmbFileModePadded)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"smb_file_mode_padded\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "SmbDirectoryModePadded", resource.SmbDirectoryModePadded))

	err = d.Set("smb_directory_mode_padded", resource.SmbDirectoryModePadded)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"smb_directory_mode_padded\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "Cluster", resource.Cluster))

	err = d.Set("cluster", resource.Cluster)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"cluster\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "ClusterId", resource.ClusterId))

	err = d.Set("cluster_id", resource.ClusterId)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"cluster_id\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "TenantId", resource.TenantId))

	err = d.Set("tenant_id", resource.TenantId)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"tenant_id\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "TenantName", resource.TenantName))

	err = d.Set("tenant_name", resource.TenantName)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"tenant_name\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "Url", resource.Url))

	err = d.Set("url", resource.Url)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"url\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "AtimeFrequency", resource.AtimeFrequency))

	err = d.Set("atime_frequency", resource.AtimeFrequency)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"atime_frequency\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "VipPools", resource.VipPools))

	err = d.Set("vip_pools", utils.FlattenListOfPrimitives(&resource.VipPools))

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"vip_pools\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "NfsMinimalProtectionLevel", resource.NfsMinimalProtectionLevel))

	err = d.Set("nfs_minimal_protection_level", resource.NfsMinimalProtectionLevel)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"nfs_minimal_protection_level\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "S3Visibility", resource.S3Visibility))

	err = d.Set("s3_visibility", utils.FlattenListOfPrimitives(&resource.S3Visibility))

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"s3_visibility\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "S3VisibilityGroups", resource.S3VisibilityGroups))

	err = d.Set("s3_visibility_groups", utils.FlattenListOfPrimitives(&resource.S3VisibilityGroups))

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"s3_visibility_groups\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "AppleSid", resource.AppleSid))

	err = d.Set("apple_sid", resource.AppleSid)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"apple_sid\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "Protocols", resource.Protocols))

	err = d.Set("protocols", utils.FlattenListOfPrimitives(&resource.Protocols))

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"protocols\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "DataCreateDelete", resource.DataCreateDelete))

	err = d.Set("data_create_delete", resource.DataCreateDelete)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"data_create_delete\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "DataModify", resource.DataModify))

	err = d.Set("data_modify", resource.DataModify)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"data_modify\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "DataRead", resource.DataRead))

	err = d.Set("data_read", resource.DataRead)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"data_read\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "LogFullPath", resource.LogFullPath))

	err = d.Set("log_full_path", resource.LogFullPath)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"log_full_path\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "LogHostname", resource.LogHostname))

	err = d.Set("log_hostname", resource.LogHostname)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"log_hostname\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "LogUsername", resource.LogUsername))

	err = d.Set("log_username", resource.LogUsername)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"log_username\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "LogDeleted", resource.LogDeleted))

	err = d.Set("log_deleted", resource.LogDeleted)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"log_deleted\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "CountViews", resource.CountViews))

	err = d.Set("count_views", resource.CountViews)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"count_views\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "EnableSnapshotLookup", resource.EnableSnapshotLookup))

	err = d.Set("enable_snapshot_lookup", resource.EnableSnapshotLookup)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"enable_snapshot_lookup\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "EnableListingOfSnapshotDir", resource.EnableListingOfSnapshotDir))

	err = d.Set("enable_listing_of_snapshot_dir", resource.EnableListingOfSnapshotDir)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"enable_listing_of_snapshot_dir\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "S3SpecialCharsSupport", resource.S3SpecialCharsSupport))

	err = d.Set("s3_special_chars_support", resource.S3SpecialCharsSupport)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"s3_special_chars_support\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "SmbIsCa", resource.SmbIsCa))

	err = d.Set("smb_is_ca", resource.SmbIsCa)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"smb_is_ca\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "NfsCaseInsensitive", resource.NfsCaseInsensitive))

	err = d.Set("nfs_case_insensitive", resource.NfsCaseInsensitive)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"nfs_case_insensitive\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "EnableAccessToSnapshotDirInSubdirs", resource.EnableAccessToSnapshotDirInSubdirs))

	err = d.Set("enable_access_to_snapshot_dir_in_subdirs", resource.EnableAccessToSnapshotDirInSubdirs)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"enable_access_to_snapshot_dir_in_subdirs\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "EnableVisibilityOfSnapshotDir", resource.EnableVisibilityOfSnapshotDir))

	err = d.Set("enable_visibility_of_snapshot_dir", resource.EnableVisibilityOfSnapshotDir)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"enable_visibility_of_snapshot_dir\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "NfsEnforceTls", resource.NfsEnforceTls))

	err = d.Set("nfs_enforce_tls", resource.NfsEnforceTls)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"nfs_enforce_tls\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "VippoolPermissions", resource.VippoolPermissions))

	err = d.Set("vippool_permissions", utils.FlattenListOfModelsToList(ctx, resource.VippoolPermissions))

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occurred setting value to \"vippool_permissions\"",
			Detail:   err.Error(),
		})
	}

	err = datasource_config.IdFunc(ctx, client, resource.Id, d)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to set Id",
			Detail:   err.Error(),
		})
		return diags
	}
	return diags
}
