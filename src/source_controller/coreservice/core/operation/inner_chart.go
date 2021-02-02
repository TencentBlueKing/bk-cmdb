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
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/storage/driver/mongodb"
)

const (
	AuditLogInstanceOpDetailModelIDField = "operation_detail.bk_obj_id"
	AuditLogResourceIDField              = "resource_id"
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

// ModelInstCount 统计各模型的实例数量
func (m *operationManager) ModelInstCount(kit *rest.Kit, wg *sync.WaitGroup) error {
	defer wg.Done()

	// 查询模型 （排除top模型：biz\set\module\process\host\plat）
	innerObject := []string{common.BKInnerObjIDHost, common.BKInnerObjIDApp, common.BKInnerObjIDSet, common.BKInnerObjIDModule, common.BKInnerObjIDProc, common.BKInnerObjIDPlat}
	cond := mapstr.MapStr{}
	cond[common.BKObjIDField] = mapstr.MapStr{common.BKDBNIN: innerObject}
	modelInfos := []map[string]interface{}{}
	if err := mongodb.Client().Table(common.BKTableNameObjDes).Find(cond).Fields(common.BKObjIDField).All(kit.Ctx, &modelInfos); err != nil {
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
				common.BKObjIDField: util.GetStrByInterface(modelInfo[common.BKObjIDField]),
			}
			count, err = mongodb.Client().Table(common.BKTableNameBaseInst).Find(condition).Count(kit.Ctx)
			if err != nil {
				blog.Errorf("model %s's instance count fail, err: %v, rid: %v", modelInfo[common.BKObjIDField], err, kit.Rid)
				return err
			}

			modelInstNumber = append(modelInstNumber, metadata.StringIDCount{
				ID:    util.GetStrByInterface(modelInfo[common.BKObjIDField]),
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

	innerObject := []string{common.BKInnerObjIDHost, common.BKInnerObjIDApp, common.BKInnerObjIDSet, common.BKInnerObjIDModule, common.BKInnerObjIDProc, common.BKInnerObjIDPlat}
	cond := mapstr.MapStr{}
	cond[common.BKObjIDField] = mapstr.MapStr{common.BKDBNIN: innerObject}
	fields := []string{common.BKObjIDField, common.BKObjNameField}
	modelData := []map[string]interface{}{}
	if err = mongodb.Client().Table(common.BKTableNameObjDes).Find(cond).Fields(fields...).All(kit.Ctx, &modelData); nil != err {
		blog.Errorf("request(%s): it is failed to find all models by the condition (%#v), error info is %s", kit.Rid, cond, err.Error())
		return err
	}

	modelInstData := metadata.ModelInstChange{}
	// 把bk_obj_id换成bk_obj_name
	for _, model := range modelData {
		objName := util.GetStrByInterface(model[common.BKObjNameField])
		objID := util.GetStrByInterface(model[common.BKObjIDField])
		modelInstData[objName] = &metadata.InstChangeCount{}
		for key, value := range modelInstChange {
			if key == objID {
				modelInstData[objName] = value
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
	}

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
	return nil
}

// BizHostCount 统计业务下主机数量
func (m *operationManager) BizHostCount(kit *rest.Kit) ([]metadata.StringIDCount, error) {
	// 获取需要统计的业务ID和业务名对应关系 （排除归档业务和资源池）
	cond := mapstr.MapStr{
		"bk_data_status":      M{common.BKDBNE: "disabled"},
		common.BKDefaultField: M{common.BKDBNE: common.DefaultAppFlag},
	}
	bizInfos := []map[string]interface{}{}
	fields := []string{common.BKAppIDField, common.BKAppNameField}
	if err := mongodb.Client().Table(common.BKTableNameBaseApp).Find(cond).Fields(fields...).All(kit.Ctx, &bizInfos); err != nil {
		blog.Errorf("BizHostCount failed, find err: %v, cond:%#v, rid: %v ", err, cond, kit.Rid)
		return nil, err
	}
	if len(bizInfos) == 0 {
		return []metadata.StringIDCount{}, nil
	}

	bizIDs := make([]int64, len(bizInfos))
	bizIDNameMap := make(map[int64]string)
	for idx, bizInfo := range bizInfos {
		bizID, err := util.GetInt64ByInterface(bizInfo[common.BKAppIDField])
		if err != nil {
			blog.Errorf("BizHostCount failed, GetInt64ByInterface err: %v, bizInfo:%#v, rid: %v ", err, bizInfo, kit.Rid)
			return nil, err
		}
		bizIDs[idx] = bizID
		bizName := util.GetStrByInterface(bizInfo[common.BKAppNameField])
		bizIDNameMap[bizID] = bizName
	}

	// 统计各个业务下的主机数量
	bizCountArr := make([]metadata.IntIDCount, 0)
	pipeline := []M{
		{common.BKDBMatch: M{common.BKAppIDField: M{common.BKDBIN: bizIDs}}},
		{common.BKDBGroup: M{"_id": "$bk_biz_id", "uniqueHosts": M{common.BKDBAddToSet: "$bk_host_id"}}},
		{common.BKDBProject: M{
			"_id":   1,
			"count": M{common.BKDBSize: "$uniqueHosts"},
		},
		},
	}
	if err := mongodb.Client().Table(common.BKTableNameModuleHostConfig).AggregateAll(kit.Ctx, pipeline, &bizCountArr); err != nil {
		blog.Errorf("BizHostCount failed, err: %v, rid: %v", err, kit.Rid)
		return nil, err
	}
	if len(bizCountArr) == 0 {
		return []metadata.StringIDCount{}, nil
	}
	bizHostCntMap := make(map[int64]int64)
	for _, bizCount := range bizCountArr {
		bizHostCntMap[bizCount.ID] = bizCount.Count
	}

	ret := make([]metadata.StringIDCount, 0)
	for bizID, bizName := range bizIDNameMap {
		ret = append(ret, metadata.StringIDCount{
			ID:    bizName,
			Count: bizHostCntMap[bizID],
		})
	}

	return ret, nil
}

func (m *operationManager) HostCloudChartData(kit *rest.Kit, inputParam metadata.ChartConfig) (interface{}, error) {
	commonCount := make([]metadata.IntIDCount, 0)
	groupField := fmt.Sprintf("$%s", inputParam.Field)

	respData := make([]metadata.StringIDCount, 0)
	opt := mapstr.MapStr{}
	cloudMapping := make([]metadata.CloudMapping, 0)
	if err := mongodb.Client().Table(common.BKTableNameBasePlat).Find(opt).All(kit.Ctx, &cloudMapping); err != nil {
		blog.Errorf("hostCloudChartData, search cloud mapping fail, err: %v, rid: %v", err, kit.Rid)
		return nil, err
	}
	pipeline := []M{{common.BKDBGroup: M{"_id": groupField, "count": M{common.BKDBSum: 1}}}}
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

// StatisticOperationLog
// Note: 根据 bk_obj_id 分类统计模型实例操作次数，统计时间为前一天零点，到当天零点
func (m *operationManager) StatisticOperationLog(kit *rest.Kit) (*metadata.StatisticInstOperation, error) {
	zeroTime := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(),
		0, 0, 0, 0, time.Now().Location())
	old1Time := zeroTime.AddDate(0, 0, -1)
	// 每次分页查询的步长
	limit := 2000

	// 查询模型实例创建操作审记，并根据 bk_obj_id 分类统计数量
	createInstCountMap := make(map[string]int64, 0)
	createCond := mapstr.MapStr{
		common.BKAuditTypeField:    metadata.ModelInstanceType,
		common.BKResourceTypeField: metadata.ModelInstanceRes,
		common.BKActionField:       metadata.AuditCreate,
		common.BKOperationTimeField: M{
			common.BKDBGTE: old1Time,
			common.BKDBLT:  zeroTime,
		},
	}
	fields := []string{AuditLogInstanceOpDetailModelIDField}
	for start := 0; ; start = start + limit {
		// 查询操作审计
		createAuditLogs := []map[string]interface{}{}
		if err := mongodb.Client().Table(common.BKTableNameAuditLog).Find(createCond).Fields(fields...).Start(uint64(start)).
			Limit(uint64(limit)).All(kit.Ctx, &createAuditLogs); err != nil {
			blog.Errorf("ModelInstanceAuditLogCount: query auditLog, createCond: %+v, err: %v, rid: %s", createCond, err, kit.Rid)
			return nil, err
		}
		// 判断当前类型是否查完
		if len(createAuditLogs) == 0 {
			break
		}

		for _, adtlog := range createAuditLogs {
			detail := adtlog[common.BKOperationDetailField].(map[string]interface{})
			moduleID := util.GetStrByInterface(detail[common.BKObjIDField])
			if _, ok := createInstCountMap[moduleID]; ok {
				createInstCountMap[moduleID] = createInstCountMap[moduleID] + 1
			} else {
				createInstCountMap[moduleID] = 1
			}
		}
	}

	// 查询模型实例删除操作审记，并根据 bk_obj_id 分类统计数量
	deleteCond := mapstr.MapStr{
		common.BKAuditTypeField:    metadata.ModelInstanceType,
		common.BKResourceTypeField: metadata.ModelInstanceRes,
		common.BKActionField:       metadata.AuditDelete,
		common.BKOperationTimeField: M{
			common.BKDBGTE: old1Time,
			common.BKDBLT:  zeroTime,
		},
	}
	fields = []string{AuditLogInstanceOpDetailModelIDField}
	deleteInstCountMap := make(map[string]int64, 0)
	for start := 0; ; start = start + limit {
		// 查询操作审计
		deleteAuditLogs := []map[string]interface{}{}
		if err := mongodb.Client().Table(common.BKTableNameAuditLog).Find(deleteCond).Fields(fields...).Start(uint64(start)).
			Limit(uint64(limit)).All(kit.Ctx, &deleteAuditLogs); err != nil {
			blog.Errorf("ModelInstanceAuditLogCount: query auditLog, deleteCond: %+v, err: %v, rid: %s", deleteCond, err, kit.Rid)
			return nil, err
		}

		// 判断当前类型是否查完
		if len(deleteAuditLogs) == 0 {
			break
		}

		for _, adtlog := range deleteAuditLogs {
			detail := adtlog[common.BKOperationDetailField].(map[string]interface{})
			moduleID := util.GetStrByInterface(detail[common.BKObjIDField])
			if _, ok := deleteInstCountMap[moduleID]; ok {
				deleteInstCountMap[moduleID] = deleteInstCountMap[moduleID] + 1
			} else {
				deleteInstCountMap[moduleID] = 1
			}
		}
	}

	// 查询模型实例更新操作审记，并根据 bk_obj_id， 分类统计数量
	updateCond := mapstr.MapStr{
		common.BKAuditTypeField:    metadata.ModelInstanceType,
		common.BKResourceTypeField: metadata.ModelInstanceRes,
		common.BKActionField:       metadata.AuditUpdate,
		common.BKOperationTimeField: M{
			common.BKDBGTE: old1Time,
			common.BKDBLT:  zeroTime,
		},
	}
	fields = []string{AuditLogInstanceOpDetailModelIDField, AuditLogResourceIDField}
	updateInstCountMap := make(map[string]map[int64]int64, 0)
	for start := 0; ; start = start + limit {
		// 查询操作审计
		updateAuditLogs := []map[string]interface{}{}
		if err := mongodb.Client().Table(common.BKTableNameAuditLog).Find(updateCond).Fields(fields...).Start(uint64(start)).
			Limit(uint64(limit)).All(kit.Ctx, &updateAuditLogs); err != nil {
			blog.Errorf("ModelInstanceAuditLogCount: query auditLog, updateCond: %+v, err: %v, rid: %s", updateCond, err, kit.Rid)
			return nil, err
		}

		// 判断当前类型是否查完
		if len(updateAuditLogs) == 0 {
			break
		}

		for _, adtlog := range updateAuditLogs {
			detail := adtlog[common.BKOperationDetailField].(map[string]interface{})

			moduleID := util.GetStrByInterface(detail[common.BKObjIDField])
			instID, err := util.GetInt64ByInterface(adtlog[AuditLogResourceIDField])
			if err != nil {
				return nil, err
			}
			if _, ok := updateInstCountMap[moduleID]; !ok {
				updateInstCountMap[moduleID] = make(map[int64]int64, 0)
			}
			if _, ok := updateInstCountMap[moduleID][instID]; !ok {
				updateInstCountMap[moduleID][instID] = updateInstCountMap[moduleID][instID] + 1
			} else {
				updateInstCountMap[moduleID][instID] = 1
			}
		}
	}

	createInstCount := make([]metadata.StringIDCount, 0)
	deleteInstCount := make([]metadata.StringIDCount, 0)
	updateInstCount := make([]metadata.UpdateInstCount, 0)
	for key, value := range createInstCountMap {
		createInstCount = append(createInstCount, metadata.StringIDCount{
			ID:    key,
			Count: value,
		})
	}
	for key, value := range deleteInstCountMap {
		deleteInstCount = append(deleteInstCount, metadata.StringIDCount{
			ID:    key,
			Count: value,
		})
	}
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
