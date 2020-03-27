package test

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"strings"
	"time"

	"configcenter/src/apimachinery"
	"configcenter/src/apimachinery/discovery"
	"configcenter/src/apimachinery/util"
	"configcenter/src/common"
	"configcenter/src/common/backbone/service_mange/zk"
	"configcenter/src/storage/dal/mongo/local"
	"configcenter/src/test/run"
	testutil "configcenter/src/test/util"

	. "github.com/onsi/gomega"
)

var clientSet apimachinery.ClientSetInterface
var tConfig TestConfig
var header http.Header
var reportUrl string
var reportDir string

type TestConfig struct {
	ZkAddr         string
	Concurrent     int
	SustainSeconds int
	MongoURI       string
}

func init() {
	flag.StringVar(&tConfig.ZkAddr, "zk-addr", "127.0.0.1:2181", "zk discovery addresses, comma separated.")
	flag.IntVar(&tConfig.Concurrent, "concurrent", 100, "concurrent request during the load test.")
	flag.IntVar(&tConfig.SustainSeconds, "sustain-seconds", 10, "the load test sustain time in seconds ")
	flag.StringVar(&tConfig.MongoURI, "mongo-addr", "mongodb://127.0.0.1:27017/cmdb", "mongodb URI")
	flag.StringVar(&reportUrl, "report-url", "http://127.0.0.1:8080/", "html report base url")
	flag.StringVar(&reportDir, "report-dir", "report", "report directory")
	flag.Parse()

	run.Concurrent = tConfig.Concurrent
	run.SustainSeconds = tConfig.SustainSeconds

	RegisterFailHandler(testutil.Fail)
	fmt.Println("before suit")
	js, _ := json.MarshalIndent(tConfig, "", "    ")
	fmt.Printf("test config: %s\n", run.SetRed(string(js)))
	client := zk.NewZkClient(tConfig.ZkAddr, 5*time.Second)
	Expect(client.Start()).Should(BeNil())
	Expect(client.Ping()).Should(BeNil())
	disc, err := discovery.NewServiceDiscovery(client)
	Expect(err).Should(BeNil())
	c := &util.APIMachineryConfig{
		QPS:       20000,
		Burst:     10000,
		TLSConfig: nil,
	}
	clientSet, err = apimachinery.NewApiMachinery(c, disc)
	Expect(err).Should(BeNil())
	// wait for get the apiserver address.
	time.Sleep(1 * time.Second)
	fmt.Println("**** initialize clientSet success ***")
}

func GetClientSet() apimachinery.ClientSetInterface {
	return clientSet
}

func GetHeader() http.Header {
	header = make(http.Header)
	header.Add(common.BKHTTPOwnerID, "0")
	header.Add(common.BKSupplierIDField, "0")
	header.Add(common.BKHTTPHeaderUser, "admin")
	header.Add("Content-Type", "application/json")
	return header
}

func ClearDatabase() {
	// clientSet.AdminServer().ClearDatabase(context.Background(), GetHeader())
	db, err := local.NewMgo(tConfig.MongoURI, time.Minute)
	Expect(err).Should(BeNil())
	for _, tableName := range common.AllTables {
		db.DropTable(tableName)
	}
	db.Close()
	clientSet.AdminServer().Migrate(context.Background(), "0", "community", header)
}

func GetReportUrl() string {
	if !strings.HasSuffix(reportUrl, "/") {
		reportUrl = reportUrl + "/"
	}
	return reportUrl
}

func GetReportDir() string {
	if !strings.HasSuffix(reportDir, "/") {
		reportDir = reportDir + "/"
	}
	return reportDir
}
