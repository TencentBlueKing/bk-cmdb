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
	"fmt"
	"net/http"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
)

// GetCustomObjects get objects which are custom(without mainline objects).
func (l *Logics) GetCustomObjects(ctx context.Context, header http.Header) ([]metadata.Object, error) {
	resp, err := l.CoreAPI.CoreService().Model().ReadModel(ctx, header, &metadata.QueryCondition{
		Fields: []string{common.BKObjIDField, common.BKObjNameField, common.BKFieldID},
		Page:   metadata.BasePage{Limit: common.BKNoLimit},
		Condition: map[string]interface{}{
			common.BKIsPre: false,
			common.BKClassificationIDField: map[string]interface{}{
				common.BKDBNE: "bk_biz_topo",
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("get custom models failed, err: %+v", err)
	}
	
	if len(resp.Info) == 0 {
		blog.Info("get custom models failed, no custom model is found")
		return nil, fmt.Errorf("no custom model is found")
	}

	return resp.Info, nil
}
