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

package parser

import (
	"net/http"
	"regexp"

	"configcenter/src/ac/meta"
)

/*
 http.MethodPost,  "/create/operation/chart"
 http.MethodDelete,  "/delete/operation/chart/{id}"
 http.MethodPost,  "/update/operation/chart"
 http.MethodGet,  "/search/operation/chart"
 http.MethodPost,  "/search/operation/chart/data"
*/
var OperationStatisticAuthConfigs = []AuthConfig{
	{
		Name:           "CreateOperationStatisticRegex",
		Description:    "创建运营统计",
		Regex:          regexp.MustCompile(`^/api/v3/create/operation/chart/?$`),
		HTTPMethod:     http.MethodPost,
		BizIDGetter:    nil,
		ResourceType:   meta.OperationStatistic,
		ResourceAction: meta.Update,
	},
	{
		Name:           "DeleteOperationStatisticRegex",
		Description:    "删除运营统计",
		Regex:          regexp.MustCompile(`^/api/v3/delete/operation/chart/([0-9]+)/?$`),
		HTTPMethod:     http.MethodDelete,
		BizIDGetter:    nil,
		ResourceType:   meta.OperationStatistic,
		ResourceAction: meta.Update,
	},
	{
		Name:           "UpdateOperationStatisticRegex",
		Description:    "更新运营统计",
		Regex:          regexp.MustCompile(`^/api/v3/update/operation/chart/?$`),
		HTTPMethod:     http.MethodPost,
		BizIDGetter:    nil,
		ResourceType:   meta.OperationStatistic,
		ResourceAction: meta.Update,
	},
	{
		Name:           "SearchOperationStatisticChartRegex",
		Description:    "查看运营统计图表配置",
		Regex:          regexp.MustCompile(`^/api/v3/findmany/operation/chart/?$`),
		HTTPMethod:     http.MethodGet,
		BizIDGetter:    nil,
		ResourceType:   meta.OperationStatistic,
		ResourceAction: meta.Find,
	},
	{
		Name:           "SearchOperationStatisticDataRegex",
		Description:    "查看运营统计数据",
		Regex:          regexp.MustCompile(`^/api/v3/find/operation/chart/data/?$`),
		HTTPMethod:     http.MethodPost,
		BizIDGetter:    nil,
		ResourceType:   meta.OperationStatistic,
		ResourceAction: meta.Find,
	},
	{
		Name:           "UpdateOperationStatisticPositionRegex",
		Description:    "更新运营统计图表位置",
		Regex:          regexp.MustCompile(`^/api/v3/update/operation/chart/position/?$`),
		HTTPMethod:     http.MethodPost,
		BizIDGetter:    nil,
		ResourceType:   meta.OperationStatistic,
		ResourceAction: meta.Update,
	},
}

func (ps *parseStream) OperationStatistic() *parseStream {
	return ParseStreamWithFramework(ps, OperationStatisticAuthConfigs)
}
