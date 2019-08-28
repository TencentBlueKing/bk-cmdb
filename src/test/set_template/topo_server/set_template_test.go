package topo_server_test

import (
	"context"

	"configcenter/src/common/metadata"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("set template test", func() {
	bizID := int64(3)
	ctx := context.Background()
	It("create set template", func() {
		option := metadata.CreateSetTemplateOption{
			Name:               "eereeede",
			ServiceTemplateIDs: nil,
		}
		rsp, err := topoServerClient.SetTemplate().CreateSetTemplate(ctx, header, bizID, option)
		Expect(err).NotTo(HaveOccurred())
		Expect(rsp.Result).To(Equal(true))
		Expect(rsp.Data).To(ContainElement("eereeede"))
		Expect(rsp.Data.ID).To(Not(Equal(int64(0))))
	})
})
