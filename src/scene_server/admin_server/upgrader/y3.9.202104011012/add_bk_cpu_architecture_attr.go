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

package y3_9_202104011012

import (
	"context"
	"fmt"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	mCommon "configcenter/src/scene_server/admin_server/common"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
)

// addHostBkCPUArchitectureAttr add bk_cpu_architecture attribute for host
func addHostBkCPUArchitectureAttr(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	attrFilter := map[string]interface{}{
		common.BKObjIDField:      common.BKInnerObjIDHost,
		common.BKPropertyIDField: "bk_cpu_architecture",
		common.BkSupplierAccount: conf.OwnerID,
	}

	cnt, err := db.Table(common.BKTableNameObjAttDes).Find(attrFilter).Count(ctx)
	if err != nil {
		blog.Errorf("count bk_cpu_architecture attribute failed, err: %v", err)
		return err
	}

	if cnt > 0 {
		return nil
	}

	attrGrpFilter := map[string]interface{}{
		common.BKObjIDField:           common.BKInnerObjIDHost,
		common.BKPropertyGroupIDField: mCommon.HostAutoFields,
	}

	cnt, err = db.Table(common.BKTableNamePropertyGroup).Find(attrGrpFilter).Count(ctx)
	if err != nil {
		blog.Errorf("count auto attribute group failed, err: %v", err)
		return err
	}

	if cnt == 0 {
		return fmt.Errorf("attribute group %s does not exist", mCommon.HostAutoFields)
	}

	attrID, err := db.NextSequence(ctx, common.BKTableNameObjAttDes)
	if err != nil {
		blog.Errorf("get new attribute id for bk_cpu_architecture failed, err: %v", err)
		return err
	}

	attrIndexFilter := map[string]interface{}{
		common.BKObjIDField: common.BKInnerObjIDHost,
	}
	sort := common.BKPropertyIndexField + ":-1"
	lastAttr := new(Attribute)

	if err := db.Table(common.BKTableNameObjAttDes).Find(attrIndexFilter).Sort(sort).One(ctx, lastAttr); err != nil {
		blog.Errorf("get host attribute max property index id failed, err: %v", err)
		return err
	}

	now := metadata.Now()
	attr := &Attribute{
		ID:            int64(attrID),
		OwnerID:       conf.OwnerID,
		ObjectID:      common.BKInnerObjIDHost,
		PropertyID:    "bk_cpu_architecture",
		PropertyName:  "CPU架构",
		PropertyGroup: mCommon.HostAutoFields,
		PropertyIndex: lastAttr.PropertyIndex + 1,
		Placeholder:   "选择CPU架构类型，如X86或ARM",
		IsEditable:    true,
		IsPre:         true,
		IsRequired:    false,
		PropertyType:  common.FieldTypeEnum,
		Option: []EnumVal{
			{ID: "x86", Name: "X86", Type: "text", IsDefault: true},
			{ID: "arm", Name: "ARM", Type: "text"},
		},
		Creator:    conf.User,
		CreateTime: &now,
		LastTime:   &now,
	}

	return db.Table(common.BKTableNameObjAttDes).Insert(ctx, attr)
}

// Attribute attribute definition
type Attribute struct {
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
	Description       string         `field:"description" json:"description" bson:"description"`
	Creator           string         `field:"creator" json:"creator" bson:"creator"`
	CreateTime        *metadata.Time `json:"create_time" bson:"create_time"`
	LastTime          *metadata.Time `json:"last_time" bson:"last_time"`
}

type EnumVal struct {
	ID        string `bson:"id"           json:"id"`
	Name      string `bson:"name"         json:"name"`
	Type      string `bson:"type"         json:"type"`
	IsDefault bool   `bson:"is_default"   json:"is_default"`
}
