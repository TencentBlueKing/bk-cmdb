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

package errors

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"configcenter/src/common/blog"
)

// ccErrorHelper CC 错误处理接口的实现
type ccErrorHelper struct {
	errCode map[string]ErrorCode // the key is a language code, for：en,cn,jp,us etc.
}

// CreateDefaultCCErrorIf create the default cc error interface instance
func (cli *ccErrorHelper) CreateDefaultCCErrorIf(language string) DefaultCCErrorIf {
	return &ccDefaultErrorHelper{
		language:  language,
		errorStr:  cli.Error,
		errorStrf: cli.Errorf,
	}
}

// Errorf returns an error that adapt to the error interface which not accepts arguments
func (cli *ccErrorHelper) Error(language string, errCode int) error {
	return &ccError{code: errCode, callback: func() string {
		return cli.errorStr(language, errCode)
	}}
}

// Errorf returns an error that adapt to the error interface which accepts arguments
func (cli *ccErrorHelper) Errorf(language string, ErrorCode int, args ...interface{}) error {
	return &ccError{code: ErrorCode, callback: func() string {
		return cli.errorStrf(language, ErrorCode, args...)
	}}
}

// load load language package file from dir
func (cli *ccErrorHelper) Load(errcode map[string]ErrorCode) {
	// blog.V(3).Infof("loaded error resource: %#v", errcode)
	cli.errCode = errcode
}

func LoadErrorResourceFromDir(dir string) (map[string]ErrorCode, error) {
	// read all language file from dir
	var errCode = map[string]ErrorCode{}
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

		// analysis error package file
		fmt.Printf("loading error resource from %s\n", path)
		data, rerr := ioutil.ReadFile(path)
		if nil != rerr {
			return rerr
		}

		res := ErrorCode{}
		jsErr := json.Unmarshal(data, &res)
		if nil != jsErr {
			blog.Errorf("LoadErrorResourceFromDir error: %v, file: %s", jsErr, path)
			return jsErr
		}

		// will not check language validity, subject to language package file
		for key, val := range res {
			if _, ok := errCode[language[0]][key]; ok && key != "" {
				fmt.Printf("the error code[%s] repeated\n", key)
			}

			if nil == errCode[language[0]] {
				errCode[language[0]] = ErrorCode{}
			}
			errCode[language[0]][key] = val
		}

		return nil

	})

	if walkerr != nil {
		return nil, walkerr
	}

	return errCode, nil
}

func (cli *ccErrorHelper) GetErrorCode() map[string]ErrorCode {
	return cli.errCode
}

// getErrorCode get error code manager
func (cli *ccErrorHelper) getErrorCode(language string) ErrorCode {
	codemgr, ok := cli.errCode[language]
	if !ok && language != defaultLanguage {
		// when the specified language not found, find it from default language package
		codemgr, ok = cli.errCode[defaultLanguage]
		if !ok {
			return nil
		}
	}
	return codemgr
}

// getErrorString get errors string interface
func (cli *ccErrorHelper) getErrorStr(codemgr ErrorCode, errCode int) string {

	reset := false
RESET:
	errstr, errOk := codemgr[strconv.Itoa(errCode)]
	if !errOk {
		// when the specified language not found, find it from default language package
		if !reset {
			reset = true
			codemgr = cli.getErrorCode(defaultLanguage)
			if nil != codemgr {
				goto RESET
			}
		}
		return fmt.Sprintf(UnknownTheCodeStrf, errCode)
	}

	return errstr
}

// errorStr 错误码转换成错误信息，此方法适合不需要动态填充参数的错误信息
func (cli *ccErrorHelper) errorStr(language string, errCode int) string {

	// find language package form resource cache
	codemgr := cli.getErrorCode(language)
	if nil == codemgr {
		return fmt.Sprintf(UnknownTheLanguageStrf, language)
	}

	// find error string from language language package
	return cli.getErrorStr(codemgr, errCode)
}

// errorStrf retruns the error message string by code within language, format should define within error resource file
func (cli *ccErrorHelper) errorStrf(language string, errCode int, args ...interface{}) string {

	// find language from resource cache
	codemgr := cli.getErrorCode(language)
	if nil == codemgr {
		return fmt.Sprintf(UnknownTheLanguageStrf, language)
	}

	// find error string within the language
	errstr := cli.getErrorStr(codemgr, errCode)

	// format outputs, format args should define within resource file,
	// we will not check validity the formate here
	return fmt.Sprintf(errstr, args...)
}
