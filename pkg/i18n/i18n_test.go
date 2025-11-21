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
	"os"
	"path/filepath"
	"testing"

	"github.com/TencentBlueKing/bk-cmdb/pkg/errors"
	"github.com/TencentBlueKing/bk-cmdb/pkg/kit"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func makeEmptyDir(t *testing.T) string {
	root := t.TempDir()
	base := root
	mustMkdirAll(t, filepath.Join(base, "en"))
	mustMkdirAll(t, filepath.Join(base, "zh-cn"))

	writeFile(t, filepath.Join(base, "en", "error.json"), `{}`)
	writeFile(t, filepath.Join(base, "en", "sys.json"), `{}`)
	writeFile(t, filepath.Join(base, "zh-cn", "error.json"), `{}`)
	writeFile(t, filepath.Join(base, "zh-cn", "sys.json"), `{}`)
	return root
}

func makeTestFiles(t *testing.T, root string) string {
	base := root
	mustMkdirAll(t, filepath.Join(base, "en"))
	mustMkdirAll(t, filepath.Join(base, "zh-cn"))

	writeFile(t, filepath.Join(base, "en", "error.json"), `{ "Test_INVALID_REQUEST": "invalid request",
"Test_UNKNOWN":"unknown error","Test_INVALID_ARGUMENT": "invalid argument" }`)

	writeFile(t, filepath.Join(base, "en", "sys.json"), `{ "hello": "hello world", 
"meeting": "i have a meeting with %s", "test": "i test %d times","mike": "mike", "same": "same as above" }`)

	writeFile(t, filepath.Join(base, "zh-cn", "error.json"), `
{ "Test_INVALID_ARGUMENT": "参数无效","Test_INVALID_REQUEST": "无效请求","Test_UNKNOWN":"未知错误"}`)

	writeFile(t, filepath.Join(base, "zh-cn", "sys.json"), `
{ "hello": "你好", "meeting": "我和%s有个会议", "mike": "迈克", "test": "我测试%d次","same": "与上述相同" }`)
	return root
}

