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

package data

import (
	"configcenter/src/common"
	mCommon "configcenter/src/scene_server/admin_server/common"
	"configcenter/src/source_controller/api/metadata"
)

var (
	groupBaseInfo string = mCommon.BaseInfo
)

var Distribution = "community" // could be community or enterprise

/*
	&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "", PropertyName: "", IsRequired: , IsOnly: , PropertyGroup: , PropertyType: , Option: ""},
*/

func AppRow() []*metadata.ObjectAttDes {
	objID := common.BKInnerObjIDApp

	groupAppRole := mCommon.AppRole

	lifeCycleOption := mCommon.GetOptionStr([]interface{}{common.KvMap{"name": "测试中", "type": "text"}, common.KvMap{"name": "已上线", "type": "text", "is_default": true}, common.KvMap{"name": "停运", "type": "text"}})
	languageOption := mCommon.GetOptionStr([]interface{}{common.KvMap{"name": "中文", "type": "text"}, common.KvMap{"name": "English", "type": "text"}})
	dataRows := []*metadata.ObjectAttDes{
		//&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_biz_id", PropertyName: "业务ID", IsSystem: true, IsRequired: bkFalse, IsOnly: false, PropertyType: common.FiledTypeInt, Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_biz_name", PropertyName: "业务名", IsRequired: true, IsOnly: true, Editable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FiledTypeSingleChar, Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "life_cycle", PropertyName: "生命周期", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FiledTypeEnum, Option: lifeCycleOption},

		//role
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: common.BKMaintainersField, PropertyName: "运维人员", IsRequired: true, IsOnly: false, Editable: true, PropertyGroup: groupAppRole, PropertyType: common.FiledTypeUser, Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: common.BKProductPMField, PropertyName: "产品人员", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: groupAppRole, PropertyType: common.FiledTypeUser, Option: ""},
		//&metadata.ObjectAttDes{ObjectID: "app", PropertyID: "bk_supplier_account", PropertyName: "域", IsRequired: false, IsOnly: false, PropertyType: common.FiledTypeSingleChar, Option: ""},
		//&metadata.ObjectAttDes{ObjectID: "app", PropertyID: common.BKInstParentStr, PropertyName: "", IsRequired: false, IsOnly: false, PropertyType: common.FiledTypeInt, Option: ""},

		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: common.BKTesterField, PropertyName: "测试人员", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: groupAppRole, PropertyType: common.FiledTypeUser, Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_biz_developer", PropertyName: "开发人员", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: groupAppRole, PropertyType: common.FiledTypeUser, Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: common.BKOperatorField, PropertyName: "操作人员", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: groupAppRole, PropertyType: common.FiledTypeUser, Option: ""},

		//&metadata.ObjectAttDes{ObjKeyId: "app", PropertyID: "OnwerID", PropertyName: "域", IsRequired: common.BKFalse, IsOnly: common.BKFalse, PropertyType: common.FiledTypeSingleChar, Option: ""},
		//&metadata.ObjectAttDes{ObjKeyId: "app", PropertyID: "ParentID", PropertyName: "", IsRequired: common.BKFalse, IsOnly: common.BKFalse, PropertyType: common.FiledTypeInt, Option: ""},
	}

	if Distribution == common.RevisionEnterprise {
		dataRows = append(dataRows,
			&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "time_zone", PropertyName: "时区", IsRequired: true, IsOnly: false, Editable: false, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeTimeZone, Option: "", IsReadOnly: true},
			&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "language", PropertyName: "语言", IsRequired: true, IsOnly: false, PropertyGroup: groupBaseInfo, PropertyType: common.FiledTypeEnum, Option: languageOption, IsReadOnly: true},
		)
	}

	return dataRows

}

