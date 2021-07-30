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

package blueking

import (
	"fmt"
	"os/exec"
	"strconv"

	"configcenter/src/common/blog"
	"configcenter/src/thirdparty/monitor/config"
)

type GseCmdline struct{}

// NewGseCmdline new a GseCmdline instance
func NewGseCmdline() (*GseCmdline, error) {
	gseCmdline := new(GseCmdline)
	if !gseCmdline.isAvailable() {
		return nil, fmt.Errorf("gsecmdline is not available")
	}
	return gseCmdline, nil
}

// Report send the data to bk-monitor
func (g *GseCmdline) Report(data string) error {
	cmd := exec.Command(config.MonitorCfg.GsecmdlinePath, "-d", strconv.FormatInt(config.MonitorCfg.DataID, 10),
		"-D", "-s", data, "-S", config.MonitorCfg.DomainSocketPath)
	if err := cmd.Start(); err != nil {
		blog.Errorf("Report failed, command err:%v, data:%s", err, data)
		return fmt.Errorf("GseCmdline Report failed")
	}
	if err := cmd.Wait(); err != nil {
		blog.Errorf("Report failed, wait error:%v, data:%s", err, data)
		return fmt.Errorf("GseCmdline Report failed")
	}

	return nil
}

func (g *GseCmdline) isAvailable() bool {
	cmd := exec.Command(config.MonitorCfg.GsecmdlinePath, "-v")
	tryCnt := 3
	for i := 0; i < tryCnt; i++ {
		if err := cmd.Run(); err != nil {
			blog.Errorf("check gsecmdline is available failed, command err: %v", err)
			continue
		}
		return true
	}

	return false
}
