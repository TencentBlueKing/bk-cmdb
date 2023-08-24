/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 THL A29 Limited,
 * a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 * We undertake not to change the open source license (MIT license) applicable
 * to the current version of the project delivered to anyone in the future.
 */

package core

import (
	"configcenter/src/common/backbone"
	"configcenter/src/common/language"
	"configcenter/src/common/mapstr"
)

// Client used to process excel-related data
type Client struct {
	Engine *backbone.Engine
}

type hostSetInfo struct {
	setIDs     []int64
	hostSetMap map[int64][]int64
}

type topoInstData struct {
	parentIDs         []int64
	instIdParentIDMap map[int64]int64
	instIdNameMap     map[int64]string
}

// TopoBriefMsg topo brief message
type TopoBriefMsg struct {
	ObjID string
	Name  string
}

// ImportedParam import instance parameter
type ImportedParam struct {
	Language   language.CCLanguageIf
	ObjID      string
	Req        mapstr.MapStr
	Instances  map[int]map[string]interface{}
	HandleType HandleType
}

// SameIPRes same ip resource
type SameIPRes struct {
	V4Map map[string]struct{}
	V6Map map[string]struct{}
}

// SimpleHost simple host
type SimpleHost struct {
	Ip      string
	Ipv6    string
	AgentID string
}
