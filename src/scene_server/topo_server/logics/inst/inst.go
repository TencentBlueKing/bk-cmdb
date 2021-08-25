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

package inst

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"configcenter/src/ac/extensions"
	"configcenter/src/apimachinery"
	"configcenter/src/common"
	"configcenter/src/common/auditlog"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/language"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

// InstOperationInterface inst operation methods
type InstOperationInterface interface {
	// CreateInst create instance by object and create message
	CreateInst(kit *rest.Kit, objID string, data mapstr.MapStr) (*mapstr.MapStr, error)
	// CreateManyInstance batch create instance by object and create message
	CreateManyInstance(kit *rest.Kit, objID string, data []mapstr.MapStr) (*metadata.CreateManyCommInstResultDetail,
		error)
	// CreateInstBatch batch create instance by excel
	CreateInstBatch(kit *rest.Kit, objID string, batchInfo *metadata.InstBatchInfo) (*BatchResult, error)
	// DeleteInst delete instance by objectid and condition
	DeleteInst(kit *rest.Kit, objectID string, cond mapstr.MapStr, needCheckHost bool) error
	// DeleteMainlineInstWithID delete mainline instance by it's bk_inst_id
	DeleteMainlineInstWithID(kit *rest.Kit, obj metadata.Object, instID int64) error
	// DeleteInstByInstID batch delete instance by inst id
	DeleteInstByInstID(kit *rest.Kit, objectID string, instID []int64, needCheckHost bool) error
	// FindInst search instance by condition
	FindInst(kit *rest.Kit, objID string, cond *metadata.QueryInput) (*metadata.InstResult, error)
	// UpdateInst update instance by condition
	UpdateInst(kit *rest.Kit, cond, data mapstr.MapStr, objID string) error
	// SearchObjectInstances searches object instances.
	SearchObjectInstances(kit *rest.Kit, objID string, input *metadata.CommonSearchFilter) (
		*metadata.CommonSearchResult, error)
	// CountObjectInstances counts object instances num.
	CountObjectInstances(kit *rest.Kit, objID string, input *metadata.CommonCountFilter) (*metadata.CommonCountResult,
		error)
	// FindInstChildTopo find instance's child topo
	FindInstChildTopo(kit *rest.Kit, obj metadata.Object, instID int64, query *metadata.QueryInput) (int,
		[]*CommonInstTopo, error)
	// FindInstParentTopo find instance's parent topo
	FindInstParentTopo(kit *rest.Kit, obj metadata.Object, instID int64, query *metadata.QueryInput) (int,
		[]*CommonInstTopo, error)
	// FindInstTopo find instance all topo which include it's child and parent
	FindInstTopo(kit *rest.Kit, obj metadata.Object, instID int64, query *metadata.QueryInput) (int,
		[]CommonInstanceTopo, error)
	// IsMainline check if object is mainline object
	IsMainline(kit *rest.Kit, objID string) (bool, error)
}

// NewInstOperation create a new inst operation instance
func NewInstOperation(client apimachinery.ClientSetInterface,
	authManager *extensions.AuthManager) InstOperationInterface {

	return &commonInst{
		clientSet:   client,
		authManager: authManager,
	}
}

// BatchResult batch create instance by excel result
type BatchResult struct {
	Errors         []string `json:"error"`
	Success        []string `json:"success"`
	SuccessCreated []int64  `json:"success_created"`
	SuccessUpdated []int64  `json:"success_updated"`
	UpdateErrors   []string `json:"update_error"`
}

// CommonInstTopo common inst topo
type CommonInstTopo struct {
	metadata.InstNameAsst
	Count    int                     `json:"count"`
	Children []metadata.InstNameAsst `json:"children"`
}

// CommonInstanceTopo set of CommonInstTopo
type CommonInstanceTopo struct {
	Prev []*CommonInstTopo `json:"prev"`
	Next []*CommonInstTopo `json:"next"`
	Curr interface{}       `json:"curr"`
}

// ObjectWithInsts a struct include object msg and insts array
type ObjectWithInsts struct {
	Object metadata.Object
	Insts  []mapstr.MapStr
}

// ObjectAssoPair a struct include object msg and association
type ObjectAssoPair struct {
	Object      metadata.Object
	Association metadata.Association
}

type commonInst struct {
	clientSet   apimachinery.ClientSetInterface
	language    language.CCLanguageIf
	authManager *extensions.AuthManager
}

// CreateInst create instance by object and create message
func (c *commonInst) CreateInst(kit *rest.Kit, objID string, data mapstr.MapStr) (*mapstr.MapStr, error) {

	if err := c.validObject(kit, objID, data); err != nil {
		blog.Errorf("check object (%s) if is mainline object failed, err: %v, rid: %s", objID, err, kit.Rid)
		return nil, err
	}

	data.Set(common.BKObjIDField, objID)
	data.Set(common.BkSupplierAccount, kit.SupplierAccount)

	instCond := &metadata.CreateModelInstance{Data: data}
	rsp, err := c.clientSet.CoreService().Instance().CreateInstance(kit.Ctx, kit.Header, objID, instCond)
	if err != nil {
		blog.Errorf("failed to create object instance, err: %v, rid: %s", err, kit.Rid)
		return nil, err
	}

	if rsp.Created.ID == 0 {
		blog.Errorf("failed to create object instance, return nothing, rid: %s", kit.Rid)
		return nil, kit.CCError.Error(common.CCErrTopoInstCreateFailed)
	}

	input := &metadata.QueryInput{Condition: mapstr.MapStr{metadata.GetInstIDFieldByObjID(objID): rsp.Created.ID}}
	inst, err := c.FindInst(kit, objID, input)
	if err != nil {
		blog.Errorf("search instance by inst_id(%s) failed, err: %v, rid: %s", rsp.Created.ID, err, kit.Rid)
		return nil, err
	}

	if len(inst.Info) != 1 {
		blog.Errorf("search instance by inst_id(%s) failed, get %d instance, rid: %s", rsp.Created.ID,
			len(inst.Info), kit.Rid)
		return nil, fmt.Errorf("get %d instance", len(inst.Info))
	}

	// for audit log.
	generateAuditParameter := auditlog.NewGenerateAuditCommonParameter(kit, metadata.AuditCreate)
	audit := auditlog.NewInstanceAudit(c.clientSet.CoreService())
	auditLog, err := audit.GenerateAuditLog(generateAuditParameter, objID, []mapstr.MapStr{inst.Info[0]})
	if err != nil {
		blog.Errorf(" creat inst, generate audit log failed, err: %v, rid: %s", err, kit.Rid)
		return nil, err
	}

	err = audit.SaveAuditLog(kit, auditLog...)
	if err != nil {
		blog.Errorf("create inst, save audit log failed, err: %v, rid: %s", err, kit.Rid)
		return nil, kit.CCError.Error(common.CCErrAuditSaveLogFailed)
	}

	return &inst.Info[0], nil
}

