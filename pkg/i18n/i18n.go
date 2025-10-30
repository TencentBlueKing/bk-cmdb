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
	"embed"
	"encoding/json/v2"
	"fmt"
	"io/fs"
	"os"
	"strings"
	"sync"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"golang.org/x/text/message/catalog"

	"github.com/TencentBlueKing/bk-cmdb/pkg/errors"
	"github.com/TencentBlueKing/bk-cmdb/pkg/log"
)

const translationsRoot = "resource"

// embedFS embed translation files
//
//go:embed resource
var embedFS embed.FS

// defaultManager default i18n manager
var defaultManager Interface

// SetDefaultManager set default i18n manager
func SetDefaultManager(m Interface) {
	defaultManager = m
}

// GetDefaultManager get default i18n manager
func GetDefaultManager() Interface {
	return defaultManager
}

// i18n define implementation required components for language package.
type i18n struct {
	// languageDir file path for loading language, it requires folders for various language packs, the naming of
	//language pack folders must comply with the system language definition specification.
	languageDir string
	// languageBaseInterface language base interface, support base language translation service.
	languageBaseInterface
}

// Interface i18n interface for multilingual internationalization, starting from the scenario, it can be divided
// into two types: implementing error translation and built-in system translation
type Interface interface {
	// RespError translate error info, translate for error message by error code which is pre-determined
	RespError(ctx context.Context, err error) *cerr.RespError
	// Sys translate key, return key if not found
	Sys(ctx context.Context, key string, args ...any) string
	// Validate get tag from request and check if it is supported language
	Validate(lang LanguageType) (LanguageType, error)
}

// Validate get tag from request and check if it is allowed
func (i *i18n) Validate(lang LanguageType) (LanguageType, error) {
	ok := i.isSupportedLang(lang)
	if !ok {
		e := fmt.Errorf("unsupported language: %s", lang)
		return i.getDefaultLang(), e
	}

	return lang, nil
}

// RespError translate response error message by error code
func (i *i18n) RespError(ctx context.Context, err error) *cerr.RespError {
	if err == nil {
		return nil
	}

	respErr := cerr.GetDefaultErrorManager().ConvToRespError(err)
	respErr.Message = i.T(ctx, string(respErr.Code))

	return respErr
}

// Sys translate key, return key if not found
func (i *i18n) Sys(ctx context.Context, key string, args ...any) string {
	return i.T(ctx, key, args...)
}

type multilingual struct {
	// mu for hot update of language configuration or content
	mu sync.RWMutex
	// languagePrinter stores the printer for each supported language
	languagePrinter map[LanguageType]*message.Printer
	// languages stores all supported languages
	languages   []LanguageType
	defaultLang LanguageType
	builder     *catalog.Builder
}

type languageBaseInterface interface {
	T(ctx context.Context, key string, args ...any) string
	isSupportedLang(lang LanguageType) bool
	getDefaultLang() LanguageType
}

func (l *multilingual) isSupportedLang(lang LanguageType) bool {
	l.mu.RLock()
	defer l.mu.RUnlock()
	_, ok := l.languagePrinter[lang]
	return ok
}

// T general translator, translate key, return key if not found
func (l *multilingual) T(ctx context.Context, key string, args ...any) string {
	lang := GetLangFromCtx(ctx)
	if p, ok := l.languagePrinter[lang]; ok {
		return p.Sprintf(key, args...)
	}
	log.Warn(ctx, "translate printer not found", "key", key, "lang", lang)

	return key
}

func (l *multilingual) getDefaultLang() LanguageType {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.defaultLang
}

func (l *multilingual) getLanguages() []LanguageType {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.languages
}

// loadFromFS load from file system
func (l *multilingual) loadFromFS(ctx context.Context, fsys fs.FS) (map[LanguageType]map[string]struct{}, error) {

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

func (l *multilingual) loadTranslations(ctx context.Context, lang LanguageType, fsys fs.FS) (map[string]struct{},
	error) {
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
			if tag == language.Make(string(l.defaultLang)) {
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

func newMultilingualManager(ctx context.Context, opts Options) (*multilingual, error) {
	if opts.DefaultLang == "" {
		opts.DefaultLang = DefaultLanguage
	}

	l := &multilingual{
		languagePrinter: make(map[LanguageType]*message.Printer),
		languages:       getAllLanguages(),
		defaultLang:     opts.DefaultLang,
		builder:         catalog.NewBuilder(catalog.Fallback(language.Make(string(opts.DefaultLang)))),
	}

	// load different sources
	paths := make([]string, 0, len(opts.languageDir)+1)
	sources := make([]fs.FS, 0, len(opts.languageDir)+1)
	if !isEmbedFSEmpty(embedFS) {
		sub, err := fs.Sub(embedFS, translationsRoot)
		if err != nil {
			log.Error(ctx, "fs.Sub on embed failed", "dir", translationsRoot, log.E(err))
			return nil, err
		}
		sources = append(sources, sub)
		paths = append(paths, "embed path")
	}

	if opts.languageDir != "" {
		diskFs := os.DirFS(opts.languageDir)
		sources = append(sources, diskFs)
		paths = append(paths, opts.languageDir)
	}

	languageKeyMap := make(map[LanguageType]map[string]struct{})
	for idx, src := range sources {
		fileLanguageKeyMap, err := l.loadFromFS(ctx, src)
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
	for _, lang := range l.getLanguages() {
		l.languagePrinter[lang] = message.NewPrinter(language.Make(string(lang)), message.Catalog(l.builder))
	}

	return l, nil

}

// NewI18nManager return a new i18n manager
func NewI18nManager(ctx context.Context, opts Options) (Interface, error) {
	baseLangManager, err := newMultilingualManager(ctx, opts)
	if err != nil {
		log.Error(ctx, "new base language manager failed", log.E(err))
		return nil, err
	}

	i := &i18n{
		languageDir:           opts.languageDir,
		languageBaseInterface: baseLangManager,
	}

	return i, nil
}

// langCtxKey define lang context key for translator
type langCtxKey struct{}

var langKey langCtxKey

// GetLangFromCtx get translator from context
func GetLangFromCtx(ctx context.Context) LanguageType {
	if v := ctx.Value(langKey); v != nil {
		if l, ok := v.(LanguageType); ok {
			return l
		}
	}
	return DefaultLanguage
}

// ContextWithLang set tag to context
func ContextWithLang(ctx context.Context, lang LanguageType) context.Context {
	return context.WithValue(ctx, langKey, lang)
}

// Options i18n options
type Options struct {
	// DefaultLang If the required language or key does not exist, the default language will be used.
	DefaultLang LanguageType
	// languageDir Direct files to be loaded
	languageDir string
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

func isEmbedFSEmpty(fs embed.FS) bool {
	entries, err := fs.ReadDir(".")
	if err != nil {
		return true
	}
	return len(entries) == 0
}
