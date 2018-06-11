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

package object

import (
	"context"
	"fmt"
	"net/http"

	"configcenter/src/common/core/cc/api"
	obj "configcenter/src/scene_server/topo_server/actions/object"
)

func (t *object) CreateModel(ctx context.Context, h http.Header, model *obj.MainLineObject) (resp *api.BKAPIRsp, err error) {
	resp = new(api.BKAPIRsp)
	subPath := "/model/mainline"

	err = t.client.Post().
		WithContext(ctx).
		Body(model).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}
func (t *object) DeleteModel(ctx context.Context, ownerID string, objID string, h http.Header) (resp *api.BKAPIRsp, err error) {
	resp = new(api.BKAPIRsp)
	subPath := fmt.Sprintf("/model/mainline/owners/%s/objectids/%s", ownerID, objID)

	err = t.client.Delete().
		WithContext(ctx).
		Body(nil).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}
func (t *object) SelectModel(ctx context.Context, ownerID string, h http.Header) (resp *api.BKAPIRsp, err error) {
	resp = new(api.BKAPIRsp)
	subPath := fmt.Sprintf("/model/%s", ownerID)

	err = t.client.Get().
		WithContext(ctx).
		Body(nil).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}
func (t *object) SelectModelByClsID(ctx context.Context, ownerID string, clsID string, objID string, h http.Header) (resp *api.BKAPIRsp, err error) {
	resp = new(api.BKAPIRsp)
	subPath := fmt.Sprintf("/model/%s/%s/%s", ownerID, clsID, objID)

	err = t.client.Get().
		WithContext(ctx).
		Body(nil).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}
func (t *object) SelectInst(ctx context.Context, ownerID string, appID string, h http.Header) (resp *api.BKAPIRsp, err error) {
	resp = new(api.BKAPIRsp)
	subPath := fmt.Sprintf("/inst/%s/%s", ownerID, appID)

	err = t.client.Get().
		WithContext(ctx).
		Body(nil).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}
func (t *object) SelectInstChild(ctx context.Context, ownerID string, objID string, appID string, instID string, h http.Header) (resp *api.BKAPIRsp, err error) {
	resp = new(api.BKAPIRsp)
	subPath := fmt.Sprintf("/inst/child/%s/%s/%s/%s", ownerID, objID, appID, instID)

	err = t.client.Get().
		WithContext(ctx).
		Body(nil).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}
