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

package model

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/universalsql"
	"configcenter/src/common/universalsql/mongo"
	"configcenter/src/common/util"
	"configcenter/src/source_controller/coreservice/core"
	"configcenter/src/storage/dal"
)

var (
	// notAddAttrModel 不允许新加属性的模型
	notAddAttrModel = map[string]bool{
		common.BKInnerObjIDPlat: true,
		common.BKInnerObjIDProc: true,
	}

	// RequiredFieldUnchangeableModels 模型的属性描述，的必填字段不允许修改
	// example: 禁止如下修改
	// db.getCollection('cc_ObjAttDes').update(
	//     {bk_obj_id: {$in: ['biz', 'host', 'set', 'module', 'plat', 'process']}},
	//     {$set: {isrequired: true}}
	// )
	RequiredFieldUnchangeableModels = map[string]bool{
		common.BKInnerObjIDApp:    true,
		common.BKInnerObjIDHost:   true,
		common.BKInnerObjIDSet:    true,
		common.BKInnerObjIDModule: true,
		common.BKInnerObjIDPlat:   true,
		common.BKInnerObjIDProc:   true,
	}
)

func (m *modelAttribute) count(ctx core.ContextParams, cond universalsql.Condition) (cnt uint64, err error) {
	cnt, err = m.dbProxy.Table(common.BKTableNameObjAttDes).Find(cond.ToMapStr()).Count(ctx)
	return cnt, err
}

func (m *modelAttribute) save(ctx core.ContextParams, attribute metadata.Attribute) (id uint64, err error) {

	id, err = m.dbProxy.NextSequence(ctx, common.BKTableNameObjAttDes)
	if err != nil {
		return id, ctx.Error.New(common.CCErrObjectDBOpErrno, err.Error())
	}

	index, err := m.GetAttrLastIndex(ctx, attribute)
	if err != nil {
		return id, err
	}

	attribute.PropertyIndex = index
	attribute.ID = int64(id)
	attribute.OwnerID = ctx.SupplierAccount

	if nil == attribute.CreateTime {
		attribute.CreateTime = &metadata.Time{}
		attribute.CreateTime.Time = time.Now()
	}

	if nil == attribute.LastTime {
		attribute.LastTime = &metadata.Time{}
		attribute.LastTime.Time = time.Now()
	}

	if err = m.saveCheck(ctx, attribute); err != nil {
		return 0, err
	}

	err = m.dbProxy.Table(common.BKTableNameObjAttDes).Insert(ctx, attribute)
	return id, err
}

func (m *modelAttribute) checkUnique(ctx core.ContextParams, isCreate bool, objID, propertyID, propertyName string, meta metadata.Metadata) error {
	cond := map[string]interface{}{
		common.BKObjIDField: objID,
	}
	orCond := make([]map[string]interface{}, 0)

	if isCreate {
		nameFieldCond := map[string]interface{}{common.BKPropertyNameField: propertyName}
		idFieldCond := map[string]interface{}{common.BKPropertyIDField: propertyID}
		orCond = append(orCond, nameFieldCond, idFieldCond)
	} else {
		// update attribute. not change name, 无需判断
		if propertyName == "" {
			return nil
		}
		cond[common.BKPropertyIDField] = map[string]interface{}{common.BKDBNE: propertyID}
		cond[common.BKPropertyNameField] = propertyName
	}

	isExist, bizID := meta.Label.Get(common.BKAppIDField)
	if isExist {
		orCond = append(orCond, metadata.BizLabelNotExist, map[string]interface{}{metadata.MetadataBizField: bizID})
	}

	if len(orCond) > 0 {
		cond[common.BKDBOR] = orCond
	}
	condMap := util.SetModOwner(cond, ctx.SupplierAccount)

	resultAttrs := make([]metadata.Attribute, 0)
	err := m.dbProxy.Table(common.BKTableNameObjAttDes).Find(condMap).All(ctx, &resultAttrs)
	blog.V(5).Infof("checkUnique db cond:%#v, result:%#v, rid:%s", condMap, resultAttrs, ctx.ReqID)
	if err != nil {
		blog.ErrorJSON("checkUnique select error. err:%s, cond:%s, rid:%s", err.Error(), condMap, ctx.ReqID)
		return ctx.Error.Error(common.CCErrCommDBSelectFailed)
	}
	for _, attrItem := range resultAttrs {
		if attrItem.PropertyID == propertyID {
			return ctx.Error.Errorf(common.CCErrCommDuplicateItem, ctx.Lang.Language("model_attr_bk_property_id"))
		}
		if attrItem.PropertyName == propertyName {
			return ctx.Error.Errorf(common.CCErrCommDuplicateItem, ctx.Lang.Language("model_attr_bk_property_name"))
		}
	}

	return nil
}

