package vastdata

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func (rs *Group) PrepareImportResourceState(ctx context.Context, rest *VMSRest, req resource.ImportStateRequest) error {
	return nil
}

func (rs *Group) ImportResourceState(ctx context.Context, rest *VMSRest, req resource.ImportStateRequest) (DisplayableRecord, error) {
	return nil, nil
}

func (rs *Group) PrepareCreateResource(ctx context.Context, rest *VMSRest, req resource.CreateRequest) error {
	return nil

}

func (rs *Group) CreateResource(ctx context.Context, rest *VMSRest, req resource.CreateRequest) (DisplayableRecord, error) {
	return nil, nil
}

func (rs *Group) PrepareReadResource(ctx context.Context, rest *VMSRest, req resource.ReadRequest) error {
	return nil
}

func (rs *Group) ReadResource(ctx context.Context, rest *VMSRest, req resource.ReadRequest) (DisplayableRecord, error) {
	return nil, nil
}

func (rs *Group) PrepareUpdateResource(ctx context.Context, rest *VMSRest, req resource.UpdateRequest) error {
	return nil
}

func (rs *Group) UpdateResource(ctx context.Context, rest *VMSRest, req resource.UpdateRequest) (DisplayableRecord, error) {
	return nil, nil
}

func (rs *Group) PrepareDeleteResource(ctx context.Context, rest *VMSRest, req resource.DeleteRequest) error {
	return nil
}

func (rs *Group) DeleteResource(ctx context.Context, rest *VMSRest, req resource.DeleteRequest) (DisplayableRecord, error) {
	return nil, nil
}
