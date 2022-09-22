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
	"configcenter/framework/core/config"
	"configcenter/framework/core/discovery"
	"configcenter/pkg/common"
	"configcenter/pkg/http/httpclient"
)

// CCV3Interface TODO
type CCV3Interface interface {
	ModuleGetter
	SetGetter
	HostGetter
	ModelGetter
	BusinessGetter
	ClassificationGetter
	AttributeGetter
	CommonInstGetter
	GroupGetter
}

// Client the http client
type Client struct {
	httpCli         *httpclient.HttpClient
	disc            discovery.DiscoverInterface
	address         string
	supplierAccount string
	user            string
}

// New TODO
func New(conf config.Config, disc discovery.DiscoverInterface) *Client {
	var c = &Client{}
	c.httpCli = httpclient.NewHttpClient()
	c.httpCli.SetHeader("Content-Type", "application/json")
	c.httpCli.SetHeader("Accept", "application/json")

	c.disc = disc
	c.supplierAccount = conf.Get("logics.supplierAccount")
	c.user = conf.Get("logics.user")
	c.SetAddress(conf.Get("logics.ccaddress"))
	return c
}

// Host TODO
func (cli *Client) Host() HostInterface {
	return newHost(cli)
}

// Model TODO
func (cli *Client) Model() ModelInterface {
	return newModel(cli)
}

// Classification TODO
func (cli *Client) Classification() ClassificationInterface {
	return newClassification(cli)
}

// Attribute TODO
func (cli *Client) Attribute() AttributeInterface {
	return newAttribute(cli)
}

// CommonInst TODO
func (cli *Client) CommonInst() CommonInstInterface {
	return newCommonInst(cli)
}

// Group TODO
func (cli *Client) Group() GroupInterface {
	return newGroup(cli)
}

// Business TODO
func (cli *Client) Business() BusinessInterface {
	return newBusiness(cli)
}

// Module TODO
func (cli *Client) Module() ModuleInterface {
	return newModule(cli)
}

// Set TODO
func (cli *Client) Set() SetInterface {
	return newSet(cli)
}

// SetAddress set a new address
func (cli *Client) SetAddress(address string) {
	cli.address = address
}

// SetSupplierAccount set a new supplieraccount
func (cli *Client) SetSupplierAccount(supplierAccount string) {
	if 0 != len(supplierAccount) {
		cli.httpCli.SetHeader(common.BKHTTPOwnerID, supplierAccount)
	} else {
		cli.httpCli.SetHeader(common.BKHTTPOwnerID, cli.supplierAccount)
	}

}

// SetUser set a new user
func (cli *Client) SetUser(user string) {
	if 0 != len(user) {
		cli.httpCli.SetHeader(common.BKHTTPHeaderUser, user)
	} else {
		cli.httpCli.SetHeader(common.BKHTTPHeaderUser, cli.user)
	}
}

// GetUser get the user
func (cli *Client) GetUser() string {
	return cli.httpCli.GetHeader(common.BKHTTPHeaderUser)
}

// GetSupplierAccount get the supplier account
func (cli *Client) GetSupplierAccount() string {
	return cli.httpCli.GetHeader(common.BKHTTPOwnerID)
}

// GetAddress get the address
func (cli *Client) GetAddress() string {
	if cli.disc != nil {
		return cli.disc.Output()
	}
	return cli.address
}