func (m *modelAttribute) checkAttributeMustNotEmpty(ctx core.ContextParams, attribute metadata.Attribute) error {
	if attribute.PropertyID == "" {
		return ctx.Error.Errorf(common.CCErrCommParamsNeedSet, metadata.AttributeFieldPropertyID)
	}
	if attribute.PropertyName == "" {
		return ctx.Error.Errorf(common.CCErrCommParamsNeedSet, metadata.AttributeFieldPropertyName)
	}
	if attribute.PropertyType == "" {
		return ctx.Error.Errorf(common.CCErrCommParamsNeedSet, metadata.AttributeFieldPropertyType)
	}
	return nil
}

func (m *modelAttribute) checkAttributeValidity(ctx core.ContextParams, attribute metadata.Attribute) error {
	if common.AttributeIDMaxLength < utf8.RuneCountInString(attribute.PropertyID) {
		return ctx.Error.Errorf(common.CCErrCommValExceedMaxFailed, ctx.Lang.Language("model_attr_bk_property_id"), common.AttributeIDMaxLength)
	} else if attribute.PropertyID != "" {
		match, err := regexp.MatchString(common.FieldTypeStrictCharRegexp, attribute.PropertyID)
		if nil != err || !match {
			return ctx.Error.Errorf(common.CCErrCommParamsIsInvalid, metadata.AttributeFieldPropertyID)
		}
	}

	if common.AttributeNameMaxLength < utf8.RuneCountInString(attribute.PropertyName) {
		return ctx.Error.Errorf(common.CCErrCommValExceedMaxFailed, ctx.Lang.Language("model_attr_bk_property_name"), common.AttributeNameMaxLength)
	} else if attribute.PropertyName != "" {
		attribute.PropertyName = strings.TrimSpace(attribute.PropertyName)
		match, err := regexp.MatchString(common.FieldTypeSingleCharRegexp, attribute.PropertyName)
		if nil != err || !match {
			return ctx.Error.Errorf(common.CCErrCommParamsIsInvalid, metadata.AttributeFieldPropertyName)
		}
	}

	if attribute.Placeholder != "" {
		if common.AttributePlaceHolderMaxLength < utf8.RuneCountInString(attribute.Placeholder) {
			return ctx.Error.Errorf(common.CCErrCommValExceedMaxFailed, ctx.Lang.Language("model_attr_placeholder"), common.AttributePlaceHolderMaxLength)
		}
		attribute.Placeholder = strings.TrimSpace(attribute.Placeholder)
		match, err := regexp.MatchString(common.FieldTypeLongCharRegexp, attribute.Placeholder)
		if nil != err || !match {
			return ctx.Error.Errorf(common.CCErrCommParamsIsInvalid, metadata.AttributeFieldPlaceHolder)
		}
	}

	if attribute.Unit != "" {
		if common.AttributeUnitMaxLength < utf8.RuneCountInString(attribute.Unit) {
			return ctx.Error.Errorf(common.CCErrCommValExceedMaxFailed, ctx.Lang.Language("model_attr_uint"), common.AttributeUnitMaxLength)
		}
		attribute.Unit = strings.TrimSpace(attribute.Unit)
		match, err := regexp.MatchString(common.FieldTypeSingleCharRegexp, attribute.Unit)
		if nil != err || !match {
			return ctx.Error.Errorf(common.CCErrCommParamsIsInvalid, metadata.AttributeFieldUnit)
		}
	}

	if attribute.PropertyType != "" {
		switch attribute.PropertyType {
		case common.FieldTypeSingleChar, common.FieldTypeLongChar, common.FieldTypeInt, common.FieldTypeFloat, common.FieldTypeEnum,
			common.FieldTypeDate, common.FieldTypeTime, common.FieldTypeUser, common.FieldTypeTimeZone, common.FieldTypeBool, common.FieldTypeList:
		default:
			return ctx.Error.Errorf(common.CCErrCommParamsIsInvalid, metadata.AttributeFieldPropertyType)
		}
	}

	if opt, ok := attribute.Option.(string); ok && opt != "" {
		if common.AttributeOptionMaxLength < utf8.RuneCountInString(opt) {
			return ctx.Error.Errorf(common.CCErrCommValExceedMaxFailed, ctx.Lang.Language("model_attr_option_regex"), common.AttributeOptionMaxLength)
		}
	}

	return nil
}

