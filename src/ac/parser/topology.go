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

	"configcenter/src/ac/iam"
	"configcenter/src/ac/meta"
	"configcenter/src/common"
	"configcenter/src/common/json"
	"configcenter/src/common/metadata"
	"configcenter/src/kube/types"
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
		cloudArea().
		businessSet().
		kube()

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

const (
	findReducedBusinessListPattern      = `/api/v3/biz/with_reduced`
	findSimplifiedBusinessListPattern   = `/api/v3/biz/simplify`
	updatemanyBizPropertyPattern        = `/api/v3/updatemany/biz/property`
	deletemanyBizPropertyPattern        = `/api/v3/deletemany/biz`
	updatePlatformSettingIdleSetPattern = `/api/v3/topo/update/biz/idle_set`
	deletePlatformSettingModulePattern  = `/api/v3/topo/delete/biz/extra_moudle`
)

func (ps *parseStream) business() *parseStream {
	if ps.shouldReturn() {
		return ps
	}

	// find reduced business list for the user with any business resources
	if ps.hitPattern(findReducedBusinessListPattern, http.MethodGet) {
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
	if ps.hitPattern(findSimplifiedBusinessListPattern, http.MethodGet) {
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

	// batch update business properties
	if ps.hitPattern(updatemanyBizPropertyPattern, http.MethodPut) {
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

	// batch delete archived businesses
	if ps.hitPattern(deletemanyBizPropertyPattern, http.MethodPost) {
		input := metadata.DeleteBizParam{}
		body, err := ps.RequestCtx.getRequestBody()
		if err != nil {
			ps.err = err
			return ps
		}
		if err := json.Unmarshal(body, &input); err != nil {
			ps.err = fmt.Errorf("unmarshal request body failed, err: %+v", err)
			return ps
		}
		ps.Attribute.Resources = make([]meta.ResourceAttribute, 0)
		for _, bizID := range input.BizID {
			iamResource := meta.ResourceAttribute{
				Basic: meta.Basic{
					Type: meta.Business,
					// delete archived business use archive action
					Action:     meta.Archive,
					InstanceID: bizID,
				},
			}
			ps.Attribute.Resources = append(ps.Attribute.Resources, iamResource)
		}
		return ps
	}

	// find simplified business list with limited fields return
	if ps.hitPattern(updatePlatformSettingIdleSetPattern, http.MethodPost) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:   meta.ConfigAdmin,
					Action: meta.Update,
				},
			},
		}
		return ps
	}

	if ps.hitPattern(deletePlatformSettingModulePattern, http.MethodPost) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:   meta.ConfigAdmin,
					Action: meta.Update,
				},
			},
		}
		return ps
	}

	return ps
}

const (
	createBizSetPattern                  = `/api/v3/create/biz_set`
	deleteBizSetPattern                  = `/api/v3/deletemany/biz_set`
	updateBizSetPattern                  = `/api/v3/updatemany/biz_set`
	findBizInBizSetPattern               = `/api/v3/find/biz_set/biz_list`
	findBizSetTopoPattern                = `/api/v3/find/biz_set/topo_path`
	findmanyBusinessSetPattern           = `/api/v3/findmany/biz_set`
	findReducedBusinessSetListPattern    = `/api/v3/findmany/biz_set/with_reduced`
	previewBusinessSetPattern            = `/api/v3/find/biz_set/preview`
	findSimplifiedBusinessSetListPattern = `/api/v3/findmany/biz_set/simplify`
)

var (
	// search biz resources by biz set regex, authorize by biz set access permission, **only for ui**
	listSetInBizSetRegexp     = regexp.MustCompile(`^/api/v3/findmany/set/biz_set/[0-9]+/biz/[0-9]+/?$`)
	listModuleInBizSetRegexp  = regexp.MustCompile(`^/api/v3/findmany/module/biz_set/[0-9]+/biz/[0-9]+/set/[0-9]+/?$`)
	findBizSetTopoPathRegexp  = regexp.MustCompile(`^/api/v3/find/topopath/biz_set/[0-9]+/biz/[0-9]+/?$`)
	countTopoHostAndSrvRegexp = regexp.MustCompile(`^/api/v3/count/topoinst/host_service_inst/biz_set/[0-9]+/?$`)
)

