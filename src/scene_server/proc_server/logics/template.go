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

package logics

import (
	"context"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
)

func (lgc *Logics) GetTemplateAttributes(ctx context.Context, ownerID string) ([]metadata.Header, error) {
	params := map[string]interface{}{
		common.BKOwnerIDField: ownerID,
		common.BKObjIDField:   common.BKInnerObjIDConfigTemp,
	}

	result, err := lgc.CoreAPI.ObjectController().Meta().SelectObjectAttWithParams(ctx, lgc.header, params)
	if err != nil {
		blog.Errorf("GetTemplateAttributes SelectObjectAttWithParams http do error,err:%s,query:%+v,rid:%s", err.Error(), params, lgc.rid)
		return nil, lgc.ccErr.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !result.Result {
		blog.Errorf("GetTemplateAttributes SelectObjectAttWithParams http response error,err code:%d,err msg:%s,query:%+v,rid:%s", result.Code, result.ErrMsg, params, lgc.rid)
		return nil, lgc.ccErr.New(result.Code, result.ErrMsg)
	}

	headers := make([]metadata.Header, 0)
	for _, p := range result.Data {
		if p.PropertyID == common.BKChildStr {
			continue
		}
		headers = append(headers, metadata.Header{
			PropertyID:   p.PropertyID,
			PropertyName: p.PropertyName,
		})
	}

	return headers, nil
}

func (lgc *Logics) GetTemplateInstanceDetails(ctx context.Context, tempID int64) (map[string]interface{}, error) {

	params := metadata.QueryInput{
		Condition: map[string]interface{}{
			common.BKOwnerIDField:   lgc.ownerID,
			common.BKTemlateIDField: tempID,
		},
	}
	result, err := lgc.CoreAPI.ObjectController().Instance().SearchObjects(ctx, common.BKInnerObjIDConfigTemp, lgc.header, &params)
	if err != nil {
		blog.Errorf("GetTemplateInstanceDetails SelectObjectAttWithParams http do error,err:%s,query:%+v,rid:%s", err.Error(), params, lgc.rid)
		return nil, lgc.ccErr.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !result.Result {
		blog.Errorf("GetTemplateInstanceDetails SelectObjectAttWithParams http response error,err code:%d,err msg:%s,query:%+v,rid:%s", result.Code, result.ErrMsg, params, lgc.rid)
		return nil, lgc.ccErr.New(result.Code, result.ErrMsg)
	}

	return result.Data.Info[0], nil

}