// CreateManyInstance batch create instance by object and create message
func (c *commonInst) CreateManyInstance(kit *rest.Kit, objID string, data []mapstr.MapStr) (
	*metadata.CreateManyCommInstResultDetail, error) {

	if len(data) == 0 {
		blog.Errorf("details cannot be empty, rid: %s", kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommInstDataNil, "details")
	}
	if len(data) > 200 {
		blog.Errorf("details cannot more than 200, details number: %s, rid: %s", len(data), kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommXXExceedLimit, "details", 200)
	}

	params := &metadata.CreateManyModelInstance{Datas: data}
	res, err := c.clientSet.CoreService().Instance().CreateManyInstance(kit.Ctx, kit.Header, objID, params)
	if err != nil {
		blog.Errorf("failed to save the object(%s) instances, err: %v, rid: %s", objID, err, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommHTTPDoRequestFailed, err.Error())
	}

	successIDs := make([]int64, 0)
	resp := metadata.NewManyCommInstResultDetail()
	for _, item := range res.Created {
		resp.SuccessCreated[item.OriginIndex] = int64(item.ID)
		successIDs = append(successIDs, int64(item.ID))
	}

	for _, item := range res.Repeated {
		errMsg, err := item.Data.String("err_msg")
		if err != nil {
			blog.Errorf("get result repeated data failed, err: %s, rid: %s", err.Error(), kit.Rid)
			return nil, err
		}
		resp.Error[item.OriginIndex] = errMsg
	}

	for _, item := range res.Exceptions {
		resp.Error[item.OriginIndex] = item.Message
	}

	if len(successIDs) == 0 {
		return resp, nil
	}

	// generate audit log of instance.
	cond := map[string]interface{}{
		common.BKInstIDField: map[string]interface{}{
			common.BKDBIN: successIDs,
		},
	}
	audit := auditlog.NewInstanceAudit(c.clientSet.CoreService())
	generateAuditParameter := auditlog.NewGenerateAuditCommonParameter(kit, metadata.AuditCreate)
	auditLog, err := audit.GenerateAuditLogByCondGetData(generateAuditParameter, objID, cond)
	if err != nil {
		blog.Errorf("create many instances, generate audit log failed, err: %v, rid: %s",
			err, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrAuditGenerateLogFailed, err.Error())
	}

	// save audit log.
	err = audit.SaveAuditLog(kit, auditLog...)
	if err != nil {
		blog.Errorf("creat many instances, save audit log failed, err: %v, rid: %s", err, kit.Rid)
		return nil, kit.CCError.Error(common.CCErrAuditSaveLogFailed)
	}
	return resp, nil
}

// CreateInstBatch batch create instance by excel
func (c *commonInst) CreateInstBatch(kit *rest.Kit, objID string, batchInfo *metadata.InstBatchInfo) (
	*BatchResult, error) {

	isMainline, err := c.IsMainline(kit, objID)
	if err != nil {
		blog.Errorf("check if object(%s) is mainline object failed, err: %v, rid: %s", objID, err, kit.Rid)
		return nil, err
	}

	if isMainline {
		return nil, kit.CCError.CCErrorf(common.CCErrCommForbiddenOperateMainlineInstanceWithCommonAPI)
	}

	if batchInfo.InputType != common.InputTypeExcel {
		return &BatchResult{}, fmt.Errorf("unexpected input_type: %s", batchInfo.InputType)
	}
	if len(batchInfo.BatchInfo) == 0 {
		return &BatchResult{}, fmt.Errorf("BatchInfo empty")
	}

	// 1. 检查实例与URL参数指定的模型一致
	for line, inst := range batchInfo.BatchInfo {
		objID, exist := inst[common.BKObjIDField]
		if exist == true && objID != objID {
			blog.Errorf("create object[%s] instance batch failed, bk_obj_id field conflict with url field,"+
				"rid: %s", objID, kit.Rid)
			return nil, kit.CCError.Errorf(common.CCErrorTopoObjectInstanceObjIDFieldConflictWithURL, line)
		}
	}

	updatedInstanceIDs := make([]int64, 0)
	createdInstanceIDs := make([]int64, 0)
	colIdxErrMap := map[int]string{}
	colIdxList := make([]int, 0)
	idFieldname := metadata.GetInstIDFieldByObjID(objID)
	results := &BatchResult{}
	for colIdx, colInput := range batchInfo.BatchInfo {
		if colInput == nil {
			// ignore empty excel line
			continue
		}

		delete(colInput, "import_from")

		// 实例id 为空，表示要新建实例
		// 实例ID已经赋值，更新数据.  (已经赋值, value not equal 0 or nil)

		// 是否存在实例ID字段
		if _, exist := colInput[idFieldname]; !exist {
			errStr := c.language.CreateDefaultCCLanguageIf(util.GetLanguage(kit.Header)).Languagef(
				"import_row_int_error_str", colIdx, fmt.Errorf("%s does not exist", idFieldname))
			colIdxList = append(colIdxList, int(colIdx))
			colIdxErrMap[int(colIdx)] = errStr
			continue
		}

		instID, err := util.GetInt64ByInterface(colInput[idFieldname])
		if err != nil {
			errStr := c.language.CreateDefaultCCLanguageIf(util.GetLanguage(kit.Header)).Languagef(
				"import_row_int_error_str", colIdx, err.Error())
			colIdxList = append(colIdxList, int(colIdx))
			colIdxErrMap[int(colIdx)] = errStr
			continue
		}

		// 实例ID字段是否设置值
		if instID != 0 {
			filter := mapstr.MapStr{idFieldname: instID}
			// remove unchangeable fields.
			delete(colInput, idFieldname)

			// to update.
			if err = c.UpdateInst(kit, filter, colInput, objID); err != nil {
				blog.Errorf("failed to update the object(%s) inst data (%#v), err: %v, rid: %s", objID, colInput,
					err, kit.Rid)
				errStr := c.language.CreateDefaultCCLanguageIf(util.GetLanguage(kit.Header)).Languagef(
					"import_row_int_error_str", colIdx, err.Error())
				colIdxList = append(colIdxList, int(colIdx))
				colIdxErrMap[int(colIdx)] = errStr
				continue
			}

			updatedInstanceIDs = append(updatedInstanceIDs, instID)
			results.Success = append(results.Success, strconv.FormatInt(colIdx, 10))
			continue
		}

		// set data
		// call CoreService.CreateInstance
		instCond := &metadata.CreateModelInstance{Data: colInput}
		rsp, err := c.clientSet.CoreService().Instance().CreateInstance(kit.Ctx, kit.Header, objID, instCond)
		if err != nil {
			blog.Errorf("failed to create object instance, err: %v, rid: %s", err, kit.Rid)
			errStr := c.language.CreateDefaultCCLanguageIf(util.GetLanguage(kit.Header)).Languagef(
				"import_row_int_error_str", colIdx, err.Error())
			colIdxList = append(colIdxList, int(colIdx))
			colIdxErrMap[int(colIdx)] = errStr
			continue
		}

		results.Success = append(results.Success, strconv.FormatInt(colIdx, 10))

		if rsp.Created.ID == 0 {
			blog.Errorf("instances created success, but get id failed, err: %+v, rid: %s", err, kit.Rid)
			continue
		}

		createdInstanceIDs = append(createdInstanceIDs, int64(rsp.Created.ID))
	}

	// generate audit log of instance.
	cond := map[string]interface{}{
		idFieldname: map[string]interface{}{common.BKDBIN: createdInstanceIDs},
	}
	audit := auditlog.NewInstanceAudit(c.clientSet.CoreService())
	generateAuditParameter := auditlog.NewGenerateAuditCommonParameter(kit, metadata.AuditCreate)
	auditLog, err := audit.GenerateAuditLogByCondGetData(generateAuditParameter, objID, cond)
	if err != nil {
		blog.Errorf(" creat inst, generate audit log failed, err: %v, rid: %s", err, kit.Rid)
		return nil, err
	}

	// save audit log.
	err = audit.SaveAuditLog(kit, auditLog...)
	if err != nil {
		blog.Errorf("creat inst, save audit log failed, err: %v, rid: %s", err, kit.Rid)
		return nil, kit.CCError.Error(common.CCErrAuditSaveLogFailed)
	}

	results.SuccessCreated = createdInstanceIDs
	results.SuccessUpdated = updatedInstanceIDs
	sort.Strings(results.Success)

	//sort error
	sort.Ints(colIdxList)
	for colIdx := range colIdxList {
		results.Errors = append(results.Errors, colIdxErrMap[colIdxList[colIdx]])
	}

	return results, nil
}

