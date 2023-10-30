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
	"errors"
	"fmt"
	"strconv"

	"configcenter/src/ac/meta"
	"configcenter/src/scene_server/auth_server/sdk/types"
)

var genIamResFuncMap = map[meta.ResourceType]func(ActionID, TypeID, *meta.ResourceAttribute) ([]types.Resource, error){
	meta.Business:                 genBusinessResource,
	meta.BizSet:                   genBizSetResource,
	meta.Project:                  genProjectResource,
	meta.DynamicGrouping:          genDynamicGroupingResource,
	meta.EventWatch:               genResourceWatch,
	meta.ProcessServiceTemplate:   genServiceTemplateResource,
	meta.ProcessTemplate:          genServiceTemplateResource,
	meta.SetTemplate:              genSetTemplateResource,
	meta.OperationStatistic:       genOperationStatisticResource,
	meta.AuditLog:                 genAuditLogResource,
	meta.CloudAreaInstance:        genPlat,
	meta.HostApply:                genHostApplyResource,
	meta.CloudAccount:             genCloudAccountResource,
	meta.CloudResourceTask:        genCloudResourceTaskResource,
	meta.ResourcePoolDirectory:    genResourcePoolDirectoryResource,
	meta.ProcessServiceInstance:   genProcessServiceInstanceResource,
	meta.Process:                  genProcessServiceInstanceResource,
	meta.ModelModule:              genBusinessTopologyResource,
	meta.ModelSet:                 genBusinessTopologyResource,
	meta.MainlineInstance:         genBusinessTopologyResource,
	meta.MainlineInstanceTopology: genBusinessTopologyResource,
	meta.Model:                    genModelResource,
	meta.ModelAssociation:         genModelResource,
	meta.ModelUnique:              genModelRelatedResource,
	meta.ModelClassification:      genModelClassificationResource,
	meta.AssociationType:          genAssociationTypeResource,
	meta.ModelInstanceTopology:    genSkipResource,
	meta.MainlineModelTopology:    genSkipResource,
	meta.UserCustom:               genSkipResource,
	meta.ConfigAdmin:              genGlobalConfigResource,
	meta.MainlineModel:            genBusinessLayerResource,
	meta.ModelTopology:            genModelTopologyViewResource,
	meta.HostInstance:             genHostInstanceResource,
	meta.ProcessServiceCategory:   genProcessServiceCategoryResource,
	meta.FieldTemplate:            genFieldTemplateResource,
}

// GenIamResource TODO
func GenIamResource(act ActionID, rscType TypeID, a *meta.ResourceAttribute) ([]types.Resource, error) {
	// skip actions do not need to relate to resources
	if act == Skip {
		return genSkipResource(act, rscType, a)
	}

	switch a.Basic.Type {
	case meta.ModelAttributeGroup:
		if a.BusinessID > 0 {
			return genBizModelAttributeResource(act, rscType, a)
		} else {
			return genModelRelatedResource(act, rscType, a)
		}
	case meta.ModelAttribute:
		if a.BusinessID > 0 {
			return genBizModelAttributeResource(act, rscType, a)
		} else {
			return genModelAttributeResource(act, rscType, a)
		}
	case meta.SystemBase:
		return make([]types.Resource, 0), nil
	case meta.FulltextSearch:
		return make([]types.Resource, 0), nil
	case meta.KubeCluster, meta.KubeNode, meta.KubeNamespace, meta.KubeWorkload, meta.KubeDeployment,
		meta.KubeStatefulSet, meta.KubeDaemonSet, meta.KubeGameStatefulSet, meta.KubeGameDeployment, meta.KubeCronJob,
		meta.KubeJob, meta.KubePodWorkload, meta.KubePod, meta.KubeContainer:
		return genKubeResource(act, rscType, a)
	}

	genIamResourceFunc, exists := genIamResFuncMap[a.Basic.Type]
	if exists {
		return genIamResourceFunc(act, rscType, a)
	}

	if IsCMDBSysInstance(a.Basic.Type) {
		return genSysInstanceResource(act, rscType, a)
	}
	return nil, fmt.Errorf("gen id failed: unsupported resource type: %s", a.Type)
}

