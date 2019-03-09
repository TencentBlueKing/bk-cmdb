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
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"time"

	"configcenter/src/common"

	"github.com/mongodb/mongo-go-driver/bson"
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

func (d *Document) Encode(value interface{}) error {
	if nil == value {
		return nil
	}
	out, err := bson.Marshal(value)
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
		err := decodeBsonArray(d, result)
		if err != nil {
			return fmt.Errorf("Decode Document array error: %v", err)
		}
		return nil
	default:
		if len(d) <= 0 {
			return nil
		}
		out, err := bson.Marshal(d[0])
		if nil != err {
			return fmt.Errorf("Decode Documents error when marshal2: %v, source is %s", err, d[0])
		}
		err = bson.Unmarshal(out, result)
		if nil != err {
			return fmt.Errorf("Decode Documents error when unmarshal: %v, source is %s", err, out)
		}
		return nil
	}
}

func (d *Documents) Encode(value interface{}) error {
	if nil == value {
		return nil
	}
	valuev := reflect.ValueOf(value)
	for valuev.CanAddr() {
		valuev = valuev.Elem()
	}
	switch valuev.Kind() {
	case reflect.Slice:
		err := decodeBsonArray(value, d)
		if err != nil {
			return fmt.Errorf("Encode Document array error: %v", err)
		}
		return nil
	default:
		out, err := bson.Marshal(value)
		if nil != err {
			return fmt.Errorf("Encode Documents when marshal error: %v, source is %#v", err, value)
		}
		*d = []Document{Document{}}
		err = bson.Unmarshal(out, &(*d)[0])
		if err != nil {
			return fmt.Errorf("Encode Documents when unmarshal error: %v, source is %v", err, bson.Raw(out))
		}
	}
	return nil
}

func decodeBsonArray(inArr, outArr interface{}) error {
	in := struct{ Data interface{} }{Data: inArr}
	bsonraw, err := bson.Marshal(in)
	if err != nil {
		return fmt.Errorf("[decodeBsonArray] marshal error: %v, source: %#v", err, in)
	}

	out := struct{ Data []bson.Raw }{}
	err = bson.Unmarshal(bsonraw, &out)
	if err != nil {
		return fmt.Errorf("[decodeBsonArray] unmarshal error: %v, source: %v", err, bson.Raw(bsonraw))
	}

	resultv := reflect.ValueOf(outArr)
	if resultv.Kind() != reflect.Ptr || resultv.Elem().Kind() != reflect.Slice {
		return errors.New("result argument must be a slice address")
	}
	slicev := resultv.Elem()
	slicev = slicev.Slice(0, slicev.Cap())
	elemt := slicev.Type().Elem()
	idx := 0
	for _, dataItem := range out.Data {
		if slicev.Len() == idx {
			elemp := reflect.New(elemt)
			if err := bson.Unmarshal(dataItem, elemp.Interface()); nil != err {
				return fmt.Errorf("[decodeBsonArray] unmarshal item error: %v, source: %v", err, bson.Raw(dataItem))
			}
			slicev = reflect.Append(slicev, elemp.Elem())
			slicev = slicev.Slice(0, slicev.Cap())
			idx++
			continue
		}

		if err := bson.Unmarshal(dataItem, slicev.Index(idx).Addr().Interface()); nil != err {
			return fmt.Errorf("[decodeBsonArray] unmarshal element error: %v, source: %v", err, bson.Raw(dataItem))
		}
		idx++
	}
	resultv.Elem().Set(slicev.Slice(0, idx))

	return nil
}

const (
	CommandRDBOperation              = "RDB"
	CommandWatchTransactionOperation = "WatchTransaction"
)

type Page struct {
	Limit uint64 `json:"limit,omitempty" bson:"limit,omitempty"`
	Start uint64 `json:"start,omitempty" bson:"start,omitempty"`
	Sort  string `json:"sort,omitempty" bson:"sort,omitempty"`
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
