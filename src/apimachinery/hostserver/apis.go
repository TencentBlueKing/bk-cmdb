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
	"fmt"
	"net/http"

	"configcenter/src/common/metadata"
	"configcenter/src/common/paraparse"
)

func (hs *hostServer) DeleteHostBatch(ctx context.Context, h http.Header, dat interface{}) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "/hosts/batch"

	err = hs.client.Delete().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (hs *hostServer) GetHostInstanceProperties(ctx context.Context, ownerID string, hostID string, h http.Header) (resp *metadata.HostInstancePropertiesResult, err error) {
	subPath := fmt.Sprintf("/hosts/%s/%s", ownerID, hostID)

	err = hs.client.Get().
		WithContext(ctx).
		Body(nil).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (hs *hostServer) HostSnapInfo(ctx context.Context, hostID string, h http.Header, dat interface{}) (resp *metadata.HostSnapResult, err error) {
	subPath := fmt.Sprintf("/hosts/snapshot/%s", hostID)

	err = hs.client.Get().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (hs *hostServer) AddHost(ctx context.Context, h http.Header, dat interface{}) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "/hosts/add"

	err = hs.client.Post().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (hs *hostServer) AddHostFromAgent(ctx context.Context, h http.Header, dat interface{}) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "/host/add/agent"

	err = hs.client.Post().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (hs *hostServer) GetHostFavourites(ctx context.Context, h http.Header, dat interface{}) (resp *metadata.GetHostFavoriteResult, err error) {
	resp = new(metadata.GetHostFavoriteResult)
	subPath := "hosts/favorites/search"

	err = hs.client.Post().
		WithContext(ctx).
		Body(nil).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (hs *hostServer) AddHostFavourite(ctx context.Context, h http.Header, dat *metadata.FavouriteParms) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "hosts/favorites"

	err = hs.client.Post().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (hs *hostServer) UpdateHostFavouriteByID(ctx context.Context, id string, h http.Header) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := fmt.Sprintf("hosts/favorites/%s", id)

	err = hs.client.Put().
		WithContext(ctx).
		Body(nil).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (hs *hostServer) DeleteHostFavouriteByID(ctx context.Context, id string, h http.Header) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := fmt.Sprintf("hosts/favorites/%s", id)

	err = hs.client.Delete().
		WithContext(ctx).
		Body(nil).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (hs *hostServer) IncrHostFavouritesCount(ctx context.Context, id string, h http.Header) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := fmt.Sprintf("/hosts/favorites/%s/incr", id)

	err = hs.client.Put().
		WithContext(ctx).
		Body(nil).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (hs *hostServer) AddHistory(ctx context.Context, h http.Header, dat map[string]interface{}) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "/hosts/history"

	err = hs.client.Post().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (hs *hostServer) GetHistorys(ctx context.Context, start string, limit string, h http.Header) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := fmt.Sprintf("/hosts/history/%s/%s", start, limit)

	err = hs.client.Get().
		WithContext(ctx).
		Body(nil).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (hs *hostServer) AddHostMultiAppModuleRelation(ctx context.Context, h http.Header, dat *metadata.CloudHostModuleParams) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "/hosts/modules/biz/mutiple"

	err = hs.client.Post().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (hs *hostServer) HostModuleRelation(ctx context.Context, h http.Header, params map[string]interface{}) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "/hosts/modules"

	err = hs.client.Post().
		WithContext(ctx).
		Body(params).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (hs *hostServer) MoveHost2EmptyModule(ctx context.Context, h http.Header, dat *metadata.DefaultModuleHostConfigParams) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "/hosts/modules/idle"

	err = hs.client.Post().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (hs *hostServer) MoveHost2FaultModule(ctx context.Context, h http.Header, dat *metadata.DefaultModuleHostConfigParams) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "/hosts/modules/fault"

	err = hs.client.Post().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (hs *hostServer) MoveHostToResourcePool(ctx context.Context, h http.Header, dat *metadata.DefaultModuleHostConfigParams) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "/hosts/modules/resource"

	err = hs.client.Post().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (hs *hostServer) AssignHostToApp(ctx context.Context, h http.Header, dat *metadata.DefaultModuleHostConfigParams) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "/modules/resource/idle"

	err = hs.client.Post().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (hs *hostServer) AssignHostToAppModule(ctx context.Context, h http.Header, dat *metadata.HostToAppModule) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "/host/add/module"

	err = hs.client.Post().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (hs *hostServer) SaveUserCustom(ctx context.Context, h http.Header, dat interface{}) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "/usercustom"

	err = hs.client.Post().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (hs *hostServer) GetUserCustom(ctx context.Context, h http.Header) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "/usercustom/user/search"

	err = hs.client.Post().
		WithContext(ctx).
		Body(nil).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (hs *hostServer) GetDefaultCustom(ctx context.Context, h http.Header) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "/usercustom/default/search"

	err = hs.client.Post().
		WithContext(ctx).
		Body(nil).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (hs *hostServer) GetAgentStatus(ctx context.Context, appID string, h http.Header) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := fmt.Sprintf("getAgentStatus/%s", appID)

	err = hs.client.Get().
		WithContext(ctx).
		Body(nil).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (hs *hostServer) UpdateHost(ctx context.Context, appID string, h http.Header, dat interface{}) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := fmt.Sprintf("/openapi/host/%s", appID)

	err = hs.client.Put().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (hs *hostServer) UpdateHostByAppID(ctx context.Context, appID string, h http.Header, dat interface{}) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := fmt.Sprintf("/host/updateHostByAppID/%s", appID)

	err = hs.client.Put().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (hs *hostServer) GetHostListByAppidAndField(ctx context.Context, appID string, field string, h http.Header) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := fmt.Sprintf("/host/getHostListByAppidAndField/%s/%s", appID, field)

	err = hs.client.Get().
		WithContext(ctx).
		Body(nil).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (hs *hostServer) HostSearchByIP(ctx context.Context, h http.Header, dat interface{}) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "/gethostlistbyip"

	err = hs.client.Post().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (hs *hostServer) HostSearchByModuleID(ctx context.Context, h http.Header, dat interface{}) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "/getmodulehostlist"

	err = hs.client.Post().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (hs *hostServer) HostSearchBySetID(ctx context.Context, h http.Header, dat interface{}) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "/getsethostlist"

	err = hs.client.Post().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (hs *hostServer) HostSearchByAppID(ctx context.Context, h http.Header, dat interface{}) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "/getapphostlist"

	err = hs.client.Post().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (hs *hostServer) HostSearchByProperty(ctx context.Context, h http.Header, dat map[string]interface{}) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "/gethostsbyproperty"

	err = hs.client.Post().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (hs *hostServer) GetIPAndProxyByCompany(ctx context.Context, h http.Header, dat map[string]interface{}) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "/getIPAndProxyByCompany"

	err = hs.client.Post().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (hs *hostServer) UpdateCustomProperty(ctx context.Context, h http.Header, dat map[string]interface{}) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "/openapi/updatecustomproperty"

	err = hs.client.Put().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (hs *hostServer) CloneHostProperty(ctx context.Context, h http.Header, dat *metadata.CloneHostPropertyParams) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "/hosts/property/clone"

	err = hs.client.Put().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (hs *hostServer) MoveSetHost2IdleModule(ctx context.Context, h http.Header, dat *metadata.SetHostConfigParams) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "/hosts/modules/idle/set"

	err = hs.client.Put().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (hs *hostServer) GetHostAppByCompanyId(ctx context.Context, h http.Header, dat map[string]interface{}) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "/openapi/host/getHostAppByCompanyId"

	err = hs.client.Post().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (hs *hostServer) DelHostInApp(ctx context.Context, h http.Header, dat interface{}) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "/openapi/host/delhostinapp"

	err = hs.client.Delete().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (hs *hostServer) GetGitServerIp(ctx context.Context, h http.Header, dat interface{}) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "/openapi/host/getGitServerIp"

	err = hs.client.Post().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (hs *hostServer) GetPlat(ctx context.Context, h http.Header) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "/plat"

	err = hs.client.Get().
		WithContext(ctx).
		Body(nil).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (hs *hostServer) CreatePlat(ctx context.Context, h http.Header, dat map[string]interface{}) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "/plat"

	err = hs.client.Post().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (hs *hostServer) DelPlat(ctx context.Context, cloudID string, h http.Header) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := fmt.Sprintf("/plat/%s", cloudID)

	err = hs.client.Delete().
		WithContext(ctx).
		Body(nil).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (hs *hostServer) SearchHost(ctx context.Context, h http.Header, dat *params.HostCommonSearch) (resp *metadata.SearchHostResult, err error) {
	resp = new(metadata.SearchHostResult)
	subPath := "/hosts/search"

	err = hs.client.Post().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (hs *hostServer) SearchHostWithAsstDetail(ctx context.Context, h http.Header, dat *params.HostCommonSearch) (resp *metadata.SearchHostResult, err error) {
	resp = new(metadata.SearchHostResult)
	subPath := "/hosts/search/asstdetail"

	err = hs.client.Post().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (hs *hostServer) UpdateHostBatch(ctx context.Context, h http.Header, dat interface{}) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "/hosts/batch"

	err = hs.client.Put().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (hs *hostServer) AddUserCustomQuery(ctx context.Context, h http.Header, dat map[string]interface{}) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "/userapi"

	err = hs.client.Post().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (hs *hostServer) UpdateUserCustomQuery(ctx context.Context, businessID string, id string, h http.Header, dat map[string]interface{}) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := fmt.Sprintf("/userapi/%s/%s", businessID, id)

	err = hs.client.Put().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (hs *hostServer) DeleteUserCustomQuery(ctx context.Context, businessID string, id string, h http.Header) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := fmt.Sprintf("/userapi/%s/%s", businessID, id)

	err = hs.client.Delete().
		WithContext(ctx).
		Body(nil).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (hs *hostServer) GetUserCustomQuery(ctx context.Context, businessID string, h http.Header, dat *metadata.QueryInput) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := fmt.Sprintf("/userapi/search/%s", businessID)

	err = hs.client.Post().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (hs *hostServer) GetUserCustomQueryDetail(ctx context.Context, businessID string, id string, h http.Header) (resp *metadata.UserCustomQueryDetailResult, err error) {
	resp = new(metadata.UserCustomQueryDetailResult)
	subPath := fmt.Sprintf("/userapi/detail/%s/%s", businessID, id)

	err = hs.client.Get().
		WithContext(ctx).
		Body(nil).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (hs *hostServer) GetUserCustomQueryResult(ctx context.Context, businessID, id, start, limit string, h http.Header) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := fmt.Sprintf("/userapi/data/%s/%s/%s/%s", businessID, id, start, limit)

	err = hs.client.Get().
		WithContext(ctx).
		Body(nil).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}
