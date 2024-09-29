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
	"fmt"

	"configcenter/pkg/synchronize/types"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	commonutil "configcenter/src/common/util"
	"configcenter/src/source_controller/transfer-service/app/options"
	"configcenter/src/source_controller/transfer-service/sync/util"
	"configcenter/src/storage/driver/mongodb"
	"configcenter/src/storage/driver/mongodb/instancemapping"
)

type objInstLogics struct {
	*dataWithIDLogics[mapstr.MapStr]
}

func newObjInstLogics(conf *resLogicsConfig) *objInstLogics {
	return &objInstLogics{
		dataWithIDLogics: newDataWithIDLogics(conf, &dataWithIDLgc[mapstr.MapStr]{
			idField: common.BKInstIDField,
			table: func(subRes string) string {
				// right now only sync object instance from table of the default supplier account
				return common.GetObjectInstTableName(subRes, common.BKDefaultOwnerID)
			},
			parseData: parseMapStr,
			getID:     getMapStrID,
			getRelatedIDs: func(subRes string, data mapstr.MapStr) (map[types.ResType][]int64, error) {
				_, exists := data[common.BKAppIDField]
				if exists {
					return getMapStrRelBizIDInfo(subRes, data)
				}
				return make(map[types.ResType][]int64), nil
			},
		}),
	}
}

// InsertData insert data
func (o *objInstLogics) InsertData(kit *util.Kit, subRes string, data any) error {
	dataArr, ok := data.([]DataWithID[mapstr.MapStr])
	if !ok {
		return fmt.Errorf("data type %T is invalid", data)
	}

	if len(dataArr) == 0 {
		return nil
	}

	insertData := make([]mapstr.MapStr, len(dataArr))
	mappings := make([]metadata.ObjectMapping, len(dataArr))
	for i, info := range dataArr {
		insertData[i] = info.Data

		mappings = append(mappings, metadata.ObjectMapping{
			ID:       info.ID,
			ObjectID: subRes,
			OwnerID:  commonutil.GetStrByInterface(info.Data[common.BkSupplierAccount]),
		})
	}

	table := o.table(subRes)
	err := mongodb.Client().Table(table).Insert(kit.Ctx, insertData)
	if err != nil && !mongodb.Client().IsDuplicatedError(err) {
		blog.Errorf("insert %s data(%+v) failed, err: %v, rid: %s", table, insertData, err, kit.Rid)
		return err
	}

	if err = instancemapping.Create(kit.Ctx, mappings); err != nil {
		blog.Errorf("create object instance mappings(%+v) failed, err: %v, rid: %s", mappings, err, kit.Rid)
		return err
	}
	return nil
}

// DeleteData delete data
func (o *objInstLogics) DeleteData(kit *util.Kit, subRes string, data any) error {
	var ids []int64
	switch val := data.(type) {
	case []int64:
		ids = val
	case []DataWithID[mapstr.MapStr]:
		ids = make([]int64, len(val))
		for i, info := range val {
			ids[i] = info.ID
		}
	default:
		return fmt.Errorf("data type %T is invalid", data)
	}

	if len(ids) == 0 {
		return nil
	}

	if err := instancemapping.Delete(kit.Ctx, ids); err != nil {
		blog.Errorf("delete object instance mapping failed, err: %v, inst ids: %+v, rid: %s", err, ids, kit.Rid)
		return err
	}

	cond := mapstr.MapStr{
		common.BKInstIDField: mapstr.MapStr{common.BKDBIN: ids},
	}

	table := o.table(subRes)
	err := mongodb.Client().Table(table).Delete(kit.Ctx, cond)
	if err != nil {
		blog.Errorf("delete %s data failed, err: %v, cond: %+v, rid: %s", table, err, cond, kit.Rid)
		return err
	}
	return o.dataWithIDLogics.DeleteData(kit, subRes, data)
}

var instAsstLgc = &dataWithIDLgc[metadata.InstAsst]{
	idField: common.BKFieldID,
	table: func(subRes string) string {
		// right now only sync object association from table of the default supplier account
		return common.GetObjectInstAsstTableName(subRes, common.BKDefaultOwnerID)
	},
	parseData: func(data metadata.InstAsst, _, _ *options.InnerDataIDConf) (metadata.InstAsst, error) {
		return data, nil
	},
	getID: func(data metadata.InstAsst, idField string) (int64, error) {
		return data.ID, nil
	},
	getRelatedIDs: func(subRes string, data metadata.InstAsst) (map[types.ResType][]int64, error) {
		idMap := make(map[types.ResType][]int64)
		idMap[getObjResType(data.ObjectID)] = append(idMap[getObjResType(data.ObjectID)], data.InstID)
		idMap[getObjResType(data.AsstObjectID)] = append(idMap[getObjResType(data.AsstObjectID)], data.AsstInstID)
		return idMap, nil
	},
}

var quotedInstLgc = &dataWithIDLgc[mapstr.MapStr]{
	idField: common.BKFieldID,
	table: func(subRes string) string {
		// right now only sync quoted instance from table of the default supplier account
		return common.GetInstTableName(subRes, common.BKDefaultOwnerID)
	},
	parseData: parseMapStr,
	getID:     getMapStrID,
	getRelatedIDs: func(subRes string, data mapstr.MapStr) (map[types.ResType][]int64, error) {
		srcObjID := metadata.GetModelQuoteSrcObjID(subRes)
		instID, err := commonutil.GetInt64ByInterface(data[common.BKInstIDField])
		if err != nil {
			return nil, err
		}
		return map[types.ResType][]int64{getObjResType(srcObjID): {instID}}, nil
	},
}
