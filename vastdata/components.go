// Copyright (c) HashiCorp, Inc.

package provider

import (
	"context"
	"errors"
	"fmt"
	"reflect"
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

// A setter wraps a manager with one hook.
//
// findHook returns the outermost implementation of T (first match). A later
// setter for the same exclusive hook (Create/Read/Update/Delete, …) is the one
// that runs.
//
// forEachHook walks the full wrapper chain and invokes every layer that
// actually provides T, including the original manager. Errors from all
// invocations are joined with errors.Join and returned at the end.
//
// Wrappers embed the inner manager, so they type-assert to every interface the
// inner already implements (method promotion). hookProvider.providedHook
// identifies which interface a wrapper layer itself adds, so forEachHook can
// skip those false matches. A nil function leaves the manager unchanged.

type hookUnwrapper interface {
	unwrapHook() any
}

// hookProvider is implemented by With* wrappers. providedHook returns a typed
// nil pointer of the interface this layer adds, e.g. (*AfterDeleteResource)(nil).
type hookProvider interface {
	providedHook() any
}

func findHook[T any](m any) (T, bool) {
	var zero T
	for i := 0; i < 64 && m != nil; i++ {
		if v, ok := m.(T); ok {
			return v, true
		}
		u, ok := m.(hookUnwrapper)
		if !ok {
			return zero, false
		}
		m = u.unwrapHook()
	}
	return zero, false
}

func forEachHook[T any](m any, fn func(T) error) error {
	var errs []error
	want := reflect.TypeFor[*T]()
	for i := 0; i < 64 && m != nil; i++ {
		if p, ok := m.(hookProvider); ok {
			if reflect.TypeOf(p.providedHook()) == want {
				if err := fn(m.(T)); err != nil {
					errs = append(errs, err)
				}
			}
		} else if _, wrapped := m.(hookUnwrapper); !wrapped {
			if v, ok := m.(T); ok {
				if err := fn(v); err != nil {
					errs = append(errs, err)
				}
			}
		}
		u, ok := m.(hookUnwrapper)
		if !ok {
			break
		}
		m = u.unwrapHook()
	}
	return errors.Join(errs...)
}

// hasOwnHook reports whether the original manager, the innermost value of the
// wrapper chain, implements T itself. Wrappers added by With* setters do not count.
//
// The default edit-only-fields handling in Create and Update is skipped when the
// resource handles them in its own After hook. An extra After hook attached from
// outside (for example by the Crossplane provider) must not disable that default.
func hasOwnHook[T any](m any) bool {
	for i := 0; i < 64 && m != nil; i++ {
		if _, wrapper := m.(hookProvider); !wrapper {
			if _, ok := m.(T); ok {
				return true
			}
		}
		u, ok := m.(hookUnwrapper)
		if !ok {
			return false
		}
		m = u.unwrapHook()
	}
	return false
}

func WithPrepareImportResourceState(m ResourceManager, fn func(context.Context, *VMSRest) error) ResourceManager {
	if m == nil || fn == nil {
		return m
	}
	return prepareImportResourceStateHook{ResourceManager: m, fn: fn}
}

type prepareImportResourceStateHook struct {
	ResourceManager
	fn func(context.Context, *VMSRest) error
}

func (h prepareImportResourceStateHook) unwrapHook() any   { return h.ResourceManager }
func (h prepareImportResourceStateHook) providedHook() any { return (*PrepareImportResourceState)(nil) }
func (h prepareImportResourceStateHook) PrepareImportResourceState(ctx context.Context, rest *VMSRest) error {
	return h.fn(ctx, rest)
}

func WithPrepareCreateResource(m ResourceManager, fn func(context.Context, *VMSRest) error) ResourceManager {
	if m == nil || fn == nil {
		return m
	}
	return prepareCreateResourceHook{ResourceManager: m, fn: fn}
}

type prepareCreateResourceHook struct {
	ResourceManager
	fn func(context.Context, *VMSRest) error
}

func (h prepareCreateResourceHook) unwrapHook() any   { return h.ResourceManager }
func (h prepareCreateResourceHook) providedHook() any { return (*PrepareCreateResource)(nil) }
func (h prepareCreateResourceHook) PrepareCreateResource(ctx context.Context, rest *VMSRest) error {
	return h.fn(ctx, rest)
}

func WithGetSubResources(m ResourceManager, fn func(context.Context, *VMSRest, Record, *version.Version) (Record, error)) ResourceManager {
	if m == nil || fn == nil {
		return m
	}
	return getSubResourcesHook{ResourceManager: m, fn: fn}
}

type getSubResourcesHook struct {
	ResourceManager
	fn func(context.Context, *VMSRest, Record, *version.Version) (Record, error)
}

func (h getSubResourcesHook) unwrapHook() any   { return h.ResourceManager }
func (h getSubResourcesHook) providedHook() any { return (*GetSubResources)(nil) }
func (h getSubResourcesHook) GetSubResources(ctx context.Context, rest *VMSRest, record Record, clusterVersion *version.Version) (Record, error) {
	return h.fn(ctx, rest, record, clusterVersion)
}

func WithPrepareReadResource(m ResourceManager, fn func(context.Context, *VMSRest) error) ResourceManager {
	if m == nil || fn == nil {
		return m
	}
	return prepareReadResourceHook{ResourceManager: m, fn: fn}
}

type prepareReadResourceHook struct {
	ResourceManager
	fn func(context.Context, *VMSRest) error
}

func (h prepareReadResourceHook) unwrapHook() any   { return h.ResourceManager }
func (h prepareReadResourceHook) providedHook() any { return (*PrepareReadResource)(nil) }
func (h prepareReadResourceHook) PrepareReadResource(ctx context.Context, rest *VMSRest) error {
	return h.fn(ctx, rest)
}

func WithPrepareUpdateResource(m ResourceManager, fn func(context.Context, PrepareUpdateResource, *VMSRest) error) ResourceManager {
	if m == nil || fn == nil {
		return m
	}
	return prepareUpdateResourceHook{ResourceManager: m, fn: fn}
}

type prepareUpdateResourceHook struct {
	ResourceManager
	fn func(context.Context, PrepareUpdateResource, *VMSRest) error
}

func (h prepareUpdateResourceHook) unwrapHook() any   { return h.ResourceManager }
func (h prepareUpdateResourceHook) providedHook() any { return (*PrepareUpdateResource)(nil) }
func (h prepareUpdateResourceHook) PrepareUpdateResource(ctx context.Context, plan PrepareUpdateResource, rest *VMSRest) error {
	return h.fn(ctx, plan, rest)
}

func WithPrepareDeleteResource(m ResourceManager, fn func(context.Context, *VMSRest) error) ResourceManager {
	if m == nil || fn == nil {
		return m
	}
	return prepareDeleteResourceHook{ResourceManager: m, fn: fn}
}

type prepareDeleteResourceHook struct {
	ResourceManager
	fn func(context.Context, *VMSRest) error
}

func (h prepareDeleteResourceHook) unwrapHook() any   { return h.ResourceManager }
func (h prepareDeleteResourceHook) providedHook() any { return (*PrepareDeleteResource)(nil) }
func (h prepareDeleteResourceHook) PrepareDeleteResource(ctx context.Context, rest *VMSRest) error {
	return h.fn(ctx, rest)
}

func WithModifyPlanResource(m ResourceManager, fn func(context.Context, resource.ModifyPlanRequest, *resource.ModifyPlanResponse) error) ResourceManager {
	if m == nil || fn == nil {
		return m
	}
	return modifyPlanResourceHook{ResourceManager: m, fn: fn}
}

type modifyPlanResourceHook struct {
	ResourceManager
	fn func(context.Context, resource.ModifyPlanRequest, *resource.ModifyPlanResponse) error
}

func (h modifyPlanResourceHook) unwrapHook() any   { return h.ResourceManager }
func (h modifyPlanResourceHook) providedHook() any { return (*ModifyPlanResource)(nil) }
func (h modifyPlanResourceHook) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) error {
	return h.fn(ctx, req, resp)
}

