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
	"github.com/vast-data/terraform-provider-vastdata/resources"
	utils "github.com/vast-data/terraform-provider-vastdata/utils"

	"github.com/hashicorp/terraform-plugin-log/tflogtest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api_latest "github.com/vast-data/terraform-provider-vastdata/codegen/latest"
	vast_client "github.com/vast-data/terraform-provider-vastdata/vast-client"
)

var _ = Describe(" GlobalSnapshot", func() {
	var ReadContext schema.ReadContextFunc
	var DeleteContext schema.DeleteContextFunc
	var CreateContext schema.CreateContextFunc
	var UpdateContext schema.UpdateContextFunc
	var Importer schema.ResourceImporter
	//	var ResourceSchema map[string]*schema.Schema
	//An empty resource data to be populated per test
	var GlobalSnapshotResourceData *schema.ResourceData
	var model_json = `
                         {
   "guid": "string",
   "loanee_root_path": "string",
   "loanee_snapsho": "string",
   "loanee_snapshot_id": 100,
   "name": "string",
   "remote_target_id": 100,
   "source_cluster": "string",
   "source_path": "string",
   "target_cluster": "string",
   "tenant_id": 100
}
                         `
	var server *ghttp.Server
	var client vast_client.JwtSession
	GlobalSnapshotResource := resources.ResourceGlobalSnapshot()
	ReadContext = GlobalSnapshotResource.ReadContext
	DeleteContext = GlobalSnapshotResource.DeleteContext
	CreateContext = GlobalSnapshotResource.CreateContext
	UpdateContext = GlobalSnapshotResource.UpdateContext
	Importer = *GlobalSnapshotResource.Importer
	//	ResourceSchema = GlobalSnapshotResource.Schema

	BeforeEach(func() {
		GlobalSnapshotResourceData = GlobalSnapshotResource.TestResourceData()
		GlobalSnapshotResourceData.SetId("100")
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
			It("Resource:GlobalSnapshot ,Reads Data", func() {
				var ctx context.Context = context.Background()
				var output bytes.Buffer

				ctx = tflogtest.RootLogger(ctx, &output)
				server.AppendHandlers(ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/api/globalsnapstreams/100"),
					ghttp.RespondWith(200, model_json),
				),
				)
				e := client.Start()
				Expect(e).To(BeNil())
				o := new(map[string]interface{})
				json.Unmarshal([]byte(model_json), o)

				z := make(map[string]string)
				utils.MapToTFSchema(*o, &z, "")
				d := ReadContext(ctx, GlobalSnapshotResourceData, client)
				Expect(d).To(BeNil())
				attributes := GlobalSnapshotResourceData.State().Attributes
				for k, v := range z {
					Expect(attributes).To(HaveKeyWithValue(k, v))
				}

			})
		})
	},
	)
	Describe("Validating Resource Delete Context", func() {
		Context("Delete A resource", func() {
			It("Resource:GlobalSnapshot ,Deletes the resource", func() {
				var ctx context.Context = context.Background()
				var output bytes.Buffer

				ctx = tflogtest.RootLogger(ctx, &output)
				server.AppendHandlers(ghttp.CombineHandlers(
					ghttp.VerifyRequest("DELETE", "/api/globalsnapstreams/100//"),
					ghttp.RespondWith(200, "DELETED"),
				),
				)
				e := client.Start()
				Expect(e).To(BeNil())
				d := DeleteContext(ctx, GlobalSnapshotResourceData, client)
				Expect(d).To(BeNil())
			})
		})
	},
	)
	Describe("Validating Resource Creation Context", func() {
		Context("Create A resource", func() {
			It("Resource:GlobalSnapshot ,Creates the resource", func() {
				var ctx context.Context = context.Background()
				var output bytes.Buffer

				ctx = tflogtest.RootLogger(ctx, &output)
				server.AppendHandlers(ghttp.CombineHandlers(
					ghttp.VerifyRequest("POST", "/api/globalsnapstreams/"),
					ghttp.RespondWith(200, model_json),
				),
				)
				server.AppendHandlers(ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/api/globalsnapstreams/0"), //since this is a test http server and will not return id upon POST (creation) so json will use the zero value
					ghttp.RespondWith(200, model_json),
				),
				)

				e := client.Start()
				Expect(e).To(BeNil())
				resource := api_latest.GlobalSnapshot{}
				json.Unmarshal([]byte(model_json), &resource)
				resources.ResourceGlobalSnapshotReadStructIntoSchema(ctx, resource, GlobalSnapshotResourceData)
				GlobalSnapshotResourceData.SetId("100")
				d := CreateContext(ctx, GlobalSnapshotResourceData, client)
				Expect(d).To(BeNil())

			})
		})
	},
	)
	Describe("Validating Resource Update Context", func() {
		Context("Update A resource", func() {
			It("Resource:GlobalSnapshot ,Update the resource", func() {
				var ctx context.Context = context.Background()
				var output bytes.Buffer
				var new_guid = "11111-11111-11111-11111-11111"

				ctx = tflogtest.RootLogger(ctx, &output)
				server.AppendHandlers(ghttp.CombineHandlers(
					ghttp.VerifyRequest("POST", "/api/globalsnapstreams/"),
					ghttp.RespondWith(200, model_json),
				),
				)
				server.AppendHandlers(ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/api/globalsnapstreams/0"), //since this is a test http server and will not return id upon POST (creation) so json will use the zero value
					ghttp.RespondWith(200, model_json),
				),
				)

				server.AppendHandlers(ghttp.CombineHandlers(
					ghttp.VerifyRequest("PATCH", "/api/globalsnapstreams//0/"), //since this is a test http server and will not return id upon POST (creation) so json will use the zero value
					ghttp.RespondWith(200, model_json),
				),
				)

				e := client.Start()
				Expect(e).To(BeNil())
				//First we create a resource than we change it and see if it was updated
				resource := api_latest.GlobalSnapshot{}
				json.Unmarshal([]byte(model_json), &resource)
				resources.ResourceGlobalSnapshotReadStructIntoSchema(ctx, resource, GlobalSnapshotResourceData)
				d := CreateContext(ctx, GlobalSnapshotResourceData, client)
				Expect(d).To(BeNil())
				//We update the guid as this is a fieled that always exists
				resource.Guid = new_guid
				b, err := json.Marshal(&resource)
				Expect(err).To(BeNil())

				server.AppendHandlers(ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/api/globalsnapstreams/0"), //the new_guid is returned and it should change the value of the resource
					ghttp.RespondWith(200, string(b)),
				),
				)

				d = UpdateContext(ctx, GlobalSnapshotResourceData, client)
				Expect(d).To(BeNil())
				Expect(GlobalSnapshotResourceData.Get("guid")).To(Equal(new_guid))

			})
		})
	},
	)
	Describe("Validating Resource Importer", func() {
		Context("Import A resource", func() {
			It("Resource:GlobalSnapshot ,Imports the resource", func() {
				var ctx context.Context = context.Background()
				var output bytes.Buffer
				var guid = "11111-11111-11111-11111-11111"

				resource := api_latest.GlobalSnapshot{}
				json.Unmarshal([]byte(model_json), &resource)
				GlobalSnapshotResourceData.SetId(guid)
				resource.Guid = guid
				b, err := json.Marshal(&resource)
				Expect(err).To(BeNil())

				server.AppendHandlers(ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/api/globalsnapstreams/", fmt.Sprintf("guid=%s", guid)), //since this is a test http server and will not return id upon POST (creation) so json will use the zero value
					ghttp.RespondWith(200, `[`+string(b)+`]`),
				),
				)

				ctx = tflogtest.RootLogger(ctx, &output)
				e := client.Start()
				Expect(e).To(BeNil())
				Importer.StateContext(ctx, GlobalSnapshotResourceData, client)
				//				Expect(d).To(BeNil())

				o := new(map[string]interface{})
				json.Unmarshal(b, o)
				z := make(map[string]string)
				utils.MapToTFSchema(*o, &z, "")
				attributes := GlobalSnapshotResourceData.State().Attributes
				for k, v := range z {
					Expect(attributes).To(HaveKeyWithValue(k, v))
				}

			})
		})
	},
	)
})
