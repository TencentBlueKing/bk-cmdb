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
	"fmt"
	"strings"

	"github.com/bitly/go-simplejson"

	"configcenter/src/common"
	"configcenter/src/common/util"
)

func ParseHostSnap(data string) (map[string]interface{}, error) {
	if "" == data {
		return nil, nil
	}
	js, err := simplejson.NewJson([]byte(data))
	if nil != err {
		return nil, err
	}
	js = js.Get("data")
	if nil == js {
		return nil, nil
	}
	ret := make(map[string]interface{})

	//cpu
	cpuUsageArr, _ := js.Get("cpu").Get("per_usage").Array()
	cpuNum := len(cpuUsageArr)
	cpuUsageFload, _ := js.Get("cpu").Get("total_usage").Float64()
	cpuUsage := int((cpuUsageFload)*100 + 0.5)

	//disk
	diskInfos, _ := js.Get("disk").Get("usage").Array() //!empty($data['disk']['usage']) ? $data['disk']['usage'] : array();
	var diskTotal int64 = 0
	var diskUsage int64 = 0
	var diskUsed int64 = 0

	for _, diskInfoI := range diskInfos {
		disk, ok := diskInfoI.(map[string]interface{})
		if ok {
			total, _ := util.GetInt64ByInterface(disk["total"])
			used, _ := util.GetInt64ByInterface(disk["used"])
			diskTotal += total
			diskUsed += used
		}

	}
	var unitGB int64 = 1024 * 1024 * 1024
	var unitMB int64 = 1024 * 1024
	diskTotal = diskTotal / unitGB
	diskUsed = diskUsed / unitGB
	if 0 != diskTotal {
		diskUsage = (10000 * diskUsed / diskTotal) //获取使用百分比 保留两位小数
	} else {
		diskUsage = 0
	}

	//iptable info
	iptables, _ := js.Get("env").Get("iptables").String()

	//hosts info
	hosts, _ := js.Get("env").Get("host").String()

	//crontab info
	cronInfos, err := js.Get("env").Get("crontab").Array()

	var crontabs common.KvMap = make(common.KvMap)
	if nil != err {
		for _, cronI := range cronInfos {
			cron, ok := cronI.(map[string]string)
			if ok {
				user, ok := cron["user"]
				if "" == user || !ok {
					user = "root"
				}
				content, _ := cron["content"]
				crontabs[user] = content
			}

		}
	}

	//route info
	route, _ := js.Get("env").Get("route").String()

	//mem info
	memInfo := js.Get("mem").Get("meminfo")
	var memTotal, memUsed, memUsage int64
	if nil != memInfo {
		memTotal, _ = util.GetInt64ByInterface(memInfo.Get("total").Interface())
		memUsed, _ = util.GetInt64ByInterface(memInfo.Get("used").Interface())
		memUsageF, _ := memInfo.Get("usedPercent").Float64()
		memUsage = int64(100*memUsageF + 0.5)
		memTotal = (memTotal + unitMB - 1) / unitMB
		memUsed = (memUsed + unitMB - 1) / unitMB

	}
	//系统负载信息
	load := js.Get("load").Get("load_avg")
	strLoadavg := ""
	if nil != load {
		load1 := load.Get("load1").Interface()
		load5 := load.Get("load5").Interface()
		load15 := load.Get("load15").Interface()
		strLoadavg = fmt.Sprintf("%v %v %v", load1, load5, load15)
		strLoadavg = strings.Replace(strLoadavg, "nil", "", -1)
	}

	ret["Cpu"] = cpuNum
	ret["cpuUsage"] = cpuUsage
	ret["Mem"] = memTotal
	ret["memUsage"] = memUsage
	ret["memUsed"] = memUsed

	ret["Disk"] = diskTotal
	ret["diskUsage"] = diskUsage
	if "" != hosts {
		ret["hosts"] = strings.Split(hosts, "\n")
	}
	if "" != iptables {
		ret["iptables"] = strings.Split(iptables, "\n")
	}
	if 0 < len(crontabs) {
		ret["crontab"] = crontabs

	}
	//not empty
	if "" != route {
		ret["route"] = strings.Split(route, "\n")
	}

	if "" != strLoadavg {
		ret["loadavg"] = strLoadavg
	}

	//os info
	ret["HostName"], _ = js.Get("system").Get("info").Get("hostname").String()
	ret["OsName"], _ = js.Get("system").Get("info").Get("os").String()
	ret["bootTime"], _ = util.GetIntByInterface(js.Get("system").Get("info").Get("bootTime").Interface())
	ret["upTime"], _ = js.Get("datetime").String()
	ret["timezone_number"], _ = util.GetIntByInterface(js.Get("timezone").Interface())

	//time zone info
	city, _ := js.Get("city").String()
	country, _ := js.Get("country").String()

	ret["timezone"] = country + "/" + city
	ret["rcvRate"], ret["sendRate"], err = getSnapNetInfo(js.Get("net").Get("dev"), unitMB)
	// $dataNetDev = !empty($data['net']['dev']) ? $data['net']['dev'] : array();*/
	return ret, nil

}

func getSnapNetInfo(netinfosI *simplejson.Json, unitMB int64) (int64, int64, error) {
	netinfos, err := netinfosI.Array()
	var rcvRate int64 = 0
	var sendRate int64 = 0

	if nil == err {

		for _, netinfoI := range netinfos {

			netinfo, ok := netinfoI.(map[string]interface{})
			if ok {
				name, ok := netinfo["name"].(string)
				if !ok {
					continue
				}
				if 0 <= strings.Index(name, "lo") {
					continue
				}
				rcvRateVal, _ := util.GetInt64ByInterface(netinfo["speedRecv"])
				sendRateVal, _ := util.GetInt64ByInterface(netinfo["speedSent"])
				rcvRate += rcvRateVal
				sendRate += sendRateVal
			}
		}
	}
	rcvRate = 100 * rcvRate / unitMB
	sendRate = 100 * sendRate / unitMB
	return rcvRate, sendRate, err
}
