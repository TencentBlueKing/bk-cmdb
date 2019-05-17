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
    "encoding/json"

    tchttp "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/http"
)

type ArrivedMallInfo struct {

	// 到场时间
	ArrivedTime *string `json:"ArrivedTime" name:"ArrivedTime"`

	// 出场时间
	LeaveTime *string `json:"LeaveTime" name:"LeaveTime"`

	// 停留时间，秒
	StaySecond *uint64 `json:"StaySecond" name:"StaySecond"`

	// 到场抓拍图片
	InCapPic *string `json:"InCapPic" name:"InCapPic"`

	// 出场抓拍图片
	OutCapPic *string `json:"OutCapPic" name:"OutCapPic"`

	// 轨迹编码
	TraceId *string `json:"TraceId" name:"TraceId"`
}

type CameraPersonInfo struct {

	// 临时id，还未生成face id时返回
	TempId *string `json:"TempId" name:"TempId"`

	// 人脸face id
	FaceId *int64 `json:"FaceId" name:"FaceId"`

	// 确定当次返回的哪个id有效，1-FaceId，2-TempId
	IdType *int64 `json:"IdType" name:"IdType"`

	// 当次抓拍到的人脸图片base编码
	FacePic *string `json:"FacePic" name:"FacePic"`

	// 当次抓拍时间戳
	Time *int64 `json:"Time" name:"Time"`
}

type CreateAccountRequest struct {
	*tchttp.BaseRequest

	// 集团ID
	CompanyId *string `json:"CompanyId" name:"CompanyId"`

	// 账号名；需要是手机号
	Name *string `json:"Name" name:"Name"`

	// 密码；需要是(`~!@#$%^&*()_+=-）中的至少两种且八位以上
	Password *string `json:"Password" name:"Password"`

	// 客户门店编码
	ShopCode *string `json:"ShopCode" name:"ShopCode"`

	// 备注说明; 30个字符以内
	Remark *string `json:"Remark" name:"Remark"`
}

