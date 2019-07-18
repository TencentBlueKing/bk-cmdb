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

import "configcenter/src/common"

// ccDefaultErrorHelper regular language code helper
type ccDefaultErrorHelper struct {
	language    string
	errorStr    func(language string, ErrorCode int) error
	errorStrf   func(language string, ErrorCode int, args ...interface{}) error
	ccErrorStr  func(language string, ErrorCode int) CCErrorCoder
	ccErrorStrf func(language string, ErrorCode int, args ...interface{}) CCErrorCoder
}

func (cli *ccDefaultErrorHelper) New(errorCode int, msg string) error {
	return &ccError{
		code: errorCode,
		callback: func() string {
			return msg
		},
	}
}

func NewCCError(errorCode int, msg string) CCErrorCoder {
	err := &ccError{
		code: errorCode,
		callback: func() string {
			return msg
		},
	}
	return err
}

var CCHttpError = &ccError{
	code: common.CCErrCommHTTPDoRequestFailed,
	callback: func() string {
		return "http request failed"
	},
}

// Error returns an error for specific language
func (cli *ccDefaultErrorHelper) Error(errCode int) error {
	return cli.errorStr(cli.language, errCode)
}

// Error returns an error with args for specific language
func (cli *ccDefaultErrorHelper) Errorf(errCode int, args ...interface{}) error {
	return cli.errorStrf(cli.language, errCode, args...)
}

// CCError returns an error for specific language
func (cli *ccDefaultErrorHelper) CCError(errCode int) CCErrorCoder {
	return cli.ccErrorStr(cli.language, errCode)
}

// CCError returns an error with args for specific language
func (cli *ccDefaultErrorHelper) CCErrorf(errCode int, args ...interface{}) CCErrorCoder {
	return cli.ccErrorStrf(cli.language, errCode, args...)
}
