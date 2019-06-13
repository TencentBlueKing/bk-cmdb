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
	"fmt"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/source_controller/coreservice/core"
)

func (m *operationManager) TimerFreshData(params core.ContextParams) {
	m.ModelInst(params)
	m.ModelInstChange(params)
	m.BizHostCountChange(params)
}

func (m *operationManager) ModelInst(ctx core.ContextParams) {
	modelInstCount := make([]metadata.StringIDCount, 0)

	pipeline := []M{{"$group": M{"_id": "$bk_obj_id", "count": M{"$sum": 1}}}}
	if err := m.dbProxy.Table(common.BKTableNameBaseInst).AggregateAll(ctx, pipeline, &modelInstCount); err != nil {
		blog.Errorf("model's instance count aggregate fail, err: %v", err)
		return
	}

	opt := mapstr.MapStr{}
	modelInfo := make([]metadata.Object, 0)
	if err := m.dbProxy.Table(common.BKTableNameObjDes).Find(opt).All(ctx, &modelInfo); err != nil {
		blog.Errorf("search model info fail ,err: %v", err)
		return
	}

	modelInstNumber := make([]mapstr.MapStr, 0)
	for _, countInfo := range modelInstCount {
		for _, model := range modelInfo {
			if countInfo.Id == model.ObjectID {
				info := mapstr.MapStr{}
				info["id"] = model.ObjectName
				info["count"] = countInfo.Count
				modelInstNumber = append(modelInstNumber, info)
			}
		}
	}

	data := metadata.ChartData{
		ReportType: common.ModelInstChart,
		Data:       modelInstNumber,
		OwnerID:    "0",
	}
	condition := mapstr.MapStr{}
	condition[common.OperationReportType] = common.ModelInstChart
	if err := m.dbProxy.Table(common.BKTableNameChartData).Delete(ctx, condition); err != nil {
		blog.Errorf("delete model instance change data fail, err: %v", err)
		return
	}

	if err := m.dbProxy.Table(common.BKTableNameChartData).Insert(ctx, data); err != nil {
		blog.Errorf("insert model instance change data fail, err: %v", err)
		return
	}

	return
}

func (m *operationManager) ModelInstChange(ctx core.ContextParams) {
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
		ReportType: common.ModelInstChangeChart,
		Data:       modelInstChange,
		OwnerID:    "0",
	}

	opt := mapstr.MapStr{}
	opt[common.OperationReportType] = common.ModelInstChangeChart
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
	bizHost, err := m.SearchBizHost(ctx)
	if err != nil {
		blog.Errorf("search biz host count fail, err: %v", err)
		return
	}

	opt := mapstr.MapStr{}
	bizInfo := make([]metadata.BizInst, 0)
	if err := m.dbProxy.Table(common.BKTableNameBaseApp).Find(opt).All(ctx, &bizInfo); err != nil {
		blog.Errorf("search biz info fail ,err: %v", err)
		return
	}

	condition := mapstr.MapStr{}
	condition[common.OperationReportType] = common.HostChangeBizChart
	bizHostChange := make([]metadata.HostChangeChartData, 0)
	if err := m.dbProxy.Table(common.BKTableNameChartData).Find(condition).All(ctx, &bizHostChange); err != nil {
		blog.Errorf("get host change data fail, err: %v", err)
		return
	}

	firstBizHostChange := metadata.HostChangeChartData{}
	now := time.Now().String()
	for _, info := range bizHost {
		for _, biz := range bizInfo {
			if info.Id != biz.BizID {
				continue
			}
			if len(bizHostChange) > 0 {
				_, ok := bizHostChange[0].Data[biz.BizName]
				if ok {
					bizHostChange[0].Data[biz.BizName] = append(bizHostChange[0].Data[biz.BizName], metadata.BizHostChart{
						Id:    now,
						Count: info.Count,
					})
				} else {
					bizHostChange[0].Data = map[string][]metadata.BizHostChart{}
					bizHostChange[0].Data[biz.BizName] = append(bizHostChange[0].Data[biz.BizName], metadata.BizHostChart{
						Id:    now,
						Count: info.Count,
					})
				}
			} else {
				firstBizHostChange.OwnerID = "0"
				firstBizHostChange.ReportType = common.HostChangeBizChart
				_, ok := firstBizHostChange.Data[biz.BizName]
				if ok {
					firstBizHostChange.Data[biz.BizName] = append(firstBizHostChange.Data[biz.BizName], metadata.BizHostChart{
						Id:    now,
						Count: info.Count,
					})
				} else {
					firstBizHostChange.Data = map[string][]metadata.BizHostChart{}
					firstBizHostChange.Data[biz.BizName] = append(firstBizHostChange.Data[biz.BizName], metadata.BizHostChart{
						Id:    now,
						Count: info.Count,
					})
				}
			}
		}
	}

	if len(bizHostChange) > 0 {
		blog.Debug("update info : %v", bizHostChange[0])
		if err := m.dbProxy.Table(common.BKTableNameChartData).Update(ctx, condition, bizHostChange[0]); err != nil {
			blog.Errorf("update biz host change fail, err: %v", err)
			return
		}
	} else {
		if err := m.dbProxy.Table(common.BKTableNameChartData).Insert(ctx, firstBizHostChange); err != nil {
			blog.Errorf("update biz host change fail, err: %v", err)
			return
		}
	}

}