func WithImportResourceState(m ResourceManager, fn func(resource.ImportStateRequest, context.Context, *VMSRest) error) ResourceManager {
	if m == nil || fn == nil {
		return m
	}
	return importResourceStateHook{ResourceManager: m, fn: fn}
}

type importResourceStateHook struct {
	ResourceManager
	fn func(resource.ImportStateRequest, context.Context, *VMSRest) error
}

func (h importResourceStateHook) unwrapHook() any   { return h.ResourceManager }
func (h importResourceStateHook) providedHook() any { return (*ImportResourceState)(nil) }
func (h importResourceStateHook) ImportResourceState(req resource.ImportStateRequest, ctx context.Context, rest *VMSRest) error {
	return h.fn(req, ctx, rest)
}

func WithCreateResource(m ResourceManager, fn func(context.Context, *VMSRest) (DisplayableRecord, error)) ResourceManager {
	if m == nil || fn == nil {
		return m
	}
	return createResourceHook{ResourceManager: m, fn: fn}
}

type createResourceHook struct {
	ResourceManager
	fn func(context.Context, *VMSRest) (DisplayableRecord, error)
}

func (h createResourceHook) unwrapHook() any   { return h.ResourceManager }
func (h createResourceHook) providedHook() any { return (*CreateResource)(nil) }
func (h createResourceHook) CreateResource(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	return h.fn(ctx, rest)
}

