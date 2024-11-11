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
	"configcenter/src/storage/driver/mongodb"
)

func (m *modelManager) isExists(kit *rest.Kit, cond universalsql.Condition) (oneModel *metadata.Object, exists bool, err error) {

	oneModel = &metadata.Object{}
	err = mongodb.Client().Table(common.BKTableNameObjDes).Find(cond.ToMapStr()).One(kit.Ctx, oneModel)
	if err != nil && !mongodb.Client().IsNotFoundError(err) {
		blog.Errorf("execute database findOne operation failed, err: %v, cond: %v, rid: %s", err, cond,
			kit.Rid)
		return oneModel, exists, kit.CCError.New(common.CCErrObjectDBOpErrno, err.Error())
	}
	exists = !mongodb.Client().IsNotFoundError(err)
	return oneModel, exists, nil
}

func (m *modelManager) isValid(kit *rest.Kit, objID string) error {
	checkCond, _ := mongo.NewConditionFromMapStr(make(map[string]interface{}))
	checkCond.Element(&mongo.Eq{Key: metadata.ModelFieldObjectID, Val: objID})

	cnt, err := mongodb.Client().Table(common.BKTableNameObjDes).Find(checkCond.ToMapStr()).Count(kit.Ctx)
	if err != nil {
		blog.Errorf("count operation on the table (%s) by the condition (%#v) failed, err: %v, rid: %s",
			common.BKTableNameObjDes, checkCond.ToMapStr(), err, kit.Rid)
		return err
	}

	if cnt == 0 {
		blog.Errorf("object [%s] has not been created, rid: %s", objID, kit.Rid)
		return kit.CCError.Errorf(common.CCErrCommParamsIsInvalid, objID)
	}

	return nil
}

func (m *modelManager) deleteModelAndAttributes(kit *rest.Kit, targetObjIDS []string) (uint64, error) {

	// delete the attributes of the model
	deleteAttributeCond := mongo.NewCondition()
	deleteAttributeCond.Element(&mongo.In{Key: metadata.AttributeFieldObjectID, Val: targetObjIDS})
	cnt, err := m.modelAttribute.delete(kit, deleteAttributeCond, true)
	if nil != err {
		blog.Errorf("delete the attribute failed, err: %v, cond: %v, rid: %s", err, deleteAttributeCond,
			kit.Rid)
		return cnt, err
	}

	// delete the model self
	deleteModelCond, _ := mongo.NewConditionFromMapStr(make(map[string]interface{}))
	deleteModelCond.Element(&mongo.In{Key: metadata.ModelFieldObjectID, Val: targetObjIDS})

	cnt, err = m.delete(kit, deleteModelCond)
	if nil != err {
		blog.Errorf("delete models failed, err: %v, cond: %v, rid: %s", err, deleteModelCond, kit.Rid)
		return 0, kit.CCError.New(common.CCErrObjectDBOpErrno, err.Error())
	}

	return cnt, nil
}
