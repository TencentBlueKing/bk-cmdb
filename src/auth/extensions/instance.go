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

/*
 * instance represent common instances here
 */

func (am *AuthManager) CollectInstancesByModelID(ctx context.Context, header http.Header, objectID string) ([]InstanceSimplify, error) {
	cond := metadata.QueryCondition{
		Condition: condition.CreateCondition().Field(common.BKObjIDField).Eq(objectID).ToMapStr(),
	}
	result, err := am.clientSet.CoreService().Instance().ReadInstance(ctx, header, common.BKInnerObjIDObject, &cond)
	if err != nil {
		blog.V(3).Infof("get instances by model id %s failed, err: %+v", objectID, err)
		return nil, fmt.Errorf("get instances by model id %s failed, err: %+v", objectID, err)
	}
	instances := make([]InstanceSimplify, 0)
	for _, cls := range result.Data.Info {
		instance := InstanceSimplify{}
		_, err = instance.Parse(cls)
		if err != nil {
			return nil, fmt.Errorf("get classication by object failed, err: %+v", err)
		}
		instances = append(instances, instance)
	}
	return instances, nil
}

func (am *AuthManager) collectInstancesByInstanceIDs(ctx context.Context, header http.Header, instanceIDs ...string) ([]InstanceSimplify, error) {
	// unique ids so that we can be aware of invalid id if query result length not equal ids's length
	instanceIDs = util.StrArrayUnique(instanceIDs)

	cond := metadata.QueryCondition{
		Condition: condition.CreateCondition().Field(common.BKInstIDField).In(instanceIDs).ToMapStr(),
	}
	result, err := am.clientSet.CoreService().Instance().ReadInstance(ctx, header, common.BKInnerObjIDObject, &cond)
	if err != nil {
		blog.V(3).Infof("get instances by id failed, err: %+v", err)
		return nil, fmt.Errorf("get instances by id failed, err: %+v", err)
	}
	instances := make([]InstanceSimplify, 0)
	for _, cls := range result.Data.Info {
		instance := InstanceSimplify{}
		_, err = instance.Parse(cls)
		if err != nil {
			return nil, fmt.Errorf("get classication by object failed, err: %+v", err)
		}
		instances = append(instances, instance)
	}
	return instances, nil
}

func (am *AuthManager) collectInstancesByRawIDs(ctx context.Context, header http.Header, ids ...int64) ([]InstanceSimplify, error) {
	// unique ids so that we can be aware of invalid id if query result length not equal ids's length
	ids = util.IntArrayUnique(ids)
	cond := metadata.QueryCondition{
		Condition: condition.CreateCondition().Field(common.BKFieldID).In(ids).ToMapStr(),
	}
	result, err := am.clientSet.CoreService().Instance().ReadInstance(ctx, header, common.BKInnerObjIDObject, &cond)
	if err != nil {
		blog.V(3).Infof("get classification by id failed, err: %+v", err)
		return nil, fmt.Errorf("get classification by id failed, err: %+v", err)
	}
	instances := make([]InstanceSimplify, 0)
	for _, cls := range result.Data.Info {
		classification := InstanceSimplify{}
		_, err = classification.Parse(cls)
		if err != nil {
			return nil, fmt.Errorf("get classication by object failed, err: %+v", err)
		}
		instances = append(instances, classification)
	}
	return instances, nil
}

func (am *AuthManager) extractBusinessIDFromInstances(classifications ...InstanceSimplify) (int64, error) {
	var businessID int64
	for idx, instance := range classifications {
		bizID := instance.BizID
		// we should ignore metadata.LabelBusinessID field not found error
		if idx > 0 && bizID != businessID {
			return 0, fmt.Errorf("authorization failed, get multiple business ID from objects")
		}
		businessID = bizID
	}
	return businessID, nil
}

