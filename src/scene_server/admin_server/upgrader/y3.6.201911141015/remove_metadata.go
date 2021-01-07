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

package y3_6_201911141015

import (
	"context"
	"fmt"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
)

// removeMetadataOnMainlineInstance 当前仅移除Set和Module的metadata字段
func removeMetadataOnMainlineInstance(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	if err := fixBizIDFieldWithMetadataOnSet(ctx, db, conf); err != nil {
		blog.Errorf("removeMetadataOnMainlineInstance failed, fixBizIDFieldWithMetadataOnSet failed, err: %s", err.Error())
		return fmt.Errorf("fixBizIDFieldWithMetadataOnSet failed, err: %s", err.Error())
	}

	if err := fixBizIDFieldWithMetadataOnModule(ctx, db, conf); err != nil {
		blog.Errorf("removeMetadataOnMainlineInstance failed, fixBizIDFieldWithMetadataOnModule failed, err: %s", err.Error())
		return fmt.Errorf("fixBizIDFieldWithMetadataOnModule failed, err: %s", err.Error())
	}

	// custom mainline instance don't have bk_biz_id field, it's dangerous to remove metadata
	/*
		customMainlineModels, err := getCustomMainlineModels(ctx, db, conf)
		if err != nil {
			blog.Errorf("removeMetadataOnMainlineInstance failed, getCustomMainlineModels failed, err: %s", err.Error())
			return fmt.Errorf("getCustomMainlineModels failed, err: %s", err.Error())
		}

		for _, mainlineModel := range customMainlineModels {
			if err := fixBizIDFieldWithMetadataOnCustomMainline(ctx, db, conf, mainlineModel); err != nil {
				blog.Errorf("removeMetadataOnMainlineInstance failed, fixBizIDFieldWithMetadataOnModule failed, err: %s", err.Error())
				return fmt.Errorf("fixBizIDFieldWithMetadataOnModule failed, err: %s", err.Error())
			}
		}
	*/

	// remove metadata field on table set
	if err := db.Table(common.BKTableNameBaseSet).DropColumn(ctx, common.MetadataField); err != nil {
		blog.Errorf("drop metadata field on set table failed, err: %+v", err)
		return fmt.Errorf("drop metadata field on set table failed, err: %v", err)
	}

	// remove metadata field on table module
	if err := db.Table(common.BKTableNameBaseModule).DropColumn(ctx, common.MetadataField); err != nil {
		blog.Errorf("drop metadata field on module table failed, err: %+v", err)
		return fmt.Errorf("drop metadata field on module table failed, err: %v", err)
	}

	return nil
}

func getCustomMainlineModels(ctx context.Context, db dal.RDB, conf *upgrader.Config) ([]string, error) {
	filter := map[string]interface{}{
		common.AssociationKindIDField: common.AssociationKindMainline,
	}

	type Association struct {
		// describe which object this association is defined for.
		ObjectID string `field:"bk_obj_id" json:"bk_obj_id" bson:"bk_obj_id"`
		// describe where the Object associate with.
		AsstObjID string `field:"bk_asst_obj_id" json:"bk_asst_obj_id" bson:"bk_asst_obj_id"`
		// the association kind used by this association.
		AsstKindID string `field:"bk_asst_id" json:"bk_asst_id" bson:"bk_asst_id"`
	}
	associations := make([]Association, 0)
	err := db.Table(common.BKTableNameObjAsst).Find(filter).All(ctx, associations)
	if err != nil {
		blog.Errorf("getCustomMainlineModels failed, db select failed, err: %s", err)
		return nil, fmt.Errorf("db select failed, %s", err.Error())
	}
	mainlineModels := make([]string, 0)
	for _, item := range associations {
		mainlineModels = append(mainlineModels, item.ObjectID)
		mainlineModels = append(mainlineModels, item.AsstObjID)
	}
	mainlineModels = util.StrArrayUnique(mainlineModels)

	// filter custom mainline
	builtInMainline := []string{
		common.BKInnerObjIDModule,
		common.BKInnerObjIDSet,
		common.BKInnerObjIDHost,
		common.BKInnerObjIDApp,
	}
	customMainlineModels := make([]string, 0)
	for _, model := range mainlineModels {
		if util.InStrArr(builtInMainline, model) {
			continue
		}
		customMainlineModels = append(customMainlineModels, model)
	}
	return customMainlineModels, nil
}

