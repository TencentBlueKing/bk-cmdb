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
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/universalsql"
	"configcenter/src/common/universalsql/mongo"
	"configcenter/src/common/util"
	"configcenter/src/storage/dal/types"
	"configcenter/src/storage/driver/mongodb"
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

func (m *modelAttribute) Count(kit *rest.Kit, cond universalsql.Condition) (cnt uint64, err error) {
	cnt, err = mongodb.Client().Table(common.BKTableNameObjAttDes).Find(cond.ToMapStr()).Count(kit.Ctx)
	return cnt, err
}

func (m *modelAttribute) save(kit *rest.Kit, attribute metadata.Attribute) (id uint64, err error) {

	id, err = mongodb.Client().NextSequence(kit.Ctx, common.BKTableNameObjAttDes)
	if err != nil {
		return id, kit.CCError.New(common.CCErrObjectDBOpErrno, err.Error())
	}

	index, err := m.GetAttrLastIndex(kit, attribute)
	if err != nil {
		return id, err
	}

	attribute.PropertyIndex = index
	attribute.ID = int64(id)
	attribute.OwnerID = kit.SupplierAccount

	if nil == attribute.CreateTime {
		attribute.CreateTime = &metadata.Time{}
		attribute.CreateTime.Time = time.Now()
	}

	if nil == attribute.LastTime {
		attribute.LastTime = &metadata.Time{}
		attribute.LastTime.Time = time.Now()
	}

	if err = m.saveCheck(kit, attribute); err != nil {
		return 0, err
	}

	err = mongodb.Client().Table(common.BKTableNameObjAttDes).Insert(kit.Ctx, attribute)
	return id, err
}

func (m *modelAttribute) checkUnique(kit *rest.Kit, isCreate bool, objID, propertyID, propertyName string, modelBizID int64) error {
	cond := map[string]interface{}{
		common.BKObjIDField: objID,
	}

	andCond := make([]map[string]interface{}, 0)
	if isCreate {
		nameFieldCond := map[string]interface{}{common.BKPropertyNameField: propertyName}
		idFieldCond := map[string]interface{}{common.BKPropertyIDField: propertyID}
		andCond = append(andCond, map[string]interface{}{
			common.BKDBOR: []map[string]interface{}{nameFieldCond, idFieldCond},
		})
	} else {
		// update attribute. not change name, 无需判断
		if propertyName == "" {
			return nil
		}
		cond[common.BKPropertyIDField] = map[string]interface{}{common.BKDBNE: propertyID}
		cond[common.BKPropertyNameField] = propertyName
	}

	if modelBizID > 0 {
		// search special business model and global shared model
		andCond = append(andCond, map[string]interface{}{
			common.BKDBOR: []map[string]interface{}{
				{common.BKAppIDField: modelBizID},
				{common.BKAppIDField: 0},
				{common.BKAppIDField: map[string]interface{}{common.BKDBExists: false}},
			},
		})
	}

	if len(andCond) > 0 {
		cond[common.BKDBAND] = andCond
	}
	util.SetModOwner(cond, kit.SupplierAccount)

	resultAttrs := make([]metadata.Attribute, 0)
	err := mongodb.Client().Table(common.BKTableNameObjAttDes).Find(cond).All(kit.Ctx, &resultAttrs)
	blog.V(5).Infof("checkUnique db cond:%#v, result:%#v, rid:%s", cond, resultAttrs, kit.Rid)
	if err != nil {
		blog.ErrorJSON("checkUnique select error. err:%s, cond:%s, rid:%s", err.Error(), cond, kit.Rid)
		return kit.CCError.Error(common.CCErrCommDBSelectFailed)
	}

	language := util.GetLanguage(kit.Header)
	lang := m.language.CreateDefaultCCLanguageIf(language)
	for _, attrItem := range resultAttrs {
		if attrItem.PropertyID == propertyID {
			return kit.CCError.Errorf(common.CCErrCommDuplicateItem, lang.Language("model_attr_bk_property_id"))
		}
		if attrItem.PropertyName == propertyName {
			return kit.CCError.Errorf(common.CCErrCommDuplicateItem, lang.Language("model_attr_bk_property_name"))
		}
	}

	return nil
}

