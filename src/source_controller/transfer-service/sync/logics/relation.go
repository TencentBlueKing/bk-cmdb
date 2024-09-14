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
	"reflect"

	"configcenter/pkg/synchronize/types"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/source_controller/transfer-service/app/options"
	"configcenter/src/source_controller/transfer-service/sync/util"
	"configcenter/src/storage/driver/mongodb"
)

type relationLogics[T any] struct {
	*resLogicsConfig
	*relationLgc[T]
}

type relationLgc[T any] struct {
	table func(subRes string) string
	// idFields are two associated id fields that uniquely define the relation
	idFields [2]string
	// parseData parse data
	parseData func(data T, srcIDConf, destIDConf *options.InnerDataIDConf) (T, error)
	// getIDs get two associated ids matching the idFields
	getIDs func(data T, idFields [2]string) ([2]int64, error)
	// getRelatedIDs get all related ids
	getRelatedIDs func(subRes string, data T) (map[types.ResType][]int64, error)
}

func newRelationLogics[T any](conf *resLogicsConfig, lgc *relationLgc[T]) *relationLogics[T] {
	return &relationLogics[T]{
		resLogicsConfig: conf,
		relationLgc:     lgc,
	}
}

// RelationData is relation data type
type RelationData[T any] struct {
	// IDs are two associated ids that uniquely define the relation data
	IDs  [2]int64
	Data T
}

// ParseDataArr parse data array to actual type
func (l *relationLogics[T]) ParseDataArr(env, subRes string, data any, rid string) (any, error) {
	arr, err := convertDataArr[T](data, rid)
	if err != nil {
		return DataWithID[T]{}, err
	}

	return l.convertToRelationDataArr(true, env, subRes, arr, rid), nil
}

func (l *relationLogics[T]) convertToRelationDataArr(isSrc bool, srcEnv, subRes string, arr []T,
	rid string) []RelationData[T] {

	res := make([]RelationData[T], 0)
	for _, val := range arr {
		// convert src data into dest data
		if isSrc {
			var err error
			val, err = l.parseData(val, l.srcInnerIDMap[srcEnv], l.metadata.InnerIDInfo)
			if err != nil {
				blog.Errorf("parse %s data failed, skip it, err: %v, data: %+v, rid: %s", l.resType, err, val, rid)
				continue
			}
		}

		// check if the data matches the id rule, do not sync if not matches
		ids, err := l.getIDs(val, l.idFields)
		if err != nil {
			blog.Errorf("get %s related ids failed, skip it, err: %v, data: %+v, rid: %s", l.resType, err, val, rid)
			continue
		}

		idMap, err := l.getRelatedIDs(subRes, val)
		if err != nil {
			blog.Errorf("get relation(%+v) related ids failed, skip it, err: %v, rid: %s", val, err, rid)
			continue
		}

		if !util.MatchIDRule(l.idRuleMap, idMap, srcEnv) {
			continue
		}

		res = append(res, RelationData[T]{
			IDs:  ids,
			Data: val,
		})
	}
	return res
}