// DeleteInst delete instance by objectid and condition
func (c *commonInst) DeleteInst(kit *rest.Kit, objectID string, cond mapstr.MapStr, needCheckHost bool) error {
	return c.deleteInstByCond(kit, objectID, cond, needCheckHost)
}

// DeleteMainlineInstWithID delete mainline instance by it's bk_inst_id
func (c *commonInst) DeleteMainlineInstWithID(kit *rest.Kit, obj metadata.Object, instID int64) error {

	// if this instance has been bind to a instance by the association, then this instance should not be deleted.
	cnt, err := c.clientSet.CoreService().Association().CountInstanceAssociations(kit.Ctx, kit.Header, obj.ObjectID,
		&metadata.Condition{
			Condition: mapstr.MapStr{common.BKDBOR: []mapstr.MapStr{
				{common.BKObjIDField: obj.ObjectID, common.BKInstIDField: instID},
				{common.BKAsstObjIDField: obj.ObjectID, common.BKAsstInstIDField: instID},
			}},
		})
	if err != nil {
		blog.Errorf("count association by object(%s) failed, err: %s, rid: %s", obj.ObjectID, err, kit.Rid)
		return err
	}

	if cnt.Count > 0 {
		return kit.CCError.CCError(common.CCErrorInstHasAsst)
	}

	// delete this instance now.
	delCond := mapstr.MapStr{obj.GetInstIDFieldName(): instID}
	if obj.IsCommon() {
		delCond.Set(common.BKObjIDField, obj.ObjectID)
	}

	// generate audit log.
	audit := auditlog.NewInstanceAudit(c.clientSet.CoreService())
	generateAuditParameter := auditlog.NewGenerateAuditCommonParameter(kit, metadata.AuditDelete)
	auditLog, err := audit.GenerateAuditLogByCondGetData(generateAuditParameter, obj.GetObjectID(), delCond)
	if err != nil {
		blog.Errorf(" delete inst, generate audit log failed, err: %v, rid: %s", err, kit.Rid)
		return err
	}

	// to delete.
	ops := metadata.DeleteOption{
		Condition: delCond,
	}
	_, err = c.clientSet.CoreService().Instance().DeleteInstance(kit.Ctx, kit.Header, obj.ObjectID, &ops)
	if err != nil {
		blog.Errorf("failed to delete the object(%s) inst by the condition(%#v), err: %v",
			obj.ObjectID, ops, err)
		return err
	}

	// save audit log.
	if err := audit.SaveAuditLog(kit, auditLog...); err != nil {
		blog.Errorf("delete inst, save audit log failed, err: %v, rid: %s", err, kit.Rid)
		return kit.CCError.Error(common.CCErrAuditSaveLogFailed)
	}

	return nil
}

// DeleteInstByInstID batch delete instance by inst id
func (c *commonInst) DeleteInstByInstID(kit *rest.Kit, objectID string, instID []int64, needCheckHost bool) error {
	cond := map[string]interface{}{
		common.GetInstIDField(objectID): map[string]interface{}{common.BKDBIN: instID},
	}
	if metadata.IsCommon(objectID) {
		cond[common.BKObjIDField] = objectID
	}

	return c.deleteInstByCond(kit, objectID, cond, needCheckHost)
}

