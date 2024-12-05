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
	"context"
	"fmt"

	"configcenter/pkg/cache/general"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/source_controller/cacheservice/cache/general/types"
	"configcenter/src/storage/driver/mongodb"
	"configcenter/src/storage/driver/mongodb/instancemapping"
)

func init() {
	addCache(newCacheWithIDAndSubRes[mapstr.MapStr](general.ObjInstKey, common.BKInstIDField,
		[]string{common.BKObjIDField}, getObjInstTable, parseObjInstData))
	addCache(newCacheWithIDAndSubRes[mapstr.MapStr](general.MainlineInstKey, common.BKInstIDField,
		[]string{common.BKObjIDField}, getObjInstTable, parseObjInstData))
}

// getObjInstTable get object instance table by objID and tenant account
// NOTE: obj with "0" tenant can have inst with other suppliers, so we need to search the obj for its actual tenant
func getObjInstTable(ctx context.Context, filter *types.BasicFilter, rid string) (string, error) {
	cond := mapstr.MapStr{
		common.BKObjIDField: filter.SubRes,
	}

	obj := new(metadata.Object)
	err := mongodb.Client().Table(common.BKTableNameObjDes).Find(cond).Fields(common.TenantID).One(ctx, &obj)
	if err != nil {
		blog.Errorf("get object tenant account by cond(%+v) failed, err: %v, rid: %s", cond, err, rid)
		return "", err
	}
	return common.GetInstTableName(filter.SubRes, obj.TenantID), nil
}

func parseObjInstData(data dataWithTable[mapstr.MapStr]) (*basicInfo, error) {
	instID, err := util.GetInt64ByInterface(data.Data[common.BKInstIDField])
	if err != nil {
		return nil, fmt.Errorf("parse id %+v failed, err: %v", data.Data[common.BKInstIDField], err)
	}

	kit := rest.NewKit()
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
		tenant: instObjMappings[0].TenantID,
	}, nil
}
