package cloud_server_test

import (
	"testing"

	"configcenter/src/test"
	"configcenter/src/test/reporter"
	"configcenter/src/test/util"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var header = test.GetHeader()
var cloudServerClient = test.GetClientSet().CloudServer()
var hostServerClient = test.GetClientSet().HostServer()

func TestCloudServer(t *testing.T) {
	RegisterFailHandler(util.Fail)
	reporters := []Reporter{
		reporter.NewHtmlReporter(test.GetReportDir()+"cloudserver.html", test.GetReportUrl(), true),
	}
	RunSpecsWithDefaultAndCustomReporters(t, "CloudServer Suite", reporters)
}

var _ = BeforeSuite(func() {
	test.ClearDatabase()
})
