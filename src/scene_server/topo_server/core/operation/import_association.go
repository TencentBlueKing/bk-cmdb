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

	"configcenter/src/ac/extensions"
	"configcenter/src/ac/meta"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	"configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/language"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

func (assoc *association) ImportInstAssociation(ctx context.Context, kit *rest.Kit, objID string,
	importData map[int]metadata.ExcelAssociation, asstObjectUniqueIDMap map[string]int64, objectUniqueID int64,
	languageIf language.CCLanguageIf) (resp metadata.ResponeImportAssociationData, err error) {
	ia := NewImportAssociation(ctx, assoc, kit, objID, importData, asstObjectUniqueIDMap, objectUniqueID,
		assoc.authManager, languageIf.CreateDefaultCCLanguageIf(util.GetLanguage(kit.Header)))
	err = ia.ParsePrimaryKey()
	if err != nil {
		return resp, err
	}

	errIdxMsgMap := ia.ImportAssociation()
	if len(errIdxMsgMap) > 0 {
		err = kit.CCError.Error(common.CCErrorTopoImportAssociation)
	}
	for row, msg := range errIdxMsgMap {
		resp.ErrMsgMap = append(resp.ErrMsgMap, metadata.RowMsgData{
			Row: row,
			Msg: msg,
		})
	}

	return resp, err
}

func (assoc *association) FindAssociationByObjectAssociationID(ctx context.Context, kit *rest.Kit, objID string,
	asstIDArr []string) ([]metadata.Association, errors.CCError) {

	input := &metadata.QueryCondition{}
	input.Condition = map[string]interface{}{
		common.BKDBOR: []map[string]interface{}{
			{common.BKObjIDField: objID},
			{common.BKAsstObjIDField: objID},
		},
		common.AssociationObjAsstIDField: map[string]interface{}{common.BKDBIN: asstIDArr},
	}
	input.Page.Limit = common.BKNoLimit
	resp, err := assoc.clientSet.CoreService().Association().ReadModelAssociation(ctx, kit.Header, input)
	if err != nil {
		blog.ErrorJSON("find object by association http do error. err: %s, input: %s, rid: %s",
			err.Error(), input, kit.Rid)
		return nil, err
	}

	return resp.Info, nil
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
	importData map[int]metadata.ExcelAssociation
	// 模型使用的唯一校验相关的信息
	asstObjectUniqueIDMap map[string]int64
	objectUniqueID        int64
	kit                   *rest.Kit
	language              language.DefaultCCLanguageIf

	// map[AssociationName]Association alias  map[association flag]Association
	asstIDInfoMap map[string]*metadata.Association
	// asst obj info  map[objID]map[property name] attribute
	asstObjIDProperty map[string]map[string]metadata.Attribute
	// 当前操作模型使用的唯一校验，用来解决自关联使用不同的唯一校验
	objIDProperty map[string]metadata.Attribute

	parseImportDataErr map[int]string
	//map[objID][]condition.Condition， 查询与当前操作模型有关联关系的实例参数
	queryAsstInstCondArr map[string][]mapstr.MapStr
	//[]condition.Condition, 查询当前操作模型的的实例参数
	queryInstCondArr []mapstr.MapStr

	// map[objID][instance id]strings.Join([]string{property name, property value}, "=")[]importAssociationInst
	asstInstIDAttrKeyValMap map[string]map[string][]*importAssociationInst
	// map[instance id]strings.Join([]string{property name, property value}, "=")[]importAssociationInst
	instIDAttrKeyValMap map[string][]*importAssociationInst
	//http header http request id
	rid string

	authManager *extensions.AuthManager
}

type importAssociationInterface interface {
	ParsePrimaryKey() error
	ImportAssociation() map[int]string
}

