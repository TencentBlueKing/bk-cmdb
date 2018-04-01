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

package example

import (
	"configcenter/src/common/http/httpclient"
	"fmt"
	"net/http"
)

const (
	// APIServerAddress api server address
	APIServerAddress = "http://127.0.0.1:8080/api/v1"
)

var client = &mockClient{client: httpclient.NewHttpClient()}

type mockClient struct {
	client *httpclient.HttpClient
}

func (cli *mockClient) GET(uri string, header http.Header, data []byte) ([]byte, error) {
	return cli.client.GET(fmt.Sprintf("%s%s", APIServerAddress, uri), header, data)

}

func (cli *mockClient) POST(uri string, header http.Header, data []byte) ([]byte, error) {
	return cli.client.POST(fmt.Sprintf("%s%s", APIServerAddress, uri), header, data)
}

func (cli *mockClient) DELETE(uri string, header http.Header, data []byte) ([]byte, error) {
	return cli.client.POST(fmt.Sprintf("%s%s", APIServerAddress, uri), header, data)
}

func (cli *mockClient) PUT(uri string, header http.Header, data []byte) ([]byte, error) {
	return cli.client.PUT(fmt.Sprintf("%s%s", APIServerAddress, uri), header, data)
}
