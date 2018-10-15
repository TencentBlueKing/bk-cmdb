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
package nodeman

import (
	"context"
	"net/http"
	"time"

	"configcenter/src/apimachinery/rest"
	"configcenter/src/thirdpartyclient/esbserver/esbutil"
)

type NodeManClientInterface interface {
	SearchPackage(ctx context.Context, h http.Header, processname string) (resp *SearchPluginPackageResult, err error)
	SearchProcess(ctx context.Context, h http.Header, processname string) (resp *SearchPluginProcessResult, err error)
	SearchProcessInfo(ctx context.Context, h http.Header, processname string) (resp *SearchPluginProcessInfoResult, err error)
	UpgradePlugin(ctx context.Context, h http.Header, bizID string, data *UpgradePluginRequest) (resp *UpgradePluginResult, err error)
	SearchTask(ctx context.Context, h http.Header, bizID string, taskID string) (resp *SearchTaskResult, err error)
	SearchPluginHost(ctx context.Context, h http.Header, processname string) (resp *SearchPluginHostResult, err error)
}

func NewNodeManClientInterface(client rest.ClientInterface, config *esbutil.EsbConfigServ) NodeManClientInterface {
	return &nodeman{
		client: client,
		config: config,
	}
}

type nodeman struct {
	config *esbutil.EsbConfigServ
	client rest.ClientInterface
}

type ESBBaseResult struct {
	Message   string `json:"message"`
	Code      string `json:"code"`
	Result    bool   `json:"result"`
	RequestID string `json:"request_id"`
}

type SearchPluginPackageResult struct {
	ESBBaseResult
	Data []PluginPackage `json:"data"`
}

type SearchPluginProcessResult struct {
	ESBBaseResult
	Data PluginProcess `json:"data"`
}

type SearchPluginProcessInfoResult struct {
	ESBBaseResult
	Data PluginProcessInfo `json:"data"`
}

// PluginPackage define
// {
//    "pkg_mtime": "2018-09-13 22:17:19",
//    "module": "gse_plugin",
//    "project": "basereport",
//    "pkg_size": 4560103,
//    "version": "10.1.5",
//    "pkg_name": "basereport-10.1.5.tgz",
//    "location": "http://10.167.77.15/download",
//    "pkg_ctime": "2018-09-13 22:17:19",
//    "pkg_path": "/data/bkee/miniweb/download",
//    "id": 1,
//    "md5": "b17b48158fe73a4b58937b8ef096b94e"
// }
type PluginPackage struct {
	PkgMtime string `json:"pkg_mtime"`
	Module   string `json:"module"`
	Project  string `json:"project"`
	PkgSize  int    `json:"pkg_size"`
	Version  string `json:"version"`
	PkgName  string `json:"pkg_name"`
	Location string `json:"location"`
	PkgCtime string `json:"pkg_ctime"`
	PkgPath  string `json:"pkg_path"`
	ID       int    `json:"id"`
	MD5      string `json:"md5"`
}

// PluginProcess define
// {
//    "category": "official",
//    "config_file": "basereport.conf",
//    "name": "basereport",
//    "scenario": "CMDB上的实时数据，蓝鲸监控里的主机监控中的基础性能数据",
//    "is_binary": 1,
//    "use_db": 0,
//    "config_format": "yaml",
//    "id": 1
// }
type PluginProcess struct {
	Category     string `json:"category"`
	ConfigFile   string `json:"config_file"`
	Name         string `json:"name"`
	Scenario     string `json:"scenario"`
	IsBinary     int    `json:"is_binary"`
	UseDB        int    `json:"use_db"`
	ConfigFormat string `json:"config_format"`
	ID           int    `json:"id"`
}

