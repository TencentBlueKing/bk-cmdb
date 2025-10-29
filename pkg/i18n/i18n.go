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

// Translator i18n interface
type Translator interface {
	ErrorTranslator
	SysTranslator
	BaseTranslator
}

// BaseTranslator support base translator
type BaseTranslator interface {
	T(ctx context.Context, key string, args ...any) string
	ParseAllowed(code string) (language.Tag, error)
	FallBack() language.Tag
	PickTagStrict(r *http.Request) (language.Tag, error)
}

// tagCtxKey define tag for translator
type tagCtxKey struct{}

var tagKey tagCtxKey

// DefaultSysTranslator default system translator
type DefaultSysTranslator struct {
	base BaseTranslator
}

// Sys translate key, return key if not found
func (s *DefaultSysTranslator) Sys(ctx context.Context, key string, args ...any) string {
	return s.base.T(ctx, key, args...)
}

// DefaultErrorTranslator default error translator
type DefaultErrorTranslator struct {
	base BaseTranslator
}

// Error translate error
func (e *DefaultErrorTranslator) Error(ctx context.Context, err *ccError.RespError) *ccError.RespError {
	if err == nil {
		return nil
	}

	err.Message = e.base.T(ctx, string(err.Code))

	return err
}

// ErrorTranslator error translator
type ErrorTranslator interface {
	Error(ctx context.Context, err *ccError.RespError) *ccError.RespError
}

// SysTranslator system translator
type SysTranslator interface {
	Sys(ctx context.Context, key string, args ...any) string
}

// TranslatorManager translator manager
type TranslatorManager struct {
	base BaseTranslator
	err  ErrorTranslator
	sys  SysTranslator
}

// InitTranslatorManager init translator manager with default error and system translator
func InitTranslatorManager(ctx context.Context, opts Options) (*TranslatorManager, error) {
	baseManager, err := NewBaseManager(ctx, opts)
	if err != nil {
		log.Error(ctx, "new base manager failed", "err", err)
		return nil, err
	}

	m := &TranslatorManager{
		base: baseManager,
		err:  &DefaultErrorTranslator{base: baseManager},
		sys:  &DefaultSysTranslator{base: baseManager},
	}

	return m, nil
}

// NewTranslatorManager new translator manager with base translator for customization
func NewTranslatorManager(base BaseTranslator, opts ...TranslatorOption) *TranslatorManager {
	m := &TranslatorManager{
		base: base,
		err:  &DefaultErrorTranslator{base: base},
		sys:  &DefaultSysTranslator{base: base},
	}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

// TranslatorOption translator option
type TranslatorOption func(manager *TranslatorManager)

// WithErrorTranslator set error translator
func WithErrorTranslator(e ErrorTranslator) TranslatorOption {
	return func(m *TranslatorManager) { m.err = e }
}

// WithSysTranslator set system translator
func WithSysTranslator(s SysTranslator) TranslatorOption {
	return func(m *TranslatorManager) { m.sys = s }
}

// T translate key, base translate function
func (m *TranslatorManager) T(ctx context.Context, key string, args ...any) string {
	return m.base.T(ctx, key, args...)
}

// ParseAllowed parse allowed language
func (m *TranslatorManager) ParseAllowed(code string) (language.Tag, error) {
	return m.base.ParseAllowed(code)
}

// FallBack get fallback language
func (m *TranslatorManager) FallBack() language.Tag {
	return m.base.FallBack()
}

// PickTagStrict pick tag from request, and check if it is allowed
func (m *TranslatorManager) PickTagStrict(r *http.Request) (language.Tag, error) {
	return m.base.PickTagStrict(r)
}

// Error translate error
func (m *TranslatorManager) Error(ctx context.Context, err *ccError.RespError) *ccError.RespError {
	return m.err.Error(ctx, err)
}

// Sys translate system key
func (m *TranslatorManager) Sys(ctx context.Context, key string, args ...any) string {
	return m.sys.Sys(ctx, key, args...)
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
