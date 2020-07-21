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

package y3_9_202002131522

import (
	"context"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	mCommon "configcenter/src/scene_server/admin_server/common"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
)

var (
	groupBaseInfo = mCommon.BaseInfo
)

func initPlatAttr(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	objID := common.BKInnerObjIDPlat
	dataRows := []*Attribute{
		{ObjectID: objID, PropertyID: "bk_status", PropertyName: "状态", IsRequired: false, IsOnly: false, IsEditable: false, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeEnum, Option: statusEnum},
		{ObjectID: objID, PropertyID: "bk_status_detail", PropertyName: "状态详情", IsRequired: false, IsOnly: false, IsEditable: false, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeSingleChar, Option: ""},
		{ObjectID: objID, PropertyID: "bk_cloud_vendor", PropertyName: "云厂商", IsRequired: false, IsOnly: false, IsEditable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeEnum, Option: cloudVendorEnum},
		{ObjectID: objID, PropertyID: "bk_vpc_id", PropertyName: "VPC唯一标识", IsRequired: false, IsOnly: false, IsEditable: false, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeSingleChar, Option: ""},
		{ObjectID: objID, PropertyID: "bk_vpc_name", PropertyName: "VPC名称", IsRequired: false, IsOnly: false, IsEditable: false, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeSingleChar, Option: ""},
		{ObjectID: objID, PropertyID: "bk_account_id", PropertyName: "云账户ID", IsRequired: false, IsOnly: false, IsEditable: false, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeInt, Option: ""},
		{ObjectID: objID, PropertyID: "bk_region", PropertyName: "VPC所属地域", IsRequired: false, IsOnly: false, IsEditable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeSingleChar, Option: ""},
		{ObjectID: objID, PropertyID: "bk_creator", PropertyName: "创建者", IsRequired: false, IsOnly: false, IsEditable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeSingleChar, Option: ""},
		{ObjectID: objID, PropertyID: "bk_last_editor", PropertyName: "最后修改人", IsRequired: false, IsOnly: false, IsEditable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeSingleChar, Option: ""},
	}

	now := time.Now()
	uniqueFields := []string{common.BKObjIDField, common.BKPropertyIDField, common.BKOwnerIDField}
	for _, r := range dataRows {
		r.OwnerID = conf.OwnerID
		r.IsPre = true
		r.IsReadOnly = false
		r.Creator = common.CCSystemOperatorUserName
		r.CreateTime = &now
		r.LastTime = &now
		r.LastEditor = common.CCSystemOperatorUserName
		r.Description = ""

		if _, _, err := upgrader.Upsert(ctx, db, common.BKTableNameObjAttDes, r, "id", uniqueFields, []string{}); err != nil {
			blog.ErrorJSON("initPlatAttr failed, Upsert err: %s, attribute: %#v, ", err, r)
			return err
		}
	}

	return nil
}

var statusEnum = []metadata.EnumVal{
	{ID: "1", Name: "正常", Type: "text", IsDefault: true},
	{ID: "2", Name: "异常", Type: "text"},
}

var cloudVendorEnum = []metadata.EnumVal{
	{ID: "1", Name: "亚马逊云", Type: "text"},
	{ID: "2", Name: "腾讯云", Type: "text"},
}

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
	LastEditor        string      `json:"bk_last_editor" bson:"bk_last_editor"`
	LastTime          *time.Time  `json:"last_time" bson:"last_time"`
}
