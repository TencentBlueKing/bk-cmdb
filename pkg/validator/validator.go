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

// Package validator use for valiate struct
package validator

import (
	"context"
	"errors"
	"reflect"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	validator "github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"github.com/samber/lo"

	"github.com/TencentBlueKing/bk-cmdb/pkg/util"
)

var (
	validate     *validator.Validate
	defaultTrans ut.Translator
	uni          *ut.UniversalTranslator
)

// ValidationError 校验错误
type ValidationError struct {
	rawErr error
}

// Validator 实现了 Validate 接口自定义调用
type Validator interface {
	Validate(ctx context.Context) error
}

func (e *ValidationError) getTranslator() ut.Translator {
	return defaultTrans
}

// Error error iface
func (e *ValidationError) Error() string {
	errs := e.Unwrap()
	if len(errs) == 0 {
		return e.rawErr.Error()
	}

	// 只返回单个错误
	return errs[0].Error()
}

// Unwrap multierr iface
func (e *ValidationError) Unwrap() []error {
	if _, ok := e.rawErr.(*validator.InvalidValidationError); ok {
		return []error{e.rawErr}
	}

	ves, ok := e.rawErr.(validator.ValidationErrors)
	if !ok {
		return []error{e.rawErr}
	}

	errs := make([]error, 0, len(ves))
	for _, ve := range ves {
		errs = append(errs, errors.New(ve.Translate(e.getTranslator())))
	}

	return errs
}

// Struct 通过 validate tag 校验结构体, Validate 校验需要传入指针类型
func Struct(ctx context.Context, s any) error {
	err := validate.StructCtx(ctx, s)
	if err != nil {
		return &ValidationError{err}
	}

	// 实现了 Validate 接口自定义调用
	if v, ok := s.(Validator); ok {
		return v.Validate(ctx)
	}

	return nil
}

// readableTagName 返回可读的json/req校验字段名称, 唯一性由codec校验
func readableTagName(field reflect.StructField) string {
	name := util.GetTagName(field, "json")
	if name != "" && name != "-" {
		return name
	}

	name = util.GetTagName(field, "req")
	if name != "" && name != "-" {
		return name
	}

	return ""
}

func init() {
	validate = validator.New(validator.WithRequiredStructEnabled())
	validate.RegisterTagNameFunc(readableTagName)

	// 默认使用英文
	en := en.New()
	uni = ut.New(en, en)
	defaultTrans, _ = uni.GetTranslator("en")
	lo.Must0(en_translations.RegisterDefaultTranslations(validate, defaultTrans))
}
