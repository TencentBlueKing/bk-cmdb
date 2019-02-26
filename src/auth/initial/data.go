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
	"configcenter/src/auth/authcenter"
)

var expectSystem = authcenter.System{
	SystemID:           "cmdb",
	SystemName:         "蓝鲸智云配置平台",
	Desc:               "蓝鲸配置平台（CMDB）",
	ReleatedScopeTypes: "biz;system",
	Managers:           "system",
	Creator:            "system",
	Updater:            "system",
}

var expectResource = []authcenter.Resource{
	{
		ResourceTypeID:       "",
		ResourceTypeName:     "",
		ParentResourceTypeID: "",
		Actions: []authcenter.Action{
			{
				ActionID:          "",
				ActionName:        "",
				IsRelatedResource: false,
			},
		},
	},
}
