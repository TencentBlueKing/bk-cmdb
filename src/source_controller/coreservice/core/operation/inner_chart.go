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
	"configcenter/src/common/util"
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

	opt := mapstr.MapStr{}
	modelInstNumber := make([]metadata.IDStringCountInt64, 0)
	modelInfo := make([]metadata.Object, 0)
	if err := m.dbProxy.Table(common.BKTableNameObjDes).Find(opt).All(ctx, &modelInfo); err != nil {
		blog.Errorf("count model's instance, search model info fail ,err: %v, rid: %v", err, ctx.ReqID)
		return
	}

	instCount, err := m.dbProxy.Table(common.BKTableNameBaseInst).Find(opt).Count(ctx)
	if err != nil {
		blog.Errorf("model's instance count aggregate fail, err: %v, rid: %v", err, ctx.ReqID)
		return
	}
	if instCount == 0 {
		if err := m.ObjectInstEmpty(ctx, modelInfo); err != nil {
			blog.Errorf("model's instance count aggregate fail, err: %v, rid: %v", err, ctx.ReqID)
			return
		}
		return
	}

	pipeline := []M{{"$group": M{"_id": "$bk_obj_id", "count": M{"$sum": 1}}}}
	if err := m.dbProxy.Table(common.BKTableNameBaseInst).AggregateAll(ctx, pipeline, &modelInstCount); err != nil {
		blog.Errorf("model's instance count aggregate fail, err: %v, rid: %v", err, ctx.ReqID)
		return
	}

	allModels := make([]metadata.ObjectIDName, 0)
	matchedModels := make([]string, 0)
	for _, model := range modelInfo {
		allModels = append(allModels, metadata.ObjectIDName{ObjectID: model.ObjectID, ObjectName: model.ObjectName})
		for _, countInfo := range modelInstCount {
			if countInfo.Id == model.ObjectID {
				info := metadata.IDStringCountInt64{
					Id:    model.ObjectName,
					Count: countInfo.Count,
				}
				modelInstNumber = append(modelInstNumber, info)
				matchedModels = append(matchedModels, model.ObjectName)
			}
		}
	}

	for _, model := range allModels {
		if !util.InStrArr(matchedModels, model.ObjectName) && !util.IsInnerObject(model.ObjectID) {
			info := metadata.IDStringCountInt64{
				Id:    model.ObjectName,
				Count: 0,
			}
			modelInstNumber = append(modelInstNumber, info)
		}
	}

	data := metadata.ChartData{
		ReportType: common.ModelInstChart,
		Data:       modelInstNumber,
		OwnerID:    "0",
	}
	condition := mapstr.MapStr{}
	condition[common.OperationReportType] = common.ModelInstChart
	if err := m.UpdateInnerChartData(ctx, condition, data); err != nil {
		blog.Errorf("update inner chart ModelInst data fail, err: %v, rid: %v", err, ctx.ReqID)
		return
	}

	return
}

