/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package logics

import (
	"context"
	"fmt"
	"net/http"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

func (lgc *Logics) SearchReportSummary(header http.Header, param metadata.ParamSearchNetcollectReport) ([]*metadata.NetcollectReportSummary, error) {
	// search reports
	param.CloudID = -1
	_, reports, err := lgc.SearchReport(header, param)
	if err != nil {
		blog.Errorf("[NetDevice][SearchReportSummary] SearchReport by %+v failed %v", param, err)
		return nil, err
	}
	summarym := map[int64]*metadata.NetcollectReportSummary{}

	// combine details
	for _, report := range reports {
		summary, ok := summarym[report.CloudID]
		if !ok {
			summary = &metadata.NetcollectReportSummary{
				CloudID:    report.CloudID,
				CloudName:  report.CloudName,
				Statistics: map[string]int{},
			}
			summarym[report.CloudID] = summary
		}

		summary.Statistics["associations"] += len(report.Associations)
		summary.Statistics[report.ObjectName]++

		if report.LastTime.Time.Sub(summary.LastTime.Time) > 0 {
			summary.LastTime = report.LastTime
		}
	}

	summarys := []*metadata.NetcollectReportSummary{}
	for key := range summarym {
		summarys = append(summarys, summarym[key])
	}

	return summarys, nil
}

func (lgc *Logics) buildSearchCond(header http.Header, param metadata.ParamSearchNetcollectReport) (condition.Condition, error) {
	cond := condition.CreateCondition()
	if param.CloudID >= 0 {
		cond.Field(common.BKCloudIDField).Eq(param.CloudID)
	}
	cloudIDs := []int64{}
	if param.CloudName != "" || param.Query != "" {
		cloudCond := condition.CreateCondition()
		cloudCond.Field(common.BKCloudNameField).Like(param.CloudName)
		clouds, err := lgc.findInst(header, common.BKInnerObjIDPlat, &metadata.QueryInput{Condition: cloudCond.ToMapStr()})
		if err != nil {
			return nil, err
		}
		for _, cloud := range clouds {
			id, err := cloud.Int64(common.BKCloudIDField)
			if err != nil {
				return nil, err
			}
			cloudIDs = append(cloudIDs, id)
		}
	}
	if param.CloudName != "" {
		cond.Field(common.BKCloudIDField).In(cloudIDs)
	}
	if param.Action != "" {
		cond.Field("action").Eq(param.Action)
	}
	if param.ObjectID != "" {
		cond.Field(common.BKObjIDField).Eq(param.ObjectID)
	}
	if param.InnerIP != "" {
		cond.Field(common.BKHostInnerIPField).Like(param.InnerIP)
	}
	if param.Query != "" {
		cond.Field(common.BKDBOR).Eq([]map[string]interface{}{
			{
				common.BKHostInnerIPField: map[string]interface{}{
					common.BKDBLIKE: param.InnerIP,
				},
			},
			{
				common.BKCloudIDField: map[string]interface{}{
					common.BKDBIN: cloudIDs,
				},
			},
		})
	}
	if len(param.LastTime) >= 2 {
		cond.Field(common.LastTimeField).Gte(param.LastTime[0])
		cond.Field(common.LastTimeField).Lte(param.LastTime[1])
	}
	return cond, nil
}

func (lgc *Logics) SearchReport(header http.Header, param metadata.ParamSearchNetcollectReport) (uint64, []metadata.NetcollectReport, error) {
	cond, err := lgc.buildSearchCond(header, param)
	if err != nil {
		blog.Errorf("[NetDevice][SearchReport] build SearchReport condition for %+v, failed: %v", param, err)
		return 0, nil, err
	}
	count, err := lgc.Instance.Table(common.BKTableNameNetcollectReport).Find(cond.ToMapStr()).Count(lgc.ctx)
	if err != nil {
		blog.Errorf("[NetDevice][SearchReport] GetCntByCondition %+v failed: %v", cond.ToMapStr(), err)
		return 0, nil, err
	}

	// search reports
	reports := []metadata.NetcollectReport{}
	err = lgc.Instance.Table(common.BKTableNameNetcollectReport).Find(cond.ToMapStr()).Sort(param.Page.Sort).Start(uint64(param.Page.Start)).Limit(uint64(param.Page.Limit)).All(lgc.ctx, &reports)
	if err != nil {
		blog.Errorf("[NetDevice][SearchReport] GetMutilByCondition %+v failed: %v", cond.ToMapStr(), err)
		return 0, nil, err
	}

	// search details
	objIDs := []string{}
	cloudIDs := []int64{}
	for _, report := range reports {
		objIDs = append(objIDs, report.ObjectID)
		cloudIDs = append(cloudIDs, report.CloudID)

		for _, asst := range report.Associations {
			objIDs = append(objIDs, asst.AsstObjectID)
		}
	}

	cloudCond := condition.CreateCondition()
	cloudCond.Field(common.BKCloudIDField).In(cloudIDs)
	cloudMap, err := lgc.findInstMap(header, common.BKInnerObjIDPlat, &metadata.QueryInput{Condition: cloudCond.ToMapStr()})
	if err != nil {
		blog.Errorf("[NetDevice][SearchReport] find clouds by %+v failed: %v", cloudCond, err)
		return 0, nil, err
	}

	objMap, err := lgc.findObjectMap(header, objIDs...)
	if err != nil {
		blog.Errorf("[NetDevice][SearchReport] findObjectMap by %+v failed: %v", objIDs, err)
		return 0, nil, err
	}
	attrsMap, err := lgc.findAttrsMap(header, objIDs...)
	if err != nil {
		blog.Errorf("[NetDevice][SearchReport] findAttrsMap by %+v failed: %v", objIDs, err)
		return 0, nil, err
	}

	for index := range reports {
		if object, ok := objMap[reports[index].ObjectID]; ok {
			reports[index].ObjectName = object.ObjectName
		}
		if clodInst, ok := cloudMap[reports[index].CloudID]; ok {
			cloudname, err := clodInst.String(common.BKCloudNameField)
			if err != nil {
				blog.Errorf("[NetDevice][SearchReport] bk_cloud_name field invalied: %v", err)
			}
			reports[index].CloudName = cloudname
		}

		cond := condition.CreateCondition()
		cond.Field(common.BKInstNameField).Eq(reports[index].InstKey)
		objType := common.GetObjByType(reports[index].ObjectID)
		if objType == common.BKInnerObjIDObject {
			cond.Field(common.BKObjIDField).Eq(reports[index].ObjectID)
		}
		insts, err := lgc.findInst(header, reports[index].ObjectID, &metadata.QueryInput{Condition: cond.ToMapStr()})
		if err != nil {
			blog.Errorf("[NetDevice][SearchReport] find inst by %+v for %v failed %v", cond.ToMapStr(), reports[index].ObjectID, err)
			return 0, nil, err
		}
		for attrIndex := range reports[index].Attributes {
			attribute := &reports[index].Attributes[attrIndex]
			if property, ok := attrsMap[attrMapKey(reports[index].ObjectID, attribute.PropertyID)]; ok {
				attribute.PropertyName = property.PropertyName
				attribute.IsRequired = property.IsRequired
			}
			if len(insts) > 0 {
				attribute.PreValue = insts[0][attribute.PropertyID]
				reports[index].Action = metadata.ReporctActionUpdate
			} else {
				reports[index].Action = metadata.ReporctActionCreate
			}
		}

		for asstIndex := range reports[index].Associations {
			asst := &reports[index].Associations[asstIndex]
			asst.Action = metadata.ReporctActionCreate
			if object, ok := objMap[asst.AsstObjectID]; ok {
				asst.AsstObjectName = object.ObjectName
			}
		}

	}

	return count, reports, nil
}

func attrMapKey(objectID string, propertyID string) string {
	return objectID + ":" + propertyID
}

func (lgc *Logics) findAttrsMap(header http.Header, objIDs ...string) (map[string]metadata.Attribute, error) {
	attrs, err := lgc.findAttrs(header, objIDs...)
	if err != nil {
		return nil, err
	}

	attrsMap := map[string]metadata.Attribute{}
	for _, attr := range attrs {
		attrsMap[attrMapKey(attr.ObjectID, attr.PropertyID)] = attr
	}
	return attrsMap, nil
}

func (lgc *Logics) findAttrs(header http.Header, objIDs ...string) ([]metadata.Attribute, error) {
	cond := condition.CreateCondition()
	cond.Field(common.BKObjIDField).In(objIDs)
	resp, err := lgc.CoreAPI.ObjectController().Meta().SelectObjectAttWithParams(context.Background(), header, cond.ToMapStr())
	if err != nil {
		blog.Errorf("[NetDevice][findAttrs] for %v error: %v", objIDs, err)
		return nil, err
	}
	if !resp.Result {
		blog.Errorf("[NetDevice][findAttrs] for %v error: %v", objIDs, resp)
		return nil, err
	}
	return resp.Data, nil
}

func (lgc *Logics) findObjectMap(header http.Header, objIDs ...string) (map[string]metadata.Object, error) {
	objs, err := lgc.findObjectIn(header, objIDs...)
	if err != nil {
		return nil, err
	}

	objectMap := map[string]metadata.Object{}
	for _, obj := range objs {
		objectMap[obj.ObjectID] = obj
	}
	return objectMap, nil
}

func (lgc *Logics) findObjectIn(header http.Header, objIDs ...string) ([]metadata.Object, error) {
	cond := condition.CreateCondition()
	cond.Field(common.BKObjIDField).In(objIDs)
	return lgc.findObject(header, cond.ToMapStr())
}

func (lgc *Logics) findObject(header http.Header, cond interface{}) ([]metadata.Object, error) {
	resp, err := lgc.CoreAPI.ObjectController().Meta().SelectObjects(context.Background(), header, cond)
	if err != nil {
		blog.Errorf("[NetDevice][findObject] by %+v error: %v", cond, err)
		return nil, err
	}
	if !resp.Result {
		blog.Errorf("[NetDevice][findObject] by %+v error: %v", cond, resp)
		return nil, err
	}
	return resp.Data, nil
}

func (lgc *Logics) findInstMap(header http.Header, objectID string, query *metadata.QueryInput) (map[int64]mapstr.MapStr, error) {
	insts, err := lgc.findInst(header, objectID, query)
	if err != nil {
		return nil, err
	}

	instMap := map[int64]mapstr.MapStr{}
	for _, inst := range insts {
		id, err := inst.Int64(common.GetInstIDField(objectID))
		if err != nil {
			return nil, err
		}
		instMap[id] = inst
	}
	return instMap, nil
}

func (lgc *Logics) findInst(header http.Header, objectID string, query *metadata.QueryInput) ([]mapstr.MapStr, error) {
	resp, err := lgc.CoreAPI.ObjectController().Instance().SearchObjects(context.Background(), objectID, header, query)
	if err != nil {
		blog.Errorf("[NetDevice][findInst] by %+v error: %v", query, err)
		return nil, err
	}
	if !resp.Result {
		blog.Errorf("[NetDevice][findInst] by %+v error: %+v", query, resp)
		return nil, err
	}
	return resp.Data.Info, nil
}

func (lgc *Logics) ConfirmReport(header http.Header, reports []metadata.NetcollectReport) *metadata.RspNetcollectConfirm {
	result := metadata.RspNetcollectConfirm{Errors: []string{}}
	for index := range reports {
		report := &reports[index]
		if len(report.Attributes) > 0 {
			attrCount, err := lgc.confirmAttributes(header, report)
			if err != nil {
				result.ChangeAttributeFailure += attrCount
				result.Errors = append(result.Errors, err.Error())
				lgc.saveHistory(report, false)
				continue
			}
			result.ChangeAttributeSuccess += attrCount
			lgc.saveHistory(report, true)
		}
		if len(report.Associations) > 0 {
			successCount, errs := lgc.confirmAssociations(header, report)
			result.ChangeAssociationsFailure += len(errs)
			result.ChangeAssociationsSuccess += successCount
			if len(errs) > 0 {
				for _, err := range errs {
					result.Errors = append(result.Errors, err.Error())
				}
				lgc.saveHistory(report, false)
				continue
			}
			lgc.saveHistory(report, true)
		}
	}
	return &result
}

func (lgc *Logics) confirmAttributes(header http.Header, report *metadata.NetcollectReport) (int, error) {
	data := mapstr.MapStr{}
	attrCount := 0
	for _, attr := range report.Attributes {
		if attr.Method == metadata.ReporctMethodAccept {
			data.Set(attr.PropertyID, attr.CurValue)
			attrCount++
		}
	}

	if len(data) <= 0 {
		blog.Warnf("[NetDevice][ConfirmReport] empty data, continue next")
		return attrCount, nil
	}

	objType := common.GetObjByType(report.ObjectID)
	if objType == common.BKInnerObjIDHost {
		hostdata := metadata.HostList{
			InputType: metadata.CollectType,
			HostInfo: map[int64]map[string]interface{}{
				1: data,
			},
		}

		resp, err := lgc.CoreAPI.HostServer().AddHost(context.Background(), header, &hostdata)
		if err != nil {
			blog.Errorf("[NetDevice][ConfirmReport] update inst error: %v, %+v", err, data)
			return attrCount, err
		}
		if !resp.Result {
			blog.Errorf("[NetDevice][ConfirmReport] update inst error: %v, %+v", resp, data)
			return attrCount, fmt.Errorf(resp.ErrMsg)
		}
	}

	cond := condition.CreateCondition()
	if objType == common.BKInnerObjIDObject {
		cond.Field(common.GetInstNameField(report.ObjectID)).Eq(report.InstKey)
		cond.Field(common.BKObjIDField).Eq(report.ObjectID)
	}

	if objType == common.BKInnerObjIDHost {
		cond.Field(common.BKCloudIDField).Eq(report.CloudID)
		cond.Field(common.BKHostInnerIPField).Eq(report.InstKey)
	}
	insts, err := lgc.findInst(header, report.ObjectID, &metadata.QueryInput{Condition: cond.ToMapStr()})
	if err != nil {
		blog.Errorf("[NetDevice][ConfirmReport] find inst failed %v", err)
		return attrCount, err
	}
	blog.V(4).Infof("[NetDevice][ConfirmReport] find inst result: %#v, condition: %#v", insts, cond.ToMapStr())
	if len(insts) > 0 {
		instID, err := insts[0].Int64(common.GetInstIDField(report.ObjectID))
		if err != nil {
			blog.Errorf("[NetDevice][ConfirmReport] find inst failed, instID not found from %+v", insts[0])
			return 0, err
		}

		resp, err := lgc.CoreAPI.TopoServer().Instance().UpdateInst(context.Background(), util.GetOwnerID(header), report.ObjectID, instID, header, data)
		if err != nil {
			blog.Errorf("[NetDevice][ConfirmReport] update inst error: %v, %+v", err, data)
			return attrCount, err
		}
		if !resp.Result {
			blog.Errorf("[NetDevice][ConfirmReport] update inst error: %v, %+v", resp, data)
			return attrCount, fmt.Errorf(resp.ErrMsg)
		}
	} else {
		resp, err := lgc.CoreAPI.TopoServer().Instance().CreateInst(context.Background(), util.GetOwnerID(header), report.ObjectID, header, data)
		if err != nil {
			blog.Errorf("[NetDevice][ConfirmReport] create inst to %s error: %v, %+v", report.ObjectID, err, data)
			return attrCount, err
		}
		if !resp.Result {
			blog.Errorf("[NetDevice][ConfirmReport] create inst to %s error: %v, %+v", report.ObjectID, resp, data)
			return attrCount, fmt.Errorf(resp.ErrMsg)
		}
	}
	return attrCount, nil
}

func (lgc *Logics) confirmAssociations(header http.Header, report *metadata.NetcollectReport) (successCount int, errs []error) {
	objType := common.GetObjByType(report.ObjectID)
	cond := condition.CreateCondition()
	if objType == common.BKInnerObjIDObject {
		cond.Field(common.GetInstNameField(report.ObjectID)).Eq(report.InstKey)
		cond.Field(common.BKObjIDField).Eq(report.ObjectID)
	}
	if objType == common.BKInnerObjIDHost {
		cond.Field(common.BKCloudIDField).Eq(report.CloudID)
		cond.Field(common.BKHostInnerIPField).Eq(report.InstKey)
	}

	insts, err := lgc.findInst(header, report.ObjectID, &metadata.QueryInput{Condition: cond.ToMapStr()})
	if err != nil {
		blog.Errorf("[NetDevice][ConfirmReport] find inst %+v failed %v", cond.ToMapStr(), err)
		return 0, append(errs, err)
	}
	if len(insts) <= 0 {
		blog.Errorf("[NetDevice][ConfirmReport] find inst failed, inst not found by %+v", cond.ToMapStr())
		return 0, append(errs, fmt.Errorf("inst not found"))
	}
	instID, err := insts[0].Int64(common.GetInstIDField(report.ObjectID))
	if err != nil {
		blog.Errorf("[NetDevice][ConfirmReport] find inst failed, instID not found from %+v", insts[0])
		return 0, append(errs, fmt.Errorf("inst not found"))
	}

	for _, asst := range report.Associations {
		asstObjType := common.GetObjByType(asst.AsstObjectID)
		asstCond := condition.CreateCondition()
		if asstObjType == common.BKInnerObjIDObject {
			asstCond.Field(common.GetInstNameField(asst.AsstObjectID)).Eq(asst.AsstInstName)
			asstCond.Field(common.BKObjIDField).Eq(asst.AsstObjectID)
		}
		if asstObjType == common.BKInnerObjIDHost {
			asstCond.Field(common.BKCloudIDField).Eq(report.CloudID)
			asstCond.Field(common.BKHostInnerIPField).Eq(asst.AsstInstName)
		}
		asstInsts, err := lgc.findInst(header, asst.AsstObjectID, &metadata.QueryInput{Condition: asstCond.ToMapStr()})
		if err != nil {
			blog.Errorf("[NetDevice][ConfirmReport] find inst by %+v failed %v", asstCond.ToMapStr(), err)
			errs = append(errs, err)
			continue
		}
		blog.V(4).Infof("[NetDevice][ConfirmReport] find inst result: %#v, condition: %#v", asstInsts, asstCond.ToMapStr())
		if len(asstInsts) > 0 {
			asstInstID, err := asstInsts[0].Int64(common.GetInstIDField(asst.AsstObjectID))
			if err != nil {
				blog.Errorf("[NetDevice][ConfirmReport] propertyID %s not exist in %#v ", common.GetInstIDField(asst.AsstObjectID), asstInsts[0])
				errs = append(errs, err)
				continue
			}

			updateBody := map[string]interface{}{
				asst.AsstPropertyID: asstInstID,
			}

			resp, err := lgc.CoreAPI.TopoServer().Instance().UpdateInst(context.Background(), util.GetUser(header), report.ObjectID, instID, header, updateBody)
			if err != nil {
				blog.Errorf("[NetDevice][ConfirmReport] update inst error: %v, %+v", err, updateBody)
				errs = append(errs, err)
				continue
			}
			if !resp.Result {
				blog.Errorf("[NetDevice][ConfirmReport] update inst error: %v, %+v", resp.ErrMsg, updateBody)
				errs = append(errs, fmt.Errorf(resp.ErrMsg))
				continue
			}
			successCount++
		}
	}
	return successCount, errs
}

func (lgc *Logics) saveHistory(report *metadata.NetcollectReport, success bool) error {
	history := metadata.NetcollectHistory{NetcollectReport: *report, Success: success}
	err := lgc.Instance.Table(common.BKTableNameNetcollectHistory).Insert(lgc.ctx, history)
	if err != nil {
		blog.Errorf("[NetDevice][ConfirmReport] save history %+v failed: %v", history, err)
	}
	return err
}

func (lgc *Logics) SearchHistory(header http.Header, param metadata.ParamSearchNetcollectReport) (uint64, []metadata.NetcollectHistory, error) {
	historys := []metadata.NetcollectHistory{}
	cond, err := lgc.buildSearchCond(header, param)
	if err != nil {
		blog.Errorf("[NetDevice][SearchHistory] build SearchHistory condition for %+v failed: %v", param, err)
		return 0, nil, err
	}
	count, err := lgc.Instance.Table(common.BKTableNameNetcollectHistory).Find(cond.ToMapStr()).Count(lgc.ctx)
	if err != nil {
		blog.Errorf("[NetDevice][SearchHistory] GetCntByCondition for %+v failed: %v", cond.ToMapStr(), err)
		return 0, nil, err
	}

	// search historys
	err = lgc.Instance.Table(common.BKTableNameNetcollectHistory).Find(cond.ToMapStr()).Sort(param.Page.Sort).Start(uint64(param.Page.Start)).Limit(uint64(param.Page.Limit)).All(lgc.ctx, &historys)
	if err != nil {
		blog.Errorf("[NetDevice][SearchHistory] GetMutilByCondition for %+v failed: %v", cond.ToMapStr(), err)
		return 0, nil, err
	}

	return count, historys, nil
}