// PluginProcessInfo define
// {
// 	   "stop_cmd": "./stop.sh basereport",
// 	   "pid_path": "/var/run/gse/basereport.pid
// 	   "log_path": "/var/log/gse",
// 	   "reload_cmd": "./reload.sh basereport",
// 	   "data_path": "/var/lib/gse",
// 	   "install_path": "/usr/local/gse/plugins/
// 	   "module": "gse_plugin",
// 	   "project": "basereport",
// 	   "start_cmd": "./start.sh basereport",
// 	   "id": 1,
// 	   "restart_cmd": "./restart.sh basereport"
// }
type PluginProcessInfo struct {
	StopCmd     string `json:"stop_cmd"`
	PidPath     string `json:"pid_path"`
	LogPath     string `json:"log_path"`
	ReloadCmd   string `json:"reload_cmd"`
	DataPath    string `json:"data_path"`
	InstallPath string `json:"install_path"`
	Module      string `json:"module"`
	Project     string `json:"project"`
	StartCmd    string `json:"start_cmd"`
	ID          int    `json:"id"`
	RestartCmd  string `json:"restart_cmd"`
}

// UpgradePluginRequest define
// {
//     "creator": "aaaaa",
//     "bk_cloud_id": "0",                     # 云区域id
//     "node_type": "PLUGIN",                  # 操作对象  PLUGIN
//     "op_type": "UPDATE",                    # 操作   更新: UPDATE
//     "global_params": {
//         "plugin": {
//             "id": 10,
//             "name": "xxx",                    # 插件名称
//             "scenario": "网络设备采集",
//             "category": "official",            # 插件分类
//             "config_file": "xxx.conf",
//             "config_file_format": "json",
//             "use_db": 0,
//             "is_binary": 1
//         },
//         "package": {
//             "id": 9,
//             "pkg_name": "xxx-yyy-1.0.0-x86_64.tgz",
//             "version": "1.0.0",
//             "module": "gse_plugin",
//             "project": "xxx",
//             "pkg_size": 5142819,
//             "pkg_path": "/x/y/z/",
//             "md5": "5bbcfedc7c9cb6a5fbf1a3459fd12c24",
//             "pkg_mtime": "2018-09-19 15:51:10",
//             "pkg_ctime": "2018-09-19 15:51:10",
//             "location": "http://127.0.0.1/download"
//         },
//         "control": {
//             "id": 6,
//             "module": "gse_plugin",
//             "project": "xxx",
//             "install_path": "",
//             "log_path": "",
//             "data_path": "",
//             "pid_path": "",
//             "start_cmd": "",
//             "stop_cmd": "",
//             "restart_cmd": "",
//             "reload_cmd": ""
//         },
//         "option": {
//             "keep_config": 0,            # 是否保留原配置文件  1: 保留(勾选)  0：不保留(不勾选)
//             "no_restart":  0,            # 更新后是否重启      1: 不重启(勾选)  0：重启(不勾选)
//             "no_delegate": 0             # 下发后不托管        1：不托管(勾选)  0：托管(不勾选)
//         },
//         "upgrade_type": "APPEND",       # 覆盖方式  "APPEND": 增量更新(仅覆盖)  "OVERRIDE": 覆盖更新(先删除原目录后覆盖)
//         "configs": [{
//             "inner_ips": ["127.0.0.1"],            # 支持多台机器使用同一配置文件 机器ip必须在hosts参数中存在 否则不操作
//             "content": ""
//         }]
//     },
//     "hosts": [{"inner_ips": "127.0.0.1"}]
// }
type UpgradePluginRequest struct {
	Creator      string `json:"creator"`
	BkCloudID    string `json:"bk_cloud_id"`
	NodeType     string `json:"node_type"` // # 操作对象  PLUGIN
	OpType       string `json:"op_type"`   // # 操作   更新: UPDATE
	GlobalParams struct {
		Plugin  *PluginProcess     `json:"plugin"`
		Package *PluginPackage     `json:"package"`
		Control *PluginProcessInfo `json:"control"`
		Option  struct {
			KeepConfig int `json:"keep_config"` // # 是否保留原配置文件  1: 保留(勾选)  0：不保留(不勾选)
			NoRestart  int `json:"no_restart"`  // # 更新后是否重启     1: 不重启(勾选)  0：重启(不勾选)
			NoDelegate int `json:"no_delegate"` // # 下发后不托管       1：不托管(勾选)  0：托管(不勾选)
		}
		UpgradeType string                `json:"upgrade_type"` //  # 覆盖方式  "APPEND": 增量更新(仅覆盖)  "OVERRIDE": 覆盖更新(先删除原目录后覆盖)
		Configs     []UpgradePluginConfig `json:"configs"`
	} `json:"global_params"`
	Hosts []UpgradePluginConfig `json:"hosts"`
}

