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
package settemplate

import (
	"context"
	"fmt"
	"net/http"

	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/metadata"
)

func (st *SetTemplate) CreateSetTemplate(ctx context.Context, header http.Header, bizID int64, option metadata.CreateSetTemplateOption) (*metadata.SetTemplateResult, error) {
	ret := new(metadata.SetTemplateResult)
	subPath := fmt.Sprintf("/create/topo/set_template/bk_biz_id/%d/", bizID)

	err := st.client.Post().
		WithContext(ctx).
		Body(option).
		SubResource(subPath).
		WithHeaders(header).
		Do().
		Into(ret)

	if err != nil {
		blog.Errorf("CreateSetTemplate failed, http request failed, err: %+v", err)
		return nil, errors.CCHttpError
	}
	if ret.Result == false || ret.Code != 0 {
		return nil, errors.NewCCError(ret.Code, ret.ErrMsg)
	}

	return ret, nil
}

func (st *SetTemplate) UpdateSetTemplate(ctx context.Context, header http.Header, bizID int64, setTemplateID int64, option metadata.UpdateSetTemplateOption) (*metadata.SetTemplateResult, error) {
	ret := new(metadata.SetTemplateResult)
	subPath := fmt.Sprintf("/update/topo/set_template/%d/bk_biz_id/%d/", setTemplateID, bizID)

	err := st.client.Post().
		WithContext(ctx).
		Body(option).
		SubResource(subPath).
		WithHeaders(header).
		Do().
		Into(ret)

	if err != nil {
		blog.Errorf("UpdateSetTemplate failed, http request failed, err: %+v", err)
		return nil, errors.CCHttpError
	}
	if ret.Result == false || ret.Code != 0 {
		return nil, errors.NewCCError(ret.Code, ret.ErrMsg)
	}

	return ret, nil
}

func (st *SetTemplate) DeleteSetTemplate(ctx context.Context, header http.Header, bizID int64, option metadata.CreateSetTemplateOption) error {
	ret := struct {
		metadata.BaseResp
	}{}
	subPath := fmt.Sprintf("/deletemany/topo/set_template/bk_biz_id/%d/", bizID)

	err := st.client.Delete().
		WithContext(ctx).
		Body(option).
		SubResource(subPath).
		WithHeaders(header).
		Do().
		Into(&ret)

	if err != nil {
		blog.Errorf("DeleteSetTemplate failed, http request failed, err: %+v", err)
		return errors.CCHttpError
	}
	if ret.Result == false || ret.Code != 0 {
		return errors.NewCCError(ret.Code, ret.ErrMsg)
	}

	return err
}
