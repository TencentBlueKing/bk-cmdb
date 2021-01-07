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
	"configcenter/src/framework/core/config"
	"configcenter/src/framework/core/discovery"
)

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

func New(conf config.Config, disc discovery.DiscoverInterface) *Client {
	var c = &Client{}
	c.httpCli = httpclient.NewHttpClient()
	c.httpCli.SetHeader("Content-Type", "application/json")
	c.httpCli.SetHeader("Accept", "application/json")

	c.disc = disc
	c.supplierAccount = conf.Get("core.supplierAccount")
	c.user = conf.Get("core.user")
	c.SetAddress(conf.Get("core.ccaddress"))
	return c
}

func (cli *Client) Host() HostInterface {
	return newHost(cli)
}
func (cli *Client) Model() ModelInterface {
	return newModel(cli)
}
func (cli *Client) Classification() ClassificationInterface {
	return newClassification(cli)
}
func (cli *Client) Attribute() AttributeInterface {
	return newAttribute(cli)
}
func (cli *Client) CommonInst() CommonInstInterface {
	return newCommonInst(cli)
}
func (cli *Client) Group() GroupInterface {
	return newGroup(cli)
}
func (cli *Client) Business() BusinessInterface {
	return newBusiness(cli)
}
func (cli *Client) Module() ModuleInterface {
	return newModule(cli)
}
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
		//fmt.Println("client owner:", supplierAccount)
		//panic(supplierAccount)
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
