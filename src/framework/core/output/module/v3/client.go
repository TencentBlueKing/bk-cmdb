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

package v3

import (
	"configcenter/src/common"
	"configcenter/src/common/http/httpclient"
)

// Client the http client
type Client struct {
	httpCli         *httpclient.HttpClient
	address         string
	supplierAccount string
	user            string
}

var client = &Client{}

func init() {

	client.httpCli = httpclient.NewHttpClient()
	client.httpCli.SetHeader("Content-Type", "application/json")
	client.httpCli.SetHeader("Accept", "application/json")
}

// GetClient get the client instance
func GetClient() *Client {
	return client
}

// GetV3Client get the v3 client
func GetV3Client() *Client {

	return client
}

// SetAddress set a new address
func (cli *Client) SetAddress(address string) {
	cli.address = address
}

// SetSupplierAccount set a new supplieraccount
func (cli *Client) SetSupplierAccount(supplierAccount string) {
	cli.supplierAccount = supplierAccount
	cli.httpCli.SetHeader(common.BKHTTPOwnerID, supplierAccount)
}

// SetUser set a new user
func (cli *Client) SetUser(user string) {
	cli.user = user
	cli.httpCli.SetHeader(common.BKHTTPHeaderUser, user)
}

// GetUser get the user
func (cli *Client) GetUser() string {
	return cli.user
}

// GetSupplierAccount get the supplier account
func (cli *Client) GetSupplierAccount() string {
	return cli.supplierAccount
}

// GetAddress get the address
func (cli *Client) GetAddress() string {
	return cli.address
}
