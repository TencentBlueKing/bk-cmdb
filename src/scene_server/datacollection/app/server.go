/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package app

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"sync"
	"time"

	"configcenter/src/ac/extensions"
	"configcenter/src/ac/iam"
	apiutil "configcenter/src/apimachinery/util"
	"configcenter/src/common"
	"configcenter/src/common/auth"
	"configcenter/src/common/backbone"
	cc "configcenter/src/common/backbone/configcenter"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/metadata"
	"configcenter/src/common/types"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/datacollection/app/options"
	"configcenter/src/scene_server/datacollection/collections"
	"configcenter/src/scene_server/datacollection/collections/hostsnap"
	"configcenter/src/scene_server/datacollection/collections/middleware"
	"configcenter/src/scene_server/datacollection/collections/netcollect"
	svc "configcenter/src/scene_server/datacollection/service"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/dal/kafka"
	"configcenter/src/storage/dal/mongo"
	"configcenter/src/storage/dal/mongo/local"
	"configcenter/src/storage/dal/redis"
	"configcenter/src/thirdparty/esbserver"
	"configcenter/src/thirdparty/esbserver/esbutil"

	"github.com/Shopify/sarama"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	// ESBPrefix is prefix of configs variable for ESB.
	ESBPrefix = "esb"

	// snapPorterName is name of snap porter.
	snapPorterName = "hostsnap"

	// middlewarePorterName is name of middleware porter.
	middlewarePorterName = "middleware"

	// netCollectPorterName is name of netcollect porter.
	netCollectPorterName = "netcollect"

	// defaultInitWaitDuration is default duration for new DataCollection init.
	defaultInitWaitDuration = time.Second

	// defaultDBConnectTimeout is default connect timeout of cc db.
	defaultDBConnectTimeout = 5 * time.Second

	// defaultAppInitWaitDuration is default wait duration for app db init.
	defaultAppInitWaitDuration = 10 * time.Second
)

// DataCollectionConfig is configs for DataCollection app.
type DataCollectionConfig struct {
	// MongoDB mongodb configs.
	MongoDB mongo.Config

	// CCRedis CC main redis configs.
	CCRedis redis.Config

	// SnapRedis snap redis configs.
	SnapRedis redis.Config

	// DiscoverRedis discover redis configs.
	DiscoverRedis redis.Config

	// NetCollectRedis net collection redis configs.
	NetCollectRedis redis.Config

	// SnapKafka snap kafka configs.
	SnapKafka kafka.Config

	// ESB blueking ESB configs.
	Esb esbutil.EsbConfig

	// DefaultAppName default name of this app.
	DefaultAppName string

	// Auth is auth config
	Auth iam.AuthConfig

	// SnapReportMode hostsnap report mode
	SnapReportMode string
}

// DataCollection is data collection server.
type DataCollection struct {
	ctx    context.Context
	engine *backbone.Engine

	defaultAppID    string
	snapshotBizName string

	// config for this DataCollection app.
	config *DataCollectionConfig

	// service main service instance.
	service *svc.Service

	// make host configs update action safe.
	hostConfigUpdateMu sync.Mutex

	// db is cc main database.
	db dal.RDB

	// redisCli is cc main cache redis client.
	redisCli redis.Client

	// snapRedisCli is snap redis client.
	snapRedisCli redis.Client

	// disRedisCli is discover redis client.
	disRedisCli redis.Client

	// netRedisCli is net collect redis client.
	netRedisCli redis.Client

	// snapConsumerGroup is snap kafka consumer group client.
	snapConsumerGroup sarama.ConsumerGroup

	// authManager is auth manager.
	authManager *extensions.AuthManager

	// porterManager is porters manager.
	porterManager *collections.PorterManager

	// registry is prometheus registry.
	registry prometheus.Registerer

	// hash collections hash object, that updates target nodes in dynamic mode,
	// and calculates node base on hash key of data.
	hash *collections.Hash
}

