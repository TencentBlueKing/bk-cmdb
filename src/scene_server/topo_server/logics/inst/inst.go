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

// Package inst TODO
package inst

import (
	"fmt"
	"strconv"
	"strings"

	"configcenter/pkg/filter"
	filtertools "configcenter/pkg/tools/filter"
	"configcenter/src/ac/extensions"
	"configcenter/src/apimachinery"
	"configcenter/src/common"
	"configcenter/src/common/auditlog"
	"configcenter/src/common/blog"
	httpheader "configcenter/src/common/http/header"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/language"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	gparams "configcenter/src/common/paraparse"
	"configcenter/src/common/util"
)

// InstOperationInterface inst operation methods
type InstOperationInterface interface {
	// CreateInst create instance by object and create message
	CreateInst(kit *rest.Kit, objID string, data mapstr.MapStr) (mapstr.MapStr, error)
	// CreateManyInstance batch create instance by object and create message
	CreateManyInstance(kit *rest.Kit, objID string, data []mapstr.MapStr) (*metadata.CreateManyCommInstResultDetail,
		error)
	// BatchCreateInstance batch create instance, if one of instances fails to create, an error is returned.
	BatchCreateInstance(kit *rest.Kit, objID string, data []mapstr.MapStr) (*metadata.BatchCreateInstRespData, error)
	// CreateInstBatch batch create instance by excel
	CreateInstBatch(kit *rest.Kit, objID string, batchInfo *metadata.InstBatchInfo) (*metadata.ImportInstRes, error)
	// DeleteInst delete instance by objectid and condition
	DeleteInst(kit *rest.Kit, objectID string, cond mapstr.MapStr, needCheckHost bool) error
	// DeleteInstByInstID batch delete instance by inst id
	DeleteInstByInstID(kit *rest.Kit, objectID string, instID []int64, needCheckHost bool) error
	// FindInst search instance by condition
	FindInst(kit *rest.Kit, objID string, cond *metadata.QueryCondition) (*metadata.InstResult, error)
	// FindInstByAssociationInst deprecated function.
	FindInstByAssociationInst(kit *rest.Kit, objID string, asstParamCond *AssociationParams) (*metadata.InstResult,
		error)
	// UpdateInst update instance by condition
	UpdateInst(kit *rest.Kit, cond, data mapstr.MapStr, objID string) error
	// SearchObjectInstances searches object instances.
	SearchObjectInstances(kit *rest.Kit, objID string, input *metadata.SearchInstanceFilter) (
		*metadata.CommonSearchResult, error)
	// CountObjectInstances counts object instances num.
	CountObjectInstances(kit *rest.Kit, objID string, input *metadata.CountInstanceFilter) (*metadata.CommonCountResult,
		error)
	// FindInstChildTopo find instance's child topo
	FindInstChildTopo(kit *rest.Kit, objID string, instID int64) (int, []*metadata.CommonInstTopo, error)
	// FindInstTopo find instance all topo which include it's child and parent
	FindInstTopo(kit *rest.Kit, obj metadata.Object, instID int64) (int, []metadata.CommonInstTopoV2, error)
	// SetProxy proxy the interface
	SetProxy(instAssoc AssociationOperationInterface)
}

// NewInstOperation create a new inst operation instance
func NewInstOperation(client apimachinery.ClientSetInterface, lang language.CCLanguageIf,
	authManager *extensions.AuthManager) InstOperationInterface {

	return &commonInst{
		clientSet:   client,
		language:    lang,
		authManager: authManager,
	}
}

// ObjectWithInsts a struct include object msg and insts array
type ObjectWithInsts struct {
	Object metadata.Object
	Insts  []mapstr.MapStr
}

// ObjectAssoPair a struct include object msg and association
type ObjectAssoPair struct {
	Object    metadata.Object
	AssocName string
}

// ConditionItem subcondition
type ConditionItem struct {
	Field    string      `json:"field,omitempty"`
	Operator string      `json:"operator,omitempty"`
	Value    interface{} `json:"value,omitempty"`
}

// AssociationParams  association params
type AssociationParams struct {
	Page      metadata.BasePage          `json:"page,omitempty"`
	Fields    map[string][]string        `json:"fields,omitempty"`
	Condition map[string][]ConditionItem `json:"condition,omitempty"`
	// 非必填，只能用来查时间，且与Condition是与关系
	TimeCondition *metadata.TimeCondition `json:"time_condition,omitempty"`
}

type commonInst struct {
	clientSet   apimachinery.ClientSetInterface
	language    language.CCLanguageIf
	authManager *extensions.AuthManager
	asst        AssociationOperationInterface
}

// SetProxy proxy the interface
func (c *commonInst) SetProxy(instAssoc AssociationOperationInterface) {
	c.asst = instAssoc
}

// CreateInst create instance by object and create message
func (c *commonInst) CreateInst(kit *rest.Kit, objID string, data mapstr.MapStr) (mapstr.MapStr, error) {

	if err := c.validObject(kit, objID, data); err != nil {
		blog.Errorf("check object (%s) if is mainline object failed, err: %v, rid: %s", objID, err, kit.Rid)
		return nil, err
	}

	if metadata.IsCommon(objID) {
		data.Set(common.BKObjIDField, objID)
	}

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

	input := &metadata.QueryCondition{Condition: mapstr.MapStr{metadata.GetInstIDFieldByObjID(objID): rsp.Created.ID}}
	inst, err := c.FindInst(kit, objID, input)
	if err != nil {
		blog.Errorf("search instance by inst_id(%s) failed, err: %v, rid: %s", rsp.Created.ID, err, kit.Rid)
		return nil, err
	}

	if len(inst.Info) != 1 {
		blog.Errorf("search instance by inst_id(%s) failed, get %d instance, rid: %s", rsp.Created.ID,
			len(inst.Info), kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrTopoInstSelectFailed)
	}

	// for audit log.
	generateAuditParameter := auditlog.NewGenerateAuditCommonParameter(kit, metadata.AuditCreate)
	audit := auditlog.NewInstanceAudit(c.clientSet.CoreService())
	auditLog, err := audit.GenerateAuditLog(generateAuditParameter, objID, inst.Info)
	if err != nil {
		blog.Errorf(" creat inst, generate audit log failed, err: %v, rid: %s", err, kit.Rid)
		return nil, err
	}

	err = audit.SaveAuditLog(kit, auditLog...)
	if err != nil {
		blog.Errorf("create inst, save audit log failed, err: %v, rid: %s", err, kit.Rid)
		return nil, kit.CCError.Error(common.CCErrAuditSaveLogFailed)
	}

	return inst.Info[0], nil
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
	auditLog, rawErr := audit.GenerateAuditLogByCondGetData(generateAuditParameter, objID, cond)
	if rawErr != nil {
		blog.Errorf("create many instances, generate audit log failed, err: %v, rid: %s",
			rawErr, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrAuditGenerateLogFailed, rawErr.Error())
	}

	// save audit log.
	rawErr = audit.SaveAuditLog(kit, auditLog...)
	if rawErr != nil {
		blog.Errorf("creat many instances, save audit log failed, err: %v, rid: %s", rawErr, kit.Rid)
		return nil, kit.CCError.Error(common.CCErrAuditSaveLogFailed)
	}
	return resp, nil
}

