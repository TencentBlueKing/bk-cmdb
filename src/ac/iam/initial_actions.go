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
		SystemID: systemID,
		ID:       Business,
	}
)

// GenerateActions generate all the actions registered to IAM.
func GenerateActions() []ResourceAction {
	resourceActionList := make([]ResourceAction, 0)
	// add business resource actions
	resourceActionList = append(resourceActionList, genBusinessHostActions()...)
	resourceActionList = append(resourceActionList, genBusinessCustomQueryActions()...)
	resourceActionList = append(resourceActionList, genBusinessModelActions()...)
	resourceActionList = append(resourceActionList, genBusinessModelGroupActions()...)
	resourceActionList = append(resourceActionList, genBusinessServiceCategoryActions()...)
	resourceActionList = append(resourceActionList, genBusinessServiceInstanceActions()...)
	resourceActionList = append(resourceActionList, genBusinessServiceTemplateActions()...)
	resourceActionList = append(resourceActionList, genBusinessSetTemplateActions()...)
	resourceActionList = append(resourceActionList, genBusinessTopologyActions()...)

	// add public resource actions

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
				SystemID: systemID,
				ID:       BizHostInstance,
			},
		},
	}}

	relatedResource := []RelateResourceType{{
		SystemID:           systemID,
		ID:                 BizHostInstance,
		NameAlias:          "",
		NameAliasEn:        "",
		Scope:              nil,
		InstanceSelections: hostSelection,
	}}

	actions := make([]ResourceAction, 0)

	// create business's host action
	actions = append(actions, ResourceAction{
		ID:                   CreateBusinessHost,
		Name:                 "新建业务主机",
		NameEn:               "create business's host",
		Type:                 Create,
		RelatedResourceTypes: relatedResource,
		RelatedActions:       nil,
		Version:              0,
	})

	// edit business's host actions
	actions = append(actions, ResourceAction{
		ID:                   EditBusinessHost,
		Name:                 "编辑业务主机",
		NameEn:               "edit business's host",
		Type:                 Edit,
		RelatedResourceTypes: relatedResource,
		RelatedActions:       nil,
		Version:              0,
	})

	// remove business host actions
	actions = append(actions, ResourceAction{
		ID:                   RemoveBusinessHost,
		Name:                 "编辑业务主机",
		NameEn:               "edit business's host",
		Type:                 Edit,
		RelatedResourceTypes: relatedResource,
		RelatedActions:       nil,
		Version:              0,
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
					SystemID: systemID,
					ID:       BizCustomQuery,
				},
			},
		},
	}

	relatedResource := []RelateResourceType{
		{
			SystemID:           systemID,
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
		RelatedResourceTypes: relatedResource,
		RelatedActions:       nil,
		Version:              0,
	})

	actions = append(actions, ResourceAction{
		ID:                   EditBusinessCustomQuery,
		Name:                 "编辑动态分组",
		NameEn:               "edit custom query",
		Type:                 Edit,
		RelatedResourceTypes: relatedResource,
		RelatedActions:       nil,
		Version:              0,
	})

	actions = append(actions, ResourceAction{
		ID:                   DeleteBusinessCustomQuery,
		Name:                 "删除动态分组",
		NameEn:               "delete custom query",
		Type:                 Delete,
		RelatedResourceTypes: relatedResource,
		RelatedActions:       nil,
		Version:              0,
	})

	actions = append(actions, ResourceAction{
		ID:                   FindBusinessCustomQuery,
		Name:                 "查询动态分组",
		NameEn:               "find custom query",
		Type:                 View,
		RelatedResourceTypes: relatedResource,
		RelatedActions:       nil,
		Version:              0,
	})

	return actions
}

