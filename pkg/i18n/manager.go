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
	"io/fs"
	"os"
	"strings"
	"sync"

	"golang.org/x/text/language"
	"golang.org/x/text/message/catalog"

	"github.com/TencentBlueKing/bk-cmdb/pkg/logger"
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
	// mu for hot update of language configuration or content
	mu      sync.RWMutex
	builder *catalog.Builder
	// fallback is used by default when a language is missing from the translation
	fallback language.Tag
	matcher  language.Matcher
}

// NewManager return a new i18n manager with loading different sources
func NewManager(ctx context.Context, opts Options) (*Manager, error) {
	if opts.Fallback == "" {
		opts.Fallback = DefaultLanguage
	}

	i := &Manager{
		fallback: language.Make(string(opts.Fallback)),
	}
	i.builder = catalog.NewBuilder(catalog.Fallback(i.fallback))
	// load different sources
	paths := make([]string, 0, len(opts.AttachedFS)+1)
	sources := make([]fs.FS, 0, len(opts.AttachedFS)+1)
	if !isEmbedFSEmpty(embedFS) {
		sub, err := fs.Sub(embedFS, translationsRoot)
		if err != nil {
			logger.Error(ctx, "fs.Sub on embed failed", "dir", translationsRoot, logger.E(err))
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
			logger.Error(ctx, "load i18n from file system failed", "path", paths[idx], logger.E(err))
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
			logger.Warn(ctx, "lang key not same with default", "defaultLang", DefaultLanguage, "lang", lang)
		}
	}

	// initialize matcher
	languages := make([]language.Tag, 0)
	for _, langEntry := range getAllLanguages() {
		languages = append(languages, language.Make(string(langEntry)))
	}
	i.matcher = language.NewMatcher(languages)

	return i, nil
}

// cmpKeyWithDefault compare key with default language key
func cmpKeyWithDefault(ctx context.Context, defaultLang, lang map[string]struct{}) bool {
	if len(defaultLang) != len(lang) {
		logger.Warn(ctx, "default lang key count not equal with lang", "defaultLangLen", len(defaultLang),
			"langLen", len(lang))
		return false
	}
	isPassed := true
	for k := range defaultLang {
		if _, ok := lang[k]; !ok {
			logger.Warn(ctx, "key in defaultLang not found in lang", "key", k)
			isPassed = false
		}
	}
	return isPassed
}

func (i *Manager) loadTranslations(ctx context.Context, lang LanguageType, fsys fs.FS) (map[string]struct{}, error) {
	keyMap := make(map[string]struct{})

	tag := language.Make(string(lang))
	root := string(lang)
	err := fs.WalkDir(fsys, root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			logger.Error(ctx, "walk dir entry failed", "path", path, logger.E(err))
			return err
		}

		if d.IsDir() {
			return nil
		}

		name := d.Name()
		if !strings.HasSuffix(strings.ToLower(name), ".json") {
			return nil
		}

		b, readErr := fs.ReadFile(fsys, path)
		if readErr != nil {
			logger.Error(ctx, "read i18n json failed", "path", path, logger.E(readErr))
			return readErr
		}

		var m map[string]string
		if unmarshalErr := json.Unmarshal(b, &m); unmarshalErr != nil {
			logger.Error(ctx, "unmarshal i18n json failed", "path", path, logger.E(unmarshalErr))
			return unmarshalErr
		}

		for k, v := range m {
			if setErr := i.builder.SetString(tag, k, v); setErr != nil {
				logger.Error(ctx, "set string failed", "key", k, logger.E(setErr))
				return setErr
			}
			if tag == i.fallback {
				if setErr := i.builder.SetString(language.Und, k, v); setErr != nil {
					logger.Error(ctx, "set string failed", "key", k, logger.E(setErr))
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
func (i *Manager) loadFromFS(ctx context.Context, fsys fs.FS) (map[LanguageType]map[string]struct{}, error) {

	languageKeyMap := make(map[LanguageType]map[string]struct{})
	languages := getAllLanguages()
	for _, lang := range languages {
		keyMap, err := i.loadTranslations(ctx, lang, fsys)
		if err != nil {
			logger.Error(ctx, "load i18n from file system failed", "lang", lang, logger.E(err))
			return languageKeyMap, err
		}
		languageKeyMap[lang] = keyMap
	}

	return languageKeyMap, nil
}

// Match return language Tag which best matches given tags
func (i *Manager) Match(ctx context.Context, tags ...language.Tag) language.Tag {
	i.mu.RLock()
	defer i.mu.RUnlock()
	if i.matcher == nil {
		return i.fallback
	}

	tag, _, _ := i.matcher.Match(tags...)
	if tag == language.Und {
		logger.Error(ctx, "match language failed", "tags", tags)
	}
	return tag
}

// Catalog return the catalog builder
func (i *Manager) Catalog() catalog.Catalog {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.builder
}

// Fallback return the fallback language Tag
func (i *Manager) Fallback() language.Tag {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.fallback
}

func isEmbedFSEmpty(fs embed.FS) bool {
	entries, err := fs.ReadDir(".")
	if err != nil {
		return true
	}
	return len(entries) == 0
}
