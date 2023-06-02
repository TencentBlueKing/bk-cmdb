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

package fieldtmpl

import (
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/storage/driver/mongodb"
)

// ListObjFieldTmplRel list field template and object relations.
func (s *service) ListObjFieldTmplRel(cts *rest.Contexts) {
	opt := new(metadata.CommonQueryOption)
	if err := cts.DecodeInto(opt); err != nil {
		cts.RespAutoError(err)
		return
	}

	if rawErr := opt.Validate(); rawErr.ErrCode != 0 {
		cts.RespAutoError(rawErr.ToCCError(cts.Kit.CCError))
		return
	}

	filter, err := opt.ToMgo()
	if err != nil {
		cts.RespAutoError(cts.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, err.Error()))
		return
	}

	table := common.BKTableNameObjFieldTemplateRelation
	filter = util.SetQueryOwner(filter, cts.Kit.SupplierAccount)

	if opt.Page.EnableCount {
		count, err := mongodb.Client().Table(table).Find(filter).Count(cts.Kit.Ctx)
		if err != nil {
			blog.Errorf("count obj field tmpl relation failed, err: %v, filter: %+v, rid: %v", err, filter, cts.Kit.Rid)
			cts.RespAutoError(cts.Kit.CCError.CCError(common.CCErrCommDBSelectFailed))
			return
		}

		cts.RespEntity(metadata.ObjFieldTmplRelInfo{Count: count})
		return
	}

	relations := make([]metadata.ObjFieldTemplateRelation, 0)
	err = mongodb.Client().Table(table).Find(filter).Start(uint64(opt.Page.Start)).Limit(uint64(opt.Page.Limit)).
		Sort(opt.Page.Sort).Fields(opt.Fields...).All(cts.Kit.Ctx, &relations)
	if err != nil {
		blog.Errorf("list obj field tmpl relations failed, err: %v, filter: %+v, rid: %s", err, filter, cts.Kit.Rid)
		cts.RespAutoError(cts.Kit.CCError.CCError(common.CCErrCommDBSelectFailed))
		return
	}

	cts.RespEntity(metadata.ObjFieldTmplRelInfo{Info: relations})
}
