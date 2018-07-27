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

package logics

import (
	"fmt"
	"net/http"

	simplejson "github.com/bitly/go-simplejson"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	lang "configcenter/src/common/language"
	"configcenter/src/common/util"
	webCommon "configcenter/src/web_server/common"
)

// Property object fields
type Property struct {
	ID                     string
	Name                   string
	PropertyType           string
	Option                 interface{}
	IsPre                  bool
	IsRequire              bool
	Group                  string
	Index                  int
	ExcelColIndex          int
	NotObjPropery          bool //Not an attribute of the object, indicating that the field to be exported is needed for export,
	AsstObjPrimaryProperty []Property
	IsOnly                 bool
	AsstObjID              string
}

// PropertyGroup property group
type PropertyGroup struct {
	Name  string
	Index int
	ID    string
}

// GetObjFieldIDs get object fields
func GetObjFieldIDs(objID, url string, filterFields []string, header http.Header) (map[string]Property, error) {

	fields, err := getObjFieldIDs(objID, url, header)
	if nil != err {
		return nil, err
	}
	groups, err := getObjectGroup(objID, url, header)
	if nil != err {
		return nil, err
	}

	ret := make(map[string]Property)
	index := 0

	for _, group := range groups {
		for _, field := range fields {
			if field.Group == group.ID {
				if util.InStrArr(filterFields, field.ID) {
					continue
				}
				switch field.PropertyType {
				case common.FieldTypeSingleAsst:
					fallthrough
				case common.FieldTypeMultiAsst:

					field.AsstObjPrimaryProperty, err = getAsstObjectPrimaryFieldByObjID(field.AsstObjID, url, header)
					if nil != err {
						blog.Errorf("get associate object fields error: error:%s", err.Error())
						return nil, fmt.Errorf("get associate object fields error: error:%s", err.Error())
					}
				}
				field.ExcelColIndex = index
				ret[field.ID] = field
				index++
			}
		}
	}

	return ret, nil
}

func getObjectGroup(objID, url string, header http.Header) ([]PropertyGroup, error) {
	///api/v3/objectatt/group/property/owner/0/object/host
	url = fmt.Sprintf("%s/api/%s/objectatt/group/property/owner/%s/object/%s", url, webCommon.API_VERSION, util.GetActionOnwerIDByHTTPHeader(header), objID)
	conds := common.KvMap{common.BKObjIDField: objID, common.BKOwnerIDField: common.BKDefaultOwnerID, "page": common.KvMap{"start": 0, "limit": common.BKNoLimit, "sort": common.BKPropertyGroupIndexField}}
	result, err := httpRequest(url, conds, header)
	if nil != err {
		return nil, err
	}
	blog.Info("get %s fields group  url:%s", objID, url)
	blog.Info("get %s fields group return:%s", objID, result)
	js, err := simplejson.NewJson([]byte(result))
	if nil != err {
		blog.Errorf("get %s fields group  url:%s return:%s, err:%s", objID, url, result, err.Error())
		return nil, err
	}
	fields, err := js.Get("data").Array()
	if nil != err {
		blog.Errorf("get %s fields group  url:%s return:%s data not array, err:%s", objID, url, result, err.Error())
		return nil, err
	}
	ret := []PropertyGroup{}
	for _, field := range fields {
		mapField, _ := field.(map[string]interface{})
		propertyGroup := PropertyGroup{}
		propertyGroup.Index, _ = util.GetIntByInterface(mapField[common.BKPropertyGroupIndexField])
		propertyGroup.Name, _ = mapField[common.BKPropertyGroupNameField].(string)
		propertyGroup.ID, _ = mapField[common.BKPropertyGroupIDField].(string)
		ret = append(ret, propertyGroup)
	}
	blog.V(3).Infof("getObjectGroup count:%d", len(ret))
	return ret, nil

}

func getAsstObjectPrimaryFieldByObjID(objID string, url string, header http.Header) ([]Property, error) {

	fields, err := getObjFieldIDsBySort(objID, url, common.BKPropertyIDField, header)
	if nil != err {
		return nil, err
	}
	var ret []Property
	for _, field := range fields {
		if true == field.IsOnly {
			ret = append(ret, field)
		}
	}
	return ret, nil

}

func getObjFieldIDs(objID, url string, header http.Header) ([]Property, error) {
	sort := fmt.Sprintf("-%s,bk_property_index", common.BKIsRequiredField)
	return getObjFieldIDsBySort(objID, url, sort, header)

}

