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

func TestCreateManyModelClassification(t *testing.T) {
	mockAPI := NewMockApiMachinery()

	t.Log("test with struct mock do output")
	var resp metadata.CreatedManyOptionResult

	err := json.Unmarshal([]byte(getCreatedManyResult()), &resp)
	if err != nil {
		t.Error("get response data error")
		return
	}

	mockAPI.MockDo(resp).CoreService().Model().CreateManyModelClassification(nil, nil, &metadata.CreateManyModelClassifiaction{})
	rtn, err := mockAPI.CoreService().Model().CreateManyModelClassification(nil, nil, &metadata.CreateManyModelClassifiaction{})
	if err != nil {
		t.Errorf("get  logics service create many model classification result failed, err: %v", err)
		return
	}

	if !reflect.DeepEqual(*rtn, resp) {
		t.Error("test with struct mock do output.")
		return
	}
	t.Log("test with struct mock do output success.")
	return
}

func TestCreateModelClassification(t *testing.T) {
	mockAPI := NewMockApiMachinery()

	t.Log("test with struct mock do output")
	var resp metadata.CreatedOneOptionResult

	err := json.Unmarshal([]byte(getCreatedResult()), &resp)
	if err != nil {
		t.Error("get response data error")
		return
	}

	mockAPI.MockDo(resp).CoreService().Model().CreateModelClassification(nil, nil, &metadata.CreateOneModelClassification{})
	rtn, err := mockAPI.CoreService().Model().CreateModelClassification(nil, nil, &metadata.CreateOneModelClassification{})
	if err != nil {
		t.Errorf("get  logics service create  model classification result failed, err: %v", err)
		return
	}

	if !reflect.DeepEqual(*rtn, resp) {
		t.Error("test with struct mock do output.")
		return
	}
	t.Log("test with struct mock do output success.")
	return
}

func TestSetManyModelClassification(t *testing.T) {
	mockAPI := NewMockApiMachinery()

	t.Log("test with struct mock do output")
	var resp metadata.SetOptionResult

	err := json.Unmarshal([]byte(getSetResult()), &resp)
	if err != nil {
		t.Error("get response data error")
		return
	}

	mockAPI.MockDo(resp).CoreService().Model().SetManyModelClassification(nil, nil, &metadata.SetManyModelClassification{})
	rtn, err := mockAPI.CoreService().Model().SetManyModelClassification(nil, nil, &metadata.SetManyModelClassification{})
	if err != nil {
		t.Errorf("get  logics service set many model classification result failed, err: %v", err)
		return
	}

	if !reflect.DeepEqual(*rtn, resp) {
		t.Error("test with struct mock do output.")
		return
	}
	t.Log("test with struct mock do output success.")
	return
}

func TestSetModelClassification(t *testing.T) {
	mockAPI := NewMockApiMachinery()

	t.Log("test with struct mock do output")
	var resp metadata.SetOptionResult

	err := json.Unmarshal([]byte(getSetResult()), &resp)
	if err != nil {
		t.Error("get response data error")
		return
	}

	mockAPI.MockDo(resp).CoreService().Model().SetModelClassification(nil, nil, &metadata.SetOneModelClassification{})
	rtn, err := mockAPI.CoreService().Model().SetModelClassification(nil, nil, &metadata.SetOneModelClassification{})
	if err != nil {
		t.Errorf("get  logics service set model classification result failed, err: %v", err)
		return
	}

	if !reflect.DeepEqual(*rtn, resp) {
		t.Error("test with struct mock do output.")
		return
	}
	t.Log("test with struct mock do output success.")
	return
}

func TestUpdateModelClassification(t *testing.T) {
	mockAPI := NewMockApiMachinery()

	t.Log("test with struct mock do output")
	var resp metadata.UpdatedOptionResult

	err := json.Unmarshal([]byte(getUpdatedResult()), &resp)
	if err != nil {
		t.Errorf("get response data error, err:%s", err.Error())
		return
	}

	mockAPI.MockDo(resp).CoreService().Model().UpdateModelClassification(nil, nil, &metadata.UpdateOption{})
	rtn, err := mockAPI.CoreService().Model().UpdateModelClassification(nil, nil, &metadata.UpdateOption{})
	if err != nil {
		t.Errorf("get  logics service update model classification result failed, err: %v", err)
		return
	}

	if !reflect.DeepEqual(*rtn, resp) {
		t.Error("test with struct mock do output.")
		return
	}
	t.Log("test with struct mock do output success.")
	return
}

