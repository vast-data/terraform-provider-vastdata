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

var _ = Describe(" QosPolicy", func() {
	var ReadContext schema.ReadContextFunc
	var DeleteContext schema.DeleteContextFunc
	var CreateContext schema.CreateContextFunc
	var UpdateContext schema.UpdateContextFunc
	var Importer schema.ResourceImporter
	//	var ResourceSchema map[string]*schema.Schema
	//An empty resource data to be populated per test
	var QosPolicyResourceData *schema.ResourceData
	var model_json = `
                         {
   "capacity_limits": {
      "max_reads_bw_mbps_per_gb_capacity": 100,
      "max_reads_iops_per_gb_capacity": 100,
      "max_writes_bw_mbps_per_gb_capacity": 100,
      "max_writes_iops_per_gb_capacity": 100
   },
   "guid": "string",
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
	QosPolicyResource := resources.ResourceQosPolicy()
	ReadContext = QosPolicyResource.ReadContext
	DeleteContext = QosPolicyResource.DeleteContext
	CreateContext = QosPolicyResource.CreateContext
	UpdateContext = QosPolicyResource.UpdateContext
	Importer = *QosPolicyResource.Importer
	//	ResourceSchema = QosPolicyResource.Schema

	BeforeEach(func() {
		QosPolicyResourceData = QosPolicyResource.TestResourceData()
		QosPolicyResourceData.SetId("100")
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
			It("Resource:QosPolicy ,Reads Data", func() {
				var ctx context.Context = context.Background()
				var output bytes.Buffer

				ctx = tflogtest.RootLogger(ctx, &output)
				server.AppendHandlers(ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/api/qospolicies/100"),
					ghttp.RespondWith(200, model_json),
				),
				)
				e := client.Start()
				Expect(e).To(BeNil())
				o := new(map[string]interface{})
				json.Unmarshal([]byte(model_json), o)

				z := make(map[string]string)
				utils.MapToTFSchema(*o, &z, "")
				d := ReadContext(ctx, QosPolicyResourceData, client)
				Expect(d).To(BeNil())
				attributes := QosPolicyResourceData.State().Attributes
				for k, v := range z {
					Expect(attributes).To(HaveKeyWithValue(k, v))
				}

			})
		})
	},
	)
	Describe("Validating Resource Delete Context", func() {
		Context("Delete A resource", func() {
			It("Resource:QosPolicy ,Deletes the resource", func() {
				var ctx context.Context = context.Background()
				var output bytes.Buffer

				ctx = tflogtest.RootLogger(ctx, &output)
				server.AppendHandlers(ghttp.CombineHandlers(
					ghttp.VerifyRequest("DELETE", "/api/qospolicies/100//"),
					ghttp.RespondWith(200, "DELETED"),
				),
				)
				e := client.Start()
				Expect(e).To(BeNil())
				d := DeleteContext(ctx, QosPolicyResourceData, client)
				Expect(d).To(BeNil())
			})
		})
	},
	)
	Describe("Validating Resource Creation Context", func() {
		Context("Create A resource", func() {
			It("Resource:QosPolicy ,Creates the resource", func() {
				var ctx context.Context = context.Background()
				var output bytes.Buffer

				ctx = tflogtest.RootLogger(ctx, &output)
				server.AppendHandlers(ghttp.CombineHandlers(
					ghttp.VerifyRequest("POST", "/api/qospolicies/"),
					ghttp.RespondWith(200, model_json),
				),
				)
				server.AppendHandlers(ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/api/qospolicies/0"), //since this is a test http server and will not return id upon POST (creation) so json will use the zero value
					ghttp.RespondWith(200, model_json),
				),
				)

				e := client.Start()
				Expect(e).To(BeNil())
				resource := api_latest.QosPolicy{}
				json.Unmarshal([]byte(model_json), &resource)
				resources.ResourceQosPolicyReadStructIntoSchema(ctx, resource, QosPolicyResourceData)
				QosPolicyResourceData.SetId("100")
				d := CreateContext(ctx, QosPolicyResourceData, client)
				Expect(d).To(BeNil())

			})
		})
	},
	)
	Describe("Validating Resource Update Context", func() {
		Context("Update A resource", func() {
			It("Resource:QosPolicy ,Update the resource", func() {
				var ctx context.Context = context.Background()
				var output bytes.Buffer
				var new_guid = "11111-11111-11111-11111-11111"

				ctx = tflogtest.RootLogger(ctx, &output)
				server.AppendHandlers(ghttp.CombineHandlers(
					ghttp.VerifyRequest("POST", "/api/qospolicies/"),
					ghttp.RespondWith(200, model_json),
				),
				)
				server.AppendHandlers(ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/api/qospolicies/0"), //since this is a test http server and will not return id upon POST (creation) so json will use the zero value
					ghttp.RespondWith(200, model_json),
				),
				)

				server.AppendHandlers(ghttp.CombineHandlers(
					ghttp.VerifyRequest("PATCH", "/api/qospolicies//0/"), //since this is a test http server and will not return id upon POST (creation) so json will use the zero value
					ghttp.RespondWith(200, model_json),
				),
				)

				e := client.Start()
				Expect(e).To(BeNil())
				//First we create a resource than we change it and see if it was updated
				resource := api_latest.QosPolicy{}
				json.Unmarshal([]byte(model_json), &resource)
				resources.ResourceQosPolicyReadStructIntoSchema(ctx, resource, QosPolicyResourceData)
				d := CreateContext(ctx, QosPolicyResourceData, client)
				Expect(d).To(BeNil())
				//We update the guid as this is a fieled that always exists
				resource.Guid = new_guid
				b, err := json.Marshal(&resource)
				Expect(err).To(BeNil())

				server.AppendHandlers(ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/api/qospolicies/0"), //the new_guid is returned and it should change the value of the resource
					ghttp.RespondWith(200, string(b)),
				),
				)

				d = UpdateContext(ctx, QosPolicyResourceData, client)
				Expect(d).To(BeNil())
				Expect(QosPolicyResourceData.Get("guid")).To(Equal(new_guid))

			})
		})
	},
	)
	Describe("Validating Resource Importer", func() {
		Context("Import A resource", func() {
			It("Resource:QosPolicy ,Imports the resource", func() {
				var ctx context.Context = context.Background()
				var output bytes.Buffer
				var guid = "11111-11111-11111-11111-11111"

				resource := api_latest.QosPolicy{}
				json.Unmarshal([]byte(model_json), &resource)
				QosPolicyResourceData.SetId(guid)
				resource.Guid = guid
				b, err := json.Marshal(&resource)
				Expect(err).To(BeNil())

				server.AppendHandlers(ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/api/qospolicies/", fmt.Sprintf("guid=%s", guid)), //since this is a test http server and will not return id upon POST (creation) so json will use the zero value
					ghttp.RespondWith(200, `[`+string(b)+`]`),
				),
				)

				ctx = tflogtest.RootLogger(ctx, &output)
				e := client.Start()
				Expect(e).To(BeNil())
				Importer.StateContext(ctx, QosPolicyResourceData, client)
				//				Expect(d).To(BeNil())

				o := new(map[string]interface{})
				json.Unmarshal(b, o)
				z := make(map[string]string)
				utils.MapToTFSchema(*o, &z, "")
				attributes := QosPolicyResourceData.State().Attributes
				for k, v := range z {
					Expect(attributes).To(HaveKeyWithValue(k, v))
				}

			})
		})
	},
	)
})
