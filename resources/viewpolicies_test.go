package resources_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
	"github.com/vast-data/terraform-provider-vastdata.git/resources"
	utils "github.com/vast-data/terraform-provider-vastdata.git/utils"

	"github.com/hashicorp/terraform-plugin-log/tflogtest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api_latest "github.com/vast-data/terraform-provider-vastdata.git/codegen/latest"
	vast_client "github.com/vast-data/terraform-provider-vastdata.git/vast-client"
)

var _ = Describe(" ViewPolicy", func() {
	var ReadContext schema.ReadContextFunc
	var DeleteContext schema.DeleteContextFunc
	var CreateContext schema.CreateContextFunc
	var UpdateContext schema.UpdateContextFunc
	var Importer schema.ResourceImporter
	//	var ResourceSchema map[string]*schema.Schema
	//An empty resource data to be populated per test
	var ViewPolicyResourceData *schema.ResourceData
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
   "sync": "string",
   "sync_time": "string",
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
   ]
}
                         `
	var server *ghttp.Server
	var client vast_client.JwtSession
	ViewPolicyResource := resources.ResourceViewPolicy()
	ReadContext = ViewPolicyResource.ReadContext
	DeleteContext = ViewPolicyResource.DeleteContext
	CreateContext = ViewPolicyResource.CreateContext
	UpdateContext = ViewPolicyResource.UpdateContext
	Importer = *ViewPolicyResource.Importer
	//	ResourceSchema = ViewPolicyResource.Schema

	BeforeEach(func() {
		ViewPolicyResourceData = ViewPolicyResource.TestResourceData()
		ViewPolicyResourceData.SetId("100")
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
	Describe("Validating Resource Read Context", func() {
		Context("Read Data into a ResourceData", func() {
			It("Resource:ViewPolicy ,Reads Data", func() {
				var ctx context.Context = context.Background()
				var output bytes.Buffer

				ctx = tflogtest.RootLogger(ctx, &output)
				server.AppendHandlers(ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/api/viewpolicies/100"),
					ghttp.RespondWith(200, model_json),
				),
				)
				e := client.Start()
				Expect(e).To(BeNil())
				o := new(map[string]interface{})
				json.Unmarshal([]byte(model_json), o)

				z := make(map[string]string)
				utils.MapToTFSchema(*o, &z, "")
				d := ReadContext(ctx, ViewPolicyResourceData, client)
				Expect(d).To(BeNil())
				attributes := ViewPolicyResourceData.State().Attributes
				for k, v := range z {
					Expect(attributes).To(HaveKeyWithValue(k, v))
				}

			})
		})
	},
	)
	Describe("Validating Resource Delete Context", func() {
		Context("Delete A resource", func() {
			It("Resource:ViewPolicy ,Deletes the resource", func() {
				var ctx context.Context = context.Background()
				var output bytes.Buffer

				ctx = tflogtest.RootLogger(ctx, &output)
				server.AppendHandlers(ghttp.CombineHandlers(
					ghttp.VerifyRequest("DELETE", "/api/viewpolicies/100//"),
					ghttp.RespondWith(200, "DELETED"),
				),
				)
				e := client.Start()
				Expect(e).To(BeNil())
				d := DeleteContext(ctx, ViewPolicyResourceData, client)
				Expect(d).To(BeNil())
			})
		})
	},
	)
	Describe("Validating Resource Creation Context", func() {
		Context("Create A resource", func() {
			It("Resource:ViewPolicy ,Creates the resource", func() {
				var ctx context.Context = context.Background()
				var output bytes.Buffer

				ctx = tflogtest.RootLogger(ctx, &output)
				server.AppendHandlers(ghttp.CombineHandlers(
					ghttp.VerifyRequest("POST", "/api/viewpolicies/"),
					ghttp.RespondWith(200, model_json),
				),
				)
				server.AppendHandlers(ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/api/viewpolicies/0"), //since this is a test http server and will not return id upon POST (creation) so json will use the zero value
					ghttp.RespondWith(200, model_json),
				),
				)

				e := client.Start()
				Expect(e).To(BeNil())
				resource := api_latest.ViewPolicy{}
				json.Unmarshal([]byte(model_json), &resource)
				resources.ResourceViewPolicyReadStructIntoSchema(ctx, resource, ViewPolicyResourceData)
				ViewPolicyResourceData.SetId("100")
				d := CreateContext(ctx, ViewPolicyResourceData, client)
				Expect(d).To(BeNil())

			})
		})
	},
	)
	Describe("Validating Resource Update Context", func() {
		Context("Update A resource", func() {
			It("Resource:ViewPolicy ,Update the resource", func() {
				var ctx context.Context = context.Background()
				var output bytes.Buffer
				var new_guid = "11111-11111-11111-11111-11111"

				ctx = tflogtest.RootLogger(ctx, &output)
				server.AppendHandlers(ghttp.CombineHandlers(
					ghttp.VerifyRequest("POST", "/api/viewpolicies/"),
					ghttp.RespondWith(200, model_json),
				),
				)
				server.AppendHandlers(ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/api/viewpolicies/0"), //since this is a test http server and will not return id upon POST (creation) so json will use the zero value
					ghttp.RespondWith(200, model_json),
				),
				)

				server.AppendHandlers(ghttp.CombineHandlers(
					ghttp.VerifyRequest("PATCH", "/api/viewpolicies//0/"), //since this is a test http server and will not return id upon POST (creation) so json will use the zero value
					ghttp.RespondWith(200, model_json),
				),
				)

				e := client.Start()
				Expect(e).To(BeNil())
				//First we create a resource than we change it and see if it was updated
				resource := api_latest.ViewPolicy{}
				json.Unmarshal([]byte(model_json), &resource)
				resources.ResourceViewPolicyReadStructIntoSchema(ctx, resource, ViewPolicyResourceData)
				d := CreateContext(ctx, ViewPolicyResourceData, client)
				Expect(d).To(BeNil())
				//We update the guid as this is a fieled that always exists
				resource.Guid = new_guid
				b, err := json.Marshal(&resource)
				Expect(err).To(BeNil())

				server.AppendHandlers(ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/api/viewpolicies/0"), //the new_guid is returned and it should change the value of the resource
					ghttp.RespondWith(200, string(b)),
				),
				)

				d = UpdateContext(ctx, ViewPolicyResourceData, client)
				Expect(d).To(BeNil())
				Expect(ViewPolicyResourceData.Get("guid")).To(Equal(new_guid))

			})
		})
	},
	)
	Describe("Validating Resource Importer", func() {
		Context("Import A resource", func() {
			It("Resource:ViewPolicy ,Imports the resource", func() {
				var ctx context.Context = context.Background()
				var output bytes.Buffer
				var guid = "11111-11111-11111-11111-11111"

				resource := api_latest.ViewPolicy{}
				json.Unmarshal([]byte(model_json), &resource)
				ViewPolicyResourceData.SetId(guid)
				resource.Guid = guid
				b, err := json.Marshal(&resource)
				Expect(err).To(BeNil())

				server.AppendHandlers(ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/api/viewpolicies/", fmt.Sprintf("guid=%s", guid)), //since this is a test http server and will not return id upon POST (creation) so json will use the zero value
					ghttp.RespondWith(200, `[`+string(b)+`]`),
				),
				)

				ctx = tflogtest.RootLogger(ctx, &output)
				e := client.Start()
				Expect(e).To(BeNil())
				Importer.StateContext(ctx, ViewPolicyResourceData, client)
				//				Expect(d).To(BeNil())

				o := new(map[string]interface{})
				json.Unmarshal(b, o)
				z := make(map[string]string)
				utils.MapToTFSchema(*o, &z, "")
				attributes := ViewPolicyResourceData.State().Attributes
				for k, v := range z {
					Expect(attributes).To(HaveKeyWithValue(k, v))
				}

			})
		})
	},
	)
})
