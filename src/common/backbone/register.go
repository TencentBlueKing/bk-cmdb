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

	"configcenter/src/common/registerdiscover"
	"configcenter/src/common/types"
	"configcenter/src/framework/core/errors"
)

type ServiceRegisterInterface interface {
	// Ping ping register and discover to verify accessibility
	Ping() error
	// Register add service info to register and discover
	Register(path string, c types.ServerInfo) error
	// Unregister delete service register key from register and discover
	Unregister() error
	// Cancel stop service register and discover
	Cancel()
}

// NewServiceRegister create a service register object
func NewServiceRegister(rd *registerdiscover.RegDiscv) (ServiceRegisterInterface, error) {
	return &serviceRegister{rd: rd}, nil
}

type serviceRegister struct {
	rd     *registerdiscover.RegDiscv
	regKey string
}

// Ping ping register and discover to verify accessibility
func (s *serviceRegister) Ping() error {
	return s.rd.Ping()
}

// Register add service info to register and discover
func (s *serviceRegister) Register(path string, c types.ServerInfo) error {
	if c.RegisterIP == "0.0.0.0" {
		return errors.New("register ip can not be 0.0.0.0")
	}

	js, err := json.Marshal(c)
	if err != nil {
		return err
	}

	s.regKey = path

	return s.rd.RegisterAndKeepAlive(path, string(js))
}

// Unregister delete service register key from register and discover
func (s *serviceRegister) Unregister() error {
	return s.rd.Delete(s.regKey)
}

// Cancel stop service register and discover
func (s *serviceRegister) Cancel() {
	s.rd.Cancel()
}
