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
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/universalsql/mongo"
	"configcenter/src/common/util"
	"configcenter/src/source_controller/coreservice/core"
	"configcenter/src/storage/dal"
)

var _ core.ModelAttributeGroup = nil

type modelAttributeGroup struct {
	model   *modelManager
	dbProxy dal.RDB
}

func (g *modelAttributeGroup) CreateModelAttributeGroup(ctx core.ContextParams, objID string, inputParam metadata.CreateModelAttributeGroup) (*metadata.CreateOneDataResult, error) {
	dataResult := &metadata.CreateOneDataResult{}

	// NOW: Temporarily blocked
	// if err := g.model.isValid(ctx, objID); nil != err {
	// 	blog.Errorf("request(%s): it is failed to create a model attribute group (%#v), error info is %s", ctx.ReqID, inputParam.Data, err.Error())
	// 	return dataResult, err
	// }
	inputParam.Data.ObjectID = objID
	inputParam.Data.OwnerID = ctx.SupplierAccount

	_, isExists, err := g.groupIDIsExists(ctx, objID, inputParam.Data.GroupID, inputParam.Data.Metadata)
	if nil != err {
		blog.Errorf("request(%s): it is to failed to check the group ID (%s) if it is exists, error info is %s", ctx.ReqID, inputParam.Data.GroupID, err.Error())
		return dataResult, err
	}
	if isExists {
		blog.Errorf("request(%s): it is to failed to create a new group (%#v), because of the groupID (%s) is exists", ctx.ReqID, inputParam.Data, inputParam.Data.GroupID)
		return dataResult, ctx.Error.Errorf(common.CCErrCommDuplicateItem, inputParam.Data.GroupID)
	}

	_, isExists, err = g.groupNameIsExists(ctx, objID, inputParam.Data.GroupName, inputParam.Data.Metadata)
	if nil != err {
		blog.Errorf("request(%s): it is to failed to check the group name (%s) if it is exists, error info is %s", ctx.ReqID, inputParam.Data.GroupName, err.Error())
		return dataResult, err
	}
	if isExists {
		blog.Errorf("request(%s): it is to failed to check the group name (%s) if it is exists, error info is %#v", ctx.ReqID, inputParam.Data.GroupName, err)
		return dataResult, ctx.Error.Errorf(common.CCErrCommDuplicateItem, inputParam.Data.GroupName)
	}
	id, err := g.save(ctx, inputParam.Data)
	if nil != err {
		blog.Errorf("request(%s): it is to failed to create a new model attribute group (%#v), error info is %s", ctx.ReqID, inputParam.Data, err.Error())
		return dataResult, err
	}
	dataResult.Created.ID = id
	return dataResult, err
}