func (r *CreateAccountRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *CreateAccountRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type CreateAccountResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *CreateAccountResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *CreateAccountResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type CreateFacePictureRequest struct {
	*tchttp.BaseRequest

	// 集团ID
	CompanyId *string `json:"CompanyId" name:"CompanyId"`

	// 人物类型（0表示普通顾客，1 白名单，2 表示黑名单，101表示集团白名单，102表示集团黑名单）
	PersonType *int64 `json:"PersonType" name:"PersonType"`

	// 图片BASE编码
	Picture *string `json:"Picture" name:"Picture"`

	// 图片名称
	PictureName *string `json:"PictureName" name:"PictureName"`

	// 店铺ID，如果不填表示操作集团身份库
	ShopId *int64 `json:"ShopId" name:"ShopId"`

	// 是否强制更新：为ture时会为用户创建一个新的指定PersonType的身份;目前这个参数已废弃，可不传
	IsForceUpload *bool `json:"IsForceUpload" name:"IsForceUpload"`
}

func (r *CreateFacePictureRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *CreateFacePictureRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type CreateFacePictureResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 人物ID
		PersonId *int64 `json:"PersonId" name:"PersonId"`

		// 0.正常建档 1.重复身份 2.未检测到人脸 3.检测到多个人脸 4.人脸大小过小 5.人脸质量不达标 6.其他错误
		Status *int64 `json:"Status" name:"Status"`

		// 图片url
		PictureUrl *string `json:"PictureUrl" name:"PictureUrl"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *CreateFacePictureResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *CreateFacePictureResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DailyTracePoint struct {

	// 轨迹日期
	TraceDate *string `json:"TraceDate" name:"TraceDate"`

	// 轨迹点序列
	TracePointSet []*PersonTracePoint `json:"TracePointSet" name:"TracePointSet" list`
}

type DeletePersonFeatureRequest struct {
	*tchttp.BaseRequest

	// 公司ID
	CompanyId *string `json:"CompanyId" name:"CompanyId"`

	// 门店ID
	ShopId *int64 `json:"ShopId" name:"ShopId"`

	// 顾客ID
	PersonId *int64 `json:"PersonId" name:"PersonId"`
}

func (r *DeletePersonFeatureRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DeletePersonFeatureRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DeletePersonFeatureResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DeletePersonFeatureResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DeletePersonFeatureResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeCameraPersonRequest struct {
	*tchttp.BaseRequest

	// 优mall集团id，通过"指定身份标识获取客户门店列表"接口获取
	CompanyId *string `json:"CompanyId" name:"CompanyId"`

	// 优mall店铺id，通过"指定身份标识获取客户门店列表"接口获取
	ShopId *int64 `json:"ShopId" name:"ShopId"`

	// 摄像头id
	CameraId *int64 `json:"CameraId" name:"CameraId"`

	// 拉取开始时间戳，单位秒
	StartTime *int64 `json:"StartTime" name:"StartTime"`

	// 拉取结束时间戳，单位秒，不超过StartTime+10秒，超过默认为StartTime+10
	EndTime *int64 `json:"EndTime" name:"EndTime"`

	// pos机id
	PosId *string `json:"PosId" name:"PosId"`

	// 拉取图片数，默认为1，最大为3
	Num *int64 `json:"Num" name:"Num"`

	// 是否需要base64的图片，0-不需要，1-需要，默认0
	IsNeedPic *int64 `json:"IsNeedPic" name:"IsNeedPic"`
}

func (r *DescribeCameraPersonRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeCameraPersonRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeCameraPersonResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 集团id
		CompanyId *string `json:"CompanyId" name:"CompanyId"`

		// 店铺id
		ShopId *int64 `json:"ShopId" name:"ShopId"`

		// 摄像机id
		CameraId *int64 `json:"CameraId" name:"CameraId"`

		// pos机id
		PosId *string `json:"PosId" name:"PosId"`

		// 抓取的顾客信息
		Infos []*CameraPersonInfo `json:"Infos" name:"Infos" list`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeCameraPersonResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeCameraPersonResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeClusterPersonArrivedMallRequest struct {
	*tchttp.BaseRequest

	// 卖场编码
	MallId *string `json:"MallId" name:"MallId"`

	// 客户编码
	PersonId *string `json:"PersonId" name:"PersonId"`

	// 查询开始时间
	StartTime *string `json:"StartTime" name:"StartTime"`

	// 查询结束时间
	EndTime *string `json:"EndTime" name:"EndTime"`
}

func (r *DescribeClusterPersonArrivedMallRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeClusterPersonArrivedMallRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeClusterPersonArrivedMallResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 卖场系统编码
		MallId *string `json:"MallId" name:"MallId"`

		// 卖场客户编码
		MallCode *string `json:"MallCode" name:"MallCode"`

		// 客户编码
		PersonId *string `json:"PersonId" name:"PersonId"`

		// 到场信息
		ArrivedMallSet []*ArrivedMallInfo `json:"ArrivedMallSet" name:"ArrivedMallSet" list`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeClusterPersonArrivedMallResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeClusterPersonArrivedMallResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeClusterPersonTraceRequest struct {
	*tchttp.BaseRequest

	// 卖场编码
	MallId *string `json:"MallId" name:"MallId"`

	// 客户编码
	PersonId *string `json:"PersonId" name:"PersonId"`

	// 查询开始时间
	StartTime *string `json:"StartTime" name:"StartTime"`

	// 查询结束时间
	EndTime *string `json:"EndTime" name:"EndTime"`
}

func (r *DescribeClusterPersonTraceRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeClusterPersonTraceRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeClusterPersonTraceResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 卖场系统编码
		MallId *string `json:"MallId" name:"MallId"`

		// 卖场用户编码
		MallCode *string `json:"MallCode" name:"MallCode"`

		// 客户编码
		PersonId *string `json:"PersonId" name:"PersonId"`

		// 轨迹序列
		TracePointSet []*DailyTracePoint `json:"TracePointSet" name:"TracePointSet" list`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeClusterPersonTraceResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeClusterPersonTraceResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeFaceIdByTempIdRequest struct {
	*tchttp.BaseRequest

	// 优mall集团id，通过"指定身份标识获取客户门店列表"接口获取
	CompanyId *string `json:"CompanyId" name:"CompanyId"`

	// 优mall店铺id，通过"指定身份标识获取客户门店列表"接口获取
	ShopId *int64 `json:"ShopId" name:"ShopId"`

	// 临时id
	TempId *string `json:"TempId" name:"TempId"`

	// 摄像头id
	CameraId *int64 `json:"CameraId" name:"CameraId"`

	// pos机id
	PosId *string `json:"PosId" name:"PosId"`

	// 图片url过期时间：在当前时间+PictureExpires秒后，图片url无法继续正常访问；单位s；默认值1*24*60*60（1天）
	PictureExpires *int64 `json:"PictureExpires" name:"PictureExpires"`
}

func (r *DescribeFaceIdByTempIdRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeFaceIdByTempIdRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeFaceIdByTempIdResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 集团id
		CompanyId *string `json:"CompanyId" name:"CompanyId"`

		// 店铺id
		ShopId *int64 `json:"ShopId" name:"ShopId"`

		// 摄像机id
		CameraId *int64 `json:"CameraId" name:"CameraId"`

		// pos机id
		PosId *string `json:"PosId" name:"PosId"`

		// 请求的临时id
		TempId *string `json:"TempId" name:"TempId"`

		// 临时id对应的face id
		FaceId *int64 `json:"FaceId" name:"FaceId"`

		// 顾客属性信息
		PersonInfo *PersonInfo `json:"PersonInfo" name:"PersonInfo"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeFaceIdByTempIdResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeFaceIdByTempIdResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeHistoryNetworkInfoRequest struct {
	*tchttp.BaseRequest

	// 请求时间戳
	Time *int64 `json:"Time" name:"Time"`

	// 优mall集团id，通过"指定身份标识获取客户门店列表"接口获取
	CompanyId *string `json:"CompanyId" name:"CompanyId"`

	// 优mall店铺id，通过"指定身份标识获取客户门店列表"接口获取，为0则拉取集团全部店铺当前
	ShopId *int64 `json:"ShopId" name:"ShopId"`

	// 拉取开始日期，格式：2018-09-05
	StartDay *string `json:"StartDay" name:"StartDay"`

	// 拉取结束日期，格式L:2018-09-05，超过StartDay 90天，按StartDay+90天算
	EndDay *string `json:"EndDay" name:"EndDay"`

	// 拉取条数，默认10
	Limit *int64 `json:"Limit" name:"Limit"`

	// 拉取偏移，返回offset之后的数据
	Offset *int64 `json:"Offset" name:"Offset"`
}

func (r *DescribeHistoryNetworkInfoRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeHistoryNetworkInfoRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeHistoryNetworkInfoResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 网络状态数据
		InstanceSet *NetworkHistoryInfo `json:"InstanceSet" name:"InstanceSet"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeHistoryNetworkInfoResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeHistoryNetworkInfoResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeNetworkInfoRequest struct {
	*tchttp.BaseRequest

	// 请求时间戳
	Time *int64 `json:"Time" name:"Time"`

	// 优mall集团id，通过"指定身份标识获取客户门店列表"接口获取
	CompanyId *string `json:"CompanyId" name:"CompanyId"`

	// 优mall店铺id，通过"指定身份标识获取客户门店列表"接口获取，不填则拉取集团全部店铺当前
	ShopId *int64 `json:"ShopId" name:"ShopId"`
}

func (r *DescribeNetworkInfoRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeNetworkInfoRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeNetworkInfoResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 网络状态详情
		InstanceSet *NetworkLastInfo `json:"InstanceSet" name:"InstanceSet"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeNetworkInfoResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeNetworkInfoResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribePersonArrivedMallRequest struct {
	*tchttp.BaseRequest

	// 卖场编码
	MallId *string `json:"MallId" name:"MallId"`

	// 客户编码
	PersonId *string `json:"PersonId" name:"PersonId"`

	// 查询开始时间
	StartTime *string `json:"StartTime" name:"StartTime"`

	// 查询结束时间
	EndTime *string `json:"EndTime" name:"EndTime"`
}

func (r *DescribePersonArrivedMallRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribePersonArrivedMallRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribePersonArrivedMallResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 卖场系统编码
		MallId *string `json:"MallId" name:"MallId"`

		// 卖场用户编码
		MallCode *string `json:"MallCode" name:"MallCode"`

		// 客户编码
		PersonId *string `json:"PersonId" name:"PersonId"`

		// 到场轨迹
		ArrivedMallSet []*ArrivedMallInfo `json:"ArrivedMallSet" name:"ArrivedMallSet" list`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribePersonArrivedMallResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribePersonArrivedMallResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribePersonInfoByFacePictureRequest struct {
	*tchttp.BaseRequest

	// 优mall集团id，通过"指定身份标识获取客户门店列表"接口获取
	CompanyId *string `json:"CompanyId" name:"CompanyId"`

	// 优mall店铺id，通过"指定身份标识获取客户门店列表"接口获取
	ShopId *int64 `json:"ShopId" name:"ShopId"`

	// 人脸图片BASE编码
	Picture *string `json:"Picture" name:"Picture"`
}

func (r *DescribePersonInfoByFacePictureRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribePersonInfoByFacePictureRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribePersonInfoByFacePictureResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 集团id
		CompanyId *string `json:"CompanyId" name:"CompanyId"`

		// 店铺id
		ShopId *int64 `json:"ShopId" name:"ShopId"`

		// 顾客face id
		PersonId *int64 `json:"PersonId" name:"PersonId"`

		// 顾客底图url
		PictureUrl *string `json:"PictureUrl" name:"PictureUrl"`

		// 顾客类型（0表示普通顾客，1 白名单，2 表示黑名单，101表示集团白名单，102表示集团黑名单）
		PersonType *int64 `json:"PersonType" name:"PersonType"`

		// 顾客首次进店时间
		FirstVisitTime *string `json:"FirstVisitTime" name:"FirstVisitTime"`

		// 顾客历史到访次数
		VisitTimes *int64 `json:"VisitTimes" name:"VisitTimes"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribePersonInfoByFacePictureResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribePersonInfoByFacePictureResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribePersonInfoRequest struct {
	*tchttp.BaseRequest

	// 公司ID
	CompanyId *string `json:"CompanyId" name:"CompanyId"`

	// 门店ID
	ShopId *uint64 `json:"ShopId" name:"ShopId"`

	// 起始ID，第一次拉取时StartPersonId传0，后续送入的值为上一页最后一条数据项的PersonId
	StartPersonId *uint64 `json:"StartPersonId" name:"StartPersonId"`

	// 偏移量：分页控制参数，第一页传0，第n页Offset=(n-1)*Limit
	Offset *uint64 `json:"Offset" name:"Offset"`

	// Limit:每页的数据项，最大100，超过100会被强制指定为100
	Limit *uint64 `json:"Limit" name:"Limit"`

	// 图片url过期时间：在当前时间+PictureExpires秒后，图片url无法继续正常访问；单位s；默认值1*24*60*60（1天）
	PictureExpires *uint64 `json:"PictureExpires" name:"PictureExpires"`

	// 身份类型(0表示普通顾客，1 白名单，2 表示黑名单）
	PersonType *uint64 `json:"PersonType" name:"PersonType"`
}

func (r *DescribePersonInfoRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribePersonInfoRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribePersonInfoResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 公司ID
		CompanyId *string `json:"CompanyId" name:"CompanyId"`

		// 门店ID
		ShopId *uint64 `json:"ShopId" name:"ShopId"`

		// 总数
		TotalCount *uint64 `json:"TotalCount" name:"TotalCount"`

		// 用户信息
		PersonInfoSet []*PersonInfo `json:"PersonInfoSet" name:"PersonInfoSet" list`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribePersonInfoResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribePersonInfoResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribePersonRequest struct {
	*tchttp.BaseRequest

	// 卖场编码
	MallId *string `json:"MallId" name:"MallId"`

	// 查询偏移
	Offset *uint64 `json:"Offset" name:"Offset"`

	// 查询数量，默认20，最大查询数量100
	Limit *uint64 `json:"Limit" name:"Limit"`
}

func (r *DescribePersonRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribePersonRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribePersonResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 总计客户数量
		TotalCount *uint64 `json:"TotalCount" name:"TotalCount"`

		// 客户信息
		PersonSet []*PersonProfile `json:"PersonSet" name:"PersonSet" list`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribePersonResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribePersonResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribePersonTraceDetailRequest struct {
	*tchttp.BaseRequest

	// 卖场编码
	MallId *string `json:"MallId" name:"MallId"`

	// 客户编码
	PersonId *string `json:"PersonId" name:"PersonId"`

	// 轨迹编码
	TraceId *string `json:"TraceId" name:"TraceId"`
}

func (r *DescribePersonTraceDetailRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribePersonTraceDetailRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribePersonTraceDetailResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 卖场编码
		MallId *string `json:"MallId" name:"MallId"`

		// 客户编码
		PersonId *string `json:"PersonId" name:"PersonId"`

		// 轨迹编码
		TraceId *string `json:"TraceId" name:"TraceId"`

		// 轨迹点坐标序列
		CoordinateSet []*PersonCoordinate `json:"CoordinateSet" name:"CoordinateSet" list`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribePersonTraceDetailResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribePersonTraceDetailResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribePersonTraceRequest struct {
	*tchttp.BaseRequest

	// 卖场编码
	MallId *string `json:"MallId" name:"MallId"`

	// 客户编码
	PersonId *string `json:"PersonId" name:"PersonId"`

	// 查询开始时间
	StartTime *string `json:"StartTime" name:"StartTime"`

	// 查询结束时间
	EndTime *string `json:"EndTime" name:"EndTime"`
}

func (r *DescribePersonTraceRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribePersonTraceRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribePersonTraceResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 卖场系统编码
		MallId *string `json:"MallId" name:"MallId"`

		// 卖场用户编码
		MallCode *string `json:"MallCode" name:"MallCode"`

		// 客户编码
		PersonId *string `json:"PersonId" name:"PersonId"`

		// 轨迹列表
		TraceRouteSet []*PersonTraceRoute `json:"TraceRouteSet" name:"TraceRouteSet" list`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribePersonTraceResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribePersonTraceResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribePersonVisitInfoRequest struct {
	*tchttp.BaseRequest

	// 公司ID
	CompanyId *string `json:"CompanyId" name:"CompanyId"`

	// 门店ID
	ShopId *uint64 `json:"ShopId" name:"ShopId"`

	// 偏移量：分页控制参数，第一页传0，第n页Offset=(n-1)*Limit
	Offset *uint64 `json:"Offset" name:"Offset"`

	// Limit:每页的数据项，最大100，超过100会被强制指定为100
	Limit *uint64 `json:"Limit" name:"Limit"`

	// 开始日期，格式yyyy-MM-dd，已废弃，请使用StartDateTime
	StartDate *string `json:"StartDate" name:"StartDate"`

	// 结束日期，格式yyyy-MM-dd，已废弃，请使用EndDateTime
	EndDate *string `json:"EndDate" name:"EndDate"`

	// 图片url过期时间：在当前时间+PictureExpires秒后，图片url无法继续正常访问；单位s；默认值1*24*60*60（1天）
	PictureExpires *uint64 `json:"PictureExpires" name:"PictureExpires"`

	// 开始时间，格式yyyy-MM-dd HH:mm:ss
	StartDateTime *string `json:"StartDateTime" name:"StartDateTime"`

	// 结束时间，格式yyyy-MM-dd HH:mm:ss
	EndDateTime *string `json:"EndDateTime" name:"EndDateTime"`
}

func (r *DescribePersonVisitInfoRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribePersonVisitInfoRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribePersonVisitInfoResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 公司ID
		CompanyId *string `json:"CompanyId" name:"CompanyId"`

		// 门店ID
		ShopId *uint64 `json:"ShopId" name:"ShopId"`

		// 总数
		TotalCount *uint64 `json:"TotalCount" name:"TotalCount"`

		// 用户到访明细
		PersonVisitInfoSet []*PersonVisitInfo `json:"PersonVisitInfoSet" name:"PersonVisitInfoSet" list`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribePersonVisitInfoResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribePersonVisitInfoResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeShopHourTrafficInfoRequest struct {
	*tchttp.BaseRequest

	// 公司ID
	CompanyId *string `json:"CompanyId" name:"CompanyId"`

	// 门店ID
	ShopId *uint64 `json:"ShopId" name:"ShopId"`

	// 开始日期，格式：yyyy-MM-dd
	StartDate *string `json:"StartDate" name:"StartDate"`

	// 结束日期，格式：yyyy-MM-dd
	EndDate *string `json:"EndDate" name:"EndDate"`

	// 偏移量：分页控制参数，第一页传0，第n页Offset=(n-1)*Limit
	Offset *uint64 `json:"Offset" name:"Offset"`

	// Limit:每页的数据项，最大100，超过100会被强制指定为100
	Limit *uint64 `json:"Limit" name:"Limit"`
}

func (r *DescribeShopHourTrafficInfoRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeShopHourTrafficInfoRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeShopHourTrafficInfoResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 公司ID
		CompanyId *string `json:"CompanyId" name:"CompanyId"`

		// 门店ID
		ShopId *uint64 `json:"ShopId" name:"ShopId"`

		// 查询结果总数
		TotalCount *uint64 `json:"TotalCount" name:"TotalCount"`

		// 分时客流信息
		ShopHourTrafficInfoSet []*ShopHourTrafficInfo `json:"ShopHourTrafficInfoSet" name:"ShopHourTrafficInfoSet" list`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeShopHourTrafficInfoResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeShopHourTrafficInfoResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeShopInfoRequest struct {
	*tchttp.BaseRequest

	// 偏移量：分页控制参数，第一页传0，第n页Offset=(n-1)*Limit
	Offset *uint64 `json:"Offset" name:"Offset"`

	// Limit:每页的数据项，最大100，超过100会被强制指定为100
	Limit *uint64 `json:"Limit" name:"Limit"`
}

func (r *DescribeShopInfoRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeShopInfoRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeShopInfoResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 门店总数
		TotalCount *uint64 `json:"TotalCount" name:"TotalCount"`

		// 门店列表信息
		ShopInfoSet []*ShopInfo `json:"ShopInfoSet" name:"ShopInfoSet" list`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeShopInfoResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeShopInfoResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeShopTrafficInfoRequest struct {
	*tchttp.BaseRequest

	// 公司ID
	CompanyId *string `json:"CompanyId" name:"CompanyId"`

	// 门店ID
	ShopId *uint64 `json:"ShopId" name:"ShopId"`

	// 开始日期，格式yyyy-MM-dd
	StartDate *string `json:"StartDate" name:"StartDate"`

	// 介绍日期，格式yyyy-MM-dd
	EndDate *string `json:"EndDate" name:"EndDate"`

	// 偏移量：分页控制参数，第一页传0，第n页Offset=(n-1)*Limit
	Offset *uint64 `json:"Offset" name:"Offset"`

	// Limit:每页的数据项，最大100，超过100会被强制指定为100
	Limit *uint64 `json:"Limit" name:"Limit"`
}

func (r *DescribeShopTrafficInfoRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeShopTrafficInfoRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeShopTrafficInfoResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 公司ID
		CompanyId *string `json:"CompanyId" name:"CompanyId"`

		// 门店ID
		ShopId *uint64 `json:"ShopId" name:"ShopId"`

		// 查询结果总数
		TotalCount *uint64 `json:"TotalCount" name:"TotalCount"`

		// 客流信息列表
		ShopDayTrafficInfoSet []*ShopDayTrafficInfo `json:"ShopDayTrafficInfoSet" name:"ShopDayTrafficInfoSet" list`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeShopTrafficInfoResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeShopTrafficInfoResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeTrajectoryDataRequest struct {
	*tchttp.BaseRequest

	// 集团ID
	CompanyId *string `json:"CompanyId" name:"CompanyId"`

	// 店铺ID
	ShopId *int64 `json:"ShopId" name:"ShopId"`

	// 开始日期，格式yyyy-MM-dd
	StartDate *string `json:"StartDate" name:"StartDate"`

	// 结束日期，格式yyyy-MM-dd
	EndDate *string `json:"EndDate" name:"EndDate"`

	// 限制返回数据的最大条数，最大 400（负数代为 400）
	Limit *int64 `json:"Limit" name:"Limit"`

	// 顾客性别顾虑，0是男，1是女，其它代表不分性别
	Gender *int64 `json:"Gender" name:"Gender"`
}

func (r *DescribeTrajectoryDataRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeTrajectoryDataRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeTrajectoryDataResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 集团ID
		CompanyId *string `json:"CompanyId" name:"CompanyId"`

		// 店铺ID
		ShopId *int64 `json:"ShopId" name:"ShopId"`

		// 总人数
		TotalPerson *int64 `json:"TotalPerson" name:"TotalPerson"`

		// 总动迹数目
		TotalTrajectory *int64 `json:"TotalTrajectory" name:"TotalTrajectory"`

		// 返回动迹中的总人数
		Person *int64 `json:"Person" name:"Person"`

		// 返回动迹的数目
		Trajectory *int64 `json:"Trajectory" name:"Trajectory"`

		// 返回动迹的具体信息
		Data []*TrajectorySunData `json:"Data" name:"Data" list`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeTrajectoryDataResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeTrajectoryDataResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeZoneFlowAgeInfoByZoneIdRequest struct {
	*tchttp.BaseRequest

	// 集团ID
	CompanyId *string `json:"CompanyId" name:"CompanyId"`

	// 店铺ID
	ShopId *int64 `json:"ShopId" name:"ShopId"`

	// 区域ID
	ZoneId *int64 `json:"ZoneId" name:"ZoneId"`

	// 开始日期，格式yyyy-MM-dd
	StartDate *string `json:"StartDate" name:"StartDate"`

	// 结束日期，格式yyyy-MM-dd
	EndDate *string `json:"EndDate" name:"EndDate"`
}

func (r *DescribeZoneFlowAgeInfoByZoneIdRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeZoneFlowAgeInfoByZoneIdRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeZoneFlowAgeInfoByZoneIdResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 集团ID
		CompanyId *string `json:"CompanyId" name:"CompanyId"`

		// 店铺ID
		ShopId *int64 `json:"ShopId" name:"ShopId"`

		// 区域ID
		ZoneId *int64 `json:"ZoneId" name:"ZoneId"`

		// 区域名称
		ZoneName *string `json:"ZoneName" name:"ZoneName"`

		// 当前年龄段占比
		Data []*float64 `json:"Data" name:"Data" list`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeZoneFlowAgeInfoByZoneIdResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeZoneFlowAgeInfoByZoneIdResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeZoneFlowAndStayTimeRequest struct {
	*tchttp.BaseRequest

	// 集团ID
	CompanyId *string `json:"CompanyId" name:"CompanyId"`

	// 店铺ID
	ShopId *int64 `json:"ShopId" name:"ShopId"`

	// 开始日期，格式yyyy-MM-dd
	StartDate *string `json:"StartDate" name:"StartDate"`

	// 结束日期，格式yyyy-MM-dd
	EndDate *string `json:"EndDate" name:"EndDate"`
}

func (r *DescribeZoneFlowAndStayTimeRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeZoneFlowAndStayTimeRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeZoneFlowAndStayTimeResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 集团id
		CompanyId *string `json:"CompanyId" name:"CompanyId"`

		// 店铺id
		ShopId *int64 `json:"ShopId" name:"ShopId"`

		// 各区域人流数目和停留时长
		Data []*ZoneFlowAndAvrStayTime `json:"Data" name:"Data" list`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeZoneFlowAndStayTimeResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeZoneFlowAndStayTimeResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeZoneFlowDailyByZoneIdRequest struct {
	*tchttp.BaseRequest

	// 集团ID
	CompanyId *string `json:"CompanyId" name:"CompanyId"`

	// 店铺ID
	ShopId *int64 `json:"ShopId" name:"ShopId"`

	// 区域ID
	ZoneId *int64 `json:"ZoneId" name:"ZoneId"`

	// 开始日期，格式yyyy-MM-dd
	StartDate *string `json:"StartDate" name:"StartDate"`

	// 结束日期，格式yyyy-MM-dd
	EndDate *string `json:"EndDate" name:"EndDate"`
}

func (r *DescribeZoneFlowDailyByZoneIdRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeZoneFlowDailyByZoneIdRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeZoneFlowDailyByZoneIdResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 集团id
		CompanyId *string `json:"CompanyId" name:"CompanyId"`

		// 店铺id
		ShopId *int64 `json:"ShopId" name:"ShopId"`

		// 区域ID
		ZoneId *int64 `json:"ZoneId" name:"ZoneId"`

		// 区域名称
		ZoneName *string `json:"ZoneName" name:"ZoneName"`

		// 每日人流量
		Data []*ZoneDayFlow `json:"Data" name:"Data" list`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeZoneFlowDailyByZoneIdResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeZoneFlowDailyByZoneIdResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeZoneFlowGenderAvrStayTimeByZoneIdRequest struct {
	*tchttp.BaseRequest

	// 集团ID
	CompanyId *string `json:"CompanyId" name:"CompanyId"`

	// 店铺ID
	ShopId *int64 `json:"ShopId" name:"ShopId"`

	// 区域ID
	ZoneId *int64 `json:"ZoneId" name:"ZoneId"`

	// 开始日期，格式yyyy-MM-dd
	StartDate *string `json:"StartDate" name:"StartDate"`

	// 结束日期，格式yyyy-MM-dd
	EndDate *string `json:"EndDate" name:"EndDate"`
}

func (r *DescribeZoneFlowGenderAvrStayTimeByZoneIdRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeZoneFlowGenderAvrStayTimeByZoneIdRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeZoneFlowGenderAvrStayTimeByZoneIdResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 集团ID
		CompanyId *string `json:"CompanyId" name:"CompanyId"`

		// 店铺ID
		ShopId *int64 `json:"ShopId" name:"ShopId"`

		// 区域ID
		ZoneId *int64 `json:"ZoneId" name:"ZoneId"`

		// 区域名称
		ZoneName *string `json:"ZoneName" name:"ZoneName"`

		// 不同年龄段男女停留时间（返回格式为数组，从第 1 个到最后一个数据，年龄段分别为 0-17，18 - 23,  24 - 30, 31 - 40, 41 - 50, 51 - 60, 61 - 100）
		Data []*ZoneAgeGroupAvrStayTime `json:"Data" name:"Data" list`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeZoneFlowGenderAvrStayTimeByZoneIdResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeZoneFlowGenderAvrStayTimeByZoneIdResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeZoneFlowGenderInfoByZoneIdRequest struct {
	*tchttp.BaseRequest

	// 集团ID
	CompanyId *string `json:"CompanyId" name:"CompanyId"`

	// 店铺ID
	ShopId *int64 `json:"ShopId" name:"ShopId"`

	// 区域ID
	ZoneId *int64 `json:"ZoneId" name:"ZoneId"`

	// 开始日期，格式yyyy-MM-dd
	StartDate *string `json:"StartDate" name:"StartDate"`

	// 结束日期，格式yyyy-MM-dd
	EndDate *string `json:"EndDate" name:"EndDate"`
}

func (r *DescribeZoneFlowGenderInfoByZoneIdRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeZoneFlowGenderInfoByZoneIdRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeZoneFlowGenderInfoByZoneIdResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 集团ID
		CompanyId *string `json:"CompanyId" name:"CompanyId"`

		// 店铺ID
		ShopId *int64 `json:"ShopId" name:"ShopId"`

		// 区域ID
		ZoneId *int64 `json:"ZoneId" name:"ZoneId"`

		// 区域名称
		ZoneName *string `json:"ZoneName" name:"ZoneName"`

		// 男性占比
		MalePercent *float64 `json:"MalePercent" name:"MalePercent"`

		// 女性占比
		FemalePercent *float64 `json:"FemalePercent" name:"FemalePercent"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeZoneFlowGenderInfoByZoneIdResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeZoneFlowGenderInfoByZoneIdResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeZoneFlowHourlyByZoneIdRequest struct {
	*tchttp.BaseRequest

	// 集团ID
	CompanyId *string `json:"CompanyId" name:"CompanyId"`

	// 店铺ID
	ShopId *int64 `json:"ShopId" name:"ShopId"`

	// 区域ID
	ZoneId *int64 `json:"ZoneId" name:"ZoneId"`

	// 开始日期，格式yyyy-MM-dd
	StartDate *string `json:"StartDate" name:"StartDate"`

	// 结束日期，格式yyyy-MM-dd
	EndDate *string `json:"EndDate" name:"EndDate"`
}

func (r *DescribeZoneFlowHourlyByZoneIdRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeZoneFlowHourlyByZoneIdRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeZoneFlowHourlyByZoneIdResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 集团ID
		CompanyId *string `json:"CompanyId" name:"CompanyId"`

		// 店铺ID
		ShopId *int64 `json:"ShopId" name:"ShopId"`

		// 区域ID
		ZoneId *int64 `json:"ZoneId" name:"ZoneId"`

		// 区域名称
		ZoneName *string `json:"ZoneName" name:"ZoneName"`

		// 各个分时人流量
		Data []*ZoneHourFlow `json:"Data" name:"Data" list`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeZoneFlowHourlyByZoneIdResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeZoneFlowHourlyByZoneIdResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeZoneTrafficInfoRequest struct {
	*tchttp.BaseRequest

	// 公司ID
	CompanyId *string `json:"CompanyId" name:"CompanyId"`

	// 店铺ID
	ShopId *uint64 `json:"ShopId" name:"ShopId"`

	// 开始日期，格式yyyy-MM-dd
	StartDate *string `json:"StartDate" name:"StartDate"`

	// 结束日期，格式yyyy-MM-dd
	EndDate *string `json:"EndDate" name:"EndDate"`

	// 偏移量：分页控制参数，第一页传0，第n页Offset=(n-1)*Limit
	Offset *uint64 `json:"Offset" name:"Offset"`

	// Limit:每页的数据项，最大100，超过100会被强制指定为100
	Limit *uint64 `json:"Limit" name:"Limit"`
}

func (r *DescribeZoneTrafficInfoRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeZoneTrafficInfoRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeZoneTrafficInfoResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 公司ID
		CompanyId *string `json:"CompanyId" name:"CompanyId"`

		// 门店ID
		ShopId *uint64 `json:"ShopId" name:"ShopId"`

		// 查询结果总数
		TotalCount *uint64 `json:"TotalCount" name:"TotalCount"`

		// 区域客流信息列表
		ZoneTrafficInfoSet []*ZoneTrafficInfo `json:"ZoneTrafficInfoSet" name:"ZoneTrafficInfoSet" list`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeZoneTrafficInfoResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeZoneTrafficInfoResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type GenderAgeTrafficDetail struct {

	// 性别: 0男1女
	Gender *uint64 `json:"Gender" name:"Gender"`

	// 年龄区间，枚举值：0-17、18-23、24-30、31-40、41-50、51-60、>60
	AgeGap *string `json:"AgeGap" name:"AgeGap"`

	// 客流量
	TrafficCount *uint64 `json:"TrafficCount" name:"TrafficCount"`
}

type HourTrafficInfoDetail struct {

	// 小时 取值为：0，1，2，3，4，5，6，7，8，9，10，11，12，13，14，15，16，17，18，19，20，21，22，23
	Hour *uint64 `json:"Hour" name:"Hour"`

	// 分时客流量
	HourTrafficTotalCount *uint64 `json:"HourTrafficTotalCount" name:"HourTrafficTotalCount"`
}

type ModifyPersonTagInfoRequest struct {
	*tchttp.BaseRequest

	// 优mall集团id，通过"指定身份标识获取客户门店列表"接口获取
	CompanyId *string `json:"CompanyId" name:"CompanyId"`

	// 优mall店铺id，通过"指定身份标识获取客户门店列表"接口获取，为0则拉取集团全部店铺当前
	ShopId *int64 `json:"ShopId" name:"ShopId"`

	// 需要设置的顾客信息，批量设置最大为10个
	Tags []*PersonTagInfo `json:"Tags" name:"Tags" list`
}

func (r *ModifyPersonTagInfoRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ModifyPersonTagInfoRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ModifyPersonTagInfoResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *ModifyPersonTagInfoResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ModifyPersonTagInfoResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ModifyPersonTypeRequest struct {
	*tchttp.BaseRequest

	// 集团ID
	CompanyId *string `json:"CompanyId" name:"CompanyId"`

	// 门店ID
	ShopId *uint64 `json:"ShopId" name:"ShopId"`

	// 顾客ID
	PersonId *uint64 `json:"PersonId" name:"PersonId"`

	// 身份类型(0表示普通顾客，1 白名单，2 表示黑名单）
	PersonType *uint64 `json:"PersonType" name:"PersonType"`

	// 身份子类型:
	// PersonType=0时(普通顾客)，0普通顾客
	// PersonType=1时(白名单)，0店员，1商场人员，2其他类型人员，3区域经理，4注册用户，5VIP用户
	// PersonType=2时(黑名单)，0普通黑名单，1小偷)
	PersonSubType *uint64 `json:"PersonSubType" name:"PersonSubType"`
}

func (r *ModifyPersonTypeRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ModifyPersonTypeRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ModifyPersonTypeResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *ModifyPersonTypeResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ModifyPersonTypeResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type NetworkAndShopInfo struct {

	// 集团id
	CompanyId *string `json:"CompanyId" name:"CompanyId"`

	// 店铺id
	ShopId *int64 `json:"ShopId" name:"ShopId"`

	// 店铺省份
	Province *string `json:"Province" name:"Province"`

	// 店铺城市
	City *string `json:"City" name:"City"`

	// 店铺名
	ShopName *string `json:"ShopName" name:"ShopName"`

	// 上传带宽，单位Mb/s，-1：未知
	Upload *float64 `json:"Upload" name:"Upload"`

	// 下载带宽，单位Mb/s，-1：未知
	Download *float64 `json:"Download" name:"Download"`

	// 最小延迟，单位ms，-1：未知
	MinRtt *float64 `json:"MinRtt" name:"MinRtt"`

	// 平均延迟，单位ms，-1：未知
	AvgRtt *float64 `json:"AvgRtt" name:"AvgRtt"`

	// 最大延迟，单位ms，-1：未知
	MaxRtt *float64 `json:"MaxRtt" name:"MaxRtt"`

	// 平均偏差延迟，单位ms，-1：未知
	MdevRtt *float64 `json:"MdevRtt" name:"MdevRtt"`

	// 丢包率百分比，-1：未知
	Loss *float64 `json:"Loss" name:"Loss"`

	// 更新时间戳
	UpdateTime *int64 `json:"UpdateTime" name:"UpdateTime"`

	// 上报网络状态设备
	Mac *string `json:"Mac" name:"Mac"`
}

type NetworkHistoryInfo struct {

	// 总数
	Count *int64 `json:"Count" name:"Count"`

	// 集团id
	CompanyId *string `json:"CompanyId" name:"CompanyId"`

	// 店铺id
	ShopId *int64 `json:"ShopId" name:"ShopId"`

	// 店铺省份
	Province *string `json:"Province" name:"Province"`

	// 店铺城市
	City *string `json:"City" name:"City"`

	// 店铺名称
	ShopName *string `json:"ShopName" name:"ShopName"`

	// 网络信息
	Infos []*NetworkInfo `json:"Infos" name:"Infos" list`
}

type NetworkInfo struct {

	// 上传带宽，单位Mb/s，-1：未知
	Upload *float64 `json:"Upload" name:"Upload"`

	// 下载带宽，单位Mb/s，-1：未知
	Download *float64 `json:"Download" name:"Download"`

	// 最小延迟，单位ms，-1：未知
	MinRtt *float64 `json:"MinRtt" name:"MinRtt"`

	// 平均延迟，单位ms，-1：未知
	AvgRtt *float64 `json:"AvgRtt" name:"AvgRtt"`

	// 最大延迟，单位ms，-1：未知
	MaxRtt *float64 `json:"MaxRtt" name:"MaxRtt"`

	// 平均偏差延迟，单位ms，-1：未知
	MdevRtt *float64 `json:"MdevRtt" name:"MdevRtt"`

	// 丢包率百分比，-1：未知
	Loss *float64 `json:"Loss" name:"Loss"`

	// 更新时间戳
	UpdateTime *int64 `json:"UpdateTime" name:"UpdateTime"`

	// 上报网络状态设备
	Mac *string `json:"Mac" name:"Mac"`
}

type NetworkLastInfo struct {

	// 总数
	Count *int64 `json:"Count" name:"Count"`

	// 网络状态
	Infos []*NetworkAndShopInfo `json:"Infos" name:"Infos" list`
}

type PersonCoordinate struct {

	// CAD图X坐标
	CADX *float64 `json:"CADX" name:"CADX"`

	// CAD图Y坐标
	CADY *float64 `json:"CADY" name:"CADY"`

	// 抓拍时间点
	CapTime *string `json:"CapTime" name:"CapTime"`

	// 抓拍图片
	CapPic *string `json:"CapPic" name:"CapPic"`

	// 卖场区域类型
	MallAreaType *int64 `json:"MallAreaType" name:"MallAreaType"`

	// 坐标编号
	PosId *int64 `json:"PosId" name:"PosId"`

	// 门店编号
	ShopId *int64 `json:"ShopId" name:"ShopId"`

	// 事件
	Event *string `json:"Event" name:"Event"`
}

type PersonInfo struct {

	// 用户ID
	PersonId *uint64 `json:"PersonId" name:"PersonId"`

	// 人脸图片Base64内容，已弃用，返回默认空值
	PersonPicture *string `json:"PersonPicture" name:"PersonPicture"`

	// 性别：0男1女
	Gender *int64 `json:"Gender" name:"Gender"`

	// 年龄
	Age *int64 `json:"Age" name:"Age"`

	// 身份类型（0表示普通顾客，1 白名单，2 表示黑名单）
	PersonType *int64 `json:"PersonType" name:"PersonType"`

	// 人脸图片Url，在有效期内可以访问下载
	PersonPictureUrl *string `json:"PersonPictureUrl" name:"PersonPictureUrl"`

	// 身份子类型:
	// PersonType=0时(普通顾客)，0普通顾客
	// PersonType=1时(白名单)，0店员，1商场人员，2其他类型人员，3区域经理，4注册用户，5VIP用户
	// PersonType=2时(黑名单)，0普通黑名单，1小偷)
	PersonSubType *int64 `json:"PersonSubType" name:"PersonSubType"`
}

type PersonProfile struct {

	// 客人编码
	PersonId *string `json:"PersonId" name:"PersonId"`

	// 性别
	Gender *uint64 `json:"Gender" name:"Gender"`

	// 年龄
	Age *uint64 `json:"Age" name:"Age"`

	// 首次到场时间
	FirstArrivedTime *string `json:"FirstArrivedTime" name:"FirstArrivedTime"`

	// 来访次数
	ArrivedCount *uint64 `json:"ArrivedCount" name:"ArrivedCount"`

	// 客户图片
	PicUrl *string `json:"PicUrl" name:"PicUrl"`

	// 置信度
	Similarity *float64 `json:"Similarity" name:"Similarity"`
}

type PersonTagInfo struct {

	// 顾客原类型
	OldType *int64 `json:"OldType" name:"OldType"`

	// 顾客新类型
	NewType *int64 `json:"NewType" name:"NewType"`

	// 顾客face id
	PersonId *int64 `json:"PersonId" name:"PersonId"`
}

type PersonTracePoint struct {

	// 卖场区域编码
	MallAreaId *uint64 `json:"MallAreaId" name:"MallAreaId"`

	// 门店编码
	ShopId *uint64 `json:"ShopId" name:"ShopId"`

	// 卖场区域类型
	MallAreaType *uint64 `json:"MallAreaType" name:"MallAreaType"`

	// 轨迹事件
	TraceEventType *uint64 `json:"TraceEventType" name:"TraceEventType"`

	// 轨迹事件发生时间点
	TraceEventTime *string `json:"TraceEventTime" name:"TraceEventTime"`

	// 抓拍图片
	CapPic *string `json:"CapPic" name:"CapPic"`
}

type PersonTraceRoute struct {

	// 轨迹编码
	TraceId *string `json:"TraceId" name:"TraceId"`

	// 轨迹点序列
	TracePointSet []*PersonTracePoint `json:"TracePointSet" name:"TracePointSet" list`
}

type PersonVisitInfo struct {

	// 用户ID
	PersonId *uint64 `json:"PersonId" name:"PersonId"`

	// 用户到访ID
	VisitId *uint64 `json:"VisitId" name:"VisitId"`

	// 到访时间：Unix时间戳
	InTime *uint64 `json:"InTime" name:"InTime"`

	// 抓拍到的头像Base64内容，已弃用，返回默认空值
	CapturedPicture *string `json:"CapturedPicture" name:"CapturedPicture"`

	// 口罩类型：0不戴口罩，1戴口罩
	MaskType *uint64 `json:"MaskType" name:"MaskType"`

	// 眼镜类型：0不戴眼镜，1普通眼镜 , 2墨镜
	GlassType *uint64 `json:"GlassType" name:"GlassType"`

	// 发型：0 短发,  1长发
	HairType *uint64 `json:"HairType" name:"HairType"`

	// 抓拍到的头像Url，在有效期内可以访问下载
	CapturedPictureUrl *string `json:"CapturedPictureUrl" name:"CapturedPictureUrl"`

	// 抓拍头像的场景图信息
	SceneInfo *SceneInfo `json:"SceneInfo" name:"SceneInfo"`
}

type RegisterCallbackRequest struct {
	*tchttp.BaseRequest

	// 集团id，通过"指定身份标识获取客户门店列表"接口获取
	CompanyId *string `json:"CompanyId" name:"CompanyId"`

	// 通知回调地址，完整url，示例（http://youmall.tencentcloudapi.com/）
	BackUrl *string `json:"BackUrl" name:"BackUrl"`

	// 请求时间戳
	Time *uint64 `json:"Time" name:"Time"`

	// 是否需要顾客图片，1-需要图片，其它-不需要图片
	NeedFacePic *uint64 `json:"NeedFacePic" name:"NeedFacePic"`
}

func (r *RegisterCallbackRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *RegisterCallbackRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type RegisterCallbackResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *RegisterCallbackResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *RegisterCallbackResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type SceneInfo struct {

	// 场景图
	ScenePictureURL *string `json:"ScenePictureURL" name:"ScenePictureURL"`

	// 抓拍头像左上角X坐标在场景图中的像素点位置
	HeadX *int64 `json:"HeadX" name:"HeadX"`

	// 抓拍头像左上角Y坐标在场景图中的像素点位置
	HeadY *int64 `json:"HeadY" name:"HeadY"`

	// 抓拍头像在场景图中占有的像素宽度
	HeadWidth *int64 `json:"HeadWidth" name:"HeadWidth"`

	// 抓拍头像在场景图中占有的像素高度
	HeadHeight *int64 `json:"HeadHeight" name:"HeadHeight"`
}

type ShopDayTrafficInfo struct {

	// 日期
	Date *string `json:"Date" name:"Date"`

	// 客流量
	DayTrafficTotalCount *uint64 `json:"DayTrafficTotalCount" name:"DayTrafficTotalCount"`

	// 性别年龄分组下的客流信息
	GenderAgeTrafficDetailSet []*GenderAgeTrafficDetail `json:"GenderAgeTrafficDetailSet" name:"GenderAgeTrafficDetailSet" list`
}

type ShopHourTrafficInfo struct {

	// 日期，格式yyyy-MM-dd
	Date *string `json:"Date" name:"Date"`

	// 分时客流详细信息
	HourTrafficInfoDetailSet []*HourTrafficInfoDetail `json:"HourTrafficInfoDetailSet" name:"HourTrafficInfoDetailSet" list`
}

type ShopInfo struct {

	// 公司ID
	CompanyId *string `json:"CompanyId" name:"CompanyId"`

	// 门店ID
	ShopId *uint64 `json:"ShopId" name:"ShopId"`

	// 门店名称
	ShopName *string `json:"ShopName" name:"ShopName"`

	// 客户门店编码
	ShopCode *string `json:"ShopCode" name:"ShopCode"`

	// 省
	Province *string `json:"Province" name:"Province"`

	// 市
	City *string `json:"City" name:"City"`

	// 公司名称
	CompanyName *string `json:"CompanyName" name:"CompanyName"`
}

type TrajectorySunData struct {

	// 区域动线，形如 x-x-x-x-x，其中 x 为区域 ID
	Zones *string `json:"Zones" name:"Zones"`

	// 该动线出现次数
	Count *int64 `json:"Count" name:"Count"`

	// 该动线平均停留时间（秒）
	AvgStayTime *int64 `json:"AvgStayTime" name:"AvgStayTime"`
}

type ZoneAgeGroupAvrStayTime struct {

	// 男性平均停留时间
	MaleAvrStayTime *float64 `json:"MaleAvrStayTime" name:"MaleAvrStayTime"`

	// 女性平均停留时间
	FemaleAvrStayTime *float64 `json:"FemaleAvrStayTime" name:"FemaleAvrStayTime"`
}

type ZoneDayFlow struct {

	// 日期，如 2018-08-6
	Day *string `json:"Day" name:"Day"`

	// 客流量
	FlowCount *int64 `json:"FlowCount" name:"FlowCount"`
}

type ZoneFlowAndAvrStayTime struct {

	// 区域id
	ZoneId *int64 `json:"ZoneId" name:"ZoneId"`

	// 区域名称
	ZoneName *string `json:"ZoneName" name:"ZoneName"`

	// 人流量
	FlowCount *uint64 `json:"FlowCount" name:"FlowCount"`

	// 平均停留时长
	AvrStayTime *uint64 `json:"AvrStayTime" name:"AvrStayTime"`
}

type ZoneHourFlow struct {

	// 分时 0~23
	Hour *int64 `json:"Hour" name:"Hour"`

	// 客流量
	FlowCount *int64 `json:"FlowCount" name:"FlowCount"`
}

type ZoneTrafficInfo struct {

	// 日期
	Date *string `json:"Date" name:"Date"`

	// 门店区域客流详细信息
	ZoneTrafficInfoDetailSet []*ZoneTrafficInfoDetail `json:"ZoneTrafficInfoDetailSet" name:"ZoneTrafficInfoDetailSet" list`
}

type ZoneTrafficInfoDetail struct {

	// 区域ID
	ZoneId *uint64 `json:"ZoneId" name:"ZoneId"`

	// 区域名称
	ZoneName *string `json:"ZoneName" name:"ZoneName"`

	// 客流量
	TrafficTotalCount *uint64 `json:"TrafficTotalCount" name:"TrafficTotalCount"`

	// 平均停留时间
	AvgStayTime *uint64 `json:"AvgStayTime" name:"AvgStayTime"`
}
