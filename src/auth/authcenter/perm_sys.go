/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2019 THL A29 Limited, a Tencent company. All rights reserved.
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
	"strings"
)

var expectSystem = System{
	SystemID:           SystemIDCMDB,
	SystemName:         SystemNameCMDB,
	Desc:               "蓝鲸配置平台（CMDB）",
	ReleatedScopeTypes: strings.Join([]string{ScopeTypeIDBiz, ScopeTypeIDSystem}, ";"),
	Managers:           "system;admin",
	Creator:            "system",
	Updater:            "system",
}

var expectSystemResourceType = []ResourceType{
	{
		ResourceTypeID:       SysModelGroup,
		ResourceTypeName:     "模型分组",
		ParentResourceTypeID: "",
		Actions: []Action{
			{
				ActionID:          Create,
				ActionName:        "新建",
				IsRelatedResource: false,
			},
			{
				ActionID:          Edit,
				ActionName:        "编辑",
				IsRelatedResource: true,
			},
			{
				ActionID:          Delete,
				ActionName:        "删除",
				IsRelatedResource: true,
			},
		},
	},
	{
		ResourceTypeID:       SysModel,
		ResourceTypeName:     "模型",
		ParentResourceTypeID: SysModelGroup,
		Actions: []Action{
			{
				ActionID:          Create,
				ActionName:        "新建",
				IsRelatedResource: false,
			},
			{
				ActionID:          Edit,
				ActionName:        "编辑",
				IsRelatedResource: true,
			},
			{
				ActionID:          Delete,
				ActionName:        "删除",
				IsRelatedResource: true,
			},
		},
	},
	{
		ResourceTypeID:       SysInstance,
		ResourceTypeName:     "实例",
		ParentResourceTypeID: SysModel,
		Actions: []Action{
			{
				ActionID:          Create,
				ActionName:        "新建",
				IsRelatedResource: false,
			},
			{
				ActionID:          Edit,
				ActionName:        "编辑",
				IsRelatedResource: true,
			},
			{
				ActionID:          Delete,
				ActionName:        "删除",
				IsRelatedResource: true,
			},
			{
				ActionID:          Get,
				ActionName:        "查询",
				IsRelatedResource: true,
			},
		},
	},
	{
		ResourceTypeID:       SysAssociationType,
		ResourceTypeName:     "关联类型",
		ParentResourceTypeID: "",
		Actions: []Action{
			{
				ActionID:          Create,
				ActionName:        "新建",
				IsRelatedResource: false,
			},
			{
				ActionID:          Edit,
				ActionName:        "编辑",
				IsRelatedResource: true,
			},
			{
				ActionID:          Delete,
				ActionName:        "删除",
				IsRelatedResource: true,
			},
		},
	},
	{
		ResourceTypeID:       SysBusinessInstance,
		ResourceTypeName:     "业务",
		ParentResourceTypeID: "",
		Actions: []Action{
			{
				ActionID:          Create,
				ActionName:        "新建",
				IsRelatedResource: false,
			},
			{
				ActionID:          Edit,
				ActionName:        "编辑",
				IsRelatedResource: true,
			},
			{
				ActionID:          Archive,
				ActionName:        "归档",
				IsRelatedResource: true,
			},
			{
				ActionID:          Get,
				ActionName:        "查询",
				IsRelatedResource: true,
			},
		},
	},
	{
		ResourceTypeID:       SysHostInstance,
		ResourceTypeName:     "主机",
		ParentResourceTypeID: "",
		Actions: []Action{
			{
				ActionID:          Create,
				ActionName:        "新建",
				IsRelatedResource: false,
			},
			{
				ActionID:          Edit,
				ActionName:        "编辑",
				IsRelatedResource: true,
			},
			{
				ActionID:          Delete,
				ActionName:        "编辑",
				IsRelatedResource: true,
			},
			{
				ActionID:          Get,
				ActionName:        "查询",
				IsRelatedResource: true,
			},
			{
				ActionID:          ModuleTransfer,
				ActionName:        "分配到业务",
				IsRelatedResource: true,
			},
		},
	},
	{
		ResourceTypeID:       SysEventPushing,
		ResourceTypeName:     "事件推送",
		ParentResourceTypeID: "",
		Actions: []Action{
			{
				ActionID:          Create,
				ActionName:        "新建",
				IsRelatedResource: false,
			},
			{
				ActionID:          Edit,
				ActionName:        "编辑",
				IsRelatedResource: true,
			},
			{
				ActionID:          Delete,
				ActionName:        "编辑",
				IsRelatedResource: true,
			},
			{
				ActionID:          Get,
				ActionName:        "查询",
				IsRelatedResource: true,
			},
		},
	},
	{
		ResourceTypeID:       SysSystemBase,
		ResourceTypeName:     "系统基础",
		ParentResourceTypeID: "",
		Actions: []Action{
			{
				ActionID:          ModelTopologyOperation,
				ActionName:        "拓扑层级管理",
				IsRelatedResource: false,
			},
			{
				ActionID:          AdminEntrance,
				ActionName:        "管理页面入口",
				IsRelatedResource: false,
			},
			{
				ActionID:          ModelTopologyView,
				ActionName:        "模型拓扑视图",
				IsRelatedResource: false,
			},
		},
	},
}

