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
	"configcenter/src/common/metadata"
	"configcenter/src/common/universalsql/mongo"
	"configcenter/src/source_controller/coreservice/core"
)

func (m *associationModel) isValid(ctx core.ContextParams, inputParam metadata.CreateModelAssociation) error {

	if 0 == len(inputParam.Spec.AssociationName) {
		blog.Errorf("request(%s): it is failed to create a new model association, because of the associationID (%s) is not set", ctx.ReqID, metadata.AssociationFieldAsstID)
		return ctx.Error.Errorf(common.CCErrCommParamsNeedSet, metadata.AssociationFieldAsstID)
	}

	if 0 == len(inputParam.Spec.ObjectID) {
		blog.Errorf("request(%s): it is failed to create a new model association, because of the objectID (%s) is not set", ctx.ReqID, metadata.AssociationFieldObjectID)
		return ctx.Error.Errorf(common.CCErrCommParamsNeedSet, metadata.AssociationFieldObjectID)
	}

	if 0 == len(inputParam.Spec.AsstObjID) {
		blog.Errorf("request(%s): it is failed to create a new model association, because of the AssoObjectID (%s) is not set", ctx.ReqID, metadata.AssociationFieldAssociationObjectID)
		return ctx.Error.Errorf(common.CCErrCommParamsNeedSet, metadata.AssociationFieldAssociationObjectID)
	}

	if 0 == len(inputParam.Spec.AsstKindID) {
		blog.Errorf("request(%s): it is failed to create a new model association, because of the AssoObjectID (%s) is not set", ctx.ReqID, metadata.AssociationFieldAssociationKind)
		return ctx.Error.Errorf(common.CCErrCommParamsNeedSet, metadata.AssociationFieldAssociationObjectID)
	}

	return nil
}

func (m *associationModel) isExistsAssociationID(ctx core.ContextParams, associationID string) (bool, error) {

	existsCheckCond := mongo.NewCondition()
	existsCheckCond.Element(&mongo.Eq{Key: metadata.AssociationFieldAsstID, Val: associationID})
	existsCheckCond.Element(&mongo.Eq{Key: metadata.AssociationFieldSupplierAccount, Val: ctx.SupplierAccount})

	cnt, err := m.count(ctx, existsCheckCond)
	if nil != err {
		blog.Errorf("request(%s): it is to failed to check whether the associationID (%s) is exists, error info is %s", ctx.ReqID, associationID, err.Error())
		return false, err
	}
	return 0 != cnt, err
}

func (m *associationModel) isExistsAssociationObjectWithAnotherObject(ctx core.ContextParams, targetObjectID, anotherObjectID string, AssociationKind string) (bool, error) {

	existsCheckCond := mongo.NewCondition()
	existsCheckCond.Element(&mongo.Eq{Key: metadata.AssociationFieldSupplierAccount, Val: ctx.SupplierAccount})
	existsCheckCond.Element(&mongo.Eq{Key: metadata.AssociationFieldObjectID, Val: targetObjectID})
	existsCheckCond.Element(&mongo.Eq{Key: metadata.AssociationFieldAssociationObjectID, Val: anotherObjectID})
	existsCheckCond.Element(&mongo.Eq{Key: metadata.AssociationFieldAssociationKind, Val: AssociationKind})

	cnt, err := m.count(ctx, existsCheckCond)
	if nil != err {
		blog.Errorf("request(%s): it is to failed to check whether the association (%s=>%s) is exists by the condition (%#v), error info is %s", ctx.ReqID, targetObjectID, anotherObjectID, existsCheckCond.ToMapStr(), err.Error())
		return false, err
	}
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
