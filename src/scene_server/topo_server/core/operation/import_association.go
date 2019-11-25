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

package operation

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"configcenter/src/auth/extensions"
	"configcenter/src/auth/meta"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/topo_server/core/types"
)

func (assoc *association) ImportInstAssociation(ctx context.Context, params types.ContextParams, objID string, importData map[int]metadata.ExcelAssocation) (resp metadata.ResponeImportAssociationData, err error) {

	ia := NewImportAssociation(ctx, assoc, params, objID, importData, assoc.authManager)
	err = ia.ParsePrimaryKey()
	if err != nil {
		return resp, err
	}

	errIdxMsgMap := ia.ImportAssociation()
	if len(errIdxMsgMap) > 0 {
		err = params.Err.Error(common.CCErrorTopoImportAssociation)
	}
	for row, msg := range errIdxMsgMap {
		resp.ErrMsgMap = append(resp.ErrMsgMap, metadata.RowMsgData{
			Row: row,
			Msg: msg,
		})
	}

	return resp, err
}

type importAssociationInst struct {
	instID int64
	//strings.Joion([]string{property name, property value}, "=")|true
	attrNameVal map[string]bool
}
type importAssociation struct {
	objID      string
	cli        *association
	ctx        context.Context
	importData map[int]metadata.ExcelAssocation
	params     types.ContextParams

	// map[AssociationName]Association alias  map[association flag]Association
	asstIDInfoMap map[string]*metadata.Association
	// asst obj info  map[objID]map[property name] attribute
	asstObjIDProperty map[string]map[string]metadata.Attribute

	parseImportDataErr map[int]string
	//map[objID][]condition.Condition
	queryInstConds map[string][]mapstr.MapStr

	// map[objID][instcnade id]strings.Joion([]string{property name, property value}, "=")[]importAssociationInst
	instIDAttrKeyValMap map[string]map[string][]*importAssociationInst
	//http header http request id
	rid string

	authManager *extensions.AuthManager
}

type importAssociationInterface interface {
	ParsePrimaryKey() error
	ImportAssociation() map[int]string
}

func NewImportAssociation(ctx context.Context, cli *association, params types.ContextParams, objID string, importData map[int]metadata.ExcelAssocation, authManager *extensions.AuthManager) importAssociationInterface {

	return &importAssociation{
		objID:      objID,
		cli:        cli,
		ctx:        ctx,
		importData: importData,
		params:     params,

		asstIDInfoMap:       make(map[string]*metadata.Association, 0),
		asstObjIDProperty:   make(map[string]map[string]metadata.Attribute, 0),
		parseImportDataErr:  make(map[int]string),
		queryInstConds:      make(map[string][]mapstr.MapStr),
		instIDAttrKeyValMap: make(map[string]map[string][]*importAssociationInst),

		rid: util.GetHTTPCCRequestID(params.Header),

		authManager: authManager,
	}
}

func (ia *importAssociation) ImportAssociation() map[int]string {
	ia.importAssociation()

	return ia.parseImportDataErr
}

func (ia *importAssociation) ParsePrimaryKey() error {
	err := ia.getAssociationInfo()
	if err != nil {
		return err
	}

	err = ia.getAssociationObjProperty()
	if err != nil {
		return err
	}

	ia.parseImportDataPrimary()
	err = ia.getInstDataByConds()
	if err != nil {
		return err
	}

	return nil

}