func (m *modelAttribute) update(ctx core.ContextParams, data mapstr.MapStr, cond universalsql.Condition) (cnt uint64, err error) {
	cnt, err = m.checkUpdate(ctx, data, cond)
	if err != nil {
		blog.ErrorJSON("checkUpdate error. data:%s, cond:%s, rid:%s", data, cond, ctx.ReqID)
		return cnt, err
	}
	err = m.dbProxy.Table(common.BKTableNameObjAttDes).Update(ctx, cond.ToMapStr(), data)
	if nil != err {
		blog.Errorf("request(%s): database operation is failed, error info is %s", ctx.ReqID, err.Error())
		return 0, err
	}

	return cnt, err
}

func (m *modelAttribute) search(ctx core.ContextParams, cond universalsql.Condition) (resultAttrs []metadata.Attribute, err error) {
	resultAttrs = []metadata.Attribute{}
	err = m.dbProxy.Table(common.BKTableNameObjAttDes).Find(cond.ToMapStr()).All(ctx, &resultAttrs)
	return resultAttrs, err
}

func (m *modelAttribute) searchWithSort(ctx core.ContextParams, cond metadata.QueryCondition) (resultAttrs []metadata.Attribute, err error) {
	resultAttrs = []metadata.Attribute{}

	instHandler := m.dbProxy.Table(common.BKTableNameObjAttDes).Find(cond.Condition)
	for _, sort := range cond.SortArr {
		field := sort.Field
		if sort.IsDsc {
			field = "-" + field
		}
		instHandler = instHandler.Sort(field)
	}
	err = instHandler.Start(uint64(cond.Limit.Offset)).Limit(uint64(cond.Limit.Limit)).All(ctx, &resultAttrs)

	return resultAttrs, err
}

func (m *modelAttribute) searchReturnMapStr(ctx core.ContextParams, cond universalsql.Condition) (resultAttrs []mapstr.MapStr, err error) {

	resultAttrs = []mapstr.MapStr{}
	err = m.dbProxy.Table(common.BKTableNameObjAttDes).Find(cond.ToMapStr()).All(ctx, &resultAttrs)
	return resultAttrs, err
}

