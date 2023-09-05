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

package operator

import (
	"configcenter/pkg/excel"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/language"
	"configcenter/src/web_server/service/excel/core"
)

// BaseOp excel base operator
type BaseOp struct {
	excel    *excel.Excel
	client   *core.Client
	objID    string
	kit      *rest.Kit
	language language.CCLanguageIf
}

type BuildOpFunc func(op *BaseOp) error

// NewBaseOp create a base operator
func NewBaseOp(opts ...BuildOpFunc) (*BaseOp, error) {
	op := new(BaseOp)
	for _, opt := range opts {
		if err := opt(op); err != nil {
			return nil, err
		}
	}

	return op, nil
}

// FilePath set operator file path
func FilePath(filePath string) BuildOpFunc {
	return func(op *BaseOp) error {
		var err error
		op.excel, err = excel.NewExcel(excel.FilePath(filePath), excel.OpenOrCreate())
		if err != nil {
			return err
		}

		return nil
	}
}

// Client set client
func Client(client *core.Client) BuildOpFunc {
	return func(op *BaseOp) error {
		op.client = client
		return nil
	}
}

// ObjID set operator object id
func ObjID(objID string) BuildOpFunc {
	return func(op *BaseOp) error {
		op.objID = objID
		return nil
	}
}

// Kit set operator kit
func Kit(kit *rest.Kit) BuildOpFunc {
	return func(op *BaseOp) error {
		op.kit = kit
		return nil
	}
}

// Language set operator language
func Language(language language.CCLanguageIf) BuildOpFunc {
	return func(op *BaseOp) error {
		op.language = language
		return nil
	}
}

// GetExcel get excel
func (op *BaseOp) GetExcel() *excel.Excel {
	return op.excel
}

// GetClient get client
func (op *BaseOp) GetClient() *core.Client {
	return op.client
}

// GetObjID get objID
func (op *BaseOp) GetObjID() string {
	return op.objID
}

// GetKit get kit
func (op *BaseOp) GetKit() *rest.Kit {
	return op.kit
}

// GetLang get language
func (op *BaseOp) GetLang() language.CCLanguageIf {
	return op.language
}
