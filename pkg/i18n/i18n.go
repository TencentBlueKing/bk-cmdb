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
	"net/http"

	"golang.org/x/text/language"

	ccError "github.com/TencentBlueKing/bk-cmdb/pkg/errors"
	"github.com/TencentBlueKing/bk-cmdb/pkg/log"
)

// defaultManager default i18n manager
var defaultManager Translator

// SetDefaultManager set default i18n manager
func SetDefaultManager(m Translator) {
	defaultManager = m
}

// GetDefaultManager get default i18n manager
func GetDefaultManager() Translator {
	return defaultManager
}

// tagCtxKey define tag for translator
type tagCtxKey struct{}

var tagKey tagCtxKey

// Translator interface
type Translator interface {
	T(ctx context.Context, key string, args ...any) string
	ParseAllowed(code string) (language.Tag, error)
	FallBack() language.Tag
	PickTagStrict(r *http.Request) (language.Tag, error)
	ErrorTranslator
	SysTranslator
}

// ErrorTranslator error translator
type ErrorTranslator interface {
	Error(ctx context.Context, err *ccError.RespError) *ccError.RespError
}

// SysTranslator system translator
type SysTranslator interface {
	Sys(ctx context.Context, key string, args ...any) string
}

// Error translate error
func (m *Manager) Error(ctx context.Context, err *ccError.RespError) *ccError.RespError {
	if err == nil {
		return nil
	}

	err.Message = m.T(ctx, string(err.Code))

	return err
}

// Sys translate key, return key if not found
func (m *Manager) Sys(ctx context.Context, key string, args ...any) string {
	return m.T(ctx, key, args...)
}

// FallBack get fallback language
func (m *Manager) FallBack() language.Tag {
	return language.Make(string(DefaultLanguage))
}

// T general translator, translate key, return key if not found
func (m *Manager) T(ctx context.Context, key string, args ...any) string {
	lang := GetTagFromCtx(ctx)
	if p, ok := m.languagePrinter[lang]; ok {
		return p.Sprintf(key, args...)
	}
	log.Warn(ctx, "translate printer not found", "key", key, "lang", lang)

	return key
}

// PickTagStrict get tag from request and check if it is allowed
func (m *Manager) PickTagStrict(r *http.Request) (language.Tag, error) {

	if c, err := r.Cookie(HTTPCookieLanguage); err == nil && c.Value != "" {
		t, e := m.ParseAllowed(c.Value)
		if e != nil {
			return language.Tag{}, e
		}
		return t, nil
	}

	if h := r.Header.Get(BKHTTPLanguage); h != "" {
		t, e := m.ParseAllowed(h)
		if e != nil {
			return language.Tag{}, e
		}
		return t, nil
	}

	return m.FallBack(), nil
}

// CtxWithLanguageTag set language Tag for context
func CtxWithLanguageTag(ctx context.Context, tag language.Tag) context.Context {
	return context.WithValue(ctx, tagKey, tag)
}

// GetTagFromCtx get translator from context
func GetTagFromCtx(ctx context.Context) language.Tag {
	if v := ctx.Value(tagKey); v != nil {
		if l, ok := v.(language.Tag); ok {
			return l
		}
	}
	return language.Make(string(DefaultLanguage))
}

// ContextWithTag set tag to context
func ContextWithTag(ctx context.Context, tag language.Tag) context.Context {
	return context.WithValue(ctx, tagKey, tag)
}
