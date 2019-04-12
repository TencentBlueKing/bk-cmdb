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

func (m *modelManager) count(ctx core.ContextParams, cond universalsql.Condition) (uint64, error) {

	cnt, err := m.dbProxy.Table(common.BKTableNameObjDes).Find(cond.ToMapStr()).Count(ctx)
	if nil != err {
		blog.Errorf("request(%s): it is failed to execute database count operation by the condition (%#v), error info is %s", ctx.ReqID, cond.ToMapStr(), err.Error())
		return 0, ctx.Error.Errorf(common.CCErrObjectDBOpErrno, err.Error())
	}

	return cnt, err
}

func (m *modelManager) save(ctx core.ContextParams, model *metadata.Object) (id uint64, err error) {

	id, err = m.dbProxy.NextSequence(ctx, common.BKTableNameObjDes)
	if err != nil {
		blog.Errorf("request(%s): it is failed to make sequence id on the table (%s), error info is %s", ctx.ReqID, common.BKTableNameObjDes, err.Error())
		return id, ctx.Error.New(common.CCErrObjectDBOpErrno, err.Error())
	}
	model.ID = int64(id)
	model.OwnerID = ctx.SupplierAccount

	if nil == model.LastTime {
		model.LastTime = &metadata.Time{}
		model.LastTime.Time = time.Now()
	}
	if nil == model.CreateTime {
		model.CreateTime = &metadata.Time{}
		model.CreateTime.Time = time.Now()
	}

	err = m.dbProxy.Table(common.BKTableNameObjDes).Insert(ctx, model)
	return id, err
}

func (m *modelManager) update(ctx core.ContextParams, data mapstr.MapStr, cond universalsql.Condition) (cnt uint64, err error) {

	cnt, err = m.count(ctx, cond)
	if nil != err {
		return 0, err
	}

	if 0 == cnt {
		return 0, nil
	}

	data.Set(metadata.ModelFieldLastTime, time.Now())

	err = m.dbProxy.Table(common.BKTableNameObjDes).Update(ctx, cond.ToMapStr(), data)
	if nil != err {
		blog.Errorf("request(%s): it is failed to execute database update operation on the table (%s), error info is %s", ctx.ReqID, common.BKTableNameObjDes, err.Error())
		return 0, ctx.Error.New(common.CCErrObjectDBOpErrno, err.Error())
	}

	return cnt, err
}

func (m *modelManager) search(ctx core.ContextParams, cond universalsql.Condition) ([]metadata.Object, error) {

	dataResult := []metadata.Object{}
	if err := m.dbProxy.Table(common.BKTableNameObjDes).Find(cond.ToMapStr()).All(ctx, &dataResult); nil != err {
		blog.Errorf("request(%s): it is failed to find all models by the condition (%#v), error info is %s", ctx.ReqID, cond.ToMapStr(), err.Error())
		return dataResult, ctx.Error.New(common.CCErrObjectDBOpErrno, err.Error())
	}

	return dataResult, nil
}

func (m *modelManager) searchReturnMapStr(ctx core.ContextParams, cond universalsql.Condition) ([]mapstr.MapStr, error) {

	dataResult := []mapstr.MapStr{}
	if err := m.dbProxy.Table(common.BKTableNameObjDes).Find(cond.ToMapStr()).All(ctx, &dataResult); nil != err {
		blog.Errorf("request(%s): it is failed to find all models by the condition (%#v), error info is %s", ctx.ReqID, cond.ToMapStr(), err.Error())
		return dataResult, ctx.Error.New(common.CCErrObjectDBOpErrno, err.Error())
	}

	return dataResult, nil
}

func (m *modelManager) delete(ctx core.ContextParams, cond universalsql.Condition) (uint64, error) {

	cnt, err := m.count(ctx, cond)
	if nil != err {
		return 0, err
	}

	if 0 == cnt {
		return 0, nil
	}

	if err = m.dbProxy.Table(common.BKTableNameObjDes).Delete(ctx, cond.ToMapStr()); nil != err {
		blog.Errorf("request(%s): it is failed to execute a deletion operation on the table (%s), error info is %s", ctx.ReqID, common.BKTableNameObjDes, err.Error())
		return 0, ctx.Error.New(common.CCErrObjectDBOpErrno, err.Error())
	}

	return cnt, nil
}

func (m *modelManager) cascadeDelete(ctx core.ContextParams, cond universalsql.Condition) (uint64, error) {

	modelItems, err := m.search(ctx, cond)
	if nil != err {
		blog.Errorf("request(%s): it is failed to execute a cascade model deletion operation by the condition (%#v), error info is %s", ctx.ReqID, cond.ToMapStr(), err.Error())
		return 0, err
	}

	targetObjIDS := []string{}
	for _, modelItem := range modelItems {
		targetObjIDS = append(targetObjIDS, modelItem.ObjectID)
	}

	// cascade delete the other resource
	if err := m.dependent.CascadeDeleteAssociation(ctx, targetObjIDS); nil != err {
		blog.Errorf("request(%s): it is failed to execute a cascade model association deletion operation by the modelIDS(%#v), error info is %s", ctx.ReqID, targetObjIDS, err.Error())
		return 0, err
	}

	cnt, err := m.deleteModelAndAttributes(ctx, targetObjIDS)
	if nil != err {
		blog.Errorf("request(%s): it is failed to delete the models (%#v) and the model's attributes ,error info is %s", ctx.ReqID, targetObjIDS, err.Error())
		return 0, ctx.Error.New(common.CCErrObjectDBOpErrno, err.Error())
	}

	return cnt, nil

}
