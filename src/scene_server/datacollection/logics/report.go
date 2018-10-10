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
	"net/http"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
)

func (lgc *Logics) SearchReportSummary(header http.Header, param metadata.ParamSearchNetcollectReport) ([]*metadata.NetcollectReportSummary, error) {
	// search reports
	param.CloudID = -1
	_, reports, err := lgc.SearchReport(header, param)
	if err != nil {
		blog.Errorf("[NetDevice][SearchReportSummary] SearchReport failed %v", err)
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

		if report.LastTime != nil && summary.LastTime != nil && report.LastTime.Time.Sub(summary.LastTime.Time) > 0 {
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

func (lgc *Logics) SearchReport(header http.Header, param metadata.ParamSearchNetcollectReport) (int64, []metadata.NetcollectReport, error) {
	cond, err := lgc.buildSearchCond(header, param)
	if err != nil {
		blog.Errorf("[NetDevice][SearchReport] build SearchReport condition failed: %v", err)
		return 0, nil, err
	}
	count, err := lgc.Instance.GetCntByCondition(common.BKTableNameNetcollectReport, cond.ToMapStr())
	if err != nil {
		blog.Errorf("[NetDevice][SearchReport] GetCntByCondition failed: %v", err)
		return 0, nil, err
	}

	// search reports
	reports := []metadata.NetcollectReport{}
	err = lgc.Instance.GetMutilByCondition(common.BKTableNameNetcollectReport, nil, cond.ToMapStr(), &reports, param.Page.Sort, param.Page.Start, param.Page.Limit)
	if err != nil {
		blog.Errorf("[NetDevice][SearchReport] GetMutilByCondition failed: %v", err)
		return 0, nil, err
	}

	// search details
	objIDs := []string{}
	cloudIDs := []int64{}
	for _, report := range reports {
		objIDs = append(objIDs, report.ObjectID)
		cloudIDs = append(cloudIDs, report.CloudID)
	}

	cloudCond := condition.CreateCondition()
	cloudCond.Field(common.BKCloudIDField).In(cloudIDs)

	cloudMap, err := lgc.findInstMap(header, common.BKInnerObjIDPlat, &metadata.QueryInput{Condition: cloudCond.ToMapStr()})
	if err != nil {
		blog.Errorf("[NetDevice][SearchReport] find clouds failed: %v", err)
		return 0, nil, err
	}

	objMap, err := lgc.findObjectMap(header, objIDs...)
	if err != nil {
		blog.Errorf("[NetDevice][SearchReport] findObjectMap failed: %v", err)
		return 0, nil, err
	}
	attrsMap, err := lgc.findAttrsMap(header, objIDs...)
	if err != nil {
		blog.Errorf("[NetDevice][SearchReport] findAttrsMap failed: %v", err)
		return 0, nil, err
	}

	for index := range reports {
		if object, ok := objMap[reports[index].ObjectID]; ok {
			reports[index].ObjectName = object.ObjectName
		}
		if clodInst, ok := cloudMap[reports[index].CloudID]; ok {
			clodname, err := clodInst.String(common.BKCloudNameField)
			if err != nil {
				blog.Errorf("[NetDevice][SearchReport] bk_cloud_name field invalied: %v", err)
			}
			reports[index].CloudName = clodname
		}

		cond := condition.CreateCondition()
		cond.Field(common.BKInstNameField).Eq(reports[index].InstKey)
		objType := common.GetObjByType(reports[index].ObjectID)
		if objType == common.BKINnerObjIDObject {
			cond.Field(common.BKObjIDField).Eq(reports[index].ObjectID)
		}
		insts, err := lgc.findInst(header, objType, &metadata.QueryInput{Condition: cond.ToMapStr()})
		if err != nil {
			blog.Errorf("[NetDevice][SearchReport] find inst failed %v", err)
			return 0, nil, err
		}
		if len(insts) > 0 {
			reports[index].Action = metadata.ReporctActionUpdate
			inst := insts[0]
			for _, attribute := range reports[index].Attributes {
				attribute.PreValue = inst[attribute.PropertyID]
				if property, ok := attrsMap[reports[index].ObjectID+":"+attribute.PropertyID]; ok {
					attribute.PropertyName = property.PropertyName
					attribute.IsRequired = property.IsRequired
				}
			}
		} else {
			reports[index].Action = metadata.ReporctActionCreate
		}

	}

	return int64(count), reports, nil
}

func (lgc *Logics) findAttrsMap(header http.Header, objIDs ...string) (map[string]metadata.Attribute, error) {
	attrs, err := lgc.findAttrs(header, objIDs...)
	if err != nil {
		return nil, err
	}

	attrsMap := map[string]metadata.Attribute{}
	for _, attr := range attrs {
		attrsMap[attr.ObjectID+":"+attr.PropertyID] = attr
	}
	return attrsMap, nil
}

func (lgc *Logics) findAttrs(header http.Header, objIDs ...string) ([]metadata.Attribute, error) {
	cond := condition.CreateCondition()
	cond.Field(common.BKObjIDField).In(objIDs)
	resp, err := lgc.CoreAPI.ObjectController().Meta().SelectObjectAttWithParams(context.Background(), header, cond.ToMapStr())
	if err != nil {
		blog.Infof("[NetDevice][findAttrs] error: %v", err)
		return nil, err
	}
	if !resp.Result {
		blog.Infof("[NetDevice][findAttrs] error: %v", resp.ErrMsg)
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
		blog.Infof("[NetDevice][findObject] error: %v", err)
		return nil, err
	}
	if !resp.Result {
		blog.Infof("[NetDevice][findObject] error: %v", resp.ErrMsg)
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
		blog.Infof("[NetDevice][findInst] error: %v", err)
		return nil, err
	}
	if !resp.Result {
		blog.Infof("[NetDevice][findInst] error: %v", resp.ErrMsg)
		return nil, err
	}
	return resp.Data.Info, nil
}

func (lgc *Logics) ConfirmReport(header http.Header, reports []metadata.NetcollectReport) *metadata.RspNetcollectConfirm {
	result := metadata.RspNetcollectConfirm{}
	for index := range reports {
		report := &reports[index]
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
			continue
		}

		objType := common.GetObjByType(report.ObjectID)
		cond := condition.CreateCondition()

		cond.Field(common.GetInstNameField(report.ObjectID)).Eq(report.InstKey)
		if objType == common.BKINnerObjIDObject {
			cond.Field(common.BKObjIDField).Eq(report.ObjectID)
			data.Set(common.BKObjIDField, report.ObjectID)
		}

		insts, err := lgc.findInst(header, report.ObjectID, &metadata.QueryInput{Condition: cond.ToMapStr()})
		if err != nil {
			blog.Errorf("[NetDevice][ConfirmReport] find inst failed %v", err)
			lgc.saveHistory(report, false)
			continue
		}
		blog.V(4).Infof("[NetDevice][ConfirmReport] find inst result: %#v, condition: %#v", insts, cond.ToMapStr())
		if len(insts) > 0 {
			updateBody := map[string]interface{}{
				"data":      data,
				"condition": cond.ToMapStr(),
			}
			resp, err := lgc.CoreAPI.ObjectController().Instance().UpdateObject(context.Background(), report.ObjectID, header, updateBody)
			if err != nil {
				blog.Errorf("[NetDevice][ConfirmReport] update inst error: %v, %+v", err, updateBody)
				result.ChangeAttributeFailure += attrCount
				result.Errors = append(result.Errors, err.Error())
				lgc.saveHistory(report, false)
				continue
			}
			if !resp.Result {
				blog.Errorf("[NetDevice][ConfirmReport] update inst error: %v, %+v", resp.ErrMsg, updateBody)
				result.ChangeAttributeFailure += attrCount
				result.Errors = append(result.Errors, resp.ErrMsg)
				lgc.saveHistory(report, false)
				continue
			}
			result.ChangeAttributeSuccess += attrCount
			lgc.saveHistory(report, true)
		} else {
			resp, err := lgc.CoreAPI.ObjectController().Instance().CreateObject(context.Background(), report.ObjectID, header, data)
			if err != nil {
				blog.Errorf("[NetDevice][ConfirmReport] create inst error: %v, %+v", err, data)
				result.ChangeAttributeFailure += attrCount
				result.Errors = append(result.Errors, err.Error())
				lgc.saveHistory(report, false)
				continue
			}
			if !resp.Result {
				blog.Errorf("[NetDevice][ConfirmReport] create inst error: %v, %+v", resp.ErrMsg, data)
				result.ChangeAttributeFailure += attrCount
				result.Errors = append(result.Errors, resp.ErrMsg)
				lgc.saveHistory(report, false)
				continue
			}
			result.ChangeAttributeSuccess += attrCount
			lgc.saveHistory(report, true)
		}
	}
	return &result
}

func (lgc *Logics) saveHistory(report *metadata.NetcollectReport, success bool) error {
	history := metadata.NetcollectHistory{NetcollectReport: *report, Success: success}
	_, err := lgc.Instance.Insert(common.BKTableNameNetcollectHistory, history)
	if err != nil {
		blog.Errorf("[NetDevice][ConfirmReport] save history failed: %v", err)
	}
	return err
}

func (lgc *Logics) SearchHistory(header http.Header, param metadata.ParamSearchNetcollectReport) (int64, []metadata.NetcollectHistory, error) {
	historys := []metadata.NetcollectHistory{}
	cond, err := lgc.buildSearchCond(header, param)
	if err != nil {
		blog.Errorf("[NetDevice][SearchHistory] build SearchHistory condition failed: %v", err)
		return 0, nil, err
	}
	count, err := lgc.Instance.GetCntByCondition(common.BKTableNameNetcollectHistory, cond.ToMapStr())
	if err != nil {
		blog.Errorf("[NetDevice][SearchHistory] GetCntByCondition failed: %v", err)
		return 0, nil, err
	}

	// search historys
	err = lgc.Instance.GetMutilByCondition(common.BKTableNameNetcollectHistory, nil, cond.ToMapStr(), &historys, param.Page.Sort, param.Page.Start, param.Page.Limit)
	if err != nil {
		blog.Errorf("[NetDevice][SearchHistory] GetMutilByCondition failed: %v", err)
		return 0, nil, err
	}

	return int64(count), historys, nil
}
