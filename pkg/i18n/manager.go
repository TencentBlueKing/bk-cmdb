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
	// Fallback If the required language or key does not exist, the default language will be used.
	Fallback LanguageType
	// AttachedFS Direct files to be loaded
	AttachedFS []string
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

// ParseAllowed parse language code and check if supported
func (m *Manager) ParseAllowed(code string) (language.Tag, error) {
	if code == "" {
		return language.Tag{}, fmt.Errorf("empty language")
	}

	if m.isSupportedTag(LanguageType(code)) {
		return language.Make(code), nil
	}
	return language.Tag{}, fmt.Errorf("unsupported language: %s", code)
}

// NewBaseManager return a new i18n manager with loading different sources
func NewBaseManager(ctx context.Context, opts Options) (*Manager, error) {
	if opts.Fallback == "" {
		opts.Fallback = DefaultLanguage
	}

	i := &Manager{
		fallback:        language.Make(string(opts.Fallback)),
		languagePrinter: make(map[language.Tag]*message.Printer),
		languages:       getAllLanguages(),
		attachedFS:      opts.AttachedFS,
	}
	i.builder = catalog.NewBuilder(catalog.Fallback(i.fallback))
	// load different sources
	paths := make([]string, 0, len(opts.AttachedFS)+1)
	sources := make([]fs.FS, 0, len(opts.AttachedFS)+1)
	if !isEmbedFSEmpty(embedFS) {
		sub, err := fs.Sub(embedFS, translationsRoot)
		if err != nil {
			log.Error(ctx, "fs.Sub on embed failed", "dir", translationsRoot, log.E(err))
			return nil, err
		}
		sources = append(sources, sub)
		paths = append(paths, "embed path")
	}

	for _, path := range opts.AttachedFS {
		if path == "" {
			continue
		}
		diskFs := os.DirFS(path)
		sources = append(sources, diskFs)
		paths = append(paths, path)
	}

	languageKeyMap := make(map[LanguageType]map[string]struct{})
	for idx, src := range sources {
		fileLanguageKeyMap, err := i.loadFromFS(ctx, src)
		if err != nil {
			log.Error(ctx, "load i18n from file system failed", "path", paths[idx], log.E(err))
			return nil, err
		}

		for lang, keyMap := range fileLanguageKeyMap {
			if _, ok := languageKeyMap[lang]; !ok {
				languageKeyMap[lang] = make(map[string]struct{})
			}
			for k := range keyMap {
				languageKeyMap[lang][k] = struct{}{}
			}
		}
	}

	for lang, keyMap := range languageKeyMap {
		if lang == DefaultLanguage {
			continue
		}
		if !cmpKeyWithDefault(ctx, languageKeyMap[DefaultLanguage], keyMap) {
			log.Warn(ctx, "lang key not same with default", "defaultLang", DefaultLanguage, "lang", lang)
		}
	}

	// initialize matcher
	languages := make([]language.Tag, 0)
	for _, lang := range i.languages {
		tag := language.Make(string(lang))
		languages = append(languages, tag)
		i.languagePrinter[tag] = message.NewPrinter(tag, message.Catalog(i.builder))
	}
	i.matcher = language.NewMatcher(languages)

	return i, nil
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

func (m *Manager) loadTranslations(ctx context.Context, lang LanguageType, fsys fs.FS) (map[string]struct{}, error) {
	keyMap := make(map[string]struct{})

	tag := language.Make(string(lang))
	root := string(lang)
	err := fs.WalkDir(fsys, root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Error(ctx, "walk dir entry failed", "path", path, log.E(err))
			return err
		}

		if d.IsDir() {
			return nil
		}

		name := strings.ToLower(d.Name())
		if name != "error.json" && name != "sys.json" {
			return nil
		}

		b, readErr := fs.ReadFile(fsys, path)
		if readErr != nil {
			log.Error(ctx, "read i18n json failed", "path", path, log.E(readErr))
			return readErr
		}

		var jsonMap map[string]string
		if unmarshalErr := json.Unmarshal(b, &jsonMap); unmarshalErr != nil {
			log.Error(ctx, "unmarshal i18n json failed", "path", path, log.E(unmarshalErr))
			return unmarshalErr
		}

		for k, v := range jsonMap {
			if setErr := m.builder.SetString(tag, k, v); setErr != nil {
				log.Error(ctx, "set string failed", "key", k, log.E(setErr))
				return setErr
			}
			if tag == m.fallback {
				if setErr := m.builder.SetString(language.Und, k, v); setErr != nil {
					log.Error(ctx, "set string failed", "key", k, log.E(setErr))
					return setErr
				}
			}
			keyMap[k] = struct{}{}
		}
		return nil
	})

	if err != nil {
		return keyMap, err
	}
	return keyMap, nil
}

// loadFromFS load from file system
func (m *Manager) loadFromFS(ctx context.Context, fsys fs.FS) (map[LanguageType]map[string]struct{}, error) {

	languageKeyMap := make(map[LanguageType]map[string]struct{})
	for _, lang := range m.languages {
		keyMap, err := m.loadTranslations(ctx, lang, fsys)
		if err != nil {
			log.Error(ctx, "load i18n from file system failed", "lang", lang, log.E(err))
			return languageKeyMap, err
		}
		languageKeyMap[lang] = keyMap
	}

	return languageKeyMap, nil
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

// FallBack return the fallback language Tag
func (m *Manager) FallBack() language.Tag {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.fallback
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

func isEmbedFSEmpty(fs embed.FS) bool {
	entries, err := fs.ReadDir(".")
	if err != nil {
		return true
	}
	return len(entries) == 0
}
