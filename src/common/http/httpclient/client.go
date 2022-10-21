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

package httpclient

import (
	"bytes"
	"context"
	"crypto/tls"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"configcenter/src/apimachinery/util"
	"configcenter/src/common/ssl"
)

// HttpClient TODO
type HttpClient struct {
	caFile   string
	certFile string
	keyFile  string
	header   map[string]string
	httpCli  *http.Client
}

// NewHttpClient TODO
func NewHttpClient() *HttpClient {
	return &HttpClient{
		httpCli: &http.Client{},
		header:  make(map[string]string),
	}
}

// GetClient TODO
func (client *HttpClient) GetClient() *http.Client {
	return client.httpCli
}

// SetTlsNoVerity TODO
func (client *HttpClient) SetTlsNoVerity() error {
	tlsConf := ssl.ClientTLSConfNoVerify()

	trans := client.NewTransPort()
	trans.TLSClientConfig = tlsConf
	client.httpCli.Transport = trans

	return nil
}

// SetTlsVerityServer TODO
func (client *HttpClient) SetTlsVerityServer(caFile string) error {
	client.caFile = caFile

	// load ca cert
	tlsConf, err := ssl.ClientTslConfVerityServer(caFile)
	if err != nil {
		return err
	}

	client.SetTlsVerityConfig(tlsConf)

	return nil
}

// SetTLSVerify set tls verify config
func (client *HttpClient) SetTLSVerify(c *util.TLSClientConfig) error {
	// load cert
	tlsConf, err := ssl.ClientTLSConfVerity(c.CAFile, c.CertFile, c.KeyFile, c.Password)
	if err != nil {
		return err
	}
	tlsConf.InsecureSkipVerify = c.InsecureSkipVerify

	client.SetTlsVerityConfig(tlsConf)

	return nil
}

// SetTlsVerityConfig TODO
func (client *HttpClient) SetTlsVerityConfig(tlsConf *tls.Config) {
	trans := client.NewTransPort()
	trans.TLSClientConfig = tlsConf
	client.httpCli.Transport = trans
}

// NewTransPort TODO
func (client *HttpClient) NewTransPort() *http.Transport {
	return &http.Transport{
		Proxy:               http.ProxyFromEnvironment,
		TLSHandshakeTimeout: 5 * time.Second,
		Dial: (&net.Dialer{
			Timeout:   5 * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
		ResponseHeaderTimeout: 30 * time.Second,
	}
}

// SetTimeOut TODO
func (client *HttpClient) SetTimeOut(timeOut time.Duration) {
	client.httpCli.Timeout = timeOut
}

// SetHeader TODO
func (client *HttpClient) SetHeader(key, value string) {
	client.header[key] = value
}

// GetHeader TODO
func (client *HttpClient) GetHeader(key string) string {
	val, _ := client.header[key]
	return val
}

// GET TODO
func (client *HttpClient) GET(url string, header http.Header, data []byte) ([]byte, error) {
	return client.Request(url, "GET", header, data)

}

// POST TODO
func (client *HttpClient) POST(url string, header http.Header, data []byte) ([]byte, error) {
	return client.Request(url, "POST", header, data)
}

// DELETE TODO
func (client *HttpClient) DELETE(url string, header http.Header, data []byte) ([]byte, error) {
	return client.Request(url, "DELETE", header, data)
}

// PUT TODO
func (client *HttpClient) PUT(url string, header http.Header, data []byte) ([]byte, error) {
	return client.Request(url, "PUT", header, data)
}

// GETEx TODO
func (client *HttpClient) GETEx(url string, header http.Header, data []byte) (int, []byte, error) {
	return client.RequestEx(url, "GET", header, data)
}

// POSTEx TODO
func (client *HttpClient) POSTEx(url string, header http.Header, data []byte) (int, []byte, error) {
	return client.RequestEx(url, "POST", header, data)
}

// DELETEEx TODO
func (client *HttpClient) DELETEEx(url string, header http.Header, data []byte) (int, []byte, error) {
	return client.RequestEx(url, "DELETE", header, data)
}

// PUTEx TODO
func (client *HttpClient) PUTEx(url string, header http.Header, data []byte) (int, []byte, error) {
	return client.RequestEx(url, "PUT", header, data)
}

// Request TODO
func (client *HttpClient) Request(url, method string, header http.Header, data []byte) ([]byte, error) {
	var req *http.Request
	var errReq error
	if data != nil {
		req, errReq = http.NewRequest(method, url, bytes.NewReader(data))
	} else {
		req, errReq = http.NewRequest(method, url, nil)
	}

	if errReq != nil {
		return nil, errReq
	}

	req.Close = true

	if header != nil {
		req.Header = header
	}

	for key, value := range client.header {
		req.Header.Set(key, value)
	}

	rsp, err := client.httpCli.Do(req)
	if err != nil {
		return nil, err
	}

	/*if rsp.StatusCode >= http.StatusBadRequest {
		return 0, nil, fmt.Errorf("statuscode:%d, status:%s", rsp.StatusCode, rsp.Status)
	}*/

	defer rsp.Body.Close()

	body, err := ioutil.ReadAll(rsp.Body)

	return body, err
}

// RequestEx TODO
func (client *HttpClient) RequestEx(url, method string, header http.Header, data []byte) (int, []byte, error) {
	var req *http.Request
	var errReq error
	if data != nil {
		req, errReq = http.NewRequest(method, url, bytes.NewReader(data))
	} else {
		req, errReq = http.NewRequest(method, url, nil)
	}

	if errReq != nil {
		return 0, nil, errReq
	}

	req.Close = true

	if header != nil {
		req.Header = header
	}

	for key, value := range client.header {
		req.Header.Set(key, value)
	}

	rsp, err := client.httpCli.Do(req)
	if err != nil {
		return 0, nil, err
	}

	defer rsp.Body.Close()

	body, err := ioutil.ReadAll(rsp.Body)

	return rsp.StatusCode, body, err
}

// DoWithTimeout TODO
func (client *HttpClient) DoWithTimeout(timeout time.Duration, req *http.Request) (*http.Response, error) {
	ctx, _ := context.WithTimeout(req.Context(), timeout)
	req = req.WithContext(ctx)
	return client.httpCli.Do(req)
}
