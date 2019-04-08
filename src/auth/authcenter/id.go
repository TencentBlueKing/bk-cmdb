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

package authcenter

import (
	"fmt"
	"strconv"

	"configcenter/src/auth/meta"
)

func GenerateResourceID(resourceType ResourceTypeID, attribute *meta.ResourceAttribute) ([]RscTypeAndID, error) {
	switch attribute.Basic.Type {
	case meta.Business:
		return businessResourceID(resourceType, attribute)
	case meta.Model:
		return modelResourceID(resourceType, attribute)
	case meta.ModelModule:
		return modelModuleResourceID(resourceType, attribute)
	case meta.ModelSet:
		return modelSetResourceID(resourceType, attribute)
	case meta.MainlineModel:
		return mainlineModelResourceID(resourceType, attribute)
	case meta.MainlineModelTopology:
		return mainlineModelTopologyResourceID(resourceType, attribute)
	case meta.MainlineInstanceTopology:
		return mainlineInstanceTopologyResourceID(resourceType, attribute)
	case meta.AssociationType:
		return associationTypeResourceID(resourceType, attribute)
	case meta.ModelAssociation:
		return modelAssociationResourceID(resourceType, attribute)
	case meta.ModelInstanceAssociation:
		return modelInstanceAssociationResourceID(resourceType, attribute)
	case meta.ModelInstance:
		return modelInstanceResourceID(resourceType, attribute)
	case meta.ModelInstanceTopology:
		return modelInstanceTopologyResourceID(resourceType, attribute)
	case meta.ModelTopology:
		return modelTopologyResourceID(resourceType, attribute)
	case meta.ModelClassification:
		return modelClassificationResourceID(resourceType, attribute)
	case meta.ModelAttributeGroup:
		return modelAttributeGroupResourceID(resourceType, attribute)
	case meta.ModelAttribute:
		return modelAttributeResourceID(resourceType, attribute)
	case meta.ModelUnique:
		return modelUniqueResourceID(resourceType, attribute)
	case meta.UserCustom:
		return hostUserCustomResourceID(resourceType, attribute)
	case meta.HostFavorite:
		return hostFavoriteResourceID(resourceType, attribute)
	case meta.Process:
		return processResourceID(resourceType, attribute)
	case meta.NetDataCollector:
		return netDataCollectorResourceID(resourceType, attribute)
	case meta.EventPushing:
		return eventSubscribeResourceID(resourceType, attribute)
	case meta.HostInstance:
		return hostInstanceResourceID(resourceType, attribute)
	case meta.DynamicGrouping:
		return dynamicGroupingResourceID(resourceType, attribute)
	case meta.AuditLog:
		return auditLogResourceID(resourceType, attribute)
	case meta.SystemBase:
		return nil, nil
	default:
		return nil, fmt.Errorf("gen id failed: unsupported resource type: %s", attribute.Type)
	}
}

// generate business related resource id.
func businessResourceID(resourceType ResourceTypeID, attribute *meta.ResourceAttribute) ([]RscTypeAndID, error) {
	if attribute.InstanceID <= 0 {
		return nil, nil
	}
	id := RscTypeAndID{
		ResourceType: resourceType,
		ResourceID:   strconv.FormatInt(attribute.InstanceID, 10),
	}

	return []RscTypeAndID{id}, nil
}

// generate model's resource id, works for app model and model management
// resource type in auth center.
func modelResourceID(resourceType ResourceTypeID, attribute *meta.ResourceAttribute) ([]RscTypeAndID, error) {
	if attribute.InstanceID <= 0 {
		return nil, nil
	}
	id := RscTypeAndID{
		ResourceType: resourceType,
	}
	id.ResourceID = strconv.FormatInt(attribute.InstanceID, 10)

	return []RscTypeAndID{id}, nil
}

// generate module resource id.
func modelModuleResourceID(resourceType ResourceTypeID, attribute *meta.ResourceAttribute) ([]RscTypeAndID, error) {

	return []RscTypeAndID{
		{
			ResourceType: resourceType,
			ResourceID:   fmt.Sprintf("module:%d", attribute.InstanceID),
		},
	}, nil
}

