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

package host

import (
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/source_controller/coreservice/core"
	"configcenter/src/source_controller/coreservice/core/host/identifier"
)

func (hm *hostManager) Identifier(ctx core.ContextParams, input *metadata.SearchHostIdentifierParam) ([]metadata.HostIdentifier, error) {
	identifier := identifier.NewIdentifier(hm.DbProxy)
	host, err := identifier.Identifier(ctx, input.HostIDs)
	if err != nil {
		blog.ErrorJSON("Identifier get host identifier error. input:%s, err:%s, rid:%s", input, err.Error(), ctx.ReqID)
		return nil, err
	}
	return host, nil
}
