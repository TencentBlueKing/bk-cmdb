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
	"fmt"

	"configcenter/src/apimachinery/discovery"
	"configcenter/src/apimachinery/flowctrl"
	"github.com/prometheus/client_golang/prometheus"
)

type APIMachineryConfig struct {
	// request's qps value
	QPS int64
	// request's burst value
	Burst     int64
	TLSConfig *TLSClientConfig
}

type Capability struct {
	Client   HttpClient
	Discover discovery.Interface
	Throttle flowctrl.RateLimiter
	Mock     MockInfo
	Reg      prometheus.Registerer
}

type MockInfo struct {
	Mocked      bool
	SetMockData bool
	MockData    interface{}
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

func NewTLSClientConfigFromConfig(prefix string, config map[string]string) (TLSClientConfig, error) {
	tlsConfig := TLSClientConfig{}

	skipVerifyKey := fmt.Sprintf("%s.insecure_skip_verify", prefix)
	skipVerifyVal, ok := config[skipVerifyKey]
	if ok == true {
		if skipVerifyVal == "true" {
			tlsConfig.InsecureSkipVerify = true
		}
	}

	certFileKey := fmt.Sprintf("%s.cert_file", prefix)
	certFileVal, ok := config[certFileKey]
	if ok == true {
		tlsConfig.CertFile = certFileVal
	}

	keyFileKey := fmt.Sprintf("%s.key_file", prefix)
	keyFileVal, ok := config[keyFileKey]
	if ok == true {
		tlsConfig.KeyFile = keyFileVal
	}

	caFileKey := fmt.Sprintf("%s.ca_file", prefix)
	caFileVal, ok := config[caFileKey]
	if ok == true {
		tlsConfig.CAFile = caFileVal
	}

	passwordKey := fmt.Sprintf("%s.password", prefix)
	passwordVal, ok := config[passwordKey]
	if ok == true {
		tlsConfig.Password = passwordVal
	}

	return tlsConfig, nil
}
