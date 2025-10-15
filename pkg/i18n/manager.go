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

// embedI18NFS embed translation files
//
//go:embed resource
var embedI18NFS embed.FS

// Options i18n options
type Options struct {
	// Fallback If the required language or key does not exist, the default language will be used.
	Fallback LanguageType
	// AttachedFS Direct files to be loaded
	AttachedFS []string
}

// I18NManager i18n manager
type I18NManager struct {
	// mu for hot update of language configuration or content
	mu      sync.RWMutex
	builder *catalog.Builder
	// fallback is used by default when a language is missing from the translation
	fallback language.Tag
	matcher  language.Matcher
}

// NewI18NManager return a new i18n manager with loading different sources
func NewI18NManager(ctx context.Context, opts Options) (*I18NManager, error) {
	if opts.Fallback == "" {
		opts.Fallback = DefaultLanguage
	}

	i := &I18NManager{
		fallback: language.Make(string(opts.Fallback)),
	}
	i.builder = catalog.NewBuilder(catalog.Fallback(i.fallback))
	// load different sources
	paths := make([]string, 0, len(opts.AttachedFS)+1)
	sources := make([]fs.FS, 0, len(opts.AttachedFS)+1)
	if !isEmbedFSEmpty(embedI18NFS) {
		sub, err := fs.Sub(embedI18NFS, translationsRoot)
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

	for idx, src := range sources {
		if err := i.loadFromFS(src); err != nil {
			logger.Error(ctx, "load i18n from file system failed", "path", paths[idx], logger.E(err))
			return nil, err
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

func (i *I18NManager) setValues(ctx context.Context, lang LanguageType, fsys fs.FS) (int, error) {
	count := 0

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

		count += len(m)
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
		}
		return nil
	})

	if err != nil {
		return 0, err
	}
	return count, nil
}

// loadFromFS load from file system
func (i *I18NManager) loadFromFS(fsys fs.FS) error {

	ctx := context.Background()
	defaultCount, err := i.setValues(ctx, DefaultLanguage, fsys)
	if err != nil {
		logger.Error(ctx, "load i18n from file system failed", logger.E(err))
		return err
	}
	languages := getAllLanguages()
	for _, lang := range languages {
		if lang == DefaultLanguage {
			continue
		}
		count, err := i.setValues(ctx, lang, fsys)
		if err != nil {
			logger.Error(ctx, "load i18n from file system failed", logger.E(err))
			return err
		}
		if count != defaultCount {
			logger.Warn(ctx, "default language count is not equal to language count", "lang", lang,
				"count", count, "defaultCount", defaultCount)
		}
	}

	return nil
}

// Match return language Tag which best matches given tags
func (i *I18NManager) Match(tags ...language.Tag) language.Tag {
	i.mu.RLock()
	defer i.mu.RUnlock()
	if i.matcher == nil {
		return i.fallback
	}

	tag, _, _ := i.matcher.Match(tags...)
	if tag == language.Und {
		logger.Error(context.Background(), "match language failed", "tags", tags)
	}
	return tag
}

// Catalog return the catalog builder
func (i *I18NManager) Catalog() catalog.Catalog {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.builder
}

// Fallback return the fallback language Tag
func (i *I18NManager) Fallback() language.Tag {
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
