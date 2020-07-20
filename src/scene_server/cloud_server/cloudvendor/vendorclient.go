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
	// NewVendorClient 创建云厂商客户端
	NewVendorClient(secretID, secretKey string) VendorClient
	// GetRegions 获取地域列表
	GetRegions() ([]*metadata.Region, error)
	// GetVpcs 获取vpc列表
	GetVpcs(region string, opt *ccom.VpcOpt) (*metadata.VpcsInfo, error)
	// GetInstances 获取实例列表
	GetInstances(region string, opt *ccom.InstanceOpt) (*metadata.InstancesInfo, error)
	// GetInstancesTotalCnt 获取实例总个数
	GetInstancesTotalCnt(region string, opt *ccom.InstanceOpt) (int64, error)
}

// Register 注册云厂商客户端
func Register(vendorName string, client VendorClient) {
	vendorClients[vendorName] = client
}

// GetVendorClient 获取云厂商客户端
func GetVendorClient(conf metadata.CloudAccountConf) (VendorClient, error) {
	var client VendorClient
	var ok bool
	if client, ok = vendorClients[conf.VendorName]; !ok {
		return nil, fmt.Errorf("vendor %s is not supported", conf.VendorName)
	}
	cli := client.NewVendorClient(conf.SecretID, conf.SecretKey)
	return cli, nil
}
