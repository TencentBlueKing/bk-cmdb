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

const (
	tcMinPageSize int64 = 1
	tcMaxPageSize int64 = 100
)

// NewVendorClient 创建云厂商客户端
func (c *tcClient) NewVendorClient(secretID, secretKey string) VendorClient {
	return &tcClient{
		vendorName: metadata.TencentCloud,
		secretID:   secretID,
		secretKey:  secretKey,
	}
}

// GetRegions 获取地域列表
// API文档：https://cloud.tencent.com/document/api/213/15708
func (c *tcClient) GetRegions() ([]*metadata.Region, error) {
	credential := c.newCredential(c.secretID, c.secretKey)
	client, err := cvm.NewClient(credential, regions.Guangzhou, profile.NewClientProfile())
	if err != nil {
		return nil, err
	}

	request := c.newDescribeRegionsRequest()
	resp, err := client.DescribeRegions(request)
	if err != nil {
		return nil, err
	}

	regionSet := make([]*metadata.Region, 0)
	for _, region := range resp.Response.RegionSet {
		regionSet = append(regionSet, &metadata.Region{
			RegionId:    *region.Region,
			RegionName:  *region.RegionName,
			RegionState: *region.RegionState,
		})
	}

	return regionSet, nil
}

// GetVpcs 获取vpc列表
// API文档：https://cloud.tencent.com/document/api/215/15778
func (c *tcClient) GetVpcs(region string, opt *ccom.VpcOpt) (*metadata.VpcsInfo, error) {
	credential := c.newCredential(c.secretID, c.secretKey)
	client, err := tcVpc.NewClient(credential, region, profile.NewClientProfile())
	if err != nil {
		return nil, err
	}

	if opt == nil {
		opt = ccom.GetDefaultVpcOpt()
	}
	vpcsInfo := new(metadata.VpcsInfo)
	loopCnt := 0
	var totalCnt int64 = 0
	request := c.newDescribeVpcsRequest(opt)
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
		if opt.Limit == int64(len(vpcsInfo.VpcSet)) || len(vpcsInfo.VpcSet) == int(totalCnt) {
			break
		}
		// 设置分页请求参数
		offset := fmt.Sprintf("%d", len(vpcsInfo.VpcSet))
		request.Offset = &offset
		loopCnt++
		if loopCnt > ccom.MaxLoopCnt {
			blog.Errorf("DescribeVpcs loopCnt:%d, bigger than MaxLoopCnt, TotalCount:%d",
				loopCnt, resp.Response.TotalCount)
			return nil, ccom.ErrorLoopCnt
		}
	}
	vpcsInfo.Count = totalCnt

	return vpcsInfo, nil
}

// GetInstances 获取实例列表
// API文档：https://cloud.tencent.com/document/api/213/15728
func (c *tcClient) GetInstances(region string, opt *ccom.InstanceOpt) (*metadata.InstancesInfo, error) {
	credential := c.newCredential(c.secretID, c.secretKey)
	client, err := cvm.NewClient(credential, region, profile.NewClientProfile())
	if err != nil {
		return nil, err
	}

	if opt == nil {
		opt = ccom.GetDefaultInstanceOpt()
	}
	instancesInfo := new(metadata.InstancesInfo)
	loopCnt := 0
	var totalCnt int64 = 0
	request := c.newDescribeInstancesRequest(opt)
	// 在limit小于全部数据量的情况下，获取limit数量的数据，否则获取全部数据
	for {
		resp, err := client.DescribeInstances(request)
		if err != nil {
			return nil, err
		}

		for _, inst := range resp.Response.InstanceSet {
			privateIP := ""
			if len(inst.PrivateIpAddresses) > 0 {
				privateIP = *inst.PrivateIpAddresses[0]
			}
			publicIP := ""
			if len(inst.PublicIpAddresses) > 0 {
				publicIP = *inst.PublicIpAddresses[0]
			}
			instancesInfo.InstanceSet = append(instancesInfo.InstanceSet, &metadata.Instance{
				InstanceId:    *inst.InstanceId,
				PrivateIp:     privateIP,
				PublicIp:      publicIP,
				InstanceState: ccom.CovertInstState(*inst.InstanceState),
				VpcId:         *inst.VirtualPrivateCloud.VpcId,
			})
		}
		totalCnt = int64(*resp.Response.TotalCount)
		// 在获取到limit数量或者全部数据的情况下，退出循环
		if opt.Limit == int64(len(instancesInfo.InstanceSet)) || len(instancesInfo.InstanceSet) == int(totalCnt) {
			break
		}
		// 设置分页请求参数
		offset := int64(len(instancesInfo.InstanceSet))
		request.Offset = &offset
		loopCnt++
		if loopCnt > ccom.MaxLoopCnt {
			blog.Errorf("DescribeInstances loopCnt:%d, bigger than MaxLoopCnt, TotalCount:%d",
				loopCnt, resp.Response.TotalCount)
			return nil, ccom.ErrorLoopCnt
		}
	}
	instancesInfo.Count = totalCnt

	return instancesInfo, nil
}

