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

package y3_10_202209231617

import (
	"context"
	"strconv"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	comm "configcenter/src/scene_server/admin_server/common"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
)

// addCloudHostIdentifierAttr add cloud host identifier attribute, hosts with the cloud host identifier attribute set cannot be transferred across biz
func addCloudHostIdentifierAttr(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	identifierAttr := &attribute{
		OwnerID:       conf.OwnerID,
		ObjectID:      common.BKInnerObjIDHost,
		PropertyID:    "bk_cloud_host_identifier",
		PropertyName:  "云主机标识",
		PropertyGroup: comm.BaseInfo,
		IsEditable:    false,
		IsPre:         true,
		IsRequired:    false,
		IsReadOnly:    true,
		IsOnly:        false,
		IsSystem:      true,
		IsAPI:         true,
		PropertyType:  common.FieldTypeBool,
		Description:   "云主机标识",
		Creator:       conf.User,
	}

	// check if the cloud host identifier attribute exists
	attrFilter := map[string]interface{}{
		common.BKObjIDField:      common.BKInnerObjIDHost,
		common.BKPropertyIDField: identifierAttr.PropertyID,
	}

	cnt, err := db.Table(common.BKTableNameObjAttDes).Find(attrFilter).Count(ctx)
	if err != nil {
		blog.Errorf("check if attribute exists failed, filter: %v, err: %v", attrFilter, err)
		return err
	}

	if cnt > 0 {
		return nil
	}

	// add cloud host identifier attribute, generate its id and property index
	newAttrID, err := db.NextSequence(ctx, common.BKTableNameObjAttDes)
	if err != nil {
		blog.Errorf("get new attributes id failed, err: %v", err)
		return err
	}

	attrIdxFilter := map[string]interface{}{
		common.BKObjIDField: common.BKInnerObjIDHost,
	}
	sort := common.BKPropertyIndexField + ":-1"
	maxIdxAttr := new(attribute)
	if err := db.Table(common.BKTableNameObjAttDes).Find(attrIdxFilter).Sort(sort).One(ctx, maxIdxAttr); err != nil {
		blog.Errorf("get max host attribute index failed, filter: %v, err: %v", attrIdxFilter, err)
		return err
	}

	identifierAttr.ID = int64(newAttrID)
	identifierAttr.PropertyIndex = maxIdxAttr.PropertyIndex + 1

	now := time.Now()
	identifierAttr.CreateTime = now
	identifierAttr.LastTime = now

	if err := db.Table(common.BKTableNameObjAttDes).Insert(ctx, identifierAttr); err != nil {
		blog.Errorf("insert host attribute(%#v) failed, err: %v", identifierAttr, err)
		return err
	}

	return nil
}

// updateCloudVendorAttr update cloud vendor attribute, add other nodeman cloud vendor option
func updateCloudVendorAttr(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	// get the cloud host identifier and cloud vendor attribute
	cloudVendors := []string{"AWS", "腾讯云", "GCP", "Azure", "企业私有云", "SalesForce", "Oracle Cloud", "IBM Cloud",
		"阿里云", "中国电信", "UCloud", "美团云", "金山云", "百度云", "华为云", "首都云"}

	enumOption := make([]enumVal, len(cloudVendors))
	for index, name := range cloudVendors {
		enumOption[index] = enumVal{
			ID:   strconv.Itoa(index + 1),
			Name: name,
			Type: "text",
		}
	}

	cond := map[string]interface{}{
		common.BKObjIDField: mapstr.MapStr{
			common.BKDBIN: []string{common.BKInnerObjIDHost, common.BKInnerObjIDPlat},
		},
		common.BKPropertyIDField: common.BKCloudVendor,
	}

	updateData := mapstr.MapStr{common.BKOptionField: enumOption}

	if err := db.Table(common.BKTableNameObjAttDes).Update(ctx, cond, updateData); err != nil {
		blog.Errorf("update cloud vendor attribute failed, cond: %#v, data: %#v, err: %v", cond, updateData, err)
		return err
	}

	return nil
}

// attribute definition
type attribute struct {
	BizID         int64       `bson:"bk_biz_id"`
	ID            int64       `bson:"id"`
	OwnerID       string      `bson:"bk_supplier_account"`
	ObjectID      string      `bson:"bk_obj_id"`
	PropertyID    string      `bson:"bk_property_id"`
	PropertyName  string      `bson:"bk_property_name"`
	PropertyGroup string      `bson:"bk_property_group"`
	PropertyIndex int64       `bson:"bk_property_index"`
	Unit          string      `bson:"unit"`
	Placeholder   string      `bson:"placeholder"`
	IsEditable    bool        `bson:"editable"`
	IsPre         bool        `bson:"ispre"`
	IsRequired    bool        `bson:"isrequired"`
	IsReadOnly    bool        `bson:"isreadonly"`
	IsOnly        bool        `bson:"isonly"`
	IsSystem      bool        `bson:"bk_issystem"`
	IsAPI         bool        `bson:"bk_isapi"`
	PropertyType  string      `bson:"bk_property_type"`
	Option        interface{} `bson:"option"`
	Description   string      `bson:"description"`
	Creator       string      `bson:"creator"`
	CreateTime    time.Time   `bson:"create_time"`
	LastTime      time.Time   `bson:"last_time"`
}

// enumVal enum option val
type enumVal struct {
	ID        string `bson:"id"`
	Name      string `bson:"name"`
	Type      string `bson:"type"`
	IsDefault bool   `bson:"is_default"`
}
