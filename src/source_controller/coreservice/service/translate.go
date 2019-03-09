/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.,
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the ",License",); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an ",AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package service

import (
	"configcenter/src/common"
	"configcenter/src/common/language"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/validator" //TODO: need to be removed
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

func (s *coreService) TranslateObjectName(defLang language.DefaultCCLanguageIf, obj *metadata.Object) string {
	return util.FirstNotEmptyString(defLang.Language("object_"+obj.ObjectID), obj.ObjectName, obj.ObjectID)
}
func (s *coreService) TranslateInstName(defLang language.DefaultCCLanguageIf, obj *metadata.Object) string {
	return util.FirstNotEmptyString(defLang.Language("inst_"+obj.ObjectID), obj.ObjectName, obj.ObjectID)
}

func (s *coreService) TranslatePropertyName(defLang language.DefaultCCLanguageIf, att *metadata.Attribute) string {
	return util.FirstNotEmptyString(defLang.Language(att.ObjectID+"_property_"+att.PropertyID), att.PropertyName, att.PropertyID)
}

func (s *coreService) TranslateDescription(defLang language.DefaultCCLanguageIf, att *metadata.Attribute) string {
	return util.FirstNotEmptyString(defLang.Language(att.ObjectID+"_description_"+att.PropertyID), att.Description)
}

func (s *coreService) TranslateEnumName(defLang language.DefaultCCLanguageIf, att *metadata.Attribute, val interface{}) interface{} {
	options := validator.ParseEnumOption(val)
	for index := range options {
		options[index].Name = util.FirstNotEmptyString(defLang.Language(att.ObjectID+"_property_"+att.PropertyID+"_enum_"+options[index].ID), options[index].Name, options[index].ID)
	}
	return options
}

func (s *coreService) TranslatePropertyGroupName(defLang language.DefaultCCLanguageIf, att *metadata.Group) string {
	return util.FirstNotEmptyString(defLang.Language(att.ObjectID+"_property_group_"+att.GroupID), att.GroupName, att.GroupID)
}

func (s *coreService) TranslateClassificationName(defLang language.DefaultCCLanguageIf, att *metadata.Classification) string {
	return util.FirstNotEmptyString(defLang.Language("classification_"+att.ClassificationID), att.ClassificationName, att.ClassificationID)
}
