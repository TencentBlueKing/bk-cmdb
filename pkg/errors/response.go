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

// Detail detail error info
type Detail struct {
	// internal custom error code
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

// RespError response error info for out layer
type RespError struct {
	// Code for show
	Code ErrorCodeType `json:"code"`
	// Message for show, can be translated
	Message string   `json:"message"`
	System  string   `json:"system,omitempty"`
	Details []Detail `json:"details,omitempty"`
	Data    any      `json:"data,omitempty"`
	// DetailError if existed, unwrap error for details
	DetailError error `json:"-"`
}

// HttpErrorManager error manager
type HttpErrorManager struct {
	System string
}

// NewErrorManager return error manager with system
func NewErrorManager(system string) *HttpErrorManager {
	return &HttpErrorManager{
		System: system,
	}
}

// ErrorResponseHandler default error interfaces
type ErrorResponseHandler interface {
	ConvToRespError(err error, opts ...convOpt) *RespError
	NewRespError(code ErrorCodeType, data ...any) *RespError
	UnwrapDetails(err error) []Detail
	WrapValidationErrors(err error) error
}

// convOpt convert error to response error option
type convOpt func(re *RespError, srcErr error)

// ConvToRespError convert error to response error with convert options
func (m *HttpErrorManager) ConvToRespError(err error, opts ...convOpt) *RespError {
	if err == nil {
		return nil
	}

	// add system or detail info
	var re *RespError
	if errors.As(err, &re) {
		if re.System == "" {
			re.System = m.System
		}
		if len(re.Details) == 0 && re.DetailError != nil {
			re.Details = m.UnwrapDetails(re.DetailError)
		}
		for _, opt := range opts {
			opt(re, err)
		}
		return re
	}

	code := UNKNOWN
	msg := err.Error()
	var codeErr CodeError
	if errors.As(err, &codeErr) {
		code = ErrorCodeType(codeErr.GetCode())
		msg = codeErr.Error()
	}

	re = &RespError{
		Code:        code,
		Message:     msg,
		System:      m.System,
		Details:     m.UnwrapDetails(err),
		DetailError: err,
	}

	for _, opt := range opts {
		opt(re, err)
	}
	return re
}

// WithCode set code for response error
func WithCode(code ErrorCodeType) convOpt {
	return func(re *RespError, _ error) {
		re.Code = code
	}
}

// WithMessage set message for response error
func WithMessage(msg string) convOpt {
	return func(re *RespError, _ error) {
		re.Message = msg
	}
}

// WithData set data for response error
func WithData(vals ...any) convOpt {
	return func(re *RespError, _ error) {
		re.Data = getValues(vals...)
	}
}

// WithDetailErr set detail error for response error
func WithDetailErr(detailErr error) convOpt {
	return func(re *RespError, _ error) {
		re.DetailError = detailErr
	}
}

// NewRespError new response error
func (m *HttpErrorManager) NewRespError(code ErrorCodeType, data ...any) *RespError {
	return &RespError{
		Code:   code,
		System: m.System,
		Data:   getValues(data...),
	}
}

// CodeError return response error message
func (r *RespError) Error() string {
	return r.Message
}

// UnwrapDetails unwrap error for details, which is nested or joined
func (m *HttpErrorManager) UnwrapDetails(err error) []Detail {
	if err == nil {
		return nil
	}

	type (
		singleUnwrapper interface{ Unwrap() error }
		multiUnwrapper  interface{ Unwrap() []error }
	)

	var details []Detail

	queue := []error{err}
	visited := make(map[error]struct{})
	for i := 0; i < len(queue); i++ {
		e := queue[i]
		if e == nil {
			continue
		}
		if _, ok := visited[e]; ok {
			continue
		}
		visited[e] = struct{}{}

		d := Detail{Message: e.Error()}
		var re *RespError
		var codeErr CodeError
		if errors.As(e, &re) {
			d.Code = string(re.Code)
		} else if errors.As(e, &codeErr) {
			d.Code = codeErr.GetCode()
		}

		details = append(details, d)

		switch uw := any(e).(type) {
		case multiUnwrapper:
			children := uw.Unwrap()
			for _, child := range children {
				if child != nil {
					queue = append(queue, child)
				}
			}
			details = details[0 : len(details)-1]
		case singleUnwrapper:
			if child := uw.Unwrap(); child != nil {
				queue = append(queue, child)
			}
		}
	}

	return details
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
		return &multiValidationErr{children: children}
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
