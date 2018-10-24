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

package metadata

import (
	types "configcenter/src/common/mapstr"
	"testing"
)

func TestAttribute(t *testing.T) {
	m, err := types.NewFromInterface(map[string]interface{}{"id": 0, "bk_supplier_account": "bk_supplier_account"})
	attr := &Attribute{}
	attr, err = attr.Parse(m)
	if str, _ := attr.ToMapStr().String("bk_supplier_account"); str != "bk_supplier_account" || err != nil {
		t.Fail()
	}

}
