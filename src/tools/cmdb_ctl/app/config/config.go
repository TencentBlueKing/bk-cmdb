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

package config

import (
	"errors"
	"os"
	"strings"

	"configcenter/src/common/zkclient"

	"github.com/spf13/cobra"
)

var Conf *Config

type Config struct {
	ZkConf   *zkclient.ZkConf
	AddrPort string
}

// AddFlags add flags
func (c *Config) AddFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVar(&c.ZkConf.ZkAddr, "zkaddr", os.Getenv("ZK_ADDR"), "the ip address and port for the zookeeper hosts, separated by comma, corresponding environment variable is ZK_ADDR")
	cmd.PersistentFlags().StringVar(&c.ZkConf.ZkUser, "zkuser", os.Getenv("ZK_USER"), "the zookeeper auth user, corresponding environment variable is ZK_USER")
	cmd.PersistentFlags().StringVar(&c.ZkConf.ZkUser, "zkpwd", os.Getenv("ZK_PWD"), "the zookeeper auth password, corresponding environment variable is ZK_PWD")
	cmd.PersistentFlags().StringVar(&c.AddrPort, "addrport", os.Getenv("ADDR_PORT"), "the ip address and port for the hosts to apply command, separated by comma, corresponding environment variable is ADDR_PORT")
}

type Service struct {
	ZkCli    *zkclient.ZkClient
	Addrport []string
}

func NewService(zkConf *zkclient.ZkConf, addrport string) (*Service, error) {
	if zkConf.ZkAddr == "" {
		return nil, errors.New("zkaddr must set via flag or environment variable")
	}
	if zkConf.ZkUser == "" {
		return nil, errors.New("zkuser must set via flag or environment variable")
	}
	if zkConf.ZkPwd == "" {
		return nil, errors.New("zkpwd must set via flag or environment variable")
	}
	service := &Service{
		ZkCli:    zkclient.NewZkClient(zkConf),
		Addrport: strings.Split(addrport, ","),
	}
	if err := service.ZkCli.Connect(); err != nil {
		return nil, err
	}
	return service, nil
}
