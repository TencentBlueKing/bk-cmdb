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

package parser

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	"configcenter/src/ac/meta"
	"configcenter/src/common"
)

func (ps *parseStream) topology() *parseStream {
	if ps.shouldReturn() {
		return ps
	}

	ps.business().
		mainline().
		object().
		objectAttributeGroup().
		objectModule().
		objectSet().
		audit().
		fullTextSearch().
		cloudArea()

	return ps
}

var (
	createBusinessRegexp = regexp.MustCompile(`^/api/v3/biz/[^\s/]+/?$`)
	updateBusinessRegexp = regexp.MustCompile(`^/api/v3/biz/[^\s/]+/[0-9]+/?$`)
	// deleteBusinessRegexp             = regexp.MustCompile(`^/api/v3/biz/[^\s/]+/[0-9]+/?$`)
	findBusinessRegexp               = regexp.MustCompile(`^/api/v3/biz/search/[^\s/]+/?$`)
	findResourcePoolBusinessRegexp   = regexp.MustCompile(`^/api/v3/biz/default/[^\s/]+/search/?$`)
	createResourcePoolBusinessRegexp = regexp.MustCompile(`^/api/v3/biz/default/[^\s/]+/?$`)
	updateBusinessStatusRegexp       = regexp.MustCompile(`^/api/v3/biz/status/[^\s/]+/[^\s/]+/[0-9]+/?$`)
)

const findReducedBusinessList = `/api/v3/biz/with_reduced`
const findSimplifiedBusinessList = `/api/v3/biz/simplify`

func (ps *parseStream) business() *parseStream {
	if ps.shouldReturn() {
		return ps
	}

	// find reduced business list for the user with any business resources
	if ps.hitPattern(findReducedBusinessList, http.MethodGet) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:   meta.Business,
					Action: meta.SkipAction,
				},
			},
		}
		return ps
	}

	// find simplified business list with limited fields return
	if ps.hitPattern(findSimplifiedBusinessList, http.MethodGet) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:   meta.Business,
					Action: meta.SkipAction,
				},
			},
		}
		return ps
	}

	// create business, this is not a normalize api.
	// TODO: update this api format.
	if ps.hitRegexp(createBusinessRegexp, http.MethodPost) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:   meta.Business,
					Action: meta.Create,
				},
			},
		}
		return ps
	}

	// 创建主机池业务
	if ps.hitRegexp(createResourcePoolBusinessRegexp, http.MethodPost) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:   meta.Business,
					Action: meta.Create,
				},
			},
		}
		return ps
	}

	// update business, this is not a normalize api.
	// TODO: update this api format.
	if ps.hitRegexp(updateBusinessRegexp, http.MethodPut) {
		if len(ps.RequestCtx.Elements) != 5 {
			ps.err = errors.New("invalid update business request uri")
			return ps
		}

		bizID, err := strconv.ParseInt(ps.RequestCtx.Elements[4], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("udpate business, but got invalid business id %s", ps.RequestCtx.Elements[4])
			return ps
		}

		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:       meta.Business,
					Action:     meta.Update,
					InstanceID: bizID,
				},
			},
		}
		return ps
	}

	// update business enable status, this is not a normalize api.
	// TODO: update this api format.
	if ps.hitRegexp(updateBusinessRegexp, http.MethodPut) {
		if len(ps.RequestCtx.Elements) != 7 {
			ps.err = errors.New("invalid update business enable status request uri")
			return ps
		}

		bizID, err := strconv.ParseInt(ps.RequestCtx.Elements[6], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("udpate business enable status, but got invalid business id %s", ps.RequestCtx.Elements[4])
			return ps
		}

		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:       meta.Business,
					Action:     meta.Update,
					InstanceID: bizID,
				},
			},
		}
		return ps
	}

	// delete business, this is not a normalize api.
	// TODO: update this api format
	if ps.hitRegexp(updateBusinessRegexp, http.MethodDelete) {
		if len(ps.RequestCtx.Elements) != 5 {
			ps.err = errors.New("invalid delete business request uri")
			return ps
		}

		bizID, err := strconv.ParseInt(ps.RequestCtx.Elements[4], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("delete business, but got invalid business id %s", ps.RequestCtx.Elements[4])
			return ps
		}

		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:       meta.Business,
					Action:     meta.Delete,
					InstanceID: bizID,
				},
			},
		}
		return ps
	}

	// find business, this is not a normalize api.
	// TODO: update this api format
	if ps.hitRegexp(findBusinessRegexp, http.MethodPost) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:   meta.Business,
					Action: meta.SkipAction,
				},
				// we don't know if one or more business is to find, so we assume it's a find many
				// business operation.
			},
		}
		return ps
	}

	// find resource pool business
	if ps.hitRegexp(findResourcePoolBusinessRegexp, http.MethodPost) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:   meta.Business,
					Action: meta.SkipAction,
				},
				// we don't know if one or more business is to find, so we assume it's a find many
				// business operation.
			},
		}
		return ps
	}

	// update business status to `disabled` or `enable`
	if ps.hitRegexp(updateBusinessStatusRegexp, http.MethodPut) {
		bizID, err := strconv.ParseInt(ps.RequestCtx.Elements[6], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("delete business, but got invalid business id %s", ps.RequestCtx.Elements[4])
			return ps
		}
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:       meta.Business,
					Action:     meta.Archive,
					InstanceID: bizID,
				},
				// we don't know if one or more business is to find, so we assume it's a find many
				// business operation.
			},
		}
		return ps
	}

	return ps
}

