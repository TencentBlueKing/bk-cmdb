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
	"maps"
	"os"
	"sync"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"golang.org/x/text/message/catalog"

	"github.com/TencentBlueKing/bk-cmdb/pkg/constant"
	"github.com/TencentBlueKing/bk-cmdb/pkg/errors"
	"github.com/TencentBlueKing/bk-cmdb/pkg/kit"
	"github.com/TencentBlueKing/bk-cmdb/pkg/log"
)

// Interface i18n interface for multilingual internationalization, starting from the scenario, it can be divided
// into two types: implementing error translation and built-in system translation.
type Interface interface {
	// RespError translate error info, translate for error message by error code which is pre-determined
	RespError(kt *kit.Kit, err error) *cerr.RespError
	// Sys translate key, return key if not found
	Sys(kt *kit.Kit, key string, args ...any) string
	// Validate get tag from request and check if it is supported language
	Validate(lang LanguageType) error
	// DefaultLang get default language if lang is not set
	DefaultLang() LanguageType
	// Reload languages from directory for dynamic changes, reloading the translation content will overwrite the
	// original translation content. The translation content and the default language used can be updated dynamically.
	// If a language directory is provided, it will be loaded; if no language directory is provided, the initial
	// configuration directory will be loaded.
	Reload(kt *kit.Kit, opt *Options) error
}

// embedDir embeds static resource files for internationalization.
//
// The embedded file system includes all files under the 'resource' directory at compile time. Paths are prefixed with
// "resource/" when accessing files.
//
// Usage example:
//
//	file, err := embedDir.ReadFile("resource/en/error.json")
//
// Important notes:
// - Only files in the 'resource' directory are embedded (subdirectories included)
// - Paths are case-sensitive and must match exactly
// - Embedded content is read-only at runtime
//
//go:embed resource
var embedDir embed.FS

var (
	i18nTranslator Interface
	setOnce        sync.Once
)

// initTranslator set default i18n manager
func initTranslator(m Interface) {
	setOnce.Do(func() { i18nTranslator = m })
}

// RespError translate error info, translate for error message by error code which is pre-determined
func RespError(kt *kit.Kit, err error) *cerr.RespError {
	return i18nTranslator.RespError(kt, err)
}

// Sys translate key, return key if not found
func Sys(kt *kit.Kit, key string, args ...any) string {
	return i18nTranslator.Sys(kt, key, args...)
}

// Validate get tag from request and check if it is supported language
func Validate(lang LanguageType) error {
	return i18nTranslator.Validate(lang)
}

// DefaultLang get default language if lang is not set
func DefaultLang() LanguageType {
	return i18nTranslator.DefaultLang()
}

// Reload languages from directory for dynamic changes
func Reload(kt *kit.Kit, opt *Options) error {
	return i18nTranslator.Reload(kt, opt)
}

// renders stores language renders for different domains
type renders struct {
	sys *message.Printer
	err *message.Printer
}

// builders stores language builders for different domains
type builders struct {
	sys *catalog.Builder
	err *catalog.Builder
}

// translations stores translations for different domains
type translations struct {
	// sys built-in common translation content
	sys map[string]string
	// err built-in error translation content
	err map[string]string
}

// i18n define implementation required components for language package.
type i18n struct {
	// languageDir file path for loading language, it requires folders for various language packs, the naming of
	// language pack folders must comply with the system language definition specification.
	languageDir string
	// languages stores all supported languages.
	languages map[LanguageType]struct{}
	// defaultLang when no language is provided in the request, the default language is used. The default language
	// supports dynamic loading.
	defaultLang LanguageType
	// render stores the render for each supported language, use render to translate.
	render map[LanguageType]renders
	// lock for hot update of language configuration or content.
	lock sync.RWMutex
}

// initTranslations initialize the i18n, register translations and build language render.
func (i *i18n) initTranslations(ctx context.Context, opts *Options) error {
	// load translations
	translations, err := i.loadTranslations(ctx, opts.LanguageDir)
	if err != nil {
		return err
	}

	// verify translations
	if err := i.verifyTranslations(ctx, opts, translations); err != nil {
		log.Error(ctx, "verify translations failed", log.E(err))
		return fmt.Errorf("verify translations failed: %w", err)
	}

	builders, err := i.registerTranslations(ctx, translations, opts)
	if err != nil {
		log.Error(ctx, "register translations failed", log.E(err))
		return fmt.Errorf("register translations failed: %w", err)
	}

	// build language translator
	i.buildTranslator(builders, translations, opts)
	return nil
}