func getObjFieldIDsBySort(objID, url, sort string, header http.Header) ([]Property, error) {
	url = fmt.Sprintf("%s/api/%s/object/attr/search", url, webCommon.API_VERSION)
	conds := common.KvMap{
		common.BKObjIDField:   objID,
		common.BKOwnerIDField: common.BKDefaultOwnerID,
		"page": common.KvMap{
			"start": 0,
			"limit": common.BKNoLimit,
			"sort":  sort,
		},
	}
	result, err := httpRequest(url, conds, header)
	if nil != err {
		return nil, err
	}
	blog.Info("get %s fields  url:%s", objID, url)
	blog.Info("get %s fields return:%s", objID, result)
	js, err := simplejson.NewJson([]byte(result))
	if nil != err {
		blog.Errorf("get %s fields  url:%s return:%s", objID, url, result)
		return nil, err
	}
	fields, err := js.Get("data").Array()
	if nil != err {
		blog.Errorf("get %s fields  url:%s return:%s data not array, error:%s", objID, url, result, err.Error())
		return nil, err
	}
	ret := []Property{}

	for _, field := range fields {
		mapField, _ := field.(map[string]interface{})

		fieldIsOnly, ok := mapField[common.BKIsOnly].(bool)
		if false == ok {
			return nil, fmt.Errorf("%s not foud", common.BKIsOnly)
		}
		fieldName, _ := mapField[common.BKPropertyNameField].(string)
		fieldID, _ := mapField[common.BKPropertyIDField].(string)
		fieldType, _ := mapField[common.BKPropertyTypeField].(string)
		fieldIsRequire, _ := mapField[common.BKIsRequiredField].(bool)
		fieldIsOption, _ := mapField[common.BKOptionField]
		fieldIsPre, _ := mapField[common.BKIsPre].(bool)
		fieldGroup, _ := mapField[common.BKPropertyGroupField].(string)
		fieldIndex, _ := util.GetIntByInterface(mapField["bk_property_index"])
		fieldAsstObjID, _ := mapField[common.BKAsstObjIDField].(string)

		ret = append(ret, Property{
			ID:           fieldID,
			Name:         fieldName,
			PropertyType: fieldType,
			IsRequire:    fieldIsRequire,
			IsPre:        fieldIsPre,
			Option:       fieldIsOption,
			Group:        fieldGroup,
			Index:        fieldIndex,
			IsOnly:       fieldIsOnly,
			AsstObjID:    fieldAsstObjID,
		})
	}
	blog.V(3).Infof("getObjFieldIDsBySort ret count:%d", len(ret))

	return ret, nil
}

// getPropertyTypeAliasName  return propertyType name, whether to export,
func getPropertyTypeAliasName(propertyType string, defLang lang.DefaultCCLanguageIf) (string, bool) {
	var skip bool
	name := defLang.Language("field_type_" + propertyType)
	switch propertyType {
	case common.FieldTypeSingleChar:
	case common.FieldTypeLongChar:
	case common.FieldTypeInt:
	case common.FieldTypeEnum:
	case common.FieldTypeDate:
	case common.FieldTypeTime:
	case common.FieldTypeUser:
	case common.FieldTypeSingleAsst:
	case common.FieldTypeMultiAsst:
	case common.FieldTypeBool:
	case common.FieldTypeTimeZone:

	}
	if "" == name {
		name = propertyType
	}
	return name, skip
}

// addSystemField add system field, get property not return property fields
func addSystemField(fields map[string]Property, objID string, defLang lang.DefaultCCLanguageIf) {
	for key, field := range fields {
		field.ExcelColIndex = field.ExcelColIndex + 1
		fields[key] = field
	}

	idProperty := Property{
		ID:            "",
		Name:          "",
		PropertyType:  common.FieldTypeInt,
		Group:         "defalut",
		ExcelColIndex: 0,
	}

	switch objID {
	case common.BKInnerObjIDHost:
		idProperty.ID = common.BKHostIDField
		idProperty.Name = defLang.Languagef("host_property_bk_host_id")
		fields[idProperty.ID] = idProperty
	case common.BKINnerObjIDObject:
		idProperty.ID = common.BKInstIDField
		idProperty.Name = defLang.Languagef("common_property_bk_inst_id")
		fields[idProperty.ID] = idProperty
	}

}
