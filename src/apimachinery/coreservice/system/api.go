/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 Tencent. All rights reserved.
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

package system

import (
	"context"
	"net/http"

	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

// GetUserConfig TODO
func (s *system) GetUserConfig(ctx context.Context, h http.Header) (*metadata.ResponseSysUserConfigData,
	errors.CCErrorCoder) {
	rid := util.ExtractRequestIDFromContext(ctx)

	resp := new(metadata.ReponseSysUserConfig)

	subPath := "/find/system/user_config"

	httpDoErr := s.client.Post().
		WithContext(ctx).
		Body(nil).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	if httpDoErr != nil {
		blog.Errorf("AddLabel failed, http request failed, err: %+v, rid: %s", httpDoErr, rid)
		return nil, errors.CCHttpError
	}
	if resp.Result == false || resp.Code != 0 {
		return nil, errors.New(resp.Code, resp.ErrMsg)
	}

	return &resp.Data, nil
}

// SearchConfigAdmin TODO
func (s *system) SearchConfigAdmin(ctx context.Context, h http.Header) (resp *metadata.ConfigAdminResult, err error) {
	resp = new(metadata.ConfigAdminResult)
	subPath := "/find/system/config_admin"

	err = s.client.Get().
		WithContext(ctx).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)

	return
}

// SearchPlatformSetting find platform config.
func (s *system) SearchPlatformSetting(ctx context.Context, h http.Header) (resp *metadata.PlatformSettingResult,
	err error) {
	resp = new(metadata.PlatformSettingResult)
	subPath := "/find/system_config/platform_setting"

	err = s.client.Get().
		WithContext(ctx).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)

	return
}

// UpdatePlatformSetting update platform config.
func (s *system) UpdatePlatformSetting(ctx context.Context, h http.Header, input *metadata.PlatformSettingConfig) (
	resp *metadata.BaseResp, err error) {

	resp = new(metadata.BaseResp)
	subPath := "/update/system_config/platform_setting"
	err = s.client.Put().
		WithContext(ctx).
		Body(input).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)

	return
}