// verifyTranslations check the legality of the loaded translation and default language.
func (i *i18n) verifyTranslations(ctx context.Context, opts *Options, trans map[LanguageType]translations) error {

	// check whether the default language is supported after load all languages
	if _, exist := trans[opts.DefaultLang]; !exist {
		log.Error(ctx, "verify default language not exist", "defaultLang", opts.DefaultLang)
		return fmt.Errorf("default language %s not exist", opts.DefaultLang)
	}

	defaultSys := trans[opts.DefaultLang].sys
	defaultErr := trans[opts.DefaultLang].err
	for lang, dom := range trans {
		if lang == opts.DefaultLang {
			continue
		}

		if !cmpKeyWithDefault(ctx, defaultSys, dom.sys) {
			log.Error(ctx, "verify load translations, sys keys not same with default lang", "defaultLang",
				opts.DefaultLang, "lang", lang)
			return fmt.Errorf("sys keys of lang %s key not same with default", lang)
		}

		if !cmpKeyWithDefault(ctx, defaultErr, dom.err) {
			log.Error(ctx, "error keys not same with default", "defaultLang",
				opts.DefaultLang, "lang", lang)
			return fmt.Errorf("error keys of lang %s key not same with default", lang)
		}

	}

	return nil
}

// Reload languages from file system for dynamic loading.
func (i *i18n) Reload(kt *kit.Kit, opt *Options) error {
	log.Info(kt, "start reload i18n translations", "languageDir", opt.LanguageDir, "defaultLang",
		opt.DefaultLang)

	if err := opt.validate(false); err != nil {
		log.Error(kt, "validate option failed", log.E(err))
		return fmt.Errorf("validate option failed: %w", err)
	}

	// set default language for reload scene
	if len(opt.DefaultLang) == 0 {
		opt.DefaultLang = i.DefaultLang()
	}
	if err := i.initTranslations(kt, opt); err != nil {
		log.Error(kt, "reload i18n translations failed", "languageDir", opt.LanguageDir, "defaultLang",
			opt.DefaultLang, log.E(err))
		return err
	}

	log.Info(kt, "reload i18n translations success", "languageDir", opt.LanguageDir, "defaultLang",
		opt.DefaultLang)
	return nil
}

func (i *i18n) loadTranslations(ctx context.Context, langDir string) (map[LanguageType]translations, error) {
	sources := make([]fs.FS, 0)
	// step1: Append loading directory from embed file for default initialization content.
	if !isEmbedFSEmpty(embedDir) {
		sub, err := fs.Sub(embedDir, languageDir)
		if err != nil {
			log.Error(ctx, "fs.Sub on embed failed", "dir", languageDir, log.E(err))
			return nil, err
		}
		sources = append(sources, sub)
	}

	// step2: Append loading directory from the input folder directory for the translation content dynamically loaded.
	if len(langDir) > 0 {
		// If a translation directory is provided, load it from the translation directory.
		sources = append(sources, os.DirFS(langDir))
	} else if len(i.languageDir) > 0 {
		// If no translation directory is provided, load it from the default configuration directory.
		sources = append(sources, os.DirFS(i.languageDir))
	}

	// step3: In order load language key from file system, first load embed file, then load from the input directory.
	// Dynamically loaded content will replace the embed translation content.
	languageDomainMap := make(map[LanguageType]translations)
	for _, src := range sources {
		fileLangMap, err := i.readLangFs(ctx, src)
		if err != nil {
			log.Error(ctx, "load i18n from file system failed", "path", src, log.E(err))
			return nil, err
		}

		for lang, dom := range fileLangMap {
			item, ok := languageDomainMap[lang]
			if !ok {
				item = translations{
					sys: make(map[string]string),
					err: make(map[string]string),
				}
			}
			maps.Copy(item.sys, dom.sys)
			maps.Copy(item.err, dom.err)
			languageDomainMap[lang] = item
		}
	}

	return languageDomainMap, nil
}