func WithLookupForCreate(m ResourceManager, fn func(context.Context, *VMSRest) (DisplayableRecord, error)) ResourceManager {
	if m == nil || fn == nil {
		return m
	}
	return lookupForCreateHook{ResourceManager: m, fn: fn}
}

type lookupForCreateHook struct {
	ResourceManager
	fn func(context.Context, *VMSRest) (DisplayableRecord, error)
}

func (h lookupForCreateHook) unwrapHook() any   { return h.ResourceManager }
func (h lookupForCreateHook) providedHook() any { return (*LookupForCreate)(nil) }
func (h lookupForCreateHook) LookupForCreate(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	return h.fn(ctx, rest)
}

func WithReadResource(m ResourceManager, fn func(context.Context, *VMSRest) (DisplayableRecord, error)) ResourceManager {
	if m == nil || fn == nil {
		return m
	}
	return readResourceHook{ResourceManager: m, fn: fn}
}

type readResourceHook struct {
	ResourceManager
	fn func(context.Context, *VMSRest) (DisplayableRecord, error)
}

func (h readResourceHook) unwrapHook() any   { return h.ResourceManager }
func (h readResourceHook) providedHook() any { return (*ReadResource)(nil) }
func (h readResourceHook) ReadResource(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	return h.fn(ctx, rest)
}

func WithUpdateResource(m ResourceManager, fn func(context.Context, UpdateResource, *VMSRest) (DisplayableRecord, error)) ResourceManager {
	if m == nil || fn == nil {
		return m
	}
	return updateResourceHook{ResourceManager: m, fn: fn}
}

type updateResourceHook struct {
	ResourceManager
	fn func(context.Context, UpdateResource, *VMSRest) (DisplayableRecord, error)
}

func (h updateResourceHook) unwrapHook() any   { return h.ResourceManager }
func (h updateResourceHook) providedHook() any { return (*UpdateResource)(nil) }
func (h updateResourceHook) UpdateResource(ctx context.Context, plan UpdateResource, rest *VMSRest) (DisplayableRecord, error) {
	return h.fn(ctx, plan, rest)
}

func WithDeleteResource(m ResourceManager, fn func(context.Context, *VMSRest) error) ResourceManager {
	if m == nil || fn == nil {
		return m
	}
	return deleteResourceHook{ResourceManager: m, fn: fn}
}

type deleteResourceHook struct {
	ResourceManager
	fn func(context.Context, *VMSRest) error
}

func (h deleteResourceHook) unwrapHook() any   { return h.ResourceManager }
func (h deleteResourceHook) providedHook() any { return (*DeleteResource)(nil) }
func (h deleteResourceHook) DeleteResource(ctx context.Context, rest *VMSRest) error {
	return h.fn(ctx, rest)
}

func WithAfterImportResourceState(m ResourceManager, fn func(context.Context, *VMSRest) error) ResourceManager {
	if m == nil || fn == nil {
		return m
	}
	return afterImportResourceStateHook{ResourceManager: m, fn: fn}
}

type afterImportResourceStateHook struct {
	ResourceManager
	fn func(context.Context, *VMSRest) error
}

func (h afterImportResourceStateHook) unwrapHook() any   { return h.ResourceManager }
func (h afterImportResourceStateHook) providedHook() any { return (*AfterImportResourceState)(nil) }
func (h afterImportResourceStateHook) AfterImportResourceState(ctx context.Context, rest *VMSRest) error {
	return h.fn(ctx, rest)
}

func WithAfterCreateResource(m ResourceManager, fn func(context.Context, *VMSRest, Record) error) ResourceManager {
	if m == nil || fn == nil {
		return m
	}
	return afterCreateResourceHook{ResourceManager: m, fn: fn}
}

type afterCreateResourceHook struct {
	ResourceManager
	fn func(context.Context, *VMSRest, Record) error
}

