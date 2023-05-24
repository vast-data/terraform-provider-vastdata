package vast_versions_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"reflect"

	version_4_7_0 "github.com/vast-data/terraform-provider-vastdata.git/codegen/4.7.0"
	vast_versions "github.com/vast-data/terraform-provider-vastdata.git/vast_versions"
)

var _ = Describe("Test Vast Version", func() {
	Describe("Test GetVersionedType functionality", func() {
		Context("Reading Values", func() {
			It("Reads the appropriate value", func() {
				t, b := vast_versions.GetVersionedType("4.7.0", "Dns")
				Expect(b).To(BeTrue())
				Expect(t).To(Equal(reflect.TypeOf((*version_4_7_0.Dns)(nil)).Elem()))
			})
			It("Dont Get Confuse Between versions", func() {
				t, b := vast_versions.GetVersionedType("4.7.0", "Dns")
				Expect(b).To(BeTrue())
				i, j := vast_versions.GetVersionedType("4.6.0", "Dns")
				Expect(j).To(BeTrue())
				Expect(t).NotTo(Equal(i))
			})
			It("Dont provide anything on wrong version", func() {
				t, b := vast_versions.GetVersionedType("4.0.0", "Dummy")
				Expect(b).To(BeFalse())
				Expect(t).To(BeNil())
			})
			It("Dont provide anything on wrong type", func() {
				t, b := vast_versions.GetVersionedType("4.6.0", "Dummy")
				Expect(b).To(BeFalse())
				Expect(t).To(BeNil())
			})

		})
	})

})