func fixBizIDFieldWithMetadataOnSet(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	type Set struct {
		BizID    int64             `bson:"bk_biz_id" json:"bk_biz_id" mapstructure:"bk_biz_id"`
		SetID    int64             `bson:"bk_set_id" json:"bk_set_id" mapstructure:"bk_set_id"`
		SetName  string            `bson:"bk_set_name" json:"bk_set_name" mapstructure:"bk_set_name"`
		Metadata metadata.Metadata `bson:"metadata" json:"metadata" mapstructure:"metadata"`
	}
	setFilter := map[string]interface{}{
		"metadata.label.bk_biz_id": map[string]interface{}{
			common.BKDBExists: true,
		},
		common.BKDBOR: []map[string]interface{}{
			{
				common.BKAppIDField: 0,
			},
			{
				common.BKAppIDField: map[string]interface{}{
					common.BKDBExists: false,
				},
			},
		},
	}
	sets := make([]Set, 0)
	if err := db.Table(common.BKTableNameBaseSet).Find(setFilter).All(ctx, &sets); err != nil {
		return err
	}

	for _, set := range sets {
		bizID, err := set.Metadata.ParseBizID()
		if err != nil {
			blog.Warnf("parse bizID from metadata failed, metadata: %+v, err: %s", set.Metadata, err.Error())
			continue
		}
		filter := map[string]interface{}{
			common.BKSetIDField: set.SetID,
		}
		doc := map[string]interface{}{
			common.BKAppIDField: bizID,
		}
		if err := db.Table(common.BKTableNameBaseSet).Update(ctx, filter, doc); err != nil {
			blog.Errorf("update set failed, filter: %s, doc: %s, err: %s", filter, doc, err.Error())
			return fmt.Errorf("update set failed, setID: %d, bizID: %d, err: %s", set.SetID, bizID, err.Error())
		}
	}

	return nil
}

func fixBizIDFieldWithMetadataOnModule(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	type Module struct {
		BizID      int64             `bson:"bk_biz_id" json:"bk_biz_id" field:"bk_biz_id" mapstructure:"bk_biz_id"`
		ModuleID   int64             `bson:"bk_module_id" json:"bk_module_id" field:"bk_module_id" mapstructure:"bk_module_id"`
		ModuleName string            `bson:"bk_module_name" json:"bk_module_name" field:"bk_module_name" mapstructure:"bk_module_name"`
		Metadata   metadata.Metadata `bson:"metadata" json:"metadata" mapstructure:"metadata"`
	}
	moduleFilter := map[string]interface{}{
		"metadata.label.bk_biz_id": map[string]interface{}{
			common.BKDBExists: true,
		},
		common.BKDBOR: []map[string]interface{}{
			{
				common.BKAppIDField: 0,
			},
			{
				common.BKAppIDField: map[string]interface{}{
					common.BKDBExists: false,
				},
			},
		},
	}
	modules := make([]Module, 0)
	if err := db.Table(common.BKTableNameBaseModule).Find(moduleFilter).All(ctx, &modules); err != nil {
		return err
	}

	for _, module := range modules {
		bizID, err := module.Metadata.ParseBizID()
		if err != nil {
			blog.Warnf("parse bizID from metadata failed, metadata: %+v, err: %s", module.Metadata, err.Error())
			continue
		}
		filter := map[string]interface{}{
			common.BKModuleIDField: module.ModuleID,
		}
		doc := map[string]interface{}{
			common.BKAppIDField: bizID,
		}
		if err := db.Table(common.BKTableNameBaseModule).Update(ctx, filter, doc); err != nil {
			blog.Errorf("update module failed, filter: %s, doc: %s, err: %s", filter, doc, err.Error())
			return fmt.Errorf("update module failed, moduleID: %d, bizID: %d, err: %s", module.ModuleID, bizID, err.Error())
		}
	}

	return nil
}

func fixBizIDFieldWithMetadataOnCustomMainline(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	type MainlineInstance struct {
		BizID    int64             `bson:"bk_biz_id" json:"bk_biz_id" field:"bk_biz_id" mapstructure:"bk_biz_id"`
		InstID   int64             `bson:"bk_inst_id" json:"bk_inst_id" field:"bk_inst_id" mapstructure:"bk_inst_id"`
		InstName string            `bson:"bk_inst_name" json:"bk_inst_name" field:"bk_inst_name" mapstructure:"bk_inst_name"`
		Metadata metadata.Metadata `bson:"metadata" json:"metadata" mapstructure:"metadata"`
	}
	instanceFilter := map[string]interface{}{
		"metadata.label.bk_biz_id": map[string]interface{}{
			common.BKDBExists: true,
		},
		common.BKDBOR: []map[string]interface{}{
			{
				common.BKAppIDField: 0,
			},
			{
				common.BKAppIDField: map[string]interface{}{
					common.BKDBExists: false,
				},
			},
		},
	}
	instances := make([]MainlineInstance, 0)
	if err := db.Table(common.BKTableNameBaseInst).Find(instanceFilter).All(ctx, &instances); err != nil {
		return err
	}

	for _, instance := range instances {
		bizID, err := instance.Metadata.ParseBizID()
		if err != nil {
			blog.Warnf("parse bizID from metadata failed, metadata: %+v, err: %s", instance.Metadata, err.Error())
			continue
		}
		filter := map[string]interface{}{
			common.BKInstIDField: instance.InstID,
		}
		doc := map[string]interface{}{
			common.BKAppIDField: bizID,
		}
		if err := db.Table(common.BKTableNameBaseInst).Update(ctx, filter, doc); err != nil {
			blog.Errorf("update mainline instance failed, filter: %s, doc: %s, err: %s", filter, doc, err.Error())
			return fmt.Errorf("update mainline instance failed, instanceID: %d, bizID: %d, err: %s", instance.InstID, bizID, err.Error())
		}
	}

	return nil
}
