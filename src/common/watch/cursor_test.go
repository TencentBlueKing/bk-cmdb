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

package watch

import (
	"testing"

	"configcenter/src/storage/stream/types"
)

const cursorSample = "MQ0yDTVlYjM4NTk3NDc3MGExMThmNDkyMmFiZQ0xNTg4ODUzNjUyDTA="

func TestCursorEncode(t *testing.T) {
	cursor := Cursor{
		ClusterTime: types.TimeStamp{
			Sec:  uint32(1588853652),
			Nano: 0,
		},
		Oid:  "5eb385974770a118f4922abe",
		Type: Host,
	}
	encode, err := cursor.Encode()
	if err != nil {
		t.Errorf("encode cursor failed, err: %v", err)
		return
	}

	if encode != cursorSample {
		t.Errorf("encode cursor failed")
		return
	}
}

func TestCursorDecode(t *testing.T) {
	cursor := new(Cursor)
	if err := cursor.Decode(cursorSample); err != nil {
		t.Errorf("decode cursor failed, err: %v", err)
		return
	}

	if cursor.Oid != "5eb385974770a118f4922abe" {
		t.Errorf("decode cursor, got invalid oid: %s", cursor.Oid)
		return
	}

	if cursor.ClusterTime.Sec != uint32(1588853652) {
		t.Errorf("decode cursor, got invalid cluster time sec: %d", cursor.ClusterTime.Sec)
		return
	}

	if cursor.ClusterTime.Nano != uint32(0) {
		t.Errorf("decode cursor, got invalid cluster time nano: %d", cursor.ClusterTime.Nano)
		return
	}

	if cursor.Type != Host {
		t.Errorf("decode cursor, got invalid cursor type: %s", cursor.Type)
		return
	}
}
