/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package language

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"configcenter/src/common/blog"
)

// ccErrorHelper CC 错误处理接口的实现
type ccLanguageHelper struct {
	lang map[string]LanguageMap // the key is a language code, for：en,cn,jp,us etc.
}

// CreateDefaultCCLanguageIf create the default cc error interface instance
func (cli *ccLanguageHelper) CreateDefaultCCLanguageIf(language string) DefaultCCLanguageIf {
	return &ccDefaultLanguageHelper{
		languageType: language,
		languageStr:  cli.Language,
		languageStrf: cli.Languagef,
	}
}

// Language returns an language that adapt to the Language interface which not accepts arguments
func (cli *ccLanguageHelper) Language(language string, key string) string {

	return cli.languageStr(language, key)

}

// Languagef returns an langauge that adapt to the language interface which accepts arguments
func (cli *ccLanguageHelper) Languagef(language string, key string, args ...interface{}) string {
	return cli.languageStrf(language, key, args...)
}

// load load language package file from dir
func (cli *ccLanguageHelper) Load(lang map[string]LanguageMap) {
	// blog.V(3).Infof("loaded language resource: %#v", lang)
	cli.lang = lang
}

// LoadLanguageResourceFromDir  load language resource from file
func LoadLanguageResourceFromDir(dir string) (map[string]LanguageMap, error) {
	blog.Infof("loading language from %s\n", dir)
	// read all language file from dir
	var langMap = map[string]LanguageMap{}
	walkerr := filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		if nil != err {
			return err
		}

		if f.IsDir() {
			return nil
		}

		if !strings.HasSuffix(f.Name(), ".json") {
			return nil
		}

		items := strings.Split(path, string(os.PathSeparator))
		language := items[len(items)-2 : len(items)-1]

		// analysis language package file
		data, rerr := ioutil.ReadFile(path)
		if nil != rerr {
			return fmt.Errorf("read language file %v, error: %v", path, rerr)
		}

		res := LanguageMap{}
		jsErr := json.Unmarshal(data, &res)
		if nil != jsErr {
			return fmt.Errorf("unmarshal language file %v, error: %v", path, jsErr)
		}

		// will not check language validity, subject to language package file
		for key, val := range res {
			if _, ok := langMap[language[0]][key]; ok && key != "" {
				fmt.Printf("the language code[%s] repeated\n", key)
			}

			if nil == langMap[language[0]] {
				langMap[language[0]] = LanguageMap{}
			}
			langMap[language[0]][key] = val
		}
		return nil

	})
	// blog.Infof("loaded language from dir %v", langMap)

	if walkerr != nil {
		return nil, walkerr
	}

	return langMap, nil
}

func (cli *ccLanguageHelper) GetLang() map[string]LanguageMap {
	return cli.lang
}

// getLanguageKey get error code manager
func (cli *ccLanguageHelper) getLanguageKey(language string) LanguageMap {
	codemgr, ok := cli.lang[language]
	if !ok && language != defaultLanguage {
		// when the specified language not found, find it from default language package
		codemgr, ok = cli.lang[defaultLanguage]
		if !ok {
			return nil
		}
	}
	return codemgr
}

// getLanguageStr get errors string interface
func (cli *ccLanguageHelper) getLanguageStr(codemgr LanguageMap, key string) string {

	errstr, errOk := codemgr[key]
	if !errOk {
		// when the specified language not found, find it from default language package
		codemgr = cli.getLanguageKey(defaultLanguage)
		if nil != codemgr {
			errstr = codemgr[key]
		}
	}

	return errstr
}

var replayHolderReg = regexp.MustCompile(`\[(.*?)\]`)

// errorStr 错误码转换成错误信息，此方法适合不需要动态填充参数的错误信息
func (cli *ccLanguageHelper) languageStr(language, key string) string {

	// find language package form resource cache
	codemgr := cli.getLanguageKey(language)

	if nil == codemgr {
		return fmt.Sprintf(UnknownTheLanguageStrf, language)
	}

	ms := replayHolderReg.FindAllString(key, -1)
	// blog.Infof("key %s match %v", key, ms)
	if len(ms) > 0 {
		fmt.Printf("ms: %s\n", ms)
		key = replayHolderReg.ReplaceAllString(key, "[]")
		fmt.Printf("key: {%s}\n", key)
		text := cli.getLanguageStr(codemgr, key)
		if text != "" {
			mm := []interface{}{}
			for _, s := range ms {
				mm = append(mm, strings.TrimSuffix(strings.TrimPrefix(s, "["), "]"))
			}
			return fmt.Sprintf(text, mm...)
		}
	}
	// find error string from language language package
	return cli.getLanguageStr(codemgr, key)
}

// errorStrf retruns the error message string by code within language, format should define within error resource file
func (cli *ccLanguageHelper) languageStrf(language, key string, args ...interface{}) string {

	// find language from resource cache
	codemgr := cli.getLanguageKey(language)
	if nil == codemgr {
		return fmt.Sprintf(UnknownTheLanguageStrf, language)
	}

	// find error string within the language
	errstr := cli.getLanguageStr(codemgr, key)

	// format outputs, format args should define within resource file,
	// we will not check validity the formate here
	return fmt.Sprintf(errstr, args...)
}
