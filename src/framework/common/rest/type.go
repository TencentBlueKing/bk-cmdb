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

package rest

import (
	"net/http"

	"configcenter/src/apimachinery/flowctrl"
	"configcenter/src/framework/clientset/discovery"
	chttp "configcenter/src/framework/common/http"
)

// http request verb type
type VerbType string

const (
	PUT    VerbType = http.MethodPut
	POST   VerbType = http.MethodPost
	GET    VerbType = http.MethodGet
	DELETE VerbType = http.MethodDelete
	PATCH  VerbType = http.MethodPatch
)

type Capability struct {
	Client   chttp.HttpClient
	Discover discovery.Interface
	Throttle flowctrl.RateLimiter
	Mock     MockInfo
}

type MockInfo struct {
	Mocked      bool
	SetMockData bool
	MockData    interface{}
}
