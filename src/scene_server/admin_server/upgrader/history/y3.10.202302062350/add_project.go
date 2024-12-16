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

package y3_10_202302062350

import (
	"context"
	"errors"
	"fmt"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	mCommon "configcenter/src/scene_server/admin_server/common"
	"configcenter/src/scene_server/admin_server/upgrader/history"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/dal/types"

	"go.mongodb.org/mongo-driver/bson"
)

var dataRows = map[string]*attribute{
	common.BKFieldID: {
		ObjectID:      common.BKInnerObjIDProject,
		PropertyID:    common.BKFieldID,
		PropertyName:  "ID",
		IsAPI:         true,
		IsRequired:    false,
		IsOnly:        false,
		PropertyGroup: mCommon.BaseInfo,
		PropertyType:  common.FieldTypeInt,
		Creator:       common.CCSystemOperatorUserName,
		Option:        metadata.PrevIntOption{},
	},
	common.BKProjectIDField: {
		ObjectID:      common.BKInnerObjIDProject,
		PropertyID:    common.BKProjectIDField,
		PropertyName:  "项目ID",
		IsAPI:         true,
		IsRequired:    false,
		IsOnly:        false,
		PropertyGroup: mCommon.BaseInfo,
		PropertyType:  common.FieldTypeSingleChar,
		Creator:       common.CCSystemOperatorUserName,
		Option:        "",
	},
	common.BKProjectNameField: {
		ObjectID:      common.BKInnerObjIDProject,
		PropertyID:    common.BKProjectNameField,
		PropertyName:  "项目名称",
		IsRequired:    true,
		IsOnly:        true,
		IsEditable:    true,
		PropertyGroup: mCommon.BaseInfo,
		PropertyType:  common.FieldTypeSingleChar,
		Creator:       common.CCSystemOperatorUserName,
		Option:        "",
	},
	common.BKProjectCodeField: {
		ObjectID:      common.BKInnerObjIDProject,
		PropertyID:    common.BKProjectCodeField,
		PropertyName:  "项目英文名",
		IsRequired:    true,
		IsOnly:        true,
		PropertyGroup: mCommon.BaseInfo,
		PropertyType:  common.FieldTypeSingleChar,
		Creator:       common.CCSystemOperatorUserName,
		Option:        "",
	},
	common.BKProjectDescField: {
		ObjectID:      common.BKInnerObjIDProject,
		PropertyID:    common.BKProjectDescField,
		PropertyName:  "项目描述",
		IsRequired:    false,
		IsOnly:        false,
		IsEditable:    true,
		PropertyGroup: mCommon.BaseInfo,
		PropertyType:  common.FieldTypeLongChar,
		Creator:       common.CCSystemOperatorUserName,
		Option:        "",
	},
	common.BKProjectTypeField: {
		ObjectID:      common.BKInnerObjIDProject,
		PropertyID:    common.BKProjectTypeField,
		PropertyName:  "项目类型",
		IsRequired:    true,
		IsOnly:        false,
		IsEditable:    true,
		PropertyGroup: mCommon.BaseInfo,
		PropertyType:  common.FieldTypeEnum,
		Creator:       common.CCSystemOperatorUserName,
		Option: []metadata.EnumVal{
			{ID: "mobile_game", Name: "手游", Type: "text"},
			{ID: "pc_game", Name: "端游", Type: "text"},
			{ID: "web_game", Name: "页游", Type: "text"},
			{ID: "platform_prod", Name: "平台产品", Type: "text"},
			{ID: "support_prod", Name: "支撑产品", Type: "text"},
			{ID: "other", Name: "其他", Type: "text", IsDefault: true},
		},
	},
	common.BKProjectSecLvlField: {
		ObjectID:      common.BKInnerObjIDProject,
		PropertyID:    common.BKProjectSecLvlField,
		PropertyName:  "保密级别",
		IsRequired:    false,
		IsOnly:        false,
		IsEditable:    true,
		PropertyGroup: mCommon.BaseInfo,
		PropertyType:  common.FieldTypeEnum,
		Creator:       common.CCSystemOperatorUserName,
		Option: []metadata.EnumVal{
			{ID: "public", Name: "公开", Type: "text", IsDefault: true},
			{ID: "private", Name: "私有", Type: "text"},
			{ID: "classified", Name: "机密", Type: "text"},
		},
	},
	common.BKProjectOwnerField: {
		ObjectID:      common.BKInnerObjIDProject,
		PropertyID:    common.BKProjectOwnerField,
		PropertyName:  "项目负责人",
		IsRequired:    true,
		IsOnly:        false,
		IsEditable:    true,
		IsMultiple:    true,
		PropertyGroup: mCommon.BaseInfo,
		PropertyType:  common.FieldTypeUser,
		Creator:       common.CCSystemOperatorUserName,
		Option:        "",
	},
	common.BKProjectTeamField: {
		ObjectID:      common.BKInnerObjIDProject,
		PropertyID:    common.BKProjectTeamField,
		PropertyName:  "所属团队",
		IsRequired:    false,
		IsOnly:        false,
		IsEditable:    true,
		IsMultiple:    true,
		PropertyGroup: mCommon.BaseInfo,
		PropertyType:  common.FieldTypeOrganization,
		Creator:       common.CCSystemOperatorUserName,
		Option:        "",
	},
	common.BKProjectStatusField: {
		ObjectID:      common.BKInnerObjIDProject,
		PropertyID:    common.BKProjectStatusField,
		PropertyName:  "项目状态",
		IsRequired:    false,
		IsOnly:        false,
		IsEditable:    true,
		PropertyGroup: mCommon.BaseInfo,
		PropertyType:  common.FieldTypeEnum,
		Creator:       common.CCSystemOperatorUserName,
		Option: []metadata.EnumVal{
			{ID: "enable", Name: "启用", Type: "text", IsDefault: true},
			{ID: "disabled", Name: "未启用", Type: "text"},
		},
	},
	common.BKProjectIconField: {
		ObjectID:      common.BKInnerObjIDProject,
		PropertyID:    common.BKProjectIconField,
		PropertyName:  "项目图标",
		IsRequired:    false,
		IsOnly:        false,
		IsEditable:    true,
		PropertyGroup: mCommon.BaseInfo,
		PropertyType:  common.FieldTypeLongChar,
		Creator:       common.CCSystemOperatorUserName,
		Option:        "",
	},
}

