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

var expectBizResourceType = []ResourceType{
	// {
	// 	ResourceTypeID:       BizModelGroup,
	// 	ResourceTypeName:     "模型分组",
	// 	ParentResourceTypeID: "",
	// 	Share:                true,
	// 	Actions: []Action{
	// 		{
	// 			ActionID:          Create,
	// 			ActionName:        "新建",
	// 			IsRelatedResource: false,
	// 		},
	// 		{
	// 			ActionID:          Edit,
	// 			ActionName:        "编辑",
	// 			IsRelatedResource: true,
	// 		},
	// 		{
	// 			ActionID:          Delete,
	// 			ActionName:        "删除",
	// 			IsRelatedResource: true,
	// 		},
	// 	},
	// },
	{
		ResourceTypeID:   BizModel,
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
	// {
	// 	ResourceTypeID:       BizInstance,
	// 	ResourceTypeName:     "实例",
	// 	ParentResourceTypeID: BizModel,
	// 	Share:                true,
	// 	Actions: []Action{
	// 		{
	// 			ActionID:          Create,
	// 			ActionName:        "新建",
	// 			IsRelatedResource: false,
	// 		},
	// 		{
	// 			ActionID:          Edit,
	// 			ActionName:        "编辑",
	// 			IsRelatedResource: true,
	// 		},
	// 		{
	// 			ActionID:          Delete,
	// 			ActionName:        "删除",
	// 			IsRelatedResource: true,
	// 		},
	// 		{
	// 			ActionID:          Get,
	// 			ActionName:        "查询",
	// 			IsRelatedResource: true,
	// 		},
	// 	},
	// },
	{
		ResourceTypeID:       BizCustomQuery,
		ResourceTypeName:     "动态分组",
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
			{
				ActionID:          Get,
				ActionName:        "查询",
				IsRelatedResource: true,
			},
		},
	},
	{
		ResourceTypeID:       BizHostInstance,
		ResourceTypeName:     "主机",
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
				ActionName:        "删除/归还",
				IsRelatedResource: true,
			},
		},
	},
	{
		ResourceTypeID:       BizTopology,
		ResourceTypeName:     "业务拓扑",
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
				IsRelatedResource: false,
			},
			{
				ActionID:          Delete,
				ActionName:        "删除",
				IsRelatedResource: false,
			},
		},
	},
	{
		ResourceTypeID:       BizProcessServiceTemplate,
		ResourceTypeName:     "服务模板",
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
		ResourceTypeID:       BizProcessServiceCategory,
		ResourceTypeName:     "服务分类",
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
				IsRelatedResource: false,
			},
			{
				ActionID:          Delete,
				ActionName:        "删除",
				IsRelatedResource: false,
			},
		},
	},
	{
		ResourceTypeID:       BizProcessServiceInstance,
		ResourceTypeName:     "服务实例",
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
				IsRelatedResource: false,
			},
			{
				ActionID:          Delete,
				ActionName:        "删除",
				IsRelatedResource: false,
			},
		},
	},
	// {
	// 	ResourceTypeID:       BizAuditLog,
	// 	ResourceTypeName:     "操作审计",
	// 	ParentResourceTypeID: "",
	// 	Share:                false,
	// 	Actions: []Action{
	// 		{
	// 			ActionID:          Get,
	// 			ActionName:        "查询",
	// 			IsRelatedResource: false,
	// 		},
	// 	},
	// },
	{
		ResourceTypeID:       BizSetTemplate,
		ResourceTypeName:     "集群模板",
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
}
