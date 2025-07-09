/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 Tencent. All rights reserved.
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

package language

// ccDefaultLanguageHelper regular language code helper
type ccDefaultLanguageHelper struct {
	languageType string
	languageStr  func(language, key string) string
	languageStrf func(language, key string, args ...interface{}) string
}

// Language returns an content for specific language
func (cli *ccDefaultLanguageHelper) Language(key string) string {
	ret := cli.languageStr(cli.languageType, key)
	return ret
}

// Languagef returns an content with args for specific language
func (cli *ccDefaultLanguageHelper) Languagef(key string, args ...interface{}) string {
	return cli.languageStrf(cli.languageType, key, args...)
}
