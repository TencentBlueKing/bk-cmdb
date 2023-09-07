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
	"sync"

	"configcenter/pkg/filter"
	filtertools "configcenter/pkg/tools/filter"
	"configcenter/src/ac/meta"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
)

// ListFieldTemplateAttr list field template attributes.
func (s *service) ListFieldTemplateAttr(cts *rest.Contexts) {
	opt := new(metadata.ListFieldTmplAttrOption)
	if err := cts.DecodeInto(opt); err != nil {
		cts.RespAutoError(err)
		return
	}

	if rawErr := opt.Validate(); rawErr.ErrCode != 0 {
		cts.RespAutoError(rawErr.ToCCError(cts.Kit.CCError))
		return
	}

	// check if user has the permission of the field template
	if authResp, authorized := s.auth.Authorize(cts.Kit, meta.ResourceAttribute{Basic: meta.Basic{
		Type: meta.FieldTemplate, Action: meta.Find, InstanceID: opt.TemplateID}}); !authorized {
		cts.RespNoAuth(authResp)
		return
	}

	attrFilter, err := filtertools.And(filtertools.GenAtomFilter(common.BKTemplateID, filter.Equal, opt.TemplateID),
		opt.Filter)
	if err != nil {
		blog.Errorf("list field template attributes failed, err: %v, opt: %+v, rid: %s", err, opt, cts.Kit.Rid)
		cts.RespAutoError(err)
		return
	}

	listOpt := &metadata.CommonQueryOption{
		CommonFilterOption: metadata.CommonFilterOption{Filter: attrFilter},
		Page:               opt.Page,
		Fields:             opt.Fields,
	}

	// list field template attributes
	res, err := s.clientSet.CoreService().FieldTemplate().ListFieldTemplateAttr(cts.Kit.Ctx, cts.Kit.Header, listOpt)
	if err != nil {
		blog.Errorf("list field template attributes failed, err: %v, opt: %+v, rid: %s", err, opt, cts.Kit.Rid)
		cts.RespAutoError(err)
		return
	}

	cts.RespEntity(res)
}

// CountFieldTemplateAttr count field templates' attributes
func (s *service) CountFieldTemplateAttr(ctx *rest.Contexts) {
	opt := new(metadata.CountFieldTmplResOption)
	if err := ctx.DecodeInto(opt); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if rawErr := opt.Validate(); rawErr.ErrCode != 0 {
		ctx.RespAutoError(rawErr.ToCCError(ctx.Kit.CCError))
		return
	}

	// count field template's attributes
	countInfos := make([]metadata.FieldTmplResCount, len(opt.TemplateIDs))

	var wg sync.WaitGroup
	var lock sync.Mutex
	var firstErr errors.CCErrorCoder
	pipeline := make(chan struct{}, 10)

	for idx := range opt.TemplateIDs {
		if firstErr != nil {
			break
		}

		pipeline <- struct{}{}
		wg.Add(1)

		go func(idx int) {
			defer func() {
				wg.Done()
				<-pipeline
			}()

			attrOpt := &metadata.CommonQueryOption{
				CommonFilterOption: metadata.CommonFilterOption{
					Filter: filtertools.GenAtomFilter(common.BKTemplateID, filter.Equal, opt.TemplateIDs[idx]),
				},
				Page: metadata.BasePage{EnableCount: true},
			}

			attrRes, err := s.clientSet.CoreService().FieldTemplate().ListFieldTemplateAttr(ctx.Kit.Ctx, ctx.Kit.Header,
				attrOpt)
			if err != nil {
				blog.Errorf("count field template attribute failed, err: %v, opt: %+v, rid: %s", err, opt, ctx.Kit.Rid)
				if firstErr == nil {
					firstErr = err
				}
				return
			}

			lock.Lock()
			countInfos[idx] = metadata.FieldTmplResCount{
				TemplateID: opt.TemplateIDs[idx],
				Count:      int(attrRes.Count),
			}
			lock.Unlock()
		}(idx)
	}

	wg.Wait()

	if firstErr != nil {
		ctx.RespAutoError(firstErr)
		return
	}

	ctx.RespEntity(countInfos)
}
