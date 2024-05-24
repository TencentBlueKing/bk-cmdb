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

// Package idrule package
package idrule

import (
	"context"
	"net/http"

	"configcenter/src/common/errors"
	"configcenter/src/common/metadata"
)

// UpdateInstIDRule update instance id rule field
func (a *idRule) UpdateInstIDRule(ctx context.Context, h http.Header,
	opt *metadata.UpdateInstIDRuleOption) errors.CCErrorCoder {

	resp := new(metadata.BaseResp)

	err := a.client.Put().
		WithContext(ctx).
		Body(opt).
		SubResourcef("/update/inst/id_rule").
		WithHeaders(h).
		Do().
		Into(resp)

	if err != nil {
		return errors.CCHttpError
	}

	if err := resp.CCError(); err != nil {
		return err
	}

	return nil
}
