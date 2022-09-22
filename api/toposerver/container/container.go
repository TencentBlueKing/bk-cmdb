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

	types2 "configcenter/pkg/kube/types"

	"configcenter/pkg/blog"
	"configcenter/pkg/errors"
	"configcenter/pkg/metadata"
)

// BatchCreateNode 批量创建node
func (st *Container) BatchCreateNode(ctx context.Context, header http.Header, bizID int64,
	data *types2.CreateNodesOption) (*types2.CreateNodesResult, errors.CCErrorCoder) {
	ret := new(types2.CreateNodesResult)
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
	return ret, nil
}

// BatchCreatePod batch create pod.
func (st *Container) BatchCreatePod(ctx context.Context, header http.Header, bizID int64,
	data *types2.CreatePodsOption) ([]types2.Pod, errors.CCErrorCoder) {
	ret := new(types2.CreatePodsResult)
	subPath := "/kube/createmany/pod/%d/instance"

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
	*types2.ResponseCluster, errors.CCErrorCoder) {
	// ret := new(table.ResponseCluster)
	ret := struct {
		metadata.BaseResp
		Data types2.ResponseCluster `json:"data"`
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
	return &ret.Data, nil
}

// SearchNode search node.
func (st *Container) SearchNode(ctx context.Context, header http.Header, input *metadata.QueryCondition) (
	*types2.SearchNodeRsp, errors.CCErrorCoder) {
	ret := struct {
		metadata.BaseResp
		Data types2.SearchNodeRsp `json:"data"`
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
	return &ret.Data, nil
}

// UpdateNodeFields update node fields.
func (st *Container) UpdateNodeFields(ctx context.Context, header http.Header, supplierAccount string, bizID int64,
	data *types2.UpdateNodeOption) errors.CCErrorCoder {
	ret := new(metadata.UpdatedCount)
	subPath := "/kube/updatemany/node/%s/%d/instance"
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
	data *types2.UpdateClusterOption) errors.CCErrorCoder {
	ret := new(metadata.UpdatedCount)
	subPath := "/kube/updatemany/cluster/%s/%d/instance"
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
	data *types2.ClusterBaseFields) (*types2.CreateClusterResult, errors.CCErrorCoder) {
	ret := new(types2.CreateClusterResult)
	subPath := "/kube/create/cluster/%d/instance"
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
	option *types2.DeleteClusterOption) errors.CCErrorCoder {
	ret := new(metadata.DeletedCount)
	subPath := "/kube/delete/cluster/%d/instance"

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
	option *types2.BatchDeleteNodeOption) errors.CCErrorCoder {
	ret := new(metadata.DeletedCount)
	subPath := "/kube/deletemany/node/%d/instance"

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
