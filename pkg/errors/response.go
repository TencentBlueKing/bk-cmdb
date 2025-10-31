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

package cerr

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

var errorManager ErrorResponseHandler

// SetDefaultErrorManager set default error manager
func SetDefaultErrorManager(m *HttpErrorManager) {
	errorManager = m
}

// GetDefaultErrorManager GetDefaultManager get default error manager
func GetDefaultErrorManager() ErrorResponseHandler {
	return errorManager
}

// RespError response error info for out layer
type RespError struct {
	// Code for show
	Code ErrorCode `json:"code"`
	// Message for show, can be translated
	Message string   `json:"message"`
	Details []string `json:"details,omitempty"`
	Data    any      `json:"data,omitempty"`
	// DetailError if existed, unwrap error for details
	DetailError error `json:"-"`
}

// HttpErrorManager error manager
type HttpErrorManager struct{}

// managerOpt for new error manager client options
type managerOpt func(re *HttpErrorManager)

// NewErrorManager return error manager
func NewErrorManager(opts ...managerOpt) *HttpErrorManager {
	manager := &HttpErrorManager{}
	for _, opt := range opts {
		opt(manager)
	}
	return manager
}

// ErrorResponseHandler default error interfaces
type ErrorResponseHandler interface {
	ConvToRespError(err error, opts ...convOpt) *RespError
	NewRespError(code ErrorCode, data ...any) *RespError
	UnwrapDetails(err error) []string
	WrapValidationErrors(err error) error
}

// convOpt convert error to response error option
type convOpt func(re *RespError)

// ConvToRespError convert error to response error with convert options
func (m *HttpErrorManager) ConvToRespError(err error, opts ...convOpt) *RespError {
	if err == nil {
		return nil
	}

	// add system or detail info
	var re *RespError
	if errors.As(err, &re) {
		if len(re.Details) == 0 && re.DetailError != nil {
			re.Details = m.UnwrapDetails(re.DetailError)
		}
		for _, opt := range opts {
			opt(re)
		}
		return re
	}

	code := UNKNOWN
	var codeErr CodeError
	if errors.As(err, &codeErr) {
		code = codeErr.GetCode()
	}

	re = &RespError{
		Code:        code,
		Details:     m.UnwrapDetails(err),
		DetailError: err,
	}

	for _, opt := range opts {
		opt(re)
	}
	return re
}

// WithCode set code for response error
func WithCode(code ErrorCode) convOpt {
	return func(re *RespError) {
		re.Code = code
	}
}

// WithMessage set message for response error
func WithMessage(msg string) convOpt {
	return func(re *RespError) {
		re.Message = msg
	}
}

// WithData set data for response error
func WithData(vals ...any) convOpt {
	return func(re *RespError) {
		re.Data = getValues(vals...)
	}
}

// WithDetailErr set detail error for response error
func WithDetailErr(detailErr error) convOpt {
	return func(re *RespError) {
		re.DetailError = detailErr
	}
}

// NewRespError new response error
func (m *HttpErrorManager) NewRespError(code ErrorCode, data ...any) *RespError {
	return &RespError{
		Code: code,
		Data: getValues(data...),
	}
}

// CodeError return response error message
func (r *RespError) Error() string {
	return r.Message
}

// GetCode return error code
func (r *RespError) GetCode() string {
	return r.Message
}

type (
	multiUnwrapper interface{ Unwrap() []error }
)

// UnwrapDetails unwrap error for details, which is  joined
func (m *HttpErrorManager) UnwrapDetails(err error) []string {
	if err == nil {
		return nil
	}
	return getDetails(err)
}

func getDetails(err error) []string {
	if err == nil {
		return nil
	}

	if uw, ok := err.(multiUnwrapper); ok {
		var out []string
		for _, child := range uw.Unwrap() {
			out = append(out, getDetails(child)...)
		}
		return out
	}

	return []string{err.Error()}
}

// WrapValidationErrors wrap validation errors
func (m *HttpErrorManager) WrapValidationErrors(err error) error {
	if err == nil {
		return nil
	}
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		children := make([]error, 0, len(ve))
		for _, fe := range ve {
			children = append(children, &fieldErr{fieldE: fe})
		}
		return errors.Join(children...)
	}
	return err
}

func getValues(vals ...any) any {
	if len(vals) == 0 {
		return nil
	}
	if len(vals) == 1 {
		return vals[0]
	}
	return vals
}
