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
package container

import (
	"context"
	"net/http"

	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/metadata"
	"configcenter/src/kube/types"
)

func (st *Container) BatchCreateNode(ctx context.Context, header http.Header, bizID int64,
	data *types.CreateNodesReq) ([]int64, errors.CCErrorCoder) {
	ret := new(types.CreateNodesResult)
	subPath := "/kube/createmany/node/%d/instance"

	err := st.client.Post().
		WithContext(ctx).
		Body(data).
		SubResourcef(subPath, bizID).
		WithHeaders(header).
		Do().
		Into(ret)

	if err != nil {
		blog.Errorf("batch create node failed, http request failed, err: %+v", err)
		return nil, errors.CCHttpError
	}
	if ret.CCError() != nil {
		return nil, ret.CCError()
	}

	return ret.Info, nil

}

// SearchCluster create cluster.
func (st *Container) SearchCluster(ctx context.Context, header http.Header, input *metadata.QueryCondition) (
	*types.ResponseCluster, errors.CCErrorCoder) {
	//ret := new(table.ResponseCluster)
	ret := struct {
		metadata.BaseResp
		Data types.ResponseCluster `json:"data"`
	}{}

	subPath := "/kube/search/cluster/instances"
	err := st.client.Post().
		WithContext(ctx).
		Body(input).
		SubResourcef(subPath).
		WithHeaders(header).
		Do().
		Into(&ret)

	if err != nil {
		blog.Errorf("search cluster failed, http request failed, err: %+v", err)
		return nil, errors.CCHttpError
	}
	blog.Errorf("5555555555555555555 ret: %+v", ret)
	return &ret.Data, nil
}

// SearchNode search node.
func (st *Container) SearchNode(ctx context.Context, header http.Header, input *metadata.QueryCondition) (
	*types.ResponseNode, errors.CCErrorCoder) {
	//ret := new(table.ResponseCluster)
	ret := struct {
		metadata.BaseResp
		Data types.ResponseNode `json:"data"`
	}{}

	subPath := "/kube/search/node/instances"
	err := st.client.Post().
		WithContext(ctx).
		Body(input).
		SubResourcef(subPath).
		WithHeaders(header).
		Do().
		Into(&ret)

	if err != nil {
		blog.Errorf("search node failed, http request failed, err: %+v", err)
		return nil, errors.CCHttpError
	}
	blog.Errorf("5555555555555555555 ret: %+v", ret)
	return &ret.Data, nil
}

// CreateCluster create cluster.
func (st *Container) CreateCluster(ctx context.Context, header http.Header, bizID int64,
	data *types.ClusterBaseFields) (*types.CreateClusterResult, errors.CCErrorCoder) {
	ret := new(types.CreateClusterResult)
	subPath := "/kube/create/cluster/%d/instance"
	blog.Errorf("0000000000000 name: %+v, uid: %+v", *data.Name, *data.Uid)
	err := st.client.Post().
		WithContext(ctx).
		Body(data).
		SubResourcef(subPath, bizID).
		WithHeaders(header).
		Do().
		Into(ret)

	if err != nil {
		blog.Errorf("create cluster failed, http request failed, err: %+v", err)
		return nil, errors.CCHttpError
	}
	if ret.CCError() != nil {
		return nil, ret.CCError()
	}

	return ret, nil
}

// DeleteCluster delete cluster.
func (st *Container) DeleteCluster(ctx context.Context, header http.Header, bizID int64,
	option *types.DeleteClusterOption) errors.CCErrorCoder {
	ret := new(types.CreateClusterResult)
	subPath := "/kube/delete/cluster/{bk_biz_id}/instance"

	err := st.client.Post().
		WithContext(ctx).
		Body(option).
		SubResourcef(subPath, bizID).
		WithHeaders(header).
		Do().
		Into(ret)

	if err != nil {
		blog.Errorf("delete cluster failed, http request failed, err: %v", err)
		return errors.CCHttpError
	}
	if ret.CCError() != nil {
		return ret.CCError()
	}

	return nil
}
