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
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/universalsql"
	"configcenter/src/source_controller/coreservice/core"
)

func (m *modelClassification) count(ctx core.ContextParams, cond universalsql.Condition) (cnt uint64, err error) {

	cnt, err = m.dbProxy.Table(common.BKTableNameObjClassifiction).Find(cond.ToMapStr()).Count(ctx)
	if nil != err {
		blog.Errorf("request(%s): it is failed to execute a database count operation on the table(%s) by the condition(%#v), error info is %s", ctx.ReqID, common.BKTableNameObjClassifiction, cond.ToMapStr(), err.Error())
		return 0, err
	}
	return cnt, err
}

func (m *modelClassification) save(ctx core.ContextParams, classification metadata.Classification) (id uint64, err error) {

	id, err = m.dbProxy.NextSequence(ctx, common.BKTableNameObjClassifiction)
	if nil != err {
		blog.Errorf("request(%s): it is failed to create a new sequence id on the table(%s) of the database, error info is %s", ctx.ReqID, common.BKTableNameObjClassifiction, err.Error())
		return id, ctx.Error.New(common.CCErrObjectDBOpErrno, err.Error())
	}

	classification.ID = int64(id)
	classification.OwnerID = ctx.SupplierAccount

	err = m.dbProxy.Table(common.BKTableNameObjClassifiction).Insert(ctx, classification)
	return id, err
}

func (m *modelClassification) update(ctx core.ContextParams, data mapstr.MapStr, cond universalsql.Condition) (cnt uint64, err error) {

	cnt, err = m.count(ctx, cond)
	if nil != err {
		return cnt, err
	}

	if 0 == cnt {
		return cnt, nil
	}

	data.Remove(metadata.ClassFieldClassificationID)
	err = m.dbProxy.Table(common.BKTableNameObjClassifiction).Update(ctx, cond.ToMapStr(), data)
	if nil != err {
		blog.Errorf("request(%s): it is failed to execute a database update operation on the table(%s) by the condition(%#v) , error info is %s", ctx.ReqID, common.BKTableNameObjClassifiction, cond.ToMapStr(), err.Error())
		return 0, err
	}
	return cnt, err
}

func (m *modelClassification) delete(ctx core.ContextParams, cond universalsql.Condition) (cnt uint64, err error) {

	cnt, err = m.count(ctx, cond)
	if nil != err {
		return cnt, err
	}

	if 0 == cnt {
		return 0, err
	}

	err = m.dbProxy.Table(common.BKTableNameObjClassifiction).Delete(ctx, cond.ToMapStr())
	if nil != err {
		blog.Errorf("request(%s): it is failed to execute a database deletion operation on the table(%s) by the condition(%#v), error info is %s", ctx.ReqID, common.BKTableNameObjClassifiction, cond.ToMapStr(), err.Error())
		return 0, err
	}

	return cnt, err
}

func (m *modelClassification) search(ctx core.ContextParams, cond universalsql.Condition) ([]metadata.Classification, error) {

	results := make([]metadata.Classification, 0)
	err := m.dbProxy.Table(common.BKTableNameObjClassifiction).Find(cond.ToMapStr()).All(ctx, &results)
	return results, err
}

func (m *modelClassification) searchReturnMapStr(ctx core.ContextParams, cond universalsql.Condition) ([]mapstr.MapStr, error) {

	results := make([]mapstr.MapStr, 0)
	err := m.dbProxy.Table(common.BKTableNameObjClassifiction).Find(cond.ToMapStr()).All(ctx, &results)
	return results, err
}
