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

type GetInvokeTxRequest struct {
	*tchttp.BaseRequest

	// 模块名，固定字段：transaction
	Module *string `json:"Module" name:"Module"`

	// 操作名，固定地段：invoke
	Operation *string `json:"Operation" name:"Operation"`

	// 区块链网络ID，可在区块链网络详情或列表中获取
	ClusterId *string `json:"ClusterId" name:"ClusterId"`

	// 业务所属通道名称，可在通道详情或列表中获取
	ChannelName *string `json:"ChannelName" name:"ChannelName"`

	// 执行该查询交易的节点名称，可以在通道详情中获取该通道上的节点名称极其所属组织名称
	PeerName *string `json:"PeerName" name:"PeerName"`

	// 执行该查询交易的节点所属组织名称，可以在通道详情中获取该通道上的节点名称极其所属组织名称
	PeerGroup *string `json:"PeerGroup" name:"PeerGroup"`

	// 事务ID
	TxId *string `json:"TxId" name:"TxId"`
}

func (r *GetInvokeTxRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *GetInvokeTxRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type GetInvokeTxResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 状态码
		TxValidationCode *int64 `json:"TxValidationCode" name:"TxValidationCode"`

		// 消息
		TxValidationMsg *string `json:"TxValidationMsg" name:"TxValidationMsg"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *GetInvokeTxResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *GetInvokeTxResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type InvokeRequest struct {
	*tchttp.BaseRequest

	// 模块名，固定字段：transaction
	Module *string `json:"Module" name:"Module"`

	// 操作名，固定地段：invoke
	Operation *string `json:"Operation" name:"Operation"`

	// 区块链网络ID，可在区块链网络详情或列表中获取
	ClusterId *string `json:"ClusterId" name:"ClusterId"`

	// 业务所属智能合约名称，可在智能合约详情或列表中获取
	ChaincodeName *string `json:"ChaincodeName" name:"ChaincodeName"`

	// 业务所属通道名称，可在通道详情或列表中获取
	ChannelName *string `json:"ChannelName" name:"ChannelName"`

	// 对该笔交易进行背书的节点列表（包括节点名称和节点所属组织名称，详见数据结构一节），可以在通道详情中获取该通道上的节点名称极其所属组织名称
	Peers []*PeerSet `json:"Peers" name:"Peers" list`

	// 该笔交易需要调用的智能合约中的函数名称
	FuncName *string `json:"FuncName" name:"FuncName"`

	// 被调用的函数参数列表
	Args []*string `json:"Args" name:"Args" list`

	// 同步调用标识，可选参数，值为0或者不传表示使用同步方法调用，调用后会等待交易执行后再返回执行结果；值为1时表示使用异步方式调用Invoke，执行后会立即返回交易对应的Txid，后续需要通过GetInvokeTx这个API查询该交易的执行结果。（对于逻辑较为简单的交易，可以使用同步模式；对于逻辑较为复杂的交易，建议使用异步模式，否则容易导致API因等待时间过长，返回等待超时）
	AsyncFlag *uint64 `json:"AsyncFlag" name:"AsyncFlag"`
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

		// 交易编号
		Txid *string `json:"Txid" name:"Txid"`

		// 交易执行结果
		Events *string `json:"Events" name:"Events"`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
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

type PeerSet struct {

	// 节点名称
	PeerName *string `json:"PeerName" name:"PeerName"`

	// 组织名称
	OrgName *string `json:"OrgName" name:"OrgName"`
}

type QueryRequest struct {
	*tchttp.BaseRequest

	// 模块名，固定字段：transaction
	Module *string `json:"Module" name:"Module"`

	// 操作名，固定地段：query
	Operation *string `json:"Operation" name:"Operation"`

	// 区块链网络ID，可在区块链网络详情或列表中获取
	ClusterId *string `json:"ClusterId" name:"ClusterId"`

	// 业务所属智能合约名称，可在智能合约详情或列表中获取
	ChaincodeName *string `json:"ChaincodeName" name:"ChaincodeName"`

	// 业务所属通道名称，可在通道详情或列表中获取
	ChannelName *string `json:"ChannelName" name:"ChannelName"`

	// 执行该查询交易的节点列表（包括节点名称和节点所属组织名称，详见数据结构一节），可以在通道详情中获取该通道上的节点名称极其所属组织名称
	Peers []*PeerSet `json:"Peers" name:"Peers" list`

	// 该笔交易查询需要调用的智能合约中的函数名称
	FuncName *string `json:"FuncName" name:"FuncName"`

	// 被调用的函数参数列表
	Args []*string `json:"Args" name:"Args" list`
}

func (r *QueryRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *QueryRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type QueryResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 查询结果数据
		Data []*string `json:"Data" name:"Data" list`

		// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *QueryResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *QueryResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}