func (i *i18n) registerTranslations(ctx context.Context, languageKeyMap map[LanguageType]translations, opts *Options) (
	*builders, error) {

	defaultLang := language.Make(string(opts.DefaultLang))
	builders := &builders{
		sys: catalog.NewBuilder(catalog.Fallback(defaultLang)),
		err: catalog.NewBuilder(catalog.Fallback(defaultLang)),
	}

	setTransValues := func(ctx context.Context, keyMap map[string]string, builder *catalog.Builder,
		lang LanguageType) error {

		tag := language.Make(string(lang))
		for k, v := range keyMap {
			if setErr := builder.SetString(tag, k, v); setErr != nil {
				log.Error(ctx, "set sys string failed", "key", k, "lang", tag, log.E(setErr))
				return setErr
			}
		}
		return nil
	}

	for lang, dom := range languageKeyMap {
		if err := setTransValues(ctx, dom.sys, builders.sys, lang); err != nil {
			return nil, err
		}

		if err := setTransValues(ctx, dom.err, builders.err, lang); err != nil {
			return nil, err
		}
	}
	return builders, nil
}

func (i *i18n) buildTranslator(builder *builders, translations map[LanguageType]translations, opts *Options) {
	tempLangRender := make(map[LanguageType]renders, len(translations))
	languages := make(map[LanguageType]struct{})
	for lang := range translations {
		tag := language.Make(string(lang))
		tempLangRender[lang] = renders{
			sys: message.NewPrinter(tag, message.Catalog(builder.sys)),
			err: message.NewPrinter(tag, message.Catalog(builder.err)),
		}
		languages[lang] = struct{}{}
	}

	i.lock.Lock()
	defer i.lock.Unlock()

	if len(opts.LanguageDir) > 0 {
		i.languageDir = opts.LanguageDir
	}
	i.render = tempLangRender
	i.languages = languages
	i.defaultLang = opts.DefaultLang
}

func (i *i18n) rend(lang LanguageType) (renders, bool) {
	i.lock.RLock()
	defer i.lock.RUnlock()
	if r, ok := i.render[lang]; ok {
		return r, true
	}
	return renders{}, false
}

// Init new i18n client.
func Init(ctx context.Context, opts *Options) error {
	if err := opts.validate(true); err != nil {
		return err
	}

	// set default lang
	if len(opts.DefaultLang) == 0 {
		opts.DefaultLang = constant.DefaultLanguage
	}
	// new i18n client
	i18n := &i18n{
		defaultLang: opts.DefaultLang,
		render:      make(map[LanguageType]renders),
		languageDir: opts.LanguageDir,
	}

	// initTranslations i18n client
	err := i18n.initTranslations(ctx, opts)
	if err != nil {
		log.Error(ctx, "initTranslations i18n client failed", log.E(err))
		return err
	}

	initTranslator(i18n)
	return nil
}

// Validate get tag from request and check if it is supported by system.
func (i *i18n) Validate(lang LanguageType) error {
	ok := i.isSupportedLang(lang)
	if !ok {
		e := fmt.Errorf("unsupported language: %s", lang)
		return e
	}

	return nil
}

// RespError translate response error message by error code
func (i *i18n) RespError(kt *kit.Kit, err error) *cerr.RespError {
	if err == nil {
		err = cerr.NewError(cerr.Unknown, "unknown error")
	}

	respErr := cerr.ErrorClient().ConvToRespError(err)
	lang := LanguageType(kt.Lang)

	if len(lang) == 0 {
		lang = i.DefaultLang()
	}

	if bundle, ok := i.rend(lang); ok && bundle.err != nil {
		respErr.Message = bundle.err.Sprintf(string(respErr.Code))
		return respErr
	}

	log.Warn(kt, "translate error rend not found", "code", respErr.Code, "lang", lang)
	respErr.Message = string(respErr.Code)

	return respErr
}

// Sys translate key, return key if not found
func (i *i18n) Sys(kt *kit.Kit, key string, args ...any) string {
	lang := LanguageType(kt.Lang)
	if len(lang) == 0 {
		lang = i.DefaultLang()
	}

	if bundle, ok := i.rend(lang); ok && bundle.sys != nil {
		return bundle.sys.Sprintf(key, args...)
	}

	log.Warn(kt, "translate sys rend not found", "key", key, "lang", lang)
	return key
}

func (i *i18n) isSupportedLang(lang LanguageType) bool {
	i.lock.RLock()
	defer i.lock.RUnlock()
	if _, exist := i.languages[lang]; exist {
		return true
	}
	return false
}

