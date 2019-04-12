/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package types

// MsgHeader message header
type MsgHeader struct {
	OPCode    OPCode
	TxnID     string
	RequestID string
}

// OPCode operation code type
type OPCode uint32

const (
	// OPInsertCode insert operation code
	OPInsertCode OPCode = iota + 1
	// OPUpdateCode update operation code
	OPUpdateCode
	// OPDeleteCode delete operation code
	OPDeleteCode
	// OPFindCode query operation code
	OPFindCode
	// OPFindAndModifyCode find and modify operation code
	OPFindAndModifyCode
	// OPCountCode count operation code
	OPCountCode
	// OPAggregateCode aggregate operation code
	OPAggregateCode
	// OPStartTransactionCode start a transaction code
	OPStartTransactionCode OPCode = 666
	// OPCommitCode transaction commit operation code
	OPCommitCode OPCode = 667
	// OPAbortCode transaction abort operation code
	OPAbortCode OPCode = 668
)

func (c OPCode) String() string {
	switch c {
	case OPInsertCode:
		return "OPInsert"
	case OPUpdateCode:
		return "OPUpdate"
	case OPDeleteCode:
		return "OPDelete"
	case OPFindCode:
		return "OPFind"
	case OPFindAndModifyCode:
		return "OPFindAndModify"
	case OPCountCode:
		return "OPCount"
	case OPStartTransactionCode:
		return "OPStartTransaction"
	case OPCommitCode:
		return "OPCommitTransaction"
	case OPAbortCode:
		return "OPAbortTransaction"
	case OPAggregateCode:
		return "OPAggregate"
	default:
		return "UNKNOW"
	}
}

// OPInsertOperation insert operation request structure
type OPInsertOperation struct {
	MsgHeader            // 标准报文头
	Collection string    // "dbname.collectionname"
	DOCS       Documents // 要插入集合的文档
}

// OPPipelineOperation insert operation request structure
type OPAggregateOperation struct {
	MsgHeader            // 标准报文头
	Collection string    // "dbname.collectionname"
	Pipiline   Documents // 要插入集合的文档
}

// OPUpdateOperation update operation request structure
type OPUpdateOperation struct {
	MsgHeader           // 标准报文头
	Collection string   // "dbname.collectionname"
	DOC        Document // 指定要执行的更新
	Selector   Document // 文档查询条件
}

// OPDeleteOperation delete operation request structure
type OPDeleteOperation struct {
	MsgHeader           // 标准报文头
	Collection string   // "dbname.collectionname"
	Selector   Document // 文档查询条件
}

// OPFindOperation find operation request structure
type OPFindOperation struct {
	MsgHeader           // 标准报文头
	Collection string   // "dbname.collectionname"
	Projection Document // ""
	Selector   Document // 文档查询条件
	Start      uint64   // start index
	Limit      uint64   // limit index
	Sort       string   // sort string
}

// OPCountOperation count operation request structure
type OPCountOperation struct {
	MsgHeader           // 标准报文头
	Collection string   // "dbname.collectionname"
	Selector   Document // 文档查询条件
}

// OPFindAndModifyOperation find and modify operation request structure
type OPFindAndModifyOperation struct {
	MsgHeader           // 标准报文头
	Collection string   // "dbname.collectionname"
	DOC        Document // 指定要执行的更新
	Selector   Document // 文档查询条件
	Upsert     bool
	Remove     bool
	ReturnNew  bool
}

// OPStartTransactionOperation transaction request structure
type OPStartTransactionOperation struct {
	MsgHeader
}

// OPCommitOperation commit operation request structure
type OPCommitOperation struct {
	MsgHeader
}

// OPAbortOperation abort operation request structure
type OPAbortOperation struct {
	MsgHeader
}

// ReplyHeader the rpc message header structure
type ReplyHeader struct {
	MsgHeader
	Processor string
	Success   bool
	Code      int
	Message   string
}

// OPReply the operation reply message header structure
type OPReply struct {
	ReplyHeader           // 标准报文头
	Count       uint64    // 文档查询结果数
	Docs        Documents // 文档查询结果
}
