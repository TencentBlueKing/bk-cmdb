/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 Tencent. All rights reserved.
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

package operator

import "strings"

const (
	// IamIDKey is the authorized resource's instance id in the authorization plan condition.
	IamIDKey = "id"
	// IamAttrPrefix is the prefix of the authorized resource's own attributes.
	// It avoids colliding with IamIDKey and ancestor resource type ids (e.g. "biz").
	IamAttrPrefix = "attr."
)

// Field is a compare node's field in the authorization plan condition.
// There are three kinds:
//  1. IamIDKey: the authorized resource's instance id
//  2. ancestor: the field is the ancestor resource type id, the value is the ancestor instance ids with 'in' operator
//  3. the authorized resource's own attribute, which must be prefixed with IamAttrPrefix
type Field string

// String returns the raw field name.
func (f Field) String() string {
	return string(f)
}

// IsID returns whether the field is the authorized resource's instance id.
func (f Field) IsID() bool {
	return string(f) == IamIDKey
}

// IsSelfAttr returns whether the field is the authorized resource's own attribute.
func (f Field) IsSelfAttr() bool {
	return strings.HasPrefix(string(f), IamAttrPrefix)
}

// SelfAttr returns the resource's own attribute name without IamAttrPrefix.
// Empty string means the field is not a self attribute, or the name after prefix is empty.
func (f Field) SelfAttr() string {
	if !f.IsSelfAttr() {
		return ""
	}
	return strings.TrimPrefix(string(f), IamAttrPrefix)
}

// IsAncestor returns whether the field is an ancestor resource type id.
func (f Field) IsAncestor() bool {
	return !f.IsID() && !f.IsSelfAttr() && f != ""
}

// AttrField builds the condition field of the authorized resource's own attribute.
func AttrField(attr string) Field {
	return Field(IamAttrPrefix + attr)
}
