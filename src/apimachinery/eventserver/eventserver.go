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

package eventserver

import (
	"context"
	"fmt"
	"net/http"

	"configcenter/src/apimachinery/rest"
	"configcenter/src/apimachinery/util"
	"configcenter/src/common/metadata"
)

type EventServerClientInterface interface {
	Query(ctx context.Context, ownerID string, appID string, h http.Header, dat metadata.ParamSubscriptionSearch) (resp *metadata.Response, err error)
	Ping(ctx context.Context, h http.Header, dat interface{}) (resp *metadata.Response, err error)
	Telnet(ctx context.Context, h http.Header, dat interface{}) (resp *metadata.Response, err error)
	Subscribe(ctx context.Context, ownerID string, appID string, h http.Header, subscription *metadata.Subscription) (resp *metadata.Response, err error)
	UnSubscribe(ctx context.Context, ownerID string, appID string, subscribeID string, h http.Header) (resp *metadata.Response, err error)
	Rebook(ctx context.Context, ownerID string, appID string, subscribeID string, h http.Header, subscription *metadata.Subscription) (resp *metadata.Response, err error)
}

func NewEventServerClientInterface(c *util.Capability, version string) EventServerClientInterface {
	base := fmt.Sprintf("/event/%s", version)

	return &eventServer{
		client: rest.NewRESTClient(c, base),
	}
}

type eventServer struct {
	client rest.ClientInterface
}
