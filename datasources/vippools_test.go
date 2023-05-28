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

var _ = Describe(" VipPool", func() {
	var ReadContext schema.ReadContextFunc
	//An empty resource data to be populated per test
	var VipPoolDataSourceData *schema.ResourceData
	var model_json = `
                         {
   "active_interfaces": 100,
   "cluster": "string",
   "cnode_ids": [
      1,
      2,
      3,
      4,
      5,
      6
   ],
   "domain_name": "string",
   "guid": "string",
   "gw_ip": "string",
   "gw_ipv6": "string",
   "id": 100,
   "ip_ranges": [
      [
         "string-0",
         "string-1"
      ],
      [
         "string-0",
         "string-1"
      ],
      [
         "string-0",
         "string-1"
      ],
      [
         "string-0",
         "string-1"
      ],
      [
         "string-0",
         "string-1"
      ],
      [
         "string-0",
         "string-1"
      ],
      [
         "string-0",
         "string-1"
      ],
      [
         "string-0",
         "string-1"
      ],
      [
         "string-0",
         "string-1"
      ],
      [
         "string-0",
         "string-1"
      ]
   ],
   "name": "string",
   "peer_asn": 100,
   "port_membership": "string",
   "role": "string",
   "state": "string",
   "subnet_cidr": 100,
   "subnet_cidr_ipv6": 100,
   "sync": "string",
   "sync_time": "string",
   "url": "string",
   "vast_asn": 100,
   "vlan": 100
}
                         `
	var server *ghttp.Server
	var client vast_client.JwtSession
	VipPoolDataSource := datasources.DataSourceVipPool()
	ReadContext = VipPoolDataSource.ReadContext

	BeforeEach(func() {
		VipPoolDataSourceData = VipPoolDataSource.TestResourceData()
		VipPoolDataSourceData.SetId("100")
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
			It("Datasource:VipPool ,reads", func() {
				var ctx context.Context = context.Background()
				var output bytes.Buffer
				var guid = "11111-11111-11111-11111-11111"

				resource := api_latest.VipPool{}
				json.Unmarshal([]byte(model_json), &resource)
				VipPoolDataSourceData.SetId(guid)
				resource.Guid = guid
				b, err := json.Marshal(&resource)
				Expect(err).To(BeNil())
				values := url.Values{}

				values.Add("name", fmt.Sprintf("%v", resource.Name))
				VipPoolDataSourceData.Set("name", resource.Name)

				server.AppendHandlers(ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/api/vippools/", values.Encode()), //since this is a test http server and will not return id upon POST (creation) so json will use the zero value
					ghttp.RespondWith(200, `[`+string(b)+`]`),
				),
				)

				ctx = tflogtest.RootLogger(ctx, &output)
				e := client.Start()
				Expect(e).To(BeNil())
				d := ReadContext(ctx, VipPoolDataSourceData, client)
				Expect(d).To(BeNil())

				o := new(map[string]interface{})
				json.Unmarshal(b, o)
				z := make(map[string]string)
				utils.MapToTFSchema(*o, &z, "")
				attributes := VipPoolDataSourceData.State().Attributes
				for k, v := range z {
					Expect(attributes).To(HaveKeyWithValue(k, v))
				}

			})
		})
	},
	)
})
