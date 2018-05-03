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
	"configcenter/src/framework/core/output/module/v3"
	//"configcenter/src/framework/core/types"
	"fmt"
	"testing"
)

func TestSearchInst(t *testing.T) {

	cli := v3.GetV3Client()
	cli.SetSupplierAccount("0")
	cli.SetUser("build_user")
	cli.SetAddress("http://test.apiserver:8080")

	cond := common.CreateCondition().Field("bk_obj_id").Eq("test1")
	dataMap, err := cli.SearchInst(cond)

	if nil != err {
		t.Errorf("failed to search, error info is %s", err.Error())
	}

	for _, item := range dataMap {
		fmt.Println("item:", item)
	}

}
