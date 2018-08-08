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

package mongobyc

// #include <stdlib.h>
// #include "mongo.h"
import "C"

import (
	"encoding/json"
	"fmt"
	"reflect"
	"unsafe"

	"configcenter/src/common"
	"configcenter/src/common/mapstr"
	"configcenter/src/txnframe/mongobyc/findopt"
)

// TransformFindOpts transform find options
func TransformFindOpts(opts *findopt.Opts) *findopt.Opts {

	if 0 >= opts.Limit {
		opts.Limit = common.BKNoLimit
	}
	if 0 >= opts.Skip {
		opts.Skip = 0
	}
	if 0 == len(opts.Sort) {
		opts.Sort = "_id"
	}

	if !opts.Fields.Exists("_id") {
		opts.Fields.Set("_id", 0)
	}
	return opts
}

// TransformError transform bson error into golang error instance
func TransformError(err C.bson_error_t) error {
	var strErr []byte
	for idx := range err.message {
		if 0 == err.message[idx] {
			break
		}
		strErr = append(strErr, byte(err.message[idx]))
	}
	return fmt.Errorf("%s", string(strErr))
}

// TransformDocument transform document into bson
func TransformDocument(doc interface{}) (*C.bson_t, error) {

	switch docType := doc.(type) {
	case string:
		bdata := &bson{
			data: docType,
		}
		return bdata.ToDocument()

	default:
		data, err := json.Marshal(docType)
		if nil != err {
			return nil, err
		}
		bdata := &bson{
			data: string(data),
		}
		return bdata.ToDocument()
	}

}

// TransformBsonIntoGoString transform bson into go string
func TransformBsonIntoGoString(reply *C.bson_t) string {
	str := C.bson_as_json(reply, nil)
	defer C.bson_free(unsafe.Pointer(str))
	return C.GoString(str)
}

// TransformMapStrIntoResult transform data into result
func TransformMapStrIntoResult(datas []mapstr.MapStr, result interface{}) {

	resultv := reflect.ValueOf(result)
	if resultv.Kind() != reflect.Ptr || resultv.Elem().Kind() != reflect.Slice {
		panic("result argument must be a slice address")
	}
	slicev := resultv.Elem()
	slicev = slicev.Slice(0, slicev.Cap())
	elemt := slicev.Type().Elem()
	idx := 0
	for _, dataItem := range datas {
		if slicev.Len() == idx {
			elemp := reflect.New(elemt)
			if err := dataItem.MarshalJSONInto(elemp.Interface()); nil != err {
				panic(err)
			}
			slicev = reflect.Append(slicev, elemp.Elem())
			slicev = slicev.Slice(0, slicev.Cap())

			continue
		}

		if err := dataItem.MarshalJSONInto(slicev.Index(idx).Addr().Interface()); nil != err {
			panic(err)
		}
		idx++
	}
	resultv.Elem().Set(slicev.Slice(0, idx))
}
