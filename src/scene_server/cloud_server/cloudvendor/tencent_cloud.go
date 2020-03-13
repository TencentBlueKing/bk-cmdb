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

	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	ccom "configcenter/src/scene_server/cloud_server/common"

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

const tcMaxPageSize int64 = 100

// 设置账号密码
func (c *tcClient) SetCredential(secretID, secretKey string) {
	c.secretID = secretID
	c.secretKey = secretKey
}

// 获取地域列表
// API文档：https://cloud.tencent.com/document/api/213/15708
func (c *tcClient) GetRegions(opt *ccom.RequestOpt) (*metadata.RegionsInfo, error) {
	credential := c.newCredential(c.secretID, c.secretKey)
	client, err := cvm.NewClient(credential, regions.Guangzhou, profile.NewClientProfile())
	if err != nil {
		return nil, err
	}

	request := c.NewDescribeRegionsRequest(opt)
	resp, err := client.DescribeRegions(request)
	if err != nil {
		return nil, err
	}

	regionsInfo := new(metadata.RegionsInfo)
	for _, region := range resp.Response.RegionSet {
		regionsInfo.RegionSet = append(regionsInfo.RegionSet, &metadata.Region{
			RegionId:    *region.Region,
			RegionName:  *region.RegionName,
			RegionState: *region.RegionState,
		})
	}
	regionsInfo.Count = int64(*resp.Response.TotalCount)

	return regionsInfo, nil
}

// 获取vpc列表
// API文档：https://cloud.tencent.com/document/api/215/15778
func (c *tcClient) GetVpcs(region string, opt *ccom.RequestOpt) (*metadata.VpcsInfo, error) {
	credential := c.newCredential(c.secretID, c.secretKey)
	client, err := tcVpc.NewClient(credential, region, profile.NewClientProfile())
	if err != nil {
		return nil, err
	}

	vpcsInfo := new(metadata.VpcsInfo)
	loopCnt := 0
	var totalCnt int64 = 0
	request := c.NewDescribeVpcsRequest(opt)
	// 在limit小于全部数据量的情况下，获取limit数量的数据，否则获取全部数据
	for {
		resp, err := client.DescribeVpcs(request)
		if err != nil {
			return nil, err
		}
		for _, vpc := range resp.Response.VpcSet {
			vpcsInfo.VpcSet = append(vpcsInfo.VpcSet, &metadata.Vpc{
				VpcId:   *vpc.VpcId,
				VpcName: *vpc.VpcName,
			})
		}
		totalCnt = int64(*resp.Response.TotalCount)
		// 在获取到limit数量或者全部数据的情况下，退出循环
		if opt == nil || opt.Limit == nil || *opt.Limit == int64(len(vpcsInfo.VpcSet)) || len(vpcsInfo.VpcSet) == int(totalCnt) {
			break
		}
		// 设置分页请求参数
		offset := fmt.Sprintf("%d", len(vpcsInfo.VpcSet))
		request.Offset = &offset
		loopCnt++
		if loopCnt > ccom.MaxLoopCnt {
			blog.Errorf("DescribeVpcs loopCnt:%d, bigger than MaxLoopCnt, TotalCount:%d", loopCnt, resp.Response.TotalCount)
			return nil, ccom.ErrorLoopCnt
		}
	}
	vpcsInfo.Count = totalCnt

	return vpcsInfo, nil
}

