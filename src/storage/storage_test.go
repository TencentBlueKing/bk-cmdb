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

package storage

import (
	_ "configcenter/src/source_controller/common/metadata"
	"fmt"
	"testing"
	_ "time"
)

func TestType1(t *testing.T) {
	db, _ := NewDB("127.0.0.1", "27017", "user", "pwd", "cmdb", "mongodb")
	db.Open()
	condition := make(map[string]interface{})
	condition["ApplicationID"] = 17
	//	host["HostName"] = "vm2"
	//	host["InnerIP"] = "127.0.0.1"
	//	db.Insert("cc_HostBase", host)
	var result interface{}
	//fields := []string{"ApplicationName"}
	fields := make([]string, 0)
	db.GetOneByCondition("cc_ApplicationBase", fields, nil, &result)
	fmt.Println(result)
}
