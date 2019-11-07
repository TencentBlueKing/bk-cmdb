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

package host

import (
	"context"
	"fmt"
	"net/http"

	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

// TransferToInnerModule  transfer host to inner module  eg:idle module and fault module
func (h *host) TransferToInnerModule(ctx context.Context, header http.Header, input *metadata.TransferHostToInnerModule) (resp *metadata.OperaterException, err error) {
	resp = new(metadata.OperaterException)
	subPath := "/set/module/host/relation/inner/module"

	err = h.client.Post().
		WithContext(ctx).
		Body(input).
		SubResource(subPath).
		WithHeaders(header).
		Do().
		Into(resp)
	return
}

// TransferHostModule  transfer host to  module
func (h *host) TransferToNormalModule(ctx context.Context, header http.Header, input *metadata.HostsModuleRelation) (resp *metadata.OperaterException, err error) {
	resp = new(metadata.OperaterException)
	subPath := "/set/module/host/relation/module"

	err = h.client.Post().
		WithContext(ctx).
		Body(input).
		SubResource(subPath).
		WithHeaders(header).
		Do().
		Into(resp)
	return
}

// RemoveFromModule 将主机从模块中移出
// 如果主机属于n+1个模块（n>0），操作之后，主机属于n个模块
// 如果主机属于1个模块, 且非空闲机模块，操作之后，主机属于空闲机模块
// 如果主机属于空闲机模块，操作失败
// 如果主机属于故障机模块，操作失败
// 如果主机不在参数指定的模块中，操作失败
func (h *host) RemoveFromModule(ctx context.Context, header http.Header, input *metadata.RemoveHostsFromModuleOption) (resp *metadata.OperaterException, err error) {
	resp = new(metadata.OperaterException)
	subPath := "/delete/host/host_module_relations"

	err = h.client.Delete().
		WithContext(ctx).
		Body(input).
		SubResource(subPath).
		WithHeaders(header).
		Do().
		Into(resp)
	return
}

// TransferHostCrossBusiness  transfer host to other business module
func (h *host) TransferToAnotherBusiness(ctx context.Context, header http.Header, input *metadata.TransferHostsCrossBusinessRequest) (resp *metadata.OperaterException, err error) {
	resp = new(metadata.OperaterException)
	subPath := "/set/module/host/relation/cross/business"

	err = h.client.Post().
		WithContext(ctx).
		Body(input).
		SubResource(subPath).
		WithHeaders(header).
		Do().
		Into(resp)
	return
}

// DeleteHost delete host
func (h *host) DeleteHostFromSystem(ctx context.Context, header http.Header, input *metadata.DeleteHostRequest) (resp *metadata.OperaterException, err error) {
	resp = new(metadata.OperaterException)
	subPath := "/delete/host"

	err = h.client.Delete().
		WithContext(ctx).
		Body(input).
		SubResource(subPath).
		WithHeaders(header).
		Do().
		Into(resp)
	return
}

// GetHostModuleRelation get host module relation
func (h *host) GetHostModuleRelation(ctx context.Context, header http.Header, input *metadata.HostModuleRelationRequest) (resp *metadata.HostConfig, err error) {
	resp = new(metadata.HostConfig)
	subPath := "/read/module/host/relation"

	err = h.client.Post().
		WithContext(ctx).
		Body(input).
		SubResource(subPath).
		WithHeaders(header).
		Do().
		Into(resp)
	return
}

// FindIdentifier  query host identifier
func (h *host) FindIdentifier(ctx context.Context, header http.Header, input *metadata.SearchHostIdentifierParam) (resp *metadata.SearchHostIdentifierResult, err error) {
	resp = new(metadata.SearchHostIdentifierResult)
	subPath := "/read/host/indentifier"

	err = h.client.Post().
		WithContext(ctx).
		Body(input).
		SubResource(subPath).
		WithHeaders(header).
		Do().
		Into(resp)
	return
}

func (h *host) GetHostByID(ctx context.Context, header http.Header, hostID string) (resp *metadata.HostInstanceResult, err error) {
	resp = new(metadata.HostInstanceResult)
	subPath := fmt.Sprintf("/find/host/%s", hostID)

	err = h.client.Get().
		WithContext(ctx).
		SubResource(subPath).
		WithHeaders(header).
		Do().
		Into(resp)
	return resp, err
}

func (h *host) GetHosts(ctx context.Context, header http.Header, opt *metadata.QueryInput) (resp *metadata.GetHostsResult, err error) {
	resp = new(metadata.GetHostsResult)
	subPath := "/findmany/hosts/search"

	err = h.client.Post().
		Body(opt).
		WithContext(ctx).
		SubResource(subPath).
		WithHeaders(header).
		Do().
		Into(resp)
	return resp, err
}

func (h *host) GetHostSnap(ctx context.Context, header http.Header, hostID string) (resp *metadata.GetHostSnapResult, err error) {
	resp = new(metadata.GetHostSnapResult)
	subPath := fmt.Sprintf("/find/host/snapshot/%s", hostID)

	err = h.client.Get().
		WithContext(ctx).
		SubResource(subPath).
		WithHeaders(header).
		Do().
		Into(resp)
	return resp, err
}

func (h *host) LockHost(ctx context.Context, header http.Header, input *metadata.HostLockRequest) (resp *metadata.HostLockResponse, err error) {
	resp = new(metadata.HostLockResponse)
	subPath := "/find/host/lock"

	err = h.client.Post().
		Body(input).
		WithContext(ctx).
		SubResource(subPath).
		WithHeaders(header).
		Do().
		Into(resp)
	return resp, err
}

func (h *host) UnlockHost(ctx context.Context, header http.Header, input *metadata.HostLockRequest) (resp *metadata.HostLockResponse, err error) {
	resp = new(metadata.HostLockResponse)
	subPath := "/delete/host/lock"

	err = h.client.Delete().
		Body(input).
		WithContext(ctx).
		SubResource(subPath).
		WithHeaders(header).
		Do().
		Into(resp)
	return resp, err
}

func (h *host) QueryHostLock(ctx context.Context, header http.Header, input *metadata.QueryHostLockRequest) (resp *metadata.HostLockQueryResponse, err error) {
	resp = new(metadata.HostLockQueryResponse)
	subPath := "/findmany/host/lock/search"

	err = h.client.Post().
		Body(input).
		WithContext(ctx).
		SubResource(subPath).
		WithHeaders(header).
		Do().
		Into(resp)
	return resp, err
}

func (h *host) AddUserConfig(ctx context.Context, header http.Header, dat *metadata.UserConfig) (resp *metadata.IDResult, err error) {
	resp = new(metadata.IDResult)
	subPath := "/create/userapi"

	err = h.client.Post().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(header).
		Do().
		Into(resp)
	return
}

func (h *host) UpdateUserConfig(ctx context.Context, businessID string, id string, header http.Header, dat map[string]interface{}) (resp *metadata.BaseResp, err error) {
	resp = new(metadata.BaseResp)
	subPath := fmt.Sprintf("/update/userapi/%s/%s", businessID, id)

	err = h.client.Put().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(header).
		Do().
		Into(resp)
	return
}

func (h *host) DeleteUserConfig(ctx context.Context, businessID string, id string, header http.Header) (resp *metadata.BaseResp, err error) {
	resp = new(metadata.BaseResp)
	subPath := fmt.Sprintf("/delete/userapi/%s/%s", businessID, id)

	err = h.client.Delete().
		WithContext(ctx).
		Body(nil).
		SubResource(subPath).
		WithHeaders(header).
		Do().
		Into(resp)
	return
}

func (h *host) GetUserConfig(ctx context.Context, header http.Header, opt *metadata.QueryInput) (resp *metadata.GetUserConfigResult, err error) {
	resp = new(metadata.GetUserConfigResult)
	subPath := "/findmany/userapi/search"

	err = h.client.Post().
		WithContext(ctx).
		Body(opt).
		SubResource(subPath).
		WithHeaders(header).
		Do().
		Into(resp)
	return
}

func (h *host) GetUserConfigDetail(ctx context.Context, businessID string, id string, header http.Header) (resp *metadata.GetUserConfigDetailResult, err error) {
	resp = new(metadata.GetUserConfigDetailResult)
	subPath := fmt.Sprintf("/find/userapi/detail/%s/%s", businessID, id)

	err = h.client.Get().
		WithContext(ctx).
		Body(nil).
		SubResource(subPath).
		WithHeaders(header).
		Do().
		Into(resp)
	return
}

func (h *host) AddUserCustom(ctx context.Context, user string, header http.Header, dat map[string]interface{}) (resp *metadata.BaseResp, err error) {
	resp = new(metadata.BaseResp)
	subPath := fmt.Sprintf("/create/usercustom/%s", user)

	err = h.client.Post().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(header).
		Do().
		Into(resp)
	return
}

func (h *host) UpdateUserCustomByID(ctx context.Context, user string, id string, header http.Header, dat map[string]interface{}) (resp *metadata.BaseResp, err error) {
	resp = new(metadata.BaseResp)
	subPath := fmt.Sprintf("/update/usercustom/%s/%s", user, id)

	err = h.client.Put().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(header).
		Do().
		Into(resp)
	return
}

func (h *host) GetUserCustomByUser(ctx context.Context, user string, header http.Header) (resp *metadata.GetUserCustomResult, err error) {
	resp = new(metadata.GetUserCustomResult)
	subPath := fmt.Sprintf("/find/usercustom/user/search/%s", user)

	err = h.client.Get().
		WithContext(ctx).
		SubResource(subPath).
		WithHeaders(header).
		Do().
		Into(resp)
	return
}

func (h *host) GetDefaultUserCustom(ctx context.Context, user string, header http.Header) (resp *metadata.GetUserCustomResult, err error) {
	resp = new(metadata.GetUserCustomResult)
	subPath := fmt.Sprintf("/find/usercustom/default/search/%s", user)

	err = h.client.Post().
		WithContext(ctx).
		Body(nil).
		SubResource(subPath).
		WithHeaders(header).
		Do().
		Into(resp)
	return
}

func (h *host) AddHostFavourite(ctx context.Context, user string, header http.Header, dat *metadata.FavouriteParms) (resp *metadata.IDResult, err error) {
	resp = new(metadata.IDResult)
	subPath := fmt.Sprintf("/create/hosts/favorites/%s", user)

	err = h.client.Post().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(header).
		Do().
		Into(resp)
	return
}

func (h *host) UpdateHostFavouriteByID(ctx context.Context, user string, id string, header http.Header, dat map[string]interface{}) (resp *metadata.BaseResp, err error) {
	resp = new(metadata.BaseResp)
	subPath := fmt.Sprintf("/update/hosts/favorites/%s/%s", user, id)

	err = h.client.Put().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(header).
		Do().
		Into(resp)
	return
}

func (h *host) DeleteHostFavouriteByID(ctx context.Context, user string, id string, header http.Header) (resp *metadata.BaseResp, err error) {
	resp = new(metadata.BaseResp)
	subPath := fmt.Sprintf("/delete/hosts/favorites/%s/%s", user, id)

	err = h.client.Delete().
		WithContext(ctx).
		Body(nil).
		SubResource(subPath).
		WithHeaders(header).
		Do().
		Into(resp)
	return
}

func (h *host) ListHostFavourites(ctx context.Context, user string, header http.Header, dat *metadata.QueryInput) (resp *metadata.GetHostFavoriteResult, err error) {
	resp = new(metadata.GetHostFavoriteResult)
	subPath := fmt.Sprintf("/findmany/hosts/favorites/search/%s", user)

	err = h.client.Post().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(header).
		Do().
		Into(resp)
	return
}

func (h *host) GetHostFavouriteByID(ctx context.Context, user string, id string, header http.Header) (resp *metadata.GetHostFavoriteWithIDResult, err error) {
	resp = new(metadata.GetHostFavoriteWithIDResult)
	subPath := fmt.Sprintf("/find/hosts/favorites/search/%s/%s", user, id)

	err = h.client.Get().
		WithContext(ctx).
		SubResource(subPath).
		WithHeaders(header).
		Do().
		Into(resp)
	return
}

func (h *host) GetHostModulesIDs(ctx context.Context, header http.Header, dat *metadata.ModuleHostConfigParams) (resp *metadata.GetHostModuleIDsResult, err error) {
	resp = new(metadata.GetHostModuleIDsResult)
	subPath := "/findmany/meta/hosts/modules/search"

	err = h.client.Post().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(header).
		Do().
		Into(resp)
	return
}

func (h *host) ListHosts(ctx context.Context, header http.Header, option metadata.ListHosts) (metadata.ListHostResult, error) {
	type Result struct {
		metadata.BaseResp `json:",inline"`
		Data              metadata.ListHostResult `json:"data"`
	}
	result := Result{}
	subPath := "/findmany/hosts/list_hosts"

	err := h.client.Post().
		WithContext(ctx).
		Body(option).
		SubResource(subPath).
		WithHeaders(header).
		Do().
		Into(&result)
	if err != nil {
		return result.Data, err
	}
	if result.Code > 0 || result.Result == false {
		return result.Data, errors.New(result.Code, result.ErrMsg)
	}
	return result.Data, nil
}

func (h *host) UpdateHostCloudAreaField(ctx context.Context, header http.Header, option metadata.UpdateHostCloudAreaFieldOption) errors.CCErrorCoder {
	rid := util.GetHTTPCCRequestID(header)

	result := metadata.BaseResp{}
	subPath := "/updatemany/hosts/cloudarea_field"

	err := h.client.Put().
		WithContext(ctx).
		Body(option).
		SubResource(subPath).
		WithHeaders(header).
		Do().
		Into(&result)
	if err != nil {
		blog.Errorf("UpdateHostCloudAreaField failed, http request failed, err: %+v, rid: %s", err, rid)
		return errors.CCHttpError
	}
	if result.Code > 0 || result.Result == false {
		return errors.New(result.Code, result.ErrMsg)
	}
	return nil
}
