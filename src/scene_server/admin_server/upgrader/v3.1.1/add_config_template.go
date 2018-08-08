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

package v3v1v1beta1

import (
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	mCommon "configcenter/src/scene_server/admin_server/common"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/scene_server/validator"
	"configcenter/src/storage"
)

// default group
var (
	groupBaseInfo = mCommon.BaseInfo
)

func addPresetObjects(db storage.DI, conf *upgrader.Config) (err error) {
	err = addObjDesData(db, conf)
	if err != nil {
		return err
	}

	err = addObjAttDescData(db, conf)
	if err != nil {
		return err
	}

	return nil
}

func addObjAttDescData(db storage.DI, conf *upgrader.Config) error {
	tablename := common.BKTableNameObjAttDes
	blog.Infof("add data for  %s table ", tablename)
	rows := getObjAttDescData(conf.OwnerID)
	for _, row := range rows {
		_, _, err := upgrader.Upsert(db, tablename, row, "id", []string{common.BKObjIDField, common.BKPropertyIDField, common.BKOwnerIDField}, []string{})
		if nil != err {
			blog.Errorf("add data for  %s table error  %s", tablename, err)
			return err
		}
	}
	selector := map[string]interface{}{
		common.BKObjIDField: map[string]interface{}{
			common.BKDBIN: []string{"config_template",
				"template_version",
			},
		},
		common.BKPropertyIDField: "bk_name",
	}

	db.DelByCondition(tablename, selector)

	return nil
}

func addObjDesData(db storage.DI, conf *upgrader.Config) error {
	tablename := common.BKTableNameObjDes
	blog.Errorf("add data for  %s table ", tablename)
	rows := getObjectDesData(conf.OwnerID)
	for _, row := range rows {
		if _, _, err := upgrader.Upsert(db, tablename, row, "id", []string{common.BKObjIDField, common.BKClassificationIDField, common.BKOwnerIDField}, []string{"id"}); err != nil {
			blog.Errorf("add data for  %s table error  %s", tablename, err)
			return err
		}
	}

	return nil
}

func getObjectDesData(ownerID string) []*metadata.ObjectDes {

	dataRows := []*metadata.ObjectDes{
		&metadata.ObjectDes{ObjCls: "bk_host_manage", ObjectID: common.BKInnerObjIDConfigTemp, ObjectName: "配置文件模板", IsPre: true, ObjIcon: "icon-cc-process", Position: `{"bk_host_manage":{"x":-550,"y":-750}}`},
		&metadata.ObjectDes{ObjCls: "bk_host_manage", ObjectID: common.BKInnerObjIDTempVersion, ObjectName: "模板版本", IsPre: true, ObjIcon: "icon-cc-process", Position: `{"bk_host_manage":{"x":-600,"y":-850}}`},
	}
	t := time.Now()
	for _, r := range dataRows {
		r.CreateTime = &t
		r.LastTime = &t
		r.IsPaused = false
		r.Creator = common.CCSystemOperatorUserName
		r.OwnerID = ownerID
		r.Description = ""
		r.Modifier = ""
	}

	return dataRows
}

func getObjAttDescData(ownerID string) []*metadata.Attribute {

	predataRows := ConfigTemplateRow()
	predataRows = append(predataRows, TempVersionRow()...)

	t := new(time.Time)
	*t = time.Now()
	for _, r := range predataRows {
		r.OwnerID = ownerID
		r.IsPre = true
		if false != r.IsEditable {
			r.IsEditable = true
		}
		r.IsReadOnly = false
		r.CreateTime = t
		r.Creator = common.CCSystemOperatorUserName
		r.LastTime = r.CreateTime
		r.Description = ""
	}

	return predataRows
}

