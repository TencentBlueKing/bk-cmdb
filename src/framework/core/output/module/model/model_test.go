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
	"configcenter/src/framework/core/output/module/model"
	//"configcenter/src/framework/core/types"
	"testing"
)

func TestSearchModel(t *testing.T) {

	client.NewForConfig(config.Config{"core.supplierAccount": "0", "core.user": "build_user", "core.ccaddress": "http://test.apiserver:8080"}, nil)

	items, err := model.FindClassificationsByCondition(common.CreateCondition().Field("bk_classification_id").Eq("bk_host_manage"))
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

		modelIterator, err := item.FindModelsByCondition(common.CreateCondition().Field("bk_obj_id").Eq("host"))

		if nil != err {
			t.Errorf("failed to search model, %s", err.Error())
			break
		}

		if nil == modelIterator {
			break
		}

		modelIterator.ForEach(func(modelItem model.Model) error {

			t.Logf("the model:%+v", modelItem.GetName())

			attrs, _ := modelItem.Attributes()
			for _, attr := range attrs {
				t.Logf("the attribute:%+v", attr.GetName())
			}

			return nil
		})

	}
}
