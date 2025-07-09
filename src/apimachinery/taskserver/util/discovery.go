// Package taskserverutil TODO
/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 Tencent. All rights reserved.
 * Licensed under the MIT License (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 * We undertake not to change the open source license (MIT license) applicable
 * to the current version of the project delivered to anyone in the future.
 */
package taskserverutil

import (
	"sync"
)

// TaskQueueConfig TODO
type TaskQueueConfig struct {
	Name  string
	Addrs func() ([]string, error)
}

// TaskQueueConfigServ TODO
type TaskQueueConfigServ struct {
	addrs map[string]func() ([]string, error)
	sync.RWMutex
}

// TaskQueueServDiscoveryInterace TODO
type TaskQueueServDiscoveryInterace interface {
	GetServers() ([]string, error)
	GetServersChan() chan []string
}

var (
	taskQueue = &TaskQueueConfigServ{
		addrs: make(map[string]func() ([]string, error), 0),
	}
)

// NewTaskServerConfigServ TODO
//  NewTaskServerConfigServ
func NewTaskServerConfigServ(srvChan chan TaskQueueConfig) *TaskQueueConfigServ {
	go func() {
		if nil == srvChan {
			return
		}
		for {
			config := <-srvChan
			taskQueue.Lock()
			taskQueue.addrs[config.Name] = config.Addrs

			taskQueue.Unlock()
		}
	}()

	return taskQueue
}

// UpdateTaskServerConfigServ TODO
func UpdateTaskServerConfigServ(name string, f func() ([]string, error)) {
	go func() {

		taskQueue.Lock()

		taskQueue.addrs[name] = f

		taskQueue.Unlock()

	}()

}

type taskQueueConfig struct {
	flag string
}

// NewSyncrhonizeConfig TODO
func NewSyncrhonizeConfig(flag string) TaskQueueServDiscoveryInterace {
	return &taskQueueConfig{
		flag: flag,
	}
}

// GetEsbServDiscoveryInterace TODO
func (s *taskQueueConfig) GetEsbServDiscoveryInterace(flag string) TaskQueueServDiscoveryInterace {
	// mabye will deal some logics about server
	return s
}

// GetServers TODO
func (s *taskQueueConfig) GetServers() ([]string, error) {
	// mabye will deal some logics about server
	taskQueue.RLock()
	defer taskQueue.RUnlock()
	return taskQueue.addrs[s.flag]()
}

// GetServersChan TODO
func (s *taskQueueConfig) GetServersChan() chan []string {
	return nil
}
