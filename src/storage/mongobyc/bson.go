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

import "unsafe"

type bson struct {
	data string
}

func (b *bson) ToDocument() (*C.bson_t, error) {

	var err C.bson_error_t
	innerData := C.CString(b.data)
	defer C.free(unsafe.Pointer(innerData))
	bson := C.bson_new_from_json((*C.uint8_t)(unsafe.Pointer(innerData)), -1, &err)
	if nil == bson {
		return nil, TransformError(err)
	}
	return bson, nil
}
