// // Copyright (c) HashiCorp, Inc.
package provider

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
)

var S3LifeCycleRuleSchemaRef = is.NewSchemaReference(
	http.MethodPost,
	"s3lifecyclerules",
	http.MethodGet,
	"s3lifecyclerules",
)

type S3LifeCycleRule struct {
	tfstate *is.TFState
}

func (m *S3LifeCycleRule) NewResourceManager(raw map[string]attr.Value, schema any) ResourceManager {
	return &S3LifeCycleRule{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef:          S3LifeCycleRuleSchemaRef,
			SearchFilterFields: []string{"view_id"},
		},
	)}
}

func (m *S3LifeCycleRule) NewDatasourceManager(raw map[string]attr.Value, schema any) DataSourceManager {
	return &S3LifeCycleRule{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef:          S3LifeCycleRuleSchemaRef,
			SearchFilterFields: []string{"view_id"},
		}),
	}
}

func (m *S3LifeCycleRule) TfState() *is.TFState {
	return m.tfstate
}

func (m *S3LifeCycleRule) API(rest *VMSRest) VastResourceAPIWithContext {
	return rest.S3LifeCycleRules
}

func (m *S3LifeCycleRule) ReadDatasource(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	return m.readByViewID(ctx, rest)
}

func (m *S3LifeCycleRule) ReadResource(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	return m.readByViewID(ctx, rest)
}

// readByViewID performs the actual lookup, translating view_id → view__id so that
// the VAST backend correctly filters s3lifecyclerules by their owning view.
func (m *S3LifeCycleRule) readByViewID(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	searchParams := getSearchParams(ctx, m.tfstate, nil)
	rawParams := params{}
	for k, v := range searchParams {
		if k == "view_id" {
			rawParams["view__id"] = v
		} else {
			rawParams[k] = v
		}
	}
	return rest.S3LifeCycleRules.GetWithContext(ctx, rawParams)
}

// LookupForCreate supplies a custom pre-flight GET for the create path.
// The default lookup queries with "view_id" (single underscore), which the VAST
// API ignores, so rules with the same name but different view_ids would be
// incorrectly merged. Using readByViewID ensures "view__id" (double underscore)
// is sent, correctly scoping the lookup to the target view.
func (m *S3LifeCycleRule) LookupForCreate(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	return m.readByViewID(ctx, rest)
}
