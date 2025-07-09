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

package cryptor

import (
	"testing"
)

var aesCryp Cryptor

func init() {
	key := "123456781234567812345678"
	aesCryp = NewAesEncrpytor(key)
}

func TestAES(t *testing.T) {
	plainText := "hello world"
	cryptedText, err := aesCryp.Encrypt(plainText)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("cryptedText is %s\n", cryptedText)

	result, err := aesCryp.Decrypt(cryptedText)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("result is %s\n", result)

	if result != plainText {
		t.Fatal("AES encrypt & decrypt fail")
	}
}
