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

var _ = Describe(" Ldap", func() {
	var ReadContext schema.ReadContextFunc
	//An empty resource data to be populated per test
	var LdapDataSourceData *schema.ResourceData
	var model_json = `
                         {
   "active_directory": "string",
   "binddn": "string",
   "bindpw": "string",
   "domain_name": "string",
   "gid_number": "string",
   "group_login_name": "string",
   "group_searchbase": "string",
   "guid": "string",
   "id": 100,
   "mail_property_name": "string",
   "match_user": "string",
   "method": "string",
   "port": 100,
   "posix_account": "string",
   "posix_attributes_source": "string",
   "posix_group": "string",
   "query_groups_mode": "string",
   "searchbase": "string",
   "tls_certificate": "string",
   "uid": "string",
   "uid_member": "string",
   "uid_member_value_property_name": "string",
   "uid_number": "string",
   "urls": [
      "A",
      "B",
      "C",
      "D",
      "E"
   ],
   "user_login_name": "string",
   "username_property_name": "string"
}
                         `
	var server *ghttp.Server
	var client vast_client.JwtSession
	LdapDataSource := datasources.DataSourceLdap()
	ReadContext = LdapDataSource.ReadContext

	BeforeEach(func() {
		LdapDataSourceData = LdapDataSource.TestResourceData()
		LdapDataSourceData.SetId("100")
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
			It("Datasource:Ldap ,reads", func() {
				var ctx context.Context = context.Background()
				var output bytes.Buffer
				var guid = "11111-11111-11111-11111-11111"

				resource := api_latest.Ldap{}
				json.Unmarshal([]byte(model_json), &resource)
				LdapDataSourceData.SetId(guid)
				resource.Guid = guid
				b, err := json.Marshal(&resource)
				Expect(err).To(BeNil())
				values := url.Values{}

				values.Add("domain_name", fmt.Sprintf("%v", resource.DomainName))
				LdapDataSourceData.Set("domain_name", resource.DomainName)

				server.AppendHandlers(ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "ldaps", values.Encode()), //since this is a test http server and will not return id upon POST (creation) so json will use the zero value
					ghttp.RespondWith(200, `[`+string(b)+`]`),
				),
				)

				ctx = tflogtest.RootLogger(ctx, &output)
				e := client.Start()
				Expect(e).To(BeNil())
				d := ReadContext(ctx, LdapDataSourceData, client)
				Expect(d).To(BeNil())

				o := new(map[string]interface{})
				json.Unmarshal(b, o)
				z := make(map[string]string)
				utils.MapToTFSchema(*o, &z, "")
				attributes := LdapDataSourceData.State().Attributes
				for k, v := range z {
					Expect(attributes).To(HaveKeyWithValue(k, v))
				}

			})
		})
	},
	)
})
