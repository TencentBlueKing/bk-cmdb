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

	iamtypes "configcenter/src/ac/iam/types"
	"configcenter/src/ac/meta"
	"configcenter/src/scene_server/auth_server/sdk/types"
	"configcenter/src/thirdparty/apigw/iam"
)

var genIamResFuncMap = map[meta.ResourceType]func(iamtypes.ActionID, iamtypes.TypeID,
	*meta.ResourceAttribute) ([]iam.Resource, error){
	meta.Business:                 genBusinessResource,
	meta.BizSet:                   genBizSetResource,
	meta.Project:                  genProjectResource,
	meta.DynamicGrouping:          genDynamicGroupingResource,
	meta.EventWatch:               genResourceWatch,
	meta.ProcessServiceTemplate:   genServiceTemplateResource,
	meta.ProcessTemplate:          genServiceTemplateResource,
	meta.SetTemplate:              genSetTemplateResource,
	meta.AuditLog:                 genAuditLogResource,
	meta.CloudAreaInstance:        genPlat,
	meta.HostApply:                genHostApplyResource,
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
	meta.FullSyncCond:             genSkipResource,
	meta.GeneralCache:             genGeneralCacheResource,
	meta.TenantSet:                genTenantSetResource,
}

// GenIamResource TODO
func GenIamResource(act iamtypes.ActionID, rscType iamtypes.TypeID, a *meta.ResourceAttribute) ([]iam.Resource, error) {
	// skip actions do not need to relate to resources
	if act == iamtypes.Skip {
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
		return make([]iam.Resource, 0), nil
	case meta.FulltextSearch:
		return make([]iam.Resource, 0), nil
	case meta.IDRuleIncrID:
		return make([]iam.Resource, 0), nil
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
func genBusinessResource(act iamtypes.ActionID, typ iamtypes.TypeID,
	attribute *meta.ResourceAttribute) ([]iam.Resource, error) {
	r := iam.Resource{
		System:    iamtypes.SystemIDCMDB,
		Type:      iam.IamResourceType(typ),
		Attribute: nil,
	}

	// create business do not related to instance authorize
	if act == iamtypes.CreateBusiness {
		return make([]iam.Resource, 0), nil
	}

	// compatible for authorize any
	if attribute.InstanceID > 0 {
		r.ID = strconv.FormatInt(attribute.InstanceID, 10)
	}

	return []iam.Resource{r}, nil
}

// genBizSetResource TODO
// generate biz set related resource id.
func genBizSetResource(act iamtypes.ActionID, typ iamtypes.TypeID, attribute *meta.ResourceAttribute) ([]iam.Resource,
	error) {

	r := iam.Resource{
		System:    iamtypes.SystemIDCMDB,
		Type:      iam.IamResourceType(typ),
		Attribute: nil,
	}

	// create biz set do not related to instance authorize
	if act == iamtypes.CreateBizSet {
		return make([]iam.Resource, 0), nil
	}

	// compatible for authorize any
	if attribute.InstanceID > 0 {
		r.ID = strconv.FormatInt(attribute.InstanceID, 10)
	}

	return []iam.Resource{r}, nil
}

func genProjectResource(act iamtypes.ActionID, typ iamtypes.TypeID, attribute *meta.ResourceAttribute) ([]iam.Resource,
	error) {

	r := iam.Resource{
		System:    iamtypes.SystemIDCMDB,
		Type:      iam.IamResourceType(typ),
		Attribute: nil,
	}

	// create project do not related to instance authorize
	if act == iamtypes.CreateProject {
		return make([]iam.Resource, 0), nil
	}

	// compatible for authorize any
	if attribute.InstanceID > 0 {
		r.ID = strconv.FormatInt(attribute.InstanceID, 10)
	}

	return []iam.Resource{r}, nil
}

func genDynamicGroupingResource(act iamtypes.ActionID, typ iamtypes.TypeID, att *meta.ResourceAttribute) (
	[]iam.Resource, error) {

	r := iam.Resource{
		System:    iamtypes.SystemIDCMDB,
		Attribute: nil,
	}

	if att.BusinessID <= 0 {
		return nil, errors.New("biz id can not be 0")
	}

	// do not related to instance authorize
	if act == iamtypes.CreateBusinessCustomQuery || act == iamtypes.ViewBusinessResource {
		r.Type = iam.IamResourceType(iamtypes.Business)
		r.ID = strconv.FormatInt(att.BusinessID, 10)
		return []iam.Resource{r}, nil
	}

	r.Type = iam.IamResourceType(typ)
	if len(att.InstanceIDEx) > 0 {
		r.ID = att.InstanceIDEx
	}

	// authorize based on business
	r.Attribute = map[string]interface{}{
		types.IamPathKey: []string{fmt.Sprintf("/%s,%d/", iamtypes.Business, att.BusinessID)},
	}

	return []iam.Resource{r}, nil
}

func genProcessServiceCategoryResource(_ iamtypes.ActionID, _ iamtypes.TypeID,
	att *meta.ResourceAttribute) ([]iam.Resource, error) {

	r := iam.Resource{
		System:    iamtypes.SystemIDCMDB,
		Attribute: nil,
	}

	// do not related to instance authorize
	r.Type = iam.IamResourceType(iamtypes.Business)
	if att.BusinessID > 0 {
		r.ID = strconv.FormatInt(att.BusinessID, 10)
	}

	return []iam.Resource{r}, nil
}

func genResourceWatch(act iamtypes.ActionID, typ iamtypes.TypeID, att *meta.ResourceAttribute) ([]iam.Resource, error) {

	switch act {
	case iamtypes.WatchCommonInstanceEvent:
		r := iam.Resource{
			System:    iamtypes.SystemIDCMDB,
			Attribute: nil,
		}

		// do not related to instance authorize
		r.Type = iam.IamResourceType(iamtypes.SysModelEvent)
		if att.InstanceID > 0 {
			r.ID = strconv.FormatInt(att.InstanceID, 10)
		}
		return []iam.Resource{r}, nil

	case iamtypes.WatchMainlineInstanceEvent:
		r := iam.Resource{
			System:    iamtypes.SystemIDCMDB,
			Attribute: nil,
		}

		// do not related to instance authorize
		r.Type = iam.IamResourceType(iamtypes.MainlineModelEvent)
		if att.InstanceID > 0 {
			r.ID = strconv.FormatInt(att.InstanceID, 10)
		}
		return []iam.Resource{r}, nil

	case iamtypes.WatchInstAsstEvent:
		r := iam.Resource{
			System:    iamtypes.SystemIDCMDB,
			Attribute: nil,
		}

		// do not related to instance authorize
		r.Type = iam.IamResourceType(iamtypes.InstAsstEvent)
		if att.InstanceID > 0 {
			r.ID = strconv.FormatInt(att.InstanceID, 10)
		}
		return []iam.Resource{r}, nil

	case iamtypes.WatchKubeWorkloadEvent:
		r := iam.Resource{
			System: iamtypes.SystemIDCMDB,
		}

		r.Type = iam.IamResourceType(iamtypes.KubeWorkloadEvent)
		if att.InstanceIDEx != "" {
			r.ID = att.InstanceIDEx
		}
		return []iam.Resource{r}, nil

	default:
		return make([]iam.Resource, 0), nil
	}
}

func genServiceTemplateResource(act iamtypes.ActionID, typ iamtypes.TypeID,
	att *meta.ResourceAttribute) ([]iam.Resource, error) {

	r := iam.Resource{
		System:    iamtypes.SystemIDCMDB,
		Attribute: nil,
	}

	if act == iamtypes.CreateBusinessServiceTemplate {
		// do not related to instance authorize
		if att.BusinessID <= 0 {
			return nil, errors.New("biz id can not be 0")
		}
		r.Type = iam.IamResourceType(iamtypes.Business)
		r.ID = strconv.FormatInt(att.BusinessID, 10)
		return []iam.Resource{r}, nil
	}

	r.Type = iam.IamResourceType(typ)
	if att.InstanceID > 0 {
		r.ID = strconv.FormatInt(att.InstanceID, 10)
	}

	return []iam.Resource{r}, nil
}

func genSetTemplateResource(act iamtypes.ActionID, typ iamtypes.TypeID, att *meta.ResourceAttribute) ([]iam.Resource,
	error) {

	r := iam.Resource{
		System:    iamtypes.SystemIDCMDB,
		Attribute: nil,
	}

	if act == iamtypes.CreateBusinessSetTemplate {
		// do not related to instance authorize
		if att.BusinessID <= 0 {
			return nil, errors.New("biz id can not be 0")
		}
		r.Type = iam.IamResourceType(iamtypes.Business)
		r.ID = strconv.FormatInt(att.BusinessID, 10)
		return []iam.Resource{r}, nil
	}

	r.Type = iam.IamResourceType(typ)
	if att.InstanceID > 0 {
		r.ID = strconv.FormatInt(att.InstanceID, 10)
	}

	return []iam.Resource{r}, nil
}

func genAuditLogResource(_ iamtypes.ActionID, typ iamtypes.TypeID, _ *meta.ResourceAttribute) ([]iam.Resource, error) {
	return make([]iam.Resource, 0), nil
}

func genPlat(act iamtypes.ActionID, typ iamtypes.TypeID, att *meta.ResourceAttribute) ([]iam.Resource, error) {
	r := iam.Resource{
		System:    iamtypes.SystemIDCMDB,
		Type:      iam.IamResourceType(typ),
		Attribute: nil,
	}

	if act == iamtypes.CreateCloudArea || act == iamtypes.ViewCloudArea {
		return make([]iam.Resource, 0), nil
	}

	if att.InstanceID > 0 {
		r.ID = strconv.FormatInt(att.InstanceID, 10)
	}
	return []iam.Resource{r}, nil

}

func genHostApplyResource(_ iamtypes.ActionID, _ iamtypes.TypeID, att *meta.ResourceAttribute) ([]iam.Resource, error) {

	r := iam.Resource{
		System:    iamtypes.SystemIDCMDB,
		Attribute: nil,
	}

	r.Type = iam.IamResourceType(iamtypes.Business)
	if att.BusinessID > 0 {
		r.ID = strconv.FormatInt(att.BusinessID, 10)
	}

	return []iam.Resource{r}, nil
}

func genResourcePoolDirectoryResource(act iamtypes.ActionID, typ iamtypes.TypeID,
	att *meta.ResourceAttribute) ([]iam.Resource, error) {
	r := iam.Resource{
		System:    iamtypes.SystemIDCMDB,
		Type:      iam.IamResourceType(typ),
		Attribute: nil,
	}

	if act == iamtypes.CreateResourcePoolDirectory {
		return make([]iam.Resource, 0), nil
	}

	if att.InstanceID > 0 {
		r.ID = strconv.FormatInt(att.InstanceID, 10)
	}
	return []iam.Resource{r}, nil
}

func genProcessServiceInstanceResource(_ iamtypes.ActionID, typ iamtypes.TypeID,
	att *meta.ResourceAttribute) ([]iam.Resource, error) {

	r := iam.Resource{
		System:    iamtypes.SystemIDCMDB,
		Type:      iam.IamResourceType(typ),
		Attribute: nil,
	}

	// do not related to exact service instance authorize
	r.Type = iam.IamResourceType(iamtypes.Business)
	if att.BusinessID > 0 {
		r.ID = strconv.FormatInt(att.BusinessID, 10)
	}

	return []iam.Resource{r}, nil
}

func genBusinessTopologyResource(_ iamtypes.ActionID, typ iamtypes.TypeID,
	att *meta.ResourceAttribute) ([]iam.Resource, error) {

	r := iam.Resource{
		System:    iamtypes.SystemIDCMDB,
		Type:      iam.IamResourceType(typ),
		Attribute: nil,
	}

	// do not related to exact instance authorize
	r.Type = iam.IamResourceType(iamtypes.Business)
	if att.BusinessID > 0 {
		r.ID = strconv.FormatInt(att.BusinessID, 10)
	}

	return []iam.Resource{r}, nil
}

func genModelResource(act iamtypes.ActionID, typ iamtypes.TypeID, att *meta.ResourceAttribute) ([]iam.Resource, error) {
	r := iam.Resource{
		System:    iamtypes.SystemIDCMDB,
		Type:      iam.IamResourceType(typ),
		Attribute: nil,
	}

	// do not related to instance authorize
	if act == iamtypes.CreateSysModel {
		r.Type = iam.IamResourceType(iamtypes.SysModelGroup)
		// create model authorized based on it's model group
		if len(att.Layers) > 0 {
			r.ID = strconv.FormatInt(att.Layers[0].InstanceID, 10)
			return []iam.Resource{r}, nil
		}
		return []iam.Resource{r}, nil
	}

	if att.InstanceID > 0 {
		r.ID = strconv.FormatInt(att.InstanceID, 10)
	}

	return []iam.Resource{r}, nil
}

func genModelRelatedResource(_ iamtypes.ActionID, typ iamtypes.TypeID, att *meta.ResourceAttribute) ([]iam.Resource,
	error) {

	r := iam.Resource{
		System:    iamtypes.SystemIDCMDB,
		Type:      iam.IamResourceType(typ),
		Attribute: nil,
	}

	if len(att.Layers) == 0 {
		return nil, NotEnoughLayer
	}

	r.Type = iam.IamResourceType(iamtypes.SysModel)
	r.ID = strconv.FormatInt(att.Layers[0].InstanceID, 10)
	return []iam.Resource{r}, nil

}

func genModelClassificationResource(act iamtypes.ActionID, typ iamtypes.TypeID,
	att *meta.ResourceAttribute) ([]iam.Resource, error) {

	r := iam.Resource{
		System:    iamtypes.SystemIDCMDB,
		Type:      iam.IamResourceType(typ),
		Attribute: nil,
	}

	// create model group do not related to instance authorize
	if act == iamtypes.CreateModelGroup {
		return make([]iam.Resource, 0), nil
	}

	if att.InstanceID > 0 {
		r.ID = strconv.FormatInt(att.InstanceID, 10)
	}

	return []iam.Resource{r}, nil
}

func genAssociationTypeResource(act iamtypes.ActionID, typ iamtypes.TypeID,
	att *meta.ResourceAttribute) ([]iam.Resource, error) {

	r := iam.Resource{
		System:    iamtypes.SystemIDCMDB,
		Type:      iam.IamResourceType(typ),
		Attribute: nil,
	}

	if act == iamtypes.CreateAssociationType {
		return make([]iam.Resource, 0), nil
	}

	if att.InstanceID > 0 {
		r.ID = strconv.FormatInt(att.InstanceID, 10)
	}

	return []iam.Resource{r}, nil
}

func genModelAttributeResource(_ iamtypes.ActionID, _ iamtypes.TypeID, att *meta.ResourceAttribute) ([]iam.Resource,
	error) {

	r := iam.Resource{
		System:    iamtypes.SystemIDCMDB,
		Type:      iam.IamResourceType(iamtypes.SysModel),
		Attribute: nil,
	}

	if len(att.Layers) == 0 {
		return nil, NotEnoughLayer
	}

	r.ID = strconv.FormatInt(att.Layers[0].InstanceID, 10)

	return []iam.Resource{r}, nil
}

func genSkipResource(_ iamtypes.ActionID, _ iamtypes.TypeID, _ *meta.ResourceAttribute) ([]iam.Resource, error) {
	return make([]iam.Resource, 0), nil
}

func genGlobalConfigResource(_ iamtypes.ActionID, _ iamtypes.TypeID, _ *meta.ResourceAttribute) ([]iam.Resource,
	error) {
	return make([]iam.Resource, 0), nil
}

func genBusinessLayerResource(_ iamtypes.ActionID, typ iamtypes.TypeID, _ *meta.ResourceAttribute) ([]iam.Resource,
	error) {
	return make([]iam.Resource, 0), nil
}

func genModelTopologyViewResource(_ iamtypes.ActionID, typ iamtypes.TypeID,
	_ *meta.ResourceAttribute) ([]iam.Resource,
	error) {
	return make([]iam.Resource, 0), nil
}

func getHostTransferResource(types []iam.IamResourceType, a *meta.ResourceAttribute) ([]iam.Resource, error) {

	if len(a.Layers) != 2 {
		return nil, NotEnoughLayer
	}
	resources := make([]iam.Resource, 2)
	resources[0] = iam.Resource{
		System: iamtypes.SystemIDCMDB,
		Type:   types[0],
		ID:     strconv.FormatInt(a.Layers[0].InstanceID, 10),
	}
	resources[1] = iam.Resource{
		System: iamtypes.SystemIDCMDB,
		Type:   types[1],
		ID:     strconv.FormatInt(a.Layers[1].InstanceID, 10),
	}

	return resources, nil
}

func genHostInstanceResource(act iamtypes.ActionID, typ iamtypes.TypeID, a *meta.ResourceAttribute) ([]iam.Resource,
	error) {

	// find host instances
	if act == iamtypes.Skip {
		r := iam.Resource{System: iamtypes.SystemIDCMDB, Type: iam.IamResourceType(typ), Attribute: nil}
		return []iam.Resource{r}, nil
	}

	// transfer resource pool's host to it's another directory.
	if act == iamtypes.ResourcePoolHostTransferToDirectory {
		types := []iam.IamResourceType{iam.IamResourceType(iamtypes.SysHostRscPoolDirectory),
			iam.IamResourceType(iamtypes.SysResourcePoolDirectory)}
		return getHostTransferResource(types, a)
	}

	// transfer host in resource pool to business
	if act == iamtypes.ResourcePoolHostTransferToBusiness {
		types := []iam.IamResourceType{iam.IamResourceType(iamtypes.SysHostRscPoolDirectory),
			iam.IamResourceType(iamtypes.Business)}
		return getHostTransferResource(types, a)
	}

	// transfer host from business to resource pool
	if act == iamtypes.BusinessHostTransferToResourcePool {
		types := []iam.IamResourceType{iam.IamResourceType(iamtypes.Business),
			iam.IamResourceType(iamtypes.SysResourcePoolDirectory)}
		return getHostTransferResource(types, a)
	}

	// transfer host from one business to another
	if act == iamtypes.HostTransferAcrossBusiness {
		types := []iam.IamResourceType{iam.IamResourceType(iamtypes.BusinessForHostTrans),
			iam.IamResourceType(iamtypes.Business)}
		return getHostTransferResource(types, a)
	}

	// import host
	if act == iamtypes.CreateResourcePoolHost {
		r := iam.Resource{System: iamtypes.SystemIDCMDB, Type: iam.IamResourceType(iamtypes.SysResourcePoolDirectory)}
		if len(a.Layers) > 0 {
			r.ID = strconv.FormatInt(a.Layers[0].InstanceID, 10)
		}
		return []iam.Resource{r}, nil
	}

	// edit or delete resource pool host instances
	if act == iamtypes.EditResourcePoolHost || act == iamtypes.DeleteResourcePoolHost {
		r := iam.Resource{System: iamtypes.SystemIDCMDB, Type: iam.IamResourceType(typ)}
		if a.InstanceID > 0 {
			r.ID = strconv.FormatInt(a.InstanceID, 10)
		}
		if len(a.Layers) > 0 {
			r.Attribute = map[string]interface{}{types.IamPathKey: []string{
				fmt.Sprintf("/%s,%d/", iamtypes.SysHostRscPoolDirectory, a.Layers[0].InstanceID)}}
		}
		return []iam.Resource{r}, nil
	}

	// edit business host
	if act == iamtypes.EditBusinessHost {
		r := iam.Resource{System: iamtypes.SystemIDCMDB, Type: iam.IamResourceType(typ)}
		if a.InstanceID > 0 {
			r.ID = strconv.FormatInt(a.InstanceID, 10)
		}
		if len(a.Layers) > 0 {
			r.Attribute = map[string]interface{}{types.IamPathKey: []string{fmt.Sprintf("/%s,%d/", iamtypes.Business,
				a.Layers[0].InstanceID)}}
		}
		return []iam.Resource{r}, nil
	}

	// find business host
	if act == iamtypes.ViewBusinessResource {
		r := iam.Resource{System: iamtypes.SystemIDCMDB, Type: iam.IamResourceType(iamtypes.Business),
			ID: strconv.FormatInt(a.BusinessID, 10)}
		return []iam.Resource{r}, nil
	}
	return []iam.Resource{}, nil
}

func genBizModelAttributeResource(_ iamtypes.ActionID, _ iamtypes.TypeID,
	att *meta.ResourceAttribute) ([]iam.Resource, error) {

	r := iam.Resource{
		System:    iamtypes.SystemIDCMDB,
		Type:      iam.IamResourceType(iamtypes.Business),
		Attribute: nil,
	}

	if att.BusinessID > 0 {
		r.ID = strconv.FormatInt(att.BusinessID, 10)
	}

	return []iam.Resource{r}, nil
}

func genSysInstanceResource(act iamtypes.ActionID, typ iamtypes.TypeID, att *meta.ResourceAttribute) ([]iam.Resource,
	error) {
	r := iam.Resource{
		System:    iamtypes.SystemIDCMDB,
		Type:      iam.IamResourceType(typ),
		Attribute: nil,
	}

	// create action do not related to instance authorize
	if att.Action == meta.Create || att.Action == meta.CreateMany || att.Action == meta.FindMany ||
		att.Action == meta.Find {
		return make([]iam.Resource, 0), nil
	}

	if att.InstanceID > 0 {
		r.ID = strconv.FormatInt(att.InstanceID, 10)
	}

	return []iam.Resource{r}, nil
}

func genFieldTemplateResource(act iamtypes.ActionID, typ iamtypes.TypeID,
	att *meta.ResourceAttribute) ([]iam.Resource, error) {

	r := iam.Resource{
		System: iamtypes.SystemIDCMDB,
		Type:   iam.IamResourceType(iamtypes.FieldGroupingTemplate),
	}

	// create action do not related to instance authorize
	if att.Action == meta.Create || att.Action == meta.CreateMany {
		return make([]iam.Resource, 0), nil
	}

	if att.InstanceID > 0 {
		r.ID = strconv.FormatInt(att.InstanceID, 10)
	}

	return []iam.Resource{r}, nil
}

func genKubeResource(act iamtypes.ActionID, typ iamtypes.TypeID, att *meta.ResourceAttribute) ([]iam.Resource,
	error) {
	if act == iamtypes.ViewBusinessResource {
		if att.BusinessID <= 0 {
			return nil, errors.New("biz id can not be 0")
		}

		return []iam.Resource{{
			System: iamtypes.SystemIDCMDB,
			Type:   iam.IamResourceType(iamtypes.Business),
			ID:     strconv.FormatInt(att.BusinessID, 10),
		}}, nil
	}

	return make([]iam.Resource, 0), nil
}

func genGeneralCacheResource(act iamtypes.ActionID, typ iamtypes.TypeID, att *meta.ResourceAttribute) ([]iam.Resource,
	error) {

	r := iam.Resource{
		System: iamtypes.SystemIDCMDB,
		Type:   iam.IamResourceType(iamtypes.GeneralCache),
	}

	if len(att.InstanceIDEx) > 0 {
		r.ID = att.InstanceIDEx
	}

	return []iam.Resource{r}, nil
}

func genTenantSetResource(act iamtypes.ActionID, typ iamtypes.TypeID,
	attribute *meta.ResourceAttribute) ([]iam.Resource, error) {

	r := iam.Resource{
		System: iamtypes.SystemIDCMDB,
		Type:   iam.IamResourceType(typ),
	}

	// compatible for authorize any
	if attribute.InstanceID > 0 {
		r.ID = strconv.FormatInt(attribute.InstanceID, 10)
	}

	return []iam.Resource{r}, nil
}
