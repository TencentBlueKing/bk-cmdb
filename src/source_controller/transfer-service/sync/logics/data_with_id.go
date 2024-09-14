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
	commonutil "configcenter/src/common/util"
	"configcenter/src/source_controller/transfer-service/app/options"
	"configcenter/src/source_controller/transfer-service/sync/util"
	"configcenter/src/storage/driver/mongodb"

	"github.com/google/go-cmp/cmp"
)

type dataWithIDLogics[T any] struct {
	*resLogicsConfig
	*dataWithIDLgc[T]
}

type dataWithIDLgc[T any] struct {
	idField       string
	table         func(subRes string) string
	parseData     func(data T, srcIDConf, destIDConf *options.InnerDataIDConf) (T, error)
	getID         func(data T, idField string) (int64, error)
	getRelatedIDs func(subRes string, data T) (map[types.ResType][]int64, error)
}

func newDataWithIDLogics[T any](conf *resLogicsConfig, lgc *dataWithIDLgc[T]) *dataWithIDLogics[T] {
	return &dataWithIDLogics[T]{
		resLogicsConfig: conf,
		dataWithIDLgc:   lgc,
	}
}

// DataWithID is data with id type
type DataWithID[T any] struct {
	ID   int64
	Data T
}

// ParseDataArr parse data array to actual type
func (l *dataWithIDLogics[T]) ParseDataArr(env, subRes string, data any, rid string) (any, error) {
	arr, err := convertDataArr[T](data, rid)
	if err != nil {
		return DataWithID[T]{}, err
	}

	return l.convertToDataWithIDArr(true, env, subRes, arr, rid), nil
}

func (l *dataWithIDLogics[T]) convertToDataWithIDArr(isSrc bool, srcEnv, subRes string, arr []T,
	rid string) []DataWithID[T] {

	res := make([]DataWithID[T], 0)
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
		id, err := l.getID(val, l.idField)
		if err != nil {
			blog.Errorf("get %s data id failed, skip it, err: %v, data: %+v, rid: %s", l.resType, err, val, rid)
			continue
		}

		idMap := make(map[types.ResType][]int64)
		if l.getRelatedIDs != nil {
			relIDMap, err := l.getRelatedIDs(subRes, val)
			if err != nil {
				blog.Errorf("get data(%+v) related ids failed, skip it, err: %v, rid: %s", val, err, rid)
				continue
			}
			idMap = relIDMap
		}
		idMap[l.resType] = []int64{id}

		if !util.MatchIDRule(l.idRuleMap, idMap, srcEnv) {
			continue
		}

		res = append(res, DataWithID[T]{
			ID:   id,
			Data: val,
		})
	}
	return res
}

// ListData list data
func (l *dataWithIDLogics[T]) ListData(kit *util.Kit, opt *types.ListDataOpt) (*types.ListDataRes, error) {
	// generate id condition by start and end options
	idCond := mapstr.MapStr{common.BKDBGT: 0}
	if len(opt.Start) > 0 {
		idCond[common.BKDBGT] = opt.Start[l.idField]
	}

	if len(opt.End) > 0 && opt.End[l.idField] != types.InfiniteEndID {
		idCond[common.BKDBLTE] = opt.End[l.idField]
	}

	cond := l.metadata.AddListCond(l.resType, mapstr.MapStr{
		l.idField: idCond,
	})

	// list data from db
	dataArr := make([]T, 0)
	table := l.table(opt.SubRes)
	err := mongodb.Client().Table(table).Find(cond).Sort(l.idField).Limit(common.BKMaxLimitSize).All(kit.Ctx, &dataArr)
	if err != nil {
		blog.Errorf("list %s data failed, err: %v, cond: %+v, rid: %s", l.resType, err, cond, kit.Rid)
		return nil, err
	}

	if len(dataArr) == 0 {
		return &types.ListDataRes{
			IsAll:     true,
			Data:      make([]T, 0),
			NextStart: make(map[string]int64),
		}, nil
	}

	// get last data id as the next start id, set lastID=startID+1 if all data ids are invalid
	var lastID int64
	for i := len(dataArr) - 1; i >= 0; i-- {
		lastID, err = l.getID(dataArr[i], l.idField)
		if err != nil {
			blog.Errorf("parse %s data failed, err: %v, data: %+v, rid: %s", l.resType, err, dataArr[i], kit.Rid)
			continue
		}
		break
	}
	if lastID == 0 {
		lastID = opt.Start[l.idField] + 1
	}

	return &types.ListDataRes{
		IsAll:     len(dataArr) < common.BKMaxLimitSize,
		Data:      dataArr,
		NextStart: map[string]int64{l.idField: lastID},
	}, nil
}

