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
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	ccom "configcenter/src/scene_server/cloud_server/common"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func init() {
	Register(metadata.AWS, &awsClient{vendorName: metadata.AWS})
}

type awsClient struct {
	vendorName string
	secretID   string
	secretKey  string
}

const (
	awsMinPageSize int64 = 5
	awsMaxPageSize int64 = 1000
)

var regionIdNameMap = map[string]string{
	"us-east-1":      "美国东部（弗吉尼亚北部）",
	"us-east-2":      "美国东部（俄亥俄州）",
	"us-west-1":      "美国西部（加利福尼亚北部）",
	"us-west-2":      "美国西部（俄勒冈）",
	"ap-east-1":      "亚太地区（香港）",
	"ap-south-1":     "亚太地区（孟买）",
	"ap-southeast-1": "亚太区域（新加坡）",
	"ap-northeast-2": "亚太区域（首尔）",
	"ap-northeast-3": "亚太区域 （大阪当地）",
	"ap-southeast-2": "亚太区域（悉尼）",
	"ap-northeast-1": "亚太区域（东京）",
	"ca-central-1":   "加拿大 （中部）",
	"eu-central-1":   "欧洲（法兰克福）",
	"eu-west-1":      "欧洲（爱尔兰）",
	"eu-west-2":      "欧洲（伦敦）",
	"eu-west-3":      "欧洲（巴黎）",
	"eu-north-1":     "欧洲（斯德哥尔摩）",
	"me-south-1":     "中东（巴林）",
	"sa-east-1":      "南美洲（圣保罗）",
}

// NewVendorClient 创建云厂商客户端
func (c *awsClient) NewVendorClient(secretID, secretKey string) VendorClient {
	return &awsClient{
		vendorName: metadata.AWS,
		secretID:   secretID,
		secretKey:  secretKey,
	}
}

// GetRegions 获取地域列表
// API文档：https://docs.aws.amazon.com/AWSEC2/latest/APIReference/API_DescribeRegions.html
func (c *awsClient) GetRegions() ([]*metadata.Region, error) {
	sess, err := c.newSession("us-west-1")
	if err != nil {
		return nil, err
	}
	ec2Svc := ec2.New(sess)

	input := c.newDescribeRegionsInput()
	resp, err := ec2Svc.DescribeRegions(input)
	if err != nil {
		return nil, err
	}

	regionSet := make([]*metadata.Region, 0)
	for _, region := range resp.Regions {
		regionSet = append(regionSet, &metadata.Region{
			RegionId:   *region.RegionName,
			RegionName: regionIdNameMap[*region.RegionName],
		})
	}
	return regionSet, nil
}

// GetVpcs 获取vpc列表
// API文档：https://docs.aws.amazon.com/AWSEC2/latest/APIReference/API_DescribeVpcs.html
// 单个用户在单个地域下的vpc最大配额是5个
func (c *awsClient) GetVpcs(region string, opt *ccom.VpcOpt) (*metadata.VpcsInfo, error) {
	sess, err := c.newSession(region)
	if err != nil {
		return nil, err
	}
	ec2Svc := ec2.New(sess)

	vpcsInfo := new(metadata.VpcsInfo)
	loopCnt := 0
	if opt == nil {
		opt = ccom.GetDefaultVpcOpt()
	}
	input := c.newDescribeVpcsInput(opt)
	// 在limit小于全部数据量的情况下，获取limit数量的数据，否则获取全部数据
	for {
		output, err := ec2Svc.DescribeVpcs(input)
		if err != nil {
			return nil, err
		}
		for _, vpc := range output.Vpcs {
			vpcsInfo.VpcSet = append(vpcsInfo.VpcSet, &metadata.Vpc{
				VpcId:   *vpc.VpcId,
				VpcName: c.getVpcName(vpc),
			})
		}
		// 在获取到limit数量或者全部数据的情况下，退出循环
		if opt.Limit == int64(len(vpcsInfo.VpcSet)) || output.NextToken == nil || *output.NextToken == "" {
			break
		}
		loopCnt++
		// 设置分页请求参数
		input.NextToken = output.NextToken
		if loopCnt > ccom.MaxLoopCnt {
			blog.Errorf("DescribeVpcs loopCnt:%d, bigger than MaxLoopCnt, len(vpcsInfo.VpcSet):%d",
				loopCnt, vpcsInfo.VpcSet)
			return nil, ccom.ErrorLoopCnt
		}
	}
	vpcsInfo.Count = int64(len(vpcsInfo.VpcSet))

	return vpcsInfo, nil
}

// GetInstances 获取实例列表
// API文档：https://docs.aws.amazon.com/AWSEC2/latest/APIReference/API_DescribeInstances.html
func (c *awsClient) GetInstances(region string, opt *ccom.InstanceOpt) (*metadata.InstancesInfo, error) {
	instancesInfo := new(metadata.InstancesInfo)
	instances, isAll, err := c.getInstances(region, opt)
	if err != nil {
		return nil, err
	}
	instancesInfo.InstanceSet = instances

	totalCnt := int64(len(instances))
	// 如果查到的不是全量，则去获取实例总数
	if !isAll {
		instOpt := &ccom.InstanceOpt{ccom.BaseOpt{
			Limit: ccom.MaxLimit,
		}}
		totalCnt, err = c.GetInstancesTotalCnt(region, instOpt)
		if err != nil {
			return nil, err
		}
	}
	instancesInfo.Count = totalCnt

	return instancesInfo, nil
}

