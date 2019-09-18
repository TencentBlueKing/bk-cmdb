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

package logics

import (
	"encoding/json"
	"fmt"
	"strings"

	"configcenter/src/common/blog"
)

type HostSnap struct {
	Data *SnapData `json:"data"`
}

type SnapData struct {
	Load     *Load   `json:"load"`
	CPU      *CPU    `json:"cpu"`
	Env      *Env    `json:"env"`
	Disk     *Disk   `json:"disk"`
	Mem      *Mem    `json:"mem"`
	Net      *Net    `json:"net"`
	System   *System `json:"system"`
	DateTime string  `json:"datetime"`
	TimeZone int     `json:"timezone"`
	Country  string  `json:"country"`
	City     string  `json:"city"`
}

// Load
type Load struct {
	LoadAvg *LoadAvg `json:"load_avg"`
}

type LoadAvg struct {
	Load1  float64 `json:"load1"`
	Load5  float64 `json:"load5"`
	Load15 float64 `json:"load15"`
}

// CPU
type CPU struct {
	PerUsage   []float64 `json:"per_usage"`
	TotalUsage float64   `json:"total_usage"`
}

// Env
type Env struct {
	IPTables string    `json:"iptables"`
	Host     string    `json:"host"`
	Crontab  []Crontab `json:"crontab"`
	Route    string    `json:"route"`
}

type Crontab struct {
	User    string `json:"user"`
	Content string `json:"content"`
}

// Disk
type Disk struct {
	Usage []DiskInfo `json:"usage"`
}

type DiskInfo struct {
	Total uint64 `json:"total"`
	Used  uint64 `json:"used"`
}

// Mem
type Mem struct {
	Meminfo *Meminfo `json:"meminfo"`
}

type Meminfo struct {
	Total       uint64  `json:"total"`
	Used        uint64  `json:"used"`
	UsedPercent float64 `json:"usedPercent"`
}

// Net
type Net struct {
	Dev []DevInfo `json:"dev"`
}

type DevInfo struct {
	Name      string `json:"name"`
	SpeedSent uint64 `json:"speedSent"`
	SpeedRecv uint64 `json:"speedRecv"`
}

// System
type System struct {
	Info SystemInfo `json:"info"`
}

type SystemInfo struct {
	*SystemStat
}

type SystemStat struct {
	HostName string `json:"hostname"`
	OS       string `json:"os"`
	BootTime uint64 `json:"bootTime"`
}

// ParseHostSnap parse hostsnap from jsonstring to map[string]interface{}
func ParseHostSnap(data string) (map[string]interface{}, error) {
	if "" == data {
		return nil, nil
	}

	hostSnap := HostSnap{}
	err := json.Unmarshal([]byte(data), &hostSnap)
	if err != nil {
		blog.Errorf("ParseHostSnap json Unmarshal err:%v", err)
		return nil, err
	}

	snap := hostSnap.Data
	if snap == nil {
		blog.Infof("ParseHostSnap snap is nil")
		return nil, nil
	}

	ret := make(map[string]interface{})

	// cpu
	if cpu := snap.CPU; cpu != nil {
		ret["Cpu"] = len(cpu.PerUsage)
		ret["cpuUsage"] = int(cpu.TotalUsage*100 + 0.5)
	}

	var unitGB uint64 = 1024 * 1024 * 1024
	var unitMB uint64 = 1024 * 1024

	// disk
	if disk := snap.Disk; disk != nil {
		var diskTotal, diskUsed, diskUsage uint64
		for _, d := range disk.Usage {
			diskTotal += d.Total
			diskUsed += d.Used
		}

		diskTotal = diskTotal / unitGB
		diskUsed = diskUsed / unitGB
		if 0 != diskTotal {
			// 获取使用百分比 保留两位小数
			diskUsage = 10000 * diskUsed / diskTotal
		} else {
			diskUsage = 0
		}
		ret["Disk"] = diskTotal
		ret["diskUsage"] = diskUsage
	}

	// env
	if env := snap.Env; env != nil {
		// hosts info
		if env.Host != "" {
			ret["hosts"] = strings.Split(env.Host, "\n")
		}

		// iptables info
		if env.IPTables != "" {
			ret["iptables"] = strings.Split(env.IPTables, "\n")
		}

		// crontab info
		crontabs := make(map[string]string)
		for _, cron := range env.Crontab {
			user := cron.User
			if user == "" {
				user = "root"
			}
			crontabs[user] = cron.Content
		}
		if len(crontabs) > 0 {
			ret["crontab"] = crontabs
		}

		// route info
		if env.Route != "" {
			ret["route"] = strings.Split(env.Route, "\n")
		}
	}

	// mem info
	if mem := snap.Mem; mem != nil {
		if info := mem.Meminfo; info != nil {
			ret["memUsage"] = uint64(100*info.UsedPercent + 0.5)
			ret["Mem"] = (info.Total + unitMB - 1) / unitMB
			ret["memUsed"] = (info.Used + unitMB - 1) / unitMB
		}
	}

	// load info
	if load := snap.Load; load != nil {
		if loadAvg := load.LoadAvg; loadAvg != nil {
			ret["loadavg"] = fmt.Sprintf("%.2f %.2f %.2f", loadAvg.Load1, loadAvg.Load5, loadAvg.Load15)
		}
	}

	// system info
	if system := snap.System; system != nil {
		if stat := system.Info.SystemStat; stat != nil {
			ret["HostName"] = stat.HostName
			ret["OsName"] = stat.OS
			ret["bootTime"] = stat.BootTime
		}
	}

	ret["upTime"] = snap.DateTime
	ret["timezone_number"] = snap.TimeZone

	// time zone info
	if snap.Country != "" && snap.City != "" {
		ret["timezone"] = snap.Country + "/" + snap.City
	} else if snap.Country != "" {
		ret["timezone"] = snap.Country
	} else if snap.City != "" {
		ret["timezone"] = snap.City
	}

	// net info
	if net := snap.Net; net != nil {
		ret["rcvRate"], ret["sendRate"], _ = calculateNetSpeed(net.Dev, unitMB)
	}

	return ret, nil
}

// getSnapNetInfo calculates net rcvRate and sendRate
func calculateNetSpeed(netInfo []DevInfo, unitMB uint64) (uint64, uint64, error) {
	var rcvRate uint64 = 0
	var sendRate uint64 = 0

	for _, info := range netInfo {
		if 0 <= strings.Index(info.Name, "lo") {
			continue
		}
		rcvRate += info.SpeedRecv
		sendRate += info.SpeedSent
	}

	rcvRate = 100 * rcvRate / unitMB
	sendRate = 100 * sendRate / unitMB
	return rcvRate, sendRate, nil
}
