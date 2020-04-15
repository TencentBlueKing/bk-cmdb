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
package y3_8_202004141131

import (
	"context"
	"fmt"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
)

func updateEnablePortAttribute(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	filter := map[string]interface{}{
		common.BKPropertyIDField: "bk_port_enable",
		common.BKObjIDField:      common.BKInnerObjIDProc,
	}
	doc := map[string]interface{}{common.BKPropertyIDField: common.BKProcPortEnable}
	err := db.Table(common.BKTableNameObjAttDes).Update(ctx, filter, doc)
	if err != nil {
		blog.Errorf("update process attribute %s failed. err: %s", common.BKProcPortEnable, err.Error())
		return fmt.Errorf("update process attribute %s failed. err: %s", common.BKProcPortEnable, err.Error())
	}
	return nil
}

func updateProcessAndProcTemplateEnablePortAttribute(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	if err := db.Table(common.BKTableNameBaseProcess).RenameColumn(ctx, "bk_port_enable", common.BKProcPortEnable); err != nil {
		blog.Errorf("update process bk_port_enable field to %s failed, err: %s, %s", common.BKProcPortEnable, err.Error())
		return fmt.Errorf("update process bk_port_enable field to %s failed, err: %s", common.BKProcPortEnable, err.Error())
	}

	if err := db.Table(common.BKTableNameProcessTemplate).RenameColumn(ctx, "property.bk_port_enable", "property."+common.BKProcPortEnable); err != nil {
		blog.Errorf("update process template bk_port_enable field to %s failed, err: %s", common.BKProcPortEnable, err.Error())
		return fmt.Errorf("update process template bk_port_enable field to %s failed, err: %s", common.BKProcPortEnable, err.Error())
	}

	return nil
}
