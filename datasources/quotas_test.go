package datasources_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
	"github.com/vast-data/terraform-provider-vastdata/datasources"
	utils "github.com/vast-data/terraform-provider-vastdata/utils"

	"github.com/hashicorp/terraform-plugin-log/tflogtest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api_latest "github.com/vast-data/terraform-provider-vastdata/codegen/latest"
	vast_client "github.com/vast-data/terraform-provider-vastdata/vast-client"
)

var _ = Describe(" Quota", func() {
	var ReadContext schema.ReadContextFunc
	//An empty resource data to be populated per test
	var QuotaDataSourceData *schema.ResourceData
	var model_json = `
                         {
   "cluster": "string",
   "cluster_id": 100,
   "default_email": "string",
   "default_group_quota": {
      "grace_period": "string",
      "hard_limit": 100,
      "hard_limit_inodes": 100,
      "quota_system_id": 100,
      "sof_limit_inodes": 100,
      "soft_limit": 100
   },
   "default_user_quota": {
      "grace_period": "string",
      "hard_limit": 100,
      "hard_limit_inodes": 100,
      "quota_system_id": 100,
      "sof_limit_inodes": 100,
      "soft_limit": 100
   },
   "grace_period": "string",
   "group_quotas": [
      {
         "entity": {
            "email": "string",
            "identifier": "string",
            "identifier_type": "string",
            "name": "string",
            "vast_id": 100
         },
         "grace_period": "string",
         "hard_limit": 100,
         "hard_limit_inodes": 100,
         "quota_system_id": 100,
         "soft_limit": 100,
         "soft_limit_inodes": 100,
         "time_to_block": "string",
         "used_capacity": 100,
         "used_inodes": 100
      }
   ],
   "guid": "string",
   "hard_limit": 100,
   "hard_limit_inodes": 100,
   "id": 100,
   "name": "string",
   "num_blocked_users": 100,
   "num_exceeded_users": 100,
   "path": "string",
   "percent_capacity": 100,
   "percent_inodes": 100,
   "pretty_grace_period": "string",
   "pretty_state": "string",
   "soft_limit": 100,
   "soft_limit_inodes": 100,
   "state": "string",
   "system_id": 100,
   "tenant_id": 100,
   "tenant_name": "string",
   "time_to_block": "string",
   "used_capacity": 100,
   "used_capacity_tb": 10.5,
   "used_effective_capacity": 100,
   "used_effective_capacity_tb": 10.5,
   "used_inodes": 100,
   "user_quotas": [
      {
         "entity": {
            "email": "string",
            "identifier": "string",
            "identifier_type": "string",
            "name": "string",
            "vast_id": 100
         },
         "grace_period": "string",
         "hard_limit": 100,
         "hard_limit_inodes": 100,
         "quota_system_id": 100,
         "soft_limit": 100,
         "soft_limit_inodes": 100,
         "time_to_block": "string",
         "used_capacity": 100,
         "used_inodes": 100
      }
   ]
}
                         `
	var server *ghttp.Server
	var client vast_client.JwtSession
	QuotaDataSource := datasources.DataSourceQuota()
	ReadContext = QuotaDataSource.ReadContext

	BeforeEach(func() {
		QuotaDataSourceData = QuotaDataSource.TestResourceData()
		QuotaDataSourceData.SetId("100")
		server = ghttp.NewTLSServer()
		host_port := strings.Split(server.Addr(), ":")
		host := host_port[0]
		_port := host_port[1]
		port, _ := strconv.ParseUint(_port, 10, 64)
		client = vast_client.NewJwtSession(host, "user", "pwd", port, true)
		server.AppendHandlers(ghttp.CombineHandlers(
			ghttp.VerifyRequest("POST", "/api/token/"),
			ghttp.VerifyJSON("{\"username\":\"user\",\"password\":\"pwd\"}"),
			ghttp.RespondWith(200, `{"access":"femcew2d332f2e2e322e2qqw#2","":"32dm0932kde,ml;sd,s;l,322332"}`),
		))

	},
	)
	Describe("Validating Datasource Read", func() {
		Context("Read A datasource", func() {
			It("Datasource:Quota ,reads", func() {
				var ctx context.Context = context.Background()
				var output bytes.Buffer
				var guid = "11111-11111-11111-11111-11111"

				resource := api_latest.Quota{}
				json.Unmarshal([]byte(model_json), &resource)
				QuotaDataSourceData.SetId(guid)
				resource.Guid = guid
				b, err := json.Marshal(&resource)
				Expect(err).To(BeNil())
				values := url.Values{}

				values.Add("name", fmt.Sprintf("%v", resource.Name))
				QuotaDataSourceData.Set("name", resource.Name)

				request_url := `[{"url":"https://` + server.Addr() + `quotas100"}]`
				server.AppendHandlers(ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "quotas", values.Encode()),
					ghttp.RespondWith(200, request_url),
				),
				)

				server.AppendHandlers(ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "quotas100"),
					ghttp.RespondWith(200, string(b)),
				),
				)

				ctx = tflogtest.RootLogger(ctx, &output)
				e := client.Start()
				Expect(e).To(BeNil())
				d := ReadContext(ctx, QuotaDataSourceData, client)
				Expect(d).To(BeNil())

				o := new(map[string]interface{})
				json.Unmarshal(b, o)
				z := make(map[string]string)
				utils.MapToTFSchema(*o, &z, "")
				attributes := QuotaDataSourceData.State().Attributes
				for k, v := range z {
					Expect(attributes).To(HaveKeyWithValue(k, v))
				}

			})
		})
	},
	)
})
