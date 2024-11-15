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
	"strconv"
	"strings"

	"configcenter/pkg/cache/general"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/source_controller/cacheservice/cache/general/types"
	"configcenter/src/storage/driver/mongodb"

	"github.com/tidwall/gjson"
)

func init() {
	cache := newCacheWithID[metadata.HostMapStr](general.HostKey, true, common.BKTableNameBaseHost,
		common.BKHostIDField, func(data metadata.HostMapStr, idField string) (*basicInfo, error) {
			id, err := util.GetInt64ByInterface(data[idField])
			if err != nil {
				return nil, fmt.Errorf("parse host id %+v failed, err: %v", data[idField], err)
			}
			return &basicInfo{
				id:     id,
				tenant: util.GetStrByInterface(data[common.TenantID]),
			}, nil
		})

	cache.uniqueKeyLogics = map[general.UniqueKeyType]uniqueKeyLogics{
		general.IPCloudIDType: {
			genKey:  genHostIPCloudIDKey,
			getData: genHostByIPCloudIDKey,
		},
		general.AgentIDType: {
			genKey:  genHostAgentIDKey,
			getData: genHostByAgentIDKey,
		},
	}

	addCache(cache)
}

func genHostIPCloudIDKey(data any, info *basicInfo) ([]string, error) {
	var cloudID int64
	ips := make([]string, 0)
	var addressType string

	switch val := data.(type) {
	case metadata.HostMapStr:
		var err error
		cloudID, err = util.GetInt64ByInterface(val[common.BKCloudIDField])
		if err != nil {
			return nil, fmt.Errorf("parse cloud id %+v failed, err: %v", val[common.BKCloudIDField], err)
		}

		ip := util.GetStrByInterface(val[common.BKHostInnerIPField])
		ips = strings.Split(ip, ",")
		addressType = util.GetStrByInterface(val[common.BKAddressingField])
	case types.WatchEventData:
		elements := gjson.GetMany(string(val.Data), common.BKCloudIDField, common.BKHostInnerIPField,
			common.BKAddressingField)

		cloudID = elements[0].Int()
		ipVal := elements[1].Array()
		for _, ip := range ipVal {
			ips = append(ips, ip.String())
		}
		addressType = elements[2].String()
	default:
		return nil, fmt.Errorf("data type %T is invalid", data)
	}

	keys := make([]string, 0)
	if addressType != "" && addressType != common.BKAddressingStatic {
		return keys, nil
	}

	for _, ip := range ips {
		keys = append(keys, general.IPCloudIDKey(ip, cloudID))
	}
	return keys, nil
}

func genHostByIPCloudIDKey(ctx context.Context, opt *getDataByKeysOpt, rid string) ([]any, error) {
	if len(opt.Keys) == 0 {
		return make([]any, 0), nil
	}

	ctx = util.SetDBReadPreference(ctx, common.SecondaryPreferredMode)

	cloudIDIPMap := make(map[int64][]string)
	for i, key := range opt.Keys {
		pair := strings.Split(key, ":")
		if len(pair) != 2 {
			blog.Errorf("host ip cloud id key %s is invalid, index: %d, rid: %s", key, i, rid)
			return nil, fmt.Errorf("ip cloud id key %s is invalid", key)
		}

		cloudID, err := strconv.ParseInt(pair[1], 10, 64)
		if err != nil {
			blog.Errorf("parse cloud id (index: %d, key: %s) failed, err: %v, rid: %s", i, key, err, rid)
			return nil, err
		}

		cloudIDIPMap[cloudID] = append(cloudIDIPMap[cloudID], pair[0])
	}

	condArr := make([]mapstr.MapStr, 0)
	for cloudID, ips := range cloudIDIPMap {
		condArr = append(condArr, mapstr.MapStr{
			common.BKCloudIDField: cloudID,
			common.BKHostInnerIPField: mapstr.MapStr{
				common.BKDBIN: ips,
			},
		})
	}

	cond := mapstr.MapStr{
		common.BKDBOR:            condArr,
		common.BKAddressingField: common.BKAddressingStatic,
	}

	hosts := make([]metadata.HostMapStr, 0)
	if err := mongodb.Client().Table(common.BKTableNameBaseHost).Find(cond).All(ctx, &hosts); err != nil {
		blog.Errorf("get host data by cond(%+v) failed, err: %v, rid: %s", cond, err, rid)
		return nil, err
	}

	return convertToAnyArr(hosts), nil
}

func genHostAgentIDKey(data any, info *basicInfo) ([]string, error) {
	var agentID string

	switch val := data.(type) {
	case metadata.HostMapStr:
		agentID = util.GetStrByInterface(val[common.BKAgentIDField])
	case types.WatchEventData:
		agentID = gjson.Get(string(val.Data), common.BKAgentIDField).String()
	default:
		return nil, fmt.Errorf("data type %T is invalid", data)
	}

	if agentID == "" {
		return make([]string, 0), nil
	}

	return []string{general.AgentIDKey(agentID)}, nil
}

func genHostByAgentIDKey(ctx context.Context, opt *getDataByKeysOpt, rid string) ([]any, error) {
	if len(opt.Keys) == 0 {
		return make([]any, 0), nil
	}

	ctx = util.SetDBReadPreference(ctx, common.SecondaryPreferredMode)

	cond := mapstr.MapStr{
		common.BKAgentIDField: mapstr.MapStr{common.BKDBType: "string", common.BKDBIN: opt.Keys},
	}

	hosts := make([]metadata.HostMapStr, 0)
	if err := mongodb.Client().Table(common.BKTableNameBaseHost).Find(cond).All(ctx, &hosts); err != nil {
		blog.Errorf("get host data by cond(%+v) failed, err: %v, rid: %s", cond, err, rid)
		return nil, err
	}

	return convertToAnyArr(hosts), nil
}
