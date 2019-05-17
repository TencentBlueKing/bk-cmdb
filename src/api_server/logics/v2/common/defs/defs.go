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

package defs

import "configcenter/src/common"

// 2.0与3.0的对象映射
var ObjMap = map[string]map[string]string{
	"1": map[string]string{
		"ObjectID":      common.BKInnerObjIDApp,
		"IDName":        "ApplicationID",
		"IDDisplayName": "业务ID",
	},
	"2": map[string]string{
		"ObjectID":      common.BKInnerObjIDSet,
		"IDName":        "SetID",
		"IDDisplayName": "集群ID",
	},
	"3": map[string]string{
		"ObjectID":      common.BKInnerObjIDModule,
		"IDName":        "ModuleID",
		"IDDisplayName": "模块ID",
	},
	"4": map[string]string{
		"ObjectID":      common.BKInnerObjIDHost,
		"IDName":        "HostID",
		"IDDisplayName": "主机ID",
	},
}

//env 1：测试 2：体验 3：正式，默认为3
var SetEnvMap = map[string]string{
	"1": "测试",
	"2": "体验",
	"3": "正式",
}

//服务状态，包含0：关闭，1：开启，默认为1
var SetStatusMap = map[string]string{
	"0": "关闭",
	"1": "开放",
}

var RoleMap = map[string]string{
	"Maintainers": common.BKMaintainersField,
	"ProductPm":   common.BKProductPMField,
	"Tester":      common.BKTesterField,
	"Developer":   common.BKDeveloperField,
}