func (m *operationManager) ModelInstChange(ctx core.ContextParams) {
	result, err := m.StatisticOperationLog(ctx)
	if err != nil {
		blog.Errorf("aggregate: count update object fail, err: %v, rid: %v", err, ctx.ReqID)
		return
	}

	cond := mapstr.MapStr{}
	modelData := make([]metadata.Object, 0)
	if err := m.dbProxy.Table(common.BKTableNameObjDes).Find(cond).All(ctx, &modelData); nil != err {
		blog.Errorf("request(%s): it is failed to find all models by the condition (%#v), error info is %s", ctx.ReqID, cond, err.Error())
		return
	}

	modelInstChange := metadata.ModelInstChange{}
	for _, createInst := range result.Create {
		if _, ok := modelInstChange[createInst.Id]; !ok {
			modelInstChange[createInst.Id] = &metadata.InstChangeCount{}
		}
		modelInstChange[createInst.Id].Create = createInst.Count
	}

	for _, deleteInst := range result.Delete {
		if _, ok := modelInstChange[deleteInst.Id]; !ok {
			modelInstChange[deleteInst.Id] = &metadata.InstChangeCount{}
		}
		modelInstChange[deleteInst.Id].Delete = deleteInst.Count
	}

	// 同一个实例更新多次，模型下实例变更数，只需要记录一次
	for _, updateInst := range result.Update {
		if _, ok := modelInstChange[updateInst.Id.ObjID]; ok {
			modelInstChange[updateInst.Id.ObjID].Update += 1
		} else {
			modelInstChange[updateInst.Id.ObjID] = &metadata.InstChangeCount{}
			modelInstChange[updateInst.Id.ObjID].Update = 1
		}
	}

	modelInstData := metadata.ModelInstChange{}
	// 把bk_obj_id换成bk_obj_name
	allModels := make([]metadata.ObjectIDName, 0)
	matchedModels := make([]string, 0)
	for _, model := range modelData {
		allModels = append(allModels, metadata.ObjectIDName{ObjectID: model.ObjectID, ObjectName: model.ObjectName})
		for key, value := range modelInstChange {
			if key == model.ObjectID {
				matchedModels = append(matchedModels, model.ObjectName)
				modelInstData[model.ObjectName] = value
			}
		}
	}

	for _, model := range allModels {
		if !util.InStrArr(matchedModels, model.ObjectName) && !util.IsInnerObject(model.ObjectID) {
			modelInstData[model.ObjectName] = &metadata.InstChangeCount{}
		}
	}

	data := metadata.ChartData{
		ReportType: common.ModelInstChangeChart,
		Data:       modelInstData,
		OwnerID:    "0",
	}
	condition := mapstr.MapStr{}
	condition[common.OperationReportType] = common.ModelInstChangeChart
	if err := m.UpdateInnerChartData(ctx, condition, data); err != nil {
		blog.Errorf("update inner chart ModelInstChange data fail, err: %v, rid: %v", err, ctx.ReqID)
		return
	}

	return
}

