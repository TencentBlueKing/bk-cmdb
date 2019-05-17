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

package v20180709

import (
    "encoding/json"

    tchttp "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/http"
)

type BillDetail struct {

	// 产品名称
	BusinessCodeName *string `json:"BusinessCodeName" name:"BusinessCodeName"`

	// 子产品名称
	ProductCodeName *string `json:"ProductCodeName" name:"ProductCodeName"`

	// 计费模式
	PayModeName *string `json:"PayModeName" name:"PayModeName"`

	// 项目
	ProjectName *string `json:"ProjectName" name:"ProjectName"`

	// 区域
	RegionName *string `json:"RegionName" name:"RegionName"`

	// 可用区
	ZoneName *string `json:"ZoneName" name:"ZoneName"`

	// 资源实例ID
	ResourceId *string `json:"ResourceId" name:"ResourceId"`

	// 实例名称
	ResourceName *string `json:"ResourceName" name:"ResourceName"`

	// 交易类型
	ActionTypeName *string `json:"ActionTypeName" name:"ActionTypeName"`

	// 订单ID
	OrderId *string `json:"OrderId" name:"OrderId"`

	// 交易ID
	BillId *string `json:"BillId" name:"BillId"`

	// 扣费时间
	PayTime *string `json:"PayTime" name:"PayTime"`

	// 开始使用时间
	FeeBeginTime *string `json:"FeeBeginTime" name:"FeeBeginTime"`

	// 结束使用时间
	FeeEndTime *string `json:"FeeEndTime" name:"FeeEndTime"`

	// 组件列表
	ComponentSet []*BillDetailComponent `json:"ComponentSet" name:"ComponentSet" list`

	// 支付者UIN
	PayerUin *string `json:"PayerUin" name:"PayerUin"`

	// 使用者UIN
	OwnerUin *string `json:"OwnerUin" name:"OwnerUin"`

	// 操作者UIN
	OperateUin *string `json:"OperateUin" name:"OperateUin"`
}

type BillDetailComponent struct {

	// 组件名称
	ComponentCodeName *string `json:"ComponentCodeName" name:"ComponentCodeName"`

	// 组件类型名称
	ItemCodeName *string `json:"ItemCodeName" name:"ItemCodeName"`

	// 组件刊例价
	SinglePrice *string `json:"SinglePrice" name:"SinglePrice"`

	// 组件指定价
	SpecifiedPrice *string `json:"SpecifiedPrice" name:"SpecifiedPrice"`

	// 价格单位
	PriceUnit *string `json:"PriceUnit" name:"PriceUnit"`

	// 组件用量
	UsedAmount *string `json:"UsedAmount" name:"UsedAmount"`

	// 组件用量单位
	UsedAmountUnit *string `json:"UsedAmountUnit" name:"UsedAmountUnit"`

	// 使用时长
	TimeSpan *string `json:"TimeSpan" name:"TimeSpan"`

	// 时长单位
	TimeUnitName *string `json:"TimeUnitName" name:"TimeUnitName"`

	// 组件原价
	Cost *string `json:"Cost" name:"Cost"`

	// 折扣率
	Discount *string `json:"Discount" name:"Discount"`

	// 优惠类型
	ReduceType *string `json:"ReduceType" name:"ReduceType"`

	// 优惠后总价
	RealCost *string `json:"RealCost" name:"RealCost"`

	// 代金券支付金额
	VoucherPayAmount *string `json:"VoucherPayAmount" name:"VoucherPayAmount"`

	// 现金支付金额
	CashPayAmount *string `json:"CashPayAmount" name:"CashPayAmount"`

	// 赠送账户支付金额
	IncentivePayAmount *string `json:"IncentivePayAmount" name:"IncentivePayAmount"`
}

