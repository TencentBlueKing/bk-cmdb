package test

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"configcenter/src/apimachinery"
	"configcenter/src/apimachinery/discovery"
	"configcenter/src/apimachinery/util"
	"configcenter/src/common"
	"configcenter/src/common/backbone/service_mange/zk"
	"configcenter/src/storage/dal/mongo"
	"configcenter/src/storage/dal/mongo/local"
	"configcenter/src/test/run"
	testutil "configcenter/src/test/util"

	. "github.com/onsi/gomega"
)

var clientSet apimachinery.ClientSetInterface
var tConfig TestConfig
var reportUrl string
var reportDir string
var db *local.Mongo

type TestConfig struct {
	ZkAddr         string
	Concurrent     int
	SustainSeconds float64
	TotalRequest   int64
	DBWriteKBSize  int
	MongoURI       string
	MongoRsName    string
	RedisCfg       RedisConfig
}

type RedisConfig struct {
	RedisAdress string
	RedisPort   string
	RedisPasswd string
}

func init() {
	testing.Init()
	flag.StringVar(&tConfig.ZkAddr, "zk-addr", "127.0.0.1:2181", "zk discovery addresses, comma separated.")
	flag.IntVar(&tConfig.Concurrent, "concurrent", 100, "concurrent request during the load test.")
	flag.Float64Var(&tConfig.SustainSeconds, "sustain-seconds", 10, "the load test sustain time in seconds ")
	flag.Int64Var(&tConfig.TotalRequest, "total-request", 0, "the load test total request,it has higher priority than SustainSeconds")
	flag.IntVar(&tConfig.DBWriteKBSize, "write-size", 1, "MongoDB write size , unit is KB.")
	flag.StringVar(&tConfig.RedisCfg.RedisAdress, "redis-addr", "127.0.0.1:6379", "redis host address with port")
	flag.StringVar(&tConfig.RedisCfg.RedisPasswd, "redis-passwd", "cc", "redis password")
	flag.StringVar(&tConfig.MongoURI, "mongo-addr", "mongodb://127.0.0.1:27017/cmdb", "mongodb URI")
	flag.StringVar(&tConfig.MongoRsName, "mongo-rs-name", "rs0", "mongodb replica set name")
	flag.StringVar(&reportUrl, "report-url", "http://127.0.0.1:8080/", "html report base url")
	flag.StringVar(&reportDir, "report-dir", "report", "report directory")
	flag.Parse()

	run.Concurrent = tConfig.Concurrent
	run.SustainSeconds = tConfig.SustainSeconds
	run.TotalRequest = tConfig.TotalRequest

	RegisterFailHandler(testutil.Fail)
	fmt.Println("before suit")
	js, _ := json.MarshalIndent(tConfig, "", "    ")
	fmt.Printf("test config: %s\n", run.SetRed(string(js)))
	client := zk.NewZkClient(tConfig.ZkAddr, 40*time.Second)
	var err error
	mongoConfig := local.MongoConf{
		MaxOpenConns: mongo.DefaultMaxOpenConns,
		MaxIdleConns: mongo.MinimumMaxIdleOpenConns,
		URI:          tConfig.MongoURI,
		RsName:       tConfig.MongoRsName,
	}
	db, err = local.NewMgo(mongoConfig, time.Minute)
	Expect(err).Should(BeNil())
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

func GetTestConfig() TestConfig {
	return tConfig
}

func GetHeader() http.Header {
	header := make(http.Header)
	header.Add(common.BKHTTPOwnerID, "0")
	header.Add(common.BKHTTPHeaderUser, "admin")
	header.Add("Content-Type", "application/json")
	return header
}

func ClearDatabase() {
	fmt.Println("********Clear Database*************")
	// clientSet.AdminServer().ClearDatabase(context.Background(), GetHeader())
	mongoConfig := local.MongoConf{
		MaxOpenConns: mongo.DefaultMaxOpenConns,
		MaxIdleConns: mongo.MinimumMaxIdleOpenConns,
		URI:          tConfig.MongoURI,
		RsName:       tConfig.MongoRsName,
	}
	db, err := local.NewMgo(mongoConfig, time.Minute)
	Expect(err).Should(BeNil())
	tables, err := db.ListTables(context.Background())
	Expect(err).Should(BeNil())
	for _, tableName := range tables {
		db.DropTable(context.Background(), tableName)
	}
	db.Close()
	clientSet.AdminServer().Migrate(context.Background(), "0", "community", GetHeader())
	clientSet.AdminServer().RunSyncDBIndex(context.Background(), GetHeader())
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

func GetDB() *local.Mongo {
	return db
}
