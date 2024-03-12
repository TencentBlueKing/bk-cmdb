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

package notice

import (
	"context"
	"fmt"
	"net/http"

	"configcenter/src/thirdparty/apigw"
)

// GetCurAnn get current announcements
func (n *notice) GetCurAnn(ctx context.Context, h http.Header, params map[string]string) ([]CurAnnData, error) {
	h.Set(apigw.AuthKey, n.service.Auth)
	resp := new(GetCurAnnResp)
	subPath := "/prod/apigw/v1/announcement/get_current_announcements"

	err := n.service.Client.Get().
		WithContext(ctx).
		WithParams(params).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)

	if err != nil {
		return nil, err
	}

	if resp.Code != 0 {
		return nil, fmt.Errorf("code: %d, message: %s", resp.Code, resp.Message)
	}

	return resp.Data, nil
}

// RegApp register application
func (n *notice) RegApp(ctx context.Context, h http.Header) (*RegAppData, error) {
	h.Set(apigw.AuthKey, n.service.Auth)
	resp := new(RegAppResp)
	subPath := "/prod/apigw/v1/register"

	err := n.service.Client.Post().
		WithContext(ctx).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)

	if err != nil {
		return nil, err
	}

	if resp.Code != 0 {
		return nil, fmt.Errorf("code: %d, message: %s", resp.Code, resp.Message)
	}

	return resp.Data, nil
}
