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

var (
	businessChain = ResourceChain{
		SystemID: SystemIDCMDB,
		ID:       Business,
	}

	businessResource = RelateResourceType{
		SystemID:    SystemIDCMDB,
		ID:          Business,
		NameAlias:   "",
		NameAliasEn: "",
		Scope:       nil,
		InstanceSelections: []InstanceSelection{{
			Name:              "业务",
			NameEn:            "business",
			ResourceTypeChain: []ResourceChain{businessChain},
		}},
	}

	resourcePoolDirResource = RelateResourceType{
		SystemID:    SystemIDCMDB,
		ID:          SysResourcePoolDirectory,
		NameAlias:   "",
		NameAliasEn: "",
		Scope:       nil,
		InstanceSelections: []InstanceSelection{{
			Name:   "资源池目录",
			NameEn: "Resource Pool Directory",
			ResourceTypeChain: []ResourceChain{{
				SystemID: SystemIDCMDB,
				ID:       SysResourcePoolDirectory,
			}},
		}},
	}
)

// GenerateActions generate all the actions registered to IAM.
func GenerateActions() []ResourceAction {
	resourceActionList := make([]ResourceAction, 0)
	// add business resource actions
	resourceActionList = append(resourceActionList, genBusinessHostActions()...)
	resourceActionList = append(resourceActionList, genBusinessCustomQueryActions()...)
	resourceActionList = append(resourceActionList, genBusinessCustomFieldActions()...)
	resourceActionList = append(resourceActionList, genBusinessServiceCategoryActions()...)
	resourceActionList = append(resourceActionList, genBusinessServiceInstanceActions()...)
	resourceActionList = append(resourceActionList, genBusinessServiceTemplateActions()...)
	resourceActionList = append(resourceActionList, genBusinessSetTemplateActions()...)
	resourceActionList = append(resourceActionList, genBusinessTopologyActions()...)
	resourceActionList = append(resourceActionList, genBusinessHostApplyActions()...)

	// add public resource actions
	resourceActionList = append(resourceActionList, genResourcePoolHostActions()...)
	resourceActionList = append(resourceActionList, genResourcePoolDirectoryActions()...)
	resourceActionList = append(resourceActionList, genBusinessActions()...)
	resourceActionList = append(resourceActionList, genCloudAreaActions()...)
	resourceActionList = append(resourceActionList, genInstanceActions()...)
	resourceActionList = append(resourceActionList, genEventPushingActions()...)
	resourceActionList = append(resourceActionList, genCloudAccountActions()...)
	resourceActionList = append(resourceActionList, genCloudResourceTaskActions()...)
	resourceActionList = append(resourceActionList, genModelActions()...)
	resourceActionList = append(resourceActionList, genAssociationTypeActions()...)
	resourceActionList = append(resourceActionList, genModelGroupActions()...)
	resourceActionList = append(resourceActionList, genBusinessLayerActions()...)
	resourceActionList = append(resourceActionList, genModelTopologyViewActions()...)
	resourceActionList = append(resourceActionList, genOperationStatisticActions()...)
	resourceActionList = append(resourceActionList, genAuditLogActions()...)

	return resourceActionList
}

func genBusinessHostActions() []ResourceAction {
	hostSelection := []InstanceSelection{{
		Name:   "业务主机",
		NameEn: "Business's Hosts",
		ResourceTypeChain: []ResourceChain{
			// select the business at first.
			businessChain,
			// then select the host instances.
			{
				SystemID: SystemIDCMDB,
				ID:       BizHostInstance,
			},
		},
	}}

	relatedResource := []RelateResourceType{{
		SystemID:           SystemIDCMDB,
		ID:                 BizHostInstance,
		NameAlias:          "",
		NameAliasEn:        "",
		Scope:              nil,
		InstanceSelections: hostSelection,
	}}

	actions := make([]ResourceAction, 0)

	// edit business's host actions
	actions = append(actions, ResourceAction{
		ID:                   EditBusinessHost,
		Name:                 "编辑业务主机",
		NameEn:               "edit business's host",
		Type:                 Edit,
		RelatedResourceTypes: relatedResource,
		RelatedActions:       nil,
		Version:              1,
	})

	transferToResourcePoolRelatedResource := append(relatedResource, resourcePoolDirResource)
	// business host transfer to resource pool actions
	actions = append(actions, ResourceAction{
		ID:                   BusinessHostTransferToResourcePool,
		Name:                 "业务主机归还资源池",
		NameEn:               "transfer business's host to resource pool",
		Type:                 Edit,
		RelatedResourceTypes: transferToResourcePoolRelatedResource,
		RelatedActions:       nil,
		Version:              1,
	})

	return actions
}

