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

	"configcenter/src/common/ssl"
)

type HttpClient struct {
	caFile   string
	certFile string
	keyFile  string
	header   map[string]string
	httpCli  *http.Client
}

func NewHttpClient() *HttpClient {
	return &HttpClient{
		httpCli: &http.Client{},
		header:  make(map[string]string),
	}
}

func (client *HttpClient) GetClient() *http.Client {
	return client.httpCli
}

func (client *HttpClient) SetTlsNoVerity() error {
	tlsConf := ssl.ClientTslConfNoVerity()

	trans := client.NewTransPort()
	trans.TLSClientConfig = tlsConf
	client.httpCli.Transport = trans

	return nil
}

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

func (client *HttpClient) SetTlsVerity(caFile, certFile, keyFile, passwd string) error {
	client.caFile = caFile
	client.certFile = certFile
	client.keyFile = keyFile

	// load cert
	tlsConf, err := ssl.ClientTLSConfVerity(caFile, certFile, keyFile, passwd)
	if err != nil {
		return err
	}

	client.SetTlsVerityConfig(tlsConf)

	return nil
}

func (client *HttpClient) SetTlsVerityConfig(tlsConf *tls.Config) {
	trans := client.NewTransPort()
	trans.TLSClientConfig = tlsConf
	client.httpCli.Transport = trans
}

func (client *HttpClient) NewTransPort() *http.Transport {
	return &http.Transport{
		TLSHandshakeTimeout: 5 * time.Second,
		Dial: (&net.Dialer{
			Timeout:   5 * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
		ResponseHeaderTimeout: 30 * time.Second,
	}
}

func (client *HttpClient) SetTimeOut(timeOut time.Duration) {
	client.httpCli.Timeout = timeOut
}

func (client *HttpClient) SetHeader(key, value string) {
	client.header[key] = value
}

func (client *HttpClient) GetHeader(key string) string {
	val, _ := client.header[key]
	return val
}

func (client *HttpClient) GET(url string, header http.Header, data []byte) ([]byte, error) {
	return client.Request(url, "GET", header, data)

}

func (client *HttpClient) POST(url string, header http.Header, data []byte) ([]byte, error) {
	return client.Request(url, "POST", header, data)
}

func (client *HttpClient) DELETE(url string, header http.Header, data []byte) ([]byte, error) {
	return client.Request(url, "DELETE", header, data)
}

func (client *HttpClient) PUT(url string, header http.Header, data []byte) ([]byte, error) {
	return client.Request(url, "PUT", header, data)
}

func (client *HttpClient) GETEx(url string, header http.Header, data []byte) (int, []byte, error) {
	return client.RequestEx(url, "GET", header, data)
}

func (client *HttpClient) POSTEx(url string, header http.Header, data []byte) (int, []byte, error) {
	return client.RequestEx(url, "POST", header, data)
}

func (client *HttpClient) DELETEEx(url string, header http.Header, data []byte) (int, []byte, error) {
	return client.RequestEx(url, "DELETE", header, data)
}

func (client *HttpClient) PUTEx(url string, header http.Header, data []byte) (int, []byte, error) {
	return client.RequestEx(url, "PUT", header, data)
}

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

func (client *HttpClient) DoWithTimeout(timeout time.Duration, req *http.Request) (*http.Response, error) {
	ctx, _ := context.WithTimeout(req.Context(), timeout)
	req = req.WithContext(ctx)
	return client.httpCli.Do(req)
}