// NewDataCollection creates a new DataCollection object.
func NewDataCollection(ctx context.Context, op *options.ServerOption) (*DataCollection, error) {
	// build server info.
	svrInfo, err := types.NewServerInfo(op.ServConf)
	if err != nil {
		return nil, fmt.Errorf("build server info, %+v", err)
	}

	// new DataCollection instance.
	newDataCollection := &DataCollection{ctx: ctx}

	engine, err := backbone.NewBackbone(ctx, &backbone.BackboneParameter{
		ConfigUpdate: newDataCollection.OnHostConfigUpdate,
		ConfigPath:   op.ServConf.ExConfig,
		Regdiscv:     op.ServConf.RegDiscover,
		SrvInfo:      svrInfo,
	})
	if err != nil {
		return nil, fmt.Errorf("build backbone, %+v", err)
	}

	// set global cc errors.
	errors.SetGlobalCCError(engine.CCErr)

	// set hash.
	newDataCollection.hash = collections.NewHash(svrInfo.RegisterIP, svrInfo.Port, engine.Discovery())

	// set backbone engine.
	newDataCollection.engine = engine
	newDataCollection.service = svc.NewService(ctx, engine)
	newDataCollection.registry = engine.Metric().Registry()

	return newDataCollection, nil
}

// Engine returns engine of the DataCollection instance.
func (c *DataCollection) Engine() *backbone.Engine {
	return c.engine
}

// Service returns main service of the DataCollection instance.
func (c *DataCollection) Service() *svc.Service {
	return c.service
}

// OnHostConfigUpdate is callback for updating configs.
func (c *DataCollection) OnHostConfigUpdate(prev, curr cc.ProcessConfig) {
	c.hostConfigUpdateMu.Lock()
	defer c.hostConfigUpdateMu.Unlock()

	if len(curr.ConfigData) > 0 {
		// NOTE: allow to update configs with empty values?
		// NOTE: what is prev used for? build a compare logic here?

		if c.config == nil {
			c.config = &DataCollectionConfig{}
		}

		blog.V(3).Infof("DataCollection| on host config update event: \n%s", string(curr.ConfigData))
		// ESB configs.
		c.config.Esb.Addrs, _ = cc.String("esb.addr")
		c.config.Esb.AppCode, _ = cc.String("esb.appCode")
		c.config.Esb.AppSecret, _ = cc.String("esb.appSecret")
	}
}

// initConfigs inits configs for new DataCollection server.
func (c *DataCollection) initConfigs() error {
	for {
		// wait and parse configs that async updated by backbone engine.
		c.hostConfigUpdateMu.Lock()
		if c.config == nil {
			c.hostConfigUpdateMu.Unlock()

			blog.Info("DataCollection| can't find configs to run the new datacollection server, try again later!")
			time.Sleep(defaultInitWaitDuration)
			continue
		}

		// ready to init new datacollection instance.
		c.hostConfigUpdateMu.Unlock()
		break
	}

	var err error
	blog.Info("DataCollection| found configs to run the new datacollection server now!")

	// set snapshot biz name
	if err := c.setSnapshotBizName(); err != nil {
		return err
	}

	// mongodb.
	c.config.MongoDB, err = c.engine.WithMongo()
	if err != nil {
		return fmt.Errorf("init mongodb configs, %+v", err)
	}

	// cc main redis.
	c.config.CCRedis, err = c.engine.WithRedis()
	if err != nil {
		return fmt.Errorf("init cc redis configs, %+v", err)
	}

	// snap redis.
	c.config.SnapRedis, err = c.engine.WithRedis("redis.snap")
	if err != nil {
		return fmt.Errorf("init snap redis configs, %+v", err)
	}

	// discover redis.
	c.config.DiscoverRedis, err = c.engine.WithRedis("redis.discover")
	if err != nil {
		return fmt.Errorf("init discover redis configs, %+v", err)
	}

	// netcollect redis.
	c.config.NetCollectRedis, err = c.engine.WithRedis("redis.netcollect")
	if err != nil {
		return fmt.Errorf("init netcollect redis configs, %+v", err)
	}

	c.config.SnapReportMode, _ = cc.String("datacollection.hostsnap.reportMode")
	if metadata.GseConfigReportMode(c.config.SnapReportMode) == metadata.GseConfigReportModeKafka {
		c.config.SnapKafka, _ = cc.Kafka("kafka.snap")
		if err := c.config.SnapKafka.Check(); err != nil {
			blog.Errorf("kafka config is error, err: %v", err)
			return err
		}
	}

	c.config.Auth, err = iam.ParseConfigFromKV("authServer", nil)
	if err != nil {
		blog.Warnf("parse auth center config failed: %v", err)
	}

	return nil
}

