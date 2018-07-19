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

package inst

import (
	"context"
	"fmt"
	"net/http"

	"configcenter/src/common/metadata"
)

// TODO: config this body data struct.
func (t *instanceClient) CreateInst(ctx context.Context, ownerID string, objID string, h http.Header, dat interface{}) (resp *metadata.CreateInstResult, err error) {
	resp = new(metadata.CreateInstResult)
	subPath := fmt.Sprintf("/inst/%s/%s", ownerID, objID)

	err = t.client.Post().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (t *instanceClient) DeleteInst(ctx context.Context, ownerID string, objID string, instID int64, h http.Header) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := fmt.Sprintf("/inst/%s/%s/%d", ownerID, objID, instID)

	err = t.client.Delete().
		WithContext(ctx).
		Body(nil).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (t *instanceClient) UpdateInst(ctx context.Context, ownerID string, objID string, instID int64, h http.Header, dat map[string]interface{}) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := fmt.Sprintf("/inst/%s/%s/%d", ownerID, objID, instID)

	err = t.client.Put().
		WithContext(ctx).
		Body(dat).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (t *instanceClient) SelectInsts(ctx context.Context, ownerID string, objID string, h http.Header, s *metadata.SearchParams) (resp *metadata.SearchInstResult, err error) {
	resp = new(metadata.SearchInstResult)
	subPath := fmt.Sprintf("/inst/search/%s/%s", ownerID, objID)

	err = t.client.Post().
		WithContext(ctx).
		Body(s).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (t *instanceClient) SelectInstsAndAsstDetail(ctx context.Context, ownerID string, objID string, h http.Header, s *metadata.SearchParams) (resp *metadata.SearchInstResult, err error) {
	resp = new(metadata.SearchInstResult)
	subPath := fmt.Sprintf("/inst/search/owner/%s/object/%s/detail", ownerID, objID)

	err = t.client.Post().
		WithContext(ctx).
		Body(s).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (t *instanceClient) InstSearch(ctx context.Context, ownerID string, objID string, h http.Header, s *metadata.SearchParams) (resp *metadata.SearchInstResult, err error) {
	resp = new(metadata.SearchInstResult)
	subPath := fmt.Sprintf("/inst/search/owner/%s/object/%s", ownerID, objID)

	err = t.client.Post().
		WithContext(ctx).
		Body(s).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (t *instanceClient) SelectInstsByAssociation(ctx context.Context, ownerID string, objID string, h http.Header, p *metadata.AssociationParams) (resp *metadata.SearchInstResult, err error) {
	resp = new(metadata.SearchInstResult)
	subPath := fmt.Sprintf("/inst/association/search/owner/%s/object/%s", ownerID, objID)

	err = t.client.Post().
		WithContext(ctx).
		Body(p).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (t *instanceClient) SelectInst(ctx context.Context, ownerID string, objID string, instID string, h http.Header, p *metadata.SearchParams) (resp *metadata.SearchInstResult, err error) {
	resp = new(metadata.SearchInstResult)
	subPath := fmt.Sprintf("/inst/search/owner/%s/%s/%s", ownerID, objID, instID)

	err = t.client.Post().
		WithContext(ctx).
		Body(p).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (t *instanceClient) SelectTopo(ctx context.Context, ownerID string, objID string, instID string, h http.Header, p *metadata.SearchParams) (resp *metadata.SearchInstResult, err error) {
	resp = new(metadata.SearchInstResult)
	subPath := fmt.Sprintf("/inst/search/topo/owner/%s/object/%s/inst/%s", ownerID, objID, instID)

	err = t.client.Post().
		WithContext(ctx).
		Body(p).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

func (t *instanceClient) SelectAssociationTopo(ctx context.Context, ownerID string, objID string, instID string, h http.Header, p *metadata.SearchParams) (resp *metadata.SearchInstResult, err error) {
	resp = new(metadata.SearchInstResult)
	subPath := fmt.Sprintf("/inst/association/topo/search/owner/%sobject/%s/inst/%s", ownerID, objID, instID)

	err = t.client.Post().
		WithContext(ctx).
		Body(p).
		SubResource(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}
