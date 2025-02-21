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

package hostserver

import (
	"context"
	"net/http"

	"configcenter/src/common/errors"
	"configcenter/src/common/metadata"
	"configcenter/src/kube/types"
)

// DeleteHostBatch TODO
func (hs *hostServer) DeleteHostBatch(ctx context.Context, h http.Header, dat interface{}) (resp *metadata.Response,
	err error) {
	resp = new(metadata.Response)
	subPath := "/hosts/batch"

	err = hs.client.Delete().
		WithContext(ctx).
		Body(dat).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

// GetHostInstanceProperties TODO
func (hs *hostServer) GetHostInstanceProperties(ctx context.Context, hostID string,
	h http.Header) (resp *metadata.HostInstancePropertiesResult, err error) {
	subPath := "/hosts/%s/%s"

	resp = new(metadata.HostInstancePropertiesResult)
	err = hs.client.Get().
		WithContext(ctx).
		Body(nil).
		// url参数已废弃，此处"0"仅作占位符，不代表实际租户
		SubResourcef(subPath, "0", hostID).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

// AddHost TODO
func (hs *hostServer) AddHost(ctx context.Context, h http.Header, dat interface{}) (resp *metadata.Response,
	err error) {
	resp = new(metadata.Response)
	subPath := "/hosts/add"

	err = hs.client.Post().
		WithContext(ctx).
		Body(dat).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

// AddHostToResourcePool TODO
func (hs *hostServer) AddHostToResourcePool(ctx context.Context, h http.Header,
	dat metadata.AddHostToResourcePoolHostList) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "/hosts/add/resource"

	err = hs.client.Post().
		WithContext(ctx).
		Body(dat).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

// AddHostFromAgent TODO
func (hs *hostServer) AddHostFromAgent(ctx context.Context, h http.Header, dat interface{}) (resp *metadata.Response,
	err error) {
	resp = new(metadata.Response)
	subPath := "/host/add/agent"

	err = hs.client.Post().
		WithContext(ctx).
		Body(dat).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

// SyncHost TODO
func (hs *hostServer) SyncHost(ctx context.Context, h http.Header, data interface{}) (resp *metadata.Response,
	err error) {
	resp = new(metadata.Response)
	subPath := "/hosts/sync/new/host"

	err = hs.client.Post().
		WithContext(ctx).
		Body(data).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

// GetHostFavourites TODO
func (hs *hostServer) GetHostFavourites(ctx context.Context, h http.Header,
	dat interface{}) (resp *metadata.GetHostFavoriteResult, err error) {
	resp = new(metadata.GetHostFavoriteResult)
	subPath := "hosts/favorites/search"

	err = hs.client.Post().
		WithContext(ctx).
		Body(dat).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

// AddHostFavourite TODO
func (hs *hostServer) AddHostFavourite(ctx context.Context, h http.Header,
	dat *metadata.FavouriteParms) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "hosts/favorites"

	err = hs.client.Post().
		WithContext(ctx).
		Body(dat).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

// UpdateHostFavouriteByID TODO
func (hs *hostServer) UpdateHostFavouriteByID(ctx context.Context, id string, h http.Header,
	data *metadata.FavouriteParms) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "hosts/favorites/%s"

	err = hs.client.Put().
		WithContext(ctx).
		Body(data).
		SubResourcef(subPath, id).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

// DeleteHostFavouriteByID TODO
func (hs *hostServer) DeleteHostFavouriteByID(ctx context.Context, id string, h http.Header) (resp *metadata.Response,
	err error) {
	resp = new(metadata.Response)
	subPath := "hosts/favorites/%s"

	err = hs.client.Delete().
		WithContext(ctx).
		Body(nil).
		SubResourcef(subPath, id).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

// IncrHostFavouritesCount TODO
func (hs *hostServer) IncrHostFavouritesCount(ctx context.Context, id string, h http.Header) (resp *metadata.Response,
	err error) {
	resp = new(metadata.Response)
	subPath := "/hosts/favorites/%s/incr"

	err = hs.client.Put().
		WithContext(ctx).
		Body(nil).
		SubResourcef(subPath, id).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

// AddHistory TODO
func (hs *hostServer) AddHistory(ctx context.Context, h http.Header,
	dat map[string]interface{}) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "/hosts/history"

	err = hs.client.Post().
		WithContext(ctx).
		Body(dat).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

// GetHistorys TODO
func (hs *hostServer) GetHistorys(ctx context.Context, start string, limit string,
	h http.Header) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "/hosts/history/%s/%s"

	err = hs.client.Get().
		WithContext(ctx).
		Body(nil).
		SubResourcef(subPath, start, limit).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

// AddHostMultiAppModuleRelation TODO
func (hs *hostServer) AddHostMultiAppModuleRelation(ctx context.Context, h http.Header,
	dat *metadata.CloudHostModuleParams) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "/hosts/modules/biz/mutiple"

	err = hs.client.Post().
		WithContext(ctx).
		Body(dat).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

// TransferHostModule TODO
func (hs *hostServer) TransferHostModule(ctx context.Context, h http.Header,
	params map[string]interface{}) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "/hosts/modules"

	err = hs.client.Post().
		WithContext(ctx).
		Body(params).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

// MoveHost2EmptyModule TODO
func (hs *hostServer) MoveHost2EmptyModule(ctx context.Context, h http.Header,
	dat *metadata.DefaultModuleHostConfigParams) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "/hosts/modules/idle"

	err = hs.client.Post().
		WithContext(ctx).
		Body(dat).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

// MoveHost2FaultModule TODO
func (hs *hostServer) MoveHost2FaultModule(ctx context.Context, h http.Header,
	dat *metadata.DefaultModuleHostConfigParams) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "/hosts/modules/fault"

	err = hs.client.Post().
		WithContext(ctx).
		Body(dat).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

// MoveHostToResourcePool TODO
func (hs *hostServer) MoveHostToResourcePool(ctx context.Context, h http.Header,
	dat *metadata.DefaultModuleHostConfigParams) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "/hosts/modules/resource"

	err = hs.client.Post().
		WithContext(ctx).
		Body(dat).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

// TransferHostAcrossBusiness transfer hosts across biz
func (hs *hostServer) TransferHostAcrossBusiness(ctx context.Context, header http.Header,
	option *metadata.TransferHostAcrossBusinessParameter) errors.CCErrorCoder {

	resp := new(metadata.CreateBatchResult)
	subPath := "/hosts/modules/across/biz"

	err := hs.client.Post().
		WithContext(ctx).
		Body(option).
		SubResourcef(subPath).
		WithHeaders(header).
		Do().
		Into(resp)

	if err != nil {
		return errors.CCHttpError
	}

	if err := resp.CCError(); err != nil {
		return err
	}

	return nil
}

// AssignHostToApp TODO
func (hs *hostServer) AssignHostToApp(ctx context.Context, h http.Header,
	dat *metadata.DefaultModuleHostConfigParams) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "/hosts/modules/resource/idle"

	err = hs.client.Post().
		WithContext(ctx).
		Body(dat).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

// SaveUserCustom TODO
func (hs *hostServer) SaveUserCustom(ctx context.Context, h http.Header, dat interface{}) (resp *metadata.Response,
	err error) {
	resp = new(metadata.Response)
	subPath := "/usercustom"

	err = hs.client.Post().
		WithContext(ctx).
		Body(dat).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

// GetUserCustom TODO
func (hs *hostServer) GetUserCustom(ctx context.Context, h http.Header) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "/usercustom/user/search"

	err = hs.client.Post().
		WithContext(ctx).
		Body(nil).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

// GetDefaultCustom TODO
func (hs *hostServer) GetDefaultCustom(ctx context.Context, h http.Header) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "/usercustom/default/search"

	err = hs.client.Post().
		WithContext(ctx).
		Body(nil).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

// CloneHostProperty TODO
func (hs *hostServer) CloneHostProperty(ctx context.Context, h http.Header,
	dat *metadata.CloneHostPropertyParams) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "/hosts/property/clone"

	err = hs.client.Put().
		WithContext(ctx).
		Body(dat).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

// MoveSetHost2IdleModule TODO
func (hs *hostServer) MoveSetHost2IdleModule(ctx context.Context, h http.Header,
	dat *metadata.SetHostConfigParams) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "/hosts/modules/idle/set"

	err = hs.client.Post().
		WithContext(ctx).
		Body(dat).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

// SearchHostWithNoAuth host query without auth.
func (hs *hostServer) SearchHostWithNoAuth(ctx context.Context, h http.Header, dat *metadata.HostCommonSearch) (
	resp *metadata.SearchHostResult, err error) {

	resp = new(metadata.SearchHostResult)
	subPath := "/findmany/hosts/search/noauth"

	err = hs.client.Post().
		WithContext(ctx).
		Body(dat).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

// SearchHostWithBiz host search with business information
func (hs *hostServer) SearchHostWithBiz(ctx context.Context, h http.Header, dat *metadata.HostCommonSearch) (
	resp *metadata.SearchHostResult, err error) {
	resp = new(metadata.SearchHostResult)
	subPath := "/findmany/hosts/search/with_biz"

	err = hs.client.Post().
		WithContext(ctx).
		Body(dat).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

// SearchHostWithAsstDetail TODO
func (hs *hostServer) SearchHostWithAsstDetail(ctx context.Context, h http.Header,
	dat *metadata.HostCommonSearch) (resp *metadata.SearchHostResult, err error) {
	resp = new(metadata.SearchHostResult)
	subPath := "/hosts/search/asstdetail"

	err = hs.client.Post().
		WithContext(ctx).
		Body(dat).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

// UpdateHostBatch TODO
func (hs *hostServer) UpdateHostBatch(ctx context.Context, h http.Header, dat interface{}) (resp *metadata.Response,
	err error) {
	resp = new(metadata.Response)
	subPath := "/hosts/batch"

	err = hs.client.Put().
		WithContext(ctx).
		Body(dat).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

// UpdateHostPropertyBatch TODO
func (hs *hostServer) UpdateHostPropertyBatch(ctx context.Context, h http.Header,
	data map[string]interface{}) errors.CCErrorCoder {

	resp := new(metadata.Response)
	subPath := "/hosts/property/batch"

	err := hs.client.Put().
		WithContext(ctx).
		Body(data).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)

	if err != nil {
		return errors.CCHttpError
	}

	if ccErr := resp.CCError(); ccErr != nil {
		return ccErr
	}

	return nil
}

// CreateDynamicGroup is dynamic group create action api machinery.
func (hs *hostServer) CreateDynamicGroup(ctx context.Context, header http.Header,
	data *metadata.DynamicGroup) (resp *metadata.IDResult, err error) {

	resp = new(metadata.IDResult)
	subPath := "/dynamicgroup"

	err = hs.client.Post().
		WithContext(ctx).
		Body(data).
		SubResourcef(subPath).
		WithHeaders(header).
		Do().
		Into(resp)
	return
}

// UpdateDynamicGroup is dynamic group update action api machinery.
func (hs *hostServer) UpdateDynamicGroup(ctx context.Context, bizID int64, id string,
	header http.Header, data *metadata.DynamicGroup) (resp *metadata.BaseResp, err error) {

	resp = new(metadata.BaseResp)
	subPath := "/dynamicgroup/%d/%s"

	err = hs.client.Put().
		WithContext(ctx).
		Body(data).
		SubResourcef(subPath, bizID, id).
		WithHeaders(header).
		Do().
		Into(resp)
	return
}

// DeleteDynamicGroup is dynamic group delete action api machinery.
func (hs *hostServer) DeleteDynamicGroup(ctx context.Context, bizID int64, id string,
	header http.Header) (resp *metadata.BaseResp, err error) {

	resp = new(metadata.BaseResp)
	subPath := "/dynamicgroup/%d/%s"

	err = hs.client.Delete().
		WithContext(ctx).
		Body(nil).
		SubResourcef(subPath, bizID, id).
		WithHeaders(header).
		Do().
		Into(resp)
	return
}

// GetDynamicGroup is dynamic group query detail action api machinery.
func (hs *hostServer) GetDynamicGroup(ctx context.Context, bizID int64, id string,
	header http.Header) (resp *metadata.GetDynamicGroupResult, err error) {

	resp = new(metadata.GetDynamicGroupResult)
	subPath := "/dynamicgroup/%d/%s"

	err = hs.client.Get().
		WithContext(ctx).
		Body(nil).
		SubResourcef(subPath, bizID, id).
		WithHeaders(header).
		Do().
		Into(resp)
	return
}

// SearchDynamicGroup is dynamic group search action api machinery.
func (hs *hostServer) SearchDynamicGroup(ctx context.Context, bizID int64, header http.Header,
	data *metadata.QueryCondition) (resp *metadata.SearchDynamicGroupResult, err error) {

	resp = new(metadata.SearchDynamicGroupResult)
	subPath := "/dynamicgroup/search/%d"

	err = hs.client.Post().
		WithContext(ctx).
		Body(data).
		SubResourcef(subPath, bizID).
		WithHeaders(header).
		Do().
		Into(resp)
	return
}

// ExecuteDynamicGroup is dynamic group execute action base on conditions api machinery.
func (hs *hostServer) ExecuteDynamicGroup(ctx context.Context, bizID int64, id string, header http.Header,
	data *metadata.ExecuteOption) (resp *metadata.ResponseInstData, err error) {

	resp = new(metadata.ResponseInstData)
	subPath := "/dynamicgroup/data/%d/%s"

	err = hs.client.Post().
		WithContext(ctx).
		Body(data).
		SubResourcef(subPath, bizID, id).
		WithHeaders(header).
		Do().
		Into(resp)
	return
}

// HostSearch TODO
func (hs *hostServer) HostSearch(ctx context.Context, h http.Header,
	params *metadata.HostCommonSearch) (resp *metadata.QueryInstResult, err error) {

	resp = new(metadata.QueryInstResult)
	subPath := "hosts/search"

	err = hs.client.Post().
		WithContext(ctx).
		Body(params).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

// ListBizHostsTopo TODO
func (hs *hostServer) ListBizHostsTopo(ctx context.Context, h http.Header, bizID int64,
	params *metadata.ListHostsWithNoBizParameter) (resp *metadata.SuccessResponse, err error) {

	resp = new(metadata.SuccessResponse)
	subPath := "/hosts/app/%d/list_hosts_topo"

	err = hs.client.Post().
		WithContext(ctx).
		Body(params).
		SubResourcef(subPath, bizID).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

// CreateCloudArea TODO
func (hs *hostServer) CreateCloudArea(ctx context.Context, h http.Header,
	data map[string]interface{}) (resp *metadata.CreatedOneOptionResult, err error) {

	resp = new(metadata.CreatedOneOptionResult)
	subPath := "/create/cloudarea"

	err = hs.client.Post().
		WithContext(ctx).
		Body(data).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

// CreateManyCloudArea TODO
func (hs *hostServer) CreateManyCloudArea(ctx context.Context, h http.Header,
	data map[string]interface{}) (resp *metadata.CreateManyCloudAreaResult, err error) {

	resp = new(metadata.CreateManyCloudAreaResult)
	subPath := "/createmany/cloudarea"

	err = hs.client.Post().
		WithContext(ctx).
		Body(data).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

// UpdateCloudArea TODO
func (hs *hostServer) UpdateCloudArea(ctx context.Context, h http.Header, cloudID int64,
	data map[string]interface{}) (resp *metadata.Response, err error) {

	resp = new(metadata.Response)
	subPath := "/update/cloudarea/%d"

	err = hs.client.Put().
		WithContext(ctx).
		Body(data).
		SubResourcef(subPath, cloudID).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

// SearchCloudArea TODO
func (hs *hostServer) SearchCloudArea(ctx context.Context, h http.Header,
	params map[string]interface{}) (resp *metadata.SearchResp, err error) {

	resp = new(metadata.SearchResp)
	subPath := "/findmany/cloudarea"

	err = hs.client.Post().
		WithContext(ctx).
		Body(params).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

// DeleteCloudArea TODO
func (hs *hostServer) DeleteCloudArea(ctx context.Context, h http.Header, cloudID int64) (resp *metadata.Response,
	err error) {

	resp = new(metadata.Response)
	subPath := "/delete/cloudarea/%d"

	err = hs.client.Delete().
		WithContext(ctx).
		Body(nil).
		SubResourcef(subPath, cloudID).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

// FindCloudAreaHostCount TODO
func (hs *hostServer) FindCloudAreaHostCount(ctx context.Context, header http.Header,
	option metadata.CloudAreaHostCount) (resp *metadata.CloudAreaHostCountResult, err error) {
	resp = new(metadata.CloudAreaHostCountResult)
	subPath := "/findmany/cloudarea/hostcount"

	err = hs.client.Post().
		WithContext(ctx).
		Body(option).
		SubResourcef(subPath).
		WithHeaders(header).
		Do().
		Into(resp)
	return
}

// SearchHostWithKube search host with k8s condition
func (hs *hostServer) SearchKubeHost(ctx context.Context, h http.Header, req types.SearchHostOption) (
	*metadata.SearchHost, errors.CCErrorCoder) {

	result := new(metadata.SearchHostResult)
	subPath := "/hosts/kube/search"

	err := hs.client.Post().
		WithContext(ctx).
		Body(req).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(result)

	if err != nil {
		return nil, errors.CCHttpError
	}

	if ccErr := result.CCError(); ccErr != nil {
		return nil, ccErr
	}

	return result.Data, nil
}

// AddCloudHostToBiz add cloud host to biz idle module
func (hs *hostServer) AddCloudHostToBiz(ctx context.Context, h http.Header, opt *metadata.AddCloudHostToBizParam) (
	*metadata.RspIDs, errors.CCErrorCoder) {

	resp := new(metadata.CreateBatchResult)
	subPath := "/createmany/cloud_hosts"

	err := hs.client.Post().
		WithContext(ctx).
		Body(opt).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)

	if err != nil {
		return nil, errors.CCHttpError
	}

	if err := resp.CCError(); err != nil {
		return nil, err
	}

	return &resp.Data, nil
}

// DeleteCloudHostFromBiz delete cloud hosts from biz
func (hs *hostServer) DeleteCloudHostFromBiz(ctx context.Context, header http.Header,
	option *metadata.DeleteCloudHostFromBizParam) errors.CCErrorCoder {

	resp := new(metadata.CreateBatchResult)
	subPath := "/deletemany/cloud_hosts"

	err := hs.client.Delete().
		WithContext(ctx).
		Body(option).
		SubResourcef(subPath).
		WithHeaders(header).
		Do().
		Into(resp)

	if err != nil {
		return errors.CCHttpError
	}

	if err := resp.CCError(); err != nil {
		return err
	}

	return nil
}

// BindAgent bind agent to host
func (hs *hostServer) BindAgent(ctx context.Context, h http.Header,
	params *metadata.BindAgentParam) errors.CCErrorCoder {

	resp := new(metadata.BaseResp)
	subPath := "/host/bind/agent"

	err := hs.client.Post().
		WithContext(ctx).
		Body(params).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)

	if err != nil {
		return errors.CCHttpError
	}
	if resp.CCError() != nil {
		return resp.CCError()
	}
	return nil
}

// UnbindAgent unbind agent to host
func (hs *hostServer) UnbindAgent(ctx context.Context, h http.Header,
	params *metadata.UnbindAgentParam) errors.CCErrorCoder {

	resp := new(metadata.BaseResp)
	subPath := "/host/unbind/agent"

	err := hs.client.Post().
		WithContext(ctx).
		Body(params).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)

	if err != nil {
		return errors.CCHttpError
	}
	if resp.CCError() != nil {
		return resp.CCError()
	}
	return nil
}

// AddHostToBizIdle add host to biz idle module
func (hs *hostServer) AddHostToBizIdle(ctx context.Context, header http.Header, option *metadata.HostListParam) (
	*metadata.HostIDsResp, errors.CCErrorCoder) {

	resp := new(metadata.CreateHostBatchResult)
	subPath := "/hosts/add/business_idle"

	err := hs.client.Post().
		WithContext(ctx).
		Body(option).
		SubResourcef(subPath).
		WithHeaders(header).
		Do().
		Into(resp)

	if err != nil {
		return nil, errors.CCHttpError
	}

	if err := resp.CCError(); err != nil {
		return nil, err
	}

	return &resp.Data, nil
}

// CountBizHostCPU count biz host CPU
func (hs *hostServer) CountBizHostCPU(ctx context.Context, header http.Header, option *metadata.CountHostCPUReq) (
	*metadata.BizHostCpuCount, errors.CCErrorCoder) {

	resp := new(metadata.BizHostCpuCountResult)
	subPath := "/host/count/cpu"

	err := hs.client.Post().
		WithContext(ctx).
		Body(option).
		SubResourcef(subPath).
		WithHeaders(header).
		Do().
		Into(resp)

	if err != nil {
		return nil, errors.CCHttpError
	}

	if err := resp.CCError(); err != nil {
		return nil, err
	}

	return &resp.Data[0], nil
}

// UpdateHostsAllProperty updatemany hosts all property
func (hs *hostServer) UpdateHostsAllProperty(ctx context.Context, header http.Header,
	option *metadata.UpdateHostOpt) errors.CCErrorCoder {

	resp := new(metadata.BaseResp)
	subPath := "/updatemany/hosts/all/property"

	err := hs.client.Put().
		WithContext(ctx).
		Body(option).
		SubResourcef(subPath).
		WithHeaders(header).
		Do().
		Into(resp)

	if err != nil {
		return errors.CCHttpError
	}

	if err := resp.CCError(); err != nil {
		return err
	}

	return nil
}

// FindModuleHostRelation findmany module host relation
func (hs *hostServer) FindModuleHostRelation(ctx context.Context, header http.Header, bizID int64,
	option *metadata.FindModuleHostRelationParameter) (*metadata.FindModuleHostRelationResult, errors.CCErrorCoder) {

	resp := new(metadata.FindModuleHostRelationResp)
	subPath := "/findmany/module_relation/bk_biz_id/%d"

	err := hs.client.Post().
		WithContext(ctx).
		Body(option).
		SubResourcef(subPath, bizID).
		WithHeaders(header).
		Do().
		Into(resp)

	if err != nil {
		return nil, errors.CCHttpError
	}

	if err := resp.CCError(); err != nil {
		return nil, err
	}

	return &resp.Data, nil
}

// FindHostByServiceTmpl findmany host by service templates
func (hs *hostServer) FindHostByServiceTmpl(ctx context.Context, header http.Header, bizID int64,
	option *metadata.FindHostsBySrvTplOpt) (*metadata.SearchHost, errors.CCErrorCoder) {

	resp := new(metadata.SearchHostResult)
	subPath := "/findmany/hosts/by_service_templates/biz/%d"

	err := hs.client.Post().
		WithContext(ctx).
		Body(option).
		SubResourcef(subPath, bizID).
		WithHeaders(header).
		Do().
		Into(resp)

	if err != nil {
		return nil, errors.CCHttpError
	}

	if err := resp.CCError(); err != nil {
		return nil, err
	}

	return resp.Data, nil
}

// FindHostBySetTmpl findmany host by set templates
func (hs *hostServer) FindHostBySetTmpl(ctx context.Context, header http.Header, bizID int64,
	option *metadata.FindHostsBySetTplOpt) (*metadata.SearchHost, errors.CCErrorCoder) {

	resp := new(metadata.SearchHostResult)
	subPath := "/findmany/hosts/by_set_templates/biz/%d"

	err := hs.client.Post().
		WithContext(ctx).
		Body(option).
		SubResourcef(subPath, bizID).
		WithHeaders(header).
		Do().
		Into(resp)

	if err != nil {
		return nil, errors.CCHttpError
	}

	if err := resp.CCError(); err != nil {
		return nil, err
	}

	return resp.Data, nil
}

// ListResourcePoolHosts list resource pool hosts
func (hs *hostServer) ListResourcePoolHosts(ctx context.Context, header http.Header,
	option *metadata.ListHostsParameter) (*metadata.ListHostResult, errors.CCErrorCoder) {

	resp := new(metadata.ListHostResp)
	subPath := "/hosts/list_resource_pool_hosts"

	err := hs.client.Post().
		WithContext(ctx).
		Body(option).
		SubResourcef(subPath).
		WithHeaders(header).
		Do().
		Into(resp)

	if err != nil {
		return nil, errors.CCHttpError
	}

	if err := resp.CCError(); err != nil {
		return nil, err
	}

	return resp.Data, nil
}

// ListHostsWithoutApp list hosts without app
func (hs *hostServer) ListHostsWithoutApp(ctx context.Context, header http.Header,
	option *metadata.ListHostsWithNoBizParameter) (*metadata.ListHostResult, errors.CCErrorCoder) {

	resp := new(metadata.ListHostResp)
	subPath := "/hosts/list_hosts_without_app"

	err := hs.client.Post().
		WithContext(ctx).
		Body(option).
		SubResourcef(subPath).
		WithHeaders(header).
		Do().
		Into(resp)

	if err != nil {
		return nil, errors.CCHttpError
	}

	if err := resp.CCError(); err != nil {
		return nil, err
	}

	return resp.Data, nil
}

// FindHostByTopoInst findmany host by topo inst
func (hs *hostServer) FindHostByTopoInst(ctx context.Context, header http.Header, bizID int64,
	option *metadata.FindHostsByTopoOpt) (*metadata.SearchHost, errors.CCErrorCoder) {

	resp := new(metadata.SearchHostResult)
	subPath := "/findmany/hosts/by_topo/biz/%d"

	err := hs.client.Post().
		WithContext(ctx).
		Body(option).
		SubResourcef(subPath, bizID).
		WithHeaders(header).
		Do().
		Into(resp)

	if err != nil {
		return nil, errors.CCHttpError
	}

	if err := resp.CCError(); err != nil {
		return nil, err
	}

	return resp.Data, nil
}

// FindHostDetailTopo findmany host detail topo
func (hs *hostServer) FindHostDetailTopo(ctx context.Context, header http.Header,
	option *metadata.ListHostsDetailAndTopoOption) (*metadata.HostMainlineTopoResult, errors.CCErrorCoder) {

	resp := new(metadata.HostMainlineTopoResp)
	subPath := "/findmany/hosts/detail_topo"

	err := hs.client.Post().
		WithContext(ctx).
		Body(option).
		SubResourcef(subPath).
		WithHeaders(header).
		Do().
		Into(resp)

	if err != nil {
		return nil, errors.CCHttpError
	}

	if err := resp.CCError(); err != nil {
		return nil, err
	}

	return resp.Data, nil
}

// FindHostRelationWithTopo findmany host detail topo
func (hs *hostServer) FindHostRelationWithTopo(ctx context.Context, header http.Header,
	option *metadata.FindHostRelationWtihTopoOpt) (*metadata.HostConfigResult, errors.CCErrorCoder) {

	resp := new(metadata.HostConfigResp)
	subPath := "/findmany/hosts/relation/with_topo"

	err := hs.client.Post().
		WithContext(ctx).
		Body(option).
		SubResourcef(subPath).
		WithHeaders(header).
		Do().
		Into(resp)

	if err != nil {
		return nil, errors.CCHttpError
	}

	if err := resp.CCError(); err != nil {
		return nil, err
	}

	return &resp.Data, nil
}

// FindHostServiceTmpl findmany host service template
func (hs *hostServer) FindHostServiceTmpl(ctx context.Context, header http.Header, option *metadata.HostIDReq) (
	*metadata.HostSrvTmplResp, errors.CCErrorCoder) {

	resp := new(metadata.HostSrvTmplResp)
	subPath := "/findmany/hosts/service_template"

	err := hs.client.Post().
		WithContext(ctx).
		Body(option).
		SubResourcef(subPath).
		WithHeaders(header).
		Do().
		Into(resp)

	if err != nil {
		return nil, errors.CCHttpError
	}

	if err := resp.CCError(); err != nil {
		return nil, err
	}

	return resp, nil
}

// FindHostTotalMainlineTopo findmany host total mainline topo
func (hs *hostServer) FindHostTotalMainlineTopo(ctx context.Context, header http.Header, bizID int64,
	option *metadata.FindHostTotalTopo) (*metadata.HostMainlineTopoResult, errors.CCErrorCoder) {

	resp := new(metadata.HostMainlineTopoResp)
	subPath := "/findmany/hosts/total_mainline_topo/biz/%d"

	err := hs.client.Post().
		WithContext(ctx).
		Body(option).
		SubResourcef(subPath).
		WithHeaders(header).
		Do().
		Into(resp)

	if err != nil {
		return nil, errors.CCHttpError
	}

	if err := resp.CCError(); err != nil {
		return nil, err
	}

	return resp.Data, nil
}

// UpdateHostCloudArea update host cloud area
func (hs *hostServer) UpdateHostCloudArea(ctx context.Context, header http.Header,
	option *metadata.UpdateHostCloudAreaFieldOption) errors.CCErrorCoder {

	resp := new(metadata.BaseResp)
	subPath := "/updatemany/hosts/cloudarea_field"

	err := hs.client.Put().
		WithContext(ctx).
		Body(option).
		SubResourcef(subPath).
		WithHeaders(header).
		Do().
		Into(resp)

	if err != nil {
		return errors.CCHttpError
	}

	if err := resp.CCError(); err != nil {
		return err
	}

	return nil
}

// UpdateHostToRecycle update host to recycle module
func (hs *hostServer) UpdateHostToRecycle(ctx context.Context, header http.Header,
	option *metadata.DefaultModuleHostConfigParams) errors.CCErrorCoder {

	resp := new(metadata.HostSrvTmplResp)
	subPath := "/hosts/modules/recycle"

	err := hs.client.Post().
		WithContext(ctx).
		Body(option).
		SubResourcef(subPath).
		WithHeaders(header).
		Do().
		Into(resp)

	if err != nil {
		return errors.CCHttpError
	}

	if err := resp.CCError(); err != nil {
		return err
	}

	return nil
}

// FindHostModules find host modules
func (hs *hostServer) FindHostModules(ctx context.Context, header http.Header,
	option *metadata.HostModuleRelationParameter) (*metadata.HostModuleResp, errors.CCErrorCoder) {

	resp := new(metadata.HostModuleResp)
	subPath := "/hosts/modules/read"

	err := hs.client.Post().
		WithContext(ctx).
		Body(option).
		SubResourcef(subPath).
		WithHeaders(header).
		Do().
		Into(resp)

	if err != nil {
		return nil, errors.CCHttpError
	}

	if err := resp.CCError(); err != nil {
		return nil, err
	}

	return resp, nil
}

// FindHostTopoRelation find host topo relation
func (hs *hostServer) FindHostTopoRelation(ctx context.Context, header http.Header,
	option *metadata.HostModuleRelationRequest) (*metadata.HostConfigData, errors.CCErrorCoder) {

	resp := new(metadata.HostConfig)
	subPath := "/host/topo/relation/read"

	err := hs.client.Post().
		WithContext(ctx).
		Body(option).
		SubResourcef(subPath).
		WithHeaders(header).
		Do().
		Into(resp)

	if err != nil {
		return nil, errors.CCHttpError
	}

	if err := resp.CCError(); err != nil {
		return nil, err
	}

	return &resp.Data, nil
}

// TransferHostResourceDirectory transfer host resource directory
func (hs *hostServer) TransferHostResourceDirectory(ctx context.Context, header http.Header,
	option *metadata.TransferHostResourceDirectory) errors.CCErrorCoder {

	resp := new(metadata.BaseResp)
	subPath := "/host/transfer/resource/directory"

	err := hs.client.Post().
		WithContext(ctx).
		Body(option).
		SubResourcef(subPath).
		WithHeaders(header).
		Do().
		Into(resp)

	if err != nil {
		return errors.CCHttpError
	}

	if err := resp.CCError(); err != nil {
		return err
	}

	return nil
}