type UpgradePluginConfig struct {
	InnerIPs []string `json:"inner_ips"` // # 支持多台机器使用同一配置文件 机器ip必须在hosts参数中存在 否则不操作
	Content  string   `json:"content,omitempty"`
}

type UpgradePluginResult struct {
	ESBBaseResult
	Data interface{} `json:"data"`
}

// SearchTaskResult define
// {
// 	"message": "success",
// 	"code": "OK",
// 	"data": {
// 	  "bk_biz_id": "2",
// 	  "host_count": 1,
// 	  "op_target": {
// 		"category": "official",
// 		"config_file": "netdevicebeat.conf",
// 		"name": "netdevicebeat",
// 		"scenario": "CMDB上的网络设备采集功能",
// 		"is_binary": 1,
// 		"use_db": 0,
// 		"config_format": "json",
// 		"id": 5
// 	  },
// 	  "os_count": {
// 		"WINDOWS": 0,
// 		"AIX": 0,
// 		"LINUX": 1
// 	  },
// 	  "creator": "drizztchen",
// 	  "job_type_desc": "更新PLUGIN",
// 	  "start_time": "2018-09-20 20:15:00",
// 	  "job_type": "UPDATE_PLUGIN",
// 	  "bk_cloud_id": "0",
// 	  "hosts": [
// 		{
// 		  "status": "FAILED",
// 		  "step": "任务执行失败(更新)",
// 		  "host": {
// 			"bk_biz_id": "2",
// 			"bk_cloud_id": "0",
// 			"outer_ip": null,
// 			"node_type": "AGENT",
// 			"inner_ip": "10.235.46.112",
// 			"has_cygwin": false,
// 			"os_type": "LINUX",
// 			"id": 1
// 		  },
// 		  "job_id": "2",
// 		  "err_code": "JOB_TIMEOUT"
// 		}
// 	  ],
// 	  "end_time": null,
// 	  "status_count": {
// 		"running_count": 0,
// 		"failed_count": 1,
// 		"success_count": 0
// 	  },
// 	  "global_params": {
// 		"control": {
// 		  "stop_cmd": "./stop.sh",
// 		  "project": "netdevicebeat",
// 		  "pid_path": "/var/run/gse",
// 		  "log_path": "/var/log/gse",
// 		  "data_path": "/var/run/gse",
// 		  "install_path": "/usr/local/gse",
// 		  "module": "gse_plugin",
// 		  "reload_cmd": "./reload.sh",
// 		  "start_cmd": "./start.sh",
// 		  "id": "5",
// 		  "restart_cmd": "./restart.sh"
// 		},
// 		"option": {
// 		  "no_restart": false,
// 		  "no_delegate": false,
// 		  "keep_config": false
// 		},
// 		"plugin": {
// 		  "category": "official",
// 		  "config_file": "netdevicebeat.conf",
// 		  "config_file_format": "json",
// 		  "scenario": "CMDB上的网络设备采集功能",
// 		  "id": "5",
// 		  "use_db": false,
// 		  "is_binary": true,
// 		  "name": "netdevicebeat"
// 		},
// 		"package": {
// 		  "pkg_mtime": "2018-09-19 15:51:10",
// 		  "module": "gse_plugin",
// 		  "project": "netdevicebeat",
// 		  "pkg_size": 5142819,
// 		  "version": "1.0.0",
// 		  "pkg_name": "netdevicebeat-linux-1.0.0-x86_64.tgz",
// 		  "location": "http://10.167.77.15/download",
// 		  "pkg_ctime": "2018-09-19 15:51:10",
// 		  "pkg_path": "/data/bkee/miniweb/download",
// 		  "id": "5",
// 		  "md5": "5bbcfedc7c9cb6a5fbf1a3459fd12c24"
// 		},
// 		"upgrade_type": "APPEND",
// 		"configs": [
// 		  {
// 			"content": "IyMjIyMjIyMjIyMjIyMjIyMjIyBOZXRkZXZpY2ViZWF0IENvbmZpZ3VyYXRpb24gRXhhbXBsZSAjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyM                                                                                                         jCgojIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyBOZXRkZXZpY2ViZWF0ICMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjCgpuZXRkZXZpY2Vi                                                                                                         ZWF0OiB7CiAgInBlcmlvZCI6ICIxcyIsCiAgImRhdGFpZCI6IDEwMTQsCiAgInNjYW5fcmFuZ2UiOiBbCiAgICAgICAgICAiMTkyLjE2OC4xLjEiLAogICAgICAgICAgI                                                                                                         jE5Mi4xNjguMS4yLTE5Mi4xNjguMS4yMDAiLAogICAgICAgICAgIjE5Mi4xNjguMS4xLzMyIiwKICBdLAogICJzbm1wIjogewogICAgICAiY29tbXVuaXR5IjogInB1Ym                                                                                                         xpYyIsCiAgICAgICJ0aW1lX291dCI6IDMwMCwKICAgICAgInJldHJ5IjogNSwKICAgICAgInBvcnQiOiAxNjEsCiAgICAgICJtYXhfb2lkcyI6IDEwLAogIH0sCiAgInB                                                                                                         pbmdfdGltZW91dCI6IDMwMCwKICAicGluZ19yZXRyeSI6IDMsCiAgIndvcmtlciI6IDMsCiAgIm1ldHJpY3MiOiBbCiAgICAgIHsKICAgICAgICAgICJhY3Rpb24iOiAi                                                                                                         Z2V0IiwKICAgICAgICAgICJwZXJpb2QiOiAiMjRIIiwKICAgICAgICAgICJvaWQiOiAiMS4zLjYuMS4yLjEuMi4yLjEuMi4xNiIKICAgICAgfSwKICAgICAgewogICAgI                                                                                                         CAgICAgImFjdGlvbiI6ICJ3YWxrIiwKICAgICAgICAgICJwZXJpb2QiOiAiMjRIIiwKICAgICAgICAgICJvaWQiOiAiMS4zLjYuMS4yLjEuMy4xLjEiCiAgICAgIH0KIC                                                                                                         BdLAogICJyZXBvcnQiOiB7CiAgICAgICJkZWJ1ZyI6IHRydWUKICB9Cn0KCiM9PT09PT09PT09PT09PT09PT09PT09PT09PT09PT09PSBPdXRwdXRzID09PT09PT09PT0                                                                                                         9PT09PT09PT09PT09PT09PT09PT09PT09PT0KIyBDb25maWd1cmUgb3V0cHV0cyB0byBnc2VhZ2VudApvdXRwdXQuZ3NlOiB7ZW5kcG9pbnQ6IC92YXIvcnVuL2lwYy5z                                                                                                         dGF0ZS5yZXBvcnR9CiMgbGludXggYWdlbnQgZGVmYXVsdCBpcGMgZmlsZSBjb25maWcKIyB3aW5kb3dzIGFnZW50IHNvY2tldCBjb25maWcKI2VuZHBvaW50OiAiL3Vzc                                                                                                         i9sb2NhbC9nc2UvZ3NlYWdlbnQvaXBjLnN0YXRlLnJlcG9ydCIKI2VuZHBvaW50OiAiMTI3LjAuMC4xOjQ3MDAwIgoKIyB0cnkgdGltZXMgYW5kIGludGVydmFsIHdoZW                                                                                                         4gZGlzY29ubmVjdCB0byBhZ2VudAojcmV0cnl0aW1lczogMwojcmV0cnlpbnRlcnZhbDogM3MKCiMgc2VuZCB0aW1lb3V0CiN3cml0ZXRpbWVvdXQ6IDVzCgojPT09PT0                                                                                                         9PT09PT09PT09PT09PT09PT09PT09PT09PT0gTG9nZ2luZyA9PT09PT09PT09PT09PT09PT09PT09PT09PT09PT09PT09PT09PQojIFRoZXJlIGFyZSB0aHJlZSBvcHRp                                                                                                         b25zIGZvciB0aGUgbG9nIG91dHB1dDogc3lzbG9nLCBmaWxlLCBzdGRlcnIuCiMgVW5kZXIgV2luZG93cyBzeXN0ZW1zLCB0aGUgbG9nIGZpbGVzIGFyZSBwZXIgZGVmY                                                                                                         XVsdCBzZW50IHRvIHRoZSBmaWxlIG91dHB1dCwKIyB1bmRlciBhbGwgb3RoZXIgc3lzdGVtIHBlciBkZWZhdWx0IHRvIHN5c2xvZy4KCiMgU2V0cyBsb2cgbGV2ZWwuIF                                                                                                         RoZSBkZWZhdWx0IGxvZyBsZXZlbCBpcyBpbmZvLgojIEF2YWlsYWJsZSBsb2cgbGV2ZWxzIGFyZTogY3JpdGljYWwsIGVycm9yLCB3YXJuaW5nLCBpbmZvLCBkZWJ1Zwo                                                                                                         jbG9nZ2luZy5sZXZlbDogaW5mbwoKIyBFbmFibGUgZGVidWcgb3V0cHV0IGZvciBzZWxlY3RlZCBjb21wb25lbnRzLiBUbyBlbmFibGUgYWxsIHNlbGVjdG9ycyB1c2Ug                                                                                                         WyIqIl0KIyBPdGhlciBhdmFpbGFibGUgc2VsZWN0b3JzIGFyZSAiYmVhdCIsICJwdWJsaXNoIiwgInNlcnZpY2UiCiMgTXVsdGlwbGUgc2VsZWN0b3JzIGNhbiBiZSBja                                                                                                         GFpbmVkLgojbG9nZ2luZy5zZWxlY3RvcnM6IFsgXQoKIyBTZW5kIGFsbCBsb2dnaW5nIG91dHB1dCB0byBzeXNsb2cuIFRoZSBkZWZhdWx0IGlzIGZhbHNlLgojbG9nZ2                                                                                                         luZy50b19zeXNsb2c6IHRydWUKCiMgSWYgZW5hYmxlZCwgZmlsZWJlYXQgcGVyaW9kaWNhbGx5IGxvZ3MgaXRzIGludGVybmFsIG1ldHJpY3MgdGhhdCBoYXZlIGNoYW5                                                                                                         nZWQKIyBpbiB0aGUgbGFzdCBwZXJpb2QuIEZvciBlYWNoIG1ldHJpYyB0aGF0IGNoYW5nZWQsIHRoZSBkZWx0YSBmcm9tIHRoZSB2YWx1ZSBhdAojIHRoZSBiZWdpbm5p                                                                                                         bmcgb2YgdGhlIHBlcmlvZCBpcyBsb2dnZWQuIEFsc28sIHRoZSB0b3RhbCB2YWx1ZXMgZm9yCiMgYWxsIG5vbi16ZXJvIGludGVybmFsIG1ldHJpY3MgYXJlIGxvZ2dlZ                                                                                                         CBvbiBzaHV0ZG93bi4gVGhlIGRlZmF1bHQgaXMgdHJ1ZS4KI2xvZ2dpbmcubWV0cmljcy5lbmFibGVkOiB0cnVlCgojIFRoZSBwZXJpb2QgYWZ0ZXIgd2hpY2ggdG8gbG                                                                                                         9nIHRoZSBpbnRlcm5hbCBtZXRyaWNzLiBUaGUgZGVmYXVsdCBpcyAzMHMuCiNsb2dnaW5nLm1ldHJpY3MucGVyaW9kOiAzMHMKCiMgTG9nZ2luZyB0byByb3RhdGluZyB                                                                                                         maWxlcyBmaWxlcy4gU2V0IGxvZ2dpbmcudG9fZmlsZXMgdG8gZmFsc2UgdG8gZGlzYWJsZSBsb2dnaW5nIHRvCiMgZmlsZXMuCmxvZ2dpbmcudG9fZmlsZXM6IHRydWUK                                                                                                         bG9nZ2luZy5maWxlczoKICAjIENvbmZpZ3VyZSB0aGUgcGF0aCB3aGVyZSB0aGUgbG9ncyBhcmUgd3JpdHRlbi4gVGhlIGRlZmF1bHQgaXMgdGhlIGxvZ3MgZGlyZWN0b                                                                                                         3J5CiAgIyB1bmRlciB0aGUgaG9tZSBwYXRoICh0aGUgYmluYXJ5IGxvY2F0aW9uKS4KICAjcGF0aDogL3Zhci9sb2cvZmlsZWJlYXQKCiAgIyBUaGUgbmFtZSBvZiB0aG                                                                                                         UgZmlsZXMgd2hlcmUgdGhlIGxvZ3MgYXJlIHdyaXR0ZW4gdG8uCiAgI25hbWU6IGZpbGViZWF0CgogICMgQ29uZmlndXJlIGxvZyBmaWxlIHNpemUgbGltaXQuIElmIGx                                                                                                         pbWl0IGlzIHJlYWNoZWQsIGxvZyBmaWxlIHdpbGwgYmUKICAjIGF1dG9tYXRpY2FsbHkgcm90YXRlZAogICNyb3RhdGVldmVyeWJ5dGVzOiAxMDQ4NTc2MCAjID0gMTBN                                                                                                         QgoKICAjIE51bWJlciBvZiByb3RhdGVkIGxvZyBmaWxlcyB0byBrZWVwLiBPbGRlc3QgZmlsZXMgd2lsbCBiZSBkZWxldGVkIGZpcnN0LgogICNrZWVwZmlsZXM6IDcKC                                                                                                         nBhdGguZGF0YTogL3Zhci9saWIvZ3NlCnBhdGgubG9nczogL3Zhci9sb2cvZ3NlCnBhdGgucGlkOiAvdmFyL3J1bi9nc2UK",
// 			"inner_ips": [
// 			  "10.235.46.112"
// 			]
// 		  }
// 		]
// 	  },
// 	  "id": 2
// 	},
// 	"result": true,
// 	"request_id": "28dc639c2b16495791f204df5161a73a"
//   }
type SearchTaskResult struct {
	ESBBaseResult
	Data Task `json:"data"`
}

