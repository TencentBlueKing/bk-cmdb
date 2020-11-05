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
	"configcenter/src/common/http/rest"
	"configcenter/src/common/json"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/storage/driver/mongodb"
)

func (m *operationManager) TimerFreshData(kit *rest.Kit) error {
	wg := &sync.WaitGroup{}
	wg.Add(3)
	go func(wg *sync.WaitGroup) {
		if err := m.ModelInstCount(kit, wg); err != nil {
			blog.Errorf("TimerFreshData, count model's instance, search model info fail ,err: %v, rid: %v", err)
			return
		}
	}(wg)

	go func(wg *sync.WaitGroup) {
		if err := m.ModelInstChange(kit, wg); err != nil {
			blog.Errorf("TimerFreshData, model inst change count fail, err: %v", err)
			return
		}
	}(wg)

	go func(wg *sync.WaitGroup) {
		if err := m.BizHostCountChange(kit, wg); err != nil {
			blog.Errorf("TimerFreshData fail, biz host change count fail, err: %v", err)
			return
		}
	}(wg)

	wg.Wait()
	return nil
}

// ModelInstCount 统计模型实例数量
func (m *operationManager) ModelInstCount(kit *rest.Kit, wg *sync.WaitGroup) error {
	defer wg.Done()

	// 查询模型 （排除top模型：biz\set\module\process\host\plat）
	innerObject := []string{common.BKInnerObjIDHost, common.BKInnerObjIDApp, common.BKInnerObjIDSet, common.BKInnerObjIDModule, common.BKInnerObjIDProc, common.BKInnerObjIDPlat}
	cond := mapstr.MapStr{}
	cond[common.BKObjIDField] = mapstr.MapStr{common.BKDBNIN: innerObject}
	modelInfos := make([]metadata.Object, 0)
	if err := mongodb.Client().Table(common.BKTableNameObjDes).Find(cond).All(kit.Ctx, &modelInfos); err != nil {
		blog.Errorf("count model's instance, search model info fail ,err: %v, rid: %v", err, kit.Rid)
		return err
	}

	// 判断是否有模型实例
	condition := mapstr.MapStr{}
	count, err := mongodb.Client().Table(common.BKTableNameBaseInst).Find(condition).Count(kit.Ctx)
	if err != nil {
		blog.Errorf("model's instance count fail, err: %v, rid: %v", err, kit.Rid)
		return err
	}

	// 如果有模型实例，根据查询出来的模型，计数
	modelInstNumber := make([]metadata.StringIDCount, 0)
	if count > 0 {
		for _, modelInfo := range modelInfos {
			condition := mapstr.MapStr{
				common.BKObjIDField: modelInfo.ObjectID,
			}
			count, err = mongodb.Client().Table(common.BKTableNameBaseInst).Find(condition).Count(kit.Ctx)
			if err != nil {
				blog.Errorf("model %s's instance count fail, err: %v, rid: %v", modelInfo.ObjectID, err, kit.Rid)
				return err
			}

			modelInstNumber = append(modelInstNumber, metadata.StringIDCount{
				ID:    modelInfo.ObjectID,
				Count: int64(count),
			})
		}
	}

	if err := m.UpdateInnerChartData(kit, common.ModelInstChart, modelInstNumber); err != nil {
		blog.Errorf("update inner chart ModelInstCount data fail, err: %v, rid: %v", err, kit.Rid)
		return err
	}
	return nil
}

