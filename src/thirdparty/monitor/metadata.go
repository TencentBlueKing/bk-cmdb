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

// MonitorData is the common monitor data struct
type MonitorData struct {
	// RequestID used to trace the log
	RequestID string `json:"request_id"`
	// RequestID used to record the event name
	EventName string `json:"event_name"`
	// EventContent used to record the event content
	EventContent string `json:"event_content"`
	// module name, like coreservice, hostserver
	Module string `json:"module"`
}
