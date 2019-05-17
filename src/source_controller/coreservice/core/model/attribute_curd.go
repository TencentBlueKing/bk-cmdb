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
	"time"

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

	err = m.dbProxy.Table(common.BKTableNameObjAttDes).Insert(ctx, attribute)
	return id, err
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
