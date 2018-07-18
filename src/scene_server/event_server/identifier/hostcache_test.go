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

package identifier

import (
	"encoding/json"
	"fmt"
	"testing"

	"configcenter/src/scene_server/event_server/types"
)

func TestIdentifier(t *testing.T) {
	i := HostIdentifier{
		HostID:          1,
		HostName:        "etcd-1",
		SupplierID:      0,
		SupplierAccount: "0",
		CloudID:         0,
		CloudName:       "default area",
		InnerIP:         "192.168.1.7",
		OuterIP:         "192.168.1.10",
		OSType:          "1",
		OSName:          "linux centos",
		Memory:          64131,
		CPU:             24,
		Disk:            271,
		Module: map[string]*Module{
			"module10": {
				BizID:      10,
				BizName:    "biz10",
				SetID:      10,
				SetName:    "set10",
				ModuleID:   10,
				ModuleName: "module10",
				SetStatus:  "1",
				SetEnv:     "1",
			},
		},
	}

	e := types.EventInst{}
	e.Data = nil
	e.EventType = types.EventTypeRelation
	e.ObjType = "hostidentifier"
	e.Action = types.EventActionUpdate

	d := types.EventData{CurData: i}
	e.Data = append(e.Data, d)

	out, _ := json.Marshal(e)
	fmt.Printf("%s\n", out)
}
