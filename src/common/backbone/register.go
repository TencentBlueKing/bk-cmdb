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

package backbone

import (
	"encoding/json"

	"configcenter/src/common/backbone/service_mange/zk"
	"configcenter/src/common/registerdiscover"
	"configcenter/src/common/types"
	"configcenter/src/framework/core/errors"
)

type ServiceRegisterInterface interface {
	// Ping to ping server
	Ping() error
	// register local server info, it can only be called for once.
	Register(path string, c types.ServerInfo) error
	// Cancel to stop server register and discover
	Cancel()
	// ClearRegisterPath to delete server register path from zk
	ClearRegisterPath() error
}

func NewServiceRegister(client *zk.ZkClient) (ServiceRegisterInterface, error) {
	s := new(serviceRegister)
	s.client = registerdiscover.NewRegDiscoverEx(client)
	return s, nil
}

type serviceRegister struct {
	client *registerdiscover.RegDiscover
}

func (s *serviceRegister) Register(path string, c types.ServerInfo) error {
	if c.IP == "0.0.0.0" {
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
