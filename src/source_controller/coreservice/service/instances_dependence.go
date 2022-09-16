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

package service

import (
	"fmt"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/universalsql/mongo"
	"configcenter/src/common/util"
)

// IsInstAsstExist used to check if the  instances  asst exist
func (s *coreService) IsInstAsstExist(kit *rest.Kit, objID string, instID uint64) (exists bool, err error) {
	// to many call. can use $or. but universalsql parse or condtion error.

	cond := mongo.NewCondition()
	cond.Element(&mongo.Eq{Key: common.BKObjIDField, Val: objID}, &mongo.Eq{Key: common.BKInstIDField, Val: instID})
	countCond := &metadata.Condition{Condition: cond.ToMapStr()}
	objInstsRst, err := s.core.AssociationOperation().CountInstanceAssociations(kit, objID, countCond)
	if err != nil {
		blog.Errorf("search instance association err: %v, rid: %s", err, kit.Rid)
		return false, err
	}
	if objInstsRst.Count > 0 {
		return true, nil
	}

	cond = mongo.NewCondition()
	cond.Element(&mongo.Eq{Key: common.BKAsstObjIDField, Val: objID}, &mongo.Eq{Key: common.BKAsstInstIDField, Val: instID})
	countCond = &metadata.Condition{Condition: cond.ToMapStr()}
	objAsstInstsRst, err := s.core.AssociationOperation().CountInstanceAssociations(kit, objID, countCond)
	if err != nil {
		blog.Errorf("search instance to association err: %v, rid: %s", err, kit.Rid)
		return false, err
	}
	if objAsstInstsRst.Count > 0 {
		return true, nil
	}

	return false, nil
}

// DeleteInstAsst used to delete inst asst
func (s *coreService) DeleteInstAsst(kit *rest.Kit, objID string, instID uint64) error {
	cond := mongo.NewCondition()
	cond.Element(&mongo.Eq{Key: common.BKObjIDField, Val: objID}, &mongo.Eq{Key: common.BKInstIDField, Val: instID})
	deleteCond := metadata.DeleteOption{Condition: cond.ToMapStr()}
	_, err := s.core.AssociationOperation().DeleteInstanceAssociation(kit, objID, deleteCond)
	if nil != err {
		blog.Errorf("delete instance association error %v, rid: %s", err, kit.Rid)
		return err
	}
	cond = mongo.NewCondition()
	cond.Element(&mongo.Eq{Key: common.BKAsstObjIDField, Val: objID}, &mongo.Eq{Key: common.BKAsstInstIDField, Val: instID})
	deleteCond = metadata.DeleteOption{Condition: cond.ToMapStr()}
	_, err = s.core.AssociationOperation().DeleteInstanceAssociation(kit, objID, deleteCond)
	if nil != err {
		blog.Errorf("delete instance to association error %v, rid: %s", err, kit.Rid)
		return err
	}
	return nil
}

// SelectObjectAttWithParams select object att with params
func (s *coreService) SelectObjectAttWithParams(kit *rest.Kit, objID string, bizIDs []int64) (
	attributeArr []metadata.Attribute, err error) {

	attributeArr = make([]metadata.Attribute, 0)
	cond := mongo.NewCondition()
	cond.Element(&mongo.Eq{Key: common.BKObjIDField, Val: objID})
	queryCond := metadata.QueryCondition{
		Condition: cond.ToMapStr(),
	}

	bizCond := make(mapstr.MapStr)
	if len(bizIDs) > 1 {
		if err := util.AddModelWithMultipleBizIDCondition(bizCond, bizIDs); err != nil {
			return nil, err
		}

	} else if len(bizIDs) == 1 {
		util.AddModelBizIDCondition(bizCond, bizIDs[0])

	} else {
		blog.Errorf("bizIDs params must be set, rid: %s", kit.Rid)
		return nil, fmt.Errorf("biz ids params must be set")

	}

	queryCond.Condition.Merge(bizCond)
	result, err := s.core.ModelOperation().SearchModelAttributes(kit, objID, queryCond)
	if err != nil {
		blog.Errorf("select object att with params error %v, rid: %s", err, kit.Rid)
		return nil, err
	}
	return result.Info, nil
}

// SelectObjectAttributes select object attributes
func (s *coreService) SelectObjectAttributes(kit *rest.Kit, objID string, bizIDs []int64) ([]metadata.Attribute,
	error) {

	// query global attributes in model, all instances has these attributes and needs to be validated
	orCond := []map[string]interface{}{
		{common.BKAppIDField: 0},
		{common.BKAppIDField: mapstr.MapStr{common.BKDBExists: false}},
	}

	// if the biz ids are defined, query the biz attributes together with global attributes for validation
	if len(bizIDs) > 0 {
		orCond = append(orCond, map[string]interface{}{
			common.BKAppIDField: map[string]interface{}{common.BKDBIN: bizIDs}},
		)
	}

	queryCond := metadata.QueryCondition{
		Condition: map[string]interface{}{
			common.BKObjIDField: objID,
			common.BKDBOR:       orCond,
		},
	}

	result, err := s.core.ModelOperation().SearchModelAttributes(kit, objID, queryCond)
	if err != nil {
		blog.Errorf("select object(%s) attributes failed, err: %v, rid: %s", objID, err, kit.Rid)
		return nil, err
	}
	return result.Info, nil
}

// SearchUnique search unique attribute
func (s *coreService) SearchUnique(kit *rest.Kit, objID string) (uniqueAttr []metadata.ObjectUnique, err error) {
	condMap := util.SetQueryOwner(make(map[string]interface{}), kit.SupplierAccount)
	cond, _ := mongo.NewConditionFromMapStr(condMap)
	cond.Element(&mongo.Eq{Key: common.BKObjIDField, Val: objID})
	queryCond := metadata.QueryCondition{
		Condition: cond.ToMapStr(),
	}
	result, err := s.core.ModelOperation().SearchModelAttrUnique(kit, queryCond)
	return result.Info, err
}

// UpdateModelInstance TODO
func (s *coreService) UpdateModelInstance(kit *rest.Kit, objID string, param metadata.UpdateOption) (*metadata.UpdatedCount, error) {
	return s.core.InstanceOperation().UpdateModelInstance(kit, objID, param)
}