func (m *operationManager) BizHostCountChange(ctx core.ContextParams) {
	bizHost, err := m.SearchBizHost(ctx)

	opt := mapstr.MapStr{}
	bizInfo := make([]metadata.BizInst, 0)
	if err := m.dbProxy.Table(common.BKTableNameBaseApp).Find(opt).All(ctx, &bizInfo); err != nil {
		blog.Errorf("biz's host count, search biz info fail ,err: %v, rid: %v", err, ctx.ReqID)
		return
	}

	if bizHost == nil {
		blog.V(3).Info("table cc_ModuleHOstConfig is empty, rid: %v", ctx.ReqID)

		m.BizHostEmpty(ctx, bizInfo)
		return
	}
	if err != nil {
		blog.Errorf("search biz's host count fail, err: %v, rid: %v", err, ctx.ReqID)
		return
	}

	condition := mapstr.MapStr{}
	condition[common.OperationReportType] = common.HostChangeBizChart
	bizHostChange := make([]metadata.HostChangeChartData, 0)
	if err := m.dbProxy.Table(common.BKTableNameChartData).Find(condition).All(ctx, &bizHostChange); err != nil {
		blog.Errorf("get host change data fail, err: %v, rid: %v", err, ctx.ReqID)
		return
	}

	firstBizHostChange := metadata.HostChangeChartData{}
	now := time.Now()

	for _, info := range bizHost {
		for _, biz := range bizInfo {
			if info.Id != biz.BizID {
				continue
			}
			if len(bizHostChange) > 0 {
				subHour := now.Sub(bizHostChange[0].UpdateTime).Hours()
				if subHour < 24 {
					blog.V(3).Info("Less than 24 hours since the last update, return")
					return
				}
				if len(bizHostChange[0].Data) > 0 {
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
				bizHostChange[0].UpdateTime = time.Now()
			} else {
				if len(firstBizHostChange.Data) > 0 {
					firstBizHostChange.Data[biz.BizName] = append(firstBizHostChange.Data[biz.BizName], metadata.BizHostChart{
						Id:    now,
						Count: info.Count,
					})
				} else {
					firstBizHostChange.OwnerID = "0"
					firstBizHostChange.ReportType = common.HostChangeBizChart
					firstBizHostChange.Data = map[string][]metadata.BizHostChart{}
					firstBizHostChange.Data[biz.BizName] = append(firstBizHostChange.Data[biz.BizName], metadata.BizHostChart{
						Id:    now,
						Count: info.Count,
					})
				}
				firstBizHostChange.UpdateTime = time.Now()
			}
		}
	}

	if len(bizHostChange) > 0 {
		if err := m.dbProxy.Table(common.BKTableNameChartData).Update(ctx, condition, bizHostChange[0]); err != nil {
			blog.Errorf("update biz host change chart fail, err: %v, rid: %v", err, ctx.ReqID)
			return
		}
	} else {
		if err := m.dbProxy.Table(common.BKTableNameChartData).Insert(ctx, firstBizHostChange); err != nil {
			blog.Errorf("update biz host change fail, err: %v, rid: %v", err, ctx.ReqID)
			return
		}
	}
}

func (m *operationManager) SearchBizHost(ctx core.ContextParams) ([]metadata.IntIDCount, error) {
	bizHostCount := make([]metadata.IntIDCount, 0)

	cond := mapstr.MapStr{}
	hostCount, err := m.dbProxy.Table(common.BKTableNameModuleHostConfig).Find(cond).Count(ctx)
	if err != nil {
		blog.Errorf("aggregate: biz' host count fail, err: %v, rid: %v", err, ctx.ReqID)
		return nil, err
	}

	if hostCount > 0 {
		pipeline := []M{{"$group": M{"_id": "$bk_biz_id", "count": M{"$sum": 1}}}}
		if err := m.dbProxy.Table(common.BKTableNameModuleHostConfig).AggregateAll(ctx, pipeline, &bizHostCount); err != nil {
			blog.Errorf("aggregate: biz' host count fail, err: %v, rid: %v", err, ctx.ReqID)
			return nil, err
		}
		return bizHostCount, nil
	}

	return nil, nil
}

func (m *operationManager) HostCloudChartData(ctx core.ContextParams, inputParam metadata.ChartConfig) (interface{}, error) {
	commonCount := make([]metadata.IntIDCount, 0)
	filterCondition := fmt.Sprintf("$%v", inputParam.Field)

	pipeline := []M{{"$group": M{"_id": filterCondition, "count": M{"$sum": 1}}}}
	if err := m.dbProxy.Table(common.BKTableNameBaseHost).AggregateAll(ctx, pipeline, &commonCount); err != nil {
		blog.Errorf("hostCloudChartData, aggregate: model's instance count fail, err: %v, rid: %v", err, ctx.ReqID)
		return nil, err
	}

	opt := mapstr.MapStr{}
	cloudMapping := make([]metadata.CloudMapping, 0)
	if err := m.dbProxy.Table(common.BKTableNameBasePlat).Find(opt).All(ctx, &cloudMapping); err != nil {
		blog.Errorf("hostCloudChartData, search cloud mapping fail, err: %v, rid: %v", err, ctx.ReqID)
		return nil, err
	}

	respData := make([]metadata.StringIDCount, 0)
	matched := make([]string, 0)
	for _, data := range commonCount {
		for _, cloud := range cloudMapping {
			if data.Id == cloud.CloudID {
				info := metadata.StringIDCount{
					Id:    cloud.CloudName,
					Count: data.Count,
				}
				matched = append(matched, cloud.CloudName)
				respData = append(respData, info)
			}
		}
	}

	for _, cloud := range cloudMapping {
		if !util.InStrArr(matched, cloud.CloudName) {
			info := metadata.StringIDCount{
				Id:    cloud.CloudName,
				Count: 0,
			}
			respData = append(respData, info)
		}
	}
	return respData, nil
}

func (m *operationManager) HostBizChartData(ctx core.ContextParams, inputParam metadata.ChartConfig) (interface{}, error) {
	bizHost, err := m.SearchBizHost(ctx)
	if err != nil {
		blog.Error("search biz's host data fail, err: %v, rid: %v", err, ctx.ReqID)
		return nil, err
	}

	respData := make([]metadata.StringIDCount, 0)
	opt := mapstr.MapStr{"bk_data_status": M{"$ne": "disabled"}}
	bizInfo := make([]metadata.BizInst, 0)
	if err := m.dbProxy.Table(common.BKTableNameBaseApp).Find(opt).All(ctx, &bizInfo); err != nil {
		blog.Errorf("HostBizChartData, get biz info fail, err: %v, rid: %v ", err, ctx.ReqID)
		return nil, err
	}

	if bizHost == nil {
		blog.V(3).Info("table cc_ModuleHOstConfig is empty, rid: %v", ctx.ReqID)

		for _, biz := range bizInfo {
			info := metadata.StringIDCount{
				Id:    biz.BizName,
				Count: 0,
			}
			respData = append(respData, info)
		}

		return respData, err
	}

	allBiz := make([]string, 0)
	matchedBiz := make([]string, 0)
	for _, biz := range bizInfo {
		allBiz = append(allBiz, biz.BizName)
		for _, host := range bizHost {
			if host.Id == biz.BizID {
				info := metadata.StringIDCount{
					Id:    biz.BizName,
					Count: host.Count,
				}
				respData = append(respData, info)
				matchedBiz = append(matchedBiz, biz.BizName)
			}
		}
	}

	for _, biz := range allBiz {
		if !util.InStrArr(matchedBiz, biz) {
			info := metadata.StringIDCount{
				Id:    biz,
				Count: 0,
			}
			respData = append(respData, info)
		}
	}
	return respData, nil
}

func (m *operationManager) UpdateInnerChartData(ctx core.ContextParams, opt mapstr.MapStr, data metadata.ChartData) error {
	if err := m.dbProxy.Table(common.BKTableNameChartData).Delete(ctx, opt); err != nil {
		blog.Errorf("delete model instance change chart data fail, err: %v, rid: %v", err, ctx.ReqID)
		return err
	}

	if err := m.dbProxy.Table(common.BKTableNameChartData).Insert(ctx, data); err != nil {
		blog.Errorf("insert model instance change chart data fail, err: %v, rid: %v", err, ctx.ReqID)
		return err
	}
	return nil
}

func (m *operationManager) StatisticOperationLog(ctx core.ContextParams) (*metadata.StatisticInstOperation, error) {
	lastTime := time.Now().AddDate(0, 0, -30)

	opt := mapstr.MapStr{}
	opt["op_desc"] = "create object"
	createCount, err := m.dbProxy.Table(common.BKTableNameOperationLog).Find(opt).Count(ctx)
	if err != nil {
		blog.Errorf("aggregate: count create object fail, err: %v, rid", err, ctx.ReqID)
		return nil, err
	}
	createInstCount := make([]metadata.StringIDCount, 0)
	if createCount > 0 {
		createPipe := []M{{"$match": M{"op_desc": "create object", "op_time": M{"$gte": lastTime}}}, {"$group": M{"_id": "$content.cur_data.bk_obj_id", "count": M{"$sum": 1}}}}
		if err := m.dbProxy.Table(common.BKTableNameOperationLog).AggregateAll(ctx, createPipe, &createInstCount); err != nil {
			blog.Errorf("aggregate: count create object fail, err: %v, rid", err, ctx.ReqID)
			return nil, err
		}
	}

	opt["op_desc"] = "delete object"
	deleteCount, err := m.dbProxy.Table(common.BKTableNameOperationLog).Find(opt).Count(ctx)
	if err != nil {
		blog.Errorf("aggregate: count create object fail, err: %v, rid", err, ctx.ReqID)
		return nil, err
	}
	deleteInstCount := make([]metadata.StringIDCount, 0)
	if deleteCount > 0 {
		deletePipe := []M{{"$match": M{"op_desc": "delete object", "op_time": M{"$gte": lastTime}}}, {"$group": M{"_id": "$content.pre_data.bk_obj_id", "count": M{"$sum": 1}}}}
		if err := m.dbProxy.Table(common.BKTableNameOperationLog).AggregateAll(ctx, deletePipe, &deleteInstCount); err != nil {
			blog.Errorf("aggregate: count delete object fail, err: %v, rid: %v", err, ctx.ReqID)
			return nil, err
		}
	}

	opt["op_desc"] = "update object"
	updateCount, err := m.dbProxy.Table(common.BKTableNameOperationLog).Find(opt).Count(ctx)
	if err != nil {
		blog.Errorf("aggregate: count create object fail, err: %v, rid", err, ctx.ReqID)
		return nil, err
	}
	updateInstCount := make([]metadata.UpdateInstCount, 0)
	if updateCount > 0 {
		updatePipe := []M{{"$match": M{"op_desc": "update object", "op_time": M{"$gte": lastTime}}}, {"$group": M{"_id": M{"bk_obj_id": "$content.cur_data.bk_obj_id", "inst_id": "$content.cur_data.bk_inst_id"}, "count": M{"$sum": 1}}}}
		if err := m.dbProxy.Table(common.BKTableNameOperationLog).AggregateAll(ctx, updatePipe, &updateInstCount); err != nil {
			blog.Errorf("aggregate: count update object fail, err: %v, rid: %v", err, ctx.ReqID)
			return nil, err
		}
	}

	result := &metadata.StatisticInstOperation{
		Create: createInstCount,
		Delete: deleteInstCount,
		Update: updateInstCount,
	}

	return result, nil
}

// BizHostEmpty cc_ModuleHOstConfig为空的情况下, 统计业务下主机为0
func (m *operationManager) BizHostEmpty(ctx core.ContextParams, bizInfo []metadata.BizInst) {
	firstBizHostChange := metadata.HostChangeChartData{}
	now := time.Now()

	for _, biz := range bizInfo {
		if len(firstBizHostChange.Data) > 0 {
			firstBizHostChange.Data[biz.BizName] = append(firstBizHostChange.Data[biz.BizName], metadata.BizHostChart{
				Id:    now,
				Count: 0,
			})
		} else {
			firstBizHostChange.OwnerID = "0"
			firstBizHostChange.ReportType = common.HostChangeBizChart
			firstBizHostChange.Data = map[string][]metadata.BizHostChart{}
			firstBizHostChange.Data[biz.BizName] = append(firstBizHostChange.Data[biz.BizName], metadata.BizHostChart{
				Id:    now,
				Count: 0,
			})
		}
		firstBizHostChange.UpdateTime = time.Now()
	}

	if err := m.dbProxy.Table(common.BKTableNameChartData).Insert(ctx, firstBizHostChange); err != nil {
		blog.Errorf("update biz host change fail, err: %v, rid: %v", err, ctx.ReqID)
		return
	}
}

func (m *operationManager) ObjectInstEmpty(ctx core.ContextParams, modelInfo []metadata.Object) error {
	modelInstNumber := make([]metadata.IDStringCountInt64, 0)
	for _, model := range modelInfo {
		info := metadata.IDStringCountInt64{
			Id:    model.ObjectName,
			Count: 0,
		}
		modelInstNumber = append(modelInstNumber, info)
	}
	data := metadata.ChartData{
		ReportType: common.ModelInstChart,
		Data:       modelInstNumber,
		OwnerID:    "0",
	}
	condition := mapstr.MapStr{}
	condition[common.OperationReportType] = common.ModelInstChart
	if err := m.UpdateInnerChartData(ctx, condition, data); err != nil {
		blog.Errorf("update inner chart ModelInst data fail, err: %v, rid: %v", err, ctx.ReqID)
		return err
	}
	return nil
}
