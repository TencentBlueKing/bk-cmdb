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
	"encoding/json/v2"
	"fmt"
	"github.com/golang/protobuf/descriptor"
	"golang.org/x/text/message"
	"golang.org/x/text/message/catalog"
	"io/fs"
	"os"
	"strings"
	"sync"

	"golang.org/x/text/language"

	ccError "github.com/TencentBlueKing/bk-cmdb/pkg/errors"
	"github.com/TencentBlueKing/bk-cmdb/pkg/log"
)

// defaultManager default i18n manager
var defaultManager I18nInterface

// SetDefaultManager set default i18n manager
func SetDefaultManager(m I18nInterface) {
	defaultManager = m
}

// GetDefaultManager get default i18n manager
func GetDefaultManager() I18nInterface {
	return defaultManager
}

// I18nInterface i18n interface
type I18nInterface interface {
	Error(ctx context.Context, err *ccError.RespError) *ccError.RespError
	Sys(ctx context.Context, key string, args ...any) string
	Validator(code LanguageType) (language.Tag, error)
}

type i18n struct {
	// langAbsPath file path for loading language
	langAbsPath string
	matcher     language.Matcher
	languageBaseInterface
}

// Validator parse language code and check if supported
func (i *i18n) Validator(code LanguageType) (language.Tag, error) {
	if code == "" {
		return language.Tag{}, fmt.Errorf("empty language")
	}

	if i.isSupportedTag(code) {
		return language.Make(string(code)), nil
	}
	return language.Tag{}, fmt.Errorf("unsupported language: %s", code)
}

type languageBase struct {
	builder *catalog.Builder
	// mu for hot update of language configuration or content
	mu sync.RWMutex
	// languagePrinter stores the printer for each supported language
	languagePrinter map[language.Tag]*message.Printer
	// languages stores all supported languages
	languages   []LanguageType
	defaultLang language.Tag
}

type languageBaseInterface interface {
	T(ctx context.Context, key string, args ...any) string
	isSupportedTag(lang LanguageType) bool
	getDefaultLang() language.Tag
	setBuilder(b *catalog.Builder)
}

func (l *languageBase) isSupportedTag(lang LanguageType) bool {
	l.mu.RLock()
	defer l.mu.RUnlock()
	_, ok := l.languagePrinter[language.Make(string(lang))]
	return ok
}

// T general translator, translate key, return key if not found
func (l *languageBase) T(ctx context.Context, key string, args ...any) string {
	lang := GetTagFromCtx(ctx)
	if p, ok := l.languagePrinter[lang]; ok {
		return p.Sprintf(key, args...)
	}
	log.Warn(ctx, "translate printer not found", "key", key, "lang", lang)

	return key
}

func (l *languageBase) getDefaultLang() language.Tag {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.defaultLang
}

func (l *languageBase) setBuilder(b *catalog.Builder) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.builder = b
}

// loadFromFS load from file system
func (l *languageBase) loadFromFS(ctx context.Context, fsys fs.FS) (map[LanguageType]map[string]struct{}, error) {

	languageKeyMap := make(map[LanguageType]map[string]struct{})
	for _, lang := range l.languages {
		keyMap, err := l.loadTranslations(ctx, lang, fsys)
		if err != nil {
			log.Error(ctx, "load i18n from file system failed", "lang", lang, log.E(err))
			return languageKeyMap, err
		}
		languageKeyMap[lang] = keyMap
	}

	return languageKeyMap, nil
}

func (l *languageBase) loadTranslations(ctx context.Context, lang LanguageType, fsys fs.FS) (map[string]struct{}, error) {
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
			if setErr := l.builder.SetString(tag, k, v); setErr != nil {
				log.Error(ctx, "set string failed", "key", k, log.E(setErr))
				return setErr
			}
			if tag == l.defaultLang {
				if setErr := l.builder.SetString(language.Und, k, v); setErr != nil {
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

// NewManager return a new i18n manager with loading different sources
func NewManager(ctx context.Context, opts Options) (*Manager, error) {
	if opts.DefaultLang == "" {
		opts.DefaultLang = DefaultLanguage
	}

	i := &i18n{
		langAbsPath: opts.langAbsDir,
		languageBaseInterface: &languageBase{
			languagePrinter: make(map[language.Tag]*message.Printer),
			languages:       getAllLanguages(),
			defaultLang:     language.Make(string(opts.DefaultLang)),
		},
	}

	i.setBuilder(catalog.NewBuilder(catalog.Fallback(i.getDefaultLang())))
	// load different sources
	paths := make([]string, 0, len(opts.langAbsDir)+1)
	sources := make([]fs.FS, 0, len(opts.langAbsDir)+1)
	if !isEmbedFSEmpty(embedFS) {
		sub, err := fs.Sub(embedFS, translationsRoot)
		if err != nil {
			log.Error(ctx, "fs.Sub on embed failed", "dir", translationsRoot, log.E(err))
			return nil, err
		}
		sources = append(sources, sub)
		paths = append(paths, "embed path")
	}

	if opts.langAbsDir != "" {
		diskFs := os.DirFS(opts.langAbsDir)
		sources = append(sources, diskFs)
		paths = append(paths, opts.langAbsDir)
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

// ContextWithTag set tag to context
func ContextWithTag(ctx context.Context, tag language.Tag) context.Context {
	return context.WithValue(ctx, tagKey, tag)
}
