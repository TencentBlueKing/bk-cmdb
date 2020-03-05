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
	"configcenter/src/common/metadata"

	tcCommon "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/regions"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
	tcVpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
)

func init() {
	Register(metadata.TencentCloud, &tcClient{vendorName: metadata.TencentCloud})
}

type tcClient struct {
	vendorName string
	secretID   string
	secretKey  string
}

func (c *tcClient) SetCredential(secretID, secretKey string) {
	c.secretID = secretID
	c.secretKey = secretKey
}

func (c *tcClient) GetRegions() ([]string, error) {
	credential := c.newCredential(c.secretID, c.secretKey)
	client, err := cvm.NewClient(credential, regions.Guangzhou, profile.NewClientProfile())
	if err != nil {
		return nil, err
	}

	rsp, err := client.DescribeRegions(nil)
	if err != nil {
		return nil, err
	}

	regions := make([]string, 0)
	for _, region := range rsp.Response.RegionSet {
		regions = append(regions, *region.Region)
	}

	return regions, nil
}

func (c *tcClient) GetVpcs(region string) ([]*metadata.Vpc, error) {
	credential := c.newCredential(c.secretID, c.secretKey)
	client, err := tcVpc.NewClient(credential, region, profile.NewClientProfile())
	if err != nil {
		return nil, err
	}

	resp, err := client.DescribeVpcs(nil)
	if err != nil {
		return nil, err
	}
	vpcs := make([]*metadata.Vpc, 0)
	for _, vpc := range resp.Response.VpcSet {
		vpcs = append(vpcs, &metadata.Vpc{
			VpcId:   *vpc.VpcId,
			VpcName: *vpc.VpcName,
		})
	}

	return vpcs, nil
}

func (c *tcClient) GetInstances(region string) ([]*metadata.Instance, error) {
	credential := c.newCredential(c.secretID, c.secretKey)
	client, err := cvm.NewClient(credential, region, profile.NewClientProfile())
	if err != nil {
		return nil, err
	}
	rsp, err := client.DescribeInstances(nil)
	if err != nil {
		return nil, err
	}

	instances := make([]*metadata.Instance, 0)
	for _, inst := range rsp.Response.InstanceSet {
		instance := &metadata.Instance{
			InstanceId:    *inst.InstanceId,
			InstanceName:  *inst.InstanceName,
			PrivateIp:     *inst.PrivateIpAddresses[0],
			PublicIp:      *inst.PublicIpAddresses[0],
			InstanceState: *inst.InstanceState,
			VpcId:         *inst.VirtualPrivateCloud.VpcId,
		}
		instances = append(instances, instance)

	}
	return instances, nil
}

func (c *tcClient) newCredential(secretID, secretKey string) *tcCommon.Credential {
	return tcCommon.NewCredential(secretID, secretKey)
}
