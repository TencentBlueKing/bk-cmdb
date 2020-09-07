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

package hostsnap

import (
	"fmt"
	"strconv"
	"strings"

	"configcenter/src/common/json"

	"github.com/tidwall/gjson"
)

// ParseHostSnap parse host snapshot info
func ParseHostSnap(snapData *gjson.Result) (*string, error) {
	ret := ""
	if snapData.String() == "" {
		return &ret, nil
	}

	snap := snapData.Get("data")
	if !snap.Exists() {
		return &ret, nil
	}

	// cpu
	cpuNum := 0
	cpuUsage := 0
	if cpu := snap.Get("cpu"); cpu.Exists() {
		cpuNum = len(cpu.Get("per_usage").Array())
		cpuUsage = int(cpu.Get("total_usage").Float()*100 + 0.5)
	}

	// disk
	diskTotal := uint64(0)
	diskUsed := uint64(0)
	diskUsage := uint64(0)
	if disk := snap.Get("disk"); disk.Exists() {
		for _, diskInfo := range disk.Get("usage").Array() {
			diskTotal += diskInfo.Get("total").Uint()
			diskUsed += diskInfo.Get("used").Uint()
		}
		// unit is GB
		diskTotal = diskTotal >> 10 >> 10 >> 10
		diskUsed = diskUsed >> 10 >> 10 >> 10
		if 0 != diskTotal {
			// get the percentage with two bits reserved.
			diskUsage = 10000 * diskUsed / diskTotal
		} else {
			diskUsage = 0
		}
	}

	// mem info
	memUsage := uint64(0)
	memTotal := uint64(0)
	memUsed := uint64(0)
	if info := snap.Get("mem.meminfo"); info.Exists() {
		memUsage = uint64(100*info.Get("usedPercent").Float() + 0.5)
		// unit is MB
		var unitMB uint64 = 1 << 10 << 10
		memTotal = (info.Get("total").Uint() + unitMB - 1) >> 10 >> 10
		memUsed = (info.Get("used").Uint() + unitMB - 1) >> 10 >> 10
	}

	// load info
	loadAvgStr := ""
	if loadAvg := snap.Get("load.load_avg"); loadAvg.Exists() {
		loadAvgStr = fmt.Sprintf("%.2f %.2f %.2f",
			loadAvg.Get("load1").Float(), loadAvg.Get("load5").Float(), loadAvg.Get("load15").Float())
	}

	// system info
	hostName := ""
	osName := ""
	bootTime := uint64(0)
	if info := snap.Get("system.info"); info.Exists() {
		hostName = info.Get("hostname").String()
		osName = info.Get("os").String()
		bootTime = info.Get("bootTime").Uint()
	}

	upTime := snap.Get("datetime").String()
	timezoneNum := snap.Get("timezone").Int()

	// time zone info
	timezone := ""
	country := snap.Get("country").String()
	city := snap.Get("city").String()
	if country != "" && city != "" {
		timezone = country + "/" + city
	} else if country != "" {
		timezone = country
	} else if city != "" {
		timezone = city
	}

	// net info
	interf := make([]byte, 0)
	rcvRate := uint64(0)
	sendRate := uint64(0)
	if net := snap.Get("net"); net.Exists() {
		rcvRate, sendRate, _ = calculateNetSpeed(net.Get("dev").Array())
		interf = getInterface(net.Get("interface").Array())
	}

	raw := strings.Builder{}
	raw.WriteByte('{')
	raw.WriteString("\"Cpu\":")
	raw.WriteString(strconv.Itoa(cpuNum))
	raw.WriteString(",")
	raw.WriteString("\"cpuUsage\":")
	raw.WriteString(strconv.Itoa(cpuUsage))
	raw.WriteString(",")
	raw.WriteString("\"Disk\":")
	raw.WriteString(strconv.FormatUint(diskTotal, 10))
	raw.WriteString(",")
	raw.WriteString("\"diskUsage\":")
	raw.WriteString(strconv.FormatUint(diskUsage, 10))
	raw.WriteString(",")
	raw.WriteString("\"memUsage\":")
	raw.WriteString(strconv.FormatUint(memUsage, 10))
	raw.WriteString(",")
	raw.WriteString("\"Mem\":")
	raw.WriteString(strconv.FormatUint(memTotal, 10))
	raw.WriteString(",")
	raw.WriteString("\"memUsed\":")
	raw.WriteString(strconv.FormatUint(memUsed, 10))
	raw.WriteString(",")
	raw.WriteString("\"loadavg\":")
	raw.Write([]byte("\"" + loadAvgStr + "\""))
	raw.WriteString(",")
	raw.WriteString("\"HostName\":")
	raw.Write([]byte("\"" + hostName + "\""))
	raw.WriteString(",")
	raw.WriteString("\"OsName\":")
	raw.Write([]byte("\"" + osName + "\""))
	raw.WriteString(",")
	raw.WriteString("\"bootTime\":")
	raw.WriteString(strconv.FormatUint(bootTime, 10))
	raw.WriteString(",")
	raw.WriteString("\"upTime\":")
	raw.Write([]byte("\"" + upTime + "\""))
	raw.WriteString(",")
	raw.WriteString("\"timezone_number\":")
	raw.WriteString(strconv.FormatInt(timezoneNum, 10))
	raw.WriteString(",")
	raw.WriteString("\"timezone\":")
	raw.Write([]byte("\"" + timezone + "\""))
	raw.WriteString(",")
	raw.WriteString("\"rcvRate\":")
	raw.WriteString(strconv.FormatUint(rcvRate, 10))
	raw.WriteString(",")
	raw.WriteString("\"sendRate\":")
	raw.WriteString(strconv.FormatUint(sendRate, 10))
	raw.WriteString(",")
	raw.WriteString("\"bk_all_ips\":")
	raw.Write(interf)
	raw.WriteByte('}')

	ret = raw.String()
	return &ret, nil
}

// calculateNetSpeed calculates net rcvRate and sendRate
func calculateNetSpeed(devInfo []gjson.Result) (uint64, uint64, error) {
	var rcvRate uint64 = 0
	var sendRate uint64 = 0

	for _, info := range devInfo {
		if 0 <= strings.Index(info.Get("name").String(), "lo") {
			continue
		}
		rcvRate += info.Get("speedRecv").Uint()
		sendRate += info.Get("speedSent").Uint()
	}

	// unit is Mb/s
	rcvRate = (100 * rcvRate) >> 10 >> 10
	sendRate = (100 * sendRate) >> 10 >> 10
	return rcvRate, sendRate, nil
}

// getInterface get all the ip and mac info
func getInterface(interfaceInfo []gjson.Result) []byte {
	allIPs := make(map[string]interface{})
	interfaceArr := make([]interface{}, 0)

	for _, info := range interfaceInfo {
		addrs := make([]map[string]string, 0)
		for _, addr := range info.Get("addrs").Array() {
			// ignore the loopback ip
			if strings.Contains(addr.Get("addr").String(), "127.0.0.1") || strings.Contains(addr.Get("addr").String(), "::1") {
				continue
			}
			// append ip address without mask
			addrs = append(addrs, map[string]string{"ip": strings.Split(addr.Get("addr").String(), "/")[0]})
		}
		if len(addrs) != 0 {
			interfaceArr = append(interfaceArr, map[string]interface{}{
				"mac":   info.Get("hardwareaddr").String(),
				"addrs": addrs,
			})
		}
	}

	allIPs["interface"] = interfaceArr

	interf, _ := json.Marshal(allIPs)
	return interf
}
