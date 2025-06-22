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

var _ = Describe(" NonLocalUser", func() {
	var ReadContext schema.ReadContextFunc
	//An empty resource data to be populated per test
	var NonLocalUserDataSourceData *schema.ResourceData
	var model_json = `
                         {
   "context": "string",
   "id": "string",
   "s3_policies_ids": [
      1,
      2,
      3,
      4,
      5,
      6
   ],
   "tenant_id": 100,
   "uid": 100,
   "username": "string"
}
                         `
	var server *ghttp.Server
	var client *vast_client.VMSSession
	NonLocalUserDataSource := datasources.DataSourceNonLocalUser()
	ReadContext = NonLocalUserDataSource.ReadContext

	BeforeEach(func() {
		NonLocalUserDataSourceData = NonLocalUserDataSource.TestResourceData()
		NonLocalUserDataSourceData.SetId("100")
		server = ghttp.NewTLSServer()
		host_port := strings.Split(server.Addr(), ":")
		host := host_port[0]
		_port := host_port[1]
		port, _ := strconv.ParseUint(_port, 10, 64)
		config := &vast_client.RestClientConfig{
			Host:      host,
			Port:      port,
			Username:  "user",
			Password:  "password",
			SslVerify: false,
		}
		client = vast_client.NewSession(context.TODO(), config)
		server.AppendHandlers(ghttp.CombineHandlers(
			ghttp.VerifyRequest("POST", "/api/token/"),
			ghttp.VerifyJSON("{\"username\":\"user\",\"password\":\"pwd\"}"),
			ghttp.RespondWith(200, `{"access":"femcew2d332f2e2e322e2qqw#2","":"32dm0932kde,ml;sd,s;l,322332"}`),
		))

	},
	)
	Describe("Validating Datasource Read", func() {
		Context("Read A datasource", func() {
			It("Datasource:NonLocalUser ,reads", func() {
				var ctx context.Context = context.Background()
				var output bytes.Buffer
				var guid = "11111-11111-11111-11111-11111"

				resource := api_latest.NonLocalUser{}
				json.Unmarshal([]byte(model_json), &resource)
				NonLocalUserDataSourceData.SetId(guid)
				resource.Guid = guid
				b, err := json.Marshal(&resource)
				Expect(err).To(BeNil())
				values := url.Values{}

				values.Add("tenant_id", fmt.Sprintf("%v", resource.TenantId))
				NonLocalUserDataSourceData.Set("tenant_id", resource.TenantId)

				values.Add("username", fmt.Sprintf("%v", resource.Username))
				NonLocalUserDataSourceData.Set("username", resource.Username)

				values.Add("context", fmt.Sprintf("%v", resource.Context))
				NonLocalUserDataSourceData.Set("context", resource.Context)

				server.AppendHandlers(ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "users/query", values.Encode()), //since this is a test http server and will not return id upon POST (creation) so json will use the zero value
					ghttp.RespondWith(200, `[`+string(b)+`]`),
				),
				)

				ctx = tflogtest.RootLogger(ctx, &output)
				e := client.Start()
				Expect(e).To(BeNil())
				d := ReadContext(ctx, NonLocalUserDataSourceData, client)
				Expect(d).To(BeNil())

				o := new(map[string]interface{})
				json.Unmarshal(b, o)
				z := make(map[string]string)
				utils.MapToTFSchema(*o, &z, "")
				attributes := NonLocalUserDataSourceData.State().Attributes
				for k, v := range z {
					Expect(attributes).To(HaveKeyWithValue(k, v))
				}

			})
		})
	},
	)
})