// ListData list data
func (l *relationLogics[T]) ListData(kit *util.Kit, opt *types.ListDataOpt) (*types.ListDataRes, error) {
	// generate id condition by start and end options
	andConds := make([]mapstr.MapStr, 0)
	if len(opt.Start) > 0 {
		andConds = append(andConds, mapstr.MapStr{
			common.BKDBOR: []mapstr.MapStr{{
				l.idFields[0]: opt.Start[l.idFields[0]],
				l.idFields[1]: mapstr.MapStr{common.BKDBGT: opt.Start[l.idFields[1]]},
			}, {
				l.idFields[0]: mapstr.MapStr{common.BKDBGT: opt.Start[l.idFields[0]]},
			}},
		})
	}

	if len(opt.End) > 0 {
		andConds = append(andConds, mapstr.MapStr{
			common.BKDBOR: []mapstr.MapStr{{
				l.idFields[0]: opt.End[l.idFields[0]],
				l.idFields[1]: mapstr.MapStr{common.BKDBLTE: opt.End[l.idFields[1]]},
			}, {
				l.idFields[0]: mapstr.MapStr{common.BKDBLT: opt.End[l.idFields[0]]},
			}},
		})
	}

	cond := make(mapstr.MapStr)
	if len(andConds) > 0 {
		cond[common.BKDBAND] = andConds
	}
	cond = l.metadata.AddListCond(l.resType, cond)

	// list data from db
	dataArr := make([]T, 0)
	table := l.table(opt.SubRes)
	err := mongodb.Client().Table(table).Find(cond).Sort(fmt.Sprintf("%s:1,%s:1", l.idFields[0], l.idFields[1])).
		Limit(common.BKMaxLimitSize).All(kit.Ctx, &dataArr)
	if err != nil {
		blog.Errorf("list %s relation failed, err: %v, cond: %+v, rid: %s", l.resType, err, cond, kit.Rid)
		return nil, err
	}

	if len(dataArr) == 0 {
		return &types.ListDataRes{
			IsAll:     true,
			Data:      make([]T, 0),
			NextStart: make(map[string]int64),
		}, nil
	}

	// get last data ids as the next start ids, set the associated lastID=startID+1 if all data ids are invalid
	var lastIDs [2]int64
	for i := len(dataArr) - 1; i >= 0; i-- {
		lastIDs, err = l.getIDs(dataArr[i], l.idFields)
		if err != nil {
			blog.Errorf("parse %s relation failed, err: %v, data: %+v, rid: %s", l.resType, err, dataArr[i], kit.Rid)
			continue
		}
		break
	}
	if lastIDs[0] == 0 || lastIDs[1] == 0 {
		lastIDs = [2]int64{opt.Start[l.idFields[0]], opt.Start[l.idFields[1]] + 1}
	}

	return &types.ListDataRes{
		IsAll: len(dataArr) < common.BKMaxLimitSize,
		Data:  dataArr,
		NextStart: map[string]int64{
			l.idFields[0]: lastIDs[0],
			l.idFields[1]: lastIDs[1],
		},
	}, nil
}

// CompareData compare src data with dest data, returns diff data and remaining src data
func (l *relationLogics[T]) CompareData(kit *util.Kit, subRes string, src *types.FullSyncTransData,
	dest *types.ListDataRes) (*types.CompDataRes, error) {

	srcDataArr, ok := src.Data.([]RelationData[T])
	if !ok {
		return nil, fmt.Errorf("src data type %T is invalid", src.Data)
	}
	destDataInfo, ok := dest.Data.([]T)
	if !ok {
		return nil, fmt.Errorf("dest data type %T is invalid", dest.Data)
	}
	// convert data to RelationData type
	destDataArr := l.convertToRelationDataArr(false, src.Name, subRes, destDataInfo, kit.Rid)

	// separate src data into id->data map that are in the interval and remaining data that are not in the interval
	srcDataMap := make(map[[2]int64]RelationData[T])
	remainingSrcData := make([]RelationData[T], 0)
	for _, srcData := range srcDataArr {
		if dest.IsAll || srcData.IDs[0] < dest.NextStart[l.idFields[0]] ||
			(srcData.IDs[0] == dest.NextStart[l.idFields[0]] && srcData.IDs[1] <= dest.NextStart[l.idFields[1]]) {
			srcDataMap[srcData.IDs] = srcData
			continue
		}

		remainingSrcData = append(remainingSrcData, srcData)
	}

	// cross compare src data with dest data
	updateData, insertData := make([]RelationData[T], 0), make([]RelationData[T], 0)
	deleteIDMap := make(map[int64][]int64)
	for _, destData := range destDataArr {
		// delete the dest data if src data not exists or if the src data is not the same with it(delete then insert it)
		srcData, ok := srcDataMap[destData.IDs]
		if !ok || !reflect.DeepEqual(srcData.Data, destData.Data) {
			deleteIDMap[destData.IDs[0]] = append(deleteIDMap[destData.IDs[0]], destData.IDs[1])
			continue
		}

		delete(srcDataMap, destData.IDs)
	}

	for _, data := range srcDataMap {
		insertData = append(insertData, data)
	}

	return &types.CompDataRes{
		Insert:       insertData,
		Update:       updateData,
		Delete:       deleteIDMap,
		RemainingSrc: remainingSrcData,
	}, nil
}

