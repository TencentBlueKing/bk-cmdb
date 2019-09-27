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

package taskserver

import (
	"fmt"
	"sync"

	"configcenter/src/apimachinery/rest"
	"configcenter/src/apimachinery/taskserver/queue"
	"configcenter/src/apimachinery/taskserver/task"
	"configcenter/src/apimachinery/util"

	"configcenter/src/apimachinery/flowctrl"
	taskUtil "configcenter/src/apimachinery/taskserver/util"
)

type TaskServerClientInterface interface {
	Task() task.TaskClientInterface
	Queue(flag string) queue.TaskQueueClientInterface
}

func NewProcServerClientInterface(c *util.Capability, version string) TaskServerClientInterface {
	base := fmt.Sprintf("/task/%s", version)
	return &taskServer{
		client:      rest.NewRESTClient(c, base),
		capability:  c,
		queueClient: make(map[string]queue.TaskQueueClientInterface, 0),
	}
}

type taskServer struct {
	client rest.ClientInterface
	sync.RWMutex
	queueClient map[string]queue.TaskQueueClientInterface
	capability  *util.Capability
}

func (ts *taskServer) Task() task.TaskClientInterface {
	return task.NewTaskClientInterface(ts.client)
}

func (ts *taskServer) Queue(flag string) queue.TaskQueueClientInterface {
	ts.RLock()
	srv, ok := ts.queueClient[flag]
	ts.RUnlock()
	if nil == srv || !ok {
		ts.Lock()
		ts.queueClient[flag] = queue.NewSychronizeClientInterface(ts.getSrvClent(flag))
		srv = ts.queueClient[flag]
		ts.Unlock()
	}
	return srv

}

func (ts *taskServer) getSrvClent(flag string) rest.ClientInterface {
	flowcontrol := flowctrl.NewRateLimiter(ts.capability.Throttle.QPS(), ts.capability.Throttle.Burst())
	config := taskUtil.NewSyncrhonizeConfig(flag)

	capability := &util.Capability{
		Client:   ts.capability.Client,
		Discover: config,
		Throttle: flowcontrol,
	}

	return rest.NewRESTClient(capability, "/")
}
