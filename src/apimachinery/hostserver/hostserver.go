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

    "configcenter/src/apimachinery/rest"
    "configcenter/src/apimachinery/util"
    "configcenter/src/scene_server/host_server/host_service/actions/hosts"
    "configcenter/src/common/core/cc/api"
    "configcenter/src/common/paraparse"
    "configcenter/src/source_controller/common/commondata"
    "fmt"
)

type HostServerClientInterface interface {
    DeleteHostBatch(ctx context.Context, h util.Headers, dat interface{}) (resp *api.BKAPIRsp, err error)
    GetHostDetailByID(ctx context.Context, hostID string, h util.Headers) (resp *api.BKAPIRsp, err error)
    HostSnapInfo(ctx context.Context, hostID string, h util.Headers, dat interface{}) (resp *api.BKAPIRsp, err error)
    AddHost(ctx context.Context, h util.Headers, dat interface{}) (resp *api.BKAPIRsp, err error)
    AddHostFromAgent(ctx context.Context, h util.Headers, dat interface{}) (resp *api.BKAPIRsp, err error)
    GetHostFavourites(ctx context.Context, h util.Headers, dat interface{}) (resp *api.BKAPIRsp, err error)
    AddHostFavourite(ctx context.Context, h util.Headers, dat *hosts.FavouriteParms) (resp *api.BKAPIRsp, err error)
    UpdateHostFavouriteByID(ctx context.Context, id string, h util.Headers) (resp *api.BKAPIRsp, err error)
    DeleteHostFavouriteByID(ctx context.Context, id string, h util.Headers) (resp *api.BKAPIRsp, err error)
    IncrHostFavouritesCount(ctx context.Context, id string, h util.Headers) (resp *api.BKAPIRsp, err error)
    AddHistory(ctx context.Context, h util.Headers, dat map[string]interface{}) (resp *api.BKAPIRsp, err error)
    GetHistorys(ctx context.Context, start string, limit string, h util.Headers) (resp *api.BKAPIRsp, err error)
    AddHostMutiltAppModuleRelation(ctx context.Context, h util.Headers, dat *hosts.CloudHostModuleParams) (resp *api.BKAPIRsp, err error)
    HostModuleRelation(ctx context.Context, h util.Headers) (resp *api.BKAPIRsp, err error)
    MoveHost2EmptyModule(ctx context.Context, h util.Headers, dat *hosts.DefaultModuleHostConfigParams) (resp *api.BKAPIRsp, err error)
    MoveHost2FaultModule(ctx context.Context, h util.Headers, dat *hosts.DefaultModuleHostConfigParams) (resp *api.BKAPIRsp, err error)
    MoveHostToResourcePool(ctx context.Context, h util.Headers, dat *hosts.DefaultModuleHostConfigParams) (resp *api.BKAPIRsp, err error)
    AssignHostToApp(ctx context.Context, h util.Headers, dat *hosts.DefaultModuleHostConfigParams) (resp *api.BKAPIRsp, err error)
    AssignHostToAppModule(ctx context.Context, h util.Headers, dat interface{}) (resp *api.BKAPIRsp, err error)
    SaveUserCustom(ctx context.Context, h util.Headers, dat interface{}) (resp *api.BKAPIRsp, err error)
    GetUserCustom(ctx context.Context, h util.Headers) (resp *api.BKAPIRsp, err error)
    GetDefaultCustom(ctx context.Context, h util.Headers) (resp *api.BKAPIRsp, err error)
    GetAgentStatus(ctx context.Context, appID string, h util.Headers) (resp *api.BKAPIRsp, err error)
    UpdateHost(ctx context.Context, appID string, h util.Headers, dat interface{}) (resp *api.BKAPIRsp, err error)
    UpdateHostByAppID(ctx context.Context, appID string, h util.Headers, dat interface{}) (resp *api.BKAPIRsp, err error)
    GetHostListByAppidAndField(ctx context.Context, appID string, field string, h util.Headers) (resp *api.BKAPIRsp, err error)
    HostSearchByIP(ctx context.Context, h util.Headers, dat interface{}) (resp *api.BKAPIRsp, err error)
    HostSearchByModuleID(ctx context.Context, h util.Headers, dat interface{}) (resp *api.BKAPIRsp, err error)
    HostSearchBySetID(ctx context.Context, h util.Headers, dat interface{}) (resp *api.BKAPIRsp, err error)
    HostSearchByAppID(ctx context.Context, h util.Headers, dat interface{}) (resp *api.BKAPIRsp, err error)
    HostSearchByProperty(ctx context.Context, h util.Headers, dat map[string]interface{}) (resp *api.BKAPIRsp, err error)
    GetIPAndProxyByCompany(ctx context.Context, h util.Headers, dat map[string]interface{}) (resp *api.BKAPIRsp, err error)
    UpdateCustomProperty(ctx context.Context, h util.Headers, dat map[string]interface{}) (resp *api.BKAPIRsp, err error)
    CloneHostProperty(ctx context.Context, h util.Headers) (resp *api.BKAPIRsp, err error)
    GetHostAppByCompanyId(ctx context.Context, h util.Headers, dat map[string]interface{}) (resp *api.BKAPIRsp, err error)
    DelHostInApp(ctx context.Context, h util.Headers, dat interface{}) (resp *api.BKAPIRsp, err error)
    GetGitServerIp(ctx context.Context, h util.Headers, dat interface{}) (resp *api.BKAPIRsp, err error)
    GetPlat(ctx context.Context, h util.Headers) (resp *api.BKAPIRsp, err error)
    CreatePlat(ctx context.Context, h util.Headers, dat map[string]interface{}) (resp *api.BKAPIRsp, err error)
    DelPlat(ctx context.Context, cloudID string, h util.Headers) (resp *api.BKAPIRsp, err error)
    HostSearch(ctx context.Context, h util.Headers, dat *params.HostCommonSearch) (resp *api.BKAPIRsp, err error)
    HostSearchWithAsstDetail(ctx context.Context, h util.Headers, dat *params.HostCommonSearch) (resp *api.BKAPIRsp, err error)
    UpdateHostBatch(ctx context.Context, h util.Headers, dat interface{}) (resp *api.BKAPIRsp, err error)
    
    AddUserCustomQuery(ctx context.Context, h util.Headers, dat map[string]interface{}) (resp *api.BKAPIRsp, err error)
    UpdateUserCustomQuery(ctx context.Context, businessID string, id string, h util.Headers, dat map[string]interface{}) (resp *api.BKAPIRsp, err error)
    DeleteUserCustomQuery(ctx context.Context, businessID string, id string, h util.Headers) (resp *api.BKAPIRsp, err error)
    GetUserCustomQuery(ctx context.Context, businessID string, h util.Headers, dat *commondata.ObjQueryInput) (resp *api.BKAPIRsp, err error)
    GetUserCustomQueryDetail(ctx context.Context, businessID string, id string, h util.Headers) (resp *api.BKAPIRsp, err error)
    GetUserCustomQueryResult(ctx context.Context, businessID, id, start, limit string, h util.Headers) (resp *api.BKAPIRsp, err error)
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

