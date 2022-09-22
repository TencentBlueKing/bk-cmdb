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

package iam

import (
	types2 "configcenter/cmd/scene_server/auth_server/sdk/types"
	meta2 "configcenter/pkg/ac/meta"
	"errors"
	"fmt"
	"strconv"
)

// GenIamResource TODO
func GenIamResource(act ActionID, rscType TypeID, a *meta2.ResourceAttribute) ([]types2.Resource, error) {
	// skip actions do not need to relate to resources
	if act == Skip {
		return genSkipResource(act, rscType, a)
	}

	switch a.Basic.Type {
	case meta2.Business:
		return genBusinessResource(act, rscType, a)
	case meta2.BizSet:
		return genBizSetResource(act, rscType, a)
	case meta2.DynamicGrouping:
		return genDynamicGroupingResource(act, rscType, a)
	case meta2.EventWatch:
		return genResourceWatch(act, rscType, a)
	case meta2.ProcessServiceTemplate, meta2.ProcessTemplate:
		return genServiceTemplateResource(act, rscType, a)
	case meta2.SetTemplate:
		return genSetTemplateResource(act, rscType, a)
	case meta2.OperationStatistic:
		return genOperationStatisticResource(act, rscType, a)
	case meta2.AuditLog:
		return genAuditLogResource(act, rscType, a)
	case meta2.CloudAreaInstance:
		return genPlat(act, rscType, a)
	case meta2.HostApply:
		return genHostApplyResource(act, rscType, a)
	case meta2.CloudAccount:
		return genCloudAccountResource(act, rscType, a)
	case meta2.CloudResourceTask:
		return genCloudResourceTaskResource(act, rscType, a)
	case meta2.ResourcePoolDirectory:
		return genResourcePoolDirectoryResource(act, rscType, a)
	case meta2.ProcessServiceInstance, meta2.Process:
		return genProcessServiceInstanceResource(act, rscType, a)
	case meta2.ModelModule, meta2.ModelSet, meta2.MainlineInstance, meta2.MainlineInstanceTopology:
		return genBusinessTopologyResource(act, rscType, a)
	case meta2.Model, meta2.ModelAssociation:
		return genModelResource(act, rscType, a)
	case meta2.ModelUnique:
		return genModelRelatedResource(act, rscType, a)
	case meta2.ModelAttributeGroup:
		if a.BusinessID > 0 {
			return genBizModelAttributeResource(act, rscType, a)
		} else {
			return genModelRelatedResource(act, rscType, a)
		}
	case meta2.ModelClassification:
		return genModelClassificationResource(act, rscType, a)
	case meta2.AssociationType:
		return genAssociationTypeResource(act, rscType, a)
	case meta2.ModelAttribute:
		if a.BusinessID > 0 {
			return genBizModelAttributeResource(act, rscType, a)
		} else {
			return genModelAttributeResource(act, rscType, a)
		}
	case meta2.ModelInstanceTopology, meta2.MainlineModelTopology, meta2.UserCustom:
		return genSkipResource(act, rscType, a)
	case meta2.ConfigAdmin:
		return genGlobalConfigResource(act, rscType, a)
	case meta2.MainlineModel:
		return genBusinessLayerResource(act, rscType, a)
	case meta2.ModelTopology:
		return genModelTopologyViewResource(act, rscType, a)
	case meta2.HostInstance:
		return genHostInstanceResource(act, rscType, a)
	case meta2.SystemBase:
		return make([]types2.Resource, 0), nil
	case meta2.ProcessServiceCategory:
		return genProcessServiceCategoryResource(act, rscType, a)
	case meta2.KubeCluster, meta2.KubeNode, meta2.KubeNamespace, meta2.KubeWorkload, meta2.KubeDeployment,
		meta2.KubeStatefulSet, meta2.KubeDaemonSet, meta2.KubeGameStatefulSet, meta2.KubeGameDeployment, meta2.KubeCronJob,
		meta2.KubeJob, meta2.KubePodWorkload, meta2.KubePod, meta2.KubeContainer:
		return make([]types2.Resource, 0), nil
	default:
		if IsCMDBSysInstance(a.Basic.Type) {
			return genSysInstanceResource(act, rscType, a)
		}
	}

	return nil, fmt.Errorf("gen id failed: unsupported resource type: %s", a.Type)
}

