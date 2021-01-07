// Copyright (c) 2017-2018 THL A29 Limited, a Tencent company. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package v20180228

import (
    "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
    tchttp "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/http"
    "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
)

const APIVersion = "2018-02-28"

type Client struct {
    common.Client
}

// Deprecated
func NewClientWithSecretId(secretId, secretKey, region string) (client *Client, err error) {
    cpf := profile.NewClientProfile()
    client = &Client{}
    client.Init(region).WithSecretId(secretId, secretKey).WithProfile(cpf)
    return
}

func NewClient(credential *common.Credential, region string, clientProfile *profile.ClientProfile) (client *Client, err error) {
    client = &Client{}
    client.Init(region).
        WithCredential(credential).
        WithProfile(clientProfile)
    return
}


func NewCreateAccountRequest() (request *CreateAccountRequest) {
    request = &CreateAccountRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("youmall", APIVersion, "CreateAccount")
    return
}

func NewCreateAccountResponse() (response *CreateAccountResponse) {
    response = &CreateAccountResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 创建集团门店管理员账号
func (c *Client) CreateAccount(request *CreateAccountRequest) (response *CreateAccountResponse, err error) {
    if request == nil {
        request = NewCreateAccountRequest()
    }
    response = NewCreateAccountResponse()
    err = c.Send(request, response)
    return
}

func NewCreateFacePictureRequest() (request *CreateFacePictureRequest) {
    request = &CreateFacePictureRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("youmall", APIVersion, "CreateFacePicture")
    return
}

func NewCreateFacePictureResponse() (response *CreateFacePictureResponse) {
    response = &CreateFacePictureResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 通过上传指定规格的人脸图片，创建黑名单用户或者白名单用户。
func (c *Client) CreateFacePicture(request *CreateFacePictureRequest) (response *CreateFacePictureResponse, err error) {
    if request == nil {
        request = NewCreateFacePictureRequest()
    }
    response = NewCreateFacePictureResponse()
    err = c.Send(request, response)
    return
}

func NewDeletePersonFeatureRequest() (request *DeletePersonFeatureRequest) {
    request = &DeletePersonFeatureRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("youmall", APIVersion, "DeletePersonFeature")
    return
}

func NewDeletePersonFeatureResponse() (response *DeletePersonFeatureResponse) {
    response = &DeletePersonFeatureResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 删除顾客特征，仅支持删除黑名单或者白名单用户特征。
func (c *Client) DeletePersonFeature(request *DeletePersonFeatureRequest) (response *DeletePersonFeatureResponse, err error) {
    if request == nil {
        request = NewDeletePersonFeatureRequest()
    }
    response = NewDeletePersonFeatureResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeCameraPersonRequest() (request *DescribeCameraPersonRequest) {
    request = &DescribeCameraPersonRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("youmall", APIVersion, "DescribeCameraPerson")
    return
}

func NewDescribeCameraPersonResponse() (response *DescribeCameraPersonResponse) {
    response = &DescribeCameraPersonResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 通过指定设备ID和指定时段，获取该时段内中收银台摄像设备抓取到顾客头像及身份ID
func (c *Client) DescribeCameraPerson(request *DescribeCameraPersonRequest) (response *DescribeCameraPersonResponse, err error) {
    if request == nil {
        request = NewDescribeCameraPersonRequest()
    }
    response = NewDescribeCameraPersonResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeClusterPersonArrivedMallRequest() (request *DescribeClusterPersonArrivedMallRequest) {
    request = &DescribeClusterPersonArrivedMallRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("youmall", APIVersion, "DescribeClusterPersonArrivedMall")
    return
}

func NewDescribeClusterPersonArrivedMallResponse() (response *DescribeClusterPersonArrivedMallResponse) {
    response = &DescribeClusterPersonArrivedMallResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 输出开始时间到结束时间段内的进出场数据。按天聚合的情况下，每天多次进出场算一次，以最初进场时间为进场时间，最后离场时间为离场时间。停留时间为多次进出场的停留时间之和。
func (c *Client) DescribeClusterPersonArrivedMall(request *DescribeClusterPersonArrivedMallRequest) (response *DescribeClusterPersonArrivedMallResponse, err error) {
    if request == nil {
        request = NewDescribeClusterPersonArrivedMallRequest()
    }
    response = NewDescribeClusterPersonArrivedMallResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeClusterPersonTraceRequest() (request *DescribeClusterPersonTraceRequest) {
    request = &DescribeClusterPersonTraceRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("youmall", APIVersion, "DescribeClusterPersonTrace")
    return
}

func NewDescribeClusterPersonTraceResponse() (response *DescribeClusterPersonTraceResponse) {
    response = &DescribeClusterPersonTraceResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 输出开始时间到结束时间段内的进出场数据。按天聚合的情况下，每天多次进出场算一次，以最初进场时间为进场时间，最后离场时间为离场时间。
func (c *Client) DescribeClusterPersonTrace(request *DescribeClusterPersonTraceRequest) (response *DescribeClusterPersonTraceResponse, err error) {
    if request == nil {
        request = NewDescribeClusterPersonTraceRequest()
    }
    response = NewDescribeClusterPersonTraceResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeFaceIdByTempIdRequest() (request *DescribeFaceIdByTempIdRequest) {
    request = &DescribeFaceIdByTempIdRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("youmall", APIVersion, "DescribeFaceIdByTempId")
    return
}

func NewDescribeFaceIdByTempIdResponse() (response *DescribeFaceIdByTempIdResponse) {
    response = &DescribeFaceIdByTempIdResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 通过DescribeCameraPerson接口上报的收银台身份ID查询顾客的FaceID。查询最佳时间为收银台上报的次日1点后。
func (c *Client) DescribeFaceIdByTempId(request *DescribeFaceIdByTempIdRequest) (response *DescribeFaceIdByTempIdResponse, err error) {
    if request == nil {
        request = NewDescribeFaceIdByTempIdRequest()
    }
    response = NewDescribeFaceIdByTempIdResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeHistoryNetworkInfoRequest() (request *DescribeHistoryNetworkInfoRequest) {
    request = &DescribeHistoryNetworkInfoRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("youmall", APIVersion, "DescribeHistoryNetworkInfo")
    return
}

func NewDescribeHistoryNetworkInfoResponse() (response *DescribeHistoryNetworkInfoResponse) {
    response = &DescribeHistoryNetworkInfoResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 返回当前门店历史网络状态数据
func (c *Client) DescribeHistoryNetworkInfo(request *DescribeHistoryNetworkInfoRequest) (response *DescribeHistoryNetworkInfoResponse, err error) {
    if request == nil {
        request = NewDescribeHistoryNetworkInfoRequest()
    }
    response = NewDescribeHistoryNetworkInfoResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeNetworkInfoRequest() (request *DescribeNetworkInfoRequest) {
    request = &DescribeNetworkInfoRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("youmall", APIVersion, "DescribeNetworkInfo")
    return
}

func NewDescribeNetworkInfoResponse() (response *DescribeNetworkInfoResponse) {
    response = &DescribeNetworkInfoResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 返回当前门店最新网络状态数据
func (c *Client) DescribeNetworkInfo(request *DescribeNetworkInfoRequest) (response *DescribeNetworkInfoResponse, err error) {
    if request == nil {
        request = NewDescribeNetworkInfoRequest()
    }
    response = NewDescribeNetworkInfoResponse()
    err = c.Send(request, response)
    return
}

func NewDescribePersonRequest() (request *DescribePersonRequest) {
    request = &DescribePersonRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("youmall", APIVersion, "DescribePerson")
    return
}

func NewDescribePersonResponse() (response *DescribePersonResponse) {
    response = &DescribePersonResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 查询指定某一卖场的用户信息
func (c *Client) DescribePerson(request *DescribePersonRequest) (response *DescribePersonResponse, err error) {
    if request == nil {
        request = NewDescribePersonRequest()
    }
    response = NewDescribePersonResponse()
    err = c.Send(request, response)
    return
}

func NewDescribePersonArrivedMallRequest() (request *DescribePersonArrivedMallRequest) {
    request = &DescribePersonArrivedMallRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("youmall", APIVersion, "DescribePersonArrivedMall")
    return
}

func NewDescribePersonArrivedMallResponse() (response *DescribePersonArrivedMallResponse) {
    response = &DescribePersonArrivedMallResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 输出开始时间到结束时间段内的进出场数据。不做按天聚合的情况下，每次进出场，产生一条进出场数据。
// 
func (c *Client) DescribePersonArrivedMall(request *DescribePersonArrivedMallRequest) (response *DescribePersonArrivedMallResponse, err error) {
    if request == nil {
        request = NewDescribePersonArrivedMallRequest()
    }
    response = NewDescribePersonArrivedMallResponse()
    err = c.Send(request, response)
    return
}

func NewDescribePersonInfoRequest() (request *DescribePersonInfoRequest) {
    request = &DescribePersonInfoRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("youmall", APIVersion, "DescribePersonInfo")
    return
}

func NewDescribePersonInfoResponse() (response *DescribePersonInfoResponse) {
    response = &DescribePersonInfoResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 指定门店获取所有顾客详情列表，包含客户ID、图片、年龄、性别
func (c *Client) DescribePersonInfo(request *DescribePersonInfoRequest) (response *DescribePersonInfoResponse, err error) {
    if request == nil {
        request = NewDescribePersonInfoRequest()
    }
    response = NewDescribePersonInfoResponse()
    err = c.Send(request, response)
    return
}

func NewDescribePersonInfoByFacePictureRequest() (request *DescribePersonInfoByFacePictureRequest) {
    request = &DescribePersonInfoByFacePictureRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("youmall", APIVersion, "DescribePersonInfoByFacePicture")
    return
}

func NewDescribePersonInfoByFacePictureResponse() (response *DescribePersonInfoByFacePictureResponse) {
    response = &DescribePersonInfoByFacePictureResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 通过上传人脸图片检索系统face id、顾客身份信息及底图
func (c *Client) DescribePersonInfoByFacePicture(request *DescribePersonInfoByFacePictureRequest) (response *DescribePersonInfoByFacePictureResponse, err error) {
    if request == nil {
        request = NewDescribePersonInfoByFacePictureRequest()
    }
    response = NewDescribePersonInfoByFacePictureResponse()
    err = c.Send(request, response)
    return
}

func NewDescribePersonTraceRequest() (request *DescribePersonTraceRequest) {
    request = &DescribePersonTraceRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("youmall", APIVersion, "DescribePersonTrace")
    return
}

func NewDescribePersonTraceResponse() (response *DescribePersonTraceResponse) {
    response = &DescribePersonTraceResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 输出开始时间到结束时间段内的进出场数据。
func (c *Client) DescribePersonTrace(request *DescribePersonTraceRequest) (response *DescribePersonTraceResponse, err error) {
    if request == nil {
        request = NewDescribePersonTraceRequest()
    }
    response = NewDescribePersonTraceResponse()
    err = c.Send(request, response)
    return
}

func NewDescribePersonTraceDetailRequest() (request *DescribePersonTraceDetailRequest) {
    request = &DescribePersonTraceDetailRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("youmall", APIVersion, "DescribePersonTraceDetail")
    return
}

func NewDescribePersonTraceDetailResponse() (response *DescribePersonTraceDetailResponse) {
    response = &DescribePersonTraceDetailResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 查询客户单次到场轨迹明细
func (c *Client) DescribePersonTraceDetail(request *DescribePersonTraceDetailRequest) (response *DescribePersonTraceDetailResponse, err error) {
    if request == nil {
        request = NewDescribePersonTraceDetailRequest()
    }
    response = NewDescribePersonTraceDetailResponse()
    err = c.Send(request, response)
    return
}

func NewDescribePersonVisitInfoRequest() (request *DescribePersonVisitInfoRequest) {
    request = &DescribePersonVisitInfoRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("youmall", APIVersion, "DescribePersonVisitInfo")
    return
}

func NewDescribePersonVisitInfoResponse() (response *DescribePersonVisitInfoResponse) {
    response = &DescribePersonVisitInfoResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 获取门店指定时间范围内的所有用户到访信息记录，支持的时间范围：过去365天，含当天。
func (c *Client) DescribePersonVisitInfo(request *DescribePersonVisitInfoRequest) (response *DescribePersonVisitInfoResponse, err error) {
    if request == nil {
        request = NewDescribePersonVisitInfoRequest()
    }
    response = NewDescribePersonVisitInfoResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeShopHourTrafficInfoRequest() (request *DescribeShopHourTrafficInfoRequest) {
    request = &DescribeShopHourTrafficInfoRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("youmall", APIVersion, "DescribeShopHourTrafficInfo")
    return
}

func NewDescribeShopHourTrafficInfoResponse() (response *DescribeShopHourTrafficInfoResponse) {
    response = &DescribeShopHourTrafficInfoResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 按小时提供查询日期范围内门店的每天每小时累计客流人数数据，支持的时间范围：过去365天，含当天。
func (c *Client) DescribeShopHourTrafficInfo(request *DescribeShopHourTrafficInfoRequest) (response *DescribeShopHourTrafficInfoResponse, err error) {
    if request == nil {
        request = NewDescribeShopHourTrafficInfoRequest()
    }
    response = NewDescribeShopHourTrafficInfoResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeShopInfoRequest() (request *DescribeShopInfoRequest) {
    request = &DescribeShopInfoRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("youmall", APIVersion, "DescribeShopInfo")
    return
}

func NewDescribeShopInfoResponse() (response *DescribeShopInfoResponse) {
    response = &DescribeShopInfoResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 根据客户身份标识获取客户下所有的门店信息列表
func (c *Client) DescribeShopInfo(request *DescribeShopInfoRequest) (response *DescribeShopInfoResponse, err error) {
    if request == nil {
        request = NewDescribeShopInfoRequest()
    }
    response = NewDescribeShopInfoResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeShopTrafficInfoRequest() (request *DescribeShopTrafficInfoRequest) {
    request = &DescribeShopTrafficInfoRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("youmall", APIVersion, "DescribeShopTrafficInfo")
    return
}

func NewDescribeShopTrafficInfoResponse() (response *DescribeShopTrafficInfoResponse) {
    response = &DescribeShopTrafficInfoResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 按天提供查询日期范围内门店的单日累计客流人数，支持的时间范围：过去365天，含当天。
func (c *Client) DescribeShopTrafficInfo(request *DescribeShopTrafficInfoRequest) (response *DescribeShopTrafficInfoResponse, err error) {
    if request == nil {
        request = NewDescribeShopTrafficInfoRequest()
    }
    response = NewDescribeShopTrafficInfoResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeTrajectoryDataRequest() (request *DescribeTrajectoryDataRequest) {
    request = &DescribeTrajectoryDataRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("youmall", APIVersion, "DescribeTrajectoryData")
    return
}

func NewDescribeTrajectoryDataResponse() (response *DescribeTrajectoryDataResponse) {
    response = &DescribeTrajectoryDataResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 获取动线轨迹信息
func (c *Client) DescribeTrajectoryData(request *DescribeTrajectoryDataRequest) (response *DescribeTrajectoryDataResponse, err error) {
    if request == nil {
        request = NewDescribeTrajectoryDataRequest()
    }
    response = NewDescribeTrajectoryDataResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeZoneFlowAgeInfoByZoneIdRequest() (request *DescribeZoneFlowAgeInfoByZoneIdRequest) {
    request = &DescribeZoneFlowAgeInfoByZoneIdRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("youmall", APIVersion, "DescribeZoneFlowAgeInfoByZoneId")
    return
}

func NewDescribeZoneFlowAgeInfoByZoneIdResponse() (response *DescribeZoneFlowAgeInfoByZoneIdResponse) {
    response = &DescribeZoneFlowAgeInfoByZoneIdResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 获取指定区域人流各年龄占比
func (c *Client) DescribeZoneFlowAgeInfoByZoneId(request *DescribeZoneFlowAgeInfoByZoneIdRequest) (response *DescribeZoneFlowAgeInfoByZoneIdResponse, err error) {
    if request == nil {
        request = NewDescribeZoneFlowAgeInfoByZoneIdRequest()
    }
    response = NewDescribeZoneFlowAgeInfoByZoneIdResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeZoneFlowAndStayTimeRequest() (request *DescribeZoneFlowAndStayTimeRequest) {
    request = &DescribeZoneFlowAndStayTimeRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("youmall", APIVersion, "DescribeZoneFlowAndStayTime")
    return
}

func NewDescribeZoneFlowAndStayTimeResponse() (response *DescribeZoneFlowAndStayTimeResponse) {
    response = &DescribeZoneFlowAndStayTimeResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 获取区域人流和停留时间
func (c *Client) DescribeZoneFlowAndStayTime(request *DescribeZoneFlowAndStayTimeRequest) (response *DescribeZoneFlowAndStayTimeResponse, err error) {
    if request == nil {
        request = NewDescribeZoneFlowAndStayTimeRequest()
    }
    response = NewDescribeZoneFlowAndStayTimeResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeZoneFlowDailyByZoneIdRequest() (request *DescribeZoneFlowDailyByZoneIdRequest) {
    request = &DescribeZoneFlowDailyByZoneIdRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("youmall", APIVersion, "DescribeZoneFlowDailyByZoneId")
    return
}

func NewDescribeZoneFlowDailyByZoneIdResponse() (response *DescribeZoneFlowDailyByZoneIdResponse) {
    response = &DescribeZoneFlowDailyByZoneIdResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 获取指定区域每日客流量
func (c *Client) DescribeZoneFlowDailyByZoneId(request *DescribeZoneFlowDailyByZoneIdRequest) (response *DescribeZoneFlowDailyByZoneIdResponse, err error) {
    if request == nil {
        request = NewDescribeZoneFlowDailyByZoneIdRequest()
    }
    response = NewDescribeZoneFlowDailyByZoneIdResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeZoneFlowGenderAvrStayTimeByZoneIdRequest() (request *DescribeZoneFlowGenderAvrStayTimeByZoneIdRequest) {
    request = &DescribeZoneFlowGenderAvrStayTimeByZoneIdRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("youmall", APIVersion, "DescribeZoneFlowGenderAvrStayTimeByZoneId")
    return
}

func NewDescribeZoneFlowGenderAvrStayTimeByZoneIdResponse() (response *DescribeZoneFlowGenderAvrStayTimeByZoneIdResponse) {
    response = &DescribeZoneFlowGenderAvrStayTimeByZoneIdResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 获取指定区域不同年龄段男女平均停留时间
func (c *Client) DescribeZoneFlowGenderAvrStayTimeByZoneId(request *DescribeZoneFlowGenderAvrStayTimeByZoneIdRequest) (response *DescribeZoneFlowGenderAvrStayTimeByZoneIdResponse, err error) {
    if request == nil {
        request = NewDescribeZoneFlowGenderAvrStayTimeByZoneIdRequest()
    }
    response = NewDescribeZoneFlowGenderAvrStayTimeByZoneIdResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeZoneFlowGenderInfoByZoneIdRequest() (request *DescribeZoneFlowGenderInfoByZoneIdRequest) {
    request = &DescribeZoneFlowGenderInfoByZoneIdRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("youmall", APIVersion, "DescribeZoneFlowGenderInfoByZoneId")
    return
}

func NewDescribeZoneFlowGenderInfoByZoneIdResponse() (response *DescribeZoneFlowGenderInfoByZoneIdResponse) {
    response = &DescribeZoneFlowGenderInfoByZoneIdResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 获取指定区域性别占比
func (c *Client) DescribeZoneFlowGenderInfoByZoneId(request *DescribeZoneFlowGenderInfoByZoneIdRequest) (response *DescribeZoneFlowGenderInfoByZoneIdResponse, err error) {
    if request == nil {
        request = NewDescribeZoneFlowGenderInfoByZoneIdRequest()
    }
    response = NewDescribeZoneFlowGenderInfoByZoneIdResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeZoneFlowHourlyByZoneIdRequest() (request *DescribeZoneFlowHourlyByZoneIdRequest) {
    request = &DescribeZoneFlowHourlyByZoneIdRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("youmall", APIVersion, "DescribeZoneFlowHourlyByZoneId")
    return
}

func NewDescribeZoneFlowHourlyByZoneIdResponse() (response *DescribeZoneFlowHourlyByZoneIdResponse) {
    response = &DescribeZoneFlowHourlyByZoneIdResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 获取指定区域分时客流量
func (c *Client) DescribeZoneFlowHourlyByZoneId(request *DescribeZoneFlowHourlyByZoneIdRequest) (response *DescribeZoneFlowHourlyByZoneIdResponse, err error) {
    if request == nil {
        request = NewDescribeZoneFlowHourlyByZoneIdRequest()
    }
    response = NewDescribeZoneFlowHourlyByZoneIdResponse()
    err = c.Send(request, response)
    return
}

func NewDescribeZoneTrafficInfoRequest() (request *DescribeZoneTrafficInfoRequest) {
    request = &DescribeZoneTrafficInfoRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("youmall", APIVersion, "DescribeZoneTrafficInfo")
    return
}

func NewDescribeZoneTrafficInfoResponse() (response *DescribeZoneTrafficInfoResponse) {
    response = &DescribeZoneTrafficInfoResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 按天提供查询日期范围内，客户指定门店下的所有区域（优Mall部署时已配置区域）的累计客流人次和平均停留时间。支持的时间范围：过去365天，含当天。
func (c *Client) DescribeZoneTrafficInfo(request *DescribeZoneTrafficInfoRequest) (response *DescribeZoneTrafficInfoResponse, err error) {
    if request == nil {
        request = NewDescribeZoneTrafficInfoRequest()
    }
    response = NewDescribeZoneTrafficInfoResponse()
    err = c.Send(request, response)
    return
}

func NewModifyPersonTagInfoRequest() (request *ModifyPersonTagInfoRequest) {
    request = &ModifyPersonTagInfoRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("youmall", APIVersion, "ModifyPersonTagInfo")
    return
}

func NewModifyPersonTagInfoResponse() (response *ModifyPersonTagInfoResponse) {
    response = &ModifyPersonTagInfoResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 标记到店顾客的身份类型，例如黑名单、白名单等
func (c *Client) ModifyPersonTagInfo(request *ModifyPersonTagInfoRequest) (response *ModifyPersonTagInfoResponse, err error) {
    if request == nil {
        request = NewModifyPersonTagInfoRequest()
    }
    response = NewModifyPersonTagInfoResponse()
    err = c.Send(request, response)
    return
}

func NewModifyPersonTypeRequest() (request *ModifyPersonTypeRequest) {
    request = &ModifyPersonTypeRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("youmall", APIVersion, "ModifyPersonType")
    return
}

func NewModifyPersonTypeResponse() (response *ModifyPersonTypeResponse) {
    response = &ModifyPersonTypeResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 修改顾客身份类型接口
func (c *Client) ModifyPersonType(request *ModifyPersonTypeRequest) (response *ModifyPersonTypeResponse, err error) {
    if request == nil {
        request = NewModifyPersonTypeRequest()
    }
    response = NewModifyPersonTypeResponse()
    err = c.Send(request, response)
    return
}

func NewRegisterCallbackRequest() (request *RegisterCallbackRequest) {
    request = &RegisterCallbackRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("youmall", APIVersion, "RegisterCallback")
    return
}

func NewRegisterCallbackResponse() (response *RegisterCallbackResponse) {
    response = &RegisterCallbackResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 调用本接口在优Mall中注册自己集团的到店通知回调接口地址，接口协议为HTTP或HTTPS。注册后，若集团有特殊身份（例如老客）到店通知，优Mall后台将主动将到店信息push给该接口
func (c *Client) RegisterCallback(request *RegisterCallbackRequest) (response *RegisterCallbackResponse, err error) {
    if request == nil {
        request = NewRegisterCallbackRequest()
    }
    response = NewRegisterCallbackResponse()
    err = c.Send(request, response)
    return
}
