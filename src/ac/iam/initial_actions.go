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
	businessResource = RelateResourceType{
		SystemID:    SystemIDCMDB,
		ID:          Business,
		NameAlias:   "",
		NameAliasEn: "",
		Scope:       nil,
		InstanceSelections: []RelatedInstanceSelection{{
			SystemID: SystemIDCMDB,
			ID:       BusinessSelection,
		}},
	}

	resourcePoolDirResource = RelateResourceType{
		SystemID:    SystemIDCMDB,
		ID:          SysResourcePoolDirectory,
		NameAlias:   "",
		NameAliasEn: "",
		Scope:       nil,
		InstanceSelections: []RelatedInstanceSelection{{
			SystemID: SystemIDCMDB,
			ID:       SysResourcePoolDirectorySelection,
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
	resourceActionList = append(resourceActionList, genModelInstanceActions()...)
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
	resourceActionList = append(resourceActionList, genEventWatchActions()...)
	resourceActionList = append(resourceActionList, genConfigAdminActions()...)

	return resourceActionList
}

func genBusinessHostActions() []ResourceAction {
	hostSelection := []RelatedInstanceSelection{{
		SystemID: SystemIDCMDB,
		ID:       BizHostInstanceSelection,
	}}

	relatedResource := []RelateResourceType{{
		SystemID:           SystemIDCMDB,
		ID:                 Host,
		NameAlias:          "",
		NameAliasEn:        "",
		Scope:              nil,
		InstanceSelections: hostSelection,
	}}

	actions := make([]ResourceAction, 0)

	// edit business's host actions
	actions = append(actions, ResourceAction{
		ID:                   EditBusinessHost,
		Name:                 "业务主机编辑",
		NameEn:               "Edit Business Hosts",
		Type:                 Edit,
		RelatedResourceTypes: relatedResource,
		RelatedActions:       []ActionID{ViewBusinessResource},
		Version:              1,
	})

	// business host transfer to resource pool actions
	actions = append(actions, ResourceAction{
		ID:                   BusinessHostTransferToResourcePool,
		Name:                 "主机归还主机池",
		NameEn:               "Return Hosts To Pool",
		Type:                 Edit,
		RelatedResourceTypes: []RelateResourceType{businessResource, resourcePoolDirResource},
		RelatedActions:       []ActionID{ViewBusinessResource},
		Version:              1,
	})

	businessHostResource := RelateResourceType{
		SystemID:    SystemIDCMDB,
		ID:          BusinessForHostTrans,
		NameAlias:   "",
		NameAliasEn: "",
		Scope:       nil,
		InstanceSelections: []RelatedInstanceSelection{{
			SystemID: SystemIDCMDB,
			ID:       BusinessHostTransferSelection,
		}},
	}

	// business host transfer to another business actions
	actions = append(actions, ResourceAction{
		ID:                   HostTransferAcrossBusiness,
		Name:                 "主机转移到其他业务",
		NameEn:               "Assigned Host To Other Business",
		Type:                 Edit,
		RelatedResourceTypes: []RelateResourceType{businessHostResource, businessResource},
		RelatedActions:       []ActionID{ViewBusinessResource},
		Version:              1,
	})

	return actions
}

func genBusinessCustomQueryActions() []ResourceAction {
	selection := []RelatedInstanceSelection{{
		SystemID: SystemIDCMDB,
		ID:       BizCustomQuerySelection,
	}}

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
		Name:                 "动态分组新建",
		NameEn:               "Create Dynamic Grouping",
		Type:                 Create,
		RelatedResourceTypes: []RelateResourceType{businessResource},
		RelatedActions:       []ActionID{ViewBusinessResource},
		Version:              1,
	})

	actions = append(actions, ResourceAction{
		ID:                   EditBusinessCustomQuery,
		Name:                 "动态分组编辑",
		NameEn:               "Edit Dynamic Grouping",
		Type:                 Edit,
		RelatedResourceTypes: relatedResource,
		RelatedActions:       []ActionID{ViewBusinessResource},
		Version:              1,
	})

	actions = append(actions, ResourceAction{
		ID:                   DeleteBusinessCustomQuery,
		Name:                 "动态分组删除",
		NameEn:               "Delete Dynamic Grouping",
		Type:                 Delete,
		RelatedResourceTypes: relatedResource,
		RelatedActions:       []ActionID{ViewBusinessResource},
		Version:              1,
	})

	return actions
}