func (m *modelAttribute) checkAttributeMustNotEmpty(kit *rest.Kit, attribute metadata.Attribute) error {
	if attribute.PropertyID == "" {
		return kit.CCError.Errorf(common.CCErrCommParamsNeedSet, metadata.AttributeFieldPropertyID)
	}
	if attribute.PropertyName == "" {
		return kit.CCError.Errorf(common.CCErrCommParamsNeedSet, metadata.AttributeFieldPropertyName)
	}
	if attribute.PropertyType == "" {
		return kit.CCError.Errorf(common.CCErrCommParamsNeedSet, metadata.AttributeFieldPropertyType)
	}
	return nil
}

func (m *modelAttribute) checkAttributeValidity(kit *rest.Kit, attribute metadata.Attribute) error {
	language := util.GetLanguage(kit.Header)
	lang := m.language.CreateDefaultCCLanguageIf(language)
	if attribute.PropertyID != "" {
		attribute.PropertyID = strings.TrimSpace(attribute.PropertyID)
		if common.AttributeIDMaxLength < utf8.RuneCountInString(attribute.PropertyID) {
			return kit.CCError.Errorf(common.CCErrCommValExceedMaxFailed, lang.Language("model_attr_bk_property_id"), common.AttributeIDMaxLength)
		}

		if !SatisfyMongoFieldLimit(attribute.PropertyID) {
			blog.Errorf("attribute.PropertyID:%s not SatisfyMongoFieldLimit", attribute.PropertyID)
			return kit.CCError.Errorf(common.CCErrCommParamsIsInvalid, metadata.AttributeFieldPropertyID)
		}

		if !attribute.IsPre {
			if strings.HasPrefix(attribute.PropertyID, "bk_") || strings.HasPrefix(attribute.PropertyID, "_bk") {
				return kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, metadata.AttributeFieldPropertyID)
			}
		}
	}

	if attribute.PropertyName = strings.TrimSpace(attribute.PropertyName); common.AttributeNameMaxLength < utf8.RuneCountInString(attribute.PropertyName) {
		return kit.CCError.Errorf(common.CCErrCommValExceedMaxFailed, lang.Language("model_attr_bk_property_name"), common.AttributeNameMaxLength)
	}

	if attribute.Placeholder != "" {
		attribute.Placeholder = strings.TrimSpace(attribute.Placeholder)

		if common.AttributePlaceHolderMaxLength < utf8.RuneCountInString(attribute.Placeholder) {
			return kit.CCError.Errorf(common.CCErrCommValExceedMaxFailed, lang.Language("model_attr_placeholder"), common.AttributePlaceHolderMaxLength)
		}
		match, err := regexp.MatchString(common.FieldTypeLongCharRegexp, attribute.Placeholder)
		if nil != err || !match {
			return kit.CCError.Errorf(common.CCErrCommParamsIsInvalid, metadata.AttributeFieldPlaceHolder)

		}
	}

	if attribute.Unit != "" {
		attribute.Unit = strings.TrimSpace(attribute.Unit)
		if common.AttributeUnitMaxLength < utf8.RuneCountInString(attribute.Unit) {
			return kit.CCError.Errorf(common.CCErrCommValExceedMaxFailed, lang.Language("model_attr_uint"), common.AttributeUnitMaxLength)
		}
	}

	if attribute.PropertyType != "" {
		switch attribute.PropertyType {
		case common.FieldTypeSingleChar, common.FieldTypeLongChar, common.FieldTypeInt, common.FieldTypeFloat, common.FieldTypeEnum,
			common.FieldTypeDate, common.FieldTypeTime, common.FieldTypeUser, common.FieldTypeOrganization, common.FieldTypeTimeZone, common.FieldTypeBool, common.FieldTypeList:
		default:
			return kit.CCError.Errorf(common.CCErrCommParamsIsInvalid, metadata.AttributeFieldPropertyType)
		}
	}

	if opt, ok := attribute.Option.(string); ok && opt != "" {
		if common.AttributeOptionMaxLength < utf8.RuneCountInString(opt) {
			return kit.CCError.Errorf(common.CCErrCommValExceedMaxFailed, lang.Language("model_attr_option_regex"), common.AttributeOptionMaxLength)
		}
	}

	return nil
}

