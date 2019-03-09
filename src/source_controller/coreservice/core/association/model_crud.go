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

package association

import (
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/universalsql"
	"configcenter/src/source_controller/coreservice/core"
)

func (m *associationModel) count(ctx core.ContextParams, cond universalsql.Condition) (cnt uint64, err error) {

	cnt, err = m.dbProxy.Table(common.BKTableNameObjAsst).Find(cond.ToMapStr()).Count(ctx)
	if nil != err {
		blog.Errorf("request(%s): it is failed to execute database count operation on the table (%s) by the condition (%#v), error info is %s", ctx.ReqID, common.BKTableNameObjAsst, cond.ToMapStr(), err.Error())
		return 0, err
	}
	return cnt, err
}

func (m *associationModel) isExists(ctx core.ContextParams, cond universalsql.Condition) (oneResult *metadata.Association, exists bool, err error) {

	oneResult = &metadata.Association{}
	err = m.dbProxy.Table(common.BKTableNameObjAsst).Find(cond.ToMapStr()).One(ctx, oneResult)
	if nil != err && !m.dbProxy.IsNotFoundError(err) {
		blog.Errorf("request(%s): it is faield to execute database findone operation on the table (%s) by the condition (%#v), error info is %s", ctx.ReqID, common.BKTableNameObjAsst, cond.ToMapStr(), err.Error())
		return oneResult, false, ctx.Error.New(common.CCErrObjectDBOpErrno, err.Error())
	}

	return oneResult, !m.dbProxy.IsNotFoundError(err), nil
}

func (m *associationModel) save(ctx core.ContextParams, assoParam *metadata.Association) (id uint64, err error) {

	id, err = m.dbProxy.NextSequence(ctx, common.BKTableNameObjAsst)
	if nil != err {
		blog.Errorf("request(%s): it is failed to make a sequence ID on the table (%s), error info is %s", ctx.ReqID, common.BKTableNameObjAsst, err.Error())
		return id, err
	}

	assoParam.ID = int64(id)
	assoParam.OwnerID = ctx.SupplierAccount
	err = m.dbProxy.Table(common.BKTableNameObjAsst).Insert(ctx, assoParam)
	if nil != err {
		blog.Errorf("request(%s): it is failed to execute database insert operation on the table (%s), error info is %s", ctx.ReqID, common.BKTableNameObjAsst, err.Error())
		return 0, err
	}
	return id, err
}

func (m *associationModel) update(ctx core.ContextParams, data mapstr.MapStr, cond universalsql.Condition) (cnt uint64, err error) {

	cnt, err = m.count(ctx, cond)
	if nil != err {
		return 0, err
	}

	if 0 >= cnt {
		return 0, err
	}

	err = m.dbProxy.Table(common.BKTableNameObjAsst).Update(ctx, cond.ToMapStr(), data)
	if nil != err {
		blog.Errorf("request(%s): it is failed to execute database upate some data (%v) on the table (%s) by the condition (%#v)", ctx.ReqID, data, common.BKTableNameObjAsst, cond.ToMapStr(), err.Error())
		return 0, err
	}
	return cnt, err
}

func (m *associationModel) delete(ctx core.ContextParams, cond universalsql.Condition) (cnt uint64, err error) {

	cnt, err = m.count(ctx, cond)
	if nil != err {
		return 0, err
	}

	if 0 >= cnt {
		return 0, err
	}

	err = m.dbProxy.Table(common.BKTableNameObjAsst).Delete(ctx, cond.ToMapStr())
	if nil != err {
		blog.Errorf("request(%s): it is to delete some data on the table (%s) by the condition (%#v), error info is %s", ctx.ReqID, common.BKTableNameObjAsst, cond.ToMapStr(), err.Error())
		return 0, err
	}
	return cnt, err
}

func (m *associationModel) search(ctx core.ContextParams, cond universalsql.Condition) ([]metadata.Association, error) {

	dataResult := []metadata.Association{}
	err := m.dbProxy.Table(common.BKTableNameObjAsst).Find(cond.ToMapStr()).All(ctx, &dataResult)
	if nil != err {
		blog.Errorf("request(%s): it is to search some data on the table (%s) by the condition (%v), error info is %s", ctx.ReqID, common.BKTableNameAsstDes, cond.ToMapStr(), err.Error())
		return dataResult, err
	}
	return dataResult, err
}

func (m *associationModel) searchReturnMapStr(ctx core.ContextParams, cond universalsql.Condition) ([]mapstr.MapStr, error) {
	dataResult := []mapstr.MapStr{}
	err := m.dbProxy.Table(common.BKTableNameObjAsst).Find(cond.ToMapStr()).All(ctx, &dataResult)
	if nil != err {
		blog.Errorf("request(%s): it is to search data on the table (%s) by the condition (%#v), error info is %s", ctx.ReqID, common.BKTableNameAsstDes, cond.ToMapStr(), err.Error())
		return dataResult, err
	}
	return dataResult, err
}
