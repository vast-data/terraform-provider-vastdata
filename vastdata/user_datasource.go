package vastdata

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func (rs *User) PrepareReadDatasource(ctx context.Context, rest *VMSRest, req datasource.ReadRequest) error {
	return nil
}

func (rs *User) ReadDatasource(ctx context.Context, rest *VMSRest, req datasource.ReadRequest) (DisplayableRecord, error) {
	return nil, nil
}
