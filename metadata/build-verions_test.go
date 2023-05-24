package metadata_test

import (
	"github.com/vast-data/terraform-provider-vastdata.git/metadata"
	version "github.com/hashicorp/go-version"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Build Version Testing", func() {
	Describe("Test build-version functionality", func() {
		Context("Testing Cluster Version", func() {
			It("Get Cluster Build Version", func() {
				b := metadata.GetBuildVersion()
				v, _ := version.NewVersion("5.0.0")
				Expect(b).To(Equal(*v))

			})
			It("Compare To Equal Version", func() {
				metadata.UpdateClusterVersion("5.0.0")
				i := metadata.ClusterVersionCompare()
				Expect(i).To(Equal(metadata.CLUSTER_VERSION_EQUALS))

			})
			It("Compares To Higher Version", func() {
				metadata.UpdateClusterVersion("6.0.0")
				i := metadata.ClusterVersionCompare()
				Expect(i).To(Equal(metadata.CLUSTER_VERSION_GRATER))

			})
			It("Compares To Lower Version", func() {
				metadata.UpdateClusterVersion("1.0.0")
				i := metadata.ClusterVersionCompare()
				Expect(i).To(Equal(metadata.CLUSTER_VERSION_LOWER))

			})
			It("Get The Build Version String", func() {
				i := metadata.BuildVersionString()
				Expect(i).To(Equal("5.0.0"))

			})

		})

	})
})