func genBusinessModelActions() []ResourceAction {

	selection := []InstanceSelection{
		{
			Name:   "业务模型",
			NameEn: "Business's Model",
			ResourceTypeChain: []ResourceChain{
				businessChain,
				{
					SystemID: systemID,
					ID:       BizModel,
				},
			},
		},
	}

	relatedResource := []RelateResourceType{
		{
			SystemID:           systemID,
			ID:                 BizModel,
			NameAlias:          "",
			NameAliasEn:        "",
			Scope:              nil,
			InstanceSelections: selection,
		},
	}

	actions := make([]ResourceAction, 0)
	actions = append(actions, ResourceAction{
		ID:                   CreateBusinessModel,
		Name:                 "创建业务模型",
		NameEn:               "create business's model",
		Type:                 Create,
		RelatedResourceTypes: relatedResource,
		RelatedActions:       nil,
		Version:              0,
	})

	actions = append(actions, ResourceAction{
		ID:                   EditBusinessModel,
		Name:                 "编辑业务模型",
		NameEn:               "edit business's model",
		Type:                 Edit,
		RelatedResourceTypes: relatedResource,
		RelatedActions:       nil,
		Version:              0,
	})

	actions = append(actions, ResourceAction{
		ID:                   DeleteBusinessModel,
		Name:                 "删除业务模型",
		NameEn:               "delete business's model",
		Type:                 Delete,
		RelatedResourceTypes: relatedResource,
		RelatedActions:       nil,
		Version:              0,
	})

	return actions
}

func genBusinessModelGroupActions() []ResourceAction {

	selection := []InstanceSelection{
		{
			Name:   "业务模型分组",
			NameEn: "Business Model's Group",
			ResourceTypeChain: []ResourceChain{
				businessChain,
				{
					SystemID: systemID,
					ID:       BizModelGroup,
				},
			},
		},
	}

	relatedResource := []RelateResourceType{
		{
			SystemID:           systemID,
			ID:                 BizModelGroup,
			NameAlias:          "",
			NameAliasEn:        "",
			Scope:              nil,
			InstanceSelections: selection,
		},
	}

	actions := make([]ResourceAction, 0)
	actions = append(actions, ResourceAction{
		ID:                   CreateBusinessModelGroup,
		Name:                 "创建业务模型分组",
		NameEn:               "create business's model group",
		Type:                 Create,
		RelatedResourceTypes: relatedResource,
		RelatedActions:       nil,
		Version:              0,
	})

	actions = append(actions, ResourceAction{
		ID:                   EditBusinessModelGroup,
		Name:                 "编辑业务模型分组",
		NameEn:               "edit business's model group",
		Type:                 Edit,
		RelatedResourceTypes: relatedResource,
		RelatedActions:       nil,
		Version:              0,
	})

	actions = append(actions, ResourceAction{
		ID:                   DeleteBusinessModelGroup,
		Name:                 "删除业务模型分组",
		NameEn:               "delete business's model group",
		Type:                 Delete,
		RelatedResourceTypes: relatedResource,
		RelatedActions:       nil,
		Version:              0,
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
					SystemID: systemID,
					ID:       BizProcessServiceCategory,
				},
			},
		},
	}

	relatedResource := []RelateResourceType{
		{
			SystemID:           systemID,
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
		RelatedResourceTypes: relatedResource,
		RelatedActions:       nil,
		Version:              0,
	})

	actions = append(actions, ResourceAction{
		ID:                   EditBusinessServiceCategory,
		Name:                 "编辑服务分类",
		NameEn:               "edit service category",
		Type:                 Edit,
		RelatedResourceTypes: relatedResource,
		RelatedActions:       nil,
		Version:              0,
	})

	actions = append(actions, ResourceAction{
		ID:                   DeleteBusinessServiceCategory,
		Name:                 "删除服务分类",
		NameEn:               "delete service category",
		Type:                 Delete,
		RelatedResourceTypes: relatedResource,
		RelatedActions:       nil,
		Version:              0,
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
					SystemID: systemID,
					ID:       BizProcessServiceInstance,
				},
			},
		},
	}

	relatedResource := []RelateResourceType{
		{
			SystemID:           systemID,
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
		RelatedResourceTypes: relatedResource,
		RelatedActions:       nil,
		Version:              0,
	})

	actions = append(actions, ResourceAction{
		ID:                   EditBusinessServiceInstance,
		Name:                 "编辑服务实例",
		NameEn:               "edit service instance",
		Type:                 Edit,
		RelatedResourceTypes: relatedResource,
		RelatedActions:       nil,
		Version:              0,
	})

	actions = append(actions, ResourceAction{
		ID:                   DeleteBusinessServiceInstance,
		Name:                 "删除服务实例",
		NameEn:               "delete service instance",
		Type:                 Delete,
		RelatedResourceTypes: relatedResource,
		RelatedActions:       nil,
		Version:              0,
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
					SystemID: systemID,
					ID:       BizProcessServiceTemplate,
				},
			},
		},
	}

	relatedResource := []RelateResourceType{
		{
			SystemID:           systemID,
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
		RelatedResourceTypes: relatedResource,
		RelatedActions:       nil,
		Version:              0,
	})

	actions = append(actions, ResourceAction{
		ID:                   EditBusinessServiceTemplate,
		Name:                 "编辑服务模板",
		NameEn:               "edit service template",
		Type:                 Edit,
		RelatedResourceTypes: relatedResource,
		RelatedActions:       nil,
		Version:              0,
	})

	actions = append(actions, ResourceAction{
		ID:                   DeleteBusinessServiceTemplate,
		Name:                 "删除服务模板",
		NameEn:               "delete service template",
		Type:                 Delete,
		RelatedResourceTypes: relatedResource,
		RelatedActions:       nil,
		Version:              0,
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
					SystemID: systemID,
					ID:       BizSetTemplate,
				},
			},
		},
	}

	relatedResource := []RelateResourceType{
		{
			SystemID:           systemID,
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
		RelatedResourceTypes: relatedResource,
		RelatedActions:       nil,
		Version:              0,
	})

	actions = append(actions, ResourceAction{
		ID:                   EditBusinessSetTemplate,
		Name:                 "编辑集群模板",
		NameEn:               "edit set template",
		Type:                 Edit,
		RelatedResourceTypes: relatedResource,
		RelatedActions:       nil,
		Version:              0,
	})

	actions = append(actions, ResourceAction{
		ID:                   DeleteBusinessSetTemplate,
		Name:                 "删除集群模板",
		NameEn:               "delete set template",
		Type:                 Delete,
		RelatedResourceTypes: relatedResource,
		RelatedActions:       nil,
		Version:              0,
	})

	return actions
}