func genBusinessCustomQueryActions() []ResourceAction {
	selection := []InstanceSelection{
		{
			Name:   "动态分组",
			NameEn: "Custom Query",
			ResourceTypeChain: []ResourceChain{
				businessChain,
				{
					SystemID: SystemIDCMDB,
					ID:       BizCustomQuery,
				},
			},
		},
	}

	relatedResource := []RelateResourceType{
		{
			SystemID:           SystemIDCMDB,
			ID:                 BizCustomQuery,
			NameAlias:          "",
			NameAliasEn:        "",
			Scope:              nil,
			InstanceSelections: selection,
		},
	}

	actions := make([]ResourceAction, 0)
	actions = append(actions, ResourceAction{
		ID:                   CreateBusinessCustomQuery,
		Name:                 "创建动态分组",
		NameEn:               "create custom query",
		Type:                 Create,
		RelatedResourceTypes: []RelateResourceType{businessResource},
		RelatedActions:       nil,
		Version:              1,
	})

	actions = append(actions, ResourceAction{
		ID:                   EditBusinessCustomQuery,
		Name:                 "编辑动态分组",
		NameEn:               "edit custom query",
		Type:                 Edit,
		RelatedResourceTypes: relatedResource,
		RelatedActions:       nil,
		Version:              1,
	})

	actions = append(actions, ResourceAction{
		ID:                   DeleteBusinessCustomQuery,
		Name:                 "删除动态分组",
		NameEn:               "delete custom query",
		Type:                 Delete,
		RelatedResourceTypes: relatedResource,
		RelatedActions:       nil,
		Version:              1,
	})

	actions = append(actions, ResourceAction{
		ID:                   FindBusinessCustomQuery,
		Name:                 "查询动态分组",
		NameEn:               "find custom query",
		Type:                 View,
		RelatedResourceTypes: relatedResource,
		RelatedActions:       nil,
		Version:              1,
	})

	return actions
}

func genBusinessCustomFieldActions() []ResourceAction {
	actions := make([]ResourceAction, 0)

	actions = append(actions, ResourceAction{
		ID:                   EditBusinessCustomField,
		Name:                 "编辑业务自定义字段",
		NameEn:               "edit business's custom field",
		Type:                 Edit,
		RelatedResourceTypes: []RelateResourceType{businessResource},
		RelatedActions:       nil,
		Version:              1,
	})

	return actions
}

func genBusinessServiceCategoryActions() []ResourceAction {
	selection := []InstanceSelection{
		{
			Name:   "服务分类",
			NameEn: "Service Category",
			ResourceTypeChain: []ResourceChain{
				businessChain,
				{
					SystemID: SystemIDCMDB,
					ID:       BizProcessServiceCategory,
				},
			},
		},
	}

	relatedResource := []RelateResourceType{
		{
			SystemID:           SystemIDCMDB,
			ID:                 BizProcessServiceCategory,
			NameAlias:          "",
			NameAliasEn:        "",
			Scope:              nil,
			InstanceSelections: selection,
		},
	}

	actions := make([]ResourceAction, 0)
	actions = append(actions, ResourceAction{
		ID:                   CreateBusinessServiceCategory,
		Name:                 "创建服务分类",
		NameEn:               "create service category",
		Type:                 Create,
		RelatedResourceTypes: []RelateResourceType{businessResource},
		RelatedActions:       nil,
		Version:              1,
	})

	actions = append(actions, ResourceAction{
		ID:                   EditBusinessServiceCategory,
		Name:                 "编辑服务分类",
		NameEn:               "edit service category",
		Type:                 Edit,
		RelatedResourceTypes: relatedResource,
		RelatedActions:       nil,
		Version:              1,
	})

	actions = append(actions, ResourceAction{
		ID:                   DeleteBusinessServiceCategory,
		Name:                 "删除服务分类",
		NameEn:               "delete service category",
		Type:                 Delete,
		RelatedResourceTypes: relatedResource,
		RelatedActions:       nil,
		Version:              1,
	})

	return actions
}