func addProjectObjectRow(ctx context.Context, db dal.RDB, ownerID string) error {
	filter := mapstr.MapStr{common.BKObjIDField: common.BKInnerObjIDProject}
	model := new(metadata.Object)
	err := db.Table(common.BKTableNameObjDes).Find(filter).
		Fields(common.BKFieldID, common.BKObjIDField, common.BKObjNameField, common.CreatorField).One(ctx, model)
	if err != nil && !db.IsNotFoundError(err) {
		blog.Errorf("count project object failed, err: %v", err)
		return err
	}
	if model.ID != 0 {
		if model.ObjectID == common.BKInnerObjIDProject && model.ObjectName == "项目" &&
			model.Creator == common.CCSystemOperatorUserName {
			return nil
		}
		blog.Errorf("the model %s already exists, but does not conform to the specification, object name: %s, "+
			"creator: %s", common.BKInnerObjIDProject, model.ObjectName, model.Creator)
		return fmt.Errorf("model %s failed to create", common.BKInnerObjIDProject)
	}

	t := metadata.Now()
	dataRows := Object{
		ObjCls:      metadata.ClassificationOrganizationID,
		ObjectID:    common.BKInnerObjIDProject,
		ObjectName:  "项目",
		IsPre:       true,
		ObjIcon:     "icon-cc-project",
		CreateTime:  &t,
		LastTime:    &t,
		IsPaused:    false,
		Creator:     common.CCSystemOperatorUserName,
		OwnerID:     ownerID,
		Description: "",
		Modifier:    "",
	}
	uniqueKeys := []string{common.BKObjIDField}
	_, _, err = history.Upsert(ctx, db, common.BKTableNameObjDes, dataRows, "id", uniqueKeys, []string{"id"})
	if err != nil {
		blog.Errorf("add data for %s table failed, err: %v", common.BKTableNameObjDes, err)
		return err
	}
	return nil
}