func (m *modelAttribute) update(kit *rest.Kit, data mapstr.MapStr, cond universalsql.Condition) (cnt uint64, err error) {
	cnt, err = m.checkUpdate(kit, data, cond)
	if err != nil {
		blog.ErrorJSON("checkUpdate error. data:%s, cond:%s, rid:%s", data, cond, kit.Rid)
		return cnt, err
	}
	err = mongodb.Client().Table(common.BKTableNameObjAttDes).Update(kit.Ctx, cond.ToMapStr(), data)
	if nil != err {
		blog.Errorf("request(%s): database operation is failed, error info is %s", kit.Rid, err.Error())
		return 0, err
	}

	return cnt, err
}

func (m *modelAttribute) newSearch(kit *rest.Kit, cond mapstr.MapStr) (resultAttrs []metadata.Attribute, err error) {
	resultAttrs = []metadata.Attribute{}
	err = mongodb.Client().Table(common.BKTableNameObjAttDes).Find(cond).All(kit.Ctx, &resultAttrs)
	return resultAttrs, err
}

func (m *modelAttribute) search(kit *rest.Kit, cond universalsql.Condition) (resultAttrs []metadata.Attribute, err error) {
	resultAttrs = []metadata.Attribute{}
	err = mongodb.Client().Table(common.BKTableNameObjAttDes).Find(cond.ToMapStr()).All(kit.Ctx, &resultAttrs)
	return resultAttrs, err
}

func (m *modelAttribute) searchWithSort(kit *rest.Kit, cond metadata.QueryCondition) (resultAttrs []metadata.Attribute, err error) {
	resultAttrs = []metadata.Attribute{}

	instHandler := mongodb.Client().Table(common.BKTableNameObjAttDes).Find(cond.Condition)
	err = instHandler.Start(uint64(cond.Page.Start)).Limit(uint64(cond.Page.Limit)).Sort(cond.Page.Sort).All(kit.Ctx, &resultAttrs)

	return resultAttrs, err
}

func (m *modelAttribute) searchReturnMapStr(kit *rest.Kit, cond universalsql.Condition) (resultAttrs []mapstr.MapStr, err error) {

	resultAttrs = []mapstr.MapStr{}
	err = mongodb.Client().Table(common.BKTableNameObjAttDes).Find(cond.ToMapStr()).All(kit.Ctx, &resultAttrs)
	return resultAttrs, err
}

