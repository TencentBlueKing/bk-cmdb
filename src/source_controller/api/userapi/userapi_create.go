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
 
package userapi

import (
	"configcenter/src/common"
	"configcenter/src/common/base"
	httpClient "configcenter/src/source_controller/api/client"
	"fmt"
)

type Client struct {
	httpClient.Client
}

// NewClient 创建审计日志操作接口
func NewClient(address string) *Client {

	cli := &Client{}
	cli.SetAddress(address)

	cli.Base = base.BaseLogic{}
	cli.Base.CreateHttpClient()

	return cli
}

//Create  新建user api列表
func (cli *Client) Create(input interface{}) (int, *common.APIRsp, error) {

	url := fmt.Sprintf("%s/host/v1/userapi", cli.GetAddress())
	return cli.GetRequestInfoEx(common.HTTPCreate, input, url)
}
