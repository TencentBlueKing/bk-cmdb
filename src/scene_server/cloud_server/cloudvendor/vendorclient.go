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
)

var vendorClients = make(map[string]VendorClient, 0)

type VendorClient interface {
	SetCredential(secretID, secretKey string)
	GetRegions() ([]string, error)
	GetVpcs(region string) ([]*metadata.Vpc, error)
	GetInstances(region string) ([]*metadata.Instance, error)
}

// Register 注册云厂商客户端
func Register(vendorName string, client VendorClient) {
	vendorClients[vendorName] = client
}

// GetVendorClient 获取云厂商客户端
func GetVendorClient(vendorName, secretID, secretKey string) (VendorClient, error) {
	var client VendorClient
	var ok bool
	if client, ok = vendorClients[vendorName]; !ok {
		return nil, fmt.Errorf("vendor %s is not supported", vendorName)
	}
	client.SetCredential(secretID, secretKey)
	return client, nil
}
