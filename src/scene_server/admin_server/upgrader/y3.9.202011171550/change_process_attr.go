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

package y3_9_202011171550

import (
	"context"
	"fmt"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	mCommon "configcenter/src/scene_server/admin_server/common"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/dal/mongo/local"
)

func changeProcessAttrs(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	//【启动优先级】增加提示信息： 批量启动进程依据优先级从小到大排序操作，停止进程按反序操作
	filter := map[string]interface{}{
		common.BKObjIDField:      common.BKInnerObjIDProc,
		common.BKPropertyIDField: "priority",
	}

	doc := map[string]interface{}{
		"placeholder": "批量启动进程依据优先级从小到大排序操作，停止进程按反序操作",
	}

	if err := db.Table(common.BKTableNameObjAttDes).Update(ctx, filter, doc); err != nil {
		blog.Errorf("update failed, filter:%#v, doc:%#v, err:%v", filter, doc, err)
		return err
	}

	// 【操作超时时长】增加单位 “s”, 增加提示信息：执行命令的超时时长
	filter = map[string]interface{}{
		common.BKObjIDField:      common.BKInnerObjIDProc,
		common.BKPropertyIDField: "timeout",
	}

	doc = map[string]interface{}{
		"unit":        "s",
		"placeholder": "执行命令的超时时长",
	}

	if err := db.Table(common.BKTableNameObjAttDes).Update(ctx, filter, doc); err != nil {
		blog.Errorf("update failed, filter:%#v, doc:%#v, err:%v", filter, doc, err)
		return err
	}

	//【工作路径】提示信息改为：执行进程启停等操作的工作路径
	filter = map[string]interface{}{
		common.BKObjIDField:      common.BKInnerObjIDProc,
		common.BKPropertyIDField: "work_path",
	}

	doc = map[string]interface{}{
		"placeholder": "执行进程启停等操作的工作路径",
	}

	if err := db.Table(common.BKTableNameObjAttDes).Update(ctx, filter, doc); err != nil {
		blog.Errorf("update failed, filter:%#v, doc:%#v, err:%v", filter, doc, err)
		return err
	}

	// 删除"拉起间隔"，"功能ID"属性
	filter = map[string]interface{}{
		common.BKObjIDField: common.BKInnerObjIDProc,
		common.BKPropertyIDField: map[string]interface{}{
			"$in": []string{"auto_time_gap", "bk_func_id"},
		},
	}

	if err := db.Table(common.BKTableNameObjAttDes).Delete(ctx, filter); err != nil {
		blog.Errorf("update failed, filter:%#v, err:%v", filter, err)
		return err
	}

	return nil
}

func addProcessAttr(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	//  增加【启动等待时长】字段
	now := time.Now()
	row := attribute{
		BizID:             0,
		OwnerID:           conf.OwnerID,
		ObjectID:          common.BKInnerObjIDProc,
		PropertyID:        "bk_start_check_secs",
		PropertyName:      "启动等待时长",
		PropertyGroup:     mCommon.ProcMgrGroupID,
		PropertyGroupName: mCommon.ProcMgrGroupName,
		PropertyIndex:     0,
		Unit:              "s",
		Placeholder:       "在执行启动命令后，等待多久检查PID存活的状态",
		IsEditable:        true,
		IsPre:             true,
		IsRequired:        false,
		IsReadOnly:        false,
		IsOnly:            false,
		IsSystem:          false,
		IsAPI:             false,
		PropertyType:      common.FieldTypeInt,
		Option:            metadata.IntOption{Min: "1", Max: "600"},
		Description:       common.CCSystemOperatorUserName,
		Creator:           common.CCSystemOperatorUserName,
		LastEditor:        common.CCSystemOperatorUserName,
		CreateTime:        &now,
		LastTime:          &now,
	}
	uniqueFields := []string{common.BKObjIDField, common.BKPropertyIDField, common.BKOwnerIDField}
	if _, _, err := upgrader.Upsert(ctx, db, common.BKTableNameObjAttDes, row, "id", uniqueFields, []string{}); err != nil {
		blog.ErrorJSON("addCloudHostAttr failed, Upsert err: %s, attribute: %#v, ", err, row)
		return err
	}
	return nil
}

