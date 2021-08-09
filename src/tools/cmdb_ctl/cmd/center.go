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
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"configcenter/src/tools/cmdb_ctl/app/config"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(NewCenterCommand())
}

type centerService struct {
	service *config.Service
	key     string
}

func newCenterService(addr string, key string) (*centerService, error) {
	if key == "" {
		return nil, errors.New("key must be set")
	}
	service, err := config.NewService(addr)
	if err != nil {
		return nil, err
	}
	return &centerService{
		service: service,
		key:     key,
	}, nil
}

type centerConf struct {
	key string
}

func (c *centerConf) addFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVar(&c.key, "key", "", "the center resource key")
}

func NewCenterCommand() *cobra.Command {
	conf := new(centerConf)

	cmd := &cobra.Command{
		Use:   "center",
		Short: "center operations",
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}

	subCmds := make([]*cobra.Command, 0)

	subCmds = append(subCmds, &cobra.Command{
		Use:   "get",
		Short: "get value of specified node",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGetCmd(conf)
		},
	})

	subCmds = append(subCmds, &cobra.Command{
		Use:   "del",
		Short: "delete specified node",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDelCmd(conf)
		},
	})

	value := new(string)
	setCmd := &cobra.Command{
		Use:   "set",
		Short: "set value of specified node",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSetCmd(conf, *value)
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

func runGetCmd(c *centerConf) error {
	srv, err := newCenterService(config.Conf.Addr, c.key)
	if err != nil {
		return err
	}
	data, err := srv.service.Cli.Get(srv.key)
	if err != nil {
		return err
	}
	var pretty bytes.Buffer
	err = json.Indent(&pretty, []byte(data), "", "\t")
	if err != nil {
		fmt.Fprintln(os.Stdout, data)
		return nil
	}
	fmt.Fprintln(os.Stdout, pretty.String())
	return nil
}

func runDelCmd(c *centerConf) error {
	srv, err := newCenterService(config.Conf.Addr, c.key)
	if err != nil {
		return err
	}
	return srv.service.Cli.Delete(srv.key)
}

func runSetCmd(c *centerConf, value string) error {
	srv, err := newCenterService(config.Conf.Addr, c.key)
	if err != nil {
		return err
	}
	return srv.service.Cli.Put(srv.key, value)
}
