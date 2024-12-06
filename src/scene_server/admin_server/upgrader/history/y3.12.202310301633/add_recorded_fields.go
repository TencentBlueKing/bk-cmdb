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

package y3_12_202310301633

import (
	"context"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	mCommon "configcenter/src/scene_server/admin_server/common"
	"configcenter/src/scene_server/admin_server/upgrader/history"
	"configcenter/src/storage/dal"
)

var recordedAttrs = []attribute{
	{
		PropertyID:    common.BKCreatedAt,
		PropertyName:  "创建时间",
		IsOnly:        false,
		IsEditable:    false,
		IsRequired:    false,
		IsMultiple:    false,
		IsPre:         true,
		PropertyGroup: mCommon.BaseInfo,
		PropertyType:  common.FieldTypeTime,
		Creator:       common.CCSystemOperatorUserName,
		Option:        "",
	},
	{
		PropertyID:    common.BKCreatedBy,
		PropertyName:  "创建人",
		IsOnly:        false,
		IsEditable:    false,
		IsRequired:    false,
		IsMultiple:    false,
		IsPre:         true,
		PropertyGroup: mCommon.BaseInfo,
		PropertyType:  common.FieldTypeUser,
		Creator:       common.CCSystemOperatorUserName,
		Option:        "",
	},
	{
		PropertyID:    common.BKUpdatedAt,
		PropertyName:  "更新时间",
		IsOnly:        false,
		IsEditable:    false,
		IsRequired:    false,
		IsMultiple:    false,
		IsPre:         true,
		PropertyGroup: mCommon.BaseInfo,
		PropertyType:  common.FieldTypeTime,
		Creator:       common.CCSystemOperatorUserName,
		Option:        "",
	},
	{
		PropertyID:    common.BKUpdatedBy,
		PropertyName:  "更新人",
		IsOnly:        false,
		IsEditable:    false,
		IsRequired:    false,
		IsMultiple:    false,
		IsPre:         true,
		PropertyGroup: mCommon.BaseInfo,
		PropertyType:  common.FieldTypeUser,
		Creator:       common.CCSystemOperatorUserName,
		Option:        "",
	},
}

func addRecordedFields(ctx context.Context, db dal.RDB, conf *history.Config) error {
	objs := make([]metadata.Object, 0)
	if err := db.Table(common.BKTableNameObjDes).Find(nil).Fields(common.BKObjIDField).All(ctx, &objs); err != nil {
		blog.Errorf("get all objIDs failed, err: %v, rid: %s", err)
		return err
	}

	nowTime := metadata.Now()
	uniqueFields := []string{common.BKObjIDField, common.BKPropertyIDField, common.BKAppIDField}
	for _, obj := range objs {
		for _, attr := range recordedAttrs {
			attr.ObjectID = obj.ObjectID
			attr.OwnerID = conf.TenantID
			attr.CreateTime = &nowTime
			attr.LastTime = &nowTime

			_, _, err := history.Upsert(ctx, db, common.BKTableNameObjAttDes, attr, "id", uniqueFields, []string{})
			if err != nil {
				blog.Errorf("add object attribute failed, attribute: %v, err: %v", attr, err)
				return err
			}
		}
	}

	return nil
}

type attribute struct {
	BizID             int64          `field:"bk_biz_id" json:"bk_biz_id" bson:"bk_biz_id"`
	ID                int64          `field:"id" json:"id" bson:"id"`
	OwnerID           string         `field:"bk_supplier_account" json:"bk_supplier_account" bson:"bk_supplier_account"`
	ObjectID          string         `field:"bk_obj_id" json:"bk_obj_id" bson:"bk_obj_id"`
	PropertyID        string         `field:"bk_property_id" json:"bk_property_id" bson:"bk_property_id"`
	PropertyName      string         `field:"bk_property_name" json:"bk_property_name" bson:"bk_property_name"`
	PropertyGroup     string         `field:"bk_property_group" json:"bk_property_group" bson:"bk_property_group"`
	PropertyGroupName string         `field:"bk_property_group_name,ignoretomap" json:"bk_property_group_name" bson:"-"`
	PropertyIndex     int64          `field:"bk_property_index" json:"bk_property_index" bson:"bk_property_index"`
	Unit              string         `field:"unit" json:"unit" bson:"unit"`
	Placeholder       string         `field:"placeholder" json:"placeholder" bson:"placeholder"`
	IsEditable        bool           `field:"editable" json:"editable" bson:"editable"`
	IsPre             bool           `field:"ispre" json:"ispre" bson:"ispre"`
	IsRequired        bool           `field:"isrequired" json:"isrequired" bson:"isrequired"`
	IsReadOnly        bool           `field:"isreadonly" json:"isreadonly" bson:"isreadonly"`
	IsOnly            bool           `field:"isonly" json:"isonly" bson:"isonly"`
	IsSystem          bool           `field:"bk_issystem" json:"bk_issystem" bson:"bk_issystem"`
	IsAPI             bool           `field:"bk_isapi" json:"bk_isapi" bson:"bk_isapi"`
	PropertyType      string         `field:"bk_property_type" json:"bk_property_type" bson:"bk_property_type"`
	Option            interface{}    `field:"option" json:"option" bson:"option"`
	Default           interface{}    `field:"default" json:"default,omitempty" bson:"default"`
	IsMultiple        bool           `field:"ismultiple" json:"ismultiple,omitempty" bson:"ismultiple"`
	Description       string         `field:"description" json:"description" bson:"description"`
	TemplateID        int64          `field:"bk_template_id" json:"bk_template_id" bson:"bk_template_id"`
	Creator           string         `field:"creator" json:"creator" bson:"creator"`
	CreateTime        *metadata.Time `json:"create_time" bson:"create_time"`
	LastTime          *metadata.Time `json:"last_time" bson:"last_time"`
}