func (m *modelAttribute) delete(kit *rest.Kit, cond universalsql.Condition) (cnt uint64, err error) {

	resultAttrs := make([]metadata.Attribute, 0)
	fields := []string{common.BKFieldID, common.BKPropertyIDField, common.BKObjIDField, common.BKAppIDField}

	condMap := util.SetQueryOwner(cond.ToMapStr(), kit.SupplierAccount)
	err = mongodb.Client().Table(common.BKTableNameObjAttDes).Find(condMap).Fields(fields...).All(kit.Ctx, &resultAttrs)
	if nil != err {
		blog.Errorf("request(%s): database count operation is failed, error info is %s", kit.Rid, err.Error())
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

	if err := m.cleanAttributeFieldInInstances(kit.Ctx, kit.SupplierAccount, resultAttrs); err != nil {
		blog.ErrorJSON("delete object attributes with cond: %s, but delete these attribute in instance failed, err: %v, rid:%s", condMap, err, kit.Rid)
		return 0, err
	}

	exist, err := m.checkAttributeInUnique(kit, objIDArrMap)
	if err != nil {
		blog.ErrorJSON("check attribute in unique error. err:%s, input:%s, rid:%s", err.Error(), condMap, kit.Rid)
		return 0, err
	}
	// delete field in module unique. not allow delete
	if exist {
		blog.ErrorJSON("delete field in unique. delete cond:%s, field:%s, rid:%s", condMap, resultAttrs, kit.Rid)
		return 0, kit.CCError.Error(common.CCErrCoreServiceNotAllowUniqueAttr)
	}

	err = mongodb.Client().Table(common.BKTableNameObjAttDes).Delete(kit.Ctx, condMap)
	if nil != err {
		blog.Errorf("request(%s): database deletion operation is failed, error info is %s", kit.Rid, err.Error())
		return 0, err
	}

	return cnt, err
}

type bizObjectFields struct {
	bizID  int64
	fields []string
}

// remove attribute filed in this object's instances
func (m *modelAttribute) cleanAttributeFieldInInstances(ctx context.Context, ownerID string, attrs []metadata.Attribute) error {
	// this operation may take a long time, do not use transaction
	ctx = context.Background()

	objectFields := make(map[string][]bizObjectFields, 0)
	hostApplyFields := make(map[int64][]int64)

	// TODO: now, we only support set, module, host model's biz attribute clean operation.
	for _, attr := range attrs {
		biz := attr.BizID
		if biz != 0 {
			if !isBizObject(attr.ObjectID) {
				return fmt.Errorf("unsupported object %s's clean instance field operation", attr.ObjectID)
			}
		}

		_, exist := objectFields[attr.ObjectID]
		if !exist {
			objectFields[attr.ObjectID] = make([]bizObjectFields, 0)
		}
		objectFields[attr.ObjectID] = append(objectFields[attr.ObjectID], bizObjectFields{
			bizID:  biz,
			fields: []string{attr.PropertyID},
		})

		if attr.ObjectID == common.BKInnerObjIDHost {
			hostApplyFields[biz] = append(hostApplyFields[biz], attr.ID)
		}
	}

	// delete these attribute's fields in the model instance
	var hitError error
	wg := sync.WaitGroup{}
	for object, objFields := range objectFields {
		if len(objFields) == 0 {
			// no fields need to be removed, skip directly.
			continue
		}

		for _, objField := range objFields {
			fields := objField.fields
			existConds := make([]map[string]interface{}, len(fields))

			for index, field := range fields {
				existConds[index] = map[string]interface{}{
					field: map[string]interface{}{
						common.BKDBExists: true,
					},
				}
			}

			cond := map[string]interface{}{
				common.BKDBOR: existConds,
			}

			if objField.bizID > 0 {
				if !isBizObject(object) {
					return fmt.Errorf("unsupported object %s's clean instance field operation", object)
				}

				if object == common.BKInnerObjIDHost {
					if err := m.cleanHostAttributeField(ctx, ownerID, objField); err != nil {
						return err
					}
					continue
				}

				cond[common.BKAppIDField] = objField.bizID
			} else {
				if isBizObject(object) {
					if object == common.BKInnerObjIDHost {
						ele := bizObjectFields{
							bizID:  0,
							fields: fields,
						}
						if err := m.cleanHostAttributeField(ctx, ownerID, ele); err != nil {
							return err
						}
						continue
					}
				} else {
					cond[common.BKObjIDField] = object
				}
			}

			cond = util.SetQueryOwner(cond, ownerID)

			collectionName := common.GetInstTableName(object)
			wg.Add(1)
			go func(collName string, filter types.Filter, fields []string) {
				defer wg.Done()

				instCount, err := mongodb.Client().Table(collName).Find(filter).Count(ctx)
				if err != nil {
					blog.Error("count instances with the attribute to delete failed, table: %s, cond: %v, fields: %v, err: %v", collectionName, filter, fields, err)
					hitError = err
					return
				}

				instIDField := common.GetInstIDField(object)
				for start := uint64(0); start < instCount; start += pageSize {
					insts := make([]map[string]interface{}, 0)
					err := mongodb.Client().Table(collName).Find(filter).Start(0).Limit(pageSize).Fields(instIDField).All(ctx, &insts)
					if err != nil {
						blog.Error("get instance ids with the attribute to delete failed, table: %s, cond: %v, fields: %v, err: %v", collectionName, filter, fields, err)
						hitError = err
						return
					}

					if len(insts) == 0 {
						return
					}

					instIDs := make([]int64, len(insts))
					for index, inst := range insts {
						instID, err := util.GetInt64ByInterface(inst[instIDField])
						if err != nil {
							blog.Error("get instance id failed, inst: %+v, err: %v", inst, err)
							hitError = err
							return
						}
						instIDs[index] = instID
					}

					instFilter := map[string]interface{}{
						instIDField: map[string]interface{}{
							common.BKDBIN: instIDs,
						},
					}

					if err := mongodb.Client().Table(collName).DropColumns(ctx, instFilter, fields); err != nil {
						blog.Error("delete object's attribute from instance failed, table: %s, cond: %v, fields: %v, err: %v", collectionName, instFilter, fields, err)
						hitError = err
						return
					}
				}
			}(collectionName, cond, fields)

		}
	}
	// wait for all the public object routine is done.
	wg.Wait()
	if hitError != nil {
		return hitError
	}

	// wait for all the public object routine is done.
	wg.Wait()
	if hitError != nil {
		return hitError
	}

	// step 3: clean host apply fields
	if err := m.cleanHostApplyField(ctx, ownerID, hostApplyFields); err != nil {
		return err
	}

	return nil
}

const pageSize = 2000

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

	hostCount, err := mongodb.Client().Table(common.BKTableNameModuleHostConfig).Find(cond).Count(ctx)
	if err != nil {
		return err
	}

	type hostInst struct {
		HostID int64 `bson:"bk_host_id"`
	}

	fields := info.fields
	existConds := make([]map[string]interface{}, len(fields))

	for index, field := range fields {
		existConds[index] = map[string]interface{}{
			field: map[string]interface{}{
				common.BKDBExists: true,
			},
		}
	}

	for start := uint64(0); start < hostCount; start += pageSize {
		hostList := make([]hostInst, 0)
		err := mongodb.Client().Table(common.BKTableNameModuleHostConfig).Find(cond).Start(start).Limit(pageSize).Fields(common.BKHostIDField).All(ctx, &hostList)
		if err != nil {
			return err
		}

		if len(hostList) == 0 {
			return nil
		}

		ids := make([]int64, len(hostList))
		for index, host := range hostList {
			ids[index] = host.HostID
		}

		hostFilter := mapstr.MapStr{
			common.BKHostIDField: mapstr.MapStr{common.BKDBIN: ids},
			common.BKDBOR:        existConds,
		}
		if err := mongodb.Client().Table(common.BKTableNameBaseHost).DropColumns(ctx, hostFilter, info.fields); err != nil {
			return fmt.Errorf("clean host biz attribute %v failed, err: %v", info.fields, err)
		}
	}

	return nil

}

