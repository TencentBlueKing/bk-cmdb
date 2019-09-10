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
	"testing"

	"github.com/bitly/go-simplejson"
)

func TestSimpleJson(t *testing.T) {
	js, err := simplejson.NewJson([]byte(data))
	if nil != err {
		t.Error("NewJson error:", err.Error())
	}
	js = js.Get("data")
	aa := js.Get("aa").Get("bb").Get("cc")
	t.Log("aa:", aa, aa == nil, aa.Interface() == nil)
	t.Log(js.Get("load"))
	js.Set("load", nil)
	t.Log(js.Get("load"))
	alist, err := js.Get("aaa").Array()
	t.Logf("alist:%v, err:%v\n", alist, err)
	for i, v := range alist {
		t.Log(i, v)
	}
	afloat, err := js.Get("aaa").Float64()
	t.Logf("afloat:%v, err:%v\n", afloat, err)
	ainterface := js.Get("aaa").Interface()
	t.Logf("ainterface:%v\n", ainterface)

}

func TestParseHostSnap(t *testing.T) {
	ret, _ := ParseHostSnap(data)
	str := ""
	for k, v := range ret {
		str += fmt.Sprintf("%s:%+v\n", k, v)
	}
	t.Log(str)
	if ret["timezone"] != "CHina/Shenzhen" {
		t.Logf("timezone err, timezone:%v", ret["timezone"])
		t.Fail()
	}
	if ret["loadavg"] != "2.60 3.54 3.78" {
		t.Logf("loadavg err, loadavg:%v", ret["loadavg"])
		t.Fail()
	}
}

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
          "packetsRecv": 0,
          "errin": 0,
          "errout": 0,
          "dropin": 0,
          "dropout": 0,
          "fifoin": 0,
          "fifoout": 0
        }
      ]
    }
  }
}`

var data string = `{
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
    "country": "CHina",
    "city": "Shenzhen",
    "cpu": {
      "cpuinfo": [
        {
          "cpu": 0,
          "vendorId": "GenuineIntel",
          "family": "6",
          "model": "45",
          "stepping": 7,
          "physicalId": "0",
          "coreId": "0",
          "cores": 1,
          "modelName": "Intel(R) Xeon(R) CPU E5-2420 0 @ 1.90GHz",
          "mhz": 1901,
          "cacheSize": 15360,
          "flags": [
            "fpu",
            "vme",
            "de",
            "pse",
            "tsc",
            "msr",
            "pae",
            "mce",
            "cx8",
            "apic",
            "sep",
            "mtrr",
            "pge",
            "mca",
            "cmov",
            "pat",
            "pse36",
            "clflush",
            "dts",
            "acpi",
            "mmx",
            "fxsr",
            "sse",
            "sse2",
            "ss",
            "ht",
            "tm",
            "pbe",
            "syscall",
            "nx",
            "pdpe1gb",
            "rdtscp",
            "lm",
            "constant_tsc",
            "arch_perfmon",
            "pebs",
            "bts",
            "rep_good",
            "nopl",
            "xtopology",
            "nonstop_tsc",
            "aperfmperf",
            "eagerfpu",
            "pni",
            "pclmulqdq",
            "dtes64",
            "monitor",
            "ds_cpl",
            "vmx",
            "smx",
            "est",
            "tm2",
            "ssse3",
            "cx16",
            "xtpr",
            "pdcm",
            "pcid",
            "dca",
            "sse4_1",
            "sse4_2",
            "x2apic",
            "popcnt",
            "tsc_deadline_timer",
            "aes",
            "xsave",
            "avx",
            "lahf_lm",
            "epb",
            "tpr_shadow",
            "vnmi",
            "flexpriority",
            "ept",
            "vpid",
            "xsaveopt",
            "dtherm",
            "ida",
            "arat",
            "pln",
            "pts"
          ],
          "microcode": "0x710"
        },
        {
          "cpu": 1,
          "vendorId": "GenuineIntel",
          "family": "6",
          "model": "45",
          "stepping": 7,
          "physicalId": "0",
          "coreId": "1",
          "cores": 1,
          "modelName": "Intel(R) Xeon(R) CPU E5-2420 0 @ 1.90GHz",
          "mhz": 1901,
          "cacheSize": 15360,
          "flags": [
            "fpu",
            "vme",
            "de",
            "pse",
            "tsc",
            "msr",
            "pae",
            "mce",
            "cx8",
            "apic",
            "sep",
            "mtrr",
            "pge",
            "mca",
            "cmov",
            "pat",
            "pse36",
            "clflush",
            "dts",
            "acpi",
            "mmx",
            "fxsr",
            "sse",
            "sse2",
            "ss",
            "ht",
            "tm",
            "pbe",
            "syscall",
            "nx",
            "pdpe1gb",
            "rdtscp",
            "lm",
            "constant_tsc",
            "arch_perfmon",
            "pebs",
            "bts",
            "rep_good",
            "nopl",
            "xtopology",
            "nonstop_tsc",
            "aperfmperf",
            "eagerfpu",
            "pni",
            "pclmulqdq",
            "dtes64",
            "monitor",
            "ds_cpl",
            "vmx",
            "smx",
            "est",
            "tm2",
            "ssse3",
            "cx16",
            "xtpr",
            "pdcm",
            "pcid",
            "dca",
            "sse4_1",
            "sse4_2",
            "x2apic",
            "popcnt",
            "tsc_deadline_timer",
            "aes",
            "xsave",
            "avx",
            "lahf_lm",
            "epb",
            "tpr_shadow",
            "vnmi",
            "flexpriority",
            "ept",
            "vpid",
            "xsaveopt",
            "dtherm",
            "ida",
            "arat",
            "pln",
            "pts"
          ],
          "microcode": "0x710"
        }
      ],
      "per_usage": [
        19.412862718581678,
        20.893561103643428
      ],
      "total_usage": 18.45950545867876,
      "per_stat": [
        {
          "cpu": "cpu0",
          "user": 378058.95,
          "system": 94173.56,
          "idle": 3225693.63,
          "nice": 0.02,
          "iowait": 170.03,
          "irq": 8.43,
          "softirq": 7235.68,
          "steal": 0,
          "guest": 0,
          "guestNice": 0,
          "stolen": 0
        },
        {
          "cpu": "cpu1",
          "user": 370195.8,
          "system": 121713.48,
          "idle": 3223651.06,
          "nice": 0.04,
          "iowait": 145.2,
          "irq": 43.1,
          "softirq": 14556.02,
          "steal": 0,
          "guest": 0,
          "guestNice": 0,
          "stolen": 0
        }
      ],
      "total_stat": {
        "cpu": "cpu-total",
        "user": 7754440.71,
        "system": 2121913.33,
        "idle": 79647735.18,
        "nice": 0.61,
        "iowait": 2370.92,
        "irq": 267.39,
        "softirq": 126155.49,
        "steal": 0,
        "guest": 0,
        "guestNice": 0,
        "stolen": 0
      }
    },
    "env": {
      "crontab": [
        {
          "user": "root",
          "content": "*/10 * * * * /data/bkee/bkdata/dataapi/bin/update_cc_cache.sh\n\\"
        }
      ],
      "host": "::1 localhost localhost.localdomain localhost6 localhost6.localdomain6\n"
    },
    "disk": {
      "diskstat": {
        "sda1": {
          "major": 8,
          "minor": 1,
          "readCount": 67286,
          "mergedReadCount": 3748,
          "writeCount": 47531110,
          "mergedWriteCount": 84670386,
          "readBytes": 1475599360,
          "writeBytes": 585211965440,
          "readSectors": 2882030,
          "writeSectors": 1142992120,
          "readTime": 274652,
          "writeTime": 6977448,
          "iopsInProgress": 0,
          "ioTime": 1978944,
          "weightedIoTime": 7211520,
          "name": "sda1",
          "serialNumber": "36234567890abcde01916b0e132b97a5c",
          "speedIORead": 0,
          "speedByteRead": 0,
          "speedIOWrite": 13.583333333333334,
          "speedByteWrite": 172032,
          "util": 0.0007333333333333332,
          "avgrq_sz": 24.736196319018404,
          "avgqu_sz": 0.001,
          "await": 0.0736196319018405,
          "svctm": 0.053987730061349694
        },
        "sda2": {
          "major": 8,
          "minor": 2,
          "readCount": 97009,
          "mergedReadCount": 555778,
          "writeCount": 65119,
          "mergedWriteCount": 1002427,
          "readBytes": 2675396608,
          "writeBytes": 4426887168,
          "readSectors": 5225384,
          "writeSectors": 8646264,
          "readTime": 484272,
          "writeTime": 5987632,
          "iopsInProgress": 0,
          "ioTime": 252388,
          "weightedIoTime": 6503264,
          "name": "sda2",
          "serialNumber": "36234567890abcde01916b0e132b97a5c",
          "speedIORead": 0,
          "speedByteRead": 0,
          "speedIOWrite": 0,
          "speedByteWrite": 0,
          "util": 0,
          "avgrq_sz": 0,
          "avgqu_sz": 0,
          "await": 0,
          "svctm": 0
        }
      },
      "partition": [
        {
          "device": "/dev/root",
          "mountpoint": "/",
          "fstype": "ext4",
          "opts": "rw,noatime,data=ordered"
        },
        {
          "device": "/dev/sda4",
          "mountpoint": "/data",
          "fstype": "ext4",
          "opts": "rw,noatime,data=ordered"
        },
        {
          "device": "/dev/sda3",
          "mountpoint": "/usr/local",
          "fstype": "ext4",
          "opts": "rw,noatime,data=ordered"
        }
      ],
      "usage": [
        {
          "path": "/",
          "fstype": "ext2/ext3",
          "total": 21003583488,
          "free": 12865982464,
          "used": 7047081984,
          "usedPercent": 33.55180789995296,
          "inodesTotal": 1310720,
          "inodesUsed": 102779,
          "inodesFree": 1207941,
          "inodesUsedPercent": 7.8414154052734375
        },
        {
          "path": "/data",
          "fstype": "ext2/ext3",
          "total": 249789050880,
          "free": 217206951936,
          "used": 19869900800,
          "usedPercent": 7.954672444608313,
          "inodesTotal": 15499264,
          "inodesUsed": 105122,
          "inodesFree": 15394142,
          "inodesUsedPercent": 0.678238657009778
        },
        {
          "path": "/usr/local",
          "fstype": "ext2/ext3",
          "total": 21003583488,
          "free": 18666192896,
          "used": 1246871552,
          "usedPercent": 5.936470568045574,
          "inodesTotal": 1310720,
          "inodesUsed": 1404,
          "inodesFree": 1309316,
          "inodesUsedPercent": 0.10711669921875
        }
      ]
    },
    "mem": {
      "meminfo": {
        "total": 67246661632,
        "available": 16814030848,
        "used": 50432630784,
        "usedPercent": 74.99648244248473,
        "free": 3770302464,
        "active": 53595607040,
        "inactive": 5542129664,
        "wired": 0,
        "buffers": 759742464,
        "cached": 10508464128,
        "writeback": 0,
        "dirty": 2056192,
        "writebacktmp": 0
      },
      "vmstat": {
        "total": 2139090944,
        "used": 1053192192,
        "free": 1085898752,
        "usedPercent": 49.23550328489446,
        "sin": 2672574464,
        "sout": 4426887168
      }
    },
    "net": {
      "interface": [
        {
          "mtu": 65536,
          "name": "lo",
          "hardwareaddr": "",
          "flags": [
            "up",
            "loopback"
          ],
          "addrs": [
            {
              "addr": ""
            }
          ]
        },
        {
          "mtu": 1500,
          "name": "eth0",
          "hardwareaddr": "",
          "flags": [
            "broadcast",
            "multicast"
          ],
          "addrs": []
        },
        {
          "mtu": 1500,
          "name": "eth1",
          "hardwareaddr": "",
          "flags": [
            "up",
            "broadcast",
            "multicast"
          ],
          "addrs": [
            {
              "addr": ""
            }
          ]
        }
      ],
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
          "packetsRecv": 0,
          "errin": 0,
          "errout": 0,
          "dropin": 0,
          "dropout": 0,
          "fifoin": 0,
          "fifoout": 0
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
          "packetsRecv": 44093455831,
          "errin": 0,
          "errout": 0,
          "dropin": 0,
          "dropout": 0,
          "fifoin": 0,
          "fifoout": 0
        },
        {
          "name": "l",
          "speedSent": 13570,
          "speedRecv": 13570,
          "speedPacketsSent": 96,
          "speedPacketsRecv": 96,
          "bytesSent": 158794029730,
          "bytesRecv": 158794029730,
          "packetsSent": 311070174,
          "packetsRecv": 311070174,
          "errin": 0,
          "errout": 0,
          "dropin": 0,
          "dropout": 0,
          "fifoin": 0,
          "fifoout": 0
        }
      ],
      "netstat": {
        "established": 6016,
        "syncSent": 0,
        "synRecv": 0,
        "finWait1": 0,
        "finWait2": 0,
        "timeWait": 1,
        "close": 0,
        "closeWait": 24,
        "lastAck": 0,
        "listen": 46,
        "closing": 0
      },
      "protocolstat": {
        "udp": {
          "inCsumErrors": 0,
          "inDatagrams": 4696,
          "inErrors": 0,
          "noPorts": 0,
          "outDatagrams": 4698,
          "rcvbufErrors": 0,
          "sndbufErrors": 0
        }
      }
    },
    "system": {
      "info": {
        "hostname": "bkdata-111",
        "uptime": 3737847,
        "bootTime": 1510553867,
        "procs": 794,
        "os": "linux",
        "platform": "centos",
        "platformFamily": "rhel",
        "platformVersion": "7.2",
        "kernelVersion": "3.10.106-1-tlinux2-0044",
        "virtualizationSystem": "",
        "virtualizationRole": "",
        "hostid": "0C7CA213-D21D-B211-B859-03242340",
        "systemtype": "64-bit"
      },
      "docker": {
        "Client": {
          "Version": "",
          "ApiVersion": "",
          "GoVersion": ""
        },
        "Server": {
          "Version": "",
          "ApiVersion": "",
          "GoVersion": ""
        }
      }
    }
  },
  "ip": "127.0.0.1",
  "bizid": 0,
  "cloudid": 0
}`
