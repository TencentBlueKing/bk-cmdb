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

// Package cerr support errors
package cerr

import (
	"errors"
	"fmt"
)

// CodeError interface for errors
type CodeError interface {
	Error() string
	GetCode() ErrorCode
}

// ccError cc error type for internal call
type ccError struct {
	code ErrorCode
	err  error
}

// CodeError implementation of errors interface
func (cli *ccError) Error() string {
	return cli.err.Error()
}

// GetCode returns errors code
func (cli *ccError) GetCode() ErrorCode {
	return cli.code
}

// Unwrap unwrap single error
func (cli *ccError) Unwrap() error {
	return cli.err
}

// NewError create new error with code and msg, use for internal error
func NewError(code ErrorCode, msg string) CodeError {
	return &ccError{
		code: code,
		err:  errors.New(msg),
	}
}

// Wrap error with code
func Wrap(code ErrorCode, err error) error {
	return &ccError{
		code: code,
		err:  err,
	}
}

// Errorf create new error with code and format msg, use for internal error
func Errorf(code ErrorCode, msg string, args ...any) error {
	return &ccError{
		code: code,
		err:  fmt.Errorf(msg, args...),
	}
}