func TestDeleteModelClassification(t *testing.T) {
	mockAPI := NewMockApiMachinery()

	t.Log("test with struct mock do output")
	var resp metadata.DeletedOptionResult

	err := json.Unmarshal([]byte(getDeletedResult()), &resp)
	if err != nil {
		t.Error("get response data error")
		return
	}

	mockAPI.MockDo(resp).CoreService().Model().DeleteModelClassification(nil, nil, &metadata.DeleteOption{})
	rtn, err := mockAPI.CoreService().Model().DeleteModelClassification(nil, nil, &metadata.DeleteOption{})
	if err != nil {
		t.Errorf("get  logics service delete model classification result failed, err: %v", err)
		return
	}

	if !reflect.DeepEqual(*rtn, resp) {
		t.Error("test with struct mock do output.")
		return
	}
	t.Log("test with struct mock do output success.")
	return
}

func TestDeleteModelClassificationCascade(t *testing.T) {
	mockAPI := NewMockApiMachinery()

	t.Log("test with struct mock do output")
	var resp metadata.DeletedOptionResult

	err := json.Unmarshal([]byte(getDeletedResult()), &resp)
	if err != nil {
		t.Error("get response data error")
		return
	}

	mockAPI.MockDo(resp).CoreService().Model().DeleteModelClassificationCascade(nil, nil, &metadata.DeleteOption{})
	rtn, err := mockAPI.CoreService().Model().DeleteModelClassificationCascade(nil, nil, &metadata.DeleteOption{})
	if err != nil {
		t.Errorf("get  logics service delete model classification cascade result failed, err: %v", err)
		return
	}

	if !reflect.DeepEqual(*rtn, resp) {
		t.Error("test with struct mock do output.")
		return
	}
	t.Log("test with struct mock do output success.")
	return
}

func TestReadModelClassification(t *testing.T) {
	mockAPI := NewMockApiMachinery()

	t.Log("test with struct mock do output")
	var resp metadata.ReadModelClassifitionResult

	err := json.Unmarshal([]byte(getReadModelClassificationResult()), &resp)
	if err != nil {
		t.Error("get response data error")
		return
	}

	mockAPI.MockDo(resp).CoreService().Model().ReadModelClassification(nil, nil, &metadata.QueryCondition{})
	rtn, err := mockAPI.CoreService().Model().ReadModelClassification(nil, nil, &metadata.QueryCondition{})
	if err != nil {
		t.Errorf("get  logics service read model classification result failed, err: %v", err)
		return
	}

	if !reflect.DeepEqual(*rtn, resp) {
		t.Error("test with struct mock do output.")
		return
	}
	t.Log("test with struct mock do output success.")
	return
}

func TestCreateModel(t *testing.T) {
	mockAPI := NewMockApiMachinery()

	t.Log("test with struct mock do output")
	var resp metadata.CreatedOneOptionResult

	err := json.Unmarshal([]byte(getCreatedResult()), &resp)
	if err != nil {
		t.Error("get response data error")
		return
	}

	mockAPI.MockDo(resp).CoreService().Model().CreateModel(nil, nil, &metadata.CreateModel{})
	rtn, err := mockAPI.CoreService().Model().CreateModel(nil, nil, &metadata.CreateModel{})
	if err != nil {
		t.Errorf("get  logics service create  model  result failed, err: %v", err)
		return
	}

	if !reflect.DeepEqual(*rtn, resp) {
		t.Error("test with struct mock do output.")
		return
	}
	t.Log("test with struct mock do output success.")
	return
}

