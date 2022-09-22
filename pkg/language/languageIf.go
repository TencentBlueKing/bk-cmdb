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

package language

// DefaultCCLanguageIf defines default language interface
type DefaultCCLanguageIf interface {
	// Language returns an content with key
	Language(key string) string
	// Languagef TODO
	// Errorf returns an content with key
	Languagef(key string, args ...interface{}) string
}

// CCLanguageIf defines error information conversion
type CCLanguageIf interface {
	// CreateDefaultCCLanguageIf create new language error interface instance
	CreateDefaultCCLanguageIf(language string) DefaultCCLanguageIf
	// Language returns an content with key
	Language(language, key string) string
	// Languagef returns an content with key
	Languagef(language, key string, args ...interface{}) string

	Load(res map[string]LanguageMap)
}
