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

package mainline

import (
	"context"
	"fmt"
	"net/http"

	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

// SearchMainlineModelTopo get topo tree of model on mainline
func (m *topoManager) SearchMainlineModelTopo(ctx context.Context, header http.Header, withDetail bool) (*metadata.TopoModelNode, error) {
	rid := util.ExtractRequestIDFromContext(ctx)
	obj, err := NewModelMainline()
	if err != nil {
		blog.Errorf("new model mainline failed, err: %+v, rid: %s", err, rid)
		return nil, fmt.Errorf("new model mainline failed, err: %+v", err)
	}
	return obj.GetRoot(ctx, header, withDetail)
}
