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
	"sync"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/source_controller/coreservice/core"
)

func (m *operationManager) TimerFreshData(params core.ContextParams) error {

	wg := &sync.WaitGroup{}
	wg.Add(3)
	go func(wg *sync.WaitGroup) {
		if err := m.ModelInst(params, wg); err != nil {
			blog.Errorf("TimerFreshData, count model's instance, search model info fail ,err: %v, rid: %v", err)
			return
		}
	}(wg)

	go func(wg *sync.WaitGroup) {
		if err := m.ModelInstChange(params, wg); err != nil {
			blog.Errorf("TimerFreshData, model inst change count fail, err: %v", err)
			return
		}
	}(wg)

	go func(wg *sync.WaitGroup) {
		if err := m.BizHostCountChange(params, wg); err != nil {
			blog.Errorf("TimerFreshData fail, biz host change count fail, err: %v", err)
			return
		}
	}(wg)

	wg.Wait()
	return nil
}

func (m *operationManager) ModelInst(ctx core.ContextParams, wg *sync.WaitGroup) error {
	defer wg.Done()
	modelInstCount := make([]metadata.StringIDCount, 0)

	innerObject := []string{"host", "biz", "set", "cloud", "module", "process", "plat"}
	cond := mapstr.MapStr{}
	cond["bk_obj_id"] = mapstr.MapStr{"$nin": innerObject}
	modelInstNumber := make([]metadata.IDStringCountInt64, 0)
	modelInfo := make([]metadata.Object, 0)
	if err := m.dbProxy.Table(common.BKTableNameObjDes).Find(cond).All(ctx, &modelInfo); err != nil {
		blog.Errorf("count model's instance, search model info fail ,err: %v, rid: %v", err, ctx.ReqID)
		return err
	}
	pipeline := []M{{"$group": M{"_id": "$bk_obj_id", "count": M{"$sum": 1}}}}
	if err := m.dbProxy.Table(common.BKTableNameBaseInst).AggregateAll(ctx, pipeline, &modelInstCount); err != nil {
		blog.Errorf("model's instance count aggregate fail, err: %v, rid: %v", err, ctx.ReqID)
		return err
	}

	for _, model := range modelInfo {
		info := metadata.IDStringCountInt64{
			Id:    model.ObjectName,
			Count: 0,
		}
		for _, instCount := range modelInstCount {
			if instCount.Id == model.ObjectID {
				info.Count = instCount.Count
			}
		}
		modelInstNumber = append(modelInstNumber, info)
	}

	if err := m.UpdateInnerChartData(ctx, common.ModelInstChart, modelInstNumber); err != nil {
		blog.Errorf("update inner chart ModelInst data fail, err: %v, rid: %v", err, ctx.ReqID)
		return err
	}

	return nil
}

func (m *operationManager) ModelInstChange(ctx core.ContextParams, wg *sync.WaitGroup) error {
	defer wg.Done()
	operationLog, err := m.StatisticOperationLog(ctx)
	if err != nil {
		blog.Errorf("aggregate: count update object fail, err: %v, rid: %v", err, ctx.ReqID)
		return err
	}

	innerObject := []string{"host", "biz", "set", "cloud", "module", "process", "plat"}
	cond := mapstr.MapStr{}
	cond[common.BKObjIDField] = mapstr.MapStr{"$nin": innerObject}
	modelData := make([]metadata.Object, 0)
	if err = m.dbProxy.Table(common.BKTableNameObjDes).Find(cond).All(ctx, &modelData); nil != err {
		blog.Errorf("request(%s): it is failed to find all models by the condition (%#v), error info is %s", ctx.ReqID, cond, err.Error())
		return err
	}

	modelInstChange := metadata.ModelInstChange{}
	for _, createInst := range operationLog.Create {
		if _, ok := modelInstChange[createInst.Id]; !ok {
			modelInstChange[createInst.Id] = &metadata.InstChangeCount{}
		}
		modelInstChange[createInst.Id].Create = createInst.Count
	}

	for _, deleteInst := range operationLog.Delete {
		if _, ok := modelInstChange[deleteInst.Id]; !ok {
			modelInstChange[deleteInst.Id] = &metadata.InstChangeCount{}
		}
		modelInstChange[deleteInst.Id].Delete = deleteInst.Count
	}

	// 同一个实例更新多次，模型下实例变更数，只需要记录一次
	for _, updateInst := range operationLog.Update {
		if _, ok := modelInstChange[updateInst.Id.ObjID]; ok {
			modelInstChange[updateInst.Id.ObjID].Update += 1
		} else {
			modelInstChange[updateInst.Id.ObjID] = &metadata.InstChangeCount{}
			modelInstChange[updateInst.Id.ObjID].Update = 1
		}
	}

	modelInstData := metadata.ModelInstChange{}
	// 把bk_obj_id换成bk_obj_name
	for _, model := range modelData {
		modelInstData[model.ObjectName] = &metadata.InstChangeCount{}
		for key, value := range modelInstChange {
			if key == model.ObjectID {
				modelInstData[model.ObjectName] = value
			}
		}
	}

	if err = m.UpdateInnerChartData(ctx, common.ModelInstChangeChart, modelInstData); err != nil {
		blog.Errorf("update inner chart ModelInstChange data fail, err: %v, rid: %v", err, ctx.ReqID)
		return err
	}

	return nil
}

