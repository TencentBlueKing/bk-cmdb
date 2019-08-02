package host_server_test

import (
	"configcenter/src/test"
	"context"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var header = test.GetHeader()
var hostServerClient = test.GetClientSet().HostServer()

func TestHostServer(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "HostServer Suite")
}

var _ = Describe("user custom test", func() {
	It("search default user custom", func() {
		rsp, err := hostServerClient.GetUserCustom(context.Background(), header)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.Result).To(Equal(true))
	})

	It("save user custom", func() {
		input := map[string]interface{}{
			"index_v2_classify_navigation": []string{"bk_middleware"},
		}
		rsp, err := hostServerClient.SaveUserCustom(context.Background(), header, input)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.Result).To(Equal(true))
	})

	It("search user custom", func() {
		rsp, err := hostServerClient.GetUserCustom(context.Background(), header)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.Result).To(Equal(true))
		data := rsp.Data.(map[string][]string)
		Expect(data["index_v2_classify_navigation"][0]).To(Equal("bk_middleware"))
	})
})
