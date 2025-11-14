/*
 * TencentBlueKing is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - CMDB) available.
 * Copyright (C) 2025 Tencent. All rights reserved.
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

package authpb

import "github.com/TencentBlueKing/bk-cmdb/pkg/auth/meta"

// ConvertToPbBasic converts meta.Basic to Basic
func ConvertToPbBasic(basic *meta.Basic) *Basic {
	return &Basic{
		Type:   string(basic.Type),
		Action: string(basic.Action),
		Name:   basic.Name,
		Id:     basic.ID,
	}
}

// ConvertToPBAuthAttr converts meta.ResourceAttribute to ResourceAttribute
func ConvertToPBAuthAttr(attr *meta.ResourceAttribute) *ResourceAttribute {
	pbAttr := &ResourceAttribute{
		Basic:  ConvertToPbBasic(attr.Basic),
		Layers: make([]*Basic, len(attr.Layers)),
	}

	for i, layer := range attr.Layers {
		pbAttr.Layers[i] = ConvertToPbBasic(&layer)
	}

	return pbAttr
}