// BatchCreateInstance batch create instance, if one of instances fails to create, an error is returned.
func (c *commonInst) BatchCreateInstance(kit *rest.Kit, objID string, data []mapstr.MapStr) (
	*metadata.BatchCreateInstRespData, error) {

	params := &metadata.BatchCreateModelInstOption{Data: data}
	res, err := c.clientSet.CoreService().Instance().BatchCreateInstance(kit.Ctx, kit.Header, objID, params)
	if err != nil {
		blog.Errorf("failed to save the object(%s) instances, err: %v, rid: %s", objID, err, kit.Rid)
		return nil, err
	}

	if len(res.IDs) == 0 {
		return res, nil
	}

	// generate audit log of instance.
	for i, id := range res.IDs {
		data[i][common.BKModuleIDField] = id
	}
	audit := auditlog.NewInstanceAudit(c.clientSet.CoreService())
	generateAuditParameter := auditlog.NewGenerateAuditCommonParameter(kit, metadata.AuditCreate)
	auditLog, rawErr := audit.GenerateAuditLog(generateAuditParameter, objID, data)
	if rawErr != nil {
		blog.Errorf("create many instances, generate audit log failed, err: %v, rid: %s",
			rawErr, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrAuditGenerateLogFailed, rawErr.Error())
	}

	// save audit log.
	rawErr = audit.SaveAuditLog(kit, auditLog...)
	if rawErr != nil {
		blog.Errorf("creat many instances, save audit log failed, err: %v, rid: %s", rawErr, kit.Rid)
		return nil, kit.CCError.Error(common.CCErrAuditSaveLogFailed)
	}

	return res, nil
}

// createInstBatch batch create instance by excel
func (c *commonInst) createInstBatch(kit *rest.Kit, objID string, batchInfo *metadata.InstBatchInfo) (
	*metadata.ImportInstRes, error) {

	relRes, err := c.getObjRelationDestMsg(kit, objID)
	if err != nil {
		blog.Errorf("get object relation failed, err: %v, rid: %s", err, kit.Rid)
		return nil, err
	}

	ccLang := c.language.CreateDefaultCCLanguageIf(httpheader.GetLanguage(kit.Header))
	result := new(metadata.ImportInstRes)

	for idx, inst := range batchInfo.BatchInfo {
		if inst == nil {
			// ignore empty excel line
			continue
		}
		delete(inst, "import_from")

		// 实例id 为空，表示要新建实例
		// 实例ID已经赋值，更新数据.  (已经赋值, value not equal 0 or nil)

		instID, exist := inst[common.BKInstIDField]
		if exist && (instID == "" || instID == nil) {
			exist = false
		}

		// use new transaction, need a new header
		kit.Header = kit.NewHeader()
		_ = c.clientSet.CoreService().Txn().AutoRunTxn(kit.Ctx, kit.Header, func() error {
			tableData, err := metadata.GetTableData(inst, relRes)
			if err != nil {
				result.Errors = append(result.Errors, ccLang.Languagef("import_row_int_error_str", idx,
					err.Error()))
				return err
			}

			if exist {
				instID, err := util.GetInt64ByInterface(instID)
				if err != nil {
					result.Errors = append(result.Errors, ccLang.Languagef("import_row_int_error_str", idx,
						err.Error()))
					return err
				}

				if err := c.updateInstByExcel(kit, objID, tableData, instID, inst); err != nil {
					blog.Errorf("update instance by excel failed, err: %v, rid: %s", err, kit.Rid)
					result.Errors = append(result.Errors, ccLang.Languagef("import_row_int_error_str", idx,
						err.Error()))
					return err
				}

				result.Success = append(result.Success, idx)
				return nil
			}

			if err := c.createInstByExcel(kit, objID, tableData, inst); err != nil {
				blog.Errorf("create instance by excel failed, err: %v, rid: %s", err, kit.Rid)
				result.Errors = append(result.Errors, ccLang.Languagef("import_row_int_error_str", idx, err.Error()))
				return err
			}

			result.Success = append(result.Success, idx)
			return nil
		})
	}

	return result, nil
}

func (c *commonInst) updateInstByExcel(kit *rest.Kit, objID string, tableData *metadata.TableData, instID int64,
	inst mapstr.MapStr) error {

	input := &metadata.QueryCondition{Condition: mapstr.MapStr{common.BKInstIDField: instID}}
	preInst, err := c.FindInst(kit, objID, input)
	if err != nil {
		blog.Errorf("find instance failed, objID: %s, input: %v, err: %v, rid: %s", objID, input, err, kit.Rid)
		return err
	}

	if len(preInst.Info) == 0 {
		blog.Errorf("instance does not exist, objID: %s, instID: %v, err: %v, rid: %s", objID, instID, err, kit.Rid)
		return fmt.Errorf("instance: %d does not exist", instID)
	}

	// update instance table field type data
	if tableData != nil {
		properties, err := c.updateTableData(kit, tableData, instID)
		if err != nil {
			blog.ErrorJSON("update table data failed, data: %s, err: %s, rid: %s", inst, err, kit.Rid)
			return err
		}

		for property, ids := range properties {
			inst.Set(property, ids)
		}
	}

	if err := c.UpdateInst(kit, mapstr.MapStr{common.BKInstIDField: instID}, inst, objID); err != nil {
		blog.Errorf("failed to update the object(%s) inst data (%#v), err: %v, rid: %s", objID, inst, err, kit.Rid)
		return err
	}

	audit := auditlog.NewInstanceAudit(c.clientSet.CoreService())
	auditParam := auditlog.NewGenerateAuditCommonParameter(kit, metadata.AuditUpdate).WithUpdateFields(preInst.Info[0])
	auditLog, err := audit.GenerateAuditLog(auditParam, objID, []mapstr.MapStr{inst})
	if err := audit.SaveAuditLog(kit, auditLog...); err != nil {
		blog.Errorf("save audit failed, err: %v, rid: %s", err, kit.Rid)
		return err
	}

	return nil
}

