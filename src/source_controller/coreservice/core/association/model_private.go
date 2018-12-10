/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.,
 * Copyright (C) 2017,-2018 THL A29 Limited, a Tencent company. All rights reserved.
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
	"configcenter/src/common/metadata"
	"configcenter/src/common/universalsql/mongo"
	"configcenter/src/source_controller/coreservice/core"
)

func (m *associationModel) isExistsAssociationID(ctx core.ContextParams, associationID string) (bool, error) {

	existsCheckCond := mongo.NewCondition()
	existsCheckCond.Element(&mongo.Eq{Key: metadata.AssociationFieldAsstID, Val: associationID})
	existsCheckCond.Element(&mongo.Eq{Key: metadata.AssociationFieldSupplierAccount, Val: ctx.SupplierAccount})

	cnt, err := m.count(ctx, existsCheckCond)
	return 0 != cnt, err
}

func (m *associationModel) isExistsAssociationObjectWithAnotherObject(ctx core.ContextParams, targetObjectID, anotherObjectID string) (bool, error) {

	existsCheckCond := mongo.NewCondition()
	existsCheckCond.Element(&mongo.Eq{Key: metadata.AssociationFieldSupplierAccount, Val: ctx.SupplierAccount})
	existsCheckCond.Element(&mongo.Eq{Key: metadata.AssociationFieldObjectID, Val: targetObjectID})
	existsCheckCond.Element(&mongo.Eq{Key: metadata.AssociationFieldAssociationObjectID, Val: anotherObjectID})

	cnt, err := m.count(ctx, existsCheckCond)
	return 0 != cnt, err
}

func (m *associationModel) usedInSomeInstanceAssociation(ctx core.ContextParams, associationIDS []string) (bool, error) {
	// TODO: need to implement
	return false, nil
}

func (m *associationModel) cascadeInstanceAssociation(ctx core.ContextParams, associationIDS []string) error {
	// TODO: need to implement
	return nil
}