func genBusinessCustomFieldActions() []ResourceAction {
	actions := make([]ResourceAction, 0)

	actions = append(actions, ResourceAction{
		ID:                   EditBusinessCustomField,
		Name:                 "业务自定义字段编辑",
		NameEn:               "Edit Custom Field",
		Type:                 Edit,
		RelatedResourceTypes: []RelateResourceType{businessResource},
		RelatedActions:       []ActionID{ViewBusinessResource},
		Version:              1,
	})

	return actions
}

func genBusinessServiceCategoryActions() []ResourceAction {
	actions := make([]ResourceAction, 0)

	actions = append(actions, ResourceAction{
		ID:                   CreateBusinessServiceCategory,
		Name:                 "服务分类新建",
		NameEn:               "Create Service Category",
		Type:                 Create,
		RelatedResourceTypes: []RelateResourceType{businessResource},
		RelatedActions:       []ActionID{ViewBusinessResource, EditBusinessServiceCategory, DeleteBusinessServiceCategory},
		Version:              1,
	})

	actions = append(actions, ResourceAction{
		ID:                   EditBusinessServiceCategory,
		Name:                 "服务分类编辑",
		NameEn:               "Edit Service Category",
		Type:                 Edit,
		RelatedResourceTypes: []RelateResourceType{businessResource},
		RelatedActions:       []ActionID{ViewBusinessResource},
		Version:              1,
	})

	actions = append(actions, ResourceAction{
		ID:                   DeleteBusinessServiceCategory,
		Name:                 "服务分类删除",
		NameEn:               "Delete Service Category",
		Type:                 Delete,
		RelatedResourceTypes: []RelateResourceType{businessResource},
		RelatedActions:       []ActionID{ViewBusinessResource},
		Version:              1,
	})

	return actions
}

func genBusinessServiceInstanceActions() []ResourceAction {
	actions := make([]ResourceAction, 0)

	actions = append(actions, ResourceAction{
		ID:                   CreateBusinessServiceInstance,
		Name:                 "服务实例新建",
		NameEn:               "Create Service Instance",
		Type:                 Create,
		RelatedResourceTypes: []RelateResourceType{businessResource},
		RelatedActions:       []ActionID{ViewBusinessResource, EditBusinessServiceInstance, DeleteBusinessServiceInstance},
		Version:              1,
	})

	actions = append(actions, ResourceAction{
		ID:                   EditBusinessServiceInstance,
		Name:                 "服务实例编辑",
		NameEn:               "Edit Service Instance",
		Type:                 Edit,
		RelatedResourceTypes: []RelateResourceType{businessResource},
		RelatedActions:       []ActionID{ViewBusinessResource},
		Version:              1,
	})

	actions = append(actions, ResourceAction{
		ID:                   DeleteBusinessServiceInstance,
		Name:                 "服务实例删除",
		NameEn:               "Delete Service Instance",
		Type:                 Delete,
		RelatedResourceTypes: []RelateResourceType{businessResource},
		RelatedActions:       []ActionID{ViewBusinessResource},
		Version:              1,
	})

	return actions
}

func genBusinessServiceTemplateActions() []ResourceAction {
	selection := []RelatedInstanceSelection{{
		SystemID:       SystemIDCMDB,
		ID:             BizProcessServiceTemplateSelection,
		IgnoreAuthPath: true,
	}}

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
		Name:                 "服务模板新建",
		NameEn:               "Create Service Template",
		Type:                 Create,
		RelatedResourceTypes: []RelateResourceType{businessResource},
		RelatedActions:       []ActionID{ViewBusinessResource},
		Version:              1,
	})

	actions = append(actions, ResourceAction{
		ID:                   EditBusinessServiceTemplate,
		Name:                 "服务模板编辑",
		NameEn:               "Edit Service Template",
		Type:                 Edit,
		RelatedResourceTypes: relatedResource,
		RelatedActions:       []ActionID{ViewBusinessResource},
		Version:              1,
	})

	actions = append(actions, ResourceAction{
		ID:                   DeleteBusinessServiceTemplate,
		Name:                 "服务模板删除",
		NameEn:               "Delete Service Template",
		Type:                 Delete,
		RelatedResourceTypes: relatedResource,
		RelatedActions:       []ActionID{ViewBusinessResource},
		Version:              1,
	})

	return actions
}

