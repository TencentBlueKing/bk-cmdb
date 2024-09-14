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
	"configcenter/src/source_controller/transfer-service/app/options"
)

var (
	serviceInstLgc = &dataWithIDLgc[metadata.ServiceInstance]{
		idField: common.BKFieldID,
		table: func(_ string) string {
			return common.BKTableNameServiceInstance
		},
		parseData: func(data metadata.ServiceInstance, _, _ *options.InnerDataIDConf) (metadata.ServiceInstance,
			error) {
			// do not sync service template id
			data.ServiceTemplateID = 0
			return data, nil
		},
		getID: func(data metadata.ServiceInstance, _ string) (int64, error) {
			return data.ID, nil
		},
		getRelatedIDs: func(subRes string, data metadata.ServiceInstance) (map[types.ResType][]int64, error) {
			return map[types.ResType][]int64{types.Biz: {data.BizID}, types.Host: {data.HostID}}, nil
		},
	}

	procLgc = &dataWithIDLgc[mapstr.MapStr]{
		idField: common.BKProcessIDField,
		table: func(_ string) string {
			return common.BKTableNameBaseProcess
		},
		parseData:     parseMapStr,
		getID:         getMapStrID,
		getRelatedIDs: getMapStrRelBizIDInfo,
	}

	procRelLgc = &relationLgc[metadata.ProcessInstanceRelation]{
		idFields: [2]string{common.BKProcessIDField, common.BKServiceInstanceIDField},
		table: func(_ string) string {
			return common.BKTableNameProcessInstanceRelation
		},
		parseData: func(data metadata.ProcessInstanceRelation, _, _ *options.InnerDataIDConf) (
			metadata.ProcessInstanceRelation, error) {
			// do not sync process template id
			data.ProcessTemplateID = 0
			return data, nil
		},
		getIDs: func(data metadata.ProcessInstanceRelation, idFields [2]string) ([2]int64, error) {
			return [2]int64{data.ProcessID, data.ServiceInstanceID}, nil
		},
		getRelatedIDs: func(subRes string, data metadata.ProcessInstanceRelation) (map[types.ResType][]int64, error) {
			return map[types.ResType][]int64{
				types.Process: {data.ProcessID},
				types.Biz:     {data.BizID},
				types.Host:    {data.HostID},
			}, nil
		},
	}
)
