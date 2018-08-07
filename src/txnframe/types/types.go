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

import (
	"time"
)

type Tansaction struct {
	TxnID      string     // 事务ID,uuid
	RequestID  string     // 请求ID,可选项
	Processor  string     // 处理进程号，结构为"IP:PORT-PID"用于识别事务session被存于那个TM多活实例
	Status     TxStatus   // 事务状态，作为定时补偿判断条件，这个字段需要加索引
	CreateTime *time.Time // 创建时间，作为定时补偿判断条件和统计信息存在，这个字段需要加索引
	LastTime   *time.Time // 修改时间，作为统计信息存在
}

// TxStatus describe
type TxStatus int

// TxStatus enumerations
const (
	TxStatusOnProgress TxStatus = iota + 1
	TxStatusCommited
	TxStatusAborted
)

type Document interface{}
