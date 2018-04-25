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
	"configcenter/src/common/language"
)

func TranslateOpLanguage(data interface{}, defLang language.DefaultCCLanguageIf) interface{} {

	mapData, ok := data.(map[string]interface{})
	if false == ok {
		return data
	}
	if nil == mapData {
		return data
	}

	info, ok := mapData["info"].([]map[string]interface{})

	if false == ok {
		return data
	}

	if nil == info {
		return data
	}
	for index, row := range info {
		opDesc, ok := row[common.BKOpDescField].(string)
		if false == ok {
			continue
		}
		newDesc := defLang.Language("auditlog_" + opDesc)
		if "" == newDesc {
			return opDesc
		}
		info[index][common.BKOpDescField] = opDesc
	}

	mapData["info"] = info
	return mapData

}