// businessSet TODO
// NOCC:golint/fnsize(business set操作需要放到一个函数中)
func (ps *parseStream) businessSet() *parseStream {
	if ps.shouldReturn() {
		return ps
	}

	if ps.hitPattern(findBizInBizSetPattern, http.MethodPost) {
		bizSetIDVal, err := ps.RequestCtx.getValueFromBody("bk_biz_set_id")
		if err != nil {
			ps.err = err
			return ps
		}

		bizSetID := bizSetIDVal.Int()
		if bizSetID <= 0 {
			ps.err = errors.New("invalid biz set id")
			return ps
		}

		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:       meta.BizSet,
					Action:     meta.AccessBizSet,
					InstanceID: bizSetID,
				},
			},
		}
		return ps
	}

	// update biz set, authorize if user has update permission to all of the specified biz set ids
	if ps.hitPattern(updateBizSetPattern, http.MethodPut) {
		bizSetIDsVal, err := ps.RequestCtx.getValueFromBody("bk_biz_set_ids")
		if err != nil {
			ps.err = err
			return ps
		}

		bizSetIDArr := bizSetIDsVal.Array()
		if len(bizSetIDArr) == 0 {
			ps.err = errors.New("bk_biz_set_ids is not set")
			return ps
		}

		for _, bizSetIDVal := range bizSetIDArr {
			bizSetID := bizSetIDVal.Int()
			if bizSetID <= 0 {
				ps.err = errors.New("invalid biz set id")
				return ps
			}

			ps.Attribute.Resources = []meta.ResourceAttribute{
				{
					Basic: meta.Basic{
						Type:       meta.BizSet,
						Action:     meta.Update,
						InstanceID: bizSetID,
					},
				},
			}
		}
		return ps
	}

	if ps.hitPattern(deleteBizSetPattern, http.MethodPost) {
		bizSetIDsVal, err := ps.RequestCtx.getValueFromBody("bk_biz_set_ids")
		if err != nil {
			ps.err = err
			return ps
		}

		bizSetIDArr := bizSetIDsVal.Array()
		if len(bizSetIDArr) == 0 {
			ps.err = errors.New("bk_biz_set_ids is not set")
			return ps
		}
		if len(bizSetIDArr) > 100 {
			ps.err = errors.New("bk_biz_set_ids exceeds maximum length 100")
			return ps
		}

		for _, bizSetIDVal := range bizSetIDArr {
			bizSetID := bizSetIDVal.Int()
			if bizSetID <= 0 {
				ps.err = errors.New("invalid biz set id")
				return ps
			}

			ps.Attribute.Resources = []meta.ResourceAttribute{
				{
					Basic: meta.Basic{
						Type:       meta.BizSet,
						Action:     meta.Delete,
						InstanceID: bizSetID,
					},
				},
			}
		}
		return ps
	}

	if ps.hitPattern(findBizSetTopoPattern, http.MethodPost) {
		bizSetIDVal, err := ps.RequestCtx.getValueFromBody("bk_biz_set_id")
		if err != nil {
			ps.err = err
			return ps
		}

		bizSetID := bizSetIDVal.Int()
		if bizSetID <= 0 {
			ps.err = errors.New("invalid biz set id")
			return ps
		}

		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:       meta.BizSet,
					Action:     meta.AccessBizSet,
					InstanceID: bizSetID,
				},
			},
		}
		return ps
	}

	// find many business set list for the user with any business set resources
	if ps.hitPattern(findmanyBusinessSetPattern, http.MethodPost) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:   meta.BizSet,
					Action: meta.SkipAction,
				},
			},
		}
		return ps
	}

	// NOTE: find many business set for front-end use alone.
	if ps.hitPattern(findSimplifiedBusinessSetListPattern, http.MethodGet) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:   meta.BizSet,
					Action: meta.SkipAction,
				},
			},
		}
		return ps
	}

	// find reduced business set list for the user with any business set resources
	if ps.hitPattern(findReducedBusinessSetListPattern, http.MethodGet) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:   meta.BizSet,
					Action: meta.SkipAction,
				},
			},
		}
		return ps
	}

	// create business set, this is not a normalize api.
	if ps.hitPattern(createBizSetPattern, http.MethodPost) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:   meta.BizSet,
					Action: meta.Create,
				},
			},
		}
		return ps
	}

	// preview business set
	if ps.hitPattern(previewBusinessSetPattern, http.MethodPost) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:   meta.BizSet,
					Action: meta.SkipAction,
				},
			},
		}
		return ps
	}

	if ps.hitRegexp(listSetInBizSetRegexp, http.MethodPost) {
		if len(ps.RequestCtx.Elements) != 8 {
			ps.err = fmt.Errorf("get invalid url elements length %d", len(ps.RequestCtx.Elements))
			return ps
		}

		bizSetID, err := strconv.ParseInt(ps.RequestCtx.Elements[5], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("get invalid business set id %s, err: %v", ps.RequestCtx.Elements[5], err)
			return ps
		}

		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:       meta.BizSet,
					Action:     meta.AccessBizSet,
					InstanceID: bizSetID,
				},
			},
		}

		return ps
	}

	if ps.hitRegexp(listModuleInBizSetRegexp, http.MethodPost) {
		if len(ps.RequestCtx.Elements) != 10 {
			ps.err = fmt.Errorf("get invalid url elements length %d", len(ps.RequestCtx.Elements))
			return ps
		}

		bizSetID, err := strconv.ParseInt(ps.RequestCtx.Elements[5], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("get invalid business set id %s, err: %v", ps.RequestCtx.Elements[5], err)
			return ps
		}

		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:       meta.BizSet,
					Action:     meta.AccessBizSet,
					InstanceID: bizSetID,
				},
			},
		}

		return ps
	}

	if ps.hitRegexp(findBizSetTopoPathRegexp, http.MethodPost) {
		if len(ps.RequestCtx.Elements) != 8 {
			ps.err = fmt.Errorf("get invalid url elements length %d", len(ps.RequestCtx.Elements))
			return ps
		}

		bizSetID, err := strconv.ParseInt(ps.RequestCtx.Elements[5], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("get invalid business set id %s, err: %v", ps.RequestCtx.Elements[5], err)
			return ps
		}

		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:       meta.BizSet,
					Action:     meta.AccessBizSet,
					InstanceID: bizSetID,
				},
			},
		}

		return ps
	}

	if ps.hitRegexp(countTopoHostAndSrvRegexp, http.MethodPost) {
		if len(ps.RequestCtx.Elements) != 7 {
			ps.err = fmt.Errorf("get invalid url elements length %d", len(ps.RequestCtx.Elements))
			return ps
		}

		bizSetID, err := strconv.ParseInt(ps.RequestCtx.Elements[6], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("get invalid business set id %s, err: %v", ps.RequestCtx.Elements[6], err)
			return ps
		}

		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:       meta.BizSet,
					Action:     meta.AccessBizSet,
					InstanceID: bizSetID,
				},
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

