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

package y3_14_202410100930

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/admin_server/upgrader/history"
	"configcenter/src/storage/dal"
)

func updatePlatformConfigAdmin(ctx context.Context, db dal.RDB, conf *history.Config) error {
	cond := map[string]interface{}{
		"_id": common.ConfigAdminID,
	}
	platCfg := make(map[string]string)
	err := db.Table(common.BKTableNameSystem).Find(cond).Fields(common.ConfigAdminValueField).One(ctx, &platCfg)
	if err != nil {
		blog.Errorf("get system config failed, db find err: %v", err)
		return err
	}

	if platCfg[common.ConfigAdminValueField] == "" {
		blog.Errorf("get system config failed, db config is empty")
		return err
	}

	dbCfg := new(platformSettingConfig)
	if err := json.Unmarshal([]byte(platCfg[common.ConfigAdminValueField]), dbCfg); err != nil {
		blog.Errorf("get dbConfig failed, unmarshal err: %v, config: %v", err, dbCfg)
		return err
	}
	if dbCfg.Backend.SnapshotBizID != 0 {
		return nil
	}

	if dbCfg.Backend.SnapshotBizName == "" {
		dbCfg.Backend.SnapshotBizName = common.BKAppName
	}
	bizInfo := new(metadata.BizBasicInfo)
	err = db.Table(common.BKTableNameBaseApp).Find(map[string]string{common.BKAppNameField: dbCfg.Backend.
		SnapshotBizName}).Fields(common.BKAppIDField).One(ctx, &bizInfo)
	if err != nil {
		blog.Errorf("get bizID failed, db find err: %v", err)
		return err
	}
	if bizInfo.BizID == 0 {
		blog.Errorf("can not get biz: %s", dbCfg.Backend.SnapshotBizName)
		return errors.New("can not get biz in db")
	}

	dbCfg.Backend.SnapshotBizID = bizInfo.BizID
	dbCfg.Backend.SnapshotBizName = ""

	bytes, err := json.Marshal(dbCfg)
	if err != nil {
		return err
	}

	updateCond := map[string]interface{}{
		"_id": common.ConfigAdminID,
	}
	data := map[string]interface{}{
		common.ConfigAdminValueField: string(bytes),
		common.LastTimeField:         time.Now(),
	}
	err = db.Table(common.BKTableNameSystem).Update(ctx, updateCond, data)
	if err != nil {
		blog.Errorf("update system config failed, err: %v", err)
		return err
	}
	return nil
}

// adminBackendCfg admin backend config
type adminBackendCfg struct {
	MaxBizTopoLevel int64  `json:"max_biz_topo_level"`
	SnapshotBizID   int64  `json:"snapshot_biz_id"`
	SnapshotBizName string `json:"snapshot_biz_name,omitempty"`
}

// platformSettingConfig platform setting config
type platformSettingConfig struct {
	Backend             adminBackendCfg             `json:"backend"`
	ValidationRules     metadata.ValidationRulesCfg `json:"validation_rules"`
	BuiltInSetName      metadata.ObjectString       `json:"set"`
	BuiltInModuleConfig metadata.GlobalModule       `json:"idle_pool"`
	IDGenerator         metadata.IDGeneratorConf    `json:"id_generator"`
}