func (h afterCreateResourceHook) unwrapHook() any   { return h.ResourceManager }
func (h afterCreateResourceHook) providedHook() any { return (*AfterCreateResource)(nil) }
func (h afterCreateResourceHook) AfterCreateResource(ctx context.Context, rest *VMSRest, record Record) error {
	return h.fn(ctx, rest, record)
}

func WithAfterReadResource(m ResourceManager, fn func(context.Context, *VMSRest) error) ResourceManager {
	if m == nil || fn == nil {
		return m
	}
	return afterReadResourceHook{ResourceManager: m, fn: fn}
}

type afterReadResourceHook struct {
	ResourceManager
	fn func(context.Context, *VMSRest) error
}

func (h afterReadResourceHook) unwrapHook() any   { return h.ResourceManager }
func (h afterReadResourceHook) providedHook() any { return (*AfterReadResource)(nil) }
func (h afterReadResourceHook) AfterReadResource(ctx context.Context, rest *VMSRest) error {
	return h.fn(ctx, rest)
}

func WithAfterUpdateResource(m ResourceManager, fn func(context.Context, AfterUpdateResource, *VMSRest, Record) error) ResourceManager {
	if m == nil || fn == nil {
		return m
	}
	return afterUpdateResourceHook{ResourceManager: m, fn: fn}
}

type afterUpdateResourceHook struct {
	ResourceManager
	fn func(context.Context, AfterUpdateResource, *VMSRest, Record) error
}

func (h afterUpdateResourceHook) unwrapHook() any   { return h.ResourceManager }
func (h afterUpdateResourceHook) providedHook() any { return (*AfterUpdateResource)(nil) }
func (h afterUpdateResourceHook) AfterUpdateResource(ctx context.Context, plan AfterUpdateResource, rest *VMSRest, record Record) error {
	return h.fn(ctx, plan, rest, record)
}

func WithAfterDeleteResource(m ResourceManager, fn func(context.Context, *VMSRest) error) ResourceManager {
	if m == nil || fn == nil {
		return m
	}
	return afterDeleteResourceHook{ResourceManager: m, fn: fn}
}

type afterDeleteResourceHook struct {
	ResourceManager
	fn func(context.Context, *VMSRest) error
}

func (h afterDeleteResourceHook) unwrapHook() any   { return h.ResourceManager }
func (h afterDeleteResourceHook) providedHook() any { return (*AfterDeleteResource)(nil) }
func (h afterDeleteResourceHook) AfterDeleteResource(ctx context.Context, rest *VMSRest) error {
	return h.fn(ctx, rest)
}

func WithPrepareReadDatasource(m DataSourceManager, fn func(context.Context, *VMSRest) error) DataSourceManager {
	if m == nil || fn == nil {
		return m
	}
	return prepareReadDatasourceHook{DataSourceManager: m, fn: fn}
}

type prepareReadDatasourceHook struct {
	DataSourceManager
	fn func(context.Context, *VMSRest) error
}

func (h prepareReadDatasourceHook) unwrapHook() any   { return h.DataSourceManager }
func (h prepareReadDatasourceHook) providedHook() any { return (*PrepareReadDatasource)(nil) }
func (h prepareReadDatasourceHook) PrepareReadDatasource(ctx context.Context, rest *VMSRest) error {
	return h.fn(ctx, rest)
}

func WithReadDatasource(m DataSourceManager, fn func(context.Context, *VMSRest) (DisplayableRecord, error)) DataSourceManager {
	if m == nil || fn == nil {
		return m
	}
	return readDatasourceHook{DataSourceManager: m, fn: fn}
}

type readDatasourceHook struct {
	DataSourceManager
	fn func(context.Context, *VMSRest) (DisplayableRecord, error)
}

func (h readDatasourceHook) unwrapHook() any   { return h.DataSourceManager }
func (h readDatasourceHook) providedHook() any { return (*ReadDatasource)(nil) }
func (h readDatasourceHook) ReadDatasource(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	return h.fn(ctx, rest)
}

func WithAfterReadDatasource(m DataSourceManager, fn func(context.Context, *VMSRest) error) DataSourceManager {
	if m == nil || fn == nil {
		return m
	}
	return afterReadDatasourceHook{DataSourceManager: m, fn: fn}
}

type afterReadDatasourceHook struct {
	DataSourceManager
	fn func(context.Context, *VMSRest) error
}