func (g *modelAttributeGroup) SetModelAttributeGroup(ctx core.ContextParams, objID string, inputParam metadata.SetModelAttributeGroup) (*metadata.SetDataResult, error) {

	dataResult := &metadata.SetDataResult{
		Created:    []metadata.CreatedDataResult{},
		Updated:    []metadata.UpdatedDataResult{},
		Exceptions: []metadata.ExceptionResult{},
	}

	if err := g.model.isValid(ctx, objID); nil != err {
		blog.Errorf("request(%s): it is failed to set a model attribute group (%#v), error info is %s", ctx.ReqID, inputParam.Data, err.Error())
		return dataResult, err
	}
	inputParam.Data.ObjectID = objID
	inputParam.Data.OwnerID = ctx.SupplierAccount

	_, isExists, err := g.groupNameIsExists(ctx, objID, inputParam.Data.GroupName, inputParam.Data.Metadata)
	if nil != err {
		blog.Errorf("request(%s): it is to failed to check the group name (%s) if it is exists, error info is %s", ctx.ReqID, inputParam.Data.GroupName, err.Error())
		return dataResult, err
	}
	if isExists {
		return dataResult, ctx.Error.Errorf(common.CCErrCommDuplicateItem, inputParam.Data.GroupName)
	}

	existsGroup, isExists, err := g.groupIDIsExists(ctx, objID, inputParam.Data.GroupID, inputParam.Data.Metadata)
	if nil != err {
		blog.Errorf("request(%s): it is to failed to check the group ID (%s) if it is exists, error info is %s", ctx.ReqID, inputParam.Data.GroupID, err.Error())
		return dataResult, err
	}

	if !isExists {

		id, err := g.save(ctx, inputParam.Data)
		if nil != err {
			blog.Errorf("request(%s): it is to failed to create a new model attribute group (%#v), error info is %s", ctx.ReqID, inputParam.Data, err.Error())
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
	cond.Element(&mongo.Eq{Key: metadata.GroupFieldSupplierAccount, Val: ctx.SupplierAccount})
	cond.Element(&mongo.Eq{Key: metadata.GroupFieldID, Val: existsGroup.ID})

	cnt, err := g.update(ctx, mapstr.NewFromStruct(inputParam.Data, "field"), cond)
	if nil != err {
		blog.Errorf("request(%s): it is failed to update the model attribute group (%#v) by the condition (%#v), err: %s", ctx.ReqID, g, cond, err)
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

func (g *modelAttributeGroup) UpdateModelAttributeGroup(ctx core.ContextParams, objID string, inputParam metadata.UpdateOption) (*metadata.UpdatedCount, error) {

	if err := g.model.isValid(ctx, objID); nil != err {
		blog.Errorf("request(%s): it is failed to create a model attribute group (%#v), error info is %s", ctx.ReqID, inputParam.Data, err.Error())
		return &metadata.UpdatedCount{}, err
	}

	cond, err := mongo.NewConditionFromMapStr(util.SetModOwner(inputParam.Condition.ToMapInterface(), ctx.SupplierAccount))
	if nil != err {
		blog.Errorf("request(%s): it is failed to convert the condition (%#v) from mapstr to condition, error info is %s", ctx.ReqID, inputParam.Condition, err.Error())
		return &metadata.UpdatedCount{}, err
	}
	cond.Element(&mongo.Eq{Key: metadata.GroupFieldObjectID, Val: objID})

	inputParam.Data.Remove(metadata.GroupFieldGroupID)
	inputParam.Data.Remove(metadata.GroupFieldObjectID)
	inputParam.Data.Remove(metadata.GroupFieldSupplierAccount)
	inputParam.Data.Remove(metadata.GroupFieldIsPre)

	if name, exists := inputParam.Data.Get("bk_group_name"); exists {
		name := name.(string)
		queryCond := metadata.QueryCondition{
			Condition: cond.ToMapStr(),
		}
		resp, err := g.SearchModelAttributeGroupByCondition(ctx, queryCond)
		if nil != err {
			blog.Errorf("request(%s): it is to failed to check the group name (%s) if it is exists, error info is %s", ctx.ReqID, name, err.Error())
			return &metadata.UpdatedCount{}, err
		}
		for _, item := range resp.Info {
			if item.GroupName == name {
				continue
			}
			_, exists, err := g.groupNameIsExists(ctx, item.ObjectID, name, metadata.Metadata{Label: metadata.Label{}})
			if nil != err {
				blog.Errorf("request(%s): it is to failed to check the group name (%s) if it is exists, error info is %s", ctx.ReqID, name, err.Error())
				return &metadata.UpdatedCount{}, err
			}
			if exists {
				blog.Errorf("request(%s): it is to failed to update the group name, because (%s) exists", ctx.ReqID, name)
				return &metadata.UpdatedCount{}, ctx.Error.Errorf(common.CCErrCommDuplicateItem, name)
			}
		}
	}
	cnt, err := g.update(ctx, inputParam.Data, cond)
	if nil != err {
		blog.Errorf("request(%s): it is failed to update the data (%s) by the condition (%#v), error info is %s", ctx.ReqID, inputParam.Data, err.Error())
		return &metadata.UpdatedCount{}, err
	}

	return &metadata.UpdatedCount{Count: cnt}, nil
}

func (g *modelAttributeGroup) UpdateModelAttributeGroupByCondition(ctx core.ContextParams, inputParam metadata.UpdateOption) (*metadata.UpdatedCount, error) {

	cond, err := mongo.NewConditionFromMapStr(util.SetModOwner(inputParam.Condition.ToMapInterface(), ctx.SupplierAccount))
	if nil != err {
		blog.Errorf("request(%s): it is failed to convert the condition (%#v) from mapstr to condition, error info is %s", ctx.ReqID, inputParam.Condition, err.Error())
		return &metadata.UpdatedCount{}, err
	}

	inputParam.Data.Remove(metadata.GroupFieldGroupID)
	inputParam.Data.Remove(metadata.GroupFieldObjectID)
	inputParam.Data.Remove(metadata.GroupFieldSupplierAccount)
	inputParam.Data.Remove(metadata.GroupFieldIsPre)

	if name, exists := inputParam.Data.Get("bk_group_name"); exists {
		name := name.(string)
		queryCond := metadata.QueryCondition{
			Condition: cond.ToMapStr(),
		}
		resp, err := g.SearchModelAttributeGroupByCondition(ctx, queryCond)
		if nil != err {
			blog.Errorf("request(%s): it is to failed to check the group name (%s) if it is exists, error info is %s", ctx.ReqID, name, err.Error())
			return &metadata.UpdatedCount{}, err
		}
		for _, item := range resp.Info {
			if item.GroupName == name {
				continue
			}
			_, exists, err := g.groupNameIsExists(ctx, item.ObjectID, name, metadata.Metadata{Label: metadata.Label{}})
			if nil != err {
				blog.Errorf("request(%s): it is to failed to check the group name (%s) if it is exists, error info is %s", ctx.ReqID, name, err.Error())
				return &metadata.UpdatedCount{}, err
			}
			if exists {
				blog.Errorf("request(%s): it is to failed to update the group name, because (%s) exists", ctx.ReqID, name)
				return &metadata.UpdatedCount{}, ctx.Error.Errorf(common.CCErrCommDuplicateItem, name)
			}
		}
	}

	cnt, err := g.update(ctx, inputParam.Data, cond)
	if nil != err {
		blog.Errorf("request(%s): it is failed to update the data (%s) by the condition (%#v), error info is %s", ctx.ReqID, inputParam.Data, err.Error())
		return &metadata.UpdatedCount{}, err
	}

	return &metadata.UpdatedCount{Count: cnt}, nil
}

func (g *modelAttributeGroup) SearchModelAttributeGroup(ctx core.ContextParams, objID string, inputParam metadata.QueryCondition) (*metadata.QueryModelAttributeGroupDataResult, error) {

	dataResult := &metadata.QueryModelAttributeGroupDataResult{
		Info: []metadata.Group{},
	}

	// TODO: Paging queries needs to be implemented
	// NOW: Temporarily blocked
	// if err := g.model.isValid(ctx, objID); nil != err {
	// 	blog.Errorf("request(%s): it is failed to query a model by the condition(%#v), error info is %s", ctx.ReqID, inputParam.Condition, err.Error())
	// 	return dataResult, err
	// }

	cond, err := mongo.NewConditionFromMapStr(util.SetQueryOwner(inputParam.Condition.ToMapInterface(), ctx.SupplierAccount))
	if nil != err {
		blog.Errorf("request(%s): it is failed to convert the condition (%#v) from mapstr to condition, error info is %s", ctx.ReqID, inputParam.Condition, err.Error())
		return dataResult, err
	}
	cond.Element(&mongo.Eq{Key: metadata.GroupFieldObjectID, Val: objID})

	totalCount, err := g.count(ctx, cond)
	if nil != err {
		blog.Errorf("request(%s): it is failed to count by the condition (%#v), error info is %s ", ctx.ReqID, cond.ToMapStr(), err.Error())
		return dataResult, err
	}

	grps, err := g.search(ctx, cond)
	if nil != err {
		blog.Errorf("request(%s): it is failed to query a model by the condition(%#v), error info is %s", ctx.ReqID, cond.ToMapStr(), err.Error())
		return dataResult, err
	}

	dataResult.Count = totalCount
	dataResult.Info = grps
	return dataResult, nil

}

func (g *modelAttributeGroup) SearchModelAttributeGroupByCondition(ctx core.ContextParams, inputParam metadata.QueryCondition) (*metadata.QueryModelAttributeGroupDataResult, error) {

	cond, err := mongo.NewConditionFromMapStr(util.SetQueryOwner(inputParam.Condition.ToMapInterface(), ctx.SupplierAccount))
	if nil != err {
		blog.Errorf("request(%s): it is failed to convert the condition (%#v) from mapstr to condition, error info is %s", ctx.ReqID, inputParam.Condition, err.Error())
		return &metadata.QueryModelAttributeGroupDataResult{}, err
	}

	totalCount, err := g.count(ctx, cond)
	if nil != err {
		blog.Errorf("request(%s): it is failed to count by the condition (%#v), error info is %s ", ctx.ReqID, cond.ToMapStr(), err.Error())
		return &metadata.QueryModelAttributeGroupDataResult{}, err
	}

	grps, err := g.search(ctx, cond)
	if nil != err {
		blog.Errorf("request(%s): it is failed to query a model by the condition(%#v), error info is %s", ctx.ReqID, cond.ToMapStr(), err.Error())
		return &metadata.QueryModelAttributeGroupDataResult{}, err
	}

	return &metadata.QueryModelAttributeGroupDataResult{Count: totalCount, Info: grps}, nil

}

// desperated only for old api
func (g *modelAttributeGroup) DeleteModelAttributeGroupByCondition(ctx core.ContextParams, inputParam metadata.DeleteOption) (*metadata.DeletedCount, error) {

	cond, err := mongo.NewConditionFromMapStr(util.SetModOwner(inputParam.Condition.ToMapInterface(), ctx.SupplierAccount))
	if nil != err {
		blog.Errorf("request(%s): it is failed to convert the condition (%#v) from mapstr to condition, error info is %s", ctx.ReqID, inputParam.Condition, err.Error())
		return &metadata.DeletedCount{}, err
	}

	grps, err := g.search(ctx, cond)
	if nil != err {
		blog.Errorf("request(%s): it is failed to query model attribute groups by the condition (%#v), error info is %s", ctx.ReqID, cond.ToMapStr(), err.Error())
		return &metadata.DeletedCount{}, err
	}

	for _, grp := range grps {
		hasAttrs, err := g.hasAttributes(ctx, grp.ObjectID, []string{grp.GroupID})
		if nil != err {
			blog.Errorf("request(%s): it is failed to check whether the group(%s) has some attributes for the model(%s), error info is %s", ctx.ReqID, grp.GroupID, grp.ObjectID, err.Error())
			return &metadata.DeletedCount{}, err
		}
		if hasAttrs {
			blog.Errorf("request(%s): the group(%s) has some attributes, forbidden to delete", ctx.ReqID, grp.GroupID)
			return &metadata.DeletedCount{}, ctx.Error.Error(common.CCErrCoreServiceModelAttributeGroupHasSomeAttributes)
		}
	}

	cnt, err := g.delete(ctx, cond)
	if nil != err {
		blog.Errorf("request(%s): it is failed to delete the group by the condition(%#v), error info is %s", ctx.ReqID, cond.ToMapStr(), err.Error())
		return &metadata.DeletedCount{}, err
	}

	return &metadata.DeletedCount{Count: cnt}, nil
}

func (g *modelAttributeGroup) DeleteModelAttributeGroup(ctx core.ContextParams, objID string, inputParam metadata.DeleteOption) (*metadata.DeletedCount, error) {

	if err := g.model.isValid(ctx, objID); nil != err {
		blog.Errorf("request(%s): it is failed to delete a model by the condition(%#v), error info is %s", ctx.ReqID, inputParam.Condition, err.Error())
		return &metadata.DeletedCount{}, err
	}

	cond, err := mongo.NewConditionFromMapStr(util.SetModOwner(inputParam.Condition.ToMapInterface(), ctx.SupplierAccount))
	if nil != err {
		blog.Errorf("request(%s): it is failed to convert the condition (%#v) from mapstr to condition, error info is %s", ctx.ReqID, inputParam.Condition, err.Error())
		return &metadata.DeletedCount{}, err
	}
	cond.Element(&mongo.Eq{Key: metadata.GroupFieldObjectID, Val: objID})

	grps, err := g.search(ctx, cond)
	if nil != err {
		blog.Errorf("request(%s): it is failed to query model attribute groups by the condition (%#v), error info is %s", ctx.ReqID, cond.ToMapStr(), err.Error())
		return &metadata.DeletedCount{}, err
	}

	grpIDS := make([]string, 0)
	for _, grp := range grps {
		grpIDS = append(grpIDS, grp.GroupID)
	}

	isExists, err := g.hasAttributes(ctx, objID, grpIDS)
	if nil != err {
		blog.Errorf("request(%s): it is failed to check whether the group IDS (%#v) has some attributes for the model(%s), error info is %s", ctx.ReqID, grpIDS, objID, err.Error())
		return &metadata.DeletedCount{}, err
	}

	if isExists {
		return &metadata.DeletedCount{}, ctx.Error.Error(common.CCErrCoreServiceModelAttributeGroupHasSomeAttributes)
	}

	cnt, err := g.delete(ctx, cond)
	if nil != err {
		blog.Errorf("request(%s): it is failed to delete the group by the condition(%#v), error info is %s", ctx.ReqID, cond.ToMapStr(), err.Error())
		return &metadata.DeletedCount{}, err
	}

	return &metadata.DeletedCount{Count: cnt}, nil
}