// initModules inits modules for new DataCollection server.
func (c *DataCollection) initModules() error {
	// create mongodb client.
	mgoCli, err := local.NewMgo(c.config.MongoDB.GetMongoConf(), defaultDBConnectTimeout)
	if err != nil {
		return fmt.Errorf("create new mongodb client, %+v", err)
	}
	c.db = mgoCli
	c.service.SetDB(mgoCli)
	blog.Info("DataCollection| init modules, create mongo client success[%+v]", c.config.MongoDB.GetMongoConf())

	// create blueking ESB client.
	tlsConfig, err := apiutil.NewTLSClientConfigFromConfig("esb")
	if err != nil {
		return err
	}
	apiMachineryConfig := &apiutil.APIMachineryConfig{
		QPS:       1000,
		Burst:     1000,
		TLSConfig: &tlsConfig,
	}

	esb, err := esbserver.NewEsb(apiMachineryConfig, nil, /* you can update it by a chan here */
		&c.config.Esb, c.engine.Metric().Registry())
	if err != nil {
		return fmt.Errorf("create ESB client, %+v", err)
	}
	blog.Info("DataCollection| init modules, create ESB success[%+v]", c.config.Esb)

	// build logics comm.
	c.service.SetLogics(mgoCli, esb)

	// connect to cc main redis.
	redisCli, err := redis.NewFromConfig(c.config.CCRedis)
	if err != nil {
		return fmt.Errorf("connect to cc main redis, %+v", err)
	}
	c.redisCli = redisCli
	c.service.SetCache(redisCli)
	blog.Infof("DataCollection| init modules, connected to cc main redis, %+v", c.config.CCRedis)

	// connect to snap redis.
	if c.config.SnapRedis.Enable != "false" {
		snapCli, err := redis.NewFromConfig(c.config.SnapRedis)
		if err != nil {
			return fmt.Errorf("connect to snap redis, %+v", err)
		}
		c.snapRedisCli = snapCli
		c.service.SetSnapCli(snapCli)
		blog.Infof("DataCollection| init modules, connected to snap redis, %+v", c.config.SnapRedis)
	}

	// connect to discover redis.
	if c.config.DiscoverRedis.Enable != "false" {
		disCli, err := redis.NewFromConfig(c.config.DiscoverRedis)
		if nil != err {
			return fmt.Errorf("connect to discover redis, %+v", err)
		}
		c.disRedisCli = disCli
		c.service.SetDiscoverCli(disCli)
		blog.Infof("DataCollection| init modules, connected to discover redis, %+v", c.config.DiscoverRedis)
	}

	// connect to net collect redis.
	if c.config.NetCollectRedis.Enable != "false" {
		netCli, err := redis.NewFromConfig(c.config.NetCollectRedis)
		if nil != err {
			return fmt.Errorf("connect to netcollect redis, %+v", err)
		}
		c.netRedisCli = netCli
		c.service.SetNetCollectCli(netCli)
		blog.Infof("DataCollection| init modules, connected to netcollect redis, %+v", c.config.NetCollectRedis)
	}

	// connect to snap kafka.
	if c.config.SnapReportMode == "kafka" {
		config := sarama.NewConfig()
		config.Consumer.Return.Errors = false
		config.Consumer.Offsets.AutoCommit.Enable = false // 禁用自动提交，改为手动
		config.Consumer.Offsets.Initial = sarama.OffsetOldest
		if c.config.SnapKafka.User != "" && c.config.SnapKafka.Password != "" {
			config.Net.SASL.Enable = true
			config.Net.SASL.User = c.config.SnapKafka.User
			config.Net.SASL.Password = c.config.SnapKafka.Password
			config.Net.SASL.Handshake = true
			config.Net.SASL.SCRAMClientGeneratorFunc = func() sarama.SCRAMClient {
				return &kafka.XDGSCRAMClient{HashGeneratorFcn: kafka.SHA512}
			}
			config.Net.SASL.Mechanism = sarama.SASLTypeSCRAMSHA512
		}
		var err error
		c.snapConsumerGroup, err = sarama.NewConsumerGroup(c.config.SnapKafka.Brokers, c.config.SnapKafka.GroupID,
			config)
		if err != nil {
			blog.Errorf("create kafka consumer group error, err: %v", err)
			return err
		}
	}

	iamCli := new(iam.IAM)
	if auth.EnableAuthorize() {
		blog.Info("enable auth center access")
		iamCli, err = iam.NewIAM(c.config.Auth, c.engine.Metric().Registry())
		if err != nil {
			return fmt.Errorf("new iam client failed: %v", err)
		}
	} else {
		blog.Infof("disable auth center access")
	}
	c.authManager = extensions.NewAuthManager(c.engine.CoreAPI, iamCli)

	return nil
}

