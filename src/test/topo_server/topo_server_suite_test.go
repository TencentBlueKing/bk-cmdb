package topo_server_test

import (
	"testing"

	"configcenter/src/test"
	"configcenter/src/test/reporter"
	"configcenter/src/test/util"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var header = test.GetHeader()
var clientSet = test.GetClientSet()
var topoServerClient = clientSet.TopoServer()
var procServerClient = clientSet.ProcServer()
var apiServerClient = clientSet.ApiServer()
var instClient = topoServerClient.Instance()
var asstClient = topoServerClient.Association()
var objectClient = topoServerClient.Object()
var serviceClient = clientSet.ProcServer().Service()

func TestTopoServer(t *testing.T) {
	RegisterFailHandler(util.Fail)
	reporters := []Reporter{
		reporter.NewHtmlReporter(test.GetReportDir()+"toposerver.html", test.GetReportUrl(), true),
	}
	RunSpecsWithDefaultAndCustomReporters(t, "TopoServer Suite", reporters)
}

var _ = BeforeSuite(func() {
	test.ClearDatabase()
})
