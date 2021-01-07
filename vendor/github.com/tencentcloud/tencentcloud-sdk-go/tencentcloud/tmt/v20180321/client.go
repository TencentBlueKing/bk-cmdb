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
    "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
    tchttp "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/http"
    "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
)

const APIVersion = "2018-03-21"

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


func NewImageTranslateRequest() (request *ImageTranslateRequest) {
    request = &ImageTranslateRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("tmt", APIVersion, "ImageTranslate")
    return
}

func NewImageTranslateResponse() (response *ImageTranslateResponse) {
    response = &ImageTranslateResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 提供中文到英文、英文到中文两种语言的图片翻译服务，可自动识别图片中的文本内容并翻译成目标语言，识别后的文本按行翻译，后续会提供可按段落翻译的版本
func (c *Client) ImageTranslate(request *ImageTranslateRequest) (response *ImageTranslateResponse, err error) {
    if request == nil {
        request = NewImageTranslateRequest()
    }
    response = NewImageTranslateResponse()
    err = c.Send(request, response)
    return
}

func NewLanguageDetectRequest() (request *LanguageDetectRequest) {
    request = &LanguageDetectRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("tmt", APIVersion, "LanguageDetect")
    return
}

func NewLanguageDetectResponse() (response *LanguageDetectResponse) {
    response = &LanguageDetectResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 可自动识别文本内容的语言种类，轻量高效，无需额外实现判断方式，使面向客户的服务体验更佳。 
func (c *Client) LanguageDetect(request *LanguageDetectRequest) (response *LanguageDetectResponse, err error) {
    if request == nil {
        request = NewLanguageDetectRequest()
    }
    response = NewLanguageDetectResponse()
    err = c.Send(request, response)
    return
}

func NewSpeechTranslateRequest() (request *SpeechTranslateRequest) {
    request = &SpeechTranslateRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("tmt", APIVersion, "SpeechTranslate")
    return
}

func NewSpeechTranslateResponse() (response *SpeechTranslateResponse) {
    response = &SpeechTranslateResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 本接口提供音频内文字识别 + 翻译功能，目前开放中到英的语音翻译服务。
// 待识别和翻译的音频文件可以是 pcm、mp3、amr和speex 格式，音频内语音清晰，采用流式传输和翻译的方式。
func (c *Client) SpeechTranslate(request *SpeechTranslateRequest) (response *SpeechTranslateResponse, err error) {
    if request == nil {
        request = NewSpeechTranslateRequest()
    }
    response = NewSpeechTranslateResponse()
    err = c.Send(request, response)
    return
}

func NewTextTranslateRequest() (request *TextTranslateRequest) {
    request = &TextTranslateRequest{
        BaseRequest: &tchttp.BaseRequest{},
    }
    request.Init().WithApiInfo("tmt", APIVersion, "TextTranslate")
    return
}

func NewTextTranslateResponse() (response *TextTranslateResponse) {
    response = &TextTranslateResponse{
        BaseResponse: &tchttp.BaseResponse{},
    }
    return
}

// 提供中文到英文、英文到中文的等多种语言的文本内容翻译服务， 经过大数据语料库、多种解码算法、翻译引擎深度优化，在新闻文章、生活口语等不同语言场景中都有深厚积累，翻译结果专业评价处于行业顶级水平。
func (c *Client) TextTranslate(request *TextTranslateRequest) (response *TextTranslateResponse, err error) {
    if request == nil {
        request = NewTextTranslateRequest()
    }
    response = NewTextTranslateResponse()
    err = c.Send(request, response)
    return
}
