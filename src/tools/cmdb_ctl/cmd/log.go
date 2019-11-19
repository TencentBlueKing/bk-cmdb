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
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"configcenter/src/common/types"
	"configcenter/src/tools/cmdb_ctl/app/config"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(NewLogCommand())
}

type logConf struct {
	v        string
	def      bool
	addrPort string
}

func NewLogCommand() *cobra.Command {
	conf := new(logConf)

	cmd := &cobra.Command{
		Use:   "log",
		Short: "set log level",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runLog(conf)
		},
	}

	conf.addFlags(cmd)

	return cmd
}

func (c *logConf) addFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&c.v, "set-v", "", "set log level for V logs")
	cmd.Flags().BoolVar(&c.def, "set-default", false, "set log level to default value")
	cmd.Flags().StringVar(&c.addrPort, "addrport", "", "the ip address and port for the hosts to apply command, separated by comma")
}

type logService struct {
	service  *config.Service
	addrport []string
}

func newLogService(zkaddr string, addrport string) (*logService, error) {
	if addrport == "" {
		return nil, errors.New("addrport must set via flag or environment variable")
	}
	service, err := config.NewZkService(zkaddr)
	if err != nil {
		return nil, err
	}
	return &logService{
		service:  service,
		addrport: strings.Split(addrport, ","),
	}, nil
}

func runLog(c *logConf) error {
	if c.v == "" && !c.def {
		return fmt.Errorf("use set-v or set-default flag to set log level")
	}
	if c.v != "" && c.def {
		return fmt.Errorf("can't set log level to v and default at the same time")
	}

	srv, err := newLogService(config.Conf.ZkAddr, c.addrPort)
	if err != nil {
		return err
	}

	if c.v != "" {
		v, err := strconv.ParseInt(c.v, 0, 32)
		if err != nil {
			return err
		}
		return srv.setV(int32(v))
	}
	if c.def {
		return srv.setDefault()
	}
	return nil
}

func (s *logService) setV(v int32) error {
	for _, addr := range s.addrport {
		if err := s.service.ZkCli.Ping(); err != nil {
			if err = s.service.ZkCli.Connect(); err != nil {
				return err
			}
		}
		logVPath := fmt.Sprintf("%s/%s/%s/v", types.CC_SERVNOTICE_BASEPATH, "log", addr)
		logVData, err := s.service.ZkCli.Get(logVPath)
		if err != nil {
			return err
		}
		data := make(map[string]int32)
		err = json.Unmarshal([]byte(logVData), &data)
		if err != nil {
			return err
		}
		data["v"] = v
		dat, err := json.Marshal(data)
		if err != nil {
			return err
		}
		if err = s.service.ZkCli.Update(logVPath, string(dat)); err != nil {
			return err
		}
	}
	return nil
}

func (s *logService) setDefault() error {
	for _, addr := range s.addrport {
		if err := s.service.ZkCli.Ping(); err != nil {
			if err = s.service.ZkCli.Connect(); err != nil {
				return err
			}
		}
		logVPath := fmt.Sprintf("%s/%s/%s/v", types.CC_SERVNOTICE_BASEPATH, "log", addr)
		logVData, err := s.service.ZkCli.Get(logVPath)
		if err != nil {
			return err
		}
		data := make(map[string]int32)
		err = json.Unmarshal([]byte(logVData), &data)
		if err != nil {
			return err
		}
		data["v"] = data["defaultV"]
		dat, err := json.Marshal(data)
		if err != nil {
			return err
		}
		if err = s.service.ZkCli.Update(logVPath, string(dat)); err != nil {
			return err
		}
	}
	return nil
}
