package host_server_test

import (
	"testing"

	"configcenter/src/test"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var header = test.GetHeader()
var hostServerClient = test.GetClientSet().HostServer()
var apiServerClient = test.GetClientSet().ApiServer()
var instClient = test.GetClientSet().TopoServer().Instance()

func TestHostServer(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "HostServer Suite")
}

var _ = BeforeSuite(func() {
	test.ClearDatabase()
})

// var _ = AfterSuite(func() {
// 	test.ClearDatabase()
// })