// 获取实例列表
// API文档：https://cloud.tencent.com/document/api/213/15728
func (c *tcClient) GetInstances(region string, opt *ccom.RequestOpt) (*metadata.InstancesInfo, error) {
	credential := c.newCredential(c.secretID, c.secretKey)
	client, err := cvm.NewClient(credential, region, profile.NewClientProfile())
	if err != nil {
		return nil, err
	}

	instancesInfo := new(metadata.InstancesInfo)
	loopCnt := 0
	var totalCnt int64 = 0
	request := c.NewDescribeInstancesRequest(opt)
	// 在limit小于全部数据量的情况下，获取limit数量的数据，否则获取全部数据
	for {
		resp, err := client.DescribeInstances(request)
		if err != nil {
			return nil, err
		}

		for _, inst := range resp.Response.InstanceSet {
			instancesInfo.InstanceSet = append(instancesInfo.InstanceSet, &metadata.Instance{
				InstanceId:    *inst.InstanceId,
				InstanceName:  *inst.InstanceName,
				PrivateIp:     *inst.PrivateIpAddresses[0],
				PublicIp:      *inst.PublicIpAddresses[0],
				InstanceState: ccom.CovertInstState(*inst.InstanceState),
				VpcId:         *inst.VirtualPrivateCloud.VpcId,
				OsName:        *inst.OsName,
			})
		}
		totalCnt = int64(*resp.Response.TotalCount)
		// 在获取到limit数量或者全部数据的情况下，退出循环
		if opt == nil || opt.Limit == nil || *opt.Limit == int64(len(instancesInfo.InstanceSet)) || len(instancesInfo.InstanceSet) == int(totalCnt) {
			break
		}
		// 设置分页请求参数
		offset := int64(len(instancesInfo.InstanceSet))
		request.Offset = &offset
		loopCnt++
		if loopCnt > ccom.MaxLoopCnt {
			blog.Errorf("DescribeInstances loopCnt:%d, bigger than MaxLoopCnt, TotalCount:%d", loopCnt, resp.Response.TotalCount)
			return nil, ccom.ErrorLoopCnt
		}
	}
	instancesInfo.Count = totalCnt

	return instancesInfo, nil
}

func (c *tcClient) newCredential(secretID, secretKey string) *tcCommon.Credential {
	return tcCommon.NewCredential(secretID, secretKey)
}

// 获取地域请求条件
func (c *tcClient) NewDescribeRegionsRequest(opt *ccom.RequestOpt) *cvm.DescribeRegionsRequest {
	request := cvm.NewDescribeRegionsRequest()
	return request
}

// 获取vpc请求条件
func (c *tcClient) NewDescribeVpcsRequest(opt *ccom.RequestOpt) *tcVpc.DescribeVpcsRequest {
	request := tcVpc.NewDescribeVpcsRequest()
	if opt == nil {
		return request
	}
	if len(opt.Filters) > 0 {
		request.Filters = make([]*tcVpc.Filter, 0)
		for i, _ := range opt.Filters {
			filter := &tcVpc.Filter{Name: opt.Filters[i].Name, Values: opt.Filters[i].Values}
			request.Filters = append(request.Filters, filter)
		}
	}
	if opt.Limit != nil {
		limit := fmt.Sprintf("%d", *opt.Limit)
		if *opt.Limit > tcMaxPageSize {
			limit = fmt.Sprintf("%d", tcMaxPageSize)
		}
		request.Limit = &limit
	}
	if opt.Offset != nil {
		offset := fmt.Sprintf("%d", *opt.Offset)
		request.Offset = &offset
	}

	return request
}

// 获取实例请求条件
func (c *tcClient) NewDescribeInstancesRequest(opt *ccom.RequestOpt) *cvm.DescribeInstancesRequest {
	request := cvm.NewDescribeInstancesRequest()
	if opt == nil {
		return request
	}
	if len(opt.Filters) > 0 {
		request.Filters = make([]*cvm.Filter, 0)
		for i, _ := range opt.Filters {
			filter := &cvm.Filter{Name: opt.Filters[i].Name, Values: opt.Filters[i].Values}
			request.Filters = append(request.Filters, filter)
		}
	}
	if opt.Limit != nil {
		limit := *opt.Limit
		if *opt.Limit > tcMaxPageSize {
			limit = tcMaxPageSize
		}
		request.Limit = &limit
	}
	request.Offset = opt.Offset
	return request
}
