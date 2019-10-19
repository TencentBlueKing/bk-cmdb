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

	It(types.CCModuleMigrate+" healthz test", func() {
		isHealthy, err := healthzClient.HealthCheck(types.CCModuleMigrate)
		Expect(err).NotTo(HaveOccurred())
		Expect(isHealthy).To(Equal(true))
	})

	It(types.CCModuleProc+" healthz test", func() {
		isHealthy, err := healthzClient.HealthCheck(types.CCModuleProc)
		Expect(err).NotTo(HaveOccurred())
		Expect(isHealthy).To(Equal(true))
	})

	It(types.CCModuleTop+" healthz test", func() {
		isHealthy, err := healthzClient.HealthCheck(types.CCModuleTop)
		Expect(err).NotTo(HaveOccurred())
		Expect(isHealthy).To(Equal(true))
	})

	It(types.CCModuleEventServer+" healthz test", func() {
		isHealthy, err := healthzClient.HealthCheck(types.CCModuleEventServer)
		Expect(err).NotTo(HaveOccurred())
		Expect(isHealthy).To(Equal(true))
	})

	It(types.CCModuleAPIServer+" healthz test", func() {
		isHealthy, err := healthzClient.HealthCheck(types.CCModuleAPIServer)
		Expect(err).NotTo(HaveOccurred())
		Expect(isHealthy).To(Equal(true))
	})

	It(types.CCModuleCoerService+" healthz test", func() {
		isHealthy, err := healthzClient.HealthCheck(types.CCModuleCoerService)
		Expect(err).NotTo(HaveOccurred())
		Expect(isHealthy).To(Equal(true))
	})
})
