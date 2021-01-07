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
 
package util

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestGetCurrentTimeStr(t *testing.T) {
	now := time.Now()
	val := GetCurrentTimeStr()
	valTime, err := time.ParseInLocation("2006-01-02 15:04:05", val, time.Local)
	require.NoError(t, err)
	require.InDelta(t, now.Unix(), valTime.Unix(), 1)
}

func TestConvParamsTime(t *testing.T) {
	//strJSON := `{"page":{"start":0,"limit":10,"sort":"bk_host_id"},"pattern":"","bk_biz_id":2,"ip":{"flag":"bk_host_innerip|bk_host_outerip","exact":0,"data":[]},"condition":[{"bk_obj_id":"host","fields":[],"condition":[{"create_time":["2018-03-04","2018-03-17"]}]},{"bk_obj_id":"biz","fields":[],"condition":[{"field":"default","operator":"$ne","value":1}]},{"bk_obj_id":"module","fields":[],"condition":[]},{"bk_obj_id":"set","fields":[],"condition":[]}]}`
	strJSON := `{"bk_host_id":{"$in":[99,100,101,102,103,104]},"create_time":{"$in":["2018-03-16 02:45:28","2018-03-16"]}}`
	var a interface{}
	err := json.Unmarshal([]byte(strJSON), &a)
	if nil != err {
		t.Error(err.Error())
	}
	fmt.Println("====================")
	a = ConvParamsTime(a)
	fmt.Println(a)

}
