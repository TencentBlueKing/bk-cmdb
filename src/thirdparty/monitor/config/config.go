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

// Package config TODO
package config

import (
	"errors"
	"net"

	"configcenter/src/common"
	"configcenter/src/common/blog"
)

const (
	// QueueSizeMax TODO
	QueueSizeMax = 1000
	// QueueSizeMin TODO
	QueueSizeMin = 1
	// QueueSizeDefault TODO
	QueueSizeDefault = 100

	// QPSMax TODO
	QPSMax = 50
	// QPSMin TODO
	QPSMin = 1
	// QPSDefault TODO
	QPSDefault = 10

	// BurstMax TODO
	BurstMax = 100
	// BurstMin TODO
	BurstMin = 1
	// BurstDefault TODO
	BurstDefault = 20
)

var (
	// MonitorCfg TODO
	MonitorCfg = new(MonitorConfig)
)

// MonitorConfig is the config of monitor
type MonitorConfig struct {
	// EnableMonitor enable monitor or not
	EnableMonitor bool
	// PluginName current plugin name
	PluginName string
	// QueueSize is queue size to cache the collected data
	QueueSize int64
	// QPS used for rate limit
	QPS int64
	// Burst used for rate limit, represent the capacity
	Burst int64
	// SourceIP is the source ip address to report data
	SourceIP string

	BluekingPluginConfig
}

// BluekingPluginConfig is the config of blueking monitor plugin
type BluekingPluginConfig struct {
	// DataID is identifier for report data
	DataID int64
	// BkMonitorReportUrl blueking monitor report url
	BkMonitorReportUrl string
	// AccessToken blueking monitor report access token
	AccessToken string
}

// SetMonitorSourceIP set monitor source ip
func SetMonitorSourceIP() {
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

// CheckAndCorrectCfg check the config, correct it if config is wrong
func CheckAndCorrectCfg() error {
	if MonitorCfg.QueueSize < QueueSizeMin || MonitorCfg.QueueSize > QueueSizeMax {
		MonitorCfg.QueueSize = QueueSizeDefault
	}

	if MonitorCfg.QPS < QPSMin || MonitorCfg.QPS > QPSMax {
		MonitorCfg.QPS = QPSDefault
	}

	if MonitorCfg.Burst < BurstMin || MonitorCfg.Burst > BurstMax {
		MonitorCfg.Burst = BurstDefault
	}

	switch MonitorCfg.PluginName {
	case common.BKBluekingMonitorPlugin:
		if MonitorCfg.DataID == 0 {
			blog.Errorf("init monitor failed, config monitor.dataID is not set")
			return errors.New("config monitor.dataID is not set")
		}

		if MonitorCfg.BkMonitorReportUrl == "" {
			blog.Errorf("init monitor failed, config monitor.bkMonitorReportUrl is not set")
			return errors.New("config monitor.bkMonitorReportUrl is not set")
		}

		if MonitorCfg.AccessToken == "" {
			blog.Errorf("init monitor failed, config monitor.accessToken is not set")
			return errors.New("config monitor.accessToken is not set")
		}
	}

	return nil
}
