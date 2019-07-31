package topo_server_test

import (
	"configcenter/src/test"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var header = test.GetHeader()
var clientSet = test.GetClientSet()
var topoServerClient = clientSet.TopoServer()
var apiServerClient = clientSet.ApiServer()

func TestTopoServer(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "TopoServer Suite")
}

var _ = BeforeSuite(func() {
	test.ClearDatabase()
})

var _ = AfterSuite(func() {
	test.ClearDatabase()
})
