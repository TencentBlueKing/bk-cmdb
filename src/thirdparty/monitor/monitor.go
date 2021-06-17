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
	"errors"
	"fmt"
	"time"

	"configcenter/src/apimachinery/flowctrl"
	"configcenter/src/common"
	cc "configcenter/src/common/backbone/configcenter"
	"configcenter/src/common/blog"
	"configcenter/src/common/types"
	"configcenter/src/thirdparty/monitor/config"
	"configcenter/src/thirdparty/monitor/meta"
	"configcenter/src/thirdparty/monitor/plugins"
)

type Monitor struct {
	// queue is the plugin used to process the collected data
	plugin plugins.Plugin
	// queue used to cache the data which will be reported
	queue chan meta.Content
}

var monitor *Monitor

// Collect is the monitor entrance used by other services
func Collect(c meta.Content) {
	if monitor == nil {
		return
	}

	if !config.MonitorCfg.EnableMonitor {
		blog.V(4).InfoDepthf(1, "Collect skipped, the monitor is not enabled")
		return
	}

	monitor.collect(c)
}

// startMonitor create monitor instance and start monitor
func startMonitor() error {
	plugin, err := plugins.GetPlugin(config.MonitorCfg.PluginName)
	if err != nil {
		blog.Errorf("start monitor failed, GetPlugin err: %s, pluginName:%s", err, config.MonitorCfg.PluginName)
		return err
	}

	monitor = &Monitor{
		plugin: plugin,
		queue:  make(chan meta.Content, config.MonitorCfg.QueueSize),
	}
	monitor.start()
	blog.Infof("start monitor successfully, plugin is %#v", monitor.plugin)
	return nil
}

// collect process the content collected
func (m *Monitor) collect(c meta.Content) {
	m.pushToQueue(c)
}

// pushToQueue push data to queue
// throw the data away if exceed the queue length
func (m *Monitor) pushToQueue(c meta.Content) {
	select {
	case m.queue <- c:
	default:
	}
}

// start run the monitor framework
func (m *Monitor) start() {
	if !config.MonitorCfg.EnableMonitor {
		return
	}

	go m.reportLoop()
}

// reportLoop report data continuously and control the report rate
func (m *Monitor) reportLoop() {
	// control the report rate
	throttle := flowctrl.NewRateLimiter(config.MonitorCfg.QPS, config.MonitorCfg.Burst)
	for content := range m.queue {
		if throttle.TryAccept() {
			m.plugin.Report(content)
		}
	}
}

// InitMonitor init monitor config and monitor instance
func InitMonitor() error {

	var err error
	// no need init for adminserver
	if common.GetIdentification() == types.CC_MODULE_MIGRATE {
		return nil
	}

	maxCnt := 100
	cnt := 0
	for !cc.IsExist("monitor") && cnt < maxCnt {
		blog.Infof("waiting monitor config to be init")
		cnt++
		time.Sleep(time.Millisecond * 300)
	}

	if cnt == maxCnt {
		blog.Infof("init monitor failed, no monitor config is found, the config 'monitor' must exist")
		return fmt.Errorf("init monitor failed, no monitor config is found, the config 'monitor' must exist")
	}

	config.MonitorCfg.EnableMonitor, err = cc.Bool("monitor.enableMonitor")
	if err != nil {
		blog.Errorf("init monitor failed,monitor.enableMonitor err: %v", err)
		return errors.New("config monitor.enableMonitor is not found")
	}
	//如果不需要进行监控上报，那么后续没有必要再检查配置
	if !config.MonitorCfg.EnableMonitor {
		return nil
	}

	config.MonitorCfg.PluginName, err = cc.String("monitor.pluginName")
	if err != nil {
		blog.Errorf("init monitor failed, err: %v", err)
		return errors.New("config monitor.pluginName is not found")
	}

	dataID, err := cc.Int64("monitor.dataID")
	if err != nil {
		blog.Errorf("init monitor failed, err: %v", err)
		return errors.New("config monitor.dataID is not found")
	}
	config.MonitorCfg.DataID = dataID

	queueSize, err := cc.Int64("monitor.queueSize")
	if err != nil {
		blog.Errorf("init monitor failed, err: %v", err)
		return errors.New("config monitor.queueSize is not found")
	}
	config.MonitorCfg.QueueSize = queueSize

	qps, err := cc.Int64("monitor.rateLimiter.qps")
	if err != nil {
		blog.Errorf("init monitor failed, err: %v", err)
		return errors.New("config monitor.qps is not found")
	}
	config.MonitorCfg.QPS = qps

	burst, err := cc.Int64("monitor.rateLimiter.burst")
	if err != nil {
		blog.Errorf("init monitor failed, err: %v", err)
		return errors.New("config monitor.burst is not found")
	}
	config.MonitorCfg.Burst = burst

	config.MonitorCfg.GsecmdlinePath, err = cc.String("monitor.gsecmdlinePath")
	if err != nil {
		blog.Errorf("init monitor failed, err: %v", err)
		return errors.New("config monitor.gsecmdlinePath is not found")
	}

	config.MonitorCfg.DomainSocketPath, err = cc.String("monitor.domainSocketPath")
	if err != nil {
		blog.Errorf("init monitor failed, err: %v", err)
		return errors.New("config monitor.domainSocketPath is not found")
	}

	err = config.CheckAndCorrectCfg()
	if err != nil {
		blog.Errorf("init monitor check cfg failed,  err: %v", err)
		return err
	}

	config.SetMonitorSourceIP()
	if err := startMonitor(); err != nil {
		blog.Errorf("init monitor failed, startMonitor err: %v", err)
		return err
	}
	blog.InfoJSON("init monitor successfully, cfg: %s", config.MonitorCfg)

	return nil
}