var expectModelGroupResourceInst = RegisterInfo{
	CreatorID:   "system",
	CreatorType: "user",
	Resources: []ResourceEntity{
		{
			ResourceType: SysModelGroup,
			ResourceID: []ResourceID{
				{ResourceType: SysModelGroup, ResourceID: "bk_middleware"},
			},
			ResourceName: "中间件",
			ScopeInfo: ScopeInfo{
				ScopeType: "system",
				ScopeID:   SystemIDCMDB,
			},
		},
		{
			ResourceType: SysModelGroup,
			ResourceID: []ResourceID{
				{ResourceType: SysModelGroup, ResourceID: "bk_network"},
			},
			ResourceName: "网络",
			ScopeInfo: ScopeInfo{
				ScopeType: "system",
				ScopeID:   SystemIDCMDB,
			},
		},
	},
}

var expectModelResourceInst = RegisterInfo{
	CreatorID:   "system",
	CreatorType: "user",
	Resources: []ResourceEntity{
		{
			ResourceType: SysModel,
			ResourceID: []ResourceID{
				{ResourceType: SysModelGroup, ResourceID: "bk_network"},
				{ResourceType: SysModel, ResourceID: "bk_switch"},
			},
			ResourceName: "交换机",
			ScopeInfo: ScopeInfo{
				ScopeType: "system",
				ScopeID:   SystemIDCMDB,
			},
		},
		{
			ResourceType: SysModel,
			ResourceID: []ResourceID{
				{ResourceType: SysModelGroup, ResourceID: "bk_network"},
				{ResourceType: SysModel, ResourceID: "bk_firewall"},
			},
			ResourceName: "防火墙",
			ScopeInfo: ScopeInfo{
				ScopeType: "system",
				ScopeID:   SystemIDCMDB,
			},
		},
		{
			ResourceType: SysModel,
			ResourceID: []ResourceID{
				{ResourceType: SysModelGroup, ResourceID: "bk_network"},
				{ResourceType: SysModel, ResourceID: "bk_load_blance"},
			},
			ResourceName: "负载均衡",
			ScopeInfo: ScopeInfo{
				ScopeType: "system",
				ScopeID:   SystemIDCMDB,
			},
		},
		{
			ResourceType: SysModel,
			ResourceID: []ResourceID{
				{ResourceType: SysModelGroup, ResourceID: "bk_network"},
				{ResourceType: SysModel, ResourceID: "bk_router"},
			},
			ResourceName: "路由器",
			ScopeInfo: ScopeInfo{
				ScopeType: "system",
				ScopeID:   SystemIDCMDB,
			},
		},
		{
			ResourceType: SysModel,
			ResourceID: []ResourceID{
				{ResourceType: SysModelGroup, ResourceID: "bk_middleware"},
				{ResourceType: SysModel, ResourceID: "bk_apache"},
			},
			ResourceName: "apache",
			ScopeInfo: ScopeInfo{
				ScopeType: "system",
				ScopeID:   SystemIDCMDB,
			},
		},
		{
			ResourceType: SysModel,
			ResourceID: []ResourceID{
				{ResourceType: SysModelGroup, ResourceID: "bk_middleware"},
				{ResourceType: SysModel, ResourceID: "bk_weblogic"},
			},

			ResourceName: "weblogic",
			ScopeInfo: ScopeInfo{
				ScopeType: "system",
				ScopeID:   SystemIDCMDB,
			},
		},
		{
			ResourceType: SysModel,
			ResourceID: []ResourceID{
				{ResourceType: SysModelGroup, ResourceID: "bk_middleware"},
				{ResourceType: SysModel, ResourceID: "bk_tomcat"},
			},
			ResourceName: "bk_tomcat",
			ScopeInfo: ScopeInfo{
				ScopeType: "system",
				ScopeID:   SystemIDCMDB,
			},
		},
	},
}