func genBusinessServiceInstanceActions() []ResourceAction {
	selection := []InstanceSelection{
		{
			Name:   "服务实例",
			NameEn: "Service Instance",
			ResourceTypeChain: []ResourceChain{
				businessChain,
				{
					SystemID: SystemIDCMDB,
					ID:       BizProcessServiceInstance,
				},
			},
		},
	}

	relatedResource := []RelateResourceType{
		{
			SystemID:           SystemIDCMDB,
			ID:                 BizProcessServiceInstance,
			NameAlias:          "",
			NameAliasEn:        "",
			Scope:              nil,
			InstanceSelections: selection,
		},
	}

	actions := make([]ResourceAction, 0)
	actions = append(actions, ResourceAction{
		ID:                   CreateBusinessServiceInstance,
		Name:                 "创建服务实例",
		NameEn:               "create service instance",
		Type:                 Create,
		RelatedResourceTypes: []RelateResourceType{businessResource},
		RelatedActions:       nil,
		Version:              1,
	})

	actions = append(actions, ResourceAction{
		ID:                   EditBusinessServiceInstance,
		Name:                 "编辑服务实例",
		NameEn:               "edit service instance",
		Type:                 Edit,
		RelatedResourceTypes: relatedResource,
		RelatedActions:       nil,
		Version:              1,
	})

	actions = append(actions, ResourceAction{
		ID:                   DeleteBusinessServiceInstance,
		Name:                 "删除服务实例",
		NameEn:               "delete service instance",
		Type:                 Delete,
		RelatedResourceTypes: relatedResource,
		RelatedActions:       nil,
		Version:              1,
	})

	return actions
}

func genBusinessServiceTemplateActions() []ResourceAction {
	selection := []InstanceSelection{
		{
			Name:   "服务模板",
			NameEn: "Service Template",
			ResourceTypeChain: []ResourceChain{
				businessChain,
				{
					SystemID: SystemIDCMDB,
					ID:       BizProcessServiceTemplate,
				},
			},
		},
	}

	relatedResource := []RelateResourceType{
		{
			SystemID:           SystemIDCMDB,
			ID:                 BizProcessServiceTemplate,
			NameAlias:          "",
			NameAliasEn:        "",
			Scope:              nil,
			InstanceSelections: selection,
		},
	}

	actions := make([]ResourceAction, 0)
	actions = append(actions, ResourceAction{
		ID:                   CreateBusinessServiceTemplate,
		Name:                 "创建服务模板",
		NameEn:               "create service template",
		Type:                 Create,
		RelatedResourceTypes: []RelateResourceType{businessResource},
		RelatedActions:       nil,
		Version:              1,
	})

	actions = append(actions, ResourceAction{
		ID:                   EditBusinessServiceTemplate,
		Name:                 "编辑服务模板",
		NameEn:               "edit service template",
		Type:                 Edit,
		RelatedResourceTypes: relatedResource,
		RelatedActions:       nil,
		Version:              1,
	})

	actions = append(actions, ResourceAction{
		ID:                   DeleteBusinessServiceTemplate,
		Name:                 "删除服务模板",
		NameEn:               "delete service template",
		Type:                 Delete,
		RelatedResourceTypes: relatedResource,
		RelatedActions:       nil,
		Version:              1,
	})

	return actions
}

func genBusinessSetTemplateActions() []ResourceAction {
	selection := []InstanceSelection{
		{
			Name:   "集群模板",
			NameEn: "Set Template",
			ResourceTypeChain: []ResourceChain{
				businessChain,
				{
					SystemID: SystemIDCMDB,
					ID:       BizSetTemplate,
				},
			},
		},
	}

	relatedResource := []RelateResourceType{
		{
			SystemID:           SystemIDCMDB,
			ID:                 BizSetTemplate,
			NameAlias:          "",
			NameAliasEn:        "",
			Scope:              nil,
			InstanceSelections: selection,
		},
	}

	actions := make([]ResourceAction, 0)
	actions = append(actions, ResourceAction{
		ID:                   CreateBusinessSetTemplate,
		Name:                 "创建集群模板",
		NameEn:               "create set template",
		Type:                 Create,
		RelatedResourceTypes: []RelateResourceType{businessResource},
		RelatedActions:       nil,
		Version:              1,
	})

	actions = append(actions, ResourceAction{
		ID:                   EditBusinessSetTemplate,
		Name:                 "编辑集群模板",
		NameEn:               "edit set template",
		Type:                 Edit,
		RelatedResourceTypes: relatedResource,
		RelatedActions:       nil,
		Version:              1,
	})

	actions = append(actions, ResourceAction{
		ID:                   DeleteBusinessSetTemplate,
		Name:                 "删除集群模板",
		NameEn:               "delete set template",
		Type:                 Delete,
		RelatedResourceTypes: relatedResource,
		RelatedActions:       nil,
		Version:              1,
	})

	return actions
}

