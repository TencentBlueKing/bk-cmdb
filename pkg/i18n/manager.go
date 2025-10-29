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
	"embed"
	"encoding/json/v2"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"strings"
	"sync"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"golang.org/x/text/message/catalog"

	"github.com/TencentBlueKing/bk-cmdb/pkg/log"
)

const translationsRoot = "resource"

// embedFS embed translation files
//
//go:embed resource
var embedFS embed.FS

// Options i18n options
type Options struct {
	// DefaultLang If the required language or key does not exist, the default language will be used.
	DefaultLang LanguageType
	// langAbsDir Direct files to be loaded
	langAbsDir string
}

// Manager i18n manager
type Manager struct {
	// attachedFS files for dynamic loading
	attachedFS []string
	// mu for hot update of language configuration or content
	mu      sync.RWMutex
	builder *catalog.Builder
	// fallback is used by default when a language is missing from the translation
	fallback language.Tag
	matcher  language.Matcher
	// languagePrinter stores the printer for each supported language
	languagePrinter map[language.Tag]*message.Printer
	// languages stores all supported languages
	languages []LanguageType
}

// isSupportedTag check if supported
func (m *Manager) isSupportedTag(lang LanguageType) bool {
	_, ok := m.languagePrinter[language.Make(string(lang))]
	return ok
}

// Validator parse language code and check if supported
func (m *Manager) Validator(code LanguageType) (language.Tag, error) {
	if code == "" {
		return language.Tag{}, fmt.Errorf("empty language")
	}

	if m.isSupportedTag(code) {
		return language.Make(string(code)), nil
	}
	return language.Tag{}, fmt.Errorf("unsupported language: %s", code)
}

// cmpKeyWithDefault compare key with default language key
func cmpKeyWithDefault(ctx context.Context, defaultLang, lang map[string]struct{}) bool {
	if len(defaultLang) != len(lang) {
		log.Warn(ctx, "default lang key count not equal with lang", "defaultLangLen", len(defaultLang),
			"langLen", len(lang))
		return false
	}
	isPassed := true
	for k := range defaultLang {
		if _, ok := lang[k]; !ok {
			log.Warn(ctx, "key in defaultLang not found in lang", "key", k)
			isPassed = false
		}
	}
	return isPassed
}

// Match return language Tag which best matches given tags
func (m *Manager) Match(ctx context.Context, tags ...language.Tag) language.Tag {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if m.matcher == nil {
		return m.fallback
	}

	tag, _, _ := m.matcher.Match(tags...)
	if tag == language.Und {
		log.Error(ctx, "match language failed", "tags", tags)
	}
	return tag
}

// Catalog return the catalog builder
func (m *Manager) Catalog() catalog.Catalog {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.builder
}

// FallBack return the defaultLang language Tag
func (m *Manager) FallBack() language.Tag {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.fallback
}

// PickTagStrict get tag from request and check if it is allowed
func (m *Manager) PickTagStrict(r *http.Request) (language.Tag, error) {

	if c, err := r.Cookie(HTTPCookieLanguage); err == nil && c.Value != "" {
		t, e := m.Validator(c.Value)
		if e != nil {
			return language.Tag{}, e
		}
		return t, nil
	}

	if h := r.Header.Get(BKHTTPLanguage); h != "" {
		t, e := m.Validator(h)
		if e != nil {
			return language.Tag{}, e
		}
		return t, nil
	}

	return m.FallBack(), nil
}

func isEmbedFSEmpty(fs embed.FS) bool {
	entries, err := fs.ReadDir(".")
	if err != nil {
		return true
	}
	return len(entries) == 0
}
