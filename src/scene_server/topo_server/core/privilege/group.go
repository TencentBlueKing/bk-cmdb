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

package privilege

import (
	"encoding/json"

	types "configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
)

// UserGroupInterface the permission user groups methods
type UserGroupInterface interface {
	CreateUserGroup(supplierAccount string, userGroup *metadata.UserGroup) error
	DeleteUserGroup(supplierAccount, groupID string) error
	UpdateUserGroup(supplierAccount, groupID string, data types.MapStr) error
	SearchUserGroup(supplierAccount string) ([]*metadata.UserGroup, error)
}

// userGroup the permission user group definitions
type userGroup struct {
	userGroup metadata.UserGroup
}

// MarshalJSON marshal the data into json
func (u *userGroup) MarshalJSON() ([]byte, error) {
	return json.Marshal(u.userGroup)
}

func (u *userGroup) CreateUserGroup(supplierAccount string, userGroup *metadata.UserGroup) error {
	return nil
}

func (u *userGroup) DeleteUserGroup(supplierAccount, groupID string) error {
	return nil
}

func (u *userGroup) UpdateUserGroup(supplierAccount, groupID string, data types.MapStr) error {
	return nil
}

func (u *userGroup) SearchUserGroup(supplierAccount string) ([]*metadata.UserGroup, error) {
	return nil, nil
}
