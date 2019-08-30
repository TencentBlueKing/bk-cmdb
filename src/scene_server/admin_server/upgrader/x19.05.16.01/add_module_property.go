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

package x19_05_16_01

import (
	"context"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
)

func addModuleProperty(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
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

	now := time.Now()
	serviceCategoryIDProperty := Attribute{
		ID:                0,
		OwnerID:           conf.OwnerID,
		ObjectID:          common.BKInnerObjIDModule,
		PropertyID:        "service_category_id",
		PropertyName:      "服务分类ID",
		PropertyGroup:     "default",
		PropertyGroupName: "default",
		PropertyIndex:     0,
		Unit:              "",
		Placeholder:       "",
		IsEditable:        true,
		IsPre:             true,
		IsRequired:        true,
		IsReadOnly:        false,
		IsOnly:            false,
		IsSystem:          false,
		IsAPI:             true,
		PropertyType:      common.FieldTypeInt,
		Option:            "",
		Description:       "服务分类，外键到 cc_ServiceCategory",
		Creator:           common.CCSystemOperatorUserName,
		CreateTime:        &now,
		LastTime:          &now,
	}

	uniqueFields := []string{common.BKObjIDField, common.BKPropertyIDField, common.BKOwnerIDField}
	_, _, err := upgrader.Upsert(ctx, db, common.BKTableNameObjAttDes, serviceCategoryIDProperty, "id", uniqueFields, []string{})
	if nil != err {
		blog.Errorf("[upgrade v19.05.16.01] addModuleProperty service_category_id failed, err: %+v", err)
		return err
	}

	serviceTemplateIDProperty := Attribute{
		ID:                0,
		OwnerID:           conf.OwnerID,
		ObjectID:          common.BKInnerObjIDModule,
		PropertyID:        "service_template_id",
		PropertyName:      "服务模板ID",
		PropertyGroup:     "default",
		PropertyGroupName: "default",
		PropertyIndex:     0,
		Unit:              "",
		Placeholder:       "",
		IsEditable:        true,
		IsPre:             true,
		IsRequired:        true,
		IsReadOnly:        false,
		IsOnly:            false,
		IsSystem:          false,
		IsAPI:             true,
		PropertyType:      common.FieldTypeInt,
		Option:            "",
		Description:       "服务模板, 外键到 cc_ServiceTemplate",
		Creator:           common.CCSystemOperatorUserName,
		CreateTime:        &now,
		LastTime:          &now,
	}

	if _, _, err := upgrader.Upsert(ctx, db, common.BKTableNameObjAttDes, serviceTemplateIDProperty, "id", uniqueFields, []string{}); err != nil {
		blog.Errorf("[upgrade v19.05.16.01] addModuleProperty service_template_id failed, err: %+v", err)
		return err
	}
	return nil
}
