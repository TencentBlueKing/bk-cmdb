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

package y3_14_202502101200

import (
	"context"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
)

func addHostOsKernelVersionField(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	nowTime := time.Now()
	attr := &attribute{
		ObjectID:      common.BKInnerObjIDHost,
		BizID:         0,
		PropertyID:    common.BKOsKernelVersionField,
		PropertyName:  "操作系统内核版本",
		IsEditable:    true,
		IsPre:         true,
		IsRequired:    false,
		IsOnly:        false,
		PropertyGroup: "auto",
		PropertyType:  common.FieldTypeSingleChar,
		Option:        "",
		Creator:       common.CCSystemOperatorUserName,
		OwnerID:       "0",
		Placeholder:   "可人工设置。若需自动发现，请在主机上安装GseAgent和采集器。",
		CreateTime:    nowTime,
		LastTime:      nowTime,
	}

	cond := map[string]interface{}{
		common.BKObjIDField:         common.BKInnerObjIDHost,
		common.BKPropertyGroupField: "auto",
	}
	propIndex := make(map[string]int64, 0)
	if err := db.Table(common.BKTableNameObjAttDes).Find(cond).Fields(common.BKPropertyIndexField).
		Sort(common.BKPropertyIndexField+":-1").One(ctx, &propIndex); err != nil {
		blog.Errorf("get property index failed, err: %v", err)
		return err
	}
	attr.PropertyIndex = propIndex[common.BKPropertyIndexField] + 1

	uniqueFields := []string{common.BKObjIDField, common.BKPropertyIDField, common.BKAppIDField}
	_, _, err := upgrader.Upsert(ctx, db, common.BKTableNameObjAttDes, attr, "id", uniqueFields, []string{})
	if err != nil {
		blog.Errorf("add host os kernel version field failed, err: %v", err)
		return err
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
	Default           interface{} `field:"default" json:"default,omitempty" bson:"default"`
	IsMultiple        bool        `field:"ismultiple" json:"ismultiple,omitempty" bson:"ismultiple"`
	Description       string      `field:"description" json:"description" bson:"description"`
	TemplateID        int64       `field:"bk_template_id" json:"bk_template_id" bson:"bk_template_id"`
	Creator           string      `field:"creator" json:"creator" bson:"creator"`
	CreateTime        time.Time   `json:"create_time" bson:"create_time"`
	LastTime          time.Time   `json:"last_time" bson:"last_time"`
}
