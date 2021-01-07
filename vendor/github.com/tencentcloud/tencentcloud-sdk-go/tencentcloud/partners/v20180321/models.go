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

type AgentAuditedClient struct {

	// 代理商账号ID
	Uin *string `json:"Uin" name:"Uin"`

	// 代客账号ID
	ClientUin *string `json:"ClientUin" name:"ClientUin"`

	// 代客审核通过时间戳
	AgentTime *string `json:"AgentTime" name:"AgentTime"`

	// 代客类型，可能值为a/b/c
	ClientFlag *string `json:"ClientFlag" name:"ClientFlag"`

	// 代客备注
	ClientRemark *string `json:"ClientRemark" name:"ClientRemark"`

	// 代客名称（首选实名认证名称）
	ClientName *string `json:"ClientName" name:"ClientName"`

	// 认证类型, 0：个人，1：企业；其他：未认证
	AuthType *string `json:"AuthType" name:"AuthType"`

	// 代客APPID
	AppId *string `json:"AppId" name:"AppId"`

	// 上月消费金额
	LastMonthAmt *uint64 `json:"LastMonthAmt" name:"LastMonthAmt"`

	// 本月消费金额
	ThisMonthAmt *uint64 `json:"ThisMonthAmt" name:"ThisMonthAmt"`

	// 是否欠费,0：不欠费；1：欠费
	HasOverdueBill *uint64 `json:"HasOverdueBill" name:"HasOverdueBill"`
}

type AgentBillElem struct {

	// 代理商账号ID
	Uin *string `json:"Uin" name:"Uin"`

	// 订单号，仅对预付费账单有意义
	OrderId *string `json:"OrderId" name:"OrderId"`

	// 代客账号ID
	ClientUin *string `json:"ClientUin" name:"ClientUin"`

	// 代客备注名称
	ClientRemark *string `json:"ClientRemark" name:"ClientRemark"`

	// 支付时间
	PayTime *string `json:"PayTime" name:"PayTime"`

	// 云产品名称
	GoodsType *string `json:"GoodsType" name:"GoodsType"`

	// 预付费/后付费
	PayMode *string `json:"PayMode" name:"PayMode"`

	// 支付月份
	SettleMonth *string `json:"SettleMonth" name:"SettleMonth"`

	// 支付金额，单位分
	Amt *uint64 `json:"Amt" name:"Amt"`

	// agentpay：代付；selfpay：自付
	PayerMode *string `json:"PayerMode" name:"PayerMode"`
}

type AgentClientElem struct {

	// 代理商账号ID
	Uin *string `json:"Uin" name:"Uin"`

	// 代客账号ID
	ClientUin *string `json:"ClientUin" name:"ClientUin"`

	// 代客申请时间戳
	ApplyTime *uint64 `json:"ApplyTime" name:"ApplyTime"`

	// 代客类型，可能值为a/b/c
	ClientFlag *string `json:"ClientFlag" name:"ClientFlag"`

	// 代客邮箱，打码显示
	Mail *string `json:"Mail" name:"Mail"`

	// 代客手机，打码显示
	Phone *string `json:"Phone" name:"Phone"`

	// 0表示不欠费，1表示欠费
	HasOverdueBill *uint64 `json:"HasOverdueBill" name:"HasOverdueBill"`
}

type AgentPayDealsRequest struct {
	*tchttp.BaseRequest

	// 订单所有者uin
	OwnerUin *string `json:"OwnerUin" name:"OwnerUin"`

	// 代付标志，1：代付；0：自付
	AgentPay *uint64 `json:"AgentPay" name:"AgentPay"`

	// 订单号数组
	DealNames []*string `json:"DealNames" name:"DealNames" list`
}

