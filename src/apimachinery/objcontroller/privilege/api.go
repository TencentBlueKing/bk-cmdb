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
	"context"

	"configcenter/src/apimachinery/rest"
	"configcenter/src/apimachinery/util"
	"configcenter/src/common/core/cc/api"
)

type PrivilegeInterface interface {
	CreateUserGroupPrivi(ctx context.Context, groupID string, h util.Headers, dat interface{}) (resp *api.BKAPIRsp, err error)
	UpdateUserGroupPrivi(ctx context.Context, groupID string, h util.Headers, dat interface{}) (resp *api.BKAPIRsp, err error)
	GetUserGroupPrivi(ctx context.Context, groupID string, h util.Headers) (resp *api.BKAPIRsp, err error)
	CreateRolePri(ctx context.Context, objID string, propertyID string, h util.Headers, role []string) (resp *api.BKAPIRsp, err error)
	GetRolePri(ctx context.Context, objID string, propertyID string, h util.Headers) (resp *api.BKAPIRsp, err error)
	UpdateRolePri(ctx context.Context, objID string, propertyID string, h util.Headers, role []string) (resp *api.BKAPIRsp, err error)
	GetSystemFlag(ctx context.Context, flag string, h util.Headers) (resp *api.BKAPIRsp, err error)
}

func NewPrivilegeInterface(client rest.ClientInterface) PrivilegeInterface {
	return &privilege{client: client}
}

type privilege struct {
	client rest.ClientInterface
}