func (c *commonInst) createInstByExcel(kit *rest.Kit, objID string, tableData *metadata.TableData,
	inst mapstr.MapStr) error {

	// add host table field type data
	if tableData != nil {
		properties, err := c.addTableData(kit, tableData, 0)
		if err != nil {
			blog.ErrorJSON("add table data failed, data: %s, err: %s, rid: %s", tableData, err, kit.Rid)
			return err
		}

		for property, ids := range properties {
			inst.Set(property, ids)
		}
	}

	inst.Set(common.BKObjIDField, objID)
	instCond := &metadata.CreateModelInstance{Data: inst}
	rsp, err := c.clientSet.CoreService().Instance().CreateInstance(kit.Ctx, kit.Header, objID, instCond)
	if err != nil {
		blog.Errorf("failed to create object instance, err: %v, rid: %s", err, kit.Rid)
		return err
	}
	inst[common.BKInstIDField] = rsp.Created.ID

	// to generate audit log.
	audit := auditlog.NewInstanceAudit(c.clientSet.CoreService())
	generateAuditParameter := auditlog.NewGenerateAuditCommonParameter(kit, metadata.AuditCreate)
	auditLog, err := audit.GenerateAuditLog(generateAuditParameter, objID, []mapstr.MapStr{inst})
	if err != nil {
		blog.Errorf("save audit failed, err: %v, rid: %s", err, kit.Rid)
		return err
	}

	err = audit.SaveAuditLog(kit, auditLog...)
	if err != nil {
		blog.Errorf("create inst, save audit log failed, err: %v, rid: %s", err, kit.Rid)
		return err
	}

	return nil
}

// TODO 与hostserver的getHostRelationDestMsg方法有重复逻辑，后续重构调整
// getObjRelationDestMsg get object relation, it can only get bk_property_id and dest_model field
func (c *commonInst) getObjRelationDestMsg(kit *rest.Kit, objID string) ([]metadata.ModelQuoteRelation, error) {
	opt := &metadata.CommonQueryOption{
		CommonFilterOption: metadata.CommonFilterOption{
			Filter: &filter.Expression{
				RuleFactory: &filter.CombinedRule{
					Condition: filter.And,
					Rules: []filter.RuleFactory{
						&filter.AtomRule{
							Field:    common.BKSrcModelField,
							Operator: filter.OpFactory(filter.Equal),
							Value:    objID,
						},
					},
				},
			},
		},
		Page: metadata.BasePage{
			Limit: common.BKNoLimit,
		},
		Fields: []string{common.BKPropertyIDField, common.BKDestModelField},
	}

	relRes, err := c.clientSet.CoreService().ModelQuote().ListModelQuoteRelation(kit.Ctx, kit.Header, opt)
	if err != nil {
		return nil, err
	}
	return relRes.Info, nil
}

// TODO 与hostserver的addTableData方法重复，后续重构处理
func (c *commonInst) addTableData(kit *rest.Kit, tableData *metadata.TableData, id int64) (map[string][]uint64, error) {
	audit := auditlog.NewQuotedInstAuditLog(c.clientSet.CoreService())
	genAuditParams := auditlog.NewGenerateAuditCommonParameter(kit, metadata.AuditCreate)
	var auditLogs []metadata.AuditLog
	result := make(map[string][]uint64)
	for destModel, table := range tableData.ModelData {
		if id != 0 {
			for idx := range table {
				table[idx].Set(common.BKInstIDField, id)
			}
		}
		ids, err := c.clientSet.CoreService().ModelQuote().BatchCreateQuotedInstance(kit.Ctx, kit.Header, destModel,
			table)
		if err != nil {
			return nil, err
		}

		propertyID := tableData.ModelPropertyRel[destModel]
		// generate and save audit logs
		for i := range table {
			table[i][common.BKFieldID] = ids[i]
			result[propertyID] = append(result[propertyID], ids[i])
		}

		auditLog, ccErr := audit.GenerateAuditLog(genAuditParams, destModel, tableData.SrcModel, propertyID, table)
		if ccErr != nil {
			return nil, ccErr
		}
		auditLogs = append(auditLogs, auditLog...)
	}

	if err := audit.SaveAuditLog(kit, auditLogs...); err != nil {
		return nil, kit.CCError.Error(common.CCErrAuditSaveLogFailed)
	}

	return result, nil
}

// TODO 与hostserver的updateTableData方法重复，后续重构处理
// updateTableData update table data
func (c *commonInst) updateTableData(kit *rest.Kit, tableData *metadata.TableData, id int64) (
	map[string][]uint64, error) {

	audit := auditlog.NewQuotedInstAuditLog(c.clientSet.CoreService())
	genAuditParams := auditlog.NewGenerateAuditCommonParameter(kit, metadata.AuditDelete)
	filterOpt := metadata.CommonFilterOption{
		Filter: filtertools.GenAtomFilter(common.BKInstIDField, filter.Equal, id),
	}
	listOpt := &metadata.CommonQueryOption{
		CommonFilterOption: filterOpt,
		Page:               metadata.BasePage{Limit: common.BKMaxPageSize},
	}

	var auditLogs []metadata.AuditLog
	for destModel := range tableData.ModelData {
		listRes, err := c.clientSet.CoreService().ModelQuote().ListQuotedInstance(kit.Ctx, kit.Header, destModel,
			listOpt)
		if err != nil {
			return nil, err
		}

		err = c.clientSet.CoreService().ModelQuote().BatchDeleteQuotedInstance(kit.Ctx, kit.Header, destModel,
			&filterOpt)
		if err != nil {
			return nil, err
		}

		auditLog, ccErr := audit.GenerateAuditLog(genAuditParams, destModel, tableData.SrcModel,
			tableData.ModelPropertyRel[destModel], listRes.Info)
		if ccErr != nil {
			return nil, err
		}

		auditLogs = append(auditLogs, auditLog...)
	}

	// save audit logs
	ccErr := audit.SaveAuditLog(kit, auditLogs...)
	if ccErr != nil {
		return nil, kit.CCError.Error(common.CCErrAuditSaveLogFailed)
	}

	properties, err := c.addTableData(kit, tableData, id)
	if err != nil {
		return nil, err
	}

	return properties, nil
}