func (ia *importAssociation) importAssociation() {
	for idx, asstInfo := range ia.importData {
		_, ok := ia.parseImportDataErr[idx]
		if ok {
			continue
		}
		asstID, ok := ia.asstIDInfoMap[asstInfo.ObjectAsstID]
		if !ok {
			ia.parseImportDataErr[idx] = ia.params.Lang.Languagef("import_association_id_not_found", asstInfo.ObjectAsstID)
			continue
		}
		srcInstID, err := ia.getInstIDByPrimaryKey(ia.objID, asstInfo.SrcPrimary)
		if err != nil {
			ia.parseImportDataErr[idx] = err.Error()
			continue
		}
		dstInstID, err := ia.getInstIDByPrimaryKey(asstID.AsstObjID, asstInfo.DstPrimary)
		if err != nil {
			ia.parseImportDataErr[idx] = err.Error()
			continue
		}
		err = ia.authManager.AuthorizeByInstanceID(ia.ctx, ia.params.Header, meta.Update, ia.objID, srcInstID)
		if err != nil {
			ia.parseImportDataErr[idx] = err.Error()
			continue
		}
		err = ia.authManager.AuthorizeByInstanceID(ia.ctx, ia.params.Header, meta.Update, asstID.AsstObjID, dstInstID)
		if err != nil {
			ia.parseImportDataErr[idx] = err.Error()
			continue
		}
		switch asstInfo.Operate {
		case metadata.ExcelAssocationOperateAdd:

			conds := condition.CreateCondition()
			conds.Field(common.AssociationObjAsstIDField).Eq(asstInfo.ObjectAsstID)
			conds.Field(common.BKObjIDField).Eq(ia.objID)
			conds.Field(common.BKInstIDField).Eq(srcInstID)
			conds.Field(common.AssociatedObjectIDField).Eq(asstID.AsstObjID)
			isExist, err := ia.isExistInstAsst(idx, conds, dstInstID, asstID.Mapping)
			if err != nil {
				ia.parseImportDataErr[idx] = err.Error()
				continue
			}
			if isExist {
				continue
			}

			ia.addSrcAssociation(idx, asstID.AssociationName, srcInstID, dstInstID)
		case metadata.ExcelAssocationOperateDelete:
			conds := condition.CreateCondition()
			conds.Field(common.AssociationObjAsstIDField).Eq(asstInfo.ObjectAsstID)
			conds.Field(common.BKObjIDField).Eq(ia.objID)
			conds.Field(common.BKInstIDField).Eq(srcInstID)
			conds.Field(common.AssociatedObjectIDField).Eq(asstID.AsstObjID)
			conds.Field(common.BKAsstInstIDField).Eq(dstInstID)
			ia.delSrcAssociation(idx, conds)
		default:
			ia.parseImportDataErr[idx] = ia.params.Lang.Language("import_association_operate_not_found")
		}

	}
}

