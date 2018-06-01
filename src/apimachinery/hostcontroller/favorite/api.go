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

package favorite

import (
    "context"
    
    "configcenter/src/apimachinery/rest"
    "configcenter/src/apimachinery/util"
    "configcenter/src/common/core/cc/api"
    "configcenter/src/source_controller/common/commondata"
)

type FavoriteInterface interface {
    AddHostFavourite(ctx context.Context, user string, h util.Headers, dat map[string]interface{}) (resp *api.BKAPIRsp, err error)
    UpdateHostFavouriteByID(ctx context.Context, user string, id string, h util.Headers, dat map[string]interface{}) (resp *api.BKAPIRsp, err error)
    DeleteHostFavouriteByID(ctx context.Context, user string, id string,h util.Headers) (resp *api.BKAPIRsp, err error)
    GetHostFavourites(ctx context.Context, user string, h util.Headers, dat commondata.ObjQueryInput) (resp *api.BKAPIRsp, err error)
    GetHostFavouriteByID(ctx context.Context, user string, id string,h util.Headers) (resp *api.BKAPIRsp, err error)
}

func NewFavoriteInterface(client rest.ClientInterface) FavoriteInterface {
    return &favorites{client:client}
}

type favorites struct {
    client rest.ClientInterface
}