// 删除进程实例中的"拉起间隔"，"功能ID"字段
func deleteProcessInstsFields(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	filter := map[string]interface{}{
		"$or": []map[string]interface{}{
			{
				"auto_time_gap": map[string]interface{}{
					"$exists": true,
				},
			},
			{
				"bk_func_id": map[string]interface{}{
					"$exists": true,
				},
			},
		},
	}
	processIDs, err := db.Table(common.BKTableNameBaseProcess).Distinct(ctx, common.BKProcessIDField, filter)
	if err != nil {
		blog.ErrorJSON("deleteProcessInstsFields failed, Distinct err: %s, filter: %#v, ", err, filter)
		return err
	}

	doc := map[string]interface{}{
		"$unset": map[string]string{
			"auto_time_gap": "",
			"bk_func_id":    "",
		},
	}

	mongo, ok := db.(*local.Mongo)
	if !ok {
		return fmt.Errorf("db is not *local.Mongo type")
	}

	pageSize := 1000
	length := len(processIDs)
	for startIdx := 0; startIdx < length; startIdx += pageSize {
		endIdx := startIdx + pageSize
		if endIdx > length {
			endIdx = length
		}

		filter := map[string]interface{}{
			common.BKProcessIDField: map[string]interface{}{
				"$in": processIDs[startIdx:endIdx],
			},
		}
		if _, err := mongo.GetDBClient().Database(mongo.GetDBName()).Collection(common.BKTableNameBaseProcess).UpdateMany(ctx, filter, doc); err != nil {
			blog.Errorf("update process fields failed, filter:%#v, doc:%#v, err:%v", filter, doc, err)
			return err
		}
	}

	return nil
}

// 删除进程模版实例中的"拉起间隔"，"功能ID"字段
func deleteProcessTemplateInstsFields(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	filter := map[string]interface{}{
		"$or": []map[string]interface{}{
			{
				"property.auto_time_gap": map[string]interface{}{
					"$exists": true,
				},
			},
			{
				"property.bk_func_id": map[string]interface{}{
					"$exists": true,
				},
			},
		},
	}
	templateIDs, err := db.Table(common.BKTableNameProcessTemplate).Distinct(ctx, common.BKFieldID, filter)
	if err != nil {
		blog.ErrorJSON("deleteProcessTemplateInstsFields failed, Distinct err: %s, filter: %#v, ", err, filter)
		return err
	}

	doc := map[string]interface{}{
		"$unset": map[string]string{
			"property.auto_time_gap": "",
			"property.bk_func_id":    "",
		},
	}

	mongo, ok := db.(*local.Mongo)
	if !ok {
		return fmt.Errorf("db is not *local.Mongo type")
	}

	pageSize := 1000
	length := len(templateIDs)
	for startIdx := 0; startIdx < length; startIdx += pageSize {
		endIdx := startIdx + pageSize
		if endIdx > length {
			endIdx = length
		}

		filter := map[string]interface{}{
			common.BKFieldID: map[string]interface{}{
				"$in": templateIDs[startIdx:endIdx],
			},
		}
		if _, err := mongo.GetDBClient().Database(mongo.GetDBName()).Collection(common.BKTableNameProcessTemplate).UpdateMany(ctx, filter, doc); err != nil {
			blog.Errorf("update process template fields failed, filter:%#v, doc:%#v, err:%v", filter, doc, err)
			return err
		}
	}

	return nil
}

type attribute struct {
	BizID             int64       `field:"bk_biz_id" json:"bk_biz_id" bson:"bk_biz_id"`
	ID                int64       `field:"id" json:"id" bson:"id"`
	OwnerID           string      `field:"bk_supplier_account" json:"bk_supplier_account" bson:"bk_supplier_account"`
	ObjectID          string      `field:"bk_obj_id" json:"bk_obj_id" bson:"bk_obj_id"`
	PropertyID        string      `field:"bk_property_id" json:"bk_property_id" bson:"bk_property_id"`
	PropertyName      string      `field:"bk_property_name" json:"bk_property_name" bson:"bk_property_name"`
	PropertyGroup     string      `field:"bk_property_group" json:"bk_property_group" bson:"bk_property_group"`
	PropertyGroupName string      `field:"bk_property_group_name,ignoretomap" json:"bk_property_group_name" bson:"-"`
	PropertyIndex     int64       `field:"bk_property_index" json:"bk_property_index" bson:"bk_property_index"`
	Unit              string      `field:"unit" json:"unit" bson:"unit"`
	Placeholder       string      `field:"placeholder" json:"placeholder" bson:"placeholder"`
	IsEditable        bool        `field:"editable" json:"editable" bson:"editable"`
	IsPre             bool        `field:"ispre" json:"ispre" bson:"ispre"`
	IsRequired        bool        `field:"isrequired" json:"isrequired" bson:"isrequired"`
	IsReadOnly        bool        `field:"isreadonly" json:"isreadonly" bson:"isreadonly"`
	IsOnly            bool        `field:"isonly" json:"isonly" bson:"isonly"`
	IsSystem          bool        `field:"bk_issystem" json:"bk_issystem" bson:"bk_issystem"`
	IsAPI             bool        `field:"bk_isapi" json:"bk_isapi" bson:"bk_isapi"`
	PropertyType      string      `field:"bk_property_type" json:"bk_property_type" bson:"bk_property_type"`
	Option            interface{} `field:"option" json:"option" bson:"option"`
	Description       string      `field:"description" json:"description" bson:"description"`
	Creator           string      `field:"creator" json:"creator" bson:"creator"`
	LastEditor        string      `json:"bk_last_editor" bson:"bk_last_editor"`
	CreateTime        *time.Time  `json:"create_time" bson:"create_time"`
	LastTime          *time.Time  `json:"last_time" bson:"last_time"`
}