func TestSetModel(t *testing.T) {
	mockAPI := NewMockApiMachinery()

	t.Log("test with struct mock do output")
	var resp metadata.SetOptionResult

	err := json.Unmarshal([]byte(getSetResult()), &resp)
	if err != nil {
		t.Error("get response data error")
		return
	}

	mockAPI.MockDo(resp).CoreService().Model().SetModel(nil, nil, &metadata.SetModel{})
	rtn, err := mockAPI.CoreService().Model().SetModel(nil, nil, &metadata.SetModel{})
	if err != nil {
		t.Errorf("get  logics service set model result failed, err: %v", err)
		return
	}

	if !reflect.DeepEqual(*rtn, resp) {
		t.Error("test with struct mock do output.")
		return
	}
	t.Log("test with struct mock do output success.")
	return
}

func TestUpdateModel(t *testing.T) {
	mockAPI := NewMockApiMachinery()

	t.Log("test with struct mock do output")
	var resp metadata.UpdatedOptionResult

	err := json.Unmarshal([]byte(getUpdatedResult()), &resp)
	if err != nil {
		t.Errorf("get response data error, err:%s", err.Error())
		return
	}

	mockAPI.MockDo(resp).CoreService().Model().UpdateModel(nil, nil, &metadata.UpdateOption{})
	rtn, err := mockAPI.CoreService().Model().UpdateModel(nil, nil, &metadata.UpdateOption{})
	if err != nil {
		t.Errorf("get  logics service update model result failed, err: %v", err)
		return
	}

	if !reflect.DeepEqual(*rtn, resp) {
		t.Error("test with struct mock do output.")
		return
	}
	t.Log("test with struct mock do output success.")
	return
}

func TestDeleteModel(t *testing.T) {
	mockAPI := NewMockApiMachinery()

	t.Log("test with struct mock do output")
	var resp metadata.DeletedOptionResult

	err := json.Unmarshal([]byte(getDeletedResult()), &resp)
	if err != nil {
		t.Error("get response data error")
		return
	}

	mockAPI.MockDo(resp).CoreService().Model().DeleteModel(nil, nil, &metadata.DeleteOption{})
	rtn, err := mockAPI.CoreService().Model().DeleteModel(nil, nil, &metadata.DeleteOption{})
	if err != nil {
		t.Errorf("get  logics service delete model result failed, err: %v", err)
		return
	}

	if !reflect.DeepEqual(*rtn, resp) {
		t.Error("test with struct mock do output.")
		return
	}
	t.Log("test with struct mock do output success.")
	return
}

func TestDeleteModelCascade(t *testing.T) {
	mockAPI := NewMockApiMachinery()

	t.Log("test with struct mock do output")
	var resp metadata.DeletedOptionResult

	err := json.Unmarshal([]byte(getDeletedResult()), &resp)
	if err != nil {
		t.Error("get response data error")
		return
	}

	mockAPI.MockDo(resp).CoreService().Model().DeleteModelCascade(nil, nil, 0)
	rtn, err := mockAPI.CoreService().Model().DeleteModelCascade(nil, nil, 0)
	if err != nil {
		t.Errorf("get  logics service delete model cascade result failed, err: %v", err)
		return
	}

	if !reflect.DeepEqual(*rtn, resp) {
		t.Error("test with struct mock do output.")
		return
	}
	t.Log("test with struct mock do output success.")
	return
}

func TestReadModel(t *testing.T) {
	mockAPI := NewMockApiMachinery()

	t.Log("test with struct mock do output")
	var resp metadata.ReadModelResult

	err := json.Unmarshal([]byte(getReadModelResult()), &resp)
	if err != nil {
		t.Error("get response data error")
		return
	}

	mockAPI.MockDo(resp).CoreService().Model().ReadModel(nil, nil, &metadata.QueryCondition{})
	rtn, err := mockAPI.CoreService().Model().ReadModel(nil, nil, &metadata.QueryCondition{})
	if err != nil {
		t.Errorf("get  logics service read model result failed, err: %v", err)
		return
	}

	if !reflect.DeepEqual(*rtn, resp) {
		t.Error("test with struct mock do output.")
		return
	}
	t.Log("test with struct mock do output success.")
	return
}

