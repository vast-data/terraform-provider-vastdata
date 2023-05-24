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
	"github.com/vast-data/terraform-provider-vastdata.git/datasources"
	utils "github.com/vast-data/terraform-provider-vastdata.git/utils"

	"github.com/hashicorp/terraform-plugin-log/tflogtest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api_latest "github.com/vast-data/terraform-provider-vastdata.git/codegen/latest"
	vast_client "github.com/vast-data/terraform-provider-vastdata.git/vast-client"
)

var _ = Describe(" User", func() {
	var ReadContext schema.ReadContextFunc
	//An empty resource data to be populated per test
	var UserDataSourceData *schema.ResourceData
	var model_json = `
                         {
   "access_keys": [
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
   "gids": [
      1,
      2,
      3,
      4,
      5,
      6
   ],
   "group_count": 100,
   "groups": [
      "A",
      "B",
      "C",
      "D",
      "E"
   ],
   "guid": "string",
   "id": 100,
   "leading_gid": 100,
   "leading_group_gid": 100,
   "leading_group_name": "string",
   "name": "string",
   "primary_group_sid": "string",
   "s3_policies_ids": [
      1,
      2,
      3,
      4,
      5,
      6
   ],
   "sid": "string",
   "sids": [
      "A",
      "B",
      "C",
      "D",
      "E"
   ],
   "uid": 100
}
                         `
	var server *ghttp.Server
	var client vast_client.JwtSession
	UserDataSource := datasources.DataSourceUser()
	ReadContext = UserDataSource.ReadContext

	BeforeEach(func() {
		UserDataSourceData = UserDataSource.TestResourceData()
		UserDataSourceData.SetId("100")
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
			It("Datasource:User ,reads", func() {
				var ctx context.Context = context.Background()
				var output bytes.Buffer
				var guid = "11111-11111-11111-11111-11111"

				resource := api_latest.User{}
				json.Unmarshal([]byte(model_json), &resource)
				UserDataSourceData.SetId(guid)
				resource.Guid = guid
				b, err := json.Marshal(&resource)
				Expect(err).To(BeNil())
				values := url.Values{}

				values.Add("name", fmt.Sprintf("%v", resource.Name))
				UserDataSourceData.Set("name", resource.Name)

				server.AppendHandlers(ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/api/users/", values.Encode()), //since this is a test http server and will not return id upon POST (creation) so json will use the zero value
					ghttp.RespondWith(200, `[`+string(b)+`]`),
				),
				)

				ctx = tflogtest.RootLogger(ctx, &output)
				e := client.Start()
				Expect(e).To(BeNil())
				d := ReadContext(ctx, UserDataSourceData, client)
				Expect(d).To(BeNil())

				o := new(map[string]interface{})
				json.Unmarshal(b, o)
				z := make(map[string]string)
				utils.MapToTFSchema(*o, &z, "")
				attributes := UserDataSourceData.State().Attributes
				for k, v := range z {
					Expect(attributes).To(HaveKeyWithValue(k, v))
				}

			})
		})
	},
	)
})
