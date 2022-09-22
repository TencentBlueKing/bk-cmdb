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
	"testing"
)

// TestParseHostSnap test for function ParseHostSnap
func TestParseHostSnap(t *testing.T) {
	checkItems := map[string]interface{}{
		"Mem":             uint64(64132),
		"HostName":        "bkdata-111",
		"bootTime":        uint64(1510553867),
		"rcvRate":         uint64(346),
		"Cpu":             2,
		"Disk":            uint64(271),
		"memUsage":        uint64(7500),
		"upTime":          "2017-12-26 20:35:14",
		"hosts":           []string{"::1 localhost localhost.localdomain localhost6 localhost6.localdomain6", "127.0.0.1"},
		"loadavg":         "2.60 3.54 3.78",
		"OsName":          "linux",
		"sendRate":        uint64(110),
		"diskUsage":       uint64(959),
		"crontab":         map[string]string{"root": "*/10 * * * * /data/bkee/bkdata/dataapi/bin/update_cc_cache.sh"},
		"timezone":        "China/Shenzhen",
		"cpuUsage":        int(1846),
		"memUsed":         uint64(48097),
		"timezone_number": 8,
	}
	ret, _ := ParseHostSnap(data)

	// print result
	str := ""
	for k, v := range ret {
		str += fmt.Sprintf("%s:%+v\n", k, v)
	}
	t.Log(str)

	// check result
	for k, _ := range checkItems {
		if k == "hosts" {
			if ret[k].([]string)[0] != checkItems[k].([]string)[0] || ret[k].([]string)[1] != checkItems[k].([]string)[1] {
				t.Logf("check error, for key:%v, the ret[%s] is %#v, but checkItems[%s] is %#v", k, k, ret[k], k, checkItems[k])
				t.Fail()
			}
			continue
		}
		if k == "crontab" {
			if ret[k].(map[string]string)["root"] != checkItems[k].(map[string]string)["root"] {
				t.Logf("check error, for key:%v, the ret[%s] is %#v, but checkItems[%s] is %#v", k, k, ret[k], k, checkItems[k])
				t.Fail()
			}
			continue
		}
		if ret[k] != checkItems[k] {
			t.Logf("check error, for key:%v, the ret[%s] is %+v, but checkItems[%s] is %+v", k, k, ret[k], k, checkItems[k])
			t.Fail()
		}
	}
}

// TestJson test for Json operate to object HostSnap
func TestJson(t *testing.T) {
	snap := HostSnap{}
	err := json.Unmarshal([]byte(data), &snap)
	if err != nil {
		t.Logf("Unmarshal err:%v", err)
	}
	t.Logf("snap:%+v", snap)
	t.Logf("sys:%+v", snap.Data.System.Info.SystemStat)
}

// TestParseHostSnap1 test for finding bug
func TestParseHostSnap1(t *testing.T) {
	ret1, _ := ParseHostSnap(data1)
	t.Logf("loadavg:%v", ret1["loadavg"])
}

var data1 string = `{
  "data": {
    "load": {
      "load_avg": {
      }
    },
    "cpu": {
      "per_usage": [
        19.412862718581678,
        20.893561103643428
      ],
      "total_usage": 18.45950545867876
    },
    "net": {
      "dev": [
        {
          "name": "eth0",
          "speedSent": 1111110,
          "speedRecv": 1111110,
          "speedPacketsSent": 111110,
          "speedPacketsRecv": 222220,
          "bytesSent": 333330,
          "bytesRecv": 4444440,
          "packetsSent": 0,
          "packetsRecv": 0
        }
      ]
    }
  }
}`

var data string = `{
  "bizid": 0,
  "cloudid": 0,
  "data": {
    "load": {
      "load_avg": {
        "load1": 2.6,
        "load5": 3.54,
        "load15": 3.78
      }
    },
    "timezone": 8,
    "datetime": "2017-12-26 20:35:14",
    "utctime": "2017-12-26 12:35:14",
    "country": "China",
    "city": "Shenzhen",
    "cpu": {
      "per_usage": [
        19.412862718581678,
        20.893561103643428
      ],
      "total_usage": 18.45950545867876
    },
    "env": {
      "crontab": [
        {
          "user": "root",
          "content": "*/10 * * * * /data/bkee/bkdata/dataapi/bin/update_cc_cache.sh"
        }
      ],
      "host": "::1 localhost localhost.localdomain localhost6 localhost6.localdomain6\n127.0.0.1"
    },
    "disk": {
      "usage": [
        {
          "path": "/",
          "fstype": "ext2/ext3",
          "total": 21003583488,
          "free": 12865982464,
          "used": 7047081984,
          "usedPercent": 33.55180789995296
        },
        {
          "path": "/data",
          "fstype": "ext2/ext3",
          "total": 249789050880,
          "free": 217206951936,
          "used": 19869900800,
          "usedPercent": 7.954672444608313
        },
        {
          "path": "/usr/local",
          "fstype": "ext2/ext3",
          "total": 21003583488,
          "free": 18666192896,
          "used": 1246871552,
          "usedPercent": 5.936470568045574
        }
      ]
    },
    "mem": {
      "meminfo": {
        "total": 67246661632,
        "available": 16814030848,
        "used": 50432630784,
        "usedPercent": 74.99648244248473,
        "free": 3770302464
      }
    },
    "net": {
      "dev": [
        {
          "name": "eth0",
          "speedSent": 0,
          "speedRecv": 0,
          "speedPacketsSent": 0,
          "speedPacketsRecv": 0,
          "bytesSent": 0,
          "bytesRecv": 0,
          "packetsSent": 0,
          "packetsRecv": 0
        },
        {
          "name": "eth1",
          "speedSent": 1158566,
          "speedRecv": 3636136,
          "speedPacketsSent": 7099,
          "speedPacketsRecv": 12360,
          "bytesSent": 3257536071332,
          "bytesRecv": 7616250441907,
          "packetsSent": 27005508120,
          "packetsRecv": 44093455831
        },
        {
          "name": "lo",
          "speedSent": 13570,
          "speedRecv": 13570,
          "speedPacketsSent": 96,
          "speedPacketsRecv": 96,
          "bytesSent": 158794029730,
          "bytesRecv": 158794029730,
          "packetsSent": 311070174,
          "packetsRecv": 311070174
        }
      ]
    },
    "system": {
      "info": {
        "hostname": "bkdata-111",
        "uptime": 3737847,
        "bootTime": 1510553867,
        "procs": 794,
        "os": "linux",
        "platform": "centos",
        "hostid": "0C7CA213-D21D-B211-B859-03242340"
      }
    }
  },
  "ip": "127.0.0.1"
}`
