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
	"sync"
)

// ErrorRespConvertor convert error to response error interface
type ErrorRespConvertor interface {
	// ConvToRespError convert error to response error
	ConvToRespError(err error, opts ...ConvOpt) *RespError
}

type (
	// httpRespErrorConvertor http response error convertor to covert error to response error
	httpRespErrorConvertor struct{}
	// ConvOpt convert error to response error option
	ConvOpt func(re *RespError)
	// clientOpt func(re *httpRespErrorConvertor) error
	clientOpt func(re *httpRespErrorConvertor) error
)

// Init init error client
func Init(opts ...clientOpt) error {
	client := &httpRespErrorConvertor{}
	for _, opt := range opts {
		if err := opt(client); err != nil {
			return err
		}
	}

	initErrorClient(client)
	return nil
}

// ConvToRespError convert error to response error with convert options
func (m *httpRespErrorConvertor) ConvToRespError(err error, opts ...ConvOpt) *RespError {
	if err == nil {
		err = NewError(Unknown, "unknown error")
	}

	// add system or detail info
	var re *RespError
	if errors.As(err, &re) {
		if len(re.Details) == 0 && re.DetailError != nil {
			re.Details = unwrapDetails(re.DetailError)
		}
		for _, opt := range opts {
			opt(re)
		}
		return re
	}

	code := Unknown
	var codeErr CodeError
	if errors.As(err, &codeErr) {
		code = codeErr.GetCode()
	}

	re = &RespError{
		Code:        code,
		Details:     unwrapDetails(err),
		DetailError: err,
	}

	for _, opt := range opts {
		opt(re)
	}
	return re
}

var (
	errRespConvertor ErrorRespConvertor
	setOnce          sync.Once
)

// setDefaultErrorManager set default error manager
func initErrorClient(m *httpRespErrorConvertor) {
	setOnce.Do(func() { errRespConvertor = m })
}

// ErrorClient GetDefaultManager get default error manager
func ErrorClient() ErrorRespConvertor {
	return errRespConvertor
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

// NewRespError new response error
func (m *httpRespErrorConvertor) NewRespError(code ErrorCode, data ...any) *RespError {
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
func (r *RespError) GetCode() ErrorCode {
	return r.Code
}

// WithCode set code for response error
func WithCode(code ErrorCode) ConvOpt {
	return func(re *RespError) {
		re.Code = code
	}
}

// WithMessage set message for response error
func WithMessage(msg string) ConvOpt {
	return func(re *RespError) {
		re.Message = msg
	}
}

// WithData set data for response error
func WithData(vals ...any) ConvOpt {
	return func(re *RespError) {
		re.Data = getValues(vals...)
	}
}

// WithDetailErr set detail error for response error
func WithDetailErr(detailErr error) ConvOpt {
	return func(re *RespError) {
		re.DetailError = detailErr
	}
}

type (
	multiUnwrapper interface{ Unwrap() []error }
)

// UnwrapDetails unwrap error for details, which is  joined
func unwrapDetails(err error) []string {
	if err == nil {
		return []string{}
	}
	return getDetails(err)
}

func getDetails(err error) []string {
	if err == nil {
		return []string{}
	}

	var re *ccError
	if errors.As(err, &re) {
		err = re.err
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

func getValues(vals ...any) any {
	if len(vals) == 0 {
		return nil
	}
	if len(vals) == 1 {
		return vals[0]
	}
	return vals
}