func (am *AuthManager) MakeResourcesByInstances(header http.Header, action meta.Action, businessID int64, instances ...InstanceSimplify) []meta.ResourceAttribute {
	resources := make([]meta.ResourceAttribute, 0)
	for _, instance := range instances {
		resource := meta.ResourceAttribute{
			Basic: meta.Basic{
				Action:     action,
				Type:       meta.ModelInstance,
				Name:       instance.Name,
				InstanceID: instance.ID,
			},
			SupplierAccount: util.GetOwnerID(header),
			BusinessID:      businessID,
		}

		resources = append(resources, resource)
	}
	return resources
}

func (am *AuthManager) AuthorizeByInstances(ctx context.Context, header http.Header, action meta.Action, instances ...InstanceSimplify) error {

	// extract business id
	bizID, err := am.extractBusinessIDFromInstances(instances...)
	if err != nil {
		return fmt.Errorf("authorize instances failed, extract business id from instance failed, err: %+v", err)
	}

	// make auth resources
	resources := am.MakeResourcesByInstances(header, action, bizID, instances...)

	return am.authorize(ctx, header, bizID, resources...)
}

func (am *AuthManager) UpdateRegisteredInstances(ctx context.Context, header http.Header, instances ...InstanceSimplify) error {
	// extract business id
	bizID, err := am.extractBusinessIDFromInstances(instances...)
	if err != nil {
		return fmt.Errorf("authorize instances failed, extract business id from instances failed, err: %+v", err)
	}

	// make auth resources
	resources := am.MakeResourcesByInstances(header, meta.EmptyAction, bizID, instances...)

	for _, resource := range resources {
		if err := am.Authorize.UpdateResource(ctx, &resource); err != nil {
			return err
		}
	}

	return nil
}

func (am *AuthManager) UpdateRegisteredInstanceByID(ctx context.Context, header http.Header, ids ...int64) error {
	instances, err := am.collectInstancesByRawIDs(ctx, header, ids...)
	if err != nil {
		return fmt.Errorf("update registered classifications failed, get classfication by id failed, err: %+v", err)
	}
	return am.UpdateRegisteredInstances(ctx, header, instances...)
}

func (am *AuthManager) UpdateRegisteredInstanceByRawID(ctx context.Context, header http.Header, ids ...int64) error {
	instances, err := am.collectInstancesByRawIDs(ctx, header, ids...)
	if err != nil {
		return fmt.Errorf("update registered classifications failed, get classfication by id failed, err: %+v", err)
	}
	return am.UpdateRegisteredInstances(ctx, header, instances...)
}

func (am *AuthManager) DeregisterInstanceByRawID(ctx context.Context, header http.Header, ids ...int64) error {
	instances, err := am.collectClassificationsByRawIDs(ctx, header, ids...)
	if err != nil {
		return fmt.Errorf("deregister instance failed, get instance by id failed, err: %+v", err)
	}
	return am.DeregisterClassification(ctx, header, instances...)
}

func (am *AuthManager) RegisterInstancesByID(ctx context.Context, header http.Header, ids ...int64) error {
	instances, err := am.collectInstancesByRawIDs(ctx, header, ids...)
	if err != nil {
		return fmt.Errorf("register instances failed, get instance by id failed, err: %+v", err)
	}
	return am.RegisterInstances(ctx, header, instances...)
}

func (am *AuthManager) RegisterInstances(ctx context.Context, header http.Header, instances ...InstanceSimplify) error {

	// extract business id
	bizID, err := am.extractBusinessIDFromInstances(instances...)
	if err != nil {
		return fmt.Errorf("register classifications failed, extract business id from classification failed, err: %+v", err)
	}

	// make auth resources
	resources := am.MakeResourcesByInstances(header, meta.EmptyAction, bizID, instances...)

	return am.Authorize.RegisterResource(ctx, resources...)
}

func (am *AuthManager) DeregisterInstances(ctx context.Context, header http.Header, instances ...InstanceSimplify) error {

	// extract business id
	bizID, err := am.extractBusinessIDFromInstances(instances...)
	if err != nil {
		return fmt.Errorf("deregister classifications failed, extract business id from classification failed, err: %+v", err)
	}

	// make auth resources
	resources := am.MakeResourcesByInstances(header, meta.EmptyAction, bizID, instances...)

	return am.Authorize.DeregisterResource(ctx, resources...)
}
