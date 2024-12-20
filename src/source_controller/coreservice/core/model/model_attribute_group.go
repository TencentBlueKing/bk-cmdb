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
	"configcenter/src/common/universalsql/mongo"
	"configcenter/src/common/util"
	"configcenter/src/source_controller/coreservice/core"
	"configcenter/src/storage/driver/mongodb"
)

var _ core.ModelAttributeGroup = nil

type modelAttributeGroup struct {
	model *modelManager
}

// CreateModelAttributeGroup TODO
func (g *modelAttributeGroup) CreateModelAttributeGroup(kit *rest.Kit, objID string,
	inputParam metadata.CreateModelAttributeGroup) (*metadata.CreateOneDataResult, error) {
	dataResult := &metadata.CreateOneDataResult{}

	// NOW: Temporarily blocked
	// if err := g.model.isValid(kit, objID); err != nil {
	// 	blog.Errorf("failed to create a model attribute group (%#v), error: %v, rid: %s", kit.Rid, inputParam.Data,
	// 	err.CCError())
	// 	return dataResult, err
	// }
	inputParam.Data.ObjectID = objID

	_, isExists, err := g.groupIDIsExists(kit, objID, inputParam.Data.GroupID, inputParam.Data.BizID)
	if err != nil {
		blog.Errorf("failed to check the group ID %s if it is exists, error: %v, rid: %s", inputParam.Data.GroupID,
			err, kit.Rid)
		return dataResult, err
	}
	if isExists {
		blog.Errorf("failed to create a new group: %#v, because of the groupID %s is exists, rid: %s",
			inputParam.Data, inputParam.Data.GroupID, kit.Rid)
		return dataResult, kit.CCError.Errorf(common.CCErrCommDuplicateItem, inputParam.Data.GroupID)
	}

	_, isExists, err = g.groupNameIsExists(kit, objID, inputParam.Data.GroupName, inputParam.Data.BizID)
	if err != nil {
		blog.Errorf("failed to check the group name (%s) if it is exists, error: %v, rid: %s",
			inputParam.Data.GroupName, err, kit.Rid)
		return dataResult, err
	}
	if isExists {
		blog.Errorf("failed to check the group name (%s) if it is exists, error: %v, rid: %s",
			inputParam.Data.GroupName, err, kit.Rid)
		return dataResult, kit.CCError.Errorf(common.CCErrCommDuplicateItem, inputParam.Data.GroupName)
	}
	id, err := g.save(kit, inputParam.Data)
	if err != nil {
		blog.Errorf("failed to create a new model attribute group, data: %v, err: %v, rid: %s", inputParam.Data, err,
			kit.Rid)
		if mongodb.IsDuplicatedError(err) {
			dupErr := kit.CCError.CCErrorf(common.CCErrCommDuplicateItem, mongodb.GetDuplicateKey(err))
			return nil, dupErr
		}
		return dataResult, err
	}
	dataResult.Created.ID = id
	return dataResult, err
}

// SetModelAttributeGroup set the model attribute group
func (g *modelAttributeGroup) SetModelAttributeGroup(kit *rest.Kit, objID string,
	inputParam metadata.SetModelAttributeGroup) (*metadata.SetDataResult, error) {

	dataResult := &metadata.SetDataResult{
		Created:    []metadata.CreatedDataResult{},
		Updated:    []metadata.UpdatedDataResult{},
		Exceptions: []metadata.ExceptionResult{},
	}

	if err := g.model.isValid(kit, objID); err != nil {
		blog.Errorf("failed to set a model attribute group (%#v), error: %v, rid: %s", inputParam.Data, err, kit.Rid)
		return dataResult, err
	}
	inputParam.Data.ObjectID = objID

	_, isExists, err := g.groupNameIsExists(kit, objID, inputParam.Data.GroupName, inputParam.Data.BizID)
	if err != nil {
		blog.Errorf("failed to check the group name (%s) if it is exists, error: %v, rid: %s",
			inputParam.Data.GroupName, err, kit.Rid)
		return dataResult, err
	}
	if isExists {
		return dataResult, kit.CCError.Errorf(common.CCErrCommDuplicateItem, inputParam.Data.GroupName)
	}

	existsGroup, isExists, err := g.groupIDIsExists(kit, objID, inputParam.Data.GroupID, inputParam.Data.BizID)
	if err != nil {
		blog.Errorf("failed to check the group ID (%s) if it is exists, error: %v, rid: %s",
			inputParam.Data.GroupID, err, kit.Rid)
		return dataResult, err
	}

	if !isExists {

		id, err := g.save(kit, inputParam.Data)
		if err != nil {
			blog.Errorf("failed to create a new model attribute group (%#v), error: %v, rid: %s", inputParam.Data, err,
				kit.Rid)
			return &metadata.SetDataResult{}, err
		}

		dataResult.CreatedCount.Count++
		dataResult.Created = []metadata.CreatedDataResult{
			{
				ID: id,
			}}

		return dataResult, nil
	}

	cond := mongo.NewCondition()
	cond.Element(&mongo.Eq{Key: metadata.GroupFieldGroupID, Val: inputParam.Data.GroupID})
	cond.Element(&mongo.Eq{Key: metadata.GroupFieldID, Val: existsGroup.ID})

	cnt, err := g.update(kit, mapstr.NewFromStruct(inputParam.Data, "field"), cond)
	if err != nil {
		blog.Errorf("failed to update the model attribute group (%#v) by the condition (%#v), err: %s", g, cond, err,
			kit.Rid)
		return dataResult, err
	}
	dataResult.UpdatedCount.Count = cnt
	dataResult.Updated = []metadata.UpdatedDataResult{
		{
			ID: uint64(existsGroup.ID),
		},
	}
	return dataResult, nil
}

