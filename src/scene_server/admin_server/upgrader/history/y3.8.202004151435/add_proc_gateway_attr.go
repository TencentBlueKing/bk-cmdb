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

package y3_8_202004151435

import (
	"context"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	mCommon "configcenter/src/scene_server/admin_server/common"
	"configcenter/src/scene_server/admin_server/upgrader/history"
	"configcenter/src/storage/dal"
)

// addProcNetworkProxyGroup 添加外网代理信息分组
func addProcNetworkProxyGroup(ctx context.Context, db dal.RDB, conf *history.Config) error {
	group := Group{
		ObjectID:   common.BKInnerObjIDProc,
		GroupID:    mCommon.ProcNetworkProxyInfo,
		GroupName:  mCommon.ProcNetworkProxyInfoName,
		GroupIndex: 5,
		OwnerID:    conf.TenantID,
		IsDefault:  true,
		IsPre:      true,
		IsCollapse: true,
	}

	uniqueFields := []string{common.BKObjIDField, common.BKPropertyGroupIDField, "bk_supplier_account"}
	err := history.Insert(ctx, db, common.BKTableNamePropertyGroup, group, "id", uniqueFields)
	if err != nil {
		if db.IsNotFoundError(err) == false {
			blog.ErrorJSON("addProcNetworkProxyGroup failed, Insert err: %s, group: %#v, ", err, group)

			return err
		}
	}

	return nil
}

// addProcNetworkProxyAttrs 添加外网代理信息相关字段
func addProcNetworkProxyAttrs(ctx context.Context, db dal.RDB, conf *history.Config) error {
	objID := common.BKInnerObjIDProc
	dataRows := []*Attribute{
		{ObjectID: objID, PropertyID: "bk_gateway_ip", PropertyName: "网关IP", IsRequired: false, IsOnly: false,
			IsEditable: true, PropertyGroup: mCommon.ProcNetworkProxyInfo, PropertyType: common.FieldTypeLongChar,
			Option: ""},
		{ObjectID: objID, PropertyID: "bk_gateway_port", PropertyName: "网关端口", IsRequired: false, IsOnly: false,
			IsEditable: true, PropertyGroup: mCommon.ProcNetworkProxyInfo, PropertyType: common.FieldTypeSingleChar,
			Option:      common.PatternMultiplePortRange,
			Placeholder: `单个端口：8080 </br>多个连续端口：8080-8089 </br>多个不连续端口：8080-8089,8199`},
		{ObjectID: objID, PropertyID: "bk_gateway_protocol", PropertyName: "网关协议", IsRequired: false, IsOnly: false,
			IsEditable: true, PropertyGroup: mCommon.ProcNetworkProxyInfo, PropertyType: common.FieldTypeEnum,
			Option: []metadata.EnumVal{{ID: "1", Name: "TCP", Type: "text"}, {ID: "2", Name: "UDP", Type: "text"}}},
		{ObjectID: objID, PropertyID: "bk_gateway_city", PropertyName: "网关所在城市", IsRequired: false, IsOnly: false,
			IsEditable: true, PropertyGroup: mCommon.ProcNetworkProxyInfo, PropertyType: common.FieldTypeSingleChar,
			Option: ""},
	}

	now := time.Now()
	uniqueFields := []string{common.BKObjIDField, common.BKPropertyIDField, "bk_supplier_account"}
	for _, r := range dataRows {
		r.OwnerID = conf.TenantID
		r.IsPre = true
		r.IsReadOnly = false
		r.CreateTime = &now
		r.LastTime = &now
		r.Creator = common.CCSystemOperatorUserName
		r.LastEditor = common.CCSystemOperatorUserName
		r.Description = ""

		if err := history.Insert(ctx, db, common.BKTableNameObjAttDes, r, "id", uniqueFields); err != nil {
			blog.ErrorJSON("addProcNetworkProxyAttrs failed, Upsert err: %s, attribute: %#v, ", err, r)
			return err
		}
	}

	return nil
}

// Attribute TODO
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

// Group group metadata definition
type Group struct {
	BizID      int64  `field:"bk_biz_id" json:"bk_biz_id" bson:"bk_biz_id"`
	ID         int64  `field:"id" json:"id" bson:"id"`
	GroupID    string `field:"bk_group_id" json:"bk_group_id" bson:"bk_group_id"`
	GroupName  string `field:"bk_group_name" json:"bk_group_name" bson:"bk_group_name"`
	GroupIndex int64  `field:"bk_group_index" json:"bk_group_index" bson:"bk_group_index"`
	ObjectID   string `field:"bk_obj_id" json:"bk_obj_id" bson:"bk_obj_id"`
	OwnerID    string `field:"bk_supplier_account" json:"bk_supplier_account" bson:"bk_supplier_account"`
	IsDefault  bool   `field:"bk_isdefault" json:"bk_isdefault" bson:"bk_isdefault"`
	IsPre      bool   `field:"ispre" json:"ispre" bson:"ispre"`
	IsCollapse bool   `field:"is_collapse" json:"is_collapse" bson:"is_collapse"`
}
