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

var _ = Describe(" QosPolicy", func() {
	var ReadContext schema.ReadContextFunc
	//An empty resource data to be populated per test
	var QosPolicyDataSourceData *schema.ResourceData
	var model_json = `
                         {
   "capacity_limits": {
      "max_reads_bw_mbps_per_gb_capacity": 100,
      "max_reads_iops_per_gb_capacity": 100,
      "max_writes_bw_mbps_per_gb_capacity": 100,
      "max_writes_iops_per_gb_capacity": 100
   },
   "guid": "string",
   "id": 100,
   "io_size_bytes": 100,
   "mode": "string",
   "name": "string",
   "static_limits": {
      "max_reads_bw_mbps": 100,
      "max_reads_iops": 100,
      "max_writes_bw_mbps": 100,
      "max_writes_iops": 100,
      "min_reads_bw_mbps": 100,
      "min_reads_iops": 100,
      "min_writes_bw_mbps": 100,
      "min_writes_iops": 100
   }
}
                         `
	var server *ghttp.Server
	var client vast_client.JwtSession
	QosPolicyDataSource := datasources.DataSourceQosPolicy()
	ReadContext = QosPolicyDataSource.ReadContext

	BeforeEach(func() {
		QosPolicyDataSourceData = QosPolicyDataSource.TestResourceData()
		QosPolicyDataSourceData.SetId("100")
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
			It("Datasource:QosPolicy ,reads", func() {
				var ctx context.Context = context.Background()
				var output bytes.Buffer
				var guid = "11111-11111-11111-11111-11111"

				resource := api_latest.QosPolicy{}
				json.Unmarshal([]byte(model_json), &resource)
				QosPolicyDataSourceData.SetId(guid)
				resource.Guid = guid
				b, err := json.Marshal(&resource)
				Expect(err).To(BeNil())
				values := url.Values{}

				values.Add("name", fmt.Sprintf("%v", resource.Name))
				QosPolicyDataSourceData.Set("name", resource.Name)

				server.AppendHandlers(ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/api/qospolicies/", values.Encode()), //since this is a test http server and will not return id upon POST (creation) so json will use the zero value
					ghttp.RespondWith(200, `[`+string(b)+`]`),
				),
				)

				ctx = tflogtest.RootLogger(ctx, &output)
				e := client.Start()
				Expect(e).To(BeNil())
				d := ReadContext(ctx, QosPolicyDataSourceData, client)
				Expect(d).To(BeNil())

				o := new(map[string]interface{})
				json.Unmarshal(b, o)
				z := make(map[string]string)
				utils.MapToTFSchema(*o, &z, "")
				attributes := QosPolicyDataSourceData.State().Attributes
				for k, v := range z {
					Expect(attributes).To(HaveKeyWithValue(k, v))
				}

			})
		})
	},
	)
})