const (
	findBriefTopologyNodeRelation = "/api/v3/find/topo/biz/brief_node_relation"
	findHostTopoPath              = "/api/v3/find/host/topopath"
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

	if ps.hitPattern(findBriefTopologyNodeRelation, http.MethodPost) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:   meta.MainlineInstance,
					Action: meta.SkipAction,
				},
			},
		}

		return ps
	}

	if ps.hitPattern(findHostTopoPath, http.MethodPost) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:   meta.MainlineInstanceTopology,
					Action: meta.Find,
				},
			},
		}

		return ps
	}
	return ps
}

const (
	objectStatistics = "/api/v3/object/statistics"
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
	createModuleRegexp = regexp.MustCompile(`^/api/v3/module/[0-9]+/[0-9]+/?$`)
	deleteModuleRegexp = regexp.MustCompile(`^/api/v3/module/[0-9]+/[0-9]+/[0-9]+/?$`)
	updateModuleRegexp = regexp.MustCompile(`^/api/v3/module/[0-9]+/[0-9]+/[0-9]+/?$`)

	// NOCC:tosa/linelength(ignore length)
	updateModuleHostApplyStatusRegexp = regexp.MustCompile(`^/api/v3/module/host_apply_enable_status/bk_biz_id/([0-9]+)/?$`)

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
						Type:       meta.ModelSet,
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
						Type:       meta.ModelSet,
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
						Type:       meta.ModelSet,
						Name:       "set",
						InstanceID: setID,
					},
				},
			},
		}
		return ps
	}
	if ps.hitRegexp(updateModuleHostApplyStatusRegexp, http.MethodPut) {
		if len(ps.RequestCtx.Elements) != 6 {
			ps.err = errors.New("update module host apply enabled status, but got invalid url")
			return ps
		}

		bizID, err := strconv.ParseInt(ps.RequestCtx.Elements[5], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("update module host apply enabled status, but got invalid business id %s",
				ps.RequestCtx.Elements[5])
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
	searchInstAudit   = `/api/v3/find/inst_audit`
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

	if ps.hitPattern(searchInstAudit, http.MethodPost) {
		query := new(metadata.InstAuditQueryInput)
		body, err := ps.RequestCtx.getRequestBody()
		if err != nil {
			ps.err = err
			return ps
		}
		if err := json.Unmarshal(body, query); err != nil {
			ps.err = fmt.Errorf("unmarshal request body failed, err: %+v", err)
			return ps
		}

		isMainline, err := ps.isMainlineModel(query.Condition.ObjID)
		if err != nil {
			ps.err = fmt.Errorf("check object is mainline failed, err: %v, rid: %s", err, ps.RequestCtx.Rid)
			return ps
		}

		// authorize logic reference: https://github.com/Tencent/bk-cmdb/issues/5758
		if isMainline {
			if query.Condition.BizID == 0 {
				ps.err = fmt.Errorf("bk_biz_id is invalid, rid: %s", ps.RequestCtx.Rid)
				return ps
			}

			resPoolBizID, err := ps.getResourcePoolBusinessID()
			if err != nil {
				ps.err = fmt.Errorf("get resource pool failed, err: %v, rid: %s", err, ps.RequestCtx.Rid)
				return ps
			}

			if query.Condition.ObjID == common.BKInnerObjIDHost && query.Condition.BizID == resPoolBizID {

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

			ps.Attribute.Resources = []meta.ResourceAttribute{
				{
					Basic: meta.Basic{
						Type:       meta.Business,
						InstanceID: query.Condition.BizID,
						Action:     meta.ViewBusinessResource,
					},
				},
			}

			return ps
		}

		model, err := ps.getOneModel(map[string]interface{}{common.BKObjIDField: query.Condition.ObjID})
		if err != nil {
			ps.err = err
			return ps
		}

		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:   iam.GenCMDBDynamicResType(model.ID),
					Action: meta.SkipAction,
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

const (
	findNodePathForHostPattern = "/api/v3/find/kube/host_node_path"
)

var (
	findKubeAttrsRegexp = regexp.MustCompile(`^/api/v3/find/kube/[^\s/]+/attributes$`)

	createKubeClusterRegexp     = regexp.MustCompile(`^/api/v3/create/kube/cluster/bk_biz_id/([0-9]+)$`)
	deleteKubeClustersRegexp    = regexp.MustCompile(`^/api/v3/delete/kube/cluster/bk_biz_id/([0-9]+)$`)
	findKubeClusterRegexp       = regexp.MustCompile(`^/api/v3/findmany/kube/cluster/bk_biz_id/([0-9]+)$`)
	updatemanyKubeClusterRegexp = regexp.MustCompile(`^/api/v3/updatemany/kube/cluster/bk_biz_id/([0-9]+)$`)

	createKubeNodeRegexp     = regexp.MustCompile(`^/api/v3/createmany/kube/node/bk_biz_id/([0-9]+)$`)
	findKubeNodeRegexp       = regexp.MustCompile(`^/api/v3/findmany/kube/node/bk_biz_id/([0-9]+)$`)
	deleteKubeNodeRegexp     = regexp.MustCompile(`^/api/v3/deletemany/kube/node/bk_biz_id/([0-9]+)$`)
	updatemanyKubeNodeRegexp = regexp.MustCompile(`^/api/v3/updatemany/kube/node/bk_biz_id/([0-9]+)$`)

	findKubeTopoPathRegexp  = regexp.MustCompile(`^/api/v3/find/kube/topo_path/bk_biz_id/([0-9]+)$`)
	findKubeTopoCountRegexp = regexp.MustCompile(`^/api/v3/find/kube/([0-9]+)/topo_node/[^\s/]+/count$`)
	createKubePodsRegexp    = regexp.MustCompile(`^/api/v3/createmany/kube/pod`)

	createNamespaceRegexp = regexp.MustCompile(`^/api/v3/createmany/kube/namespace/bk_biz_id/([0-9]+)/?$`)
	updateNamespaceRegexp = regexp.MustCompile(`^/api/v3/updatemany/kube/namespace/bk_biz_id/([0-9]+)/?$`)
	deleteNamespaceRegexp = regexp.MustCompile(`^/api/v3/deletemany/kube/namespace/bk_biz_id/([0-9]+)/?$`)
	findNamespaceRegexp   = regexp.MustCompile(`^/api/v3/findmany/kube/namespace/bk_biz_id/([0-9]+)/?$`)

	createWorkloadRegexp = regexp.MustCompile(`^/api/v3/createmany/kube/workload/[^\s/]+/[0-9]+/?$`)
	updateWorkloadRegexp = regexp.MustCompile(`^/api/v3/updatemany/kube/workload/[^\s/]+/[0-9]+/?$`)
	deleteWorkloadRegexp = regexp.MustCompile(`^/api/v3/deletemany/kube/workload/[^\s/]+/[0-9]+/?$`)
	findWorkloadRegexp   = regexp.MustCompile(`^/api/v3/findmany/kube/workload/[^\s/]+/[0-9]+/?$`)

	findPodPathRegexp = regexp.MustCompile(`^/api/v3/find/kube/pod_path/bk_biz_id/([0-9]+)/?$`)
	findPodRegexp     = regexp.MustCompile(`^/api/v3/findmany/kube/pod/bk_biz_id/([0-9]+)/?$`)

	findContainerRegexp = regexp.MustCompile(`^/api/v3/findmany/kube/container/bk_biz_id/([0-9]+)/?$`)
)

// NOCC:golint/fnsize(整体属于 container 操作需要放在一起)
func (ps *parseStream) kube() *parseStream {
	if ps.shouldReturn() {
		return ps
	}

	if ps.hitRegexp(findKubeAttrsRegexp, http.MethodGet) {
		if len(ps.RequestCtx.Elements) != 6 {
			ps.err = fmt.Errorf("get invalid url elements length %d", len(ps.RequestCtx.Elements))
			return ps
		}

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

	if ps.hitRegexp(createKubeClusterRegexp, http.MethodPost) {
		if len(ps.RequestCtx.Elements) != 7 {
			ps.err = fmt.Errorf("get invalid url elements length %d", len(ps.RequestCtx.Elements))
			return ps
		}

		bizID, err := strconv.ParseInt(ps.RequestCtx.Elements[6], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("get invalid business set id %s, err: %v", ps.RequestCtx.Elements[6], err)
			return ps
		}
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:   meta.KubeCluster,
					Action: meta.Create,
				},
			},
		}
		return ps
	}

	if ps.hitRegexp(deleteKubeClustersRegexp, http.MethodDelete) {
		if len(ps.RequestCtx.Elements) != 7 {
			ps.err = fmt.Errorf("get invalid url elements length %d", len(ps.RequestCtx.Elements))
			return ps
		}

		bizID, err := strconv.ParseInt(ps.RequestCtx.Elements[6], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("get invalid business set id %s, err: %v", ps.RequestCtx.Elements[6], err)
			return ps
		}
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:   meta.KubeCluster,
					Action: meta.Delete,
				},
			},
		}
		return ps
	}

	if ps.hitRegexp(updatemanyKubeClusterRegexp, http.MethodPut) {
		if len(ps.RequestCtx.Elements) != 7 {
			ps.err = fmt.Errorf("get invalid url elements length %d", len(ps.RequestCtx.Elements))
			return ps
		}

		bizID, err := strconv.ParseInt(ps.RequestCtx.Elements[6], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("get invalid business set id %s, err: %v", ps.RequestCtx.Elements[6], err)
			return ps
		}
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:   meta.KubeCluster,
					Action: meta.UpdateMany,
				},
			},
		}
		return ps
	}

	if ps.hitRegexp(findKubeClusterRegexp, http.MethodPost) {
		if len(ps.RequestCtx.Elements) != 7 {
			ps.err = fmt.Errorf("get invalid url elements length %d", len(ps.RequestCtx.Elements))
			return ps
		}

		bizID, err := strconv.ParseInt(ps.RequestCtx.Elements[6], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("get invalid business set id %s, err: %v", ps.RequestCtx.Elements[6], err)
			return ps
		}
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:   meta.KubeCluster,
					Action: meta.FindMany,
				},
			},
		}
		return ps
	}

	if ps.hitRegexp(createKubeNodeRegexp, http.MethodPost) {
		if len(ps.RequestCtx.Elements) != 7 {
			ps.err = fmt.Errorf("get invalid url elements length %d", len(ps.RequestCtx.Elements))
			return ps
		}

		bizID, err := strconv.ParseInt(ps.RequestCtx.Elements[6], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("get invalid business set id %s, err: %v", ps.RequestCtx.Elements[6], err)
			return ps
		}
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:   meta.KubeNode,
					Action: meta.CreateMany,
				},
			},
		}
		return ps
	}

	if ps.hitRegexp(createKubePodsRegexp, http.MethodPost) {

		option := new(types.CreatePodsOption)
		body, err := ps.RequestCtx.getRequestBody()
		if err != nil {
			ps.err = err
			return ps
		}
		if err := json.Unmarshal(body, option); err != nil {
			ps.err = fmt.Errorf("unmarshal request body failed, err: %+v", err)
			return ps
		}

		for _, data := range option.Data {
			ps.Attribute.Resources = []meta.ResourceAttribute{
				{
					BusinessID: data.BizID,
					Basic: meta.Basic{
						Type:   meta.KubePod,
						Action: meta.CreateMany,
					},
				},
			}
		}

		return ps
	}

	if ps.hitRegexp(findKubeNodeRegexp, http.MethodPost) {
		if len(ps.RequestCtx.Elements) != 7 {
			ps.err = fmt.Errorf("get invalid url elements length %d", len(ps.RequestCtx.Elements))
			return ps
		}

		bizID, err := strconv.ParseInt(ps.RequestCtx.Elements[6], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("get invalid business set id %s, err: %v", ps.RequestCtx.Elements[6], err)
			return ps
		}
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:   meta.KubeNode,
					Action: meta.FindMany,
				},
			},
		}

		return ps
	}

	if ps.hitRegexp(deleteKubeNodeRegexp, http.MethodDelete) {
		if len(ps.RequestCtx.Elements) != 7 {
			ps.err = fmt.Errorf("get invalid url elements length %d", len(ps.RequestCtx.Elements))
			return ps
		}

		bizID, err := strconv.ParseInt(ps.RequestCtx.Elements[6], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("get invalid business set id %s, err: %v", ps.RequestCtx.Elements[6], err)
			return ps
		}
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:   meta.KubeNode,
					Action: meta.DeleteMany,
				},
			},
		}

		return ps
	}

	if ps.hitRegexp(updatemanyKubeNodeRegexp, http.MethodPut) {
		if len(ps.RequestCtx.Elements) != 7 {
			ps.err = fmt.Errorf("get invalid url elements length %d", len(ps.RequestCtx.Elements))
			return ps
		}

		bizID, err := strconv.ParseInt(ps.RequestCtx.Elements[6], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("get invalid business set id %s, err: %v", ps.RequestCtx.Elements[6], err)
			return ps
		}
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:   meta.KubeNode,
					Action: meta.UpdateMany,
				},
			},
		}
		return ps
	}

	if ps.hitRegexp(findKubeTopoPathRegexp, http.MethodPost) {
		if len(ps.RequestCtx.Elements) != 7 {
			ps.err = fmt.Errorf("get invalid url elements length %d", len(ps.RequestCtx.Elements))
			return ps
		}

		bizID, err := strconv.ParseInt(ps.RequestCtx.Elements[6], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("get invalid business set id %s, err: %v", ps.RequestCtx.Elements[6], err)
			return ps
		}
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:       meta.Business,
					InstanceID: bizID,
					Action:     meta.ViewBusinessResource,
				},
			},
		}
		return ps
	}

	if ps.hitRegexp(findKubeTopoCountRegexp, http.MethodPost) {
		if len(ps.RequestCtx.Elements) != 8 {
			ps.err = fmt.Errorf("get invalid url elements length %d", len(ps.RequestCtx.Elements))
			return ps
		}
		bizID, err := strconv.ParseInt(ps.RequestCtx.Elements[4], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("get invalid business id %s, err: %v", ps.RequestCtx.Elements[4], err)
			return ps
		}

		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:       meta.Business,
					InstanceID: bizID,
					Action:     meta.ViewBusinessResource,
				},
			},
		}

		return ps
	}
	if ps.shouldReturn() {
		return ps
	}

	if ps.hitPattern(findNodePathForHostPattern, http.MethodPost) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:   meta.KubeNode,
					Action: meta.Find,
				},
			},
		}
		return ps
	}

	if ps.hitRegexp(createNamespaceRegexp, http.MethodPost) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:   meta.KubeNamespace,
					Action: meta.Create,
				},
			},
		}
		return ps
	}

	if ps.hitRegexp(updateNamespaceRegexp, http.MethodPut) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:   meta.KubeNamespace,
					Action: meta.Update,
				},
			},
		}
		return ps
	}

	if ps.hitRegexp(deleteNamespaceRegexp, http.MethodDelete) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:   meta.KubeNamespace,
					Action: meta.Delete,
				},
			},
		}
		return ps
	}

	if ps.hitRegexp(findNamespaceRegexp, http.MethodPost) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:   meta.KubeNamespace,
					Action: meta.Find,
				},
			},
		}
		return ps
	}

	if ps.hitRegexp(createWorkloadRegexp, http.MethodPost) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:   meta.KubeWorkload,
					Action: meta.Create,
				},
			},
		}
		return ps
	}

	if ps.hitRegexp(updateWorkloadRegexp, http.MethodPut) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:   meta.KubeWorkload,
					Action: meta.Update,
				},
			},
		}
		return ps
	}

	if ps.hitRegexp(deleteWorkloadRegexp, http.MethodDelete) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:   meta.KubeWorkload,
					Action: meta.Delete,
				},
			},
		}
		return ps
	}

	if ps.hitRegexp(findWorkloadRegexp, http.MethodPost) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:   meta.KubeWorkload,
					Action: meta.Find,
				},
			},
		}
		return ps
	}

	if ps.hitRegexp(findPodPathRegexp, http.MethodPost) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:   meta.KubePod,
					Action: meta.Find,
				},
			},
		}
		return ps
	}

	if ps.hitRegexp(findPodRegexp, http.MethodPost) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:   meta.KubePod,
					Action: meta.Find,
				},
			},
		}
		return ps
	}

	if ps.hitRegexp(findContainerRegexp, http.MethodPost) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:   meta.KubeContainer,
					Action: meta.Find,
				},
			},
		}
		return ps
	}

	return ps
}
