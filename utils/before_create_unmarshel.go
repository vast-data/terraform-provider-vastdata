package utils

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type BeforeCreateUnmarshalFunc func(context.Context, []byte, *schema.ResourceData) ([]byte, error)

func VastDatabaseSchemaBeforeCreateUnmarshel(ctx context.Context, response []byte, d *schema.ResourceData) ([]byte, error) {
	var vastdatabse_schema map[string]interface{} = map[string]interface{}{}
	vastdatabse_schema["name"] = d.Get("name")
	vastdatabse_schema["database_name"] = d.Get("database_name")
	vastdatabse_schema["identifier"] = fmt.Sprintf("%v/%v", d.Get("name"), d.Get("database_name"))
	random_uuid, err := uuid.GenerateUUID()
	if err != nil {
		return nil, err
	}
	vastdatabse_schema["Id"] = random_uuid
	vastdatabse_schema["Guid"] = random_uuid
	return json.Marshal(vastdatabse_schema)
}

func VastDatabaseTableBeforeCreateUnmarshel(ctx context.Context, response []byte, d *schema.ResourceData) ([]byte, error) {
	var m map[string]interface{} = map[string]interface{}{}
	random_uuid, uuid_err := uuid.GenerateUUID()
	if uuid_err != nil {
		return nil, uuid_err
	}
	m["Id"] = random_uuid
	m["Guid"] = random_uuid
	tflog.Debug(ctx, fmt.Sprintf("About to convert to byte stream : %v", m))
	return json.Marshal(m)
}
