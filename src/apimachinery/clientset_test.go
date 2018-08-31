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

package apimachinery

import (
	"reflect"
	"testing"

	"configcenter/src/common/metadata"
)

func TestStructMockDo(t *testing.T) {
	mockAPI := NewMockApiMachinery()

	t.Log("test with struct mock do output")
	resp := metadata.Response{
		BaseResp: metadata.SuccessBaseResp,
	}
	mockAPI.MockDo(resp).TopoServer().Object().CreateObject(nil, nil, metadata.Object{})
	rtn, err := mockAPI.TopoServer().Object().CreateObject(nil, nil, metadata.Object{})
	if err != nil {
		t.Errorf("get top server create object result failed, err: %v", err)
		return
	}

	if !reflect.DeepEqual(*rtn, resp) {
		t.Error("test with struct mock do output.")
		return
	}
	t.Log("test with struct mock do output success.")
	return
}

func TestPointerMockDo(t *testing.T) {
	mockAPI := NewMockApiMachinery()

	t.Log("test with pointer mock do output")
	resp := metadata.Response{
		BaseResp: metadata.SuccessBaseResp,
	}
	mockAPI.MockDo(&resp).TopoServer().Object().CreateObject(nil, nil, metadata.Object{})
	rtn, err := mockAPI.TopoServer().Object().CreateObject(nil, nil, metadata.Object{})
	if err != nil {
		t.Errorf("get top server create object result failed, err: %v", err)
		return
	}

	if !reflect.DeepEqual(*rtn, resp) {
		t.Error("test with pointer mock do output.")
		return
	}
	t.Log("test with pointer mock do output success.")
	return
}

func TestMapMockDo(t *testing.T) {
	mockAPI := NewMockApiMachinery()

	t.Log("test with map mock do output")
	resp := map[string]interface{}{
		"result":        true,
		"bk_error_code": 0,
		"bk_error_msg":  "success",
	}
	rtn, err := mockAPI.MockDo(resp).TopoServer().Object().CreateObject(nil, nil, metadata.Object{})
	if err != nil {
		t.Errorf("test with map mock do output. err: %v", err)
	}

	sResp := metadata.Response{
		BaseResp: metadata.SuccessBaseResp,
	}
	if !reflect.DeepEqual(*rtn, sResp) {
		t.Error("test with map mock do output.")
		return
	}
	t.Log("test with map mock do output success.")
	return
}

func TestStringMockDo(t *testing.T) {
	mockAPI := NewMockApiMachinery()

	t.Log("test with string mock do output")
	resp := `{
        "result": true,
        "bk_error_code": 0,
        "bk_error_msg": "success"
    }`

	rtn, err := mockAPI.MockDo(resp).TopoServer().Object().CreateObject(nil, nil, metadata.Object{})
	if err != nil {
		t.Errorf("test with map mock do output. err: %v", err)
	}

	sResp := metadata.Response{
		BaseResp: metadata.SuccessBaseResp,
	}
	if !reflect.DeepEqual(*rtn, sResp) {
		t.Error("test with string mock do output.")
		return
	}
	t.Log("test with string mock do output success.")
	return
}
