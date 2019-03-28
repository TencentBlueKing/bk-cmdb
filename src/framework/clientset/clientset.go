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

package clientset

import (
	"fmt"

	"configcenter/src/apimachinery/flowctrl"
	"configcenter/src/framework/clientset/business"
	"configcenter/src/framework/clientset/discovery"
	"configcenter/src/framework/clientset/host"
	"configcenter/src/framework/clientset/model"
	"configcenter/src/framework/common/http"
	"configcenter/src/framework/common/rest"
)

type ClientConfig struct {
	// comma separated.
	ZkAddr string
	TLS    http.TLSConfig
}

func NewV3Client(c ClientConfig) (V3ClientSet, error) {

	disc, err := discovery.DiscoveryAPIServer(c.ZkAddr)
	if err != nil {
		return nil, fmt.Errorf("service discovery api failed, err: %v", err)
	}

	cli, err := http.NewClient(&c.TLS)
	if err != nil {
		return nil, fmt.Errorf("new http client failed, err: %v", err)
	}

	cap := &rest.Capability{
		Discover: disc,
		Client:   cli,
		Throttle: flowctrl.NewRateLimiter(1000, 2000),
	}

	return &v3{
		client: rest.NewRESTClient(cap, "/api/v3"),
	}, nil
}

type V3ClientSet interface {
	Business() business.Interface
	Host() host.Interface
	Model() model.Interface
}

type v3 struct {
	client rest.ClientInterface
}

func (v *v3) Business() business.Interface {
	return business.NewBusinessClient(v.client)
}

func (v *v3) Host() host.Interface {
	return host.NewHostCtrl(v.client)
}

func (v *v3) Model() model.Interface {
	return model.NewModelClient(v.client)
}
