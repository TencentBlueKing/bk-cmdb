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

package initial

import (
	"strings"

	"configcenter/src/auth/authcenter"
	"configcenter/src/auth/meta"
)

var expectSystem = authcenter.System{
	SystemID:           meta.SystemIDCMDB,
	SystemName:         meta.SystemNameCMDB,
	Desc:               "蓝鲸配置平台（CMDB）",
	ReleatedScopeTypes: strings.Join([]string{meta.ScopeTypeIDBiz, meta.ScopeTypeIDSystem}, ";"),
	Managers:           "system",
	Creator:            "system",
	Updater:            "system",
}

var expectSystemResourceType = []authcenter.ResourceType{
	{
		ResourceTypeID:       meta.ModelClassification.String(),
		ResourceTypeName:     "模型分组",
		ParentResourceTypeID: "",
		Actions: []authcenter.Action{
			{
				ActionID:          authcenter.Create,
				ActionName:        "新建",
				IsRelatedResource: false,
			},
			{
				ActionID:          authcenter.Edit,
				ActionName:        "编辑",
				IsRelatedResource: true,
			},
			{
				ActionID:          authcenter.Delete,
				ActionName:        "删除",
				IsRelatedResource: true,
			},
		},
	},
	{
		ResourceTypeID:       meta.Model.String(),
		ResourceTypeName:     "模型",
		ParentResourceTypeID: meta.ModelClassification.String(),
		Actions: []authcenter.Action{
			{
				ActionID:          authcenter.Create,
				ActionName:        "新建",
				IsRelatedResource: false,
			},
			{
				ActionID:          authcenter.Edit,
				ActionName:        "编辑",
				IsRelatedResource: true,
			},
			{
				ActionID:          authcenter.Delete,
				ActionName:        "删除",
				IsRelatedResource: true,
			},
		},
	},
	{
		ResourceTypeID:       meta.ModelInstance.String(),
		ResourceTypeName:     "实例",
		ParentResourceTypeID: meta.Model.String(),
		Actions: []authcenter.Action{
			{
				ActionID:          authcenter.Create,
				ActionName:        "新建",
				IsRelatedResource: false,
			},
			{
				ActionID:          authcenter.Edit,
				ActionName:        "编辑",
				IsRelatedResource: true,
			},
			{
				ActionID:          authcenter.Delete,
				ActionName:        "删除",
				IsRelatedResource: true,
			},
			{
				ActionID:          authcenter.Get,
				ActionName:        "查询",
				IsRelatedResource: true,
			},
		},
	},
	{
		ResourceTypeID:       meta.AssociationType.String(),
		ResourceTypeName:     "关联类型",
		ParentResourceTypeID: "",
		Actions: []authcenter.Action{
			{
				ActionID:          authcenter.Create,
				ActionName:        "新建",
				IsRelatedResource: false,
			},
			{
				ActionID:          authcenter.Edit,
				ActionName:        "编辑",
				IsRelatedResource: true,
			},
			{
				ActionID:          authcenter.Delete,
				ActionName:        "删除",
				IsRelatedResource: true,
			},
		},
	},
	{
		ResourceTypeID:       meta.Business.String(),
		ResourceTypeName:     "业务",
		ParentResourceTypeID: "",
		Actions: []authcenter.Action{
			{
				ActionID:          authcenter.Create,
				ActionName:        "新建",
				IsRelatedResource: false,
			},
			{
				ActionID:          authcenter.Edit,
				ActionName:        "编辑",
				IsRelatedResource: true,
			},
			{
				ActionID:          authcenter.Archive,
				ActionName:        "归档",
				IsRelatedResource: true,
			},
			{
				ActionID:          authcenter.Get,
				ActionName:        "查询",
				IsRelatedResource: true,
			},
		},
	},
	{
		ResourceTypeID:       meta.Host,
		ResourceTypeName:     "主机",
		ParentResourceTypeID: "",
		Actions: []authcenter.Action{
			{
				ActionID:          authcenter.Create,
				ActionName:        "新建",
				IsRelatedResource: false,
			},
			{
				ActionID:          authcenter.Edit,
				ActionName:        "编辑",
				IsRelatedResource: true,
			},
			{
				ActionID:          authcenter.Delete,
				ActionName:        "编辑",
				IsRelatedResource: true,
			},
			{
				ActionID:          authcenter.Get,
				ActionName:        "查询",
				IsRelatedResource: true,
			},
			{
				ActionID:          authcenter.ModuleTransfer,
				ActionName:        "分配到业务",
				IsRelatedResource: true,
			},
		},
	},
	{
		ResourceTypeID:       meta.EventPushing,
		ResourceTypeName:     "事件推送",
		ParentResourceTypeID: "",
		Actions: []authcenter.Action{
			{
				ActionID:          authcenter.Create,
				ActionName:        "新建",
				IsRelatedResource: false,
			},
			{
				ActionID:          authcenter.Edit,
				ActionName:        "编辑",
				IsRelatedResource: true,
			},
			{
				ActionID:          authcenter.Delete,
				ActionName:        "编辑",
				IsRelatedResource: true,
			},
			{
				ActionID:          authcenter.Get,
				ActionName:        "查询",
				IsRelatedResource: true,
			},
		},
	},
	{
		ResourceTypeID:       meta.SystemFunctionality,
		ResourceTypeName:     "系统基础",
		ParentResourceTypeID: "",
		Actions: []authcenter.Action{
			{
				ActionID:          authcenter.TopoLayerManage,
				ActionName:        "拓扑层级管理",
				IsRelatedResource: false,
			},
			{
				ActionID:          authcenter.AdminEntrance,
				ActionName:        "管理页面入口",
				IsRelatedResource: false,
			},
		},
	},
}