func genBusinessTopologyActions() []ResourceAction {
	actions := make([]ResourceAction, 0)

	actions = append(actions, ResourceAction{
		ID:                   CreateBusinessTopology,
		Name:                 "创建业务拓扑",
		NameEn:               "create business topology",
		Type:                 Create,
		RelatedResourceTypes: []RelateResourceType{businessResource},
		RelatedActions:       nil,
		Version:              1,
	})

	actions = append(actions, ResourceAction{
		ID:                   EditBusinessTopology,
		Name:                 "编辑业务拓扑",
		NameEn:               "edit business topology",
		Type:                 Edit,
		RelatedResourceTypes: []RelateResourceType{businessResource},
		RelatedActions:       nil,
		Version:              1,
	})

	actions = append(actions, ResourceAction{
		ID:                   DeleteBusinessTopology,
		Name:                 "删除业务拓扑",
		NameEn:               "delete business topology",
		Type:                 Delete,
		RelatedResourceTypes: []RelateResourceType{businessResource},
		RelatedActions:       nil,
		Version:              1,
	})

	return actions
}

func genBusinessHostApplyActions() []ResourceAction {
	actions := make([]ResourceAction, 0)

	actions = append(actions, ResourceAction{
		ID:                   EditBusinessHostApply,
		Name:                 "编辑主机属性自动应用",
		NameEn:               "edit business's host apply rule",
		Type:                 Edit,
		RelatedResourceTypes: []RelateResourceType{businessResource},
		RelatedActions:       nil,
		Version:              1,
	})

	return actions
}

func genResourcePoolHostActions() []ResourceAction {
	hostSelection := []InstanceSelection{{
		Name:   "资源池主机",
		NameEn: "Resource Pool Host",
		ResourceTypeChain: []ResourceChain{
			{
				SystemID: SystemIDCMDB,
				ID:       SysResourcePoolDirectory,
			},
			{
				SystemID: SystemIDCMDB,
				ID:       SysHostInstance,
			},
		},
	}}

	relatedResource := []RelateResourceType{{
		SystemID:           SystemIDCMDB,
		ID:                 SysHostInstance,
		NameAlias:          "",
		NameAliasEn:        "",
		Scope:              nil,
		InstanceSelections: hostSelection,
	}}

	actions := make([]ResourceAction, 0)

	actions = append(actions, ResourceAction{
		ID:                   CreateResourcePoolHost,
		Name:                 "创建资源池主机",
		NameEn:               "create resource pool host",
		Type:                 Create,
		RelatedResourceTypes: []RelateResourceType{resourcePoolDirResource},
		RelatedActions:       nil,
		Version:              1,
	})

	actions = append(actions, ResourceAction{
		ID:                   EditResourcePoolHost,
		Name:                 "编辑资源池主机",
		NameEn:               "edit resource pool host",
		Type:                 Edit,
		RelatedResourceTypes: relatedResource,
		RelatedActions:       nil,
		Version:              1,
	})

	actions = append(actions, ResourceAction{
		ID:                   DeleteResourcePoolHost,
		Name:                 "删除资源池主机",
		NameEn:               "delete resource pool host",
		Type:                 Delete,
		RelatedResourceTypes: relatedResource,
		RelatedActions:       nil,
		Version:              1,
	})

	transferToBusinessRelatedResource := append(relatedResource, businessResource)
	actions = append(actions, ResourceAction{
		ID:                   ResourcePoolHostTransferToBusiness,
		Name:                 "资源池主机分配到业务",
		NameEn:               "transfer resource pool host to business",
		Type:                 Edit,
		RelatedResourceTypes: transferToBusinessRelatedResource,
		RelatedActions:       nil,
		Version:              1,
	})

	transferToDirectoryRelatedResource := append(relatedResource, resourcePoolDirResource)
	actions = append(actions, ResourceAction{
		ID:                   ResourcePoolHostTransferToDirectory,
		Name:                 "资源池主机分配到目录",
		NameEn:               "transfer resource pool host to resource pool directory",
		Type:                 Edit,
		RelatedResourceTypes: transferToDirectoryRelatedResource,
		RelatedActions:       nil,
		Version:              1,
	})

	return actions
}