var (
	findMainlineIdleFaultModuleRegexp               = regexp.MustCompile(`^/api/v3/topo/internal/[^\s/]+/[0-9]+/?$`)
	findMainlineIdleFaultModuleWithStatisticsRegexp = regexp.MustCompile(`^/api/v3/topo/internal/[^\s/]+/[0-9]+/with_statistics/?$`)
	findBriefBizTopoRegexp                          = regexp.MustCompile(`^/api/v3/find/topo/tree/brief/biz/[0-9]+/?$`)
)

func (ps *parseStream) mainline() *parseStream {
	if ps.shouldReturn() {
		return ps
	}

	// find internal mainline idle and fault module operation.
	if ps.hitRegexp(findMainlineIdleFaultModuleRegexp, http.MethodGet) {
		if len(ps.RequestCtx.Elements) != 6 {
			ps.err = errors.New("find mainline idle and fault module, but got invalid url")
			return ps
		}

		bizID, err := strconv.ParseInt(ps.RequestCtx.Elements[5], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("find mainline idle and fault module, but got invalid business id %s", ps.RequestCtx.Elements[5])
			return ps
		}
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:   meta.MainlineInstance,
					Action: meta.Find,
				},
			},
		}

		return ps
	}

	// find internal mainline idle and fault module with statistics operation.
	if ps.hitRegexp(findMainlineIdleFaultModuleWithStatisticsRegexp, http.MethodGet) {
		if len(ps.RequestCtx.Elements) != 7 {
			ps.err = errors.New("find mainline idle and fault module with statistics, but got invalid url")
			return ps
		}

		bizID, err := strconv.ParseInt(ps.RequestCtx.Elements[5], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("find mainline idle and fault module with statistics, but got invalid business id %s", ps.RequestCtx.Elements[5])
			return ps
		}
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:   meta.MainlineInstance,
					Action: meta.Find,
				},
			},
		}

		return ps
	}

	// find brief biz topo
	if ps.hitRegexp(findBriefBizTopoRegexp, http.MethodPost) {
		if len(ps.RequestCtx.Elements) != 8 {
			ps.err = errors.New("find brief biz topo, but got invalid url")
			return ps
		}

		bizID, err := strconv.ParseInt(ps.RequestCtx.Elements[7], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("find brief biz topo, but got invalid business id %s", ps.RequestCtx.Elements[5])
			return ps
		}
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:   meta.MainlineInstance,
					Action: meta.Find,
				},
			},
		}

		return ps
	}

	return ps
}

const (
	objectStatistics         = "/api/v3/object/statistics"
)

func (ps *parseStream) object() *parseStream {
	if ps.shouldReturn() {
		return ps
	}

	// 统计模型使用情况
	if ps.hitPattern(objectStatistics, http.MethodGet) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Action: meta.SkipAction,
				},
			},
		}
		return ps
	}

	return ps
}

var (
	removeAttributeAwayFromGroupRegexp = regexp.MustCompile(`^/api/v3/objectatt/group/owner/[^\s/]+/object/[^\s/]+/propertyids/[^\s/]+/groupids/[^\s/]+/?$`)
)

