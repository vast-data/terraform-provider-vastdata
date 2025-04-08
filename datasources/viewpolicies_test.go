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

var _ = Describe(" ViewPolicy", func() {
	var ReadContext schema.ReadContextFunc
	//An empty resource data to be populated per test
	var ViewPolicyDataSourceData *schema.ResourceData
	var model_json = `
                         {
   "access_flavor": "string",
   "allowed_characters": "string",
   "atime_frequency": "string",
   "auth_source": "string",
   "cluster": "string",
   "cluster_id": 100,
   "count_views": 100,
   "flavor": "string",
   "gid_inheritance": "string",
   "guid": "string",
   "id": 100,
   "name": "string",
   "nfs_all_squash": [
      "A",
      "B",
      "C",
      "D",
      "E"
   ],
   "nfs_minimal_protection_level": "string",
   "nfs_no_squash": [
      "A",
      "B",
      "C",
      "D",
      "E"
   ],
   "nfs_read_only": [
      "A",
      "B",
      "C",
      "D",
      "E"
   ],
   "nfs_read_write": [
      "A",
      "B",
      "C",
      "D",
      "E"
   ],
   "nfs_root_squash": [
      "A",
      "B",
      "C",
      "D",
      "E"
   ],
   "path_length": "string",
   "protocols": [
      "A",
      "B",
      "C",
      "D",
      "E"
   ],
   "read_only": [
      "A",
      "B",
      "C",
      "D",
      "E"
   ],
   "read_write": [
      "A",
      "B",
      "C",
      "D",
      "E"
   ],
   "s3_bucket_full_control": "string",
   "s3_bucket_listing": "string",
   "s3_bucket_read": "string",
   "s3_bucket_read_acp": "string",
   "s3_bucket_write": "string",
   "s3_bucket_write_acp": "string",
   "s3_object_full_control": "string",
   "s3_object_read": "string",
   "s3_object_read_acp": "string",
   "s3_object_write": "string",
   "s3_object_write_acp": "string",
   "s3_read_only": [
      "A",
      "B",
      "C",
      "D",
      "E"
   ],
   "s3_read_write": [
      "A",
      "B",
      "C",
      "D",
      "E"
   ],
   "s3_visibility": [
      "A",
      "B",
      "C",
      "D",
      "E"
   ],
   "s3_visibility_groups": [
      "A",
      "B",
      "C",
      "D",
      "E"
   ],
   "smb_directory_mode": 100,
   "smb_directory_mode_padded": "string",
   "smb_file_mode": 100,
   "smb_file_mode_padded": "string",
   "smb_read_only": [
      "A",
      "B",
      "C",
      "D",
      "E"
   ],
   "smb_read_write": [
      "A",
      "B",
      "C",
      "D",
      "E"
   ],
   "tenant_id": 100,
   "tenant_name": "string",
   "trash_access": [
      "A",
      "B",
      "C",
      "D",
      "E"
   ],
   "url": "string",
   "vip_pools": [
      1,
      2,
      3,
      4,
      5,
      6
   ],
   "vippool_permissions": [
      {
         "vippool_id": "string",
         "vippool_permissions": "string"
      }
   ]
}
                         `
	var server *ghttp.Server
	var client *vast_client.VMSSession
	ViewPolicyDataSource := datasources.DataSourceViewPolicy()
	ReadContext = ViewPolicyDataSource.ReadContext

	BeforeEach(func() {
		ViewPolicyDataSourceData = ViewPolicyDataSource.TestResourceData()
		ViewPolicyDataSourceData.SetId("100")
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
			It("Datasource:ViewPolicy ,reads", func() {
				var ctx context.Context = context.Background()
				var output bytes.Buffer
				var guid = "11111-11111-11111-11111-11111"

				resource := api_latest.ViewPolicy{}
				json.Unmarshal([]byte(model_json), &resource)
				ViewPolicyDataSourceData.SetId(guid)
				resource.Guid = guid
				b, err := json.Marshal(&resource)
				Expect(err).To(BeNil())
				values := url.Values{}

				values.Add("name", fmt.Sprintf("%v", resource.Name))
				ViewPolicyDataSourceData.Set("name", resource.Name)

				server.AppendHandlers(ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "viewpolicies", values.Encode()), //since this is a test http server and will not return id upon POST (creation) so json will use the zero value
					ghttp.RespondWith(200, `[`+string(b)+`]`),
				),
				)

				ctx = tflogtest.RootLogger(ctx, &output)
				e := client.Start()
				Expect(e).To(BeNil())
				d := ReadContext(ctx, ViewPolicyDataSourceData, client)
				Expect(d).To(BeNil())

				o := new(map[string]interface{})
				json.Unmarshal(b, o)
				z := make(map[string]string)
				utils.MapToTFSchema(*o, &z, "")
				attributes := ViewPolicyDataSourceData.State().Attributes
				for k, v := range z {
					Expect(attributes).To(HaveKeyWithValue(k, v))
				}

			})
		})
	},
	)
})
