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
	"configcenter/pkg/synchronize/types"
	"configcenter/src/common"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	commonutil "configcenter/src/common/util"
	"configcenter/src/source_controller/transfer-service/app/options"
)

var (
	hostLgc = &dataWithIDLgc[metadata.HostMapStr]{
		idField: common.BKHostIDField,
		table: func(subRes string) string {
			return common.BKTableNameBaseHost
		},
		parseData: func(data metadata.HostMapStr, _, _ *options.InnerDataIDConf) (metadata.HostMapStr, error) {
			// convert host special string field(ip and operator) to array
			data, err := metadata.ConvertHostSpecialStringToArray(data)
			if err != nil {
				return nil, err
			}

			mapstrData, err := parseMapStr(mapstr.MapStr(data), nil, nil)
			if err != nil {
				return nil, err
			}
			return metadata.HostMapStr(mapstrData), nil
		},
		getID: func(data metadata.HostMapStr, idField string) (int64, error) {
			return commonutil.GetInt64ByInterface(data[idField])
		},
	}

	hostRelLgc = &relationLgc[metadata.ModuleHost]{
		idFields: [2]string{common.BKHostIDField, common.BKModuleIDField},
		table: func(_ string) string {
			return common.BKTableNameModuleHostConfig
		},
		parseData: func(data metadata.ModuleHost, srcIDConf, destIDConf *options.InnerDataIDConf) (metadata.ModuleHost,
			error) {

			if srcIDConf == nil || destIDConf == nil || srcIDConf.HostPool == nil || destIDConf.HostPool == nil {
				return data, nil
			}

			// convert src host pool biz & set & module id to dest env id
			if data.AppID == srcIDConf.HostPool.Biz {
				data.AppID = destIDConf.HostPool.Biz
			}
			if data.SetID == srcIDConf.HostPool.Set {
				data.SetID = destIDConf.HostPool.Set
			}
			if data.ModuleID == srcIDConf.HostPool.Module {
				data.ModuleID = destIDConf.HostPool.Module
			}

			return data, nil
		},
		getIDs: func(data metadata.ModuleHost, idFields [2]string) ([2]int64, error) {
			return [2]int64{data.HostID, data.ModuleID}, nil
		},
		getRelatedIDs: func(subRes string, data metadata.ModuleHost) (map[types.ResType][]int64, error) {
			return map[types.ResType][]int64{types.Biz: {data.AppID}, types.Host: {data.HostID}}, nil
		},
	}
)
