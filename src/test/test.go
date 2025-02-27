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

	"configcenter/pkg/tenant"
	"configcenter/src/apimachinery"
	"configcenter/src/apimachinery/adminserver"
	"configcenter/src/apimachinery/discovery"
	"configcenter/src/apimachinery/util"
	"configcenter/src/common"
	"configcenter/src/common/backbone/service_mange/zk"
	"configcenter/src/common/cryptor"
	headerutil "configcenter/src/common/http/header/util"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	kubetypes "configcenter/src/kube/types"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/dal/mongo"
	"configcenter/src/storage/dal/mongo/local"
	"configcenter/src/storage/dal/mongo/sharding"
	"configcenter/src/test/run"
	testutil "configcenter/src/test/util"

	. "github.com/onsi/gomega"
)

var clientSet apimachinery.ClientSetInterface
var adminClient adminserver.AdminServerClientInterface
var tConfig TestConfig
var reportUrl string
var reportDir string
var db dal.Dal

// TestConfig TODO
type TestConfig struct {
	ZkAddr         string
	Concurrent     int
	SustainSeconds float64
	TotalRequest   int64
	DBWriteKBSize  int
	MongoURI       string
	MongoRsName    string
	CryptoConf     string
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
	flag.StringVar(&tConfig.CryptoConf, "crypto-config", "", "mongodb crypto config in json format")
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

	cryptoConf := new(cryptor.Config)
	if tConfig.CryptoConf != "" {
		err = json.Unmarshal([]byte(tConfig.CryptoConf), cryptoConf)
		Expect(err).Should(BeNil())
	}
	crypto, err := cryptor.NewCrypto(cryptoConf)
	Expect(err).Should(BeNil())

	db, err = sharding.NewShardingMongo(mongoConfig, time.Minute, crypto)
	Expect(err).Should(BeNil())
	Expect(client.Start()).Should(BeNil())
	Expect(client.Ping()).Should(BeNil())
	disc, err := discovery.NewServiceDiscovery(client, "")
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

// TestTenantID is the tenant id for test
// TODO change to test tenant, right now use default tenant for compatible
const TestTenantID = common.BKDefaultTenantID

// GetHeader get header for test
func GetHeader() http.Header {
	return headerutil.GenCommonHeader(common.CCSystemOperatorUserName, TestTenantID, "")
}

// ClearDatabase TODO
func ClearDatabase() {
	fmt.Println("********Clear Database*************")
	for _, tableName := range common.PlatformTables() {
		err := db.Shard(sharding.NewShardOpts().WithIgnoreTenant()).DropTable(context.Background(), tableName)
		Expect(err).Should(BeNil())
	}

	err := tenant.ExecForAllTenants(func(tenantID string) error {
		shardOpts := sharding.NewShardOpts().WithTenant(tenantID)
		tables, err := db.Shard(shardOpts).ListTables(context.Background())
		if err != nil {
			return err
		}
		for _, tableName := range tables {
			err = db.Shard(shardOpts).DropTable(context.Background(), tableName)
			if err != nil {
				return err
			}
		}
		return nil
	})
	Expect(err).Should(BeNil())

	err = adminClient.Migrate(context.Background(), GetHeader())
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

// GetDal get dal
func GetDal() dal.Dal {
	return db
}

// GetDB get db client for test tenant
func GetDB() dal.DB {
	return db.Shard(sharding.NewShardOpts().WithTenant(TestTenantID))
}

// DeleteAllBizs delete all non-default bizs, used to clean biz data without ClearDatabase
func DeleteAllBizs() {
	ctx := context.Background()

	DeleteAllHosts()
	DeleteAllObjects()

	biz := make([]metadata.BizInst, 0)
	bizCond := mapstr.MapStr{common.BKAppNameField: mapstr.MapStr{common.BKDBNIN: []string{"资源池", "蓝鲸"}}}
	err := tenant.ExecForAllTenants(func(tenantID string) error {
		tenantBiz := make([]metadata.BizInst, 0)
		err := db.Shard(sharding.NewShardOpts().WithTenant(tenantID)).Table(common.BKTableNameBaseApp).Find(bizCond).Fields(common.BKAppIDField).
			All(ctx, &tenantBiz)
		if err != nil {
			return err
		}
		biz = append(biz, tenantBiz...)
		return nil
	})
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
		err = tenant.ExecForAllTenants(func(tenantID string) error {
			return db.Shard(sharding.NewShardOpts().WithTenant(tenantID)).Table(table).Delete(ctx,
				delCond)
		})
		Expect(err).NotTo(HaveOccurred())
	}
}

// DeleteAllHosts delete all hosts and their related data, used to clean host data without ClearDatabase
func DeleteAllHosts() {
	ctx := context.Background()

	tableNames := []string{common.BKTableNameBaseHost, common.BKTableNameModuleHostConfig,
		common.BKTableNameServiceInstance, common.BKTableNameProcessInstanceRelation, common.BKTableNameBaseProcess}

	for _, table := range tableNames {
		err := tenant.ExecForAllTenants(func(tenantID string) error {
			return db.Shard(sharding.NewShardOpts().WithTenant(tenantID)).Table(table).Delete(ctx, mapstr.New())
		})
		Expect(err).NotTo(HaveOccurred())
	}
}

// DeleteAllObjects delete all non-default objects, used to clean object data without ClearDatabase
func DeleteAllObjects() {
	ctx := context.Background()

	innerObjs := []string{common.BKInnerObjIDBizSet, common.BKInnerObjIDApp, common.BKInnerObjIDSet,
		common.BKInnerObjIDModule, common.BKInnerObjIDHost, common.BKInnerObjIDProc, common.BKInnerObjIDPlat,
		common.BKInnerObjIDProject}

	delCond := mapstr.MapStr{common.BKObjIDField: mapstr.MapStr{common.BKDBNIN: innerObjs}}
	objects := make([]metadata.Object, 0)

	err := tenant.ExecForAllTenants(func(tenantID string) error {
		shardOpts := sharding.NewShardOpts().WithTenant(tenantID)
		err := db.Shard(shardOpts).Table(common.BKTableNameObjDes).Find(delCond).
			Fields(common.BKObjIDField).All(ctx, &objects)
		if err != nil {
			return err
		}

		if len(objects) == 0 {
			return nil
		}

		objIDs := make([]string, len(objects))
		for i, obj := range objects {
			err = db.Shard(shardOpts).DropTable(ctx, common.GetInstTableName(obj.ObjectID, tenantID))
			if err != nil {
				return err
			}
			err = db.Shard(shardOpts).DropTable(ctx, common.GetObjectInstAsstTableName(obj.ObjectID, tenantID))
			if err != nil {
				return err
			}
			objIDs[i] = obj.ObjectID
		}

		objCond := mapstr.MapStr{common.BKObjIDField: mapstr.MapStr{common.BKDBIN: objIDs}}
		objTables := []string{common.BKTableNameObjDes, common.BKTableNameObjAttDes, common.BKTableNameObjUnique,
			common.BKTableNameObjectBaseMapping, common.BKTableNamePropertyGroup}
		for _, table := range objTables {
			err = db.Shard(shardOpts).Table(table).Delete(ctx, objCond)
			if err != nil {
				return err
			}
		}

		// compensate for mainline object association
		mainlineCond := mapstr.MapStr{
			common.AssociationKindIDField: common.AssociationKindMainline,
			common.BKAsstObjIDField:       common.BKInnerObjIDApp,
		}
		mainlineData := mapstr.MapStr{common.BKObjIDField: common.BKInnerObjIDSet}
		err = db.Shard(shardOpts).Table(common.BKTableNameObjAsst).Update(ctx, mainlineCond, mainlineData)
		if err != nil {
			return err
		}

		asstObjCond := mapstr.MapStr{common.BKAsstObjIDField: mapstr.MapStr{common.BKDBIN: objIDs}}
		err = db.Shard(shardOpts).Table(common.BKTableNameObjAsst).Delete(ctx,
			mapstr.MapStr{common.BKDBOR: []mapstr.MapStr{objCond, asstObjCond}})

		var idRuleNames []string
		for _, obj := range objIDs {
			idRuleNames = append(idRuleNames, "id_rule:incr_id:"+obj)
		}
		idGenerateCond := mapstr.MapStr{"_id": mapstr.MapStr{common.BKDBIN: idRuleNames}}
		err = db.Shard(sharding.NewShardOpts().WithIgnoreTenant()).Table(common.BKTableNameIDgenerator).Delete(ctx,
			idGenerateCond)

		if err != nil {
			return err
		}
		return nil
	})
	Expect(err).NotTo(HaveOccurred())
}

// GetCloudID get any one cloud id
func GetCloudID() int64 {
	defaultPlat := new(metadata.CloudArea)
	cond := make(map[string]interface{})
	err := GetDB().Table(common.BKTableNameBasePlat).Find(cond).Fields(common.BKCloudIDField).One(
		context.Background(), defaultPlat)
	Expect(err).NotTo(HaveOccurred())

	return defaultPlat.CloudID
}

// GetDefaultCategory get default service category
func GetDefaultCategory() int64 {
	subDefaultCategory := new(metadata.ServiceCategory)
	cond := map[string]interface{}{
		common.BKFieldName:     common.DefaultServiceCategoryName,
		common.BKParentIDField: mapstr.MapStr{common.BKDBNE: 0},
	}
	err := GetDB().Table(common.BKTableNameServiceCategory).Find(cond).One(context.Background(), subDefaultCategory)
	Expect(err).NotTo(HaveOccurred())

	return subDefaultCategory.ID
}

// GetResBizID get resource pool biz id
func GetResBizID() int64 {
	biz := new(metadata.BizInst)
	err := GetDB().Table(common.BKTableNameBaseApp).Find(map[string]interface{}{
		common.BKAppNameField: common.DefaultAppName}).One(context.Background(), biz)
	Expect(err).NotTo(HaveOccurred())

	return biz.BizID
}