func (m *modelAttribute) delete(ctx core.ContextParams, cond universalsql.Condition) (cnt uint64, err error) {

	resultAttrs := make([]metadata.Attribute, 0)
	fields := []string{common.BKFieldID, common.BKPropertyIDField, common.BKObjIDField, common.MetadataField}

	condMap := util.SetModOwner(cond.ToMapStr(), ctx.SupplierAccount)
	err = m.dbProxy.Table(common.BKTableNameObjAttDes).Find(condMap).Fields(fields...).All(ctx, &resultAttrs)
	if nil != err {
		blog.Errorf("request(%s): database count operation is failed, error info is %s", ctx.ReqID, err.Error())
		return 0, err
	}

	cnt = uint64(len(resultAttrs))
	if cnt == 0 {
		return cnt, nil
	}

	objIDArrMap := make(map[string][]int64, 0)
	for _, attr := range resultAttrs {
		objIDArrMap[attr.ObjectID] = append(objIDArrMap[attr.ObjectID], attr.ID)
	}

	if err := m.cleanAttributeFieldInInstances(ctx, ctx.SupplierAccount, resultAttrs); err != nil {
		blog.ErrorJSON("delete object attributes with cond: %s, but delete these attribute in instance failed, err: %v, rid:%s", condMap, err, ctx.ReqID)
		return 0, err
	}

	exist, err := m.checkAttributeInUnique(ctx, objIDArrMap)
	if err != nil {
		blog.ErrorJSON("check attribute in unique error. err:%s, input:%s, rid:%s", err.Error(), condMap, ctx.ReqID)
		return 0, err
	}
	// delete field in module unique. not allow delete
	if exist {
		blog.ErrorJSON("delete field in unique. delete cond:%s, field:%s, rid:%s", condMap, resultAttrs, ctx.ReqID)
		return 0, ctx.Error.Error(common.CCErrCoreServiceNotAllowUniqueAttr)
	}

	err = m.dbProxy.Table(common.BKTableNameObjAttDes).Delete(ctx, condMap)
	if nil != err {
		blog.Errorf("request(%s): database deletion operation is failed, error info is %s", ctx.ReqID, err.Error())
		return 0, err
	}

	return cnt, err
}

type bizObjectFields struct {
	bizID  int64
	object string
	fields []string
}

// remove attribute filed in this object's instances
func (m *modelAttribute) cleanAttributeFieldInInstances(ctx context.Context, ownerID string, attrs []metadata.Attribute) error {

	objPublicFields := make(map[string][]string)
	objBizFields := make([]bizObjectFields, 0)

	// TODO: now, we only support set, module, host model's biz attribute clean operation.
	for _, attr := range attrs {

		biz, err := metadata.BizIDFromMetadata(attr.Metadata)
		if err != nil {
			return err
		}

		if biz != 0 {
			if !isBizObject(attr.ObjectID) {
				return fmt.Errorf("unsupported object %s's clean instance field operation", attr.ObjectID)
			}

			// this is a business attribute
			hit := false
			for index, ele := range objBizFields {
				if ele.object == attr.ObjectID && ele.bizID == biz {
					hit = true
					objBizFields[index].fields = append(objBizFields[index].fields, attr.PropertyID)
				}
			}

			if !hit {
				objBizFields = append(objBizFields, bizObjectFields{
					bizID:  biz,
					object: attr.ObjectID,
					fields: []string{attr.PropertyID},
				})
			}

		} else {
			// this is a public attribute
			_, exist := objPublicFields[attr.ObjectID]
			if !exist {
				objPublicFields[attr.ObjectID] = make([]string, 0)
			}
			objPublicFields[attr.ObjectID] = append(objPublicFields[attr.ObjectID], attr.PropertyID)
		}
	}

	// delete these attribute's filed in the model instance
	// step 1: handle object's public attribute
	var hitError error
	wg := sync.WaitGroup{}
	for object, fields := range objPublicFields {

		if len(fields) == 0 {
			// no fields need to be removed, skip directly.
			continue
		}

		var cond mapstr.MapStr
		if isBizObject(object) {
			if object == common.BKInnerObjIDHost {
				ele := bizObjectFields{
					bizID:  0,
					object: object,
					fields: fields,
				}
				if err := m.cleanHostAttributeField(ctx, ownerID, ele); err != nil {
					return err
				}
				return nil
			}

			cond = mapstr.MapStr{}
		} else {
			cond = mapstr.MapStr{
				common.BKObjIDField: object,
			}
		}
		cond = util.SetQueryOwner(cond, ownerID)

		collectionName := common.GetInstTableName(object)
		wg.Add(1)
		go func(collName string, filter dal.Filter, fields []string) {
			defer wg.Done()
			if err := m.dbProxy.Table(collName).DropColumns(ctx, filter, fields); err != nil {
				blog.Error("delete object's attribute from instance failed, table: %s, cond: %v, fields: %v, err: %v", collectionName, filter, fields, err)
				hitError = err
				return
			}
		}(collectionName, cond, fields)
	}
	// wait for all the public object routine is done.
	wg.Wait()
	if hitError != nil {
		return hitError
	}

	// step 2: handle object's biz attribute
	wg = sync.WaitGroup{}
	for _, ele := range objBizFields {
		if len(ele.fields) == 0 {
			// no fields need to be removed, skip directly.
			continue
		}
		if !isBizObject(ele.object) {
			return fmt.Errorf("unsupported object %s's clean instance field operation", ele.object)
		}

		if ele.object == common.BKInnerObjIDHost {
			if err := m.cleanHostAttributeField(ctx, ownerID, ele); err != nil {
				return err
			}
			return nil
		}

		cond := mapstr.MapStr{
			common.BKAppIDField: ele.bizID,
		}
		cond = util.SetQueryOwner(cond, ownerID)

		collectionName := common.GetInstTableName(ele.object)
		wg.Add(1)
		go func(collName string, filter dal.Filter, fields []string) {
			defer wg.Done()
			if err := m.dbProxy.Table(collName).DropColumns(ctx, filter, fields); err != nil {
				blog.Error("delete object's attribute from instance failed, table: %s, cond: %v, fields: %v, err: %v", collectionName, filter, fields, err)
				hitError = err
				return
			}
		}(collectionName, cond, ele.fields)
	}
	// wait for all the public object routine is done.
	wg.Wait()
	if hitError != nil {
		return hitError
	}

	return nil
}

