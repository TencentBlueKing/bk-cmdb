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

package service

import (
	"context"
	"fmt"
	"reflect"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/language"
	meta "configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/validator"
	"configcenter/src/storage/dal"
)

var defaultNameLanguagePkg = map[string]map[string][]string{
	common.BKInnerObjIDModule: {
		"1": {"inst_module_idle", common.BKModuleNameField, common.BKModuleIDField},
		"2": {"inst_module_fault", common.BKModuleNameField, common.BKModuleIDField},
	},
	common.BKInnerObjIDApp: {
		"1": {"inst_biz_default", common.BKAppNameField, common.BKAppIDField},
	},
	common.BKInnerObjIDSet: {
		"1": {"inst_set_default", common.BKSetNameField, common.BKSetIDField},
	},
}

//GetCntByCondition get count by condition
func (cli *Service) GetCntByCondition(ctx context.Context, db dal.RDB, objType string, condition interface{}) (uint64, error) {
	tName := common.GetInstTableName(objType)
	cnt, err := db.Table(tName).Find(condition).Count(ctx)
	if nil != err {
		return 0, err
	}
	return cnt, nil
}

// GetHostByCondition query
func (cli *Service) GetHostByCondition(ctx context.Context, db dal.RDB, fields []string, condition, result interface{}, sort string, skip, limit int) error {
	return db.Table(common.BKTableNameBaseHost).Find(condition).Limit(uint64(limit)).Start(uint64(skip)).Sort(sort).All(ctx, result)
}

//GetObjectByCondition get object by condition
func (cli *Service) GetObjectByCondition(ctx context.Context, db dal.RDB, defLang language.DefaultCCLanguageIf, objType string, fields []string, condition, result interface{}, sort string, skip, limit int) error {
	tName := common.GetInstTableName(objType)
	if err := db.Table(tName).Find(condition).Fields(fields...).Limit(uint64(limit)).Start(uint64(skip)).Sort(sort).All(ctx, result); err != nil {
		blog.Errorf("failed to query the inst , error info %s", err.Error())
		return err
	}

	// translate language for default name
	if m, ok := defaultNameLanguagePkg[objType]; nil != defLang && ok {
		switch result.(type) {
		case *[]map[string]interface{}:
			results := *result.(*[]map[string]interface{})
			for index := range results {
				l := m[fmt.Sprint(results[index]["default"])]
				if len(l) >= 3 {
					results[index][l[1]] = util.FirstNotEmptyString(defLang.Language(l[0]), fmt.Sprint(results[index][l[1]]), fmt.Sprint(results[index][l[2]]))
				}
			}
		default:
			blog.Infof("GetObjectByCondition translate error: %v", reflect.TypeOf(result))
		}
	}

	return nil
}

func (cli *Service) DelObjByCondition(ctx context.Context, db dal.RDB, objType string, condition interface{}) error {
	tName := common.GetInstTableName(objType)
	err := db.Table(tName).Delete(ctx, condition)
	if nil != err {
		return err
	}
	return nil
}

//UpdateObjByCondition update object by condition
func (cli *Service) UpdateObjByCondition(ctx context.Context, db dal.RDB, objType string, data interface{}, condition interface{}) error {
	tName := common.GetInstTableName(objType)
	err := db.Table(tName).Update(ctx, condition, data)
	if nil != err {
		return err
	}
	return nil
}

//CreateObjectIntoDB add new object
func (cli *Service) CreateObjectIntoDB(ctx context.Context, db dal.RDB, objType string, input interface{}, idName *string) (int, error) {
	tName := common.GetInstTableName(objType)
	objID, err := db.NextSequence(ctx, tName)
	if err != nil {
		return 0, err
	}
	inputc := input.(map[string]interface{})
	*idName = common.GetInstIDField(objType)
	inputc[*idName] = objID
	err = db.Table(tName).Insert(ctx, inputc)
	return int(objID), err
}

//GetObjectByID get object by id
func (cli *Service) GetObjectByID(ctx context.Context, db dal.RDB, objType string, fields []string, id int, result interface{}, sort string) error {
	tName := common.GetInstTableName(objType)
	condition := make(map[string]interface{}, 1)
	condition[common.GetInstIDField(objType)] = id
	if tName == common.BKTableNameBaseInst && objType != common.BKInnerObjIDObject {
		condition[common.BKObjIDField] = objType
	}
	err := db.Table(tName).Find(condition).Fields(fields...).One(ctx, result)
	return err
}

func (cli *Service) TranslateObjectName(defLang language.DefaultCCLanguageIf, obj *meta.Object) string {
	return util.FirstNotEmptyString(defLang.Language("object_"+obj.ObjectID), obj.ObjectName, obj.ObjectID)
}
func (cli *Service) TranslateInstName(defLang language.DefaultCCLanguageIf, obj *meta.Object) string {
	return util.FirstNotEmptyString(defLang.Language("inst_"+obj.ObjectID), obj.ObjectName, obj.ObjectID)
}

func (cli *Service) TranslatePropertyName(defLang language.DefaultCCLanguageIf, att *meta.Attribute) string {
	return util.FirstNotEmptyString(defLang.Language(att.ObjectID+"_property_"+att.PropertyID), att.PropertyName, att.PropertyID)
}

func (cli *Service) TranslatePlaceholder(defLang language.DefaultCCLanguageIf, att *meta.Attribute) string {
	return util.FirstNotEmptyString(defLang.Language(att.ObjectID+"_placeholder_"+att.PropertyID), att.Placeholder)
}

func (cli *Service) TranslateEnumName(defLang language.DefaultCCLanguageIf, att *meta.Attribute, val interface{}) interface{} {
	options := validator.ParseEnumOption(val)
	for index := range options {
		options[index].Name = util.FirstNotEmptyString(defLang.Language(att.ObjectID+"_property_"+att.PropertyID+"_enum_"+options[index].ID), options[index].Name, options[index].ID)
	}
	return options
}

func (cli *Service) TranslatePropertyGroupName(defLang language.DefaultCCLanguageIf, att *meta.Group) string {
	return util.FirstNotEmptyString(defLang.Language(att.ObjectID+"_property_group_"+att.GroupID), att.GroupName, att.GroupID)
}

func (cli *Service) TranslateClassificationName(defLang language.DefaultCCLanguageIf, att *meta.Classification) string {
	return util.FirstNotEmptyString(defLang.Language("classification_"+att.ClassificationID), att.ClassificationName, att.ClassificationID)
}

func (cli *Service) TranslateAssociationKind(defLang language.DefaultCCLanguageIf, kind *meta.AssociationKind) {
	kind.AssociationKindName = util.FirstNotEmptyString(defLang.Language("unique_kind_name_"+kind.AssociationKindID), kind.AssociationKindName)
	kind.SourceToDestinationNote = util.FirstNotEmptyString(defLang.Language("unique_kind_src_to_dest_"+kind.AssociationKindID), kind.SourceToDestinationNote)
	kind.DestinationToSourceNote = util.FirstNotEmptyString(defLang.Language("unique_kind_dest_to_src_"+kind.AssociationKindID), kind.DestinationToSourceNote)
}
