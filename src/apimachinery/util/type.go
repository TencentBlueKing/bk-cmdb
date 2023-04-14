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
	"crypto/tls"
	"fmt"
	"time"

	"configcenter/src/apimachinery/discovery"
	"configcenter/src/apimachinery/flowctrl"
	cc "configcenter/src/common/backbone/configcenter"
	"configcenter/src/common/ssl"

	"github.com/prometheus/client_golang/prometheus"
)

// APIMachineryConfig TODO
type APIMachineryConfig struct {
	// request's qps value
	QPS int64
	// request's burst value
	Burst     int64
	TLSConfig *TLSClientConfig
	ExtraConf *ExtraClientConfig
}

// Capability TODO
type Capability struct {
	Client     HttpClient
	Discover   discovery.Interface
	Throttle   flowctrl.RateLimiter
	Mock       MockInfo
	MetricOpts MetricOption
	// the max tolerance api request latency time, if exceeded this time, then
	// this request will be logged and warned.
	ToleranceLatencyTime time.Duration
}

// MetricOption TODO
type MetricOption struct {
	// prometheus metric register
	Register prometheus.Registerer
	// if not set, use default buckets value
	DurationBuckets []float64
}

// MockInfo TODO
type MockInfo struct {
	Mocked      bool
	SetMockData bool
	MockData    interface{}
}

// TLSClientConfig TODO
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

// NewTLSClientConfigFromConfig new config about tls client config
func NewTLSClientConfigFromConfig(prefix string) (TLSClientConfig, error) {
	tlsConfig := TLSClientConfig{}

	skipVerifyKey := fmt.Sprintf("%s.insecureSkipVerify", prefix)
	if val, err := cc.String(skipVerifyKey); err == nil {
		skipVerifyVal := val
		if skipVerifyVal == "true" {
			tlsConfig.InsecureSkipVerify = true
		}
	}

	certFileKey := fmt.Sprintf("%s.certFile", prefix)
	if val, err := cc.String(certFileKey); err == nil {
		tlsConfig.CertFile = val
	}

	keyFileKey := fmt.Sprintf("%s.keyFile", prefix)
	if val, err := cc.String(keyFileKey); err == nil {
		tlsConfig.KeyFile = val
	}

	caFileKey := fmt.Sprintf("%s.caFile", prefix)
	if val, err := cc.String(caFileKey); err == nil {
		tlsConfig.CAFile = val
	}

	passwordKey := fmt.Sprintf("%s.password", prefix)
	if val, err := cc.String(passwordKey); err == nil {
		tlsConfig.Password = val
	}

	return tlsConfig, nil
}

// ExtraClientConfig extra http client configuration
type ExtraClientConfig struct {
	// ResponseHeaderTimeout the amount of time to wait for a server's response headers
	ResponseHeaderTimeout time.Duration
}

// GetClientTLSConfig get client tls config
func GetClientTLSConfig(prefix string) (*tls.Config, error) {
	tlsConf := &tls.Config{
		InsecureSkipVerify: true,
	}

	config, err := NewTLSClientConfigFromConfig(prefix)
	if err != nil {
		return nil, err
	}

	if len(config.CAFile) != 0 && len(config.CertFile) != 0 && len(config.KeyFile) != 0 {
		tlsConf, err = ssl.ClientTLSConfVerity(config.CAFile, config.CertFile, config.KeyFile, config.Password)
		if err != nil {
			return nil, err
		}
		tlsConf.InsecureSkipVerify = config.InsecureSkipVerify
	}

	return tlsConf, nil
}
