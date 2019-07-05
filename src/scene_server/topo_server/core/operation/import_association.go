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
	"sort"
	"strconv"
	"strings"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/topo_server/core/types"
)

func (a *association) ImportInstAssociation(ctx context.Context, params types.ContextParams, objID string, importData map[int]metadata.ExcelAssocation) (resp metadata.ResponeImportAssociationData, err error) {

	ia := NewImportAssociation(ctx, a, params, objID, importData)
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

	// map[objID][kesname1,keyname2]
	objUniques map[string]map[string]bool

	//http header http request id
	rid string
}

type importAssociationInterface interface {
	ParsePrimaryKey() error
	ImportAssociation() map[int]string
}

func NewImportAssociation(ctx context.Context, cli *association, params types.ContextParams, objID string, importData map[int]metadata.ExcelAssocation) importAssociationInterface {

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

	err = ia.getAssociationObjUnique()
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
	queryInput := &metadata.SearchAssociationObjectRequest{
		Condition: cond.ToMapStr(),
	}

	rsp, err := ia.cli.clientSet.ObjectController().Association().SearchObject(ia.ctx, ia.params.Header, queryInput)
	if nil != err {
		blog.Errorf("[getAssociationInfo] failed to request the object controller , error info is %s, input:%+v, rid:%s", err.Error(), queryInput, ia.rid)
		return ia.params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if !rsp.Result {
		blog.Errorf("[getAssociationInfo] failed to search the inst association, error info is %s, input:%+v, rid:%s", rsp.ErrMsg, queryInput, ia.rid)
		return ia.params.Err.New(rsp.Code, rsp.ErrMsg)
	}

	for _, inst := range rsp.Data {
		ia.asstIDInfoMap[inst.AssociationName] = inst
	}

	return nil
}

func (ia *importAssociation) getAssociationObjProperty() error {
	var objIDArr []string
	for _, info := range ia.asstIDInfoMap {
		objIDArr = append(objIDArr, info.AsstObjID)
	}
	objIDArr = append(objIDArr, ia.objID)

	cond := condition.CreateCondition()
	cond.Field(common.BKObjIDField).In(objIDArr)

	rsp, err := ia.cli.clientSet.ObjectController().Meta().SelectObjectAttWithParams(context.Background(), ia.params.Header, cond.ToMapStr())
	if nil != err {
		blog.Errorf("[getAssociationInfo] failed to  search attribute , error info is %s, input:%+v, rid:%s", err.Error(), cond, ia.rid)
		return ia.params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if !rsp.Result {
		blog.Errorf("[getAssociationInfo] failed to search attribute, error code:%s, error messge: %s, input:%+v, rid:%s", rsp.Code, rsp.ErrMsg, cond, ia.rid)
		return ia.params.Err.New(rsp.Code, rsp.ErrMsg)
	}

	for _, attr := range rsp.Data {
		_, ok := ia.asstObjIDProperty[attr.ObjectID]
		if !ok {
			ia.asstObjIDProperty[attr.ObjectID] = make(map[string]metadata.Attribute)
		}
		ia.asstObjIDProperty[attr.ObjectID][attr.PropertyName] = attr
	}

	return nil

}

func (ia *importAssociation) getAssociationObjUnique() error {
	var objIDArr = []string{ia.objID}
	for _, info := range ia.asstIDInfoMap {
		objIDArr = append(objIDArr, info.AsstObjID)
	}

	ia.objUniques = map[string]map[string]bool{}
	for _, objID := range objIDArr {
		rsp, err := ia.cli.clientSet.ObjectController().Unique().Search(context.Background(), ia.params.Header, objID)
		if nil != err {
			blog.Errorf("[getAssociationInfo] failed to  search attribute , error info is %s, input:%+v, rid:%s", err.Error(), objID, ia.rid)
			return ia.params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
		}

		if !rsp.Result {
			blog.Errorf("[getAssociationInfo] failed to search attribute, error code:%s, error messge: %s, input:%+v, rid:%s", rsp.Code, rsp.ErrMsg, objID, ia.rid)
			return ia.params.Err.New(rsp.Code, rsp.ErrMsg)
		}

		for _, unique := range rsp.Data {
			keynames := []string{}
			for _, keyitem := range unique.Keys {
				for _, property := range ia.asstObjIDProperty[objID] {
					if uint64(property.ID) == keyitem.ID {
						keynames = append(keynames, property.PropertyName)
					}
				}
			}
			sort.Strings(keynames)
			if _, ok := ia.objUniques[objID]; !ok {
				ia.objUniques[objID] = map[string]bool{}
			}
			ia.objUniques[objID][strings.Join(keynames, ",")] = true
		}
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
			continue
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

	keys := []string{}
	for _, primary := range primaryArr {
		primary = strings.TrimSpace(primary)
		keyValArr := strings.Split(primary, common.ExcelAsstPrimaryKeyJoinChar)
		if len(keyValArr) != 2 {
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
		keys = append(keys, keyValArr[0])
	}
	sort.Strings(keys)
	if !ia.objUniques[objID][strings.Join(keys, ",")] {
		var key = ""
		for key = range ia.objUniques[objID] {
			break
		}
		return nil, fmt.Errorf(ia.params.Lang.Languagef("import_asst_obj_property_str_primary_count_len", objID, item, key))
	}

	return keyValMap, nil

}

func (ia *importAssociation) getInstDataByConds() error {

	for objID, valArr := range ia.queryInstConds {

		instIDKey := metadata.GetInstIDFieldByObjID(objID)
		conds := condition.CreateCondition()
		if !util.IsInnerObject(objID) {
			conds.Field(common.BKObjIDField).Eq(objID)
		}
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
	queryInput := &metadata.QueryInput{}

	queryInput.Condition = conds.ToMapStr()
	queryInput.Fields = strings.Join(fields, ",")

	instSearchResult, err := ia.cli.clientSet.ObjectController().Instance().SearchObjects(ia.ctx, objID, ia.params.Header, queryInput)
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
	input := &metadata.SearchAssociationInstRequest{
		Condition: cond.ToMapStr(),
	}

	result, err := ia.cli.clientSet.ObjectController().Association().SearchInst(ia.ctx, ia.params.Header, input)
	if err != nil {
		ia.parseImportDataErr[idx] = err.Error()
		return
	}

	if !result.Result {
		ia.parseImportDataErr[idx] = result.ErrMsg
		return
	}

	if len(result.Data) == 0 {
		ia.parseImportDataErr[idx] = "can not find this association."
		return
	}

	asso := *result.Data[0]
	rsp, err := ia.cli.clientSet.ObjectController().Association().DeleteInst(ia.ctx, ia.params.Header, asso.ID)
	if err != nil {
		ia.parseImportDataErr[idx] = err.Error()
		return
	}

	if !rsp.Result {
		ia.parseImportDataErr[idx] = rsp.ErrMsg
		return
	}
}

func (ia *importAssociation) addSrcAssociation(idx int, asstFlag string, instID, assInstID int64) {
	_, ok := ia.parseImportDataErr[idx]
	if ok {
		return
	}
	inst := &metadata.CreateAssociationInstRequest{}
	inst.ObjectAsstId = asstFlag
	inst.InstId = instID
	inst.AsstInstId = assInstID
	rsp, err := ia.cli.clientSet.ObjectController().Association().CreateInst(ia.ctx, ia.params.Header, inst)
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
	input := &metadata.SearchAssociationInstRequest{
		Condition: cond.ToMapStr(),
	}
	rsp, err := ia.cli.clientSet.ObjectController().Association().SearchInst(ia.ctx, ia.params.Header, input)
	if err != nil {
		return false, err
	}
	if !rsp.Result {
		ia.parseImportDataErr[idx] = rsp.ErrMsg
		return false, ia.params.Err.New(rsp.Code, rsp.ErrMsg)
	}

	if len(rsp.Data) == 0 {
		return false, nil
	}
	if rsp.Data[0].AsstInstID != dstInstID &&
		asstMapping == metadata.OneToOneMapping {
		return false, ia.params.Err.Errorf(common.CCErrCommDuplicateItem)
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
	case common.FieldTypeForeignKey:
		if attr.PropertyID == common.BKCloudIDField {
			return util.GetInt64ByInterface(val)
		}
		fallthrough
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