func TestCreateModelAttrs(t *testing.T) {
	mockAPI := NewMockApiMachinery()

	t.Log("test with struct mock do output")
	var resp metadata.CreatedManyOptionResult

	err := json.Unmarshal([]byte(getCreatedManyResult()), &resp)
	if err != nil {
		t.Error("get response data error")
		return
	}

	mockAPI.MockDo(resp).CoreService().Model().CreateModelAttrs(nil, nil, "", &metadata.CreateModelAttributes{})
	rtn, err := mockAPI.CoreService().Model().CreateModelAttrs(nil, nil, "", &metadata.CreateModelAttributes{})
	if err != nil {
		t.Errorf("get  logics service create  model attribute result failed, err: %v", err)
		return
	}

	if !reflect.DeepEqual(*rtn, resp) {
		t.Error("test with struct mock do output.")
		return
	}
	t.Log("test with struct mock do output success.")
	return
}

func TestUpdateModelAttrs(t *testing.T) {
	mockAPI := NewMockApiMachinery()

	t.Log("test with struct mock do output")
	var resp metadata.UpdatedOptionResult

	err := json.Unmarshal([]byte(getUpdatedResult()), &resp)
	if err != nil {
		t.Errorf("get response data error, err:%s", err.Error())
		return
	}

	mockAPI.MockDo(resp).CoreService().Model().UpdateModelAttrs(nil, nil, "", &metadata.UpdateOption{})
	rtn, err := mockAPI.CoreService().Model().UpdateModelAttrs(nil, nil, "", &metadata.UpdateOption{})
	if err != nil {
		t.Errorf("get  logics service update model attribute result failed, err: %v", err)
		return
	}

	if !reflect.DeepEqual(*rtn, resp) {
		t.Error("test with struct mock do output.")
		return
	}
	t.Log("test with struct mock do output success.")
	return
}

func TestSetModelAttrs(t *testing.T) {
	mockAPI := NewMockApiMachinery()

	t.Log("test with struct mock do output")
	var resp metadata.SetOptionResult

	err := json.Unmarshal([]byte(getSetResult()), &resp)
	if err != nil {
		t.Error("get response data error")
		return
	}

	mockAPI.MockDo(resp).CoreService().Model().SetModelAttrs(nil, nil, "", &metadata.SetModelAttributes{})
	rtn, err := mockAPI.CoreService().Model().SetModelAttrs(nil, nil, "", &metadata.SetModelAttributes{})
	if err != nil {
		t.Errorf("get  logics service set model attribute result failed, err: %v", err)
		return
	}

	if !reflect.DeepEqual(*rtn, resp) {
		t.Error("test with struct mock do output.")
		return
	}
	t.Log("test with struct mock do output success.")
	return
}

func TestDeleteModelAttr(t *testing.T) {
	mockAPI := NewMockApiMachinery()

	t.Log("test with struct mock do output")
	var resp metadata.DeletedOptionResult

	err := json.Unmarshal([]byte(getDeletedResult()), &resp)
	if err != nil {
		t.Error("get response data error")
		return
	}

	mockAPI.MockDo(resp).CoreService().Model().DeleteModelAttr(nil, nil, "", &metadata.DeleteOption{})
	rtn, err := mockAPI.CoreService().Model().DeleteModelAttr(nil, nil, "", &metadata.DeleteOption{})
	if err != nil {
		t.Errorf("get  logics service delete model attribute result failed, err: %v", err)
		return
	}

	if !reflect.DeepEqual(*rtn, resp) {
		t.Error("test with struct mock do output.")
		return
	}
	t.Log("test with struct mock do output success.")
	return
}

func TestReadModelAttr(t *testing.T) {
	mockAPI := NewMockApiMachinery()

	t.Log("test with struct mock do output")
	var resp metadata.QueryConditionResult

	err := json.Unmarshal([]byte(getReadModelAttrResult()), &resp)
	if err != nil {
		t.Error("get response data error")
		return
	}

	mockAPI.MockDo(resp).CoreService().Model().ReadModelAttr(nil, nil, "", &metadata.QueryCondition{})
	rtn, err := mockAPI.CoreService().Model().ReadModelAttr(nil, nil, "", &metadata.QueryCondition{})
	if err != nil {
		t.Errorf("get  logics service read model attribute result failed, err: %v", err)
		return
	}

	if !reflect.DeepEqual(*rtn, resp) {
		t.Error("test with struct mock do output.")
		return
	}
	t.Log("test with struct mock do output success.")
	return
}