func genBusinessSetTemplateActions() []ResourceAction {
	selection := []RelatedInstanceSelection{{
		SystemID:       SystemIDCMDB,
		ID:             BizSetTemplateSelection,
		IgnoreAuthPath: true,
	}}

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
		Name:                 "集群模板新建",
		NameEn:               "Create Set Template",
		Type:                 Create,
		RelatedResourceTypes: []RelateResourceType{businessResource},
		RelatedActions:       []ActionID{ViewBusinessResource},
		Version:              1,
	})

	actions = append(actions, ResourceAction{
		ID:                   EditBusinessSetTemplate,
		Name:                 "集群模板编辑",
		NameEn:               "Edit Set Template",
		Type:                 Edit,
		RelatedResourceTypes: relatedResource,
		RelatedActions:       []ActionID{ViewBusinessResource},
		Version:              1,
	})

	actions = append(actions, ResourceAction{
		ID:                   DeleteBusinessSetTemplate,
		Name:                 "集群模板删除",
		NameEn:               "Delete Set Template",
		Type:                 Delete,
		RelatedResourceTypes: relatedResource,
		RelatedActions:       []ActionID{ViewBusinessResource},
		Version:              1,
	})

	return actions
}

func genBusinessTopologyActions() []ResourceAction {
	actions := make([]ResourceAction, 0)

	actions = append(actions, ResourceAction{
		ID:                   CreateBusinessTopology,
		Name:                 "业务拓扑新建",
		NameEn:               "Create Business Topo",
		Type:                 Create,
		RelatedResourceTypes: []RelateResourceType{businessResource},
		RelatedActions:       []ActionID{ViewBusinessResource, EditBusinessTopology, DeleteBusinessTopology},
		Version:              1,
	})

	actions = append(actions, ResourceAction{
		ID:                   EditBusinessTopology,
		Name:                 "业务拓扑编辑",
		NameEn:               "Edit Business Topo",
		Type:                 Edit,
		RelatedResourceTypes: []RelateResourceType{businessResource},
		RelatedActions:       []ActionID{ViewBusinessResource},
		Version:              1,
	})

	actions = append(actions, ResourceAction{
		ID:                   DeleteBusinessTopology,
		Name:                 "业务拓扑删除",
		NameEn:               "Delete Business Topo",
		Type:                 Delete,
		RelatedResourceTypes: []RelateResourceType{businessResource},
		RelatedActions:       []ActionID{ViewBusinessResource},
		Version:              1,
	})

	return actions
}

func genBusinessHostApplyActions() []ResourceAction {
	actions := make([]ResourceAction, 0)

	actions = append(actions, ResourceAction{
		ID:                   EditBusinessHostApply,
		Name:                 "主机自动应用编辑",
		NameEn:               "Edit Host Apply",
		Type:                 Edit,
		RelatedResourceTypes: []RelateResourceType{businessResource},
		RelatedActions:       []ActionID{ViewBusinessResource},
		Version:              1,
	})

	return actions
}