// genBusinessResource TODO
// generate business related resource id.
func genBusinessResource(act ActionID, typ TypeID, attribute *meta2.ResourceAttribute) ([]types2.Resource, error) {
	r := types2.Resource{
		System:    SystemIDCMDB,
		Type:      types2.ResourceType(typ),
		Attribute: nil,
	}

	// create business do not related to instance authorize
	if act == CreateBusiness {
		return make([]types2.Resource, 0), nil
	}

	// compatible for authorize any
	if attribute.InstanceID > 0 {
		r.ID = strconv.FormatInt(attribute.InstanceID, 10)
	}

	return []types2.Resource{r}, nil
}

// genBizSetResource TODO
// generate biz set related resource id.
func genBizSetResource(act ActionID, typ TypeID, attribute *meta2.ResourceAttribute) ([]types2.Resource, error) {
	r := types2.Resource{
		System:    SystemIDCMDB,
		Type:      types2.ResourceType(typ),
		Attribute: nil,
	}

	// create biz set do not related to instance authorize
	if act == CreateBizSet {
		return make([]types2.Resource, 0), nil
	}

	// compatible for authorize any
	if attribute.InstanceID > 0 {
		r.ID = strconv.FormatInt(attribute.InstanceID, 10)
	}

	return []types2.Resource{r}, nil
}

func genDynamicGroupingResource(act ActionID, typ TypeID, att *meta2.ResourceAttribute) ([]types2.Resource, error) {

	r := types2.Resource{
		System:    SystemIDCMDB,
		Attribute: nil,
	}

	if att.BusinessID <= 0 {
		return nil, errors.New("biz id can not be 0")
	}

	// do not related to instance authorize
	if act == CreateBusinessCustomQuery || act == ViewBusinessResource {
		r.Type = types2.ResourceType(Business)
		r.ID = strconv.FormatInt(att.BusinessID, 10)
		return []types2.Resource{r}, nil
	}

	r.Type = types2.ResourceType(typ)
	if len(att.InstanceIDEx) > 0 {
		r.ID = att.InstanceIDEx
	}

	// authorize based on business
	r.Attribute = map[string]interface{}{
		types2.IamPathKey: []string{fmt.Sprintf("/%s,%d/", Business, att.BusinessID)},
	}

	return []types2.Resource{r}, nil
}

func genProcessServiceCategoryResource(_ ActionID, _ TypeID, att *meta2.ResourceAttribute) ([]types2.Resource, error) {

	r := types2.Resource{
		System:    SystemIDCMDB,
		Attribute: nil,
	}

	// do not related to instance authorize
	r.Type = types2.ResourceType(Business)
	if att.BusinessID > 0 {
		r.ID = strconv.FormatInt(att.BusinessID, 10)
	}

	return []types2.Resource{r}, nil
}

func genResourceWatch(act ActionID, typ TypeID, att *meta2.ResourceAttribute) ([]types2.Resource, error) {
	switch act {
	case WatchCommonInstanceEvent:
		r := types2.Resource{
			System:    SystemIDCMDB,
			Attribute: nil,
		}

		// do not related to instance authorize
		r.Type = types2.ResourceType(SysModelEvent)
		if att.InstanceID > 0 {
			r.ID = strconv.FormatInt(att.InstanceID, 10)
		}
		return []types2.Resource{r}, nil

	case WatchMainlineInstanceEvent:
		r := types2.Resource{
			System:    SystemIDCMDB,
			Attribute: nil,
		}

		// do not related to instance authorize
		r.Type = types2.ResourceType(MainlineModelEvent)
		if att.InstanceID > 0 {
			r.ID = strconv.FormatInt(att.InstanceID, 10)
		}
		return []types2.Resource{r}, nil

	case WatchInstAsstEvent:
		r := types2.Resource{
			System:    SystemIDCMDB,
			Attribute: nil,
		}

		// do not related to instance authorize
		r.Type = types2.ResourceType(InstAsstEvent)
		if att.InstanceID > 0 {
			r.ID = strconv.FormatInt(att.InstanceID, 10)
		}
		return []types2.Resource{r}, nil

	case WatchKubeWorkloadEvent:
		r := types2.Resource{
			System: SystemIDCMDB,
		}

		r.Type = types2.ResourceType(KubeWorkloadEvent)
		if att.InstanceIDEx != "" {
			r.ID = att.InstanceIDEx
		}
		return []types2.Resource{r}, nil

	default:
		return make([]types2.Resource, 0), nil
	}
}

