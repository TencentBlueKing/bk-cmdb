/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 THL A29 Limited,
 * a Tencent company. All rights reserved.
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

package excel

import (
	"fmt"

	"configcenter/src/common/util"

	"github.com/xuri/excelize/v2"
)

// FieldType field type
type FieldType string

const (
	Decimal FieldType = "decimal"
	Bool    FieldType = "bool"
	Enum    FieldType = "enum"
)

// ValidationParam validation parameter
type ValidationParam struct {
	Type   FieldType
	Sqref  string
	Option interface{}
}

const (
	fieldTypeBoolTrue  = "true"
	fieldTypeBoolFalse = "false"
	enumRefSuffix      = "!$A:$A"
	errTitle           = "警告"
	errMessage         = "此值与此单元格定义的数据验证限制不匹配。"
)

func newValidation(param *ValidationParam) (*excelize.DataValidation, error) {
	validation := excelize.NewDataValidation(true)
	validation.SetSqref(param.Sqref)
	switch param.Type {
	case Decimal:
		validation.Type = string(Decimal)
	case Bool:
		if err := validation.SetDropList([]string{fieldTypeBoolTrue, fieldTypeBoolFalse}); err != nil {
			return nil, err
		}
	case Enum:
		ref, err := getEnumRefDropList(util.GetStrByInterface(param.Option))
		if err != nil {
			return nil, err
		}
		validation.SetSqrefDropList(ref)
	}

	validation.SetError(excelize.DataValidationErrorStyleStop, errTitle, errMessage)

	return validation, nil
}

func getEnumRefDropList(sheet string) (string, error) {
	return fmt.Sprintf("'%s'%s", sheet, enumRefSuffix), nil
}
