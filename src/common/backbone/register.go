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

package backbone

import (
	"encoding/json"
	"errors"

	"configcenter/src/common/backbone/service_mange/zk"
	"configcenter/src/common/registerdiscover"
	"configcenter/src/common/types"
)

// ServiceRegisterInterface TODO
type ServiceRegisterInterface interface {
	// Ping to ping server
	Ping() error
	// Register local server info, it can only be called for once.
	Register(path string, c types.ServerInfo) error
	// Cancel to stop server register and discover
	Cancel()
	// ClearRegisterPath to delete server register path from zk
	ClearRegisterPath() error
}

// NewServiceRegister TODO
func NewServiceRegister(client *zk.ZkClient) (ServiceRegisterInterface, error) {
	s := new(serviceRegister)
	s.client = registerdiscover.NewRegDiscoverEx(client)
	return s, nil
}

type serviceRegister struct {
	client *registerdiscover.RegDiscover
}

// Register TODO
func (s *serviceRegister) Register(path string, c types.ServerInfo) error {
	if c.RegisterIP == "0.0.0.0" {
		return errors.New("register ip can not be 0.0.0.0")
	}

	js, err := json.Marshal(c)
	if err != nil {
		return err
	}

	return s.client.RegisterAndWatchService(path, js)
}

// Ping to ping server
func (s *serviceRegister) Ping() error {
	return s.client.Ping()
}

// Cancel to stop server register and discover
func (s *serviceRegister) Cancel() {
	s.client.Cancel()
}

// ClearRegisterPath to delete server register path from zk
func (s *serviceRegister) ClearRegisterPath() error {
	return s.client.ClearRegisterPath()
}
