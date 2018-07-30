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
	"fmt"
	"reflect"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/language"
	meta "configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/validator"
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
func (cli *Service) GetCntByCondition(objType string, condition interface{}) (int, error) {
	tName := common.GetInstTableName(objType)
	cnt, err := cli.Instance.GetCntByCondition(tName, condition)
	if nil != err {
		return 0, err
	}
	return cnt, nil
}

// GetHostByCondition query
func (cli *Service) GetHostByCondition(fields []string, condition, result interface{}, sort string, skip, limit int) error {
	return cli.Instance.GetMutilByCondition("cc_HostBase", fields, condition, result, sort, skip, limit)
}

//GetObjectByCondition get object by condition
func (cli *Service) GetObjectByCondition(defLang language.DefaultCCLanguageIf, objType string, fields []string, condition, result interface{}, sort string, skip, limit int) error {
	tName := common.GetInstTableName(objType)
	if err := cli.Instance.GetMutilByCondition(tName, fields, condition, result, sort, skip, limit); err != nil {
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

func (cli *Service) DelObjByCondition(objType string, condition interface{}) error {
	tName := common.GetInstTableName(objType)
	err := cli.Instance.DelByCondition(tName, condition)
	if nil != err {
		return err
	}
	return nil
}

//UpdateObjByCondition update object by condition
func (cli *Service) UpdateObjByCondition(objType string, data interface{}, condition interface{}) error {
	tName := common.GetInstTableName(objType)
	err := cli.Instance.UpdateByCondition(tName, data, condition)
	if nil != err {
		return err
	}
	return nil
}

//CreateObjectIntoDB add new object
func (cli *Service) CreateObjectIntoDB(objType string, input interface{}, idName *string) (int, error) {
	tName := common.GetInstTableName(objType)
	objID, err := cli.Instance.GetIncID(tName)
	if err != nil {
		return 0, err
	}
	inputc := input.(map[string]interface{})
	*idName = common.GetInstIDField(objType)
	inputc[*idName] = objID
	cli.Instance.Insert(tName, inputc)
	return int(objID), nil
}

//GetObjectByID get object by id
func (cli *Service) GetObjectByID(objType string, fields []string, id int, result interface{}, sort string) error {
	tName := common.GetInstTableName(objType)
	condition := make(map[string]interface{}, 1)
	condition[common.GetInstIDField(objType)] = id
	if tName == common.BKTableNameBaseInst && objType != common.BKINnerObjIDObject {
		condition[common.BKObjIDField] = objType
	}
	err := cli.Instance.GetOneByCondition(tName, fields, condition, result)
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
