package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
)

func main() {
	credential := common.NewCredential(
		os.Getenv("TENCENTCLOUD_SECRET_ID"),
		os.Getenv("TENCENTCLOUD_SECRET_KEY"),
	)

	cpf := profile.NewClientProfile()
	cpf.HttpProfile.ReqMethod = "POST"
	cpf.HttpProfile.ReqTimeout = 5
	//cpf.HttpProfile.Endpoint = "cvm.ap-guangzhou.tencentcloudapi.com"
	cpf.SignMethod = "HmacSHA1"

	client, _ := cvm.NewClient(credential, "ap-guangzhou", cpf)
	request := cvm.NewRunInstancesRequest()
	request.ImageId = common.StringPtr("img-8toqc6s3")
	request.Placement = &cvm.Placement{
		Zone: common.StringPtr("ap-guangzhou-3"),
	}
	request.InstanceChargeType = common.StringPtr("POSTPAID_BY_HOUR")
	request.InstanceCount = common.Int64Ptr(1)
	request.InstanceName = common.StringPtr("API-SDK-GO")
	request.InstanceType = common.StringPtr("S2.SMALL1")
	request.InternetAccessible = &cvm.InternetAccessible{
		InternetChargeType:      common.StringPtr("BANDWIDTH_POSTPAID_BY_HOUR"),
		InternetMaxBandwidthOut: common.Int64Ptr(10),
		PublicIpAssigned:        common.BoolPtr(true),
	}
	request.LoginSettings = &cvm.LoginSettings{
		Password: common.StringPtr("passw0rdExample"),
	}
	request.SecurityGroupIds = common.StringPtrs([]string{"sg-icy671l9"})
	request.SystemDisk = &cvm.SystemDisk{
		DiskSize: common.Int64Ptr(50),
		DiskType: common.StringPtr("CLOUD_BASIC"),
	}
	request.VirtualPrivateCloud = &cvm.VirtualPrivateCloud{
		SubnetId: common.StringPtr("subnet-b1wk8b10"),
		VpcId:    common.StringPtr("vpc-8ek64x3d"),
	}

	// get response structure
	response, err := client.RunInstances(request)
	// API errors
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		fmt.Printf("An API error has returned: %s", err)
		return
	}
	// unexpected errors
	if err != nil {
		panic(err)
	}
	b, _ := json.Marshal(response.Response)
	fmt.Printf("%s", b)
}