func addObjectUnique(ctx context.Context, db dal.RDB, conf *history.Config) error {
	attrs := make([]metadata.Attribute, 0)
	cond := mapstr.MapStr{
		common.BKObjIDField: common.BKInnerObjIDProject,
		common.BKPropertyIDField: mapstr.MapStr{
			common.BKDBIN: []string{common.BKFieldID, common.BKProjectIDField, common.BKProjectNameField,
				common.BKProjectCodeField},
		},
	}
	if err := db.Table(common.BKTableNameObjAttDes).Find(cond).All(ctx, &attrs); err != nil {
		return err
	}

	uniqueIdxs := make([]metadata.ObjectUnique, 0)
	condObjUnique := mapstr.MapStr{common.BKObjIDField: common.BKInnerObjIDProject}
	if err := db.Table(common.BKTableNameObjUnique).Find(condObjUnique).Fields(common.BKObjectUniqueKeys).
		All(ctx, &uniqueIdxs); err != nil {
		return err
	}

	// if it exists, you need to determine whether the id is legal.
	if len(uniqueIdxs) > 0 {
		if len(uniqueIdxs) != len(attrs) {
			return errors.New("inconsistent number of unique indexes")
		}

		keysId := make(map[uint64]struct{})
		for _, index := range uniqueIdxs {
			for _, id := range index.Keys {
				keysId[id.ID] = struct{}{}
			}
		}
		// to prevent compliance with index scenarios, the number of indexes must be exactly the same.
		if len(keysId) != len(attrs) {
			return errors.New("invalid number of index ids")
		}
		for _, attr := range attrs {
			if _, ok := keysId[uint64(attr.ID)]; !ok {
				return fmt.Errorf("attr id: %d does not exist", attr.ID)
			}
		}
		return nil
	}

	for _, attr := range attrs {
		keys := []metadata.UniqueKey{
			{
				Kind: metadata.UniqueKeyKindProperty,
				ID:   uint64(attr.ID),
			},
		}
		unique := ObjectUnique{
			ObjID:    common.BKInnerObjIDProject,
			Keys:     keys,
			Ispre:    true,
			OwnerID:  conf.TenantID,
			LastTime: metadata.Now(),
		}

		uid, err := db.NextSequence(ctx, common.BKTableNameObjUnique)
		if err != nil {
			return err
		}
		unique.ID = uid

		if err := db.Table(common.BKTableNameObjUnique).Insert(ctx, unique); err != nil {
			return err
		}
	}

	return nil
}

func addProjectCollection(ctx context.Context, db dal.RDB) error {
	exists, err := db.HasTable(ctx, common.BKTableNameBaseProject)
	if err != nil {
		blog.Errorf("check if table %s exists failed, err: %v", common.BKTableNameBaseProject, err)
		return err
	}

	if !exists {
		if err := db.CreateTable(ctx, common.BKTableNameBaseProject); err != nil {
			return err
		}
		return nil
	}
	blog.Infof("project collection has been created")
	return nil
}

func addProjectPropertyGroup(ctx context.Context, db dal.RDB, ownerID string) error {
	rows := []*Group{
		{
			ObjectID:   common.BKInnerObjIDProject,
			GroupID:    mCommon.BaseInfo,
			GroupName:  mCommon.BaseInfoName,
			GroupIndex: 1,
			OwnerID:    ownerID,
			IsDefault:  true,
		},
	}

	for _, row := range rows {
		if _, _, err := history.Upsert(ctx, db, common.BKTableNamePropertyGroup, row, "id",
			[]string{common.BKObjIDField, common.BKPropertyGroupIDField}, []string{"id"}); err != nil {
			blog.Errorf("add data for %s table failed, err: %v", common.BKTableNamePropertyGroup, err)
			return err
		}
	}

	return nil
}

