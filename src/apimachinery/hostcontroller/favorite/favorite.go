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
	"fmt"
	"net/http"

	"configcenter/src/common/metadata"
)

func (f *favorites) AddHostFavourite(ctx context.Context, user string, h http.Header, dat *metadata.FavouriteParms) (resp *metadata.IDResult, err error) {
	resp = new(metadata.IDResult)
	subPath := fmt.Sprintf("/hosts/favorites/%s", user)

	err = f.client.Post().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (f *favorites) UpdateHostFavouriteByID(ctx context.Context, user string, id string, h http.Header, dat map[string]interface{}) (resp *metadata.BaseResp, err error) {
	resp = new(metadata.BaseResp)
	subPath := fmt.Sprintf("/hosts/favorites/%s/%s", user, id)

	err = f.client.Put().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (f *favorites) DeleteHostFavouriteByID(ctx context.Context, user string, id string, h http.Header) (resp *metadata.BaseResp, err error) {
	resp = new(metadata.BaseResp)
	subPath := fmt.Sprintf("/hosts/favorites/%s/%s", user, id)

	err = f.client.Delete().
		WithContext(ctx).
		Body(nil).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (f *favorites) GetHostFavourites(ctx context.Context, user string, h http.Header, dat *metadata.QueryInput) (resp *metadata.GetHostFavoriteResult, err error) {
	resp = new(metadata.GetHostFavoriteResult)
	subPath := fmt.Sprintf("/hosts/favorites/search/%s", user)

	err = f.client.Post().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (f *favorites) GetHostFavouriteByID(ctx context.Context, user string, id string, h http.Header) (resp *metadata.GetHostFavoriteWithIDResult, err error) {
	resp = new(metadata.GetHostFavoriteWithIDResult)
	subPath := fmt.Sprintf("/hosts/favorites/search/%s/%s", user, id)

	err = f.client.Post().
		WithContext(ctx).
		Body(nil).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}
