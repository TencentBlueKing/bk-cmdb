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

// Package parser defines the full-text search data parser
package parser

import (
	"configcenter/src/apimachinery/cacheservice"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/sync_server/logics/full-text-search/cache"

	"github.com/olivere/elastic/v7"
)

// Parser defines the es data parser
type Parser interface {
	// GenEsID generate es id
	GenEsID(coll, oid string) string
	// ParseData parse mongo data to es data
	// @param info: one mongo data related info, the first one is the data itself, others are optional extra info
	ParseData(info []mapstr.MapStr, coll string, rid string) (bool, mapstr.MapStr, error)
	// ParseWatchDeleteData parse delete data from mongodb watch
	ParseWatchDeleteData(collOidMap map[string][]string, rid string) ([]string, []elastic.BulkableRequest, bool)
}

// IndexParserMap is the map of es index alias name -> Parser
var IndexParserMap = make(map[string]Parser)

// InitParser initialize parser info
func InitParser(cs *ClientSet) error {
	// init cache data
	if err := cache.InitResourcePoolBiz(); err != nil {
		return err
	}

	// init index to parser map
	IndexParserMap = map[string]Parser{
		metadata.IndexNameBizSet:         newObjInstParser(metadata.IndexNameBizSet, cs),
		metadata.IndexNameBiz:            newObjInstParser(metadata.IndexNameBiz, cs),
		metadata.IndexNameSet:            newBizResParser(metadata.IndexNameSet, cs),
		metadata.IndexNameModule:         newBizResParser(metadata.IndexNameModule, cs),
		metadata.IndexNameHost:           newObjInstParser(metadata.IndexNameHost, cs),
		metadata.IndexNameModel:          newModelParser(metadata.IndexNameModel, cs),
		metadata.IndexNameObjectInstance: newCommonObjInstParser(metadata.IndexNameObjectInstance, cs),
	}

	return nil
}

// ClientSet is the client set of parser
type ClientSet struct {
	EsCli    *elastic.Client
	CacheCli cacheservice.Cache
}
