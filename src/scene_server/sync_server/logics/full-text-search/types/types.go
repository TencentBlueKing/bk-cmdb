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

// Package types defines the full-text search synchronization common types
package types

import (
	"configcenter/src/common"
	"configcenter/src/common/metadata"
)

var (
	// AllIndexNames is all elastic index names
	AllIndexNames = []string{metadata.IndexNameBizSet, metadata.IndexNameBiz, metadata.IndexNameSet,
		metadata.IndexNameModule, metadata.IndexNameHost, metadata.IndexNameModel, metadata.IndexNameObjectInstance}

	// IndexVersionMap is elastic alias index name to version map
	// NOTE: CHANGE the version name if you have modified the indexes metadata struct.
	IndexVersionMap = map[string]string{
		metadata.IndexNameBizSet:         "20210710",
		metadata.IndexNameBiz:            "20210710",
		metadata.IndexNameSet:            "20210710",
		metadata.IndexNameModule:         "20210710",
		metadata.IndexNameHost:           "20210710",
		metadata.IndexNameModel:          "20210710",
		metadata.IndexNameObjectInstance: "20210710",
	}

	// IndexMap is es index alias name to related index data map
	IndexMap = make(map[string][]*metadata.ESIndex)

	// IndexExtraFieldsMap is the map of es index alias name -> extra fields in the index besides common fields
	IndexExtraFieldsMap = map[string][]string{
		metadata.IndexNameBizSet: {metadata.IndexPropertyBKBizSetID},
		metadata.IndexNameSet:    {metadata.IndexPropertyBKParentID},
		metadata.IndexNameHost:   {metadata.IndexPropertyBKCloudID},
	}

	// IndexExcludeFieldsMap is the map of es index alias name -> the excluded fields of common fields
	IndexExcludeFieldsMap = map[string][]string{
		metadata.IndexNameBizSet: {metadata.IndexPropertyBKBizID},
		metadata.IndexNameHost:   {metadata.IndexPropertyBKBizID},
	}

	// IndexCollMap is the map of es index alias name -> cmdb collection
	IndexCollMap = map[string]string{
		metadata.IndexNameBizSet:         common.BKTableNameBaseBizSet,
		metadata.IndexNameBiz:            common.BKTableNameBaseApp,
		metadata.IndexNameSet:            common.BKTableNameBaseSet,
		metadata.IndexNameModule:         common.BKTableNameBaseModule,
		metadata.IndexNameHost:           common.BKTableNameBaseHost,
		metadata.IndexNameModel:          common.BKTableNameObjDes,
		metadata.IndexNameObjectInstance: common.BKTableNameBaseInst,
	}
)

// GetIndexName get actual index name by alias name
// right now one alias name is related to only one index
func GetIndexName(alias string) string {
	for _, index := range IndexMap[alias] {
		return index.Name()
	}
	return alias
}
