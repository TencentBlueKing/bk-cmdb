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
	It(types.CCModuleDataCollection.Name+" healthz test", func() {
		isHealthy, err := healthzClient.HealthCheck(types.CCModuleDataCollection.Name)
		Expect(err).NotTo(HaveOccurred())
		Expect(isHealthy).To(Equal(true))
	})

	It(types.CCModuleHost.Name+" healthz test", func() {
		isHealthy, err := healthzClient.HealthCheck(types.CCModuleHost.Name)
		Expect(err).NotTo(HaveOccurred())
		Expect(isHealthy).To(Equal(true))
	})

	It(types.CCModuleMigrate.Name+" healthz test", func() {
		isHealthy, err := healthzClient.HealthCheck(types.CCModuleMigrate.Name)
		Expect(err).NotTo(HaveOccurred())
		Expect(isHealthy).To(Equal(true))
	})

	It(types.CCModuleProc.Name+" healthz test", func() {
		isHealthy, err := healthzClient.HealthCheck(types.CCModuleProc.Name)
		Expect(err).NotTo(HaveOccurred())
		Expect(isHealthy).To(Equal(true))
	})

	It(types.CCModuleTop.Name+" healthz test", func() {
		isHealthy, err := healthzClient.HealthCheck(types.CCModuleTop.Name)
		Expect(err).NotTo(HaveOccurred())
		Expect(isHealthy).To(Equal(true))
	})

	It(types.CCModuleEventServer.Name+" healthz test", func() {
		isHealthy, err := healthzClient.HealthCheck(types.CCModuleEventServer.Name)
		Expect(err).NotTo(HaveOccurred())
		Expect(isHealthy).To(Equal(true))
	})

	It(types.CCModuleAPIServer.Name+" healthz test", func() {
		isHealthy, err := healthzClient.HealthCheck(types.CCModuleAPIServer.Name)
		Expect(err).NotTo(HaveOccurred())
		Expect(isHealthy).To(Equal(true))
	})

	It(types.CCModuleCoreService.Name+" healthz test", func() {
		isHealthy, err := healthzClient.HealthCheck(types.CCModuleCoreService.Name)
		Expect(err).NotTo(HaveOccurred())
		Expect(isHealthy).To(Equal(true))
	})
})
