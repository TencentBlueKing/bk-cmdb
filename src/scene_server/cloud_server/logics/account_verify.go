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
	"configcenter/src/common/http/rest"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	tc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/regions"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
)

func (lgc *Logics) AwsAccountVerify(kit *rest.Kit, secretID, secretKey string) (bool, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("us-west-2"),
		Credentials: credentials.NewStaticCredentials(secretID, secretKey, ""),
	})
	if err != nil {
		return false, err
	}
	ec2Svc := ec2.New(sess)

	_, err = ec2Svc.DescribeInstances(nil)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (lgc *Logics) TecentCloudVerify(kit *rest.Kit, secretID, secretKey string) (bool, error) {
	credential := tc.NewCredential(secretID, secretKey)

	cpf := profile.NewClientProfile()
	cpf.HttpProfile.ReqMethod = common.HTTPSelectGet
	cpf.HttpProfile.ReqTimeout = 10
	cpf.HttpProfile.Endpoint = "cvm.tencentcloudapi.com"
	cpf.SignMethod = "HmacSHA1"

	client, err := cvm.NewClient(credential, regions.Guangzhou, cpf)
	if err != nil {
		return false, err
	}
	instRequest := cvm.NewDescribeInstancesRequest()
	_, err = client.DescribeInstances(instRequest)
	if err != nil {
		return false, err
	}

	return true, nil
}