// ConfigTemplateRow config template structure
func ConfigTemplateRow() []*metadata.Attribute {
	objID := common.BKInnerObjIDConfigTemp
	formatOption := []validator.EnumVal{{ID: "utf8", Name: "utf8", Type: "text", IsDefault: true}, {ID: "gbk", Name: "gbk", Type: "text"}}
	rightOption := []validator.EnumVal{{ID: "644", Name: "644", Type: "text", IsDefault: true}, {ID: "755", Name: "755", Type: "text"}}
	dataRows := []*metadata.Attribute{
		//base info
		&metadata.Attribute{ObjectID: objID, PropertyID: common.BKAppIDField, PropertyName: "业务ID", IsAPI: true, IsRequired: true, IsOnly: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeInt, Option: validator.IntOption{}},
		&metadata.Attribute{ObjectID: objID, PropertyID: "template_name", PropertyName: "配置模板名", IsRequired: true, IsOnly: true, IsPre: true, IsEditable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeSingleChar, Option: ""},
		&metadata.Attribute{ObjectID: objID, PropertyID: "file_name", PropertyName: "配置文件名", IsRequired: false, IsOnly: false, IsPre: true, IsEditable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeSingleChar, Option: ""},
		&metadata.Attribute{ObjectID: objID, PropertyID: "path", PropertyName: "配置文件路径", IsRequired: false, IsOnly: false, IsPre: true, IsEditable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeSingleChar, Option: ""},
		&metadata.Attribute{ObjectID: objID, PropertyID: "user", PropertyName: "所属用户", IsRequired: false, IsOnly: false, IsPre: true, IsEditable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeSingleChar, Option: ""},
		&metadata.Attribute{ObjectID: objID, PropertyID: "group", PropertyName: "文件分组", IsRequired: false, IsOnly: false, IsPre: true, IsEditable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeSingleChar, Option: ""},
		&metadata.Attribute{ObjectID: objID, PropertyID: "format", PropertyName: "输出格式", IsRequired: false, IsOnly: false, IsPre: true, IsEditable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeEnum, Option: formatOption},
		&metadata.Attribute{ObjectID: objID, PropertyID: "right", PropertyName: "配置文件权限", IsRequired: false, IsOnly: false, IsPre: true, IsEditable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeEnum, Option: rightOption},
	}
	return dataRows

}

// TempVersionRow temp version structure
func TempVersionRow() []*metadata.Attribute {
	objID := common.BKInnerObjIDTempVersion
	statusOption := []validator.EnumVal{{ID: "draft", Name: "草稿", Type: "text", IsDefault: true}, {ID: "online", Name: "上线", Type: "text"}, {ID: "history", Name: "历史", Type: "text"}}

	dataRows := []*metadata.Attribute{
		//base info
		&metadata.Attribute{ObjectID: objID, PropertyID: "template_id", PropertyName: "模板ID", IsAPI: true, IsRequired: true, IsOnly: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeInt, Option: validator.IntOption{}},
		&metadata.Attribute{ObjectID: objID, PropertyID: common.BKAppIDField, PropertyName: "业务ID", IsAPI: true, IsRequired: true, IsOnly: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeInt, Option: validator.IntOption{}},
		&metadata.Attribute{ObjectID: objID, PropertyID: "description", PropertyName: "描述", IsRequired: false, IsOnly: true, IsPre: true, IsEditable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeSingleChar, Option: ""},
		&metadata.Attribute{ObjectID: objID, PropertyID: "operator", PropertyName: "操作人", IsRequired: false, IsOnly: false, IsPre: true, IsEditable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeSingleChar, Option: ""},
		&metadata.Attribute{ObjectID: objID, PropertyID: "create_time", PropertyName: "创建时间", IsRequired: false, IsOnly: false, IsPre: true, IsEditable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeSingleChar, Option: ""},
		&metadata.Attribute{ObjectID: objID, PropertyID: "content", PropertyName: "内容", IsRequired: true, IsOnly: false, IsPre: true, IsEditable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeLongChar, Option: ""},
		&metadata.Attribute{ObjectID: objID, PropertyID: "status", PropertyName: "状态", IsAPI: true, IsRequired: false, IsOnly: false, IsPre: true, IsEditable: true, PropertyGroup: groupBaseInfo, PropertyType: common.FieldTypeEnum, Option: statusOption},
	}
	return dataRows
}
