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

// Package fulltextsearch defines the full-text search synchronization logics
package fulltextsearch

import (
	"context"
	"fmt"
	"strconv"

	types "configcenter/pkg/types/sync/full-text-search"
	"configcenter/src/apimachinery/cacheservice"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/sync_server/logics/full-text-search/parser"
	"configcenter/src/scene_server/sync_server/logics/full-text-search/upgrader"
	"configcenter/src/storage/stream"
	"configcenter/src/thirdparty/elasticsearch"
)

var _ SyncI = new(fullTextSearch)

// SyncI defines the full-text search synchronization interface
type SyncI interface {
	SyncData(ctx context.Context, opt *types.SyncDataOption, rid string) error
	Migrate(ctx context.Context, rid string) (*types.MigrateResult, error)
}

// New full-text search sync interface instance
func New(conf *Config, cacheCli cacheservice.Cache, watcher stream.LoopInterface) (SyncI, error) {
	if !conf.EnableSync {
		return new(fullTextSearch), nil
	}

	if conf.Es.FullTextSearch != "on" {
		return new(fullTextSearch), nil
	}

	f := &fullTextSearch{
		enableSync: conf.EnableSync,
	}

	var err error
	f.esCli, err = elasticsearch.NewEsClient(conf.Es)
	if err != nil {
		blog.Errorf("create es client failed, err: %v, conf: %+v", err, conf)
		return nil, err
	}

	f.cacheCli = cacheCli

	if conf.IndexShardNum <= 0 || conf.IndexReplicaNum <= 0 {
		return nil, fmt.Errorf("index shard num %d or replica num %d is invalid", conf.IndexShardNum,
			conf.IndexReplicaNum)
	}

	indexSetting := metadata.ESIndexMetaSettings{
		Shards:   strconv.Itoa(conf.IndexShardNum),
		Replicas: strconv.Itoa(conf.IndexReplicaNum),
	}

	upgrader.InitUpgrader(f.esCli.Client, indexSetting)
	if _, err = f.Migrate(context.Background(), util.GenerateRID()); err != nil {
		blog.Errorf("migrate failed, err: %v, conf: %+v", err, conf)
		return nil, err
	}

	parserClientSet := &parser.ClientSet{
		EsCli:    f.esCli.Client,
		CacheCli: cacheCli,
	}
	if err = parser.InitParser(parserClientSet); err != nil {
		blog.Errorf("init parser failed, err: %v", err)
		return nil, err
	}

	if err = f.incrementalSync(watcher); err != nil {
		blog.Errorf("start full-text search incremental sync failed, err: %v, conf: %+v", err, conf)
		return nil, err
	}

	return f, nil
}

// Config defines full-text search sync configuration
type Config struct {
	// EnableSync defines if full-text search sync is enabled
	EnableSync bool `mapstructure:"enableSync"`
	// IndexShardNum defines the number of es index shards
	IndexShardNum int `mapstructure:"indexShardNum"`
	// IndexReplicaNum defines the number of es index replicas
	IndexReplicaNum int `mapstructure:"indexReplicaNum"`

	//  Es elasticsearch configuration
	Es *elasticsearch.EsConfig
}

// fullTextSearch implements the full-text search synchronization interface
type fullTextSearch struct {
	enableSync bool
	esCli      *elasticsearch.EsSrv
	cacheCli   cacheservice.Cache
}
