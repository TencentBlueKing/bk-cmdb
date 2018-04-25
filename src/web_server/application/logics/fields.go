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
	"configcenter/src/common"
	"configcenter/src/common/blog"
	lang "configcenter/src/common/language"
	webCommon "configcenter/src/web_server/common"
	"fmt"
	simplejson "github.com/bitly/go-simplejson"
	"net/http"
)

// Property object fields
type Property struct {
	ID           string
	Name         string
	PropertyType string
	Option       interface{}
	IsPre        bool
	IsRequire    bool
}

// GetObjFieldIDs get object fields
func GetObjFieldIDs(objID, url string, header http.Header) (map[string]Property, error) {
	url = fmt.Sprintf("%s/api/%s/object/attr/search", url, webCommon.API_VERSION)
	conds := common.KvMap{common.BKObjIDField: objID, common.BKOwnerIDField: common.BKDefaultOwnerID, "page": common.KvMap{"start": 0, "limit": common.BKNoLimit}}
	result, err := httpRequest(url, conds, header)
	if nil != err {
		return nil, err
	}
	blog.Info("get %s fields  url:%s", objID, url)
	blog.Info("get %s fields return:%s", objID, result)
	js, err := simplejson.NewJson([]byte(result))
	if nil != err {
		blog.Info("get %s fields  url:%s return:%s", objID, url, result)
		return nil, err
	}
	fields, _ := js.Get("data").Array()
	ret := make(map[string]Property)

	for _, field := range fields {
		mapField, _ := field.(map[string]interface{})

		fieldName, _ := mapField[common.BKPropertyNameField].(string)
		fieldID, _ := mapField[common.BKPropertyIDField].(string)
		fieldType, _ := mapField[common.BKPropertyTypeField].(string)
		fieldIsRequire, _ := mapField[common.BKIsRequiredField].(bool)
		fieldIsOption, _ := mapField[common.BKOptionField]
		fieldIsPre, _ := mapField[common.BKIsPre].(bool)

		ret[fieldID] = Property{
			ID:           fieldID,
			Name:         fieldName,
			PropertyType: fieldType,
			IsRequire:    fieldIsRequire,
			IsPre:        fieldIsPre,
			Option:       fieldIsOption,
		}
	}

	return ret, nil
}

// getPropertyTypeAliasName  return propertyType name, whether to export,
func getPropertyTypeAliasName(propertyType string, defLang lang.DefaultCCLanguageIf) (string, bool) {
	var skip bool
	name := defLang.Language("field_type_" + propertyType)
	switch propertyType {
	case common.FiledTypeSingleChar:
	case common.FiledTypeLongChar:
	case common.FiledTypeInt:
	case common.FiledTypeEnum:
	case common.FiledTypeDate:
	case common.FiledTypeTime:
	case common.FiledTypeUser:
	case common.FiledTypeSingleAsst:
	case common.FieldTypeMultiAsst:
	case common.FieldTypeBool:
	case common.FieldTypeTimeZone:
	default:
		name = "not found field type"
	}
	return name, skip
}
