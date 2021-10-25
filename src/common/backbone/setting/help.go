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

import "fmt"

type help struct {
	op      OperationType
	url     string
	explain string
}

var helps []*help

// AddHelp add help
func AddHelp(help *help) {
	helps = append(helps, help)
}

// GetHelp get help message
func GetHelp() string {
	var helpMessage string
	for _, help := range helps {
		helpMessage += fmt.Sprintf("%s: %s  %s\n", help.op, help.url, help.explain)
	}
	return helpMessage
}
