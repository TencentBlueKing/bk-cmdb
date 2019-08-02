package host_server_test

import (
	"configcenter/src/test"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var header = test.GetHeader()
var hostServerClient = test.GetClientSet().HostServer()
var apiServerClient = test.GetClientSet().ApiServer()

func TestHostServer(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "HostServer Suite")
}

var _ = BeforeSuite(func() {
	test.ClearDatabase()
})

var _ = AfterSuite(func() {
	test.ClearDatabase()
})
