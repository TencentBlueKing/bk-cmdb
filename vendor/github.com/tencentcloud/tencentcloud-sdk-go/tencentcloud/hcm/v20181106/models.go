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

package v20181106

import (
    "encoding/json"

    tchttp "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/http"
)

type EvaluationRequest struct {
	*tchttp.BaseRequest

	// 图片唯一标识，一张图片一个SessionId；
	SessionId *string `json:"SessionId" name:"SessionId"`

	// 图片数据，需要使用base64对图片的二进制数据进行编码；
	Image *string `json:"Image" name:"Image"`
}

func (r *EvaluationRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *EvaluationRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type EvaluationResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 图片唯一标识，一张图片一个SessionId；
		SessionId *string `json:"SessionId" name:"SessionId"`

		// 识别出的算式信息；
		Items []*Item `json:"Items" name:"Items" list`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *EvaluationResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *EvaluationResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type Item struct {

	// 识别的算式是否正确
	Item *string `json:"Item" name:"Item"`

	// 识别的算式
	ItemString *string `json:"ItemString" name:"ItemString"`

	// 识别的算式在图片上的位置信息
	ItemCoord *ItemCoord `json:"ItemCoord" name:"ItemCoord"`
}

type ItemCoord struct {

	// 算式高度
	Height *int64 `json:"Height" name:"Height"`

	// 算式宽度
	Width *int64 `json:"Width" name:"Width"`

	// 算式图的左上角横坐标
	X *int64 `json:"X" name:"X"`

	// 算式图的左上角纵坐标
	Y *int64 `json:"Y" name:"Y"`
}