// getDefaultAppID returns default appid of this DataCollection server.
func (c *DataCollection) getDefaultAppID() (string, error) {
	// query condition.
	condition := map[string]interface{}{common.BKAppNameField: c.snapshotBizName}

	// query results.
	results := []map[string]interface{}{}

	// query appid from cc db.
	if err := c.db.Table(common.BKTableNameBaseApp).Find(condition).All(c.ctx, &results); err != nil {
		return "", err
	}
	if len(results) <= 0 {
		return "", fmt.Errorf("target app not found")
	}

	id, err := util.GetInt64ByInterface(results[0][common.BKAppIDField])
	if err != nil {
		return "", fmt.Errorf("can't query default appid, unkonw id type, %+v",
			reflect.TypeOf(results[0][common.BKAppIDField]))
	}
	defaultAppID := strconv.FormatInt(id, 10)

	return defaultAppID, nil
}

func (c *DataCollection) snapMessageTopic(defaultAppID string) []string {
	return []string{
		// current snap topic name.
		fmt.Sprintf("snapshot%s", defaultAppID),

		// old snap topic name, just for compatibility.
		fmt.Sprintf("%s_snapshot", defaultAppID),
	}
}

func (c *DataCollection) discoverMessageTopic(defaultAppID string) []string {
	return []string{
		// current discover topic name.
		fmt.Sprintf("discover%s", defaultAppID),
	}
}

func (c *DataCollection) netcollectMessageTopic(defaultAppID string) []string {
	return []string{
		// current netcollect topic name.
		fmt.Sprintf("netdevice%s", defaultAppID),
	}
}