const pageSize = 500

func (m *modelAttribute) cleanHostAttributeField(ctx context.Context, ownerID string, info bizObjectFields) error {
	cond := mapstr.MapStr{}
	cond = util.SetQueryOwner(cond, ownerID)
	// biz id = 0 means all the hosts.
	// TODO: optimize when the filed is a public filed in all the host instances. handle with page
	if info.bizID != 0 {
		// find hosts in this biz
		cond = mapstr.MapStr{
			common.BKAppIDField: info.bizID,
		}
	}
	type hostInst struct {
		HostID int64 `bson:"bk_host_id"`
	}
	hostList := make([]hostInst, 0)
	err := m.dbProxy.Table(common.BKTableNameModuleHostConfig).Find(cond).Fields(common.BKHostIDField).All(ctx, &hostList)
	if err != nil {
		return err
	}
	if len(hostList) == 0 {
		// no host in this business, do not need to clean the filed.
		return nil
	}

	count := len(hostList)
	for start := 0; start < count; start += pageSize {
		end := start + pageSize
		if end > count {
			end = count
		}
		ids := make([]int64, 0)
		for index := start; index < end; index++ {
			ids = append(ids, hostList[index].HostID)
		}
		hostFilter := mapstr.MapStr{
			common.BKHostIDField: mapstr.MapStr{common.BKDBIN: ids},
		}
		if err := m.dbProxy.Table(common.BKTableNameBaseHost).DropColumns(ctx, hostFilter, info.fields); err != nil {
			return fmt.Errorf("clean host biz attribute %v failed, err: %v", info.fields, err)
		}
	}

	return nil

}

// now, we only support set, module, host model's biz attribute clean operation.
func isBizObject(objectID string) bool {
	switch objectID {
	// biz is a special object, but it can not have biz attribute obviously.
	case common.BKInnerObjIDApp:
		return true
	case common.BKInnerObjIDHost:
		return true
	case common.BKInnerObjIDModule:
		return true
	case common.BKInnerObjIDSet:
		return true
	default:
		// TODO: remove this when the common object support biz attribute and biz instance field.
		return false

	}
}

//  saveCheck 新加字段检查
func (m *modelAttribute) saveCheck(ctx core.ContextParams, attribute metadata.Attribute) error {

	if err := m.checkAddField(ctx, attribute); err != nil {
		return err
	}

	if err := m.checkAttributeMustNotEmpty(ctx, attribute); err != nil {
		return err
	}
	if err := m.checkAttributeValidity(ctx, attribute); err != nil {
		return err
	}

	// check name duplicate
	if err := m.checkUnique(ctx, true, attribute.ObjectID, attribute.PropertyID, attribute.PropertyName, attribute.Metadata); err != nil {
		blog.ErrorJSON("save attribute check unique err:%s, input:%s, rid:%s", err.Error(), attribute, ctx.ReqID)
		return err
	}

	return nil
}

