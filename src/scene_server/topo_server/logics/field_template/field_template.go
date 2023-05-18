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

package field_template

import (
	"configcenter/src/apimachinery"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
)

// FieldTemplateOperation field template operation methods
type FieldTemplateOperation interface {
	CreateFieldTemplate(kit *rest.Kit, opt *metadata.CreateFieldTmplOption) (*metadata.RspID, error)
}

// NewFieldTemplateOperation create a new field quote operation instance
func NewFieldTemplateOperation(client apimachinery.ClientSetInterface) FieldTemplateOperation {
	return &fieldTemplate{
		clientSet: client,
	}
}

type fieldTemplate struct {
	clientSet apimachinery.ClientSetInterface
}

// CreateFieldTemplate create field template(contains field template brief information, attributes and uniques)
func (f *fieldTemplate) CreateFieldTemplate(kit *rest.Kit, opt *metadata.CreateFieldTmplOption) (
	*metadata.RspID, error) {

	res, err := f.clientSet.CoreService().FieldTemplate().CreateFieldTemplate(kit.Ctx, kit.Header, &opt.FieldTemplate)
	if err != nil {
		blog.Errorf("create field template failed, err: %v, data: %v, rid: %s", err, opt, kit.Rid)
		return nil, err
	}

	for idx := range opt.Attributes {
		opt.Attributes[idx].TemplateID = res.ID
	}
	attrIDs, err := f.clientSet.CoreService().FieldTemplate().CreateFieldTemplateAttrs(kit.Ctx, kit.Header,
		opt.Attributes)
	if err != nil {
		blog.Errorf("create field template attributes failed, err: %v, data: %v, rid: %s", err, opt.Attributes, kit.Rid)
		return nil, err
	}

	propertyIDToIDMap := make(map[string]int64)
	for idx, attr := range opt.Attributes {
		propertyIDToIDMap[attr.PropertyID] = attrIDs.IDs[idx]
	}

	uniques := make([]metadata.FieldTemplateUnique, len(opt.Uniques))
	for idx, uniqueOpt := range opt.Uniques {
		uniqueOpt.TemplateID = res.ID
		unique, ccErr := uniqueOpt.Convert(propertyIDToIDMap)
		if ccErr.ErrCode != 0 {
			return nil, ccErr.ToCCError(kit.CCError)
		}

		uniques[idx] = *unique
	}

	_, err = f.clientSet.CoreService().FieldTemplate().CreateFieldTemplateUniques(kit.Ctx, kit.Header, uniques)
	if err != nil {
		blog.Errorf("create field template uniques failed, err: %v, data: %v, rid: %s", err, uniques, kit.Rid)
		return nil, err
	}

	// todo 添加审计

	return res, nil
}
