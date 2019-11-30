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
	"encoding/json"
	"reflect"
	"testing"

	"configcenter/src/common/metadata"
)

func TestCreateModelInst(t *testing.T) {
	mockAPI := NewMockApiMachinery()

	t.Log("test with struct mock do output")
	var resp metadata.CreatedOneOptionResult

	err := json.Unmarshal([]byte(getCreatedResult()), &resp)
	if err != nil {
		t.Error("get response data error")
		return
	}

	mockAPI.MockDo(resp).CoreService().Instance().CreateInstance(nil, nil, "", &metadata.CreateModelInstance{})
	rtn, err := mockAPI.CoreService().Instance().CreateInstance(nil, nil, "", &metadata.CreateModelInstance{})
	if err != nil {
		t.Errorf("get  logics service create  instance result failed, err: %v", err)
		return
	}

	if !reflect.DeepEqual(*rtn, resp) {
		t.Error("test with struct mock do output.")
		return
	}
	t.Log("test with struct mock do output success.")
	return
}

func TestCreateManyModelInstance(t *testing.T) {
	mockAPI := NewMockApiMachinery()

	t.Log("test with struct mock do output")
	var resp metadata.CreatedManyOptionResult

	err := json.Unmarshal([]byte(getCreatedManyResult()), &resp)
	if err != nil {
		t.Error("get response data error")
		return
	}

	mockAPI.MockDo(resp).CoreService().Instance().CreateManyInstance(nil, nil, "", &metadata.CreateManyModelInstance{})
	rtn, err := mockAPI.CoreService().Instance().CreateManyInstance(nil, nil, "", &metadata.CreateManyModelInstance{})
	if err != nil {
		t.Errorf("get  logics service create many instance result failed, err: %v", err)
		return
	}

	if !reflect.DeepEqual(*rtn, resp) {
		t.Error("test with struct mock do output.")
		return
	}
	t.Log("test with struct mock do output success.")
	return
}

func TestSetManyModelInstace(t *testing.T) {
	mockAPI := NewMockApiMachinery()

	t.Log("test with struct mock do output")
	var resp metadata.SetOptionResult

	err := json.Unmarshal([]byte(getSetResult()), &resp)
	if err != nil {
		t.Error("get response data error")
		return
	}

	mockAPI.MockDo(resp).CoreService().Instance().SetManyInstance(nil, nil, "", &metadata.SetManyModelInstance{})
	rtn, err := mockAPI.CoreService().Instance().SetManyInstance(nil, nil, "", &metadata.SetManyModelInstance{})
	if err != nil {
		t.Errorf("get  logics service set many model instance result failed, err: %v", err)
		return
	}

	if !reflect.DeepEqual(*rtn, resp) {
		t.Error("test with struct mock do output.")
		return
	}
	t.Log("test with struct mock do output success.")
	return
}

func TestUpdateModelInstance(t *testing.T) {
	mockAPI := NewMockApiMachinery()

	t.Log("test with struct mock do output")
	var resp metadata.UpdatedOptionResult

	err := json.Unmarshal([]byte(getUpdatedResult()), &resp)
	if err != nil {
		t.Errorf("get response data error, err:%s", err.Error())
		return
	}

	mockAPI.MockDo(resp).CoreService().Instance().UpdateInstance(nil, nil, "", &metadata.UpdateOption{})
	rtn, err := mockAPI.CoreService().Instance().UpdateInstance(nil, nil, "", &metadata.UpdateOption{})
	if err != nil {
		t.Errorf("get  logics service update model many result failed, err: %v", err)
		return
	}

	if !reflect.DeepEqual(*rtn, resp) {
		t.Error("test with struct mock do output.")
		return
	}
	t.Log("test with struct mock do output success.")
	return
}

func TestReadModelInstance(t *testing.T) {
	mockAPI := NewMockApiMachinery()

	t.Log("test with struct mock do output")
	var resp metadata.QueryConditionResult

	err := json.Unmarshal([]byte(getReadModelClassificationResult()), &resp)
	if err != nil {
		t.Error("get response data error")
		return
	}

	mockAPI.MockDo(resp).CoreService().Instance().ReadInstance(nil, nil, "", &metadata.QueryCondition{})
	rtn, err := mockAPI.CoreService().Instance().ReadInstance(nil, nil, "", &metadata.QueryCondition{})
	if err != nil {
		t.Errorf("get  logics service read model instance result failed, err: %v", err)
		return
	}

	if !reflect.DeepEqual(*rtn, resp) {
		t.Error("test with struct mock do output.")
		return
	}
	t.Log("test with struct mock do output success.")
	return
}

func TestDeleteModelInstance(t *testing.T) {
	mockAPI := NewMockApiMachinery()

	t.Log("test with struct mock do output")
	var resp metadata.DeletedOptionResult

	err := json.Unmarshal([]byte(getDeletedResult()), &resp)
	if err != nil {
		t.Error("get response data error")
		return
	}

	mockAPI.MockDo(resp).CoreService().Instance().DeleteInstance(nil, nil, "", &metadata.DeleteOption{})
	rtn, err := mockAPI.CoreService().Instance().DeleteInstance(nil, nil, "", &metadata.DeleteOption{})
	if err != nil {
		t.Errorf("get  logics service delete model instance result failed, err: %v", err)
		return
	}

	if !reflect.DeepEqual(*rtn, resp) {
		t.Error("test with struct mock do output.")
		return
	}
	t.Log("test with struct mock do output success.")
	return
}

func TestDeleteModelInstanceCascade(t *testing.T) {
	mockAPI := NewMockApiMachinery()

	t.Log("test with struct mock do output")
	var resp metadata.DeletedOptionResult

	err := json.Unmarshal([]byte(getDeletedResult()), &resp)
	if err != nil {
		t.Error("get response data error")
		return
	}

	mockAPI.MockDo(resp).CoreService().Instance().DeleteInstanceCascade(nil, nil, "", &metadata.DeleteOption{})
	rtn, err := mockAPI.CoreService().Instance().DeleteInstanceCascade(nil, nil, "", &metadata.DeleteOption{})
	if err != nil {
		t.Errorf("get  logics service delete model instance cascade result failed, err: %v", err)
		return
	}

	if !reflect.DeepEqual(*rtn, resp) {
		t.Error("test with struct mock do output.")
		return
	}
	t.Log("test with struct mock do output success.")
	return
}

func getReadModelInstanceResult() string {
	return `{
    "result": true,
    "bk_error_code": 0,
    "bk_error_msg": null,
    "data": {
        "count":1,
        "info":[{
            "bk_inst_id":1
            }]
    	}
	}`
}
