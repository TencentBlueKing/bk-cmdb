/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.,
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the ",License",); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an ",AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package cloudvendor

import (
	"os"
	"strings"
	"sync"
	"testing"

	"configcenter/src/common/metadata"
	ccom "configcenter/src/scene_server/cloud_server/common"
)

var awsTestClient VendorClient
var defaultAWSRegion string

func init() {
	secretID := os.Getenv("AWS_SECRET_ID")
	secretKey := os.Getenv("AWS_SECRET_KEY")

	// 自动检测是否使用中国区凭证
	defaultAWSRegion = "us-west-1"
	if strings.HasPrefix(secretID, "AKIAZ") {
		// 中国区凭证，使用中国区域
		defaultAWSRegion = "cn-north-1"
	}

	conf := metadata.CloudAccountConf{
		VendorName: metadata.AWS,
		SecretID:   secretID,
		SecretKey:  secretKey,
	}
	var err error
	awsTestClient, err = GetVendorClient(conf)
	if err != nil {
		panic(err.Error())
	}
}

func TestAWSGetRegions(t *testing.T) {
	regionSet, err := awsTestClient.GetRegions()
	if err != nil {
		t.Fatal(err)
	}
	for i, region := range regionSet {
		t.Logf("i:%d, region:%#v\n", i, *region)
	}
}

func TestAWSGetVpcs(t *testing.T) {
	opt := &ccom.VpcOpt{}
	vpcsInfo, err := awsTestClient.GetVpcs(defaultAWSRegion, opt)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("vpcs count:%#v\n", vpcsInfo.Count)
	for i, vpc := range vpcsInfo.VpcSet {
		t.Logf("i:%d, vpc:%#v\n", i, *vpc)
	}
}

func TestAWSGetInstances(t *testing.T) {
	opt := &ccom.InstanceOpt{}
	instancesInfo, err := awsTestClient.GetInstances(defaultAWSRegion, opt)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("instances count:%#v\n", instancesInfo.Count)
	for i, instance := range instancesInfo.InstanceSet {
		t.Logf("i:%d, instance:%#v\n", i, *instance)
	}
}

func TestAWSGetInstancesTotalCnt(t *testing.T) {
	opt := &ccom.InstanceOpt{}
	count, err := awsTestClient.GetInstancesTotalCnt(defaultAWSRegion, opt)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("instances count:%#v\n", count)
}

func TestAWSRequestOpt(t *testing.T) {
	opt := &ccom.VpcOpt{
		BaseOpt: ccom.BaseOpt{
			Filters: []*ccom.Filter{{ccom.StringPtr("tag:Name"), ccom.StringPtrs([]string{"game2"})}},
		},
	}
	vpcsInfo, err := awsTestClient.GetVpcs(defaultAWSRegion, opt)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("vpc count:%#v\n", vpcsInfo.Count)
	for i, vpc := range vpcsInfo.VpcSet {
		t.Logf("i:%d, vpc:%#v\n", i, *vpc)
	}
}

func TestAWSConcurrence(t *testing.T) {
	var wg sync.WaitGroup
	cnt := 10
	wg.Add(cnt)
	for i := 1; i <= cnt; i++ {
		go func(idx int) {
			defer wg.Done()
			opt := &ccom.VpcOpt{
				BaseOpt: ccom.BaseOpt{
					Filters: []*ccom.Filter{{ccom.StringPtr("tag:Name"), ccom.StringPtrs([]string{"game2"})}},
				},
			}
			vpcsInfo, err := awsTestClient.GetVpcs(defaultAWSRegion, opt)
			if err != nil {
				t.Logf("g%d error: %v", idx, err)
				return
			}
			t.Logf("g%d vpcs count:%#v\n", idx, vpcsInfo.Count)
			for i, vpc := range vpcsInfo.VpcSet {
				t.Logf("g%d i:%d, vpc:%#v\n", idx, i, *vpc)
			}
		}(i)
	}
	wg.Wait()
}