func genResourcePoolDirectoryActions() []ResourceAction {
	actions := make([]ResourceAction, 0)
	actions = append(actions, ResourceAction{
		ID:                   CreateResourcePoolDirectory,
		Name:                 "创建资源池目录",
		NameEn:               "create resource pool directory",
		Type:                 Create,
		RelatedResourceTypes: nil,
		RelatedActions:       nil,
		Version:              1,
	})

	actions = append(actions, ResourceAction{
		ID:                   EditResourcePoolDirectory,
		Name:                 "编辑资源池目录",
		NameEn:               "edit resource pool directory",
		Type:                 Edit,
		RelatedResourceTypes: []RelateResourceType{resourcePoolDirResource},
		RelatedActions:       nil,
		Version:              1,
	})

	actions = append(actions, ResourceAction{
		ID:                   DeleteResourcePoolDirectory,
		Name:                 "删除资源池目录",
		NameEn:               "delete resource pool directory",
		Type:                 Delete,
		RelatedResourceTypes: []RelateResourceType{resourcePoolDirResource},
		RelatedActions:       nil,
		Version:              1,
	})

	return actions
}

func genBusinessActions() []ResourceAction {
	actions := make([]ResourceAction, 0)
	actions = append(actions, ResourceAction{
		ID:                   CreateBusiness,
		Name:                 "创建业务",
		NameEn:               "create business",
		Type:                 Create,
		RelatedResourceTypes: nil,
		RelatedActions:       nil,
		Version:              1,
	})

	actions = append(actions, ResourceAction{
		ID:                   EditBusiness,
		Name:                 "编辑业务",
		NameEn:               "edit business",
		Type:                 Edit,
		RelatedResourceTypes: []RelateResourceType{businessResource},
		RelatedActions:       nil,
		Version:              1,
	})

	actions = append(actions, ResourceAction{
		ID:                   ArchiveBusiness,
		Name:                 "归档业务",
		NameEn:               "archive business",
		Type:                 Edit,
		RelatedResourceTypes: []RelateResourceType{businessResource},
		RelatedActions:       nil,
		Version:              1,
	})

	actions = append(actions, ResourceAction{
		ID:                   FindBusiness,
		Name:                 "查询业务",
		NameEn:               "find business",
		Type:                 View,
		RelatedResourceTypes: []RelateResourceType{businessResource},
		RelatedActions:       nil,
		Version:              1,
	})

	return actions
}

func genCloudAreaActions() []ResourceAction {
	selection := []InstanceSelection{
		{
			Name:   "云区域",
			NameEn: "Cloud Area",
			ResourceTypeChain: []ResourceChain{
				{
					SystemID: SystemIDCMDB,
					ID:       SysCloudArea,
				},
			},
		},
	}

	relatedResource := []RelateResourceType{
		{
			SystemID:           SystemIDCMDB,
			ID:                 SysCloudArea,
			NameAlias:          "",
			NameAliasEn:        "",
			Scope:              nil,
			InstanceSelections: selection,
		},
	}

	actions := make([]ResourceAction, 0)
	actions = append(actions, ResourceAction{
		ID:                   CreateCloudArea,
		Name:                 "创建云区域",
		NameEn:               "create cloud area",
		Type:                 Create,
		RelatedResourceTypes: nil,
		RelatedActions:       nil,
		Version:              1,
	})

	actions = append(actions, ResourceAction{
		ID:                   EditCloudArea,
		Name:                 "编辑云区域",
		NameEn:               "edit cloud area",
		Type:                 Edit,
		RelatedResourceTypes: relatedResource,
		RelatedActions:       nil,
		Version:              1,
	})

	actions = append(actions, ResourceAction{
		ID:                   DeleteCloudArea,
		Name:                 "删除云区域",
		NameEn:               "delete cloud area",
		Type:                 Delete,
		RelatedResourceTypes: relatedResource,
		RelatedActions:       nil,
		Version:              1,
	})

	return actions
}

