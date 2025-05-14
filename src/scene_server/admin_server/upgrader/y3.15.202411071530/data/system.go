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
	"fmt"

	idgen "configcenter/pkg/id-gen"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/storage/dal/mongo/local"
	"configcenter/src/storage/driver/mongodb"
)

func addSystemData(kit *rest.Kit, db local.DB) error {
	blog.Infof("start add init data for table: %s", common.BKTableNameSystem)

	if err := initPlatformSetting(kit, mongodb.Shard(kit.SysShardOpts())); err != nil {
		blog.Errorf("add id generator failed, error: %v", err)
		return err
	}

	blog.Infof("end add init data for table: %s", common.BKTableNameSystem)
	return nil
}

func initPlatformSetting(kit *rest.Kit, db local.DB) error {

	existConfig := make([]PlatformConfig, 0)
	err := db.Table(common.BKTableNameSystem).Find(mapstr.MapStr{common.BKFieldDBID: common.PlatformConfig}).Fields(
		metadata.IDGeneratorConfig).All(kit.Ctx, &existConfig)
	if err != nil {
		blog.Errorf("get config id generator failed, error: %v", err)
		return err
	}

	if len(existConfig) > 0 {
		if !cmpSame(&existConfig[0].IDGenerator, &InitIDGeneratorConfig) {
			blog.Errorf("config id generator is not same, exist: %v, insert: %v", existConfig[0], InitIDGeneratorConfig)
			return fmt.Errorf("config id generator is not same")
		}
		return nil
	}

	insertData := idGeneratorConf{
		BID:         common.PlatformConfig,
		IDGenerator: &InitIDGeneratorConfig,
	}
	err = db.Table(common.BKTableNameSystem).Insert(kit.Ctx, insertData)
	if err != nil {
		blog.Errorf("insert data for table %s failed, err: %v", common.BKTableNameSystem, err)
		return err
	}

	return nil
}

// InitIDGeneratorConfig id generator init config
var InitIDGeneratorConfig = IDGeneratorConf{
	Enabled:   false,
	Step:      1,
	InitID:    nil,
	CurrentID: nil,
}

type idGeneratorConf struct {
	BID         string           `json:"_id" bson:"_id"`
	IDGenerator *IDGeneratorConf `json:"id_generator" bson:"id_generator"`
}

func cmpSame(existData, insertData *IDGeneratorConf) bool {
	if existData.Enabled != insertData.Enabled || existData.Step != insertData.Step {
		return false
	}
	return true
}

// IDGeneratorConf is id generator config
type IDGeneratorConf struct {
	Enabled bool                       `json:"enabled" bson:"enabled"`
	Step    int                        `json:"step" bson:"step"`
	InitID  map[idgen.IDGenType]uint64 `json:"init_id,omitempty" bson:"init_id,omitempty"`
	// CurrentID is the current id of each resource, this is only used for ui display
	CurrentID map[idgen.IDGenType]uint64 `json:"current_id,omitempty" bson:"current_id,omitempty"`
}

// PlatformConfig platform config
type PlatformConfig struct {
	IDGenerator IDGeneratorConf `bson:"id_generator" json:"id_generator"`
}