func (h afterReadDatasourceHook) unwrapHook() any   { return h.DataSourceManager }
func (h afterReadDatasourceHook) providedHook() any { return (*AfterReadDatasource)(nil) }
func (h afterReadDatasourceHook) AfterReadDatasource(ctx context.Context, rest *VMSRest) error {
	return h.fn(ctx, rest)
}

func WithValidateResourceConfig(m ResourceManager, fn func(context.Context) error) ResourceManager {
	if m == nil || fn == nil {
		return m
	}
	return validateResourceConfigHook{ResourceManager: m, fn: fn}
}

type validateResourceConfigHook struct {
	ResourceManager
	fn func(context.Context) error
}

func (h validateResourceConfigHook) unwrapHook() any   { return h.ResourceManager }
func (h validateResourceConfigHook) providedHook() any { return (*ValidateResourceConfig)(nil) }
func (h validateResourceConfigHook) ValidateResourceConfig(ctx context.Context) error {
	return h.fn(ctx)
}

func WithMigrateModePassThroughUpdate(m ResourceManager) ResourceManager {
	if m == nil {
		return nil
	}
	return migrateModePassThroughUpdateHook{ResourceManager: m}
}

type migrateModePassThroughUpdateHook struct{ ResourceManager }

func (h migrateModePassThroughUpdateHook) unwrapHook() any { return h.ResourceManager }
func (h migrateModePassThroughUpdateHook) providedHook() any {
	return (*MigrateModePassThroughUpdate)(nil)
}
func (migrateModePassThroughUpdateHook) MigrateModePassThroughUpdate() {}

func WithTransformRequestBody(m ResourceManager, fn func(params) params) ResourceManager {
	if m == nil || fn == nil {
		return m
	}
	return transformRequestBodyHook{ResourceManager: m, fn: fn}
}

type transformRequestBodyHook struct {
	ResourceManager
	fn func(params) params
}

func (h transformRequestBodyHook) unwrapHook() any   { return h.ResourceManager }
func (h transformRequestBodyHook) providedHook() any { return (*TransformRequestBody)(nil) }
func (h transformRequestBodyHook) TransformRequestBody(body params) params {
	return h.fn(body)
}

func WithTransformResponseRecord(m ResourceManager, fn func(Record) Record) ResourceManager {
	if m == nil || fn == nil {
		return m
	}
	return transformResponseRecordHook{ResourceManager: m, fn: fn}
}

type transformResponseRecordHook struct {
	ResourceManager
	fn func(Record) Record
}

func (h transformResponseRecordHook) unwrapHook() any   { return h.ResourceManager }
func (h transformResponseRecordHook) providedHook() any { return (*TransformResponseRecord)(nil) }
func (h transformResponseRecordHook) TransformResponseRecord(response Record) Record {
	return h.fn(response)
}

func WithNormalizeRecordForCreateAdopt(m ResourceManager, fn func(Record) Record) ResourceManager {
	if m == nil || fn == nil {
		return m
	}
	return normalizeRecordForCreateAdoptHook{ResourceManager: m, fn: fn}
}

type normalizeRecordForCreateAdoptHook struct {
	ResourceManager
	fn func(Record) Record
}

func (h normalizeRecordForCreateAdoptHook) unwrapHook() any { return h.ResourceManager }
func (h normalizeRecordForCreateAdoptHook) providedHook() any {
	return (*NormalizeRecordForCreateAdopt)(nil)
}
func (h normalizeRecordForCreateAdoptHook) NormalizeRecordForCreateAdopt(r Record) Record {
	return h.fn(r)
}

func WithResolveRecordAfterAsyncTask(m ResourceManager, fn func(context.Context, *VMSRest, Record) (Record, error)) ResourceManager {
	if m == nil || fn == nil {
		return m
	}
	return resolveRecordAfterAsyncTaskHook{ResourceManager: m, fn: fn}
}

type resolveRecordAfterAsyncTaskHook struct {
	ResourceManager
	fn func(context.Context, *VMSRest, Record) (Record, error)
}

func (h resolveRecordAfterAsyncTaskHook) unwrapHook() any { return h.ResourceManager }
func (h resolveRecordAfterAsyncTaskHook) providedHook() any {
	return (*ResolveRecordAfterAsyncTask)(nil)
}
func (h resolveRecordAfterAsyncTaskHook) ResolveRecordAfterAsyncTask(ctx context.Context, rest *VMSRest, record Record) (Record, error) {
	return h.fn(ctx, rest, record)
}