// genBusinessResource TODO
// generate business related resource id.
func genBusinessResource(act ActionID, typ TypeID, attribute *meta.ResourceAttribute) ([]types.Resource, error) {
	r := types.Resource{
		System:    SystemIDCMDB,
		Type:      types.ResourceType(typ),
		Attribute: nil,
	}

	// create business do not related to instance authorize
	if act == CreateBusiness {
		return make([]types.Resource, 0), nil
	}

	// compatible for authorize any
	if attribute.InstanceID > 0 {
		r.ID = strconv.FormatInt(attribute.InstanceID, 10)
	}

	return []types.Resource{r}, nil
}

// genBizSetResource TODO
// generate biz set related resource id.
func genBizSetResource(act ActionID, typ TypeID, attribute *meta.ResourceAttribute) ([]types.Resource, error) {
	r := types.Resource{
		System:    SystemIDCMDB,
		Type:      types.ResourceType(typ),
		Attribute: nil,
	}

	// create biz set do not related to instance authorize
	if act == CreateBizSet {
		return make([]types.Resource, 0), nil
	}

	// compatible for authorize any
	if attribute.InstanceID > 0 {
		r.ID = strconv.FormatInt(attribute.InstanceID, 10)
	}

	return []types.Resource{r}, nil
}

func genProjectResource(act ActionID, typ TypeID, attribute *meta.ResourceAttribute) ([]types.Resource, error) {
	r := types.Resource{
		System:    SystemIDCMDB,
		Type:      types.ResourceType(typ),
		Attribute: nil,
	}

	// create project do not related to instance authorize
	if act == CreateProject {
		return make([]types.Resource, 0), nil
	}

	// compatible for authorize any
	if attribute.InstanceID > 0 {
		r.ID = strconv.FormatInt(attribute.InstanceID, 10)
	}

	return []types.Resource{r}, nil
}

func genDynamicGroupingResource(act ActionID, typ TypeID, att *meta.ResourceAttribute) ([]types.Resource, error) {

	r := types.Resource{
		System:    SystemIDCMDB,
		Attribute: nil,
	}

	if att.BusinessID <= 0 {
		return nil, errors.New("biz id can not be 0")
	}

	// do not related to instance authorize
	if act == CreateBusinessCustomQuery || act == ViewBusinessResource {
		r.Type = types.ResourceType(Business)
		r.ID = strconv.FormatInt(att.BusinessID, 10)
		return []types.Resource{r}, nil
	}

	r.Type = types.ResourceType(typ)
	if len(att.InstanceIDEx) > 0 {
		r.ID = att.InstanceIDEx
	}

	// authorize based on business
	r.Attribute = map[string]interface{}{
		types.IamPathKey: []string{fmt.Sprintf("/%s,%d/", Business, att.BusinessID)},
	}

	return []types.Resource{r}, nil
}

func genProcessServiceCategoryResource(_ ActionID, _ TypeID, att *meta.ResourceAttribute) ([]types.Resource, error) {

	r := types.Resource{
		System:    SystemIDCMDB,
		Attribute: nil,
	}

	// do not related to instance authorize
	r.Type = types.ResourceType(Business)
	if att.BusinessID > 0 {
		r.ID = strconv.FormatInt(att.BusinessID, 10)
	}

	return []types.Resource{r}, nil
}

func genResourceWatch(act ActionID, typ TypeID, att *meta.ResourceAttribute) ([]types.Resource, error) {
	switch act {
	case WatchCommonInstanceEvent:
		r := types.Resource{
			System:    SystemIDCMDB,
			Attribute: nil,
		}

		// do not related to instance authorize
		r.Type = types.ResourceType(SysModelEvent)
		if att.InstanceID > 0 {
			r.ID = strconv.FormatInt(att.InstanceID, 10)
		}
		return []types.Resource{r}, nil

	case WatchMainlineInstanceEvent:
		r := types.Resource{
			System:    SystemIDCMDB,
			Attribute: nil,
		}

		// do not related to instance authorize
		r.Type = types.ResourceType(MainlineModelEvent)
		if att.InstanceID > 0 {
			r.ID = strconv.FormatInt(att.InstanceID, 10)
		}
		return []types.Resource{r}, nil

	case WatchInstAsstEvent:
		r := types.Resource{
			System:    SystemIDCMDB,
			Attribute: nil,
		}

		// do not related to instance authorize
		r.Type = types.ResourceType(InstAsstEvent)
		if att.InstanceID > 0 {
			r.ID = strconv.FormatInt(att.InstanceID, 10)
		}
		return []types.Resource{r}, nil

	case WatchKubeWorkloadEvent:
		r := types.Resource{
			System: SystemIDCMDB,
		}

		r.Type = types.ResourceType(KubeWorkloadEvent)
		if att.InstanceIDEx != "" {
			r.ID = att.InstanceIDEx
		}
		return []types.Resource{r}, nil

	default:
		return make([]types.Resource, 0), nil
	}
}

