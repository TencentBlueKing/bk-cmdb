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

	"configcenter/src/apimachinery"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	//"configcenter/src/common/metadata"
	"configcenter/src/scene_server/topo_server/core/types"
)

// RolePermission role permission
type RolePermission interface {
	CreatePermission(supplierAccount, objID, propertyID string, data []string) error
	GetPermission(supplierAccount, objID, propertyID string) (interface{}, error)
}

// NewRole create a new role instance
func NewRole(params types.ContextParams, client apimachinery.ClientSetInterface) RolePermission {
	return &rolePermission{
		params: params,
		client: client,
	}
}

type rolePermission struct {
	params types.ContextParams
	client apimachinery.ClientSetInterface
}

func (r *rolePermission) CreatePermission(supplierAccount, objID, propertyID string, data []string) error {

	rsp, err := r.client.ObjectController().Privilege().GetRolePri(context.Background(), supplierAccount, objID, propertyID, r.params.Header)
	if nil != err {
		blog.Errorf("[permission] failed to request object controller, error info is %s", err.Error())
		return r.params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if rsp.Result {
		shouldCreate := false
		switch tmpData := rsp.Data.(type) {
		case nil:
		case map[string]interface{}:
			if 0 == len(tmpData) {
				shouldCreate = true
			}
		}

		if !shouldCreate {
			rsp, err := r.client.ObjectController().Privilege().UpdateRolePri(context.Background(), supplierAccount, objID, propertyID, r.params.Header, data)
			if nil != err {
				blog.Errorf("[permission] failed to request object controller, error info is %s", err.Error())
				return r.params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
			}

			if !rsp.Result {
				blog.Errorf("[permission] failed to update the role, error info is %s", rsp.ErrMsg)
				return r.params.Err.New(rsp.Code, rsp.ErrMsg)
			}

			return nil
		}
	}

	rsp, err = r.client.ObjectController().Privilege().CreateRolePri(context.Background(), supplierAccount, objID, propertyID, r.params.Header, data)

	if nil != err {
		blog.Errorf("[permission] failed to request object controller, error info is %s", err.Error())
		return r.params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if !rsp.Result {
		blog.Errorf("[permission] failed to create the role, error info is %s", rsp.ErrMsg)
		return r.params.Err.New(rsp.Code, rsp.ErrMsg)
	}

	return nil
}

func (r *rolePermission) GetPermission(supplierAccount, objID, propertyID string) (interface{}, error) {

	rsp, err := r.client.ObjectController().Privilege().GetRolePri(context.Background(), supplierAccount, objID, propertyID, r.params.Header)
	if nil != err {
		blog.Errorf("[permission] failed to request object controller, error info is %s", err.Error())
		return nil, r.params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if !rsp.Result {
		blog.Errorf("[permission] failed to search the role permission, error info is %s", rsp.ErrMsg)
		return nil, r.params.Err.New(rsp.Code, rsp.ErrMsg)
	}

	return rsp.Data, nil
}
