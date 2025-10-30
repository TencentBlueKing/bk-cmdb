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

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"

	ccError "github.com/TencentBlueKing/bk-cmdb/pkg/errors"
)

func makeTestFiles(t *testing.T) string {
	root := t.TempDir()
	base := root
	mustMkdirAll(t, filepath.Join(base, "en"))
	mustMkdirAll(t, filepath.Join(base, "en"))
	mustMkdirAll(t, filepath.Join(base, "zh-cn"))
	mustMkdirAll(t, filepath.Join(base, "zh-cn"))

	writeFile(t, filepath.Join(base, "en", "error.json"), `{ "INVALID_ARGUMENT": "invalid argument","UNKNOWN": 
"unknown error" }`)

	writeFile(t, filepath.Join(base, "en", "sys.json"), `{ "hello": "hello world", 
"meeting": "i have a meeting with %s", "test": "i test %d times" }`)

	writeFile(t, filepath.Join(base, "zh-cn", "error.json"), `
{ "INVALID_ARGUMENT": "参数无效","INVALID_REQUEST": "无效请求"}`)

	writeFile(t, filepath.Join(base, "zh-cn", "sys.json"), `
{ "hello": "你好", "meeting": "我和%s有个会议", "mike": "迈克", "test": "我测试%d次","same": "与上述相同" }`)
	return root
}

func mustMkdirAll(t *testing.T, dir string) {
	t.Helper()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", dir, err)
	}
}

func writeFile(t *testing.T, path string, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}

// 测试有参数自动翻译
// 测试是否支持%d等更多格式
// 测试未设置fallback情况
// 测试设置fallback 没有相关语言情况
// 测试设置fallback 有语言但没有对应翻译情况
// 测试errorCode翻译
// Test_BasicTranslate test basic translate
func Test_BasicTranslate(t *testing.T) {
	root := makeTestFiles(t)
	cxt := context.Background()
	manager, err := NewI18nManager(cxt, Options{langAbsDir: root})
	SetDefaultManager(manager)
	assert.NoError(t, err)

	languageTag := language.Make("zh-cn")
	ctx := ContextWithTag(cxt, languageTag)
	// test basic translate without parameter
	assert.Equal(t, "你好", GetDefaultManager().Sys(ctx, "hello"))
	// test basic translate with parameter
	assert.Equal(t, "我和nancy有个会议", GetDefaultManager().Sys(ctx, "meeting", "nancy"))

	languageTag = language.English
	ctx = ContextWithTag(ctx, languageTag)
	assert.Equal(t, "hello world", GetDefaultManager().Sys(ctx, "hello"))
	assert.Equal(t, "i have a meeting with nancy", GetDefaultManager().Sys(ctx, "meeting", "nancy"))

	// test translate with other format data
	languageTag = language.English
	ctx = ContextWithTag(ctx, languageTag)
	assert.Equal(t, "i test 3 times", GetDefaultManager().Sys(ctx, "test", 3))

	errorManager := ccError.NewErrorManager("cmdb")
	ccError.SetDefaultErrorManager(errorManager)
	testError := ccError.GetDefaultErrorManager().NewRespError(ccError.INVALID_REQUEST)
	testError = manager.Error(ctx, testError)
	assert.Equal(t, "invalid request", testError.Message)

	languageTag = language.Make("zh-cn")
	ctx = ContextWithTag(ctx, languageTag)
	testError = manager.Error(ctx, testError)
	assert.Equal(t, "无效请求", testError.Message)
}