// DefaultLang get default language
func (i *i18n) DefaultLang() LanguageType {
	i.lock.RLock()
	defer i.lock.RUnlock()
	return i.defaultLang
}

// readLangFs load multilingual translation files from input file system.
func (i *i18n) readLangFs(ctx context.Context, fsys fs.FS) (map[LanguageType]translations, error) {

	entries, err := fs.ReadDir(fsys, ".")
	if err != nil {
		log.Error(ctx, "read fs root failed", log.E(err))
		return nil, err
	}

	out := make(map[LanguageType]translations)
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		dirName := e.Name()
		lang := LanguageType(dirName)

		sub, err := fs.Sub(fsys, dirName)
		if err != nil {
			log.Error(ctx, "fs.Sub failed for language dir", "langDir", dirName, log.E(err))
			return nil, err
		}

		errMap, err := i.readFile(ctx, sub, "error.json")
		if err != nil {
			log.Error(ctx, "read i18n json file failed", "dir", dirName, "file", "error.json", log.E(err))
			return nil, err
		}

		sysMap, err := i.readFile(ctx, sub, "sys.json")
		if err != nil {
			log.Error(ctx, "read i18n json file failed", "dir", dirName, "file", "sys.json", log.E(err))
			return nil, err
		}
		out[lang] = translations{
			sys: sysMap,
			err: errMap,
		}
	}

	return out, nil
}

func (i *i18n) readFile(ctx context.Context, sub fs.FS, filename string) (map[string]string, error) {

	fileContent, err := fs.ReadFile(sub, filename)
	if err != nil {
		log.Error(ctx, "read i18n json file failed", "file", filename, log.E(err))
		return nil, err
	}

	jsonMap := make(map[string]string)
	if unmarshalErr := json.Unmarshal(fileContent, &jsonMap); unmarshalErr != nil {
		log.Error(ctx, "unmarshal i18n json failed", "file", filename, log.E(unmarshalErr))
		return nil, unmarshalErr
	}

	return jsonMap, nil
}

// Options i18n options
type Options struct {
	// DefaultLang If the language or key does not exist, the default language will be used. Non-required fields, if
	// not filled in, use predefined default language. Used to set or dynamically load the default language.
	DefaultLang LanguageType
	// LanguageDir language directory to load language files, it is a required field. Used for initializing or
	// dynamically loading language directories.
	LanguageDir string
	// RequireExternalDir if RequireExternalDir is true, The input of LanguageDir cannot be empty, and it must comply
	// with the directory specifications. Otherwise, languageDir is empty. Todo after config is ready, remove it.
	RequireExternalDir bool
}

// validate options, the two scenarios are initialization and reloading the language. When initializing the language,
// the configuration language directory must be provided. Todo after config is ready, remove RequireExternalDir.
// During reloading, the language directory does not have to be provided. If a language directory is provided, it will
// be loaded; if not, the initial configuration directory will be loaded.
func (o Options) validate(isInit bool) error {
	hasLangDir := len(o.LanguageDir) > 0
	hasDefault := len(o.DefaultLang) > 0

	// initTranslations language scenarios
	if isInit {
		if o.RequireExternalDir != hasLangDir {
			return fmt.Errorf("language dir and require external dir are not match, languageDir: %s, "+
				"requireExternalDir: %v", o.LanguageDir, o.RequireExternalDir)
		}
	} else {
		// reload language scenarios
		if !hasLangDir && !hasDefault {
			return fmt.Errorf("validate reload lanuages, both languageDir and defaultLang are empty")
		}
	}

	if !hasLangDir {
		return nil
	}
	return validDirectory(o.LanguageDir)
}

func validDirectory(path string) error {
	s, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("language dir %s does not exist: %w", path, err)
		}
		return fmt.Errorf("get directory stat failed, err: %w", err)
	}
	if !s.IsDir() {
		return fmt.Errorf("language path %s is not a directory", path)
	}
	return nil
}

// cmpKeyWithDefault compare key with default language key
func cmpKeyWithDefault(ctx context.Context, defaultLang, lang map[string]string) bool {
	if len(defaultLang) != len(lang) {
		log.Error(ctx, "default lang key count not equal with lang", "defaultLangLen", len(defaultLang),
			"langLen", len(lang))
		return false
	}
	isPassed := true
	for k := range defaultLang {
		if _, ok := lang[k]; !ok {
			log.Error(ctx, "key in defaultLang not found in lang", "lang", lang, "key", k)
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