func genInstanceActions() []ResourceAction {
	selection := []InstanceSelection{
		{
			Name:   "实例",
			NameEn: "Instance",
			ResourceTypeChain: []ResourceChain{
				{
					SystemID: SystemIDCMDB,
					ID:       SysModel,
				},
				{
					SystemID: SystemIDCMDB,
					ID:       SysInstance,
				},
			},
		},
	}

	relatedResource := []RelateResourceType{
		{
			SystemID:           SystemIDCMDB,
			ID:                 SysInstance,
			NameAlias:          "",
			NameAliasEn:        "",
			Scope:              nil,
			InstanceSelections: selection,
		},
	}

	actions := make([]ResourceAction, 0)
	actions = append(actions, ResourceAction{
		ID:                   CreateInstance,
		Name:                 "创建实例",
		NameEn:               "create instance",
		Type:                 Create,
		RelatedResourceTypes: nil,
		RelatedActions:       nil,
		Version:              1,
	})

	actions = append(actions, ResourceAction{
		ID:                   EditInstance,
		Name:                 "编辑实例",
		NameEn:               "edit instance",
		Type:                 Edit,
		RelatedResourceTypes: relatedResource,
		RelatedActions:       nil,
		Version:              1,
	})

	actions = append(actions, ResourceAction{
		ID:                   DeleteInstance,
		Name:                 "删除实例",
		NameEn:               "delete instance",
		Type:                 Delete,
		RelatedResourceTypes: relatedResource,
		RelatedActions:       nil,
		Version:              1,
	})

	actions = append(actions, ResourceAction{
		ID:                   FindInstance,
		Name:                 "查询实例",
		NameEn:               "find instance",
		Type:                 View,
		RelatedResourceTypes: relatedResource,
		RelatedActions:       nil,
		Version:              1,
	})

	return actions
}

func genEventPushingActions() []ResourceAction {
	selection := []InstanceSelection{
		{
			Name:   "事件推送",
			NameEn: "Event Pushing",
			ResourceTypeChain: []ResourceChain{
				{
					SystemID: SystemIDCMDB,
					ID:       SysEventPushing,
				},
			},
		},
	}

	relatedResource := []RelateResourceType{
		{
			SystemID:           SystemIDCMDB,
			ID:                 SysEventPushing,
			NameAlias:          "",
			NameAliasEn:        "",
			Scope:              nil,
			InstanceSelections: selection,
		},
	}

	actions := make([]ResourceAction, 0)
	actions = append(actions, ResourceAction{
		ID:                   CreateEventPushing,
		Name:                 "创建事件订阅",
		NameEn:               "create event subscription",
		Type:                 Create,
		RelatedResourceTypes: nil,
		RelatedActions:       nil,
		Version:              1,
	})

	actions = append(actions, ResourceAction{
		ID:                   EditEventPushing,
		Name:                 "编辑事件订阅",
		NameEn:               "edit event subscription",
		Type:                 Edit,
		RelatedResourceTypes: relatedResource,
		RelatedActions:       nil,
		Version:              1,
	})

	actions = append(actions, ResourceAction{
		ID:                   DeleteEventPushing,
		Name:                 "删除事件订阅",
		NameEn:               "delete event subscription",
		Type:                 Delete,
		RelatedResourceTypes: relatedResource,
		RelatedActions:       nil,
		Version:              1,
	})

	actions = append(actions, ResourceAction{
		ID:                   FindEventPushing,
		Name:                 "查询事件订阅",
		NameEn:               "find event subscription",
		Type:                 View,
		RelatedResourceTypes: relatedResource,
		RelatedActions:       nil,
		Version:              1,
	})

	return actions
}

func genCloudAccountActions() []ResourceAction {
	selection := []InstanceSelection{
		{
			Name:   "云账户",
			NameEn: "Cloud Account",
			ResourceTypeChain: []ResourceChain{
				{
					SystemID: SystemIDCMDB,
					ID:       SysCloudAccount,
				},
			},
		},
	}

	relatedResource := []RelateResourceType{
		{
			SystemID:           SystemIDCMDB,
			ID:                 SysCloudAccount,
			NameAlias:          "",
			NameAliasEn:        "",
			Scope:              nil,
			InstanceSelections: selection,
		},
	}

	actions := make([]ResourceAction, 0)
	actions = append(actions, ResourceAction{
		ID:                   CreateCloudAccount,
		Name:                 "创建云账户",
		NameEn:               "create cloud account",
		Type:                 Create,
		RelatedResourceTypes: nil,
		RelatedActions:       nil,
		Version:              1,
	})

	actions = append(actions, ResourceAction{
		ID:                   EditCloudAccount,
		Name:                 "编辑云账户",
		NameEn:               "edit cloud account",
		Type:                 Edit,
		RelatedResourceTypes: relatedResource,
		RelatedActions:       nil,
		Version:              1,
	})

	actions = append(actions, ResourceAction{
		ID:                   DeleteCloudAccount,
		Name:                 "删除云账户",
		NameEn:               "delete cloud account",
		Type:                 Delete,
		RelatedResourceTypes: relatedResource,
		RelatedActions:       nil,
		Version:              1,
	})

	actions = append(actions, ResourceAction{
		ID:                   FindCloudAccount,
		Name:                 "查询云账户",
		NameEn:               "find cloud account",
		Type:                 View,
		RelatedResourceTypes: relatedResource,
		RelatedActions:       nil,
		Version:              1,
	})

	return actions
}

