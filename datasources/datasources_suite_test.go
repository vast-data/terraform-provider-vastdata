package datasources_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestDatasources(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Datasources Suite")
}
