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

package logics

import (
	"configcenter/src/common"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/util"
	"configcenter/src/source_controller/transfer-service/app/options"
)

var (
	bizLgc = &dataWithIDLgc[mapstr.MapStr]{
		idField: common.BKAppIDField,
		table: func(_ string) string {
			return common.BKTableNameBaseApp
		},
		parseData: parseMapStr,
		getID:     getMapStrID,
	}

	setLgc = &dataWithIDLgc[mapstr.MapStr]{
		idField: common.BKSetIDField,
		table: func(_ string) string {
			return common.BKTableNameBaseSet
		},
		parseData: func(data mapstr.MapStr, srcIDConf, destIDConf *options.InnerDataIDConf) (mapstr.MapStr, error) {
			// do not sync set template id
			data[common.BKSetTemplateIDField] = int64(0)

			if srcIDConf == nil || destIDConf == nil || srcIDConf.HostPool == nil || destIDConf.HostPool == nil {
				return parseMapStr(data, srcIDConf, destIDConf)
			}

			// convert src host pool biz id to dest env id
			bizID, err := util.GetInt64ByInterface(data[common.BKAppIDField])
			if err != nil {
				return nil, err
			}
			if bizID == srcIDConf.HostPool.Biz {
				data[common.BKAppIDField] = destIDConf.HostPool.Biz
				data[common.BKParentIDField] = destIDConf.HostPool.Biz
			}
			return parseMapStr(data, srcIDConf, destIDConf)
		},
		getID:         getMapStrID,
		getRelatedIDs: getMapStrRelBizIDInfo,
	}

	moduleLgc = &dataWithIDLgc[mapstr.MapStr]{
		idField: common.BKModuleIDField,
		table: func(_ string) string {
			return common.BKTableNameBaseModule
		},
		parseData: func(data mapstr.MapStr, srcIDConf, destIDConf *options.InnerDataIDConf) (mapstr.MapStr, error) {
			// do not sync set & service template id
			data[common.BKSetTemplateIDField] = int64(0)
			data[common.BKServiceTemplateIDField] = int64(0)

			if srcIDConf == nil || destIDConf == nil || srcIDConf.HostPool == nil || destIDConf.HostPool == nil {
				return parseMapStr(data, srcIDConf, destIDConf)
			}

			// convert src host pool biz & set id to dest env id
			bizID, err := util.GetInt64ByInterface(data[common.BKAppIDField])
			if err != nil {
				return nil, err
			}
			if bizID == srcIDConf.HostPool.Biz {
				data[common.BKAppIDField] = destIDConf.HostPool.Biz
			}

			setID, err := util.GetInt64ByInterface(data[common.BKSetIDField])
			if err != nil {
				return nil, err
			}
			if setID == srcIDConf.HostPool.Set {
				data[common.BKSetIDField] = destIDConf.HostPool.Set
				data[common.BKParentIDField] = destIDConf.HostPool.Set
			}
			return parseMapStr(data, srcIDConf, destIDConf)
		},
		getID:         getMapStrID,
		getRelatedIDs: getMapStrRelBizIDInfo,
	}
)
