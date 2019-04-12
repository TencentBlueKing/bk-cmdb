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

package v20180416

import (
    "encoding/json"

    tchttp "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/http"
)

type Code struct {

	// 对象存储桶名称
	CosBucketName *string `json:"CosBucketName" name:"CosBucketName"`

	// 对象存储对象路径
	CosObjectName *string `json:"CosObjectName" name:"CosObjectName"`

	// 包含函数代码文件及其依赖项的 zip 格式文件，使用该接口时要求将 zip 文件的内容转成 base64 编码，最大支持20M
	ZipFile *string `json:"ZipFile" name:"ZipFile"`

	// 对象存储的地域，地域为北京时需要传入ap-beijing,北京一区时需要传递ap-beijing-1，其他的地域不需要传递。
	CosBucketRegion *string `json:"CosBucketRegion" name:"CosBucketRegion"`
}

type CreateFunctionRequest struct {
	*tchttp.BaseRequest

	// 创建的函数名称，函数名称支持26个英文字母大小写、数字、连接符和下划线，第一个字符只能以字母开头，最后一个字符不能为连接符或者下划线，名称长度2-60
	FunctionName *string `json:"FunctionName" name:"FunctionName"`

	// 函数的代码. 注意：不能同时指定Cos与ZipFile
	Code *Code `json:"Code" name:"Code"`

	// 函数处理方法名称，名称格式支持 "文件名称.方法名称" 形式，文件名称和函数名称之间以"."隔开，文件名称和函数名称要求以字母开始和结尾，中间允许插入字母、数字、下划线和连接符，文件名称和函数名字的长度要求是 2-60 个字符
	Handler *string `json:"Handler" name:"Handler"`

	// 函数描述,最大支持 1000 个英文字母、数字、空格、逗号、换行符和英文句号，支持中文
	Description *string `json:"Description" name:"Description"`

	// 函数运行时内存大小，默认为 128M，可选范围 128MB-1536MB，并且以 128MB 为阶梯
	MemorySize *int64 `json:"MemorySize" name:"MemorySize"`

	// 函数最长执行时间，单位为秒，可选值范围 1-300 秒，默认为 3 秒
	Timeout *int64 `json:"Timeout" name:"Timeout"`

	// 函数的环境变量
	Environment *Environment `json:"Environment" name:"Environment"`

	// 函数运行环境，目前仅支持 Python2.7，Python3.6，Nodejs6.10， PHP5， PHP7，Golang1 和 Java8，默认Python2.7
	Runtime *string `json:"Runtime" name:"Runtime"`

	// 函数的私有网络配置
	VpcConfig *VpcConfig `json:"VpcConfig" name:"VpcConfig"`
}

