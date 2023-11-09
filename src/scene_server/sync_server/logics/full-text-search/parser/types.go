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

package parser

import (
	"configcenter/src/common"
	"configcenter/src/common/metadata"
)

const (
	// nullMetaID default metaID
	nullMetaID = "0"
	// commonObject common object instance identifier
	commonObject = "common"
)

const (
	// deleteTableQueryScript 表格实例删除脚本的条件
	// 例： 删除disk表格中实例_id为1的行 {"field": "tables.disk.1"}
	deleteTableQueryScript = "tables.%s.%s"
	// deleteTableScript 表格实例删除脚本，
	// 例：删除disk表格中实例_id为1的行 ctx._source.tables.disk.remove('1')，如果删除后表格为空则删除表格字段
	deleteTableScript = `ctx._source.tables.%s.remove('%s');
                         if (ctx._source.tables.%s.size()==0) {ctx._source.tables.remove('%s')}`
	// updateTableScript 表格实例更新脚本（如果tables字段和表格字段不存在则先创建再更新）
	// 例：更新disk表格中实例_id为1的行的keyword为xxx ctx._source.tables.disk['1'] = ["xxx"]
	updateTableScript = `if(!ctx._source.containsKey('tables')){ctx._source['tables']=[:];}
                         if(!ctx._source.tables.containsKey('%s')){ctx._source.tables['%s']=[:];}
                         ctx._source.tables.%s['%s']=%s`
)

var (
	// extraEsFieldMap is the extra es field to cc field map
	extraEsFieldMap = map[string]string{
		metadata.IndexPropertyBKBizSetID: common.BKBizSetIDField,
		metadata.IndexPropertyBKParentID: common.BKParentIDField,
		metadata.IndexPropertyBKCloudID:  common.BKCloudIDField,
	}

	// baseCleanFields is the basic fields that should be cleaned from the keyword
	baseCleanFields = []string{common.MongoMetaID, common.CreateTimeField, common.LastTimeField, common.BKOwnerIDField}

	// indexKeywordCleanFieldsMap is the map of es index name -> the fields that should be cleaned from the keyword
	indexKeywordCleanFieldsMap = map[string][]string{
		metadata.IndexNameBizSet: {common.BKDefaultField, common.BKBizSetScopeField},
		metadata.IndexNameBiz:    {common.BKDefaultField, common.BKParentIDField},
		metadata.IndexNameSet: {common.BKAppIDField, common.BKParentIDField, common.BKSetTemplateIDField,
			common.BKDefaultField},
		metadata.IndexNameModule: {common.BKDefaultField, common.BKSetTemplateIDField, common.BKAppIDField,
			common.BKParentIDField, common.BKSetIDField, common.BKServiceCategoryIDField},
		metadata.IndexNameHost:           {common.BKOperationTimeField, common.BKParentIDField},
		metadata.IndexNameObjectInstance: {common.BKObjIDField, common.BKParentIDField},
	}

	// indexIdentifierMap is the map of es index name -> the identifier of the index, used as suffix of es id
	indexIdentifierMap = map[string]string{
		metadata.IndexNameBizSet:         common.BKInnerObjIDBizSet,
		metadata.IndexNameBiz:            common.BKInnerObjIDApp,
		metadata.IndexNameSet:            common.BKInnerObjIDSet,
		metadata.IndexNameModule:         common.BKInnerObjIDModule,
		metadata.IndexNameHost:           common.BKInnerObjIDHost,
		metadata.IndexNameModel:          common.BKInnerObjIDObject,
		metadata.IndexNameObjectInstance: commonObject,
	}
)
