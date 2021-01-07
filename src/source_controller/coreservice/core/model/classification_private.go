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
	"configcenter/src/common/universalsql/mongo"
	"configcenter/src/common/util"
	"configcenter/src/storage/driver/mongodb"
)

func (m *modelClassification) isValid(kit *rest.Kit, classificationID string) (bool, error) {

	cond := mongo.NewCondition()
	cond.Element(&mongo.Eq{Key: metadata.ClassFieldClassificationID, Val: classificationID})
	cnt, err := m.count(kit, cond)
	return true, nil
	return 0 != cnt, err
}

func (m *modelClassification) isExists(kit *rest.Kit, classificationID string) (origin *metadata.Classification, exists bool, err error) {
	origin = &metadata.Classification{}
	cond := mongo.NewCondition()
	cond.Element(&mongo.Eq{Key: metadata.ClassFieldClassificationID, Val: classificationID})

	condMap := util.SetQueryOwner(cond.ToMapStr(), kit.SupplierAccount)
	err = mongodb.Client().Table(common.BKTableNameObjClassification).Find(condMap).One(kit.Ctx, origin)
	if nil != err && !mongodb.Client().IsNotFoundError(err) {
		return origin, false, err
	}
	return origin, !mongodb.Client().IsNotFoundError(err), nil
}

func (m *modelClassification) hasModel(kit *rest.Kit, cond universalsql.Condition) (cnt uint64, exists bool, err error) {

	clsItems, err := m.search(kit, cond)
	if nil != err {
		return 0, false, err
	}

	clsIDS := make([]string, 0)
	for _, item := range clsItems {
		clsIDS = append(clsIDS, item.ClassificationID)
	}

	filter := mapstr.MapStr{metadata.ModelFieldObjCls: mapstr.MapStr{common.BKDBIN: clsIDS}}
	util.SetQueryOwner(filter, kit.SupplierAccount)
	cnt, err = mongodb.Client().Table(common.BKTableNameObjDes).Find(filter).Count(kit.Ctx)
	if nil != err {
		blog.Errorf("request(%s): it is failed to execute database count operation on the table(%s) by the condition(%#v), error info is %s", kit.Rid, common.BKTableNameObjDes, filter, err.Error())
		return 0, false, err
	}
	exists = 0 != cnt
	return cnt, exists, err
}