func (m *operationManager) BizHostCountChange(ctx core.ContextParams, wg *sync.WaitGroup) error {
	defer wg.Done()
	bizHost, err := m.SearchBizHost(ctx)
	if err != nil {
		blog.Errorf("search biz's host count fail, err: %v, rid: %v", err, ctx.ReqID)
		return err
	}

	opt := M{"bk_data_status": M{"$ne": "disabled"}, "bk_biz_id": M{"$ne": 1}}
	bizInfo := make([]metadata.BizInst, 0)
	if err := m.dbProxy.Table(common.BKTableNameBaseApp).Find(opt).All(ctx, &bizInfo); err != nil {
		blog.Errorf("biz's host count, search biz info fail ,err: %v, rid: %v", err, ctx.ReqID)
		return err
	}

	condition := mapstr.MapStr{}
	condition[common.OperationReportType] = common.HostChangeBizChart
	bizHostChange := make([]metadata.HostChangeChartData, 0)
	if err = m.dbProxy.Table(common.BKTableNameChartData).Find(condition).All(ctx, &bizHostChange); err != nil {
		blog.Errorf("get host change data fail, err: %v, rid: %v", err, ctx.ReqID)
		return err
	}

	firstBizHostChange := metadata.HostChangeChartData{}
	now := time.Now()

	for _, info := range bizHost {
		currentData := metadata.BizHostChart{}
		currentData.Id.Time = now
		currentData.Count = info.Count
		if len(bizHostChange) > 0 {
			subHour := now.Sub(bizHostChange[0].LastTime.Time).Hours()
			if subHour < 24 {
				blog.V(3).Info("Less than 24 hours since the last update, return")
				return nil
			}
			if len(bizHostChange[0].Data) > 0 {
				bizHostChange[0].Data[info.Id] = append(bizHostChange[0].Data[info.Id], currentData)
			} else {
				bizHostChange[0].Data = map[string][]metadata.BizHostChart{}
				bizHostChange[0].Data[info.Id] = append(bizHostChange[0].Data[info.Id], currentData)
			}
			bizHostChange[0].LastTime.Time = time.Now()
		} else {
			if len(firstBizHostChange.Data) > 0 {
				firstBizHostChange.Data[info.Id] = append(firstBizHostChange.Data[info.Id], currentData)
			} else {
				firstBizHostChange.OwnerID = ctx.SupplierAccount
				firstBizHostChange.ReportType = common.HostChangeBizChart
				firstBizHostChange.Data = map[string][]metadata.BizHostChart{}
				firstBizHostChange.Data[info.Id] = append(firstBizHostChange.Data[info.Id], currentData)
			}
			firstBizHostChange.LastTime.Time = time.Now()
		}
	}

	if len(bizHostChange) > 0 {
		if err = m.dbProxy.Table(common.BKTableNameChartData).Update(ctx, condition, bizHostChange[0]); err != nil {
			blog.Errorf("update biz host change chart fail, err: %v, rid: %v", err, ctx.ReqID)
			return err
		}
	} else {
		if err = m.dbProxy.Table(common.BKTableNameChartData).Insert(ctx, firstBizHostChange); err != nil {
			blog.Errorf("update biz host change fail, err: %v, rid: %v", err, ctx.ReqID)
			return err
		}
	}
	return nil
}

