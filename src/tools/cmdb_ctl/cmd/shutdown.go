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
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(NewShutdownCommand())
}

type shutdownConf struct {
	showPids bool
	pids     string
}

func NewShutdownCommand() *cobra.Command {
	conf := new(shutdownConf)

	cmd := &cobra.Command{
		Use:   "shutdown",
		Short: "graceful shutdown",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runShutdown(conf)
		},
	}

	conf.addFlags(cmd)

	return cmd
}

func (c *shutdownConf) addFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&c.pids, "pids", "", "the processes to be shutdown, separated by comma")
	cmd.Flags().BoolVar(&c.showPids, "show-pids", false, "show cmdb processes pid")
}

type shutdownService struct {
	pids []int
}

func newShutdownService() *shutdownService {
	return new(shutdownService)
}

func runShutdown(c *shutdownConf) error {
	srv := newShutdownService()
	if c.pids != "" {
		pids := strings.Split(c.pids, ",")
		return srv.shutdown(pids)
	}
	if c.showPids {
		return srv.showPids()
	}
	return nil
}

func (s *shutdownService) shutdown(pids []string) error {
	for _, pid := range pids {
		cmd := exec.Command("kill", fmt.Sprintf("-%d", int(syscall.SIGTERM)), pid)
		if err := cmd.Run(); err != nil {
			return err
		}
	}
	return nil
}

func (s *shutdownService) showPids() error {
	fmt.Println("Pid\tProcess\tAddrPort")
	cmd := exec.Command("/bin/sh", "-c", `ps -ef | grep cmdb_ | grep -v "grep" | awk '{print $2"\t"$8"\t"$9}'`)
	cmd.Stdout = os.Stdout
	return cmd.Run()
}
