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

package plugins

import (
	"time"

	"configcenter/src/apimachinery/flowctrl"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/json"
	"configcenter/src/thirdparty/monitor"
)

func init() {
	gseCmdline := monitor.NewGseCmdline()
	bkMonitor := &BKMonitor{
		gseCmdline:      gseCmdline,
		gseCmdAvailable: gseCmdline.IsAvailable(),
	}
	blog.Infof("gseCmdAvailable is %v", bkMonitor.gseCmdAvailable)
	Register(common.BKBluekingMonitorPlugin, bkMonitor)

	go bkMonitor.reportLoop()
}

// BKMonitor is a implementation of Monitor for blueking
type BKMonitor struct {
	gseCmdline      *monitor.GseCmdline
	gseCmdAvailable bool
}

func (m *BKMonitor) Collect(msg interface{}) error {
	if !monitor.MonitorCfg.EnableMonitor || monitor.MonitorCfg.DataID == 0 || !m.gseCmdAvailable {
		return nil
	}

	switch msg.(type) {
	case monitor.MonitorData:
		data := msg.(monitor.MonitorData)
		pushToQueue(&data)
	case *monitor.MonitorData:
		pushToQueue(msg.(*monitor.MonitorData))
	case []byte:
		data := new(monitor.MonitorData)
		if err := json.Unmarshal(msg.([]byte), data); err != nil {
			blog.Infof("Unmarshal failed, msg:%s, err:%v", msg, err)
			break
		}
		pushToQueue(data)

	case string:
		data := new(monitor.MonitorData)
		if err := json.Unmarshal([]byte(msg.(string)), data); err != nil {
			blog.Infof("Unmarshal failed, msg:%s, err:%v", msg, err)
			break
		}
		pushToQueue(data)
	default:
		blog.V(4).Infof("unknown msg type:%T, msg:%#v", msg, msg)
	}

	return nil
}

// queue used to cache the data which will be reported
var queue chan *monitor.MonitorData

// pushToQueue push data to queue
func pushToQueue(data *monitor.MonitorData) {
	for len(queue) >= int(monitor.MonitorCfg.QueueSize) {
		<-queue
	}
	queue <- data
}

// reportLoop report msg continuously
func (m *BKMonitor) reportLoop() {
	for !monitor.MonitorCfg.FinishInit {
		time.Sleep(time.Second)
	}
	queue = make(chan *monitor.MonitorData, monitor.MonitorCfg.QueueSize)

	if !monitor.MonitorCfg.EnableMonitor || monitor.MonitorCfg.DataID == 0 || !m.gseCmdAvailable {
		return
	}

	gNum := 10
	for i := 0; i < gNum; i++ {
		go func() {
			throttle := flowctrl.NewRateLimiter(monitor.MonitorCfg.QPS, monitor.MonitorCfg.Burst)
			for data := range queue {
				throttle.Accept()
				msg, err := convertToReportMsg(data)
				if err != nil {
					blog.Errorf("convertToReportMsg failed, err:%v, data:%#v", err, data)
				}
				m.gseCmdline.Report(msg)
			}
		}()
	}
}

// convertToReportMsg convert data to a msg used by gseCmdline
func convertToReportMsg(data *monitor.MonitorData) (string, error) {
	event := EventData{
		DataID: monitor.MonitorCfg.DataID,
		Data: []*EventMsg{
			{
				EventName: data.EventName,
				EventInfo: EventInfo{Content: data.EventContent},
				Target:    monitor.MonitorCfg.SourceIP,
				Dimension: map[string]interface{}{
					"module":     data.Module,
					"request_id": data.RequestID,
				},
				TimeStampMs: time.Now().UnixNano() / 1e6,
			},
		},
	}
	msg, err := json.Marshal(event)
	if err != nil {
		blog.Errorf("marshal error:%v, msg:%s", err, msg)
		return "", err
	}
	return string(msg), nil
}

// EventData is self-defined event in bk-monitor
type EventData struct {
	DataID int64       `json:"dataid"`
	Data   []*EventMsg `json:"data"`
}

type EventMsg struct {
	EventName   string                 `json:"event_name"`
	EventInfo   EventInfo              `json:"event"`
	Target      string                 `json:"target"`
	Dimension   map[string]interface{} `json:"dimension"`
	TimeStampMs int64                  `json:"timestamp"`
}

type EventInfo struct {
	Content string `json:"content"`
}
