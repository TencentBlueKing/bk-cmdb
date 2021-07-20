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
	HostSnapInfoBatch(ctx context.Context, h http.Header, dat interface{}) (resp *metadata.HostSnapBatchResult, err error)
	AddHost(ctx context.Context, h http.Header, dat interface{}) (resp *metadata.Response, err error)
	AddHostToResourcePool(ctx context.Context, h http.Header, dat metadata.AddHostToResourcePoolHostList) (resp *metadata.Response, err error)
	AddHostFromAgent(ctx context.Context, h http.Header, dat interface{}) (resp *metadata.Response, err error)
	SyncHost(ctx context.Context, h http.Header, data interface{}) (resp *metadata.Response, err error)

	GetHostFavourites(ctx context.Context, h http.Header, dat interface{}) (resp *metadata.GetHostFavoriteResult, err error)
	AddHostFavourite(ctx context.Context, h http.Header, dat *metadata.FavouriteParms) (resp *metadata.Response, err error)
	UpdateHostFavouriteByID(ctx context.Context, id string, h http.Header, data *metadata.FavouriteParms) (resp *metadata.Response, err error)
	DeleteHostFavouriteByID(ctx context.Context, id string, h http.Header) (resp *metadata.Response, err error)
	IncrHostFavouritesCount(ctx context.Context, id string, h http.Header) (resp *metadata.Response, err error)

	AddHistory(ctx context.Context, h http.Header, dat map[string]interface{}) (resp *metadata.Response, err error)
	GetHistorys(ctx context.Context, start string, limit string, h http.Header) (resp *metadata.Response, err error)

	AddHostMultiAppModuleRelation(ctx context.Context, h http.Header, dat *metadata.CloudHostModuleParams) (resp *metadata.Response, err error)
	TransferHostModule(ctx context.Context, h http.Header, params map[string]interface{}) (resp *metadata.Response, err error)

	MoveHost2EmptyModule(ctx context.Context, h http.Header, dat *metadata.DefaultModuleHostConfigParams) (resp *metadata.Response, err error)
	MoveHost2FaultModule(ctx context.Context, h http.Header, dat *metadata.DefaultModuleHostConfigParams) (resp *metadata.Response, err error)
	MoveHostToResourcePool(ctx context.Context, h http.Header, dat *metadata.DefaultModuleHostConfigParams) (resp *metadata.Response, err error)

	AssignHostToApp(ctx context.Context, h http.Header, dat *metadata.DefaultModuleHostConfigParams) (resp *metadata.Response, err error)
	SaveUserCustom(ctx context.Context, h http.Header, dat interface{}) (resp *metadata.Response, err error)
	GetUserCustom(ctx context.Context, h http.Header) (resp *metadata.Response, err error)
	GetDefaultCustom(ctx context.Context, h http.Header) (resp *metadata.Response, err error)
	CloneHostProperty(ctx context.Context, h http.Header, dat *metadata.CloneHostPropertyParams) (resp *metadata.Response, err error)
	MoveSetHost2IdleModule(ctx context.Context, h http.Header, dat *metadata.SetHostConfigParams) (resp *metadata.Response, err error)
	SearchHost(ctx context.Context, h http.Header, dat *params.HostCommonSearch) (resp *metadata.SearchHostResult, err error)
	SearchHostWithAsstDetail(ctx context.Context, h http.Header, dat *params.HostCommonSearch) (resp *metadata.SearchHostResult, err error)
	UpdateHostBatch(ctx context.Context, h http.Header, dat interface{}) (resp *metadata.Response, err error)
	UpdateHostPropertyBatch(ctx context.Context, h http.Header, data map[string]interface{}) (resp *metadata.Response, err error)

	// dynamic group interfaces.
	CreateDynamicGroup(ctx context.Context, header http.Header, data map[string]interface{}) (resp *metadata.IDResult, err error)
	UpdateDynamicGroup(ctx context.Context, bizID, id string, header http.Header, data map[string]interface{}) (resp *metadata.BaseResp, err error)
	DeleteDynamicGroup(ctx context.Context, bizID, id string, header http.Header) (resp *metadata.BaseResp, err error)
	GetDynamicGroup(ctx context.Context, bizID, id string, header http.Header) (resp *metadata.GetDynamicGroupResult, err error)
	SearchDynamicGroup(ctx context.Context, bizID string, header http.Header, data *metadata.QueryCondition) (resp *metadata.SearchDynamicGroupResult, err error)
	ExecuteDynamicGroup(ctx context.Context, bizID, id string, header http.Header, data map[string]interface{}) (resp *metadata.Response, err error)

	HostSearch(ctx context.Context, h http.Header, params *metadata.HostCommonSearch) (resp *metadata.QueryInstResult, err error)
	ListBizHostsTopo(ctx context.Context, h http.Header, bizID int64, params *metadata.ListHostsWithNoBizParameter) (resp *metadata.SuccessResponse, err error)

	CreateCloudArea(ctx context.Context, h http.Header, data map[string]interface{}) (resp *metadata.CreatedOneOptionResult, err error)
	CreateManyCloudArea(ctx context.Context, h http.Header, data map[string]interface{}) (resp *metadata.CreateManyCloudAreaResult, err error)
	UpdateCloudArea(ctx context.Context, h http.Header, cloudID int64, data map[string]interface{}) (resp *metadata.Response, err error)
	SearchCloudArea(ctx context.Context, h http.Header, params map[string]interface{}) (resp *metadata.SearchResp, err error)
	DeleteCloudArea(ctx context.Context, h http.Header, cloudID int64) (resp *metadata.Response, err error)
	FindCloudAreaHostCount(ctx context.Context, header http.Header, option metadata.CloudAreaHostCount) (resp *metadata.CloudAreaHostCountResult, err error)
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
