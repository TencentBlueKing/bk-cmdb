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

	"configcenter/src/auth/extensions"
	"configcenter/src/common"
	"configcenter/src/common/backbone"
	"configcenter/src/common/blog"
	"configcenter/src/scene_server/datacollection/datacollection/hostsnap"
	"configcenter/src/scene_server/datacollection/datacollection/middleware"
	"configcenter/src/scene_server/datacollection/datacollection/netcollect"
	"configcenter/src/storage/dal"

	"github.com/prometheus/client_golang/prometheus"
	"gopkg.in/redis.v5"
)

type DataCollection struct {
	*backbone.Engine
	db          dal.RDB
	ctx         context.Context
	registry    prometheus.Registerer
	AuthManager extensions.AuthManager
}

func NewDataCollection(ctx context.Context, backbone *backbone.Engine, db dal.RDB, registry prometheus.Registerer) *DataCollection {
	return &DataCollection{ctx: ctx, Engine: backbone, db: db, registry: registry}
}

func (d *DataCollection) Run(redisCli, snapCli, disCli, netCli *redis.Client) error {
	blog.Infof("data-collection start...")

	var err error
	var defaultAppID string
	for {
		defaultAppID, err = d.getDefaultAppID(d.ctx)
		if nil == err {
			break
		}
		blog.Errorf("getDefaultAppID failed: %v, please init database first, we will try 10 second later", err)
		time.Sleep(time.Second * 10)
	}

	manager := NewManager()

	if snapCli != nil {
		snapChanName := d.getSnapChanName(defaultAppID)
		hostsnapCollector := hostsnap.NewHostSnap(d.ctx, redisCli, d.db, d.AuthManager)
		snapPorter := BuildChanPorter("hostsnap", hostsnapCollector, redisCli, snapCli, snapChanName, hostsnap.MockMessage, d.registry, d.Engine)
		manager.AddPorter(snapPorter)
	}
	if disCli != nil {
		discoverChanName := d.getDiscoverChanName(defaultAppID)
		middlewareCollector := middleware.NewDiscover(d.ctx, redisCli, d.Engine, d.AuthManager)
		middlewarePorter := BuildChanPorter("middleware", middlewareCollector, redisCli, disCli, discoverChanName, middleware.MockMessage, d.registry, d.Engine)
		manager.AddPorter(middlewarePorter)
	}
	if netCli != nil {
		netDevChanName := d.getNetcollectChanName(defaultAppID)
		netCollector := netcollect.NewNetCollect(d.ctx, d.db, d.AuthManager)
		netCollectPorter := BuildChanPorter("netcollect", netCollector, redisCli, netCli, netDevChanName, netcollect.MockMessage, d.registry, d.Engine)
		manager.AddPorter(netCollectPorter)
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
