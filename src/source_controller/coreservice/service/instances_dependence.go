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

package service

import (
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"

	//	"configcenter/src/common/universalsql/mongo"
	"configcenter/src/source_controller/coreservice/core"
	"configcenter/src/source_controller/coreservice/core/association"
	"configcenter/src/source_controller/coreservice/core/instances"
	"configcenter/src/source_controller/coreservice/core/model"
	"configcenter/src/storage/dal"
)

type instancesDepend struct {
	modelOperation core.ModelOperation
	asstOperation  core.AssociationOperation
}

func NewInstancesDepend(dbProxy dal.RDB) instances.OperationDependences {
	instDepend := &instancesDepend{}
	instDepend.modelOperation = model.New(dbProxy, nil)
	instDepend.asstOperation = association.New(dbProxy, nil)
	return instDepend

}

// IsInstanceExist used to check if the  instances  asst exist
func (m *instancesDepend) IsInstAsstExist(ctx core.ContextParams, objID string, instID uint64) (exists bool, err error) {
	cond := mapstr.MapStr{common.BKObjIDField: objID,
		common.BKInstIDField: instID}
	queryCond := metadata.QueryCondition{Condition: cond}
	objInsts, err := m.asstOperation.SearchInstanceAssociation(ctx, queryCond)
	if nil != err {
		blog.Errorf("search instance association error %v", err)
		return false, err
	}
	cond = mapstr.MapStr{
		common.BKAsstObjIDField:  objID,
		common.BKAsstInstIDField: instID,
	}
	queryCond = metadata.QueryCondition{Condition: cond}
	objAsstInsts, err := m.asstOperation.SearchInstanceAssociation(ctx, queryCond)
	if nil != err {
		blog.Errorf("search instance to association error %v", err)
		return false, err
	}
	if 0 == objInsts.Count && 0 == objAsstInsts.Count {
		return true, nil
	}
	return false, nil

}

// DeleteInstAsst used to delete inst asst
func (m *instancesDepend) DeleteInstAsst(ctx core.ContextParams, objID string, instID uint64) error {
	cond := mapstr.MapStr{common.BKObjIDField: objID,
		common.BKInstIDField: instID}
	deleteCond := metadata.DeleteOption{Condition: cond}
	_, err := m.asstOperation.DeleteInstanceAssociation(ctx, deleteCond)
	if nil != err {
		blog.Errorf("delete instance association error %v", err)
		return err
	}
	cond = mapstr.MapStr{
		common.BKAsstObjIDField:  objID,
		common.BKAsstInstIDField: instID,
	}
	deleteCond = metadata.DeleteOption{Condition: cond}
	_, err = m.asstOperation.DeleteInstanceAssociation(ctx, deleteCond)
	if nil != err {
		blog.Errorf("delete instance to association error %v", err)
		return err
	}
	return nil
}

// SelectObjectAttWithParams select object att with params
func (m *instancesDepend) SelectObjectAttWithParams(ctx core.ContextParams, objID string) (attributeArr []metadata.Attribute, err error) {
	attributeArr = make([]metadata.Attribute, 0)
	condition := mapstr.MapStr{common.BKObjIDField: mapstr.MapStr{common.BKDBEQ: objID}}
	queryCond := metadata.QueryCondition{
		Condition: condition,
	}
	result, err := m.modelOperation.SearchModelAttributes(ctx, objID, queryCond)
	for _, info := range result.Info {
		attribute := metadata.Attribute{}
		err := info.ToStructByTag(&attribute, "field")
		if nil != err {
			continue
		}
		attributeArr = append(attributeArr, attribute)
	}
	return attributeArr, nil
}

// SearchUnique search unique attribute
func (m *instancesDepend) SearchUnique(ctx core.ContextParams, objID string) (uniqueAttr []metadata.ObjectUnique, err error) {
	return uniqueAttr, nil
}
