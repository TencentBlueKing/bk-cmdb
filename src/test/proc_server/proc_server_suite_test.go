package proc_server_test

import (
	"configcenter/src/test"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var header = test.GetHeader()
var clientSet = test.GetClientSet()
var procServerClient = clientSet.ProcServer()

func TestProcServer(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ProcServer Suite")
}
