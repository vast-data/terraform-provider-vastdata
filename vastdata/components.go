// Copyright (c) HashiCorp, Inc.

package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
)

// allTFComponents holds all Terraform-managed entities (both resources and data sources)
// that implement TFManager. These are later filtered based on their supported interfaces.
var allTFComponents = []TFManager{
	&User{},
	&Group{},
	&VipPool{},
	&UserKey{},
	&Tenant{},
	&ViewPolicy{},
	&View{},
	&Snapshot{},
	&S3Policy{},
	&S3LifeCycleRule{},
	&ReplicationPeer{},
	&Quota{},
	&QosPolicy{},
	&ActiveDirectory{},
	&Ldap{},
	&AdministratorManager{},
	&AdministratorRealm{},
	&AdministratorRole{},
	&Dns{},
	&GlobalSnapshot{},
	&GlobalLocalSnapshot{},
	&Nis{},
	&NonlocalGroup{},
	&NonlocalUser{},
	&NonlocalUserKey{},
	&ProtectedPath{},
	&ProtectionPolicy{},
	&KafkaBroker{},
	&BlockHostMapping{},
	&EventDefinition{},
	&EventDefinitionConfig{},
	&EncryptionGroup{},
	&EncryptionGroupControl{},
	&UserCopy{},
	&FolderReadOnly{},
	&UserTenantData{},
	&LocalS3Key{},
	&LocalProvider{},
	&ApiToken{},
	&Vms{},
	&Volume{},
	&BgpConfig{},
	&S3PolicyAttachment{},
	&TenantEncryptionGroupControl{},
	&TenantClientMetrics{},
	&TenantConfiguredIdp{},
	//&BlockHost{},
}

// GetResourceFactories returns a list of factory functions that instantiate
// Terraform resources supported by the provider.
//
// Only components implementing the ResourceManager interface will be included.
func GetResourceFactories() []func() resource.Resource {
	var factories []func() resource.Resource
	for _, f := range allTFComponents {
		if manager, ok := f.(ResourceManager); ok {
			managerFn := manager.NewResourceManager
			managerType := is.SnakeCaseName(f)

			factories = append(factories, func() resource.Resource {
				return &Resource{
					newManager:  managerFn,
					managerName: managerType,
				}
			})
		}
	}
	return factories

}

// GetDatasourceFactories returns a list of factory functions that instantiate
// Terraform data sources supported by the provider.
//
// Only components implementing the DataSourceManager interface will be included.
func GetDatasourceFactories() []func() datasource.DataSource {
	var factories []func() datasource.DataSource
	for _, f := range allTFComponents {

		if manager, ok := f.(DataSourceManager); ok {
			managerFn := manager.NewDatasourceManager
			managerType := is.SnakeCaseName(f)

			factories = append(factories, func() datasource.DataSource {
				return &Datasource{
					newManager:  managerFn,
					managerName: managerType,
				}
			})
		}
	}
	return factories

}

type ResourceFactoryFn func(raw map[string]attr.Value, schema any) ResourceManager
type DatasourceFactoryFn func(raw map[string]attr.Value, schema any) DataSourceManager

// -----------------
// VAST API interfaces
// -----------------

type VastAPIGetter interface {
	API(*VMSRest) VastResourceAPIWithContext
}

// -----------------
// Resource interfaces
// -----------------

type TFManager interface {
	VastAPIGetter
	TfState() *is.TFState
}

type PrepareImportResourceState interface {
	PrepareImportResourceState(context.Context, *VMSRest) error
}

type PrepareCreateResource interface {
	PrepareCreateResource(context.Context, *VMSRest) error
}

type PrepareReadResource interface {
	PrepareReadResource(context.Context, *VMSRest) error
}

type PrepareUpdateResource interface {
	PrepareUpdateResource(context.Context, PrepareUpdateResource, *VMSRest) error
}

type PrepareDeleteResource interface {
	PrepareDeleteResource(context.Context, *VMSRest) error
}

type ImportResourceState interface {
	ImportResourceState(context.Context, *VMSRest) (DisplayableRecord, error)
}

type CreateResource interface {
	CreateResource(context.Context, *VMSRest) (DisplayableRecord, error)
}

type ReadResource interface {
	ReadResource(context.Context, *VMSRest) (DisplayableRecord, error)
}

type UpdateResource interface {
	UpdateResource(context.Context, UpdateResource, *VMSRest) (DisplayableRecord, error)
}

type DeleteResource interface {
	DeleteResource(context.Context, *VMSRest) error
}

type AfterImportResourceState interface {
	AfterImportResourceState(context.Context, *VMSRest) error
}

type AfterCreateResource interface {
	AfterCreateResource(context.Context, *VMSRest, Record) error
}

type AfterReadResource interface {
	AfterReadResource(context.Context, *VMSRest) error
}

type AfterUpdateResource interface {
	AfterUpdateResource(context.Context, AfterUpdateResource, *VMSRest, Record) error
}

type AfterDeleteResource interface {
	AfterDeleteResource(context.Context, *VMSRest) error
}

type ResourceManager interface {
	TFManager
	NewResourceManager(raw map[string]attr.Value, schema any) ResourceManager
}

// -----------------
// Datasource interfaces
// -----------------

type PrepareReadDatasource interface {
	PrepareReadDatasource(context.Context, *VMSRest) error
}

type ReadDatasource interface {
	ReadDatasource(context.Context, *VMSRest) (DisplayableRecord, error)
}

type AfterReadDatasource interface {
	AfterReadDatasource(context.Context, *VMSRest) error
}

type DataSourceManager interface {
	TFManager
	NewDatasourceManager(raw map[string]attr.Value, schema any) DataSourceManager
}

// -----------------
// Validate (Resource) interfaces
// -----------------

type ValidateResourceConfig interface {
	ValidateResourceConfig(context.Context) error
}

// -----------------
// Transform interfaces
// -----------------

// TransformRequestBody allows a resource to modify the outgoing request body
// before it is sent in Create or Update operations. This can be used to inject
// derived fields, sanitize inputs, or restructure data.
//
// The `params` type is assumed to be a map[string]any.
type TransformRequestBody interface {
	TransformRequestBody(body params) params
}

// TransformResponseRecord allows a resource to modify the backend response
// before it is saved into Terraform state. This can be used to normalize
// fields, decode values, or inject defaults.
//
// The `Record` type is assumed to be a map[string]any.
type TransformResponseRecord interface {
	TransformResponseRecord(response Record) Record
}