func addProjectObjectAttrRow(ctx context.Context, db dal.RDB, ownerID string) error {
	filter := mapstr.MapStr{common.BKObjIDField: common.BKInnerObjIDProject}
	attrs := make([]metadata.Attribute, 0)
	// 判断是否有project的对象属性表，如果没有需要初始化
	if err := db.Table(common.BKTableNameObjAttDes).Find(filter).Fields(common.BKPropertyIDField,
		common.BKPropertyNameField, common.CreatorField).All(ctx, &attrs); err != nil && !db.IsNotFoundError(err) {
		blog.Errorf("find object attribute describe failed, err: %v", err)
		return err
	}

	if len(attrs) > 0 {
		// 如果存在的话，数量必须一致。并且必须严格校验每个属性bk_property_name和creator必须完全一致，不一致直接报错需要先处理完毕后再升级
		if len(attrs) != len(dataRows) {
			blog.Errorf("Illegal number of project model attributes, num is: %d", len(attrs))
			return errors.New("illegal number of project model attributes")
		}

		for _, attr := range attrs {
			if data, ok := dataRows[attr.PropertyID]; ok {
				if attr.PropertyName != data.PropertyName || attr.Creator != data.Creator {
					blog.Errorf("the model project attribute %s already exists, but is illegal, name: %v, creator: %v",
						attr.PropertyID, attr.PropertyName, attr.Creator)
					return fmt.Errorf("model project attribute %s is invalid", attr.PropertyID)
				}
			}
		}
		return nil
	}

	uniqueFields := []string{common.BKObjIDField, common.BKPropertyIDField, "bk_supplier_account"}
	nowTime := metadata.Now()
	for _, row := range dataRows {
		row.OwnerID = ownerID
		row.IsPre = true
		row.IsReadOnly = false
		row.CreateTime = &nowTime
		row.LastTime = &nowTime
		row.Description = ""
		_, _, err := history.Upsert(ctx, db, common.BKTableNameObjAttDes, row, "id", uniqueFields, []string{})
		if err != nil {
			blog.Errorf("add project attr failed, attribute: %v, err: %v", row, err)
			return err
		}
	}
	return nil
}

func addProjectTableIndexes(ctx context.Context, db dal.RDB) error {
	indexes := []types.Index{
		{
			Name:       common.CCLogicUniqueIdxNamePrefix + common.BKFieldID,
			Keys:       bson.D{{common.BKFieldID, 1}},
			Background: true,
			Unique:     true,
		},
		{
			Name:       common.CCLogicUniqueIdxNamePrefix + common.BKProjectIDField,
			Keys:       bson.D{{common.BKProjectIDField, 1}},
			Background: true,
			Unique:     true,
		},
		{
			Name:       common.CCLogicUniqueIdxNamePrefix + common.BKProjectNameField,
			Keys:       bson.D{{common.BKProjectNameField, 1}},
			Unique:     true,
			Background: true,
		},
		{
			Name:       common.CCLogicUniqueIdxNamePrefix + common.BKProjectCodeField,
			Keys:       bson.D{{common.BKProjectCodeField, 1}},
			Unique:     true,
			Background: true,
		},
		{
			Name:       common.CCLogicIndexNamePrefix + common.BKProjectStatusField,
			Keys:       bson.D{{common.BKProjectStatusField, 1}},
			Background: true,
		},
	}

	existIndexArr, err := db.Table(common.BKTableNameBaseProject).Indexes(ctx)
	if err != nil {
		blog.Errorf("get exist index for project table failed, err: %v", err)
		return err
	}

	existIdxMap := make(map[string]bool)
	for _, index := range existIndexArr {
		// skip the default "_id" index for the database
		if index.Name == "_id_" {
			continue
		}
		existIdxMap[index.Name] = true
	}

	needAddIndexes := make([]types.Index, 0)
	for _, index := range indexes {
		if _, exist := existIdxMap[index.Name]; exist {
			continue
		}
		needAddIndexes = append(needAddIndexes, index)
	}

	if len(needAddIndexes) != 0 {
		err = db.Table(common.BKTableNameBaseProject).BatchCreateIndexes(ctx, needAddIndexes)
		if err != nil && !db.IsDuplicatedError(err) {
			blog.Errorf("create index for project table failed, err: %v, index: %+v", err, needAddIndexes)
			return err
		}
	}
	return nil
}