// runCollectPorters runs porters for collections.
func (c *DataCollection) runCollectPorters() {
	// create porters manager.
	c.porterManager = collections.NewPorterManager()
	go c.porterManager.Run()

	// default appid.
	for {
		defaultAppID, err := c.getDefaultAppID()
		if err == nil {
			// success.
			c.defaultAppID = defaultAppID
			break
		}

		blog.Errorf("DataCollection| get default appid failed: %+v, init database first and it would "+
			"try again in %+v seconds later", err, defaultAppInitWaitDuration)
		time.Sleep(defaultAppInitWaitDuration)
	}
	blog.Info("DataCollection| get default appid id success[%s]", c.defaultAppID)

	// create and add new porters.
	if metadata.GseConfigReportMode(c.config.SnapReportMode) == metadata.GseConfigReportModeKafka {
		topic := c.snapMessageTopic(c.defaultAppID)
		analyzer := hostsnap.NewHostSnap(c.ctx, c.redisCli, c.db, c.engine, c.authManager)

		kafkaPorter := collections.NewKafkaPorter(snapPorterName, c.engine, c.ctx, analyzer, c.snapConsumerGroup,
			topic, c.registry)
		c.porterManager.AddPorter(kafkaPorter)
		blog.Info("DataCollection| create kafka hostsnap analyzer with target porter[%s] on topic[%s] success",
			snapPorterName, topic)
	}

	if c.snapRedisCli != nil {
		topic := c.snapMessageTopic(c.defaultAppID)
		analyzer := hostsnap.NewHostSnap(c.ctx, c.redisCli, c.db, c.engine, c.authManager)

		porter := collections.NewSimplePorter(snapPorterName, c.engine, c.hash, analyzer, c.snapRedisCli, topic,
			c.registry)
		c.porterManager.AddPorter(porter)
		blog.Info("DataCollection| create redis hostsnap analyzer with target porter[%s] on topic[%s] success",
			snapPorterName, topic)
	}

	if c.disRedisCli != nil {
		topic := c.discoverMessageTopic(c.defaultAppID)
		analyzer := middleware.NewDiscover(c.ctx, c.redisCli, c.engine, c.authManager)

		porter := collections.NewSimplePorter(middlewarePorterName, c.engine, c.hash, analyzer, c.disRedisCli, topic,
			c.registry)
		c.porterManager.AddPorter(porter)
		blog.Info("DataCollection| create kafka discover analyzer with target porter[%s] on topic[%s] success",
			middlewarePorterName, topic)
	}

	if c.netRedisCli != nil {
		topic := c.netcollectMessageTopic(c.defaultAppID)
		analyzer := netcollect.NewNetCollect(c.ctx, c.db, c.authManager)

		porter := collections.NewSimplePorter(netCollectPorterName, c.engine, c.hash, analyzer, c.netRedisCli, topic,
			c.registry)
		c.porterManager.AddPorter(porter)
		blog.Info("DataCollection| create redis netcollect analyzer with target porter[%s] on topic[%s] success",
			netCollectPorterName, topic)
	}
}

// Run runs a new datacollection server.
func (c *DataCollection) Run() error {
	// init configs.
	if err := c.initConfigs(); err != nil {
		return err
	}
	blog.Info("init configs success!")

	// ready to setup comms for new server instance now.
	if err := c.initModules(); err != nil {
		return err
	}
	blog.Info("init modules success!")

	// run collection porters for new datacollection instance.
	c.runCollectPorters()

	blog.Info("run collect porters success!")

	return nil
}

// Run setups a new datacollection app with a context and options and runs it as server instance.
func Run(ctx context.Context, cancel context.CancelFunc, op *options.ServerOption) error {
	// create datacollection server.
	dataCollection, err := NewDataCollection(ctx, op)
	if err != nil {
		return fmt.Errorf("create new datacollection server, %+v", err)
	}

	if err := dataCollection.Run(); err != nil {
		return err
	}

	// all modules is inited success, start the new server now.
	if err := backbone.StartServer(ctx, cancel, dataCollection.Engine(),
		dataCollection.Service().WebService(), true); err != nil {
		return err
	}
	blog.Info("DataCollection init and run success!")

	<-ctx.Done()
	blog.Info("DataCollection stopping now!")
	return nil
}

func (c *DataCollection) setSnapshotBizName() error {
	tryCnt := 30
	header := util.BuildHeader(common.CCSystemOperatorUserName, common.BKDefaultOwnerID)
	for i := 1; i <= tryCnt; i++ {
		time.Sleep(time.Second * 2)
		res, err := c.engine.CoreAPI.CoreService().System().SearchPlatformSetting(context.Background(), header)
		if err != nil {
			blog.Warnf("setSnapshotBizName failed,  try count:%d, SearchConfigAdmin err: %v", i, err)
			continue
		}
		if res.Result == false {
			blog.Warnf("setSnapshotBizName failed,  try count:%d, SearchConfigAdmin err: %s", i, res.ErrMsg)
			continue
		}
		c.snapshotBizName = res.Data.Backend.SnapshotBizName
		break
	}

	if c.snapshotBizName == "" {
		blog.Errorf("setSnapshotBizName failed, SnapshotBizName is empty, check the coreservice and the value " +
			"in table cc_System")
		return fmt.Errorf("setSnapshotBizName failed")
	}

	blog.Infof("setSnapshotBizName successfully, SnapshotBizName is %s", c.snapshotBizName)
	return nil
}
