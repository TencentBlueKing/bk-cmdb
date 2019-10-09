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

package extensions

import (
	"context"
	"fmt"
	"net/http"

	"configcenter/src/auth/meta"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

func (am *AuthManager) collectAttributesGroupByIDs(ctx context.Context, header http.Header, agIDs ...int64) ([]metadata.Group, error) {
	rid := util.ExtractRequestIDFromContext(ctx)

	// unique ids so that we can be aware of invalid id if query result length not equal ids's length
	agIDs = util.IntArrayUnique(agIDs)

	// get attribute group by objID
	cond := condition.CreateCondition().Field(common.BKFieldID).In(agIDs)
	queryCond := &metadata.QueryCondition{Condition: cond.ToMapStr()}
	resp, err := am.clientSet.CoreService().Instance().ReadInstance(ctx, header, common.BKTableNamePropertyGroup, queryCond)
	if err != nil {
		return nil, fmt.Errorf("get model attribute group by id: %+v failed, err: %+v", agIDs, err)
	}
	if len(resp.Data.Info) == 0 {
		return nil, fmt.Errorf("get model attribute group by id: %+v failed, not found", agIDs)
	}
	if len(resp.Data.Info) != len(agIDs) {
		return nil, fmt.Errorf("get model attribute group by id: %+v failed, get %d, expect %d", agIDs, len(resp.Data.Info), len(agIDs))
	}

	pgs := make([]metadata.Group, 0)
	for _, item := range resp.Data.Info {
		pg := metadata.Group{}
		_, err := pg.Parse(item)
		if err != nil {
			blog.Errorf("collectAttributesGroupByAttributeIDs %+v failed, parse attribute group %+v failed, err: %+v, rid: %s", agIDs, item, err, rid)
			return nil, fmt.Errorf("parse attribute group from db data failed, err: %+v", err)
		}
		pgs = append(pgs, pg)
	}
	return pgs, nil
}

func (am *AuthManager) makeResourceByAttributeGroup(ctx context.Context, header http.Header, action meta.Action, attributeGroups ...metadata.Group) ([]meta.ResourceAttribute, error) {
	rid := util.ExtractRequestIDFromContext(ctx)

	businessID, err := am.ExtractBusinessIDFromAttributeGroup(attributeGroups...)
	if err != nil {
		return nil, fmt.Errorf("extract business id from attribute groups failed, err: %+v", err)
	}

	objectIDs := make([]string, 0)
	for _, attributeGroup := range attributeGroups {
		objectIDs = append(objectIDs, attributeGroup.ObjectID)
	}

	objectIDs = util.StrArrayUnique(objectIDs)
	if len(objectIDs) > 1 {
		blog.Errorf("makeResourceByAttributeGroup failed, model attribute group belongs to multiple models: %+v, rid: %s", objectIDs, rid)
		return nil, fmt.Errorf("model attribute groups belongs to multiple models: %+v", objectIDs)
	}
	if len(objectIDs) == 0 {
		blog.Errorf("makeResourceByAttributeGroup failed, model id not found, attribute groups: %+v, rid: %s", attributeGroups, rid)
		return nil, fmt.Errorf("model id not found")
	}

	objects, err := am.collectObjectsByObjectIDs(ctx, header, businessID, objectIDs...)
	if err != nil {
		blog.Errorf("makeResourceByAttributeGroup failed, collectObjectsByObjectIDs failed, objectIDs: %+v, err: %+v", objectIDs, err)
		return nil, fmt.Errorf("collect object id failed, err: %+v", err)
	}

	if len(objects) == 0 {
		blog.Errorf("makeResourceByAttributeGroup failed, collectObjectsByObjectIDs no objects found, objectIDs: %+v, err: %+v, rid: %s", objectIDs, err, rid)
		return nil, fmt.Errorf("collect object by id not found")
	}

	parentResources, err := am.MakeResourcesByObjects(ctx, header, meta.EmptyAction, objects...)
	if err != nil {
		blog.Errorf("makeResourceByAttributeGroup failed, get parent resource failed, objects: %+v, err: %+v, rid: %s", objects, err, rid)
		return nil, fmt.Errorf("get parent resources failed, objects: %+v, err: %+v", objects, err)
	}
	if len(parentResources) > 1 {
		blog.Errorf("makeResourceByAttributeGroup failed, get multiple parent resource, parent resources: %+v, rid: %s", parentResources, rid)
		return nil, fmt.Errorf("get multiple parent resources, parent resources: %+v", parentResources)
	}
	if len(parentResources) == 0 {
		blog.Errorf("makeResourceByAttributeGroup failed, get parent resource empty, objects: %+v, rid: %s", objects, rid)
		return nil, fmt.Errorf("get parent resources empty, objects: %+v", objects)
	}

	parentResource := parentResources[0]

	// prepare resource layers for authorization
	resources := make([]meta.ResourceAttribute, 0)
	for _, attributeGroup := range attributeGroups {
		parentLayers := parentResource.Layers
		parentLayers = append(parentLayers, meta.Item{
			Type:       meta.Model,
			Name:       parentResource.Name,
			InstanceID: parentResource.InstanceID,
		})

		// attribute
		resource := meta.ResourceAttribute{
			Basic: meta.Basic{
				Action:     action,
				Type:       meta.ModelAttributeGroup,
				Name:       attributeGroup.GroupName,
				InstanceID: attributeGroup.ID,
			},
			SupplierAccount: util.GetOwnerID(header),
			BusinessID:      businessID,
			Layers:          parentLayers,
		}

		resources = append(resources, resource)
	}

	blog.V(5).Infof("makeResourceByAttributeGroup result: %+v, rid: %s", resources, rid)
	return resources, nil
}

func (am *AuthManager) ExtractBusinessIDFromAttributeGroup(attributeGroups ...metadata.Group) (int64, error) {
	if len(attributeGroups) == 0 {
		return 0, fmt.Errorf("no object found")
	}

	businessIDs := make([]int64, 0)
	for _, attributeGroup := range attributeGroups {
		bizID, err := metadata.BizIDFromMetadata(attributeGroup.Metadata)
		if err != nil {
			return 0, fmt.Errorf("parse business id failed, err: %+v", err)
		}
		businessIDs = append(businessIDs, bizID)
	}

	businessIDs = util.IntArrayUnique(businessIDs)
	if len(businessIDs) > 1 {
		return 0, fmt.Errorf("attribute groups belongs to multiple business: [%+v]", businessIDs)
	}

	if len(businessIDs) == 0 {
		return 0, fmt.Errorf("unexpected error, no business found with attribute groups: %+v", attributeGroups)
	}
	return businessIDs[0], nil
}

func (am *AuthManager) RegisterModelAttributeGroup(ctx context.Context, header http.Header, attributeGroups ...metadata.Group) error {
	if am.Enabled() == false {
		return nil
	}

	if len(attributeGroups) == 0 {
		return nil
	}

	if am.RegisterModelAttributeEnabled == false {
		return nil
	}

	resources, err := am.makeResourceByAttributeGroup(ctx, header, meta.EmptyAction, attributeGroups...)
	if err != nil {
		return fmt.Errorf("register model attribute group failed, err: %+v", err)
	}

	return am.Authorize.RegisterResource(ctx, resources...)
}

func (am *AuthManager) DeregisterModelAttributeGroup(ctx context.Context, header http.Header, attributeGroups ...metadata.Group) error {
	if am.Enabled() == false {
		return nil
	}

	if len(attributeGroups) == 0 {
		return nil
	}

	if am.RegisterModelAttributeEnabled == false {
		return nil
	}

	resources, err := am.makeResourceByAttributeGroup(ctx, header, meta.EmptyAction, attributeGroups...)
	if err != nil {
		return fmt.Errorf("deregister model attribute group failed, err: %+v", err)
	}

	return am.Authorize.DeregisterResource(ctx, resources...)
}

func (am *AuthManager) DeregisterModelAttributeGroupByID(ctx context.Context, header http.Header, groupIDs ...int64) error {
	if am.Enabled() == false {
		return nil
	}

	if len(groupIDs) == 0 {
		return nil
	}

	if am.RegisterModelAttributeEnabled == false {
		return nil
	}

	attributeGroups, err := am.collectAttributesGroupByIDs(ctx, header, groupIDs...)
	if err != nil {
		return fmt.Errorf("deregistered model attribute group failed, get model attribute group by id failed, err: %+v", err)
	}
	return am.DeregisterModelAttributeGroup(ctx, header, attributeGroups...)
}

// func (am *AuthManager) AuthorizeModelAttributeGroup(ctx context.Context, header http.Header, action meta.Action, attributeGroups ...metadata.Group) error {
// 	if am.Enabled() == false {
// 		return nil
// 	}
//
// 	businessID, err := am.ExtractBusinessIDFromAttributeGroup(attributeGroups...)
// 	if err != nil {
// 		return fmt.Errorf("extract business id from attribute groups failed, err: %+v", err)
// 	}
//
// 	if am.RegisterModelAttributeEnabled == false {
// 		objectIDs := make([]string, 0)
// 		for _, attributeGroup := range attributeGroups {
// 			objectIDs = append(objectIDs, attributeGroup.ObjectID)
// 		}
// 		return am.AuthorizeByObjectID(ctx, header, action, businessID, objectIDs...)
// 	}
//
// 	resources, err := am.makeResourceByAttributeGroup(ctx, header, action, attributeGroups...)
// 	if err != nil {
// 		return fmt.Errorf("authorize model attribute failed, err: %+v", err)
// 	}
//
// 	return am.batchAuthorize(ctx, header, resources...)
// }

func (am *AuthManager) UpdateRegisteredModelAttributeGroup(ctx context.Context, header http.Header, attributeGroups ...metadata.Group) error {
	if am.Enabled() == false {
		return nil
	}

	if len(attributeGroups) == 0 {
		return nil
	}

	if am.RegisterModelAttributeEnabled == false {
		return nil
	}
	resources, err := am.makeResourceByAttributeGroup(ctx, header, meta.EmptyAction, attributeGroups...)
	if err != nil {
		return fmt.Errorf("update registered model attribute failed, err: %+v", err)
	}

	for _, resource := range resources {
		err = am.Authorize.UpdateResource(ctx, &resource)
		if err != nil {
			return err
		}
	}

	return nil
}

func (am *AuthManager) UpdateRegisteredModelAttributeGroupByID(ctx context.Context, header http.Header, attributeIDs ...int64) error {
	if am.Enabled() == false {
		return nil
	}

	if len(attributeIDs) == 0 {
		return nil
	}

	if am.RegisterModelAttributeEnabled == false {
		return nil
	}
	attributeGroups, err := am.collectAttributesGroupByIDs(ctx, header, attributeIDs...)
	if err != nil {
		return fmt.Errorf("update registered model attribute group failed, get attribute by id failed, err: %+v", err)
	}
	return am.UpdateRegisteredModelAttributeGroup(ctx, header, attributeGroups...)
}

// func (am *AuthManager) AuthorizeAttributeGroupByID(ctx context.Context, header http.Header, action meta.Action, attributeIDs ...int64) error {
// 	if am.Enabled() == false {
// 		return nil
// 	}
//
// 	attributeGroups, err := am.collectAttributesGroupByIDs(ctx, header, attributeIDs...)
// 	if err != nil {
// 		return fmt.Errorf("get attributes by id failed, err: %+v", err)
// 	}
// 	return am.AuthorizeModelAttributeGroup(ctx, header, action, attributeGroups...)
// }
