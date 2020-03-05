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

package logics

import (
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/cloud_server/cloudvendor"
)

func (lgc *Logics) AccountVerify(vendorName, secretID, secretKey string) (bool, error) {
	client, err := cloudvendor.GetVendorClient(vendorName, secretID, secretKey)
	if err != nil {
		return false, err
	}

	_, err = client.GetRegions()
	if err != nil {
		return false, err
	}

	return true, nil
}

func (lgc *Logics) GetVpcInstInfo(region, vendorName, secretID, secretKey string) ([]*metadata.Vpc, error) {
	client, err := cloudvendor.GetVendorClient(vendorName, secretID, secretKey)
	if err != nil {
		return nil, err
	}
	//metadata.VpcInstanceInfo

	return client.GetVpcs(region)
}