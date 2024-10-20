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

func ResourceView() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceViewRead,
		DeleteContext: resourceViewDelete,
		CreateContext: resourceViewCreate,
		UpdateContext: resourceViewUpdate,

		Importer: &schema.ResourceImporter{
			StateContext: resourceViewImporter,
		},

		Description: ``,
		Schema:      getResourceViewSchema(),
	}
}

func getResourceViewSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{

		"guid": &schema.Schema{
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("View").GetConflictingFields("guid"),

			Computed:    true,
			Optional:    false,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) A uniqe GUID assigned to the View`,
		},

		"name": &schema.Schema{
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("View").GetConflictingFields("name"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) A uniq name given to the view`,
		},

		"path": &schema.Schema{
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("View").GetConflictingFields("path"),

			Required:    true,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) File system path. Begin with '/'. Do not include a trailing slash`,
		},

		"create_dir": &schema.Schema{
			Type:          schema.TypeBool,
			ConflictsWith: codegen_configs.GetResourceByName("View").GetConflictingFields("create_dir"),

			DiffSuppressOnRefresh: false,
			DiffSuppressFunc:      utils.DoNothingOnUpdate(),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Creates the directory specified by the path`,
		},

		"alias": &schema.Schema{
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("View").GetConflictingFields("alias"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Alias for NFS export, must start with '/' and only ASCII characters are allowed. If configured, this supersedes the exposed NFS export path`,
		},

		"bucket": &schema.Schema{
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("View").GetConflictingFields("bucket"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) S3 Bucket name`,
		},

		"policy_id": &schema.Schema{
			Type:          schema.TypeInt,
			ConflictsWith: codegen_configs.GetResourceByName("View").GetConflictingFields("policy_id"),

			Required:    true,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Associated view policy ID`,
		},

		"cluster": &schema.Schema{
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("View").GetConflictingFields("cluster"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Parent Cluster`,
		},

		"cluster_id": &schema.Schema{
			Type:          schema.TypeInt,
			ConflictsWith: codegen_configs.GetResourceByName("View").GetConflictingFields("cluster_id"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Parent Cluster ID`,
		},

		"tenant_id": &schema.Schema{
			Type:          schema.TypeInt,
			ConflictsWith: codegen_configs.GetResourceByName("View").GetConflictingFields("tenant_id"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) The tenant ID related to this view`,
		},

		"directory": &schema.Schema{
			Type:          schema.TypeBool,
			ConflictsWith: codegen_configs.GetResourceByName("View").GetConflictingFields("directory"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Create the directory if it does not exist`,
		},

		"s3_versioning": &schema.Schema{
			Type:          schema.TypeBool,
			ConflictsWith: codegen_configs.GetResourceByName("View").GetConflictingFields("s3_versioning"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Trun on S3 Versioning`,
		},

		"s3_unverified_lookup": &schema.Schema{
			Type:          schema.TypeBool,
			ConflictsWith: codegen_configs.GetResourceByName("View").GetConflictingFields("s3_unverified_lookup"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Allow S3 Unverified Lookup`,
		},

		"allow_anonymous_access": &schema.Schema{
			Type:          schema.TypeBool,
			ConflictsWith: codegen_configs.GetResourceByName("View").GetConflictingFields("allow_anonymous_access"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Allow S3 anonymous access`,
		},

		"allow_s3_anonymous_access": &schema.Schema{
			Type:          schema.TypeBool,
			ConflictsWith: codegen_configs.GetResourceByName("View").GetConflictingFields("allow_s3_anonymous_access"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Allow S3 anonymous access`,
		},

		"protocols": &schema.Schema{
			Type:          schema.TypeList,
			ConflictsWith: codegen_configs.GetResourceByName("View").GetConflictingFields("protocols"),

			DiffSuppressOnRefresh: false,
			DiffSuppressFunc:      codegen_configs.GetResourceByName("View").GetAttributeDiffFunc("protocols"),
			Computed:              true,
			Optional:              true,
			Sensitive:             false,
			Description:           `(Valid for versions: 5.0.0,5.1.0,5.2.0) Protocols exposed by this view`,

			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},

		"share": &schema.Schema{
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("View").GetConflictingFields("share"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Name of the SMB Share. Must not include the following characters: " \ / [ ] : | < > + = ; , * ?`,
		},

		"bucket_owner": &schema.Schema{
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("View").GetConflictingFields("bucket_owner"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) S3 Bucket owner`,
		},

		"bucket_creators": &schema.Schema{
			Type:          schema.TypeList,
			ConflictsWith: codegen_configs.GetResourceByName("View").GetConflictingFields("bucket_creators"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) List of bucket creators users`,

			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},

		"bucket_creators_groups": &schema.Schema{
			Type:          schema.TypeList,
			ConflictsWith: codegen_configs.GetResourceByName("View").GetConflictingFields("bucket_creators_groups"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) List of bucket creators groups`,

			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},

		"s3_locks": &schema.Schema{
			Type:          schema.TypeBool,
			ConflictsWith: codegen_configs.GetResourceByName("View").GetConflictingFields("s3_locks"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) S3 Object Lock`,
		},

		"s3_locks_retention_mode": &schema.Schema{
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("View").GetConflictingFields("s3_locks_retention_mode"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) S3 Locks retention mode`,
		},

		"s3_locks_retention_period": &schema.Schema{
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("View").GetConflictingFields("s3_locks_retention_period"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Period should be positive in format like 0d|2d|1y|2y`,
		},

		"physical_capacity": &schema.Schema{
			Type:          schema.TypeInt,
			ConflictsWith: codegen_configs.GetResourceByName("View").GetConflictingFields("physical_capacity"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Physical Capacity`,
		},

		"logical_capacity": &schema.Schema{
			Type:          schema.TypeInt,
			ConflictsWith: codegen_configs.GetResourceByName("View").GetConflictingFields("logical_capacity"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Logical Capacity`,
		},

		"nfs_interop_flags": &schema.Schema{
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("View").GetConflictingFields("nfs_interop_flags"),

			Computed:  true,
			Optional:  true,
			Sensitive: false,

			ValidateDiagFunc: utils.OneOf([]string{"BOTH_NFS3_AND_NFS4_INTEROP_DISABLED", "ONLY_NFS3_INTEROP_ENABLED", "ONLY_NFS4_INTEROP_ENABLED", "BOTH_NFS3_AND_NFS4_INTEROP_ENABLED"}),
			Description:      `(Valid for versions: 5.0.0,5.1.0,5.2.0) Indicates whether the view should support simultaneous access to NFS3/NFS4/SMB protocols. Allowed Values are [BOTH_NFS3_AND_NFS4_INTEROP_DISABLED ONLY_NFS3_INTEROP_ENABLED ONLY_NFS4_INTEROP_ENABLED BOTH_NFS3_AND_NFS4_INTEROP_ENABLED]`,
		},

		"is_remote": &schema.Schema{
			Type:          schema.TypeBool,
			ConflictsWith: codegen_configs.GetResourceByName("View").GetConflictingFields("is_remote"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) `,
		},

		"share_acl": &schema.Schema{
			Type:          schema.TypeList,
			ConflictsWith: codegen_configs.GetResourceByName("View").GetConflictingFields("share_acl"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) Share-level ACL details`,

			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{

					"enabled": &schema.Schema{
						Type:          schema.TypeBool,
						ConflictsWith: codegen_configs.GetResourceByName("ViewShareAcl").GetConflictingFields("enabled"),

						Computed:    true,
						Optional:    true,
						Sensitive:   false,
						Description: ``,
					},

					"acl": &schema.Schema{
						Type:          schema.TypeList,
						ConflictsWith: codegen_configs.GetResourceByName("ViewShareAcl").GetConflictingFields("acl"),

						Computed:    true,
						Optional:    true,
						Sensitive:   false,
						Description: ``,

						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{

								"grantee": &schema.Schema{
									Type:          schema.TypeString,
									ConflictsWith: codegen_configs.GetResourceByName("ShareAcl").GetConflictingFields("grantee"),

									Computed:  true,
									Optional:  true,
									Sensitive: false,

									ValidateDiagFunc: utils.OneOf([]string{"users", "groups"}),
									Description:      `(Valid for versions: 5.0.0,5.1.0,5.2.0)  Allowed Values are [users groups]`,
								},

								"permissions": &schema.Schema{
									Type:          schema.TypeString,
									ConflictsWith: codegen_configs.GetResourceByName("ShareAcl").GetConflictingFields("permissions"),

									Computed:  true,
									Optional:  true,
									Sensitive: false,

									ValidateDiagFunc: utils.OneOf([]string{"FULL", "CHANGE", "READ"}),
									Description:      `(Valid for versions: 5.0.0,5.1.0,5.2.0)  Allowed Values are [FULL CHANGE READ]`,
								},

								"sid_str": &schema.Schema{
									Type:          schema.TypeString,
									ConflictsWith: codegen_configs.GetResourceByName("ShareAcl").GetConflictingFields("sid_str"),

									Computed:    true,
									Optional:    true,
									Sensitive:   false,
									Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) `,
								},

								"uid_or_gid": &schema.Schema{
									Type:          schema.TypeInt,
									ConflictsWith: codegen_configs.GetResourceByName("ShareAcl").GetConflictingFields("uid_or_gid"),

									Computed:    true,
									Optional:    true,
									Sensitive:   false,
									Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) `,
								},

								"name": &schema.Schema{
									Type:          schema.TypeString,
									ConflictsWith: codegen_configs.GetResourceByName("ShareAcl").GetConflictingFields("name"),

									Required:    true,
									Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) `,
								},

								"fqdn": &schema.Schema{
									Type:          schema.TypeString,
									ConflictsWith: codegen_configs.GetResourceByName("ShareAcl").GetConflictingFields("fqdn"),

									Computed:    true,
									Optional:    true,
									Sensitive:   false,
									Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) `,
								},
							},
						},
					},
				},
			},
		},

		"qos_policy_id": &schema.Schema{
			Type:          schema.TypeInt,
			ConflictsWith: codegen_configs.GetResourceByName("View").GetConflictingFields("qos_policy_id"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.0.0,5.1.0,5.2.0) QoS Policy ID`,
		},

		"is_seamless": &schema.Schema{
			Type:          schema.TypeBool,
			ConflictsWith: codegen_configs.GetResourceByName("View").GetConflictingFields("is_seamless"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.1.0,5.2.0) Supports seamless failover between replication peers by syncing file handles between the view and remote views on the replicated path on replication peers. This enables NFSv3 client users to retain the same mount point to the view in the event of a failover of the view path to a replication peer. This feature enables NFSv3 client users to retain the same mount point to the view in the event of a failover of the view path to a replication peer. Enabling this option may cause overhead and should only be enabled when the use case is relevant. To complete the configuration for seamless failover between any two peers, a seamless view must be created on each peer.`,
		},

		"max_retention_period": &schema.Schema{
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("View").GetConflictingFields("max_retention_period"),

			Computed:  true,
			Optional:  true,
			Sensitive: false,

			ValidateDiagFunc: utils.ValidateRetention,
		},

		"min_retention_period": &schema.Schema{
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("View").GetConflictingFields("min_retention_period"),

			Computed:  true,
			Optional:  true,
			Sensitive: false,

			ValidateDiagFunc: utils.ValidateRetention,
		},

		"files_retention_mode": &schema.Schema{
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("View").GetConflictingFields("files_retention_mode"),

			Computed:  true,
			Optional:  true,
			Sensitive: false,

			ValidateDiagFunc: utils.OneOf([]string{"GOVERNANCE", "COMPLIANCE", "NONE"}),
			Description:      `(Valid for versions: 5.1.0,5.2.0) Applicable if locking is enabled. The retention mode for new files. For views enabled for NFSv3 or SMB, if locking is enabled, files_retention_mode must be set to GOVERNANCE or COMPLIANCE. If the view is enabled for S3 and not for NFSv3 or SMB, files_retention_mode can be set to NONE. If GOVERNANCE, locked files cannot be deleted or changed. The Retention settings can be shortened or extended by users with sufficient permissions. If COMPLIANCE, locked files cannot be deleted or changed. Retention settings can be extended, but not shortened, by users with sufficient permissions. If NONE (S3 only), the retention mode is not set for the view; it is set individually for each object. Allowed Values are [GOVERNANCE COMPLIANCE NONE]`,
		},

		"default_retention_period": &schema.Schema{
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("View").GetConflictingFields("default_retention_period"),

			Computed:  true,
			Optional:  true,
			Sensitive: false,

			ValidateDiagFunc: utils.ValidateRetention,
		},

		"auto_commit": &schema.Schema{
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("View").GetConflictingFields("auto_commit"),

			Computed:  true,
			Optional:  true,
			Sensitive: false,

			ValidateDiagFunc: utils.ValidateRetention,
		},

		"s3_object_ownership_rule": &schema.Schema{
			Type:          schema.TypeString,
			ConflictsWith: codegen_configs.GetResourceByName("View").GetConflictingFields("s3_object_ownership_rule"),

			Computed:  true,
			Optional:  true,
			Sensitive: false,

			ValidateDiagFunc: utils.OneOf([]string{"None", "BucketOwnerPreferred", "ObjectWriter", "BucketOwnerEnforced"}),
			Description:      `(Valid for versions: 5.1.0,5.2.0) S3 Object Ownership lets you set ownership of objects uploaded to a given bucket and to determine whether ACLs are used to control access to objects within this bucket. A bucket can be configured with one of the following object ownership rules: BucketOwnerEnforced - The bucket owner has full control over any object in the bucket ObjectWriter - The user that uploads an object has full control over this object. ACLs can be used to let other users access the object. BucketOwnerPreferred - The bucket owner has full control over new objects uploaded to the bucket by other users. ACLs can be used to control access to the objects. None - S3 Object Ownership is disabled for the bucket.  Allowed Values are [None BucketOwnerPreferred ObjectWriter BucketOwnerEnforced]`,
		},

		"locking": &schema.Schema{
			Type:          schema.TypeBool,
			ConflictsWith: codegen_configs.GetResourceByName("View").GetConflictingFields("locking"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.1.0,5.2.0) Write Once Read Many (WORM) locking enabled`,
		},

		"ignore_oos": &schema.Schema{
			Type:          schema.TypeBool,
			ConflictsWith: codegen_configs.GetResourceByName("View").GetConflictingFields("ignore_oos"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.1.0,5.2.0) Ignore oos`,
		},

		"bucket_logging": &schema.Schema{
			Type:          schema.TypeList,
			ConflictsWith: codegen_configs.GetResourceByName("View").GetConflictingFields("bucket_logging"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.2.0) `,

			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{

					"destination_id": &schema.Schema{
						Type:          schema.TypeInt,
						ConflictsWith: codegen_configs.GetResourceByName("BucketLogging").GetConflictingFields("destination_id"),

						Computed:    true,
						Optional:    true,
						Sensitive:   false,
						Description: `(Valid for versions: 5.2.0) The Logging bucket ID`,
					},

					"prefix": &schema.Schema{
						Type:          schema.TypeString,
						ConflictsWith: codegen_configs.GetResourceByName("BucketLogging").GetConflictingFields("prefix"),

						Computed:    false,
						Optional:    true,
						Sensitive:   false,
						Description: `(Valid for versions: 5.2.0) Log line prefix to add`,

						Default: "",
					},

					"key_format": &schema.Schema{
						Type:          schema.TypeString,
						ConflictsWith: codegen_configs.GetResourceByName("BucketLogging").GetConflictingFields("key_format"),

						Computed:  false,
						Optional:  true,
						Sensitive: false,

						ValidateDiagFunc: utils.OneOf([]string{"SIMPLE_PREFIX", "PARTITIONED_PREFIX_EVENT_TIME", "PARTITIONED_PREFIX_DELIVERY_TIME"}),
						Description:      `(Valid for versions: 5.2.0) The format for log object keys. SIMPLE_PREFIX=[DestinationPrefix][YYYY]-[MM]-[DD]-[hh]-[mm]-[ss]-[UniqueString], PARTITIONED_PREFIX_EVENT_TIME=[DestinationPrefix][SourceUsername]/[SourceBucket]/[YYYY]/[MM]/[DD]/[YYYY]-[MM]-[DD]-[hh]-[mm]-[ss]-[UniqueString] where the partitioning is done based on the time when the logged events occurred, PARTITIONED_PREFIX_DELIVERY_TIME=[DestinationPrefix][SourceUsername]/[SourceBucket]/[YYYY]/[MM]/[DD]/[YYYY]-[MM]-[DD]-[hh]-[mm]-[ss]-[UniqueString] where the partitioning is done based on the time when the log object has been delivered to the destination bucket. Default: SIMPLE_PREFIX Allowed Values are [SIMPLE_PREFIX PARTITIONED_PREFIX_EVENT_TIME PARTITIONED_PREFIX_DELIVERY_TIME]`,

						Default: "SIMPLE_PREFIX",
					},
				},
			},
		},

		"abac_tags": &schema.Schema{
			Type:          schema.TypeList,
			ConflictsWith: codegen_configs.GetResourceByName("View").GetConflictingFields("abac_tags"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.1.0,5.2.0) List of attribute based access control tags, this option can be used only when using SMB/NFSv4 protocols`,

			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},

		"abe_max_depth": &schema.Schema{
			Type:          schema.TypeInt,
			ConflictsWith: codegen_configs.GetResourceByName("View").GetConflictingFields("abe_max_depth"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.2.0) Restricts ABE to a specified path depth. For example, if max depth is 3, ABE does not affect paths deeper than three levels. If not specified, ABE affects all path depths.`,
		},

		"abe_protocols": &schema.Schema{
			Type:          schema.TypeList,
			ConflictsWith: codegen_configs.GetResourceByName("View").GetConflictingFields("abe_protocols"),

			Computed:    true,
			Optional:    true,
			Sensitive:   false,
			Description: `(Valid for versions: 5.2.0) The protocols for which Access-Based Enumeration (ABE) is enabled , allowed values [ NFS, SMB, NFS4, S3 ]`,

			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
	}
}

var View_names_mapping map[string][]string = map[string][]string{}

func ResourceViewReadStructIntoSchema(ctx context.Context, resource api_latest.View, d *schema.ResourceData) diag.Diagnostics {
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

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "Path", resource.Path))

	err = d.Set("path", resource.Path)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"path\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "CreateDir", resource.CreateDir))

	err = d.Set("create_dir", resource.CreateDir)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"create_dir\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "Alias", resource.Alias))

	err = d.Set("alias", resource.Alias)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"alias\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "Bucket", resource.Bucket))

	err = d.Set("bucket", resource.Bucket)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"bucket\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "PolicyId", resource.PolicyId))

	err = d.Set("policy_id", resource.PolicyId)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"policy_id\"",
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

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "Directory", resource.Directory))

	err = d.Set("directory", resource.Directory)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"directory\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "S3Versioning", resource.S3Versioning))

	err = d.Set("s3_versioning", resource.S3Versioning)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"s3_versioning\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "S3UnverifiedLookup", resource.S3UnverifiedLookup))

	err = d.Set("s3_unverified_lookup", resource.S3UnverifiedLookup)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"s3_unverified_lookup\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "AllowAnonymousAccess", resource.AllowAnonymousAccess))

	err = d.Set("allow_anonymous_access", resource.AllowAnonymousAccess)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"allow_anonymous_access\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "AllowS3AnonymousAccess", resource.AllowS3AnonymousAccess))

	err = d.Set("allow_s3_anonymous_access", resource.AllowS3AnonymousAccess)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"allow_s3_anonymous_access\"",
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

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "Share", resource.Share))

	err = d.Set("share", resource.Share)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"share\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "BucketOwner", resource.BucketOwner))

	err = d.Set("bucket_owner", resource.BucketOwner)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"bucket_owner\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "BucketCreators", resource.BucketCreators))

	err = d.Set("bucket_creators", utils.FlattenListOfPrimitives(&resource.BucketCreators))

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"bucket_creators\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "BucketCreatorsGroups", resource.BucketCreatorsGroups))

	err = d.Set("bucket_creators_groups", utils.FlattenListOfPrimitives(&resource.BucketCreatorsGroups))

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"bucket_creators_groups\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "S3Locks", resource.S3Locks))

	err = d.Set("s3_locks", resource.S3Locks)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"s3_locks\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "S3LocksRetentionMode", resource.S3LocksRetentionMode))

	err = d.Set("s3_locks_retention_mode", resource.S3LocksRetentionMode)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"s3_locks_retention_mode\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "S3LocksRetentionPeriod", resource.S3LocksRetentionPeriod))

	err = d.Set("s3_locks_retention_period", resource.S3LocksRetentionPeriod)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"s3_locks_retention_period\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "PhysicalCapacity", resource.PhysicalCapacity))

	err = d.Set("physical_capacity", resource.PhysicalCapacity)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"physical_capacity\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "LogicalCapacity", resource.LogicalCapacity))

	err = d.Set("logical_capacity", resource.LogicalCapacity)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"logical_capacity\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "NfsInteropFlags", resource.NfsInteropFlags))

	err = d.Set("nfs_interop_flags", resource.NfsInteropFlags)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"nfs_interop_flags\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "IsRemote", resource.IsRemote))

	err = d.Set("is_remote", resource.IsRemote)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"is_remote\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "ShareAcl", resource.ShareAcl))

	tflog.Debug(ctx, fmt.Sprintf("Found a pointer object %v", resource.ShareAcl))
	err = d.Set("share_acl", utils.FlattenModelAsList(ctx, resource.ShareAcl))

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"share_acl\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "QosPolicyId", resource.QosPolicyId))

	err = d.Set("qos_policy_id", resource.QosPolicyId)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"qos_policy_id\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "IsSeamless", resource.IsSeamless))

	err = d.Set("is_seamless", resource.IsSeamless)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"is_seamless\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "MaxRetentionPeriod", resource.MaxRetentionPeriod))

	err = d.Set("max_retention_period", resource.MaxRetentionPeriod)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"max_retention_period\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "MinRetentionPeriod", resource.MinRetentionPeriod))

	err = d.Set("min_retention_period", resource.MinRetentionPeriod)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"min_retention_period\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "FilesRetentionMode", resource.FilesRetentionMode))

	err = d.Set("files_retention_mode", resource.FilesRetentionMode)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"files_retention_mode\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "DefaultRetentionPeriod", resource.DefaultRetentionPeriod))

	err = d.Set("default_retention_period", resource.DefaultRetentionPeriod)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"default_retention_period\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "AutoCommit", resource.AutoCommit))

	err = d.Set("auto_commit", resource.AutoCommit)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"auto_commit\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "S3ObjectOwnershipRule", resource.S3ObjectOwnershipRule))

	err = d.Set("s3_object_ownership_rule", resource.S3ObjectOwnershipRule)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"s3_object_ownership_rule\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "Locking", resource.Locking))

	err = d.Set("locking", resource.Locking)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"locking\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "IgnoreOos", resource.IgnoreOos))

	err = d.Set("ignore_oos", resource.IgnoreOos)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"ignore_oos\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "BucketLogging", resource.BucketLogging))

	tflog.Debug(ctx, fmt.Sprintf("Found a pointer object %v", resource.BucketLogging))
	err = d.Set("bucket_logging", utils.FlattenModelAsList(ctx, resource.BucketLogging))

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"bucket_logging\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "AbacTags", resource.AbacTags))

	err = d.Set("abac_tags", utils.FlattenListOfPrimitives(&resource.AbacTags))

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"abac_tags\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "AbeMaxDepth", resource.AbeMaxDepth))

	err = d.Set("abe_max_depth", resource.AbeMaxDepth)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"abe_max_depth\"",
			Detail:   err.Error(),
		})
	}

	tflog.Info(ctx, fmt.Sprintf("%v - %v", "AbeProtocols", resource.AbeProtocols))

	err = d.Set("abe_protocols", utils.FlattenListOfPrimitives(&resource.AbeProtocols))

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured setting value to \"abe_protocols\"",
			Detail:   err.Error(),
		})
	}

	return diags

}
func resourceViewRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(vast_client.JwtSession)
	resource_config := codegen_configs.GetResourceByName("View")
	attrs := map[string]interface{}{"path": utils.GenPath("views"), "id": d.Id()}
	tflog.Debug(ctx, fmt.Sprintf("[resourceViewRead] Calling Get Function : %v for resource View", utils.GetFuncName(resource_config.GetFunc)))
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
	resource := api_latest.View{}
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
	diags = ResourceViewReadStructIntoSchema(ctx, resource, d)

	var after_read_error error
	after_read_error = resource_config.AfterReadFunc(client, ctx, d)
	if after_read_error != nil {
		return diag.FromErr(after_read_error)
	}

	return diags
}

func resourceViewDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(vast_client.JwtSession)
	resource_config := codegen_configs.GetResourceByName("View")
	attrs := map[string]interface{}{"path": utils.GenPath("views"), "id": d.Id()}

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

func resourceViewCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	names_mapping := utils.ContextKey("names_mapping")
	new_ctx := context.WithValue(ctx, names_mapping, View_names_mapping)
	var diags diag.Diagnostics
	data := make(map[string]interface{})
	client := m.(vast_client.JwtSession)
	resource_config := codegen_configs.GetResourceByName("View")
	tflog.Info(ctx, fmt.Sprintf("Creating Resource View"))
	reflect_View := reflect.TypeOf((*api_latest.View)(nil))
	utils.PopulateResourceMap(new_ctx, reflect_View.Elem(), d, &data, "", false)

	version_compare := utils.VastVersionsWarn(ctx)

	if version_compare != metadata.CLUSTER_VERSION_EQUALS {
		cluster_version := metadata.ClusterVersionString()
		t, t_exists := vast_versions.GetVersionedType(cluster_version, "View")
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
			tflog.Warn(ctx, fmt.Sprintf("Could have not found resource %s in version %s , things might not work properly", "View", cluster_version))
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
	attrs := map[string]interface{}{"path": utils.GenPath("views")}
	response, create_err := resource_config.CreateFunc(ctx, client, attrs, data, map[string]string{})
	tflog.Info(ctx, fmt.Sprintf("Server Error for  View %v", create_err))

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
	resource := api_latest.View{}
	err = json.Unmarshal(response_body, &resource)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to convert response body into View",
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
	resourceViewRead(ctx_with_resource, d, m)

	var before_create_error error
	_, before_create_error = resource_config.BeforeCreateFunc(data, client, ctx, d)
	if before_create_error != nil {
		return diag.FromErr(before_create_error)
	}

	return diags
}

func resourceViewUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	names_mapping := utils.ContextKey("names_mapping")
	new_ctx := context.WithValue(ctx, names_mapping, View_names_mapping)
	var diags diag.Diagnostics
	data := make(map[string]interface{})
	version_compare := utils.VastVersionsWarn(ctx)
	resource_config := codegen_configs.GetResourceByName("View")
	if version_compare != metadata.CLUSTER_VERSION_EQUALS {
		cluster_version := metadata.ClusterVersionString()
		t, t_exists := vast_versions.GetVersionedType(cluster_version, "View")
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
			tflog.Warn(ctx, fmt.Sprintf("Could have not found resource %s in version %s , things might not work properly", "View", cluster_version))
		}
	}

	client := m.(vast_client.JwtSession)
	tflog.Info(ctx, fmt.Sprintf("Updating Resource View"))
	reflect_View := reflect.TypeOf((*api_latest.View)(nil))
	utils.PopulateResourceMap(new_ctx, reflect_View.Elem(), d, &data, "", false)

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
	attrs := map[string]interface{}{"path": utils.GenPath("views"), "id": d.Id()}
	response, patch_err := resource_config.UpdateFunc(ctx, client, attrs, data, d, map[string]string{})
	tflog.Info(ctx, fmt.Sprintf("Server Error for  View %v", patch_err))
	if patch_err != nil {
		error_message := patch_err.Error() + " Server Response: " + utils.GetResponseBodyAsStr(response)
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Object Creation Failed",
			Detail:   error_message,
		})
		return diags
	}
	resourceViewRead(ctx, d, m)

	var after_patch_error error
	data, after_patch_error = resource_config.AfterPatchFunc(data, client, ctx, d)
	if after_patch_error != nil {
		return diag.FromErr(after_patch_error)
	}

	return diags

}

func resourceViewImporter(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {

	result := []*schema.ResourceData{}
	client := m.(vast_client.JwtSession)
	resource_config := codegen_configs.GetResourceByName("View")
	attrs := map[string]interface{}{"path": utils.GenPath("views")}
	response, err := resource_config.ImportFunc(ctx, client, attrs, d, resource_config.Importer.GetFunc())

	if err != nil {
		return result, err
	}

	resource_l := []api_latest.View{}
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

	diags := ResourceViewReadStructIntoSchema(ctx, resource, d)
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