type BillResourceSummary struct {

	// 产品
	BusinessCodeName *string `json:"BusinessCodeName" name:"BusinessCodeName"`

	// 子产品
	ProductCodeName *string `json:"ProductCodeName" name:"ProductCodeName"`

	// 计费模式
	PayModeName *string `json:"PayModeName" name:"PayModeName"`

	// 项目
	ProjectName *string `json:"ProjectName" name:"ProjectName"`

	// 地域
	RegionName *string `json:"RegionName" name:"RegionName"`

	// 可用区
	ZoneName *string `json:"ZoneName" name:"ZoneName"`

	// 资源实例ID
	ResourceId *string `json:"ResourceId" name:"ResourceId"`

	// 资源实例名称
	ResourceName *string `json:"ResourceName" name:"ResourceName"`

	// 交易类型
	ActionTypeName *string `json:"ActionTypeName" name:"ActionTypeName"`

	// 订单ID
	OrderId *string `json:"OrderId" name:"OrderId"`

	// 扣费时间
	PayTime *string `json:"PayTime" name:"PayTime"`

	// 开始使用时间
	FeeBeginTime *string `json:"FeeBeginTime" name:"FeeBeginTime"`

	// 结束使用时间
	FeeEndTime *string `json:"FeeEndTime" name:"FeeEndTime"`

	// 配置描述
	ConfigDesc *string `json:"ConfigDesc" name:"ConfigDesc"`

	// 扩展字段1
	ExtendField1 *string `json:"ExtendField1" name:"ExtendField1"`

	// 扩展字段2
	ExtendField2 *string `json:"ExtendField2" name:"ExtendField2"`

	// 原价，单位为元
	TotalCost *string `json:"TotalCost" name:"TotalCost"`

	// 折扣率
	Discount *string `json:"Discount" name:"Discount"`

	// 优惠类型
	ReduceType *string `json:"ReduceType" name:"ReduceType"`

	// 优惠后总价，单位为元
	RealTotalCost *string `json:"RealTotalCost" name:"RealTotalCost"`

	// 代金券支付金额，单位为元
	VoucherPayAmount *string `json:"VoucherPayAmount" name:"VoucherPayAmount"`

	// 现金账户支付金额，单位为元
	CashPayAmount *string `json:"CashPayAmount" name:"CashPayAmount"`

	// 赠送账户支付金额，单位为元
	IncentivePayAmount *string `json:"IncentivePayAmount" name:"IncentivePayAmount"`
}

type Deal struct {

	// 订单号
	OrderId *string `json:"OrderId" name:"OrderId"`

	// 订单状态
	Status *int64 `json:"Status" name:"Status"`

	// 支付者
	Payer *string `json:"Payer" name:"Payer"`

	// 创建时间
	CreateTime *string `json:"CreateTime" name:"CreateTime"`

	// 创建人
	Creator *string `json:"Creator" name:"Creator"`

	// 实际支付金额（分）
	RealTotalCost *int64 `json:"RealTotalCost" name:"RealTotalCost"`

	// 代金券抵扣金额（分）
	VoucherDecline *int64 `json:"VoucherDecline" name:"VoucherDecline"`

	// 项目ID
	ProjectId *int64 `json:"ProjectId" name:"ProjectId"`

	// 产品分类ID
	GoodsCategoryId *int64 `json:"GoodsCategoryId" name:"GoodsCategoryId"`

	// 产品详情
	ProductInfo []*ProductInfo `json:"ProductInfo" name:"ProductInfo" list`

	// 时长
	TimeSpan *float64 `json:"TimeSpan" name:"TimeSpan"`

	// 时间单位
	TimeUnit *string `json:"TimeUnit" name:"TimeUnit"`

	// 货币单位
	Currency *string `json:"Currency" name:"Currency"`

	// 折扣率
	Policy *float64 `json:"Policy" name:"Policy"`

	// 单价（分）
	Price *float64 `json:"Price" name:"Price"`

	// 原价（分）
	TotalCost *float64 `json:"TotalCost" name:"TotalCost"`
}

type DescribeAccountBalanceRequest struct {
	*tchttp.BaseRequest
}

