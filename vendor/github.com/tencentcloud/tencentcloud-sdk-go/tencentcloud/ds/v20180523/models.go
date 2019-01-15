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

package v20180523

import (
    "encoding/json"

    tchttp "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/http"
)

type CheckVcodeRequest struct {
	*tchttp.BaseRequest

	// 模块名
	Module *string `json:"Module" name:"Module"`

	// 操作名
	Operation *string `json:"Operation" name:"Operation"`

	// 帐号ID
	AccountResId *string `json:"AccountResId" name:"AccountResId"`

	// 合同ID
	ContractResId *string `json:"ContractResId" name:"ContractResId"`

	// 验证码
	VerifyCode *string `json:"VerifyCode" name:"VerifyCode"`
}

func (r *CheckVcodeRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *CheckVcodeRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type CheckVcodeResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *CheckVcodeResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *CheckVcodeResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type CreateContractByUploadRequest struct {
	*tchttp.BaseRequest

	// 模块名
	Module *string `json:"Module" name:"Module"`

	// 操作名
	Operation *string `json:"Operation" name:"Operation"`

	// 签署人信息
	SignInfos []*SignInfo `json:"SignInfos" name:"SignInfos" list`

	// 合同上传链接地址
	ContractFile *string `json:"ContractFile" name:"ContractFile"`

	// 合同名称
	ContractName *string `json:"ContractName" name:"ContractName"`

	// 备注
	Remarks *string `json:"Remarks" name:"Remarks"`

	// 合同发起方帐号ID
	Initiator *string `json:"Initiator" name:"Initiator"`
}

func (r *CreateContractByUploadRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *CreateContractByUploadRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type CreateContractByUploadResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 任务ID
		TaskId *int64 `json:"TaskId" name:"TaskId"`

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *CreateContractByUploadResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *CreateContractByUploadResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type CreateEnterpriseAccountRequest struct {
	*tchttp.BaseRequest

	// 模块名
	Module *string `json:"Module" name:"Module"`

	// 操作名
	Operation *string `json:"Operation" name:"Operation"`

	// 企业用户名称
	Name *string `json:"Name" name:"Name"`

	// 企业用户证件类型，8代表营业执照
	IdentType *int64 `json:"IdentType" name:"IdentType"`

	// 企业用户营业执照号码
	IdentNo *string `json:"IdentNo" name:"IdentNo"`

	// 企业联系电话
	MobilePhone *string `json:"MobilePhone" name:"MobilePhone"`

	// 经办人姓名
	TransactorName *string `json:"TransactorName" name:"TransactorName"`

	// 经办人证件类型，0代表身份证
	TransactorIdentType *int64 `json:"TransactorIdentType" name:"TransactorIdentType"`

	// 经办人证件号码
	TransactorIdentNo *string `json:"TransactorIdentNo" name:"TransactorIdentNo"`

	// 经办人手机号
	TransactorPhone *string `json:"TransactorPhone" name:"TransactorPhone"`
}

func (r *CreateEnterpriseAccountRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *CreateEnterpriseAccountRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type CreateEnterpriseAccountResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 帐号ID
		AccountResId *string `json:"AccountResId" name:"AccountResId"`

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *CreateEnterpriseAccountResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *CreateEnterpriseAccountResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type CreatePersonalAccountRequest struct {
	*tchttp.BaseRequest

	// 模块名
	Module *string `json:"Module" name:"Module"`

	// 操作名
	Operation *string `json:"Operation" name:"Operation"`

	// 个人用户姓名
	Name *string `json:"Name" name:"Name"`

	// 个人用户证件类型。0代表身份证
	IdentType *int64 `json:"IdentType" name:"IdentType"`

	// 个人用户证件号码
	IdentNo *string `json:"IdentNo" name:"IdentNo"`

	// 个人用户手机号
	MobilePhone *string `json:"MobilePhone" name:"MobilePhone"`
}

func (r *CreatePersonalAccountRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *CreatePersonalAccountRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type CreatePersonalAccountResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 账号ID
		AccountResId *string `json:"AccountResId" name:"AccountResId"`

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *CreatePersonalAccountResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *CreatePersonalAccountResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type CreateSealRequest struct {
	*tchttp.BaseRequest

	// 模块名
	Module *string `json:"Module" name:"Module"`

	// 操作名
	Operation *string `json:"Operation" name:"Operation"`

	// 帐号ID
	AccountResId *string `json:"AccountResId" name:"AccountResId"`

	// 签章链接，图片必须为png格式
	ImgUrl *string `json:"ImgUrl" name:"ImgUrl"`
}

func (r *CreateSealRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *CreateSealRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type CreateSealResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 签章ID
		SealResId *string `json:"SealResId" name:"SealResId"`

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *CreateSealResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *CreateSealResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DeleteAccountRequest struct {
	*tchttp.BaseRequest

	// 模块名
	Module *string `json:"Module" name:"Module"`

	// 操作名
	Operation *string `json:"Operation" name:"Operation"`

	// 帐号ID列表
	AccountList []*string `json:"AccountList" name:"AccountList" list`
}

func (r *DeleteAccountRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DeleteAccountRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DeleteAccountResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 删除成功帐号ID列表
		DelSuccessList []*string `json:"DelSuccessList" name:"DelSuccessList" list`

		// 删除失败帐号ID列表
		DelFailedList []*string `json:"DelFailedList" name:"DelFailedList" list`

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DeleteAccountResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DeleteAccountResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DeleteSealRequest struct {
	*tchttp.BaseRequest

	// 模块名
	Module *string `json:"Module" name:"Module"`

	// 操作名
	Operation *string `json:"Operation" name:"Operation"`

	// 帐号ID
	AccountResId *string `json:"AccountResId" name:"AccountResId"`

	// 签章ID
	SealResId *string `json:"SealResId" name:"SealResId"`
}

func (r *DeleteSealRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DeleteSealRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DeleteSealResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 签章ID
		SealResId *string `json:"SealResId" name:"SealResId"`

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DeleteSealResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DeleteSealResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeTaskStatusRequest struct {
	*tchttp.BaseRequest

	// 模块名
	Module *string `json:"Module" name:"Module"`

	// 操作名
	Operation *string `json:"Operation" name:"Operation"`

	// 任务ID
	TaskId *uint64 `json:"TaskId" name:"TaskId"`
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

		// 任务类型，010代表合同上传结果，020代表合同下载结果
		TaskType *string `json:"TaskType" name:"TaskType"`

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
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

type DownloadContractRequest struct {
	*tchttp.BaseRequest

	// 模块名
	Module *string `json:"Module" name:"Module"`

	// 操作名
	Operation *string `json:"Operation" name:"Operation"`

	// 合同ID
	ContractResId *string `json:"ContractResId" name:"ContractResId"`
}

func (r *DownloadContractRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DownloadContractRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DownloadContractResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 任务ID
		TaskId *int64 `json:"TaskId" name:"TaskId"`

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DownloadContractResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DownloadContractResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type SendVcodeRequest struct {
	*tchttp.BaseRequest

	// 模块名
	Module *string `json:"Module" name:"Module"`

	// 操作名
	Operation *string `json:"Operation" name:"Operation"`

	// 合同ID
	ContractResId *string `json:"ContractResId" name:"ContractResId"`

	// 帐号ID
	AccountResId *string `json:"AccountResId" name:"AccountResId"`
}

func (r *SendVcodeRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *SendVcodeRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type SendVcodeResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *SendVcodeResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *SendVcodeResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type SignContractByCoordinateRequest struct {
	*tchttp.BaseRequest

	// 模块名
	Module *string `json:"Module" name:"Module"`

	// 操作名
	Operation *string `json:"Operation" name:"Operation"`

	// 合同ID
	ContractResId *string `json:"ContractResId" name:"ContractResId"`

	// 帐户ID
	AccountResId *string `json:"AccountResId" name:"AccountResId"`

	// 授权时间，格式为年月日时分秒，例20160801095509
	AuthorizationTime *string `json:"AuthorizationTime" name:"AuthorizationTime"`

	// 授权IP地址
	Position *string `json:"Position" name:"Position"`

	// 签署坐标，坐标不得超过合同文件边界
	SignLocations []*SignLocation `json:"SignLocations" name:"SignLocations" list`

	// 印章ID
	SealResId *string `json:"SealResId" name:"SealResId"`
}

func (r *SignContractByCoordinateRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *SignContractByCoordinateRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type SignContractByCoordinateResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *SignContractByCoordinateResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *SignContractByCoordinateResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type SignContractByKeywordRequest struct {
	*tchttp.BaseRequest

	// 模块名
	Module *string `json:"Module" name:"Module"`

	// 操作名
	Operation *string `json:"Operation" name:"Operation"`

	// 合同ID
	ContractResId *string `json:"ContractResId" name:"ContractResId"`

	// 账户ID
	AccountResId *string `json:"AccountResId" name:"AccountResId"`

	// 授权时间，格式为年月日时分秒，例20160801095509
	AuthorizationTime *string `json:"AuthorizationTime" name:"AuthorizationTime"`

	// 授权IP地址
	Position *string `json:"Position" name:"Position"`

	// 签章ID
	SealResId *string `json:"SealResId" name:"SealResId"`

	// 签署关键字，坐标和范围不得超过合同文件边界
	SignKeyword *SignKeyword `json:"SignKeyword" name:"SignKeyword"`
}

func (r *SignContractByKeywordRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *SignContractByKeywordRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type SignContractByKeywordResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *SignContractByKeywordResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *SignContractByKeywordResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type SignInfo struct {

	// 账户ID
	AccountResId *string `json:"AccountResId" name:"AccountResId"`

	// 授权时间，格式为年月日时分秒，例20160801095509
	AuthorizationTime *string `json:"AuthorizationTime" name:"AuthorizationTime"`

	// 授权IP地址
	Location *string `json:"Location" name:"Location"`

	// 签章ID
	SealId *string `json:"SealId" name:"SealId"`

	// 签名图片，优先级比SealId高
	ImageData *string `json:"ImageData" name:"ImageData"`

	// 默认值：1  表示RSA证书， 2 表示国密证书， 参数不传时默认为1
	CertType *int64 `json:"CertType" name:"CertType"`
}

type SignKeyword struct {

	// 关键字
	Keyword *string `json:"Keyword" name:"Keyword"`

	// X轴偏移坐标
	OffsetCoordX *string `json:"OffsetCoordX" name:"OffsetCoordX"`

	// Y轴偏移坐标
	OffsetCoordY *string `json:"OffsetCoordY" name:"OffsetCoordY"`

	// 签章突破宽度
	ImageWidth *string `json:"ImageWidth" name:"ImageWidth"`

	// 签章图片高度
	ImageHeight *string `json:"ImageHeight" name:"ImageHeight"`
}

type SignLocation struct {

	// 签名域页数
	SignOnPage *string `json:"SignOnPage" name:"SignOnPage"`

	// 签名域左下角X轴坐标轴
	SignLocationLBX *string `json:"SignLocationLBX" name:"SignLocationLBX"`

	// 签名域左下角Y轴坐标轴
	SignLocationLBY *string `json:"SignLocationLBY" name:"SignLocationLBY"`

	// 签名域右上角X轴坐标轴
	SignLocationRUX *string `json:"SignLocationRUX" name:"SignLocationRUX"`

	// 签名域右上角Y轴坐标轴
	SignLocationRUY *string `json:"SignLocationRUY" name:"SignLocationRUY"`
}