type Task struct {
	BkBizID   string `json:"bk_biz_id"`
	BkCloudID string `json:"bk_cloud_id"`
	Hosts     []struct {
		Status string `json:"status"`
		Step   string `json:"step"`
		Host   struct {
			BkBizID   string `json:"bk_biz_id"`
			BkCloudID string `json:"bk_cloud_id"`
			OuterIP   string `json:"outer_ip"`
			NodeType  string `json:"node_type"`
			InnerIP   string `json:"inner_ip"`
			HasCygwin bool   `json:"has_cygwin"`
			OsType    string `json:"os_type"`
			ID        string `json:"id"`
		} `json:"host"`
		JobID   string `json:"job_id"`
		ErrCode string `json:"err_code"`
	} `json:"hosts"`
	StartTime time.Time `json:""start_time""`
	EndTime   string    `json:"end_time"`
	ID        string    `json:"id"`
}

type Host struct {
	BkBizID   int64  `json:"bk_biz_id,string"`
	BkCloudID int64  `json:"bk_cloud_id,string"`
	OuterIP   string `json:"outer_ip"`
	NodeType  string `json:"node_type"`
	InnerIP   string `json:"inner_ip"`
	HasCygwin bool   `json:"has_cygwin"`
	OsType    string `json:"os_type"`
	ID        int64  `json:"id"`
}

// SearchPluginHostResult define
// {
//     "message": "success",
//     "code": "OK",
//     "data": [
//         {
//             "status": "RUNNING",
//             "host": {
//                 "bk_biz_id": "2",
//                 "bk_cloud_id": "0",
//                 "outer_ip": null,
//                 "node_type": "AGENT",
//                 "inner_ip": "1.2.3.2",
//                 "has_cygwin": false,
//                 "os_type": "LINUX",
//                 "id": 1
//             },
//             "version": "V0.0",
//             "name": "xxxxxx",
//             "proc_type": "AGENT"
//         }
//     ],
//     "result": true,
//     "request_id": "xxxxxxxxxx"
// }
type SearchPluginHostResult struct {
	ESBBaseResult
	Data []PluginHost `json:"data"`
}

type PluginHost struct {
	Status   string `json:"status"`
	Host     Host   `json:"host"`
	Version  string `json:"version"`
	Name     string `json:"name"`
	procType string `json:"proc_type"`
}