func NewImportAssociation(ctx context.Context, cli *association, kit *rest.Kit, objID string,
	importData map[int]metadata.ExcelAssociation, asstObjectUniqueIDMap map[string]int64, objectUniqueID int64,
	authManager *extensions.AuthManager, languageIf language.DefaultCCLanguageIf) importAssociationInterface {
	return &importAssociation{
		objID:                 objID,
		cli:                   cli,
		ctx:                   ctx,
		importData:            importData,
		asstObjectUniqueIDMap: asstObjectUniqueIDMap,
		objectUniqueID:        objectUniqueID,

		kit:      kit,
		language: languageIf,

		asstIDInfoMap:           make(map[string]*metadata.Association, 0),
		asstObjIDProperty:       make(map[string]map[string]metadata.Attribute, 0),
		objIDProperty:           make(map[string]metadata.Attribute, 0),
		parseImportDataErr:      make(map[int]string),
		queryAsstInstCondArr:    make(map[string][]mapstr.MapStr),
		queryInstCondArr:        make([]mapstr.MapStr, 0),
		asstInstIDAttrKeyValMap: make(map[string]map[string][]*importAssociationInst),
		instIDAttrKeyValMap:     make(map[string][]*importAssociationInst),
		rid:                     kit.Rid,

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

	err = ia.getObjProperty()
	if err != nil {
		return err
	}

	err = ia.getAssociationObjProperty()
	if err != nil {
		return err
	}

	ia.parseImportDataPrimary()
	err = ia.getInstDataByQueryCondArr()
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
			ia.parseImportDataErr[idx] = ia.language.Languagef("import_association_id_not_found", asstInfo.ObjectAsstID)
			continue
		}

		srcInstID, dstInstID, err := int64(0), int64(0), error(nil)
		if asstID.ObjectID == ia.objID {
			srcInstID, err = ia.getObjectInstIDByPrimaryKey(asstInfo.SrcPrimary)
			if err != nil {
				ia.parseImportDataErr[idx] = err.Error()
				continue
			}
			dstInstID, err = ia.getAssociationObjectInstIDByPrimaryKey(asstID.AsstObjID, asstInfo.DstPrimary)
			if err != nil {
				ia.parseImportDataErr[idx] = err.Error()
				continue
			}
		} else {
			srcInstID, err = ia.getAssociationObjectInstIDByPrimaryKey(asstID.ObjectID, asstInfo.SrcPrimary)
			if err != nil {
				ia.parseImportDataErr[idx] = err.Error()
				continue
			}
			dstInstID, err = ia.getObjectInstIDByPrimaryKey(asstInfo.DstPrimary)
			if err != nil {
				ia.parseImportDataErr[idx] = err.Error()
				continue
			}
		}

		err = ia.authManager.AuthorizeByInstanceID(ia.ctx, ia.kit.Header, meta.Update, ia.objID, srcInstID)
		if err != nil {
			ia.parseImportDataErr[idx] = err.Error()
			continue
		}
		err = ia.authManager.AuthorizeByInstanceID(ia.ctx, ia.kit.Header, meta.Update, asstID.AsstObjID, dstInstID)
		if err != nil {
			ia.parseImportDataErr[idx] = err.Error()
			continue
		}
		switch asstInfo.Operate {
		case metadata.ExcelAssociationOperateAdd:

			conds := condition.CreateCondition()
			conds.Field(common.AssociationObjAsstIDField).Eq(asstInfo.ObjectAsstID)
			conds.Field(common.BKObjIDField).Eq(asstID.ObjectID)
			conds.Field(common.BKInstIDField).Eq(srcInstID)
			conds.Field(common.AssociatedObjectIDField).Eq(asstID.AsstObjID)
			isExist, err := ia.isExistInstAsst(idx, conds, dstInstID, asstID.ObjectID, asstID.Mapping)
			if err != nil {
				ia.parseImportDataErr[idx] = err.Error()
				continue
			}
			if isExist {
				continue
			}

			ia.addSrcAssociation(idx, asstID.AssociationName, srcInstID, dstInstID)
		case metadata.ExcelAssociationOperateDelete:
			conds := condition.CreateCondition()
			conds.Field(common.AssociationObjAsstIDField).Eq(asstInfo.ObjectAsstID)
			conds.Field(common.BKObjIDField).Eq(asstID.ObjectID)
			conds.Field(common.BKInstIDField).Eq(srcInstID)
			conds.Field(common.AssociatedObjectIDField).Eq(asstID.AsstObjID)
			conds.Field(common.BKAsstInstIDField).Eq(dstInstID)
			ia.delSrcAssociation(idx, ia.objID, conds)
		default:
			ia.parseImportDataErr[idx] = ia.language.Language("import_association_operate_not_found")
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
	or := cond.NewOR()
	or.Item(map[string]interface{}{common.BKObjIDField: ia.objID})
	or.Item(map[string]interface{}{common.BKAsstObjIDField: ia.objID})

	queryInput := &metadata.QueryCondition{Condition: cond.ToMapStr()}

	rsp, err := ia.cli.clientSet.CoreService().Association().ReadModelAssociation(ia.ctx, ia.kit.Header, queryInput)
	if nil != err {
		blog.Errorf("[getAssociationInfo] failed to request the object controller , error info is %s, input:%+v, rid:%s", err.Error(), queryInput, ia.rid)
		return ia.kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	for index := range rsp.Info {
		ia.asstIDInfoMap[rsp.Info[index].AssociationName] = &rsp.Info[index]
	}

	return nil
}

func (ia *importAssociation) getAssociationObjProperty() error {
	var objIDArr []string
	var uniqueIDArr []int64
	for objID, uniqueID := range ia.asstObjectUniqueIDMap {
		objIDArr = append(objIDArr, objID)
		uniqueIDArr = append(uniqueIDArr, uniqueID)
	}
	objIDArr = append(objIDArr, ia.objID)

	uniqueCond := condition.CreateCondition()
	uniqueCond.Field(common.BKFieldID).In(uniqueIDArr)

	uniqueQueryCond := metadata.QueryCondition{Condition: uniqueCond.ToMapStr()}
	uniqueResult, err := ia.cli.clientSet.CoreService().Model().ReadModelAttrUnique(ia.ctx, ia.kit.Header, uniqueQueryCond)
	if err != nil {
		blog.ErrorJSON("[getAssociationInfo] http do error.  search model unique , error info is %s, input:%s, rid:%s", err.Error(), uniqueQueryCond, ia.rid)
		return ia.kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	var propertyIDArr []uint64
	for _, unique := range uniqueResult.Info {
		for _, property := range unique.Keys {
			propertyIDArr = append(propertyIDArr, property.ID)
		}
	}

	cond := condition.CreateCondition()
	cond.Field(common.BKObjIDField).In(objIDArr)
	cond.Field(common.BKFieldID).In(propertyIDArr)

	attrCond := &metadata.QueryCondition{Condition: cond.ToMapStr()}
	attrCond.Fields = []string{common.BKFieldID, common.BKObjIDField, common.BKPropertyIDField, common.BKPropertyNameField}
	rsp, err := ia.cli.clientSet.CoreService().Model().ReadModelAttrByCondition(ia.ctx, ia.kit.Header, attrCond)
	if nil != err {
		blog.Errorf("[getAssociationInfo] failed to  search attribute , error info is %s, input:%+v, rid:%s", err.Error(), attrCond, ia.rid)
		return ia.kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	for _, attr := range rsp.Info {
		_, ok := ia.asstObjIDProperty[attr.ObjectID]
		if !ok {
			ia.asstObjIDProperty[attr.ObjectID] = make(map[string]metadata.Attribute)
		}
		ia.asstObjIDProperty[attr.ObjectID][attr.PropertyName] = attr
	}

	return nil

}

func (ia *importAssociation) getObjProperty() error {

	uniqueCond := condition.CreateCondition()
	uniqueCond.Field(common.BKFieldID).In(ia.objectUniqueID)

	uniqueQueryCond := metadata.QueryCondition{Condition: uniqueCond.ToMapStr()}
	uniqueResult, err := ia.cli.clientSet.CoreService().Model().ReadModelAttrUnique(ia.ctx, ia.kit.Header, uniqueQueryCond)
	if nil != err {
		blog.ErrorJSON(" http do error.  search model unique , error info is %s, input:%s, rid:%s",
			err.Error(), uniqueQueryCond, ia.rid)
		return ia.kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	var propertyIDArr []uint64
	for _, unique := range uniqueResult.Info {
		for _, property := range unique.Keys {
			propertyIDArr = append(propertyIDArr, property.ID)
		}
	}

	cond := condition.CreateCondition()
	cond.Field(common.BKObjIDField).In(ia.objID)
	cond.Field(common.BKFieldID).In(propertyIDArr)

	attrCond := &metadata.QueryCondition{Condition: cond.ToMapStr()}
	attrCond.Fields = []string{common.BKFieldID, common.BKObjIDField, common.BKPropertyIDField, common.BKPropertyNameField}
	rsp, err := ia.cli.clientSet.CoreService().Model().ReadModelAttrByCondition(ia.ctx, ia.kit.Header, attrCond)
	if nil != err {
		blog.ErrorJSON("search attribute failed, error info is %s, input:%s, rid:%s", err.Error(), attrCond, ia.rid)
		return ia.kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	for _, attr := range rsp.Info {
		ia.objIDProperty[attr.PropertyName] = attr
	}

	return nil

}

func (ia *importAssociation) parseImportDataPrimary() {

	for idx, info := range ia.importData {

		associationInst, ok := ia.asstIDInfoMap[info.ObjectAsstID]
		if !ok {
			ia.parseImportDataErr[idx] = ia.language.Languagef("import_asstid_not_found", info.ObjectAsstID)
			continue
		}

		var srcPropertyArr map[string]metadata.Attribute
		var dstPropertyArr map[string]metadata.Attribute

		isSelfObject := false
		if associationInst.ObjectID == ia.objID {
			srcPropertyArr = ia.objIDProperty
			dstPropertyArr = ia.asstObjIDProperty[associationInst.AsstObjID]
			isSelfObject = true
			if _, ok = ia.queryAsstInstCondArr[associationInst.AsstObjID]; !ok {
				ia.queryAsstInstCondArr[associationInst.AsstObjID] = make([]mapstr.MapStr, 0)
			}

		} else {
			srcPropertyArr = ia.asstObjIDProperty[associationInst.ObjectID]
			dstPropertyArr = ia.objIDProperty
			if _, ok = ia.queryAsstInstCondArr[associationInst.ObjectID]; !ok {
				ia.queryAsstInstCondArr[associationInst.ObjectID] = make([]mapstr.MapStr, 0)
			}

		}
		srcCond, err := ia.parseImportDataPrimaryItem(associationInst.ObjectID, info.SrcPrimary, srcPropertyArr)
		if err != nil {
			ia.parseImportDataErr[idx] = err.Error()
		} else {
			if isSelfObject {
				ia.queryInstCondArr = append(ia.queryInstCondArr, srcCond)
			} else {

				ia.queryAsstInstCondArr[associationInst.ObjectID] =
					append(ia.queryAsstInstCondArr[associationInst.ObjectID], srcCond)
			}

		}

		dstCond, err := ia.parseImportDataPrimaryItem(associationInst.AsstObjID, info.DstPrimary, dstPropertyArr)
		if err != nil {
			ia.parseImportDataErr[idx] = err.Error()
		} else {
			if isSelfObject {
				ia.queryAsstInstCondArr[associationInst.AsstObjID] =
					append(ia.queryAsstInstCondArr[associationInst.AsstObjID], dstCond)
			} else {
				ia.queryInstCondArr = append(ia.queryInstCondArr, dstCond)
			}

		}

	}

	return

}

func (ia *importAssociation) parseImportDataPrimaryItem(objID string, item string,
	propertyMap map[string]metadata.Attribute) (mapstr.MapStr, error) {
	keyValMap := mapstr.New()
	primaryArr := strings.Split(item, common.ExcelAsstPrimaryKeySplitChar)

	for _, primary := range primaryArr {

		primary = strings.TrimSpace(primary)
		keyValArr := strings.Split(primary, common.ExcelAsstPrimaryKeyJoinChar)
		if len(keyValArr) != 2 {
			blog.ErrorJSON("parseImportDataPrimaryItem eror. primary:%s, rid:%s", primary, ia.rid)
			return nil, fmt.Errorf(ia.language.Languagef("import_asst_obj_property_str_primary_format_error", objID, item))
		}
		attr, ok := propertyMap[keyValArr[0]]
		if !ok {
			return nil, fmt.Errorf(ia.language.Languagef("import_asst_obj_primary_property_str_not_found", objID, keyValArr[0]))
		}
		realVal, err := convStrToCCType(keyValArr[1], attr)
		if err != nil {
			return nil, fmt.Errorf(ia.language.Languagef("import_asst_obj_property_str_primary_type_error", objID, keyValArr[0]))
		}

		keyValMap[attr.PropertyID] = realVal
	}
	if len(keyValMap) != len(propertyMap) {
		blog.ErrorJSON("parseImportDataPrimaryItem error. keyVal:%s, objID:%s, objIDProperty:%s,rid:%s",
			keyValMap, objID, propertyMap[objID], ia.rid)
		return nil, fmt.Errorf(ia.language.Languagef("import_asst_obj_property_str_primary_count_len", objID, item))
	}

	return keyValMap, nil
}

func (ia *importAssociation) getInstDataByQueryCondArr() error {

	for objID, valArr := range ia.queryAsstInstCondArr {
		instArr, err := ia.getObjectInstDataByCondArr(objID, valArr, ia.asstObjIDProperty[objID])
		if err != nil {
			return err
		}

		instIDKey := metadata.GetInstIDFieldByObjID(objID)
		for _, inst := range instArr {
			ia.parseInstToImportAssociationObjectInst(objID, instIDKey, inst)
		}
	}

	instArr, err := ia.getObjectInstDataByCondArr(ia.objID, ia.queryInstCondArr, ia.objIDProperty)
	if err != nil {
		return err
	}

	instIDKey := metadata.GetInstIDFieldByObjID(ia.objID)
	for _, inst := range instArr {
		ia.parseInstToImportObjectInst(ia.objID, instIDKey, inst)
	}

	return nil
}

// 获取模型实例数据
func (ia *importAssociation) getObjectInstDataByCondArr(objID string, valArr []mapstr.MapStr,
	attrs map[string]metadata.Attribute) ([]mapstr.MapStr, error) {
	instIDKey := metadata.GetInstIDFieldByObjID(objID)
	if objID == common.BKInnerObjIDHost && len(valArr) > 0 {
		for idx, val := range valArr {
			if ok := val.Exists(common.BKCloudIDField); !ok {
				continue
			}
			intCloudID, err := val.Int64(common.BKCloudIDField)
			if err != nil {
				return nil, err
			}
			valArr[idx][common.BKCloudIDField] = intCloudID
		}
	}
	if len(valArr) == 0 {
		return nil, nil
	}
	conds := condition.CreateCondition()
	conds.NewOR().MapStrArr(valArr)
	instArr, err := ia.getInstDataByObjIDCondArr(objID, instIDKey, conds, attrs)
	if err != nil {
		return nil, err
	}

	return instArr, err
}

func (ia *importAssociation) getInstDataByObjIDCondArr(objID, instIDKey string, conds condition.Condition,
	attrs map[string]metadata.Attribute) (
	[]mapstr.MapStr, error) {

	var fields []string
	for _, attr := range attrs { //ia.asstObjIDProperty[objID] {
		fields = append(fields, attr.PropertyID)
	}

	fields = append(fields, instIDKey)
	queryInput := &metadata.QueryCondition{}
	queryInput.Condition = conds.ToMapStr()
	queryInput.Fields = fields

	instSearchResult, err := ia.cli.clientSet.CoreService().Instance().ReadInstance(ia.ctx, ia.kit.Header, objID, queryInput)
	if err != nil {
		blog.ErrorJSON("failed to  search %s instance , error info is %s, input:%s, rid:%s",
			objID, err.Error(), queryInput, ia.rid)
		return nil, ia.kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	return instSearchResult.Info, nil
}

// 导入模型关联对象实例数据
func (ia *importAssociation) parseInstToImportAssociationObjectInst(objID, instIDKey string, inst mapstr.MapStr) {

	_, ok := ia.asstInstIDAttrKeyValMap[objID]
	if !ok {
		ia.asstInstIDAttrKeyValMap[objID] = make(map[string][]*importAssociationInst)
	}

	attrs := ia.asstObjIDProperty[objID]

	instInfoArr, err := ia.parseInstToImportAssociationInstInfo(objID, instIDKey, inst, attrs)
	if err != nil {
		// 沿用已有逻辑
		return
	}

	ia.asstInstIDAttrKeyValMap[objID] = mergeInstToImportAssociationInst(ia.asstInstIDAttrKeyValMap[objID], instInfoArr)
	return
}

// 导入模型数据实例查询， 自关联的时候是src 对象
func (ia *importAssociation) parseInstToImportObjectInst(objID, instIDKey string, inst mapstr.MapStr) {

	instInfoArr, err := ia.parseInstToImportAssociationInstInfo(objID, instIDKey, inst, ia.objIDProperty)
	if err != nil {
		// 沿用已有逻辑
		return
	}

	ia.instIDAttrKeyValMap = mergeInstToImportAssociationInst(ia.instIDAttrKeyValMap, instInfoArr)
	return
}

func (ia *importAssociation) parseInstToImportAssociationInstInfo(objID, instIDKey string, inst mapstr.MapStr,
	attrs map[string]metadata.Attribute) (map[string][]*importAssociationInst, error) {
	instID, err := inst.Int64(instIDKey)
	//inst info can not found
	if err != nil {
		blog.Warnf("parseInstToImportAssociationInst get %s field from %s model error,error:%s, rid:%d ",
			instID, objID, err.Error(), ia.rid)
		return nil, err
	}

	attrNameValMap := importAssociationInst{
		instID:      instID,
		attrNameVal: make(map[string]bool),
	}

	for _, attr := range attrs {
		val, err := inst.String(attr.PropertyID)
		//inst info can not found
		if err != nil {
			blog.Warnf("get %s field from %s model error,error:%s, rid:%d ",
				attr.PropertyID, objID, err.Error(), ia.rid)
			return nil, err
		}
		attrNameValMap.attrNameVal[buildPrimaryStr(attr.PropertyName, val)] = true
	}

	instIDAttrKeyValMap := make(map[string][]*importAssociationInst, 0)
	for key := range attrNameValMap.attrNameVal {
		instIDAttrKeyValMap[key] = append(instIDAttrKeyValMap[key], &attrNameValMap)
	}

	return instIDAttrKeyValMap, nil
}

func (ia *importAssociation) delSrcAssociation(idx int, objID string, cond condition.Condition) {
	_, ok := ia.parseImportDataErr[idx]
	if ok {
		return
	}

	delOpt := &metadata.InstAsstDeleteOption{
		Opt:   metadata.DeleteOption{Condition: cond.ToMapStr()},
		ObjID: objID,
	}

	_, err := ia.cli.clientSet.CoreService().Association().DeleteInstAssociation(ia.ctx, ia.kit.Header, delOpt)
	if err != nil {
		ia.parseImportDataErr[idx] = err.Error()
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
	inst.Data.ObjectID = asstInfo.ObjectID
	inst.Data.AsstObjectID = asstInfo.AsstObjID
	inst.Data.AsstInstID = assInstID
	inst.Data.AssociationKindID = asstInfo.AsstKindID
	_, err := ia.cli.clientSet.CoreService().Association().CreateInstAssociation(ia.ctx, ia.kit.Header, &inst)
	if err != nil {
		ia.parseImportDataErr[idx] = err.Error()
	}
}

func (ia *importAssociation) isExistInstAsst(idx int, cond condition.Condition, dstInstID int64, objID string,
	asstMapping metadata.AssociationMapping) (isExit bool, err error) {

	_, ok := ia.parseImportDataErr[idx]
	if ok {
		return
	}
	if asstMapping != metadata.OneToOneMapping {
		cond.Field(common.BKAsstInstIDField).Eq(dstInstID)
	}

	queryCond := &metadata.InstAsstQueryCondition{
		Cond:  metadata.QueryCondition{Condition: cond.ToMapStr()},
		ObjID: objID,
	}
	rsp, err := ia.cli.clientSet.CoreService().Association().ReadInstAssociation(ia.ctx, ia.kit.Header, queryCond)
	if err != nil {
		ia.parseImportDataErr[idx] = err.Error()
		return false, err
	}

	if len(rsp.Info) == 0 {
		return false, nil
	}
	if rsp.Info[0].AsstInstID != dstInstID &&
		asstMapping == metadata.OneToOneMapping {
		return false, ia.kit.CCError.Errorf(common.CCErrCommDuplicateItem, "association")
	}

	return true, nil
}

func (ia *importAssociation) getAssociationObjectInstIDByPrimaryKey(objID, primary string) (int64, error) {
	primaryArr := strings.Split(primary, common.ExcelAsstPrimaryKeySplitChar)
	if len(primaryArr) == 0 {
		return 0, fmt.Errorf(ia.language.Languagef("import_instance_not_found", objID, primary))
	}

	instArr, ok := ia.asstInstIDAttrKeyValMap[objID][primaryArr[0]]
	if !ok {
		return 0, fmt.Errorf(ia.language.Languagef("import_instance_not_found", objID, primaryArr[0]))
	}

	if instID := findInst(instArr, primaryArr); instID != 0 {
		return instID, nil
	}

	return 0, fmt.Errorf(ia.language.Languagef("import_instance_not_found", objID, primary))

}

func (ia *importAssociation) getObjectInstIDByPrimaryKey(primary string) (int64, error) {
	primaryArr := strings.Split(primary, common.ExcelAsstPrimaryKeySplitChar)
	if len(primaryArr) == 0 {
		return 0, fmt.Errorf(ia.language.Languagef("import_instance_not_found", ia.objID, primary))
	}

	instArr, ok := ia.instIDAttrKeyValMap[primaryArr[0]]
	if !ok {
		return 0, fmt.Errorf(ia.language.Languagef("import_instance_not_found", ia.objID, primaryArr[0]))
	}

	if instID := findInst(instArr, primaryArr); instID != 0 {
		return instID, nil
	}

	return 0, fmt.Errorf(ia.language.Languagef("import_instance_not_found", ia.objID, primary))

}

func findInst(instArr []*importAssociationInst, primaryArr []string) int64 {
	for _, inst := range instArr {

		isEq := true
		for _, item := range primaryArr {
			if _, ok := inst.attrNameVal[item]; !ok {
				isEq = false
				break
			}
		}
		if isEq {
			return inst.instID
		}

	}

	return 0
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

func mergeInstToImportAssociationInst(src, dst map[string][]*importAssociationInst) map[string][]*importAssociationInst {
	if dst == nil {
		return src
	}
	for key, valArr := range src {
		dst[key] = append(dst[key], valArr...)
	}

	return dst
}
