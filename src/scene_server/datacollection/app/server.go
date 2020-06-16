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
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"sync"
	"time"

	"configcenter/src/auth"
	"configcenter/src/auth/authcenter"
	"configcenter/src/auth/extensions"
	"configcenter/src/common"
	enableauth "configcenter/src/common/auth"
	"configcenter/src/common/backbone"
	cc "configcenter/src/common/backbone/configcenter"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/types"
	"configcenter/src/scene_server/datacollection/app/options"
	"configcenter/src/scene_server/datacollection/collections"
	"configcenter/src/scene_server/datacollection/collections/hostsnap"
	"configcenter/src/scene_server/datacollection/collections/middleware"
	"configcenter/src/scene_server/datacollection/collections/netcollect"
	svc "configcenter/src/scene_server/datacollection/service"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/dal/mongo"
	"configcenter/src/storage/dal/mongo/local"
	dalredis "configcenter/src/storage/dal/redis"
	"configcenter/src/thirdpartyclient/esbserver"
	"configcenter/src/thirdpartyclient/esbserver/esbutil"

	"github.com/prometheus/client_golang/prometheus"
	"gopkg.in/redis.v5"
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
	CCRedis dalredis.Config

	// SnapRedis snap redis configs.
	SnapRedis dalredis.Config

	// DiscoverRedis discover redis configs.
	DiscoverRedis dalredis.Config

	// NetCollectRedis net collection redis configs.
	NetCollectRedis dalredis.Config

	// ESB blueking ESB configs.
	Esb esbutil.EsbConfig

	// AuthConfig auth configs.
	AuthConfig authcenter.AuthConfig

	// DefaultAppName default name of this app.
	DefaultAppName string
}