func genServiceTemplateResource(act ActionID, typ TypeID, att *meta.ResourceAttribute) ([]types.Resource, error) {

	r := types.Resource{
		System:    SystemIDCMDB,
		Attribute: nil,
	}

	if act == CreateBusinessServiceTemplate {
		// do not related to instance authorize
		if att.BusinessID <= 0 {
			return nil, errors.New("biz id can not be 0")
		}
		r.Type = types.ResourceType(Business)
		r.ID = strconv.FormatInt(att.BusinessID, 10)
		return []types.Resource{r}, nil
	}

	r.Type = types.ResourceType(typ)
	if att.InstanceID > 0 {
		r.ID = strconv.FormatInt(att.InstanceID, 10)
	}

	return []types.Resource{r}, nil
}

func genSetTemplateResource(act ActionID, typ TypeID, att *meta.ResourceAttribute) ([]types.Resource, error) {
	r := types.Resource{
		System:    SystemIDCMDB,
		Attribute: nil,
	}

	if act == CreateBusinessSetTemplate {
		// do not related to instance authorize
		if att.BusinessID <= 0 {
			return nil, errors.New("biz id can not be 0")
		}
		r.Type = types.ResourceType(Business)
		r.ID = strconv.FormatInt(att.BusinessID, 10)
		return []types.Resource{r}, nil
	}

	r.Type = types.ResourceType(typ)
	if att.InstanceID > 0 {
		r.ID = strconv.FormatInt(att.InstanceID, 10)
	}

	return []types.Resource{r}, nil
}

func genOperationStatisticResource(_ ActionID, typ TypeID, _ *meta.ResourceAttribute) ([]types.Resource, error) {
	return make([]types.Resource, 0), nil
}

func genAuditLogResource(_ ActionID, typ TypeID, _ *meta.ResourceAttribute) ([]types.Resource, error) {
	return make([]types.Resource, 0), nil
}

func genPlat(act ActionID, typ TypeID, att *meta.ResourceAttribute) ([]types.Resource, error) {
	r := types.Resource{
		System:    SystemIDCMDB,
		Type:      types.ResourceType(typ),
		Attribute: nil,
	}

	if act == CreateCloudArea || act == ViewCloudArea {
		return make([]types.Resource, 0), nil
	}

	if att.InstanceID > 0 {
		r.ID = strconv.FormatInt(att.InstanceID, 10)
	}
	return []types.Resource{r}, nil

}

func genHostApplyResource(_ ActionID, _ TypeID, att *meta.ResourceAttribute) ([]types.Resource, error) {

	r := types.Resource{
		System:    SystemIDCMDB,
		Attribute: nil,
	}

	r.Type = types.ResourceType(Business)
	if att.BusinessID > 0 {
		r.ID = strconv.FormatInt(att.BusinessID, 10)
	}

	return []types.Resource{r}, nil
}

func genCloudAccountResource(act ActionID, typ TypeID, att *meta.ResourceAttribute) ([]types.Resource, error) {
	r := types.Resource{
		System:    SystemIDCMDB,
		Type:      types.ResourceType(typ),
		Attribute: nil,
	}

	if act == CreateCloudAccount {
		return make([]types.Resource, 0), nil
	}

	r.ID = strconv.FormatInt(att.InstanceID, 10)
	return []types.Resource{r}, nil
}

