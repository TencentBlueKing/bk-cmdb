// Package y3_7_202002231026 TODO
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
package y3_7_202002231026

import (
	"context"
	"fmt"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/scene_server/admin_server/upgrader/history"
	"configcenter/src/storage/dal"
)

func addProcEnablePortProperty(ctx context.Context, db dal.RDB, conf *history.Config) error {

	blog.Infof("start execute y3_7_202002231026")

	propertyGroup := "proc_port"
	maxIdxCond := map[string]interface{}{
		common.BKPropertyGroupField: propertyGroup,
		common.BKObjIDField:         common.BKInnerObjIDProc,
	}
	maxIdxAttr := &Attribute{}
	err := db.Table(common.BKTableNameObjAttDes).Find(maxIdxCond).Sort(fmt.Sprintf("%s:-1",
		common.BKPropertyIndexField)).One(ctx, maxIdxAttr)
	if err != nil {
		blog.Errorf("get proerty max index value error, cond: %v, err: %v", maxIdxCond, err)
		return fmt.Errorf("get proerty max index value error, err: %v", err)
	}

	addPortEnable := Attribute{
		OwnerID:       "0",
		ObjectID:      common.BKInnerObjIDProc,
		PropertyID:    common.BKProcPortEnable,
		PropertyName:  "启用端口",
		PropertyGroup: propertyGroup,
		PropertyIndex: maxIdxAttr.PropertyIndex + 1,
		Unit:          "",
		Placeholder:   "",
		IsEditable:    true,
		IsPre:         false,
		IsReadOnly:    false,
		IsOnly:        false,
		IsSystem:      false,
		IsAPI:         false,
		PropertyType:  common.FieldTypeBool,
		Option:        true,
		Description:   "",
		Creator:       common.CCSystemOperatorUserName,
		CreateTime:    time.Now(),
		LastTime:      time.Now(),
	}

	uniqueFields := []string{common.BKObjIDField, common.BKPropertyIDField, "bk_supplier_account"}
	if err := history.Insert(ctx, db, common.BKTableNameObjAttDes, addPortEnable, "id", uniqueFields); err != nil {
		blog.Errorf("Insert err: %v, attribute: %#v", err, addPortEnable)
		return err
	}

	return nil
}

func addProcTemplatePortEnableProperty(ctx context.Context, db dal.RDB, conf *history.Config) error {
	docPrefix := "property." + common.BKProcPortEnable
	doc := map[string]interface{}{
		docPrefix + ".value":            true,
		docPrefix + ".as_default_value": true,
	}

	updateCond := map[string]interface{}{
		common.BKDBOR: []interface{}{
			// upgrade version
			map[string]interface{}{docPrefix: map[string]interface{}{common.BKDBExists: false}},
			// new install cmdb
			map[string]interface{}{docPrefix + ".as_default_value": nil},
		},
	}

	if err := db.Table(common.BKTableNameProcessTemplate).Update(ctx, updateCond, doc); err != nil {
		blog.Errorf("update process template failed, condition: %+v, doc: %+v, err: %v", common.BKProcPortEnable,
			updateCond, doc, err)
		return fmt.Errorf("db operator error, proerpty id: %v, err: %v", common.BKProcPortEnable, err)
	}

	return nil
}

func setProcInfoProtEnableDefaultValue(ctx context.Context, db dal.RDB, conf *history.Config) error {

	doc := map[string]interface{}{
		common.BKProcPortEnable: true,
	}
	updateCond := map[string]interface{}{
		common.BKDBOR: []interface{}{
			map[string]interface{}{common.BKProcPortEnable: map[string]interface{}{common.BKDBExists: false}},
			map[string]interface{}{common.BKProcPortEnable: nil},
		},
	}
	if err := db.Table(common.BKTableNameBaseProcess).Update(ctx, updateCond, doc); err != nil {
		blog.Errorf("set process information id %s default value failed, condition: %+v, doc:%+v, err: %v",
			common.BKProcPortEnable, updateCond, doc, err)
		return fmt.Errorf("set process information id %s default value, db operator error, err: %v",
			common.BKProcPortEnable, err)
	}
	return nil
}

// Attribute attribute strcut
type Attribute struct {
	Metadata          `field:"metadata" json:"metadata" bson:"metadata"`
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

	Creator    string    `field:"creator" json:"creator" bson:"creator"`
	CreateTime time.Time `json:"create_time" bson:"create_time"`
	LastTime   time.Time `json:"last_time" bson:"last_time"`
}

// Metadata  used to define the metadata for the resources
type Metadata struct {
	Label Label `field:"label" json:"label" bson:"label"`
}

// Label define the Label type used to limit the scope of application of resources
type Label map[string]string
