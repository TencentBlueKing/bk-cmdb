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
	"testing"

	"configcenter/src/common/metadata"
	ccom "configcenter/src/scene_server/cloud_server/common"
)

var awsTestClient VendorClient

func init() {
	conf := metadata.CloudAccountConf{
		metadata.AWS,
		os.Getenv("AWS_SECRET_ID"),
		os.Getenv("AWS_SECRET_KEY"),
	}
	var err error
	awsTestClient, err = GetVendorClient(conf)
	if err != nil {
		panic(err.Error())
	}
}

func TestAWSGetRegions(t *testing.T) {
	opt := &ccom.RequestOpt{}
	regionSet, err := awsTestClient.GetRegions(opt)
	if err != nil {
		t.Fatal(err)
	}
	for i, region := range regionSet {
		t.Logf("i:%d, region:%#v\n", i, *region)
	}
}

func TestAWSGetVpcs(t *testing.T) {
	opt := &ccom.RequestOpt{}
	region := "us-west-1"
	vpcsInfo, err := awsTestClient.GetVpcs(region, opt)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("vpcs count:%#v\n", vpcsInfo.Count)
	for i, vpc := range vpcsInfo.VpcSet {
		t.Logf("i:%d, vpc:%#v\n", i, *vpc)
	}
}

func TestAWSGetInstances(t *testing.T) {
	opt := &ccom.RequestOpt{}
	region := "us-west-1"
	instancesInfo, err := awsTestClient.GetInstances(region, opt)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("instances count:%#v\n", instancesInfo.Count)
	for i, instance := range instancesInfo.InstanceSet {
		t.Logf("i:%d, instance:%#v\n", i, *instance)
	}
}

func TestAWSGetInstancesTotalCnt(t *testing.T) {
	opt := &ccom.RequestOpt{}
	region := "us-west-1"
	count, err := awsTestClient.GetInstancesTotalCnt(region, opt)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("instances count:%#v\n", count)
}

func TestAWSRequestOpt(t *testing.T) {
	//opt := &ccom.RequestOpt{Limit: ccom.Int64Ptr(int64(1))}
	opt := &ccom.RequestOpt{
		Filters: []*ccom.Filter{&ccom.Filter{ccom.StringPtr("tag:Name"), ccom.StringPtrs([]string{"game2"})}},
	}
	region := "us-west-1"
	vpcsInfo, err := awsTestClient.GetVpcs(region, opt)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("vpc count:%#v\n", vpcsInfo.Count)
	for i, vpc := range vpcsInfo.VpcSet {
		t.Logf("i:%d, vpc:%#v\n",i, *vpc)
	}
}