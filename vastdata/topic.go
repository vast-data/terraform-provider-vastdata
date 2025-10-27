// // Copyright (c) HashiCorp, Inc.
package provider

import (
	"context"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	is "github.com/vast-data/terraform-provider-vastdata/vastdata/internalstate"
)

var TopicSchemaRef = is.NewSchemaReference(
	http.MethodPost,
	"topics",
	http.MethodGet,
	"topics/show",
)

type Topic struct {
	tfstate *is.TFState
}

func (m *Topic) NewResourceManager(raw map[string]attr.Value, schema any) ResourceManager {
	return &Topic{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef:      TopicSchemaRef,
			ReadOnlyFields: []string{"num_events", "tenant_id"},
		},
	)}
}

func (m *Topic) NewDatasourceManager(raw map[string]attr.Value, schema any) DataSourceManager {
	return &Topic{tfstate: is.NewTFStateMust(
		raw,
		schema,
		&is.TFStateHints{
			SchemaRef: TopicSchemaRef,
		},
	)}
}

func (m *Topic) TfState() *is.TFState {
	return m.tfstate
}

func (m *Topic) API(rest *VMSRest) VastResourceAPIWithContext {
	return rest.Topics
}

func (m *Topic) ReadDatasource(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	searchParams, _ := m.tfstate.SetIfAvailable("database_name", "name")
	record, err := rest.Topics.TopicShowWithContext_GET(ctx, searchParams)
	if isApiError(err) {
		if strings.Contains(err.(*ApiError).Body, "UNKNOWN_TOPIC_OR_PARTITION") {
			return nil, nil
		}
	}
	return record, err
}

func (m *Topic) ReadResource(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	return m.ReadDatasource(ctx, rest)
}

func (m *Topic) CreateResource(ctx context.Context, rest *VMSRest) (DisplayableRecord, error) {
	ts := m.tfstate
	record, err := m.ReadDatasource(ctx, rest)
	if err != nil {
		return nil, err
	}
	if record != nil {
		return record, nil
	}
	createParams := ts.GetCreateParams()
	if _, err := rest.Topics.TopicTopicsWithContext_POST(ctx, createParams); err != nil {
		return nil, err
	}
	// read topics after creation to get
	// updated "message_timestamp_after_max_ms" and "message_timestamp_before_max_ms"
	return m.ReadDatasource(ctx, rest)
}

func (m *Topic) UpdateResource(ctx context.Context, plan UpdateResource, rest *VMSRest) (DisplayableRecord, error) {
	planTs := plan.(*Topic).TfState()
	updateParams := planTs.GetCreateParams()
	if _, err := rest.Topics.TopicTopicsWithContext_PATCH(ctx, updateParams); err != nil {
		return nil, err
	}
	return m.ReadDatasource(ctx, rest)
}

func (m *Topic) DeleteResource(ctx context.Context, rest *VMSRest) error {
	if _, err := m.ReadDatasource(ctx, rest); err != nil {
		return err
	}
	deleteParams, _ := m.tfstate.SetIfAvailable("database_name", "name")
	err := rest.Topics.TopicDeleteWithContext_DELETE(ctx, deleteParams)
	return ignoreStatusCodes(err, http.StatusNotFound)
}
