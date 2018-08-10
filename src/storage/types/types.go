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
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/mongodb/mongo-go-driver/bson"

	"configcenter/src/common/blog"
)

type Tansaction struct {
	TxnID      string     `bson:"txn_id"`      // 事务ID,uuid
	RequestID  string     `bson:"request_id"`  // 请求ID,可选项
	Processor  string     `bson:"processor"`   // 处理进程号，结构为"IP:PORT-PID"用于识别事务session被存于那个TM多活实例
	Status     TxStatus   `bson:"status"`      // 事务状态，作为定时补偿判断条件，这个字段需要加索引
	CreateTime *time.Time `bson:"create_time"` // 创建时间，作为定时补偿判断条件和统计信息存在，这个字段需要加索引
	LastTime   *time.Time `bson:"last_time"`   // 修改时间，作为统计信息存在
}

// TxStatus describe
type TxStatus int

// TxStatus enumerations
const (
	TxStatusOnProgress TxStatus = iota + 1
	TxStatusCommited
	TxStatusAborted
	TxStatusException
)

type Document map[string]interface{}

func (d Document) Decode(result interface{}) error {
	out, err := bson.Marshal(d)
	if nil != err {
		return err
	}
	return bson.Unmarshal(out, result)
}

func (d Document) Encode(result interface{}) error {
	if nil == result {
		return nil
	}
	out, err := bson.Marshal(result)
	if nil != err {
		return err
	}
	return bson.Unmarshal(out, d)
}

type Documents []Document

func (d Documents) Decode(result interface{}) error {
	resultv := reflect.ValueOf(result)
	switch resultv.Elem().Kind() {
	case reflect.Slice:
		out, err := bson.Marshal(d)
		if nil != err {
			return err
		}
		err = bson.Unmarshal(out, result)
		if nil != err {
			blog.Errorf("Decode Document error: %s, source is %#v", err.Error(), out)
		}
		return err
	default:
		if len(d) <= 0 {
			return nil
		}
		out, err := bson.Marshal(d[0])
		if nil != err {
			return err
		}
		return bson.Unmarshal(out, result)
	}
}

func (d Documents) Encode(result interface{}) error {
	if nil == result {
		return nil
	}
	resultv := reflect.ValueOf(result)
	switch resultv.Elem().Kind() {
	case reflect.Slice:
		out, err := bson.Marshal(result)
		if nil != err {
			return err
		}
		return bson.Unmarshal(out, d)
	default:
		out, err := bson.Marshal(result)
		if nil != err {
			return err
		}
		d = make(Documents, 1)
		return bson.Unmarshal(out, d[0])
	}
}

const (
	CommandRDBOperation = "RDB"
)

type Page struct {
	Limit uint64 `json:"limit,omitempty"`
	Start uint64 `json:"start,omitempty"`
	Sort  string `json:"sort,omitempty"`
}

func ParsePage(origin interface{}) *Page {
	if origin == nil {
		return &Page{}
	}
	page, ok := origin.(map[string]interface{})
	if !ok {
		out, err := bson.Marshal(origin)
		if err != nil {
			return &Page{}
		}
		err = bson.Unmarshal(out, &page)
		if err != nil {
			return &Page{}
		}
	}
	result := Page{}
	if sort, ok := page["sort"].(string); ok {
		result.Sort = sort
	}
	if start, ok := page["start"]; ok {
		result.Start, _ = strconv.ParseUint(fmt.Sprint(start), 10, 64)
	}
	if limit, ok := page["limit"]; ok {
		result.Limit, _ = strconv.ParseUint(fmt.Sprint(limit), 10, 64)
	}
	return &result
}

type GetServerFunc func() ([]string, error)
