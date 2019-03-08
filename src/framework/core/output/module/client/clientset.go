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

package client

import (
	"configcenter/src/framework/core/config"
	"configcenter/src/framework/core/discovery"
	"configcenter/src/framework/core/output/module/client/v3"
)

var _ Interface = &Clientset{}

type Params struct {
	SupplierAccount string
	UserName        string
}

type Interface interface {
	CCV3(params Params) v3.CCV3Interface
}

type Clientset struct {
	ccv3 *v3.Client
}

func (c *Clientset) CCV3(params Params) v3.CCV3Interface {
	if c == nil {
		return nil
	}

	c.ccv3.SetSupplierAccount(params.SupplierAccount)
	c.ccv3.SetUser(params.UserName)

	return c.ccv3
}

func NewForConfig(c config.Config, disc discovery.DiscoverInterface) *Clientset {
	var cs Clientset
	cs.ccv3 = v3.New(c, disc)
	client = &cs
	return &cs
}

var client *Clientset

func GetClient() *Clientset {
	return client
}