func genCloudResourceTaskActions() []ResourceAction {
	selection := []InstanceSelection{
		{
			Name:   "云资源任务",
			NameEn: "Cloud Resource Task",
			ResourceTypeChain: []ResourceChain{
				{
					SystemID: SystemIDCMDB,
					ID:       SysCloudResourceTask,
				},
			},
		},
	}

	relatedResource := []RelateResourceType{
		{
			SystemID:           SystemIDCMDB,
			ID:                 SysCloudResourceTask,
			NameAlias:          "",
			NameAliasEn:        "",
			Scope:              nil,
			InstanceSelections: selection,
		},
	}

	actions := make([]ResourceAction, 0)
	actions = append(actions, ResourceAction{
		ID:                   CreateCloudResourceTask,
		Name:                 "创建云资源任务",
		NameEn:               "create cloud resource task",
		Type:                 Create,
		RelatedResourceTypes: nil,
		RelatedActions:       nil,
		Version:              1,
	})

	actions = append(actions, ResourceAction{
		ID:                   EditCloudResourceTask,
		Name:                 "编辑云资源任务",
		NameEn:               "edit cloud resource task",
		Type:                 Edit,
		RelatedResourceTypes: relatedResource,
		RelatedActions:       nil,
		Version:              1,
	})

	actions = append(actions, ResourceAction{
		ID:                   DeleteCloudResourceTask,
		Name:                 "删除云资源任务",
		NameEn:               "delete cloud resource task",
		Type:                 Delete,
		RelatedResourceTypes: relatedResource,
		RelatedActions:       nil,
		Version:              1,
	})

	actions = append(actions, ResourceAction{
		ID:                   FindCloudResourceTask,
		Name:                 "查询云资源任务",
		NameEn:               "find cloud resource task",
		Type:                 View,
		RelatedResourceTypes: relatedResource,
		RelatedActions:       nil,
		Version:              1,
	})

	return actions
}

func genModelActions() []ResourceAction {
	selection := []InstanceSelection{
		{
			Name:   "模型",
			NameEn: "Model",
			ResourceTypeChain: []ResourceChain{
				{
					SystemID: SystemIDCMDB,
					ID:       SysModel,
				},
			},
		},
	}

	relatedResource := []RelateResourceType{
		{
			SystemID:           SystemIDCMDB,
			ID:                 SysModel,
			NameAlias:          "",
			NameAliasEn:        "",
			Scope:              nil,
			InstanceSelections: selection,
		},
	}

	actions := make([]ResourceAction, 0)
	actions = append(actions, ResourceAction{
		ID:                   CreateModel,
		Name:                 "创建模型",
		NameEn:               "create model",
		Type:                 Create,
		RelatedResourceTypes: nil,
		RelatedActions:       nil,
		Version:              1,
	})

	actions = append(actions, ResourceAction{
		ID:                   EditModel,
		Name:                 "编辑模型",
		NameEn:               "edit model",
		Type:                 Edit,
		RelatedResourceTypes: relatedResource,
		RelatedActions:       nil,
		Version:              1,
	})

	actions = append(actions, ResourceAction{
		ID:                   DeleteModel,
		Name:                 "删除模型",
		NameEn:               "delete model",
		Type:                 Delete,
		RelatedResourceTypes: relatedResource,
		RelatedActions:       nil,
		Version:              1,
	})

	actions = append(actions, ResourceAction{
		ID:                   FindModel,
		Name:                 "查询模型",
		NameEn:               "find model",
		Type:                 View,
		RelatedResourceTypes: relatedResource,
		RelatedActions:       nil,
		Version:              1,
	})

	return actions
}