func genServiceTemplateResource(act ActionID, typ TypeID, att *meta2.ResourceAttribute) ([]types2.Resource, error) {

	r := types2.Resource{
		System:    SystemIDCMDB,
		Attribute: nil,
	}

	if act == CreateBusinessServiceTemplate {
		// do not related to instance authorize
		if att.BusinessID <= 0 {
			return nil, errors.New("biz id can not be 0")
		}
		r.Type = types2.ResourceType(Business)
		r.ID = strconv.FormatInt(att.BusinessID, 10)
		return []types2.Resource{r}, nil
	}

	r.Type = types2.ResourceType(typ)
	if att.InstanceID > 0 {
		r.ID = strconv.FormatInt(att.InstanceID, 10)
	}

	return []types2.Resource{r}, nil
}

func genSetTemplateResource(act ActionID, typ TypeID, att *meta2.ResourceAttribute) ([]types2.Resource, error) {
	r := types2.Resource{
		System:    SystemIDCMDB,
		Attribute: nil,
	}

	if act == CreateBusinessSetTemplate {
		// do not related to instance authorize
		if att.BusinessID <= 0 {
			return nil, errors.New("biz id can not be 0")
		}
		r.Type = types2.ResourceType(Business)
		r.ID = strconv.FormatInt(att.BusinessID, 10)
		return []types2.Resource{r}, nil
	}

	r.Type = types2.ResourceType(typ)
	if att.InstanceID > 0 {
		r.ID = strconv.FormatInt(att.InstanceID, 10)
	}

	return []types2.Resource{r}, nil
}

func genOperationStatisticResource(_ ActionID, typ TypeID, _ *meta2.ResourceAttribute) ([]types2.Resource, error) {
	return make([]types2.Resource, 0), nil
}

func genAuditLogResource(_ ActionID, typ TypeID, _ *meta2.ResourceAttribute) ([]types2.Resource, error) {
	return make([]types2.Resource, 0), nil
}

func genPlat(act ActionID, typ TypeID, att *meta2.ResourceAttribute) ([]types2.Resource, error) {
	r := types2.Resource{
		System:    SystemIDCMDB,
		Type:      types2.ResourceType(typ),
		Attribute: nil,
	}

	if act == CreateCloudArea {
		return make([]types2.Resource, 0), nil
	}

	if att.InstanceID > 0 {
		r.ID = strconv.FormatInt(att.InstanceID, 10)
	}
	return []types2.Resource{r}, nil

}

func genHostApplyResource(_ ActionID, _ TypeID, att *meta2.ResourceAttribute) ([]types2.Resource, error) {

	r := types2.Resource{
		System:    SystemIDCMDB,
		Attribute: nil,
	}

	r.Type = types2.ResourceType(Business)
	if att.BusinessID > 0 {
		r.ID = strconv.FormatInt(att.BusinessID, 10)
	}

	return []types2.Resource{r}, nil
}

func genCloudAccountResource(act ActionID, typ TypeID, att *meta2.ResourceAttribute) ([]types2.Resource, error) {
	r := types2.Resource{
		System:    SystemIDCMDB,
		Type:      types2.ResourceType(typ),
		Attribute: nil,
	}

	if act == CreateCloudAccount {
		return make([]types2.Resource, 0), nil
	}

	r.ID = strconv.FormatInt(att.InstanceID, 10)
	return []types2.Resource{r}, nil
}