func (m *operationManager) SearchBizHost(ctx core.ContextParams) ([]metadata.StringIDCount, error) {
	bizHostCount := make([]metadata.IntIDArrayCount, 0)

	opt := mapstr.MapStr{"bk_data_status": M{"$ne": "disabled"}, "bk_biz_id": M{"$ne": 1}}
	bizInfo := make([]metadata.BizInst, 0)
	if err := m.dbProxy.Table(common.BKTableNameBaseApp).Find(opt).All(ctx, &bizInfo); err != nil {
		blog.Errorf("HostBizChartData, get biz info fail, err: %v, rid: %v ", err, ctx.ReqID)
		return nil, err
	}

	pipeline := []M{{"$group": M{"_id": "$bk_biz_id", "count": M{"$addToSet": "$bk_host_id"}}}}
	if err := m.dbProxy.Table(common.BKTableNameModuleHostConfig).AggregateAll(ctx, pipeline, &bizHostCount); err != nil {
		blog.Errorf("aggregate: biz' host count fail, err: %v, rid: %v", err, ctx.ReqID)
		return nil, err
	}

	rData := make([]metadata.StringIDCount, 0)
	for _, biz := range bizInfo {
		info := metadata.StringIDCount{
			Id:    biz.BizName,
			Count: 0,
		}
		for _, host := range bizHostCount {
			if host.Id == biz.BizID {
				info.Count = int64(len(host.Count))
			}
		}
		rData = append(rData, info)
	}

	return rData, nil
}

func (m *operationManager) HostCloudChartData(ctx core.ContextParams, inputParam metadata.ChartConfig) (interface{}, error) {
	commonCount := make([]metadata.IntIDCount, 0)
	filterCondition := fmt.Sprintf("$%v", inputParam.Field)

	respData := make([]metadata.StringIDCount, 0)
	opt := mapstr.MapStr{}
	cloudMapping := make([]metadata.CloudMapping, 0)
	if err := m.dbProxy.Table(common.BKTableNameBasePlat).Find(opt).All(ctx, &cloudMapping); err != nil {
		blog.Errorf("hostCloudChartData, search cloud mapping fail, err: %v, rid: %v", err, ctx.ReqID)
		return nil, err
	}
	pipeline := []M{{"$group": M{"_id": filterCondition, "count": M{"$sum": 1}}}}
	if err := m.dbProxy.Table(common.BKTableNameBaseHost).AggregateAll(ctx, pipeline, &commonCount); err != nil {
		blog.Errorf("hostCloudChartData, aggregate: model's instance count fail, err: %v, rid: %v", err, ctx.ReqID)
		return nil, err
	}

	for _, cloud := range cloudMapping {
		info := metadata.StringIDCount{
			Id:    cloud.CloudName,
			Count: 0,
		}
		for _, data := range commonCount {
			if data.Id == cloud.CloudID {
				info.Count = data.Count
			}
		}
		respData = append(respData, info)
	}

	return respData, nil
}

func (m *operationManager) HostBizChartData(ctx core.ContextParams, inputParam metadata.ChartConfig) (interface{}, error) {
	bizHost, err := m.SearchBizHost(ctx)
	if err != nil {
		blog.Error("search biz's host data fail, err: %v, rid: %v", err, ctx.ReqID)
		return nil, err
	}

	opt := mapstr.MapStr{"bk_data_status": M{"$ne": "disabled"}, "bk_biz_id": M{"$ne": 1}}
	bizInfo := make([]metadata.BizInst, 0)
	if err := m.dbProxy.Table(common.BKTableNameBaseApp).Find(opt).All(ctx, &bizInfo); err != nil {
		blog.Errorf("HostBizChartData, get biz info fail, err: %v, rid: %v ", err, ctx.ReqID)
		return nil, err
	}

	return bizHost, nil
}