func (ps *parseStream) objectAttributeGroup() *parseStream {
	if ps.shouldReturn() {
		return ps
	}

	// remove a object's attribute away from a group.
	if ps.hitRegexp(removeAttributeAwayFromGroupRegexp, http.MethodDelete) {
		if len(ps.RequestCtx.Elements) != 12 {
			ps.err = errors.New("remove a object attribute away from a group, but got invalid uri")
			return ps
		}

		bizID, err := ps.RequestCtx.getBizIDFromBody()
		if err != nil {
			ps.err = err
			return ps
		}

		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:   meta.ModelAttributeGroup,
					Action: meta.Delete,
					Name:   ps.RequestCtx.Elements[11],
				},
			},
		}
		return ps
	}

	return ps
}

var (
	createModuleRegexp                = regexp.MustCompile(`^/api/v3/module/[0-9]+/[0-9]+/?$`)
	deleteModuleRegexp                = regexp.MustCompile(`^/api/v3/module/[0-9]+/[0-9]+/[0-9]+/?$`)
	updateModuleRegexp                = regexp.MustCompile(`^/api/v3/module/[0-9]+/[0-9]+/[0-9]+/?$`)
	updateModuleHostApplyStatusRegexp = regexp.MustCompile(`^/api/v3/module/host_apply_enable_status/bk_biz_id/([0-9]+)/bk_module_id/([0-9]+)/?$`)
	findModuleRegexp                  = regexp.MustCompile(`^/api/v3/module/search/[^\s/]+/[0-9]+/[0-9]+/?$`)
	findMouduleByConditionRegexp      = regexp.MustCompile(`^/api/v3/findmany/module/biz/[0-9]+/?$`)
	findMouduleBatchRegexp            = regexp.MustCompile(`^/api/v3/findmany/module/bk_biz_id/[0-9]+/?$`)
	findMouduleWithRelationRegexp     = regexp.MustCompile(`^/api/v3/findmany/module/with_relation/biz/[0-9]+/?$`)
	findModuleByServiceTemplateRegexp = regexp.MustCompile(`^/api/v3/module/bk_biz_id/(?P<bk_biz_id>[0-9]+)/service_template_id/(?P<service_template_id>[0-9]+)/?$`)
)

