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

// kubeRelated generate kube related resource auth parse stream
func (ps *parseStream) kubeRelated() *parseStream {
	if ps.shouldReturn() {
		return ps
	}

	ps.KubePod()

	return ps
}

// KubePodConfigs kube pod auth parser configs
var KubePodConfigs = []AuthConfig{
	{
		Name:           "batchDeletePod",
		Description:    "批量删除Pod",
		Pattern:        "/api/v3/deletemany/kube/pod",
		HTTPMethod:     http.MethodDelete,
		ResourceType:   meta.KubePod,
		ResourceAction: meta.Delete,
	},
}

// KubePod generate kube pod auth parse stream
func (ps *parseStream) KubePod() *parseStream {
	return ParseStreamWithFramework(ps, KubePodConfigs)
}
