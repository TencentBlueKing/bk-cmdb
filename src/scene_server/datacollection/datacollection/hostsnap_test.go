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

package datacollection

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
)

func TestNeedToUpdate(t *testing.T) {
	need := needToUpdate(map[string]interface{}{"name": "a"}, map[string]interface{}{"name": "a"})
	if need {
		t.Errorf("not neet to update but got %v", need)
	}
	need = needToUpdate(map[string]interface{}{"name": "a"}, map[string]interface{}{"name": "b"})
	if !need {
		t.Errorf("not neet to update but got %v", need)
	}
	need = needToUpdate(map[string]interface{}{"name": 1}, map[string]interface{}{"name": 1})
	if need {
		t.Errorf("not neet to update but got %v", need)
	}
}

func TestGetIPS(t *testing.T) {
	val := gjson.Parse(example)
	ips := getIPS(&val)
	assert.Equal(t, []string{"127.0.0.1", "127.0.0.2", "127.0.0.3"}, ips)
}

func TestGetSetter(t *testing.T) {
	val := gjson.Parse(example)
	actual := parseSetter(&val, "127.0.0.1", "127.0.0.2")
	expected := map[string]interface{}{
		"bk_cpu":        int64(1),
		"bk_disk":       int64(10079),
		"bk_mem":        int64(15794),
		"bk_os_version": "6.2",
		"bk_host_name":  "test-master-02",
		"bk_outer_mac":  "",
		"bk_cpu_module": "Intel(R) Xeon(R) CPU E3-1230 V2 @ 3.30GHz",
		"bk_cpu_mhz":    int64(3301),
		"bk_os_type":    "Linux",
		"bk_os_name":    "linux centos",
		"bk_mac":        "28:31:52:1d:c6:0a",
	}
	for k, v := range expected {
		if actual[k] != v {
			t.Errorf("key %s not equal, expect %v but got %v", k, v, actual[k])
		}
	}
}

