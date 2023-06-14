package datasources

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api_latest "github.com/vast-data/terraform-provider-vastdata/codegen/latest"
	utils "github.com/vast-data/terraform-provider-vastdata/utils"
	vast_client "github.com/vast-data/terraform-provider-vastdata/vast-client"
	"net/url"
	"strconv"
)

func DataSourceView() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceViewRead,
		Description: ``,
		Schema: map[string]*schema.Schema{

			"id": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `A uniqe ID used to identify the View`,
			},

			"guid": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `A uniqe GUID assigned to the View`,
			},

			"name": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `A uniq name given to the view`,
			},

			"path": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    false,
				Required:    true,
				Optional:    false,
				Description: `File system path. Begin with '/'. Do not include a trailing slash`,
			},

			"create_dir": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `Creates the directory specified by the path`,
			},

			"alias": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `Alias for NFS export, must start with '/' and only ASCII characters are allowed. If configured, this supersedes the exposed NFS export path`,
			},

			"bucket": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `S3 Bucket name`,
			},

			"policy_id": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `Associated view policy ID`,
			},

			"cluster": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `Parent Cluster`,
			},

			"cluster_id": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `Parent Cluster ID`,
			},

			"tenant_id": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Required:    false,
				Optional:    true,
				Description: `The tenant ID related to this view`,
			},

			"directory": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `Create the directory if it does not exist`,
			},

			"s3_versioning": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `Trun on S3 Versioning`,
			},

			"s3_unverified_lookup": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `Allow S3 Unverified Lookup`,
			},

			"allow_anonymous_access": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `Allow S3 anonymous access`,
			},

			"allow_s3_anonymous_access": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `Allow S3 anonymous access`,
			},

			"protocols": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `Protocols exposed by this view`,

				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"share": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `Name of the SMB Share. Must not include the following characters: " \ / [ ] : | < > + = ; , * ?`,
			},

			"bucket_owner": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `S3 Bucket owner`,
			},

			"bucket_creators": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `List of bucket creators users`,

				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"bucket_creators_groups": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `List of bucket creators groups`,

				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"s3_locks": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `S3 Object Lock`,
			},

			"s3_locks_retention_mode": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `S3 Locks retention mode`,
			},

			"s3_locks_retention_period": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `Period should be positive in format like 0d|2d|1y|2y`,
			},

			"physical_capacity": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `Physical Capacity`,
			},

			"logical_capacity": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `Logical Capacity`,
			},

			"nfs_interop_flags": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `Indicates whether the view should support simultaneous access to NFS3/NFS4/SMB protocols. Allowed Values are [BOTH_NFS3_AND_NFS4_INTEROP_DISABLED ONLY_NFS3_INTEROP_ENABLED ONLY_NFS4_INTEROP_ENABLED BOTH_NFS3_AND_NFS4_INTEROP_ENABLED]`,
			},

			"is_remote": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: ``,
			},

			"share_acl": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `Share-level ACL details`,

				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						"enabled": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Required:    false,
							Optional:    false,
							Description: ``,
						},

						"acl": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Required:    false,
							Optional:    false,
							Description: ``,

							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{

									"grantee": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Required:    false,
										Optional:    false,
										Description: ` Allowed Values are [users groups]`,
									},

									"permissions": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Required:    false,
										Optional:    false,
										Description: ` Allowed Values are [FULL]`,
									},

									"sid_str": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Required:    false,
										Optional:    false,
										Description: ``,
									},

									"uid_or_gid": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Required:    false,
										Optional:    false,
										Description: ``,
									},

									"name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    false,
										Required:    true,
										Optional:    false,
										Description: ``,
									},

									"fqdn": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Required:    false,
										Optional:    false,
										Description: ``,
									},
								},
							},
						},
					},
				},
			},

			"qos_policy_id": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: `QoS Policy ID`,
			},
		},
	}
}

func dataSourceViewRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(vast_client.JwtSession)
	values := url.Values{}

	path := d.Get("path")
	values.Add("path", fmt.Sprintf("%v", path))

	if d.HasChanges("tenant_id") {
		tenant_id := d.Get("tenant_id")
		tflog.Debug(ctx, "Using optional attribute \"tenant_id\"")
		values.Add("tenant_id", fmt.Sprintf("%v", tenant_id))
	}

	response, err := client.Get(ctx, "/api/views/", values.Encode(), map[string]string{})
	tflog.Info(ctx, response.Request.URL.String())
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured while obtaining data from the vastdata cluster",
			Detail:   err.Error(),
		})
		return diags

	}
	resource_l := []api_latest.View{}

	body, err := utils.DefaultProcessingFunc(ctx, response)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured reading data recived from VastData cluster",
			Detail:   err.Error(),
		})
		return diags

	}
	err = json.Unmarshal(body, &resource_l)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error occured while parsing data recived from VastData cluster",
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
			Summary:  "Error occured setting value to \"id\"",
			Detail:   err.Error(),
		})
	}

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

	Id := (int64)(resource.Id)
	d.SetId(strconv.FormatInt(Id, 10))
	return diags
}