// UpdateModelAttributeGroup update model attribute group
func (g *modelAttributeGroup) UpdateModelAttributeGroup(kit *rest.Kit, objID string,
	inputParam metadata.UpdateOption) (*metadata.UpdatedCount, error) {

	if err := g.model.isValid(kit, objID); err != nil {
		blog.Errorf("failed to create a model attribute group (%#v), error: %v, rid: %s", inputParam.Data, err, kit.Rid)
		return &metadata.UpdatedCount{}, err
	}

	cond, err := mongo.NewConditionFromMapStr(inputParam.Condition.ToMapInterface())
	if err != nil {
		blog.Errorf("failed to convert the condition (%#v) from mapstr to condition, error: %v, rid: %s",
			inputParam.Condition, err, kit.Rid)
		return &metadata.UpdatedCount{}, err
	}
	cond.Element(&mongo.Eq{Key: metadata.GroupFieldObjectID, Val: objID})

	inputParam.Data.Remove(metadata.GroupFieldGroupID)
	inputParam.Data.Remove(metadata.GroupFieldObjectID)
	inputParam.Data.Remove(metadata.GroupFieldIsPre)

	if name, exists := inputParam.Data.Get("bk_group_name"); exists {
		name := name.(string)
		queryCond := metadata.QueryCondition{
			Condition: cond.ToMapStr(),
		}
		resp, err := g.SearchModelAttributeGroupByCondition(kit, queryCond)
		if err != nil {
			blog.Errorf("failed to check the group name (%s) if it is exists, error: %v, rid: %s", name, err, kit.Rid)
			return &metadata.UpdatedCount{}, err
		}
		for _, item := range resp.Info {
			if item.GroupName == name {
				continue
			}
			_, exists, err := g.groupNameIsExists(kit, item.ObjectID, name, 0)
			if err != nil {
				blog.Errorf("failed to check the group name (%s) if it is exists, error: %v, rid: %s", name, err,
					kit.Rid)
				return &metadata.UpdatedCount{}, err
			}
			if exists {
				blog.Errorf("failed to update the group name, because (%s) exists, rid: %s", name, kit.Rid)
				return &metadata.UpdatedCount{}, kit.CCError.Errorf(common.CCErrCommDuplicateItem, name)
			}
		}
	}
	cnt, err := g.update(kit, inputParam.Data, cond)
	if err != nil {
		blog.Errorf("failed to update the data (%s) by the condition (%#v), error: %v, rid: %s",
			inputParam.Data, err, kit.Rid)
		return &metadata.UpdatedCount{}, err
	}

	return &metadata.UpdatedCount{Count: cnt}, nil
}

