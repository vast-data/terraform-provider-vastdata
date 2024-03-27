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

var _ = Describe(" S3LifeCycleRule", func() {
	var ReadContext schema.ReadContextFunc
	//An empty resource data to be populated per test
	var S3LifeCycleRuleDataSourceData *schema.ResourceData
	var model_json = `
                         {
   "abort_mpu_days_after_initiation": 100,
   "expiration_date": "string",
   "expiration_days": 100,
   "guid": "string",
   "id": 100,
   "max_size": 100,
   "min_size": 100,
   "name": "string",
   "newer_noncurrent_versions": 100,
   "noncurrent_days": 100,
   "prefix": "string",
   "view_id": 100,
   "view_path": "string"
}
                         `
	var server *ghttp.Server
	var client vast_client.JwtSession
	S3LifeCycleRuleDataSource := datasources.DataSourceS3LifeCycleRule()
	ReadContext = S3LifeCycleRuleDataSource.ReadContext

	BeforeEach(func() {
		S3LifeCycleRuleDataSourceData = S3LifeCycleRuleDataSource.TestResourceData()
		S3LifeCycleRuleDataSourceData.SetId("100")
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
			It("Datasource:S3LifeCycleRule ,reads", func() {
				var ctx context.Context = context.Background()
				var output bytes.Buffer
				var guid = "11111-11111-11111-11111-11111"

				resource := api_latest.S3LifeCycleRule{}
				json.Unmarshal([]byte(model_json), &resource)
				S3LifeCycleRuleDataSourceData.SetId(guid)
				resource.Guid = guid
				b, err := json.Marshal(&resource)
				Expect(err).To(BeNil())
				values := url.Values{}

				values.Add("name", fmt.Sprintf("%v", resource.Name))
				S3LifeCycleRuleDataSourceData.Set("name", resource.Name)

				server.AppendHandlers(ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "s3lifecyclerules", values.Encode()), //since this is a test http server and will not return id upon POST (creation) so json will use the zero value
					ghttp.RespondWith(200, `[`+string(b)+`]`),
				),
				)

				ctx = tflogtest.RootLogger(ctx, &output)
				e := client.Start()
				Expect(e).To(BeNil())
				d := ReadContext(ctx, S3LifeCycleRuleDataSourceData, client)
				Expect(d).To(BeNil())

				o := new(map[string]interface{})
				json.Unmarshal(b, o)
				z := make(map[string]string)
				utils.MapToTFSchema(*o, &z, "")
				attributes := S3LifeCycleRuleDataSourceData.State().Attributes
				for k, v := range z {
					Expect(attributes).To(HaveKeyWithValue(k, v))
				}

			})
		})
	},
	)
})
