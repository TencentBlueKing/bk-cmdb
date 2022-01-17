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
	"net"
	"net/http"
	"time"

	"configcenter/src/common/ssl"
	"configcenter/src/thirdparty/logplatform"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

var httpClient *http.Client

func NewClient(c *TLSClientConfig) (*http.Client, error) {
	if httpClient != nil {
		return httpClient, nil
	}
	tlsConf := new(tls.Config)
	if nil != c {
		tlsConf.InsecureSkipVerify = c.InsecureSkipVerify
		if len(c.CAFile) != 0 && len(c.CertFile) != 0 && len(c.KeyFile) != 0 {
			var err error
			tlsConf, err = ssl.ClientTLSConfVerity(c.CAFile, c.CertFile, c.KeyFile, c.Password)
			if err != nil {
				return nil, err
			}
		}
	}

	transport := &http.Transport{
		Proxy:               http.ProxyFromEnvironment,
		TLSHandshakeTimeout: 5 * time.Second,
		TLSClientConfig:     tlsConf,
		Dial: (&net.Dialer{
			Timeout:   5 * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
		MaxIdleConnsPerHost:   100,
		ResponseHeaderTimeout: 10 * time.Minute,
	}

	httpClient = &http.Client{
		Transport: transport,
	}
	return httpClient, nil
}

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// WrapperTraceClient wrapper client to record trace
func WrapperTraceClient() {
	if httpClient == nil {
		return
	}
	if logplatform.OpenTelemetryCfg.Enable {
		httpClient.Transport = otelhttp.NewTransport(httpClient.Transport)
	}
}
