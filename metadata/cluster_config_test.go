package metadata_test

import (
	"github.com/vast-data/terraform-provider-vastdata.git/metadata"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Testing Cluster Config", func() {
	Describe("Testing Cluster Config", func() {
		Context("Testing Cluster Config", func() {
			It("Set/Get A value", func() {
				metadata.SetClusterConfig("key", "dmklww1d9012ms901mi290000000m190")
				key, exists := metadata.GetClusterConfig("key")
				Expect(key).To(Equal("dmklww1d9012ms901mi290000000m190"))
				Expect(exists).To(BeTrue())

			})

		})

	})
})