func genResourcePoolHostActions() []ResourceAction {
	hostSelection := []RelatedInstanceSelection{{
		SystemID: SystemIDCMDB,
		ID:       SysHostInstanceSelection,
	}}

	relatedResource := []RelateResourceType{{
		SystemID:           SystemIDCMDB,
		ID:                 Host,
		NameAlias:          "",
		NameAliasEn:        "",
		Scope:              nil,
		InstanceSelections: hostSelection,
	}}

	actions := make([]ResourceAction, 0)

	actions = append(actions, ResourceAction{
		ID:                   CreateResourcePoolHost,
		Name:                 "主机池主机创建",
		NameEn:               "Create Pool Hosts",
		Type:                 Create,
		RelatedResourceTypes: []RelateResourceType{resourcePoolDirResource},
		RelatedActions:       nil,
		Version:              1,
	})

	actions = append(actions, ResourceAction{
		ID:                   EditResourcePoolHost,
		Name:                 "主机池主机编辑",
		NameEn:               "Edit Pool Hosts",
		Type:                 Edit,
		RelatedResourceTypes: relatedResource,
		RelatedActions:       nil,
		Version:              1,
	})

	actions = append(actions, ResourceAction{
		ID:                   DeleteResourcePoolHost,
		Name:                 "主机池主机删除",
		NameEn:               "Delete Pool Hosts",
		Type:                 Delete,
		RelatedResourceTypes: relatedResource,
		RelatedActions:       nil,
		Version:              1,
	})

	relatedHostResource := []RelateResourceType{{
		SystemID:    SystemIDCMDB,
		ID:          SysHostRscPoolDirectory,
		NameAlias:   "",
		NameAliasEn: "",
		Scope:       nil,
		InstanceSelections: []RelatedInstanceSelection{{
			SystemID: SystemIDCMDB,
			ID:       SysHostRscPoolDirectorySelection,
		}},
	}}

	transferToBusinessRelatedResource := append(relatedHostResource, businessResource)
	actions = append(actions, ResourceAction{
		ID:                   ResourcePoolHostTransferToBusiness,
		Name:                 "主机池主机分配到业务",
		NameEn:               "Assigned Pool Hosts To Business",
		Type:                 Edit,
		RelatedResourceTypes: transferToBusinessRelatedResource,
		RelatedActions:       nil,
		Version:              1,
	})

	transferToDirectoryRelatedResource := append(relatedHostResource, resourcePoolDirResource)
	actions = append(actions, ResourceAction{
		ID:                   ResourcePoolHostTransferToDirectory,
		Name:                 "主机池主机分配到目录",
		NameEn:               "Assigned Pool Hosts To Directory",
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
		Name:                 "主机池目录创建",
		NameEn:               "Create Pool Directory",
		Type:                 Create,
		RelatedResourceTypes: nil,
		RelatedActions:       nil,
		Version:              1,
	})

	actions = append(actions, ResourceAction{
		ID:                   EditResourcePoolDirectory,
		Name:                 "主机池目录编辑",
		NameEn:               "Edit Pool Directory",
		Type:                 Edit,
		RelatedResourceTypes: []RelateResourceType{resourcePoolDirResource},
		RelatedActions:       nil,
		Version:              1,
	})

	actions = append(actions, ResourceAction{
		ID:                   DeleteResourcePoolDirectory,
		Name:                 "主机池目录删除",
		NameEn:               "Delete Pool Directory",
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
		Name:                 "业务创建",
		NameEn:               "Create Business",
		Type:                 Create,
		RelatedResourceTypes: nil,
		RelatedActions:       nil,
		Version:              1,
	})

	actions = append(actions, ResourceAction{
		ID:                   EditBusiness,
		Name:                 "业务编辑",
		NameEn:               "Edit Business",
		Type:                 Edit,
		RelatedResourceTypes: []RelateResourceType{businessResource},
		RelatedActions:       []ActionID{FindBusiness},
		Version:              1,
	})

	actions = append(actions, ResourceAction{
		ID:                   ArchiveBusiness,
		Name:                 "业务归档",
		NameEn:               "Archive Business",
		Type:                 Edit,
		RelatedResourceTypes: []RelateResourceType{businessResource},
		RelatedActions:       []ActionID{FindBusiness},
		Version:              1,
	})

	actions = append(actions, ResourceAction{
		ID:                   FindBusiness,
		Name:                 "业务查询",
		NameEn:               "View Business",
		Type:                 View,
		RelatedResourceTypes: []RelateResourceType{businessResource},
		RelatedActions:       nil,
		Version:              1,
	})

	actions = append(actions, ResourceAction{
		ID:     ViewBusinessResource,
		Name:   "业务访问",
		NameEn: "View Business Resource",
		Type:   View,
		// TODO add business collection resource
		RelatedResourceTypes: []RelateResourceType{businessResource},
		RelatedActions:       nil,
		Version:              1,
	})

	return actions
}

func genCloudAreaActions() []ResourceAction {
	selection := []RelatedInstanceSelection{{
		SystemID: SystemIDCMDB,
		ID:       SysCloudAreaSelection,
	}}

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
		Name:                 "云区域创建",
		NameEn:               "Create Cloud Area",
		Type:                 Create,
		RelatedResourceTypes: nil,
		RelatedActions:       nil,
		Version:              1,
	})

	actions = append(actions, ResourceAction{
		ID:                   EditCloudArea,
		Name:                 "云区域编辑",
		NameEn:               "Edit Cloud Area",
		Type:                 Edit,
		RelatedResourceTypes: relatedResource,
		RelatedActions:       nil,
		Version:              1,
	})

	actions = append(actions, ResourceAction{
		ID:                   DeleteCloudArea,
		Name:                 "云区域删除",
		NameEn:               "Delete Cloud Area",
		Type:                 Delete,
		RelatedResourceTypes: relatedResource,
		RelatedActions:       nil,
		Version:              1,
	})

	return actions
}

