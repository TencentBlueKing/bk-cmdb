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
package unique

import (
	"context"
	"fmt"
	"net/http"

	"configcenter/src/apimachinery/rest"
	"configcenter/src/common/metadata"
)

func NewUniqueInterface(client rest.ClientInterface) UniqueInterface {
	return &Unique{client: client}
}

type Unique struct {
	client rest.ClientInterface
}

func (unique *Unique) Search(ctx context.Context, h http.Header, objectID string) (resp *metadata.SearchUniqueResult, err error) {
	resp = new(metadata.SearchUniqueResult)
	subPath := fmt.Sprintf("/object/%s/unique/action/search", objectID)

	err = unique.client.Get().
		WithContext(ctx).
		Body(nil).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (unique *Unique) Create(ctx context.Context, h http.Header, objectID string, request *metadata.CreateUniqueRequest) (resp *metadata.CreateUniqueResult, err error) {
	resp = new(metadata.CreateUniqueResult)
	subPath := fmt.Sprintf("/object/%s/unique/action/create", objectID)

	err = unique.client.Post().
		WithContext(ctx).
		Body(request).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}
func (unique *Unique) Update(ctx context.Context, h http.Header, objectID string, id uint64, request *metadata.UpdateUniqueRequest) (resp *metadata.UpdateUniqueResult, err error) {
	resp = new(metadata.UpdateUniqueResult)
	subPath := fmt.Sprintf("/object/%s/unique/%d/action/update", objectID, id)

	err = unique.client.Put().
		WithContext(ctx).
		Body(request).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}
func (unique *Unique) Delete(ctx context.Context, h http.Header, objectID string, id uint64) (resp *metadata.DeleteUniqueResult, err error) {
	resp = new(metadata.DeleteUniqueResult)
	subPath := fmt.Sprintf("/object/%s/unique/%d/action/delete", objectID, id)

	err = unique.client.Delete().
		WithContext(ctx).
		Body(nil).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}
