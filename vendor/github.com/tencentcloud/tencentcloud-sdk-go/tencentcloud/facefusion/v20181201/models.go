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

package v20181201

import (
    "encoding/json"

    tchttp "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/http"
)

type FaceFusionRequest struct {
	*tchttp.BaseRequest

	// 活动 ID，请在人脸融合控制台查看。
	ProjectId *string `json:"ProjectId" name:"ProjectId"`

	// 素材 ID，请在人脸融合控制台查看。
	ModelId *string `json:"ModelId" name:"ModelId"`

	// 图片 base64 数据。请确保人脸为正脸，无旋转。若某些手机拍摄后人脸被旋转，请使用图片的 EXIF 信息对图片进行旋转处理。图片大小不超过 500k，分辨率不超过 1080*1080。
	Image *string `json:"Image" name:"Image"`

	// 返回图像方式（url 或 base64) ，二选一。当前仅支持 url 方式，base64 方式后期开放。
	RspImgType *string `json:"RspImgType" name:"RspImgType"`

	// 0表示不需要鉴黄，1表示需要鉴黄。2018年12月1号以前创建的活动默认值为0，其他情况默认值为1.
	PornDetect *int64 `json:"PornDetect" name:"PornDetect"`

	// 0表示不需要鉴政治，1表示需要鉴政。2018年12月1号以前创建的活动默认值为0，其他情况默认值为1。鉴政接口同时会对名人明星进行识别，您可以根据实际需要过滤。
	CelebrityIdentify *int64 `json:"CelebrityIdentify" name:"CelebrityIdentify"`
}

func (r *FaceFusionRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *FaceFusionRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type FaceFusionResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// RspImgType 为 url 时，返回结果的 url， RspImgType 为 base64 时返回 base64 数据。当前仅支持 url 方式，base64 方式后期开放。
		Image *string `json:"Image" name:"Image"`

		// 鉴黄鉴政结果
		ReviewResultSet []*FuseFaceReviewResult `json:"ReviewResultSet" name:"ReviewResultSet" list`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *FaceFusionResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *FaceFusionResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type FuseFaceReviewDetail struct {

	// 鉴政使用字段, 为职业属性,其他审核结果对应上一级category
	Field *string `json:"Field" name:"Field"`

	// 人员名称
	Label *string `json:"Label" name:"Label"`

	// 对应识别label的置信度
	Confidence *float64 `json:"Confidence" name:"Confidence"`

	// 对应子类审核建议: block,review,pass
	Suggestion *string `json:"Suggestion" name:"Suggestion"`
}

type FuseFaceReviewResult struct {

	// 对应的类别名称 porn, politics, terror
	Category *string `json:"Category" name:"Category"`

	// 对应子类别状态码
	Code *string `json:"Code" name:"Code"`

	// 对应子类别状态码信息描述
	CodeDescription *string `json:"CodeDescription" name:"CodeDescription"`

	// 对应识别种类的置信度
	Confidence *float64 `json:"Confidence" name:"Confidence"`

	// 此字段为保留字段，目前统一返回pass。请使用FusefaceReviewDetail中的“Suggestion”字段
	Suggestion *string `json:"Suggestion" name:"Suggestion"`

	// 审核详细内容
	DetailSet []*FuseFaceReviewDetail `json:"DetailSet" name:"DetailSet" list`
}
