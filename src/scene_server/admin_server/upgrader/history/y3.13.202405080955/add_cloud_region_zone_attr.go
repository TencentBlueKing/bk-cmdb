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

package y3_13_202405080955

import (
	"context"
	"fmt"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	mCommon "configcenter/src/scene_server/admin_server/common"
	"configcenter/src/scene_server/admin_server/upgrader/history"
	"configcenter/src/storage/dal"
)

var attributes = []attribute{
	{
		PropertyID:   common.BKCloudRegionField,
		PropertyName: "云地域（Region）",
		IsEditable:   false,
		IsRequired:   false,
		PropertyType: common.FieldTypeSingleChar,
		Option:       "",
	},
	{
		PropertyID:   common.BKCloudZoneField,
		PropertyName: "云可用区（Zone）",
		IsEditable:   false,
		IsRequired:   false,
		PropertyType: common.FieldTypeSingleChar,
		Option:       "",
	},
}

func addAttribute(ctx context.Context, db dal.RDB, conf *history.Config) error {
	attrFilter := mapstr.MapStr{
		common.BKObjIDField: common.BKInnerObjIDHost,
		common.BKPropertyIDField: mapstr.MapStr{
			common.BKDBIN: []string{common.BKCloudRegionField, common.BKCloudZoneField},
		},
		"bk_supplier_account": conf.TenantID,
	}

	existAttrs := make([]attribute, 0)
	err := db.Table(common.BKTableNameObjAttDes).Find(attrFilter).All(ctx, &existAttrs)
	if err != nil {
		blog.Errorf("get exist cloud region and zone attribute failed, err: %v, filter: %+v", err, attrFilter)
		return err
	}

	existAttrMap := make(map[string]struct{})
	for _, attr := range existAttrs {
		if attr.Creator != conf.User {
			blog.Errorf("exist attribute(%+v) is not created by %s", attr, conf.User)
			return fmt.Errorf("exist attribute(%s) is not created by %s", attr.PropertyName, conf.User)
		}
		existAttrMap[attr.PropertyID] = struct{}{}
	}

	if len(existAttrMap) == len(attributes) {
		return nil
	}

	attrIDs, err := db.NextSequences(ctx, common.BKTableNameObjAttDes, len(attributes)-len(existAttrMap))
	if err != nil {
		blog.Errorf("get new attribute ids for cloud region and zone failed, err: %v", err)
		return err
	}

	attrIndexFilter := mapstr.MapStr{
		common.BKObjIDField: common.BKInnerObjIDHost,
	}
	sort := common.BKPropertyIndexField + ":-1"
	lastAttr := new(attribute)

	if err = db.Table(common.BKTableNameObjAttDes).Find(attrIndexFilter).Sort(sort).One(ctx, lastAttr); err != nil {
		blog.Errorf("get host attribute max property index id failed, err: %v", err)
		return err
	}

	now := metadata.Now()
	isMultiple := false

	createAttrs := make([]attribute, 0)
	createAttrIdx := 0
	for _, attr := range attributes {
		_, exists := existAttrMap[attr.PropertyID]
		if exists {
			continue
		}

		attr.ID = int64(attrIDs[createAttrIdx])
		attr.OwnerID = conf.TenantID
		attr.ObjectID = common.BKInnerObjIDHost
		attr.PropertyGroup = mCommon.BaseInfo
		attr.PropertyIndex = lastAttr.PropertyIndex + 1 + int64(createAttrIdx)
		attr.IsPre = true
		attr.Creator = conf.User
		attr.CreateTime = &now
		attr.LastTime = &now
		attr.IsMultiple = &isMultiple
		createAttrs = append(createAttrs, attr)
		createAttrIdx++
	}

	if err = db.Table(common.BKTableNameObjAttDes).Insert(ctx, createAttrs); err != nil {
		blog.Errorf("create attributes failed, err: %v, attrs: %+v", err, createAttrs)
		return err
	}
	return nil
}

type attribute struct {
	BizID             int64          `field:"bk_biz_id" json:"bk_biz_id" bson:"bk_biz_id" mapstructure:"bk_biz_id"`
	ID                int64          `field:"id" json:"id" bson:"id" mapstructure:"id"`
	OwnerID           string         `field:"bk_supplier_account" json:"bk_supplier_account" bson:"bk_supplier_account" mapstructure:"bk_supplier_account"`
	ObjectID          string         `field:"bk_obj_id" json:"bk_obj_id" bson:"bk_obj_id" mapstructure:"bk_obj_id"`
	PropertyID        string         `field:"bk_property_id" json:"bk_property_id" bson:"bk_property_id" mapstructure:"bk_property_id"`
	PropertyName      string         `field:"bk_property_name" json:"bk_property_name" bson:"bk_property_name" mapstructure:"bk_property_name"`
	PropertyGroup     string         `field:"bk_property_group" json:"bk_property_group" bson:"bk_property_group" mapstructure:"bk_property_group"`
	PropertyGroupName string         `field:"bk_property_group_name,ignoretomap" json:"bk_property_group_name" bson:"-" mapstructure:"bk_property_group_name"`
	PropertyIndex     int64          `field:"bk_property_index" json:"bk_property_index" bson:"bk_property_index" mapstructure:"bk_property_index"`
	Unit              string         `field:"unit" json:"unit" bson:"unit" mapstructure:"unit"`
	Placeholder       string         `field:"placeholder" json:"placeholder" bson:"placeholder" mapstructure:"placeholder"`
	IsEditable        bool           `field:"editable" json:"editable" bson:"editable" mapstructure:"editable"`
	IsPre             bool           `field:"ispre" json:"ispre" bson:"ispre" mapstructure:"ispre"`
	IsRequired        bool           `field:"isrequired" json:"isrequired" bson:"isrequired" mapstructure:"isrequired"`
	IsReadOnly        bool           `field:"isreadonly" json:"isreadonly" bson:"isreadonly" mapstructure:"isreadonly"`
	IsOnly            bool           `field:"isonly" json:"isonly" bson:"isonly" mapstructure:"isonly"`
	IsSystem          bool           `field:"bk_issystem" json:"bk_issystem" bson:"bk_issystem" mapstructure:"bk_issystem"`
	IsAPI             bool           `field:"bk_isapi" json:"bk_isapi" bson:"bk_isapi" mapstructure:"bk_isapi"`
	PropertyType      string         `field:"bk_property_type" json:"bk_property_type" bson:"bk_property_type" mapstructure:"bk_property_type"`
	Option            interface{}    `field:"option" json:"option" bson:"option" mapstructure:"option"`
	Default           interface{}    `field:"default" json:"default,omitempty" bson:"default" mapstructure:"default"`
	IsMultiple        *bool          `field:"ismultiple" json:"ismultiple,omitempty" bson:"ismultiple" mapstructure:"ismultiple"`
	Description       string         `field:"description" json:"description" bson:"description" mapstructure:"description"`
	TemplateID        int64          `field:"bk_template_id" json:"bk_template_id" bson:"bk_template_id" mapstructure:"bk_template_id"`
	Creator           string         `field:"creator" json:"creator" bson:"creator" mapstructure:"creator"`
	CreateTime        *metadata.Time `json:"create_time" bson:"create_time" mapstructure:"create_time"`
	LastTime          *metadata.Time `json:"last_time" bson:"last_time" mapstructure:"last_time"`
}
