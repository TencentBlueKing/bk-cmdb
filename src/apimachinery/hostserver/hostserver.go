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

	"configcenter/src/apimachinery/rest"
	"configcenter/src/apimachinery/util"
	"configcenter/src/common/metadata"
	"configcenter/src/common/paraparse"
)

type HostServerClientInterface interface {
	DeleteHostBatch(ctx context.Context, h http.Header, dat interface{}) (resp *metadata.Response, err error)
	GetHostInstanceProperties(ctx context.Context, ownerID string, hostID string, h http.Header) (resp *metadata.HostInstancePropertiesResult, err error)
	HostSnapInfo(ctx context.Context, hostID string, h http.Header, dat interface{}) (resp *metadata.HostSnapResult, err error)
	AddHost(ctx context.Context, h http.Header, dat interface{}) (resp *metadata.Response, err error)
	AddHostFromAgent(ctx context.Context, h http.Header, dat interface{}) (resp *metadata.Response, err error)

	GetHostFavourites(ctx context.Context, h http.Header, dat interface{}) (resp *metadata.GetHostFavoriteResult, err error)
	AddHostFavourite(ctx context.Context, h http.Header, dat *metadata.FavouriteParms) (resp *metadata.Response, err error)
	UpdateHostFavouriteByID(ctx context.Context, id string, h http.Header) (resp *metadata.Response, err error)
	DeleteHostFavouriteByID(ctx context.Context, id string, h http.Header) (resp *metadata.Response, err error)
	IncrHostFavouritesCount(ctx context.Context, id string, h http.Header) (resp *metadata.Response, err error)

	AddHistory(ctx context.Context, h http.Header, dat map[string]interface{}) (resp *metadata.Response, err error)
	GetHistorys(ctx context.Context, start string, limit string, h http.Header) (resp *metadata.Response, err error)

	AddHostMultiAppModuleRelation(ctx context.Context, h http.Header, dat *metadata.CloudHostModuleParams) (resp *metadata.Response, err error)
	HostModuleRelation(ctx context.Context, h http.Header, params map[string]interface{}) (resp *metadata.Response, err error)

	MoveHost2EmptyModule(ctx context.Context, h http.Header, dat *metadata.DefaultModuleHostConfigParams) (resp *metadata.Response, err error)
	MoveHost2FaultModule(ctx context.Context, h http.Header, dat *metadata.DefaultModuleHostConfigParams) (resp *metadata.Response, err error)
	MoveHostToResourcePool(ctx context.Context, h http.Header, dat *metadata.DefaultModuleHostConfigParams) (resp *metadata.Response, err error)

	AssignHostToApp(ctx context.Context, h http.Header, dat *metadata.DefaultModuleHostConfigParams) (resp *metadata.Response, err error)
	AssignHostToAppModule(ctx context.Context, h http.Header, dat *metadata.HostToAppModule) (resp *metadata.Response, err error)
	SaveUserCustom(ctx context.Context, h http.Header, dat interface{}) (resp *metadata.Response, err error)
	GetUserCustom(ctx context.Context, h http.Header) (resp *metadata.Response, err error)
	GetDefaultCustom(ctx context.Context, h http.Header) (resp *metadata.Response, err error)
	GetAgentStatus(ctx context.Context, appID string, h http.Header) (resp *metadata.Response, err error)
	UpdateHost(ctx context.Context, appID string, h http.Header, dat interface{}) (resp *metadata.Response, err error)
	UpdateHostByAppID(ctx context.Context, appID string, h http.Header, dat interface{}) (resp *metadata.Response, err error)
	GetHostListByAppidAndField(ctx context.Context, appID string, field string, h http.Header) (resp *metadata.Response, err error)
	HostSearchByIP(ctx context.Context, h http.Header, dat interface{}) (resp *metadata.Response, err error)
	HostSearchByModuleID(ctx context.Context, h http.Header, dat interface{}) (resp *metadata.Response, err error)
	HostSearchBySetID(ctx context.Context, h http.Header, dat interface{}) (resp *metadata.Response, err error)
	HostSearchByAppID(ctx context.Context, h http.Header, dat interface{}) (resp *metadata.Response, err error)
	HostSearchByProperty(ctx context.Context, h http.Header, dat map[string]interface{}) (resp *metadata.Response, err error)
	GetIPAndProxyByCompany(ctx context.Context, h http.Header, dat map[string]interface{}) (resp *metadata.Response, err error)
	UpdateCustomProperty(ctx context.Context, h http.Header, dat map[string]interface{}) (resp *metadata.Response, err error)
	CloneHostProperty(ctx context.Context, h http.Header, dat *metadata.CloneHostPropertyParams) (resp *metadata.Response, err error)
	MoveSetHost2IdleModule(ctx context.Context, h http.Header, dat *metadata.SetHostConfigParams) (resp *metadata.Response, err error)
	GetHostAppByCompanyId(ctx context.Context, h http.Header, dat map[string]interface{}) (resp *metadata.Response, err error)
	DelHostInApp(ctx context.Context, h http.Header, dat interface{}) (resp *metadata.Response, err error)
	GetGitServerIp(ctx context.Context, h http.Header, dat interface{}) (resp *metadata.Response, err error)
	GetPlat(ctx context.Context, h http.Header) (resp *metadata.Response, err error)
	CreatePlat(ctx context.Context, h http.Header, dat map[string]interface{}) (resp *metadata.Response, err error)
	DelPlat(ctx context.Context, cloudID string, h http.Header) (resp *metadata.Response, err error)
	SearchHost(ctx context.Context, h http.Header, dat *params.HostCommonSearch) (resp *metadata.SearchHostResult, err error)
	SearchHostWithAsstDetail(ctx context.Context, h http.Header, dat *params.HostCommonSearch) (resp *metadata.SearchHostResult, err error)
	UpdateHostBatch(ctx context.Context, h http.Header, dat interface{}) (resp *metadata.Response, err error)
	AddUserCustomQuery(ctx context.Context, h http.Header, dat map[string]interface{}) (resp *metadata.Response, err error)
	UpdateUserCustomQuery(ctx context.Context, businessID string, id string, h http.Header, dat map[string]interface{}) (resp *metadata.Response, err error)
	DeleteUserCustomQuery(ctx context.Context, businessID string, id string, h http.Header) (resp *metadata.Response, err error)
	GetUserCustomQuery(ctx context.Context, businessID string, h http.Header, dat *metadata.QueryInput) (resp *metadata.Response, err error)
	GetUserCustomQueryDetail(ctx context.Context, businessID string, id string, h http.Header) (resp *metadata.UserCustomQueryDetailResult, err error)
	GetUserCustomQueryResult(ctx context.Context, businessID, id, start, limit string, h http.Header) (resp *metadata.Response, err error)
}

func NewHostServerClientInterface(c *util.Capability, version string) HostServerClientInterface {
	base := fmt.Sprintf("/host/%s", version)
	return &hostServer{
		client: rest.NewRESTClient(c, base),
	}
}

type hostServer struct {
	client rest.ClientInterface
}