func SetRow() []*metadata.ObjectAttDes {
	objID := common.BKInnerObjIDSet
	serviceStatusOption := mCommon.GetOptionStr([]interface{}{common.KvMap{"name": "开放", "type": "text", "is_default": true}, common.KvMap{"name": "关闭", "type": "text"}})

	dataRows := []*metadata.ObjectAttDes{
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: common.BKAppIDField, PropertyName: "业务ID", IsAPI: true, IsRequired: false, IsOnly: true, PropertyGroup: groupBaseInfo, PropertyType: common.FiledTypeInt, Option: "{}"},
		//&metadata.ObjectAttDes{ObjectID: objID, PropertyID: common.SETID_FIELD, PropertyName: "集群ID", IsSystem: true, IsRequired: bkFalse, IsOnly: bkTrue, PropertyGroup: groupBaseInfo, PropertyType: common.FiledTypeSingleChar, Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_set_name", PropertyName: "集群名字", IsRequired: true, IsOnly: true, Editable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FiledTypeSingleChar, Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_set_desc", PropertyName: "集群描述", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FiledTypeSingleChar, Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_set_env", PropertyName: "环境类型", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FiledTypeEnum, Option: "[{\"name\":\"测试\",\"type\":\"text\"},{\"name\":\"体验\",\"type\":\"text\"},{\"name\":\"正式\",\"type\":\"text\",\"is_default\":true}]"},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_service_status", PropertyName: "服务状态", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FiledTypeEnum, Option: serviceStatusOption},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "description", PropertyName: "备注", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FiledTypeLongChar, Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_capacity", PropertyName: "设计容量", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FiledTypeInt, Option: `{"min":"1","max":"999999999"}`},

		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: common.BKChildStr, PropertyName: "", IsRequired: false, IsOnly: false, PropertyType: "", Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: common.BKInstParentStr, PropertyName: "", IsSystem: true, IsRequired: true, IsOnly: true, PropertyType: common.FiledTypeInt, Option: ""},
	}
	return dataRows
}

func ModuleRow() []*metadata.ObjectAttDes {
	objID := common.BKInnerObjIDModule
	moduleTypeOption := mCommon.GetOptionStr([]interface{}{common.KvMap{"name": "普通", "type": "text", "is_default": true}, common.KvMap{"name": "数据库", "type": "text"}})

	dataRows := []*metadata.ObjectAttDes{
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: common.BKAppIDField, PropertyName: "业务ID", IsAPI: true, IsRequired: false, IsOnly: true, PropertyGroup: groupBaseInfo, PropertyType: common.FiledTypeInt, Option: "{}"},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: common.BKSetIDField, PropertyName: "集群ID", IsAPI: true, IsRequired: false, IsOnly: true, PropertyGroup: groupBaseInfo, PropertyType: common.FiledTypeInt, Option: "{}"},
		//&metadata.ObjectAttDes{ObjectID: objID, PropertyID: common.MODULEID_FIELD, PropertyName: "模块ID", IsSystem: true, IsRequired: bkFalse, IsOnly: false, PropertyGroup: groupBaseInfo, PropertyType: common.FiledTypeInt, Option: "{}"},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: common.BKModuleNameField, PropertyName: "模块名", IsRequired: true, IsOnly: true, Editable: true, PropertyType: common.FiledTypeSingleChar, Option: ""},
		//&metadata.ObjectAttDes{ObjectID: "module", PropertyID: common.MODULENAME_FIELD, PropertyName: "模块名", IsRequired: true IsOnly: false, PropertyType: common.FiledTypeSingleChar, Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: common.BKChildStr, PropertyName: "", IsRequired: false, IsOnly: false, PropertyType: "", Option: ""},
		//&metadata.ObjectAttDes{ObjectID: "module", PropertyID: "bk_onwer_id", PropertyName: "域", IsRequired: false, IsOnly: false, PropertyType: common.FiledTypeSingleChar, Option: ""},
		//&metadata.ObjectAttDes{ObjectID: "module", PropertyID: common.BKInstParentStr, PropertyName: "", IsRequired: false, IsOnly: false, PropertyType: common.FiledTypeInt, Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_module_type", PropertyName: "模块类型", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FiledTypeEnum, Option: moduleTypeOption},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "operator", PropertyName: "主要维护人", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FiledTypeUser, Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_bak_operator", PropertyName: "备份维护人", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FiledTypeUser, Option: ""},
	}
	return dataRows
}

func PlatRow() []*metadata.ObjectAttDes {
	objID := common.BKInnerObjIDPlat
	dataRows := []*metadata.ObjectAttDes{
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: common.BKCloudNameField, PropertyName: "云区域", IsRequired: true, IsOnly: true, IsPre: true, Editable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FiledTypeSingleChar, Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: common.BKOwnerIDField, PropertyName: "供应商", IsRequired: true, IsOnly: true, IsPre: true, PropertyGroup: groupBaseInfo, PropertyType: common.FiledTypeSingleChar, Option: ""},
	}
	return dataRows
}

