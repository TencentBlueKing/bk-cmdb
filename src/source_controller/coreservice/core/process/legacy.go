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

package process

import (
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/source_controller/coreservice/core"
)

func (p *processOperation) GetProc2Module(ctx core.ContextParams, option *metadata.GetProc2ModuleOption) ([]metadata.Proc2Module, errors.CCErrorCoder) {
	filter := mapstr.NewFromStruct(option, "json")
	filter = util.SetModOwner(filter, ctx.SupplierAccount)

	result := make([]metadata.Proc2Module, 0)
	if err := p.dbProxy.Table(common.BKTableNameProcModule).Find(filter).All(ctx, &result); err != nil {
		blog.Errorf("get process2module config failed. err: %v, rid:%s", err, ctx.ReqID)
		return nil, ctx.Error.CCError(common.CCErrProcSelectProc2Module)
	}

	return result, nil
}
