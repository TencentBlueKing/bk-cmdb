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

package v3v0v8

import (
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/scene_server/validator"
	"configcenter/src/storage"
)

func addDefaultBiz(db storage.DI, conf *upgrader.Config) error {
	// add default biz
	defaultBiz := map[string]interface{}{}
	defaultBiz[common.BKAppNameField] = common.DefaultAppName
	defaultBiz[common.BKMaintainersField] = "admin"
	defaultBiz[common.BKProductPMField] = "admin"
	defaultBiz[common.BKTimeZoneField] = "Asia/Shanghai"
	defaultBiz[common.BKLanguageField] = "1" //中文
	defaultBiz[common.BKLifeCycleField] = common.DefaultAppLifeCycleNormal
	defaultBiz[common.BKOwnerIDField] = conf.OwnerID
	defaultBiz[common.BKSupplierIDField] = common.BKDefaultSupplierID
	defaultBiz[common.BKDefaultField] = common.DefaultAppFlag
	defaultBiz[common.CreateTimeField] = time.Now()
	defaultBiz[common.LastTimeField] = time.Now()
	filled := fillEmptyFields(defaultBiz, AppRow())
	bizID, _, err := upgrader.Upsert(db, "cc_ApplicationBase", defaultBiz, common.BKAppIDField, []string{common.BKOwnerIDField, common.BKAppNameField, common.BKDefaultField}, append(filled, common.BKAppIDField))
	if err != nil {
		blog.Error("add defaultBiz error ", err.Error())
		return err
	}

	// add default set
	defaultSet := make(map[string]interface{})
	defaultSet[common.BKAppIDField] = bizID
	defaultSet[common.BKInstParentStr] = bizID
	defaultSet[common.BKSetNameField] = common.DefaultResSetName
	defaultSet[common.BKDefaultField] = common.DefaultResSetFlag
	defaultSet[common.BKOwnerIDField] = conf.OwnerID
	defaultSet[common.CreateTimeField] = time.Now()
	defaultSet[common.LastTimeField] = time.Now()
	filled = fillEmptyFields(defaultSet, SetRow())
	setID, _, err := upgrader.Upsert(db, "cc_SetBase", defaultSet, common.BKSetIDField, []string{common.BKOwnerIDField, common.BKSetNameField, common.BKAppIDField, common.BKDefaultField}, append(filled, common.BKSetIDField))
	if err != nil {
		blog.Error("add defaultSet error ", err.Error())
		return err
	}

	// add default module
	defaultResModule := make(map[string]interface{})
	defaultResModule[common.BKSetIDField] = setID
	defaultResModule[common.BKInstParentStr] = setID
	defaultResModule[common.BKAppIDField] = bizID
	defaultResModule[common.BKModuleNameField] = common.DefaultResModuleName
	defaultResModule[common.BKDefaultField] = common.DefaultResModuleFlag
	defaultResModule[common.BKOwnerIDField] = conf.OwnerID
	defaultResModule[common.CreateTimeField] = time.Now()
	defaultResModule[common.LastTimeField] = time.Now()
	filled = fillEmptyFields(defaultResModule, ModuleRow())
	_, _, err = upgrader.Upsert(db, "cc_ModuleBase", defaultResModule, common.BKModuleIDField, []string{common.BKOwnerIDField, common.BKModuleNameField, common.BKAppIDField, common.BKSetIDField, common.BKDefaultField}, append(filled, common.BKModuleIDField))
	if err != nil {
		blog.Error("add defaultResModule error ", err.Error())
		return err
	}
	defaultFaultModule := make(map[string]interface{})
	defaultFaultModule[common.BKSetIDField] = setID
	defaultFaultModule[common.BKInstParentStr] = setID
	defaultFaultModule[common.BKAppIDField] = bizID
	defaultFaultModule[common.BKModuleNameField] = common.DefaultFaultModuleName
	defaultFaultModule[common.BKDefaultField] = common.DefaultFaultModuleFlag
	defaultFaultModule[common.BKOwnerIDField] = conf.OwnerID
	defaultFaultModule[common.CreateTimeField] = time.Now()
	defaultFaultModule[common.LastTimeField] = time.Now()
	filled = fillEmptyFields(defaultFaultModule, ModuleRow())
	_, _, err = upgrader.Upsert(db, "cc_ModuleBase", defaultFaultModule, common.BKModuleIDField, []string{common.BKOwnerIDField, common.BKModuleNameField, common.BKAppIDField, common.BKSetIDField, common.BKDefaultField}, append(filled, common.BKModuleIDField))
	if err != nil {
		blog.Error("add defaultFaultModule error ", err.Error())
		return err
	}

	return nil
}

func fillEmptyFields(data map[string]interface{}, rows []*metadata.Attribute) []string {
	filled := []string{}
	for _, field := range rows {
		fieldName := field.PropertyID
		fieldType := field.PropertyType
		if _, ok := data[fieldName]; ok {
			continue
		}
		option := field.Option
		switch fieldType {
		case common.FieldTypeSingleChar:
			data[fieldName] = ""
		case common.FieldTypeLongChar:
			data[fieldName] = ""
		case common.FieldTypeInt:
			data[fieldName] = nil
		case common.FieldTypeEnum:
			enumOptions := validator.ParseEnumOption(option)
			v := ""
			if len(enumOptions) > 0 {
				var defaultOption *validator.EnumVal
				for _, k := range enumOptions {
					if k.IsDefault {
						defaultOption = &k
						break
					}
				}
				if nil != defaultOption {
					v = defaultOption.ID
				}
			}
			data[fieldName] = v
		case common.FieldTypeDate:
			data[fieldName] = ""
		case common.FieldTypeTime:
			data[fieldName] = ""
		case common.FieldTypeUser:
			data[fieldName] = ""
		case common.FieldTypeMultiAsst:
			data[fieldName] = nil
		case common.FieldTypeTimeZone:
			data[fieldName] = nil
		case common.FieldTypeBool:
			data[fieldName] = false
		default:
			data[fieldName] = nil
		}
		filled = append(filled, fieldName)
	}
	filled = append(filled, common.CreateTimeField)
	data[common.CreateTimeField] = time.Now()
	data[common.LastTimeField] = time.Now()
	return filled
}