func (c *commonInst) deleteInstByCond(kit *rest.Kit, objectID string, cond mapstr.MapStr, needCheckHost bool) error {
	query := &metadata.QueryInput{
		Condition: cond,
		Limit:     common.BKNoLimit,
	}

	instRsp, err := c.FindInst(kit, objectID, query)
	if err != nil {
		return err
	}

	if len(instRsp.Info) == 0 {
		return nil
	}

	delObjInstsMap, exists, err := c.hasHost(kit, instRsp.Info, objectID, needCheckHost)
	if err != nil {
		return err
	}
	if exists {
		return kit.CCError.Error(common.CCErrTopoHasHostCheckFailed)
	}

	bizSetMap := make(map[int64][]int64)
	audit := auditlog.NewInstanceAudit(c.clientSet.CoreService())
	auditLogs := make([]metadata.AuditLog, 0)

	for objID, delInsts := range delObjInstsMap {
		delInstIDs := make([]int64, len(delInsts))
		for index, instance := range delInsts {
			instID, err := instance.Int64(common.GetInstIDField(objID))
			if err != nil {
				blog.Errorf("can not convert ID to int64, err: %v, inst: %#v, rid: %s",
					err, instance, kit.Rid)
				return kit.CCError.CCErrorf(common.CCErrCommParamsNeedInt, common.GetInstIDField(objID))
			}
			delInstIDs[index] = instID

			if objID == common.BKInnerObjIDSet {
				bizID, err := instance.Int64(common.BKAppIDField)
				if err != nil {
					blog.Errorf("can not convert biz ID to int64, err: %v, set: %#v, rid: %s",
						err, instance, kit.Rid)
					return kit.CCError.CCErrorf(common.CCErrCommParamsNeedInt, common.BKAppIDField)
				}
				bizSetMap[bizID] = append(bizSetMap[bizID], instID)
			}
		}

		// if any instance has been bind to a instance by the association, then these instances should not be deleted.
		input := &metadata.Condition{
			Condition: mapstr.MapStr{common.BKDBOR: []mapstr.MapStr{
				{
					common.BKObjIDField:  objID,
					common.BKInstIDField: mapstr.MapStr{common.BKDBIN: delInstIDs},
				},
				{
					common.BKAsstObjIDField:  objID,
					common.BKAsstInstIDField: mapstr.MapStr{common.BKDBIN: delInstIDs},
				},
			}}}

		cnt, err := c.clientSet.CoreService().Association().CountInstanceAssociations(kit.Ctx, kit.Header, objID, input)
		if err != nil {
			blog.Errorf("count instance association failed, err: %v, rid: %s", err, kit.Rid)
			return err
		}

		if cnt.Count != 0 {
			return kit.CCError.CCError(common.CCErrorInstHasAsst)
		}

		// generate audit log.
		generateAuditParameter := auditlog.NewGenerateAuditCommonParameter(kit, metadata.AuditDelete)
		auditLog, err := audit.GenerateAuditLog(generateAuditParameter, objID, delInsts)
		if err != nil {
			blog.Errorf(" delete inst, generate audit log failed, err: %v, rid: %s", err, kit.Rid)
			return err
		}
		auditLogs = append(auditLogs, auditLog...)

		// delete this instance now.
		delCond := map[string]interface{}{
			common.GetInstIDField(objID): map[string]interface{}{common.BKDBIN: delInstIDs},
		}
		if metadata.IsCommon(objID) {
			delCond[common.BKObjIDField] = objID
		}
		dc := &metadata.DeleteOption{Condition: delCond}
		_, err = c.clientSet.CoreService().Instance().DeleteInstance(kit.Ctx, kit.Header, objID, dc)
		if err != nil {
			blog.Errorf("delete inst failed, err: %v, cond: %#v, rid: %s", err, delCond, kit.Rid)
			return err
		}
	}

	// clear set template sync status for set instances
	for bizID, setIDs := range bizSetMap {
		if len(setIDs) != 0 {
			ccErr := c.clientSet.CoreService().SetTemplate().DeleteSetTemplateSyncStatus(kit.Ctx, kit.Header, bizID,
				setIDs)
			if ccErr != nil {
				blog.Errorf("failed to delete set template sync status failed, bizID: %d, setIDs: %+v, err: %v, "+
					"rid: %s", bizID, setIDs, ccErr, kit.Rid)
				return ccErr
			}
		}
	}

	err = audit.SaveAuditLog(kit, auditLogs...)
	if err != nil {
		blog.Errorf("delete inst, save audit log failed, err: %v, rid: %s", err, kit.Rid)
		return kit.CCError.Error(common.CCErrAuditSaveLogFailed)
	}

	return nil
}

// FindInst search instance by condition
func (c *commonInst) FindInst(kit *rest.Kit, objID string, cond *metadata.QueryInput) (*metadata.InstResult, error) {

	result := new(metadata.InstResult)
	switch objID {
	case common.BKInnerObjIDHost:
		rsp, err := c.clientSet.CoreService().Host().GetHosts(kit.Ctx, kit.Header, cond)
		if err != nil {
			blog.Errorf("search object(%s) inst by the condition(%#v) failed, err: %v, rid: %s", objID, cond,
				err, kit.Rid)
			return nil, err
		}

		result.Count = rsp.Count
		result.Info = rsp.Info
		return result, nil

	default:
		input := &metadata.QueryCondition{Condition: cond.Condition, TimeCondition: cond.TimeCondition}
		input.Page.Start = cond.Start
		input.Page.Limit = cond.Limit
		input.Page.Sort = cond.Sort
		input.Fields = strings.Split(cond.Fields, ",")
		rsp, err := c.clientSet.CoreService().Instance().ReadInstance(kit.Ctx, kit.Header, objID, input)
		if err != nil {
			blog.Errorf("search object(%s) inst by the condition(%#v) failed, err: %v, rid: %s", objID, cond,
				err, kit.Rid)
			return nil, err
		}

		result.Count = rsp.Count
		result.Info = rsp.Info
		return result, nil
	}
}