// checkUpdate 删除不可以更新字段，检验字段是否重复， 返回更新的行数，错误
func (m *modelAttribute) checkUpdate(ctx core.ContextParams, data mapstr.MapStr, cond universalsql.Condition) (changeRow uint64, err error) {

	dbAttributeArr, err := m.search(ctx, cond)
	if err != nil {
		blog.Errorf("request(%s): find nothing by the condition(%#v)  error(%s)", ctx.ReqID, cond.ToMapStr(), err.Error())
		return changeRow, err
	}
	if 0 == len(dbAttributeArr) {
		blog.Errorf("request(%s): find nothing by the condition(%#v)", ctx.ReqID, cond.ToMapStr())
		return changeRow, nil
	}

	// 更新的属性是否存在预定义字段。
	hasIsPreProperty := false
	for _, dbAttribute := range dbAttributeArr {
		if dbAttribute.IsPre == true {
			hasIsPreProperty = true
			break
		}
	}

	// 预定义字段，只能更新分组、分组内排序、名称、单位、提示语和option
	if hasIsPreProperty {
		hasNotAllowField := false
		_ = data.ForEach(func(key string, val interface{}) error {
			if key != metadata.AttributeFieldPropertyGroup &&
				key != metadata.AttributeFieldPropertyIndex &&
				key != metadata.AttributeFieldPropertyName &&
				key != metadata.AttributeFieldUnit &&
				key != metadata.AttributeFieldPlaceHolder &&
				key != metadata.AttributeFieldOption {
				hasNotAllowField = true
			}
			return nil
		})
		// 出现编辑预定义属性的字段
		if hasNotAllowField {
			blog.ErrorJSON("update model predefined attribute,input:%s, attr info:%s, rid:%s", cond.ToMapStr(), dbAttributeArr, ctx.ReqID)
			return changeRow, ctx.Error.Error(common.CCErrCoreServiceNotUpdatePredefinedAttrErr)
		}
	}

	if option, exists := data.Get(metadata.AttributeFieldOption); exists {
		propertyType := dbAttributeArr[0].PropertyType
		for _, dbAttribute := range dbAttributeArr {
			if dbAttribute.PropertyType != propertyType {
				blog.ErrorJSON("update option, but property type not the same, db attributes: %s, rid:%s", dbAttributeArr, ctx.ReqID)
				return changeRow, ctx.Error.Errorf(common.CCErrCommParamsInvalid, "cond")
			}
		}
		if err := util.ValidPropertyOption(propertyType, option, ctx.Error); err != nil {
			blog.ErrorJSON("valid property option failed, err: %s, data: %s, rid:%s", err, data, ctx.ReqID)
			return changeRow, err
		}
	}

	// 删除不可更新字段， 避免由于传入数据，修改字段
	// TODO: 改成白名单方式
	data.Remove(metadata.AttributeFieldPropertyID)
	data.Remove(metadata.AttributeFieldSupplierAccount)
	data.Remove(metadata.AttributeFieldPropertyType)
	data.Remove(metadata.AttributeFieldCreateTime)
	data.Remove(metadata.AttributeFieldIsPre)
	data.Set(metadata.AttributeFieldLastTime, time.Now())

	if grp, exists := data.Get(metadata.AttributeFieldPropertyGroup); exists {
		if grp == "" {
			data.Remove(metadata.AttributeFieldPropertyGroup)
		}
		// check if property group exists in object
		objIDs := make([]string, 0)
		for _, dbAttribute := range dbAttributeArr {
			objIDs = append(objIDs, dbAttribute.ObjectID)
		}
		objIDs = util.StrArrayUnique(objIDs)
		cond := map[string]interface{}{
			common.BKObjIDField: map[string]interface{}{
				common.BKDBIN: objIDs,
			},
			common.BKPropertyGroupIDField: grp,
		}
		cnt, err := m.dbProxy.Table(common.BKTableNamePropertyGroup).Find(cond).Count(ctx)
		if err != nil {
			blog.ErrorJSON("property group count failed, err: %s, condition: %s, rid: %s", err, cond, ctx.ReqID)
			return changeRow, err
		}
		if cnt != uint64(len(objIDs)) {
			blog.Errorf("property group invalid, objIDs: %s have %d property groups, rid: %s", objIDs, cnt, ctx.ReqID)
			return changeRow, ctx.Error.Errorf(common.CCErrCommParamsInvalid, metadata.AttributeFieldPropertyGroup)
		}
	}

	attribute := metadata.Attribute{}
	if err = data.MarshalJSONInto(&attribute); err != nil {
		blog.Errorf("request(%s): MarshalJSONInto(%#v), error is %v", ctx.ReqID, data, err)
		return changeRow, err
	}

	if err = m.checkAttributeValidity(ctx, attribute); err != nil {
		return changeRow, err
	}

	for _, dbAttribute := range dbAttributeArr {
		err = m.checkUnique(ctx, false, dbAttribute.ObjectID, dbAttribute.PropertyID, attribute.PropertyName, attribute.Metadata)
		if err != nil {
			blog.ErrorJSON("save attribute check unique err:%s, input:%s, rid:%s", err.Error(), dbAttribute, ctx.ReqID)
			return changeRow, err
		}
		if err = m.checkChangeField(ctx, dbAttribute, data); err != nil {
			blog.ErrorJSON("save attribute check change unique field err:%s, input:%s, rid:%s", err.Error(), dbAttribute, ctx.ReqID)
			return changeRow, err
		}
	}

	return uint64(len(dbAttributeArr)), err

}

