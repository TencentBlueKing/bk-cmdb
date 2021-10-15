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

package tool

import (
	"strconv"

	"configcenter/src/common/blog"
)

// LogDefaultFlag log default flag
const LogDefaultFlag = "default"

// logService the structure responsible for log adjustment
type LogService struct {
	defaultV int32
}

// NewLogService new logService
func NewLogService() *LogService {
	return &LogService{
		defaultV: blog.GetV(),
	}
}

// ChangeLogLevel change the log level of the service
func (s *LogService) ChangeLogLevel(val string) error {
	if val == LogDefaultFlag {
		return s.setDefault()
	}
	logLevel, err := strconv.ParseInt(val, 0, 32)
	if err != nil {
		return err
	}
	return s.setV(int32(logLevel))
}

// setV set the log level to a specific value
func (s *LogService) setV(v int32) error {
	blog.SetV(v)
	return nil
}

// setDefault set the log level to the default value
func (s *LogService) setDefault() error {
	blog.SetV(s.defaultV)
	return nil
}
