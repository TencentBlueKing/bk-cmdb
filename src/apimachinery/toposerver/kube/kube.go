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

package kube

import (
	"context"
	"net/http"

	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/metadata"
	"configcenter/src/kube/types"
)

// BatchCreateNode batch create node
func (st *Kube) BatchCreateNode(ctx context.Context, header http.Header, bizID int64,
	data *types.CreateNodesOption) (*metadata.Response, errors.CCErrorCoder) {
	ret := new(metadata.Response)
	subPath := "/createmany/kube/node/bk_biz_id/%d"

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
func (st *Kube) BatchCreatePod(ctx context.Context, header http.Header,
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

// SearchCluster search cluster.
func (st *Kube) SearchCluster(ctx context.Context, header http.Header, bizID int64, input *types.QueryClusterOption) (
	*metadata.Response, errors.CCErrorCoder) {
	ret := new(metadata.Response)

	subPath := "/findmany/kube/cluster/bk_biz_id/%d"
	err := st.client.Post().
		WithContext(ctx).
		Body(input).
		SubResourcef(subPath, bizID).
		WithHeaders(header).
		Do().
		Into(&ret)

	if err != nil {
		return nil, errors.CCHttpError
	}
	if ret.CCError() != nil {
		return nil, ret.CCError()
	}
	return ret, nil
}

// SearchNode search node.
func (st *Kube) SearchNode(ctx context.Context, header http.Header, bizID int64, input *types.QueryNodeOption) (
	*metadata.Response, errors.CCErrorCoder) {
	ret := new(metadata.Response)

	subPath := "/findmany/kube/node/bk_biz_id/%d"
	err := st.client.Post().
		WithContext(ctx).
		Body(input).
		SubResourcef(subPath, bizID).
		WithHeaders(header).
		Do().
		Into(&ret)

	if err != nil {
		return nil, errors.CCHttpError
	}
	if ret.CCError() != nil {
		return nil, ret.CCError()
	}
	return ret, nil
}

// UpdateNodeFields update node fields.
func (st *Kube) UpdateNodeFields(ctx context.Context, header http.Header, bizID int64,
	data *types.UpdateNodeOption) (*metadata.Response, errors.CCErrorCoder) {
	ret := new(metadata.Response)
	subPath := "/updatemany/kube/node/bk_biz_id/%d"
	err := st.client.Put().
		WithContext(ctx).
		Body(data).
		SubResourcef(subPath, bizID).
		WithHeaders(header).
		Do().
		Into(ret)

	if err != nil {
		return nil, errors.CCHttpError
	}
	if ret.CCError() != nil {
		return nil, ret.CCError()
	}
	return ret, nil
}

// UpdateClusterFields update cluster fields.
func (st *Kube) UpdateClusterFields(ctx context.Context, header http.Header, bizID int64,
	data *types.UpdateClusterOption) (*metadata.Response, errors.CCErrorCoder) {
	ret := new(metadata.Response)
	subPath := "/updatemany/kube/cluster/bk_biz_id/%d"
	err := st.client.Put().
		WithContext(ctx).
		Body(data).
		SubResourcef(subPath, bizID).
		WithHeaders(header).
		Do().
		Into(ret)

	if err != nil {
		return nil, errors.CCHttpError
	}
	if ret.CCError() != nil {
		return nil, ret.CCError()
	}
	return ret, nil
}

// CreateCluster create cluster.
func (st *Kube) CreateCluster(ctx context.Context, header http.Header, bizID int64,
	data *types.Cluster) (*metadata.Response, errors.CCErrorCoder) {
	ret := new(metadata.Response)
	subPath := "/create/kube/cluster/bk_biz_id/%d"
	err := st.client.Post().
		WithContext(ctx).
		Body(data).
		SubResourcef(subPath, bizID).
		WithHeaders(header).
		Do().
		Into(ret)

	if err != nil {
		return nil, errors.CCHttpError
	}
	if ret.CCError() != nil {
		return nil, ret.CCError()
	}

	return ret, nil
}

// DeleteCluster delete cluster.
func (st *Kube) DeleteCluster(ctx context.Context, header http.Header, bizID int64,
	option *types.DeleteClusterOption) (*metadata.Response, errors.CCErrorCoder) {
	ret := new(metadata.Response)
	subPath := "/delete/kube/cluster/bk_biz_id/%d"

	err := st.client.Delete().
		WithContext(ctx).
		Body(option).
		SubResourcef(subPath, bizID).
		WithHeaders(header).
		Do().
		Into(ret)

	if err != nil {
		return nil, errors.CCHttpError
	}

	if ret.CCError() != nil {
		return nil, ret.CCError()
	}

	return ret, nil
}

// BatchDeleteNode delete node.
func (st *Kube) BatchDeleteNode(ctx context.Context, header http.Header, bizID int64,
	option *types.BatchDeleteNodeOption) (*metadata.Response, errors.CCErrorCoder) {
	ret := new(metadata.Response)
	subPath := "/deletemany/kube/node/bk_biz_id/%d"

	err := st.client.Delete().
		WithContext(ctx).
		Body(option).
		SubResourcef(subPath, bizID).
		WithHeaders(header).
		Do().
		Into(ret)

	if err != nil {
		return nil, errors.CCHttpError
	}

	if ret.CCError() != nil {
		return nil, ret.CCError()
	}

	return ret, nil
}
