/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
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

// UpdateGlobalConfig update global config
func (s *system) UpdateGlobalConfig(ctx context.Context, h http.Header, typeId string,
	input interface{}) error {

	resp := new(metadata.BaseResp)
	subPath := "/update/global_config/%s"
	err := s.client.Put().
		WithContext(ctx).
		Body(input).
		SubResourcef(subPath, typeId).
		WithHeaders(h).
		Do().
		Into(resp)

	if err != nil {
		return errors.CCHttpError
	}

	if err = resp.CCError(); err != nil {
		return err
	}
	return err
}

// SearchGlobalConfig search global config.
func (s *system) SearchGlobalConfig(ctx context.Context, h http.Header,
	options *metadata.GlobalConfOptions) (*metadata.GlobalSettingConfig, error) {

	resp := new(metadata.GlobalSettingResult)
	subPath := "/find/global_config"

	err := s.client.Get().
		WithContext(ctx).
		SubResourcef(subPath).
		WithHeaders(h).
		Body(options).
		Do().
		Into(resp)

	if err != nil {
		return nil, errors.CCHttpError
	}

	if err = resp.CCError(); err != nil {
		return nil, err
	}

	return &resp.Data, nil
}

// GetHostSnapDataID get tenant related host snap data id
func (s *system) GetHostSnapDataID(ctx context.Context, h http.Header) (int64, error) {
	resp := new(metadata.DataIDResp)
	subPath := "/find/host_snap/data_id"

	err := s.client.Post().
		WithContext(ctx).
		SubResourcef(subPath).
		WithHeaders(h).
		Body(nil).
		Do().
		Into(resp)

	if err != nil {
		return 0, errors.CCHttpError
	}

	if err = resp.CCError(); err != nil {
		return 0, err
	}

	return resp.Data, nil
}

// GetTenantBySnapDataID get tenant by snap data id
func (s *system) GetTenantBySnapDataID(ctx context.Context, h http.Header, dataID int64) (string, error) {
	resp := new(metadata.TenantInfoResp)
	subPath := "/find/host_snap/tenant/data_id/%d"

	err := s.client.Post().
		WithContext(ctx).
		SubResourcef(subPath, dataID).
		WithHeaders(h).
		Body(nil).
		Do().
		Into(resp)

	if err != nil {
		return "", errors.CCHttpError
	}

	if err = resp.CCError(); err != nil {
		return "", err
	}

	return resp.Data, nil
}
