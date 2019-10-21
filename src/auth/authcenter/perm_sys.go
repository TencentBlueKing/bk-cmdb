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
	Managers:           "admin",
	Creator:            "admin",
	Updater:            "admin",
}

var expectSystemResourceType = []ResourceType{
	{
		ResourceTypeID:       SysModelGroup,
		ResourceTypeName:     "模型分组",
		ParentResourceTypeID: "",
		Share:                true,
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
		ResourceTypeID:   SysModel,
		ResourceTypeName: "模型",
		Share:            true,
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
		Share:                true,
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
		Share:                true,
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
		Share:                true,
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
		ResourceTypeID:       SysEventPushing,
		ResourceTypeName:     "事件推送",
		ParentResourceTypeID: "",
		Share:                false,
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
		ResourceTypeID:       SysSystemBase,
		ResourceTypeName:     "系统基础",
		ParentResourceTypeID: "",
		Share:                false,
		Actions: []Action{
			{
				ActionID:          ModelTopologyOperation,
				ActionName:        "编辑业务层级",
				IsRelatedResource: false,
			},
			{
				ActionID:          ModelTopologyView,
				ActionName:        "编辑模型拓扑视图",
				IsRelatedResource: false,
			},
		},
	},
	{
		ResourceTypeID:       SysHostInstance,
		ResourceTypeName:     "主机（资源池）",
		ParentResourceTypeID: "",
		Share:                true,
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
		ResourceTypeID:       SysAuditLog,
		ResourceTypeName:     "操作审计",
		ParentResourceTypeID: "",
		Share:                false,
		Actions: []Action{
			{
				ActionID:          Get,
				ActionName:        "查询",
				IsRelatedResource: false,
			},
		},
	},
	{
		ResourceTypeID:       SysOperationStatistic,
		ResourceTypeName:     "运营统计",
		ParentResourceTypeID: "",
		Share:                false,
		Actions: []Action{
			{
				ActionID:          Edit,
				ActionName:        "编辑",
				IsRelatedResource: false,
			},
			{
				ActionID:          Get,
				ActionName:        "查询",
				IsRelatedResource: false,
			},
		},
	},
}