// UpdateInst update instance by condition
func (c *commonInst) UpdateInst(kit *rest.Kit, cond, data mapstr.MapStr, objID string) error {
	// not allowed to update these fields, need to use specialized function
	data.Remove(common.BKParentIDField)
	data.Remove(common.BKAppIDField)

	inputParams := metadata.UpdateOption{
		Data:      data,
		Condition: cond,
	}

	// generate audit log of instance.
	audit := auditlog.NewInstanceAudit(c.clientSet.CoreService())
	generateAuditParameter := auditlog.NewGenerateAuditCommonParameter(kit, metadata.AuditUpdate).WithUpdateFields(data)
	auditLog, ccErr := audit.GenerateAuditLogByCondGetData(generateAuditParameter, objID, cond)
	if ccErr != nil {
		blog.Errorf(" update inst, generate audit log failed, err: %v, rid: %s", ccErr, kit.Rid)
		return ccErr
	}

	// to update.
	_, err := c.clientSet.CoreService().Instance().UpdateInstance(kit.Ctx, kit.Header, objID, &inputParams)
	if err != nil {
		blog.Errorf("update the object(%s) inst by the condition(%#v) failed, err: %v, rid: %s", objID, cond,
			err, kit.Rid)
		return err
	}

	// save audit log.
	err = audit.SaveAuditLog(kit, auditLog...)
	if err != nil {
		blog.Errorf("create inst, save audit log failed, err: %v, rid: %s", err, kit.Rid)
		return kit.CCError.Error(common.CCErrAuditSaveLogFailed)
	}
	return nil
}

// SearchObjectInstances searches object instances.
func (c *commonInst) SearchObjectInstances(kit *rest.Kit, objID string, input *metadata.CommonSearchFilter) (
	*metadata.CommonSearchResult, error) {

	// search conditions.
	cond, err := input.GetConditions()
	if err != nil {
		return nil, kit.CCError.Errorf(common.CCErrCommParamsInvalid, err)
	}

	conditions := &metadata.QueryCondition{
		Fields:         input.Fields,
		Condition:      cond,
		Page:           input.Page,
		DisableCounter: true,
	}

	// search object instances.
	resp, err := c.clientSet.CoreService().Instance().ReadInstance(kit.Ctx, kit.Header, objID, conditions)
	if err != nil {
		blog.Errorf("search object instances failed, err: %s, rid: %s", err.Error(), kit.Rid)
		return nil, err
	}

	result := &metadata.CommonSearchResult{}
	for idx := range resp.Info {
		result.Info = append(result.Info, &resp.Info[idx])
	}

	return result, nil
}

// CountObjectInstances counts object instances num.
func (c *commonInst) CountObjectInstances(kit *rest.Kit, objID string,
	input *metadata.CommonCountFilter) (*metadata.CommonCountResult, error) {

	// count conditions.
	cond, err := input.GetConditions()
	if err != nil {
		return nil, kit.CCError.Errorf(common.CCErrCommParamsInvalid, err)
	}

	conditions := &metadata.Condition{
		Condition: cond,
	}

	// count object instances num.
	resp, err := c.clientSet.CoreService().Instance().CountInstances(kit.Ctx, kit.Header, objID, conditions)
	if err != nil {
		blog.Errorf("count object instances failed, err: %s, rid: %s", err.Error(), kit.Rid)
		return nil, err
	}

	return &metadata.CommonCountResult{Count: resp.Count}, nil
}

// FindInstChildTopo find instance's child topo
func (c *commonInst) FindInstChildTopo(kit *rest.Kit, obj metadata.Object, instID int64, query *metadata.QueryInput) (
	int, []*CommonInstTopo, error) {

	return c.findInstTopo(kit, obj, instID, query, true)
}

// FindInstParentTopo find instance's parent topo
func (c *commonInst) FindInstParentTopo(kit *rest.Kit, obj metadata.Object, instID int64, query *metadata.QueryInput) (
	int, []*CommonInstTopo, error) {

	return c.findInstTopo(kit, obj, instID, query, false)
}

func (c *commonInst) findInstTopo(kit *rest.Kit, obj metadata.Object, instID int64, query *metadata.QueryInput,
	needChild bool) (int, []*CommonInstTopo, error) {

	results := make([]*CommonInstTopo, 0)
	if query == nil {
		query = &metadata.QueryInput{
			Condition: mapstr.MapStr{
				obj.GetInstIDFieldName(): instID,
			},
		}
	}

	insts, err := c.FindInst(kit, obj.ObjectID, query)
	if err != nil {
		return 0, nil, err
	}

	tmpResults := map[string]*CommonInstTopo{}
	for _, inst := range insts.Info {

		topoInsts, err := c.getObjectWithInsts(kit, obj, inst, needChild)
		if err != nil {
			return 0, nil, err
		}

		for _, topoInst := range topoInsts {
			object := topoInst.Object
			commonInst, exists := tmpResults[object.ObjectID]
			if !exists {
				commonInst = &CommonInstTopo{}
				commonInst.ObjectName = object.ObjectName
				commonInst.ObjIcon = object.ObjIcon
				commonInst.ObjID = object.ObjectID
				commonInst.Children = []metadata.InstNameAsst{}
				tmpResults[object.ObjectID] = commonInst
			}

			commonInst.Count = commonInst.Count + len(topoInst.Insts)

			for _, inst := range topoInst.Insts {

				instAsst := metadata.InstNameAsst{}
				id, err := inst.Int64(metadata.GetInstIDFieldByObjID(object.ObjectID))
				if err != nil {
					return 0, nil, err
				}

				name, err := inst.String(metadata.GetInstNameFieldName(object.ObjectID))
				if err != nil {
					return 0, nil, err
				}

				instAsst.ID = strconv.Itoa(int(id))
				instAsst.InstID = id
				instAsst.InstName = name
				instAsst.ObjectName = object.ObjectName
				instAsst.ObjIcon = object.ObjIcon
				instAsst.ObjID = object.ObjectID
				assoID, err := inst.Int64("asso_id")
				if err != nil {
					blog.Errorf("failed to get the inst id, err: %v, rid: %s", err, kit.Rid)
					return 0, nil, err
				}
				instAsst.AssoID = assoID

				tmpResults[object.ObjectID].Children = append(tmpResults[object.ObjectID].Children, instAsst)
			}
		}
	}

	for _, subResult := range tmpResults {
		results = append(results, subResult)
	}

	return len(results), results, nil
}