func genCloudResourceTaskResource(act ActionID, typ TypeID, att *meta2.ResourceAttribute) ([]types2.Resource, error) {
	r := types2.Resource{
		System:    SystemIDCMDB,
		Type:      types2.ResourceType(typ),
		Attribute: nil,
	}

	if act == CreateCloudResourceTask {
		return make([]types2.Resource, 0), nil
	}

	if att.InstanceID > 0 {
		r.ID = strconv.FormatInt(att.InstanceID, 10)
	}
	return []types2.Resource{r}, nil
}

func genResourcePoolDirectoryResource(act ActionID, typ TypeID, att *meta2.ResourceAttribute) ([]types2.Resource, error) {
	r := types2.Resource{
		System:    SystemIDCMDB,
		Type:      types2.ResourceType(typ),
		Attribute: nil,
	}

	if act == CreateResourcePoolDirectory {
		return make([]types2.Resource, 0), nil
	}

	if att.InstanceID > 0 {
		r.ID = strconv.FormatInt(att.InstanceID, 10)
	}
	return []types2.Resource{r}, nil
}

func genProcessServiceInstanceResource(_ ActionID, typ TypeID, att *meta2.ResourceAttribute) ([]types2.Resource, error) {
	r := types2.Resource{
		System:    SystemIDCMDB,
		Type:      types2.ResourceType(typ),
		Attribute: nil,
	}

	// do not related to exact service instance authorize
	r.Type = types2.ResourceType(Business)
	if att.BusinessID > 0 {
		r.ID = strconv.FormatInt(att.BusinessID, 10)
	}

	return []types2.Resource{r}, nil
}

func genBusinessTopologyResource(_ ActionID, typ TypeID, att *meta2.ResourceAttribute) ([]types2.Resource, error) {
	r := types2.Resource{
		System:    SystemIDCMDB,
		Type:      types2.ResourceType(typ),
		Attribute: nil,
	}

	// do not related to exact instance authorize
	r.Type = types2.ResourceType(Business)
	if att.BusinessID > 0 {
		r.ID = strconv.FormatInt(att.BusinessID, 10)
	}

	return []types2.Resource{r}, nil
}

func genModelResource(act ActionID, typ TypeID, att *meta2.ResourceAttribute) ([]types2.Resource, error) {
	r := types2.Resource{
		System:    SystemIDCMDB,
		Type:      types2.ResourceType(typ),
		Attribute: nil,
	}

	// do not related to instance authorize
	if act == CreateSysModel {
		// create model authorized based on it's model group
		if len(att.Layers) > 0 {
			r.Type = types2.ResourceType(SysModelGroup)
			r.ID = strconv.FormatInt(att.Layers[0].InstanceID, 10)
			return []types2.Resource{r}, nil
		}
		return []types2.Resource{r}, nil
	}

	if att.InstanceID > 0 {
		r.ID = strconv.FormatInt(att.InstanceID, 10)
	}

	return []types2.Resource{r}, nil
}

func genModelRelatedResource(_ ActionID, typ TypeID, att *meta2.ResourceAttribute) ([]types2.Resource, error) {
	r := types2.Resource{
		System:    SystemIDCMDB,
		Type:      types2.ResourceType(typ),
		Attribute: nil,
	}

	if len(att.Layers) == 0 {
		return nil, NotEnoughLayer
	}

	r.Type = types2.ResourceType(SysModel)
	r.ID = strconv.FormatInt(att.Layers[0].InstanceID, 10)
	return []types2.Resource{r}, nil

}