// UpdateModelAttributeGroupByCondition update model attribute group by condition, check if name repeat by bizID
func (g *modelAttributeGroup) UpdateModelAttributeGroupByCondition(kit *rest.Kit, inputParam metadata.UpdateOption) (
	*metadata.UpdatedCount, error) {

	cond, err := mongo.NewConditionFromMapStr(inputParam.Condition.ToMapInterface())
	if err != nil {
		blog.Errorf("failed to convert the condition (%#v) from mapstr to condition, err: %v, rid: %s",
			inputParam.Condition, err, kit.Rid)
		return &metadata.UpdatedCount{}, err
	}

	inputParam.Data.Remove(metadata.GroupFieldGroupID)
	inputParam.Data.Remove(metadata.GroupFieldObjectID)
	inputParam.Data.Remove(metadata.GroupFieldIsPre)

	bizID := int64(0)
	if paramBizID, exist := inputParam.Condition.Get(common.BKAppIDField); exist {
		if bizID, err = util.GetInt64ByInterface(paramBizID); err != nil {
			blog.Errorf("bk_biz_id type is not int64, err: %v, rid: %s", err, kit.Rid)
			return &metadata.UpdatedCount{}, err
		}
	}

	if name, exists := inputParam.Data.Get("bk_group_name"); exists {
		name := name.(string)
		queryCond := metadata.QueryCondition{
			Condition: cond.ToMapStr(),
		}
		resp, err := g.SearchModelAttributeGroupByCondition(kit, queryCond)
		if err != nil {
			blog.Errorf("failed to check the group name (%s) if it is exists, err: %v, rid: %s", name, err, kit.Rid)
			return &metadata.UpdatedCount{}, err
		}
		for _, item := range resp.Info {
			if item.GroupName == name {
				continue
			}
			_, exists, err := g.groupNameIsExists(kit, item.ObjectID, name, bizID)
			if err != nil {
				blog.Errorf("failed to check the group name (%s) if it is exists, err: %v, rid: %s", name, err,
					kit.Rid)
				return &metadata.UpdatedCount{}, err
			}
			if exists {
				blog.Errorf("failed to update the group name, because (%s) exists, rid: %s", name, kit.Rid)
				return &metadata.UpdatedCount{}, kit.CCError.Errorf(common.CCErrCommDuplicateItem, name)
			}
		}
	}

	cnt, err := g.update(kit, inputParam.Data, cond)
	if err != nil {
		blog.Errorf("failed to update the data (%s) by the condition (%#v), err: %v, rid: %s", inputParam.Data, cond,
			err, kit.Rid)
		return &metadata.UpdatedCount{}, err
	}

	return &metadata.UpdatedCount{Count: cnt}, nil
}

// SearchModelAttributeGroup searches the model attribute group
func (g *modelAttributeGroup) SearchModelAttributeGroup(kit *rest.Kit, objID string,
	inputParam metadata.QueryCondition) (*metadata.QueryModelAttributeGroupDataResult, error) {

	dataResult := &metadata.QueryModelAttributeGroupDataResult{
		Info: []metadata.Group{},
	}

	// TODO: Paging queries needs to be implemented
	// NOW: Temporarily blocked
	// if err := g.model.isValid(kit, objID); err != nil {
	// 	blog.Errorf("failed to query a model by the condition(%#v), error: %v, rid: %s", kit.Rid, inputParam.Condition,
	// 	err.CCError())
	// 	return dataResult, err
	// }

	cond, err := mongo.NewConditionFromMapStr(inputParam.Condition.ToMapInterface())
	if err != nil {
		blog.Errorf("failed to convert the condition (%#v) from mapstr to condition, error: %v, rid: %s",
			inputParam.Condition, err, kit.Rid)
		return dataResult, err
	}
	cond.Element(&mongo.Eq{Key: metadata.GroupFieldObjectID, Val: objID})

	totalCount, err := g.count(kit, cond)
	if err != nil {
		blog.Errorf("failed to count by the condition (%#v), error: %v, rid: %s ", cond.ToMapStr(), err, kit.Rid)
		return dataResult, err
	}

	grps, err := g.search(kit, cond)
	if err != nil {
		blog.Errorf("failed to query a model by the condition(%#v), error: %v, rid: %s", cond.ToMapStr(), err, kit.Rid)
		return dataResult, err
	}

	dataResult.Count = totalCount
	dataResult.Info = grps
	return dataResult, nil

}

// SearchModelAttributeGroupByCondition searches the model attribute group by condition
func (g *modelAttributeGroup) SearchModelAttributeGroupByCondition(kit *rest.Kit,
	inputParam metadata.QueryCondition) (*metadata.QueryModelAttributeGroupDataResult, error) {

	cond, err := mongo.NewConditionFromMapStr(inputParam.Condition.ToMapInterface())
	if err != nil {
		blog.Errorf("failed to convert the condition (%#v) from mapstr to condition, error: %v, rid: %s",
			inputParam.Condition, err, kit.Rid)
		return &metadata.QueryModelAttributeGroupDataResult{}, err
	}

	totalCount, err := g.count(kit, cond)
	if err != nil {
		blog.Errorf("failed to count by the condition (%#v), error: %v, rid: %s ", cond.ToMapStr(), err, kit.Rid)
		return &metadata.QueryModelAttributeGroupDataResult{}, err
	}

	grps, err := g.search(kit, cond)
	if err != nil {
		blog.Errorf("failed to query a model by the condition(%#v), error: %v, rid: %s", cond.ToMapStr(), err, kit.Rid)
		return &metadata.QueryModelAttributeGroupDataResult{}, err
	}

	return &metadata.QueryModelAttributeGroupDataResult{Count: totalCount, Info: grps}, nil

}

