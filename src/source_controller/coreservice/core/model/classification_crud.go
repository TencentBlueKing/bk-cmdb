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
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/universalsql"
	"configcenter/src/common/util"
	"configcenter/src/storage/driver/mongodb"
)

func (m *modelClassification) count(kit *rest.Kit, cond universalsql.Condition) (cnt uint64, err error) {
	filter := util.SetModOwner(cond.ToMapStr(), kit.SupplierAccount)
	cnt, err = mongodb.Client().Table(common.BKTableNameObjClassification).Find(filter).Count(kit.Ctx)
	if nil != err {
		blog.Errorf("request(%s): it is failed to execute a database count operation on the table(%s) by the condition(%#v), error info is %s", kit.Rid, common.BKTableNameObjClassification, cond.ToMapStr(), err.Error())
		return 0, err
	}
	return cnt, err
}

func (m *modelClassification) save(kit *rest.Kit, classification metadata.Classification) (id uint64, err error) {

	id, err = mongodb.Client().NextSequence(kit.Ctx, common.BKTableNameObjClassification)
	if nil != err {
		blog.Errorf("request(%s): it is failed to create a new sequence id on the table(%s) of the database, error info is %s", kit.Rid, common.BKTableNameObjClassification, err.Error())
		return id, kit.CCError.New(common.CCErrObjectDBOpErrno, err.Error())
	}

	classification.ID = int64(id)
	classification.OwnerID = kit.SupplierAccount

	err = mongodb.Client().Table(common.BKTableNameObjClassification).Insert(kit.Ctx, classification)
	return id, err
}

func (m *modelClassification) update(kit *rest.Kit, data mapstr.MapStr, cond universalsql.Condition) (cnt uint64, err error) {

	cnt, err = m.count(kit, cond)
	if nil != err {
		return cnt, err
	}

	if 0 == cnt {
		return cnt, nil
	}

	data.Remove(metadata.ClassFieldClassificationID)
	err = mongodb.Client().Table(common.BKTableNameObjClassification).Update(kit.Ctx, cond.ToMapStr(), data)
	if nil != err {
		blog.Errorf("request(%s): it is failed to execute a database update operation on the table(%s) by the condition(%#v) , error info is %s", kit.Rid, common.BKTableNameObjClassification, cond.ToMapStr(), err.Error())
		return 0, err
	}
	return cnt, err
}

func (m *modelClassification) delete(kit *rest.Kit, cond universalsql.Condition) (cnt uint64, err error) {

	cnt, err = m.count(kit, cond)
	if nil != err {
		return cnt, err
	}

	if 0 == cnt {
		return 0, err
	}

	err = mongodb.Client().Table(common.BKTableNameObjClassification).Delete(kit.Ctx, cond.ToMapStr())
	if nil != err {
		blog.Errorf("request(%s): it is failed to execute a database deletion operation on the table(%s) by the condition(%#v), error info is %s", kit.Rid, common.BKTableNameObjClassification, cond.ToMapStr(), err.Error())
		return 0, err
	}

	return cnt, err
}

func (m *modelClassification) search(kit *rest.Kit, cond universalsql.Condition) ([]metadata.Classification, error) {

	results := make([]metadata.Classification, 0)
	err := mongodb.Client().Table(common.BKTableNameObjClassification).Find(cond.ToMapStr()).All(kit.Ctx, &results)
	return results, err
}

func (m *modelClassification) searchReturnMapStr(kit *rest.Kit, cond universalsql.Condition) ([]mapstr.MapStr, error) {

	results := make([]mapstr.MapStr, 0)
	err := mongodb.Client().Table(common.BKTableNameObjClassification).Find(cond.ToMapStr()).All(kit.Ctx, &results)
	return results, err
}