func genModelClassificationResource(act ActionID, typ TypeID, att *meta2.ResourceAttribute) ([]types2.Resource, error) {
	r := types2.Resource{
		System:    SystemIDCMDB,
		Type:      types2.ResourceType(typ),
		Attribute: nil,
	}

	// create model group do not related to instance authorize
	if act == CreateModelGroup {
		return make([]types2.Resource, 0), nil
	}

	if att.InstanceID > 0 {
		r.ID = strconv.FormatInt(att.InstanceID, 10)
	}

	return []types2.Resource{r}, nil
}

func genAssociationTypeResource(act ActionID, typ TypeID, att *meta2.ResourceAttribute) ([]types2.Resource, error) {
	r := types2.Resource{
		System:    SystemIDCMDB,
		Type:      types2.ResourceType(typ),
		Attribute: nil,
	}

	if act == CreateAssociationType {
		return make([]types2.Resource, 0), nil
	}

	if att.InstanceID > 0 {
		r.ID = strconv.FormatInt(att.InstanceID, 10)
	}

	return []types2.Resource{r}, nil
}

func genModelAttributeResource(_ ActionID, _ TypeID, att *meta2.ResourceAttribute) ([]types2.Resource, error) {
	r := types2.Resource{
		System:    SystemIDCMDB,
		Type:      types2.ResourceType(SysModel),
		Attribute: nil,
	}

	if len(att.Layers) == 0 {
		return nil, NotEnoughLayer
	}

	r.ID = strconv.FormatInt(att.Layers[0].InstanceID, 10)

	return []types2.Resource{r}, nil
}

func genSkipResource(_ ActionID, _ TypeID, _ *meta2.ResourceAttribute) ([]types2.Resource, error) {
	return make([]types2.Resource, 0), nil
}

func genGlobalConfigResource(_ ActionID, _ TypeID, _ *meta2.ResourceAttribute) ([]types2.Resource, error) {
	return make([]types2.Resource, 0), nil
}

func genBusinessLayerResource(_ ActionID, typ TypeID, _ *meta2.ResourceAttribute) ([]types2.Resource, error) {
	return make([]types2.Resource, 0), nil
}

func genModelTopologyViewResource(_ ActionID, typ TypeID, _ *meta2.ResourceAttribute) ([]types2.Resource, error) {
	return make([]types2.Resource, 0), nil
}

