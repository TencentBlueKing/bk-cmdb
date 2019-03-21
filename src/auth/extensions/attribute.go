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

func (am *AuthManager) collectAttributesByAttributeIDs(ctx context.Context, header http.Header, attributeIDs ...int64) ([]metadata.Attribute, error) {
	// get model by objID
	cond := condition.CreateCondition().Field(common.BKFieldID).In(attributeIDs)
	queryCond := &metadata.QueryCondition{Condition: cond.ToMapStr()}
	resp, err := am.clientSet.CoreService().Instance().ReadInstance(ctx, header, common.BKTableNameObjAttDes, queryCond)
	if err != nil {
		return nil, fmt.Errorf("get attribute by id: %+v failed, err: %+v", attributeIDs, err)
	}
	if len(resp.Data.Info) == 0 {
		return nil, fmt.Errorf("get attribute by id: %+v failed, not found", attributeIDs)
	}
	if len(resp.Data.Info) != len(attributeIDs) {
		return nil, fmt.Errorf("get attribute by id: %+v failed, get %d, expect %d", attributeIDs, len(resp.Data.Info), len(attributeIDs))
	}

	attributes := make([]metadata.Attribute, 0)
	for _, item := range resp.Data.Info {
		attribute := metadata.Attribute{}
		attribute.Parse(item)
		attributes = append(attributes, attribute)
	}
	return attributes, nil
}

func (am *AuthManager) makeResourceByAttributes(ctx context.Context, header http.Header, action meta.Action, attributes ...metadata.Attribute) ([]meta.ResourceAttribute, error){
	objectIDs := make([]string, 0)
	for _, attribute := range attributes {
		objectIDs = append(objectIDs, attribute.ObjectID)
	}

	objects, err := am.collectObjectsByObjectIDs(ctx, header, objectIDs...)
	if err != nil {
		return nil, fmt.Errorf("register model attributes failed, get related models failed, err: %+v", err)
	}
	objectMap := map[string]metadata.Object{}
	for _, object := range objects {
		objectMap[object.ObjectID] = object
	}
	
	businessID, err := am.extractBusinessIDFromObjects(objects...)
	if err != nil {
		return nil, fmt.Errorf("make auth resource for model attribute failed, err: %+v", err)
	}

	classificationIDs := make([]string, 0)
	for _, object := range objects {
		classificationIDs = append(classificationIDs, object.ObjCls)
	}
	classifications, err := am.collectClassificationsByClassificationIDs(ctx, header, classificationIDs...)
	if err != nil {
		return nil, fmt.Errorf("register model attributes failed, get related models failed, err: %+v", err)
	}
	classificationMap := map[string]metadata.Classification{}
	for _, classification := range classifications {
		classificationMap[classification.ClassificationID] = classification
	}
	
	// step2 prepare resource layers for authorization
	resources := make([]meta.ResourceAttribute, 0)
	for _, attribute := range attributes {

		object := objectMap[attribute.ObjectID]
		
		// check obj's group id in map
		if _, exist := classificationMap[object.ObjCls]; exist == false {
			blog.V(3).Infof("authorization failed, get classification by object failed, err: bk_classification_id not exist")
			return nil, fmt.Errorf("authorization failed, get classification by object failed, err: bk_classification_id not exist")
		}

		parentLayers := meta.Layers{}
		// model group
		parentLayers = append(parentLayers, meta.Item{
			Type:       meta.Model,
			Name:       classificationMap[object.ObjCls].ClassificationID,
			InstanceID: classificationMap[object.ObjCls].ID,
		})

		// model
		parentLayers = append(parentLayers, meta.Item{
			Type:       meta.Model,
			Name:       object.ObjectID,
			InstanceID: object.ID,
		})

		// attribute
		resource := meta.ResourceAttribute{
			Basic: meta.Basic{
				Action:     action,
				Type:       meta.ModelAttribute,
				Name:       attribute.PropertyName,
				InstanceID: attribute.ID,
			},
			SupplierAccount: util.GetOwnerID(header),
			BusinessID:      businessID,
		}

		resources = append(resources, resource)
	}
	return nil, nil
}

func (am *AuthManager) RegisterModelAttribute(ctx context.Context, header http.Header, attributes ...metadata.Attribute) error {
	resources, err := am.makeResourceByAttributes(ctx, header, meta.EmptyAction, attributes...)
	if err != nil {
		fmt.Errorf("register model attribute failed, err: %+v", err)
	}
	
	return am.Authorize.RegisterResource(ctx, resources...)
}

func (am *AuthManager) DeregisterModelAttribute(ctx context.Context, header http.Header, attributes ...metadata.Attribute) error {
	resources, err := am.makeResourceByAttributes(ctx, header, meta.EmptyAction, attributes...)
	if err != nil {
		fmt.Errorf("deregister model attribute failed, err: %+v", err)
	}

	return am.Authorize.DeregisterResource(ctx, resources...)
}

func (am *AuthManager) DeregisterModelAttributeByID(ctx context.Context, header http.Header, attributeIDs ...int64) error {
	attibutes, err := am.collectAttributesByAttributeIDs(ctx, header, attributeIDs...)
	if err != nil {
		return fmt.Errorf("update registered model attribute failed, get attribute by id failed, err: %+v", err)
	}
	return am.DeregisterModelAttribute(ctx, header, attibutes...)
}

func (am *AuthManager) AuthorizeModelAttribute(ctx context.Context, header http.Header, action meta.Action, attributes ...metadata.Attribute) error {
	resources, err := am.makeResourceByAttributes(ctx, header, action, attributes...)
	if err != nil {
		fmt.Errorf("authorize model attribute failed, err: %+v", err)
	}

	return am.Authorize.RegisterResource(ctx, resources...)
}

func (am *AuthManager) UpdateRegisteredModelAttribute(ctx context.Context, header http.Header, attributes ...metadata.Attribute) error {
	resources, err := am.makeResourceByAttributes(ctx, header, meta.EmptyAction, attributes...)
	if err != nil {
		fmt.Errorf("update registered model attribute failed, err: %+v", err)
	}

	return am.Authorize.RegisterResource(ctx, resources...)
}

func (am *AuthManager) UpdateRegisteredModelAttributeByID(ctx context.Context, header http.Header, attributeIDs ...int64) error {
	attibutes, err := am.collectAttributesByAttributeIDs(ctx, header, attributeIDs...)
	if err != nil {
		return fmt.Errorf("update registered model attribute failed, get attribute by id failed, err: %+v", err)
	}
	return am.UpdateRegisteredModelAttribute(ctx, header, attibutes...)
}

func (am *AuthManager) AuthorizeByAttributeID(ctx context.Context, header http.Header, action meta.Action, attributeIDs ...int64) error {
	attributes, err := am.collectAttributesByAttributeIDs(ctx, header, attributeIDs...)
	if err != nil {
		return fmt.Errorf("get attributes by id failed, err: %+v", err)
	}

	objectIDs := make([]string, 0)
	for _, attribute := range attributes {
		objectIDs = append(objectIDs, attribute.ObjectID)
	}

	return am.AuthorizeByObjectID(ctx, header, action, objectIDs...)
}