func (m *operationManager) ModelInstChange(kit *rest.Kit, wg *sync.WaitGroup) error {
	defer wg.Done()
	operationLog, err := m.StatisticOperationLog(kit)
	if err != nil {
		blog.Errorf("aggregate: count update object fail, err: %v, rid: %v", err, kit.Rid)
		return err
	}

	innerObject := []string{common.BKInnerObjIDHost, common.BKInnerObjIDApp, common.BKInnerObjIDSet, common.BKInnerObjIDModule, common.BKInnerObjIDProc, common.BKInnerObjIDPlat}
	cond := mapstr.MapStr{}
	cond[common.BKObjIDField] = mapstr.MapStr{common.BKDBNIN: innerObject}
	modelData := make([]metadata.Object, 0)
	if err = mongodb.Client().Table(common.BKTableNameObjDes).Find(cond).All(kit.Ctx, &modelData); nil != err {
		blog.Errorf("request(%s): it is failed to find all models by the condition (%#v), error info is %s", kit.Rid, cond, err.Error())
		return err
	}

	modelInstChange := metadata.ModelInstChange{}
	for _, createInst := range operationLog.Create {
		if _, ok := modelInstChange[createInst.ID]; !ok {
			modelInstChange[createInst.ID] = &metadata.InstChangeCount{}
		}
		modelInstChange[createInst.ID].Create = createInst.Count
	}

	for _, deleteInst := range operationLog.Delete {
		if _, ok := modelInstChange[deleteInst.ID]; !ok {
			modelInstChange[deleteInst.ID] = &metadata.InstChangeCount{}
		}
		modelInstChange[deleteInst.ID].Delete = deleteInst.Count
	}

	// 同一个实例更新多次，模型下实例变更数，只需要记录一次
	for _, updateInst := range operationLog.Update {
		if _, ok := modelInstChange[updateInst.ID.ObjID]; ok {
			modelInstChange[updateInst.ID.ObjID].Update += 1
		} else {
			modelInstChange[updateInst.ID.ObjID] = &metadata.InstChangeCount{}
			modelInstChange[updateInst.ID.ObjID].Update = 1
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

	if err = m.UpdateInnerChartData(kit, common.ModelInstChangeChart, modelInstData); err != nil {
		blog.Errorf("update inner chart ModelInstChange data fail, err: %v, rid: %v", err, kit.Rid)
		return err
	}

	return nil
}

// BizHostCountChange 统计业务下主机数量
func (m *operationManager) BizHostCountChange(kit *rest.Kit, wg *sync.WaitGroup) error {
	defer wg.Done()
	bizHost, err := m.BizHostCount(kit)
	if err != nil {
		blog.Errorf("search biz's host count fail, err: %v, rid: %v", err, kit.Rid)
		return err
	}

	// Clear data over 180 days
	go m.clearDataOverDate(kit)

	// 判断今日记录是否统计，如果统计则更新今日数据，否则插入数据
	dateTemplate := "2006-01-02"
	nowStrFormat := time.Unix(time.Now().Unix(), 0).Format(dateTemplate)
	condition := mapstr.MapStr{
		common.OperationReportType: common.HostChangeBizChart,
		common.CreateTimeField:     nowStrFormat,
	}
	bizHostChange := make([]metadata.HostChangeChartData, 0)
	if err = mongodb.Client().Table(common.BKTableNameChartData).Find(condition).All(kit.Ctx, &bizHostChange); err != nil {
		blog.Errorf("get host change data fail, err: %v, rid: %v", err, kit.Rid)
		return err
	}

	if len(bizHostChange) > 0 {
		bizHostChange[0].Data = bizHost
		if err = mongodb.Client().Table(common.BKTableNameChartData).Update(kit.Ctx, condition, bizHostChange[0]); err != nil {
			blog.Errorf("update biz host change chart fail, err: %v, rid: %v", err, kit.Rid)
			return err
		}
		return nil
	} else {
		firstBizHostChange := metadata.HostChangeChartData{
			ReportType: common.HostChangeBizChart,
			Data:       bizHost,
			OwnerID:    kit.SupplierAccount,
			CreateTime: nowStrFormat,
		}
		if err = mongodb.Client().Table(common.BKTableNameChartData).Insert(kit.Ctx, firstBizHostChange); err != nil {
			blog.Errorf("update biz host change fail, err: %v, rid: %v", err, kit.Rid)
			return err
		}
	}
	return nil
}

// BizHostCount 统计业务下主机数量
func (m *operationManager) BizHostCount(kit *rest.Kit) ([]metadata.StringIDCount, error) {
	// 查询所有业务信息 （排除归档业务和资源池）
	cond := mapstr.MapStr{"bk_data_status": M{common.BKDBNE: "disabled"}, common.BKAppIDField: M{common.BKDBNE: 1}}
	bizInfos := make([]metadata.BizInst, 0)
	if err := mongodb.Client().Table(common.BKTableNameBaseApp).Find(cond).All(kit.Ctx, &bizInfos); err != nil {
		blog.Errorf("BizHostCount, get biz info fail, err: %v, rid: %v ", err, kit.Rid)
		return nil, err
	}

	// 判断是否存在业务主机
	condition := mapstr.MapStr{common.BKAppIDField: M{common.BKDBNE: 1}}
	count, err := mongodb.Client().Table(common.BKTableNameModuleHostConfig).Find(condition).Count(kit.Ctx)
	if err != nil {
		blog.Errorf("BizHostCount count: biz' host count fail, err: %v, rid: %v", err, kit.Rid)
		return nil, err
	}

	// 如果存在业务主机，进行统计计数
	bizHostData := make([]metadata.StringIDCount, 0)
	if count > 0 {
		for _, bizInfo := range bizInfos {
			condition := mapstr.MapStr{common.BKAppIDField: bizInfo.BizID}
			count, err := mongodb.Client().Table(common.BKTableNameModuleHostConfig).Find(condition).Count(kit.Ctx)
			if err != nil {
				blog.Errorf("BizHostCount count: biz %s' host count fail, err: %v, rid: %v", bizInfo.BizName, err, kit.Rid)
				return nil, err
			}

			bizHostData = append(bizHostData, metadata.StringIDCount{
				ID:    bizInfo.BizName,
				Count: int64(count),
			})
		}
	}

	return bizHostData, nil
}

func (m *operationManager) HostCloudChartData(kit *rest.Kit, inputParam metadata.ChartConfig) (interface{}, error) {
	commonCount := make([]metadata.IntIDCount, 0)
	filterCondition := fmt.Sprintf("$%s", inputParam.Field)

	respData := make([]metadata.StringIDCount, 0)
	opt := mapstr.MapStr{}
	cloudMapping := make([]metadata.CloudMapping, 0)
	if err := mongodb.Client().Table(common.BKTableNameBasePlat).Find(opt).All(kit.Ctx, &cloudMapping); err != nil {
		blog.Errorf("hostCloudChartData, search cloud mapping fail, err: %v, rid: %v", err, kit.Rid)
		return nil, err
	}
	pipeline := []M{{common.BKDBGroup: M{"_id": filterCondition, "count": M{common.BKDBSum: 1}}}}
	if err := mongodb.Client().Table(common.BKTableNameBaseHost).AggregateAll(kit.Ctx, pipeline, &commonCount); err != nil {
		blog.Errorf("hostCloudChartData, aggregate: model's instance count fail, err: %v, rid: %v", err, kit.Rid)
		return nil, err
	}

	for _, cloud := range cloudMapping {
		info := metadata.StringIDCount{
			ID:    cloud.CloudName,
			Count: 0,
		}
		for _, data := range commonCount {
			if data.ID == cloud.CloudID {
				info.Count = data.Count
			}
		}
		respData = append(respData, info)
	}

	return respData, nil
}

func (m *operationManager) HostBizChartData(kit *rest.Kit, inputParam metadata.ChartConfig) (interface{}, error) {
	bizHost, err := m.BizHostCount(kit)
	if err != nil {
		blog.Error("search biz's host data fail, err: %v, rid: %v", err, kit.Rid)
		return nil, err
	}

	opt := mapstr.MapStr{"bk_data_status": M{common.BKDBNE: "disabled"}, common.BKAppIDField: M{common.BKDBNE: 1}}
	bizInfo := make([]metadata.BizInst, 0)
	if err := mongodb.Client().Table(common.BKTableNameBaseApp).Find(opt).All(kit.Ctx, &bizInfo); err != nil {
		blog.Errorf("HostBizChartData, get biz info fail, err: %v, rid: %v ", err, kit.Rid)
		return nil, err
	}

	return bizHost, nil
}

func (m *operationManager) UpdateInnerChartData(kit *rest.Kit, reportType string, data interface{}) error {
	chartData := metadata.ChartData{
		ReportType: reportType,
		Data:       data,
		OwnerID:    kit.SupplierAccount,
		LastTime:   time.Now(),
	}
	// 此处不用update，因为第一次初始数据的时候会导致数据写不进去
	cond := M{common.OperationReportType: reportType}
	if err := mongodb.Client().Table(common.BKTableNameChartData).Delete(kit.Ctx, cond); err != nil {
		blog.Errorf("update chart %v data fail, err: %v, rid: %v", reportType, err, kit.Rid)
		return err
	}

	if err := mongodb.Client().Table(common.BKTableNameChartData).Insert(kit.Ctx, chartData); err != nil {
		blog.Errorf("update chart %v data fail, err: %v, rid: %v", reportType, err, kit.Rid)
		return err
	}
	return nil
}

// StatisticOperationLog 根据 bk_obj_id 分类统计模型实例操作次数
func (m *operationManager) StatisticOperationLog(kit *rest.Kit) (*metadata.StatisticInstOperation, error) {
	lastTime := time.Now().AddDate(0, 0, -30)

	// 查询模型实例创建操作审记，并根据 bk_obj_id 分类统计数量
	createCond := mapstr.MapStr{
		common.BKActionField: metadata.AuditCreate,
		common.BKOperationTimeField: M{
			common.BKDBGTE: lastTime,
		},
		common.BKAuditTypeField:    metadata.ModelInstanceType,
		common.BKResourceTypeField: metadata.ModelInstanceRes,
	}
	// 查询操作审计
	createAuditLogs := make([]metadata.AuditLog, 0)
	if err := mongodb.Client().Table(common.BKTableNameAuditLog).Find(createCond).All(kit.Ctx, &createAuditLogs); err != nil {
		blog.Errorf("ModelInstanceAuditLogCount: query auditLog, createCond: %+v, err: %v, rid: %s", createCond, err, kit.Rid)
		return nil, err
	}
	createInstCountMap := make(map[string]int64, 0)
	for _, adtlog := range createAuditLogs {
		bytes, err := json.Marshal(adtlog.OperationDetail)
		if err != nil {
			return nil, err
		}
		detail := &metadata.InstanceOpDetail{}
		if err = json.Unmarshal(bytes, detail); err != nil {
			return nil, err
		}

		if _, ok := createInstCountMap[detail.ModelID]; ok {
			createInstCountMap[detail.ModelID] = createInstCountMap[detail.ModelID] + 1
		} else {
			createInstCountMap[detail.ModelID] = 1
		}
	}
	createInstCount := make([]metadata.StringIDCount, len(createInstCountMap))
	for key, value := range createInstCountMap {
		createInstCount = append(createInstCount, metadata.StringIDCount{
			ID:    key,
			Count: value,
		})
	}

	// 查询模型实例删除操作审记，并根据 bk_obj_id 分类统计数量
	deleteCond := mapstr.MapStr{
		common.BKActionField: metadata.AuditDelete,
		common.BKOperationTimeField: M{
			common.BKDBGTE: lastTime,
		},
		common.BKAuditTypeField:    metadata.ModelInstanceType,
		common.BKResourceTypeField: metadata.ModelInstanceRes,
	}
	// 查询操作审计
	deleteAuditLogs := make([]metadata.AuditLog, 0)
	if err := mongodb.Client().Table(common.BKTableNameAuditLog).Find(deleteCond).All(kit.Ctx, &deleteAuditLogs); err != nil {
		blog.Errorf("ModelInstanceAuditLogCount: query auditLog, deleteCond: %+v, err: %v, rid: %s", deleteCond, err, kit.Rid)
		return nil, err
	}
	deleteInstCountMap := make(map[string]int64, 0)
	for _, adtlog := range deleteAuditLogs {
		bytes, err := json.Marshal(adtlog.OperationDetail)
		if err != nil {
			return nil, err
		}
		detail := &metadata.InstanceOpDetail{}
		if err = json.Unmarshal(bytes, detail); err != nil {
			return nil, err
		}

		if _, ok := deleteInstCountMap[detail.ModelID]; ok {
			deleteInstCountMap[detail.ModelID] = deleteInstCountMap[detail.ModelID] + 1
		} else {
			deleteInstCountMap[detail.ModelID] = 1
		}
	}
	deleteInstCount := make([]metadata.StringIDCount, len(deleteInstCountMap))
	for key, value := range deleteInstCountMap {
		deleteInstCount = append(deleteInstCount, metadata.StringIDCount{
			ID:    key,
			Count: value,
		})
	}

	// 查询模型实例更新操作审记，并根据 bk_obj_id， 分类统计数量
	updateCond := mapstr.MapStr{
		common.BKActionField: metadata.AuditUpdate,
		common.BKOperationTimeField: M{
			common.BKDBGTE: lastTime,
		},
		common.BKAuditTypeField:    metadata.ModelInstanceType,
		common.BKResourceTypeField: metadata.ModelInstanceRes,
	}
	// 查询操作审计
	updateAuditLogs := make([]metadata.AuditLog, 0)
	if err := mongodb.Client().Table(common.BKTableNameAuditLog).Find(updateCond).All(kit.Ctx, &updateAuditLogs); err != nil {
		blog.Errorf("ModelInstanceAuditLogCount: query auditLog, updateCond: %+v, err: %v, rid: %s", updateCond, err, kit.Rid)
		return nil, err
	}

	updateInstCountMap := make(map[string]map[int64]int64, 0)
	for _, adtlog := range updateAuditLogs {
		bytes, err := json.Marshal(adtlog.OperationDetail)
		if err != nil {
			return nil, err
		}
		detail := &metadata.InstanceOpDetail{}
		if err = json.Unmarshal(bytes, detail); err != nil {
			return nil, err
		}

		instID, err := util.GetInt64ByInterface(adtlog.ResourceID)
		if err != nil {
			return nil, err
		}
		if _, ok := updateInstCountMap[detail.ModelID]; !ok {
			updateInstCountMap[detail.ModelID] = make(map[int64]int64, 0)
		}
		if _, ok := updateInstCountMap[detail.ModelID][instID]; !ok {
			updateInstCountMap[detail.ModelID][instID] = updateInstCountMap[detail.ModelID][instID] + 1
		} else {
			updateInstCountMap[detail.ModelID][instID] = 1
		}
	}

	updateInstCount := make([]metadata.UpdateInstCount, len(updateInstCountMap))
	for objID, innerMap := range updateInstCountMap {
		for instID, count := range innerMap {
			updateInstCount = append(updateInstCount, metadata.UpdateInstCount{
				ID: metadata.UpdateID{
					ObjID:  objID,
					InstID: instID,
				},
				Count: count,
			})
		}
	}

	result := &metadata.StatisticInstOperation{
		Create: createInstCount,
		Delete: deleteInstCount,
		Update: updateInstCount,
	}
	return result, nil
}

// clearDataOverDate Clear biz host change data over 180 days
func (m *operationManager) clearDataOverDate(kit *rest.Kit) {
	cond := mapstr.MapStr{}
	cond[common.OperationReportType] = common.HostChangeBizChart
	bizHostChange := make([]metadata.HostChangeChartData, 0)
	if err := mongodb.Client().Table(common.BKTableNameChartData).Find(cond).All(kit.Ctx, &bizHostChange); err != nil {
		blog.Errorf("get host change data fail, err: %v, rid: %v", err, kit.Rid)
		return
	}

	shouldClear := make([]string, 0)
	now := time.Now()
	for _, info := range bizHostChange {
		dateFormat := "2006-01-02"
		loc, _ := time.LoadLocation("Asia/Shanghai")
		createTime, _ := time.ParseInLocation(dateFormat, info.CreateTime, loc)

		duration := now.Sub(createTime).Hours() / 24
		if duration > 180 {
			shouldClear = append(shouldClear, info.CreateTime)
		}
	}

	cond[common.CreateTimeField] = mapstr.MapStr{common.BKDBIN: shouldClear}
	if err := mongodb.Client().Table(common.BKTableNameChartData).Delete(kit.Ctx, cond); err != nil {
		blog.Errorf("get host change data fail, err: %v, rid: %v", err, kit.Rid)
		return
	}

	return
}
