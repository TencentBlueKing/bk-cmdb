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
	"encoding/json"
	"testing"
)

func TestConvTime(t *testing.T) {
	condStr := `{"condition":{"create_time":{"$lt":"2018-05-31 00:00:00", "cc_time_type":"1"},"$or":[{"bk_biz_maintainer":{"$regex":"admin"}},{"bk_biz_productor":{"$regex":"admin"}},{"bk_biz_tester":{"$regex":"admin"}},{"bk_biz_developer":{"$regex":"admin"}},{"operator":{"$regex":"admin"}}],"bk_supplier_account":"0","default":0},"fields":"","start":0,"limit":0,"sort":""}`
	var input ObjQueryInput
	err := json.Unmarshal([]byte(condStr), &input)
	if nil != err {
		t.Errorf("json unmarshal error:%s", err.Error())
		return
	}
	input.ConvTime()
	t.Logf("%v", input.Condition)

}