// ClassifyUpsertData classify relation upsert data into insert data and exist data
func (l *relationLogics[T]) ClassifyUpsertData(kit *util.Kit, subRes string, upsertData any) (any, any, error) {
	dataArr, ok := upsertData.([]RelationData[T])
	if !ok {
		return nil, nil, fmt.Errorf("upsert data type %T is invalid", upsertData)
	}

	insertData, updateData := make([]RelationData[T], 0), make([]RelationData[T], 0)
	if len(dataArr) == 0 {
		return insertData, updateData, nil
	}

	// generate relation condition
	idMap := make(map[int64][]int64)
	for _, info := range dataArr {
		idMap[info.IDs[0]] = append(idMap[info.IDs[0]], info.IDs[1])
	}

	orConds := make([]mapstr.MapStr, 0)
	for id, relatedIDs := range idMap {
		orConds = append(orConds, mapstr.MapStr{
			l.idFields[0]: id,
			l.idFields[1]: mapstr.MapStr{common.BKDBIN: relatedIDs},
		})
	}
	cond := mapstr.MapStr{common.BKDBOR: orConds}

	// check if relation exists, only insert the not exist relations
	table := l.table(subRes)
	relations := make([]T, 0)
	err := mongodb.Client().Table(table).Find(cond).Fields(l.idFields[:]...).All(kit.Ctx, &relations)
	if err != nil {
		blog.Errorf("get exist %s ids failed, err: %v, rid: %s", table, err, kit.Rid)
		return nil, nil, err
	}

	if len(relations) == 0 {
		return dataArr, updateData, nil
	}

	existIDMap := make(map[[2]int64]struct{})
	for _, rel := range relations {
		ids, err := l.getIDs(rel, l.idFields)
		if err != nil {
			blog.Errorf("parse exist relation(%+v) ids(%+v) failed, err: %v, rid: %s", rel, l.idFields, err, kit.Rid)
			continue
		}

		existIDMap[ids] = struct{}{}
	}

	for _, data := range dataArr {
		_, exists := existIDMap[data.IDs]
		if !exists {
			insertData = append(insertData, data)
		}
	}
	return insertData, updateData, nil
}

// InsertData insert data
func (l *relationLogics[T]) InsertData(kit *util.Kit, subRes string, data any) error {
	dataArr, ok := data.([]RelationData[T])
	if !ok {
		return fmt.Errorf("data type %T is invalid", data)
	}

	if len(dataArr) == 0 {
		return nil
	}

	insertData := make([]T, len(dataArr))
	for i, info := range dataArr {
		insertData[i] = info.Data
	}

	table := l.table(subRes)
	err := mongodb.Client().Table(table).Insert(kit.Ctx, insertData)
	if err != nil && !mongodb.Client().IsDuplicatedError(err) {
		blog.Errorf("insert %s data(%+v) failed, err: %v, rid: %s", table, insertData, err, kit.Rid)
		return err
	}
	return nil
}

// UpdateData relation can not be updated, skip these data
func (l *relationLogics[T]) UpdateData(kit *util.Kit, subRes string, data any) error {
	return nil
}

// DeleteData delete data
func (l *relationLogics[T]) DeleteData(kit *util.Kit, subRes string, data any) error {
	var idMap map[int64][]int64
	switch val := data.(type) {
	case map[int64][]int64:
		idMap = val
	case []RelationData[T]:
		idMap = make(map[int64][]int64)
		for _, info := range val {
			idMap[info.IDs[0]] = append(idMap[info.IDs[0]], info.IDs[1])
		}
	default:
		return fmt.Errorf("data type %T is invalid", data)
	}

	if len(idMap) == 0 {
		return nil
	}

	orConds := make([]mapstr.MapStr, 0)
	for id, relatedIDs := range idMap {
		orConds = append(orConds, mapstr.MapStr{
			l.idFields[0]: id,
			l.idFields[1]: mapstr.MapStr{common.BKDBIN: relatedIDs},
		})
	}

	cond := mapstr.MapStr{
		common.BKDBOR: orConds,
	}

	table := l.table(subRes)
	err := mongodb.Client().Table(table).Delete(kit.Ctx, cond)
	if err != nil {
		blog.Errorf("delete %s data failed, err: %v, cond: %+v, rid: %s", table, err, cond, kit.Rid)
		return err
	}
	return nil
}
