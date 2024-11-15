// Package test TODO
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
	"configcenter/src/apimachinery/adminserver"
	"configcenter/src/apimachinery/discovery"
	"configcenter/src/apimachinery/util"
	"configcenter/src/common"
	"configcenter/src/common/backbone/service_mange/zk"
	headerutil "configcenter/src/common/http/header/util"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	kubetypes "configcenter/src/kube/types"
	"configcenter/src/storage/dal/mongo"
	"configcenter/src/storage/dal/mongo/local"
	"configcenter/src/test/run"
	testutil "configcenter/src/test/util"

	. "github.com/onsi/gomega"
)

var clientSet apimachinery.ClientSetInterface
var adminClient adminserver.AdminServerClientInterface
var tConfig TestConfig
var reportUrl string
var reportDir string
var db *local.Mongo

// TestConfig TODO
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

// RedisConfig TODO
type RedisConfig struct {
	RedisAddress string
	RedisPort    string
	RedisPasswd  string
}

func init() {
	testing.Init()
	flag.StringVar(&tConfig.ZkAddr, "zk-addr", "127.0.0.1:2181", "zk discovery addresses, comma separated.")
	flag.IntVar(&tConfig.Concurrent, "concurrent", 100, "concurrent request during the load ")
	flag.Float64Var(&tConfig.SustainSeconds, "sustain-seconds", 10, "the load test sustain time in seconds ")
	flag.Int64Var(&tConfig.TotalRequest, "total-request", 0,
		"the load test total request,it has higher priority than SustainSeconds")
	flag.IntVar(&tConfig.DBWriteKBSize, "write-size", 1, "MongoDB write size , unit is KB.")
	flag.StringVar(&tConfig.RedisCfg.RedisAddress, "redis-addr", "127.0.0.1:6379", "redis host address with port")
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

	// initialize admin-server client set, because migrate needs a longer timeout
	adminApiConfig := &util.APIMachineryConfig{
		QPS:       20000,
		Burst:     10000,
		ExtraConf: &util.ExtraClientConfig{ResponseHeaderTimeout: 5 * time.Minute},
	}
	adminClientSet, err := apimachinery.NewApiMachinery(adminApiConfig, disc)
	Expect(err).Should(BeNil())
	adminClient = adminClientSet.AdminServer()
	// wait for get the apiserver address.
	time.Sleep(1 * time.Second)
	fmt.Println("**** initialize clientSet success ***")
}

// GetClientSet TODO
func GetClientSet() apimachinery.ClientSetInterface {
	return clientSet
}

// GetTestConfig TODO
func GetTestConfig() TestConfig {
	return tConfig
}

// GetHeader TODO
func GetHeader() http.Header {
	return headerutil.GenDefaultHeader()
}

// ClearDatabase TODO
func ClearDatabase() {
	fmt.Println("********Clear Database*************")

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
		err = db.DropTable(context.Background(), tableName)
		Expect(err).Should(BeNil())
	}
	_ = db.Close()

	err = adminClient.Migrate(context.Background(), "0", "community", GetHeader())
	Expect(err).Should(BeNil())
	err = adminClient.RunSyncDBIndex(context.Background(), GetHeader())
	Expect(err).Should(BeNil())
}

// GetReportUrl TODO
func GetReportUrl() string {
	if !strings.HasSuffix(reportUrl, "/") {
		reportUrl = reportUrl + "/"
	}
	return reportUrl
}

// GetReportDir TODO
func GetReportDir() string {
	if !strings.HasSuffix(reportDir, "/") {
		reportDir = reportDir + "/"
	}
	return reportDir
}

// GetDB TODO
func GetDB() *local.Mongo {
	return db
}

