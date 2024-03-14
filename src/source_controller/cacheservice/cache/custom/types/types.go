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

// Package types defines custom resource cache types
package types

// ListPodLabelKeyOption list pod label key options
type ListPodLabelKeyOption struct {
	BizID int64 `json:"bk_biz_id"`
}

// ListPodLabelValueOption list pod label value options
type ListPodLabelValueOption struct {
	BizID int64  `json:"bk_biz_id"`
	Key   string `json:"key"`
}

// RefreshPodLabelOption refresh pod label key and value options
type RefreshPodLabelOption struct {
	BizID int64 `json:"bk_biz_id"`
}

// ResType is the custom resource type
type ResType string

const (
	// PodLabelKeyType is the kube pod label key type
	PodLabelKeyType ResType = "pod_label_key"
	// PodLabelValueType is the kube pod label value type
	PodLabelValueType ResType = "pod_label_value"
	// SharedNsAsstBizType is the shared namespace to associated biz type
	SharedNsAsstBizType ResType = "shared_ns_asst_biz"
)

const (
	// RedisPage is the default redis operation page size
	RedisPage = 500
	// DBPage is the default db operation page size
	DBPage = 200
)
