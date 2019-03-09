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
	"strconv"

	"configcenter/src/common/core/cc/config"
	"configcenter/src/storage/dal/mongo"
	"configcenter/src/storage/dal/redis"

	"github.com/spf13/pflag"
)

//ServerOption define option of server in flags
type ServerOption struct {
	ServConf *config.CCAPIConfig
}

//NewServerOption create a ServerOption object
func NewServerOption() *ServerOption {
	s := ServerOption{
		ServConf: config.NewCCAPIConfig(),
	}

	return &s
}

//AddFlags add flags
func (s *ServerOption) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&s.ServConf.AddrPort, "addrport", "127.0.0.1:50010", "The ip address and port for the serve on")
	fs.StringVar(&s.ServConf.RegDiscover, "regdiscv", "", "hosts of register and discover server. e.g: 127.0.0.1:2181")
	fs.StringVar(&s.ServConf.ExConfig, "config", "", "The config path. e.g conf/api.conf")
}

// Config transaction server config structure
type Config struct {
	MongoDB     mongo.Config
	Redis       redis.Config
	Transaction TransactionConfig
}

// TransactionConfig transaction config structure
type TransactionConfig struct {
	Enable                    string
	TransactionLifetimeSecond string
}

// IsTransactionEnable check it the transaction is enable
func (c TransactionConfig) IsTransactionEnable() bool {
	switch c.Enable {
	case "1", "true", "":
		return true
	default:
		return false
	}
}

// GetTransactionLifetimeSecond return the lifecycle of the transaction
func (c TransactionConfig) GetTransactionLifetimeSecond() int {
	sec, _ := strconv.Atoi(c.TransactionLifetimeSecond)
	if sec == 0 {
		sec = 60
	}
	return sec
}
