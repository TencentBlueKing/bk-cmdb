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
)

func (m *modelClassification) isValid(kit *rest.Kit, classificationID string) (bool, error) {

	cond := mongo.NewCondition()
	cond.Element(&mongo.Eq{Key: metadata.ClassFieldClassificationID, Val: classificationID})
	cond.Element(&mongo.Eq{Key: metadata.ClassFieldClassificationSupplierAccount, Val: kit.SupplierAccount})

	cnt, err := m.count(kit, cond)
	return 0 != cnt, err
}

func (m *modelClassification) isExists(kit *rest.Kit, classificationID string, meta metadata.Metadata) (origin *metadata.Classification, exists bool, err error) {

	origin = &metadata.Classification{}
	cond := mongo.NewCondition()
	cond.Element(&mongo.Eq{Key: metadata.ClassFieldClassificationSupplierAccount, Val: kit.SupplierAccount})
	cond.Element(&mongo.Eq{Key: metadata.ClassFieldClassificationID, Val: classificationID})

	// ATTENTION: Currently only business dimension isolation is done,
	//           and there may be isolation requirements for other dimensions in the future.
	isExist, bizID := meta.Label.Get(common.BKAppIDField)
	if isExist {
		_, metaCond := cond.Embed(metadata.BKMetadata)
		_, labelCond := metaCond.Embed(metadata.BKLabel)
		labelCond.Element(&mongo.Eq{Key: common.BKAppIDField, Val: bizID})
	}
	condMap := util.SetQueryOwner(cond.ToMapStr(), kit.SupplierAccount)
	err = m.dbProxy.Table(common.BKTableNameObjClassification).Find(condMap).One(kit.Ctx, origin)
	if nil != err && !m.dbProxy.IsNotFoundError(err) {
		return origin, false, err
	}
	return origin, !m.dbProxy.IsNotFoundError(err), nil
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

	checkModelCond := mongo.NewCondition()
	checkModelCond.Element(mongo.Field(metadata.ModelFieldObjCls).In(clsIDS))
	checkModelCond.Element(mongo.Field(metadata.ModelFieldOwnerID).Eq(kit.SupplierAccount))

	cnt, err = m.dbProxy.Table(common.BKTableNameObjDes).Find(checkModelCond.ToMapStr()).Count(kit.Ctx)
	if nil != err {
		blog.Errorf("request(%s): it is failed to execute database count operation on the table(%s) by the condition(%#v), error info is %s", kit.Rid, common.BKTableNameObjDes, cond.ToMapStr(), err.Error())
		return 0, false, err
	}
	exists = 0 != cnt
	return cnt, exists, err
}
