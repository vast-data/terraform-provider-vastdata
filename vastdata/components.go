// Copyright (c) HashiCorp, Inc.

package provider

import (
	"context"
	"fmt"
	"sort"

	version "github.com/hashicorp/go-version"
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
	&TlsCertificate{},
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
	&TenantMetricLabelValuesBulk{},
	&TenantConfiguredIdp{},
	&S3ReplicationPeer{},
	&SamlConfig{},
	&VmsConfiguredIdps{},
	&BlockHost{},
	&Kerberos{},
	&KerberosKeytab{},
	&Cluster{},
	&ClusterEkm{},
	&ComputeCluster{},
	&ComputeClusterControl{},
	&ComputeClusterNodes{},
	&ComputeClusterPods{},
	&ComputeClusterNamespaces{},
	&ComputeClusterServices{},
	&ComputeClusterDeployments{},
	&ComputeClusterTenants{},
	&ComputeClusterDashboard{},
	&ComputeClusterReplicaSets{},
	&ComputeClusterEvents{},
	&Cnode{},
	&CnodeBgpConfig{},
	&Rack{},
	&RackBgpConfig{},
	&NicPortRelatedPorts{},
	&ManagerAuthorizedStatus{},
	&ManagerPassword{},
	&Topic{},
	&IamRole{},
	&IamRoleCredentials{},
	&Oidc{},
	&VastDbVips{},
	&TenantNfs4Delegation{},
	&TenantMetricLabels{},
	&TenantMetricLabelValues{},
	&SupportedDrives{},
	&SupportBundlesQueue{},
	&Webhook{},
	&QuotaGroup{},
	&BlobExpansion{},
	&Certificate{},
	&VmsCertificate{},
	&S3Certificate{},
	&NfsCertificate{},
	&License{},
	&CallhomeConfig{},
	&VmsNetworkSettings{},
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

// -----------------
// Runtime-agnostic registry
// -----------------

// ResourceManagerPrototypes returns every component that can be managed as a resource,
// keyed by its snake_case name (for example "vip_pool"). The returned managers are
// prototypes: use NewResourceManagerByName to obtain a usable manager.
//
// This is the same component list Terraform registers via GetResourceFactories, so any
// other runtime (e.g. Crossplane) automatically sees the same set of resources.
// The map is built once and then shared, so callers must not modify it.
var ResourceManagerPrototypes = cached(resourceManagerPrototypes)

func resourceManagerPrototypes() map[string]ResourceManager {
	out := make(map[string]ResourceManager)
	for _, f := range allTFComponents {
		if manager, ok := f.(ResourceManager); ok {
			out[is.SnakeCaseName(f)] = manager
		}
	}
	return out
}

// ResourceNames returns the registry names of every resource the provider manages,
// sorted. It is exactly the set of resources the Terraform provider exposes.
var ResourceNames = cached(resourceNames)

func resourceNames() []string {
	protos := ResourceManagerPrototypes()
	names := make([]string, 0, len(protos))
	for name := range protos {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// NewResourceManagerByName builds a manager for the named resource. Its state carries the
// generated schema and a zeroed (all-null) attribute map, ready to be seeded by the caller.
// No Terraform request/response types are involved.
func NewResourceManagerByName(ctx context.Context, name string) (ResourceManager, error) {
	prototype, ok := ResourceManagerPrototypes()[name]
	if !ok {
		return nil, fmt.Errorf("unknown resource %q", name)
	}
	return newManagerWithSchemaOnly(ctx, prototype.NewResourceManager)
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
	TfState() is.ResourceState
}

type PrepareImportResourceState interface {
	PrepareImportResourceState(context.Context, *VMSRest) error
}

type PrepareCreateResource interface {
	PrepareCreateResource(context.Context, *VMSRest) error
}

// GetSubResources fetches nested sub-endpoint data and returns it as a Record.
// The framework merges the returned Record into the main resource record
// (flat — all keys are merged directly) before FillFromRecord is called,
// so state is populated in one pass without any manual tfstate manipulation.
// Return nil Record (with nil error) to skip merging.
// GetSubResources is implemented by resources that expose nested sub-endpoints.
// clusterVersion is the connected cluster's version (from GetCachedClusterVersion),
// or nil when no sub-resource hint declares MinVastVersion.
// Implementations should use clusterVersion (when non-nil) to gate fetches rather
// than calling GetCachedClusterVersion themselves.
type GetSubResources interface {
	GetSubResources(ctx context.Context, rest *VMSRest, record Record, clusterVersion *version.Version) (Record, error)
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

type ModifyPlanResource interface {
	ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) error
}

type ImportResourceState interface {
	ImportResourceState(resource.ImportStateRequest, context.Context, *VMSRest) error
}

type CreateResource interface {
	CreateResource(context.Context, *VMSRest) (DisplayableRecord, error)
}

type LookupForCreate interface {
	LookupForCreate(context.Context, *VMSRest) (DisplayableRecord, error)
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

// -----------------
// Resource callbacks
// -----------------
//
// A callback is an extra function attached to a manager instance from outside, e.g. by the
// Crossplane provider. Each Prepare*/After* stage of Create, Read, Update and Delete first
// runs the manager's own hook when it implements one. If that hook fails, its error is
// returned and no callback runs. The callbacks then run in the order they were added, and
// the first error stops the rest.
//
// Callbacks are not hooks: a manager does not implement PrepareCreateResource just because
// a callback is added for that stage, so defaults that depend on the resource's own hook
// (e.g. edit-only fields handled by the default After* behaviour) are not affected.

type PrepareCreateResourceCallbackExecutor interface {
	AddPrepareCreateResourceCallback(func(context.Context, *VMSRest) error)
	ExecutePrepareCreateResourceCallbacks(context.Context, *VMSRest) error
}

type AfterCreateResourceCallbackExecutor interface {
	AddAfterCreateResourceCallback(func(context.Context, *VMSRest, Record) error)
	ExecuteAfterCreateResourceCallbacks(context.Context, *VMSRest, Record) error
}

type PrepareReadResourceCallbackExecutor interface {
	AddPrepareReadResourceCallback(func(context.Context, *VMSRest) error)
	ExecutePrepareReadResourceCallbacks(context.Context, *VMSRest) error
}

type AfterReadResourceCallbackExecutor interface {
	AddAfterReadResourceCallback(func(context.Context, *VMSRest) error)
	ExecuteAfterReadResourceCallbacks(context.Context, *VMSRest) error
}

// PrepareUpdateResourceCallbackExecutor callbacks receive the desired manager.
type PrepareUpdateResourceCallbackExecutor interface {
	AddPrepareUpdateResourceCallback(func(context.Context, ResourceManager, *VMSRest) error)
	ExecutePrepareUpdateResourceCallbacks(context.Context, ResourceManager, *VMSRest) error
}

// AfterUpdateResourceCallbackExecutor callbacks receive the desired manager.
type AfterUpdateResourceCallbackExecutor interface {
	AddAfterUpdateResourceCallback(func(context.Context, ResourceManager, *VMSRest, Record) error)
	ExecuteAfterUpdateResourceCallbacks(context.Context, ResourceManager, *VMSRest, Record) error
}

type PrepareDeleteResourceCallbackExecutor interface {
	AddPrepareDeleteResourceCallback(func(context.Context, *VMSRest) error)
	ExecutePrepareDeleteResourceCallbacks(context.Context, *VMSRest) error
}

type AfterDeleteResourceCallbackExecutor interface {
	AddAfterDeleteResourceCallback(func(context.Context, *VMSRest) error)
	ExecuteAfterDeleteResourceCallbacks(context.Context, *VMSRest) error
}

// ResourceCallbackExecutor is every stage's callback executor.
type ResourceCallbackExecutor interface {
	PrepareCreateResourceCallbackExecutor
	AfterCreateResourceCallbackExecutor
	PrepareReadResourceCallbackExecutor
	AfterReadResourceCallbackExecutor
	PrepareUpdateResourceCallbackExecutor
	AfterUpdateResourceCallbackExecutor
	PrepareDeleteResourceCallbackExecutor
	AfterDeleteResourceCallbackExecutor
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
// MigrateMode interfaces
// -----------------

// MigrateModePassThroughUpdate is implemented by resource managers that must
// allow updates during VASTDATA_MIGRATE_MODE.  Normally all updates are
// blocked in migrate mode because an update implies the resource is already
// in state — which is unexpected when migrating from an empty state.
//
// Non-importable resources (e.g. vastdata_user_key) are intentionally
// carried over from the old state verbatim.  When the new provider schema
// differs from the old one, Terraform plans a schema-reconciliation update.
// Implementing this interface allows that update to proceed instead of
// failing with "update operation blocked".
type MigrateModePassThroughUpdate interface {
	MigrateModePassThroughUpdate()
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

// NormalizeRecordForCreateAdopt allows a resource to normalize an existing API record
// before adopt-on-create compares it to the user's create parameters.
type NormalizeRecordForCreateAdopt interface {
	NormalizeRecordForCreateAdopt(Record) Record
}

// ResolveRecordAfterAsyncTask allows a resource to replace the API response record
// after an async task completes. Some endpoints return the async task body (with the
// task id as record id) instead of the parent resource.
type ResolveRecordAfterAsyncTask interface {
	ResolveRecordAfterAsyncTask(context.Context, *VMSRest, Record) (Record, error)
}
