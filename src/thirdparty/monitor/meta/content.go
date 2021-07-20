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

package meta

// Content is the data processed by monitor
type Content interface {
	ContentName() string
}

// Alarm a Content implement used to alarm
type Alarm struct {
	// RequestID used to trace the log
	RequestID string `json:"request_id"`
	// Type used to classify alarm
	Type MonitorType `json:"type"`
	// Detail is the alarm detail
	Detail string `json:"detail"`
	// module name, like coreservice, hostserver
	Module string `json:"module"`
	// Dimension is self-defined kv info
	Dimension map[string]string `json:"dimension"`
}

func (c *Alarm) ContentName() string {
	return "alarm"
}