// CompareData compare src data with dest data, returns diff data and remaining src data
func (l *dataWithIDLogics[T]) CompareData(kit *util.Kit, subRes string, srcInfo *types.FullSyncTransData,
	destInfo *types.ListDataRes) (*types.CompDataRes, error) {

	srcDataArr, ok := srcInfo.Data.([]DataWithID[T])
	if !ok {
		return nil, fmt.Errorf("src data type %T is invalid", srcInfo.Data)
	}
	destDataInfo, ok := destInfo.Data.([]T)
	if !ok {
		return nil, fmt.Errorf("dest data type %T is invalid", destInfo.Data)
	}
	destDataArr := l.convertToDataWithIDArr(false, srcInfo.Name, subRes, destDataInfo, kit.Rid)

	// separate src data into id->data map that are in the interval and remaining data that are not in the interval
	srcDataMap := make(map[int64]DataWithID[T])
	remainingSrcData := make([]DataWithID[T], 0)
	for _, srcData := range srcDataArr {
		if destInfo.IsAll || srcData.ID <= destInfo.NextStart[l.idField] {
			srcDataMap[srcData.ID] = srcData
			continue
		}

		remainingSrcData = append(remainingSrcData, srcData)
	}

	// cross compare src data with dest data
	updateData, insertData := make([]DataWithID[T], 0), make([]DataWithID[T], 0)
	deleteIDs := make([]int64, 0)
	for _, destData := range destDataArr {
		srcData, ok := srcDataMap[destData.ID]
		if !ok {
			deleteIDs = append(deleteIDs, destData.ID)
			continue
		}

		if !cmp.Equal(srcData.Data, destData.Data) {
			updateData = append(updateData, srcData)
		}

		delete(srcDataMap, destData.ID)
	}

	for _, data := range srcDataMap {
		insertData = append(insertData, data)
	}

	return &types.CompDataRes{
		Insert:       insertData,
		Update:       updateData,
		Delete:       deleteIDs,
		RemainingSrc: remainingSrcData,
	}, nil
}

// ClassifyUpsertData classify upsert data into insert and update data
func (l *dataWithIDLogics[T]) ClassifyUpsertData(kit *util.Kit, subRes string, upsertData any) (any, any, error) {
	dataArr, ok := upsertData.([]DataWithID[T])
	if !ok {
		return nil, nil, fmt.Errorf("upsert data type %T is invalid", upsertData)
	}

	insertData, updateData := make([]DataWithID[T], 0), make([]DataWithID[T], 0)
	if len(dataArr) == 0 {
		return insertData, updateData, nil
	}

	// get exist ids to judge if data exists
	ids := make([]int64, len(dataArr))
	for i, data := range dataArr {
		ids[i] = data.ID
	}

	cond := mapstr.MapStr{l.idField: mapstr.MapStr{common.BKDBIN: ids}}
	table := l.table(subRes)
	rawIDs, err := mongodb.Client().Table(table).Distinct(kit.Ctx, l.idField, cond)
	if err != nil {
		blog.Errorf("get exist %s ids failed, err: %v, rid: %s", table, err, kit.Rid)
		return nil, nil, err
	}

	if len(rawIDs) == 0 {
		return dataArr, updateData, nil
	}

	ids, err = commonutil.SliceInterfaceToInt64(rawIDs)
	if err != nil {
		blog.Errorf("parse raw ids(%+v) failed, err: %v, rid: %s", rawIDs, err, kit.Rid)
		return nil, nil, err
	}

	// classify upsert data by id
	existIDMap := make(map[int64]struct{})
	for _, id := range ids {
		existIDMap[id] = struct{}{}
	}

	for _, data := range dataArr {
		_, exists := existIDMap[data.ID]
		if exists {
			updateData = append(updateData, data)
			continue
		}
		insertData = append(insertData, data)
	}
	return insertData, updateData, nil
}

// InsertData insert data
func (l *dataWithIDLogics[T]) InsertData(kit *util.Kit, subRes string, data any) error {
	dataArr, ok := data.([]DataWithID[T])
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

// UpdateData update data
func (l *dataWithIDLogics[T]) UpdateData(kit *util.Kit, subRes string, data any) error {
	dataArr, ok := data.([]DataWithID[T])
	if !ok {
		return fmt.Errorf("data type %T is invalid", data)
	}

	if len(dataArr) == 0 {
		return nil
	}

	table := l.table(subRes)
	for _, info := range dataArr {
		cond := mapstr.MapStr{l.idField: info.ID}
		err := mongodb.Client().Table(table).Update(kit.Ctx, cond, info.Data)
		if err != nil && !mongodb.Client().IsDuplicatedError(err) {
			blog.Errorf("update %s data(%+v) failed, err: %v, rid: %s", table, info, err, kit.Rid)
			return err
		}
	}
	return nil
}

// DeleteData delete data
func (l *dataWithIDLogics[T]) DeleteData(kit *util.Kit, subRes string, data any) error {
	var ids []int64
	switch val := data.(type) {
	case []int64:
		ids = val
	case []DataWithID[T]:
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

	cond := mapstr.MapStr{
		l.idField: mapstr.MapStr{common.BKDBIN: ids},
	}

	table := l.table(subRes)
	err := mongodb.Client().Table(table).Delete(kit.Ctx, cond)
	if err != nil {
		blog.Errorf("delete %s data failed, err: %v, cond: %+v, rid: %s", table, err, cond, kit.Rid)
		return err
	}
	return nil
}
