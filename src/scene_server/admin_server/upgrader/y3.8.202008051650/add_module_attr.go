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

package y3_8_202008051650

import (
	"context"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	mCommon "configcenter/src/scene_server/admin_server/common"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
)

func addModuleAttr(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	objID := common.BKInnerObjIDModule
	now := time.Now()
	row := attribute{
		ObjectID: objID,
		PropertyID: "bk_parent_id",
		PropertyName: "模块实例的父节点ID",
		IsRequired: true,
		IsOnly: false,
		IsEditable: false,
		PropertyGroup: mCommon.BaseInfo,
		PropertyType: common.FieldTypeInt,
		Option: "",
		OwnerID: conf.OwnerID,
		IsPre: true,
		IsReadOnly: false,
		IsSystem: true,
		IsAPI: true,
		CreateTime: &now,
		LastTime: &now,
		Creator: common.CCSystemOperatorUserName,
		LastEditor: common.CCSystemOperatorUserName,
		Description: "",
	}
	uniqueFields := []string{common.BKObjIDField, common.BKPropertyIDField, common.BKOwnerIDField}
	if err := upgrader.Insert(ctx, db, common.BKTableNameObjAttDes, row, "id", uniqueFields); err != nil {
		blog.ErrorJSON("add module field bk_parent_id failed, Insert err: %s, attribute: %#v, ", err, row)
		return err
	}
	return nil
}

type attribute struct {
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
	LastEditor        string      `json:"bk_last_editor" bson:"bk_last_editor"`
	LastTime          *time.Time  `json:"last_time" bson:"last_time"`
}