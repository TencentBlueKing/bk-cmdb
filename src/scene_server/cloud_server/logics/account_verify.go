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
	"configcenter/src/common/http/rest"
	"github.com/aws/aws-sdk-go/service/ec2"
	tc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/regions"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
)

func (lgc *Logics) AwsAccountVerify(kit *rest.Kit, secretID, secretKey string) (bool, error) {
	sess, err := lgc.AwsNewSession(kit, "", secretID, secretKey)
	if err != nil {
		return false, nil
	}
	ec2Svc := ec2.New(sess)

	_, err = ec2Svc.DescribeInstances(nil)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (lgc *Logics) TencentCloudVerify(kit *rest.Kit, secretID, secretKey string) (bool, error) {
	credential := tc.NewCredential(secretID, secretKey)
	client, _ := cvm.NewClient(credential, regions.Guangzhou, profile.NewClientProfile())
	instRequest := cvm.NewDescribeInstancesRequest()
	_, err := client.DescribeInstances(instRequest)
	if err != nil {
		return false, err
	}

	return true, nil
}
