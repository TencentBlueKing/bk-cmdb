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
	"net/http"
	"reflect"
	"strconv"
	"time"

	"encoding/json"

	"configcenter/src/common"
	"configcenter/src/common/blog"
)

type Transaction struct {
	TxnID      string    `bson:"bk_txn_id"`     // 事务ID,uuid
	RequestID  string    `bson:"bk_request_id"` // 请求ID,可选项
	Processor  string    `bson:"processor"`     // 处理进程号，结构为"IP:PORT-PID"用于识别事务session被存于那个TM多活实例
	Status     TxStatus  `bson:"status"`        // 事务状态，作为定时补偿判断条件，这个字段需要加索引
	CreateTime time.Time `bson:"create_time"`   // 创建时间，作为定时补偿判断条件和统计信息存在，这个字段需要加索引
	LastTime   time.Time `bson:"last_time"`     // 修改时间，作为统计信息存在
}

func (t Transaction) IntoHeader(header http.Header) http.Header {
	tar := http.Header{}
	for key := range header {
		tar.Set(key, header.Get(key))
	}
	tar.Set(common.BKHTTPCCTransactionID, t.TxnID)
	return tar
}

// TxStatus describe
type TxStatus int

// TxStatus enumerations
const (
	TxStatusOnProgress TxStatus = iota + 1
	TxStatusCommitted
	TxStatusAborted
	TxStatusException
)

type Document map[string]interface{}

func (d Document) Decode(result interface{}) error {
	out, err := json.Marshal(d)
	if nil != err {
		return err
	}
	return json.Unmarshal(out, result)
}

func (d *Document) Encode(result interface{}) error {
	if nil == result {
		return nil
	}
	out, err := json.Marshal(result)
	if nil != err {
		return err
	}
	return json.Unmarshal(out, d)
}

type Documents []Document

func (d Documents) Decode(result interface{}) error {
	resultv := reflect.ValueOf(result)
	switch resultv.Elem().Kind() {
	case reflect.Slice:
		out, err := json.Marshal(d)
		if nil != err {
			return err
		}
		err = json.Unmarshal(out, result)
		if nil != err {
			blog.Errorf("Decode Document error: %s, source is %#v", err.Error(), out)
		}
		return err
	default:
		if len(d) <= 0 {
			return nil
		}
		out, err := json.Marshal(d[0])
		if nil != err {
			return err
		}
		return json.Unmarshal(out, result)
	}
}

func (d *Documents) Encode(result interface{}) error {
	if nil == result {
		return nil
	}
	resultv := reflect.ValueOf(result)
	for resultv.CanAddr() {
		resultv = resultv.Elem()
	}
	switch resultv.Kind() {
	case reflect.Slice:
		out, err := json.Marshal(result)
		if nil != err {
			return err
		}
		*d = []Document{}
		blog.Infof("Encode slice %s", out)
		return json.Unmarshal(out, d)
	default:
		out, err := json.Marshal(result)
		if nil != err {
			return err
		}
		*d = []Document{Document{}}
		return json.Unmarshal(out, &(*d)[0])
	}
}

const (
	CommandRDBOperation              = "RDB"
	CommandWatchTransactionOperation = "WatchTransaction"
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
		out, err := json.Marshal(origin)
		if err != nil {
			return &Page{}
		}
		err = json.Unmarshal(out, &page)
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
