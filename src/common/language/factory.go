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

// EmptyLanguageSetting empty language setting
var EmptyLanguageSetting = map[string]LanguageMap{}

// New create new Language instance,
// dir is directory of language description resource
func New(dir string) (CCLanguageIf, error) {

	tmp := &ccLanguageHelper{lang: make(map[string]LanguageMap)}

	langType, err := LoadLanguageResourceFromDir(dir)
	if nil != err {
		// blog.Errorf("failed to load the error resource, error info is %s", err.Error())
		return nil, err
	}
	tmp.Load(langType)

	return tmp, nil
}

// NewFromCtx  get lange helper
func NewFromCtx(lang map[string]LanguageMap) CCLanguageIf {
	tmp := &ccLanguageHelper{}
	tmp.Load(lang)
	return tmp
}
