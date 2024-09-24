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

var _ = Describe(" View", func() {
	var ReadContext schema.ReadContextFunc
	//An empty resource data to be populated per test
	var ViewDataSourceData *schema.ResourceData
	var model_json = `
                         {
   "abac_tags": [
      "A",
      "B",
      "C",
      "D",
      "E"
   ],
   "abe_max_depth": 100,
   "abe_protocols": [
      "A",
      "B",
      "C",
      "D",
      "E"
   ],
   "alias": "string",
   "auto_commit": "string",
   "bucket": "string",
   "bucket_creators": [
      "A",
      "B",
      "C",
      "D",
      "E"
   ],
   "bucket_creators_groups": [
      "A",
      "B",
      "C",
      "D",
      "E"
   ],
   "bucket_logging": {},
   "bucket_owner": "string",
   "cluster": "string",
   "cluster_id": 100,
   "default_retention_period": "string",
   "files_retention_mode": "string",
   "guid": "string",
   "id": 100,
   "logical_capacity": 100,
   "max_retention_period": "string",
   "min_retention_period": "string",
   "name": "string",
   "nfs_interop_flags": "string",
   "path": "string",
   "physical_capacity": 100,
   "policy_id": 100,
   "protocols": [
      "A",
      "B",
      "C",
      "D",
      "E"
   ],
   "qos_policy_id": 100,
   "s3_locks_retention_mode": "string",
   "s3_locks_retention_period": "string",
   "s3_object_ownership_rule": "string",
   "share": "string",
   "share_acl": {
      "acl": [
         {
            "fqdn": "string",
            "grantee": "string",
            "name": "string",
            "permissions": "string",
            "sid_str": "string",
            "uid_or_gid": 100
         }
      ]
   },
   "tenant_id": 100
}
                         `
	var server *ghttp.Server
	var client vast_client.JwtSession
	ViewDataSource := datasources.DataSourceView()
	ReadContext = ViewDataSource.ReadContext

	BeforeEach(func() {
		ViewDataSourceData = ViewDataSource.TestResourceData()
		ViewDataSourceData.SetId("100")
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
			It("Datasource:View ,reads", func() {
				var ctx context.Context = context.Background()
				var output bytes.Buffer
				var guid = "11111-11111-11111-11111-11111"

				resource := api_latest.View{}
				json.Unmarshal([]byte(model_json), &resource)
				ViewDataSourceData.SetId(guid)
				resource.Guid = guid
				b, err := json.Marshal(&resource)
				Expect(err).To(BeNil())
				values := url.Values{}

				values.Add("path", fmt.Sprintf("%v", resource.Path))
				ViewDataSourceData.Set("path", resource.Path)

				server.AppendHandlers(ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "views", values.Encode()), //since this is a test http server and will not return id upon POST (creation) so json will use the zero value
					ghttp.RespondWith(200, `[`+string(b)+`]`),
				),
				)

				ctx = tflogtest.RootLogger(ctx, &output)
				e := client.Start()
				Expect(e).To(BeNil())
				d := ReadContext(ctx, ViewDataSourceData, client)
				Expect(d).To(BeNil())

				o := new(map[string]interface{})
				json.Unmarshal(b, o)
				z := make(map[string]string)
				utils.MapToTFSchema(*o, &z, "")
				attributes := ViewDataSourceData.State().Attributes
				for k, v := range z {
					Expect(attributes).To(HaveKeyWithValue(k, v))
				}

			})
		})
	},
	)
})