func genModelInstanceActions() []ResourceAction {
	selection := []RelatedInstanceSelection{{
		SystemID: SystemIDCMDB,
		ID:       SysInstanceSelection,
	}}

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
		ID:     CreateSysInstance,
		Name:   "实例创建",
		NameEn: "Create Instance",
		Type:   Create,
		RelatedResourceTypes: []RelateResourceType{
			{
				SystemID:    SystemIDCMDB,
				ID:          SysInstanceModel,
				NameAlias:   "",
				NameAliasEn: "",
				Scope:       nil,
				InstanceSelections: []RelatedInstanceSelection{{
					SystemID: SystemIDCMDB,
					ID:       SysInstanceModelSelection,
				}},
			},
		},
		RelatedActions: nil,
		Version:        1,
	})

	actions = append(actions, ResourceAction{
		ID:                   EditSysInstance,
		Name:                 "实例编辑",
		NameEn:               "Edit Instance",
		Type:                 Edit,
		RelatedResourceTypes: relatedResource,
		RelatedActions:       nil,
		Version:              1,
	})

	actions = append(actions, ResourceAction{
		ID:                   DeleteSysInstance,
		Name:                 "实例删除",
		NameEn:               "Delete Instance",
		Type:                 Delete,
		RelatedResourceTypes: relatedResource,
		RelatedActions:       nil,
		Version:              1,
	})

	// actions = append(actions, ResourceAction{
	// 	ID:                   FindSysInstance,
	// 	Name:                 "实例查询",
	// 	NameEn:               "View Instance",
	// 	Type:                 View,
	// 	RelatedResourceTypes: relatedResource,
	// 	RelatedActions:       nil,
	// 	Version:              1,
	// })

	return actions
}

func genEventPushingActions() []ResourceAction {
	selection := []RelatedInstanceSelection{{
		SystemID: SystemIDCMDB,
		ID:       SysEventPushingSelection,
	}}

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
		Name:                 "事件订阅新建",
		NameEn:               "Create Event Subscription",
		Type:                 Create,
		RelatedResourceTypes: nil,
		RelatedActions:       nil,
		Version:              1,
	})

	actions = append(actions, ResourceAction{
		ID:                   EditEventPushing,
		Name:                 "事件订阅编辑",
		NameEn:               "Edit Event Subscription",
		Type:                 Edit,
		RelatedResourceTypes: relatedResource,
		RelatedActions:       []ActionID{FindEventPushing},
		Version:              1,
	})

	actions = append(actions, ResourceAction{
		ID:                   DeleteEventPushing,
		Name:                 "事件订阅删除",
		NameEn:               "Delete Event Subscription",
		Type:                 Delete,
		RelatedResourceTypes: relatedResource,
		RelatedActions:       []ActionID{FindEventPushing},
		Version:              1,
	})

	actions = append(actions, ResourceAction{
		ID:                   FindEventPushing,
		Name:                 "事件订阅查询",
		NameEn:               "View Event Subscription",
		Type:                 View,
		RelatedResourceTypes: relatedResource,
		RelatedActions:       nil,
		Version:              1,
	})

	return actions
}

func genCloudAccountActions() []ResourceAction {
	selection := []RelatedInstanceSelection{{
		SystemID: SystemIDCMDB,
		ID:       SysCloudAccountSelection,
	}}

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
		Name:                 "云账户新建",
		NameEn:               "Create Cloud Account",
		Type:                 Create,
		RelatedResourceTypes: nil,
		RelatedActions:       nil,
		Version:              1,
	})

	actions = append(actions, ResourceAction{
		ID:                   EditCloudAccount,
		Name:                 "云账户编辑",
		NameEn:               "Edit Cloud Account",
		Type:                 Edit,
		RelatedResourceTypes: relatedResource,
		RelatedActions:       []ActionID{FindCloudAccount},
		Version:              1,
	})

	actions = append(actions, ResourceAction{
		ID:                   DeleteCloudAccount,
		Name:                 "云账户删除",
		NameEn:               "Delete Cloud Account",
		Type:                 Delete,
		RelatedResourceTypes: relatedResource,
		RelatedActions:       []ActionID{FindCloudAccount},
		Version:              1,
	})

	actions = append(actions, ResourceAction{
		ID:                   FindCloudAccount,
		Name:                 "云账户查询",
		NameEn:               "View Cloud Account",
		Type:                 View,
		RelatedResourceTypes: relatedResource,
		RelatedActions:       nil,
		Version:              1,
	})

	return actions
}

