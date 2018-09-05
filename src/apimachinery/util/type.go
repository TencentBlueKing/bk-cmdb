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

package util

import (
	"net/http"
	"reflect"

	"configcenter/src/apimachinery/discovery"
	"configcenter/src/apimachinery/flowctrl"
)

type APIMachineryConfig struct {
	// the address of zookeeper address, comma separated.
	ZkAddr string
	// request's qps value
	QPS int64
	// request's burst value
	Burst     int64
	TLSConfig *TLSClientConfig
}

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Capability struct {
	Client   HttpClient
	Discover discovery.Interface
	Throttle flowctrl.RateLimiter
}

// Attention: all the fields must be string, or the ToHeader method will be panic.
// the struct filed tag is the key of header, and the header's value is the struct
// filed value.
type Headers struct {
	Language string `HTTP_BLUEKING_LANGUAGE`
	User     string `BK_User`
	OwnerID  string `HTTP_BLUEKING_SUPPLIER_ID`
}

func (h Headers) ToHeader() http.Header {
	header := make(http.Header)

	valueof := reflect.ValueOf(h)
	for i := 0; i < valueof.NumField(); i++ {
		k := reflect.TypeOf(h).Field(i).Tag
		v := valueof.Field(i).String()
		header[string(k)] = []string{v}
	}

	return header
}

type TLSClientConfig struct {
	// Server should be accessed without verifying the TLS certificate. For testing only.
	InsecureSkipVerify bool
	// Server requires TLS client certificate authentication
	CertFile string
	// Server requires TLS client certificate authentication
	KeyFile string
	// Trusted root certificates for server
	CAFile string
	// the password to decrypt the certificate
	Password string
}
