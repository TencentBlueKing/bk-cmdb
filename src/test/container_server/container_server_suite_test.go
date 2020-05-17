package container_server_test

import (
	"testing"

	"configcenter/src/test"
	"configcenter/src/test/reporter"
	"configcenter/src/test/util"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var header = test.GetHeader()
var containerServerClient = test.GetClientSet().ContainerServer()
var apiServerClient = test.GetClientSet().ApiServer()
var instClient = test.GetClientSet().TopoServer().Instance()

func TestContainerServer(t *testing.T) {
	RegisterFailHandler(util.Fail)
	reporters := []Reporter{
		reporter.NewHtmlReporter(test.GetReportDir()+"containerserver.html", test.GetReportUrl(), true),
	}
	RunSpecsWithDefaultAndCustomReporters(t, "ContainerServer Suite", reporters)
}

var _ BeforeSuite(func() {
	test.ClearDatabase()
})