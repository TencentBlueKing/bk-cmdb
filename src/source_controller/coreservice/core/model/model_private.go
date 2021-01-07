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
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
	"configcenter/src/common/universalsql"
	"configcenter/src/common/universalsql/mongo"
	"configcenter/src/common/util"
	"configcenter/src/storage/driver/mongodb"
)

func (m *modelManager) isExists(kit *rest.Kit, cond universalsql.Condition) (oneModel *metadata.Object, exists bool, err error) {

	oneModel = &metadata.Object{}
	err = mongodb.Client().Table(common.BKTableNameObjDes).Find(cond.ToMapStr()).One(kit.Ctx, oneModel)
	if nil != err && !mongodb.Client().IsNotFoundError(err) {
		blog.Errorf("request(%s): it is failed to execute database findOne operation on the table (%#v) by the condition (%#v), error info is %s", kit.Rid, common.BKTableNameObjDes, cond.ToMapStr(), err.Error())
		return oneModel, exists, kit.CCError.New(common.CCErrObjectDBOpErrno, err.Error())
	}
	exists = !mongodb.Client().IsNotFoundError(err)
	return oneModel, exists, nil
}

func (m *modelManager) isValid(kit *rest.Kit, objID string) error {
	checkCondMap := util.SetQueryOwner(make(map[string]interface{}), kit.SupplierAccount)
	checkCond, _ := mongo.NewConditionFromMapStr(checkCondMap)
	checkCond.Element(&mongo.Eq{Key: metadata.ModelFieldObjectID, Val: objID})

	cnt, err := mongodb.Client().Table(common.BKTableNameObjDes).Find(checkCond.ToMapStr()).Count(kit.Ctx)
	if nil != err {
		blog.Errorf("count operation on the table (%s) by the condition (%#v) failed , err: %v", common.BKTableNameObjDes, checkCond.ToMapStr(), err, kit.Rid)
		return kit.CCError.Error(common.CCErrObjectDBOpErrno)
	}

	if cnt == 0 {
		return kit.CCError.Errorf(common.CCErrCommParamsIsInvalid, objID)
	}

	return err
}

func (m *modelManager) deleteModelAndAttributes(kit *rest.Kit, targetObjIDS []string) (uint64, error) {

	// delete the attributes of the model
	deleteAttributeCond := mongo.NewCondition()
	deleteAttributeCond.Element(&mongo.In{Key: metadata.AttributeFieldObjectID, Val: targetObjIDS})
	cnt, err := m.modelAttribute.delete(kit, deleteAttributeCond)
	if nil != err {
		blog.Errorf("request(%s): it is failed to delete the attribute by the condition (%#v), error info is %s", kit.Rid, deleteAttributeCond.ToMapStr(), err.Error())
		return cnt, err
	}

	// delete the model self
	deleteModelCondMap := util.SetModOwner(make(map[string]interface{}), kit.SupplierAccount)
	deleteModelCond, _ := mongo.NewConditionFromMapStr(deleteModelCondMap)
	deleteModelCond.Element(&mongo.In{Key: metadata.ModelFieldObjectID, Val: targetObjIDS})

	cnt, err = m.delete(kit, deleteModelCond)
	if nil != err {
		blog.Errorf("request(%s): it is failed to delete some models by the condition (%#v), error info is %s", kit.Rid, deleteModelCond.ToMapStr(), err.Error())
		return 0, kit.CCError.New(common.CCErrObjectDBOpErrno, err.Error())
	}

	return cnt, nil
}
