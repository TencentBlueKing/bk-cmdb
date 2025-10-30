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
	"github.com/go-playground/validator/v10"
)

type fieldErr struct {
	fieldE validator.FieldError
}

// Error return error string
func (e *fieldErr) Error() string {
	return e.fieldE.Error()
}

// GetCode return error code
func (e *fieldErr) GetCode() string {
	return "ValidationError"
}

type multiValidationErr struct {
	children []error
}

// Error return error string
func (e *multiValidationErr) Error() string {
	return "validation error"
}

// Unwrap return children errors
func (e *multiValidationErr) Unwrap() []error { return e.children }
