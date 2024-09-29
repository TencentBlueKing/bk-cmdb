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
	"encoding/json"
	"fmt"

	"configcenter/pkg/synchronize/types"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	ccjson "configcenter/src/common/json"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/util"
	"configcenter/src/source_controller/transfer-service/app/options"
)

// convertDataArr convert data array to actual type
func convertDataArr[T any](data any, rid string) ([]T, error) {
	switch val := data.(type) {
	case []T:
		return val, nil
	case []json.RawMessage:
		arr := make([]T, len(val))
		for i, raw := range val {
			err := ccjson.Unmarshal(raw, &arr[i])
			if err != nil {
				blog.Errorf("decode json data(%s) failed, err: %v, rid: %s", raw, err, rid)
				continue
			}
		}
		return arr, nil
	default:
		blog.Errorf("data(%+v) type(%T) is invalid, rid: %s", data, data, rid)
		return nil, fmt.Errorf("data type(%T) is invalid", data)
	}
}

func parseMapStr(data mapstr.MapStr, _, _ *options.InnerDataIDConf) (mapstr.MapStr, error) {
	for k, v := range data {
		switch val := v.(type) {
		case json.Number:
			if intVal, err := val.Int64(); err == nil {
				data[k] = intVal
				continue
			}
			if floatVal, err := val.Float64(); err == nil {
				data[k] = floatVal
				continue
			}
		}
	}
	delete(data, "_id")
	return data, nil
}

func getMapStrID(data mapstr.MapStr, idField string) (int64, error) {
	return util.GetInt64ByInterface(data[idField])
}

func getMapStrRelBizIDInfo(_ string, data mapstr.MapStr) (map[types.ResType][]int64, error) {
	bizID, err := util.GetInt64ByInterface(data[common.BKAppIDField])
	if err != nil {
		return nil, err
	}
	return map[types.ResType][]int64{types.Biz: {bizID}}, nil
}

func getObjResType(objID string) types.ResType {
	switch objID {
	case common.BKInnerObjIDApp:
		return types.Biz
	case common.BKInnerObjIDSet:
		return types.Set
	case common.BKInnerObjIDModule:
		return types.Module
	case common.BKInnerObjIDHost:
		return types.Host
	case common.BKInnerObjIDProc:
		return types.Process
	default:
		return types.ObjectInstance
	}
}