func genCloudResourceTaskResource(act ActionID, typ TypeID, att *meta.ResourceAttribute) ([]types.Resource, error) {
	r := types.Resource{
		System:    SystemIDCMDB,
		Type:      types.ResourceType(typ),
		Attribute: nil,
	}

	if act == CreateCloudResourceTask {
		return make([]types.Resource, 0), nil
	}

	if att.InstanceID > 0 {
		r.ID = strconv.FormatInt(att.InstanceID, 10)
	}
	return []types.Resource{r}, nil
}

func genResourcePoolDirectoryResource(act ActionID, typ TypeID, att *meta.ResourceAttribute) ([]types.Resource, error) {
	r := types.Resource{
		System:    SystemIDCMDB,
		Type:      types.ResourceType(typ),
		Attribute: nil,
	}

	if act == CreateResourcePoolDirectory {
		return make([]types.Resource, 0), nil
	}

	if att.InstanceID > 0 {
		r.ID = strconv.FormatInt(att.InstanceID, 10)
	}
	return []types.Resource{r}, nil
}

func genProcessServiceInstanceResource(_ ActionID, typ TypeID, att *meta.ResourceAttribute) ([]types.Resource, error) {
	r := types.Resource{
		System:    SystemIDCMDB,
		Type:      types.ResourceType(typ),
		Attribute: nil,
	}

	// do not related to exact service instance authorize
	r.Type = types.ResourceType(Business)
	if att.BusinessID > 0 {
		r.ID = strconv.FormatInt(att.BusinessID, 10)
	}

	return []types.Resource{r}, nil
}

func genBusinessTopologyResource(_ ActionID, typ TypeID, att *meta.ResourceAttribute) ([]types.Resource, error) {
	r := types.Resource{
		System:    SystemIDCMDB,
		Type:      types.ResourceType(typ),
		Attribute: nil,
	}

	// do not related to exact instance authorize
	r.Type = types.ResourceType(Business)
	if att.BusinessID > 0 {
		r.ID = strconv.FormatInt(att.BusinessID, 10)
	}

	return []types.Resource{r}, nil
}

func genModelResource(act ActionID, typ TypeID, att *meta.ResourceAttribute) ([]types.Resource, error) {
	r := types.Resource{
		System:    SystemIDCMDB,
		Type:      types.ResourceType(typ),
		Attribute: nil,
	}

	// do not related to instance authorize
	if act == CreateSysModel {
		// create model authorized based on it's model group
		if len(att.Layers) > 0 {
			r.Type = types.ResourceType(SysModelGroup)
			r.ID = strconv.FormatInt(att.Layers[0].InstanceID, 10)
			return []types.Resource{r}, nil
		}
		return []types.Resource{r}, nil
	}

	if att.InstanceID > 0 {
		r.ID = strconv.FormatInt(att.InstanceID, 10)
	}

	return []types.Resource{r}, nil
}

func genModelRelatedResource(_ ActionID, typ TypeID, att *meta.ResourceAttribute) ([]types.Resource, error) {
	r := types.Resource{
		System:    SystemIDCMDB,
		Type:      types.ResourceType(typ),
		Attribute: nil,
	}

	if len(att.Layers) == 0 {
		return nil, NotEnoughLayer
	}

	r.Type = types.ResourceType(SysModel)
	r.ID = strconv.FormatInt(att.Layers[0].InstanceID, 10)
	return []types.Resource{r}, nil

}

func genModelClassificationResource(act ActionID, typ TypeID, att *meta.ResourceAttribute) ([]types.Resource, error) {
	r := types.Resource{
		System:    SystemIDCMDB,
		Type:      types.ResourceType(typ),
		Attribute: nil,
	}

	// create model group do not related to instance authorize
	if act == CreateModelGroup {
		return make([]types.Resource, 0), nil
	}

	if att.InstanceID > 0 {
		r.ID = strconv.FormatInt(att.InstanceID, 10)
	}

	return []types.Resource{r}, nil
}

func genAssociationTypeResource(act ActionID, typ TypeID, att *meta.ResourceAttribute) ([]types.Resource, error) {
	r := types.Resource{
		System:    SystemIDCMDB,
		Type:      types.ResourceType(typ),
		Attribute: nil,
	}

	if act == CreateAssociationType {
		return make([]types.Resource, 0), nil
	}

	if att.InstanceID > 0 {
		r.ID = strconv.FormatInt(att.InstanceID, 10)
	}

	return []types.Resource{r}, nil
}

