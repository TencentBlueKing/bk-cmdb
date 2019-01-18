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

package v20180801

import (
    "encoding/json"

    tchttp "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/http"
)

type AddDelayLiveStreamRequest struct {
	*tchttp.BaseRequest

	// 应用名称。
	AppName *string `json:"AppName" name:"AppName"`

	// 您的加速域名。
	DomainName *string `json:"DomainName" name:"DomainName"`

	// 流名称。
	StreamName *string `json:"StreamName" name:"StreamName"`

	// 延播时间，单位：秒，上限：600秒。
	DelayTime *uint64 `json:"DelayTime" name:"DelayTime"`
}

func (r *AddDelayLiveStreamRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *AddDelayLiveStreamRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type AddDelayLiveStreamResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *AddDelayLiveStreamResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *AddDelayLiveStreamResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type AddLiveWatermarkRequest struct {
	*tchttp.BaseRequest

	// 水印图片url。
	PictureUrl *string `json:"PictureUrl" name:"PictureUrl"`

	// 水印名称。
	WatermarkName *string `json:"WatermarkName" name:"WatermarkName"`

	// 显示位置,X轴偏移。
	XPosition *int64 `json:"XPosition" name:"XPosition"`

	// 显示位置,Y轴偏移。
	YPosition *int64 `json:"YPosition" name:"YPosition"`
}

func (r *AddLiveWatermarkRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *AddLiveWatermarkRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type AddLiveWatermarkResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 水印ID。
		WatermarkId *uint64 `json:"WatermarkId" name:"WatermarkId"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *AddLiveWatermarkResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *AddLiveWatermarkResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type CreateLiveRecordRequest struct {
	*tchttp.BaseRequest

	// 流名称。
	StreamName *string `json:"StreamName" name:"StreamName"`

	// 直播流所属应用名称。
	AppName *string `json:"AppName" name:"AppName"`

	// 推流域名。多域名推流必须设置。
	DomainName *string `json:"DomainName" name:"DomainName"`

	// 任务起始时间，中国标准时间，需要URLEncode。如 2017-01-01 10:10:01，编码为：2017-01-01+10%3a10%3a01。录制视频为精彩视频时，忽略此字段。
	StartTime *string `json:"StartTime" name:"StartTime"`

	// 任务结束时间，中国标准时间，需要URLEncode。如 2017-01-01 10:30:01，编码为：2017-01-01+10%3a30%3a01。若指定精彩视频录制，结束时间不超过当前时间+30分钟，如果超过或小于起始时间，则实际结束时间为当前时间+30分钟。
	EndTime *string `json:"EndTime" name:"EndTime"`

	// 录制类型。不区分大小写。
	// “video” : 音视频录制【默认】。
	// “audio” : 纯音频录制。
	RecordType *string `json:"RecordType" name:"RecordType"`

	// 录制文件格式。不区分大小写。其值为：
	// “flv”,“hls”,”mp4”,“aac”,”mp3”，默认“flv”。
	FileFormat *string `json:"FileFormat" name:"FileFormat"`

	// 精彩视频标志。0：普通视频【默认】；1：精彩视频。
	Highlight *int64 `json:"Highlight" name:"Highlight"`

	// A+B=C混流标志。0：非A+B=C混流录制【默认】；1：标示为A+B=C混流录制。
	MixStream *int64 `json:"MixStream" name:"MixStream"`

	// 录制流参数，当前支持以下参数： 
	// interval 录制分片时长，单位 秒，0 - 7200
	// storage_time 录制文件存储时长，单位 秒
	// eg. interval=3600&storage_time=7200
	// 注：参数需要url encode。
	StreamParam *string `json:"StreamParam" name:"StreamParam"`
}

func (r *CreateLiveRecordRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *CreateLiveRecordRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type CreateLiveRecordResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 任务ID，全局唯一标识录制任务。
		TaskId *uint64 `json:"TaskId" name:"TaskId"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *CreateLiveRecordResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *CreateLiveRecordResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type CreatePullStreamConfigRequest struct {
	*tchttp.BaseRequest

	// 源Url。
	FromUrl *string `json:"FromUrl" name:"FromUrl"`

	// 目的Url，目前限制该目标地址为腾讯域名。
	ToUrl *string `json:"ToUrl" name:"ToUrl"`

	// 区域id,1-深圳,2-上海，3-天津,4-香港。
	AreaId *int64 `json:"AreaId" name:"AreaId"`

	// 运营商id,1-电信,2-移动,3-联通,4-其他,AreaId为4的时候,IspId只能为其他。
	IspId *int64 `json:"IspId" name:"IspId"`

	// 开始时间。
	StartTime *string `json:"StartTime" name:"StartTime"`

	// 结束时间，注意：
	// 1. 结束时间必须大于开始时间；
	// 2. 结束时间和开始时间必须大于当前时间；
	// 3. 结束时间 和 开始时间 间隔必须小于七天。
	EndTime *string `json:"EndTime" name:"EndTime"`
}

func (r *CreatePullStreamConfigRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *CreatePullStreamConfigRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type CreatePullStreamConfigResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 配置成功后的id。
		ConfigId *string `json:"ConfigId" name:"ConfigId"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *CreatePullStreamConfigResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *CreatePullStreamConfigResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DeleteLiveRecordRequest struct {
	*tchttp.BaseRequest

	// 流名称。
	StreamName *string `json:"StreamName" name:"StreamName"`

	// 任务ID，全局唯一标识录制任务。
	TaskId *int64 `json:"TaskId" name:"TaskId"`
}

func (r *DeleteLiveRecordRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DeleteLiveRecordRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DeleteLiveRecordResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DeleteLiveRecordResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DeleteLiveRecordResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DeleteLiveWatermarkRequest struct {
	*tchttp.BaseRequest

	// 水印ID。
	WatermarkId *int64 `json:"WatermarkId" name:"WatermarkId"`
}

func (r *DeleteLiveWatermarkRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DeleteLiveWatermarkRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DeleteLiveWatermarkResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DeleteLiveWatermarkResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DeleteLiveWatermarkResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DeletePullStreamConfigRequest struct {
	*tchttp.BaseRequest

	// 配置id。
	ConfigId *string `json:"ConfigId" name:"ConfigId"`
}

func (r *DeletePullStreamConfigRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DeletePullStreamConfigRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DeletePullStreamConfigResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DeletePullStreamConfigResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DeletePullStreamConfigResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeLivePlayAuthKeyRequest struct {
	*tchttp.BaseRequest

	// 域名。
	DomainName *string `json:"DomainName" name:"DomainName"`
}

func (r *DescribeLivePlayAuthKeyRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeLivePlayAuthKeyRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeLivePlayAuthKeyResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 播放鉴权key信息。
		PlayAuthKeyInfo *PlayAuthKeyInfo `json:"PlayAuthKeyInfo" name:"PlayAuthKeyInfo"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeLivePlayAuthKeyResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeLivePlayAuthKeyResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeLivePushAuthKeyRequest struct {
	*tchttp.BaseRequest

	// 推流域名。
	DomainName *string `json:"DomainName" name:"DomainName"`
}

func (r *DescribeLivePushAuthKeyRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeLivePushAuthKeyRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeLivePushAuthKeyResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 推流鉴权key信息。
		PushAuthKeyInfo *PushAuthKeyInfo `json:"PushAuthKeyInfo" name:"PushAuthKeyInfo"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeLivePushAuthKeyResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeLivePushAuthKeyResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeLiveStreamOnlineInfoRequest struct {
	*tchttp.BaseRequest

	// 取得第几页。
	// 默认值：1
	PageNum *uint64 `json:"PageNum" name:"PageNum"`

	// 分页大小。
	// 
	// 最大值：100。
	// 取值范围：1~100 之前的任意整数。
	// 默认值：10
	PageSize *uint64 `json:"PageSize" name:"PageSize"`

	// 0:未开始推流 1:正在推流 2:服务出错 3:已关闭。
	Status *int64 `json:"Status" name:"Status"`

	// 流名称。
	StreamName *string `json:"StreamName" name:"StreamName"`
}

func (r *DescribeLiveStreamOnlineInfoRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeLiveStreamOnlineInfoRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeLiveStreamOnlineInfoResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 分页的页码。
		PageNum *uint64 `json:"PageNum" name:"PageNum"`

		// 每页大小
		PageSize *uint64 `json:"PageSize" name:"PageSize"`

		// 符合条件的总个数。
		TotalNum *uint64 `json:"TotalNum" name:"TotalNum"`

		// 总页数。
		TotalPage *uint64 `json:"TotalPage" name:"TotalPage"`

		// 流信息列表
		StreamInfoList []*StreamInfo `json:"StreamInfoList" name:"StreamInfoList" list`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeLiveStreamOnlineInfoResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeLiveStreamOnlineInfoResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeLiveStreamOnlineListRequest struct {
	*tchttp.BaseRequest

	// 您的加速域名。
	DomainName *string `json:"DomainName" name:"DomainName"`

	// 应用名称。
	AppName *string `json:"AppName" name:"AppName"`

	// 取得第几页，默认1。
	PageNum *uint64 `json:"PageNum" name:"PageNum"`

	// 每页大小，最大100。 
	// 取值：1~100之前的任意整数。
	// 默认值：10
	PageSize *uint64 `json:"PageSize" name:"PageSize"`
}

func (r *DescribeLiveStreamOnlineListRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeLiveStreamOnlineListRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeLiveStreamOnlineListResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 符合条件的总个数。
		TotalNum *uint64 `json:"TotalNum" name:"TotalNum"`

		// 总页数。
		TotalPage *uint64 `json:"TotalPage" name:"TotalPage"`

		// 分页的页码。
		PageNum *uint64 `json:"PageNum" name:"PageNum"`

		// 每页显示的条数。
		PageSize *uint64 `json:"PageSize" name:"PageSize"`

		// 正在推送流的信息列表
		OnlineInfo []*StreamOnlineInfo `json:"OnlineInfo" name:"OnlineInfo" list`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeLiveStreamOnlineListResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeLiveStreamOnlineListResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeLiveStreamPublishedListRequest struct {
	*tchttp.BaseRequest

	// 您的域名。
	DomainName *string `json:"DomainName" name:"DomainName"`

	// 结束时间。
	// UTC 格式，例如：2016-06-30T19:00:00Z。
	// 不超过当前时间。
	EndTime *string `json:"EndTime" name:"EndTime"`

	// 起始时间。 
	// UTC 格式，例如：2016-06-29T19:00:00Z。
	// 和当前时间相隔不超过7天。
	StartTime *string `json:"StartTime" name:"StartTime"`

	// 直播流所属应用名称。
	AppName *string `json:"AppName" name:"AppName"`

	// 取得第几页。
	// 默认值：1
	PageNum *uint64 `json:"PageNum" name:"PageNum"`

	// 分页大小。
	// 
	// 最大值：100。
	// 取值范围：1~100 之前的任意整数。
	// 默认值：10
	PageSize *uint64 `json:"PageSize" name:"PageSize"`
}

func (r *DescribeLiveStreamPublishedListRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeLiveStreamPublishedListRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeLiveStreamPublishedListResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 推流记录信息。
		PublishInfo []*StreamName `json:"PublishInfo" name:"PublishInfo" list`

		// 分页的页码。
		PageNum *uint64 `json:"PageNum" name:"PageNum"`

		// 每页大小
		PageSize *uint64 `json:"PageSize" name:"PageSize"`

		// 符合条件的总个数。
		TotalNum *uint64 `json:"TotalNum" name:"TotalNum"`

		// 总页数。
		TotalPage *uint64 `json:"TotalPage" name:"TotalPage"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeLiveStreamPublishedListResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeLiveStreamPublishedListResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeLiveStreamStateRequest struct {
	*tchttp.BaseRequest

	// 应用名称。
	AppName *string `json:"AppName" name:"AppName"`

	// 您的加速域名。
	DomainName *string `json:"DomainName" name:"DomainName"`

	// 流名称。
	StreamName *string `json:"StreamName" name:"StreamName"`
}

func (r *DescribeLiveStreamStateRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeLiveStreamStateRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeLiveStreamStateResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 流状态
		StreamState *string `json:"StreamState" name:"StreamState"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeLiveStreamStateResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeLiveStreamStateResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeLiveWatermarksRequest struct {
	*tchttp.BaseRequest
}

func (r *DescribeLiveWatermarksRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeLiveWatermarksRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeLiveWatermarksResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 水印总个数。
		TotalNum *uint64 `json:"TotalNum" name:"TotalNum"`

		// 水印信息列表。
		WatermarkList []*WatermarkInfo `json:"WatermarkList" name:"WatermarkList" list`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeLiveWatermarksResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeLiveWatermarksResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribePullStreamConfigsRequest struct {
	*tchttp.BaseRequest

	// 配置id。
	ConfigId *string `json:"ConfigId" name:"ConfigId"`
}

func (r *DescribePullStreamConfigsRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribePullStreamConfigsRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribePullStreamConfigsResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 拉流配置。
		PullStreamConfigs []*PullStreamConfig `json:"PullStreamConfigs" name:"PullStreamConfigs" list`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribePullStreamConfigsResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribePullStreamConfigsResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DropLiveStreamRequest struct {
	*tchttp.BaseRequest

	// 流名称。
	StreamName *string `json:"StreamName" name:"StreamName"`

	// 您的加速域名。
	DomainName *string `json:"DomainName" name:"DomainName"`

	// 应用名称。
	AppName *string `json:"AppName" name:"AppName"`
}

func (r *DropLiveStreamRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DropLiveStreamRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DropLiveStreamResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DropLiveStreamResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DropLiveStreamResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ForbidLiveStreamRequest struct {
	*tchttp.BaseRequest

	// 应用名称。
	AppName *string `json:"AppName" name:"AppName"`

	// 您的加速域名。
	DomainName *string `json:"DomainName" name:"DomainName"`

	// 流名称。
	StreamName *string `json:"StreamName" name:"StreamName"`

	// 恢复流的时间。UTC 格式，例如：2018-11-29T19:00:00Z。
	// 注意：默认禁播90天，且最长支持禁播90天。
	ResumeTime *string `json:"ResumeTime" name:"ResumeTime"`
}

func (r *ForbidLiveStreamRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ForbidLiveStreamRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ForbidLiveStreamResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *ForbidLiveStreamResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ForbidLiveStreamResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ModifyLivePlayAuthKeyRequest struct {
	*tchttp.BaseRequest

	// 域名。
	DomainName *string `json:"DomainName" name:"DomainName"`

	// 是否启用，0：关闭，1：启用。
	Enable *int64 `json:"Enable" name:"Enable"`

	// 鉴权key。
	AuthKey *string `json:"AuthKey" name:"AuthKey"`

	// 有效时间，单位：秒。
	AuthDelta *uint64 `json:"AuthDelta" name:"AuthDelta"`

	// 鉴权backkey。
	AuthBackKey *string `json:"AuthBackKey" name:"AuthBackKey"`
}

func (r *ModifyLivePlayAuthKeyRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ModifyLivePlayAuthKeyRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ModifyLivePlayAuthKeyResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *ModifyLivePlayAuthKeyResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ModifyLivePlayAuthKeyResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ModifyLivePushAuthKeyRequest struct {
	*tchttp.BaseRequest

	// 推流域名。
	DomainName *string `json:"DomainName" name:"DomainName"`

	// 是否启用，0：关闭，1：启用。
	Enable *int64 `json:"Enable" name:"Enable"`

	// 主鉴权key。
	MasterAuthKey *string `json:"MasterAuthKey" name:"MasterAuthKey"`

	// 备鉴权key。
	BackupAuthKey *string `json:"BackupAuthKey" name:"BackupAuthKey"`

	// 有效时间，单位：秒。
	AuthDelta *uint64 `json:"AuthDelta" name:"AuthDelta"`
}

func (r *ModifyLivePushAuthKeyRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ModifyLivePushAuthKeyRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ModifyLivePushAuthKeyResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *ModifyLivePushAuthKeyResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ModifyLivePushAuthKeyResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ModifyPullStreamConfigRequest struct {
	*tchttp.BaseRequest

	// 配置id。
	ConfigId *string `json:"ConfigId" name:"ConfigId"`

	// 源Url。
	FromUrl *string `json:"FromUrl" name:"FromUrl"`

	// 目的Url。
	ToUrl *string `json:"ToUrl" name:"ToUrl"`

	// 区域id,1-深圳,2-上海，3-天津,4-香港。如有改动，需同时传入IspId。
	AreaId *int64 `json:"AreaId" name:"AreaId"`

	// 运营商id,1-电信,2-移动,3-联通,4-其他,AreaId为4的时候,IspId只能为其他。如有改动，需同时传入AreaId。
	IspId *int64 `json:"IspId" name:"IspId"`

	// 开始时间。
	StartTime *string `json:"StartTime" name:"StartTime"`

	// 结束时间，注意：
	// 1. 结束时间必须大于开始时间；
	// 2. 结束时间和开始时间必须大于当前时间；
	// 3. 结束时间 和 开始时间 间隔必须小于七天。
	EndTime *string `json:"EndTime" name:"EndTime"`
}

func (r *ModifyPullStreamConfigRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ModifyPullStreamConfigRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ModifyPullStreamConfigResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *ModifyPullStreamConfigResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ModifyPullStreamConfigResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ModifyPullStreamStatusRequest struct {
	*tchttp.BaseRequest

	// 配置id列表。
	ConfigIds []*string `json:"ConfigIds" name:"ConfigIds" list`

	// 目标状态。0无效，2正在运行，4暂停。
	Status *string `json:"Status" name:"Status"`
}

func (r *ModifyPullStreamStatusRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ModifyPullStreamStatusRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ModifyPullStreamStatusResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *ModifyPullStreamStatusResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ModifyPullStreamStatusResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type PlayAuthKeyInfo struct {

	// 域名。
	DomainName *string `json:"DomainName" name:"DomainName"`

	// 是否启用，0：关闭，1：启用。
	Enable *int64 `json:"Enable" name:"Enable"`

	// 鉴权key。
	AuthKey *string `json:"AuthKey" name:"AuthKey"`

	// 有效时间，单位：秒。
	AuthDelta *uint64 `json:"AuthDelta" name:"AuthDelta"`

	// 鉴权BackKey。
	AuthBackKey *string `json:"AuthBackKey" name:"AuthBackKey"`
}

type PublishTime struct {

	// 推流时间
	// UTC 格式，例如：2018-06-29T19:00:00Z。
	PublishTime *string `json:"PublishTime" name:"PublishTime"`
}

type PullStreamConfig struct {

	// 拉流配置Id。
	ConfigId *string `json:"ConfigId" name:"ConfigId"`

	// 源Url。
	FromUrl *string `json:"FromUrl" name:"FromUrl"`

	// 目的Url。
	ToUrl *string `json:"ToUrl" name:"ToUrl"`

	// 区域名。
	AreaName *string `json:"AreaName" name:"AreaName"`

	// 运营商名。
	IspName *string `json:"IspName" name:"IspName"`

	// 开始时间。
	StartTime *string `json:"StartTime" name:"StartTime"`

	// 结束时间。
	EndTime *string `json:"EndTime" name:"EndTime"`

	// 0无效，1初始状态，2正在运行，3拉起失败，4暂停。
	Status *string `json:"Status" name:"Status"`
}

type PushAuthKeyInfo struct {

	// 域名。
	DomainName *string `json:"DomainName" name:"DomainName"`

	// 是否启用，0：关闭，1：启用。
	Enable *int64 `json:"Enable" name:"Enable"`

	// 主鉴权key。
	MasterAuthKey *string `json:"MasterAuthKey" name:"MasterAuthKey"`

	// 备鉴权key。
	BackupAuthKey *string `json:"BackupAuthKey" name:"BackupAuthKey"`

	// 有效时间，单位：秒。
	AuthDelta *uint64 `json:"AuthDelta" name:"AuthDelta"`
}

type ResumeDelayLiveStreamRequest struct {
	*tchttp.BaseRequest

	// 应用名称。
	AppName *string `json:"AppName" name:"AppName"`

	// 您的加速域名。
	DomainName *string `json:"DomainName" name:"DomainName"`

	// 流名称。
	StreamName *string `json:"StreamName" name:"StreamName"`
}

func (r *ResumeDelayLiveStreamRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ResumeDelayLiveStreamRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ResumeDelayLiveStreamResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *ResumeDelayLiveStreamResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ResumeDelayLiveStreamResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ResumeLiveStreamRequest struct {
	*tchttp.BaseRequest

	// 应用名称。
	AppName *string `json:"AppName" name:"AppName"`

	// 您的加速域名。
	DomainName *string `json:"DomainName" name:"DomainName"`

	// 流名称。
	StreamName *string `json:"StreamName" name:"StreamName"`
}

func (r *ResumeLiveStreamRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ResumeLiveStreamRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ResumeLiveStreamResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *ResumeLiveStreamResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ResumeLiveStreamResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type SetLiveWatermarkStatusRequest struct {
	*tchttp.BaseRequest

	// 水印ID。
	WatermarkId *int64 `json:"WatermarkId" name:"WatermarkId"`

	// 状态。0：停用，1:启用
	Status *int64 `json:"Status" name:"Status"`
}

func (r *SetLiveWatermarkStatusRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *SetLiveWatermarkStatusRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type SetLiveWatermarkStatusResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *SetLiveWatermarkStatusResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *SetLiveWatermarkStatusResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type StopLiveRecordRequest struct {
	*tchttp.BaseRequest

	// 流名称。
	StreamName *string `json:"StreamName" name:"StreamName"`

	// 任务ID，全局唯一标识录制任务。
	TaskId *int64 `json:"TaskId" name:"TaskId"`
}

func (r *StopLiveRecordRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *StopLiveRecordRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type StopLiveRecordResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *StopLiveRecordResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *StopLiveRecordResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type StreamInfo struct {

	// 直播流所属应用名称
	AppName *string `json:"AppName" name:"AppName"`

	// 创建模式
	CreateMode *string `json:"CreateMode" name:"CreateMode"`

	// 创建时间，如: 2018-07-13 14:48:23
	CreateTime *string `json:"CreateTime" name:"CreateTime"`

	// 流状态
	Status *int64 `json:"Status" name:"Status"`

	// 流id
	StreamId *string `json:"StreamId" name:"StreamId"`

	// 流名称
	StreamName *string `json:"StreamName" name:"StreamName"`

	// 水印id
	WaterMarkId *string `json:"WaterMarkId" name:"WaterMarkId"`
}

type StreamName struct {

	// 流名称。
	StreamName *string `json:"StreamName" name:"StreamName"`
}

type StreamOnlineInfo struct {

	// 流名称。
	StreamName *string `json:"StreamName" name:"StreamName"`

	// 推流时间列表
	PublishTimeList []*PublishTime `json:"PublishTimeList" name:"PublishTimeList" list`
}

type UpdateLiveWatermarkRequest struct {
	*tchttp.BaseRequest

	// 水印ID。
	WatermarkId *int64 `json:"WatermarkId" name:"WatermarkId"`

	// 水印图片url。
	PictureUrl *string `json:"PictureUrl" name:"PictureUrl"`

	// 显示位置，X轴偏移。
	XPosition *int64 `json:"XPosition" name:"XPosition"`

	// 显示位置，Y轴偏移。
	YPosition *int64 `json:"YPosition" name:"YPosition"`

	// 水印名称。
	WatermarkName *string `json:"WatermarkName" name:"WatermarkName"`
}

func (r *UpdateLiveWatermarkRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *UpdateLiveWatermarkRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type UpdateLiveWatermarkResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *UpdateLiveWatermarkResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *UpdateLiveWatermarkResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type WatermarkInfo struct {

	// 水印ID。
	WatermarkId *int64 `json:"WatermarkId" name:"WatermarkId"`

	// 水印图片url。
	PictureUrl *string `json:"PictureUrl" name:"PictureUrl"`

	// 显示位置，X轴偏移。
	XPosition *int64 `json:"XPosition" name:"XPosition"`

	// 显示位置，Y轴偏移。
	YPosition *int64 `json:"YPosition" name:"YPosition"`

	// 水印名称。
	WatermarkName *string `json:"WatermarkName" name:"WatermarkName"`

	// 当前状态。0：未使用，1:使用中。
	Status *int64 `json:"Status" name:"Status"`

	// 添加时间。
	CreateTime *string `json:"CreateTime" name:"CreateTime"`
}
