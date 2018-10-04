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

func (lgc *Logics) SearchReport(header http.Header, param metadata.ParamSearchNetcollectReport) (int64, []metadata.NetcollectReport, error) {
	cond := condition.CreateCondition()
	if param.CloudID >= 0 {
		cond.Field(common.BKCloudIDField).Eq(param.CloudID)
	}
	if param.CloudName != "" {
		cloudCond := condition.CreateCondition()
		cloudCond.Field(common.BKCloudNameField).Like(param.CloudName)
		lgc.findObject(header, cloudCond)
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
		cond.Field(common.BKObjIDField).Eq(reports[index].ObjectID)
		insts, err := lgc.findInst(header, reports[index].ObjectID, &metadata.QueryInput{Condition: cond.ToMapStr()})
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

func (lgc *Logics) searchDetail(cloudID int64, ips []string) ([]metadata.NetcollectReportAttribute, []metadata.NetcollectReportAssociation, error) {
	detailCond := condition.CreateCondition()
	detailCond.Field(common.BKCloudIDField).Eq(cloudID)
	detailCond.Field(common.BKHostInnerIPField).In(ips)

	attributes := []metadata.NetcollectReportAttribute{}
	associations := []metadata.NetcollectReportAssociation{}
	err := lgc.Instance.GetMutilByCondition(common.BKTableNameNetcollectAssociation, nil, detailCond, &attributes, "", 0, 0)
	if err != nil {
		blog.Errorf("[NetDevice][SearchReport] GetMutilByCondition failed %v", err)
		return nil, nil, err
	}
	err = lgc.Instance.GetMutilByCondition(common.BKTableNameNetcollectAssociation, nil, detailCond, &associations, "", 0, 0)
	if err != nil {
		blog.Errorf("[NetDevice][SearchReport] GetMutilByCondition failed %v", err)
		return nil, nil, err
	}

	return attributes, associations, nil
}

func (lgc *Logics) ConfirmReport(report metadata.RspNetcollectReport) {

}
