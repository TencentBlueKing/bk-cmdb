/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 THL A29 Limited,
 * a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 * We undertake not to change the open source license (MIT license) applicable
 * to the current version of the project delivered to anyone in the future.
 */

// Package sync defines cmdb data syncer logics
package sync

import (
	"context"
	"fmt"
	"sort"
	"time"

	"configcenter/pkg/synchronize/types"
	"configcenter/src/apimachinery/cacheservice"
	"configcenter/src/apimachinery/discovery"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/source_controller/transfer-service/app/options"
	"configcenter/src/source_controller/transfer-service/sync/logics"
	"configcenter/src/source_controller/transfer-service/sync/medium"
	"configcenter/src/source_controller/transfer-service/sync/metadata"
	"configcenter/src/source_controller/transfer-service/sync/watch"
	"configcenter/src/storage/dal/mongo/sharding"
	"configcenter/src/storage/driver/mongodb"
	"configcenter/src/storage/stream/task"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/tidwall/gjson"
)

// Syncer is cmdb data syncer
type Syncer struct {
	enableSync   bool
	isMaster     discovery.ServiceManageInterface
	metadata     *metadata.Metadata
	resSyncerMap map[types.ResType]*resSyncer
}

// NewSyncer new cmdb data syncer
func NewSyncer(conf *options.Config, isMaster discovery.ServiceManageInterface, task *task.Task,
	cacheCli cacheservice.CacheServiceClientInterface, reg prometheus.Registerer) (*Syncer, error) {

	if !conf.Sync.EnableSync {
		return &Syncer{enableSync: false}, nil
	}

	// check if id generator is enabled, can only start syncing when id generator is enabled
	configAdminCond := map[string]interface{}{"_id": common.ConfigAdminID}
	configAdminData := make(map[string]string)
	err := mongodb.Shard(sharding.NewShardOpts().WithIgnoreTenant()).Table(common.BKTableNameSystem).
		Find(configAdminCond).Fields(common.ConfigAdminValueField).One(context.Background(), &configAdminData)
	if err != nil {
		blog.Errorf("get config admin data failed, err: %v, cond: %+v", err, configAdminCond)
		return nil, err
	}

	if !gjson.Get(configAdminData[common.ConfigAdminValueField], "id_generator.enabled").Bool() {
		blog.Infof("config admin id generator is not enabled, do not sync cmdb data")
		return &Syncer{enableSync: false}, nil
	}

	meta, err := metadata.NewMetadata(conf.Sync.Role)
	if err != nil {
		blog.Errorf("new metadata failed, err: %v", err)
		return nil, err
	}

	transMedium, err := medium.NewTransferMedium(conf.Sync.TransMediumAddr, reg)
	if err != nil {
		blog.Errorf("new transfer medium failed, err: %v, addr: %+v", err, conf.Sync.TransMediumAddr)
		return nil, err
	}

	idRuleMap, srcInnerIDMap := parseDestExConf(conf)
	resLgcMap := logics.New(&logics.LogicsConfig{
		Metadata:      meta,
		IDRuleMap:     idRuleMap,
		SrcInnerIDMap: srcInnerIDMap,
	})

	syncer := &Syncer{
		enableSync:   true,
		isMaster:     isMaster,
		metadata:     meta,
		resSyncerMap: make(map[types.ResType]*resSyncer),
	}

	for _, resType := range types.ListAllResType() {
		lgc, exists := resLgcMap[resType]
		if !exists {
			return nil, fmt.Errorf("res type %s is invalid", resType)
		}

		syncer.resSyncerMap[resType] = &resSyncer{
			lgc:         lgc,
			name:        conf.Sync.Name,
			transMedium: transMedium,
		}
	}

	err = syncer.run(conf, task, transMedium, cacheCli)
	if err != nil {
		return nil, err
	}

	return syncer, nil
}

func parseDestExConf(conf *options.Config) (map[types.ResType]map[string][]options.IDRuleInfo,
	map[string]*options.InnerDataIDConf) {

	idRuleMap := make(map[types.ResType]map[string][]options.IDRuleInfo)
	innerDataIDMap := make(map[string]*options.InnerDataIDConf)
	if conf.DestExConf == nil {
		return idRuleMap, innerDataIDMap
	}

	// parse id rule config into map[resource]map[src env name]sorted id rules
	for _, ruleInfo := range conf.DestExConf.IDRules {
		for _, rule := range ruleInfo.Rules {
			_, exists := idRuleMap[rule.Resource]
			if !exists {
				idRuleMap[rule.Resource] = make(map[string][]options.IDRuleInfo)
			}
			idRuleMap[rule.Resource][ruleInfo.Name] = append(idRuleMap[rule.Resource][ruleInfo.Name], rule.Rules...)
		}
	}

	for _, resType := range types.ListAllResType() {
		nameRuleMap, exists := idRuleMap[resType]
		if !exists {
			idRuleMap[resType] = make(map[string][]options.IDRuleInfo)
			continue
		}

		for name, infos := range nameRuleMap {
			sort.Slice(infos, func(i, j int) bool {
				return infos[i].StartID < infos[j].StartID
			})
			nameRuleMap[name] = infos
		}
		idRuleMap[resType] = nameRuleMap
	}

	// parse inner data id config into map[src env name]inner data id info
	for i, innerIDInfo := range conf.DestExConf.InnerDataID {
		innerDataIDMap[innerIDInfo.Name] = &conf.DestExConf.InnerDataID[i]
	}

	return idRuleMap, innerDataIDMap
}

func (s *Syncer) run(conf *options.Config, task *task.Task, transMedium medium.ClientI,
	cacheCli cacheservice.CacheServiceClientInterface) error {

	switch conf.Sync.Role {
	case options.SyncRoleSrc:
		go s.loopPushFullSyncData(time.Duration(conf.Sync.SyncIntervalHours) * time.Hour)

		if !conf.Sync.EnableIncrSync {
			return nil
		}

		watcher, err := watch.New(conf.Sync.Name, task, s.isMaster, s.metadata, cacheCli, transMedium)
		if err != nil {
			blog.Errorf("new watcher failed, err: %v", err)
			return err
		}

		if err = watcher.Watch(); err != nil {
			blog.Errorf("watch src event failed, err: %v", err)
			return err
		}
	case options.SyncRoleDest:
		go s.loopPullFullSyncData()

		if !conf.Sync.EnableIncrSync {
			return nil
		}

		for _, resType := range types.ListAllResTypeForIncrSync() {
			go s.loopPullIncrSyncData(resType)
		}
	default:
		return fmt.Errorf("invalid sync role: %s", conf.Sync.Role)
	}
	return nil
}

type resSyncer struct {
	name        string
	transMedium medium.ClientI
	lgc         logics.Logics
	metadata    *metadata.Metadata
}