func (ps *parseStream) objectModule() *parseStream {
	if ps.shouldReturn() {
		return ps
	}

	// create module
	if ps.hitRegexp(createModuleRegexp, http.MethodPost) {
		if len(ps.RequestCtx.Elements) != 5 {
			ps.err = errors.New("create module, but got invalid url")
			return ps
		}

		bizID, err := strconv.ParseInt(ps.RequestCtx.Elements[3], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("create module, but got invalid business id %s", ps.RequestCtx.Elements[3])
			return ps
		}

		setID, err := strconv.ParseInt(ps.RequestCtx.Elements[4], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("create module, but got invalid set id %s", ps.RequestCtx.Elements[4])
			return ps
		}

		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:   meta.ModelModule,
					Action: meta.Create,
				},
				Layers: []meta.Item{
					{
						Type:       meta.ModelInstance,
						Name:       "set",
						InstanceID: setID,
					},
				},
			},
		}
		return ps
	}

	// delete module operation.
	if ps.hitRegexp(deleteModuleRegexp, http.MethodDelete) {
		if len(ps.RequestCtx.Elements) != 6 {
			ps.err = errors.New("delete module, but got invalid url")
			return ps
		}

		bizID, err := strconv.ParseInt(ps.RequestCtx.Elements[3], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("delete module, but got invalid business id %s", ps.RequestCtx.Elements[3])
			return ps
		}

		setID, err := strconv.ParseInt(ps.RequestCtx.Elements[4], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("delete module, but got invalid set id %s", ps.RequestCtx.Elements[4])
			return ps
		}

		moduleID, err := strconv.ParseInt(ps.RequestCtx.Elements[5], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("delete module, but got invalid module id %s", ps.RequestCtx.Elements[5])
			return ps
		}

		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:       meta.ModelModule,
					Action:     meta.Delete,
					InstanceID: moduleID,
				},
				Layers: []meta.Item{
					{
						Type:       meta.ModelInstance,
						Name:       "set",
						InstanceID: setID,
					},
				},
			},
		}
		return ps
	}

	// update module operation.
	if ps.hitRegexp(updateModuleRegexp, http.MethodPut) {
		if len(ps.RequestCtx.Elements) != 6 {
			ps.err = errors.New("update module, but got invalid url")
			return ps
		}

		bizID, err := strconv.ParseInt(ps.RequestCtx.Elements[3], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("update module, but got invalid business id %s", ps.RequestCtx.Elements[3])
			return ps
		}

		setID, err := strconv.ParseInt(ps.RequestCtx.Elements[4], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("update module, but got invalid set id %s", ps.RequestCtx.Elements[4])
			return ps
		}

		moduleID, err := strconv.ParseInt(ps.RequestCtx.Elements[5], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("update module, but got invalid module id %s", ps.RequestCtx.Elements[5])
			return ps
		}
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:       meta.ModelModule,
					Action:     meta.Update,
					InstanceID: moduleID,
				},
				Layers: []meta.Item{
					{
						Type:       meta.ModelInstance,
						Name:       "set",
						InstanceID: setID,
					},
				},
			},
		}
		return ps
	}
	if ps.hitRegexp(updateModuleHostApplyStatusRegexp, http.MethodPut) {
		if len(ps.RequestCtx.Elements) != 8 {
			ps.err = errors.New("update module host apply enabled status, but got invalid url")
			return ps
		}

		bizID, err := strconv.ParseInt(ps.RequestCtx.Elements[5], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("update module host apply enabled status, but got invalid business id %s", ps.RequestCtx.Elements[5])
			return ps
		}

		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:   meta.HostApply,
					Action: meta.Update,
				},
			},
		}
		return ps
	}

	// find module operation.
	if ps.hitRegexp(findModuleRegexp, http.MethodPost) {
		if len(ps.RequestCtx.Elements) != 7 {
			ps.err = errors.New("find module, but got invalid url")
			return ps
		}

		bizID, err := strconv.ParseInt(ps.RequestCtx.Elements[5], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("find module, but got invalid business id %s", ps.RequestCtx.Elements[5])
			return ps
		}

		setID, err := strconv.ParseInt(ps.RequestCtx.Elements[6], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("find module, but got invalid set id %s", ps.RequestCtx.Elements[6])
			return ps
		}
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:   meta.ModelModule,
					Action: meta.FindMany,
				},
				Layers: []meta.Item{
					{
						Type:       meta.ModelSet,
						InstanceID: setID,
					},
				},
			},
		}
		return ps
	}

	// find module by service template.
	if ps.hitRegexp(findModuleByServiceTemplateRegexp, http.MethodPost) {
		if len(ps.RequestCtx.Elements) != 7 {
			ps.err = errors.New("find module by service template id, but got invalid url")
			return ps
		}
		var bizIDStr string
		match := findModuleByServiceTemplateRegexp.FindStringSubmatch(ps.RequestCtx.URI)
		for i, name := range findModuleByServiceTemplateRegexp.SubexpNames() {
			if name == common.BKAppIDField {
				bizIDStr = match[i]
				break
			}
		}
		bizID, err := strconv.ParseInt(bizIDStr, 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("find module, but parse bk_biz_id failed, bizIDStr: %s, uri: %s, err: %+v", bizIDStr, ps.RequestCtx.URI, err)
			return ps
		}

		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:   meta.ModelModule,
					Action: meta.FindMany,
				},
				Layers: []meta.Item{},
			},
		}
		return ps
	}

	// find modules by condition in one biz
	if ps.hitRegexp(findMouduleByConditionRegexp, http.MethodPost) {
		if len(ps.RequestCtx.Elements) != 6 {
			ps.err = errors.New("find module batch, but got invalid url")
			return ps
		}

		bizID, err := strconv.ParseInt(ps.RequestCtx.Elements[5], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("update module batch, but got invalid biz id %s", ps.RequestCtx.Elements[5])
			return ps
		}

		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:   meta.ModelModule,
					Action: meta.FindMany,
				},
				Layers: []meta.Item{},
			},
		}
		return ps
	}

	// find module batch in one biz
	if ps.hitRegexp(findMouduleBatchRegexp, http.MethodPost) {
		if len(ps.RequestCtx.Elements) != 6 {
			ps.err = errors.New("find module batch, but got invalid url")
			return ps
		}

		bizID, err := strconv.ParseInt(ps.RequestCtx.Elements[5], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("find module batch, but got invalid biz id %s", ps.RequestCtx.Elements[5])
			return ps
		}

		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:   meta.ModelModule,
					Action: meta.FindMany,
				},
				Layers: []meta.Item{},
			},
		}
		return ps
	}

	// find module with relation in one biz
	if ps.hitRegexp(findMouduleWithRelationRegexp, http.MethodPost) {
		if len(ps.RequestCtx.Elements) != 7 {
			ps.err = errors.New("find module with relation, but got invalid url")
			return ps
		}
		bizID, err := strconv.ParseInt(ps.RequestCtx.Elements[6], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("find module with relation, but got invalid biz id %s", ps.RequestCtx.Elements[6])
			return ps
		}

		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:   meta.ModelModule,
					Action: meta.FindMany,
				},
				Layers: []meta.Item{},
			},
		}
		return ps
	}

	return ps
}

