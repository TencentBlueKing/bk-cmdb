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

type ImageRecord struct {

	// 图片翻译结果
	Value []*ItemValue `json:"Value" name:"Value" list`
}

type ImageTranslateRequest struct {
	*tchttp.BaseRequest

	// 唯一id，返回时原样返回
	SessionUuid *string `json:"SessionUuid" name:"SessionUuid"`

	// doc:文档扫描
	Scene *string `json:"Scene" name:"Scene"`

	// 图片数据的Base64字符串
	Data *string `json:"Data" name:"Data"`

	// 源语言，支持语言列表<li> zh : 中文 </li> <li> en : 英文 </li>
	Source *string `json:"Source" name:"Source"`

	// 目标语言，支持语言列表<li> zh : 中文 </li> <li> en : 英文 </li>
	Target *string `json:"Target" name:"Target"`

	// 项目id
	ProjectId *int64 `json:"ProjectId" name:"ProjectId"`
}

func (r *ImageTranslateRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ImageTranslateRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ImageTranslateResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 请求的SessionUuid返回
		SessionUuid *string `json:"SessionUuid" name:"SessionUuid"`

		// 源语言
		Source *string `json:"Source" name:"Source"`

		// 目标语言
		Target *string `json:"Target" name:"Target"`

		// 图片翻译结果，翻译结果按识别的文本每一行独立翻译，后续会推出按段落划分并翻译的版本
		ImageRecord *ImageRecord `json:"ImageRecord" name:"ImageRecord"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *ImageTranslateResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ImageTranslateResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ItemValue struct {

	// 识别出的源文
	SourceText *string `json:"SourceText" name:"SourceText"`

	// 翻译后的译文
	TargetText *string `json:"TargetText" name:"TargetText"`

	// X坐标
	X *int64 `json:"X" name:"X"`

	// Y坐标
	Y *int64 `json:"Y" name:"Y"`

	// 宽度
	W *int64 `json:"W" name:"W"`

	// 高度
	H *int64 `json:"H" name:"H"`
}

type LanguageDetectRequest struct {
	*tchttp.BaseRequest

	// 待识别的文本，文本统一使用utf-8格式编码，非utf-8格式编码字符会翻译失败
	Text *string `json:"Text" name:"Text"`

	// 项目id
	ProjectId *int64 `json:"ProjectId" name:"ProjectId"`
}

func (r *LanguageDetectRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *LanguageDetectRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type LanguageDetectResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 识别出的语言种类，参考语言列表
	// <li> zh : 中文 </li> <li> en : 英文 </li><li> jp : 日语 </li> <li> kr : 韩语 </li><li> de : 德语 </li><li> fr : 法语 </li><li> es : 西班牙文 </li> <li> it : 意大利文 </li><li> tr : 土耳其文 </li><li> ru : 俄文 </li><li> pt : 葡萄牙文 </li><li> vi : 越南文 </li><li> id : 印度尼西亚文 </li><li> ms : 马来西亚文 </li><li> th : 泰文 </li>
		Lang *string `json:"Lang" name:"Lang"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *LanguageDetectResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *LanguageDetectResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type SpeechTranslateRequest struct {
	*tchttp.BaseRequest

	// 一段完整的语音对应一个SessionUuid
	SessionUuid *string `json:"SessionUuid" name:"SessionUuid"`

	// 音频中的语言类型，支持语言列表<li> zh : 中文 </li> <li> en : 英文 </li>
	Source *string `json:"Source" name:"Source"`

	// 翻译目标语⾔言类型 ，支持的语言列表<li> zh : 中文 </li> <li> en : 英文 </li>
	Target *string `json:"Target" name:"Target"`

	// pcm : 146   amr : 33554432   mp3 : 83886080
	AudioFormat *int64 `json:"AudioFormat" name:"AudioFormat"`

	// 语音分片的序号，从0开始
	Seq *int64 `json:"Seq" name:"Seq"`

	// 是否最后一片语音分片，0-否，1-是
	IsEnd *int64 `json:"IsEnd" name:"IsEnd"`

	// 语音分片内容的base64字符串，音频内容应含有效并可识别的文本
	Data *string `json:"Data" name:"Data"`

	// 项目id，用户可自定义
	ProjectId *int64 `json:"ProjectId" name:"ProjectId"`
}

func (r *SpeechTranslateRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *SpeechTranslateRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type SpeechTranslateResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 请求的SessionUuid直接返回
		SessionUuid *string `json:"SessionUuid" name:"SessionUuid"`

		// 语音识别状态 1-进行中 0-完成
		RecognizeStatus *int64 `json:"RecognizeStatus" name:"RecognizeStatus"`

		// 识别出的源文
		SourceText *string `json:"SourceText" name:"SourceText"`

		// 翻译出的译文
		TargetText *string `json:"TargetText" name:"TargetText"`

		// 第几个语音分片
		Seq *int64 `json:"Seq" name:"Seq"`

		// 源语言
		Source *string `json:"Source" name:"Source"`

		// 目标语言
		Target *string `json:"Target" name:"Target"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *SpeechTranslateResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *SpeechTranslateResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type TextTranslateRequest struct {
	*tchttp.BaseRequest

	// 待翻译的文本，文本统一使用utf-8格式编码，非utf-8格式编码字符会翻译失败，请传入有效文本，html标记等非常规翻译文本会翻译失败
	SourceText *string `json:"SourceText" name:"SourceText"`

	// 源语言，参照Target支持语言列表
	Source *string `json:"Source" name:"Source"`

	// 目标语言，参照支持语言列表
	// <li> zh : 中文 </li> <li> en : 英文 </li><li> jp : 日语 </li> <li> kr : 韩语 </li><li> de : 德语 </li><li> fr : 法语 </li><li> es : 西班牙文 </li> <li> it : 意大利文 </li><li> tr : 土耳其文 </li><li> ru : 俄文 </li><li> pt : 葡萄牙文 </li><li> vi : 越南文 </li><li> id : 印度尼西亚文 </li><li> ms : 马来西亚文 </li><li> th : 泰文 </li><li> auto : 自动识别源语言，只能用于source字段 </li>
	Target *string `json:"Target" name:"Target"`

	// 项目id
	ProjectId *int64 `json:"ProjectId" name:"ProjectId"`
}

func (r *TextTranslateRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *TextTranslateRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type TextTranslateResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 翻译后的文本
		TargetText *string `json:"TargetText" name:"TargetText"`

		// 源语言，详见入参Target
		Source *string `json:"Source" name:"Source"`

		// 目标语言，详见入参Target
		Target *string `json:"Target" name:"Target"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *TextTranslateResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *TextTranslateResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}
