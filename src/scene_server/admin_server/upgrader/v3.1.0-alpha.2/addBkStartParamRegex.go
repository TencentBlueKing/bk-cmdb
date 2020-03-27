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

package v3v0v1alpha2

import (
	"context"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
)

func addBkStartParamRegex(ctx context.Context, db dal.RDB, conf *upgrader.Config) (err error) {
	tablename := common.BKTableNameObjAttDes
	now := time.Now()

	type Attribute struct {
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
		CreateTime        *time.Time  `json:"create_time" bson:"create_time"`
		LastTime          *time.Time  `json:"last_time" bson:"last_time"`
	}

	row := &Attribute{
		ObjectID:      common.BKInnerObjIDProc,
		PropertyID:    "bk_start_param_regex",
		PropertyName:  "启动参数匹配规则",
		IsRequired:    false,
		IsOnly:        false,
		IsEditable:    true,
		PropertyGroup: "default",
		PropertyType:  common.FieldTypeLongChar,
		Option:        "",
		OwnerID:       conf.OwnerID,
		IsPre:         true,
		IsReadOnly:    false,
		CreateTime:    &now,
		Creator:       common.CCSystemOperatorUserName,
		LastTime:      &now,
		Description:   "通过进程启动参数唯一识别进程，比如kafka和zookeeper的二进制名称为java，通过启动参数包含kafka或zookeeper来区分",
	}
	_, _, err = upgrader.Upsert(ctx, db, tablename, row, "id", []string{common.BKObjIDField, common.BKPropertyIDField, common.BKOwnerIDField}, []string{})
	if nil != err {
		blog.Errorf("[upgrade v3.1.0-alpha.2] addBkStartParamRegex  %s", err)
		return err
	}

	return nil
}

func updateLanguageField(ctx context.Context, db dal.RDB, conf *upgrader.Config) (err error) {
	condition := map[string]interface{}{
		common.BKObjIDField:      common.BKInnerObjIDApp,
		common.BKPropertyIDField: "language",
	}
	data := map[string]interface{}{
		"isreadonly": false,
		"editable":   true,
	}
	err = db.Table(common.BKTableNameObjAttDes).Update(ctx, condition, data)
	if nil != err {
		blog.Errorf("[upgrade v3.1.0-alpha.2] updateLanguageField error  %s", err.Error())
		return err
	}
	return nil
}