func genHostInstanceResource(act ActionID, typ TypeID, a *meta2.ResourceAttribute) ([]types2.Resource, error) {

	// find host instances
	if act == Skip {
		r := types2.Resource{
			System:    SystemIDCMDB,
			Type:      types2.ResourceType(typ),
			Attribute: nil,
		}
		return []types2.Resource{r}, nil
	}

	// transfer resource pool's host to it's another directory.
	if act == ResourcePoolHostTransferToDirectory {
		if len(a.Layers) != 2 {
			return nil, NotEnoughLayer
		}

		resources := make([]types2.Resource, 2)
		resources[0] = types2.Resource{
			System: SystemIDCMDB,
			Type:   types2.ResourceType(SysHostRscPoolDirectory),
			ID:     strconv.FormatInt(a.Layers[0].InstanceID, 10),
		}

		resources[1] = types2.Resource{
			System: SystemIDCMDB,
			Type:   types2.ResourceType(SysResourcePoolDirectory),
			ID:     strconv.FormatInt(a.Layers[1].InstanceID, 10),
		}

		return resources, nil
	}

	// transfer host in resource pool to business
	if act == ResourcePoolHostTransferToBusiness {
		if len(a.Layers) != 2 {
			return nil, NotEnoughLayer
		}

		resources := make([]types2.Resource, 2)
		resources[0] = types2.Resource{
			System: SystemIDCMDB,
			Type:   types2.ResourceType(SysHostRscPoolDirectory),
			ID:     strconv.FormatInt(a.Layers[0].InstanceID, 10),
		}

		resources[1] = types2.Resource{
			System: SystemIDCMDB,
			Type:   types2.ResourceType(Business),
			ID:     strconv.FormatInt(a.Layers[1].InstanceID, 10),
		}

		return resources, nil
	}

	// transfer host from business to resource pool
	if act == BusinessHostTransferToResourcePool {
		if len(a.Layers) != 2 {
			return nil, NotEnoughLayer
		}

		resources := make([]types2.Resource, 2)
		resources[0] = types2.Resource{
			System: SystemIDCMDB,
			Type:   types2.ResourceType(Business),
			ID:     strconv.FormatInt(a.Layers[0].InstanceID, 10),
		}

		resources[1] = types2.Resource{
			System: SystemIDCMDB,
			Type:   types2.ResourceType(SysResourcePoolDirectory),
			ID:     strconv.FormatInt(a.Layers[1].InstanceID, 10),
		}

		return resources, nil
	}

	// transfer host from one business to another
	if act == HostTransferAcrossBusiness {
		if len(a.Layers) != 2 {
			return nil, NotEnoughLayer
		}

		resources := make([]types2.Resource, 2)
		resources[0] = types2.Resource{
			System: SystemIDCMDB,
			Type:   types2.ResourceType(BusinessForHostTrans),
			ID:     strconv.FormatInt(a.Layers[0].InstanceID, 10),
		}

		resources[1] = types2.Resource{
			System: SystemIDCMDB,
			Type:   types2.ResourceType(Business),
			ID:     strconv.FormatInt(a.Layers[1].InstanceID, 10),
		}

		return resources, nil
	}

	// import host
	if act == CreateResourcePoolHost {
		r := types2.Resource{
			System: SystemIDCMDB,
			Type:   types2.ResourceType(SysResourcePoolDirectory),
		}
		if len(a.Layers) > 0 {
			r.ID = strconv.FormatInt(a.Layers[0].InstanceID, 10)
		}
		return []types2.Resource{r}, nil
	}

	// edit or delete resource pool host instances
	if act == EditResourcePoolHost || act == DeleteResourcePoolHost {
		r := types2.Resource{
			System: SystemIDCMDB,
			Type:   types2.ResourceType(typ),
		}
		if a.InstanceID > 0 {
			r.ID = strconv.FormatInt(a.InstanceID, 10)
		}
		if len(a.Layers) > 0 {
			r.Attribute = map[string]interface{}{
				types2.IamPathKey: []string{fmt.Sprintf("/%s,%d/", SysHostRscPoolDirectory, a.Layers[0].InstanceID)},
			}
		}
		return []types2.Resource{r}, nil
	}

	// edit business host
	if act == EditBusinessHost {
		r := types2.Resource{
			System: SystemIDCMDB,
			Type:   types2.ResourceType(typ),
		}
		if a.InstanceID > 0 {
			r.ID = strconv.FormatInt(a.InstanceID, 10)
		}
		if len(a.Layers) > 0 {
			r.Attribute = map[string]interface{}{
				types2.IamPathKey: []string{fmt.Sprintf("/%s,%d/", Business, a.Layers[0].InstanceID)},
			}
		}
		return []types2.Resource{r}, nil
	}

	return []types2.Resource{}, nil
}

func genBizModelAttributeResource(_ ActionID, _ TypeID, att *meta2.ResourceAttribute) ([]types2.Resource, error) {
	r := types2.Resource{
		System:    SystemIDCMDB,
		Type:      types2.ResourceType(Business),
		Attribute: nil,
	}

	if att.BusinessID > 0 {
		r.ID = strconv.FormatInt(att.BusinessID, 10)
	}

	return []types2.Resource{r}, nil
}

func genSysInstanceResource(act ActionID, typ TypeID, att *meta2.ResourceAttribute) ([]types2.Resource, error) {
	r := types2.Resource{
		System:    SystemIDCMDB,
		Type:      types2.ResourceType(typ),
		Attribute: nil,
	}

	// create action do not related to instance authorize
	if att.Action == meta2.Create || att.Action == meta2.CreateMany {
		return make([]types2.Resource, 0), nil
	}

	if att.InstanceID > 0 {
		r.ID = strconv.FormatInt(att.InstanceID, 10)
	}

	return []types2.Resource{r}, nil
}