func (r *DescribeAccountBalanceRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeAccountBalanceRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeAccountBalanceResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 云账户信息中的”展示可用余额”字段，单位为"分"
		Balance *int64 `json:"Balance" name:"Balance"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeAccountBalanceResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeAccountBalanceResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeBillDetailRequest struct {
	*tchttp.BaseRequest

	// 偏移量
	Offset *uint64 `json:"Offset" name:"Offset"`

	// 数量，最大值为100
	Limit *uint64 `json:"Limit" name:"Limit"`

	// 周期类型，byPayTime按扣费周期/byUsedTime按计费周期
	PeriodType *string `json:"PeriodType" name:"PeriodType"`

	// 月份，格式为yyyy-mm
	Month *string `json:"Month" name:"Month"`

	// 周期开始时间，格式为Y-m-d H:i:s，如果有该字段则Month字段无效。BeginTime和EndTime必须一起传
	BeginTime *string `json:"BeginTime" name:"BeginTime"`

	// 周期结束时间，格式为Y-m-d H:i:s，如果有该字段则Month字段无效。BeginTime和EndTime必须一起传
	EndTime *string `json:"EndTime" name:"EndTime"`
}

func (r *DescribeBillDetailRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeBillDetailRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeBillDetailResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 详情列表
		DetailSet []*BillDetail `json:"DetailSet" name:"DetailSet" list`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeBillDetailResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeBillDetailResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeBillResourceSummaryRequest struct {
	*tchttp.BaseRequest

	// 偏移量
	Offset *uint64 `json:"Offset" name:"Offset"`

	// 数量，最大值为1000
	Limit *uint64 `json:"Limit" name:"Limit"`

	// 周期类型，byUsedTime按计费周期/byPayTime按扣费周期
	PeriodType *string `json:"PeriodType" name:"PeriodType"`

	// 月份，格式为yyyy-mm
	Month *string `json:"Month" name:"Month"`
}

func (r *DescribeBillResourceSummaryRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeBillResourceSummaryRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeBillResourceSummaryResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 资源汇总列表
		ResourceSummarySet []*BillResourceSummary `json:"ResourceSummarySet" name:"ResourceSummarySet" list`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeBillResourceSummaryResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeBillResourceSummaryResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeDealsByCondRequest struct {
	*tchttp.BaseRequest

	// 开始时间
	StartTime *string `json:"StartTime" name:"StartTime"`

	// 结束时间
	EndTime *string `json:"EndTime" name:"EndTime"`

	// 一页多少条数据，默认是20条，最大不超过1000
	Limit *int64 `json:"Limit" name:"Limit"`

	// 第多少页，从0开始，默认是0
	Offset *int64 `json:"Offset" name:"Offset"`

	// 订单状态,默认为4（成功的订单）
	// 订单的状态
	// 1：未支付
	// 2：已支付3：发货中
	// 4：已发货
	// 5：发货失败
	// 6：已退款
	// 7：已关单
	// 8：订单过期
	// 9：订单已失效
	// 10：产品已失效
	// 11：代付拒绝
	// 12：支付中
	Status *int64 `json:"Status" name:"Status"`

	// 订单号
	OrderId *string `json:"OrderId" name:"OrderId"`
}

func (r *DescribeDealsByCondRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeDealsByCondRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeDealsByCondResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 订单列表
		Deals []*Deal `json:"Deals" name:"Deals" list`

		// 订单总数
		TotalCount *int64 `json:"TotalCount" name:"TotalCount"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeDealsByCondResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeDealsByCondResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type PayDealsRequest struct {
	*tchttp.BaseRequest

	// 需要支付的一个或者多个订单号
	OrderIds []*string `json:"OrderIds" name:"OrderIds" list`

	// 是否自动使用代金券,1:是,0否,默认0
	AutoVoucher *int64 `json:"AutoVoucher" name:"AutoVoucher"`

	// 代金券ID列表,目前仅支持指定一张代金券
	VoucherIds []*string `json:"VoucherIds" name:"VoucherIds" list`
}

func (r *PayDealsRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *PayDealsRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type PayDealsResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 此次操作支付成功的订单号数组
		OrderIds []*string `json:"OrderIds" name:"OrderIds" list`

		// 此次操作支付成功的资源Id数组
		ResourceIds []*string `json:"ResourceIds" name:"ResourceIds" list`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *PayDealsResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *PayDealsResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ProductInfo struct {

	// 商品详情名称标识
	Name *string `json:"Name" name:"Name"`

	// 商品详情
	Value *string `json:"Value" name:"Value"`
}