// checkAttributeInUnique 检查属性是否存在唯一校验中  objIDPropertyIDArr  属性的bk_obj_id和表中ID的集合
func (m *modelAttribute) checkAttributeInUnique(ctx core.ContextParams, objIDPropertyIDArr map[string][]int64) (bool, error) {

	cond := mongo.NewCondition()

	var orCondArr []universalsql.ConditionElement
	for objID, propertyIDArr := range objIDPropertyIDArr {
		orCondItem := mongo.NewCondition()
		orCondItem.Element(mongo.Field(common.BKObjIDField).Eq(objID))
		orCondItem.Element(mongo.Field("keys.key_id").In(propertyIDArr))
		orCondItem.Element(mongo.Field("keys.key_kind").Eq("property"))
		orCondArr = append(orCondArr, orCondItem)
	}

	cond.Or(orCondArr...)
	condMap := util.SetQueryOwner(cond.ToMapStr(), ctx.SupplierAccount)

	cnt, err := m.dbProxy.Table(common.BKTableNameObjUnique).Find(condMap).Count(ctx)
	if err != nil {
		blog.ErrorJSON("checkAttributeInUnique db select error. err:%s, cond:%s, rid:%s", err.Error(), condMap, ctx.ReqID)
		return false, ctx.Error.Error(common.CCErrCommDBSelectFailed)
	}

	if cnt > 0 {
		return true, nil
	}

	return false, nil
}

// checkAddRequireField 新加模型属性的时候，如果新加的是必填字段，需要判断是否可以新加必填字段
func (m *modelAttribute) checkAddField(ctx core.ContextParams, attribute metadata.Attribute) error {
	langObjID := m.getLangObjID(ctx, attribute.ObjectID)
	if _, ok := notAddAttrModel[attribute.ObjectID]; ok {
		//  不允许新加字段的模型
		return ctx.Error.Errorf(common.CCErrCoreServiceNotAllowAddFieldErr, langObjID)
	}

	if _, ok := RequiredFieldUnchangeableModels[attribute.ObjectID]; ok {
		if attribute.IsRequired {
			//  不允许修改必填字段的模型
			return ctx.Error.Errorf(common.CCErrCoreServiceNotAllowAddRequiredFieldErr, langObjID)
		}

	}
	return nil
}