func modelSetResourceID(resourceType ResourceTypeID, attribute *meta.ResourceAttribute) ([]RscTypeAndID, error) {

	return []RscTypeAndID{
		{
			ResourceType: resourceType,
			ResourceID:   fmt.Sprintf("set:%d", attribute.InstanceID),
		},
	}, nil
}

func mainlineModelResourceID(resourceType ResourceTypeID, attribute *meta.ResourceAttribute) ([]RscTypeAndID, error) {

	return nil, nil
}

func mainlineModelTopologyResourceID(resourceType ResourceTypeID, attribute *meta.ResourceAttribute) ([]RscTypeAndID, error) {
	return nil, nil
}

func mainlineInstanceTopologyResourceID(resourceType ResourceTypeID, attribute *meta.ResourceAttribute) ([]RscTypeAndID, error) {

	return nil, nil
}

func modelAssociationResourceID(resourceType ResourceTypeID, attribute *meta.ResourceAttribute) ([]RscTypeAndID, error) {

	return nil, nil
}

func associationTypeResourceID(resourceType ResourceTypeID, attribute *meta.ResourceAttribute) ([]RscTypeAndID, error) {
	if attribute.InstanceID <= 0 {
		return nil, nil
	}
	id := RscTypeAndID{
		ResourceType: resourceType,
		ResourceID:   strconv.FormatInt(attribute.InstanceID, 10),
	}

	return []RscTypeAndID{id}, nil
}

func modelInstanceAssociationResourceID(resourceType ResourceTypeID, attribute *meta.ResourceAttribute) ([]RscTypeAndID, error) {

	return nil, nil
}

func modelInstanceResourceID(resourceType ResourceTypeID, attribute *meta.ResourceAttribute) ([]RscTypeAndID, error) {
	if attribute.InstanceID <= 0 {
		return nil, nil
	}

	if len(attribute.Layers) < 1 {
		return nil, NotEnoughLayer
	}

	// groupType := SysModelGroup
	modelType := SysModel
	if attribute.BusinessID > 0 {
		// groupType = BizModelGroup
		modelType = BizModel
	}

	return []RscTypeAndID{
		{
			ResourceType: modelType,
			ResourceID:   strconv.FormatInt(attribute.Layers[0].InstanceID, 10),
		},
		{
			ResourceType: resourceType,
			ResourceID:   strconv.FormatInt(attribute.InstanceID, 10),
		},
	}, nil
}

func modelInstanceTopologyResourceID(resourceType ResourceTypeID, attribute *meta.ResourceAttribute) ([]RscTypeAndID, error) {

	return nil, nil
}

func modelTopologyResourceID(resourceType ResourceTypeID, attribute *meta.ResourceAttribute) ([]RscTypeAndID, error) {

	return nil, nil
}

func modelClassificationResourceID(resourceType ResourceTypeID, attribute *meta.ResourceAttribute) ([]RscTypeAndID, error) {
	if attribute.InstanceID <= 0 {
		return nil, nil
	}
	id := RscTypeAndID{
		ResourceType: resourceType,
		ResourceID:   strconv.FormatInt(attribute.InstanceID, 10),
	}
	return []RscTypeAndID{id}, nil
}

func modelAttributeGroupResourceID(resourceType ResourceTypeID, attribute *meta.ResourceAttribute) ([]RscTypeAndID, error) {
	if len(attribute.Layers) < 2 {
		return nil, NotEnoughLayer
	}
	id := RscTypeAndID{
		ResourceType: SysModel,
	}
	if attribute.BusinessID > 0 {
		id.ResourceType = BizModel
	}
	id.ResourceID = strconv.FormatInt(attribute.Layers[len(attribute.Layers)-1].InstanceID, 10)
	return []RscTypeAndID{id}, nil
}

