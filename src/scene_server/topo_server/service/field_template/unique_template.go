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
	"configcenter/pkg/filter"
	filtertools "configcenter/pkg/tools/filter"
	"configcenter/src/ac/meta"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
)

// ListFieldTemplateUnique list field template unique
func (s *service) ListFieldTemplateUnique(ctx *rest.Contexts) {
	opt := new(metadata.ListFieldTmplUniqueOption)
	if err := ctx.DecodeInto(opt); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if rawErr := opt.Validate(); rawErr.ErrCode != 0 {
		ctx.RespAutoError(rawErr.ToCCError(ctx.Kit.CCError))
		return
	}

	// check if user has the permission of the field template
	if authResp, authorized := s.auth.Authorize(ctx.Kit, meta.ResourceAttribute{Basic: meta.Basic{
		Type: meta.FieldTemplate, Action: meta.Find, InstanceID: opt.TemplateID}}); !authorized {
		ctx.RespNoAuth(authResp)
		return
	}

	uniqueFilter, err := filtertools.And(filtertools.GenAtomFilter(common.BKTemplateID, filter.Equal, opt.TemplateID),
		opt.Filter)
	if err != nil {
		blog.Errorf("list field template uniques failed, err: %v, opt: %+v, rid: %s", err, opt, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	listOpt := &metadata.CommonQueryOption{
		CommonFilterOption: metadata.CommonFilterOption{Filter: uniqueFilter},
		Page:               opt.Page,
		Fields:             opt.Fields,
	}

	// list field template uniques
	res, err := s.clientSet.CoreService().FieldTemplate().ListFieldTemplateUnique(ctx.Kit.Ctx, ctx.Kit.Header, listOpt)
	if err != nil {
		blog.Errorf("list field template uniques failed, err: %v, opt: %+v, rid: %s", err, opt, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(res)
}

// getObjectPausedAttr get the relationship between model ID and state.
func (s *service) getObjectPausedAttr(kit *rest.Kit, objIDs []int64) (map[int64]bool, error) {

	idMap := make(map[int64]struct{})
	for _, id := range objIDs {
		idMap[id] = struct{}{}
	}

	ids := make([]int64, 0)
	for id := range idMap {
		ids = append(ids, id)
	}

	query := &metadata.QueryCondition{
		Condition: map[string]interface{}{
			common.BKFieldID: mapstr.MapStr{
				common.BKDBIN: ids,
			},
		},
		Fields:         []string{common.BKFieldID, metadata.ModelFieldIsPaused},
		DisableCounter: true,
	}

	obj, err := s.clientSet.CoreService().Model().ReadModel(kit.Ctx, kit.Header, query)
	if err != nil {
		blog.Errorf("failed to find objects by query(%#v), err: %v, rid: %s", query, err, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, "object_ids")
	}

	if len(obj.Info) == 0 {
		blog.Errorf("failed to find objIDs by query(%#v), rid: %s", query, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrCommNotFound)
	}

	if len(obj.Info) != len(ids) {
		blog.Errorf("object ids are invalid, input: %+v, rid: %s", query, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, "object_ids")
	}

	idPausedMap := make(map[int64]bool)
	for _, info := range obj.Info {
		idPausedMap[info.ID] = info.IsPaused
	}
	return idPausedMap, nil
}
