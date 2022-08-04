/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.,
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the ",License",); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an ",AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package service

import (
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
)

// GetSystemUserConfig TODO
func (s *coreService) GetSystemUserConfig(ctx *rest.Contexts) {
	ctx.RespEntityWithError(s.core.SystemOperation().GetSystemUserConfig(ctx.Kit))
}

// SearchConfigAdmin TODO
func (s *coreService) SearchConfigAdmin(ctx *rest.Contexts) {
	conf, err := s.core.SystemOperation().SearchConfigAdmin(ctx.Kit)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(conf)
}

// SearchPlatformSettingConfig search platform setting.
func (s *coreService) SearchPlatformSettingConfig(ctx *rest.Contexts) {
	conf, err := s.core.SystemOperation().SearchPlatformSettingConfig(ctx.Kit)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(conf)
}

// UpdatePlatformSetting update platform setting.
func (s *coreService) UpdatePlatformSetting(ctx *rest.Contexts) {

	input := new(metadata.PlatformSettingConfig)
	if jsErr := ctx.DecodeInto(&input); nil != jsErr {
		ctx.RespAutoError(jsErr)
		return
	}

	err := s.core.SystemOperation().UpdatePlatformSettingConfig(ctx.Kit, input)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(nil)
}
