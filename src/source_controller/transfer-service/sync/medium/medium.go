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

// Package medium defines transfer medium client
package medium

import (
	"context"
	"net/http"

	"configcenter/pkg/synchronize/types"
	"configcenter/src/apimachinery/flowctrl"
	"configcenter/src/apimachinery/rest"
	"configcenter/src/apimachinery/util"
	"configcenter/src/common/blog"

	"github.com/prometheus/client_golang/prometheus"
)

// ClientI defines transfer medium client interface
type ClientI interface {
	PushSyncData(ctx context.Context, h http.Header, opt *types.PushSyncDataOpt) error
	PullSyncData(ctx context.Context, h http.Header, opt *types.PullSyncDataOpt) (*types.PullSyncDataRes, error)
}

// NewTransferMedium new transfer medium client
func NewTransferMedium(addr []string, reg prometheus.Registerer) (ClientI, error) {
	client, err := util.NewClient(nil)
	if err != nil {
		blog.Errorf("new http client failed, err: %v", err)
		return nil, err
	}

	c := &util.Capability{
		Client:     client,
		Discover:   &discovery{servers: addr},
		Throttle:   flowctrl.NewRateLimiter(500, 500),
		MetricOpts: util.MetricOption{Register: reg},
	}

	restCli := rest.NewRESTClient(c, "/")

	return &transMediumCli{client: restCli}, nil
}

// transMediumCli defines transfer medium client
type transMediumCli struct {
	client rest.ClientInterface
}
