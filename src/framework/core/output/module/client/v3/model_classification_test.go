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
 
package v3_test

import (
	"configcenter/src/framework/common"
	"configcenter/src/framework/core/config"
	"configcenter/src/framework/core/output/module/client"
	"configcenter/src/framework/core/types"
	"fmt"
	"testing"
)

func TestCreateClassification(t *testing.T) {

	cli := client.NewForConfig(config.Config{"core.supplierAccount": "0", "core.user": "build_user", "core.ccaddress": "http://test.apiserver:8080"}, nil)

	id, err := cli.CCV3().Classification().CreateClassification(types.MapStr{
		"bk_classification_id":   common.UUID(),
		"bk_classification_name": fmt.Sprintf("test_%s", common.UUID()),
	})

	if nil != err {
		t.Errorf("failed to create, error info is %s", err.Error())
	}

	t.Logf("id:%d", id)
}

func TestDeleteClassification(t *testing.T) {
	cli := client.NewForConfig(config.Config{"core.supplierAccount": "0", "core.user": "build_user", "core.ccaddress": "http://test.apiserver:8080"}, nil)

	cond := common.CreateCondition().Field("id").Eq(9)

	err := cli.CCV3().Classification().DeleteClassification(cond)

	if nil != err {
		t.Errorf("failed to create, error info is %s", err.Error())
	}

	t.Log("success")
}

func TestUpdateClassification(t *testing.T) {
	cli := client.NewForConfig(config.Config{"core.supplierAccount": "0", "core.user": "build_user", "core.ccaddress": "http://test.apiserver:8080"}, nil)

	cond := common.CreateCondition().Field("id").Eq(10)

	err := cli.CCV3().Classification().UpdateClassification(types.MapStr{"bk_classification_name": "test_update"}, cond)

	if nil != err {
		t.Errorf("failed to update, error info is %s", err.Error())
	}

	t.Log("success")
}
func TestSearchClassification(t *testing.T) {
	cli := client.NewForConfig(config.Config{"core.supplierAccount": "0", "core.user": "build_user", "core.ccaddress": "http://test.apiserver:8080"}, nil)

	cond := common.CreateCondition().Field("bk_classification_name").Like("test_")

	dataMap, err := cli.CCV3().Classification().SearchClassifications(cond)

	if nil != err {
		t.Errorf("failed to create, error info is %s", err.Error())
	}

	for _, item := range dataMap {
		t.Logf("success, data:%+v", item.String("bk_classification_name"))
	}

}

func TestSearchClassificationWithObjects(t *testing.T) {
	cli := client.NewForConfig(config.Config{"core.supplierAccount": "0", "core.user": "build_user", "core.ccaddress": "http://test.apiserver:8080"}, nil)

	cond := common.CreateCondition().Field("bk_classification_name").Like("业务")

	dataMap, err := cli.CCV3().Classification().SearchClassificationWithObjects(cond)

	if nil != err {
		t.Errorf("failed to create, error info is %s", err.Error())
	}

	t.Logf("data:%+v", dataMap)

	for _, item := range dataMap {
		data, dataErr := item.MapStrArray("bk_objects")
		if nil != dataErr {
			t.Errorf("the error is %s", dataErr.Error())
			continue
		}

		//t.Logf("success,time:%v data:%+v ", tm, data)
		for _, subItem := range data {
			tm, tmErr := subItem.Time("create_time")
			if nil != tmErr {
				t.Errorf("the error info is %s", tmErr.Error())
				continue
			}
			t.Logf("success,time:%v ", tm)
		}
	}

}
