package metadata_test

import (
	"github.com/vast-data/terraform-provider-vastdata.git/metadata"
	version "github.com/hashicorp/go-version"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Cluster Version Version", func() {
	Describe("Test cluster-version functionality", func() {
		Context("Updating Cluster Version", func() {
			It("Updates the cluster version", func() {
				err := metadata.UpdateClusterVersion("1.1.1")
				Expect(err).To(BeNil())
				err = metadata.UpdateClusterVersion("ThisIsNotAValidVersion")
				Expect(err).NotTo(BeNil())

			})
			It("Setup The right version", func() {
				err := metadata.UpdateClusterVersion("1.1.1")
				Expect(err).To(BeNil())
				v := metadata.GetClusterVersion()
				o, _ := version.NewVersion("1.1.1")
				n := o.Core()
				Expect(v).To(Equal(*n))
			})
			It("Shows the version string as it should", func() {
				err := metadata.UpdateClusterVersion("2.2.2")
				Expect(err).To(BeNil())
				c := metadata.ClusterVersionString()
				Expect(c).To(Equal("2.2.2"))
			})
		})

	})
})
