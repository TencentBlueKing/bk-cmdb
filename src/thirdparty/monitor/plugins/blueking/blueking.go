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
	"time"

	"configcenter/src/common/blog"
	"configcenter/src/common/json"
	"configcenter/src/thirdparty/monitor/config"
	"configcenter/src/thirdparty/monitor/meta"
)

// bkMonitor is a implementation of monitor Plugin for blueking
type bkMonitor struct {
	gseCmdline *GseCmdline
}

// NewBKmonitor new a bkMonitor instance
func NewBKmonitor() (*bkMonitor, error) {
	gseCmdline, err := NewGseCmdline()
	if err != nil {
		return nil, err
	}
	return &bkMonitor{
		gseCmdline: gseCmdline,
	}, nil
}

// Report is a interface implement for bkMonitor
func (m *bkMonitor) Report(c meta.Content) error {
	if config.MonitorCfg.DataID == 0 {
		blog.Errorf("Report failed, config monitor.dataID is not set")
		return fmt.Errorf("Report failed, config monitor.dataID is not set")
	}

	alarm, ok := c.(*meta.Alarm)
	if !ok {
		blog.Errorf("Report failed, the content typeis not *Alarm, but %T, value:%#v", c, c)
		return fmt.Errorf("Report failed, the content typeis not *Alarm, but %T", c)
	}

	msg, err := m.convertToReportMsg(alarm)
	if err != nil {
		blog.Errorf("Report failed, convertToReportMsg err:%v, data:%#v", err, msg)
		return err
	}

	err = m.gseCmdline.Report(msg)
	if err != nil {
		blog.Errorf("Report failed, gseCmdline Report err:%v, msg:%#v", err, msg)
		return err
	}

	return nil
}

// convertToReportMsg convert data to a msg used by gseCmdline
func (m *bkMonitor) convertToReportMsg(alarm *meta.Alarm) (string, error) {
	event := EventData{
		DataID: config.MonitorCfg.DataID,
		Data: []*EventMsg{
			{
				EventName: string(alarm.Type),
				EventInfo: EventInfo{Content: alarm.Detail},
				Target:    config.MonitorCfg.SourceIP,
				Dimension: map[string]string{
					"module":     alarm.Module,
					"request_id": alarm.RequestID,
				},
				TimeStampMs: time.Now().UnixNano() / 1e6,
			},
		},
	}
	msg, err := json.Marshal(event)
	if err != nil {
		blog.Errorf("marshal error:%v, msg:%s", err, msg)
		return "", fmt.Errorf("convertToReportMsg failed")
	}
	return string(msg), nil
}

// EventData is self-defined event in bk-monitor
type EventData struct {
	DataID int64       `json:"dataid"`
	Data   []*EventMsg `json:"data"`
}

type EventMsg struct {
	EventName   string            `json:"event_name"`
	EventInfo   EventInfo         `json:"event"`
	Target      string            `json:"target"`
	Dimension   map[string]string `json:"dimension"`
	TimeStampMs int64             `json:"timestamp"`
}

type EventInfo struct {
	Content string `json:"content"`
}
