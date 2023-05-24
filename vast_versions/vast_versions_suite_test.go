package vast_versions_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestVastVersions(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "VastVersions Suite")
}
