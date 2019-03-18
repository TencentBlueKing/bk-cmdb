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
	"configcenter/src/auth/authcenter"
	"configcenter/src/common/backbone"
	cc "configcenter/src/common/backbone/configcenter"
	"configcenter/src/common/blog"
	"configcenter/src/scene_server/auth_synchronizer/app/options"
	webservice "configcenter/src/scene_server/auth_synchronizer/pkg/service"
	"configcenter/src/storage/dal/mongo"
	"configcenter/src/storage/dal/redis"
	"encoding/json"
	"sync"
)

// SynchronizerConfig is a container to hold synchronizer config
type SynchronizerConfig struct {
	Core    *backbone.Engine
	Config  *options.Config
	Service *webservice.Service
}

var configLock sync.Mutex

func (h *SynchronizerConfig) onHostConfigUpdate(previous, current cc.ProcessConfig) {
	configLock.Lock()
	defer configLock.Unlock()
	if len(current.ConfigMap) > 0 {
		if h.Config == nil {
			h.Config = new(options.Config)
		}

		out, _ := json.MarshalIndent(current.ConfigMap, "", "  ") //ignore err, cause ConfigMap is map[string]string
		blog.V(3).Infof("config updated: \n%s", out)

		dbprefix := "mongodb"
		mongoConf := mongo.ParseConfigFromKV(dbprefix, current.ConfigMap)
		h.Config.MongoDB = mongoConf

		ccredisPrefix := "redis"
		redisConf := redis.ParseConfigFromKV(ccredisPrefix, current.ConfigMap)
		h.Config.CCRedis = redisConf

		snapPrefix := "snap-redis"
		snapredisConf := redis.ParseConfigFromKV(snapPrefix, current.ConfigMap)
		h.Config.SnapRedis.Config = snapredisConf
		h.Config.SnapRedis.Enable = current.ConfigMap[snapPrefix+".enable"]

		discoverPrefix := "discover-redis"
		discoverRedisConf := redis.ParseConfigFromKV(discoverPrefix, current.ConfigMap)
		h.Config.DiscoverRedis.Config = discoverRedisConf
		h.Config.SnapRedis.Enable = current.ConfigMap[discoverPrefix+".enable"]

		netcollectPrefix := "netcollect-redis"
		netcollectRedisConf := redis.ParseConfigFromKV(netcollectPrefix, current.ConfigMap)
		h.Config.NetcollectRedis.Config = netcollectRedisConf
		h.Config.SnapRedis.Enable = current.ConfigMap[netcollectPrefix+".enable"]

		esbPrefix := "esb"
		h.Config.Esb.Addrs = current.ConfigMap[esbPrefix+".addr"]
		h.Config.Esb.AppCode = current.ConfigMap[esbPrefix+".appCode"]
		h.Config.Esb.AppSecret = current.ConfigMap[esbPrefix+".appSecret"]

		h.Config.Auth, err = authcenter.ParseConfigFromKV("auth", current.ConfigMap)
		if err != nil {
			blog.Warnf("parse authcenter config failed: %v", err)
		}
	}
}
