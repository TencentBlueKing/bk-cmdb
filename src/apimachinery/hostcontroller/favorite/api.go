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
	"net/http"

	"configcenter/src/apimachinery/rest"
	"configcenter/src/common/metadata"
)

type FavoriteInterface interface {
	AddHostFavourite(ctx context.Context, user string, h http.Header, dat *metadata.FavouriteParms) (resp *metadata.IDResult, err error)
	UpdateHostFavouriteByID(ctx context.Context, user string, id string, h http.Header, dat map[string]interface{}) (resp *metadata.BaseResp, err error)
	DeleteHostFavouriteByID(ctx context.Context, user string, id string, h http.Header) (resp *metadata.BaseResp, err error)
	GetHostFavourites(ctx context.Context, user string, h http.Header, dat *metadata.QueryInput) (resp *metadata.GetHostFavoriteResult, err error)
	GetHostFavouriteByID(ctx context.Context, user string, id string, h http.Header) (resp *metadata.GetHostFavoriteWithIDResult, err error)
}

func NewFavoriteInterface(client rest.ClientInterface) FavoriteInterface {
	return &favorites{client: client}
}

type favorites struct {
	client rest.ClientInterface
}