// getInstances 获取实例列表以及查到的是否为全量的bool值，不包含实例总个数
func (c *awsClient) getInstances(region string, opt *ccom.InstanceOpt) ([]*metadata.Instance, bool, error) {
	sess, err := c.newSession(region)
	if err != nil {
		return nil, false, err
	}
	ec2Svc := ec2.New(sess)

	if opt == nil {
		opt = ccom.GetDefaultInstanceOpt()
	}
	instances := make([]*metadata.Instance, 0)
	loopCnt := 0
	var nextToken *string
	input := c.newDescribeInstancesInput(opt)
	// 在limit小于全部数据量的情况下，获取limit数量的数据，否则获取全部数据
	for {
		output, err := ec2Svc.DescribeInstances(input)
		if err != nil {
			return nil, false, err
		}
		for _, reservation := range output.Reservations {
			for _, inst := range reservation.Instances {
				instances = append(instances, &metadata.Instance{
					InstanceId:    *inst.InstanceId,
					PrivateIp:     *inst.PrivateIpAddress,
					PublicIp:      *inst.PublicIpAddress,
					InstanceState: ccom.CovertInstState(*inst.State.Name),
					VpcId:         *inst.VpcId,
				})
			}
		}
		nextToken = output.NextToken
		// 在获取到limit数量或者全部数据的情况下，退出循环
		if opt.Limit == int64(len(instances)) || output.NextToken == nil || *output.NextToken == "" {
			break
		}
		// 设置分页请求参数
		input.NextToken = output.NextToken
		loopCnt++
		if loopCnt > ccom.MaxLoopCnt {
			blog.Errorf("DescribeVpcs loopCnt:%d, bigger than MaxLoopCnt, len(instances):%d",
				loopCnt, instances)
			return nil, false, ccom.ErrorLoopCnt
		}
	}
	isAll := false
	if nextToken == nil || *nextToken == "" {
		isAll = true
	}
	return instances, isAll, nil
}

// GetInstancesTotalCnt 获取实例总个数
func (c *awsClient) GetInstancesTotalCnt(region string, opt *ccom.InstanceOpt) (int64, error) {
	instances, _, err := c.getInstances(region, opt)
	if err != nil {
		return int64(0), err
	}
	return int64(len(instances)), nil
}

// newSession 创建会话
func (c *awsClient) newSession(region string) (*session.Session, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(c.secretID, c.secretKey, ""),
	})
	if err != nil {
		return nil, err
	}

	return sess, nil
}

// newDescribeRegionsInput 获取地域请求条件
func (c *awsClient) newDescribeRegionsInput() *ec2.DescribeRegionsInput {
	return new(ec2.DescribeRegionsInput)
}

// newDescribeVpcsInput 获取vpc请求条件
func (c *awsClient) newDescribeVpcsInput(opt *ccom.VpcOpt) *ec2.DescribeVpcsInput {
	input := &ec2.DescribeVpcsInput{}
	c.setFilters(&input.Filters, opt.Filters)
	c.setMaxResults(&input.MaxResults, opt.Limit)
	return input
}

// newDescribeInstancesInput 获取实例请求条件
func (c *awsClient) newDescribeInstancesInput(opt *ccom.InstanceOpt) *ec2.DescribeInstancesInput {
	input := &ec2.DescribeInstancesInput{}
	c.setFilters(&input.Filters, opt.Filters)
	c.setMaxResults(&input.MaxResults, opt.Limit)
	return input
}

// setFilters 设置过滤条件
func (c *awsClient) setFilters(dst *[]*ec2.Filter, src []*ccom.Filter) {
	if src == nil || len(src) == 0 {
		return
	}
	if dst == nil {
		tmp := make([]*ec2.Filter, 0)
		dst = &tmp
	}
	for i := range src {
		filter := &ec2.Filter{Name: src[i].Name, Values: src[i].Values}
		*dst = append(*dst, filter)
	}
}

// setMaxResults 设置单次请求返回结果条数
func (c *awsClient) setMaxResults(maxResults **int64, limit int64) {
	*maxResults = &limit
	// 按API要求，设置的MaxResults的取值范围为5～1000，不在该范围内的值会报错，不在该范围的设为最大值
	if limit < awsMinPageSize || limit > awsMaxPageSize {
		limit = awsMaxPageSize
	}
}

// 获取vpc名称，没有vpc名称标签，则使用vpcid作为名称
func (c *awsClient) getVpcName(vpc *ec2.Vpc) string {
	if len(vpc.Tags) <= 0 {
		return *vpc.VpcId
	}
	for _, tag := range vpc.Tags {
		if *tag.Key == "Name" {
			return *tag.Value
		}
	}
	return *vpc.VpcId
}
