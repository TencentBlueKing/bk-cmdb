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

package cache

import (
	"fmt"

	"configcenter/pkg/cache/general"
	"configcenter/src/common"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/util"
	"configcenter/src/source_controller/cacheservice/cache/general/types"
	"configcenter/src/storage/driver/mongodb/instancemapping"
)

func init() {
	addCache(newCacheWithIDAndSubRes[mapstr.MapStr](general.ObjInstKey, common.BKInstIDField,
		[]string{common.BKObjIDField}, getObjInstTable, parseObjInstData))
	addCache(newCacheWithIDAndSubRes[mapstr.MapStr](general.MainlineInstKey, common.BKInstIDField,
		[]string{common.BKObjIDField}, getObjInstTable, parseObjInstData))
}

// getObjInstTable get object instance table by objID and tenant account
func getObjInstTable(kit *rest.Kit, filter *types.BasicFilter) (string, error) {
	return common.GetInstTableName(filter.SubRes, kit.TenantID), nil
}

func parseObjInstData(data dataWithTenant[mapstr.MapStr]) (*basicInfo, error) {
	instID, err := util.GetInt64ByInterface(data.Data[common.BKInstIDField])
	if err != nil {
		return nil, fmt.Errorf("parse id %+v failed, err: %v", data.Data[common.BKInstIDField], err)
	}

	kit := rest.NewKit().WithTenant(data.TenantID)
	instObjMappings, err := instancemapping.GetInstanceObjectMapping(kit, []int64{instID})
	if err != nil {
		return nil, fmt.Errorf("get object ids from instance ids(%d) failed, err: %v", instID, err)
	}

	if len(instObjMappings) != 1 {
		return nil, fmt.Errorf("inst id %d obj mapping(%+v) is invalid", instID, instObjMappings)
	}

	return &basicInfo{
		id:     instID,
		subRes: []string{instObjMappings[0].ObjectID},
	}, nil
}
