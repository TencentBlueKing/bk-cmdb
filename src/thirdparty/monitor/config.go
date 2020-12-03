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
	"net"
	"time"

	cc "configcenter/src/common/backbone/configcenter"
	"configcenter/src/common/blog"
)

const (
	QueueSizeMax     = 10000
	QueueSizeMin     = 1
	QueueSizeDefault = 100

	QPSMax     = 10000
	QPSMin     = 1
	QPSDefault = 10

	BurstMax     = 20000
	BurstMin     = 1
	BurstDefault = 20
)

var (
	MonitorCfg = new(MonitorConfig)
)

// MonitorConfig is the config of monitor
type MonitorConfig struct {
	FinishInit bool
	// PluginName current plugin name
	PluginName string
	// EnableMonitor enable monitor or not
	EnableMonitor bool
	// DataID is identifier for report data
	DataID int64
	// QueueSize is queue size to cache the collected data
	QueueSize int64
	// QPS used for rate limit
	QPS int64
	// Burst used for rate limit, represent the capacity
	Burst int64
	// SourceIP is the source ip address to report data
	SourceIP string
}

// InitMonitorCfg set monitor config
func InitMonitorCfg() {
	maxCnt := 100
	cnt := 0
	for MonitorCfg.PluginName == "" && cnt < maxCnt {
		MonitorCfg.PluginName, _ = cc.String("monitor.pluginName")
		MonitorCfg.EnableMonitor, _ = cc.Bool("monitor.enableMonitor")
		dataID, _ := cc.Int("monitor.dataID")
		MonitorCfg.DataID = int64(dataID)
		queueSize, _ := cc.Int("monitor.queueSize")
		MonitorCfg.QueueSize = int64(queueSize)
		qps, _ := cc.Int("monitor.qps")
		MonitorCfg.QPS = int64(qps)
		burst, _ := cc.Int("monitor.burst")
		MonitorCfg.Burst = int64(burst)
		cnt++
		blog.Infof("waiting monitor config to be init")
		time.Sleep(time.Second)
	}
	MonitorCfg.FinishInit = true
	if cnt == maxCnt {
		blog.Infof("InitMonitorCfg failed, no monitor config is found")
		return
	}
	checkAndCorrectCfg()
	setMonitorSourceIP()
	blog.Infof("InitMonitorCfg successfully, cfg: %#v", MonitorCfg)
}

// setMonitorSourceIP set monitor source ip
func setMonitorSourceIP() {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		blog.Errorf("get addrs err:%v", err)
		return
	}

	for _, address := range addrs {
		// check if the ip is loop address
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				MonitorCfg.SourceIP = ipnet.IP.String()
				blog.Infof("source ip is %s", MonitorCfg.SourceIP)
				return
			}
		}
	}
}

// checkAndCorrectCfg check the config, correct it if config is wrong
func checkAndCorrectCfg() {
	if MonitorCfg.QueueSize < QueueSizeMin || MonitorCfg.QueueSize > QueueSizeMax {
		MonitorCfg.QueueSize = QueueSizeDefault
	}

	if MonitorCfg.QPS < QPSMin || MonitorCfg.QPS > QPSMax {
		MonitorCfg.QPS = QPSDefault
	}

	if MonitorCfg.Burst < BurstMin || MonitorCfg.Burst > BurstMax {
		MonitorCfg.Burst = BurstDefault
	}
}