func genCloudResourceTaskActions() []ResourceAction {
	selection := []RelatedInstanceSelection{{
		SystemID: SystemIDCMDB,
		ID:       SysCloudResourceTaskSelection,
	}}

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
		Name:                 "云资源任务新建",
		NameEn:               "Create Cloud Resource Task",
		Type:                 Create,
		RelatedResourceTypes: nil,
		RelatedActions:       nil,
		Version:              1,
	})

	actions = append(actions, ResourceAction{
		ID:                   EditCloudResourceTask,
		Name:                 "云资源任务编辑",
		NameEn:               "Edit Cloud Resource Task",
		Type:                 Edit,
		RelatedResourceTypes: relatedResource,
		RelatedActions:       []ActionID{FindCloudResourceTask},
		Version:              1,
	})

	actions = append(actions, ResourceAction{
		ID:                   DeleteCloudResourceTask,
		Name:                 "云资源任务删除",
		NameEn:               "Delete Cloud Resource Task",
		Type:                 Delete,
		RelatedResourceTypes: relatedResource,
		RelatedActions:       []ActionID{FindCloudResourceTask},
		Version:              1,
	})

	actions = append(actions, ResourceAction{
		ID:                   FindCloudResourceTask,
		Name:                 "云资源任务查询",
		NameEn:               "View Cloud Resource Task",
		Type:                 View,
		RelatedResourceTypes: relatedResource,
		RelatedActions:       nil,
		Version:              1,
	})

	return actions
}

func genModelActions() []ResourceAction {
	selection := []RelatedInstanceSelection{{
		SystemID: SystemIDCMDB,
		ID:       SysModelSelection,
	}}

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
		ID:     CreateSysModel,
		Name:   "模型新建",
		NameEn: "Create Model",
		Type:   Create,
		RelatedResourceTypes: []RelateResourceType{
			{
				SystemID:    SystemIDCMDB,
				ID:          SysModelGroup,
				NameAlias:   "",
				NameAliasEn: "",
				Scope:       nil,
				InstanceSelections: []RelatedInstanceSelection{{
					SystemID: SystemIDCMDB,
					ID:       SysModelGroupSelection,
				}},
			},
		},
		RelatedActions: nil,
		Version:        1,
	})

	actions = append(actions, ResourceAction{
		ID:                   EditSysModel,
		Name:                 "模型编辑",
		NameEn:               "Edit Model",
		Type:                 Edit,
		RelatedResourceTypes: relatedResource,
		RelatedActions:       nil,
		Version:              1,
	})

	actions = append(actions, ResourceAction{
		ID:                   DeleteSysModel,
		Name:                 "模型删除",
		NameEn:               "Delete Model",
		Type:                 Delete,
		RelatedResourceTypes: relatedResource,
		RelatedActions:       nil,
		Version:              1,
	})

	return actions
}

func genAssociationTypeActions() []ResourceAction {
	selection := []RelatedInstanceSelection{{
		SystemID: SystemIDCMDB,
		ID:       SysAssociationTypeSelection,
	}}

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
		Name:                 "关联类型新建",
		NameEn:               "Create Association Type",
		Type:                 Create,
		RelatedResourceTypes: nil,
		RelatedActions:       nil,
		Version:              1,
	})

	actions = append(actions, ResourceAction{
		ID:                   EditAssociationType,
		Name:                 "关联类型编辑",
		NameEn:               "Edit Association Type",
		Type:                 Edit,
		RelatedResourceTypes: relatedResource,
		RelatedActions:       nil,
		Version:              1,
	})

	actions = append(actions, ResourceAction{
		ID:                   DeleteAssociationType,
		Name:                 "关联类型删除",
		NameEn:               "Delete Association Type",
		Type:                 Delete,
		RelatedResourceTypes: relatedResource,
		RelatedActions:       nil,
		Version:              1,
	})

	return actions
}

