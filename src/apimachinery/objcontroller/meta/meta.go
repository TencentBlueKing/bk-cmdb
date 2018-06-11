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

package meta

import (
    "context"
    "fmt"
    "net/http"

    "configcenter/src/common/core/cc/api"
    "configcenter/src/source_controller/api/metadata"
)

func (t *meta) SelectObjects(ctx context.Context, h http.Header, dat interface{}) (resp *api.BKAPIRsp, err error) {
    resp = new(api.BKAPIRsp)
    subPath := "/meta/objects"

    err = t.client.Post().
        WithContext(ctx).
        Body(dat).
        SubResource(subPath).
        WithHeaders(h).
        Do().
        Into(resp)
    return
}

func (t *meta) DeleteObject(ctx context.Context, objID string, h http.Header, dat map[string]interface{}) (resp *api.BKAPIRsp, err error) {
    resp = new(api.BKAPIRsp)
    subPath := fmt.Sprintf("/meta/object/%s", objID)

    err = t.client.Delete().
        WithContext(ctx).
        Body(dat).
        SubResource(subPath).
        WithHeaders(h).
        Do().
        Into(resp)
    return
}

func (t *meta) CreateObject(ctx context.Context, h http.Header, dat *metadata.ObjectAttDes) (resp *api.BKAPIRsp, err error) {
    resp = new(api.BKAPIRsp)
    subPath := "/meta/object"

    err = t.client.Post().
        WithContext(ctx).
        Body(dat).
        SubResource(subPath).
        WithHeaders(h).
        Do().
        Into(resp)
    return
}

func (t *meta) UpdateObject(ctx context.Context, objID string, h http.Header, dat interface{}) (resp *api.BKAPIRsp, err error) {
    resp = new(api.BKAPIRsp)
    subPath := fmt.Sprintf("/meta/object/%s", objID)

    err = t.client.Put().
        WithContext(ctx).
        Body(dat).
        SubResource(subPath).
        WithHeaders(h).
        Do().
        Into(resp)
    return
}

func (t *meta) SelectObjectAssociations(ctx context.Context, h http.Header, dat map[string]interface{}) (resp *api.BKAPIRsp, err error) {
    resp = new(api.BKAPIRsp)
    subPath := "/meta/objectassts"

    err = t.client.Post().
        WithContext(ctx).
        Body(dat).
        SubResource(subPath).
        WithHeaders(h).
        Do().
        Into(resp)
    return
}

func (t *meta) DeleteObjectAssociation(ctx context.Context, objID string, h http.Header, dat map[string]interface{}) (resp *api.BKAPIRsp, err error) {
    resp = new(api.BKAPIRsp)
    subPath := fmt.Sprintf("/meta/objectasst/%s", objID)

    err = t.client.Delete().
        WithContext(ctx).
        Body(dat).
        SubResource(subPath).
        WithHeaders(h).
        Do().
        Into(resp)
    return
}

func (t *meta) CreateObjectAssociation(ctx context.Context, h http.Header, dat *metadata.ObjectAsst) (resp *api.BKAPIRsp, err error) {
    resp = new(api.BKAPIRsp)
    subPath := "/meta/objectasst"

    err = t.client.Post().
        WithContext(ctx).
        Body(dat).
        SubResource(subPath).
        WithHeaders(h).
        Do().
        Into(resp)
    return
}

func (t *meta) UpdateObjectAssociation(ctx context.Context, objID string, h http.Header, dat map[string]interface{}) (resp *api.BKAPIRsp, err error) {
    resp = new(api.BKAPIRsp)
    subPath := fmt.Sprintf("/meta/objectasst/%s", objID)

    err = t.client.Put().
        WithContext(ctx).
        Body(dat).
        SubResource(subPath).
        WithHeaders(h).
        Do().
        Into(resp)
    return
}

func (t *meta) SelectObjectAttByID(ctx context.Context, objID string, h http.Header) (resp *api.BKAPIRsp, err error) {
    resp = new(api.BKAPIRsp)
    subPath := fmt.Sprintf("/meta/objectatt/%s", objID)

    err = t.client.Post().
        WithContext(ctx).
        Body(nil).
        SubResource(subPath).
        WithHeaders(h).
        Do().
        Into(resp)
    return
}

func (t *meta) SelectObjectAttWithParams(ctx context.Context, h http.Header, dat interface{}) (resp *api.BKAPIRsp, err error) {
    resp = new(api.BKAPIRsp)
    subPath := "/meta/objectatts"

    err = t.client.Post().
        WithContext(ctx).
        Body(dat).
        SubResource(subPath).
        WithHeaders(h).
        Do().
        Into(resp)
    return
}

func (t *meta) DeleteObjectAttByID(ctx context.Context, objID string, h http.Header, dat interface{}) (resp *api.BKAPIRsp, err error) {
    resp = new(api.BKAPIRsp)
    subPath := fmt.Sprintf("/meta/objectatt/%s", objID)

    err = t.client.Delete().
        WithContext(ctx).
        Body(dat).
        SubResource(subPath).
        WithHeaders(h).
        Do().
        Into(resp)
    return
}

func (t *meta) CreateObjectAtt(ctx context.Context, h http.Header, dat *metadata.ObjectAttDes) (resp *api.BKAPIRsp, err error) {
    resp = new(api.BKAPIRsp)
    subPath := "/meta/objectatt"

    err = t.client.Post().
        WithContext(ctx).
        Body(dat).
        SubResource(subPath).
        WithHeaders(h).
        Do().
        Into(resp)
    return
}

func (t *meta) UpdateObjectAttByID(ctx context.Context, objID string, h http.Header, dat map[string]interface{}) (resp *api.BKAPIRsp, err error) {
    resp = new(api.BKAPIRsp)
    subPath := fmt.Sprintf("/meta/objectatt/%s", objID)

    err = t.client.Put().
        WithContext(ctx).
        Body(dat).
        SubResource(subPath).
        WithHeaders(h).
        Do().
        Into(resp)
    return
}