package vastdata

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"math/rand"
)

func (rs *User) PrepareImportResourceState(ctx context.Context, rest *VMSRest, req resource.ImportStateRequest) error {
	return nil
}

func (rs *User) ImportResourceState(ctx context.Context, rest *VMSRest, req resource.ImportStateRequest) (DisplayableRecord, error) {
	return nil, nil
}

func (rs *User) PrepareCreateResource(ctx context.Context, rest *VMSRest, req resource.CreateRequest) error {
	if !rs.Uid.IsUnknown() {
		return nil
	}
	searchParams := params{"name": rs.Name.ValueString()}
	record, err := ignoreNotFound(rest.Users.GetWithContext(ctx, searchParams))
	if err != nil {
		return err
	}
	if record == nil {
		for attempts := 0; attempts < 20; attempts++ {
			uid := rand.Intn(4e4) + 2e4
			tflog.Info(ctx, fmt.Sprintf("Getting UID %d", uid))

			if rest.Users.MustExistsWithContext(ctx, params{"uid": uid}) {
				tflog.Info(ctx, fmt.Sprintf("User with this UID  %d already exists, trying another one", uid))
				continue
			}
			tflog.Info(ctx, fmt.Sprintf("Found free UID %d for new user", uid))
			rs.Uid = types.Int64Value(int64(uid))
			break
		}
	}
	return nil

}

func (rs *User) CreateResource(ctx context.Context, rest *VMSRest, req resource.CreateRequest) (DisplayableRecord, error) {
	return nil, nil
}

func (rs *User) PrepareReadResource(ctx context.Context, rest *VMSRest, req resource.ReadRequest) error {
	return nil
}

func (rs *User) ReadResource(ctx context.Context, rest *VMSRest, req resource.ReadRequest) (DisplayableRecord, error) {
	return nil, nil
}

func (rs *User) PrepareUpdateResource(ctx context.Context, rest *VMSRest, req resource.UpdateRequest) error {
	return nil
}

func (rs *User) UpdateResource(ctx context.Context, rest *VMSRest, req resource.UpdateRequest) (DisplayableRecord, error) {
	return nil, nil
}

func (rs *User) PrepareDeleteResource(ctx context.Context, rest *VMSRest, req resource.DeleteRequest) error {
	return nil
}

func (rs *User) DeleteResource(ctx context.Context, rest *VMSRest, req resource.DeleteRequest) (DisplayableRecord, error) {
	return nil, nil
}
