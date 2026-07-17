/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 Tencent. All rights reserved.
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

package y3_15_202607151050

import (
	"configcenter/pkg/tenant"
	tenanttmp "configcenter/pkg/types/tenant-template"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/admin_server/upgrader/tools"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/driver/mongodb"
)

var hostOsKernelVersionAttr = attribute{
	ObjectID:      common.BKInnerObjIDHost,
	BizID:         0,
	PropertyID:    common.BKOsKernelVersionField,
	PropertyName:  "操作系统内核版本",
	IsEditable:    true,
	IsPre:         true,
	IsRequired:    false,
	IsOnly:        false,
	PropertyGroup: "auto",
	PropertyType:  common.FieldTypeSingleChar,
	Option:        "",
	Creator:       common.CCSystemOperatorUserName,
	Placeholder:   "可人工设置。若需自动发现，请在主机上安装GseAgent和采集器。",
	Time:          tools.NewTime(),
}

func addHostOsKernelVersionField(kit *rest.Kit, db dal.Dal) error {
	insertOpt := &tools.InsertOptions{
		UniqueFields: []string{common.BKObjIDField, common.BKPropertyIDField, common.BKAppIDField},
		IDField:      []string{common.BKFieldID},
		IgnoreExists: true,
		AuditTypeField: &tools.AuditResType{
			AuditType:    metadata.ModelType,
			ResourceType: metadata.ModelAttributeRes,
		},
		AuditDataField: &tools.AuditDataField{
			BizIDField:   common.BKAppIDField,
			ResIDField:   common.BKFieldID,
			ResNameField: common.BKPropertyNameField,
		},
	}

	// insert the host os kernel version field for each tenant
	err := tenant.ExecForAllTenants(func(tenantID string) error {
		tenantKit := kit.NewKit().WithTenant(tenantID)
		tenantDB := db.Shard(tenantKit.ShardOpts())

		attr := hostOsKernelVersionAttr

		cond := map[string]interface{}{
			common.BKObjIDField:         common.BKInnerObjIDHost,
			common.BKPropertyGroupField: "auto",
		}
		index := make(map[string]int64)
		if err := tenantDB.Table(common.BKTableNameObjAttDes).Find(cond).Fields(common.BKPropertyIndexField).
			Sort(common.BKPropertyIndexField+":-1").One(kit.Ctx, &index); err != nil && !mongodb.IsNotFoundError(err) {
			blog.Errorf("get property index failed, err: %v", err)
			return err
		}
		attr.PropertyIndex = index[common.BKPropertyIndexField] + 1

		_, err := tools.InsertData(tenantKit, tenantDB, common.BKTableNameObjAttDes, []interface{}{attr}, insertOpt)
		if err != nil {
			blog.Errorf("upsert host os kernel version field for tenant %s failed, err: %v", tenantID, err)
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	// add the tenant template of the host os kernel version field
	idOpt := &tools.IDOptions{ResNameField: common.BKPropertyNameField, RemoveKeys: []string{common.BKFieldID}}
	if err = tools.InsertTemplateData(kit, db.Shard(kit.ShardOpts()), []interface{}{hostOsKernelVersionAttr},
		insertOpt, idOpt, tenanttmp.TemplateTypeObjAttribute); err != nil {
		blog.Errorf("insert host os kernel version field template failed, err: %v", err)
		return err
	}
	return nil
}

type attribute struct {
	ID                int64       `bson:"id"`
	BizID             int64       `bson:"bk_biz_id"`
	ObjectID          string      `bson:"bk_obj_id"`
	PropertyID        string      `bson:"bk_property_id"`
	PropertyName      string      `bson:"bk_property_name"`
	PropertyGroup     string      `bson:"bk_property_group"`
	PropertyGroupName string      `bson:"bk_property_group_name"`
	PropertyIndex     int64       `bson:"bk_property_index"`
	Unit              string      `bson:"unit"`
	Placeholder       string      `bson:"placeholder"`
	IsEditable        bool        `bson:"editable"`
	IsPre             bool        `bson:"ispre"`
	IsRequired        bool        `bson:"isrequired"`
	IsReadOnly        bool        `bson:"isreadonly"`
	IsOnly            bool        `bson:"isonly"`
	IsSystem          bool        `bson:"bk_issystem"`
	IsAPI             bool        `bson:"bk_isapi"`
	PropertyType      string      `bson:"bk_property_type"`
	Option            interface{} `bson:"option"`
	Default           interface{} `bson:"default"`
	IsMultiple        *bool       `bson:"ismultiple"`
	Description       string      `bson:"description"`
	TemplateID        int64       `bson:"bk_template_id"`
	Creator           string      `bson:"creator"`
	Time              *tools.Time `bson:",inline"`
}
