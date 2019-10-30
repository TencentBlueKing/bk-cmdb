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
	"configcenter/src/common/metadata"
	"configcenter/src/common/universalsql"
	"configcenter/src/common/universalsql/mongo"
	"configcenter/src/source_controller/coreservice/core"
)

func (m *modelManager) isExists(ctx core.ContextParams, cond universalsql.Condition) (oneModel *metadata.Object, exists bool, err error) {

	oneModel = &metadata.Object{}
	err = m.dbProxy.Table(common.BKTableNameObjDes).Find(cond.ToMapStr()).One(ctx, oneModel)
	if nil != err && !m.dbProxy.IsNotFoundError(err) {
		blog.Errorf("request(%s): it is failed to execute database findOne operation on the table (%#v) by the condition (%#v), error info is %s", ctx.ReqID, common.BKTableNameObjDes, cond.ToMapStr(), err.Error())
		return oneModel, exists, ctx.Error.New(common.CCErrObjectDBOpErrno, err.Error())
	}
	exists = !m.dbProxy.IsNotFoundError(err)
	return oneModel, exists, nil
}

func (m *modelManager) isValid(ctx core.ContextParams, objID string) error {

	checkCond := mongo.NewCondition()
	//	checkCond.Element(&mongo.Eq{Key: metadata.ModelFieldOwnerID, Val: ctx.SupplierAccount})
	checkCond.Element(&mongo.Eq{Key: metadata.ModelFieldObjectID, Val: objID})

	cnt, err := m.dbProxy.Table(common.BKTableNameObjDes).Find(checkCond.ToMapStr()).Count(ctx)
	if nil != err {
		blog.Errorf("request(%s): it is failed to execute database count operation on the table (%s) by the condition (%#v), error info is %s", ctx.ReqID, common.BKTableNameObjDes, checkCond.ToMapStr(), err.Error())
		return ctx.Error.Error(common.CCErrObjectDBOpErrno)
	}

	if cnt == 0 {
		return ctx.Error.Errorf(common.CCErrCommParamsIsInvalid, objID)
	}

	return err
}

func (m *modelManager) deleteModelAndAttributes(ctx core.ContextParams, targetObjIDS []string) (uint64, error) {

	// delete the attributes of the model
	deleteAttributeCond := mongo.NewCondition()
	deleteAttributeCond.Element(&mongo.In{Key: metadata.AttributeFieldObjectID, Val: targetObjIDS})
	cnt, err := m.modelAttribute.delete(ctx, deleteAttributeCond)
	if nil != err {
		blog.Errorf("request(%s): it is failed to delete the attribute by the condition (%#v), error info is %s", ctx.ReqID, deleteAttributeCond.ToMapStr(), err.Error())
		return cnt, err
	}

	// delete the model self
	deleteModelCond := mongo.NewCondition()
	deleteModelCond.Element(&mongo.Eq{Key: metadata.ModelFieldOwnerID, Val: ctx.SupplierAccount})
	deleteModelCond.Element(&mongo.In{Key: metadata.ModelFieldObjectID, Val: targetObjIDS})

	cnt, err = m.delete(ctx, deleteModelCond)
	if nil != err {
		blog.Errorf("request(%s): it is failed to delete some models by the condition (%#v), error info is %s", ctx.ReqID, deleteModelCond.ToMapStr(), err.Error())
		return 0, ctx.Error.New(common.CCErrObjectDBOpErrno, err.Error())
	}

	return cnt, nil
}
