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

	err = m.dbProxy.Table(common.BKTableNameObjAttDes).Insert(ctx, attribute)
	return id, err
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
		return ctx.Error.Errorf(common.CCErrCommOverLimit, attribute.PropertyID)
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
		return ctx.Error.Errorf(common.CCErrCommOverLimit, attribute.PropertyName)
	}

	if attribute.Placeholder != "" {
		if common.AttributePlaceHolderMaxLength < utf8.RuneCountInString(attribute.Placeholder) {
			return ctx.Error.Errorf(common.CCErrCommOverLimit, attribute.Placeholder)
		}
	}

	if attribute.Unit != "" {
		if 20 < utf8.RuneCountInString(attribute.Unit) {
			return ctx.Error.Errorf(common.CCErrCommOverLimit, attribute.Unit)
		}
	}

	if opt, ok := attribute.Option.(string); ok && opt != "" {
		if common.AttributeOptionMaxLength < utf8.RuneCountInString(opt) {
			return ctx.Error.Errorf(common.CCErrCommOverLimit, opt)
		}
	}

	return nil
}

func (m *modelAttribute) update(ctx core.ContextParams, data mapstr.MapStr, cond universalsql.Condition) (cnt uint64, err error) {

	cnt, err = m.count(ctx, cond)
	if 0 == cnt {
		blog.Errorf("request(%s): find nothing by the condition(%#v)", ctx.ReqID, cond.ToMapStr())
		return cnt, nil
	}

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
