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
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/universalsql/mongo"
	"configcenter/src/common/util"
)

// IsInstanceExist used to check if the  instances  asst exist
func (s *coreService) IsInstAsstExist(kit *rest.Kit, objID string, instID uint64) (exists bool, err error) {
	// to many call. can use $or. but universalsql parse or condtion error.

	cond := mongo.NewCondition()
	cond.Element(&mongo.Eq{Key: common.BKObjIDField, Val: objID}, &mongo.Eq{Key: common.BKInstIDField, Val: instID})
	queryCond := metadata.QueryCondition{Condition: cond.ToMapStr()}
	objInsts, err := s.core.AssociationOperation().SearchInstanceAssociation(kit, queryCond)
	if nil != err {
		blog.Errorf("search instance association error %v, rid: %s", err, kit.Rid)
		return false, err
	}
	cond = mongo.NewCondition()
	cond.Element(&mongo.Eq{Key: common.BKAsstObjIDField, Val: objID}, &mongo.Eq{Key: common.BKAsstInstIDField, Val: instID})
	queryCond = metadata.QueryCondition{Condition: cond.ToMapStr()}
	objAsstInsts, err := s.core.AssociationOperation().SearchInstanceAssociation(kit, queryCond)
	if nil != err {
		blog.Errorf("search instance to association error %v, rid: %s", err, kit.Rid)
		return false, err
	}
	if 0 < objInsts.Count || 0 < objAsstInsts.Count {
		return true, nil
	}
	return false, nil

}

// DeleteInstAsst used to delete inst asst
func (s *coreService) DeleteInstAsst(kit *rest.Kit, objID string, instID uint64) error {
	cond := mongo.NewCondition()
	cond.Element(&mongo.Eq{Key: common.BKObjIDField, Val: objID}, &mongo.Eq{Key: common.BKInstIDField, Val: instID})
	deleteCond := metadata.DeleteOption{Condition: cond.ToMapStr()}
	_, err := s.core.AssociationOperation().DeleteInstanceAssociation(kit, deleteCond)
	if nil != err {
		blog.Errorf("delete instance association error %v, rid: %s", err, kit.Rid)
		return err
	}
	cond = mongo.NewCondition()
	cond.Element(&mongo.Eq{Key: common.BKAsstObjIDField, Val: objID}, &mongo.Eq{Key: common.BKAsstInstIDField, Val: instID})
	deleteCond = metadata.DeleteOption{Condition: cond.ToMapStr()}
	_, err = s.core.AssociationOperation().DeleteInstanceAssociation(kit, deleteCond)
	if nil != err {
		blog.Errorf("delete instance to association error %v, rid: %s", err, kit.Rid)
		return err
	}
	return nil
}

// SelectObjectAttWithParams select object att with params
func (s *coreService) SelectObjectAttWithParams(kit *rest.Kit, objID string, bizID int64) (attributeArr []metadata.Attribute, err error) {
	attributeArr = make([]metadata.Attribute, 0)
	cond := mongo.NewCondition()
	cond.Element(&mongo.Eq{Key: common.BKObjIDField, Val: objID})
	queryCond := metadata.QueryCondition{
		Condition: cond.ToMapStr(),
	}

	bizCond := make(mapstr.MapStr)
	util.AddModelBizIDConditon(bizCond, bizID)

	queryCond.Condition.Merge(bizCond)
	result, err := s.core.ModelOperation().SearchModelAttributes(kit, objID, queryCond)
	if err != nil {
		blog.Errorf("select object att with params error %v, rid: %s", err, kit.Rid)
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

func (s *coreService) UpdateModelInstance(kit *rest.Kit, objID string, param metadata.UpdateOption) (*metadata.UpdatedCount, error) {
	return s.core.InstanceOperation().UpdateModelInstance(kit, objID, param)
}
