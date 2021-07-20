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
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/universalsql"
	"configcenter/src/storage/driver/mongodb"
)

func (m *associationModel) count(kit *rest.Kit, cond universalsql.Condition) (cnt uint64, err error) {

	cnt, err = mongodb.Client().Table(common.BKTableNameObjAsst).Find(cond.ToMapStr()).Count(kit.Ctx)
	if nil != err {
		blog.Errorf("request(%s): it is failed to execute database count operation on the table (%s) by the condition (%#v), error info is %s", kit.Rid, common.BKTableNameObjAsst, cond.ToMapStr(), err.Error())
		return 0, err
	}
	return cnt, err
}

func (m *associationModel) isExists(kit *rest.Kit, cond universalsql.Condition) (oneResult *metadata.Association, exists bool, err error) {

	oneResult = &metadata.Association{}
	err = mongodb.Client().Table(common.BKTableNameObjAsst).Find(cond.ToMapStr()).One(kit.Ctx, oneResult)
	if nil != err && !mongodb.Client().IsNotFoundError(err) {
		blog.Errorf("request(%s): it is faield to execute database findone operation on the table (%s) by the condition (%#v), error info is %s", kit.Rid, common.BKTableNameObjAsst, cond.ToMapStr(), err.Error())
		return oneResult, false, kit.CCError.New(common.CCErrObjectDBOpErrno, err.Error())
	}

	return oneResult, !mongodb.Client().IsNotFoundError(err), nil
}

func (m *associationModel) save(kit *rest.Kit, assoParam *metadata.Association) (id uint64, err error) {

	id, err = mongodb.Client().NextSequence(kit.Ctx, common.BKTableNameObjAsst)
	if nil != err {
		blog.Errorf("request(%s): it is failed to make a sequence ID on the table (%s), error info is %s", kit.Rid, common.BKTableNameObjAsst, err.Error())
		return id, err
	}

	assoParam.ID = int64(id)
	assoParam.OwnerID = kit.SupplierAccount
	err = mongodb.Client().Table(common.BKTableNameObjAsst).Insert(kit.Ctx, assoParam)
	if nil != err {
		blog.Errorf("request(%s): it is failed to execute database insert operation on the table (%s), error info is %s", kit.Rid, common.BKTableNameObjAsst, err.Error())
		return 0, err
	}
	return id, err
}

func (m *associationModel) update(kit *rest.Kit, data mapstr.MapStr, cond universalsql.Condition) (cnt uint64, err error) {

	cnt, err = m.count(kit, cond)
	if nil != err {
		return 0, err
	}

	if 0 >= cnt {
		return 0, err
	}

	err = mongodb.Client().Table(common.BKTableNameObjAsst).Update(kit.Ctx, cond.ToMapStr(), data)
	if nil != err {
		blog.Errorf("request(%s): it is failed to execute database upate some data (%v) on the table (%s) by the condition (%#v)", kit.Rid, data, common.BKTableNameObjAsst, cond.ToMapStr(), err.Error())
		return 0, err
	}
	return cnt, err
}

func (m *associationModel) delete(kit *rest.Kit, cond universalsql.Condition) (cnt uint64, err error) {

	cnt, err = m.count(kit, cond)
	if nil != err {
		return 0, err
	}

	if 0 >= cnt {
		return 0, err
	}

	err = mongodb.Client().Table(common.BKTableNameObjAsst).Delete(kit.Ctx, cond.ToMapStr())
	if nil != err {
		blog.Errorf("request(%s): it is to delete some data on the table (%s) by the condition (%#v), error info is %s", kit.Rid, common.BKTableNameObjAsst, cond.ToMapStr(), err.Error())
		return 0, err
	}
	return cnt, err
}

func (m *associationModel) search(kit *rest.Kit, cond universalsql.Condition) ([]metadata.Association, error) {

	dataResult := []metadata.Association{}
	err := mongodb.Client().Table(common.BKTableNameObjAsst).Find(cond.ToMapStr()).All(kit.Ctx, &dataResult)
	if nil != err {
		blog.Errorf("request(%s): it is to search some data on the table (%s) by the condition (%v), error info is %s", kit.Rid, common.BKTableNameObjAsst, cond.ToMapStr(), err.Error())
		return dataResult, err
	}
	return dataResult, err
}

func (m *associationModel) searchReturnMapStr(kit *rest.Kit, cond universalsql.Condition) ([]mapstr.MapStr, error) {
	dataResult := []mapstr.MapStr{}
	err := mongodb.Client().Table(common.BKTableNameObjAsst).Find(cond.ToMapStr()).All(kit.Ctx, &dataResult)
	if nil != err {
		blog.Errorf("request(%s): it is to search data on the table (%s) by the condition (%#v), error info is %s", kit.Rid, common.BKTableNameObjAsst, cond.ToMapStr(), err.Error())
		return dataResult, err
	}
	return dataResult, err
}
