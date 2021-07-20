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

	"configcenter/src/apimachinery/rest"
	"configcenter/src/common/metadata"
	"configcenter/src/thirdparty/esbserver/esbutil"
)

type NodeManClientInterface interface {
	SearchPackage(ctx context.Context, h http.Header, processname string) (resp *SearchPluginPackageResult, err error)
	SearchProcess(ctx context.Context, h http.Header, processname string) (resp *SearchPluginProcessResult, err error)
	SearchProcessInfo(ctx context.Context, h http.Header, processname string) (resp *SearchPluginProcessInfoResult, err error)
	UpgradePlugin(ctx context.Context, h http.Header, bizID string, data *UpgradePluginRequest) (resp *UpgradePluginResult, err error)
	SearchTask(ctx context.Context, h http.Header, bizID int64, taskID int64) (resp *SearchTaskResult, err error)
	SearchPluginHost(ctx context.Context, h http.Header, processname string) (resp *SearchPluginHostResult, err error)
}

func NewNodeManClientInterface(client rest.ClientInterface, config *esbutil.EsbConfigSrv) NodeManClientInterface {
	return &nodeman{
		client: client,
		config: config,
	}
}

type nodeman struct {
	config *esbutil.EsbConfigSrv
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
//    "location": "http://127.0.0.1/download",
//    "pkg_ctime": "2018-09-13 22:17:19",
//    "pkg_path": "/data/bkee/miniweb/download",
//    "os": "linux",
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
	OS       string `json:"os"`
	ID       int    `json:"id"`
	MD5      string `json:"md5"`
}

// PluginProcess  define
// {
//    "category": "official",
//    "auto_launch": true,
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
	IsBinary     bool   `json:"is_binary"`
	AutoLaunch   bool   `json:"auto_launch"`
	UseDB        bool   `json:"use_db"`
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
//     "bk_cloud_id": "0",                     云区域id
//     "node_type": "PLUGIN",                  操作对象  PLUGIN
//     "op_type": "UPDATE",                    操作   更新: UPDATE
//     "global_params": {
//         "plugin": {
//             "id": 10,
//             "name": "xxx",                    插件名称
//             "scenario": "网络设备采集",
//             "category": "official",            插件分类
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
//             "keep_config": 0,            是否保留原配置文件  1: 保留(勾选)  0：不保留(不勾选)
//             "no_restart":  0,            更新后是否重启      1: 不重启(勾选)  0：重启(不勾选)
//             "no_delegate": 0             下发后不托管        1：不托管(勾选)  0：托管(不勾选)
//         },
//         "upgrade_type": "APPEND",       覆盖方式  "APPEND": 增量更新(仅覆盖)  "OVERRIDE": 覆盖更新(先删除原目录后覆盖)
//         "configs": [{
//             "inner_ips": ["127.0.0.1"],            支持多台机器使用同一配置文件 机器ip必须在hosts参数中存在 否则不操作
//             "content": ""
//         }]
//     },
//     "hosts": [{"inner_ips": "127.0.0.1"}]
// }
type UpgradePluginRequest struct {
	Creator      string `json:"creator"`
	BkCloudID    string `json:"bk_cloud_id"`
	NodeType     string `json:"node_type"` // 操作对象  PLUGIN
	OpType       string `json:"op_type"`   // 操作   更新: UPDATE
	GlobalParams struct {
		Plugin  *PluginProcess     `json:"plugin"`
		Package *PluginPackage     `json:"package"`
		Control *PluginProcessInfo `json:"control"`
		Option  struct {
			KeepConfig int `json:"keep_config"` // 是否保留原配置文件  1: 保留(勾选)  0：不保留(不勾选)
			NoRestart  int `json:"no_restart"`  // 更新后是否重启     1: 不重启(勾选)  0：重启(不勾选)
			NoDelegate int `json:"no_delegate"` // 下发后不托管       1：不托管(勾选)  0：托管(不勾选)
		} `json:"option"`
		UpgradeType string                `json:"upgrade_type"` //  覆盖方式  "APPEND": 增量更新(仅覆盖)  "OVERRIDE": 覆盖更新(先删除原目录后覆盖)
		Configs     []UpgradePluginConfig `json:"configs"`
	} `json:"global_params"`
	Hosts []UpgradePluginHostField `json:"hosts"`
}

type UpgradePluginConfig struct {
	InnerIPs []string `json:"inner_ips"` // 支持多台机器使用同一配置文件 机器ip必须在hosts参数中存在 否则不操作
	Content  string   `json:"content,omitempty"`
}

type UpgradePluginHostField struct {
	InnerIPs string `json:"inner_ips"` // 支持多台机器使用同一配置文件 机器ip必须在hosts参数中存在 否则不操作
}

type esbUpgradePluginParams struct {
	*esbutil.EsbCommParams
	*UpgradePluginRequest
}

type UpgradePluginResult struct {
	ESBBaseResult
	Data struct {
		ID int64 `json:"id"`
	} `json:"data"`
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
// 	  "creator": "xxxx",
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
// 			"inner_ip": "127.0.0.1",
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
// 		  "location": "http://127.0.0.1/download",
// 		  "pkg_ctime": "2018-09-19 15:51:10",
// 		  "pkg_path": "/data/bkee/miniweb/download",
// 		  "id": "5",
// 		  "md5": "5bbcfedc7c9cb6a5fbf1a3459fd12c24"
// 		},
// 		"upgrade_type": "APPEND",
// 		"configs": [
// 		  {
// 			"content": "",
// 			"inner_ips": [
// 			  "127.0.0.1"
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
		Status string `json:"status"` // QUEUE: 队列等待中 RUNNING: 执行中 SUCCESS: 执行成功 FAILED: 执行失败
		Step   string `json:"step"`
		Host   struct {
			BkBizID   string `json:"bk_biz_id"`
			BkCloudID string `json:"bk_cloud_id"`
			OuterIP   string `json:"outer_ip"`
			NodeType  string `json:"node_type"`
			InnerIP   string `json:"inner_ip"`
			HasCygwin bool   `json:"has_cygwin"`
			OsType    string `json:"os_type"`
			ID        int64  `json:"id"`
		} `json:"host"`
		JobID   string `json:"job_id"`
		ErrCode string `json:"err_code"`
	} `json:"hosts"`
	StartTime metadata.Time `json:"start_time"`
	EndTime   string        `json:"end_time"`
	ID        int64         `json:"id"`
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
//                 "inner_ip": "127.0.0.1",
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
	Status   string `json:"status"` // UNREGISTER RUNNING TERMINATED
	Host     Host   `json:"host"`
	Version  string `json:"version"`
	Name     string `json:"name"`
	procType string `json:"proc_type"`
}