func genModelAttributeResource(_ ActionID, _ TypeID, att *meta.ResourceAttribute) ([]types.Resource, error) {
	r := types.Resource{
		System:    SystemIDCMDB,
		Type:      types.ResourceType(SysModel),
		Attribute: nil,
	}

	if len(att.Layers) == 0 {
		return nil, NotEnoughLayer
	}

	r.ID = strconv.FormatInt(att.Layers[0].InstanceID, 10)

	return []types.Resource{r}, nil
}

func genSkipResource(_ ActionID, _ TypeID, _ *meta.ResourceAttribute) ([]types.Resource, error) {
	return make([]types.Resource, 0), nil
}

func genGlobalConfigResource(_ ActionID, _ TypeID, _ *meta.ResourceAttribute) ([]types.Resource, error) {
	return make([]types.Resource, 0), nil
}

func genBusinessLayerResource(_ ActionID, typ TypeID, _ *meta.ResourceAttribute) ([]types.Resource, error) {
	return make([]types.Resource, 0), nil
}

func genModelTopologyViewResource(_ ActionID, typ TypeID, _ *meta.ResourceAttribute) ([]types.Resource, error) {
	return make([]types.Resource, 0), nil
}

func genHostInstanceResource(act ActionID, typ TypeID, a *meta.ResourceAttribute) ([]types.Resource, error) {

	// find host instances
	if act == Skip {
		r := types.Resource{
			System:    SystemIDCMDB,
			Type:      types.ResourceType(typ),
			Attribute: nil,
		}
		return []types.Resource{r}, nil
	}

	// transfer resource pool's host to it's another directory.
	if act == ResourcePoolHostTransferToDirectory {
		if len(a.Layers) != 2 {
			return nil, NotEnoughLayer
		}

		resources := make([]types.Resource, 2)
		resources[0] = types.Resource{
			System: SystemIDCMDB,
			Type:   types.ResourceType(SysHostRscPoolDirectory),
			ID:     strconv.FormatInt(a.Layers[0].InstanceID, 10),
		}

		resources[1] = types.Resource{
			System: SystemIDCMDB,
			Type:   types.ResourceType(SysResourcePoolDirectory),
			ID:     strconv.FormatInt(a.Layers[1].InstanceID, 10),
		}

		return resources, nil
	}

	// transfer host in resource pool to business
	if act == ResourcePoolHostTransferToBusiness {
		if len(a.Layers) != 2 {
			return nil, NotEnoughLayer
		}

		resources := make([]types.Resource, 2)
		resources[0] = types.Resource{
			System: SystemIDCMDB,
			Type:   types.ResourceType(SysHostRscPoolDirectory),
			ID:     strconv.FormatInt(a.Layers[0].InstanceID, 10),
		}

		resources[1] = types.Resource{
			System: SystemIDCMDB,
			Type:   types.ResourceType(Business),
			ID:     strconv.FormatInt(a.Layers[1].InstanceID, 10),
		}

		return resources, nil
	}

	// transfer host from business to resource pool
	if act == BusinessHostTransferToResourcePool {
		if len(a.Layers) != 2 {
			return nil, NotEnoughLayer
		}

		resources := make([]types.Resource, 2)
		resources[0] = types.Resource{
			System: SystemIDCMDB,
			Type:   types.ResourceType(Business),
			ID:     strconv.FormatInt(a.Layers[0].InstanceID, 10),
		}

		resources[1] = types.Resource{
			System: SystemIDCMDB,
			Type:   types.ResourceType(SysResourcePoolDirectory),
			ID:     strconv.FormatInt(a.Layers[1].InstanceID, 10),
		}

		return resources, nil
	}

	// transfer host from one business to another
	if act == HostTransferAcrossBusiness {
		if len(a.Layers) != 2 {
			return nil, NotEnoughLayer
		}

		resources := make([]types.Resource, 2)
		resources[0] = types.Resource{
			System: SystemIDCMDB,
			Type:   types.ResourceType(BusinessForHostTrans),
			ID:     strconv.FormatInt(a.Layers[0].InstanceID, 10),
		}

		resources[1] = types.Resource{
			System: SystemIDCMDB,
			Type:   types.ResourceType(Business),
			ID:     strconv.FormatInt(a.Layers[1].InstanceID, 10),
		}

		return resources, nil
	}

	// import host
	if act == CreateResourcePoolHost {
		r := types.Resource{
			System: SystemIDCMDB,
			Type:   types.ResourceType(SysResourcePoolDirectory),
		}
		if len(a.Layers) > 0 {
			r.ID = strconv.FormatInt(a.Layers[0].InstanceID, 10)
		}
		return []types.Resource{r}, nil
	}

	// edit or delete resource pool host instances
	if act == EditResourcePoolHost || act == DeleteResourcePoolHost {
		r := types.Resource{
			System: SystemIDCMDB,
			Type:   types.ResourceType(typ),
		}
		if a.InstanceID > 0 {
			r.ID = strconv.FormatInt(a.InstanceID, 10)
		}
		if len(a.Layers) > 0 {
			r.Attribute = map[string]interface{}{
				types.IamPathKey: []string{fmt.Sprintf("/%s,%d/", SysHostRscPoolDirectory, a.Layers[0].InstanceID)},
			}
		}
		return []types.Resource{r}, nil
	}

	// edit business host
	if act == EditBusinessHost {
		r := types.Resource{
			System: SystemIDCMDB,
			Type:   types.ResourceType(typ),
		}
		if a.InstanceID > 0 {
			r.ID = strconv.FormatInt(a.InstanceID, 10)
		}
		if len(a.Layers) > 0 {
			r.Attribute = map[string]interface{}{
				types.IamPathKey: []string{fmt.Sprintf("/%s,%d/", Business, a.Layers[0].InstanceID)},
			}
		}
		return []types.Resource{r}, nil
	}

	// find business host
	if act == ViewBusinessResource {
		r := types.Resource{
			System: SystemIDCMDB,
			Type:   types.ResourceType(Business),
			ID:     strconv.FormatInt(a.BusinessID, 10),
		}
		return []types.Resource{r}, nil
	}

	return []types.Resource{}, nil
}