var (
	createSetRegexp      = regexp.MustCompile(`^/api/v3/set/[0-9]+/?$`)
	batchCreateSetRegexp = regexp.MustCompile(`^/api/v3/set/[0-9]+/batch/?$`)
	deleteSetRegexp      = regexp.MustCompile(`^/api/v3/set/[0-9]+/[0-9]+/?$`)
	deleteManySetRegexp  = regexp.MustCompile(`^/api/v3/set/[0-9]+/batch$`)
	updateSetRegexp      = regexp.MustCompile(`^/api/v3/set/[0-9]+/[0-9]+/?$`)
	findSetRegexp        = regexp.MustCompile(`^/api/v3/set/search/[^\s/]+/[0-9]+/?$`)
	findSetBatchRegexp   = regexp.MustCompile(`^/api/v3/findmany/set/bk_biz_id/[0-9]+/?$`)
)

func (ps *parseStream) objectSet() *parseStream {
	if ps.shouldReturn() {
		return ps
	}

	// create set
	if ps.hitRegexp(createSetRegexp, http.MethodPost) {
		if len(ps.RequestCtx.Elements) != 4 {
			ps.err = errors.New("create set, but got invalid url")
			return ps
		}

		bizID, err := strconv.ParseInt(ps.RequestCtx.Elements[3], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("create set, but got invalid business id %s", ps.RequestCtx.Elements[3])
			return ps
		}
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:   meta.ModelSet,
					Action: meta.Create,
				},
			},
		}
		return ps
	}

	// batch create set
	if ps.hitRegexp(batchCreateSetRegexp, http.MethodPost) {
		if len(ps.RequestCtx.Elements) != 5 {
			ps.err = errors.New("batch create set, but got invalid url")
			return ps
		}

		bizID, err := strconv.ParseInt(ps.RequestCtx.Elements[3], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("batch create set, but got invalid business id %s", ps.RequestCtx.Elements[3])
			return ps
		}
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:   meta.ModelSet,
					Action: meta.Create,
				},
			},
		}
		return ps
	}

	// delete set operation.
	if ps.hitRegexp(deleteSetRegexp, http.MethodDelete) {
		if len(ps.RequestCtx.Elements) != 5 {
			ps.err = errors.New("delete set, but got invalid url")
			return ps
		}

		bizID, err := strconv.ParseInt(ps.RequestCtx.Elements[3], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("delete set, but got invalid business id %s", ps.RequestCtx.Elements[3])
			return ps
		}

		setID, err := strconv.ParseInt(ps.RequestCtx.Elements[4], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("delete set, but got invalid set id %s", ps.RequestCtx.Elements[4])
			return ps
		}
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:       meta.ModelSet,
					Action:     meta.Delete,
					InstanceID: setID,
				},
			},
		}
		return ps
	}

	// delete many set operation.
	if ps.hitRegexp(deleteManySetRegexp, http.MethodDelete) {
		if len(ps.RequestCtx.Elements) != 5 {
			ps.err = errors.New("delete set list, but got invalid url")
			return ps
		}

		bizID, err := strconv.ParseInt(ps.RequestCtx.Elements[3], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("delete set list, but got invalid business id %s", ps.RequestCtx.Elements[3])
			return ps
		}
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:   meta.ModelSet,
					Action: meta.DeleteMany,
				},
			},
		}
		return ps
	}

	// update set operation.
	if ps.hitRegexp(updateSetRegexp, http.MethodPut) {
		if len(ps.RequestCtx.Elements) != 5 {
			ps.err = errors.New("update set, but got invalid url")
			return ps
		}

		bizID, err := strconv.ParseInt(ps.RequestCtx.Elements[3], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("update set, but got invalid business id %s", ps.RequestCtx.Elements[3])
			return ps
		}

		setID, err := strconv.ParseInt(ps.RequestCtx.Elements[4], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("update set, but got invalid set id %s", ps.RequestCtx.Elements[4])
			return ps
		}
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:       meta.ModelSet,
					Action:     meta.Update,
					InstanceID: setID,
				},
			},
		}
		return ps
	}

	// find set operation.
	if ps.hitRegexp(findSetRegexp, http.MethodPost) {
		if len(ps.RequestCtx.Elements) != 6 {
			ps.err = errors.New("find set, but got invalid url")
			return ps
		}

		bizID, err := strconv.ParseInt(ps.RequestCtx.Elements[5], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("find set, but got invalid business id %s", ps.RequestCtx.Elements[5])
			return ps
		}
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:   meta.ModelSet,
					Action: meta.FindMany,
				},
			},
		}
		return ps
	}

	// find set operation.
	if ps.hitRegexp(findSetBatchRegexp, http.MethodPost) {
		if len(ps.RequestCtx.Elements) != 6 {
			ps.err = errors.New("find set batch, but got invalid url")
			return ps
		}

		bizID, err := strconv.ParseInt(ps.RequestCtx.Elements[5], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("find set batch, but got invalid business id %s", ps.RequestCtx.Elements[5])
			return ps
		}

		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:   meta.ModelSet,
					Action: meta.FindMany,
				},
			},
		}
		return ps
	}

	return ps
}

