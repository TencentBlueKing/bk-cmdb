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

package medium

import (
	"context"
	"fmt"
	"net/http"

	"configcenter/pkg/synchronize/types"
)

// PushSyncData push sync data
func (t *transMediumCli) PushSyncData(ctx context.Context, h http.Header, opt *types.PushSyncDataOpt) error {
	resp := new(types.TransferMediumResp[any])

	err := t.client.Post().
		WithContext(ctx).
		Body(opt).
		SubResourcef("/api/sync/publish").
		WithHeaders(h).
		Do().
		Into(resp)

	if err != nil {
		return err
	}

	if resp.Code != 0 {
		return fmt.Errorf("code: %d, message: %s", resp.Code, resp.Message)
	}

	return nil
}

// PullSyncData pull sync data
func (t *transMediumCli) PullSyncData(ctx context.Context, h http.Header, opt *types.PullSyncDataOpt) (
	*types.PullSyncDataRes, error) {

	resp := new(types.TransferMediumResp[types.PullSyncDataRes])

	err := t.client.Post().
		WithContext(ctx).
		Body(opt).
		SubResourcef("/api/sync/consume").
		WithHeaders(h).
		Do().
		Into(resp)

	if err != nil {
		return nil, err
	}

	if resp.Code != 0 {
		return nil, fmt.Errorf("code: %d, message: %s", resp.Code, resp.Message)
	}

	return &resp.Data, nil
}
