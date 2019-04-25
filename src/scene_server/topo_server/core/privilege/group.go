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
	"encoding/json"

	"configcenter/src/apimachinery"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/topo_server/core/types"
)

// UserGroupInterface the permission user groups methods
type UserGroupInterface interface {
	CreateUserGroup(supplierAccount string, userGroup *metadata.UserGroup) error
	DeleteUserGroup(supplierAccount, groupID string) error
	UpdateUserGroup(supplierAccount, groupID string, data mapstr.MapStr) error
	SearchUserGroup(supplierAccount string, cond condition.Condition) ([]metadata.UserGroup, error)
}

// NewUserGroup create a user group instance
func NewUserGroup(params types.ContextParams, client apimachinery.ClientSetInterface) UserGroupInterface {
	return &userGroup{
		params: params,
		client: client,
	}
}

// userGroup the permission user group definitions
type userGroup struct {
	params    types.ContextParams
	client    apimachinery.ClientSetInterface
	userGroup metadata.UserGroup
}

// MarshalJSON marshal the data into json
func (u *userGroup) MarshalJSON() ([]byte, error) {
	return json.Marshal(u.userGroup)
}

func (u *userGroup) checkGroupNameRepeat(supplierAccount, groupID, groupName string) error {

	if 0 == len(groupName) {
		return u.params.Err.Errorf(common.CCErrCommParamsNeedSet, "group name")
	}

	cond := condition.CreateCondition()
	if 0 != len(groupID) {
		cond.Field("group_id").NotIn([]string{groupID})
	}
	cond.Field("group_name").Eq(groupName)

	rsp, err := u.client.ObjectController().Privilege().SearchUserGroup(context.Background(), supplierAccount, u.params.Header, cond.ToMapStr())
	if nil != err {
		blog.Errorf("[permission] failed to request object controller, error info is %s", err.Error())
		return u.params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if !rsp.Result {
		blog.Errorf("[permission] failed to check the user group, error info is %s", rsp.ErrMsg)
		return u.params.Err.New(rsp.Code, rsp.ErrMsg)
	}

	if 0 < len(rsp.Data) {
		blog.Warnf("[permission] the group name (%s) repeated", groupName)
		return u.params.Err.Error(common.CCErrCommDuplicateItem)
	}

	return nil
}

func (u *userGroup) CreateUserGroup(supplierAccount string, userGroup *metadata.UserGroup) error {

	if err := u.checkGroupNameRepeat(supplierAccount, "", userGroup.GroupName); nil != err {
		return err
	}

	rspCreate, err := u.client.ObjectController().Privilege().CreateUserGroup(context.Background(), supplierAccount, u.params.Header, userGroup.ToMapStr())
	if nil != err {
		blog.Errorf("[permission] failed to request object controller, error info is %s", err.Error())
		return u.params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if !rspCreate.Result {
		blog.Errorf("[permission] failed to create the user group name(%s), error info is %s", userGroup.GroupName, rspCreate.ErrMsg)
		return u.params.Err.New(rspCreate.Code, rspCreate.ErrMsg)
	}

	return nil
}

func (u *userGroup) DeleteUserGroup(supplierAccount, groupID string) error {

	rsp, err := u.client.ObjectController().Privilege().DeleteUserGroup(context.Background(), supplierAccount, groupID, u.params.Header)
	if nil != err {
		blog.Errorf("[permission] failed to request object controller, error info is %s", err.Error())
		return u.params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if !rsp.Result {
		blog.Errorf("[permission] failed to delete the group (%s), error info is %s", groupID, rsp.ErrMsg)
		return u.params.Err.New(rsp.Code, rsp.ErrMsg)
	}

	return nil
}

func (u *userGroup) UpdateUserGroup(supplierAccount, groupID string, data mapstr.MapStr) error {

	groupName, err := data.String("group_name")
	if nil != err {
		blog.Errorf("the group name (%#v) is invalid, error info is %s", data, err.Error())
		return u.params.Err.Errorf(common.CCErrCommParamsNeedSet, "group name")
	}

	if err := u.checkGroupNameRepeat(supplierAccount, groupID, groupName); nil != err {
		return err
	}

	rsp, err := u.client.ObjectController().Privilege().UpdateUserGroup(context.Background(), supplierAccount, groupID, u.params.Header, data)
	if nil != err {
		blog.Errorf("[permission] failed to request object controller, error info is %s", err.Error())
		return u.params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if !rsp.Result {
		blog.Errorf("[permission] failed to update the group (%s), error info is %s", groupID, rsp.ErrMsg)
		return u.params.Err.New(rsp.Code, rsp.ErrMsg)
	}

	return nil
}

func (u *userGroup) SearchUserGroup(supplierAccount string, cond condition.Condition) ([]metadata.UserGroup, error) {

	rsp, err := u.client.ObjectController().Privilege().SearchUserGroup(context.Background(), supplierAccount, u.params.Header, cond.ToMapStr())
	if nil != err {
		blog.Errorf("[permission] failed to request object controller, error info is %s", err.Error())
		return nil, u.params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if !rsp.Result {
		blog.Errorf("[permission] failed to search, the condition (%#v), error info is %s", cond.ToMapStr(), rsp.ErrMsg)
		return nil, u.params.Err.New(rsp.Code, rsp.ErrMsg)
	}

	return rsp.Data, nil
}
