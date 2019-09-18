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

package datacollection

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime/debug"

	"configcenter/src/common/blog"
	"configcenter/src/common/version"
)

// Manager manage the porter goroutine
type Manager struct {
	porterC chan Porter
	porters map[string]Porter
}

func NewManager() *Manager {
	mgr := &Manager{
		porterC: make(chan Porter, 1),
		porters: map[string]Porter{},
	}
	go mgr.run()
	if version.CCRunMode != version.CCRunModeProduct {
		go mgr.mockServer()
	}
	return mgr
}

func (m *Manager) run() error {
	for porter := range m.porterC {
		m.porters[porter.Name()] = porter
		go m.porterLoop(porter)
	}

	return nil
}

type mockMesg struct {
	Name    string `json:"name"`
	Message string `json:"mesg"`
}

func (m *Manager) mockServer() {
	mockServer := http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		mockMSG := mockMesg{}
		err := json.NewDecoder(req.Body).Decode(&mockMSG)
		if err != nil {
			fmt.Fprintf(resp, "decode message error: %v", err)
			resp.WriteHeader(400)
			return
		}
		if porter, ok := m.porters[mockMSG.Name]; ok {
			if err := porter.Mock(mockMSG.Message); err != nil {
				fmt.Fprintf(resp, "mock failed: %v", err)
				resp.WriteHeader(400)
				return
			}
		} else {
			fmt.Fprintf(resp, "unknow porter: %v", mockMSG.Name)
			resp.WriteHeader(400)
			return
		}
	})
	if err := http.ListenAndServe("127.0.0.1:12140", mockServer); err != nil {
		blog.Warnf("start mock server failed: %v", err)
	}

}

func (m *Manager) AddPorter(p Porter) {
	m.porterC <- p
}

func (m *Manager) porterLoop(p Porter) {
	for {
		m.runPorter(p)
	}
}

func (m *Manager) runPorter(p Porter) {
	defer func() {
		if sysErr := recover(); sysErr != nil {
			blog.Errorf("[manager] porter [%s] panic by: %v, stack:\n %s", p.Name(), sysErr, debug.Stack())
		}
	}()
	if err := p.Run(); err != nil {
		blog.Errorf("[manager] porter [%s] return an error %v", p.Name(), err)
	}
}
