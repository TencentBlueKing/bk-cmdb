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

package example

import (
	"configcenter/src/framework/api"
	"configcenter/src/framework/core/output/module/model"
	"fmt"
	"time"
)

func init() {

	// api.RegisterInputer(host, nil)
	api.RegisterTimingInputer(modelMgr, time.Second*5, nil)
}

var modelMgr = &modelInputer{}

type modelInputer struct {
}

// Name the Inputer name.
// This information will be printed when the Inputer is abnormal, which is convenient for debugging.
func (cli *modelInputer) Name() string {
	return "model_inputer"
}

// Run the input should not be blocked
func (cli *modelInputer) Run() interface{} {

	// create a  classification
	clsMgr := api.CreateClassification("test_demo_cls")
	clsMgr.SetID("test_demo_cls_id")
	err := clsMgr.Save()
	if nil != err {
		fmt.Println("failed to save the classification, ", err)
		return nil
	}

	// create a  model
	clsModel := clsMgr.CreateModel()
	clsModel.SetID("test_demo_cls_model_id")
	clsModel.SetName("test_demo_cls_model_name")
	clsModel.SetSupplierAccount("0")
	clsModel.SetCreator("test_user")
	clsModel.SetDescription("test_desc")

	err = clsModel.Save()
	if nil != err {
		fmt.Println("failed to save the model, ", err)
		return nil
	}

	// create a attribute
	modelAttr := clsModel.CreateAttribute()
	modelAttr.SetDescrition("test_desc")
	modelAttr.SetID("test_demo_model_attr_id")
	modelAttr.SetName("test_demo_model_attr_name")
	modelAttr.SetType(model.FieldTypeLongChar)

	err = modelAttr.Save()
	if nil != err {
		fmt.Println("failed to save the model attribute, ", err)
		return nil
	}

	return nil
}

func (cli *modelInputer) Stop() error {
	return nil
}
