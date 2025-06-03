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
	iamtypes "configcenter/src/ac/iam/types"
	"configcenter/src/common/metadata"
	"configcenter/src/thirdparty/apigw/iam"
)

var (
	businessChain = iam.ResourceChain{
		SystemID: iamtypes.SystemIDCMDB,
		ID:       iamtypes.Business,
	}
)

// GenerateInstanceSelections generate all the instance selections registered to IAM.
func GenerateInstanceSelections(tenantObjects map[string][]metadata.Object) []iam.InstanceSelection {
	instSelections := GenerateStaticInstanceSelections()
	instSelections = append(instSelections, genDynamicInstanceSelections(tenantObjects)...)
	return instSelections
}

// GenerateStaticInstanceSelections TODO
func GenerateStaticInstanceSelections() []iam.InstanceSelection {
	return []iam.InstanceSelection{
		{ID: iamtypes.BusinessSelection, Name: "业务列表", NameEn: "Business List",
			ResourceTypeChain: []iam.ResourceChain{businessChain}},
		{ID: iamtypes.BizSetSelection, Name: "业务集列表", NameEn: "Business Set List",
			ResourceTypeChain: []iam.ResourceChain{{SystemID: iamtypes.SystemIDCMDB, ID: iamtypes.BizSet}}},
		{ID: iamtypes.ProjectSelection, Name: "项目列表", NameEn: "Project List",
			ResourceTypeChain: []iam.ResourceChain{{SystemID: iamtypes.SystemIDCMDB, ID: iamtypes.Project}}},
		{ID: iamtypes.BusinessHostTransferSelection, Name: "业务主机选择", NameEn: "Business",
			ResourceTypeChain: []iam.ResourceChain{{SystemID: iamtypes.SystemIDCMDB,
				ID: iamtypes.BusinessForHostTrans}}},
		{ID: iamtypes.SysResourcePoolDirectorySelection, Name: "主机池目录列表", NameEn: "Host Pool Directory List",
			ResourceTypeChain: []iam.ResourceChain{{SystemID: iamtypes.SystemIDCMDB,
				ID: iamtypes.SysResourcePoolDirectory}}},
		{ID: iamtypes.SysHostRscPoolDirectorySelection, Name: "主机池主机选择", NameEn: "Host Pool Directory",
			ResourceTypeChain: []iam.ResourceChain{{SystemID: iamtypes.SystemIDCMDB,
				ID: iamtypes.SysHostRscPoolDirectory}}},
		{ID: iamtypes.BizHostInstanceSelection, Name: "业务主机列表", NameEn: "Business Host List",
			ResourceTypeChain: []iam.ResourceChain{
				// select the business at first.
				businessChain,
				// then select the host instances.
				{SystemID: iamtypes.SystemIDCMDB, ID: iamtypes.Host}}},
		{ID: iamtypes.BizCustomQuerySelection, Name: "业务动态分组列表", NameEn: "Business Dynamic Grouping List",
			ResourceTypeChain: []iam.ResourceChain{businessChain,
				{SystemID: iamtypes.SystemIDCMDB, ID: iamtypes.BizCustomQuery}}},
		{ID: iamtypes.BizProcessServiceTemplateSelection, Name: "服务模版列表", NameEn: "Service Template List",
			ResourceTypeChain: []iam.ResourceChain{businessChain,
				{SystemID: iamtypes.SystemIDCMDB, ID: iamtypes.BizProcessServiceTemplate}}},
		{ID: iamtypes.BizSetTemplateSelection, Name: "集群模板列表", NameEn: "Set Template List",
			ResourceTypeChain: []iam.ResourceChain{businessChain,
				{SystemID: iamtypes.SystemIDCMDB, ID: iamtypes.BizSetTemplate}}},
		{ID: iamtypes.SysHostInstanceSelection, Name: "主机池主机列表", NameEn: "Host Pool List",
			ResourceTypeChain: []iam.ResourceChain{{SystemID: iamtypes.SystemIDCMDB,
				ID: iamtypes.SysHostRscPoolDirectory},
				{SystemID: iamtypes.SystemIDCMDB, ID: iamtypes.Host}}},
		{ID: iamtypes.SysCloudAreaSelection, Name: "管控区域列表", NameEn: "Cloud Area List",
			ResourceTypeChain: []iam.ResourceChain{{SystemID: iamtypes.SystemIDCMDB, ID: iamtypes.SysCloudArea}}},
		{ID: iamtypes.SysInstanceModelSelection, Name: "实例模型列表", NameEn: "Instance Model List",
			ResourceTypeChain: []iam.ResourceChain{{SystemID: iamtypes.SystemIDCMDB, ID: iamtypes.SysInstanceModel}}},
		{ID: iamtypes.SysModelSelection, Name: "模型列表", NameEn: "Model List", ResourceTypeChain: []iam.ResourceChain{
			{SystemID: iamtypes.SystemIDCMDB, ID: iamtypes.SysModel}}},
		{ID: iamtypes.SysAssociationTypeSelection, Name: "关联类型", NameEn: "Association Type",
			ResourceTypeChain: []iam.ResourceChain{{SystemID: iamtypes.SystemIDCMDB, ID: iamtypes.SysAssociationType}}},
		{ID: iamtypes.SysModelGroupSelection, Name: "模型分组", NameEn: "Model Group",
			ResourceTypeChain: []iam.ResourceChain{{SystemID: iamtypes.SystemIDCMDB, ID: iamtypes.SysModelGroup}}},
		{ID: iamtypes.SysModelEventSelection, Name: "模型事件列表", NameEn: "Model Event List",
			ResourceTypeChain: []iam.ResourceChain{{SystemID: iamtypes.SystemIDCMDB, ID: iamtypes.SysModelEvent}}},
		{ID: iamtypes.MainlineModelEventSelection, Name: "自定义拓扑层级列表", NameEn: "Custom Topo Layer Event List",
			ResourceTypeChain: []iam.ResourceChain{{SystemID: iamtypes.SystemIDCMDB,
				ID: iamtypes.MainlineModelEvent}}},
		{ID: iamtypes.InstAsstEventSelection, Name: "实例关联事件列表", NameEn: "Instance Association Event List",
			ResourceTypeChain: []iam.ResourceChain{{SystemID: iamtypes.SystemIDCMDB, ID: iamtypes.InstAsstEvent}}},
		{ID: iamtypes.KubeWorkloadEventSelection, Name: "容器工作负载事件列表", NameEn: "Kube Workload Event List",
			ResourceTypeChain: []iam.ResourceChain{{SystemID: iamtypes.SystemIDCMDB, ID: iamtypes.KubeWorkloadEvent}}},
		{ID: iamtypes.FieldGroupingTemplateSelection, Name: "字段组合模板列表", NameEn: "Field Grouping Template List",
			ResourceTypeChain: []iam.ResourceChain{{SystemID: iamtypes.SystemIDCMDB,
				ID: iamtypes.FieldGroupingTemplate}}},
		{ID: iamtypes.GeneralCacheSelection, Name: "通用缓存类型列表", NameEn: "General Resource Cache Type List",
			ResourceTypeChain: []iam.ResourceChain{{SystemID: iamtypes.SystemIDCMDB, ID: iamtypes.GeneralCache}}},
		// only for other system's biz topo instance selection usage, not related to actions
		{ID: iamtypes.BizTopoSelection, Name: "业务拓扑", NameEn: "Business Topology",
			ResourceTypeChain: []iam.ResourceChain{businessChain, {SystemID: iamtypes.SystemIDCMDB, ID: iamtypes.Set},
				{SystemID: iamtypes.SystemIDCMDB, ID: iamtypes.Module}}},
		{ID: iamtypes.TenantSetSelection, Name: "租户集列表", NameEn: "Tenant Set List",
			ResourceTypeChain: []iam.ResourceChain{{SystemID: iamtypes.SystemIDCMDB, ID: iamtypes.TenantSet}}},
	}
}
