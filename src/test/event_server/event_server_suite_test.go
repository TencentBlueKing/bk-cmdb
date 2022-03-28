package event_server_test

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
var eventServerClient = clientSet.EventServer()

func TestEventServer(t *testing.T) {
	RegisterFailHandler(util.Fail)
	reporters := []Reporter{
		reporter.NewHtmlReporter(test.GetReportDir()+"eventserver.html", test.GetReportUrl(), true),
	}
	RunSpecsWithDefaultAndCustomReporters(t, "EventServer Suite", reporters)
}

var _ = BeforeSuite(func() {
	test.ClearDatabase()
})

var _ = Describe("event server test", func() {
	// NOTE: add more tests here.
})
