/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 THL A29 Limited,
 * a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 * We undertake not to change the open source license (MIT license) applicable
 * to the current version of the project delivered to anyone in the future.
 */

package parser

import (
	"net/http"

	"configcenter/src/ac/meta"
)

// ModelQuoteAuthConfigs model quote related auth configs, skip all, authorize in topo-server.
var ModelQuoteAuthConfigs = []AuthConfig{
	{
		Name:           "BatchCreateQuotedInstance",
		Description:    "创建引用模型实例",
		Pattern:        "/api/v3/createmany/quoted/instance",
		HTTPMethod:     http.MethodPost,
		ResourceAction: meta.SkipAction,
	},
	{
		Name:           "ListQuotedInstance",
		Description:    "查询引用模型实例列表",
		Pattern:        "/api/v3/findmany/quoted/instance",
		HTTPMethod:     http.MethodPost,
		ResourceAction: meta.SkipAction,
	},
	{
		Name:           "BatchUpdateQuotedInstance",
		Description:    "更新引用模型实例",
		Pattern:        "/api/v3/updatemany/quoted/instance",
		HTTPMethod:     http.MethodPut,
		ResourceAction: meta.SkipAction,
	},
	{
		Name:           "BatchDeleteQuotedInstance",
		Description:    "删除引用模型实例",
		Pattern:        "/api/v3/deletemany/quoted/instance",
		HTTPMethod:     http.MethodDelete,
		ResourceAction: meta.SkipAction,
	},
}

func (ps *parseStream) modelQuote() *parseStream {
	return ParseStreamWithFramework(ps, ModelQuoteAuthConfigs)
}
