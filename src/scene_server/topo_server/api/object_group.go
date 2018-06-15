/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package api

import (
	"fmt"
	"net/http"

	"configcenter/src/common"
	"configcenter/src/common/condition"
	frtypes "configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/topo_server/core/types"
)

func init() {
	apiInst.initFuncs = append(apiInst.initFuncs, apiInst.initObjectGroup)
}

func (cli *topoAPI) initObjectGroup() {
	cli.actions = append(cli.actions, action{Method: http.MethodPost, Path: "/objectatt/group/new", HandlerFunc: cli.CreateObjectGroup})
	cli.actions = append(cli.actions, action{Method: http.MethodPut, Path: "/objectatt/group/update", HandlerFunc: cli.UpdateObjectGroup})
	cli.actions = append(cli.actions, action{Method: http.MethodDelete, Path: "/objectatt/group/groupid/{id}", HandlerFunc: cli.DeleteObjectGroup})
	cli.actions = append(cli.actions, action{Method: http.MethodPut, Path: "/objectatt/group/property", HandlerFunc: cli.UpdateObjectAttributeGroup})
	cli.actions = append(cli.actions, action{Method: http.MethodDelete, Path: "/objectatt/group/owner/{owner_id}/object/{object_id}/propertyids/{property_id}/groupids/{group_id}", HandlerFunc: cli.DeleteObjectAttributeGroup})
	cli.actions = append(cli.actions, action{Method: http.MethodPost, Path: "/objectatt/group/property/owner/{owner_id}/object/{object_id}", HandlerFunc: cli.SearchGroupByObject})
}

// CreateObjectGroup create a new object grouphttps://github.com/Tencent/bk-cmdb/blob/apimachinery/src/common/metadata/object_controller.go
func (cli *topoAPI) CreateObjectGroup(params types.LogicParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (interface{}, error) {

	rsp, err := cli.core.GroupOperation().CreateObjectGroup(params, data)
	if nil != err {
		return nil, err
	}

	return rsp.ToMapStr()
}

// UpdateObjectGroup update the object group information
func (cli *topoAPI) UpdateObjectGroup(params types.LogicParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (interface{}, error) {

	fmt.Println("UpdateObjectGroup")

	cond := condition.CreateCondition()
	val, exists := data.Get("id")

	if exists {
		cond.Field("id").Eq(val)
	} else {
		val, exists = data.Get(metadata.GroupFieldGroupID)
		if !exists {
			return nil, params.Err.Errorf(common.CCErrCommParamsLostField, metadata.GroupFieldGroupID)
		}
		cond.Field(metadata.GroupFieldGroupID).Eq(val)
	}

	err := cli.core.GroupOperation().UpdateObjectGroup(params, data, cond)
	if nil != err {
		return nil, err
	}

	return nil, nil
}

// DeleteObjectGroup delete the object group
func (cli *topoAPI) DeleteObjectGroup(params types.LogicParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (interface{}, error) {
	fmt.Println("DeleteObjectGroup")
	cond := condition.CreateCondition()
	val, exists := data.Get("id")

	if exists {
		cond.Field("id").Eq(val)
	} else {
		val, exists = data.Get(metadata.GroupFieldGroupID)
		if !exists {
			return nil, params.Err.Errorf(common.CCErrCommParamsLostField, metadata.GroupFieldGroupID)
		}
		cond.Field(metadata.GroupFieldGroupID).Eq(val)
	}

	err := cli.core.GroupOperation().DeleteObjectGroup(params, cond)
	if nil != err {
		return nil, err
	}

	return nil, nil
}

// UpdateObjectAttributeGroup update the object attribute belongs to group information
func (cli *topoAPI) UpdateObjectAttributeGroup(params types.LogicParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (interface{}, error) {
	fmt.Println("UpdateObjectAttributeGroup")
	cond := condition.CreateCondition()

	cond.Field(common.BKOwnerIDField).Eq(params.Header.OwnerID)
	val, exists := data.Get(metadata.ModelFieldObjectID)
	if !exists {
		return nil, params.Err.Errorf(common.CCErrCommParamsLostField, metadata.ModelFieldObjectID)
	}
	cond.Field(metadata.ModelFieldObjectID).Eq(val)
	data.Remove(metadata.ModelFieldObjectID)

	val, exists = data.Get(metadata.AttributeFieldPropertyID)
	if !exists {
		return nil, params.Err.Errorf(common.CCErrCommParamsLostField, metadata.AttributeFieldPropertyID)
	}
	cond.Field(metadata.AttributeFieldPropertyID).Eq(val)
	data.Remove(metadata.AttributeFieldPropertyID)

	err := cli.core.AttributeOperation().UpdateObjectAttribute(params, data, cond)
	if nil != err {
		return nil, err
	}

	return nil, nil
}

// DeleteObjectAttributeGroup delete the object attribute belongs to group information
func (cli *topoAPI) DeleteObjectAttributeGroup(params types.LogicParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (interface{}, error) {
	fmt.Println("DeleteObjectAttributeGroup")
	cond := condition.CreateCondition()

	cond.Field(common.BKOwnerIDField).Eq(params.Header.OwnerID)
	cond.Field(metadata.ModelFieldObjectID).Eq(pathParams("object_id"))
	cond.Field(metadata.AttributeFieldPropertyID).Eq(pathParams("property_id"))
	cond.Field(metadata.AttributeFieldPropertyGroup).Eq("group_id")

	innerData := frtypes.MapStr{}
	innerData.Set(metadata.AttributeFieldPropertyGroup, "none")
	innerData.Set(metadata.AttributeFieldPropertyIndex, -1)

	err := cli.core.AttributeOperation().UpdateObjectAttribute(params, innerData, cond)
	if nil != err {
		return nil, err
	}

	return nil, nil
}

// SearchGroupByObject search the groups by the object
func (cli *topoAPI) SearchGroupByObject(params types.LogicParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (interface{}, error) {
	fmt.Println("SearchGroupByObject")
	cond := condition.CreateCondition()
	cond.Field(metadata.ModelFieldObjectID).Eq(pathParams("object_id"))

	items, err := cli.core.GroupOperation().FindGroupByObject(params, cond)
	if nil != err {
		return nil, err
	}

	result := frtypes.MapStr{}
	result.Set("data", items)

	return result, nil
}
