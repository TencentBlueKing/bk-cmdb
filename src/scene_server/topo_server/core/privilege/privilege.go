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
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/topo_server/core/types"
)

// PermissionInterface the permission methods
type PermissionInterface interface {
	SetUserGroupPermission(supplierAccount, gourpID string, permission *metadata.PrivilegeUserGroup) error
	GetUserGroupPermission(supplierAccount, groupID string) (*metadata.GroupPrivilege, error)
	GetUserPermission(supplierAccount, userName string) (*metadata.Gprivilege, error)
}

// NewPermission create a new permission instance
func NewPermission(params types.ContextParams, client apimachinery.ClientSetInterface) PermissionInterface {

	return &userGroupPermission{
		params: params,
		client: client,
	}
}

// userGroupPermission the permission user group definitions
type userGroupPermission struct {
	params types.ContextParams
	client apimachinery.ClientSetInterface

	permission metadata.PrivilegeUserGroup
}

// MarshalJSON marshal the data into json
func (u *userGroupPermission) MarshalJSON() ([]byte, error) {
	return json.Marshal(u.permission)
}

func (u *userGroupPermission) SetUserGroupPermission(supplierAccount, groupID string, permission *metadata.PrivilegeUserGroup) error {

	rsp, err := u.client.ObjectController().Privilege().GetUserGroupPrivi(context.Background(), supplierAccount, groupID, u.params.Header)
	if nil != err {
		blog.Errorf("[privilege] failed to request object controller, error info is %s", err.Error())
		return u.params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	// create a new privilege
	if !rsp.Result {
		blog.Errorf("[privilege] failed to get user group privi, error info is %s", rsp.ErrMsg)
		return u.params.Err.New(rsp.Code, rsp.ErrMsg)
	}

	if nil == rsp.Data.Privilege || (0 == len(rsp.Data.Privilege.ModelConfig) && nil == rsp.Data.Privilege.SysConfig) {
		rsp, err := u.client.ObjectController().Privilege().CreateUserGroupPrivi(context.Background(), supplierAccount, groupID, u.params.Header, permission)
		if nil != err {
			blog.Errorf("[privilege] failed to request object controller, error info is %s", err.Error())
			return u.params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
		}
		if !rsp.Result {
			return u.params.Err.New(rsp.Code, rsp.ErrMsg)
		}
		return nil
	}

	// update privilege
	rspUpdate, err := u.client.ObjectController().Privilege().UpdateUserGroupPrivi(context.Background(), supplierAccount, groupID, u.params.Header, permission)
	if nil != err {
		blog.Errorf("[privilege] failed to request object controller, error info is %s", err.Error())
		return u.params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if !rspUpdate.Result {
		return u.params.Err.New(rspUpdate.Code, rspUpdate.ErrMsg)
	}

	return nil
}
func (u *userGroupPermission) GetUserGroupPermission(supplierAccount, groupID string) (*metadata.GroupPrivilege, error) {

	rsp, err := u.client.ObjectController().Privilege().GetUserGroupPrivi(context.Background(), supplierAccount, groupID, u.params.Header)
	if nil != err {
		blog.Errorf("[privilege] failed to request object controller, error info is %s", err.Error())
		return nil, u.params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if !rsp.Result {
		return nil, u.params.Err.New(rsp.Code, rsp.ErrMsg)
	}

	return &rsp.Data, nil
}

func (u *userGroupPermission) GetUserPermission(supplierAccount, userName string) (*metadata.Gprivilege, error) {

	gPrivilege := metadata.Gprivilege{
		IsHostCrossBiz: false,
		ModelConfig:    map[string]map[string][]string{},
	}

	// get cross biz permission
	rsp, err := u.client.ObjectController().Privilege().GetSystemFlag(context.Background(), supplierAccount, common.HostCrossBizField, u.params.Header)
	if nil != err {
		blog.Errorf("[privilege] failed to request object controller, error info is %s", err.Error())
		//		return nil, u.params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	} else {
		gPrivilege.IsHostCrossBiz = rsp.Result
	}

	// search user group permission
	cond := condition.CreateCondition()
	cond.Field(common.BKUserListField).Like(userName)
	rspSearchGroup, err := u.client.ObjectController().Privilege().SearchUserGroup(context.Background(), supplierAccount, u.params.Header, cond.ToMapStr())
	if nil != err {
		blog.Errorf("[privilege] failed to request object controller, error info is %s", err.Error())
		return nil, u.params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if !rspSearchGroup.Result {
		blog.Errorf("[privilege] failed to search group permission, error info is %s", rspSearchGroup.ErrMsg)
		return nil, u.params.Err.New(rspSearchGroup.Code, rspSearchGroup.ErrMsg)
	}

	var gglconfig []string
	var gbkconfig []string
	var modelCls []string
	modelPrivi := make(map[string][]string)
	modelClsConfig := make(map[string]string)
	// construct the result
	for _, item := range rspSearchGroup.Data {
		//item.GroupID

		grpPrivilege, err := u.client.ObjectController().Privilege().GetUserGroupPrivi(context.Background(), supplierAccount, item.GroupID, u.params.Header)
		if nil != err {
			blog.Errorf("[privilege] failed to get the user group(%s) privilege, error info is %s", item.GroupID, err.Error())
			continue
		}

		if !grpPrivilege.Result {
			blog.Errorf("[privige] failed to search the user group, error info is %s", grpPrivilege.ErrMsg)
			continue
		}

		if nil == grpPrivilege.Data.Privilege {
			grpPrivilege.Data.Privilege = &metadata.Privilege{}
			continue
		}

		if nil != grpPrivilege.Data.Privilege.SysConfig {
			sysConfig := *grpPrivilege.Data.Privilege.SysConfig
			for _, i := range sysConfig.Globalbusi {
				gglconfig = append(gglconfig, i)
			}
			for _, j := range sysConfig.BackConfig {
				gbkconfig = append(gbkconfig, j)
			}
		}

		for key, val := range grpPrivilege.Data.Privilege.ModelConfig {
			for subKey, subVal := range val {
				for _, data := range subVal {
					modelPrivi[subKey] = append(modelPrivi[subKey], data)
				}
				modelClsConfig[subKey] = key
			}
			modelCls = append(modelCls, key)

		}

	} // end for

	umodelCls := util.RemoveDuplicatesAndEmpty(modelCls)
	cls := make(map[string]map[string][]string)
	for _, i := range umodelCls {
		modelCls := make(map[string][]string)
		for j, k := range modelPrivi {
			if modelClsConfig[j] == i {
				modelCls[j] = k
			}
		}
		cls[i] = modelCls
	}

	gPrivilege.SysConfig.BackConfig = util.RemoveDuplicatesAndEmpty(gbkconfig)
	gPrivilege.SysConfig.Globalbusi = util.RemoveDuplicatesAndEmpty(gglconfig)
	gPrivilege.ModelConfig = cls
	return &gPrivilege, nil

}
