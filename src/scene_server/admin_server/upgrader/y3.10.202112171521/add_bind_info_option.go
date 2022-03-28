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

package y3_10_202112171521

import (
	"context"
	"fmt"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
)

// addProcessBindInfoOption update process bind info attribute
func addProcessBindInfoOption(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {

	procAttr := make([]metadata.Attribute, 0)
	bindInfoCond := map[string]interface{}{
		common.BKObjIDField:      common.BKInnerObjIDProc,
		common.BKPropertyIDField: common.BKProcBindInfo,
	}
	if err := db.Table(common.BKTableNameObjAttDes).Find(bindInfoCond).All(ctx, &procAttr); err != nil {
		blog.Errorf("search bind_info in table[%s] failed, cond: %v, err: %v", common.BKTableNameObjAttDes,
			bindInfoCond, err)
		return err
	}

	if len(procAttr) != 1 {
		blog.Errorf("bind_info in %s has %d, not only one", common.BKTableNameObjAttDes, len(procAttr))
		return fmt.Errorf("bind_info in %s not only one", common.BKTableNameObjAttDes)
	}

	options, err := metadata.ParseSubAttribute(ctx, procAttr[0].Option)
	if err != nil {
		blog.Errorf("parse sub-attribute failed, err: %v", err)
		return err
	}

	for _, subAttr := range options {
		if subAttr.PropertyID == "template_row_id" {
			blog.Infof("bind_info option already has field template_row_id")
			return nil
		}
	}

	bindInfoOption := metadata.SubAttribute{
		PropertyID:    "template_row_id",
		PropertyName:  "RowID",
		Placeholder:   "process template row id",
		IsEditable:    false,
		PropertyType:  common.FieldTypeInt,
		PropertyGroup: common.BKProcBindInfo,
		IsAPI:         true,
	}

	options = append(options, bindInfoOption)
	doc := map[string]interface{}{
		"option": options,
	}
	return db.Table(common.BKTableNameObjAttDes).Update(ctx, bindInfoCond, doc)
}
