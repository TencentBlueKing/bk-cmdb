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

func (c *awsClient) SetCredential(secretID, secretKey string) {
	c.secretID = secretID
	c.secretKey = secretKey
}

func (c *awsClient) GetRegions() ([]string, error) {
	sess, err := c.newSession("us-west-1")
	if err != nil {
		return nil, err
	}
	ec2Svc := ec2.New(sess)

	rsp, err := ec2Svc.DescribeRegions(nil)
	if err != nil {
		return nil, err
	}

	regions := make([]string, 0)
	for _, region := range rsp.Regions {
		regions = append(regions, *region.RegionName)
	}

	return regions, nil
}

func (c *awsClient) GetVpcs(region string) ([]*metadata.Vpc, error) {
	sess, err := c.newSession(region)
	if err != nil {
		return nil, err
	}
	ec2Svc := ec2.New(sess)

	output, err := ec2Svc.DescribeVpcs(nil)
	if err != nil {
		return nil, err
	}

	vpcs := make([]*metadata.Vpc, 0)
	for _, vpc := range output.Vpcs {
		vpcInfo := &metadata.Vpc{
			VpcId:   *vpc.VpcId,
			VpcName: *vpc.VpcId,
		}
		name := c.getVpcName(vpc)
		if name != "" {
			vpcInfo.VpcName = name
		}
		vpcs = append(vpcs, vpcInfo)
	}

	return vpcs, nil
}

func (c *awsClient) GetInstances(region string) ([]*metadata.Instance, error) {
	sess, err := c.newSession(region)
	if err != nil {
		return nil, err
	}
	ec2Svc := ec2.New(sess)

	output, err := ec2Svc.DescribeInstances(nil)
	if err != nil {
		return nil, err
	}

	instances := make([]*metadata.Instance, 0)
	for _, reservation := range output.Reservations {
		for _, inst := range reservation.Instances {
			instance := &metadata.Instance{
				InstanceId:    *inst.InstanceId,
				InstanceName:  *inst.InstanceId,
				PrivateIp:     *inst.PrivateIpAddress,
				PublicIp:      *inst.PublicIpAddress,
				InstanceState: *inst.State.Name,
				VpcId:         *inst.VpcId,
			}
			name := c.getInstanceName(inst)
			if name != "" {
				instance.InstanceName = name
			}
			instances = append(instances, instance)
		}

	}

	return instances, nil
}

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

func (c *awsClient) getVpcName(vpc *ec2.Vpc) string {
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

func (c *awsClient) getInstanceName(inst *ec2.Instance) string {
	if len(inst.Tags) <= 0 {
		return ""
	}
	for _, tag := range inst.Tags {
		if *tag.Key == "Name" {
			return *tag.Value
		}
	}
	return ""
}
