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

// Package i18n for international
package i18n

import (
	"context"
	"fmt"

	"golang.org/x/text/language"

	"github.com/TencentBlueKing/bk-cmdb/pkg/logger"
)

// defaultManager default i18n manager
var defaultManager *Manager

// tagCtxKey define tag for translator
type tagCtxKey struct{}

var tagKey tagCtxKey

// GetTagFromCtx get translator from context
func GetTagFromCtx(ctx context.Context) language.Tag {
	if v := ctx.Value(tagKey); v != nil {
		if l, ok := v.(language.Tag); ok {
			return l
		}
	}
	return language.Make(string(DefaultLanguage))
}

// T translate key, return key if not found
func T(ctx context.Context, key string, args ...any) string {
	lang := GetTagFromCtx(ctx)
	if p, ok := defaultManager.languagePrinter[lang]; ok {
		return p.Sprintf(key, args...)
	}
	logger.Warn(ctx, "translate printer not found", "key", key, "lang", lang)

	// try base language
	baseLang, _ := lang.Base()
	if p, ok := defaultManager.languagePrinter[language.Make(baseLang.String())]; ok {
		return p.Sprintf(key, args...)
	}
	logger.Warn(ctx, "translate base printer not found", "key", key, "baseLang", baseLang)

	// try default language
	if p, ok := defaultManager.languagePrinter[language.Make(string(DefaultLanguage))]; ok {
		return p.Sprintf(key, args...)
	}

	return key
}

// CtxWithLanguageTag set language Tag for context
func CtxWithLanguageTag(ctx context.Context, tag language.Tag) context.Context {
	return context.WithValue(ctx, tagKey, tag)
}

func init() {
	ctx := context.Background()
	manager, err := NewManager(ctx, Options{})
	if err != nil {
		panic(fmt.Errorf("new i18n manager failed: %w", err))
	}
	defaultManager = manager
}
