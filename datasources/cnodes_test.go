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

var _ = Describe(" Cnode", func() {
	var ReadContext schema.ReadContextFunc
	//An empty resource data to be populated per test
	var CnodeDataSourceData *schema.ResourceData
	var model_json = `
                         {
   "bmc_fw_version": "string",
   "cbox": "string",
   "cbox_id": 100,
   "cbox_uid": "string",
   "cluster": "string",
   "data_rdma_port": 100,
   "data_tcp_port": 100,
   "display_state": "string",
   "guid": "string",
   "host_label": "string",
   "hostname": "string",
   "id": 100,
   "ip": "string",
   "ip1": "string",
   "ip2": "string",
   "ipv6": "string",
   "led_status": "string",
   "mgmt_ip": "string",
   "name": "string",
   "new_name": "string",
   "os_version": "string",
   "platform_rdma_port": 100,
   "platform_tcp_port": 100,
   "sn": "string",
   "state": "string",
   "url": "string"
}
                         `
	var server *ghttp.Server
	var client vast_client.JwtSession
	CnodeDataSource := datasources.DataSourceCnode()
	ReadContext = CnodeDataSource.ReadContext

	BeforeEach(func() {
		CnodeDataSourceData = CnodeDataSource.TestResourceData()
		CnodeDataSourceData.SetId("100")
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
			It("Datasource:Cnode ,reads", func() {
				var ctx context.Context = context.Background()
				var output bytes.Buffer
				var guid = "11111-11111-11111-11111-11111"

				resource := api_latest.Cnode{}
				json.Unmarshal([]byte(model_json), &resource)
				CnodeDataSourceData.SetId(guid)
				resource.Guid = guid
				b, err := json.Marshal(&resource)
				Expect(err).To(BeNil())
				values := url.Values{}

				values.Add("name", fmt.Sprintf("%v", resource.Name))
				CnodeDataSourceData.Set("name", resource.Name)

				server.AppendHandlers(ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/api/cnodes/", values.Encode()), //since this is a test http server and will not return id upon POST (creation) so json will use the zero value
					ghttp.RespondWith(200, `[`+string(b)+`]`),
				),
				)

				ctx = tflogtest.RootLogger(ctx, &output)
				e := client.Start()
				Expect(e).To(BeNil())
				d := ReadContext(ctx, CnodeDataSourceData, client)
				Expect(d).To(BeNil())

				o := new(map[string]interface{})
				json.Unmarshal(b, o)
				z := make(map[string]string)
				utils.MapToTFSchema(*o, &z, "")
				attributes := CnodeDataSourceData.State().Attributes
				for k, v := range z {
					Expect(attributes).To(HaveKeyWithValue(k, v))
				}

			})
		})
	},
	)
})
