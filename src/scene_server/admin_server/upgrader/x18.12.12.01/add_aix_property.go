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
package x18_12_12_01

import (
	"context"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/condition"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
)

type Attribute struct {
	ID                int64       `json:"id" bson:"id"`
	OwnerID           string      `json:"bk_supplier_account" bson:"bk_supplier_account"`
	ObjectID          string      `json:"bk_obj_id" bson:"bk_obj_id"`
	PropertyID        string      `json:"bk_property_id" bson:"bk_property_id"`
	PropertyName      string      `json:"bk_property_name" bson:"bk_property_name"`
	PropertyGroup     string      `json:"bk_property_group" bson:"bk_property_group"`
	PropertyGroupName string      `json:"bk_property_group_name" bson:"-"`
	PropertyIndex     int64       `json:"bk_property_index" bson:"bk_property_index"`
	Unit              string      `json:"unit" bson:"unit"`
	Placeholder       string      `json:"placeholder" bson:"placeholder"`
	IsEditable        bool        `json:"editable" bson:"editable"`
	IsPre             bool        `json:"ispre" bson:"ispre"`
	IsRequired        bool        `json:"isrequired" bson:"isrequired"`
	IsReadOnly        bool        `json:"isreadonly" bson:"isreadonly"`
	IsOnly            bool        `json:"isonly" bson:"isonly"`
	IsSystem          bool        `json:"bk_issystem" bson:"bk_issystem"`
	IsAPI             bool        `json:"bk_isapi" bson:"bk_isapi"`
	PropertyType      string      `json:"bk_property_type" bson:"bk_property_type"`
	Option            interface{} `json:"option" bson:"option"`
	Description       string      `json:"description" bson:"description"`
	Creator           string      `json:"creator" bson:"creator"`
	CreateTime        *time.Time  `json:"create_time" bson:"create_time"`
	LastTime          *time.Time  `json:"last_time" bson:"last_time"`
}

func addAIXProperty(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {

	cond := condition.CreateCondition()
	cond.Field(common.BKOwnerIDField).Eq(common.BKDefaultOwnerID)
	cond.Field(common.BKObjIDField).Eq(common.BKInnerObjIDHost)
	cond.Field(common.BKPropertyIDField).Eq(common.BKOSTypeField)

	ostypeProperty := Attribute{}
	err := db.Table(common.BKTableNameObjAttDes).Find(cond.ToMapStr()).One(ctx, &ostypeProperty)
	if err != nil {
		return err
	}

	enumOpts, err := metadata.ParseEnumOption(ctx, ostypeProperty.Option)
	if err != nil {
		return err
	}
	for _, enum := range enumOpts {
		if enum.ID == "3" {
			return nil
		}
	}

	aixEnum := metadata.EnumVal{
		ID:   "3",
		Name: "AIX",
		Type: "text",
	}
	enumOpts = append(enumOpts, aixEnum)

	data := mapstr.MapStr{
		common.BKOptionField: enumOpts,
	}

	err = db.Table(common.BKTableNameObjAttDes).Update(ctx, cond.ToMapStr(), data)
	if err != nil {
		return err
	}
	return nil
}
