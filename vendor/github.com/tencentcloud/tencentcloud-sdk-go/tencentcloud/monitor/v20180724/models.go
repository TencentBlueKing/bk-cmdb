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

package v20180724

import (
    "encoding/json"

    tchttp "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/http"
)

type DataPoint struct {

	// 实例对象维度组合
	Dimensions []*Dimension `json:"Dimensions" name:"Dimensions" list`

	// 时间戳数组，表示那些时间点有数据，缺失的时间戳，没有数据点，可以理解为掉点了
	Timestamps []*float64 `json:"Timestamps" name:"Timestamps" list`

	// 监控值数组，该数组和Timestamps一一对应
	Values []*float64 `json:"Values" name:"Values" list`
}

type DescribeBaseMetricsRequest struct {
	*tchttp.BaseRequest

	// 业务命名空间
	Namespace *string `json:"Namespace" name:"Namespace"`

	// 指标名
	MetricName *string `json:"MetricName" name:"MetricName"`
}

func (r *DescribeBaseMetricsRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeBaseMetricsRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeBaseMetricsResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 查询得到的指标描述列表
		MetricSet []*MetricSet `json:"MetricSet" name:"MetricSet" list`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeBaseMetricsResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeBaseMetricsResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type Dimension struct {

	// 实例维度名称
	Name *string `json:"Name" name:"Name"`

	// 实例维度值
	Value *string `json:"Value" name:"Value"`
}

type DimensionsDesc struct {

	// 维度名数组
	Dimensions []*string `json:"Dimensions" name:"Dimensions" list`
}

type GetMonitorDataRequest struct {
	*tchttp.BaseRequest

	// 命名空间，每个云产品会有一个命名空间
	Namespace *string `json:"Namespace" name:"Namespace"`

	// 指标名称
	MetricName *string `json:"MetricName" name:"MetricName"`

	// 实例对象的维度组合
	Instances []*Instance `json:"Instances" name:"Instances" list`

	// 监控统计周期。默认为取值为300，单位为s
	Period *uint64 `json:"Period" name:"Period"`

	// 起始时间，如2018-09-22T19:51:23+08:00
	StartTime *string `json:"StartTime" name:"StartTime"`

	// 结束时间，默认为当前时间。 EndTime不能小于EtartTime
	EndTime *string `json:"EndTime" name:"EndTime"`
}

func (r *GetMonitorDataRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *GetMonitorDataRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type GetMonitorDataResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 统计周期
		Period *uint64 `json:"Period" name:"Period"`

		// 指标名
		MetricName *string `json:"MetricName" name:"MetricName"`

		// 数据点数组
		DataPoints []*DataPoint `json:"DataPoints" name:"DataPoints" list`

		// 开始时间
		StartTime *string `json:"StartTime" name:"StartTime"`

		// 结束时间
		EndTime *string `json:"EndTime" name:"EndTime"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *GetMonitorDataResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *GetMonitorDataResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type Instance struct {

	// 实例的维度组合
	Dimensions []*Dimension `json:"Dimensions" name:"Dimensions" list`
}

type MetricObjectMeaning struct {

	// 指标英文解释
	En *string `json:"En" name:"En"`

	// 指标中文解释
	Zh *string `json:"Zh" name:"Zh"`
}

type MetricSet struct {

	// 命名空间，每个云产品会有一个命名空间
	Namespace *string `json:"Namespace" name:"Namespace"`

	// 指标名称
	MetricName *string `json:"MetricName" name:"MetricName"`

	// 指标使用的单位
	Unit *string `json:"Unit" name:"Unit"`

	// 指标使用的单位
	UnitCname *string `json:"UnitCname" name:"UnitCname"`

	// 指标支持的统计周期，单位是秒，如60、300
	Period []*int64 `json:"Period" name:"Period" list`

	// 统计周期内指标方式
	Periods []*PeriodsSt `json:"Periods" name:"Periods" list`

	// 统计指标含义解释
	Meaning *MetricObjectMeaning `json:"Meaning" name:"Meaning"`

	// 维度描述信息
	Dimensions []*DimensionsDesc `json:"Dimensions" name:"Dimensions" list`
}

type PeriodsSt struct {

	// 周期
	Period *string `json:"Period" name:"Period"`

	// 统计方式
	StatType []*string `json:"StatType" name:"StatType" list`
}
