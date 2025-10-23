/*
 * TencentBlueKing is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - CMDB) available.
 * Copyright (C) 2025 Tencent. All rights reserved.
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

package i18n

const (
	// HTTPCookieLanguage is the blueking language cookie name
	HTTPCookieLanguage = "blueking_language"
	// BKHTTPLanguage the language key word
	BKHTTPLanguage = "blueking-language"
	// DefaultLanguage the default language
	DefaultLanguage = CN
)

// LanguageType the language type
type LanguageType string

// naming notations：https://i18ns.com/languagecode.html
// Language constant definitions must align with specifications, and remain consistent with resource directory naming
// conventions.
const (
	// CN Chinese
	CN LanguageType = "zh"
	// EN English
	EN LanguageType = "en"
)

var allLanguages = []LanguageType{CN, EN}

func getAllLanguages() []LanguageType {
	return allLanguages
}
