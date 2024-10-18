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

// Package logics defines cmdb resource sync logics
package logics

import (
	"configcenter/pkg/synchronize/types"
	"configcenter/src/source_controller/transfer-service/app/options"
	"configcenter/src/source_controller/transfer-service/sync/metadata"
	"configcenter/src/source_controller/transfer-service/sync/util"
)

// Logics is the resource sync logics interface
type Logics interface {
	ResType() types.ResType
	ParseDataArr(env, subRes string, data any, rid string) (any, error)
	ListData(kit *util.Kit, opt *types.ListDataOpt) (*types.ListDataRes, error)
	CompareData(kit *util.Kit, subRes string, srcInfo *types.FullSyncTransData, destInfo *types.ListDataRes) (
		*types.CompDataRes, error)
	ClassifyUpsertData(kit *util.Kit, subRes string, upsertData any) (any, any, error)
	InsertData(kit *util.Kit, subRes string, data any) error
	UpdateData(kit *util.Kit, subRes string, data any) error
	DeleteData(kit *util.Kit, subRes string, data any) error
}

// New creates a new resource type to resource sync logics map
func New(conf *LogicsConfig) map[types.ResType]Logics {
	lgcMap := map[types.ResType]Logics{
		types.Biz:             newDataWithIDLogics(conf.genResLgcConf(types.Biz), bizLgc),
		types.Set:             newDataWithIDLogics(conf.genResLgcConf(types.Set), setLgc),
		types.Module:          newDataWithIDLogics(conf.genResLgcConf(types.Module), moduleLgc),
		types.Host:            newDataWithIDLogics(conf.genResLgcConf(types.Host), hostLgc),
		types.HostRelation:    newRelationLogics(conf.genResLgcConf(types.HostRelation), hostRelLgc),
		types.ObjectInstance:  newObjInstLogics(conf.genResLgcConf(types.ObjectInstance)),
		types.InstAsst:        newDataWithIDLogics(conf.genResLgcConf(types.InstAsst), instAsstLgc),
		types.ServiceInstance: newDataWithIDLogics(conf.genResLgcConf(types.ServiceInstance), serviceInstLgc),
		types.Process:         newDataWithIDLogics(conf.genResLgcConf(types.Process), procLgc),
		types.ProcessRelation: newRelationLogics(conf.genResLgcConf(types.ProcessRelation), procRelLgc),
		types.QuotedInstance:  newDataWithIDLogics(conf.genResLgcConf(types.QuotedInstance), quotedInstLgc),
	}

	return lgcMap
}

// LogicsConfig is the cmdb resource sync logics config
type LogicsConfig struct {
	Metadata      *metadata.Metadata
	IDRuleMap     map[types.ResType]map[string][]options.IDRuleInfo
	SrcInnerIDMap map[string]*options.InnerDataIDConf
}

func (c *LogicsConfig) genResLgcConf(resType types.ResType) *resLogicsConfig {
	return &resLogicsConfig{
		resType:       resType,
		metadata:      c.Metadata,
		idRuleMap:     c.IDRuleMap,
		srcInnerIDMap: c.SrcInnerIDMap,
	}
}

// resLogicsConfig is the cmdb resource sync logics config
type resLogicsConfig struct {
	resType       types.ResType
	metadata      *metadata.Metadata
	idRuleMap     map[types.ResType]map[string][]options.IDRuleInfo
	srcInnerIDMap map[string]*options.InnerDataIDConf
}

// ResType get resource type
func (l *resLogicsConfig) ResType() types.ResType {
	return l.resType
}
