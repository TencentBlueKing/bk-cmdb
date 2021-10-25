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

package setting

import (
	"fmt"
	"net/http"

	"configcenter/src/common/blog"
)

type ActionType string

const (
	Get    ActionType = "get"
	GetAll ActionType = "getAll"
	Add    ActionType = "add"
	Delete ActionType = "delete"
	Update ActionType = "update"
)

type OperationType string

const (
	SettingsLog  OperationType = "log"
	SettingsHelp OperationType = "help"
)

const (
	SettingsAction    = "action"
	SettingsOperation = "op"
	SettingsLogLevel  = "v"
)

// Service dynamically adjust the structure of the runtime configuration
type Service struct {
	log *LogService
}

// NewService new a service struct
func NewService() *Service {
	return &Service{
		log: NewLogService(),
	}
}

// ServeHTTP method of dynamically adjusting runtime configuration
func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	op := r.FormValue(SettingsOperation)

	switch OperationType(op) {
	case SettingsHelp:
		fmt.Fprintf(w, "%s", GetHelp())
		return

	case SettingsLog:
		loglevel := r.FormValue(SettingsLogLevel)
		if err := s.log.ChangeLogLevel(loglevel); err != nil {
			blog.Errorf("update log level failed err: %v", err)
			fmt.Fprintf(w, "update log level failed err: %v", err)
			return
		}
		fmt.Fprintln(w, "success!")

	default:
		blog.Errorf("adjust operation error, can't find the relevant operation to adjust.")
		fmt.Fprintln(w, "adjust operation error, can't find the relevant operation to adjust!")
	}
}