// GetInstancesTotalCnt 获取实例总个数
func (c *tcClient) GetInstancesTotalCnt(region string, opt *ccom.InstanceOpt) (int64, error) {
	if opt == nil {
		opt = ccom.GetDefaultInstanceOpt()
	}
	// 直接将limit设为最小值，能最快地获取到实例总个数
	opt.Limit = tcMinPageSize
	instsInfo, err := c.GetInstances(region, opt)
	if err != nil {
		return 0, err
	}
	return instsInfo.Count, nil
}

// newCredential 创建认证信息
func (c *tcClient) newCredential(secretID, secretKey string) *tcCommon.Credential {
	return tcCommon.NewCredential(secretID, secretKey)
}

// newDescribeRegionsRequest 获取地域请求条件
func (c *tcClient) newDescribeRegionsRequest() *cvm.DescribeRegionsRequest {
	request := cvm.NewDescribeRegionsRequest()
	return request
}

// newDescribeVpcsRequest 获取vpc请求条件
func (c *tcClient) newDescribeVpcsRequest(opt *ccom.VpcOpt) *tcVpc.DescribeVpcsRequest {
	request := tcVpc.NewDescribeVpcsRequest()
	c.setVpcFilters(&request.Filters, opt.Filters)
	c.setVpcLimit(&request.Limit, opt.Limit)
	return request
}

// newDescribeInstancesRequest 获取实例请求条件
func (c *tcClient) newDescribeInstancesRequest(opt *ccom.InstanceOpt) *cvm.DescribeInstancesRequest {
	request := cvm.NewDescribeInstancesRequest()
	c.setCvmFilters(&request.Filters, opt.Filters)
	c.setCvmLimit(&request.Limit, opt.Limit)
	return request
}

// setVpcFilters 设置过滤条件
func (c *tcClient) setVpcFilters(dst *[]*tcVpc.Filter, src []*ccom.Filter) {
	if src == nil || len(src) == 0 {
		return
	}
	if dst == nil {
		tmp := make([]*tcVpc.Filter, 0)
		dst = &tmp
	}
	for i := range src {
		filter := &tcVpc.Filter{Name: src[i].Name, Values: src[i].Values}
		*dst = append(*dst, filter)
	}
}

// setVpcLimit 设置单次请求返回结果条数
func (c *tcClient) setVpcLimit(Dstlimit **string, limit int64) {
	limitStr := fmt.Sprintf("%d", limit)
	*Dstlimit = &limitStr
	// 按API要求，设置的MaxResults的取值范围为0～100，不在该范围内的值会报错，不在该范围的设为最大值
	if limit < tcMinPageSize || limit > tcMaxPageSize {
		limitStr = fmt.Sprintf("%d", tcMaxPageSize)
	}
}

// setCvmFilters 设置过滤条件
func (c *tcClient) setCvmFilters(dst *[]*cvm.Filter, src []*ccom.Filter) {
	if src == nil || len(src) == 0 {
		return
	}
	if dst == nil {
		tmp := make([]*cvm.Filter, 0)
		dst = &tmp
	}
	for i := range src {
		filter := &cvm.Filter{Name: src[i].Name, Values: src[i].Values}
		*dst = append(*dst, filter)
	}
}

// setCvmLimit 设置单次请求返回结果条数
func (c *tcClient) setCvmLimit(Dstlimit **int64, limit int64) {
	*Dstlimit = &limit
	// 按API要求，设置的MaxResults的取值范围为1～100，不在该范围内的值会报错，不在该范围的设为最大值
	if limit < tcMinPageSize || limit > tcMaxPageSize {
		limit = tcMaxPageSize
	}
}