func genModelGroupActions() []ResourceAction {
	selection := []RelatedInstanceSelection{{
		SystemID: SystemIDCMDB,
		ID:       SysModelGroupSelection,
	}}

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
		Name:                 "模型分组新建",
		NameEn:               "Create Model Group",
		Type:                 Create,
		RelatedResourceTypes: nil,
		RelatedActions:       nil,
		Version:              1,
	})

	actions = append(actions, ResourceAction{
		ID:                   EditModelGroup,
		Name:                 "模型分组编辑",
		NameEn:               "Edit Model Group",
		Type:                 Edit,
		RelatedResourceTypes: relatedResource,
		RelatedActions:       nil,
		Version:              1,
	})

	actions = append(actions, ResourceAction{
		ID:                   DeleteModelGroup,
		Name:                 "模型分组删除",
		NameEn:               "Delete Model Group",
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
		Name:                 "业务层级编辑",
		NameEn:               "Edit Business Level",
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
		Name:                 "模型拓扑视图编辑",
		NameEn:               "Edit Model Topo View",
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
		Name:                 "运营统计查询",
		NameEn:               "View Operational Statistics",
		Type:                 View,
		RelatedResourceTypes: nil,
		RelatedActions:       nil,
		Version:              1,
	})

	actions = append(actions, ResourceAction{
		ID:                   EditOperationStatistic,
		Name:                 "运营统计编辑",
		NameEn:               "Edit Operational Statistics",
		Type:                 Edit,
		RelatedResourceTypes: nil,
		RelatedActions:       []ActionID{FindOperationStatistic},
		Version:              1,
	})

	return actions
}

func genAuditLogActions() []ResourceAction {
	actions := make([]ResourceAction, 0)
	actions = append(actions, ResourceAction{
		ID:                   FindAuditLog,
		Name:                 "操作审计查询",
		NameEn:               "View Operation Audit",
		Type:                 View,
		RelatedResourceTypes: nil,
		RelatedActions:       nil,
		Version:              1,
	})
	return actions
}

func genEventWatchActions() []ResourceAction {
	actions := make([]ResourceAction, 0)
	actions = append(actions, ResourceAction{
		ID:                   WatchHostEvent,
		Name:                 "主机事件监听",
		NameEn:               "Host Event Listen",
		Type:                 View,
		RelatedResourceTypes: nil,
		RelatedActions:       nil,
		Version:              1,
	})

	actions = append(actions, ResourceAction{
		ID:                   WatchHostRelationEvent,
		Name:                 "主机关系事件监听",
		NameEn:               "Host Relation Event Listen",
		Type:                 View,
		RelatedResourceTypes: nil,
		RelatedActions:       nil,
		Version:              1,
	})

	actions = append(actions, ResourceAction{
		ID:                   WatchBizEvent,
		Name:                 "业务事件监听",
		NameEn:               "Business Event Listen",
		Type:                 View,
		RelatedResourceTypes: nil,
		RelatedActions:       nil,
		Version:              1,
	})

	actions = append(actions, ResourceAction{
		ID:                   WatchSetEvent,
		Name:                 "集群事件监听",
		NameEn:               "Set Event Listen",
		Type:                 View,
		RelatedResourceTypes: nil,
		RelatedActions:       nil,
		Version:              1,
	})

	actions = append(actions, ResourceAction{
		ID:                   WatchModuleEvent,
		Name:                 "模块数据监听",
		NameEn:               "Module Event Listen",
		Type:                 View,
		RelatedResourceTypes: nil,
		RelatedActions:       nil,
		Version:              1,
	})

	actions = append(actions, ResourceAction{
		ID:                   WatchSetTemplateEvent,
		Name:                 "集群模板数据监听",
		NameEn:               "Set Template Event Listen",
		Type:                 View,
		RelatedResourceTypes: nil,
		RelatedActions:       nil,
		Version:              1,
	})
	return actions
}

func genConfigAdminActions() []ResourceAction {
	actions := make([]ResourceAction, 0)
	actions = append(actions, ResourceAction{
		ID:                   GlobalSettings,
		Name:                 "全局设置",
		NameEn:               "Global Settings",
		Type:                 Edit,
		RelatedResourceTypes: nil,
		RelatedActions:       nil,
		Version:              1,
	})
	return actions
}