func (m *operationManager) UpdateInnerChartData(ctx core.ContextParams, reportType string, data interface{}) error {
	chartData := metadata.ChartData{
		ReportType: reportType,
		Data:       data,
		OwnerID:    ctx.SupplierAccount,
		LastTime:   time.Now(),
	}
	// 此处不用update，因为第一次初始数据的时候会导致数据写不进去
	cond := M{common.OperationReportType: reportType}
	if err := m.dbProxy.Table(common.BKTableNameChartData).Delete(ctx, cond); err != nil {
		blog.Errorf("update chart %v data fail, err: %v, rid: %v", reportType, err, ctx.ReqID)
		return err
	}

	if err := m.dbProxy.Table(common.BKTableNameChartData).Insert(ctx, chartData); err != nil {
		blog.Errorf("update chart %v data fail, err: %v, rid: %v", reportType, err, ctx.ReqID)
		return err
	}
	return nil
}

func (m *operationManager) StatisticOperationLog(ctx core.ContextParams) (*metadata.StatisticInstOperation, error) {
	lastTime := time.Now().AddDate(0, 0, -30)

	innerObject := []string{"host", "biz", "set", "cloud", "module", "process", "plat"}
	opt := mapstr.MapStr{}
	opt[common.OperationDescription] = common.CreateObject
	createCount, err := m.dbProxy.Table(common.BKTableNameOperationLog).Find(opt).Count(ctx)
	if err != nil {
		blog.Errorf("aggregate: count create object fail, err: %v, rid", err, ctx.ReqID)
		return nil, err
	}
	createInstCount := make([]metadata.StringIDCount, 0)
	if createCount > 0 {
		createPipe := []M{{"$match": M{"op_type": 1, "op_time": M{"$gte": lastTime}, "op_target": M{"$nin": innerObject}}}, {"$group": M{"_id": "$op_target", "count": M{"$sum": 1}}}}
		if err := m.dbProxy.Table(common.BKTableNameOperationLog).AggregateAll(ctx, createPipe, &createInstCount); err != nil {
			blog.Errorf("aggregate: count create object fail, err: %v, rid", err, ctx.ReqID)
			return nil, err
		}
	}

	opt[common.OperationDescription] = common.DeleteObject
	deleteCount, err := m.dbProxy.Table(common.BKTableNameOperationLog).Find(opt).Count(ctx)
	if err != nil {
		blog.Errorf("aggregate: count delete object fail, err: %v, rid", err, ctx.ReqID)
		return nil, err
	}
	deleteInstCount := make([]metadata.StringIDCount, 0)
	if deleteCount > 0 {
		deletePipe := []M{{"$match": M{"op_type": 3, "op_time": M{"$gte": lastTime}, "op_target": M{"$nin": innerObject}}}, {"$group": M{"_id": "$op_target", "count": M{"$sum": 1}}}}
		if err := m.dbProxy.Table(common.BKTableNameOperationLog).AggregateAll(ctx, deletePipe, &deleteInstCount); err != nil {
			blog.Errorf("aggregate: count delete object fail, err: %v, rid: %v", err, ctx.ReqID)
			return nil, err
		}
	}

	opt[common.OperationDescription] = common.UpdateObject
	updateCount, err := m.dbProxy.Table(common.BKTableNameOperationLog).Find(opt).Count(ctx)
	if err != nil {
		blog.Errorf("aggregate: count create object fail, err: %v, rid", err, ctx.ReqID)
		return nil, err
	}
	updateInstCount := make([]metadata.UpdateInstCount, 0)
	if updateCount > 0 {
		updatePipe := []M{{"$match": M{"op_type": 2, "op_time": M{"$gte": lastTime}, "op_target": M{"$nin": innerObject}}},
			{"$group": M{"_id": M{"bk_obj_id": "$op_target", "inst_id": "$inst_id"}, "count": M{"$sum": 1}}}}
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
