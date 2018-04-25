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

import (
	"fmt"
	"testing"
)

func TestLoad(t *testing.T) {

	cclangEr, err := New("./examples/errorres")
	if nil != err {
		t.Errorf("failed to create cc error manager, error info is %s", err.Error())
		return
	}

	testKey := "test_key_val"
	fmt.Printf("\n[%s]%s\n", testKey, cclangEr.Language("en", testKey))
	fmt.Printf("\n[%s]%s\n", testKey, cclangEr.Language("cn", testKey))

	testKey = "test_key_val1"
	fmt.Printf("\n[%s]%s\n", testKey, cclangEr.Language("en", testKey))
	fmt.Printf("\n[%s]%s\n", testKey, cclangEr.Language("cn", testKey))
	testKey = "test_key_val2"
	fmt.Printf("\n[%s]%s\n", testKey, cclangEr.Language("en", testKey))
	fmt.Printf("\n[%s]%s\n", testKey, cclangEr.Language("cn", testKey))
	testKey = "test_key_val3"
	fmt.Printf("\n[%s]%s\n", testKey, cclangEr.Language("en", testKey))
	fmt.Printf("\n[%s]%s\n", testKey, cclangEr.Language("cn", testKey))

	testKey = "test_key_format_str"
	fmt.Printf("\n[%s]%s\n", testKey, cclangEr.Languagef("en", testKey, "XXX"))
	fmt.Printf("\n[%s]%s\n\n", testKey, cclangEr.Languagef("cn", testKey, "XXX"))
	testKey = "test_key_format_int"
	fmt.Printf("\n[%s]%s\n", testKey, cclangEr.Languagef("en", testKey, 0))
	fmt.Printf("\n[%s]%s\n\n", testKey, cclangEr.Languagef("cn", testKey, 0))

	defaultLanguage := cclangEr.CreateDefaultCCLanguageIf("cn")

	testKey = "test_key_val"
	fmt.Printf("\ndefault[%s]%s\n", testKey, defaultLanguage.Language(testKey))
	testKey = "test_key_val1"
	fmt.Printf("\ndefault[%s]%s\n", testKey, defaultLanguage.Language(testKey))
	testKey = "test_key_val2"
	fmt.Printf("\ndefault[%s]%s\n", testKey, defaultLanguage.Language(testKey))
	testKey = "test_key_format_str"
	fmt.Printf("\ndefault[%s]%s\n", testKey, defaultLanguage.Languagef(testKey, "XXX"))
	testKey = "test_key_format_int"
	fmt.Printf("\ndefault[%s]%s\n\n", testKey, defaultLanguage.Languagef(testKey, 123))

	defaultLanguage = cclangEr.CreateDefaultCCLanguageIf("en")
	testKey = "test_key_val"
	fmt.Printf("\ndefault[%s]%s\n", testKey, defaultLanguage.Language(testKey))
	testKey = "test_key_val1"
	fmt.Printf("\ndefault[%s]%s\n", testKey, defaultLanguage.Language(testKey))
	testKey = "test_key_val2"
	fmt.Printf("\ndefault[%s]%s\n", testKey, defaultLanguage.Language(testKey))
	testKey = "test_key_format_str"
	fmt.Printf("\ndefault[%s]%s\n", testKey, defaultLanguage.Languagef(testKey, "XXX"))
	testKey = "test_key_format_int"
	fmt.Printf("\ndefault[%s]%s\n\n", testKey, defaultLanguage.Languagef(testKey, 123))

	testKey = "space    "
	fmt.Printf("\ndefault key:%s  content:%s \n\n", testKey, defaultLanguage.Language(testKey))
	testKey = "space    space"
	fmt.Printf("\ndefault  key:%s  content:%s \n\n", testKey, defaultLanguage.Language(testKey))
	testKey = "space    space_not_found"
	fmt.Printf("\ndefault not found   key:%s  content:%s\n\n", testKey, defaultLanguage.Language(testKey))

}
