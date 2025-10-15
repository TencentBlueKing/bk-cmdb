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

// Package i18n
package i18n

import (
	"context"

	"golang.org/x/text/language"
	"golang.org/x/text/message"

	"github.com/TencentBlueKing/bk-cmdb/pkg/logger"
)

// defaultI18NManager default i18n manager
var defaultI18NManager *I18NManager

// TranslatePrinter translate printer
type TranslatePrinter struct {
	printer *message.Printer
}

// translatorCtxKey define context key for translator
type translatorCtxKey struct{}

var translatorKey translatorCtxKey

// GetTranslatePrinter get translator from context
func GetTranslatePrinter(ctx context.Context) *TranslatePrinter {
	if v := ctx.Value(translatorKey); v != nil {
		if l, ok := v.(*TranslatePrinter); ok {
			return l
		}
	}
	return nil
}

// T translate key, return key if not found
func T(ctx context.Context, key string, args ...any) string {
	if l := GetTranslatePrinter(ctx); l != nil {
		return l.printer.Sprintf(key, args...)
	}
	return key
}

// CtxWithLanguageTag set language Tag for context
func CtxWithLanguageTag(ctx context.Context, m *I18NManager, tag language.Tag) context.Context {
	p := message.NewPrinter(tag, message.Catalog(m.Catalog()))
	printer := &TranslatePrinter{printer: p}
	return context.WithValue(ctx, translatorKey, printer)
}

func init() {
	ctx := context.Background()
	manager, err := NewI18NManager(ctx, Options{})
	if err != nil {
		logger.Error(ctx, "new i18n manager failed", err)
		return
	}
	defaultI18NManager = manager
}