var (
	searchAuditDict   = `/api/v3/find/audit_dict`
	searchAuditList   = `/api/v3/findmany/audit_list`
	searchAuditDetail = `/api/v3/find/audit`
)

func (ps *parseStream) audit() *parseStream {
	if ps.shouldReturn() {
		return ps
	}

	if ps.hitPattern(searchAuditDict, http.MethodGet) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:   meta.AuditLog,
					Action: meta.Find,
				},
			},
		}
		return ps
	}

	if ps.hitPattern(searchAuditList, http.MethodPost) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:   meta.AuditLog,
					Action: meta.FindMany,
				},
			},
		}
		return ps
	}

	if ps.hitPattern(searchAuditDetail, http.MethodPost) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:   meta.AuditLog,
					Action: meta.Find,
				},
			},
		}
		return ps
	}

	return ps
}

var (
	fullTextSearchPattern = "/api/v3/find/full_text"
)

func (ps *parseStream) fullTextSearch() *parseStream {
	if ps.shouldReturn() {
		return ps
	}

	if ps.hitPattern(fullTextSearchPattern, http.MethodPost) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Action: meta.SkipAction,
				},
			},
		}
		return ps
	}

	return ps
}

const (
	findManyCloudAreaPattern      = "/api/v3/findmany/cloudarea"
	createCloudAreaPattern        = "/api/v3/create/cloudarea"
	createManyCloudAreaPattern    = "/api/v3/createmany/cloudarea"
	findCloudAreaHostCountPattern = "/api/v3/findmany/cloudarea/hostcount"
)

var (
	updateCloudAreaRegexp = regexp.MustCompile(`^/api/v3/update/cloudarea/[0-9]+/?$`)
	deleteCloudAreaRegexp = regexp.MustCompile(`^/api/v3/delete/cloudarea/[0-9]+/?$`)
)

func (ps *parseStream) cloudArea() *parseStream {
	if ps.shouldReturn() {
		return ps
	}

	if ps.hitPattern(findManyCloudAreaPattern, http.MethodPost) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:   meta.CloudAreaInstance,
					Action: meta.SkipAction,
				},
			},
		}
		return ps
	}

	if ps.hitPattern(createCloudAreaPattern, http.MethodPost) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:   meta.CloudAreaInstance,
					Action: meta.Create,
				},
			},
		}
		return ps
	}

	if ps.hitPattern(createManyCloudAreaPattern, http.MethodPost) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:   meta.CloudAreaInstance,
					Action: meta.Create,
				},
			},
		}
		return ps
	}

	if ps.hitRegexp(updateCloudAreaRegexp, http.MethodPut) {
		id, err := strconv.ParseInt(ps.RequestCtx.Elements[4], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("parse cloud id %s failed", ps.RequestCtx.Elements[5])
			return ps
		}
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:       meta.CloudAreaInstance,
					Action:     meta.Update,
					InstanceID: id,
				},
			},
		}
		return ps
	}

	if ps.hitRegexp(deleteCloudAreaRegexp, http.MethodDelete) {

		id, err := strconv.ParseInt(ps.RequestCtx.Elements[4], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("parse cloud id %s failed", ps.RequestCtx.Elements[5])
			return ps
		}

		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:       meta.CloudAreaInstance,
					Action:     meta.Delete,
					InstanceID: id,
				},
			},
		}
		return ps
	}

	if ps.hitPattern(findCloudAreaHostCountPattern, http.MethodPost) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:   meta.CloudAreaInstance,
					Action: meta.SkipAction,
				},
			},
		}
		return ps
	}

	return ps
}
