/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.,
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the ",License",); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an ",AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package operation

import (
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/source_controller/coreservice/core"
)

var (
	FreshDataInterval int64 = 12
)

func (m *operationManager) TimerFreshData(params core.ContextParams, data interface{}) {
	timer := time.NewTicker(time.Duration(FreshDataInterval) * time.Hour)
	go func() {
		for range timer.C {
			m.SearchModelInst(params)
			m.SearchModelInstChange(params)
			m.BizHostCountChange(params)
		}
	}()
}

func (m *operationManager) SearchModelInst(ctx core.ContextParams) {
	modelInstCount := make([]metadata.StringIDCount, 0)

	pipeline := []M{{"$group": M{"_id": "$bk_obj_id", "count": M{"$sum": 1}}}}
	if err := m.dbProxy.Table(common.BKTableNameBaseInst).AggregateAll(ctx, pipeline, &modelInstCount); err != nil {
		blog.Errorf("model's instance count aggregate fail, err: %v", err)
		return
	}

	opt := mapstr.MapStr{}
	modelInfo := make([]mapstr.MapStr, 0)
	if err := m.dbProxy.Table(common.BKTableNameObjAttDes).Find(opt).All(ctx, &modelInfo); err != nil {
		blog.Errorf("search model info fail ,err: %v", err)
		return
	}

	modelInstNumber := mapstr.MapStr{}
	for _, countInfo := range modelInstCount {
		for _, model := range modelInfo {
			objID, err := model.String(common.BKObjIDField)
			if err != nil {
				blog.Errorf("model objID interface convert to string fail, err: %v", err)
				continue
			}
			if countInfo.Id == objID {
				objName, err := model.String(common.BKObjNameField)
				if err != nil {
					blog.Errorf("interface convert to string fail, err: %v", err)
					continue
				}
				modelInstNumber[objName] = countInfo.Count
			}
		}
	}

	condition := metadata.ChartData{
		ReportType: "model_inst_chart",
		Data:       modelInstNumber,
		OwnerID:    "0",
	}
	if err := m.dbProxy.Table(common.BKTableNameChartData).Insert(ctx, condition); err != nil {
		blog.Errorf("insert model instance change data fail, err: %v", err)
		return
	}

	return
}

func (m *operationManager) SearchModelInstChange(ctx core.ContextParams) {
	lastTime := time.Now().AddDate(-1, 0, 0)

	createInstCount := make([]metadata.StringIDCount, 0)
	createPipe := []M{{"$match": M{"op_desc": "create object", "op_time": M{"$gte": lastTime}}}, {"$group": M{"_id": "$content.cur_data.bk_obj_id", "count": M{"$sum": 1}}}}
	if err := m.dbProxy.Table(common.BKTableNameOperationLog).AggregateAll(ctx, createPipe, &createInstCount); err != nil {
		blog.Errorf("model's instance count aggregate fail, err: %v", err)
		return
	}

	deleteInstCount := make([]metadata.StringIDCount, 0)
	deletePipe := []M{{"$match": M{"op_desc": "delete object", "op_time": M{"$gte": lastTime}}}, {"$group": M{"_id": "$content.pre_data.bk_obj_id", "count": M{"$sum": 1}}}}
	if err := m.dbProxy.Table(common.BKTableNameOperationLog).AggregateAll(ctx, deletePipe, &deleteInstCount); err != nil {
		blog.Errorf("model's instance count aggregate fail, err: %v", err)
		return
	}

	updateInstCount := make([]metadata.UpdateInstCount, 0)
	updatePipe := []M{{"$match": M{"op_desc": "update object", "op_time": M{"$gte": lastTime}}}, {"$group": M{"_id": M{"bk_obj_id": "$content.cur_data.bk_obj_id", "inst_id": "$content.cur_data.bk_inst_id"}, "count": M{"$sum": 1}}}}
	if err := m.dbProxy.Table(common.BKTableNameOperationLog).AggregateAll(ctx, updatePipe, &updateInstCount); err != nil {
		blog.Errorf("model's instance count aggregate fail, err: %v", err)
		return
	}

	modelInstChange := metadata.ModelInstChange{}
	for _, createInst := range createInstCount {
		if _, ok := modelInstChange[createInst.Id]; ok {
			modelInstChange[createInst.Id].Create += 1
		} else {
			modelInstChange[createInst.Id] = &metadata.InstChangeCount{}
			modelInstChange[createInst.Id].Create = 1
		}
	}

	for _, deleteInst := range deleteInstCount {
		if _, ok := modelInstChange[deleteInst.Id]; ok {
			modelInstChange[deleteInst.Id].Delete += 1
		} else {
			modelInstChange[deleteInst.Id] = &metadata.InstChangeCount{}
			modelInstChange[deleteInst.Id].Delete = 1
		}
	}

	for _, updateInst := range updateInstCount {
		if _, ok := modelInstChange[updateInst.Id.ObjID]; ok {
			modelInstChange[updateInst.Id.ObjID].Update += 1
		} else {
			modelInstChange[updateInst.Id.ObjID] = &metadata.InstChangeCount{}
			modelInstChange[updateInst.Id.ObjID].Update = 1
		}
	}

	condition := metadata.ChartData{
		ReportType: "model_inst_change_chart",
		Data:       modelInstChange,
		OwnerID:    "0",
	}

	opt := mapstr.MapStr{}
	opt["report_type"] = "model_inst_change_chart"
	if err := m.dbProxy.Table(common.BKTableNameChartData).Delete(ctx, opt); err != nil {
		blog.Errorf("delete model instance change data fail, err: %v", err)
		return
	}

	if err := m.dbProxy.Table(common.BKTableNameChartData).Insert(ctx, condition); err != nil {
		blog.Errorf("insert model instance change data fail, err: %v", err)
		return
	}
}

func (m *operationManager) BizHostCountChange(ctx core.ContextParams) {
	data, err := m.SearchBizHost(ctx)
	if err != nil {
		blog.Errorf("search biz host count fail, err: %v", err)
		return
	}

	opt := mapstr.MapStr{}
	bizInfo := make([]mapstr.MapStr, 0)
	if err := m.dbProxy.Table(common.BKTableNameBaseApp).Find(opt).All(ctx, &bizInfo); err != nil {
		blog.Errorf("search biz info fail ,err: %v", err)
		return
	}

	condition := mapstr.MapStr{}
	condition["report_type"] = "host_change_biz_chart"
	bizHostChange := metadata.HostChangeChartData{}
	if err := m.dbProxy.Table(common.BKTableNameChartData).Find(condition).All(ctx, &bizHostChange); err != nil {
		blog.Errorf("get host change data fail, err: %v", err)
		return
	}

	now := time.Now().String()
	for _, info := range data {
		for _, biz := range bizInfo {
			bizID, err := biz.Int64(common.BKAppIDField)
			if err != nil {
				blog.Errorf("biz id interface convert to int64 fail, err: %v", err)
				continue
			}

			bizName, err := biz.String(common.BKAppNameField)
			if err != nil {
				blog.Errorf("biz name interface convert to string fail, err: %v", err)
				continue
			}

			if info.Id == bizID {
				_, ok := bizHostChange.Data[bizName]
				if ok {
					bizHostChange.Data[bizName][now] = info.Count
				} else {
					bizHostChange.Data[bizName] = mapstr.MapStr{}
					bizHostChange.Data[bizName][now] = info.Count
				}
			}
		}
	}

	if err := m.dbProxy.Table(common.BKTableNameChartData).Update(ctx, condition, bizHostChange); err != nil {
		blog.Errorf("update biz host change fail, err: %v", err)
		return
	}
}
