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

package auditlog

import (
	"context"
	"net/http"

	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

func (s *system) GetUserConfig(ctx context.Context, h http.Header) (*metadata.ResponseSysUserConfigData, errors.CCErrorCoder) {
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
