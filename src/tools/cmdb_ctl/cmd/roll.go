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

package cmd

import (
	"os"
	"strings"

	"configcenter/src/common/zkclient"

	"github.com/spf13/cobra"
)

type rollConf struct {
	regdiscv string
	addrport string
	showNodes        bool
	showStatus      bool
}

func NewRollCommand() *cobra.Command {
	conf := new(rollConf)

	cmd := &cobra.Command{
		Use:   "roll",
		Short: "rolling update",
		Long:`roll is for rolling update of cmdb processes, update will be applied when called without show-nodes or show-status flags
show-nodes will show cmdb nodes registered in zookeeper specified by regdiscv that can be picked to apply update
show-status will show update status of processes specified by addrport
regdiscv and addrport will be stored in environment variable for later call
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRoll(conf)
		},
	}

	conf.addFlags(cmd)

	return cmd
}

func (c *rollConf) addFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&c.regdiscv, "regdiscv", os.Getenv("ROLL_REG_DISCV"), "hosts of register and discover server, separated by ,")
	cmd.Flags().StringVar(&c.addrport, "addrport", os.Getenv("ROLL_ADDR_PORT"), "the ip address and port for the processes to update, separated by ,")
	cmd.Flags().BoolVar(&c.showNodes, "show-nodes", false, "show zookeeper nodes")
	cmd.Flags().BoolVar(&c.showStatus, "show-status", false, "show rolling update status")
}

type rollService struct {
	zkCli    *zkclient.ZkClient
	addrport []string
}

func newRollService(regdiscv string, addrport string) *rollService {
	return &rollService{
		zkCli:    zkclient.NewZkClient(strings.Split(regdiscv, ",")),
		addrport: strings.Split(addrport, ","),
	}
}

func runRoll(c *rollConf) error {
	// store these variables as environment variables so that can be reused for future call
	if err := os.Setenv("ROLL_REG_DISCV", c.regdiscv); err != nil {
		return err
	}
	if err := os.Setenv("ROLL_ADDR_PORT", c.addrport); err != nil {
		return err
	}
	return nil
}