// CreateInstBatch batch create instance by excel
func (c *commonInst) CreateInstBatch(kit *rest.Kit, objID string, batchInfo *metadata.InstBatchInfo) (
	*metadata.ImportInstRes, error) {

	// forbidden create inner model instance with common api
	if common.IsInnerModel(objID) {
		blog.Errorf("create %s instance with common create api forbidden, rid: %s", objID, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrCommForbiddenOperateInnerModelInstanceWithCommonAPI)

	}

	// forbidden create mainline instance with common api
	filter := []map[string]interface{}{{
		common.BKDBOR:                 []mapstr.MapStr{{common.BKObjIDField: objID}, {common.BKAsstObjIDField: objID}},
		common.AssociationKindIDField: common.AssociationKindMainline,
	}}
	cnt, ccErr := c.clientSet.CoreService().Count().GetCountByFilter(kit.Ctx, kit.Header, common.BKTableNameObjAsst,
		filter)
	if ccErr != nil {
		blog.Errorf("count object(%s) mainline association failed, err: %v, rid: %s", objID, ccErr, kit.Rid)
		return nil, ccErr
	}

	if cnt[0] != 0 {
		return nil, kit.CCError.CCError(common.CCErrCommForbiddenOperateMainlineInstanceWithCommonAPI)
	}

	if batchInfo.InputType != common.InputTypeExcel {
		return nil, kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, "input_type")
	}
	if len(batchInfo.BatchInfo) == 0 {
		return nil, kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, "BatchInfo")
	}

	// 检查实例与URL参数指定的模型一致
	for line, inst := range batchInfo.BatchInfo {
		objectID, exist := inst[common.BKObjIDField]
		if exist && objectID != objID {
			blog.Errorf("create object[%s] instance batch failed, bk_obj_id field conflict with url field, rid: %s",
				objID, kit.Rid)
			return nil, kit.CCError.Errorf(common.CCErrorTopoObjectInstanceObjIDFieldConflictWithURL, line)
		}
	}

	result, err := c.createInstBatch(kit, objID, batchInfo)
	if err != nil {
		blog.Errorf("create inst by export failed, err: %v, rid: %s", err, kit.Rid)
		return nil, err
	}

	return result, nil
}

