/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 THL A29 Limited,
 * a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 * We undertake not to change the open source license (MIT license) applicable
 * to the current version of the project delivered to anyone in the future.
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

// BatchCreateNode batch create nodes
func (st *Container) BatchCreateNode(ctx context.Context, header http.Header, bizID int64,
	data *types.CreateNodesOption) (*types.CreateNodesResult, errors.CCErrorCoder) {
	ret := new(types.CreateNodesResult)
	subPath := "/createmany/kube/node/%d"

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
	return ret, nil
}

// BatchCreatePod batch create pod.
func (st *Container) BatchCreatePod(ctx context.Context, header http.Header,
	data *types.CreatePodsOption) ([]types.Pod, errors.CCErrorCoder) {
	ret := new(types.CreatePodsResult)
	subPath := "/createmany/kube/pod"

	err := st.client.Post().
		WithContext(ctx).
		Body(data).
		SubResourcef(subPath).
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
	ret := struct {
		metadata.BaseResp
		Data types.ResponseCluster `json:"data"`
	}{}

	subPath := "/findmany/kube/cluster"
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
	return &ret.Data, nil
}

// SearchNode search node.
func (st *Container) SearchNode(ctx context.Context, header http.Header, input *metadata.QueryCondition) (
	*types.SearchNodeRsp, errors.CCErrorCoder) {
	ret := struct {
		metadata.BaseResp
		Data types.SearchNodeRsp `json:"data"`
	}{}

	subPath := "/findmany/kube/node"
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
	return &ret.Data, nil
}

// UpdateNodeFields update node fields.
func (st *Container) UpdateNodeFields(ctx context.Context, header http.Header, supplierAccount string, bizID int64,
	data *types.UpdateNodeOption) errors.CCErrorCoder {
	ret := new(metadata.UpdatedCount)
	subPath := "/updatemany/kube/node/%s/%d"
	err := st.client.Put().
		WithContext(ctx).
		Body(data).
		SubResourcef(subPath, supplierAccount, bizID).
		WithHeaders(header).
		Do().
		Into(ret)

	if err != nil {
		blog.Errorf("update node field failed, err: %+v", err)
		return errors.CCHttpError
	}

	return nil
}

// UpdateClusterFields update cluster fields.
func (st *Container) UpdateClusterFields(ctx context.Context, header http.Header, supplierAccount string, bizID int64,
	data *types.UpdateClusterOption) errors.CCErrorCoder {
	ret := new(metadata.UpdatedCount)
	subPath := "/updatemany/kube/cluster/%s/%d"
	err := st.client.Put().
		WithContext(ctx).
		Body(data).
		SubResourcef(subPath, supplierAccount, bizID).
		WithHeaders(header).
		Do().
		Into(ret)

	if err != nil {
		blog.Errorf("update cluster field failed, err: %+v", err)
		return errors.CCHttpError
	}

	return nil
}

// CreateCluster create cluster.
func (st *Container) CreateCluster(ctx context.Context, header http.Header, bizID int64,
	data *types.Cluster) (*types.CreateClusterResult, errors.CCErrorCoder) {
	ret := new(types.CreateClusterResult)
	subPath := "/create/kube/cluster/%d"
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
	ret := new(metadata.DeletedCount)
	subPath := "/deletemany/kube/cluster/%d"

	err := st.client.Delete().
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

	return nil
}

// BatchDeleteNode delete cluster.
func (st *Container) BatchDeleteNode(ctx context.Context, header http.Header, bizID int64,
	option *types.BatchDeleteNodeOption) errors.CCErrorCoder {
	ret := new(metadata.DeletedCount)
	subPath := "/deletemany/kube/node/%d"

	err := st.client.Delete().
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

	return nil
}