// 修改模型属性的时候，如果修改的属性包含是否为必填字段(isrequired)，需要判断该模型的必填字段是否允许被修改
func (m *modelAttribute) checkChangeField(ctx core.ContextParams, attr metadata.Attribute, attrInfo mapstr.MapStr) error {
	langObjID := m.getLangObjID(ctx, attr.ObjectID)
	if _, ok := RequiredFieldUnchangeableModels[attr.ObjectID]; ok {
		if attrInfo.Exists(metadata.AttributeFieldIsRequired) {
			// 不允许修改模型的必填字段
			val, ok := attrInfo[metadata.AttributeFieldIsRequired].(bool)
			if !ok {
				return ctx.Error.Errorf(common.CCErrCommParamsIsInvalid, metadata.AttributeFieldIsRequired)
			}
			if val != attr.IsRequired {
				return ctx.Error.Errorf(common.CCErrCoreServiceNotAllowChangeRequiredFieldErr, langObjID)
			}
		}
	}
	return nil
}

func (m *modelAttribute) getLangObjID(ctx core.ContextParams, objID string) string {
	langKey := "object_" + objID
	langObjID := ctx.Lang.Language(langKey)
	if langObjID == langKey {
		langObjID = objID
	}
	return langObjID
}

func (m *modelAttribute) buildUpdateAttrIndexReturn(ctx core.ContextParams, objID, propertyGroup string) (*metadata.UpdateAttrIndexData, error) {
	cond := mapstr.MapStr{
		common.BKObjIDField:         objID,
		common.BKPropertyGroupField: propertyGroup,
	}
	attrs := []metadata.Attribute{}
	err := m.dbProxy.Table(common.BKTableNameObjAttDes).Find(cond).All(ctx, &attrs)
	if nil != err {
		blog.Errorf("buildUpdateIndexReturn failed, request(%s): database operation is failed, error info is %s", ctx.ReqID, err.Error())
		return nil, err
	}

	count, err := m.dbProxy.Table(common.BKTableNameObjAttDes).Find(cond).Count(ctx)
	if nil != err {
		blog.Errorf("buildUpdateIndexReturn failed, request(%s): database operation is failed, error info is %s", ctx.ReqID, err.Error())
		return nil, err
	}
	info := make([]*metadata.UpdateAttributeIndex, 0)
	for _, attr := range attrs {
		idIndex := &metadata.UpdateAttributeIndex{
			Id:    attr.ID,
			Index: attr.PropertyIndex,
		}
		info = append(info, idIndex)
	}
	result := &metadata.UpdateAttrIndexData{
		Info:  info,
		Count: count,
	}

	return result, nil
}

func (m *modelAttribute) GetAttrLastIndex(ctx core.ContextParams, attribute metadata.Attribute) (int64, error) {
	opt := make(map[string]interface{})
	opt[common.BKObjIDField] = attribute.ObjectID
	opt[common.BKPropertyGroupField] = attribute.PropertyGroup
	opt = util.SetModOwner(opt, attribute.OwnerID)
	count, err := m.dbProxy.Table(common.BKTableNameObjAttDes).Find(opt).Count(ctx)
	if err != nil {
		blog.Error("GetAttrLastIndex, request(%s): database operation is failed, error info is %v", ctx.ReqID, err)
		return 0, ctx.Error.Error(common.CCErrCommDBSelectFailed)
	}
	if count <= 0 {
		return 0, nil
	}

	attrs := make([]metadata.Attribute, 0)
	sortCond := "-bk_property_index"
	if err := m.dbProxy.Table(common.BKTableNameObjAttDes).Find(opt).Sort(sortCond).Limit(1).All(ctx, &attrs); err != nil {
		blog.Error("GetAttrLastIndex, request(%s): database operation is failed, error info is %v", ctx.ReqID, err)
		return 0, ctx.Error.Error(common.CCErrCommDBSelectFailed)
	}

	if len(attrs) <= 0 {
		return 0, nil
	}
	return attrs[0].PropertyIndex + 1, nil
}