func genAssociationTypeActions() []ResourceAction {
	selection := []InstanceSelection{
		{
			Name:   "关联类型",
			NameEn: "Association Type",
			ResourceTypeChain: []ResourceChain{
				{
					SystemID: SystemIDCMDB,
					ID:       SysAssociationType,
				},
			},
		},
	}

	relatedResource := []RelateResourceType{
		{
			SystemID:           SystemIDCMDB,
			ID:                 SysAssociationType,
			NameAlias:          "",
			NameAliasEn:        "",
			Scope:              nil,
			InstanceSelections: selection,
		},
	}

	actions := make([]ResourceAction, 0)
	actions = append(actions, ResourceAction{
		ID:                   CreateAssociationType,
		Name:                 "创建关联类型",
		NameEn:               "create association type",
		Type:                 Create,
		RelatedResourceTypes: nil,
		RelatedActions:       nil,
		Version:              1,
	})

	actions = append(actions, ResourceAction{
		ID:                   EditAssociationType,
		Name:                 "编辑关联类型",
		NameEn:               "edit association type",
		Type:                 Edit,
		RelatedResourceTypes: relatedResource,
		RelatedActions:       nil,
		Version:              1,
	})

	actions = append(actions, ResourceAction{
		ID:                   DeleteAssociationType,
		Name:                 "删除关联类型",
		NameEn:               "delete association type",
		Type:                 Delete,
		RelatedResourceTypes: relatedResource,
		RelatedActions:       nil,
		Version:              1,
	})

	return actions
}

func genModelGroupActions() []ResourceAction {
	selection := []InstanceSelection{
		{
			Name:   "模型分组",
			NameEn: "Model Group",
			ResourceTypeChain: []ResourceChain{
				{
					SystemID: SystemIDCMDB,
					ID:       SysModelGroup,
				},
			},
		},
	}

	relatedResource := []RelateResourceType{
		{
			SystemID:           SystemIDCMDB,
			ID:                 SysModelGroup,
			NameAlias:          "",
			NameAliasEn:        "",
			Scope:              nil,
			InstanceSelections: selection,
		},
	}

	actions := make([]ResourceAction, 0)
	actions = append(actions, ResourceAction{
		ID:                   CreateModelGroup,
		Name:                 "创建模型分组",
		NameEn:               "create model group",
		Type:                 Create,
		RelatedResourceTypes: nil,
		RelatedActions:       nil,
		Version:              1,
	})

	actions = append(actions, ResourceAction{
		ID:                   EditModelGroup,
		Name:                 "编辑模型分组",
		NameEn:               "edit model group",
		Type:                 Edit,
		RelatedResourceTypes: relatedResource,
		RelatedActions:       nil,
		Version:              1,
	})

	actions = append(actions, ResourceAction{
		ID:                   DeleteModelGroup,
		Name:                 "删除模型分组",
		NameEn:               "delete model group",
		Type:                 Delete,
		RelatedResourceTypes: relatedResource,
		RelatedActions:       nil,
		Version:              1,
	})

	return actions
}

func genBusinessLayerActions() []ResourceAction {
	actions := make([]ResourceAction, 0)
	actions = append(actions, ResourceAction{
		ID:                   EditBusinessLayer,
		Name:                 "编辑业务层级",
		NameEn:               "edit business topology layer",
		Type:                 Edit,
		RelatedResourceTypes: nil,
		RelatedActions:       nil,
		Version:              1,
	})
	return actions
}

func genModelTopologyViewActions() []ResourceAction {
	actions := make([]ResourceAction, 0)
	actions = append(actions, ResourceAction{
		ID:                   EditModelTopologyView,
		Name:                 "编辑模型拓扑视图",
		NameEn:               "edit model topology view",
		Type:                 Edit,
		RelatedResourceTypes: nil,
		RelatedActions:       nil,
		Version:              1,
	})
	return actions
}

func genOperationStatisticActions() []ResourceAction {
	actions := make([]ResourceAction, 0)
	actions = append(actions, ResourceAction{
		ID:                   FindOperationStatistic,
		Name:                 "查询运营统计",
		NameEn:               "find operation statistic",
		Type:                 View,
		RelatedResourceTypes: nil,
		RelatedActions:       nil,
		Version:              1,
	})

	actions = append(actions, ResourceAction{
		ID:                   EditOperationStatistic,
		Name:                 "编辑运营统计",
		NameEn:               "edit operation statistic",
		Type:                 Edit,
		RelatedResourceTypes: nil,
		RelatedActions:       nil,
		Version:              1,
	})

	return actions
}

func genAuditLogActions() []ResourceAction {
	actions := make([]ResourceAction, 0)
	actions = append(actions, ResourceAction{
		ID:                   FindAuditLog,
		Name:                 "查询操作审计",
		NameEn:               "find audit log",
		Type:                 View,
		RelatedResourceTypes: nil,
		RelatedActions:       nil,
		Version:              1,
	})
	return actions
}
