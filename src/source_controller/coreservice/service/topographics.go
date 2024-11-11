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

package service

import (
	"context"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	meta "configcenter/src/common/metadata"
	"configcenter/src/storage/driver/mongodb"
)

// SearchTopoGraphics TODO
// CreateClassification create object's classification
func (s *coreService) SearchTopoGraphics(ctx *rest.Contexts) {
	selector := meta.TopoGraphics{}
	if jsErr := ctx.DecodeInto(&selector); jsErr != nil {
		ctx.RespAutoError(jsErr)
		return
	}

	cond := mapstr.MapStr{
		"scope_type": selector.ScopeType,
		"scope_id":   selector.ScopeID,
	}

	results := make([]meta.TopoGraphics, 0)
	if selErr := mongodb.Client().Table(common.BKTableNameTopoGraphics).Find(cond).All(ctx.Kit.Ctx,
		&results); selErr != nil {
		blog.Errorf("select data failed, err: %v, cond: %v, rid: %s", selErr, cond, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommDBSelectFailed))
		return
	}
	ctx.RespEntity(results)
}

// UpdateTopoGraphics TODO
func (s *coreService) UpdateTopoGraphics(ctx *rest.Contexts) {
	inputBody := struct {
		Data []meta.TopoGraphics `json:"data" field:"data" bson:"data"`
	}{}
	if jsErr := ctx.DecodeInto(&inputBody); jsErr != nil {
		ctx.RespAutoError(jsErr)
		return
	}

	for index := range inputBody.Data {
		inputBody.Data[index].SetSupplierAccount(ctx.Kit.SupplierAccount)
		cond := mapstr.MapStr{
			"scope_type": inputBody.Data[index].ScopeType,
			"scope_id":   inputBody.Data[index].ScopeID,
			"node_type":  inputBody.Data[index].NodeType,
			"bk_obj_id":  inputBody.Data[index].ObjID,
			"bk_inst_id": inputBody.Data[index].InstID,
		}

		cnt, err := mongodb.Client().Table(common.BKTableNameTopoGraphics).Find(cond).Count(ctx.Kit.Ctx)
		if err != nil {
			blog.Errorf("search data failed, err: %v, cond: %v, rid: %s", err, cond, ctx.Kit.Rid)
			ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommDBSelectFailed))
			return
		}
		if cnt == 0 {
			err = mongodb.Client().Table(common.BKTableNameTopoGraphics).Insert(ctx.Kit.Ctx, inputBody.Data[index])
			if err != nil {
				blog.Errorf("insert data failed, err: %v, inputBody: %v, rid: %s", err, inputBody, ctx.Kit.Rid)
				ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommDBInsertFailed))
				return
			}
		} else {
			if err = mongodb.Client().Table(common.BKTableNameTopoGraphics).Update(context.Background(), cond,
				inputBody.Data[index]); err != nil {
				blog.Errorf("insert data failed, err: %v, inputBody: %v, rid: %s", err, inputBody, ctx.Kit.Rid)
				ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommDBUpdateFailed))
				return
			}
		}
	}

	ctx.RespEntity(nil)
}
