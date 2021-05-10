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

package client

import (
	"context"
	"net/http"

	"configcenter/src/apimachinery/flowctrl"
	"configcenter/src/apimachinery/rest"
	"configcenter/src/apimachinery/util"
	"configcenter/src/scene_server/auth_server/sdk/operator"
	"configcenter/src/scene_server/auth_server/sdk/types"
)

type Interface interface {
	GetUserPolicy(ctx context.Context, opt *types.GetPolicyOption) (*operator.Policy, error)
	ListUserPolicies(ctx context.Context, opts *types.ListPolicyOptions) ([]*types.ActionPolicy, error)
	GetSystemToken(ctx context.Context) (string, error)
}

func NewClient(conf types.IamConfig, opt types.Options) (Interface, error) {

	if err := conf.Validate(); err != nil {
		return nil, err
	}

	client, err := util.NewClient(&conf.TLS)
	if err != nil {
		return nil, err
	}

	c := &util.Capability{
		Client: client,
		Discover: &acDiscovery{
			servers: conf.Address,
		},
		Throttle: flowctrl.NewRateLimiter(2000, 3000),
		Mock: util.MockInfo{
			Mocked: false,
		},
	}

	// add prometheus metric if possible.
	if opt.Metric != nil {
		c.Reg = opt.Metric
	}

	header := http.Header{}
	header.Set("Content-Type", "application/json")
	header.Set("Accept", "application/json")
	header.Set("X-Bk-App-Code", conf.AppCode)
	header.Set("X-Bk-App-Secret", conf.AppSecret)

	return &authClient{
		client:      rest.NewRESTClient(c, ""),
		basicHeader: header,
		config:      conf,
	}, nil
}

type authClient struct {
	// http client instance
	client rest.ClientInterface
	// http header info
	basicHeader http.Header
	// iam config
	config types.IamConfig
}

func (ac *authClient) cloneHeader(ctx context.Context) http.Header {
	h := http.Header{}
	rid, ok := ctx.Value(types.RequestIDKey).(string)
	if ok {
		h.Set(types.RequestIDHeaderKey, rid)
	}
	for key := range ac.basicHeader {
		h.Set(key, ac.basicHeader.Get(key))
	}
	return h
}
