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
 
package rdiscover

import (
	"testing"
	"encoding/json"
	"configcenter/src/common/types"
	"github.com/stretchr/testify/assert"
)

var rd *RegDiscover
var servInfo types.ServerInfo

const (
	HostServerName = "host-server-test"
	TopoServerName = "topo-server-test"
	EventServerName = "event-server-test"
	ProcServerName = "proc-server-test"
)

func init(){
	rd = &RegDiscover{}
	servInfo = types.ServerInfo{
		Port:123,
		Scheme:"https",
		Version:"v1",
		Pid:123,
	}
}

func TestDiscoverHostServ(t *testing.T){
	hostS := servInfo
	hostS.HostName = HostServerName
	hostS.IP = HostServerName
	hostInfo := &types.HostServerInfo{
		ServerInfo:hostS,
	}

	by,_ := json.Marshal(hostInfo)
	err := rd.discoverHostServ([]string{string(by)})
	assert.Nil(t,err)

	host,err := rd.GetHostServ()
	assert.Nil(t,err)
	assert.Contains(t,host,HostServerName)
}

func TestDiscoverTopoServ(t *testing.T){
	topoS := servInfo
	topoS.HostName = TopoServerName
	topoS.IP = TopoServerName
	topoInfo := &types.TopoServInfo{
		ServerInfo:topoS,
	}

	by,_ := json.Marshal(topoInfo)
	err := rd.discoverTopoServ([]string{string(by)})
	assert.Nil(t,err)

	topo,err := rd.GetTopoServ()
	assert.Nil(t,err)
	assert.Contains(t,topo,TopoServerName)
}

func TestDiscoverProcServ(t *testing.T){
	procS := servInfo
	procS.HostName = ProcServerName
	procS.IP = ProcServerName
	procInfo := &types.ProcServInfo{
		ServerInfo:procS,
	}

	by,_ := json.Marshal(procInfo)
	err := rd.discoverProcServ([]string{string(by)})
	assert.Nil(t,err)

	proc,err := rd.GetProcServ()
	assert.Nil(t,err)
	assert.Contains(t,proc,ProcServerName)
}

func TestDiscoverEventServ(t *testing.T){
	eventS := servInfo
	eventS.HostName = EventServerName
	eventS.IP = EventServerName
	eventInfo := &types.EventServInfo{
		ServerInfo:eventS,
	}

	by,_ := json.Marshal(eventInfo)
	err := rd.discoverEventServ([]string{string(by)})
	assert.Nil(t,err)

	event,err := rd.GetEventServ()
	assert.Nil(t,err)
	assert.Contains(t,event,EventServerName)
}