// DeleteInst delete instance by objectid and condition
func (c *commonInst) DeleteInst(kit *rest.Kit, objectID string, cond mapstr.MapStr, needCheckHost bool) error {
	query := &metadata.QueryCondition{
		Condition: cond,
		Page:      metadata.BasePage{Limit: common.BKNoLimit},
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

	audit := auditlog.NewInstanceAudit(c.clientSet.CoreService())
	auditLogs := make([]metadata.AuditLog, 0)

	for objID, delInsts := range delObjInstsMap {
		// generate audit log.
		generateAuditParameter := auditlog.NewGenerateAuditCommonParameter(kit, metadata.AuditDelete)
		auditLog, err := audit.GenerateAuditLog(generateAuditParameter, objID, delInsts)
		if err != nil {
			blog.Errorf("generate audit log failed, err: %v, rid: %s", err, kit.Rid)
			return err
		}

		auditLogs = append(auditLogs, auditLog...)
		if err := c.deleteInsts(kit, delInsts, objID); err != nil {
			return err
		}
	}

	err = audit.SaveAuditLog(kit, auditLogs...)
	if err != nil {
		blog.Errorf("delete inst, save audit log failed, err: %v, rid: %s", err, kit.Rid)
		return kit.CCError.Error(common.CCErrAuditSaveLogFailed)
	}

	return nil
}

func (c *commonInst) deleteInsts(kit *rest.Kit, delInsts []mapstr.MapStr, objID string) error {

	delInstIDs := make([]int64, len(delInsts))
	for index, instance := range delInsts {
		instID, err := instance.Int64(common.GetInstIDField(objID))
		if err != nil {
			blog.Errorf("can not convert ID to int64, err: %v, inst: %#v, rid: %s", err, instance, kit.Rid)
			return kit.CCError.CCErrorf(common.CCErrCommParamsNeedInt, common.GetInstIDField(objID))
		}
		delInstIDs[index] = instID
	}

	// if any instance has been bind to a instance by the association, then these instances should not be deleted.
	err := c.asst.CheckAssociations(kit, objID, delInstIDs)
	if err != nil {
		blog.Errorf("check object(%s) asst by instID(%v) failed, err: %v, rid: %s", objID, delInstIDs, err, kit.Rid)
		return err
	}

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
	return nil
}

// DeleteInstByInstID batch delete instance by inst id
func (c *commonInst) DeleteInstByInstID(kit *rest.Kit, objectID string, instID []int64, needCheckHost bool) error {
	if len(instID) == 0 {
		blog.Errorf("inst id array is empty, rid: %s", kit.Rid)
		return nil
	}

	cond := map[string]interface{}{
		common.GetInstIDField(objectID): map[string]interface{}{common.BKDBIN: instID},
	}
	if metadata.IsCommon(objectID) {
		cond[common.BKObjIDField] = objectID
	}

	return c.DeleteInst(kit, objectID, cond, needCheckHost)
}

// FindInst search instance by condition
func (c *commonInst) FindInst(kit *rest.Kit, objID string, cond *metadata.QueryCondition) (*metadata.InstResult,
	error) {

	switch objID {
	case common.BKInnerObjIDHost:
		input := &metadata.QueryInput{
			Condition:     cond.Condition,
			Fields:        strings.Join(cond.Fields, ","),
			TimeCondition: cond.TimeCondition,
			Start:         cond.Page.Start,
			Limit:         cond.Page.Limit,
			Sort:          cond.Page.Sort,
		}
		rsp, err := c.clientSet.CoreService().Host().GetHosts(kit.Ctx, kit.Header, input)
		if err != nil {
			blog.Errorf("search object(%s) inst by the input(%#v) failed, err: %v, rid: %s", objID, input, err, kit.Rid)
			return nil, err
		}

		return &metadata.InstResult{Count: rsp.Count, Info: rsp.Info}, nil

	default:
		rsp, err := c.clientSet.CoreService().Instance().ReadInstance(kit.Ctx, kit.Header, objID, cond)
		if err != nil {
			blog.Errorf("search object(%s) inst by the cond(%#v) failed, err: %v, rid: %s", objID, cond, err, kit.Rid)
			return nil, err
		}

		return &metadata.InstResult{Count: rsp.Count, Info: rsp.Info}, nil
	}
}

// FindInstByAssociationInst deprecated function.
func (c *commonInst) FindInstByAssociationInst(kit *rest.Kit, objID string,
	asstParamCond *AssociationParams) (*metadata.InstResult, error) {

	instCond := map[string]interface{}{}
	if metadata.IsCommon(objID) {
		instCond[common.BKObjIDField] = objID
	}
	targetInstIDS := make([]int64, 0)

	for keyObjID, objs := range asstParamCond.Condition {
		// Extract the ID of the instance according to the associated object.
		cond := map[string]interface{}{}
		if common.GetObjByType(keyObjID) == common.BKInnerObjIDObject {
			cond[common.BKObjIDField] = keyObjID
		}

		for _, objCondition := range objs {
			if objCondition.Operator != common.BKDBEQ {
				if objID == keyObjID {
					if objCondition.Operator == common.BKDBLIKE ||
						objCondition.Operator == common.BKDBMULTIPLELike {
						switch t := objCondition.Value.(type) {
						case string:
							instCond[objCondition.Field] = map[string]interface{}{
								objCondition.Operator: gparams.SpecialCharChange(t),
							}
						default:
							// deal self condition
							instCond[objCondition.Field] = map[string]interface{}{
								objCondition.Operator: objCondition.Value,
							}
						}
					} else if objCondition.Operator == common.BKDBLT ||
						objCondition.Operator == common.BKDBLTE ||
						objCondition.Operator == common.BKDBGT ||
						objCondition.Operator == common.BKDBGTE {

						// fix condition covered when do date range search action.
						// ISSUE: https://github.com/TencentBlueKing/bk-cmdb/issues/5302
						if _, isExist := instCond[objCondition.Field]; !isExist {
							instCond[objCondition.Field] = make(map[string]interface{})
						}
						if condValue, ok := instCond[objCondition.Field].(map[string]interface{}); ok {
							condValue[objCondition.Operator] = objCondition.Value
						}
					} else {
						// deal self condition
						instCond[objCondition.Field] = map[string]interface{}{
							objCondition.Operator: objCondition.Value,
						}
					}
				} else {
					// deal association condition
					cond[objCondition.Field] = map[string]interface{}{
						objCondition.Operator: objCondition.Value,
					}
				}
			} else {
				if objID == keyObjID {
					// deal self condition
					switch t := objCondition.Value.(type) {
					case string:
						instCond[objCondition.Field] = map[string]interface{}{
							common.BKDBEQ: t,
						}
					default:
						instCond[objCondition.Field] = objCondition.Value
					}

				} else {
					// deal association condition
					cond[objCondition.Field] = objCondition.Value
				}
			}

		}

		if objID == keyObjID {
			// no need to search the association objects
			continue
		}

		innerCond := &metadata.QueryCondition{
			Condition: cond,
			Fields:    []string{metadata.GetInstIDFieldByObjID(keyObjID)},
		}

		instRsp, err := c.FindInst(kit, keyObjID, innerCond)
		if err != nil {
			blog.Errorf("failed to search the association inst, err: %v, rid: %s", err, kit.Rid)
			return nil, err
		}

		if len(instRsp.Info) == 0 {
			continue
		}

		asstInstIDS := make([]int64, 0)
		for _, inst := range instRsp.Info {
			id, err := inst.Int64(metadata.GetInstIDFieldByObjID(keyObjID))
			if err != nil {
				blog.Errorf("get inst id failed, err: %v, rid: %s", err, kit.Rid)
				return nil, err
			}
			asstInstIDS = append(asstInstIDS, id)
		}

		queryCond := &metadata.InstAsstQueryCondition{
			Cond: metadata.QueryCondition{Condition: map[string]interface{}{
				"bk_asst_inst_id": map[string]interface{}{
					common.BKDBIN: asstInstIDS,
				},
				"bk_asst_obj_id": keyObjID,
				"bk_obj_id":      objID,
			}},
			ObjID: objID,
		}

		rsp, err := c.clientSet.CoreService().Association().ReadInstAssociation(kit.Ctx, kit.Header, queryCond)
		if nil != err {
			blog.Errorf("search inst association failed, err: %v, rid: %s", err, kit.Rid)
			return nil, err
		}

		if len(rsp.Info) == 0 {
			continue
		}

		for _, asst := range rsp.Info {
			targetInstIDS = append(targetInstIDS, asst.InstID)
		}
	}

	if len(targetInstIDS) != 0 {
		instCond[metadata.GetInstIDFieldByObjID(objID)] = map[string]interface{}{
			common.BKDBIN: targetInstIDS,
		}
	} else if len(asstParamCond.Condition) != 0 {
		if _, ok := asstParamCond.Condition[objID]; !ok {
			return &metadata.InstResult{}, nil
		}
	}

	query := &metadata.QueryCondition{
		Condition:     instCond,
		TimeCondition: asstParamCond.TimeCondition,
		Page: metadata.BasePage{
			Limit: asstParamCond.Page.Limit,
			Sort:  asstParamCond.Page.Sort,
			Start: asstParamCond.Page.Start,
		},
	}
	if fields, ok := asstParamCond.Fields[objID]; ok {
		query.Fields = fields
	}
	return c.FindInst(kit, objID, query)
}

// UpdateInst update instance by condition
func (c *commonInst) UpdateInst(kit *rest.Kit, cond, data mapstr.MapStr, objID string) error {
	// not allowed to update these fields, need to use specialized function
	data.Remove(common.BKParentIDField)
	data.Remove(common.BKAppIDField)
	// remove unchangeable fields.
	data.Remove(metadata.GetInstIDFieldByObjID(objID))

	// generate audit log of instance.
	audit := auditlog.NewInstanceAudit(c.clientSet.CoreService())
	generateAuditParameter := auditlog.NewGenerateAuditCommonParameter(kit, metadata.AuditUpdate).WithUpdateFields(data)
	auditLog, ccErr := audit.GenerateAuditLogByCondGetData(generateAuditParameter, objID, cond)
	if ccErr != nil {
		blog.Errorf(" update inst, generate audit log failed, err: %v, rid: %s", ccErr, kit.Rid)
		return ccErr
	}

	// to update.
	inputParams := metadata.UpdateOption{
		Data:      data,
		Condition: cond,
	}
	if _, err := c.clientSet.CoreService().Instance().UpdateInstance(kit.Ctx, kit.Header, objID,
		&inputParams); err != nil {
		blog.Errorf("update the object(%s) inst by the condition(%#v) failed, err: %v, rid: %s", objID, cond,
			err, kit.Rid)
		return err
	}

	// save audit log.
	err := audit.SaveAuditLog(kit, auditLog...)
	if err != nil {
		blog.Errorf("create inst, save audit log failed, err: %v, rid: %s", err, kit.Rid)
		return kit.CCError.Error(common.CCErrAuditSaveLogFailed)
	}
	return nil
}

// SearchObjectInstances searches object instances.
func (c *commonInst) SearchObjectInstances(kit *rest.Kit, objID string, input *metadata.SearchInstanceFilter) (
	*metadata.CommonSearchResult, error) {

	// search conditions.
	cond, err := input.GetCond()
	if err != nil {
		return nil, kit.CCError.Errorf(common.CCErrCommParamsInvalid, err)
	}

	conditions := &metadata.QueryCondition{
		Fields:         input.Fields,
		Condition:      cond,
		TimeCondition:  input.TimeCondition,
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
	input *metadata.CountInstanceFilter) (*metadata.CommonCountResult, error) {

	// count conditions.
	cond, err := input.GetCond()
	if err != nil {
		return nil, kit.CCError.Errorf(common.CCErrCommParamsInvalid, err)
	}
	conditions := &metadata.Condition{
		Condition:     cond,
		TimeCondition: input.TimeCondition,
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
func (c *commonInst) FindInstChildTopo(kit *rest.Kit, objID string, instID int64) (
	int, []*metadata.CommonInstTopo, error) {

	return c.findInstTopo(kit, objID, instID, true)
}

// findInstParentTopo find instance's parent topo
func (c *commonInst) findInstParentTopo(kit *rest.Kit, objID string, instID int64) (
	int, []*metadata.CommonInstTopo, error) {

	return c.findInstTopo(kit, objID, instID, false)
}

func (c *commonInst) findInstTopo(kit *rest.Kit, objID string, instID int64, needChild bool) (int,
	[]*metadata.CommonInstTopo, error) {

	instIDField := metadata.GetInstIDFieldByObjID(objID)
	if instID == 0 {
		blog.Errorf("inst id is 0, rid:%s", kit.Rid)
		return 0, nil, kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, instIDField)
	}

	tableName := common.GetInstTableName(objID, kit.TenantID)
	filter := []map[string]interface{}{{metadata.GetInstIDFieldByObjID(objID): instID}}
	cnt, ccErr := c.clientSet.CoreService().Count().GetCountByFilter(kit.Ctx, kit.Header, tableName, filter)
	if ccErr != nil {
		blog.Errorf("failed to check the inst, err: %v, rid: %s", ccErr, kit.Rid)
		return 0, nil, ccErr
	}

	if cnt[0] == 0 {
		blog.Errorf("inst of inst id(%d) is non-exist, rid: %s", instID, kit.Rid)
		return 0, nil, kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, instIDField)
	}

	tmpResults := map[string]*metadata.CommonInstTopo{}

	topoInsts, relation, err := c.getAssociatedObjectWithInsts(kit, objID, instID, needChild)
	if err != nil {
		return 0, nil, err
	}

	for _, topoInst := range topoInsts {
		object := topoInst.Object
		commonInst, exists := tmpResults[object.ObjectID]
		if !exists {
			commonInst = &metadata.CommonInstTopo{
				Children: []metadata.InstNameAsst{},
			}
			commonInst.ObjectName = object.ObjectName
			commonInst.ObjIcon = object.ObjIcon
			commonInst.ObjID = object.ObjectID
			tmpResults[object.ObjectID] = commonInst
		}

		commonInst.Count = commonInst.Count + len(topoInst.Insts)

		for _, inst := range topoInst.Insts {

			id, err := inst.Int64(metadata.GetInstIDFieldByObjID(object.ObjectID))
			if err != nil {
				return 0, nil, err
			}

			name, err := inst.String(metadata.GetInstNameFieldName(object.ObjectID))
			if err != nil {
				return 0, nil, err
			}

			instAsst := metadata.InstNameAsst{
				ID:         strconv.Itoa(int(id)),
				InstID:     id,
				InstName:   name,
				ObjectName: object.ObjectName,
				ObjIcon:    object.ObjIcon,
				ObjID:      object.ObjectID,
				AssoID:     relation[id],
			}

			tmpResults[object.ObjectID].Children = append(tmpResults[object.ObjectID].Children, instAsst)
		}
	}

	results := make([]*metadata.CommonInstTopo, 0)
	for _, subResult := range tmpResults {
		results = append(results, subResult)
	}

	return len(results), results, nil
}

// FindInstTopo find instance all topo which include it's child and parent
func (c *commonInst) FindInstTopo(kit *rest.Kit, obj metadata.Object, instID int64) (int, []metadata.CommonInstTopoV2,
	error) {

	instIDField := metadata.GetInstIDFieldByObjID(obj.ObjectID)
	instNameField := metadata.GetInstNameFieldName(obj.ObjectID)
	if instID == 0 {
		blog.Errorf("inst id is 0, rid:%s", kit.Rid)
		return 0, nil, kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, instIDField)
	}

	query := &metadata.QueryCondition{
		Condition: map[string]interface{}{instIDField: instID},
		Fields:    []string{instIDField, instNameField},
	}
	inst, err := c.FindInst(kit, obj.ObjectID, query)
	if err != nil {
		blog.Errorf("failed to find the inst, err: %v, rid: %s", err, kit.Rid)
		return 0, nil, err
	}

	if len(inst.Info) == 0 {
		blog.Errorf("inst of inst id(%d) is non-exist, rid: %s", instID, kit.Rid)
		return 0, nil, kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, instIDField)
	}

	results := make([]metadata.CommonInstTopoV2, 0)
	id, err := inst.Info[0].Int64(instIDField)
	if err != nil {
		blog.Errorf("failed to find the inst id, err: %v, rid: %s", err, kit.Rid)
		return 0, nil, err
	}

	name, err := inst.Info[0].String(instNameField)
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

	_, parentInsts, err := c.findInstParentTopo(kit, obj.ObjectID, id)
	if err != nil {
		blog.Errorf("failed to find the inst, err: %v rid: %s", err, kit.Rid)
		return 0, nil, err
	}

	_, childInsts, err := c.FindInstChildTopo(kit, obj.ObjectID, id)
	if err != nil {
		blog.Errorf("failed to find the inst, err: %v, rid: %s", err, kit.Rid)
		return 0, nil, err
	}

	results = append(results, metadata.CommonInstTopoV2{
		Prev: parentInsts,
		Next: childInsts,
		Curr: commonInst,
	})

	return len(results), results, nil
}