// FindInstTopo find instance all topo which include it's child and parent
func (c *commonInst) FindInstTopo(kit *rest.Kit, obj metadata.Object, instID int64, query *metadata.QueryInput) (int,
	[]CommonInstanceTopo, error) {

	if query == nil {
		query = &metadata.QueryInput{Condition: mapstr.MapStr{obj.GetInstIDFieldName(): instID}}
	}

	insts, err := c.FindInst(kit, obj.ObjectID, query)
	if err != nil {
		blog.Errorf("failed to find the inst, err: %v, rid: %s", err, kit.Rid)
		return 0, nil, err
	}

	results := make([]CommonInstanceTopo, 0)
	for _, inst := range insts.Info {
		id, err := inst.Int64(metadata.GetInstIDFieldByObjID(obj.ObjectID))
		if err != nil {
			blog.Errorf("failed to find the inst id, err: %v, rid: %s", err, kit.Rid)
			return 0, nil, err
		}

		name, err := inst.String(metadata.GetInstNameFieldName(obj.ObjectID))
		if err != nil {
			blog.Errorf("failed to find the inst name, err: %v, rid: %s", err, kit.Rid)
			return 0, nil, err
		}

		commonInst := metadata.InstNameAsst{}
		commonInst.ObjectName = obj.ObjectName
		commonInst.ObjID = obj.ObjectID
		commonInst.ObjIcon = obj.ObjIcon
		commonInst.InstID = id
		commonInst.ID = strconv.Itoa(int(id))
		commonInst.InstName = name

		_, parentInsts, err := c.FindInstParentTopo(kit, obj, id, nil)
		if err != nil {
			blog.Errorf("failed to find the inst, err: %v rid: %s", err, kit.Rid)
			return 0, nil, err
		}

		_, childInsts, err := c.FindInstChildTopo(kit, obj, id, nil)
		if err != nil {
			blog.Errorf("failed to find the inst, err: %v, rid: %s", err, kit.Rid)
			return 0, nil, err
		}

		results = append(results, CommonInstanceTopo{
			Prev: parentInsts,
			Next: childInsts,
			Curr: commonInst,
		})

	}

	return len(results), results, nil
}

func (c *commonInst) validMainLineParentID(kit *rest.Kit, objID string, data mapstr.MapStr) error {

	def, exist := data.Get(common.BKDefaultField)
	if exist && def.(int) != common.DefaultFlagDefaultValue {
		return nil
	}

	bizID, err := metadata.GetBizID(data)
	if err != nil {
		blog.Errorf("failed to parse the biz id, err: %v, rid: %s", err, kit.Rid)
		return kit.CCError.Errorf(common.CCErrCommParamsIsInvalid, common.BKAppIDField)
	}

	parentID, err := metadata.GetParentID(data)
	if err != nil {
		blog.Errorf("failed to parse the parent id, err: %v, rid: %s", err, kit.Rid)
		return kit.CCError.Errorf(common.CCErrCommParamsIsInvalid, common.BKParentIDField)
	}

	if err = c.validParentInstID(kit, objID, parentID, bizID); err != nil {
		blog.Errorf("parent id %d is invalid, err: %v, rid: %s", parentID, err, kit.Rid)
		return kit.CCError.Errorf(common.CCErrCommParamsIsInvalid, common.BKParentIDField)
	}
	return nil
}

func (c *commonInst) validParentInstID(kit *rest.Kit, objID string, instID int64, bizID int64) error {

	cond := &metadata.QueryCondition{
		Condition:      mapstr.MapStr{metadata.GetInstIDFieldByObjID(objID): instID},
		DisableCounter: true,
	}
	if bizID != 0 {
		cond.Condition.Set(common.BKAppIDField, bizID)
	}

	if metadata.IsCommon(objID) {
		cond.Condition.Set(common.BKObjIDField, objID)
	}

	rsp, err := c.clientSet.CoreService().Instance().ReadInstance(kit.Ctx, kit.Header, objID, cond)
	if err != nil {
		blog.Errorf("failed to read the object(%s) inst by the condition(%#v), err: %v, rid: %s", objID, cond,
			err, kit.Rid)
		return err
	}

	if len(rsp.Info) == 0 {
		return kit.CCError.Error(common.CCErrTopoInstSelectFailed)
	}

	return nil
}

func (c *commonInst) validObject(kit *rest.Kit, objID string, data mapstr.MapStr) error {

	if !metadata.IsCommon(objID) {
		blog.Errorf("object (%s) isn't common object, rid: %s", objID, kit.Rid)
		return kit.CCError.CCError(common.CCErrCommForbiddenOperateInnerModelInstanceWithCommonAPI)
	}

	input := &metadata.QueryCondition{
		Condition: mapstr.MapStr{common.BKObjIDField: objID},
		Fields:    []string{metadata.ModelFieldIsPaused},
	}
	rsp, err := c.clientSet.CoreService().Model().ReadModel(kit.Ctx, kit.Header, input)
	if err != nil {
		blog.Errorf("search object(%s) failed, err: %v, rid: %s", objID, err, kit.Rid)
		return err
	}

	if len(rsp.Info) == 0 {
		blog.Errorf("search object(%s) failed, object does not exist, rid: %s", objID, kit.Rid)
		return kit.CCError.CCError(common.CCErrTopoModuleSelectFailed)
	}

	// 暂停使用的model不允许创建实例
	if rsp.Info[0].IsPaused {
		blog.Errorf("object (%s) is paused, rid: %s", objID, kit.Rid)
		return kit.CCError.CCError(common.CCErrorTopoModelStopped)
	}

	cond := mapstr.MapStr{
		common.BKDBOR:                 []mapstr.MapStr{{common.BKObjIDField: objID}, {common.BKAsstObjIDField: objID}},
		common.AssociationKindIDField: common.AssociationKindMainline,
	}
	asst, err := c.clientSet.CoreService().Association().ReadModelAssociation(kit.Ctx, kit.Header,
		&metadata.QueryCondition{Condition: cond, DisableCounter: true})
	if err != nil {
		blog.Errorf("search object association failed, err: %v, rid: %s", err, kit.Rid)
		return err
	}

	if len(asst.Info) == 0 {
		return nil
	}

	if len(asst.Info) < 2 {
		blog.Errorf("mainline object(%s) does not get parent or children, rid: %s", objID, kit.Rid)
		return fmt.Errorf("mainline object(%s) does not get parent or children", objID)
	}

	var parentObjID string
	for _, item := range asst.Info {
		if item.ObjectID == objID {
			parentObjID = item.AsstObjID
		}
	}

	if err := c.validMainLineParentID(kit, parentObjID, data); err != nil {
		blog.Errorf("valid mainline object(%s) parentID failed, err: %v, rid: %s", objID, err, kit.Rid)
		return err
	}

	return nil
}

