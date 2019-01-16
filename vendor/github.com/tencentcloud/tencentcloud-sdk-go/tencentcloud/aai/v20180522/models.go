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

package v20180522

import (
    "encoding/json"

    tchttp "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/http"
)

type ChatRequest struct {
	*tchttp.BaseRequest

	// 聊天输入文本
	Text *string `json:"Text" name:"Text"`

	// 腾讯云项目 ID，可填 0，总长度不超过 1024 字节。
	ProjectId *int64 `json:"ProjectId" name:"ProjectId"`

	// json格式，比如 {"id":"test","gender":"male"}。记录当前与机器人交互的用户id，非必须但强烈建议传入，否则多轮聊天功能会受影响
	User *string `json:"User" name:"User"`
}

func (r *ChatRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ChatRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ChatResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 聊天输出文本
		Answer *string `json:"Answer" name:"Answer"`

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *ChatResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ChatResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type SentenceRecognitionRequest struct {
	*tchttp.BaseRequest

	// 腾讯云项目 ID，可填 0，总长度不超过 1024 字节。
	ProjectId *uint64 `json:"ProjectId" name:"ProjectId"`

	// 子服务类型。2，一句话识别。
	SubServiceType *uint64 `json:"SubServiceType" name:"SubServiceType"`

	// 引擎类型。8k：电话 8k 通用模型；16k：16k 通用模型。只支持单声道音频识别。
	EngSerViceType *string `json:"EngSerViceType" name:"EngSerViceType"`

	// 语音数据来源。0：语音 URL；1：语音数据（post body）。
	SourceType *uint64 `json:"SourceType" name:"SourceType"`

	// 识别音频的音频格式（支持mp3,wav）。
	VoiceFormat *string `json:"VoiceFormat" name:"VoiceFormat"`

	// 用户端对此任务的唯一标识，用户自助生成，用于用户查找识别结果。
	UsrAudioKey *string `json:"UsrAudioKey" name:"UsrAudioKey"`

	// 语音 URL，公网可下载。当 SourceType 值为 0 时须填写该字段，为 1 时不填；URL 的长度大于 0，小于 2048，需进行urlencode编码。音频时间长度要小于60s。
	Url *string `json:"Url" name:"Url"`

	// 语音数据，当SourceType 值为1时必须填写，为0可不写。要base64编码(采用python语言时注意读取文件应该为string而不是byte，以byte格式读取后要decode()。编码后的数据不可带有回车换行符)。音频数据要小于900k。
	Data *string `json:"Data" name:"Data"`

	// 数据长度，当 SourceType 值为1时必须填写，为0可不写（此数据长度为数据未进行base64编码时的数据长度）。
	DataLen *int64 `json:"DataLen" name:"DataLen"`
}

func (r *SentenceRecognitionRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *SentenceRecognitionRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type SentenceRecognitionResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 识别结果。
		Result *string `json:"Result" name:"Result"`

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *SentenceRecognitionResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *SentenceRecognitionResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type SimultaneousInterpretingRequest struct {
	*tchttp.BaseRequest

	// 腾讯云项目 ID，可填 0，总长度不超过 1024 字节。
	ProjectId *uint64 `json:"ProjectId" name:"ProjectId"`

	// 子服务类型。0：离线语音识别。1：实时流式识别，2，一句话识别。3：同传。
	SubServiceType *uint64 `json:"SubServiceType" name:"SubServiceType"`

	// 识别引擎类型。8k_zh： 8k 中文会场模型；16k_zh：16k 中文会场模型，8k_en： 8k 英文会场模型；16k_en：16k 英文会场模型。当前仅支持16K。
	RecEngineModelType *string `json:"RecEngineModelType" name:"RecEngineModelType"`

	// 语音数据，要base64编码。
	Data *string `json:"Data" name:"Data"`

	// 数据长度。
	DataLen *uint64 `json:"DataLen" name:"DataLen"`

	// 声音id，标识一句话。
	VoiceId *string `json:"VoiceId" name:"VoiceId"`

	// 是否是一句话的结束。
	IsEnd *uint64 `json:"IsEnd" name:"IsEnd"`

	// 声音编码的格式1:pcm，4:speex，6:silk，默认为1。
	VoiceFormat *uint64 `json:"VoiceFormat" name:"VoiceFormat"`

	// 是否需要翻译结果，1表示需要翻译，0是不需要。
	OpenTranslate *uint64 `json:"OpenTranslate" name:"OpenTranslate"`

	// 如果需要翻译，表示源语言类型，可取值：zh，en。
	SourceLanguage *string `json:"SourceLanguage" name:"SourceLanguage"`

	// 如果需要翻译，表示目标语言类型，可取值：zh，en。
	TargetLanguage *string `json:"TargetLanguage" name:"TargetLanguage"`

	// 表明当前语音分片的索引，从0开始
	Seq *uint64 `json:"Seq" name:"Seq"`
}

func (r *SimultaneousInterpretingRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *SimultaneousInterpretingRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type SimultaneousInterpretingResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 语音识别的结果
		AsrText *string `json:"AsrText" name:"AsrText"`

		// 机器翻译的结果
		NmtText *string `json:"NmtText" name:"NmtText"`

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *SimultaneousInterpretingResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *SimultaneousInterpretingResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type TextToVoiceRequest struct {
	*tchttp.BaseRequest

	// 合成语音的源文本
	Text *string `json:"Text" name:"Text"`

	// 一次请求对应一个SessionId，会原样返回，建议传入类似于uuid的字符串防止重复
	SessionId *string `json:"SessionId" name:"SessionId"`

	// 模型类型，1-默认模型
	ModelType *int64 `json:"ModelType" name:"ModelType"`

	// 音量大小，范围：[0，10]，分别对应10个等级的音量，默认为0
	Volume *float64 `json:"Volume" name:"Volume"`

	// 语速，范围：[-2，2]，分别对应不同语速：0.6倍，0.8倍，1.0倍，1.2倍，1.5倍，默认为0
	Speed *float64 `json:"Speed" name:"Speed"`

	// 项目id，用户自定义，默认为0
	ProjectId *int64 `json:"ProjectId" name:"ProjectId"`

	// 音色<li>0-女声1，亲和风格(默认)</li><li>1-男声1，成熟风格</li><li>2-男声2，成熟风格</li>
	VoiceType *int64 `json:"VoiceType" name:"VoiceType"`

	// 主语言类型<li>1-中文(包括粤语)，最大100字符</li><li>2-英文，最大支持400字符</li>
	PrimaryLanguage *uint64 `json:"PrimaryLanguage" name:"PrimaryLanguage"`

	// 音频采样率，16000：16k，8000：8k，默认16k
	SampleRate *uint64 `json:"SampleRate" name:"SampleRate"`
}

func (r *TextToVoiceRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *TextToVoiceRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type TextToVoiceResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// base编码的wav音频
		Audio *string `json:"Audio" name:"Audio"`

		// 一次请求对应一个SessionId
		SessionId *string `json:"SessionId" name:"SessionId"`

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *TextToVoiceResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *TextToVoiceResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}
