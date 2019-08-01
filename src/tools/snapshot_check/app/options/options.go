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

package options

import (
	"github.com/spf13/pflag"
)

// ServerOption define option of server in flags
type ServerOption struct {
	ServConf Config
}

// NewServerOption create a ServerOption object
func NewServerOption() *ServerOption {
	s := ServerOption{
		ServConf: Config{},
	}

	return &s
}

// AddFlags add flags
func (s *ServerOption) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&s.ServConf.RegDiscover, "regdiscv", "", "hosts of register and discover server. e.g: 127.0.0.1:2181")
	s.ServConf.DefaultAppID = *fs.Int("appID", 2, "blueking business id. e.g: 2")
	s.ServConf.TriggerInterval = *fs.Int("interval", 10, "blueking business id. e.g: 2")
}

type Config struct {
	RegDiscover     string
	DefaultAppID    int
	TriggerInterval int
}
