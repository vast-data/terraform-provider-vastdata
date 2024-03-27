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

var _ = Describe(" Quota", func() {
	var ReadContext schema.ReadContextFunc
	var DeleteContext schema.DeleteContextFunc
	var CreateContext schema.CreateContextFunc
	var UpdateContext schema.UpdateContextFunc
	var Importer schema.ResourceImporter
	//	var ResourceSchema map[string]*schema.Schema
	//An empty resource data to be populated per test
	var QuotaResourceData *schema.ResourceData
	var model_json = `
                         {
   "cluster": "string",
   "cluster_id": 100,
   "default_email": "string",
   "default_group_quota": {
      "grace_period": "string",
      "hard_limit": 100,
      "hard_limit_inodes": 100,
      "quota_system_id": 100,
      "sof_limit_inodes": 100,
      "soft_limit": 100
   },
   "default_user_quota": {
      "grace_period": "string",
      "hard_limit": 100,
      "hard_limit_inodes": 100,
      "quota_system_id": 100,
      "sof_limit_inodes": 100,
      "soft_limit": 100
   },
   "grace_period": "string",
   "group_quotas": [
      {
         "entity": {
            "email": "string",
            "identifier": "string",
            "identifier_type": "string",
            "name": "string",
            "vast_id": 100
         },
         "grace_period": "string",
         "hard_limit": 100,
         "hard_limit_inodes": 100,
         "quota_system_id": 100,
         "soft_limit": 100,
         "soft_limit_inodes": 100,
         "time_to_block": "string",
         "used_capacity": 100,
         "used_inodes": 100
      }
   ],
   "guid": "string",
   "hard_limit": 100,
   "hard_limit_inodes": 100,
   "name": "string",
   "num_blocked_users": 100,
   "num_exceeded_users": 100,
   "path": "string",
   "percent_capacity": 100,
   "percent_inodes": 100,
   "pretty_grace_period": "string",
   "pretty_state": "string",
   "soft_limit": 100,
   "soft_limit_inodes": 100,
   "state": "string",
   "system_id": 100,
   "tenant_id": 100,
   "tenant_name": "string",
   "time_to_block": "string",
   "used_capacity": 100,
   "used_capacity_tb": 10.5,
   "used_effective_capacity": 100,
   "used_effective_capacity_tb": 10.5,
   "used_inodes": 100,
   "user_quotas": [
      {
         "entity": {
            "email": "string",
            "identifier": "string",
            "identifier_type": "string",
            "name": "string",
            "vast_id": 100
         },
         "grace_period": "string",
         "hard_limit": 100,
         "hard_limit_inodes": 100,
         "quota_system_id": 100,
         "soft_limit": 100,
         "soft_limit_inodes": 100,
         "time_to_block": "string",
         "used_capacity": 100,
         "used_inodes": 100
      }
   ]
}
                         `
	var server *ghttp.Server
	var client vast_client.JwtSession
	QuotaResource := resources.ResourceQuota()
	ReadContext = QuotaResource.ReadContext
	DeleteContext = QuotaResource.DeleteContext
	CreateContext = QuotaResource.CreateContext
	UpdateContext = QuotaResource.UpdateContext
	Importer = *QuotaResource.Importer
	//	ResourceSchema = QuotaResource.Schema

	BeforeEach(func() {
		QuotaResourceData = QuotaResource.TestResourceData()
		QuotaResourceData.SetId("100")
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
			It("Resource:Quota ,Reads Data", func() {
				var ctx context.Context = context.Background()
				var output bytes.Buffer

				ctx = tflogtest.RootLogger(ctx, &output)
				server.AppendHandlers(ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "quotas100"),
					ghttp.RespondWith(200, model_json),
				),
				)
				e := client.Start()
				Expect(e).To(BeNil())
				o := new(map[string]interface{})
				json.Unmarshal([]byte(model_json), o)

				z := make(map[string]string)
				utils.MapToTFSchema(*o, &z, "")
				d := ReadContext(ctx, QuotaResourceData, client)
				Expect(d).To(BeNil())
				attributes := QuotaResourceData.State().Attributes
				for k, v := range z {
					Expect(attributes).To(HaveKeyWithValue(k, v))
				}

			})
		})
	},
	)
	Describe("Validating Resource Delete Context", func() {
		Context("Delete A resource", func() {
			It("Resource:Quota ,Deletes the resource", func() {
				var ctx context.Context = context.Background()
				var output bytes.Buffer

				ctx = tflogtest.RootLogger(ctx, &output)
				server.AppendHandlers(ghttp.CombineHandlers(
					ghttp.VerifyRequest("DELETE", "quotas100//"),
					ghttp.RespondWith(200, "DELETED"),
				),
				)
				e := client.Start()
				Expect(e).To(BeNil())
				d := DeleteContext(ctx, QuotaResourceData, client)
				Expect(d).To(BeNil())
			})
		})
	},
	)
	Describe("Validating Resource Creation Context", func() {
		Context("Create A resource", func() {
			It("Resource:Quota ,Creates the resource", func() {
				var ctx context.Context = context.Background()
				var output bytes.Buffer

				ctx = tflogtest.RootLogger(ctx, &output)
				server.AppendHandlers(ghttp.CombineHandlers(
					ghttp.VerifyRequest("POST", "quotas"),
					ghttp.RespondWith(200, model_json),
				),
				)
				server.AppendHandlers(ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "quotas0"), //since this is a test http server and will not return id upon POST (creation) so json will use the zero value
					ghttp.RespondWith(200, model_json),
				),
				)

				e := client.Start()
				Expect(e).To(BeNil())
				resource := api_latest.Quota{}
				json.Unmarshal([]byte(model_json), &resource)
				resources.ResourceQuotaReadStructIntoSchema(ctx, resource, QuotaResourceData)
				QuotaResourceData.SetId("100")
				d := CreateContext(ctx, QuotaResourceData, client)
				Expect(d).To(BeNil())

			})
		})
	},
	)
	Describe("Validating Resource Update Context", func() {
		Context("Update A resource", func() {
			It("Resource:Quota ,Update the resource", func() {
				var ctx context.Context = context.Background()
				var output bytes.Buffer
				var new_guid = "11111-11111-11111-11111-11111"

				ctx = tflogtest.RootLogger(ctx, &output)
				server.AppendHandlers(ghttp.CombineHandlers(
					ghttp.VerifyRequest("POST", "quotas"),
					ghttp.RespondWith(200, model_json),
				),
				)
				server.AppendHandlers(ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "quotas0"), //since this is a test http server and will not return id upon POST (creation) so json will use the zero value
					ghttp.RespondWith(200, model_json),
				),
				)

				server.AppendHandlers(ghttp.CombineHandlers(
					ghttp.VerifyRequest("PATCH", "quotas/0/"), //since this is a test http server and will not return id upon POST (creation) so json will use the zero value
					ghttp.RespondWith(200, model_json),
				),
				)

				e := client.Start()
				Expect(e).To(BeNil())
				//First we create a resource than we change it and see if it was updated
				resource := api_latest.Quota{}
				json.Unmarshal([]byte(model_json), &resource)
				resources.ResourceQuotaReadStructIntoSchema(ctx, resource, QuotaResourceData)
				d := CreateContext(ctx, QuotaResourceData, client)
				Expect(d).To(BeNil())
				//We update the guid as this is a fieled that always exists
				resource.Guid = new_guid
				b, err := json.Marshal(&resource)
				Expect(err).To(BeNil())

				server.AppendHandlers(ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "quotas0"), //the new_guid is returned and it should change the value of the resource
					ghttp.RespondWith(200, string(b)),
				),
				)

				d = UpdateContext(ctx, QuotaResourceData, client)
				Expect(d).To(BeNil())
				Expect(QuotaResourceData.Get("guid")).To(Equal(new_guid))

			})
		})
	},
	)
	Describe("Validating Resource Importer", func() {
		Context("Import A resource", func() {
			It("Resource:Quota ,Imports the resource", func() {
				var ctx context.Context = context.Background()
				var output bytes.Buffer
				var guid = "11111-11111-11111-11111-11111"

				resource := api_latest.Quota{}
				json.Unmarshal([]byte(model_json), &resource)
				QuotaResourceData.SetId(guid)
				resource.Guid = guid
				b, err := json.Marshal(&resource)
				Expect(err).To(BeNil())

				request_url := `[{"url":"https://` + server.Addr() + `quotas100"}]`
				server.AppendHandlers(ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "quotas", fmt.Sprintf("guid=%s", guid)),
					ghttp.RespondWith(200, request_url),
				),
				)
				server.AppendHandlers(ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "quotas100"),
					ghttp.RespondWith(200, string(b)),
				),
				)

				ctx = tflogtest.RootLogger(ctx, &output)
				e := client.Start()
				Expect(e).To(BeNil())
				Importer.StateContext(ctx, QuotaResourceData, client)
				//				Expect(d).To(BeNil())

				o := new(map[string]interface{})
				json.Unmarshal(b, o)
				z := make(map[string]string)
				utils.MapToTFSchema(*o, &z, "")
				attributes := QuotaResourceData.State().Attributes
				for k, v := range z {
					Expect(attributes).To(HaveKeyWithValue(k, v))
				}

			})
		})
	},
	)
})
