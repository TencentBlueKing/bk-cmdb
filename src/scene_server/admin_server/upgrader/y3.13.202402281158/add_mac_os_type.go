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

package y3_13_202402281158

import (
	"context"

	"configcenter/src/common"
	"configcenter/src/common/metadata"
	"configcenter/src/storage/dal"
)

// addMacOSType add host macOS type
func addMacOSType(ctx context.Context, db dal.RDB) error {
	filter := map[string]interface{}{
		common.BKOwnerIDField:    common.BKDefaultOwnerID,
		common.BKObjIDField:      common.BKInnerObjIDHost,
		common.BKPropertyIDField: common.BKOSTypeField,
		common.BKAppIDField:      0,
	}

	osType := attribute{}
	err := db.Table(common.BKTableNameObjAttDes).Find(filter).One(ctx, &osType)
	if err != nil {
		return err
	}

	enumOpts, err := metadata.ParseEnumOption(osType.Option)
	if err != nil {
		return err
	}
	for _, enum := range enumOpts {
		if enum.ID == common.HostOSTypeEnumMacOS {
			return nil
		}
	}

	macOS := metadata.EnumVal{
		ID:   common.HostOSTypeEnumMacOS,
		Name: "MacOS",
		Type: "text",
	}
	enumOpts = append(enumOpts, macOS)

	data := map[string]interface{}{
		common.BKOptionField: enumOpts,
	}

	err = db.Table(common.BKTableNameObjAttDes).Update(ctx, filter, data)
	if err != nil {
		return err
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
	IsMultiple        *bool          `field:"ismultiple" json:"ismultiple,omitempty" bson:"ismultiple"`
	Description       string         `field:"description" json:"description" bson:"description"`
	Creator           string         `field:"creator" json:"creator" bson:"creator"`
	CreateTime        *metadata.Time `json:"create_time" bson:"create_time"`
	LastTime          *metadata.Time `json:"last_time" bson:"last_time"`
}
