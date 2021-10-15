package healthz_test

import (
	"testing"

	"configcenter/src/common/types"
	"configcenter/src/test"
	"configcenter/src/test/reporter"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var header = test.GetHeader()
var clientSet = test.GetClientSet()
var healthzClient = clientSet.Healthz()

func TestHealthz(t *testing.T) {
	RegisterFailHandler(Fail)
	reporters := []Reporter{
		reporter.NewHtmlReporter(test.GetReportDir()+"healthz.html", test.GetReportUrl(), true),
	}
	RunSpecsWithDefaultAndCustomReporters(t, "Healthz Suite", reporters)
}

var _ = Describe("healthz test", func() {
	It(types.CCModuleDataCollection+" healthz test", func() {
		isHealthy, err := healthzClient.HealthCheck(types.CCModuleDataCollection)
		Expect(err).NotTo(HaveOccurred())
		Expect(isHealthy).To(Equal(true))
	})

	It(types.CCModuleHost+" healthz test", func() {
		isHealthy, err := healthzClient.HealthCheck(types.CCModuleHost)
		Expect(err).NotTo(HaveOccurred())
		Expect(isHealthy).To(Equal(true))
	})

	It(types.CCModuleAdmin+" healthz test", func() {
		isHealthy, err := healthzClient.HealthCheck(types.CCModuleAdmin)
		Expect(err).NotTo(HaveOccurred())
		Expect(isHealthy).To(Equal(true))
	})

	It(types.CCModuleProc+" healthz test", func() {
		isHealthy, err := healthzClient.HealthCheck(types.CCModuleProc)
		Expect(err).NotTo(HaveOccurred())
		Expect(isHealthy).To(Equal(true))
	})

	It(types.CCModuleTopo+" healthz test", func() {
		isHealthy, err := healthzClient.HealthCheck(types.CCModuleTopo)
		Expect(err).NotTo(HaveOccurred())
		Expect(isHealthy).To(Equal(true))
	})

	It(types.CCModuleEvent+" healthz test", func() {
		isHealthy, err := healthzClient.HealthCheck(types.CCModuleEvent)
		Expect(err).NotTo(HaveOccurred())
		Expect(isHealthy).To(Equal(true))
	})

	It(types.CCModuleApi+" healthz test", func() {
		isHealthy, err := healthzClient.HealthCheck(types.CCModuleApi)
		Expect(err).NotTo(HaveOccurred())
		Expect(isHealthy).To(Equal(true))
	})

	It(types.CCModuleCoreService+" healthz test", func() {
		isHealthy, err := healthzClient.HealthCheck(types.CCModuleCoreService)
		Expect(err).NotTo(HaveOccurred())
		Expect(isHealthy).To(Equal(true))
	})
})
