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
 
package model_test

import (
	"configcenter/src/framework/common"
	"configcenter/src/framework/core/config"
	"configcenter/src/framework/core/output/module/client"
	"configcenter/src/framework/core/output/module/client/v3"
	"configcenter/src/framework/core/output/module/model"
	//"configcenter/src/framework/core/types"
	"testing"
)

func TestSearchClassification(t *testing.T) {

	client.NewForConfig(config.Config{"core.supplierAccount": "0", "core.user": "build_user", "core.ccaddress": "http://test.apiserver:8080"}, nil)

	cond := common.CreateCondition().Field("id").Gt(1)

	items, err := model.FindClassificationsByCondition(cond)
	if nil != err {
		t.Errorf("failed to find classifications, %s", err.Error())
		return
	}

	for {
		item, err := items.Next()
		if nil != err {
			t.Errorf("failed to get next classification, %s ", err.Error())
			break
		}
		if nil == item {
			t.Log("exit")
			break
		}
		t.Logf("the classifications:%+v", item)
	}

}

func TestSearchClassificationByName(t *testing.T) {

	cli := client.GetClient().CCV3().(*v3.Client)
	cli.SetSupplierAccount("0")
	cli.SetUser("build_user")
	cli.SetAddress("http://test.apiserver:8080")

	items, err := model.FindClassificationsLikeName("tes")
	if nil != err {
		t.Errorf("failed to find classifications, %s", err.Error())
		return
	}

	for {
		item, err := items.Next()
		if nil != err {
			t.Errorf("failed to get next classification, %s ", err.Error())
			break
		}
		if nil == item {
			t.Log("exit")
			break
		}
		t.Logf("the classifications:%+v", item)
	}

}
