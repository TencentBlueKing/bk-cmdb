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

package confregdiscover

import (
	"time"
)

//DiscoverEvent if servers changed, will create a discover event
type DiscoverEvent struct { //
	Err  error
	Key  string
	Data []byte
}

// ConfRegDiscover is config register and discover
type ConfRegDiscover struct {
	confRD ConfRegDiscvIf
}

// NewConfRegDiscover used to create a object of ConfRegDiscover
// session timeout default 60 second
func NewConfRegDiscover(serv string) *ConfRegDiscover {
	confRD := &ConfRegDiscover{
        confRD: nil,
	}

	confRD.confRD = ConfRegDiscvIf(NewZkRegDiscover(serv, time.Second*60))

	return confRD
}

// NewConfRegDiscoverWithTimeOut used to create a object
func NewConfRegDiscoverWithTimeOut(serv string, timeOut time.Duration) *ConfRegDiscover {
	confRD := &ConfRegDiscover{
        confRD: nil,
	}

	confRD.confRD = ConfRegDiscvIf(NewZkRegDiscover(serv, timeOut))

	return confRD
}

// Ping to ping server
func (crd *ConfRegDiscover) Ping() error {
	return crd.confRD.Ping()
}

//Start the register and discover service
func (crd *ConfRegDiscover) Start() error {
	return crd.confRD.Start()
}

//Stop the register and discover service
func (crd *ConfRegDiscover) Stop() error {
	return crd.confRD.Stop()
}

//Write the configure data
func (crd *ConfRegDiscover) Write(key string, data []byte) error {
	return crd.confRD.Write(key, data)
}

// Read the configure data
func (crd *ConfRegDiscover) Read(path string) (string, error) {
    return crd.confRD.Read(path)
}

//DiscoverConfig discover the config wether is changed
func (crd *ConfRegDiscover) DiscoverConfig(key string) (<-chan *DiscoverEvent, error) {
	return crd.confRD.Discover(key)
}
