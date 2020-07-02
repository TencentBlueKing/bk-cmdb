package hostsnap

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestHostsnap(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Hostsnap Suite")
}