func genBizModelAttributeResource(_ ActionID, _ TypeID, att *meta.ResourceAttribute) ([]types.Resource, error) {
	r := types.Resource{
		System:    SystemIDCMDB,
		Type:      types.ResourceType(Business),
		Attribute: nil,
	}

	if att.BusinessID > 0 {
		r.ID = strconv.FormatInt(att.BusinessID, 10)
	}

	return []types.Resource{r}, nil
}

func genSysInstanceResource(act ActionID, typ TypeID, att *meta.ResourceAttribute) ([]types.Resource, error) {
	r := types.Resource{
		System:    SystemIDCMDB,
		Type:      types.ResourceType(typ),
		Attribute: nil,
	}

	// create action do not related to instance authorize
	if att.Action == meta.Create || att.Action == meta.CreateMany || att.Action == meta.FindMany ||
		att.Action == meta.Find {
		return make([]types.Resource, 0), nil
	}

	if att.InstanceID > 0 {
		r.ID = strconv.FormatInt(att.InstanceID, 10)
	}

	return []types.Resource{r}, nil
}

func genFieldTemplateResource(act ActionID, typ TypeID, att *meta.ResourceAttribute) ([]types.Resource, error) {
	r := types.Resource{
		System: SystemIDCMDB,
		Type:   types.ResourceType(FieldGroupingTemplate),
	}

	// create action do not related to instance authorize
	if att.Action == meta.Create || att.Action == meta.CreateMany {
		return make([]types.Resource, 0), nil
	}

	if att.InstanceID > 0 {
		r.ID = strconv.FormatInt(att.InstanceID, 10)
	}

	return []types.Resource{r}, nil
}

func genKubeResource(act ActionID, typ TypeID, att *meta.ResourceAttribute) ([]types.Resource, error) {
	if act == ViewBusinessResource {
		if att.BusinessID <= 0 {
			return nil, errors.New("biz id can not be 0")
		}

		return []types.Resource{{
			System: SystemIDCMDB,
			Type:   types.ResourceType(Business),
			ID:     strconv.FormatInt(att.BusinessID, 10),
		}}, nil
	}

	return make([]types.Resource, 0), nil
}