func makeDyTestFiles(t *testing.T, root string) string {
	base := makeTestFiles(t, root)
	mustMkdirAll(t, filepath.Join(base, "ko"))

	writeFile(t, filepath.Join(base, "ko", "error.json"), `{
  "Test_INVALID_REQUEST": "잘못된 요청",
  "Test_INVALID_ARGUMENT": "잘못된 인수",
  "Test_UNKNOWN": "알 수 없는 오류",
  "INVALID_REQUEST": "요청이 잘못되었거나 처리할 수 없습니다.",
  "INVALID_ARGUMENT": "요청에 잘못된 인수가 포함되어 있습니다.",
  "UNAUTHORIZED": "요청된 리소스에 대한 인증 자격 증명이 누락되었거나 유효하지 않습니다.",
  "FORBIDDEN": "사용자가 리소스에 대한 필요한 권한이 없거나 일종의 계정이 필요하거나 금지된 작업을 시도했습니다.",
  "TOO_MANY_REQUESTS": "요청 속도 제한을 초과했습니다. 나중에 다시 시도하십시오.",
  "SERVER_ERROR": "내부 서버 오류.",
  "UNKNOWN_ERROR": "예기치 않은 조건이 발생했으며 더 구체적인 메시지가 적합하지 않습니다.",
  "NOT_FOUND": "리소스를 찾을 수 없습니다.",
  "METHOD_NOT_ALLOWED": "요청에 지정된 메서드가 허용되지 않습니다."

}`)

	writeFile(t, filepath.Join(base, "ko", "sys.json"), `{
  "hello": "안녕 하세요 세계",
  "meeting": "나는 %s와 회의가 있어요",
  "test": "나는 %d번 테스트해요",
  "mike": "마이크",
  "same": "위와 동일",
  "Tom": "톰은 고양이입니다",
  "Jerry": "제리는 쥐입니다",
  "Spike": "스파이크는 개입니다"
}`)

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
// 测试未设置defaultLang情况
// 测试设置defaultLang 没有相关语言情况
// 测试设置defaultLang 有语言但没有对应翻译情况
// 测试errorCode翻译
// Test_BasicTranslate test basic translate
func Test_BasicTranslate(t *testing.T) {
	root := t.TempDir()
	root = makeTestFiles(t, root)
	kt := kit.DefaultKit()
	err := Init(kt, &Options{LanguageDir: root, DefaultLang: CN, RequireExternalDir: true})
	assert.NoError(t, err)

	kt.Lang = string(CN)
	// test basic translate without parameter
	assert.Equal(t, "你好", Sys(kt, "hello"))
	// test basic translate with parameter
	assert.Equal(t, "我和nancy有个会议", Sys(kt, "meeting", "nancy"))

	kt.Lang = string(EN)
	assert.Equal(t, "hello world", Sys(kt, "hello"))
	assert.Equal(t, "i have a meeting with nancy", Sys(kt, "meeting", "nancy"))

	// test translate with other format data
	kt.Lang = string(EN)
	assert.Equal(t, "i test 3 times", Sys(kt, "test", 3))

	err = cerr.Init()
	require.NoError(t, err)
	codeErr := cerr.NewError("Test_INVALID_REQUEST", "invalid request")
	convErr := cerr.ErrorClient().ConvToRespError(codeErr)
	testError := RespError(kt, convErr)
	assert.Equal(t, "invalid request", testError.Message)

	kt.Lang = string(CN)
	testError = RespError(kt, testError)
	assert.Equal(t, "无效请求", testError.Message)
}

// 测试动态加载已存在语言
// 测试动态加载默认语言
// 测试动态加载不支持默认语言
func Test_DynamicTranslate(t *testing.T) {
	root := makeEmptyDir(t)
	kt := kit.DefaultKit()
	err := Init(kt, &Options{LanguageDir: root, RequireExternalDir: true})
	require.NoError(t, err)

	// 未进行动态加载，返回原始key
	kt.Lang = string(CN)
	// test basic translate without parameter
	assert.Equal(t, "hello", Sys(kt, "hello"))
	// test basic translate with parameter
	assert.Equal(t, "meeting", Sys(kt, "meeting", "nancy"))

	root = makeTestFiles(t, root)
	// 测试动态加载已存在语言
	err = Reload(kt, &Options{LanguageDir: root, DefaultLang: CN})
	assert.NoError(t, err)

	// test without setting default language
	kt.Lang = ""
	assert.Equal(t, "你好", Sys(kt, "hello"))

	// test basic translate without parameter
	kt.Lang = string(CN)
	assert.Equal(t, "你好", Sys(kt, "hello"))
	// test basic translate with parameter
	assert.Equal(t, "我和nancy有个会议", Sys(kt, "meeting", "nancy"))

	kt.Lang = string(EN)
	assert.Equal(t, "hello world", Sys(kt, "hello"))
	assert.Equal(t, "i have a meeting with nancy", Sys(kt, "meeting", "nancy"))

	// test translate with other format data
	kt.Lang = string(EN)
	assert.Equal(t, "i test 3 times", Sys(kt, "test", 3))

	err = cerr.Init()
	require.NoError(t, err)
	codeErr := cerr.NewError("Test_INVALID_REQUEST", "invalid request")
	convErr := cerr.ErrorClient().ConvToRespError(codeErr)
	testError := RespError(kt, convErr)
	assert.Equal(t, "invalid request", testError.Message)

	kt.Lang = string(CN)
	testError = RespError(kt, testError)
	assert.Equal(t, "无效请求", testError.Message)

	// 测试动态加载设置默认语言
	err = Reload(kt, &Options{LanguageDir: root, DefaultLang: EN})
	assert.NoError(t, err)
	// test without setting default language
	kt.Lang = ""
	assert.Equal(t, "i have a meeting with nancy", Sys(kt, "meeting", "nancy"))

	// 测试动态加载设置默认语言后不存在key情况
	kt.Lang = string(EN)
	// 测试动态加载不支持默认语言
	err = Reload(kt, &Options{LanguageDir: root, DefaultLang: LanguageType("af")})
	assert.Error(t, err)

	// 测试动态加载当前不存在语言及翻译
	dyRoot := makeDyTestFiles(t, root)
	err = Reload(kt, &Options{LanguageDir: dyRoot, DefaultLang: CN})
	assert.NoError(t, err)

	kt.Lang = "ko"
	// 测试sys与error区分key
	assert.Equal(t, "Test_INVALID_REQUEST", Sys(kt, "Test_INVALID_REQUEST"))
	assert.Equal(t, "안녕 하세요 세계", Sys(kt, "hello"))
}