func getReadModelAttrResult() string {
	return `{
		"result": true,
		"bk_error_code": 0,
		"bk_error_msg": null,
		"data": {
			"count":1,
			"info":[{
				"bk_property_group" : "default",
				"unit" : "",
				"ispre" : true,
				"bk_property_type" : "singlechar",
				"placeholder" : "",
				"editable" : true,
				"last_time" : "2018-11-19T08:10:20.554Z",
				"bk_property_name" : "业务名",
				"option" : "",
				"description" : "",
				"bk_isapi" : false,
				"creator" : "cc_system",
				"bk_property_index" : 0,
				"isrequired" : true,
				"isreadonly" : false,
				"isonly" : true,
				"bk_issystem" : false,
				"create_time" : "2018-11-19T08:10:20.554Z"
			}]
		}
	}`
}

func getReadModelResult() string {
	return `{
		"result": true,
		"bk_error_code": 0,
		"bk_error_msg": null,
		"data": {
			 "count":1,
			 "info":[
				{
					"spec":{
						"bk_obj_id" : "",
						"ispre" : true,
						"description" : "",
						"create_time" : "2018-11-19T08:10:20.538Z",
						"modifier" : "",
						"last_time" : "2018-11-19T08:10:20.538Z",
						"id" : 1,
						"bk_obj_icon" : "icon-cc-host",
						"bk_ispaused" : false,
						"position" : "{\"bk_host_manage\":{\"x\":-600,\"y\":-650}}",
						"bk_classification_id" : "bk_host_manage",
						"bk_obj_name" : "主机",
						"bk_supplier_account" : "0",
						"creator" : "cc_system"
					},
					"attributes":[{
						"bk_property_group" : "default",
						"unit" : "",
						"ispre" : true,
						"bk_property_type" : "singlechar",
						"bk_supplier_account" : "0",
						"bk_property_id" : "bk_biz_name",
						"placeholder" : "",
						"editable" : true,
						"last_time" : "2018-11-19T08:10:20.554Z",
						"bk_property_name" : "业务名",
						"option" : "",
						"description" : "",
						"bk_isapi" : false,
						"creator" : "cc_system",
						"bk_obj_id" : "biz",
						"bk_property_index" : 0,
						"isrequired" : true,
						"isreadonly" : false,
						"isonly" : true,
						"bk_issystem" : false,
						"create_time" : "2018-11-19T08:10:20.554Z"
					}]
				}
			 ]
		}
	}`
}

func getReadModelClassificationResult() string {
	return `{
		"result": true,
		"bk_error_code": 0,
		"bk_error_msg": "",
		"data": {
			"count":1,
			"info":[
				{
					"bk_classification_icon" : "",
					"bk_classification_id" : "",
					"bk_classification_name" : "",
					"bk_classification_type" : "",
					"bk_supplier_account" : ""
				}
			 ]
		}
	}`
}

func getDeletedResult() string {
	return `{
		"result": true,
		"bk_error_code": 0,
		"bk_error_msg": "",
		"data": {
				"deleted_count":0
		}
	}`
}

func getUpdatedResult() string {
	return `{
		"result": true,
		"bk_error_code": 0,
		"bk_error_msg": "",
		"data": {
				"updated_count":0
		}
	}`
}

func getSetResult() string {
	return `{
		"result": true,
		"bk_error_code": 0,
		"bk_error_msg": "",
		"data": {
				"created_count":0,
				"updated_count":0,
				"created":[{"id":0}],
				"updated":[{"id":1}],
				"exception":[{
					"message":"error message",
					"code":0,
					"data":{}
				}]
		}
	}`
}

func getCreatedResult() string {
	return `{
		"result": true,
		"bk_error_code": 0,
		"bk_error_msg": "",
		"data": {
				"created":{"id":0}
		}
	}`
}

func getCreatedManyResult() string {
	return `{
		"result": true,
		"bk_error_code": 0,
		"bk_error_msg": "",
		"data": {
				"created":[{"id":0}],
				"repeated":[{}],
				"exception":[{
					"message":"error message",
					"code":0,
					"data":{}
				}]
		}
		
	}`

}
