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
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	tcCommon "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	tcRegions "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/regions"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
	tcVpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
)

func (lgc *Logics) GetAwsRegions(kit *rest.Kit, secretID, secretKey string) ([]string, error) {
	sess, err := lgc.AwsNewSession(kit, "", secretID, secretKey)
	if err != nil {
		blog.ErrorJSON("getAwsRegions get aws new session failed, err: %v, rid: %s", err, kit.Rid)
		return nil, err
	}
	ec2Svc := ec2.New(sess)
	rsp, err := ec2Svc.DescribeRegions(nil)
	if err != nil {
		blog.ErrorJSON("getAwsRegions, sdk api DescribeRegions failed, err: %v. rid: %s", err, kit.Rid)
		return nil, err
	}

	regions := make([]string, 0)
	for _, region := range rsp.Regions {
		regions = append(regions, *region.RegionName)
	}

	return regions, nil
}

func (lgc *Logics) GetTencentCloudRegions(kit *rest.Kit, secretID, secretKey string) ([]string, error) {
	credential := lgc.TencentCloudNewCredential(kit, secretID, secretKey)

	client, err := cvm.NewClient(credential, tcRegions.Guangzhou, profile.NewClientProfile())
	if err != nil {
		blog.ErrorJSON("getTencentCloudRegions new client failed, err: %v, rid: %s", err, kit.Rid)
		return nil, nil
	}

	regionRequest := cvm.NewDescribeRegionsRequest()
	rsp, err := client.DescribeRegions(regionRequest)
	regions := make([]string, 0)
	for _, region := range rsp.Response.RegionSet {
		regions = append(regions, *region.Region)
	}

	return regions, nil
}

func (lgc *Logics) AwsNewSession(kit *rest.Kit, region, secretID, secretKey string) (*session.Session, error) {
	if region == "" {
		region = "us-west-2"
	}
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(secretID, secretKey, ""),
	})
	if err != nil {
		return nil, err
	}

	return sess, nil
}

func (lgc *Logics) TencentCloudNewCredential(kit *rest.Kit, secretID, secretKey string) *tcCommon.Credential {
	return tcCommon.NewCredential(secretID, secretKey)
}

func (lgc *Logics) GetAwsVpc(kit *rest.Kit, secretID, secretKey string) ([]metadata.VpcInfo, error) {
	regions, err := lgc.GetAwsRegions(kit, secretID, secretKey)
	if err != nil {
		blog.ErrorJSON("getAwsVpc failed, because getAwsVpc failed, err: %v, rid: %s", err, kit.Rid)
		return nil, err
	}

	vpcs := make([]metadata.VpcInfo, 0)
	for _, region := range regions {
		sess, err := lgc.AwsNewSession(kit, region, secretID, secretKey)
		if err != nil {
			blog.ErrorJSON("getAwsVpc failed, awsNewSession failed, err: %v, rid: %s", err, kit.Rid)
			return nil, err
		}
		ec2Svc := ec2.New(sess)
		output, err := ec2Svc.DescribeVpcs(nil)
		for _, vpc := range output.Vpcs {
			vpcInfo := metadata.VpcInfo{
				VpcName: *vpc.VpcId,
				VpcID:   *vpc.VpcId,
				Region:  region,
			}
			name := lgc.getAwsVpcName(kit, vpc)
			if name != "" {
				vpcInfo.VpcName = name
			}
			vpcs = append(vpcs, vpcInfo)
		}
	}

	return vpcs, nil
}

func (lgc *Logics) getAwsVpcName(kit *rest.Kit, vpc *ec2.Vpc) string {
	if len(vpc.Tags) <= 0 {
		return ""
	}
	for _, tag := range vpc.Tags {
		if *tag.Key == "Name" {
			return *tag.Value
		}
	}

	return ""
}

func (lgc *Logics) GetTencentCloudVpc(kit *rest.Kit, secretID, secretKey string) ([]metadata.VpcInfo, error) {
	credential := lgc.TencentCloudNewCredential(kit, secretID, secretKey)
	regions, err := lgc.GetTencentCloudRegions(kit, secretID, secretKey)
	if err != nil {
		blog.ErrorJSON("getTencentCloudVpc failed, getTencentCloudRegions failed, err: %v, rid: %s", err, kit.Rid)
		return nil, err
	}

	vpcs := make([]metadata.VpcInfo, 0)
	for _, region := range regions {
		client, err := tcVpc.NewClient(credential, region, profile.NewClientProfile())
		if err != nil {
			blog.ErrorJSON("getTencentCloudVpc new client failed, err: %v, rid: %s", err, kit.Rid)
			return nil, nil
		}

		vpcReq := tcVpc.NewDescribeVpcsRequest()
		vpcResp, err := client.DescribeVpcs(vpcReq)
		for _, vpc := range vpcResp.Response.VpcSet {
			vpcs = append(vpcs, metadata.VpcInfo{
				VpcName: *vpc.VpcName,
				VpcID:   *vpc.VpcId,
				Region:  region,
			})
		}
	}

	return vpcs, nil
}

// GetAwsInstance
// session是需要region的，这样每次都只能拿一个region的instance，根据目前的需求，传filters好像没必要
func (lgc *Logics) GetAwsInstance(kit *rest.Kit, sess *session.Session, filters []*ec2.Filter) (*ec2.DescribeInstancesOutput, error) {
	ec2Svc := ec2.New(sess)
	input := &ec2.DescribeInstancesInput{
		Filters: filters,
	}
	result, err := ec2Svc.DescribeInstances(input)
	if err != nil {
		blog.ErrorJSON("GetAwsInstance failed, aws sdk api http call failed, err: %v, rid: %s", err, kit.Rid)
		return nil, err
	}

	return result, nil
}

// GetTencentCloudInstance
// client的初始化是需要region，所以每次都只能拿到一个region的instance，传filter好像没必要
func (lgc *Logics) GetTencentCloudInstance(kit *rest.Kit, client *cvm.Client, filters []*cvm.Filter) (*cvm.DescribeInstancesResponse, error) {
	instRequest := cvm.NewDescribeInstancesRequest()
	instRequest.Filters = filters
	rsp, err := client.DescribeInstances(instRequest)
	if err != nil {
		blog.ErrorJSON("GetTencentCloudInstance failed, tencent_cloud sdk api http failed, err: %v, rid: %s", err, kit.Rid)
		return nil, err
	}
	return rsp, nil
}

func (lgc *Logics) GetCloudVendorVpc(kit *rest.Kit, account metadata.CloudAccount) ([]metadata.VpcInfo, error) {
	vpc := make([]metadata.VpcInfo, 0)
	var err error

	switch account.CloudVendor {
	case metadata.AWS:
		// todo id、key需要解密
		vpc, err = lgc.GetAwsVpc(kit, account.SecretID, account.SecretKey)
	case metadata.TencentCloud:
		// todo id、key需要解密
		vpc, err = lgc.GetTencentCloudVpc(kit, account.SecretID, account.SecretKey)
	default:
		return nil, kit.CCError.CCError(common.CCErrCloudVendorNotSupport)
	}

	return vpc, err
}