func (r *AgentPayDealsRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *AgentPayDealsRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type AgentPayDealsResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *AgentPayDealsResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *AgentPayDealsResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type AgentTransferMoneyRequest struct {
	*tchttp.BaseRequest

	// 客户账号ID
	ClientUin *string `json:"ClientUin" name:"ClientUin"`

	// 转账金额，单位分
	Amount *uint64 `json:"Amount" name:"Amount"`
}

func (r *AgentTransferMoneyRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *AgentTransferMoneyRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type AgentTransferMoneyResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *AgentTransferMoneyResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *AgentTransferMoneyResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type AuditApplyClientRequest struct {
	*tchttp.BaseRequest

	// 待审核客户账号ID
	ClientUin *string `json:"ClientUin" name:"ClientUin"`

	// 审核结果，可能的取值：accept/reject
	AuditResult *string `json:"AuditResult" name:"AuditResult"`

	// 申请理由，B类客户审核通过时必须填写申请理由
	Note *string `json:"Note" name:"Note"`
}

func (r *AuditApplyClientRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *AuditApplyClientRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type AuditApplyClientResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 代理商账号ID
		Uin *string `json:"Uin" name:"Uin"`

		// 客户账号ID
		ClientUin *string `json:"ClientUin" name:"ClientUin"`

		// 审核结果，包括accept/reject/qcloudaudit（腾讯云审核）
		AuditResult *string `json:"AuditResult" name:"AuditResult"`

		// 关联时间对应的时间戳
		AgentTime *uint64 `json:"AgentTime" name:"AgentTime"`

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *AuditApplyClientResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *AuditApplyClientResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeAgentAuditedClientsRequest struct {
	*tchttp.BaseRequest

	// 客户账号ID
	ClientUin *string `json:"ClientUin" name:"ClientUin"`

	// 客户名称。由于涉及隐私，名称打码显示，故名称仅支持打码后的模糊搜索
	ClientName *string `json:"ClientName" name:"ClientName"`

	// 客户类型，a/b，类型定义参考代理商相关政策文档
	ClientFlag *string `json:"ClientFlag" name:"ClientFlag"`

	// ASC/DESC， 不区分大小写，按审核通过时间排序
	OrderDirection *string `json:"OrderDirection" name:"OrderDirection"`

	// 客户账号ID列表
	ClientUins []*string `json:"ClientUins" name:"ClientUins" list`

	// 是否欠费。0：不欠费；1：欠费
	HasOverdueBill *uint64 `json:"HasOverdueBill" name:"HasOverdueBill"`

	// 客户备注
	ClientRemark *string `json:"ClientRemark" name:"ClientRemark"`

	// 偏移量
	Offset *uint64 `json:"Offset" name:"Offset"`

	// 限制数目
	Limit *uint64 `json:"Limit" name:"Limit"`
}

func (r *DescribeAgentAuditedClientsRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeAgentAuditedClientsRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeAgentAuditedClientsResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 已审核代客列表
		AgentClientSet []*AgentAuditedClient `json:"AgentClientSet" name:"AgentClientSet" list`

		// 符合条件的代客总数
		TotalCount *uint64 `json:"TotalCount" name:"TotalCount"`

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeAgentAuditedClientsResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeAgentAuditedClientsResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeAgentBillsRequest struct {
	*tchttp.BaseRequest

	// 支付月份，如2018-02
	SettleMonth *string `json:"SettleMonth" name:"SettleMonth"`

	// 客户账号ID
	ClientUin *string `json:"ClientUin" name:"ClientUin"`

	// 支付方式，prepay/postpay
	PayMode *string `json:"PayMode" name:"PayMode"`

	// 预付费订单号
	OrderId *string `json:"OrderId" name:"OrderId"`

	// 客户备注名称
	ClientRemark *string `json:"ClientRemark" name:"ClientRemark"`

	// 偏移量
	Offset *uint64 `json:"Offset" name:"Offset"`

	// 限制数目
	Limit *uint64 `json:"Limit" name:"Limit"`
}

func (r *DescribeAgentBillsRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeAgentBillsRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeAgentBillsResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 符合查询条件列表总数量
		TotalCount *uint64 `json:"TotalCount" name:"TotalCount"`

		// 业务明细列表
		AgentBillSet []*AgentBillElem `json:"AgentBillSet" name:"AgentBillSet" list`

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeAgentBillsResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeAgentBillsResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeAgentClientsRequest struct {
	*tchttp.BaseRequest

	// 客户账号ID
	ClientUin *string `json:"ClientUin" name:"ClientUin"`

	// 客户名称。由于涉及隐私，名称打码显示，故名称仅支持打码后的模糊搜索
	ClientName *string `json:"ClientName" name:"ClientName"`

	// 客户类型，a/b，类型定义参考代理商相关政策文档
	ClientFlag *string `json:"ClientFlag" name:"ClientFlag"`

	// ASC/DESC， 不区分大小写，按申请时间排序
	OrderDirection *string `json:"OrderDirection" name:"OrderDirection"`

	// 偏移量
	Offset *uint64 `json:"Offset" name:"Offset"`

	// 限制数目
	Limit *uint64 `json:"Limit" name:"Limit"`
}

func (r *DescribeAgentClientsRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeAgentClientsRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeAgentClientsResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 待审核代客列表
		AgentClientSet []*AgentClientElem `json:"AgentClientSet" name:"AgentClientSet" list`

		// 符合条件的代客总数
		TotalCount *uint64 `json:"TotalCount" name:"TotalCount"`

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeAgentClientsResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeAgentClientsResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeClientBalanceRequest struct {
	*tchttp.BaseRequest

	// 客户(代客)账号ID
	ClientUin *string `json:"ClientUin" name:"ClientUin"`
}

func (r *DescribeClientBalanceRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeClientBalanceRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeClientBalanceResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 账户余额，单位分
		Balance *uint64 `json:"Balance" name:"Balance"`

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeClientBalanceResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeClientBalanceResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeRebateInfosRequest struct {
	*tchttp.BaseRequest

	// 返佣月份，如2018-02
	RebateMonth *string `json:"RebateMonth" name:"RebateMonth"`

	// 偏移量
	Offset *uint64 `json:"Offset" name:"Offset"`

	// 限制数目
	Limit *uint64 `json:"Limit" name:"Limit"`
}

func (r *DescribeRebateInfosRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeRebateInfosRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeRebateInfosResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 返佣信息列表
		RebateInfoSet []*RebateInfoElem `json:"RebateInfoSet" name:"RebateInfoSet" list`

		// 符合查询条件返佣信息数目
		TotalCount *uint64 `json:"TotalCount" name:"TotalCount"`

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeRebateInfosResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeRebateInfosResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ModifyClientRemarkRequest struct {
	*tchttp.BaseRequest

	// 客户备注名称
	ClientRemark *string `json:"ClientRemark" name:"ClientRemark"`

	// 客户账号ID
	ClientUin *string `json:"ClientUin" name:"ClientUin"`
}

func (r *ModifyClientRemarkRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ModifyClientRemarkRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type ModifyClientRemarkResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *ModifyClientRemarkResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *ModifyClientRemarkResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type RebateInfoElem struct {

	// 代理商账号ID
	Uin *string `json:"Uin" name:"Uin"`

	// 返佣月份，如2018-02
	RebateMonth *string `json:"RebateMonth" name:"RebateMonth"`

	// 返佣金额，单位分
	Amt *uint64 `json:"Amt" name:"Amt"`

	// 月度业绩，单位分
	MonthSales *uint64 `json:"MonthSales" name:"MonthSales"`

	// 季度业绩，单位分
	QuarterSales *uint64 `json:"QuarterSales" name:"QuarterSales"`

	// NORMAL(正常)/HAS_OVERDUE_BILL(欠费)/NO_CONTRACT(缺合同)
	ExceptionFlag *string `json:"ExceptionFlag" name:"ExceptionFlag"`
}