func addProjectPropertyOption(ctx context.Context, db dal.RDB, conf *history.Config) error {
	if err := addProjectObjectRow(ctx, db, conf.TenantID); err != nil {
		return err
	}

	if err := addProjectPropertyGroup(ctx, db, conf.TenantID); err != nil {
		return err
	}

	if err := addProjectObjectAttrRow(ctx, db, conf.TenantID); err != nil {
		return err
	}

	if err := addObjectUnique(ctx, db, conf); err != nil {
		return err
	}

	if err := addProjectCollection(ctx, db); err != nil {
		return err
	}

	if err := addProjectTableIndexes(ctx, db); err != nil {
		return err
	}

	return nil
}

type attribute struct {
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
	IsMultiple        bool           `field:"ismultiple" json:"ismultiple" bson:"ismultiple"`
	Description       string         `field:"description" json:"description" bson:"description"`
	Creator           string         `field:"creator" json:"creator" bson:"creator"`
	CreateTime        *metadata.Time `json:"create_time" bson:"create_time"`
	LastTime          *metadata.Time `json:"last_time" bson:"last_time"`
}

// Object object metadata definition
type Object struct {
	ID         int64  `field:"id" json:"id" bson:"id" mapstructure:"id"`
	ObjCls     string `field:"bk_classification_id" json:"bk_classification_id" bson:"bk_classification_id" mapstructure:"bk_classification_id"`
	ObjIcon    string `field:"bk_obj_icon" json:"bk_obj_icon" bson:"bk_obj_icon" mapstructure:"bk_obj_icon"`
	ObjectID   string `field:"bk_obj_id" json:"bk_obj_id" bson:"bk_obj_id" mapstructure:"bk_obj_id"`
	ObjectName string `field:"bk_obj_name" json:"bk_obj_name" bson:"bk_obj_name" mapstructure:"bk_obj_name"`

	// IsHidden front-end don't display the object if IsHidden is true
	IsHidden bool `field:"bk_ishidden" json:"bk_ishidden" bson:"bk_ishidden" mapstructure:"bk_ishidden"`

	IsPre         bool           `field:"ispre" json:"ispre" bson:"ispre" mapstructure:"ispre"`
	IsPaused      bool           `field:"bk_ispaused" json:"bk_ispaused" bson:"bk_ispaused" mapstructure:"bk_ispaused"`
	Position      string         `field:"position" json:"position" bson:"position" mapstructure:"position"`
	OwnerID       string         `field:"bk_supplier_account" json:"bk_supplier_account" bson:"bk_supplier_account" mapstructure:"bk_supplier_account"`
	Description   string         `field:"description" json:"description" bson:"description" mapstructure:"description"`
	Creator       string         `field:"creator" json:"creator" bson:"creator" mapstructure:"creator"`
	Modifier      string         `field:"modifier" json:"modifier" bson:"modifier" mapstructure:"modifier"`
	CreateTime    *metadata.Time `field:"create_time" json:"create_time" bson:"create_time" mapstructure:"create_time"`
	LastTime      *metadata.Time `field:"last_time" json:"last_time" bson:"last_time" mapstructure:"last_time"`
	ObjSortNumber int64          `field:"obj_sort_number" json:"obj_sort_number" bson:"obj_sort_number" mapstructure:"obj_sort_number"`
}

// ObjectUnique TODO
type ObjectUnique struct {
	ID         uint64               `json:"id" bson:"id"`
	TemplateID int64                `json:"bk_template_id" bson:"bk_template_id"`
	ObjID      string               `json:"bk_obj_id" bson:"bk_obj_id"`
	Keys       []metadata.UniqueKey `json:"keys" bson:"keys"`
	Ispre      bool                 `json:"ispre" bson:"ispre"`
	OwnerID    string               `json:"bk_supplier_account" bson:"bk_supplier_account"`
	LastTime   metadata.Time        `json:"last_time" bson:"last_time"`
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
