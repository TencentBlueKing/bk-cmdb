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

package cloudvendor

import (
	"fmt"

	"configcenter/src/common/metadata"
	ccom "configcenter/src/scene_server/cloud_server/common"
)

var vendorClients = make(map[string]VendorClient, 0)

type VendorClient interface {
	// 设置账号密码
	SetCredential(secretID, secretKey string)
	// 获取地域列表
	GetRegions(opt *ccom.RequestOpt) ([]*metadata.Region, error)
	// 获取vpc列表
	GetVpcs(region string, opt *ccom.RequestOpt) (*metadata.VpcsInfo, error)
	// 获取实例列表
	GetInstances(region string, opt *ccom.RequestOpt) (*metadata.InstancesInfo, error)
	// 获取实例总个数
	GetInstancesTotalCnt(region string, opt *ccom.RequestOpt) (int64, error)
}

// Register 注册云厂商客户端
func Register(vendorName string, client VendorClient) {
	vendorClients[vendorName] = client
}

// GetVendorClient 获取云厂商客户端
func GetVendorClient(conf metadata.CloudAccountConf) (VendorClient, error) {
	var client VendorClient
	var ok bool
	vendorName, ok := metadata.VendorNamesMap[conf.VendorName]
	if !ok {
		return nil, fmt.Errorf("vendor %s is invalid, it's not in VendorNamesMap %#v", conf.VendorName, metadata.VendorNamesMap)
	}
	if client, ok = vendorClients[vendorName]; !ok {
		return nil, fmt.Errorf("vendor %s is not supported", vendorName)
	}
	client.SetCredential(conf.SecretID, conf.SecretKey)
	return client, nil
}
