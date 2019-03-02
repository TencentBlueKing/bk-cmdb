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

package synchronizer

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
	"configcenter/src/storage/dal/redis"
)

type AuthSynchronizer struct {
	Config *options.Config
	*backbone.Engine
	db  dal.RDB
	ctx context.Context
}

func NewSynchronizer(ctx context.Context, config *options.Config, backbone *backbone.Engine) *AuthSynchronizer {
	return &AuthSynchronizer{ctx: ctx, Config: config, Engine: backbone}
}

func (d *AuthSynchronizer) Run() error {
	blog.Infof("auth synchronizer start...")

	blog.Infof("[datacollect][RUN]connecting to cc redis %+v", d.Config.CCRedis)
	rediscli, err := redis.NewFromConfig(d.Config.CCRedis)
	if nil != err {
		blog.Errorf("[AuthSynchronizer][RUN] connect cc redis failed: %v", err)
		return err
	}
	blog.Infof("[datacollect][RUN]connected to cc redis %+v", d.Config.CCRedis)

	db, err := local.NewMgo(d.Config.MongoDB.BuildURI(), time.Minute)
	if err != nil {
		blog.Errorf("[AuthSynchronizer][RUN] connect mongo failed: %v", err)
		return fmt.Errorf("connect mongo server failed %s", err.Error())
	}
	d.db = db

	var defaultAppID string
	for {
		defaultAppID, err = d.getDefaultAppID(d.ctx)
		if nil == err {
			break
		}
		blog.Errorf("getDefaultAppID faile: %v, please init database first, we will try 10 second later", err)
		time.Sleep(time.Second * 10)
	}

	man := NewManager()

	if d.Config.SnapRedis.Enable != "false" {
		blog.Infof("[datacollect][RUN]connecting to snap-redis %+v", d.Config.SnapRedis.Config)
		snapcli, err := redis.NewFromConfig(d.Config.SnapRedis.Config)
		if nil != err {
			blog.Errorf("[AuthSynchronizer][RUN] connect snap-redis failed: %v", err)
			return err
		}
		blog.Infof("[datacollect][RUN]connected to snap-redis %+v", d.Config.SnapRedis.Config)
		snapChanName := d.getSnapChanName(defaultAppID)
		hostsnapCollector := hostsnap.NewHostSnap(d.ctx, rediscli, db)
		snapPorter := BuildChanPorter("hostsnap", hostsnapCollector, rediscli, snapcli, snapChanName, hostsnap.MockMessage)
		man.AddPorter(snapPorter)
	}

	if d.Config.DiscoverRedis.Enable != "false" {
		blog.Infof("[datacollect][RUN]connecting to discover-redis %+v", d.Config.DiscoverRedis.Config)
		discli, err := redis.NewFromConfig(d.Config.DiscoverRedis.Config)
		if nil != err {
			blog.Errorf("[AuthSynchronizer][RUN] connect discover-redis failed: %v", err)
			return err
		}
		blog.Infof("[datacollect][RUN]connected to discover-redis %+v", d.Config.DiscoverRedis.Config)
		discoverChanName := d.getDiscoverChanName(defaultAppID)
		middlewareCollector := middleware.NewDiscover(d.ctx, rediscli, d.Engine)
		middlewarePorter := BuildChanPorter("middleware", middlewareCollector, rediscli, discli, discoverChanName, middleware.MockMessage)
		man.AddPorter(middlewarePorter)