func (m *modelAttribute) cleanHostApplyField(ctx context.Context, ownerID string, hostApplyFields map[int64][]int64) error {
	orCond := make([]map[string]interface{}, 0)
	for bizID, attrIDs := range hostApplyFields {
		attrCond := map[string]interface{}{
			common.BKAttributeIDField: map[string]interface{}{
				common.BKDBIN: attrIDs,
			},
		}
		// global attribute requires removing host apply rules for all biz
		if bizID != 0 {
			attrCond[common.BKAppIDField] = bizID
		}
		orCond = append(orCond, attrCond)
	}
	if len(orCond) == 0 {
		return nil
	}
	cond := make(map[string]interface{})
	cond = util.SetQueryOwner(cond, ownerID)
	cond[common.BKDBOR] = orCond
	if err := mongodb.Client().Table(common.BKTableNameHostApplyRule).Delete(ctx, cond); err != nil {
		blog.ErrorJSON("cleanHostApplyField failed, err: %s, cond: %s", err, cond)
		return err
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
func (m *modelAttribute) saveCheck(kit *rest.Kit, attribute metadata.Attribute) error {

	if err := m.checkAddField(kit, attribute); err != nil {
		return err
	}

	if err := m.checkAttributeMustNotEmpty(kit, attribute); err != nil {
		return err
	}
	if err := m.checkAttributeValidity(kit, attribute); err != nil {
		return err
	}

	// check name duplicate
	if err := m.checkUnique(kit, true, attribute.ObjectID, attribute.PropertyID, attribute.PropertyName, attribute.BizID); err != nil {
		blog.ErrorJSON("save attribute check unique err:%s, input:%s, rid:%s", err.Error(), attribute, kit.Rid)
		return err
	}

	return nil
}

// checkUpdate 删除不可以更新字段，检验字段是否重复， 返回更新的行数，错误
func (m *modelAttribute) checkUpdate(kit *rest.Kit, data mapstr.MapStr, cond universalsql.Condition) (changeRow uint64, err error) {

	dbAttributeArr, err := m.search(kit, cond)
	if err != nil {
		blog.Errorf("request(%s): find nothing by the condition(%#v)  error(%s)", kit.Rid, cond.ToMapStr(), err.Error())
		return changeRow, err
	}
	if 0 == len(dbAttributeArr) {
		blog.Errorf("request(%s): find nothing by the condition(%#v)", kit.Rid, cond.ToMapStr())
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
		_ = data.ForEach(func(key string, val interface{}) error {
			if key != metadata.AttributeFieldPropertyGroup &&
				key != metadata.AttributeFieldPropertyIndex &&
				key != metadata.AttributeFieldPropertyName &&
				key != metadata.AttributeFieldUnit &&
				key != metadata.AttributeFieldPlaceHolder &&
				key != metadata.AttributeFieldOption {
				data.Remove(key)
			}
			return nil
		})
	}

	if option, exists := data.Get(metadata.AttributeFieldOption); exists {
		propertyType := dbAttributeArr[0].PropertyType
		for _, dbAttribute := range dbAttributeArr {
			if dbAttribute.PropertyType != propertyType {
				blog.ErrorJSON("update option, but property type not the same, db attributes: %s, rid:%s", dbAttributeArr, kit.Ctx)
				return changeRow, kit.CCError.Errorf(common.CCErrCommParamsInvalid, "cond")
			}
		}
		if err := util.ValidPropertyOption(propertyType, option, kit.CCError); err != nil {
			blog.ErrorJSON("valid property option failed, err: %s, data: %s, rid:%s", err, data, kit.Ctx)
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
		cnt, err := mongodb.Client().Table(common.BKTableNamePropertyGroup).Find(cond).Count(kit.Ctx)
		if err != nil {
			blog.ErrorJSON("property group count failed, err: %s, condition: %s, rid: %s", err, cond, kit.Rid)
			return changeRow, err
		}
		if cnt != uint64(len(objIDs)) {
			blog.Errorf("property group invalid, objIDs: %s have %d property groups, rid: %s", objIDs, cnt, kit.Rid)
			return changeRow, kit.CCError.Errorf(common.CCErrCommParamsInvalid, metadata.AttributeFieldPropertyGroup)
		}
	}

	attribute := metadata.Attribute{}
	if err = data.MarshalJSONInto(&attribute); err != nil {
		blog.Errorf("request(%s): MarshalJSONInto(%#v), error is %v", kit.Rid, data, err)
		return changeRow, err
	}

	if err = m.checkAttributeValidity(kit, attribute); err != nil {
		return changeRow, err
	}

	for _, dbAttribute := range dbAttributeArr {
		err = m.checkUnique(kit, false, dbAttribute.ObjectID, dbAttribute.PropertyID, attribute.PropertyName, attribute.BizID)
		if err != nil {
			blog.ErrorJSON("save attribute check unique err:%s, input:%s, rid:%s", err.Error(), attribute, kit.Rid)
			return changeRow, err
		}
		if err = m.checkChangeField(kit, dbAttribute, data); err != nil {
			return changeRow, err
		}
	}

	return uint64(len(dbAttributeArr)), err

}

// checkAttributeInUnique 检查属性是否存在唯一校验中  objIDPropertyIDArr  属性的bk_obj_id和表中ID的集合
func (m *modelAttribute) checkAttributeInUnique(kit *rest.Kit, objIDPropertyIDArr map[string][]int64) (bool, error) {

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
	condMap := util.SetQueryOwner(cond.ToMapStr(), kit.SupplierAccount)

	cnt, err := mongodb.Client().Table(common.BKTableNameObjUnique).Find(condMap).Count(kit.Ctx)
	if err != nil {
		blog.ErrorJSON("checkAttributeInUnique db select error. err:%s, cond:%s, rid:%s", err.Error(), condMap, kit.Rid)
		return false, kit.CCError.Error(common.CCErrCommDBSelectFailed)
	}

	if cnt > 0 {
		return true, nil
	}

	return false, nil
}

// checkAddRequireField 新加模型属性的时候，如果新加的是必填字段，需要判断是否可以新加必填字段
func (m *modelAttribute) checkAddField(kit *rest.Kit, attribute metadata.Attribute) error {
	langObjID := m.getLangObjID(kit, attribute.ObjectID)
	if _, ok := notAddAttrModel[attribute.ObjectID]; ok {
		//  不允许新加字段的模型
		return kit.CCError.Errorf(common.CCErrCoreServiceNotAllowAddFieldErr, langObjID)
	}

	if _, ok := RequiredFieldUnchangeableModels[attribute.ObjectID]; ok {
		if attribute.IsRequired {
			//  不允许修改必填字段的模型
			return kit.CCError.Errorf(common.CCErrCoreServiceNotAllowAddRequiredFieldErr, langObjID)
		}

	}
	return nil
}

// 修改模型属性的时候，如果修改的属性包含是否为必填字段(isrequired)，需要判断该模型的必填字段是否允许被修改
func (m *modelAttribute) checkChangeField(kit *rest.Kit, attr metadata.Attribute, attrInfo mapstr.MapStr) error {
	langObjID := m.getLangObjID(kit, attr.ObjectID)
	if _, ok := RequiredFieldUnchangeableModels[attr.ObjectID]; ok {
		if attrInfo.Exists(metadata.AttributeFieldIsRequired) {
			// 不允许修改模型的必填字段
			val, ok := attrInfo[metadata.AttributeFieldIsRequired].(bool)
			if !ok {
				return kit.CCError.Errorf(common.CCErrCommParamsIsInvalid, metadata.AttributeFieldIsRequired)
			}
			if val != attr.IsRequired {
				return kit.CCError.Errorf(common.CCErrCoreServiceNotAllowChangeRequiredFieldErr, langObjID)
			}
		}
	}
	return nil
}

func (m *modelAttribute) getLangObjID(kit *rest.Kit, objID string) string {
	langKey := "object_" + objID
	language := util.GetLanguage(kit.Header)
	lang := m.language.CreateDefaultCCLanguageIf(language)
	langObjID := lang.Language(langKey)
	if langObjID == langKey {
		langObjID = objID
	}
	return langObjID
}

func (m *modelAttribute) buildUpdateAttrIndexReturn(kit *rest.Kit, objID, propertyGroup string) (*metadata.UpdateAttrIndexData, error) {
	cond := mapstr.MapStr{
		common.BKObjIDField:         objID,
		common.BKPropertyGroupField: propertyGroup,
	}
	attrs := []metadata.Attribute{}
	err := mongodb.Client().Table(common.BKTableNameObjAttDes).Find(cond).All(kit.Ctx, &attrs)
	if nil != err {
		blog.Errorf("buildUpdateIndexReturn failed, request(%s): database operation is failed, error info is %s", kit.Rid, err.Error())
		return nil, err
	}

	count, err := mongodb.Client().Table(common.BKTableNameObjAttDes).Find(cond).Count(kit.Ctx)
	if nil != err {
		blog.Errorf("buildUpdateIndexReturn failed, request(%s): database operation is failed, error info is %s", kit.Rid, err.Error())
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

func (m *modelAttribute) GetAttrLastIndex(kit *rest.Kit, attribute metadata.Attribute) (int64, error) {
	opt := make(map[string]interface{})
	opt[common.BKObjIDField] = attribute.ObjectID
	opt[common.BKPropertyGroupField] = attribute.PropertyGroup
	opt = util.SetModOwner(opt, attribute.OwnerID)
	count, err := mongodb.Client().Table(common.BKTableNameObjAttDes).Find(opt).Count(kit.Ctx)
	if err != nil {
		blog.Error("GetAttrLastIndex, request(%s): database operation is failed, error info is %v", kit.Rid, err)
		return 0, kit.CCError.Error(common.CCErrCommDBSelectFailed)
	}
	if count <= 0 {
		return 0, nil
	}

	attrs := make([]metadata.Attribute, 0)
	sortCond := "-bk_property_index"
	if err := mongodb.Client().Table(common.BKTableNameObjAttDes).Find(opt).Sort(sortCond).Limit(1).All(kit.Ctx, &attrs); err != nil {
		blog.Error("GetAttrLastIndex, database operation is failed, err: %v, rid: %s", err, kit.Rid)
		return 0, kit.CCError.Error(common.CCErrCommDBSelectFailed)
	}

	if len(attrs) <= 0 {
		return 0, nil
	}
	return attrs[0].PropertyIndex + 1, nil
}
