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

package data

import (
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/querybuilder"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/admin_server/service/utils"
	"configcenter/src/scene_server/admin_server/upgrader/tools"
	"configcenter/src/storage/dal/mongo/local"
)

var (
	bizSetData = bizSetInst{
		BizSetName:       "BlueKing",
		BizSetMaintainer: "admin",
		Default:          common.DefaultResBusinessSetFlag,
		Scope: bizSetScope{
			MatchAll: true,
		},
		Time: tools.NewTime(),
	}
	bizSetAudit = &utils.AuditResType{
		AuditType:    metadata.BizSetType,
		ResourceType: metadata.BizSetRes,
	}
)

func addBizSetData(kit *rest.Kit, db local.DB) error {

	needField := &utils.InsertOptions{
		UniqueFields:   []string{common.BKBizSetNameField},
		IgnoreKeys:     []string{common.BKBizSetIDField},
		IDField:        []string{common.BKBizSetIDField},
		AuditTypeField: bizSetAudit,
		AuditDataField: &utils.AuditDataField{
			ResIDField:   "bk_biz_set_id",
			ResNameField: "bk_biz_set_name",
		},
	}
	data, err := util.ConvStructToMap(bizSetData)
	if err != nil {
		blog.Errorf("convert struct to map failed, err: %v", err)
		return err
	}

	_, err = utils.InsertData(kit, db, common.BKTableNameBaseBizSet, []mapstr.MapStr{data}, needField)
	if err != nil {
		blog.Errorf("insert default biz data for table %s failed, err: %v", common.BKTableNameBaseBizSet, err)
		return err
	}

	idOptions := &tools.IDOptions{IDField: "id", RemoveKeys: []string{"bk_biz_set_id"}}
	err = tools.InsertTemplateData(kit, db, []mapstr.MapStr{data}, "biz_set", []string{"data.bk_biz_set_name"},
		idOptions)
	if err != nil {
		blog.Errorf("insert template data failed, err: %v", err)
		return err
	}

	return nil
}

type bizSetInst struct {
	BizSetID         int64       `bson:"bk_biz_set_id"`
	BizSetName       string      `bson:"bk_biz_set_name"`
	Description      string      `bson:"bk_biz_set_desc"`
	BizSetMaintainer string      `bson:"bk_biz_maintainer"`
	Scope            bizSetScope `bson:"bk_scope"`
	Default          int64       `bson:"default"`
	*tools.Time      `bson:",inline"`
}

type bizSetScope struct {
	MatchAll bool                      `bson:"match_all"`
	Filter   *querybuilder.QueryFilter `bson:"filter,omitempty"`
}
