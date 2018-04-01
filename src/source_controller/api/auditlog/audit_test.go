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
 
package auditlog

import (
	"configcenter/src/common"
	"configcenter/src/common/auditoplog"
	"configcenter/src/source_controller/common/commondata"
	"encoding/json"
	"fmt"
	"time"
	//"configcenter/src/source_controller/common/commondata"
	//"encoding/json"
	"testing"
	//"time"
)

var (
	appID   = "99999999"
	ownerID = common.BKDefaultOwnerID
	user    = "user"
	opArr   = []string{"add", "modify", "del"}
	client  = NewClient("http://127.0.0.1:60001")
)

func TestAuditAppCreate(t *testing.T) {

	contentType := "app"
	for index, op := range opArr {
		content := fmt.Sprintf("%s %s", contentType, op)
		_, err := client.AuditAppLog(0, content, content, ownerID, appID, user, auditoplog.AuditOpType(index+1))
		if nil != err {
			t.Errorf(err.Error())
		}
	}

}

func TestAuditSetCreate(t *testing.T) {

	contentType := "set"
	for index, op := range opArr {
		content := fmt.Sprintf("%s %s", contentType, op)
		_, err := client.AuditSetLog(0, content, content, ownerID, appID, user, auditoplog.AuditOpType(index+1))
		if nil != err {
			t.Errorf(err.Error())
		}
	}

	for index, op := range opArr {
		content := []auditoplog.AuditLogContext{auditoplog.AuditLogContext{Content: fmt.Sprintf("%s mutli 1 %s", contentType, op), ID:1}, auditoplog.AuditLogContext{Content: fmt.Sprintf("%s mutli 2 %s", contentType, op), ID:1}}
		_, err := client.AuditSetsLog(content, fmt.Sprintf("%s mutil %s", contentType, op), ownerID, appID, user, auditoplog.AuditOpType(index+1))
		if nil != err {
			t.Errorf(err.Error())
		}
	}

}
func TestAuditModuleCreate(t *testing.T) {

	contentType := "module"
	for index, op := range opArr {
		content := fmt.Sprintf("%s %s", contentType, op)
		_, err := client.AuditModuleLog(1, content, content, ownerID, appID, user, auditoplog.AuditOpType(index+1))
		if nil != err {
			t.Errorf(err.Error())
		}
	}

	for index, op := range opArr {
		content := []auditoplog.AuditLogContext{auditoplog.AuditLogContext{Content: fmt.Sprintf("%s mutli 1 %s", contentType, op), ID:1}, auditoplog.AuditLogContext{Content: fmt.Sprintf("%s mutli 2 %s", contentType, op), ID:2}}
		_, err := client.AuditModulesLog(content, fmt.Sprintf("%s mutil %s", contentType, op), ownerID, appID, user, auditoplog.AuditOpType(index+1))
		if nil != err {
			t.Errorf(err.Error())
		}
	}

}

func TestAuditHostCreate(t *testing.T) {

	contentType := "host"
	for index, op := range opArr {
		content := fmt.Sprintf("%s %s", contentType, op)
		_, err := client.AuditHostLog(0, content, content, "127.0.0.1", ownerID, appID, user, auditoplog.AuditOpType(index+1))
		if nil != err {
			t.Errorf(err.Error())
		}
	}

	for index, op := range opArr {
		content := []auditoplog.AuditLogExt{
			auditoplog.AuditLogExt{Content: fmt.Sprintf("%s mutli 1 %s", contentType, op), ExtKey: "127.0.1.1", ID:1},
			auditoplog.AuditLogExt{Content: fmt.Sprintf("%s mutli 2 %s", contentType, op), ExtKey: "127.0.1.2", ID:2},
		}
		_, err := client.AuditHostsLog(content, fmt.Sprintf("%s mutil %s", contentType, op), ownerID, appID, user, auditoplog.AuditOpType(index+1))
		if nil != err {
			t.Errorf(err.Error())
		}
	}

}

func TestAuditProcCreate(t *testing.T) {

	contentType := "proc"
	for index, op := range opArr {
		content := fmt.Sprintf("%s %s", contentType, op)
		_, err := client.AuditProcLog(0, content, content, ownerID, appID, user, auditoplog.AuditOpType(index+1))
		if nil != err {
			t.Errorf(err.Error())
		}
	}

	for index, op := range opArr {
		content := []auditoplog.AuditLogContext{auditoplog.AuditLogContext{Content: fmt.Sprintf("%s mutli 1 %s", contentType, op), ID:1}, auditoplog.AuditLogContext{Content: fmt.Sprintf("%s mutli 2 %s", contentType, op), ID:2}}
		_, err := client.AuditProcsLog(content, fmt.Sprintf("%s mutil %s", contentType, op), ownerID, appID, user, auditoplog.AuditOpType(index+1))
		if nil != err {
			t.Errorf(err.Error())
		}
	}

}

func TestAuditObjCreate(t *testing.T) {

	contentType := "proc"
	for index, op := range opArr {
		content := fmt.Sprintf("%s %s", contentType, op)
		_, err := client.AuditObjLog(0, content, content, "test obj", ownerID, appID, user, auditoplog.AuditOpType(index+1))
		if nil != err {
			t.Errorf(err.Error())
		}
	}

	for index, op := range opArr {
		content := []auditoplog.AuditLogContext{auditoplog.AuditLogContext{Content: fmt.Sprintf("%s mutli 1 %s", contentType, op)}, auditoplog.AuditLogContext{Content: fmt.Sprintf("%s mutli 2 %s", contentType, op)}}
		_, err := client.AuditObjsLog(content, fmt.Sprintf("%s mutil %s", contentType, op), "test objs", ownerID, appID, user, auditoplog.AuditOpType(index+1))
		if nil != err {
			t.Errorf(err.Error())
		}
	}

}

func TestAuditLog(t *testing.T) {
	var input commondata.ObjQueryInput
	input.Start = 0
	input.Limit = 100
	times := []string{"2017-12-19 17:52:24", "2018-12-30 17:52:25"}
	//2017-12-27T17:52:24
	//start, errStart := time.Parse("17-11-20", "10-09-01")
	start, errStart := time.Parse("2006-01-02 15:04:05", times[0])
	if nil != errStart {
		t.Errorf(errStart.Error())
	}

	end, errEnd := time.Parse("2006-01-02 15:04:05", times[1])
	if nil != errEnd {
		t.Errorf(errEnd.Error())
	}
	input.Condition = common.KvMap{common.BKAppIDField: 99999999, common.BKOpTargetField: "host", common.CreateTimeField: common.KvMap{"$gte": start.UTC(), "$lte": end.UTC(), commondata.CC_time_type_parse_flag: 1}}
	data, err := client.GetAuditlogs(input)

	fmt.Println(start.UTC(), end.Local())

	if nil != err {
		t.Errorf(err.Error())
	}
	ret, _ := json.Marshal(data)
	fmt.Println(string(ret))
}
