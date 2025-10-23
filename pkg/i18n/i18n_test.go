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
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"
)

func makeTestFiles(t *testing.T) string {
	root := t.TempDir()
	base := root
	mustMkdirAll(t, filepath.Join(base, "en"))
	mustMkdirAll(t, filepath.Join(base, "en"))
	mustMkdirAll(t, filepath.Join(base, "zh"))
	mustMkdirAll(t, filepath.Join(base, "zh"))

	writeFile(t, filepath.Join(base, "en", "error.json"), `{ "1199000": "Test Error" }`)

	writeFile(t, filepath.Join(base, "en", "sys.json"), `{ "hello": "hello world", 
"meeting": "i have a meeting with %s", "test": "i test %d times" }`)

	writeFile(t, filepath.Join(base, "zh", "error.json"), `
{ "1199000": "测试错误" }`)

	writeFile(t, filepath.Join(base, "zh", "sys.json"), `
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
// Test_BasicTranslate test basic translate
func Test_BasicTranslate(t *testing.T) {
	root := makeTestFiles(t)
	cxt := context.Background()
	manager, err := NewManager(cxt, Options{AttachedFS: []string{root}})
	assert.NoError(t, err)

	languageTag := language.Chinese
	ctx := CtxWithLanguageTag(cxt, languageTag)
	// test basic translate without parameter
	assert.Equal(t, "你好", Tm(ctx, manager, "hello"))
	// test basic translate with parameter
	assert.Equal(t, "我和nancy有个会议", Tm(ctx, manager, "meeting", "nancy"))

	languageTag = language.English
	ctx = CtxWithLanguageTag(ctx, languageTag)
	assert.Equal(t, "hello world", Tm(ctx, manager, "hello"))
	assert.Equal(t, "i have a meeting with nancy", Tm(ctx, manager, "meeting", "nancy"))

	// test language not supported
	languageTag = language.Japanese
	ctx = CtxWithLanguageTag(ctx, languageTag)
	assert.Equal(t, "你好", Tm(ctx, manager, "hello"))
	assert.Equal(t, "我和nancy有个会议", Tm(ctx, manager, "meeting", "nancy"))

	// test language supported but no translation
	languageTag = language.English
	ctx = CtxWithLanguageTag(ctx, languageTag)
	assert.Equal(t, "与上述相同", Tm(ctx, manager, "same"))

	// test translate with other format data
	languageTag = language.English
	ctx = CtxWithLanguageTag(ctx, languageTag)
	assert.Equal(t, "i test 3 times", Tm(ctx, manager, "test", 3))
}

// Tm translate message, use for test with test translation resources
func Tm(ctx context.Context, m *Manager, key string, args ...any) string {
	lang := GetTagFromCtx(ctx)
	if p, ok := m.languagePrinter[lang]; ok {
		return p.Sprintf(key, args...)
	}
	fmt.Printf("translate printer not found")

	// try base language
	baseLang, _ := lang.Base()
	if p, ok := m.languagePrinter[language.Make(baseLang.String())]; ok {
		return p.Sprintf(key, args...)
	}
	fmt.Printf("translate base printer not found")

	// try default language
	if p, ok := m.languagePrinter[language.Make(string(DefaultLanguage))]; ok {
		return p.Sprintf(key, args...)
	}

	return key
}
