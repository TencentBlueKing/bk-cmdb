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

package collections

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"configcenter/src/common/blog"
	"configcenter/src/common/version"
)

const (
	// defaultPortersChanTimeout is default porters channel timeout.
	defaultPortersChanTimeout = time.Second

	// defaultMockServerEndpoint is default endpoint of local mock server.
	defaultMockServerEndpoint = "127.0.0.1:12140"
)

// local mock message.
type mockMessage struct {
	// Name is porter name(cmdline collector name).
	Name string `json:"name"`

	// Message is mock message content.
	Message string `json:"mesg"`
}

// PorterManager manages the collection porters.
type PorterManager struct {
	// porters saves all runtime porters, porter name -> porter instance.
	porters map[string]Porter

	// portersChan is used for add a new porter instance when setups the manager.
	portersChan chan Porter
}

// NewPorterManager creates a new PorterManager object.
func NewPorterManager() *PorterManager {
	return &PorterManager{
		porters:     make(map[string]Porter),
		portersChan: make(chan Porter),
	}
}

// runMockServer run a local mock server.
func (mgr *PorterManager) runMockServer() {
	mockServer := http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		// decode mock request, just get the porter name(cmdline collector name),
		// and use the mock message in target analyzer.
		mock := mockMessage{}

		if err := json.NewDecoder(req.Body).Decode(&mock); err != nil {
			fmt.Fprintf(resp, "decode mock message error: %+v", err)
			resp.WriteHeader(http.StatusBadRequest)
			return
		}

		if porter, ok := mgr.porters[mock.Name]; ok {
			if err := porter.Mock(); err != nil {
				fmt.Fprintf(resp, "mock failed, %+v", err)
				resp.WriteHeader(http.StatusBadRequest)
				return
			}
		} else {
			fmt.Fprintf(resp, "unknow porter: %v", mock.Name)
			resp.WriteHeader(http.StatusBadRequest)
			return
		}
	})

	if err := http.ListenAndServe(defaultMockServerEndpoint, mockServer); err != nil {
		blog.Warnf("PorterManager| run local mock server failed, %+v", err)
	}
}

// handlePorters handles porters setup events.
func (mgr *PorterManager) handlePorters() {
	for porter := range mgr.portersChan {
		// porter is Porter interface, eg SimplePorter instance point.
		if _, isExist := mgr.porters[porter.Name()]; !isExist {
			// new porter, add and run it.
			mgr.porters[porter.Name()] = porter
			go porter.Run()
		}
	}
}

// AddPorter adds and runs a new porter.
func (mgr *PorterManager) AddPorter(p Porter) error {
	select {
	case mgr.portersChan <- p:

	case <-time.After(defaultPortersChanTimeout):
		return fmt.Errorf("add to porters channel timeout")
	}
	return nil
}

// Run runs the new PorterManager.
func (mgr *PorterManager) Run() error {
	if version.CCRunMode != version.CCRunModeProduct {
		go mgr.runMockServer()
	}

	// handle porters.
	mgr.handlePorters()

	return nil
}
