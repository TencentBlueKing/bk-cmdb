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

	"github.com/rs/xid"
)

func (lgc *Logics) SearchReportSummary(header http.Header, param metadata.ParamSearchNetcollectReport) ([]*metadata.NetcollectReportSummary, error) {
	rid := util.GetHTTPCCRequestID(header)
	// search reports
	param.CloudID = -1
	_, reports, err := lgc.SearchReport(header, param)
	if err != nil {
		blog.Errorf("[NetDevice][SearchReportSummary] SearchReport by %+v failed %v, rid: %s", param, err, rid)
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

	summaries := make([]*metadata.NetcollectReportSummary, 0)
	for key := range summarym {
		summaries = append(summaries, summarym[key])
	}

	return summaries, nil
}

func (lgc *Logics) buildSearchCond(header http.Header, param metadata.ParamSearchNetcollectReport) (condition.Condition, error) {
	cond := condition.CreateCondition()
	if param.CloudID >= 0 {
		cond.Field(common.BKCloudIDField).Eq(param.CloudID)
	}
	cloudIDs := make([]int64, 0)
	if param.CloudName != "" || param.Query != "" {
		cloudCond := condition.CreateCondition()
		cloudCond.Field(common.BKCloudNameField).Like(param.CloudName)
		clouds, err := lgc.findInst(header, common.BKInnerObjIDPlat, &metadata.QueryCondition{Condition: cloudCond.ToMapStr()})
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
	rid := util.GetHTTPCCRequestID(header)
	cond, err := lgc.buildSearchCond(header, param)
	if err != nil {
		blog.Errorf("[NetDevice][SearchReport] build SearchReport condition for %+v, failed: %v, rid: %s", param, err, rid)
		return 0, nil, err
	}
	count, err := lgc.db.Table(common.BKTableNameNetcollectReport).Find(cond.ToMapStr()).Count(lgc.ctx)
	if err != nil {
		blog.Errorf("[NetDevice][SearchReport] GetCntByCondition %+v failed: %v, rid: %s", cond.ToMapStr(), err, rid)
		return 0, nil, err
	}

	// search reports
	reports := make([]metadata.NetcollectReport, 0)
	err = lgc.db.Table(common.BKTableNameNetcollectReport).Find(cond.ToMapStr()).Sort(param.Page.Sort).Start(uint64(param.Page.Start)).Limit(uint64(param.Page.Limit)).All(lgc.ctx, &reports)
	if err != nil {
		blog.Errorf("[NetDevice][SearchReport] GetMutilByCondition %+v failed: %v, rid: %s", cond.ToMapStr(), err, rid)
		return 0, nil, err
	}

	// search details
	objIDs := make([]string, 0)
	cloudIDs := make([]int64, 0)
	for _, report := range reports {
		objIDs = append(objIDs, report.ObjectID)
		cloudIDs = append(cloudIDs, report.CloudID)

		for _, asst := range report.Associations {
			objIDs = append(objIDs, asst.AsstObjectID)
		}
	}

	cloudCond := condition.CreateCondition()
	cloudCond.Field(common.BKCloudIDField).In(cloudIDs)
	cloudMap, err := lgc.findInstMap(header, common.BKInnerObjIDPlat, &metadata.QueryCondition{Condition: cloudCond.ToMapStr()})
	if err != nil {
		blog.Errorf("[NetDevice][SearchReport] find clouds by %+v failed, err: %v, rid: %s", cloudCond, err, rid)
		return 0, nil, err
	}

	objMap, err := lgc.findObjectMap(header, objIDs...)
	if err != nil {
		blog.Errorf("[NetDevice][SearchReport] findObjectMap by %+v failed, err: %v, rid: %s", objIDs, err, rid)
		return 0, nil, err
	}
	attrsMap, err := lgc.findAttrsMap(header, objIDs...)
	if err != nil {
		blog.Errorf("[NetDevice][SearchReport] findAttrsMap by %+v failed, err: %v, rid: %s", objIDs, err, rid)
		return 0, nil, err
	}

	for index := range reports {
		if object, ok := objMap[reports[index].ObjectID]; ok {
			reports[index].ObjectName = object.ObjectName
		}
		if clodInst, ok := cloudMap[reports[index].CloudID]; ok {
			cloudname, err := clodInst.String(common.BKCloudNameField)
			if err != nil {
				blog.Errorf("[NetDevice][SearchReport] bk_cloud_name field invalid, err: %v, rid: %s", err, rid)
			}
			reports[index].CloudName = cloudname
		}

		cond := condition.CreateCondition()
		cond.Field(common.BKInstNameField).Eq(reports[index].InstKey)
		objType := common.GetObjByType(reports[index].ObjectID)
		if objType == common.BKInnerObjIDObject {
			cond.Field(common.BKObjIDField).Eq(reports[index].ObjectID)
		}
		insts, err := lgc.findInst(header, reports[index].ObjectID, &metadata.QueryCondition{Condition: cond.ToMapStr()})
		if err != nil {
			blog.Errorf("[NetDevice][SearchReport] find inst by %+v for %v failed, err: %v, rid: %s", cond.ToMapStr(), reports[index].ObjectID, err, rid)
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
	rid := util.GetHTTPCCRequestID(header)
	cond := condition.CreateCondition()
	cond.Field(common.BKObjIDField).In(objIDs)
	resp, err := lgc.CoreAPI.CoreService().Model().ReadModelAttrByCondition(context.Background(), header, &metadata.QueryCondition{Condition: cond.ToMapStr()})
	if err != nil {
		blog.Errorf("[NetDevice][findAttrs] for %v failed, err: %v, rid: %s", objIDs, err, rid)
		return nil, err
	}
	if !resp.Result {
		blog.Errorf("[NetDevice][findAttrs] for %v failed, err: %v, rid: %s", objIDs, resp, rid)
		return nil, err
	}
	return resp.Data.Info, nil
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

func (lgc *Logics) findObject(header http.Header, filter mapstr.MapStr) ([]metadata.Object, error) {
	rid := util.GetHTTPCCRequestID(header)
	cond := &metadata.QueryCondition{Condition: filter}
	resp, err := lgc.CoreAPI.CoreService().Model().ReadModel(context.Background(), header, cond)
	if err != nil {
		blog.Errorf("[NetDevice][findObject] by %+v failed, err: %v, rid: %s", filter, err, rid)
		return nil, err
	}
	if !resp.Result {
		blog.Errorf("[NetDevice][findObject] by %+v failed, err: %v, rid: %s", filter, resp, rid)
		return nil, err
	}
	models := make([]metadata.Object, 0)
	for _, info := range resp.Data.Info {
		models = append(models, info.Spec)
	}

	return models, nil
}

func (lgc *Logics) findInstMap(header http.Header, objectID string, query *metadata.QueryCondition) (map[int64]mapstr.MapStr, error) {
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

func (lgc *Logics) findInst(header http.Header, objectID string, query *metadata.QueryCondition) ([]mapstr.MapStr, error) {
	rid := util.GetHTTPCCRequestID(header)
	resp, err := lgc.CoreAPI.CoreService().Instance().ReadInstance(context.Background(), header, objectID, query)
	if err != nil {
		blog.Errorf("[NetDevice][findInst] by %+v failed, err: %v, rid: %s", query, err, rid)
		return nil, err
	}
	if !resp.Result {
		blog.Errorf("[NetDevice][findInst] by %+v failed, err: %+v, rid: %s", query, resp, rid)
		return nil, err
	}
	return resp.Data.Info, nil
}

func (lgc *Logics) findInstAssociation(header http.Header, objectID string, instID int64) ([]*metadata.InstAsst, error) {
	rid := util.GetHTTPCCRequestID(header)
	cond := condition.CreateCondition()
	or := cond.NewOR()
	or.Item(mapstr.MapStr{
		"bk_inst_id": instID,
		"bk_obj_id":  objectID,
	})
	or.Item(mapstr.MapStr{
		"bk_asst_inst_id": instID,
		"bk_asst_obj_id":  objectID,
	})

	option := &metadata.SearchAssociationInstRequest{
		Condition: cond.ToMapStr(),
	}
	resp, err := lgc.CoreAPI.ApiServer().SearchAssociationInst(context.Background(), header, option)
	if err != nil {
		blog.Errorf("[NetDevice][findInstAssociation] by %+v failed, err: %v, rid: %s", cond.ToMapStr(), err, rid)
		return nil, err
	}
	if !resp.Result {
		blog.Errorf("[NetDevice][findInstAssociation] by %+v failed, err: %+v, rid: %s", cond.ToMapStr(), resp, rid)
		return nil, err
	}
	return resp.Data, nil
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
		cond := condition.CreateCondition()
		cond.Field(common.BKObjIDField).Eq(report.ObjectID)
		cond.Field(common.BKInstKeyField).Eq(report.InstKey)

		if err := lgc.db.Table(common.BKTableNameNetcollectReport).Delete(context.Background(), cond.ToMapStr()); err != nil {
			result.Errors = append(result.Errors, err.Error())
		}

	}
	return &result
}

func (lgc *Logics) confirmAttributes(header http.Header, report *metadata.NetcollectReport) (int, error) {
	rid := util.GetHTTPCCRequestID(header)
	data := mapstr.MapStr{}
	attrCount := 0
	for _, attr := range report.Attributes {
		if attr.Method == metadata.ReporctMethodAccept {
			data.Set(attr.PropertyID, attr.CurValue)
			attrCount++
		}
	}

	if len(data) <= 0 {
		blog.Warnf("[NetDevice][ConfirmReport] empty data, continue next, rid: %s", rid)
		return attrCount, nil
	}

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
	insts, err := lgc.findInst(header, report.ObjectID, &metadata.QueryCondition{Condition: cond.ToMapStr()})
	if err != nil {
		blog.Errorf("[NetDevice][ConfirmReport] find inst failed %v, rid: %s", err, rid)
		return attrCount, err
	}
	blog.V(4).Infof("[NetDevice][ConfirmReport] find inst result: %#v, condition: %#v, rid: %s", insts, cond.ToMapStr(), rid)

	if objType == common.BKInnerObjIDHost {
		var assetID string
		if len(insts) > 0 {
			assetID, err = insts[0].String(common.BKAssetIDField)
			if err != nil {
				blog.Warnf("[NetDevice][ConfirmReport] find inst failed, bk_asset_id not found from %+v, rid: %s", insts[0], rid)
			}
		}
		if assetID == "" {
			data.Set(common.BKAssetIDField, xid.New().String())
		}

		hostdata := metadata.HostList{
			InputType: metadata.CollectType,
			HostInfo: map[int64]map[string]interface{}{
				1: data,
			},
		}

		resp, err := lgc.CoreAPI.HostServer().AddHost(context.Background(), header, &hostdata)
		if err != nil {
			blog.Errorf("[NetDevice][ConfirmReport] add host error: %v, %+v, rid: %s", err, data, rid)
			return attrCount, err
		}
		if !resp.Result {
			blog.Errorf("[NetDevice][ConfirmReport] add host error: %v, %+v, rid: %s", resp, data, rid)
			return attrCount, fmt.Errorf(resp.ErrMsg)
		}
		return attrCount, nil
	}

	if len(insts) > 0 {
		instID, err := insts[0].Int64(common.GetInstIDField(report.ObjectID))
		if err != nil {
			blog.Errorf("[NetDevice][ConfirmReport] find inst failed, instID not found from %+v, rid: %s", insts[0], rid)
			return 0, err
		}

		resp, err := lgc.CoreAPI.TopoServer().Instance().UpdateInst(context.Background(), report.ObjectID, instID, header, data)
		if err != nil {
			blog.Errorf("[NetDevice][ConfirmReport] update inst error: %v, %+v, rid: %s", err, data, rid)
			return attrCount, err
		}
		if !resp.Result {
			blog.Errorf("[NetDevice][ConfirmReport] update inst error: %v, %+v, rid: %s", resp, data, rid)
			return attrCount, fmt.Errorf(resp.ErrMsg)
		}
	} else {
		resp, err := lgc.CoreAPI.TopoServer().Instance().CreateInst(context.Background(), report.ObjectID, header, data)
		if err != nil {
			blog.Errorf("[NetDevice][ConfirmReport] create inst to %s error: %v, %+v, rid: %s", report.ObjectID, err, data, rid)
			return attrCount, err
		}
		if !resp.Result {
			blog.Errorf("[NetDevice][ConfirmReport] create inst to %s error: %v, %+v, rid: %s", report.ObjectID, resp, data, rid)
			return attrCount, fmt.Errorf(resp.ErrMsg)
		}
	}
	return attrCount, nil
}

func (lgc *Logics) confirmAssociations(header http.Header, report *metadata.NetcollectReport) (successCount int, errs []error) {
	rid := util.GetHTTPCCRequestID(header)
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

	insts, err := lgc.findInst(header, report.ObjectID, &metadata.QueryCondition{Condition: cond.ToMapStr()})
	if err != nil {
		blog.Errorf("[NetDevice][ConfirmReport] find inst %+v failed %v, rid: %s", cond.ToMapStr(), err, rid)
		return 0, append(errs, err)
	}
	if len(insts) <= 0 {
		blog.Errorf("[NetDevice][ConfirmReport] find inst failed, inst not found by %+v, rid: %s", cond.ToMapStr(), rid)
		return 0, append(errs, fmt.Errorf("inst not found"))
	}
	instID, err := insts[0].Int64(common.GetInstIDField(report.ObjectID))
	if err != nil {
		blog.Errorf("[NetDevice][ConfirmReport] find inst failed, instID not found from %+v, rid: %s", insts[0], rid)
		return 0, append(errs, fmt.Errorf("inst not found"))
	}

	instassts, err := lgc.findInstAssociation(header, report.ObjectID, instID)
	if err != nil {
		blog.Errorf("[NetDevice][ConfirmReport] find inst association failed, association not found from %+v, rid: %s", insts[0], rid)
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
		asstInsts, err := lgc.findInst(header, asst.AsstObjectID, &metadata.QueryCondition{Condition: asstCond.ToMapStr()})
		if err != nil {
			blog.Errorf("[NetDevice][ConfirmReport] find inst by %+v failed %v, rid: %s", asstCond.ToMapStr(), err, rid)
			errs = append(errs, err)
			continue
		}
		blog.V(4).Infof("[NetDevice][ConfirmReport] find inst result: %#v, condition: %#v, rid: %s", asstInsts, asstCond.ToMapStr(), rid)
		if len(asstInsts) > 0 {
			asstInstID, err := asstInsts[0].Int64(common.GetInstIDField(asst.AsstObjectID))
			if err != nil {
				blog.Errorf("[NetDevice][ConfirmReport] propertyID %s not exist in %#v, rid: %s", common.GetInstIDField(asst.AsstObjectID), asstInsts[0], rid)
				errs = append(errs, err)
				continue
			}

			if !isAssociationExists(instassts, report.ObjectID, instID, asst.AsstObjectID, asstInstID) {
				req := metadata.CreateAssociationInstRequest{
					ObjectAsstID: asst.ObjectAsstID,
					InstID:       instID,
					AsstInstID:   asstInstID,
				}
				resp, err := lgc.CoreAPI.TopoServer().Association().CreateInst(context.Background(), header, &req)
				if err != nil {
					blog.Errorf("[NetDevice][ConfirmReport] create inst association error: %v, %+v, rid: %s", err, req, rid)
					errs = append(errs, err)
					continue
				}
				if !resp.Result {
					blog.Errorf("[NetDevice][ConfirmReport] create inst association error: %v, %+v, rid: %s", resp.ErrMsg, req, rid)
					errs = append(errs, fmt.Errorf(resp.ErrMsg))
					continue
				}
			}
			successCount++
		}
	}
	return successCount, errs
}

func isAssociationExists(assts []*metadata.InstAsst, objectID string, instID int64, asstObjectID string, asstInstID int64) bool {
	for _, asst := range assts {
		if asst.ObjectID == objectID &&
			asst.InstID == instID &&
			asst.AsstObjectID == asstObjectID &&
			asst.AsstInstID == asstInstID {
			return true
		}
	}
	return false
}

func (lgc *Logics) saveHistory(report *metadata.NetcollectReport, success bool) error {
	rid := util.ExtractRequestIDFromContext(lgc.ctx)
	history := metadata.NetcollectHistory{NetcollectReport: *report, Success: success}
	err := lgc.db.Table(common.BKTableNameNetcollectHistory).Insert(lgc.ctx, history)
	if err != nil {
		blog.Errorf("[NetDevice][ConfirmReport] save history %+v failed: %v, rid: %s", history, err, rid)
	}
	return err
}

func (lgc *Logics) SearchHistory(header http.Header, param metadata.ParamSearchNetcollectReport) (uint64, []metadata.NetcollectHistory, error) {
	rid := util.GetHTTPCCRequestID(header)
	reports := make([]metadata.NetcollectHistory, 0)
	cond, err := lgc.buildSearchCond(header, param)
	if err != nil {
		blog.Errorf("[NetDevice][SearchHistory] build SearchHistory condition for %+v failed: %v, rid: %s", param, err, rid)
		return 0, nil, err
	}
	count, err := lgc.db.Table(common.BKTableNameNetcollectHistory).Find(cond.ToMapStr()).Count(lgc.ctx)
	if err != nil {
		blog.Errorf("[NetDevice][SearchHistory] GetCntByCondition for %+v failed: %v, rid: %s", cond.ToMapStr(), err, rid)
		return 0, nil, err
	}

	// search reports
	err = lgc.db.Table(common.BKTableNameNetcollectHistory).Find(cond.ToMapStr()).Sort(param.Page.Sort).Start(uint64(param.Page.Start)).Limit(uint64(param.Page.Limit)).All(lgc.ctx, &reports)
	if err != nil {
		blog.Errorf("[NetDevice][SearchHistory] GetMutilByCondition for %+v failed: %v, rid: %s", cond.ToMapStr(), err, rid)
		return 0, nil, err
	}

	// search details
	objIDs := make([]string, 0)
	cloudIDs := make([]int64, 0)
	for _, report := range reports {
		objIDs = append(objIDs, report.ObjectID)
		cloudIDs = append(cloudIDs, report.CloudID)

		for _, asst := range report.Associations {
			objIDs = append(objIDs, asst.AsstObjectID)
		}
	}

	cloudCond := condition.CreateCondition()
	cloudCond.Field(common.BKCloudIDField).In(cloudIDs)
	cloudMap, err := lgc.findInstMap(header, common.BKInnerObjIDPlat, &metadata.QueryCondition{Condition: cloudCond.ToMapStr()})
	if err != nil {
		blog.Errorf("[NetDevice][SearchReport] find clouds by %+v failed: %v, rid: %s", cloudCond, err, rid)
		return 0, nil, err
	}

	objMap, err := lgc.findObjectMap(header, objIDs...)
	if err != nil {
		blog.Errorf("[NetDevice][SearchReport] findObjectMap by %+v failed: %v, rid: %s", objIDs, err, rid)
		return 0, nil, err
	}

	for index := range reports {
		if object, ok := objMap[reports[index].ObjectID]; ok {
			reports[index].ObjectName = object.ObjectName
		}
		if clodInst, ok := cloudMap[reports[index].CloudID]; ok {
			cloudname, err := clodInst.String(common.BKCloudNameField)
			if err != nil {
				blog.Errorf("[NetDevice][SearchReport] bk_cloud_name field invalid: %v, rid: %s", err, rid)
			}
			reports[index].CloudName = cloudname
		}
	}

	return count, reports, nil
}
