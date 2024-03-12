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

var _ = Describe(" GlobalSnapshot", func() {
	var ReadContext schema.ReadContextFunc
	//An empty resource data to be populated per test
	var GlobalSnapshotDataSourceData *schema.ResourceData
	var model_json = `
                         {
   "guid": "string",
   "id": 100,
   "loanee_root_path": "string",
   "loanee_tenant_id": 100,
   "name": "string",
   "owner_root_snapshot": {},
   "owner_tenant": {},
   "remote_target_guid": "string",
   "remote_target_id": 100,
   "remote_target_path": "string"
}
                         `
	var server *ghttp.Server
	var client vast_client.JwtSession
	GlobalSnapshotDataSource := datasources.DataSourceGlobalSnapshot()
	ReadContext = GlobalSnapshotDataSource.ReadContext

	BeforeEach(func() {
		GlobalSnapshotDataSourceData = GlobalSnapshotDataSource.TestResourceData()
		GlobalSnapshotDataSourceData.SetId("100")
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
			It("Datasource:GlobalSnapshot ,reads", func() {
				var ctx context.Context = context.Background()
				var output bytes.Buffer
				var guid = "11111-11111-11111-11111-11111"

				resource := api_latest.GlobalSnapshot{}
				json.Unmarshal([]byte(model_json), &resource)
				GlobalSnapshotDataSourceData.SetId(guid)
				resource.Guid = guid
				b, err := json.Marshal(&resource)
				Expect(err).To(BeNil())
				values := url.Values{}

				values.Add("name", fmt.Sprintf("%v", resource.Name))
				GlobalSnapshotDataSourceData.Set("name", resource.Name)

				server.AppendHandlers(ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "globalsnapstreams", values.Encode()), //since this is a test http server and will not return id upon POST (creation) so json will use the zero value
					ghttp.RespondWith(200, `[`+string(b)+`]`),
				),
				)

				ctx = tflogtest.RootLogger(ctx, &output)
				e := client.Start()
				Expect(e).To(BeNil())
				d := ReadContext(ctx, GlobalSnapshotDataSourceData, client)
				Expect(d).To(BeNil())

				o := new(map[string]interface{})
				json.Unmarshal(b, o)
				z := make(map[string]string)
				utils.MapToTFSchema(*o, &z, "")
				attributes := GlobalSnapshotDataSourceData.State().Attributes
				for k, v := range z {
					Expect(attributes).To(HaveKeyWithValue(k, v))
				}

			})
		})
	},
	)
})
