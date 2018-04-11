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

// UnknownTheCodeStrf const define unknow code
const UnknownTheKeyStrf = "the key is '%v', unknown meaning"

// UnknownTheLanguageStrf define unknow language
const UnknownTheLanguageStrf = "the language code is '%s', unknown meaning"

// defaultLanguage default language package name
const defaultLanguage = "default"

// LanguageMap  mapping
type LanguageMap map[string]string

// ccError  CC custom error  defind
type ccLanguage struct {
	key      string
	callback func() string
}