func (m *operationManager) SearchBizHost(ctx core.ContextParams) ([]metadata.IntIDCount, error) {

	bizHostCount := make([]metadata.IntIDCount, 0)

	pipeline := []M{{"$group": M{"_id": "$bk_biz_id", "count": M{"$sum": 1}}}}
	if err := m.dbProxy.Table(common.BKTableNameModuleHostConfig).AggregateAll(ctx, pipeline, &bizHostCount); err != nil {
		blog.Errorf("biz' host count aggregate fail, err: %v", err)
		return nil, err
	}

	return bizHostCount, nil
}

func (m *operationManager) HostCloudChartData(ctx core.ContextParams, inputParam metadata.ChartConfig) (interface{}, error) {
	commonCount := make([]metadata.IntIDCount, 0)
	filterCondition := fmt.Sprintf("$%v", inputParam.Field)

	pipeline := []M{{"$group": M{"_id": filterCondition, "count": M{"$sum": 1}}}}
	if err := m.dbProxy.Table(common.BKTableNameBaseHost).AggregateAll(ctx, pipeline, &commonCount); err != nil {
		blog.Errorf("model's instance count aggregate fail, err: %v", err)
		return nil, err
	}

	opt := mapstr.MapStr{}
	cloudMapping := make([]metadata.CloudMapping, 0)
	if err := m.dbProxy.Table(common.BKTableNameBasePlat).Find(opt).All(ctx, &cloudMapping); err != nil {
		blog.Errorf("search chart config fail, err: %v", err)
		return nil, err
	}

	respData := make([]mapstr.MapStr, 0)
	info := mapstr.MapStr{}
	for _, data := range commonCount {
		for _, cloud := range cloudMapping {
			if data.Id == cloud.CloudID {
				info["id"] = cloud.CloudName
				info["count"] = data.Count
				respData = append(respData, info)
			}
		}
	}

	return respData, nil
}

func (m *operationManager) HostBizChartData(ctx core.ContextParams, inputParam metadata.ChartConfig) (interface{}, error) {
	bizHost, err := m.SearchBizHost(ctx)
	if err != nil {
		blog.Error("search biz host info fail, err: %v", err)
		return nil, err
	}

	opt := mapstr.MapStr{}
	bizInfo := make([]mapstr.MapStr, 0)
	if err := m.dbProxy.Table(common.BKTableNameBaseApp).Find(opt).All(ctx, &bizInfo); err != nil {
		blog.Errorf("get biz info fail, err: %v", err)
		return nil, err
	}

	respData := make([]metadata.StringIDCount, 0)

	for _, biz := range bizInfo {
		id, err := biz.Int64(common.BKAppIDField)
		if err != nil {
			blog.Error("search biz host chart data fail, interface convert to int64 fail, err: %v", err)
			continue
		}
		name, err := biz.String(common.BKAppNameField)
		if err != nil {
			blog.Error("search biz host chart data fail, interface convert to int64 fail, err: %v", err)
			continue
		}
		for _, host := range bizHost {
			if host.Id == id {
				info := metadata.StringIDCount{}
				info.Id = name
				info.Count = host.Count
				respData = append(respData, info)
			}
		}
	}

	return respData, nil
}
