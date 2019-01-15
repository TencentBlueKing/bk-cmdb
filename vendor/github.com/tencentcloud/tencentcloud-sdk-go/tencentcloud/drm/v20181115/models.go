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

package v20181115

import (
    "encoding/json"

    tchttp "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/http"
)

type CreateLicenseRequest struct {
	*tchttp.BaseRequest

	// DRM方案类型，接口取值：WIDEVINE，FAIRPLAY。
	DrmType *string `json:"DrmType" name:"DrmType"`

	// Base64编码的终端设备License Request数据。
	LicenseRequest *string `json:"LicenseRequest" name:"LicenseRequest"`

	// 内容类型，接口取值：VodVideo,LiveVideo。
	ContentType *string `json:"ContentType" name:"ContentType"`

	// 授权播放的Track列表。
	// 该值为空时，默认授权所有track播放。
	Tracks []*string `json:"Tracks" name:"Tracks" list`

	// 播放策略参数。
	PlaybackPolicy *PlaybackPolicy `json:"PlaybackPolicy" name:"PlaybackPolicy"`
}

func (r *CreateLicenseRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *CreateLicenseRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type CreateLicenseResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// Base64 编码的许可证二进制数据。
		License *string `json:"License" name:"License"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *CreateLicenseResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *CreateLicenseResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeKeysRequest struct {
	*tchttp.BaseRequest

	// 使用的DRM方案类型，接口取值WIDEVINE、FAIRPLAY、NORMALAES。
	DrmType *string `json:"DrmType" name:"DrmType"`

	// 加密的track列表，接口取值VIDEO、AUDIO。
	Tracks []*string `json:"Tracks" name:"Tracks" list`

	// 内容类型。接口取值VodVideo,LiveVideo
	ContentType *string `json:"ContentType" name:"ContentType"`

	// Base64编码的Rsa公钥，用来加密出参中的SessionKey。
	// 如果该参数为空，则出参中SessionKey为明文。
	RsaPublicKey *string `json:"RsaPublicKey" name:"RsaPublicKey"`

	// 一个加密内容的唯一标识。
	// 如果该参数为空，则后台自动生成
	ContentId *string `json:"ContentId" name:"ContentId"`
}

func (r *DescribeKeysRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeKeysRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeKeysResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 加密密钥列表
		Keys []*Key `json:"Keys" name:"Keys" list`

		// 用来加密密钥。
	// 如果入参中带有RsaPublicKey，则SessionKey为使用Rsa公钥加密后的二进制数据，Base64编码字符串。
	// 如果入参中没有RsaPublicKey，则SessionKey为原始数据的字符串形式。
		SessionKey *string `json:"SessionKey" name:"SessionKey"`

		// 内容ID
		ContentId *string `json:"ContentId" name:"ContentId"`

		// Widevine方案的Pssh数据，Base64编码。
	// Fairplay方案无该值。
		Pssh *string `json:"Pssh" name:"Pssh"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeKeysResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeKeysResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DrmOutputObject struct {

	// 输出的桶名称。
	BucketName *string `json:"BucketName" name:"BucketName"`

	// 输出的对象名称。
	ObjectName *string `json:"ObjectName" name:"ObjectName"`

	// 输出对象参数。
	Para *DrmOutputPara `json:"Para" name:"Para"`
}

type DrmOutputPara struct {

	// 内容类型。例:video，audio，mpd，m3u8
	Type *string `json:"Type" name:"Type"`

	// 语言,例: en, zh-cn
	Language *string `json:"Language" name:"Language"`
}

type DrmSourceObject struct {

	// 输入的桶名称。
	BucketName *string `json:"BucketName" name:"BucketName"`

	// 输入对象名称。
	ObjectName *string `json:"ObjectName" name:"ObjectName"`
}

type Key struct {

	// 加密track类型。
	Track *string `json:"Track" name:"Track"`

	// 密钥ID。
	KeyId *string `json:"KeyId" name:"KeyId"`

	// 原始Key使用AES-128 ECB模式和SessionKey加密的后的二进制数据，Base64编码的字符串。
	Key *string `json:"Key" name:"Key"`

	// 原始IV使用AES-128 ECB模式和SessionKey加密的后的二进制数据，Base64编码的字符串。
	Iv *string `json:"Iv" name:"Iv"`
}

type PlaybackPolicy struct {

	// 播放许可证的有效期
	LicenseDurationSeconds *uint64 `json:"LicenseDurationSeconds" name:"LicenseDurationSeconds"`

	// 开始播放后，允许最长播放时间
	PlaybackDurationSeconds *uint64 `json:"PlaybackDurationSeconds" name:"PlaybackDurationSeconds"`
}

type StartEncryptionRequest struct {
	*tchttp.BaseRequest

	// cos的end point。
	CosEndPoint *string `json:"CosEndPoint" name:"CosEndPoint"`

	// cos api密钥id。
	CosSecretId *string `json:"CosSecretId" name:"CosSecretId"`

	// cos api密钥。
	CosSecretKey *string `json:"CosSecretKey" name:"CosSecretKey"`

	// 使用的DRM方案类型，接口取值WIDEVINE,FAIRPLAY
	DrmType *string `json:"DrmType" name:"DrmType"`

	// 存储在COS上的原始内容信息
	SourceObject *DrmSourceObject `json:"SourceObject" name:"SourceObject"`

	// 加密后的内容存储到COS的对象
	OutputObjects []*DrmOutputObject `json:"OutputObjects" name:"OutputObjects" list`
}

func (r *StartEncryptionRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *StartEncryptionRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type StartEncryptionResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *StartEncryptionResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *StartEncryptionResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}
