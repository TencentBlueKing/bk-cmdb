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
package x18_10_30_02

import (
	"context"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/mapstr"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
)

type attribute struct {
	ID            int64       `json:"id" bson:"id"`
	OwnerID       string      `json:"bk_supplier_account" bson:"bk_supplier_account"`
	ObjectID      string      `json:"bk_obj_id" bson:"bk_obj_id"`
	PropertyID    string      `json:"bk_property_id" bson:"bk_property_id"`
	PropertyName  string      `json:"bk_property_name" bson:"bk_property_name"`
	PropertyGroup string      `json:"bk_property_group" bson:"bk_property_group"`
	PropertyIndex int64       `json:"bk_property_index" bson:"bk_property_index"`
	Unit          string      `json:"unit" bson:"unit"`
	Placeholder   string      `json:"placeholder" bson:"placeholder"`
	IsEditable    bool        `json:"editable" bson:"editable"`
	IsPre         bool        `json:"ispre" bson:"ispre"`
	IsRequired    bool        `json:"isrequired" bson:"isrequired"`
	IsReadOnly    bool        `json:"isreadonly" bson:"isreadonly"`
	IsOnly        bool        `json:"isonly" bson:"isonly"`
	IsSystem      bool        `json:"bk_issystem" bson:"bk_issystem"`
	IsAPI         bool        `json:"bk_isapi" bson:"bk_isapi"`
	PropertyType  string      `json:"bk_property_type" bson:"bk_property_type"`
	Option        interface{} `json:"option" bson:"option"`
	Description   string      `json:"description" bson:"description"`
	Creator       string      `json:"creator" bson:"creator"`
	CreateTime    *time.Time  `bson:"create_time"`
	LastTime      *time.Time  `bson:"last_time"`
}

func addBizSuupierID(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	attributeArr := make([]attribute, 0)
	var attrID uint64
	attrID, err := db.NextSequence(ctx, common.BKTableNameObjAttDes)
	if nil != err {
		return err
	}
	ts := time.Now().UTC()
	attributeArr = append(attributeArr, attribute{
		ID:            int64(attrID),
		OwnerID:       conf.OwnerID,
		Creator:       conf.User,
		PropertyGroup: "default",
		IsPre:         false,
		IsRequired:    true,
		ObjectID:      common.BKInnerObjIDApp,
		PropertyName:  common.BKSupplierIDField,
		PropertyID:    common.BKSupplierIDField,
		PropertyType:  common.FieldTypeInt,
		CreateTime:    &ts,
		LastTime:      &ts,
		IsAPI:         true,
	})
	for _, attr := range attributeArr {
		filter := mapstr.MapStr{
			common.BKPropertyIDField: attr.PropertyID,
			common.BKOwnerIDField:    attr.OwnerID,
		}
		cnt, err := db.Table(common.BKTableNameObjAttDes).Find(filter).Count(ctx)
		if err != nil {
			return err
		}
		if cnt > 0 {
			continue
		}
		err = db.Table(common.BKTableNameObjAttDes).Insert(ctx, attr)
		if err != nil {
			return err
		}
	}
	return nil
}