// DataCollection is data collection server.
type DataCollection struct {
	ctx    context.Context
	engine *backbone.Engine

	defaultAppID   string
	defaultAppName string

	// config for this DataCollection app.
	config *DataCollectionConfig

	// service main service instance.
	service *svc.Service

	// make host configs update action safe.
	hostConfigUpdateMu sync.Mutex

	// db is cc main database.
	db dal.RDB

	// redisCli is cc main cache redis client.
	redisCli *redis.Client

	// snapCli is snap redis client.
	snapCli *redis.Client

	// disCli is discover redis client.
	disCli *redis.Client

	// netCli is net collect redis client.
	netCli *redis.Client

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

	if len(curr.ConfigMap) > 0 {
		// NOTE: allow to update configs with empty values?
		// NOTE: what is prev used for? build a compare logic here?

		if c.config == nil {
			c.config = &DataCollectionConfig{}
		}

		if data, err := json.MarshalIndent(curr.ConfigMap, "", "  "); err == nil {
			blog.V(3).Infof("DataCollection| on host config update event: \n%s", data)
		}

		// ESB configs.
		c.config.Esb.Addrs = curr.ConfigMap["esb.addr"]
		c.config.Esb.AppCode = curr.ConfigMap["esb.appCode"]
		c.config.Esb.AppSecret = curr.ConfigMap["esb.appSecret"]

		// default app name.
		c.config.DefaultAppName = curr.ConfigMap["biz.default_app_name"]
		if len(c.config.DefaultAppName) == 0 {
			c.config.DefaultAppName = common.BKAppName
		}
		c.defaultAppName = c.config.DefaultAppName
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
	c.config.SnapRedis, err = c.engine.WithRedis("snap-redis")
	if err != nil {
		return fmt.Errorf("init snap redis configs, %+v", err)
	}

	// discover redis.
	c.config.DiscoverRedis, err = c.engine.WithRedis("discover-redis")
	if err != nil {
		return fmt.Errorf("init discover redis configs, %+v", err)
	}

	// netcollect redis.
	c.config.NetCollectRedis, err = c.engine.WithRedis("netcollect-redis")
	if err != nil {
		return fmt.Errorf("init netcollect redis configs, %+v", err)
	}

	// authorization.
	c.config.AuthConfig, err = c.engine.WithAuth()
	if err != nil {
		return fmt.Errorf("init authorization configs, %+v", err)
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
	esb, err := esbserver.NewEsb(c.engine.ApiMachineryConfig(), nil, /* you can update it by a chan here */
		&c.config.Esb, c.engine.Metric().Registry())
	if err != nil {
		return fmt.Errorf("create ESB client, %+v", err)
	}
	blog.Info("DataCollection| init modules, create ESB success[%+v]", c.config.Esb)

	// build logics comm.
	c.service.SetLogics(mgoCli, esb)

	// connect to cc main redis.
	redisCli, err := dalredis.NewFromConfig(c.config.CCRedis)
	if err != nil {
		return fmt.Errorf("connect to cc main redis, %+v", err)
	}
	c.redisCli = redisCli
	c.service.SetCache(redisCli)
	blog.Infof("DataCollection| init modules, connected to cc main redis, %+v", c.config.CCRedis)

	// connect to snap redis.
	if c.config.SnapRedis.Enable != "false" {
		snapCli, err := dalredis.NewFromConfig(c.config.SnapRedis)
		if err != nil {
			return fmt.Errorf("connect to snap redis, %+v", err)
		}
		c.snapCli = snapCli
		c.service.SetSnapCli(snapCli)
		blog.Infof("DataCollection| init modules, connected to snap redis, %+v", c.config.SnapRedis)
	}

	// connect to discover redis.
	if c.config.DiscoverRedis.Enable != "false" {
		disCli, err := dalredis.NewFromConfig(c.config.DiscoverRedis)
		if nil != err {
			return fmt.Errorf("connect to discover redis, %+v", err)
		}
		c.disCli = disCli
		c.service.SetDiscoverCli(disCli)
		blog.Infof("DataCollection| init modules, connected to discover redis, %+v", c.config.DiscoverRedis)
	}

	// connect to net collect redis.
	if c.config.NetCollectRedis.Enable != "false" {
		netCli, err := dalredis.NewFromConfig(c.config.NetCollectRedis)
		if nil != err {
			return fmt.Errorf("connect to netcollect redis, %+v", err)
		}
		c.netCli = netCli
		c.service.SetNetCollectCli(netCli)
		blog.Infof("DataCollection| init modules, connected to netcollect redis, %+v", c.config.NetCollectRedis)
	}

	// handle authorize.
	if enableauth.IsAuthed() {
		authorize, err := auth.NewAuthorize(nil, c.config.AuthConfig, c.engine.Metric().Registry())
		if err != nil {
			return fmt.Errorf("create new authorize failed, %+v", err)
		}
		c.authManager = extensions.NewAuthManager(c.engine.CoreAPI, authorize)
		blog.Infof("DataCollection| init modules, create authorize success[%+v]", c.config.AuthConfig)
	}

	return nil
}

// getDefaultAppID returns default appid of this DataCollection server.
func (c *DataCollection) getDefaultAppID() (string, error) {
	// query condition.
	condition := map[string]interface{}{common.BKAppNameField: c.defaultAppName}

	// query results.
	results := []map[string]interface{}{}

	// query appid from cc db.
	if err := c.db.Table(common.BKTableNameBaseApp).Find(condition).All(c.ctx, &results); err != nil {
		return "", err
	}
	if len(results) <= 0 {
		return "", fmt.Errorf("target app not found")
	}

	defaultAppID := ""

	switch id := results[0][common.BKAppIDField].(type) {
	case int:
		defaultAppID = strconv.Itoa(id)

	case int64:
		defaultAppID = strconv.FormatInt(id, 10)

	default:
		return "", fmt.Errorf("can't query default appid, unkonw id type, %+v", reflect.TypeOf(results[0][common.BKAppIDField]))
	}

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
		fmt.Sprintf("netdevice2"),
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

		blog.Errorf("DataCollection| get default appid failed: %+v, init database first and it would try again in %+v seconds later",
			err, defaultAppInitWaitDuration)
		time.Sleep(defaultAppInitWaitDuration)
	}
	blog.Info("DataCollection| get default appid id success[%s]", c.defaultAppID)

	// create and add new porters.
	if c.snapCli != nil {
		topic := c.snapMessageTopic(c.defaultAppID)
		analyzer := hostsnap.NewHostSnap(c.ctx, c.redisCli, c.db, c.engine, c.authManager)

		porter := collections.NewSimplePorter(snapPorterName, c.engine, c.hash, analyzer, c.snapCli, topic, c.registry)
		c.porterManager.AddPorter(porter)
		blog.Info("DataCollection| create hostsnap analyzer with target porter[%s] on topic[%s] success", snapPorterName, topic)
	}

	if c.disCli != nil {
		topic := c.discoverMessageTopic(c.defaultAppID)
		analyzer := middleware.NewDiscover(c.ctx, c.redisCli, c.engine, c.authManager)

		porter := collections.NewSimplePorter(middlewarePorterName, c.engine, c.hash, analyzer, c.disCli, topic, c.registry)
		c.porterManager.AddPorter(porter)
		blog.Info("DataCollection| create discover analyzer with target porter[%s] on topic[%s] success", middlewarePorterName, topic)
	}

	if c.netCli != nil {
		topic := c.netcollectMessageTopic(c.defaultAppID)
		analyzer := netcollect.NewNetCollect(c.ctx, c.db, c.authManager)

		porter := collections.NewSimplePorter(netCollectPorterName, c.engine, c.hash, analyzer, c.netCli, topic, c.registry)
		c.porterManager.AddPorter(porter)
		blog.Info("DataCollection| create netcollect analyzer with target porter[%s] on topic[%s] success", netCollectPorterName, topic)
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
