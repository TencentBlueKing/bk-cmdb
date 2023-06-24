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

package fieldtmpl

import (
	"configcenter/src/ac/meta"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
)

// CompareFieldTemplateAttr compare field template attributes with object attributes.
func (s *service) CompareFieldTemplateAttr(cts *rest.Contexts) {
	opt := new(metadata.CompareFieldTmplAttrOption)
	if err := cts.DecodeInto(opt); err != nil {
		cts.RespAutoError(err)
		return
	}

	if rawErr := opt.Validate(); rawErr.ErrCode != 0 {
		cts.RespAutoError(rawErr.ToCCError(cts.Kit.CCError))
		return
	}

	// check if user has the permission of the field template
	// TODO add find object auth check too after find object operation authorization is supported
	if authResp, authorized := s.auth.Authorize(cts.Kit, meta.ResourceAttribute{Basic: meta.Basic{
		Type: meta.FieldTemplate, Action: meta.Find, InstanceID: opt.TemplateID}}); !authorized {
		cts.RespNoAuth(authResp)
		return
	}

	res, _, err := s.logics.FieldTemplateOperation().CompareFieldTemplateAttr(cts.Kit, opt, true)
	if err != nil {
		cts.RespAutoError(err)
		return
	}

	cts.RespEntity(res)
}

// CompareFieldTemplateUnique compare field template uniques with object uniques.
func (s *service) CompareFieldTemplateUnique(cts *rest.Contexts) {
	opt := new(metadata.CompareFieldTmplUniqueOption)
	if err := cts.DecodeInto(opt); err != nil {
		cts.RespAutoError(err)
		return
	}

	if rawErr := opt.Validate(); rawErr.ErrCode != 0 {
		cts.RespAutoError(rawErr.ToCCError(cts.Kit.CCError))
		return
	}

	// check if user has the permission of the field template
	// TODO add find object auth check too after find object operation authorization is supported
	if authResp, authorized := s.auth.Authorize(cts.Kit, meta.ResourceAttribute{Basic: meta.Basic{
		Type: meta.FieldTemplate, Action: meta.Find, InstanceID: opt.TemplateID}}); !authorized {
		cts.RespNoAuth(authResp)
		return
	}

	res, _, err := s.logics.FieldTemplateOperation().CompareFieldTemplateUnique(cts.Kit, opt, true)
	if err != nil {
		cts.RespAutoError(err)
		return
	}

	cts.RespEntity(res)
}
