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

package v20180321

import (
    "encoding/json"

    tchttp "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/http"
)

type ApplyBlackListRequest struct {
	*tchttp.BaseRequest

	// 模块
	Module *string `json:"Module" name:"Module"`

	// 操作
	Operation *string `json:"Operation" name:"Operation"`

	// 黑名单列表
	BlackList []*SingleBlackApply `json:"BlackList" name:"BlackList" list`
}

func (r *ApplyBlackListRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ApplyBlackListRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ApplyBlackListResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *ApplyBlackListResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ApplyBlackListResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeRecordsRequest struct {
	*tchttp.BaseRequest

	// 模块
	Module *string `json:"Module" name:"Module"`

	// 操作
	Operation *string `json:"Operation" name:"Operation"`

	// 产品ID
	ProductId *string `json:"ProductId" name:"ProductId"`

	// 案件编号
	AccountNum *string `json:"AccountNum" name:"AccountNum"`

	// 被叫号码
	CalledPhone *string `json:"CalledPhone" name:"CalledPhone"`

	// 查询起始日期
	StartBizDate *string `json:"StartBizDate" name:"StartBizDate"`

	// 查询结束日期
	EndBizDate *string `json:"EndBizDate" name:"EndBizDate"`

	// 分页参数，索引，从0开始
	Offset *string `json:"Offset" name:"Offset"`

	// 分页参数，页长
	Limit *string `json:"Limit" name:"Limit"`
}

func (r *DescribeRecordsRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeRecordsRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeRecordsResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 录音列表。
		RecordList []*SingleRecord `json:"RecordList" name:"RecordList" list`

		// 录音总量。
		TotalCount *int64 `json:"TotalCount" name:"TotalCount"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeRecordsResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeRecordsResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeTaskStatusRequest struct {
	*tchttp.BaseRequest

	// 模块名
	Module *string `json:"Module" name:"Module"`

	// 操作名
	Operation *string `json:"Operation" name:"Operation"`

	// 任务ID
	TaskId *int64 `json:"TaskId" name:"TaskId"`
}

func (r *DescribeTaskStatusRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeTaskStatusRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeTaskStatusResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 任务结果
		TaskResult *string `json:"TaskResult" name:"TaskResult"`

		// 任务类型，010代表上传任务
		TaskType *string `json:"TaskType" name:"TaskType"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeTaskStatusResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeTaskStatusResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DownloadReportRequest struct {
	*tchttp.BaseRequest

	// 模块名
	Module *string `json:"Module" name:"Module"`

	// 操作名
	Operation *string `json:"Operation" name:"Operation"`

	// 报告日期
	ReportDate *string `json:"ReportDate" name:"ReportDate"`
}

func (r *DownloadReportRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DownloadReportRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DownloadReportResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 日报下载地址
		DailyReportUrl *string `json:"DailyReportUrl" name:"DailyReportUrl"`

		// 结果下载地址
		ResultReportUrl *string `json:"ResultReportUrl" name:"ResultReportUrl"`

		// 明细下载地址
		DetailReportUrl *string `json:"DetailReportUrl" name:"DetailReportUrl"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DownloadReportResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DownloadReportResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type SingleBlackApply struct {

	// 黑名单类型，01代表手机号码。
	BlackType *string `json:"BlackType" name:"BlackType"`

	// 操作类型，A为新增，D为删除。
	OperationType *string `json:"OperationType" name:"OperationType"`

	// 黑名单值，BlackType为01时，填写11位手机号码。
	BlackValue *string `json:"BlackValue" name:"BlackValue"`

	// 备注。
	BlackDescription *string `json:"BlackDescription" name:"BlackDescription"`
}

type SingleRecord struct {

	// 案件编号。
	AccountNum *string `json:"AccountNum" name:"AccountNum"`

	// 外呼日期。
	BizDate *string `json:"BizDate" name:"BizDate"`

	// 开始呼叫时间。
	CallStartTime *string `json:"CallStartTime" name:"CallStartTime"`

	// 主叫号码。
	CallerPhone *string `json:"CallerPhone" name:"CallerPhone"`

	// 呼叫方向，O为呼出，I为呼入。
	Direction *string `json:"Direction" name:"Direction"`

	// 通话时长。
	Duration *int64 `json:"Duration" name:"Duration"`

	// 产品ID。
	ProductId *string `json:"ProductId" name:"ProductId"`

	// 录音下载链接。
	RecordCosUrl *string `json:"RecordCosUrl" name:"RecordCosUrl"`
}

type UploadDataFileRequest struct {
	*tchttp.BaseRequest

	// 模块名
	Module *string `json:"Module" name:"Module"`

	// 操作名
	Operation *string `json:"Operation" name:"Operation"`

	// 文件名
	FileName *string `json:"FileName" name:"FileName"`

	// 上传类型，不填默认催收文件，催收文件为data，还款文件为repay。
	UploadModel *string `json:"UploadModel" name:"UploadModel"`

	// 文件，文件与文件地址上传只可选用一种，使用 Content-Type: multipart/form-data 协议来上传二进制流文件。
	File *binary `json:"File" name:"File"`

	// 文件上传地址
	FileUrl *string `json:"FileUrl" name:"FileUrl"`
}

func (r *UploadDataFileRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *UploadDataFileRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type UploadDataFileResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 数据ID
		DataResId *string `json:"DataResId" name:"DataResId"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *UploadDataFileResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *UploadDataFileResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type UploadFileRequest struct {
	*tchttp.BaseRequest

	// 模块名
	Module *string `json:"Module" name:"Module"`

	// 操作名
	Operation *string `json:"Operation" name:"Operation"`

	// 文件上传地址，要求地址协议为HTTPS，且URL端口必须为443
	FileUrl *string `json:"FileUrl" name:"FileUrl"`

	// 文件名
	FileName *string `json:"FileName" name:"FileName"`

	// 文件日期
	FileDate *string `json:"FileDate" name:"FileDate"`
}

func (r *UploadFileRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *UploadFileRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type UploadFileResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 任务ID
		TaskId *int64 `json:"TaskId" name:"TaskId"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *UploadFileResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *UploadFileResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}
