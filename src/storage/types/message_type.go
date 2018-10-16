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

type MsgHeader struct {
	TxnID     string
	RequestID string
	OPCode    OPCode
}

type OPCode uint32

const (
	OPInsert OPCode = iota + 1
	OPUpdate
	OPDelete
	OPFind
	OPFindAndModify
	OPCount
	OPStartTransaction OPCode = 666
	OPCommit           OPCode = 667
	OPAbort            OPCode = 668
)

func (c OPCode) String() string {
	switch c {
	case OPInsert:
		return "OPInsert"
	case OPUpdate:
		return "OPUpdate"
	case OPDelete:
		return "OPDelete"
	case OPFind:
		return "OPFind"
	case OPFindAndModify:
		return "OPFindAndModify"
	case OPCount:
		return "OPCount"
	case OPStartTransaction:
		return "OPStartTransaction"
	case OPCommit:
		return "OPCommit"
	case OPAbort:
		return "OPAbort"
	default:
		return "UNKNOW"
	}
}

type OPINSERT struct {
	MsgHeader            // 标准报文头
	Collection string    // "dbname.collectionname"
	DOCS       Documents // 要插入集合的文档
}

type OPUPDATE struct {
	MsgHeader           // 标准报文头
	Collection string   // "dbname.collectionname"
	DOC        Document // 指定要执行的更新
	Selector   Document // 文档查询条件
}

type OPDELETE struct {
	MsgHeader           // 标准报文头
	Collection string   // "dbname.collectionname"
	Selector   Document // 文档查询条件
}

type OPFIND struct {
	MsgHeader           // 标准报文头
	Collection string   // "dbname.collectionname"
	Projection Document // ""
	Selector   Document // 文档查询条件
	Start      uint64   // start index
	Limit      uint64   // limit index
	Sort       string   // sort string
}

type OPCOUNT struct {
	MsgHeader           // 标准报文头
	Collection string   // "dbname.collectionname"
	Selector   Document // 文档查询条件
}

type OPFINDANDMODIFY struct {
	MsgHeader           // 标准报文头
	Collection string   // "dbname.collectionname"
	DOC        Document // 指定要执行的更新
	Selector   Document // 文档查询条件
	Upsert     bool
	Remove     bool
	ReturnNew  bool
}

type OPSTARTTTRANSATION struct {
	MsgHeader
}

type OPCOMMIT struct {
	MsgHeader
}
type OPABORT struct {
	MsgHeader
}

type ReplyHeader struct {
	MsgHeader
	Processor string
	Success   bool
	Code      int
	Message   string
}

type OPREPLY struct {
	ReplyHeader           // 标准报文头
	Count       uint64    // 文档查询结果数
	Docs        Documents // 文档查询结果
}
