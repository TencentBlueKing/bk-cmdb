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

package x20_02_17_01

import (
	"context"
	"fmt"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
)

type ProcessTemplate struct {
	Metadata metadata.Metadata `field:"metadata" json:"metadata" bson:"metadata"`

	ID          int64  `field:"id" json:"id,omitempty" bson:"id"`
	ProcessName string `field:"bk_process_name" json:"bk_process_name" bson:"bk_process_name"`
	// the service template's, which this process template belongs to.
	ServiceTemplateID int64 `field:"service_template_id" json:"service_template_id" bson:"service_template_id"`

	// stores a process instance's data includes all the process's
	// properties's value.
	Property *metadata.ProcessProperty `field:"property" json:"property,omitempty" bson:"property"`

	Creator         string    `field:"creator" json:"creator,omitempty" bson:"creator"`
	Modifier        string    `field:"modifier" json:"modifier,omitempty" bson:"modifier"`
	CreateTime      time.Time `field:"create_time" json:"create_time,omitempty" bson:"create_time"`
	LastTime        time.Time `field:"last_time" json:"last_time,omitempty" bson:"last_time"`
	SupplierAccount string    `field:"bk_supplier_account" json:"bk_supplier_account,omitempty" bson:"bk_supplier_account"`
}

func fixProcTemplateProcName(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {

	conds := map[string]interface{}{
		"$or": []interface{}{
			map[string]string{"property.bk_func_name.value": ""},
			map[string]interface{}{"property.bk_func_name.value": map[string]interface{}{"$exists": false}},
		},
	}

	rows := make([]ProcessTemplate, 0)
	if err := db.Table(common.BKTableNameProcessTemplate).Find(conds).All(ctx, &rows); err != nil {
		blog.ErrorJSON("find process template bk_process_name empty error. err:%s", err.Error())
		return err
	}

	for _, row := range rows {
		updateCond := map[string]interface{}{
			"id": row.ID,
		}
		if row.Property == nil {
			blog.ErrorJSON("fix process template id:%v, process template property empty, raw info:%s", row.ID, row)
			return fmt.Errorf("fix process template id:%v, process template property empty", row.ID)
		}
		if row.Property.ProcessName.Value == nil {
			blog.ErrorJSON("fix process template id:%v, bk_process_name empty, raw info:%s", row.ID, row)
			return fmt.Errorf("fix process template id:%v, bk_process_name empty", row.ID)
		}
		doc := map[string]interface{}{
			"property.bk_func_name.value":            *row.Property.ProcessName.Value,
			"property.bk_func_name.as_default_value": true,
		}
		if row.ProcessName == "" {
			doc["bk_process_name"] = *row.Property.ProcessName.Value
		}
		if err := db.Table(common.BKTableNameProcessTemplate).Update(ctx, updateCond, doc); err != nil {
			blog.ErrorJSON("fix process template id:%v, update db error. condition:%s, doc:%s, err:%s", row.ID, updateCond, doc, err.Error())
			return fmt.Errorf("fix process template id:%v, update db error. err:%s", row.ID, err.Error())
		}
	}

	return nil
}