func genBusinessTopologyActions() []ResourceAction {

	selection := []InstanceSelection{
		{
			Name:   "业务拓扑",
			NameEn: "Business Topology",
			ResourceTypeChain: []ResourceChain{
				businessChain,
				{
					SystemID: systemID,
					ID:       BizTopology,
				},
			},
		},
	}

	relatedResource := []RelateResourceType{
		{
			SystemID:           systemID,
			ID:                 BizTopology,
			NameAlias:          "",
			NameAliasEn:        "",
			Scope:              nil,
			InstanceSelections: selection,
		},
	}

	actions := make([]ResourceAction, 0)
	actions = append(actions, ResourceAction{
		ID:                   CreateBusinessTopology,
		Name:                 "创建业务拓扑",
		NameEn:               "create business topology",
		Type:                 Create,
		RelatedResourceTypes: relatedResource,
		RelatedActions:       nil,
		Version:              0,
	})

	actions = append(actions, ResourceAction{
		ID:                   EditBusinessTopology,
		Name:                 "编辑业务拓扑",
		NameEn:               "edit business topology",
		Type:                 Edit,
		RelatedResourceTypes: relatedResource,
		RelatedActions:       nil,
		Version:              0,
	})

	actions = append(actions, ResourceAction{
		ID:                   DeleteBusinessTopology,
		Name:                 "删除业务拓扑",
		NameEn:               "delete business topology",
		Type:                 Delete,
		RelatedResourceTypes: relatedResource,
		RelatedActions:       nil,
		Version:              0,
	})

	return actions
}