func (ia *importAssociation) getAssociationInfo() error {
	var associationFlag []string
	for _, info := range ia.importData {
		associationFlag = append(associationFlag, info.ObjectAsstID)
	}

	cond := condition.CreateCondition()
	cond.Field(common.AssociationObjAsstIDField).In(associationFlag)
	cond.Field(common.BKObjIDField).Eq(ia.objID)
	queryInput := &metadata.QueryCondition{Condition: cond.ToMapStr()}

	rsp, err := ia.cli.clientSet.CoreService().Association().ReadModelAssociation(ia.ctx, ia.params.Header, queryInput)
	if nil != err {
		blog.Errorf("[getAssociationInfo] failed to request the object controller , error info is %s, input:%+v, rid:%s", err.Error(), queryInput, ia.rid)
		return ia.params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if !rsp.Result {
		blog.Errorf("[getAssociationInfo] failed to search the inst association, error info is %s, input:%+v, rid:%s", rsp.ErrMsg, queryInput, ia.rid)
		return ia.params.Err.New(rsp.Code, rsp.ErrMsg)
	}

	for index := range rsp.Data.Info {
		ia.asstIDInfoMap[rsp.Data.Info[index].AssociationName] = &rsp.Data.Info[index]
	}

	return nil
}

func (ia *importAssociation) getAssociationObjProperty() error {
	var objIDArr []string
	for _, info := range ia.asstIDInfoMap {
		objIDArr = append(objIDArr, info.AsstObjID)
	}
	objIDArr = append(objIDArr, ia.objID)

	uniqueCond := condition.CreateCondition()
	uniqueCond.Field(common.BKObjIDField).In(objIDArr)

	uniqueQueryCond := metadata.QueryCondition{Condition: uniqueCond.ToMapStr()}
	uniqueResult, err := ia.cli.clientSet.CoreService().Model().ReadModelAttrUnique(ia.ctx, ia.params.Header, uniqueQueryCond)
	if nil != err {
		blog.ErrorJSON("[getAssociationInfo] http do error.  search model unique , error info is %s, input:%s, rid:%s", err.Error(), uniqueQueryCond, ia.rid)
		return ia.params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if nil != err {
		blog.ErrorJSON("[getAssociationInfo]http reply error. search model unique , error info is %s, input:%s, rid:%s", err.Error(), uniqueQueryCond, ia.rid)
		return ia.params.Err.New(uniqueResult.Code, uniqueResult.ErrMsg)
	}
	var propertyIDArr []uint64
	for _, unique := range uniqueResult.Data.Info {
		if !unique.MustCheck {
			continue
		}
		for _, property := range unique.Keys {
			propertyIDArr = append(propertyIDArr, property.ID)
		}
	}

	cond := condition.CreateCondition()
	cond.Field(common.BKObjIDField).In(objIDArr)
	cond.Field(common.BKFieldID).In(propertyIDArr)

	attrCond := &metadata.QueryCondition{Condition: cond.ToMapStr()}
	rsp, err := ia.cli.clientSet.CoreService().Model().ReadModelAttrByCondition(ia.ctx, ia.params.Header, attrCond)
	if nil != err {
		blog.Errorf("[getAssociationInfo] failed to  search attribute , error info is %s, input:%+v, rid:%s", err.Error(), attrCond, ia.rid)
		return ia.params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if !rsp.Result {
		blog.Errorf("[getAssociationInfo] failed to search attribute, error code:%s, error messge: %s, input:%+v, rid:%s", rsp.Code, rsp.ErrMsg, cond, ia.rid)
		return ia.params.Err.New(rsp.Code, rsp.ErrMsg)
	}

	for _, attr := range rsp.Data.Info {
		_, ok := ia.asstObjIDProperty[attr.ObjectID]
		if !ok {
			ia.asstObjIDProperty[attr.ObjectID] = make(map[string]metadata.Attribute)
		}
		ia.asstObjIDProperty[attr.ObjectID][attr.PropertyName] = attr
	}

	return nil

}

func (ia *importAssociation) parseImportDataPrimary() {

	for idx, info := range ia.importData {

		associationInst, ok := ia.asstIDInfoMap[info.ObjectAsstID]
		if !ok {
			ia.parseImportDataErr[idx] = ia.params.Lang.Languagef("import_asstid_not_foud", info.ObjectAsstID)
			continue
		}
		srcCond, err := ia.parseImportDataPrimaryItem(ia.objID, info.SrcPrimary)
		if err != nil {
			ia.parseImportDataErr[idx] = err.Error()
		} else {
			_, ok = ia.queryInstConds[ia.objID]
			if !ok {
				ia.queryInstConds[ia.objID] = make([]mapstr.MapStr, 0)
			}
			ia.queryInstConds[ia.objID] = append(ia.queryInstConds[ia.objID], srcCond)

		}
		dstCond, err := ia.parseImportDataPrimaryItem(associationInst.AsstObjID, info.DstPrimary)
		if err != nil {
			ia.parseImportDataErr[idx] = err.Error()
		} else {
			_, ok = ia.queryInstConds[associationInst.AsstObjID]
			if !ok {
				ia.queryInstConds[associationInst.AsstObjID] = make([]mapstr.MapStr, 0)
			}
			ia.queryInstConds[associationInst.AsstObjID] = append(ia.queryInstConds[associationInst.AsstObjID], dstCond)

		}
	}

	return

}

func (ia *importAssociation) parseImportDataPrimaryItem(objID string, item string) (mapstr.MapStr, error) {

	keyValMap := mapstr.New()
	primaryArr := strings.Split(item, common.ExcelAsstPrimaryKeySplitChar)

	for _, primary := range primaryArr {

		primary = strings.TrimSpace(primary)
		keyValArr := strings.Split(primary, common.ExcelAsstPrimaryKeyJoinChar)
		if len(keyValArr) != 2 {
			blog.ErrorJSON("parseImportDataPrimaryItem eror. primary:%s, rid:%s", primary, ia.rid)
			return nil, fmt.Errorf(ia.params.Lang.Languagef("import_asst_obj_property_str_primary_format_error", objID, item))
		}
		attr, ok := ia.asstObjIDProperty[objID][keyValArr[0]]
		if !ok {
			return nil, fmt.Errorf(ia.params.Lang.Languagef("import_asst_obj_primary_property_str_not_found", objID, keyValArr[0]))
		}
		realVal, err := convStrToCCType(keyValArr[1], attr)
		if err != nil {
			return nil, fmt.Errorf(ia.params.Lang.Languagef("import_asst_obj_property_str_primary_type_error", objID, keyValArr[0]))
		}

		keyValMap[attr.PropertyID] = realVal
	}
	if len(keyValMap) != len(ia.asstObjIDProperty[objID]) {
		blog.ErrorJSON("parseImportDataPrimaryItem eror. keyVal:%s, objID:%s, objIDProperty:%s,rid:%s", keyValMap, objID, ia.asstObjIDProperty[objID], ia.rid)
		return nil, fmt.Errorf(ia.params.Lang.Languagef("import_asst_obj_property_str_primary_count_len", objID, item))
	}

	return keyValMap, nil
}

func (ia *importAssociation) getInstDataByConds() error {

	for objID, valArr := range ia.queryInstConds {

		instIDKey := metadata.GetInstIDFieldByObjID(objID)
		conds := condition.CreateCondition()
		conds.NewOR().MapStrArr(valArr)

		instArr, err := ia.getInstDataByObjIDConds(objID, instIDKey, conds)
		if err != nil {
			return err
		}
		for _, inst := range instArr {
			ia.parseInstToImportAssociationInst(objID, instIDKey, inst)
		}
	}

	return nil
}

func (ia *importAssociation) getInstDataByObjIDConds(objID, instIDKey string, conds condition.Condition) ([]mapstr.MapStr, error) {

	var fields []string
	for _, attr := range ia.asstObjIDProperty[objID] {
		fields = append(fields, attr.PropertyID)
	}

	fields = append(fields, instIDKey)
	queryInput := &metadata.QueryCondition{}
	queryInput.Condition = conds.ToMapStr()
	queryInput.Fields = fields

	instSearchResult, err := ia.cli.clientSet.CoreService().Instance().ReadInstance(ia.ctx, ia.params.Header, objID, queryInput)
	if err != nil {
		blog.Errorf("[getInstDataByObjIDConds] failed to  search %s instance , error info is %s, input:%#v, rid:%s", objID, err.Error(), queryInput, ia.rid)
		return nil, ia.params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !instSearchResult.Result {
		blog.Errorf("[getInstDataByObjIDConds] failed to search %s instance,  error code:%d, error message: %s, input:%+v, rid:%s", objID, instSearchResult.Code, instSearchResult.ErrMsg, queryInput, ia.rid)
		return nil, ia.params.Err.New(instSearchResult.Code, instSearchResult.ErrMsg)
	}
	return instSearchResult.Data.Info, nil
}

func (ia *importAssociation) parseInstToImportAssociationInst(objID, instIDKey string, inst mapstr.MapStr) {
	instID, err := inst.Int64(instIDKey)
	//inst info can not found
	if err != nil {
		blog.Warnf("parseInstToImportAssociationInst get %s field from %s model error,error:%s, rid:%d ", instID, objID, err.Error(), ia.rid)
		return
	}

	attrNameValMap := importAssociationInst{
		instID:      instID,
		attrNameVal: make(map[string]bool),
	}
	isErr := false
	for _, attr := range ia.asstObjIDProperty[objID] {
		val, err := inst.String(attr.PropertyID)
		//inst info can not found
		if err != nil {
			isErr = true
			blog.Warnf("parseInstToImportAssociationInst get %s field from %s model error,error:%s, rid:%d ", attr.PropertyID, objID, err.Error(), ia.rid)
			continue
		}
		attrNameValMap.attrNameVal[buildPrimaryStr(attr.PropertyName, val)] = true
	}
	if isErr {
		return
	}
	for key := range attrNameValMap.attrNameVal {
		_, ok := ia.instIDAttrKeyValMap[objID]
		if !ok {
			ia.instIDAttrKeyValMap[objID] = make(map[string][]*importAssociationInst)
		}
		_, ok = ia.instIDAttrKeyValMap[objID][key]
		if !ok {
			ia.instIDAttrKeyValMap[objID][key] = make([]*importAssociationInst, 0)
		}
		ia.instIDAttrKeyValMap[objID][key] = append(ia.instIDAttrKeyValMap[objID][key], &attrNameValMap)
	}
}

func (ia *importAssociation) delSrcAssociation(idx int, cond condition.Condition) {
	_, ok := ia.parseImportDataErr[idx]
	if ok {
		return
	}

	result, err := ia.cli.clientSet.CoreService().Association().DeleteInstAssociation(ia.ctx, ia.params.Header, &metadata.DeleteOption{Condition: cond.ToMapStr()})
	if err != nil {
		ia.parseImportDataErr[idx] = err.Error()
		return
	}

	if !result.Result {
		ia.parseImportDataErr[idx] = result.ErrMsg
		return
	}

}

func (ia *importAssociation) addSrcAssociation(idx int, asstFlag string, instID, assInstID int64) {
	_, ok := ia.parseImportDataErr[idx]
	if ok {
		return
	}

	asstInfo := ia.asstIDInfoMap[asstFlag]

	inst := metadata.CreateOneInstanceAssociation{}
	inst.Data.ObjectAsstID = asstFlag
	inst.Data.InstID = instID
	inst.Data.ObjectID = ia.objID
	inst.Data.AsstObjectID = asstInfo.AsstObjID
	inst.Data.AsstInstID = assInstID
	inst.Data.AssociationKindID = asstInfo.AsstKindID
	rsp, err := ia.cli.clientSet.CoreService().Association().CreateInstAssociation(ia.ctx, ia.params.Header, &inst)
	if err != nil {
		ia.parseImportDataErr[idx] = err.Error()
	}
	if !rsp.Result {
		ia.parseImportDataErr[idx] = rsp.ErrMsg
	}
}

func (ia *importAssociation) isExistInstAsst(idx int, cond condition.Condition, dstInstID int64, asstMapping metadata.AssociationMapping) (isExit bool, err error) {
	_, ok := ia.parseImportDataErr[idx]
	if ok {
		return
	}
	if asstMapping != metadata.OneToOneMapping {
		cond.Field(common.BKAsstInstIDField).Eq(dstInstID)
	}
	rsp, err := ia.cli.clientSet.CoreService().Association().ReadInstAssociation(ia.ctx, ia.params.Header, &metadata.QueryCondition{Condition: cond.ToMapStr()})
	if err != nil {
		return false, err
	}
	if !rsp.Result {
		ia.parseImportDataErr[idx] = rsp.ErrMsg
		return false, ia.params.Err.New(rsp.Code, rsp.ErrMsg)
	}

	if len(rsp.Data.Info) == 0 {
		return false, nil
	}
	if rsp.Data.Info[0].AsstInstID != dstInstID &&
		asstMapping == metadata.OneToOneMapping {
		return false, ia.params.Err.Errorf(common.CCErrCommDuplicateItem, "association")
	}

	return true, nil
}

func (ia *importAssociation) getInstIDByPrimaryKey(objID, primary string) (int64, error) {
	primaryArr := strings.Split(primary, common.ExcelAsstPrimaryKeySplitChar)
	if len(primaryArr) == 0 {
		return 0, fmt.Errorf(ia.params.Lang.Languagef("import_instance_not_foud", objID, primary))
	}

	instArr, ok := ia.instIDAttrKeyValMap[objID][primaryArr[0]]
	if !ok {
		return 0, fmt.Errorf(ia.params.Lang.Languagef("import_instance_not_foud", objID, primaryArr[0]))
	}

	for _, inst := range instArr {

		isEq := true
		for _, item := range primaryArr {
			if _, ok := inst.attrNameVal[item]; !ok {
				isEq = false
				break
			}
		}
		if isEq {
			return inst.instID, nil
		}

	}

	return 0, fmt.Errorf(ia.params.Lang.Languagef("import_instance_not_foud", objID, primary))

}

func buildPrimaryStr(name, val string) string {
	return name + common.ExcelAsstPrimaryKeyJoinChar + val
}

func convStrToCCType(val string, attr metadata.Attribute) (interface{}, error) {
	switch attr.PropertyType {
	case common.FieldTypeBool:

		return strconv.ParseBool(val)
	case common.FieldTypeEnum:
		option, optionOk := attr.Option.([]interface{})
		if !optionOk {
			return nil, fmt.Errorf("not foud")
		}
		return getEnumIDByName(val, option), nil
	case common.FieldTypeInt:
		return util.GetInt64ByInterface(val)
	case common.FieldTypeFloat:
		return util.GetFloat64ByInterface(val)

	default:
		return val, nil
	}
}

// getEnumIDByName get enum name from option
func getEnumIDByName(name string, items []interface{}) string {
	id := name
	for _, valRow := range items {
		mapVal, ok := valRow.(map[string]interface{})
		if ok {
			enumName, ok := mapVal["name"].(string)
			if true == ok {
				if enumName == name {
					id = mapVal["id"].(string)
				}
			}
		}
	}

	return id
}