var example = `{
    "bizid":0,
    "cloudid":0,
    "data":{
        "city":"Shanghai",
        "country":"Asia",
        "cpu":{
            "cpuinfo":[
                {
                    "cacheSize":8192, // KB
                    "coreID":"0",
                    "cores":1,
                    "cpu":0,
                    "family":"6",
                    "flags":[
                        "fpu"
                    ],
                    "mhz":3301,  // MHz
                    "model":"58",
                    "modelName":"Intel(R) Xeon(R) CPU E3-1230 V2 @ 3.30GHz",
                    "physicalID":"0",
                    "stepping":9,
                    "vendorID":"GenuineIntel"
                }
            ],
            "per_stat":[
                {
                    "cpu":"cpu0",
                    "guest":0,  // s
                    "guestNice":0,  // s
                    "idle":7106735.46,  // s
                    "iowait":63491.51,  // s
                    "irq":0.93,  // s
                    "nice":0.45,  // s
                    "softirq":2700.54,  // s
                    "steal":0,  // s
                    "stolen":0,  // s
                    "system":30454.17,  // s
                    "user":247246.22  // s
                }
            ],
            "per_usage":[
                0.6835611871116047  // %
            ],
            "total_stat":{
                "cpu":"cpu-total",
                "guest":0,
                "guestNice":0,
                "idle":57401853.66,
                "iowait":109319.42,
                "irq":27.86,
                "nice":1.66,
                "softirq":15461.68,
                "steal":0,
                "stolen":0,
                "system":202392.7,
                "user":1877224.77
            },
            "total_usage":1.1292372440302074
        },
        "datetime":"2017-05-26 17:30:28",
        "disk":{
            "diskstat":{
                "sda":{
                    "avgqu_sz":0.0037333333333333333,  //  平均I/O队列长度
                    "avgrq_sz":21.134328358208954,  // sectors/io
                    "await":0.835820895522388, // ms/io
                    "ioTime":94182716,  // ms
                    "iopsInProgress":0,
                    "mergedReadCount":693911,  // count
                    "mergedWriteCount":1929503164,  // count
                    "name":"sda",
                    "readBytes":109136587776,  // Byte
                    "readCount":2106909,  // count
                    "readSectors":213157398,  // count
                    "readTime":45955864,  // ms
                    "serialNumber":"ST1000NM0011_Z1N36QQ3",
                    "speedByteRead":0,  // Byte/s
                    "speedByteWrite":48332.8,  // Byte/s
                    "speedIORead":0,  // count/s
                    "speedIOWrite":4.466666666666667,  // count/s
                    "svctm":0.26865671641791045,  // ms/io
                    "util":0,  // (0~1)
                    "weightedIoTime":55158904,  // ms
                    "writeBytes":8204215267328,  // Byte
                    "writeCount":73379253,  // count
                    "writeSectors":16023857944,  // count
                    "writeTime":9042312  // ms
                }
            },
            "partition":[
                {
                    "device":"/dev/sda1",
                    "fstype":"ext3",
                    "mountpoint":"/",
                    "opts":"rw,noatime,acl,user_xattr"
                }
            ],
            "usage":[
                {
                    "free":5730779136,  // Byte
                    "fstype":"ext2/ext3",
                    "inodesFree":572476,  // count
                    "inodesTotal":655360,  // count
                    "inodesUsed":82884,  // count
                    "inodesUsedPercent":12.6470947265625,
                    "path":"/",
                    "total":10568916992,  // Byte
                    "used":4838137856,  // Byte
                    "usedPercent":45.77704470251931
                }
            ]
        },
        "env":{
            "crontab":[
                {
                    "content":"*/5 * * * * /usr/local/agenttools/agent/check_tmp_agent.sh >/dev/null 2>&1 ",
                    "user":"root"
                }
            ],
            "host":"127.0.0.1 TENCENT64 "
        },
        "load":{
            "load_avg":{
                "load1":0.06,
                "load15":0.04,
                "load5":0.05
            }
        },
        "mem":{
            "meminfo":{
                "active":5803778048,  // Byte
                "available":14617051136,  // Byte
                "buffers":460066816,  // Byte
                "cached":9954938880,  // Byte
                "dirty":348160,  // Byte
                "free":4202045440,  // Byte
                "inactive":6059380736,  // Byte
                "total":16561496064,  // Byte
                "used":1944444928,  // Byte
                "usedPercent":11.740756514302305,
                "wired":0,  // Byte
                "writeback":0,  // Byte
                "writebacktmp":0  // Byte
            },
            "vmstat":{
                "free":1268301824,  // Byte
                "sin":3736612864,  // Byte
                "sout":5425745920,  // Byte
                "total":2139086848,  // Byte
                "used":870785024,  // Byte
                "usedPercent":40.70825945258675
            }
        },
        "net":{
            "dev":[
                {
                    "bytesRecv":10797513221,  // Byte
                    "bytesSent":106912772,  // Byte
                    "dropin":0,  // count
                    "dropout":0,  // count
                    "errin":0,  // count
                    "errout":0,  // count
                    "fifoin":0,  // count
                    "fifoout":0,  // count
                    "name":"eth0",
                    "packetsRecv":137587223,  // count
                    "packetsSent":743070,  // count
                    "speedPacketsRecv":0,  // count/s
                    "speedPacketsSent":0,  // count/s
                    "speedRecv":34,  // Byte/s
                    "speedSent":0  // Byte/s
                }
            ],
            "interface":[
                {
                    "addrs":[
                        {
                            "addr":"127.0.0.2/25"
                        }
                    ],
                    "flags":[
                        "up",
                        "broadcast",
                        "multicast"
                    ],
                    "hardwareaddr":"28:31:52:1d:c6:0a",
                    "mtu":1500,
                    "name":"eth0"
                },
                {
                    "addrs":[
                        {
                            "addr":"127.0.0.1/25"
                        }
                    ],
                    "flags":[
                        "up",
                        "broadcast",
                        "multicast"
                    ],
                    "hardwareaddr":"28:31:52:1d:c6:0a",
                    "mtu":1500,
                    "name":"eth0"
                }
            ]
        },
        "system":{
            "info":{
                "bootTime":1488339525,  // 启动时间
                "hostid":"BABD2A14-D21D-B211-B859-000000821800", // 启动生成的随机id
                "hostname":"test-master-02",
                "kernelVersion":"2.6.32.43-tlinux-1.0.23-state", // 内核版本
                "os":"linux",
                "platform":"centos",
                "platformFamily":"rhel",
                "platformVersion":"6.2",
                "procs":304, // 进程数
                "uptime":7451503, // 系统启动到现在的时间
                "virtualizationRole":"", // 虚拟机角色，guest，host
                "virtualizationSystem":"" // 虚拟机系统
            }
        },
        "timezone":8,
        "utctime":"2017-05-26 09:30:28"
    },
    "dataid":6381,
    "ip":"127.0.0.1",
    "type":"BaseReportBeat"
}`
