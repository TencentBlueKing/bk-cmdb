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

package x19_05_16_01

import (
	"context"
	"fmt"

	"configcenter/src/common"
	"configcenter/src/common/condition"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
)

func upgradeServiceTemplate(ctx context.Context, db dal.RDB, conf *upgrader.Config) (err error) {
	categoryID, err := addDefaultCategory(ctx, db, conf)
	if err != nil {
		return fmt.Errorf("addDefaultCategory failed: %v", err)
	}

	allmodules := []metadata.ModuleInst{}
	if err = db.Table(common.BKTableNameBaseModule).Find(nil).All(ctx, &allmodules); err != nil {
		return err
	}

	// bizID:modulename:modules
	biz2Module := map[int64]map[string][]metadata.ModuleInst{}
	for _, module := range allmodules {
		_, ok := biz2Module[module.BizID]
		if !ok {
			biz2Module[module.BizID] = map[string][]metadata.ModuleInst{}
		}
		biz2Module[module.BizID][module.ModuleName] = append(biz2Module[module.BizID][module.ModuleName], module)
	}

	for bizID, modules := range biz2Module {
		process := metadata.Process{}
	}

	return nil
}
