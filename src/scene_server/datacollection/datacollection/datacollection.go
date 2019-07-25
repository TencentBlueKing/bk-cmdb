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

package datacollection

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/backbone"
	"configcenter/src/common/blog"
	"configcenter/src/scene_server/datacollection/app/options"
	"configcenter/src/scene_server/datacollection/datacollection/hostsnap"
	"configcenter/src/scene_server/datacollection/datacollection/middleware"
	"configcenter/src/scene_server/datacollection/datacollection/netcollect"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/dal/mongo/local"
	"configcenter/src/storage/dal/mongo/remote"
	"configcenter/src/storage/dal/redis"

	"github.com/prometheus/client_golang/prometheus"
)

type DataCollection struct {
	Config *options.Config
	*backbone.Engine
	db       dal.RDB
	ctx      context.Context
	registry prometheus.Registerer
}

func NewDataCollection(ctx context.Context, config *options.Config, backbone *backbone.Engine, registry prometheus.Registerer) *DataCollection {
	return &DataCollection{ctx: ctx, Config: config, Engine: backbone, registry: registry}
}

func (d *DataCollection) Run() error {
	blog.Infof("data-collection start...")

	blog.Infof("[data-collection][RUN]connecting to cc redis %+v", d.Config.CCRedis)
	redisCli, err := redis.NewFromConfig(d.Config.CCRedis)
	if nil != err {
		blog.Errorf("[data-collection][RUN] connect cc redis failed: %v", err)
		return err
	}
	blog.Infof("[data-collection][RUN]connected to cc redis %+v", d.Config.CCRedis)

	var db dal.RDB
	if d.Config.MongoDB.Enable == "true" {
		db, err = local.NewMgo(d.Config.MongoDB.BuildURI(), time.Minute)
	} else {
		db, err = remote.NewWithDiscover(d.Engine)
	}
	if err != nil {
		blog.Errorf("[data-collection][RUN] connect mongo failed: %v", err)
		return fmt.Errorf("connect mongo server failed %s", err.Error())
	}

	d.db = db

	var defaultAppID string
	for {
		defaultAppID, err = d.getDefaultAppID(d.ctx)
		if nil == err {
			break
		}
		blog.Errorf("getDefaultAppID failed: %v, please init database first, we will try 10 second later", err)
		time.Sleep(time.Second * 10)
	}

	man := NewManager()

	if d.Config.SnapRedis.Enable != "false" {
		blog.Infof("[data-collection][RUN]connecting to snap-redis %+v", d.Config.SnapRedis.Config)
		snapcli, err := redis.NewFromConfig(d.Config.SnapRedis.Config)
		if nil != err {
			blog.Errorf("[data-collection][RUN] connect snap-redis failed: %v", err)
			return err
		}
		blog.Infof("[data-collection][RUN]connected to snap-redis %+v", d.Config.SnapRedis.Config)
		snapChanName := d.getSnapChanName(defaultAppID)
		hostsnapCollector := hostsnap.NewHostSnap(d.ctx, redisCli, db)
		snapPorter := BuildChanPorter("hostsnap", hostsnapCollector, redisCli, snapcli, snapChanName, hostsnap.MockMessage, d.registry, d.Engine)
		man.AddPorter(snapPorter)
	}

	if d.Config.DiscoverRedis.Enable != "false" {
		blog.Infof("[data-collection][RUN]connecting to discover-redis %+v", d.Config.DiscoverRedis.Config)
		disCli, err := redis.NewFromConfig(d.Config.DiscoverRedis.Config)
		if nil != err {
			blog.Errorf("[data-collection][RUN] connect discover-redis failed: %v", err)
			return err
		}
		blog.Infof("[data-collection][RUN]connected to discover-redis %+v", d.Config.DiscoverRedis.Config)
		discoverChanName := d.getDiscoverChanName(defaultAppID)
		middlewareCollector := middleware.NewDiscover(d.ctx, redisCli, d.Engine)
		middlewarePorter := BuildChanPorter("middleware", middlewareCollector, redisCli, disCli, discoverChanName, middleware.MockMessage, d.registry, d.Engine)
		man.AddPorter(middlewarePorter)
	}

	if d.Config.NetCollectRedis.Enable != "false" {
		blog.Infof("[data-collection][RUN]connecting to netcollect-redis %+v", d.Config.NetCollectRedis.Config)
		netCli, err := redis.NewFromConfig(d.Config.NetCollectRedis.Config)
		if nil != err {
			blog.Errorf("[data-collection][RUN] connect netcollect-redis failed: %v", err)
			return err
		}
		blog.Infof("[data-collection][RUN]connected to netcollect-redis %+v", d.Config.NetCollectRedis.Config)
		netDevChanName := d.getNetcollectChanName(defaultAppID)
		netCollector := netcollect.NewNetCollect(d.ctx, db)
		netCollectPorter := BuildChanPorter("netcollect", netCollector, redisCli, netCli, netDevChanName, netcollect.MockMessage, d.registry, d.Engine)
		man.AddPorter(netCollectPorter)
	}

	blog.Infof("data-collection started")
	return nil
}

func (d *DataCollection) getNetcollectChanName(defaultAppID string) []string {
	return []string{"netdevice2"}
}

func (d *DataCollection) getDiscoverChanName(defaultAppID string) []string {
	return []string{"discover" + defaultAppID}
}

func (d *DataCollection) getSnapChanName(defaultAppID string) []string {
	return []string{
		// 瘦身后的通道名
		"snapshot" + defaultAppID,
		// 瘦身前的通道名，为增加向前兼容的而订阅这个老通道
		defaultAppID + "_snapshot",
	}
}

func (d *DataCollection) getDefaultAppID(ctx context.Context) (defaultAppID string, err error) {
	condition := map[string]interface{}{common.BKAppNameField: common.BKAppName}
	results := make([]map[string]interface{}, 0)

	if err = d.db.Table(common.BKTableNameBaseApp).Find(condition).All(ctx, &results); err != nil {
		return "", err
	}

	if len(results) <= 0 {
		return "", fmt.Errorf("default app not found")
	}

	switch id := results[0][common.BKAppIDField].(type) {
	case int:
		defaultAppID = strconv.Itoa(id)
	case int64:
		defaultAppID = strconv.FormatInt(id, 10)
	default:
		return "", fmt.Errorf("default defaultAppID type %v not support", reflect.TypeOf(results[0][common.BKAppIDField]))
	}
	return
}