func (r *CreateFunctionRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *CreateFunctionRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type CreateFunctionResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *CreateFunctionResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *CreateFunctionResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type CreateTriggerRequest struct {
	*tchttp.BaseRequest

	// 新建触发器绑定的函数名称
	FunctionName *string `json:"FunctionName" name:"FunctionName"`

	// 新建触发器名称。如果是定时触发器，名称支持英文字母、数字、连接符和下划线，最长100个字符；如果是其他触发器，见具体触发器绑定参数的说明
	TriggerName *string `json:"TriggerName" name:"TriggerName"`

	// 触发器类型，目前支持 cos 、cmq、 timers、 ckafka类型
	Type *string `json:"Type" name:"Type"`

	// 触发器对应的参数，如果是 timer 类型的触发器其内容是 Linux cron 表达式，如果是其他触发器，见具体触发器说明
	TriggerDesc *string `json:"TriggerDesc" name:"TriggerDesc"`

	// 函数的版本
	Qualifier *string `json:"Qualifier" name:"Qualifier"`
}

func (r *CreateTriggerRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *CreateTriggerRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type CreateTriggerResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *CreateTriggerResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *CreateTriggerResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DeleteFunctionRequest struct {
	*tchttp.BaseRequest

	// 要删除的函数名称
	FunctionName *string `json:"FunctionName" name:"FunctionName"`
}

func (r *DeleteFunctionRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DeleteFunctionRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DeleteFunctionResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DeleteFunctionResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DeleteFunctionResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DeleteTriggerRequest struct {
	*tchttp.BaseRequest

	// 函数的名称
	FunctionName *string `json:"FunctionName" name:"FunctionName"`

	// 要删除的触发器名称
	TriggerName *string `json:"TriggerName" name:"TriggerName"`

	// 要删除的触发器类型，目前支持 cos 、cmq、 timer、ckafka 类型
	Type *string `json:"Type" name:"Type"`

	// 如果删除的触发器类型为 COS 触发器，该字段为必填值，存放 JSON 格式的数据 {"event":"cos:ObjectCreated:*"}，数据内容和 SetTrigger 接口中该字段的格式相同；如果删除的触发器类型为定时触发器或 CMQ 触发器，可以不指定该字段
	TriggerDesc *string `json:"TriggerDesc" name:"TriggerDesc"`

	// 函数的版本信息
	Qualifier *string `json:"Qualifier" name:"Qualifier"`
}

func (r *DeleteTriggerRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DeleteTriggerRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DeleteTriggerResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DeleteTriggerResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DeleteTriggerResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type Environment struct {

	// 环境变量数组
	Variables []*Variable `json:"Variables" name:"Variables" list`
}

type Filter struct {

	// filter.RetCode=not0 表示只返回错误日志，filter.RetCode=is0 表示只返回正确日志，无输入则返回所有日志。
	RetCode *string `json:"RetCode" name:"RetCode"`
}

type Function struct {

	// 修改时间
	ModTime *string `json:"ModTime" name:"ModTime"`

	// 创建时间
	AddTime *string `json:"AddTime" name:"AddTime"`

	// 运行时
	Runtime *string `json:"Runtime" name:"Runtime"`

	// 函数名称
	FunctionName *string `json:"FunctionName" name:"FunctionName"`

	// 函数ID
	FunctionId *string `json:"FunctionId" name:"FunctionId"`

	// 命名空间
	Namespace *string `json:"Namespace" name:"Namespace"`
}

type FunctionLog struct {

	// 函数的名称
	FunctionName *string `json:"FunctionName" name:"FunctionName"`

	// 函数执行完成后的返回值
	RetMsg *string `json:"RetMsg" name:"RetMsg"`

	// 执行该函数对应的requestId
	RequestId *string `json:"RequestId" name:"RequestId"`

	// 函数开始执行时的时间点
	StartTime *string `json:"StartTime" name:"StartTime"`

	// 函数执行结果，如果是 0 表示执行成功，其他值表示失败
	RetCode *int64 `json:"RetCode" name:"RetCode"`

	// 函数调用是否结束，如果是 1 表示执行结束，其他值表示调用异常
	InvokeFinished *int64 `json:"InvokeFinished" name:"InvokeFinished"`

	// 函数执行耗时，单位为 ms
	Duration *float64 `json:"Duration" name:"Duration"`

	// 函数计费时间，根据 duration 向上取最近的 100ms，单位为ms
	BillDuration *int64 `json:"BillDuration" name:"BillDuration"`

	// 函数执行时消耗实际内存大小，单位为 Byte
	MemUsage *int64 `json:"MemUsage" name:"MemUsage"`

	// 函数执行过程中的日志输出
	Log *string `json:"Log" name:"Log"`
}

type GetFunctionLogsRequest struct {
	*tchttp.BaseRequest

	// 函数的名称
	FunctionName *string `json:"FunctionName" name:"FunctionName"`

	// 数据的偏移量，Offset+Limit不能大于10000
	Offset *int64 `json:"Offset" name:"Offset"`

	// 返回数据的长度，Offset+Limit不能大于10000
	Limit *int64 `json:"Limit" name:"Limit"`

	// 以升序还是降序的方式对日志进行排序，可选值 desc和 acs
	Order *string `json:"Order" name:"Order"`

	// 根据某个字段排序日志,支持以下字段：startTime、functionName、requestId、duration和 memUsage
	OrderBy *string `json:"OrderBy" name:"OrderBy"`

	// 日志过滤条件。可用来区分正确和错误日志，filter.retCode=not0 表示只返回错误日志，filter.retCode=is0 表示只返回正确日志，不传，则返回所有日志
	Filter *Filter `json:"Filter" name:"Filter"`

	// 函数的版本
	Qualifier *string `json:"Qualifier" name:"Qualifier"`

	// 执行该函数对应的requestId
	FunctionRequestId *string `json:"FunctionRequestId" name:"FunctionRequestId"`

	// 查询的具体日期，例如：2017-05-16 20:00:00，只能与endtime相差一天之内
	StartTime *string `json:"StartTime" name:"StartTime"`

	// 查询的具体日期，例如：2017-05-16 20:59:59，只能与startTime相差一天之内
	EndTime *string `json:"EndTime" name:"EndTime"`
}

func (r *GetFunctionLogsRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *GetFunctionLogsRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type GetFunctionLogsResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 函数日志的总数
		TotalCount *int64 `json:"TotalCount" name:"TotalCount"`

		// 函数日志信息
		Data []*FunctionLog `json:"Data" name:"Data" list`

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *GetFunctionLogsResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *GetFunctionLogsResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type GetFunctionRequest struct {
	*tchttp.BaseRequest

	// 需要获取详情的函数名称
	FunctionName *string `json:"FunctionName" name:"FunctionName"`

	// 函数的版本号
	Qualifier *string `json:"Qualifier" name:"Qualifier"`

	// 是否显示代码, TRUE表示显示代码，FALSE表示不显示代码,大于1M的入口文件不会显示
	ShowCode *string `json:"ShowCode" name:"ShowCode"`
}

func (r *GetFunctionRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *GetFunctionRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type GetFunctionResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 函数的最后修改时间
		ModTime *string `json:"ModTime" name:"ModTime"`

		// 函数的代码
		CodeInfo *string `json:"CodeInfo" name:"CodeInfo"`

		// 函数的描述信息
		Description *string `json:"Description" name:"Description"`

		// 函数的触发器列表
		Triggers []*Trigger `json:"Triggers" name:"Triggers" list`

		// 函数的入口
		Handler *string `json:"Handler" name:"Handler"`

		// 函数代码大小
		CodeSize *int64 `json:"CodeSize" name:"CodeSize"`

		// 函数的超时时间
		Timeout *int64 `json:"Timeout" name:"Timeout"`

		// 函数的版本
		FunctionVersion *string `json:"FunctionVersion" name:"FunctionVersion"`

		// 函数的最大可用内存
		MemorySize *int64 `json:"MemorySize" name:"MemorySize"`

		// 函数的运行环境
		Runtime *string `json:"Runtime" name:"Runtime"`

		// 函数的名称
		FunctionName *string `json:"FunctionName" name:"FunctionName"`

		// 函数的私有网络
		VpcConfig *VpcConfig `json:"VpcConfig" name:"VpcConfig"`

		// 是否使用GPU
		UseGpu *string `json:"UseGpu" name:"UseGpu"`

		// 函数的环境变量
		Environment *Environment `json:"Environment" name:"Environment"`

		// 代码是否正确
		CodeResult *string `json:"CodeResult" name:"CodeResult"`

		// 代码错误信息
		CodeError *string `json:"CodeError" name:"CodeError"`

		// 代码错误码
		ErrNo *int64 `json:"ErrNo" name:"ErrNo"`

		// 函数的命名空间
		Namespace *string `json:"Namespace" name:"Namespace"`

		// 函数绑定的角色
		Role *string `json:"Role" name:"Role"`

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *GetFunctionResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *GetFunctionResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type InvokeRequest struct {
	*tchttp.BaseRequest

	// 函数名称
	FunctionName *string `json:"FunctionName" name:"FunctionName"`

	// RequestResponse(同步) 和 Event(异步)，默认为同步
	InvocationType *string `json:"InvocationType" name:"InvocationType"`

	// 触发函数的版本号
	Qualifier *string `json:"Qualifier" name:"Qualifier"`

	// 运行函数时的参数，以json格式传入，最大支持的参数长度是 1M
	ClientContext *string `json:"ClientContext" name:"ClientContext"`

	// 同步调用时指定该字段，返回值会包含4K的日志，可选值为None和Tail，默认值为None。当该值为Tail时，返回参数中的logMsg字段会包含对应的函数执行日志
	LogType *string `json:"LogType" name:"LogType"`

	// 命名空间
	Namespace *string `json:"Namespace" name:"Namespace"`
}

func (r *InvokeRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *InvokeRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type InvokeResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 函数执行结果
		Result *Result `json:"Result" name:"Result"`

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *InvokeResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *InvokeResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ListFunctionsRequest struct {
	*tchttp.BaseRequest

	// 以升序还是降序的方式返回结果，可选值 ASC 和 DESC
	Order *string `json:"Order" name:"Order"`

	// 根据哪个字段进行返回结果排序,支持以下字段：AddTime, ModTime, FunctionName
	Orderby *string `json:"Orderby" name:"Orderby"`

	// 数据偏移量，默认值为 0
	Offset *int64 `json:"Offset" name:"Offset"`

	// 返回数据长度，默认值为 20
	Limit *int64 `json:"Limit" name:"Limit"`

	// 支持FunctionName模糊匹配
	SearchKey *string `json:"SearchKey" name:"SearchKey"`
}

func (r *ListFunctionsRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ListFunctionsRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ListFunctionsResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 函数列表
		Functions []*Function `json:"Functions" name:"Functions" list`

		// 总数
		TotalCount *int64 `json:"TotalCount" name:"TotalCount"`

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *ListFunctionsResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ListFunctionsResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type Result struct {

	// 表示执行过程中的日志输出，异步调用返回为空
	Log *string `json:"Log" name:"Log"`

	// 表示执行函数的返回，异步调用返回为空
	RetMsg *string `json:"RetMsg" name:"RetMsg"`

	// 表示执行函数的错误返回信息，异步调用返回为空
	ErrMsg *string `json:"ErrMsg" name:"ErrMsg"`

	// 执行函数时的内存大小，单位为Byte，异步调用返回为空
	MemUsage *int64 `json:"MemUsage" name:"MemUsage"`

	// 表示执行函数的耗时，单位是毫秒，异步调用返回为空
	Duration *float64 `json:"Duration" name:"Duration"`

	// 表示函数的计费耗时，单位是毫秒，异步调用返回为空
	BillDuration *int64 `json:"BillDuration" name:"BillDuration"`

	// 此次函数执行的Id
	FunctionRequestId *string `json:"FunctionRequestId" name:"FunctionRequestId"`

	// 0为正确，异步调用返回为空
	InvokeResult *int64 `json:"InvokeResult" name:"InvokeResult"`
}

type Trigger struct {

	// 触发器最后修改时间
	ModTime *string `json:"ModTime" name:"ModTime"`

	// 触发器类型
	Type *string `json:"Type" name:"Type"`

	// 触发器详细配置
	TriggerDesc *string `json:"TriggerDesc" name:"TriggerDesc"`

	// 触发器名称
	TriggerName *string `json:"TriggerName" name:"TriggerName"`

	// 触发器创建时间
	AddTime *string `json:"AddTime" name:"AddTime"`
}

type UpdateFunctionCodeRequest struct {
	*tchttp.BaseRequest

	// 函数处理方法名称。名称格式支持“文件名称.函数名称”形式，文件名称和函数名称之间以"."隔开，文件名称和函数名称要求以字母开始和结尾，中间允许插入字母、数字、下划线和连接符，文件名称和函数名字的长度要求 2-60 个字符
	Handler *string `json:"Handler" name:"Handler"`

	// 要修改的函数名称
	FunctionName *string `json:"FunctionName" name:"FunctionName"`

	// 对象存储桶名称
	CosBucketName *string `json:"CosBucketName" name:"CosBucketName"`

	// 对象存储对象路径
	CosObjectName *string `json:"CosObjectName" name:"CosObjectName"`

	// 包含函数代码文件及其依赖项的 zip 格式文件，使用该接口时要求将 zip 文件的内容转成 base64 编码，最大支持20M
	ZipFile *string `json:"ZipFile" name:"ZipFile"`

	// 对象存储的地域，地域为北京时需要传入ap-beijing,北京一区时需要传递ap-beijing-1，其他的地域不需要传递。
	CosBucketRegion *string `json:"CosBucketRegion" name:"CosBucketRegion"`
}

func (r *UpdateFunctionCodeRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *UpdateFunctionCodeRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type UpdateFunctionCodeResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *UpdateFunctionCodeResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *UpdateFunctionCodeResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type UpdateFunctionConfigurationRequest struct {
	*tchttp.BaseRequest

	// 要修改的函数名称
	FunctionName *string `json:"FunctionName" name:"FunctionName"`

	// 函数描述。最大支持 1000 个英文字母、数字、空格、逗号和英文句号，支持中文
	Description *string `json:"Description" name:"Description"`

	// 函数运行时内存大小，默认为 128 M，可选范 128 M-1536 M
	MemorySize *int64 `json:"MemorySize" name:"MemorySize"`

	// 函数最长执行时间，单位为秒，可选值范 1-300 秒，默认为 3 秒
	Timeout *int64 `json:"Timeout" name:"Timeout"`

	// 函数运行环境，目前仅支持 Python2.7，Python3.6，Nodejs6.10，PHP5， PHP7，Golang1 和 Java8
	Runtime *string `json:"Runtime" name:"Runtime"`

	// 函数的环境变量
	Environment *Environment `json:"Environment" name:"Environment"`

	// 函数的私有网络配置
	VpcConfig *VpcConfig `json:"VpcConfig" name:"VpcConfig"`
}

func (r *UpdateFunctionConfigurationRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *UpdateFunctionConfigurationRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type UpdateFunctionConfigurationResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *UpdateFunctionConfigurationResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *UpdateFunctionConfigurationResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type Variable struct {

	// 变量的名称
	Key *string `json:"Key" name:"Key"`

	// 变量的值
	Value *string `json:"Value" name:"Value"`
}

type VpcConfig struct {

	// 私有网络 的 id
	VpcId *string `json:"VpcId" name:"VpcId"`

	// 子网的 id
	SubnetId *string `json:"SubnetId" name:"SubnetId"`
}
