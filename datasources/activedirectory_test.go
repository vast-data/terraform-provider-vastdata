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

var _ = Describe(" ActiveDirectory", func() {
	var ReadContext schema.ReadContextFunc
	//An empty resource data to be populated per test
	var ActiveDirectoryDataSourceData *schema.ResourceData
	var model_json = `
                         {
   "domain_name": "string",
   "guid": "string",
   "id": 100,
   "ldap_id": "string",
   "machine_account_name": "string",
   "match_user": "string",
   "method": "string",
   "organizational_unit": "string",
   "port": 100,
   "posix_account": "string",
   "posix_group": "string",
   "preferred_dc_list": [
      "A",
      "B",
      "C",
      "D",
      "E"
   ],
   "query_groups_mode": "string",
   "state": "string",
   "uid": "string",
   "uid_member": "string",
   "uid_member_value_property_name": "string",
   "user_login_name": "string",
   "username_property_name": "string"
}
                         `
	var server *ghttp.Server
	var client vast_client.JwtSession
	ActiveDirectoryDataSource := datasources.DataSourceActiveDirectory()
	ReadContext = ActiveDirectoryDataSource.ReadContext

	BeforeEach(func() {
		ActiveDirectoryDataSourceData = ActiveDirectoryDataSource.TestResourceData()
		ActiveDirectoryDataSourceData.SetId("100")
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
			It("Datasource:ActiveDirectory ,reads", func() {
				var ctx context.Context = context.Background()
				var output bytes.Buffer
				var guid = "11111-11111-11111-11111-11111"

				resource := api_latest.ActiveDirectory{}
				json.Unmarshal([]byte(model_json), &resource)
				ActiveDirectoryDataSourceData.SetId(guid)
				resource.Guid = guid
				b, err := json.Marshal(&resource)
				Expect(err).To(BeNil())
				values := url.Values{}

				values.Add("machine_account_name", fmt.Sprintf("%v", resource.MachineAccountName))
				ActiveDirectoryDataSourceData.Set("machine_account_name", resource.MachineAccountName)

				server.AppendHandlers(ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/api/activedirectory/", values.Encode()), //since this is a test http server and will not return id upon POST (creation) so json will use the zero value
					ghttp.RespondWith(200, `[`+string(b)+`]`),
				),
				)

				ctx = tflogtest.RootLogger(ctx, &output)
				e := client.Start()
				Expect(e).To(BeNil())
				d := ReadContext(ctx, ActiveDirectoryDataSourceData, client)
				Expect(d).To(BeNil())

				o := new(map[string]interface{})
				json.Unmarshal(b, o)
				z := make(map[string]string)
				utils.MapToTFSchema(*o, &z, "")
				attributes := ActiveDirectoryDataSourceData.State().Attributes
				for k, v := range z {
					Expect(attributes).To(HaveKeyWithValue(k, v))
				}

			})
		})
	},
	)
})