func modelAttributeResourceID(resourceType ResourceTypeID, attribute *meta.ResourceAttribute) ([]RscTypeAndID, error) {
	if len(attribute.Layers) < 2 {
		return nil, NotEnoughLayer
	}
	id := RscTypeAndID{
		ResourceType: SysModel,
	}
	if attribute.BusinessID > 0 {
		id.ResourceType = BizModel
	}
	id.ResourceID = strconv.FormatInt(attribute.Layers[len(attribute.Layers)-1].InstanceID, 10)
	return []RscTypeAndID{id}, nil
}

func modelUniqueResourceID(resourceType ResourceTypeID, attribute *meta.ResourceAttribute) ([]RscTypeAndID, error) {
	if len(attribute.Layers) < 2 {
		return nil, NotEnoughLayer
	}
	id := RscTypeAndID{
		ResourceType: SysModel,
	}
	if attribute.BusinessID > 0 {
		id.ResourceType = BizModel
	}
	id.ResourceID = strconv.FormatInt(attribute.Layers[len(attribute.Layers)-1].InstanceID, 10)
	return []RscTypeAndID{id}, nil
}

func hostUserCustomResourceID(resourceType ResourceTypeID, attribute *meta.ResourceAttribute) ([]RscTypeAndID, error) {

	return nil, nil
}

func hostFavoriteResourceID(resourceType ResourceTypeID, attribute *meta.ResourceAttribute) ([]RscTypeAndID, error) {

	return nil, nil
}

func processResourceID(resourceType ResourceTypeID, attribute *meta.ResourceAttribute) ([]RscTypeAndID, error) {
	if attribute.InstanceID <= 0 {
		return nil, nil
	}
	id := RscTypeAndID{
		ResourceType: BizProcessInstance,
		ResourceID:   strconv.FormatInt(attribute.InstanceID, 10),
	}
	return []RscTypeAndID{id}, nil
}

func netDataCollectorResourceID(resourceType ResourceTypeID, attribute *meta.ResourceAttribute) ([]RscTypeAndID, error) {

	return nil, nil
}

func hostInstanceResourceID(resourceType ResourceTypeID, attribute *meta.ResourceAttribute) ([]RscTypeAndID, error) {
	// translate all parent layers
	resourceIDs := make([]RscTypeAndID, 0)
	for _, layer := range attribute.Layers {
		iamResourceType, err := convertResourceType(layer.Type, attribute.BusinessID)
		if err != nil {
			return nil, fmt.Errorf("convert resource type to iam resource type failed, layer: %+v, err: %+v", layer, err)
		}
		resourceID := RscTypeAndID{
			ResourceType: *iamResourceType,
			ResourceID:   strconv.FormatInt(layer.InstanceID, 10),
		}
		resourceIDs = append(resourceIDs, resourceID)
	}

	// append host resource id to end
	hostResourceID := RscTypeAndID{
		ResourceType: resourceType,
		ResourceID:   strconv.FormatInt(attribute.InstanceID, 10),
	}
	resourceIDs = append(resourceIDs, hostResourceID)
	return resourceIDs, nil
}

func eventSubscribeResourceID(resourceType ResourceTypeID, attribute *meta.ResourceAttribute) ([]RscTypeAndID, error) {
	if attribute.InstanceID <= 0 {
		return nil, nil
	}
	return []RscTypeAndID{
		{
			ResourceType: resourceType,
			ResourceID:   strconv.FormatInt(attribute.InstanceID, 10),
		},
	}, nil
}

func dynamicGroupingResourceID(resourceType ResourceTypeID, attribute *meta.ResourceAttribute) ([]RscTypeAndID, error) {
	if attribute.InstanceID <= 0 {
		return nil, nil
	}
	return []RscTypeAndID{
		{
			ResourceType: resourceType,
			ResourceID:   strconv.FormatInt(attribute.InstanceID, 10),
		},
	}, nil
}

func auditLogResourceID(resourceType ResourceTypeID, attribute *meta.ResourceAttribute) ([]RscTypeAndID, error) {
	return []RscTypeAndID{
		{
			ResourceType: resourceType,
			ResourceID:   strconv.FormatInt(attribute.InstanceID, 10),
		},
	}, nil
}