// IsMainline check if object is mainline object
func (c *commonInst) IsMainline(kit *rest.Kit, objID string) (bool, error) {
	filter := []map[string]interface{}{{
		common.BKDBOR:                 []mapstr.MapStr{{common.BKObjIDField: objID}, {common.BKAsstObjIDField: objID}},
		common.AssociationKindIDField: common.AssociationKindMainline,
	}}
	cnt, ccErr := c.clientSet.CoreService().Count().GetCountByFilter(kit.Ctx, kit.Header, common.BKTableNameObjAsst,
		filter)
	if ccErr != nil {
		blog.Errorf("count object(%s) mainline association failed, err: %v, rid: %s", objID, ccErr, kit.Rid)
		return false, ccErr
	}

	if cnt[0] != 0 {
		return true, nil
	}

	return false, nil
}

// hasHost get objID and instances map for mainline instances with its children topology, and check if they have hosts
func (c *commonInst) hasHost(kit *rest.Kit, instances []mapstr.MapStr, objID string, checkHost bool) (
	map[string][]mapstr.MapStr, bool, error) {

	if len(instances) == 0 {
		return nil, false, nil
	}

	instIDs := make([]int64, len(instances))
	for index, instance := range instances {
		instID, err := instance.Int64(common.GetInstIDField(objID))
		if err != nil {
			blog.Errorf("can not convert ID to int64, err: %v, inst: %#v, rid: %s", err, instance, kit.Rid)
			return nil, false, kit.CCError.CCErrorf(common.CCErrCommParamsNeedInt, common.GetInstIDField(objID))
		}
		instIDs[index] = instID
	}

	objInstMap := map[string][]mapstr.MapStr{objID: instances}

	var moduleIDs []int64
	if objID == common.BKInnerObjIDModule {
		moduleIDs = instIDs
	} else if objID == common.BKInnerObjIDSet {
		query := &metadata.QueryInput{
			Condition: map[string]interface{}{common.BKSetIDField: map[string]interface{}{common.BKDBIN: instIDs}},
			Limit:     common.BKNoLimit,
		}

		moduleRsp, err := c.FindInst(kit, common.BKInnerObjIDModule, query)
		if err != nil {
			blog.Errorf("find modules for set failed, err: %v, set IDs: %+v, rid: %s", err, instIDs, kit.Rid)
			return nil, false, err
		}

		if len(moduleRsp.Info) == 0 {
			return objInstMap, false, nil
		}

		objInstMap[common.BKInnerObjIDModule] = moduleRsp.Info
		moduleIDs = make([]int64, len(moduleRsp.Info))
		for index, module := range moduleRsp.Info {
			moduleID, err := module.Int64(common.BKModuleIDField)
			if err != nil {
				blog.Errorf("can not convert ID to int64, err: %v, module: %s, rid: %s", err, module, kit.Rid)
				return nil, false, kit.CCError.CCErrorf(common.CCErrCommParamsNeedInt, common.BKModuleIDField)
			}
			moduleIDs[index] = moduleID
		}
	} else {
		// get mainline object relation(excluding hosts) by mainline associations
		mainlineCond := &metadata.QueryCondition{
			Condition: map[string]interface{}{
				common.AssociationKindIDField: common.AssociationKindMainline,
				common.BKObjIDField: mapstr.MapStr{
					common.BKDBNIN: []string{common.BKInnerObjIDSet, common.BKInnerObjIDModule},
				},
			},
			DisableCounter: true,
		}
		asstRsp, err := c.clientSet.CoreService().Association().ReadModelAssociation(kit.Ctx, kit.Header, mainlineCond)
		if err != nil {
			blog.Errorf("search mainline association failed, error: %v, rid: %s", err, kit.Rid)
			return nil, false, err
		}

		objChildMap := make(map[string]string)
		isMainline := false
		for _, asst := range asstRsp.Info {
			if asst.ObjectID == common.BKInnerObjIDHost {
				continue
			}
			objChildMap[asst.AsstObjID] = asst.ObjectID
			if asst.AsstObjID == objID || asst.ObjectID == objID {
				isMainline = true
			}
		}

		if !isMainline {
			return objInstMap, false, nil
		}

		// loop through the child topology level to get all instances
		parentIDs := instIDs
		for childObjID := objChildMap[objID]; len(childObjID) != 0; childObjID = objChildMap[childObjID] {
			cond := map[string]interface{}{common.BKParentIDField: map[string]interface{}{common.BKDBIN: parentIDs}}
			if metadata.IsCommon(childObjID) {
				cond[metadata.ModelFieldObjectID] = childObjID
			}

			if childObjID == common.BKInnerObjIDSet {
				cond[common.BKDefaultField] = common.DefaultFlagDefaultValue
			}

			query := &metadata.QueryInput{
				Condition: cond,
				Limit:     common.BKNoLimit,
			}

			childRsp, err := c.FindInst(kit, childObjID, query)
			if err != nil {
				blog.Errorf("find children failed, err: %v, parent IDs: %+v, rid: %s", err, parentIDs, kit.Rid)
				return nil, false, err
			}

			if len(childRsp.Info) == 0 {
				return objInstMap, false, nil
			}

			parentIDs = make([]int64, len(childRsp.Info))
			for index, instance := range childRsp.Info {
				instID, err := instance.Int64(common.GetInstIDField(childObjID))
				if err != nil {
					blog.Errorf("can not convert ID to int64, err: %v, inst: %s, rid: %s",
						err, instance, kit.Rid)
					return nil, false, kit.CCError.CCErrorf(common.CCErrCommParamsNeedInt,
						common.GetInstIDField(childObjID))
				}
				parentIDs[index] = instID
			}

			if childObjID == common.BKInnerObjIDModule {
				moduleIDs = parentIDs
			}

			objInstMap[childObjID] = childRsp.Info
		}
	}

	// check if module contains hosts
	if checkHost && len(moduleIDs) > 0 {
		exists, err := c.innerHasHost(kit, moduleIDs)
		if err != nil {
			return nil, false, err
		}

		if exists {
			return nil, true, nil
		}
	}

	return objInstMap, false, nil
}