func (c *commonInst) validMainLineParentID(kit *rest.Kit, objID string, data mapstr.MapStr) error {
	if objID == common.BKInnerObjIDApp {
		return nil
	}

	def, exist := data.Get(common.BKDefaultField)
	if exist {
		defaultInt, err := util.GetIntByInterface(def)
		if err != nil {
			blog.Errorf("failed to parse the default value, err: %v, rid: %s", err, kit.Rid)
			return err
		}

		if defaultInt != common.DefaultFlagDefaultValue {
			return nil
		}
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

	cond := &metadata.Condition{
		Condition: map[string]interface{}{metadata.GetInstIDFieldByObjID(objID): instID},
	}
	if bizID != 0 {
		cond.Condition[common.BKAppIDField] = bizID
	}

	if metadata.IsCommon(objID) {
		cond.Condition[common.BKObjIDField] = objID
	}

	rsp, err := c.clientSet.CoreService().Instance().CountInstances(kit.Ctx, kit.Header, objID, cond)
	if err != nil {
		blog.Errorf("count object(%s) inst by the condition(%#v), err: %v, rid: %s", objID, cond, err, kit.Rid)
		return err
	}

	if rsp.Count == 0 {
		return kit.CCError.Error(common.CCErrTopoInstSelectFailed)
	}

	return nil
}

func (c *commonInst) validObject(kit *rest.Kit, objID string, data mapstr.MapStr) error {

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
		common.BKObjIDField:           objID,
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

	if err := c.validMainLineParentID(kit, asst.Info[0].AsstObjID, data); err != nil {
		blog.Errorf("valid mainline object(%s) parentID failed, err: %v, rid: %s", objID, err, kit.Rid)
		return err
	}

	return nil
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
		query := &metadata.QueryCondition{
			Condition: map[string]interface{}{common.BKSetIDField: map[string]interface{}{common.BKDBIN: instIDs}},
			Fields:    []string{common.BKModuleIDField},
			Page:      metadata.BasePage{Limit: common.BKNoLimit},
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
				blog.Errorf("can not convert ID to int64, err: %v, module: %#v, rid: %s", err, module, kit.Rid)
				return nil, false, kit.CCError.CCErrorf(common.CCErrCommParamsNeedInt, common.BKModuleIDField)
			}
			moduleIDs[index] = moduleID
		}
	} else {
		var err error
		moduleIDs, err = c.mainlineHasHost(kit, objID, objInstMap, instIDs)
		if err != nil {
			return nil, false, err
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

func (c *commonInst) mainlineHasHost(kit *rest.Kit, objID string, objInstMap map[string][]mapstr.MapStr,
	instIDs []int64) ([]int64, error) {

	// get mainline object relation(excluding hosts) by mainline associations
	mainlineCond := &metadata.QueryCondition{
		Condition: map[string]interface{}{
			common.AssociationKindIDField: common.AssociationKindMainline,
			common.BKObjIDField: mapstr.MapStr{
				common.BKDBNE: common.BKInnerObjIDHost,
			},
		},
		Fields:         []string{common.BKObjIDField, common.BKAsstObjIDField},
		DisableCounter: true,
	}
	asstRsp, err := c.clientSet.CoreService().Association().ReadModelAssociation(kit.Ctx, kit.Header, mainlineCond)
	if err != nil {
		blog.Errorf("search mainline association failed, error: %v, rid: %s", err, kit.Rid)
		return nil, err
	}

	objChildMap := make(map[string]string)
	isMainline := false
	for _, asst := range asstRsp.Info {
		objChildMap[asst.AsstObjID] = asst.ObjectID
		if asst.AsstObjID == objID || asst.ObjectID == objID {
			isMainline = true
		}
	}

	if !isMainline {
		return nil, nil
	}

	// loop through the child topology level to get all instances
	var moduleIDs []int64
	parentIDs := instIDs
	for childObjID := objChildMap[objID]; len(childObjID) != 0; childObjID = objChildMap[childObjID] {
		cond := map[string]interface{}{common.BKParentIDField: map[string]interface{}{common.BKDBIN: parentIDs}}
		if metadata.IsCommon(childObjID) {
			cond[metadata.ModelFieldObjectID] = childObjID
		}

		if childObjID == common.BKInnerObjIDSet {
			cond[common.BKDefaultField] = common.DefaultFlagDefaultValue
		}

		query := &metadata.QueryCondition{
			Condition: cond,
			Page:      metadata.BasePage{Limit: common.BKNoLimit},
		}

		childRsp, err := c.FindInst(kit, childObjID, query)
		if err != nil {
			blog.Errorf("find children failed, err: %v, parent IDs: %+v, rid: %s", err, parentIDs, kit.Rid)
			return nil, err
		}

		if len(childRsp.Info) == 0 {
			return nil, nil
		}

		parentIDs = make([]int64, len(childRsp.Info))
		for index, instance := range childRsp.Info {
			instID, err := instance.Int64(common.GetInstIDField(childObjID))
			if err != nil {
				blog.Errorf("can not convert ID to int64, err: %v, inst: %#v, rid: %s", err, instance, kit.Rid)
				return nil, kit.CCError.CCErrorf(common.CCErrCommParamsNeedInt, common.GetInstIDField(childObjID))
			}
			parentIDs[index] = instID
		}

		if childObjID == common.BKInnerObjIDModule {
			moduleIDs = parentIDs
		}

		objInstMap[childObjID] = childRsp.Info
	}
	return moduleIDs, nil
}

func (c *commonInst) innerHasHost(kit *rest.Kit, moduleIDs []int64) (bool, error) {
	if len(moduleIDs) == 0 {
		blog.Errorf("module id array is empty, rid: %s", kit.Rid)
		return false, nil
	}
	filter := []map[string]interface{}{{common.BKModuleIDField: mapstr.MapStr{common.BKDBIN: moduleIDs}}}
	rsp, err := c.clientSet.CoreService().Count().GetCountByFilter(kit.Ctx, kit.Header,
		common.BKTableNameModuleHostConfig, filter)
	if err != nil {
		blog.Errorf("searh host object relation failed, err: %v, rid: %s", err, kit.Rid)
		return false, err
	}

	return rsp[0] != 0, nil
}

// getAssociatedObjectWithInsts TODO
// GetObjectWithInsts get object with insts, get parent or child depends on needChild
func (c *commonInst) getAssociatedObjectWithInsts(kit *rest.Kit, objID string, instID int64, needChild bool) (
	[]*ObjectWithInsts, map[int64]int64, error) {

	cond := mapstr.New()
	if needChild {
		cond.Set(common.BKObjIDField, objID)
	} else {
		cond.Set(common.BKAsstObjIDField, objID)
	}

	objPairs, err := c.searchAssoObjects(kit, needChild, cond)
	if err != nil {
		blog.Errorf("failed to get the object(%s)'s parent, err: %v, rid: %s", objID, err, kit.Rid)
		return nil, nil, err
	}

	relation := make(map[int64]int64)
	result := make([]*ObjectWithInsts, 0)
	for _, objPair := range objPairs {

		queryCond := &metadata.InstAsstQueryCondition{
			Cond: metadata.QueryCondition{Condition: mapstr.MapStr{
				common.AssociationObjAsstIDField: objPair.AssocName,
			}},
			ObjID: objPair.Object.ObjectID,
		}

		if needChild {
			queryCond.Cond.Condition.Set(common.BKInstIDField, instID)
			queryCond.Cond.Condition.Set(common.BKObjIDField, objID)
			queryCond.Cond.Condition.Set(common.BKAsstObjIDField, objPair.Object.ObjectID)
		} else {
			queryCond.Cond.Condition.Set(common.BKAsstInstIDField, instID)
			queryCond.Cond.Condition.Set(common.BKObjIDField, objPair.Object.ObjectID)
			queryCond.Cond.Condition.Set(common.BKAsstObjIDField, objID)
		}

		rsp, err := c.clientSet.CoreService().Association().ReadInstAssociation(kit.Ctx, kit.Header, queryCond)
		if err != nil {
			blog.Errorf("search inst association failed , err: %v, rid: %s", err, kit.Rid)
			return nil, nil, err
		}

		// found no noe inst association with this object and association info.
		// which means that, this object association has not been instantiated.
		if len(rsp.Info) == 0 {
			continue
		}

		instIDs := make([]int64, 0)
		for _, item := range rsp.Info {
			var instID int64
			if needChild {
				instID = item.AsstInstID
			} else {
				instID = item.InstID
			}
			relation[instID] = item.ID
			instIDs = append(instIDs, instID)
		}

		innerCond := &metadata.QueryCondition{
			Condition: mapstr.MapStr{objPair.Object.GetInstIDFieldName(): mapstr.MapStr{common.BKDBIN: instIDs}},
			Fields: []string{common.GetInstIDField(objPair.Object.ObjectID),
				common.GetInstNameField(objPair.Object.ObjectID)},
		}
		if objPair.Object.IsCommon() {
			innerCond.Condition[common.BKObjIDField] = objPair.Object.ObjectID
		}

		rspItems, err := c.FindInst(kit, objPair.Object.ObjectID, innerCond)
		if err != nil {
			blog.Errorf("failed to search the insts by the condition(%#v), err: %v, rid: %s", innerCond, err, kit.Rid)
			return result, nil, err
		}

		rstObj := &ObjectWithInsts{Object: objPair.Object, Insts: rspItems.Info}
		result = append(result, rstObj)
	}

	return result, relation, nil
}

func (c *commonInst) searchAssoObjects(kit *rest.Kit, needChild bool, cond mapstr.MapStr) ([]ObjectAssoPair,
	error) {

	input := &metadata.QueryCondition{
		Condition:      cond,
		Fields:         []string{common.BKAsstObjIDField, common.BKObjIDField, common.AssociationObjAsstIDField},
		DisableCounter: true,
	}
	rsp, err := c.clientSet.CoreService().Association().ReadModelAssociation(kit.Ctx, kit.Header, input)
	if err != nil {
		blog.Errorf("search object association failed, err: %v, rid: %s", err, kit.Rid)
		return nil, err
	}

	if len(rsp.Info) == 0 {
		blog.Errorf("search object association return empty, rid: %s", kit.Rid)
		return make([]ObjectAssoPair, 0), nil
	}

	objAssoMap := make(map[string]metadata.Association, 0)
	var objIDArray []string
	for _, asst := range rsp.Info {
		if needChild {
			objIDArray = append(objIDArray, asst.AsstObjID)
			objAssoMap[asst.AsstObjID] = asst
		} else {
			objIDArray = append(objIDArray, asst.ObjectID)
			objAssoMap[asst.ObjectID] = asst
		}
	}

	queryCond := &metadata.QueryCondition{
		Condition:      mapstr.MapStr{metadata.ModelFieldObjectID: mapstr.MapStr{common.BKDBIN: objIDArray}},
		Fields:         []string{common.BKObjNameField, common.BKObjIDField, common.BKObjIconField},
		DisableCounter: true,
	}
	rspRst, err := c.clientSet.CoreService().Model().ReadModel(kit.Ctx, kit.Header, queryCond)
	if err != nil {
		blog.Errorf("failed to search the object by cond(%#v), err: %v, rid: %s", queryCond, err, kit.Rid)
		return nil, err
	}

	if len(rspRst.Info) == 0 {
		return make([]ObjectAssoPair, 0), nil
	}

	pair := make([]ObjectAssoPair, 0)
	for _, object := range rspRst.Info {
		pair = append(pair, ObjectAssoPair{
			Object:    object,
			AssocName: objAssoMap[object.ObjectID].AssociationName,
		})
	}

	return pair, nil
}
