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

package monitor

import (
	"os/exec"
	"strconv"

	"configcenter/src/common/blog"
)

type GseCmdline struct{}

func NewGseCmdline() *GseCmdline {
	return new(GseCmdline)
}

// IsAvailable judge whether the GseCmdline is available
func (g *GseCmdline) IsAvailable() bool {
	cmd := exec.Command("gsecmdline", "-h")
	tryCnt := 3
	for i := 0; i < tryCnt; i++ {
		if err := cmd.Start(); err != nil {
			blog.Errorf("IsAvailable failed, command err:%v", err)
		}
		if err := cmd.Wait(); err != nil {
			blog.Errorf("IsAvailable failed, wait error:%v", err)
		} else {
			return true
		}
	}

	return false
}

// Report send the data to bk-monitor
func (g *GseCmdline) Report(data string) error {
	cmd := exec.Command("gsecmdline", "-d", strconv.FormatInt(MonitorCfg.DataID, 10), "-j", data)
	if err := cmd.Start(); err != nil {
		blog.Errorf("Report failed, command err:%v, data:%s", err, data)
		return err
	}
	if err := cmd.Wait(); err != nil {
		blog.Errorf("Report failed, wait error:%v, data:%s", err, data)
		return err
	}

	return nil
}