func (c *commonInst) innerHasHost(kit *rest.Kit, moduleIDS []int64) (bool, error) {
	option := &metadata.HostModuleRelationRequest{
		ModuleIDArr: moduleIDS,
		Fields:      []string{common.BKHostIDField},
		Page:        metadata.BasePage{Limit: 1},
	}
	rsp, err := c.clientSet.CoreService().Host().GetHostModuleRelation(kit.Ctx, kit.Header, option)
	if err != nil {
		blog.Errorf("searh host object relation failed, err: %v, rid: %s", err, kit.Rid)
		return false, kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	return len(rsp.Info) != 0, nil
}

// GetObjectWithInsts get object with insts, get parent or child depends on needChild
func (c *commonInst) getObjectWithInsts(kit *rest.Kit, object metadata.Object,
	inst mapstr.MapStr, needChild bool) ([]*ObjectWithInsts, error) {

	result := make([]*ObjectWithInsts, 0)

	cond := mapstr.New()
	if needChild {
		cond = mapstr.MapStr{common.BKObjIDField: object.ObjectID}
	} else {
		cond = mapstr.MapStr{common.BKAsstObjIDField: object.ObjectID}
	}

	objPairs, err := c.searchAssoObjects(kit, needChild, cond)
	if err != nil {
		blog.Errorf("failed to get the object(%s)'s parent, err: %v, rid: %s", object.ObjectID, err, kit.Rid)
		return result, err
	}

	currInstID, err := inst.Int64(metadata.GetInstIDFieldByObjID(object.ObjectID))
	if err != nil {
		blog.Errorf("failed to get the inst id, err: %v, rid: %s", err, kit.Rid)
		return result, err
	}

	for _, objPair := range objPairs {

		queryCond := &metadata.InstAsstQueryCondition{
			Cond: metadata.QueryCondition{Condition: mapstr.MapStr{
				common.BKAsstInstIDField:         currInstID,
				common.AssociationObjAsstIDField: objPair.Association.AssociationName,
			}},
			ObjID: objPair.Object.ObjectID,
		}

		if needChild {
			queryCond.Cond.Condition.Set(common.BKObjIDField, object.ObjectID)
			queryCond.Cond.Condition.Set(common.BKAsstObjIDField, objPair.Object.ObjectID)
		} else {
			queryCond.Cond.Condition.Set(common.BKObjIDField, objPair.Object.ObjectID)
			queryCond.Cond.Condition.Set(common.BKAsstObjIDField, object.ObjectID)
		}

		rsp, err := c.clientSet.CoreService().Association().ReadInstAssociation(kit.Ctx, kit.Header, queryCond)
		if err != nil {
			blog.Errorf("search inst association failed , err: %v, rid: %s", err, kit.Rid)
			return nil, err
		}

		// found no noe inst association with this object and association info.
		// which means that, this object association has not been instantiated.
		if len(rsp.Info) == 0 {
			continue
		}

		relation := make(map[int64]int64)
		InstIDS := []int64{}
		for _, item := range rsp.Info {

			InstID := item.InstID
			relation[InstID] = item.ID
			InstIDS = append(InstIDS, InstID)
		}

		innerCond := mapstr.MapStr{objPair.Object.GetInstIDFieldName(): mapstr.MapStr{common.BKDBIN: InstIDS}}
		if objPair.Object.IsCommon() {
			innerCond.Set(metadata.ModelFieldObjectID, objPair.Object.ObjectID)
		}

		rspItems, err := c.FindInst(kit, objPair.Object.ObjectID, &metadata.QueryInput{Condition: innerCond})
		if err != nil {
			blog.Errorf("failed to search the insts by the condition(%#v), err: %v, rid: %s",
				innerCond, err, kit.Rid)
			return result, err
		}

		for _, item := range rspItems.Info {
			id, err := item.Int64(metadata.GetInstIDFieldByObjID(objPair.Object.ObjectID))
			if err != nil {
				blog.Errorf("failed to parse the instance id , err: %v, rid: %s", err, kit.Rid)
				return result, err
			}
			item.Set("asso_id", relation[id])
		}

		rstObj := &ObjectWithInsts{Object: objPair.Object, Insts: rspItems.Info}
		result = append(result, rstObj)

	}

	return result, nil
}

// TODO maybe can add this to model/association
func (c *commonInst) searchAssoObjects(kit *rest.Kit, isNeedChild bool, cond mapstr.MapStr) ([]ObjectAssoPair,
	error) {
	// TODO after merge can replace to SearchObjectAssociation
	rsp, err := c.clientSet.CoreService().Association().ReadModelAssociation(kit.Ctx, kit.Header,
		&metadata.QueryCondition{Condition: cond, DisableCounter: true})
	if err != nil {
		blog.Errorf("search object association failed, err: %v, rid: %s", err, kit.Rid)
		return nil, err
	}

	objAssoMap := make(map[string]metadata.Association, 0)
	pair := make([]ObjectAssoPair, 0)
	var objIDArray []string
	for _, asst := range rsp.Info {
		if isNeedChild {
			objIDArray = append(objIDArray, asst.AsstObjID)
			objAssoMap[asst.AsstObjID] = asst
		} else {
			objIDArray = append(objIDArray, asst.ObjectID)
			objAssoMap[asst.ObjectID] = asst
		}
	}

	cond = mapstr.MapStr{metadata.ModelFieldObjectID: mapstr.MapStr{common.BKDBIN: objIDArray}}

	// TODO after merge can replace by SearchObject
	rspRst, err := c.clientSet.CoreService().Model().ReadModel(kit.Ctx, kit.Header,
		&metadata.QueryCondition{Condition: cond, DisableCounter: true})
	if err != nil {
		blog.Errorf("failed to search the object by cond(%#v), err: %v, rid: %s", cond, err, kit.Rid)
		return nil, kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if len(rspRst.Info) == 0 {
		blog.Errorf("search asso object, but can not found object with cond: %v, rid: %s", cond, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrTopoModuleSelectFailed)
	}

	for _, object := range rspRst.Info {
		pair = append(pair, ObjectAssoPair{
			Object:      object,
			Association: objAssoMap[object.ObjectID],
		})
	}

	return pair, nil
}
