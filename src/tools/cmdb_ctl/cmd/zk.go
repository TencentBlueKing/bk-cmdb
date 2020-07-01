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
	"bytes"
	"configcenter/src/tools/cmdb_ctl/app/config"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(NewZkCommand())
}

type zkConf struct {
	path string
}

func NewZkCommand() *cobra.Command {
	conf := new(zkConf)

	cmd := &cobra.Command{
		Use:   "zk",
		Short: "zookeeper operations",
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}

	subCmds := make([]*cobra.Command, 0)

	subCmds = append(subCmds, &cobra.Command{
		Use:   "ls",
		Short: "list children of specified zookeeper node",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runZkLsCmd(conf)
		},
	})

	subCmds = append(subCmds, &cobra.Command{
		Use:   "get",
		Short: "get value of specified zookeeper node",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runZkGetCmd(conf)
		},
	})

	subCmds = append(subCmds, &cobra.Command{
		Use:   "del",
		Short: "delete specified zookeeper node",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runZkDelCmd(conf)
		},
	})

	value := new(string)
	setCmd := &cobra.Command{
		Use:   "set",
		Short: "set value of specified zookeeper node",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runZkSetCmd(conf, *value)
		},
	}
	setCmd.Flags().StringVar(value, "value", "", "the value to be set")
	_ = setCmd.MarkFlagRequired("value")
	subCmds = append(subCmds, setCmd)

	for _, subCmd := range subCmds {
		cmd.AddCommand(subCmd)
	}
	conf.addFlags(cmd)

	return cmd
}

func (c *zkConf) addFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVar(&c.path, "zk-path", "", "the zookeeper  resource path")
}

type zkService struct {
	service *config.Service
	path    string
}

func newZkService(zkaddr string, path string) (*zkService, error) {
	if path == "" {
		return nil, errors.New("zk-path must be set")
	}
	service, err := config.NewZkService(zkaddr)
	if err != nil {
		return nil, err
	}
	return &zkService{
		service: service,
		path:    path,
	}, nil
}

func runZkLsCmd(c *zkConf) error {
	srv, err := newZkService(config.Conf.ZkAddr, c.path)
	if err != nil {
		return err
	}
	children, err := srv.service.ZkCli.GetChildren(srv.path)
	if err != nil {
		return err
	}
	for _, child := range children {
		fmt.Fprintf(os.Stdout, "%s\n", child)
	}
	return nil
}

func runZkGetCmd(c *zkConf) error {
	srv, err := newZkService(config.Conf.ZkAddr, c.path)
	if err != nil {
		return err
	}
	data, err := srv.service.ZkCli.Get(srv.path)
	if err != nil {
		return err
	}
	var pretty bytes.Buffer
	err = json.Indent(&pretty, []byte(data), "", "\t")
	if err != nil {
		return err
	}
	fmt.Fprintln(os.Stdout, pretty.String())
	return nil
}

func runZkDelCmd(c *zkConf) error {
	srv, err := newZkService(config.Conf.ZkAddr, c.path)
	if err != nil {
		return err
	}
	return srv.service.ZkCli.Del(srv.path, -1)
}

func runZkSetCmd(c *zkConf, value string) error {
	srv, err := newZkService(config.Conf.ZkAddr, c.path)
	if err != nil {
		return err
	}
	return srv.service.ZkCli.Set(srv.path, value, -1)
}