func HostRow() []*metadata.ObjectAttDes {
	objID := common.BKInnerObjIDHost
	groupAgent := mCommon.HostAutoFields
	dataRows := []*metadata.ObjectAttDes{
		//&metadata.ObjectAttDes{ObjectID: objID, PropertyID: common.HOSTID_FIELD, PropertyName: "主机ID", IsSystem: true, IsRequired: true, IsOnly: false, PropertyGroup: groupBaseInfo, PropertyType: common.FiledTypeInt, Option: ""},
		//基本信息分组
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: common.BKHostInnerIPField, PropertyName: "内网IP", IsRequired: true, IsOnly: true, Editable: false, PropertyGroup: groupBaseInfo, PropertyType: common.FiledTypeSingleChar, Option: common.PatternMultipleIP},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: common.BKHostOuterIPField, PropertyName: "外网IP", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FiledTypeSingleChar, Option: common.PatternMultipleIP},
		// &metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_agent_status", PropertyName: "Agent状态", IsRequired: false, IsOnly: false, PropertyGroup: groupBaseInfo, PropertyType: common.FiledTypeEnum, Option: "[{\"name\":\"正常\", \"type\":\"text\"},{\"name\":\"异常\", \"type\":\"text\"},{\"name\":\"未安装\", \"type\":\"text\"}]"},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "operator", PropertyName: "主要维护人", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FiledTypeUser, Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_bak_operator", PropertyName: "备份维护人", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FiledTypeUser, Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_asset_id", PropertyName: "固资编号", IsRequired: false, IsOnly: false, Editable: false, PropertyGroup: groupBaseInfo, PropertyType: common.FiledTypeSingleChar, Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_sn", PropertyName: "设备SN", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FiledTypeSingleChar, Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_comment", PropertyName: "备注", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FiledTypeLongChar, Option: ""},

		//自动发现分组
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: common.BKCloudIDField, PropertyName: "云区域", IsRequired: false, IsOnly: true, Editable: false, PropertyGroup: groupBaseInfo, PropertyType: common.FiledTypeSingleAsst, Option: `{"value":"plat","label":"云区域"}`}, //common.FiledTypeInt, Option: "{}"},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_host_name", PropertyName: "主机名称", IsRequired: false, IsOnly: false, PropertyGroup: groupAgent, PropertyType: common.FiledTypeSingleChar, Option: ""},
		//&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_host_type", PropertyName: "主机类型", IsRequired: false, IsOnly: false, Editable: false, PropertyType: common.FiledTypeEnum, Option: "[{\"name\":\"虚拟机\",\"type\":\"text\"},{\"name\":\"实体机\",\"type\":\"text\"},{\"name\":\"容器虚拟机\",\"type\":\"text\"}]"},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_service_term", PropertyName: "质保年限", IsRequired: false, IsOnly: false, Editable: false, PropertyGroup: groupBaseInfo, PropertyType: common.FiledTypeInt, Option: "{\"min\":\"1\",\"max\":\"10\"}"},
		//&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_level", PropertyName: "重要级别", IsRequired: false, IsOnly: false, PropertyGroup: groupAgent, PropertyType: common.FiledTypeEnum, Option: "[{\"name\":\"重要\", \"type\":\"text\"},{\"name\":\"非常重要\", \"type\":\"text\"},{\"name\":\"一般\", \"type\":\"text\"},{\"name\":\"不重要\", \"type\":\"text\"}]"},
		//&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_status", PropertyName: "运行状态", IsRequired: false, IsOnly: false, PropertyGroup: groupAgent, PropertyType: common.FiledTypeEnum, Option: "[{\"name\":\"离线\", \"type\":\"text\"},{\"name\":\"在线\", \"type\":\"text\"}]"},
		//&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_current_status", PropertyName: "当前状态", IsRequired: false, IsOnly: false, PropertyGroup: groupAgent, PropertyType: common.FiledTypeEnum, Option: "[{\"name\":\"运营中\", \"type\":\"text\"},{\"name\":\"故障\", \"type\":\"text\"},{\"name\":\"备用\", \"type\":\"text\"},{\"name\":\"重装中\", \"type\":\"text\"}]"},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_sla", PropertyName: "SLA级别", IsRequired: false, IsOnly: false, PropertyGroup: groupBaseInfo, PropertyType: common.FiledTypeEnum, Option: "[{\"name\":\"L1\", \"type\":\"text\"},{\"name\":\"L2\", \"type\":\"text\"},{\"name\":\"L3\", \"type\":\"text\"}]"},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: common.BKOSTypeField, PropertyName: "操作系统类型", IsRequired: false, IsOnly: false, PropertyGroup: groupAgent, PropertyType: common.FiledTypeEnum, Option: "[{\"name\":\"Linux\", \"type\":\"text\"},{\"name\":\"Windows\", \"type\":\"text\"}]"},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: common.BKOSNameField, PropertyName: "操作系统名称", IsRequired: false, IsOnly: false, PropertyGroup: groupAgent, PropertyType: common.FiledTypeSingleChar, Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_os_version", PropertyName: "操作系统版本", IsRequired: false, IsOnly: false, PropertyGroup: groupAgent, PropertyType: common.FiledTypeSingleChar, Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_os_bit", PropertyName: "操作系统位数", IsRequired: false, IsOnly: false, PropertyGroup: groupAgent, PropertyType: common.FiledTypeSingleChar, Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_cpu", PropertyName: "CPU逻辑核心数", IsRequired: false, IsOnly: false, PropertyGroup: groupAgent, PropertyType: common.FiledTypeInt, Option: "{\"min\":\"1\",\"max\":\"1000000\"}"},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_cpu_mhz", PropertyName: "CPU频率", IsRequired: false, IsOnly: false, PropertyGroup: groupAgent, PropertyType: common.FiledTypeInt, Unit: "Hz", Option: "{\"min\":\"1\",\"max\":\"100000000\"}"},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_cpu_module", PropertyName: "CPU型号", IsRequired: false, IsOnly: false, PropertyGroup: groupAgent, PropertyType: common.FiledTypeSingleChar, Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_mem", PropertyName: "内存容量", IsRequired: false, IsOnly: false, PropertyGroup: groupAgent, PropertyType: common.FiledTypeInt, Unit: "MB", Option: "{\"min\":\"1\",\"max\":\"100000000\"}"},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_disk", PropertyName: "磁盘容量", IsRequired: false, IsOnly: false, PropertyGroup: groupAgent, PropertyType: common.FiledTypeInt, Option: "{\"min\":\"1\",\"max\":\"100000000\"}"},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_mac", PropertyName: "内网MAC地址", IsRequired: false, IsOnly: false, PropertyGroup: groupAgent, PropertyType: common.FiledTypeSingleChar, Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_outer_mac", PropertyName: "外网MAC", IsRequired: false, IsOnly: false, PropertyGroup: groupAgent, PropertyType: common.FiledTypeSingleChar, Option: ""},
		//agent 没有分组
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: common.CreateTimeField, PropertyName: "录入时间", IsRequired: false, IsOnly: false, Editable: false, PropertyGroup: mCommon.GroupNone, PropertyType: common.FiledTypeTime, Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "import_from", PropertyName: "录入方式", IsRequired: false, IsOnly: false, Editable: false, PropertyGroup: mCommon.GroupNone, PropertyType: common.FiledTypeEnum, Option: "[{\"name\":\"excel\",\"type\":\"text\"},{\"name\":\"agent\",\"type\":\"text\"},{\"name\":\"api\",\"type\":\"text\"}]"},
		// &metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_agent_version", PropertyName: "Agent版本", IsRequired: false, IsOnly: false, PropertyGroup: mCommon.Group_None, PropertyType: common.FiledTypeSingleChar, Option: ""},
	}

	return dataRows
}