// DeleteModelAttributeGroupByCondition deletes the model attribute group by condition
// desperated only for old api
func (g *modelAttributeGroup) DeleteModelAttributeGroupByCondition(kit *rest.Kit,
	inputParam metadata.DeleteOption) (*metadata.DeletedCount, error) {

	cond, err := mongo.NewConditionFromMapStr(inputParam.Condition.ToMapInterface())
	if err != nil {
		blog.Errorf("failed to convert the condition (%#v) from mapstr to condition, error: %v, rid: %s",
			inputParam.Condition, err, kit.Rid)
		return &metadata.DeletedCount{}, err
	}

	grps, err := g.search(kit, cond)
	if err != nil {
		blog.Errorf("failed to query model attribute groups by the condition (%#v), error: %v, rid: %s",
			cond.ToMapStr(), err, kit.Rid)
		return &metadata.DeletedCount{}, err
	}

	for _, grp := range grps {
		hasAttrs, err := g.hasAttributes(kit, grp.ObjectID, []string{grp.GroupID})
		if err != nil {
			blog.Errorf("failed to check whether the group(%s) has some attributes for the model(%s), error: %v, rid: %s",
				grp.GroupID, grp.ObjectID, kit.Rid)
			return &metadata.DeletedCount{}, err
		}
		if hasAttrs {
			blog.Errorf("the group(%s) has some attributes, forbidden to delete, rid: %s", grp.GroupID, kit.Rid)
			return &metadata.DeletedCount{}, kit.CCError.Error(common.CCErrCoreServiceModelAttributeGroupHasSomeAttributes)
		}
	}

	cnt, err := g.delete(kit, cond)
	if err != nil {
		blog.Errorf("failed to delete the group by the condition(%#v), error: %v, rid: %s", cond.ToMapStr(), err,
			kit.Rid)
		return &metadata.DeletedCount{}, err
	}

	return &metadata.DeletedCount{Count: cnt}, nil
}

// DeleteModelAttributeGroup deletes the model attribute group
func (g *modelAttributeGroup) DeleteModelAttributeGroup(kit *rest.Kit, objID string,
	inputParam metadata.DeleteOption) (*metadata.DeletedCount, error) {

	if err := g.model.isValid(kit, objID); err != nil {
		blog.Errorf("failed to delete a model by the condition(%#v), error: %v, rid: %s", inputParam.Condition, err,
			kit.Rid)
		return &metadata.DeletedCount{}, err
	}

	cond, err := mongo.NewConditionFromMapStr(inputParam.Condition.ToMapInterface())
	if err != nil {
		blog.Errorf("failed to convert the condition (%#v) from mapstr to condition, error: %v, rid: %s",
			inputParam.Condition, err, kit.Rid)
		return &metadata.DeletedCount{}, err
	}
	cond.Element(&mongo.Eq{Key: metadata.GroupFieldObjectID, Val: objID})

	grps, err := g.search(kit, cond)
	if err != nil {
		blog.Errorf("failed to query model attribute groups by the condition (%#v), error: %v, rid: %s",
			cond.ToMapStr(), err, kit.Rid)
		return &metadata.DeletedCount{}, err
	}

	grpIDS := make([]string, 0)
	for _, grp := range grps {
		grpIDS = append(grpIDS, grp.GroupID)
	}

	isExists, err := g.hasAttributes(kit, objID, grpIDS)
	if err != nil {
		blog.Errorf("failed to check whether the group IDS (%#v) has some attributes for the model(%s), error: %v, rid: %s",
			grpIDS, objID, err, kit.Rid)
		return &metadata.DeletedCount{}, err
	}

	if isExists {
		return &metadata.DeletedCount{}, kit.CCError.Error(common.CCErrCoreServiceModelAttributeGroupHasSomeAttributes)
	}

	cnt, err := g.delete(kit, cond)
	if err != nil {
		blog.Errorf("failed to delete the group by the condition(%#v), error: %v, rid: %s", cond.ToMapStr(), err,
			kit.Rid)
		return &metadata.DeletedCount{}, err
	}

	return &metadata.DeletedCount{Count: cnt}, nil
}
