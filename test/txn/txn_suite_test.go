package txn_test

import (
	"testing"

	"configcenter/test"
	"configcenter/test/reporter"
	"configcenter/test/util"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var clientSet = test.GetClientSet()

func TestTxn(t *testing.T) {
	RegisterFailHandler(util.Fail)
	reporters := []Reporter{
		reporter.NewHtmlReporter(test.GetReportDir()+"transaction.html", test.GetReportUrl(), true),
	}
	RunSpecsWithDefaultAndCustomReporters(t, "Transaction Suite", reporters)
}

var _ = BeforeSuite(func() {
	test.ClearDatabase()
})
