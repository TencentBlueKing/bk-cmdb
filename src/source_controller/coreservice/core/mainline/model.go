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
