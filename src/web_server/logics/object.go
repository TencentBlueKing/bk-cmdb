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
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/language"
	"configcenter/src/common/mapstr"
)

// GetObjectData get object data
func (lgc *Logics) GetObjectData(ownerID, objID string, header http.Header) ([]interface{}, error) {

	condition := mapstr.MapStr{
		"condition": []string{
			objID,
		},
	}

	result, err := lgc.Engine.CoreAPI.ApiServer().GetObjectData(context.Background(), header, condition)

	if nil != err {
		blog.Errorf("failed to parse the code, error info is %s ", err.Error())
		return nil, err
	}

	if false == result.Result {
		return nil, fmt.Errorf(result.ErrMsg)
	}

	return result.Data[objID].Attr, nil

}

func GetPropertyFieldType(lang language.DefaultCCLanguageIf) map[string]string {
	var fieldType = map[string]string{
		"bk_property_id":         lang.Language("val_type_text"), //"文本",
		"bk_property_name":       lang.Language("val_type_text"), //"文本",
		"bk_property_type":       lang.Language("val_type_text"), //"文本",
		"bk_property_group_name": lang.Language("val_type_text"), // 文本
		"option":                 lang.Language("val_type_text"), //"文本",
		"unit":                   lang.Language("val_type_text"), //"文本",
		"description":            lang.Language("val_type_text"), //"文本",
		"placeholder":            lang.Language("val_type_text"), //"文本",
		"editable":               lang.Language("val_type_bool"), //"布尔",
		"isrequired":             lang.Language("val_type_bool"), //"布尔",
		"isreadonly":             lang.Language("val_type_bool"), //"布尔",
	}
	return fieldType
}

func GetPropertyFieldDesc(lang language.DefaultCCLanguageIf) map[string]string {

	var fields = map[string]string{
		"bk_property_id":         lang.Language("web_en_name_required"),       //"英文名(必填)",
		"bk_property_name":       lang.Language("web_bk_alias_name_required"), //"中文名(必填)",
		"bk_property_type":       lang.Language("web_bk_data_type_required"),  //"数据类型(必填)",
		"bk_property_group_name": lang.Language("property_group"),             // 字段分组
		"option":                 lang.Language("property_option"),            //"数据配置",
		"unit":                   lang.Language("unit"),                       //"单位",
		"description":            lang.Language("desc"),                       //"描述",
		"placeholder":            lang.Language("placeholder"),                //"提示",
		"editable":               lang.Language("is_editable"),                //"是否可编辑",
		"isrequired":             lang.Language("property_is_required"),       //"是否必填",
		"isreadonly":             lang.Language("property_is_readonly"),       //"是否只读",
	}

	return fields
}

func ConvAttrOption(attrItems map[int]map[string]interface{}) {
	for index, attr := range attrItems {

		option, ok := attr[common.BKOptionField].(string)
		if false == ok {
			continue
		}

		if "\"\"" == option {
			option = ""
			attrItems[index][common.BKOptionField] = option
			continue
		}
		fieldType, _ := attr[common.BKPropertyTypeField].(string)
		if common.FieldTypeEnum != fieldType && common.FieldTypeInt != fieldType {
			continue
		}

		var iOption interface{}
		err := json.Unmarshal([]byte(option), &iOption)
		if nil == err {
			attrItems[index][common.BKOptionField] = iOption
		}
	}
}
