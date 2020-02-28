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
	"context"
	"fmt"
	"testing"

	"configcenter/src/common/util"
)

const (
	awsSecretID  = ""
	awsSecretKey = ""
	tcSecretID   = ""
	tcSecretKey  = ""
)

func TestGetAwsRegions(t *testing.T) {
	lgc := NewLogics(context.Background(), nil, nil, nil)
	result, err := lgc.GetAwsRegions(nil, awsSecretID, awsSecretKey)
	if err != nil {
		t.Errorf("failed")
	}

	if result == nil {
		t.Errorf("failed")
	}

	expect := []string{"eu-north-1", "ap-south-1", "eu-west-3", "eu-west-2", "eu-west-1", "ap-northeast-2", "ap-northeast-1", "sa-east-1",
		"ca-central-1", "ap-southeast-1", "ap-southeast-2", "eu-central-1", "us-east-1", "us-east-2", "us-west-1", "us-west-2"}

	failed := false
	for _, r := range expect {
		if !util.InStrArr(result, r) {
			failed = true
			break
		}
	}

	if failed {
		t.Fail()
	}
}

func TestGetTencentCloudRegions(t *testing.T) {
	lgc := NewLogics(context.Background(), nil, nil, nil)
	result, err := lgc.GetTencentCloudRegions(nil, tcSecretID, tcSecretKey)
	if err != nil {
		t.Errorf("failed")
	}

	if result == nil {
		t.Errorf("failed")
	}

	fmt.Println(result)

	expect := []string{"ap-bangkok", "ap-beijing", "ap-chengdu", "ap-chongqing", "ap-guangzhou", "ap-guangzhou-open", "ap-hongkong", "ap-mumbai", "ap-seoul",
		"ap-shanghai", "ap-shanghai-fsi", "ap-shenzhen-fsi", "ap-singapore", "eu-frankfurt", "eu-moscow", "na-ashburn", "na-siliconvalley", "na-toronto"}

	failed := false
	for _, r := range expect {
		if !util.InStrArr(result, r) {
			failed = true
			break
		}
	}

	if failed {
		t.Fail()
	}
}

func TestGetAwsVpc(t *testing.T) {
	lgc := NewLogics(context.Background(), nil, nil, nil)
	result, err := lgc.GetAwsVpc(nil, awsSecretID, awsSecretKey)
	if err != nil {
		t.Errorf("failed")
	}

	if result == nil {
		t.Errorf("failed")
	}

	fmt.Println(result)
}

func TestGetTencentCloudVpc(t *testing.T) {
	lgc := NewLogics(context.Background(), nil, nil, nil)
	result, err := lgc.GetTencentCloudVpc(nil, tcSecretID, tcSecretKey)
	if err != nil {
		t.Errorf("failed")
	}

	if result == nil {
		t.Errorf("failed")
	}

	fmt.Println(result)
}