// DeleteAllBizs delete all non-default bizs, used to clean biz data without ClearDatabase
func DeleteAllBizs() {
	ctx := context.Background()

	DeleteAllHosts()
	DeleteAllObjects()

	biz := make([]metadata.BizInst, 0)
	bizCond := mapstr.MapStr{common.BKAppNameField: mapstr.MapStr{common.BKDBNIN: []string{"资源池", "蓝鲸"}}}
	err := GetDB().Table(common.BKTableNameBaseApp).Find(bizCond).Fields(common.BKAppIDField).All(ctx, &biz)
	Expect(err).NotTo(HaveOccurred())

	if len(biz) == 0 {
		return
	}

	bizIDs := make([]int64, len(biz))
	for i, b := range biz {
		bizIDs[i] = b.BizID
	}

	tableNames := []string{common.BKTableNameBaseSet, common.BKTableNameBaseModule, common.BKTableNameSetTemplate,
		common.BKTableNameSetTemplateAttr, common.BKTableNameServiceTemplate, common.BKTableNameServiceTemplateAttr,
		common.BKTableNameProcessTemplate, common.BKTableNameServiceCategory, common.BKTableNamePropertyGroup,
		common.BKTableNameObjAttDes, kubetypes.BKTableNameBaseCluster, kubetypes.BKTableNameBaseNode,
		kubetypes.BKTableNameBaseNamespace, kubetypes.BKTableNameBaseWorkload, kubetypes.BKTableNameBaseDeployment,
		kubetypes.BKTableNameBaseStatefulSet, kubetypes.BKTableNameBaseDaemonSet, kubetypes.BKTableNameGameDeployment,
		kubetypes.BKTableNameGameStatefulSet, kubetypes.BKTableNameBaseCronJob, kubetypes.BKTableNameBaseJob,
		kubetypes.BKTableNameBasePodWorkload, kubetypes.BKTableNameBaseCustom, kubetypes.BKTableNameBasePod,
		kubetypes.BKTableNameBaseContainer, common.BKTableNameBaseApp}

	delCond := mapstr.MapStr{common.BKAppIDField: mapstr.MapStr{common.BKDBIN: bizIDs}}
	for _, table := range tableNames {
		err = GetDB().Table(table).Delete(ctx, delCond)
		Expect(err).NotTo(HaveOccurred())
	}
}

// DeleteAllHosts delete all hosts and their related data, used to clean host data without ClearDatabase
func DeleteAllHosts() {
	ctx := context.Background()

	tableNames := []string{common.BKTableNameBaseHost, common.BKTableNameModuleHostConfig,
		common.BKTableNameServiceInstance, common.BKTableNameProcessInstanceRelation, common.BKTableNameBaseProcess}

	for _, table := range tableNames {
		err := GetDB().Table(table).Delete(ctx, mapstr.New())
		Expect(err).NotTo(HaveOccurred())
	}
}

// DeleteAllObjects delete all non-default objects, used to clean object data without ClearDatabase
func DeleteAllObjects() {
	ctx := context.Background()

	innerObjs := []string{common.BKInnerObjIDBizSet, common.BKInnerObjIDApp, common.BKInnerObjIDSet,
		common.BKInnerObjIDModule, common.BKInnerObjIDHost, common.BKInnerObjIDProc, common.BKInnerObjIDPlat,
		common.BKInnerObjIDProject, common.BKInnerObjIDSwitch, common.BKInnerObjIDRouter, common.BKInnerObjIDBlance,
		common.BKInnerObjIDFirewall, common.BKInnerObjIDWeblogic, common.BKInnerObjIDTomcat, common.BKInnerObjIDApache}

	delCond := mapstr.MapStr{common.BKObjIDField: mapstr.MapStr{common.BKDBNIN: innerObjs}}
	objects := make([]metadata.Object, 0)
	err := GetDB().Table(common.BKTableNameObjDes).Find(delCond).Fields(common.BKObjIDField, common.TenantID).
		All(ctx, &objects)
	Expect(err).NotTo(HaveOccurred())

	if len(objects) == 0 {
		return
	}

	objIDs := make([]string, len(objects))
	for i, obj := range objects {
		err = db.DropTable(ctx, common.GetInstTableName(obj.ObjectID, obj.TenantID))
		Expect(err).NotTo(HaveOccurred())
		err = db.DropTable(ctx, common.GetObjectInstAsstTableName(obj.ObjectID, obj.TenantID))
		Expect(err).NotTo(HaveOccurred())
		objIDs[i] = obj.ObjectID
	}

	objCond := mapstr.MapStr{common.BKObjIDField: mapstr.MapStr{common.BKDBIN: objIDs}}
	objTables := []string{common.BKTableNameObjDes, common.BKTableNameObjAttDes, common.BKTableNameObjUnique,
		"cc_ObjectBaseMapping"}
	for _, table := range objTables {
		err = db.Table(table).Delete(ctx, objCond)
		Expect(err).NotTo(HaveOccurred())
	}

	// compensate for mainline object association
	mainlineCond := mapstr.MapStr{
		common.AssociationKindIDField: common.AssociationKindMainline,
		common.BKAsstObjIDField:       common.BKInnerObjIDApp,
	}
	mainlineData := mapstr.MapStr{common.BKObjIDField: common.BKInnerObjIDSet}
	err = db.Table(common.BKTableNameObjAsst).Update(ctx, mainlineCond, mainlineData)
	Expect(err).NotTo(HaveOccurred())

	asstObjCond := mapstr.MapStr{common.BKAsstObjIDField: mapstr.MapStr{common.BKDBIN: objIDs}}
	err = db.Table(common.BKTableNameObjAsst).Delete(ctx,
		mapstr.MapStr{common.BKDBOR: []mapstr.MapStr{objCond, asstObjCond}})
	Expect(err).NotTo(HaveOccurred())
}
