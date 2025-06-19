package vastdata

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	vast_client "github.com/vast-data/go-vast-client"
)

type VastResourceGetter interface {
	GetRestResource(*VMSRest) VastResourceAPIWithContext
}

type DataSourceState interface {
	VastResourceGetter
	PrepareReadDatasource(context.Context, *VMSRest, datasource.ReadRequest) error
	ReadDatasource(context.Context, *VMSRest, datasource.ReadRequest) (DisplayableRecord, error)
}

type ResourceState interface {
	VastResourceGetter

	//  Prepare methods are used to prepare the resource (rest params)
	PrepareImportResourceState(context.Context, *VMSRest, resource.ImportStateRequest) error
	PrepareCreateResource(context.Context, *VMSRest, resource.CreateRequest) error
	PrepareReadResource(context.Context, *VMSRest, resource.ReadRequest) error
	PrepareUpdateResource(context.Context, *VMSRest, resource.UpdateRequest) error
	PrepareDeleteResource(context.Context, *VMSRest, resource.DeleteRequest) error

	//  Resource methods are used to perform the actual resource operations
	ImportResourceState(context.Context, *VMSRest, resource.ImportStateRequest) (DisplayableRecord, error)
	CreateResource(context.Context, *VMSRest, resource.CreateRequest) (DisplayableRecord, error)
	ReadResource(context.Context, *VMSRest, resource.ReadRequest) (DisplayableRecord, error)
	UpdateResource(context.Context, *VMSRest, resource.UpdateRequest) (DisplayableRecord, error)
	DeleteResource(context.Context, *VMSRest, resource.DeleteRequest) (DisplayableRecord, error)
}

type (
	VastResourceAPI            = vast_client.VastResourceAPI
	VastResourceAPIWithContext = vast_client.VastResourceAPIWithContext
	DisplayableRecord          = vast_client.DisplayableRecord
	Renderable                 = vast_client.Renderable
	Record                     = vast_client.Record
	RecordSet                  = vast_client.RecordSet
	params                     = vast_client.Params
	VMSRest                    = vast_client.VMSRest
)

var (
	isNotFoundErr  = vast_client.IsNotFoundErr
	ignoreNotFound = vast_client.IgnoreNotFound
)