func ProcRow() []*metadata.ObjectAttDes {
	objID := common.BKInnerObjIDProc
	groupPort := mCommon.ProcPort
	// groupGsekit := mCommon.Proc_gsekit_base_info
	// groupGsekitManage := mCommon.Proc_gsekit_manage_info
	dataRows := []*metadata.ObjectAttDes{
		//base info
		//&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_process_id", PropertyName: "进程ID", IsSystem: true, IsRequired: true, IsOnly: false, PropertyGroup: groupBaseInfo, PropertyType: common.FiledTypeInt, Option: "{}"},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: common.BKAppIDField, PropertyName: "业务ID", IsAPI: true, IsRequired: true, IsOnly: true, PropertyGroup: groupBaseInfo, PropertyType: common.FiledTypeInt, Option: "{}"},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: common.BKProcessNameField, PropertyName: "进程名称", IsRequired: true, IsOnly: true, IsPre: true, Editable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FiledTypeSingleChar, Option: ""},

		//监听端口分组
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bind_ip", PropertyName: "绑定IP", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: groupPort, PropertyType: common.FiledTypeEnum, Option: "[{\"name\":\"127.0.0.1\",\"type\":\"text\"},{\"name\":\"0.0.0.0\",\"type\":\"text\"},{\"name\":\"第一内网IP\",\"type\":\"text\"},{\"name\":\"第一外网IP\",\"type\":\"text\"}]"},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "port", PropertyName: "端口", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: groupPort, PropertyType: common.FiledTypeSingleChar, Option: `^((\d+-\d+)|(\d+))(,((\d+)|(\d+-\d+)))*$`, Placeholder: `单个端口：8080 </br>多个连续端口：8080-8089 </br>多个不连续端口：8080-8089,8199`},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "protocol", PropertyName: "协议", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: groupPort, PropertyType: common.FiledTypeEnum, Option: "[{\"name\":\"TCP\",\"type\":\"text\"},{\"name\":\"UDP\",\"type\":\"text\"}]"},

		//gsekit 基础信息
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_func_id", PropertyName: "功能ID", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: mCommon.GroupNone, PropertyType: common.FiledTypeSingleChar, Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_func_name", PropertyName: "功能名称", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: mCommon.GroupNone, PropertyType: common.FiledTypeSingleChar, Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "work_path", PropertyName: "工作路径", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: mCommon.GroupNone, PropertyType: common.FiledTypeLongChar, Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "user", PropertyName: "启动用户", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: mCommon.GroupNone, PropertyType: common.FiledTypeSingleChar, Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "proc_num", PropertyName: "启动数量", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: mCommon.GroupNone, PropertyType: common.FiledTypeInt, Option: "{\"min\":\"1\",\"max\":\"1000000\"}"},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "priority", PropertyName: "启动优先级", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: mCommon.GroupNone, PropertyType: common.FiledTypeInt, Option: "{\"min\":\"1\",\"max\":\"100\"}"},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "timeout", PropertyName: "操作超时时长", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: mCommon.GroupNone, PropertyType: common.FiledTypeInt, Option: "{\"min\":\"1\",\"max\":\"1000000\"}"},

		//gsekit 进程信息
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "start_cmd", PropertyName: "启动命令", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: mCommon.GroupNone, PropertyType: common.FiledTypeLongChar, Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "stop_cmd", PropertyName: "停止命令", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: mCommon.GroupNone, PropertyType: common.FiledTypeLongChar, Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "restart_cmd", PropertyName: "重启命令", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: mCommon.GroupNone, PropertyType: common.FiledTypeLongChar, Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "face_stop_cmd", PropertyName: "强制停止命令", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: mCommon.GroupNone, PropertyType: common.FiledTypeLongChar, Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "reload_cmd", PropertyName: "进程重载命令", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: mCommon.GroupNone, PropertyType: common.FiledTypeLongChar, Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "pid_file", PropertyName: "PID文件路径", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: mCommon.GroupNone, PropertyType: common.FiledTypeLongChar, Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "auto_start", PropertyName: "是否自动拉起", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: mCommon.GroupNone, PropertyType: common.FiledTypeBool, Option: ""},
		&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "auto_time_gap", PropertyName: "拉起间隔", IsRequired: false, IsOnly: false, Editable: true, PropertyGroup: mCommon.GroupNone, PropertyType: common.FiledTypeInt, Option: "{\"min\":\"1\",\"max\":\"1000000\"}"},
		//&metadata.ObjectAttDes{ObjectID: objID, PropertyID: "bk_comment", PropertyName: "备注", IsRequired: false, IsOnly: false, PropertyType: common.FiledTypeLongChar, Option: ""},
	}
	return dataRows
}
