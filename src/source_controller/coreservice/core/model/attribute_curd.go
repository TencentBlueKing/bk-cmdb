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
	"regexp"
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

	if err = m.checkAttributeMustNotEmpty(ctx, attribute); err != nil {
		return 0, err
	}
	if err = m.checkAttributeValidity(ctx, attribute); err != nil {
		return 0, err
	}

	// check name duplicate
	err = m.checkUnique(ctx, true, attribute.ObjectID, attribute.PropertyID, attribute.PropertyName)
	if err != nil {
		blog.ErrorJSON("save atttribute check unique err:%s, input:%s, rid:%s", err.Error(), attribute, ctx.ReqID)
		return 0, err
	}
	err = m.dbProxy.Table(common.BKTableNameObjAttDes).Insert(ctx, attribute)
	return id, err
}

func (m *modelAttribute) checkUnique(ctx core.ContextParams, isCreate bool, objID, propertyID, propertyName string) error {
	cond := mongo.NewCondition()
	cond = cond.Element(mongo.Field(common.BKObjIDField).Eq(objID))

	nameFieldCond := mongo.Field(common.BKPropertyNameField).Eq(propertyName)
	if isCreate {
		idFieldCond := mongo.Field(common.BKPropertyIDField).Eq(propertyID)
		cond = cond.Or(nameFieldCond, idFieldCond)
	} else {
		// update attribute. not change name, 无需判断
		if propertyName == "" {
			return nil
		}
		
		idFieldCond := mongo.Field(common.BKPropertyIDField).Neq(propertyID)
		cond = cond.Element(nameFieldCond, idFieldCond)
	}

	condMap := util.SetModOwner(cond.ToMapStr(), ctx.SupplierAccount)

	resultAttrs := []metadata.Attribute{}
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
	return nil
}

func (m *modelAttribute) checkAttributeValidity(ctx core.ContextParams, attribute metadata.Attribute) error {
	if common.AttributeIDMaxLength < utf8.RuneCountInString(attribute.PropertyID) {
		return ctx.Error.Errorf(common.CCErrCommValExceedMaxFailed, ctx.Lang.Language("model_attr_bk_property_id"), common.AttributeIDMaxLength)
	} else if attribute.PropertyID != "" {
		match, err := regexp.MatchString(`^[a-z\d_]+$`, attribute.PropertyID)
		if nil != err {
			return ctx.Error.Errorf(common.CCErrCommParamsIsInvalid, metadata.AttributeFieldPropertyID)
		}
		if !match {
			return ctx.Error.Errorf(common.CCErrCommParamsIsInvalid, metadata.AttributeFieldPropertyID)
		}
	}

	if common.AttributeNameMaxLength < utf8.RuneCountInString(attribute.PropertyName) {
		return ctx.Error.Errorf(common.CCErrCommValExceedMaxFailed, ctx.Lang.Language("model_attr_bk_property_name"), common.AttributeNameMaxLength)
	}

	if attribute.Placeholder != "" {
		if common.AttributePlaceHolderMaxLength < utf8.RuneCountInString(attribute.Placeholder) {
			return ctx.Error.Errorf(common.CCErrCommValExceedMaxFailed, ctx.Lang.Language("model_attr_placeholder"), common.AttributePlaceHolderMaxLength)
		}
	}

	if attribute.Unit != "" {
		if common.AttributeUnitMaxLength < utf8.RuneCountInString(attribute.Unit) {
			return ctx.Error.Errorf(common.CCErrCommValExceedMaxFailed, ctx.Lang.Language("model_attr_uint"), common.AttributeUnitMaxLength)
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

	dbAttributeArr, err := m.search(ctx, cond)
	if err != nil {
		blog.Errorf("request(%s): find nothing by the condition(%#v)  error(%s)", ctx.ReqID, cond.ToMapStr(), err.Error())
		return 0, err
	}
	if 0 == len(dbAttributeArr) {
		blog.Errorf("request(%s): find nothing by the condition(%#v)", ctx.ReqID, cond.ToMapStr())
		return cnt, nil
	}
	for _, dbAttribute := range dbAttributeArr {
		if dbAttribute.IsPre == true {
			blog.Warnf("try to update pre record, record: %+v, rid: %s", dbAttribute, ctx.ReqID)
			return 0, ctx.Error.CCError(common.CCErrCommOperateBuiltInItemForbidden)
		}
	}

	// 删除不可更新字段， 避免由于传入数据，修改字段
	// TODO: 改成白名单方式
	data.Remove(metadata.AttributeFieldPropertyID)
	data.Remove(metadata.AttributeFieldSupplierAccount)
	data.Set(metadata.AttributeFieldLastTime, time.Now())

	attribute := metadata.Attribute{}
	if err = data.MarshalJSONInto(&attribute); err != nil {
		blog.Errorf("request(%s): MarshalJSONInto(%#v), error is %v", ctx.ReqID, data, err)
		return 0, err
	}

	if err = m.checkAttributeValidity(ctx, attribute); err != nil {
		return 0, err
	}

<<<<<<< HEAD
	for _, dbAttribute := range dbAttributeArr {
		err = m.checkUnique(ctx, false, dbAttribute.ObjectID, dbAttribute.PropertyID, attribute.PropertyName)
		if err != nil {
			blog.ErrorJSON("save attribute check unique err:%s, input:%s, rid:%s", err.Error(), attribute, ctx.ReqID)
			return 0, err
		}
=======
	err = m.checkUnique(ctx, false, dbAttributeArr[0].ObjectID, dbAttributeArr[0].PropertyID, attribute.PropertyName)
	if err != nil {
		blog.ErrorJSON("save atttribute check unique err:%s, input:%s, rid:%s", err.Error(), attribute, ctx.ReqID)
		return 0, err
>>>>>>> fix:import model attribute bug issue #2978
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

func (m *modelAttribute) searchReturnMapStr(ctx core.ContextParams, cond universalsql.Condition) (resultAttrs []mapstr.MapStr, err error) {

	resultAttrs = []mapstr.MapStr{}
	err = m.dbProxy.Table(common.BKTableNameObjAttDes).Find(cond.ToMapStr()).All(ctx, &resultAttrs)
	return resultAttrs, err
}

func (m *modelAttribute) delete(ctx core.ContextParams, cond universalsql.Condition) (cnt uint64, err error) {

	cnt, err = m.dbProxy.Table(common.BKTableNameObjAttDes).Find(cond.ToMapStr()).Count(ctx)
	if nil != err {
		blog.Errorf("request(%s): database count operation is failed, error info is %s", ctx.ReqID, err.Error())
		return cnt, err
	}

	if 0 == cnt {
		return cnt, nil
	}

	err = m.dbProxy.Table(common.BKTableNameObjAttDes).Delete(ctx, cond.ToMapStr())
	if nil != err {
		blog.Errorf("request(%s): database deletion operation is failed, error info is %s", ctx.ReqID, err.Error())
		return 0, err
	}

	return cnt, err
}
