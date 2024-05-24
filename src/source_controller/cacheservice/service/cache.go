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
	"strconv"
	"strings"
	"time"

	fullsynccond "configcenter/pkg/cache/full-sync-cond"
	"configcenter/pkg/cache/general"
	acmeta "configcenter/src/ac/meta"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	httpheader "configcenter/src/common/http/header"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/common/watch"
	"configcenter/src/source_controller/cacheservice/cache/biz-topo/types"
	customtypes "configcenter/src/source_controller/cacheservice/cache/custom/types"
	"configcenter/src/source_controller/cacheservice/cache/topotree"
	"configcenter/src/source_controller/cacheservice/event"
	"configcenter/src/storage/driver/mongodb"
)

// SearchHostWithInnerIPInCache This function is only used to query the host through ip+cloud in the static IP scenario
// of the host snapshot !!!
func (s *cacheService) SearchHostWithInnerIPInCache(ctx *rest.Contexts) {
	opt := new(metadata.SearchHostWithInnerIPOption)
	if err := ctx.DecodeInto(opt); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if len(opt.InnerIP) == 0 || len(strings.Split(opt.InnerIP, ",")) > 1 {
		blog.Errorf("host inner ip %s is not set or is multiple, rid: %s", opt.InnerIP, ctx.Kit.Rid)
		ctx.RespErrorCodeOnly(common.CCErrCommDBSelectFailed, "host inner ip is not set or is multiple")
		return
	}

	listOpt := &general.ListDetailByUniqueKeyOpt{
		Resource: general.Host,
		Type:     general.IPCloudIDType,
		Keys:     []string{general.IPCloudIDKey(opt.InnerIP, opt.CloudID)},
		Fields:   opt.Fields,
	}
	details, err := s.cacheSet.General.ListDetailByUniqueKey(ctx.Kit, listOpt, true)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	if len(details) != 1 {
		blog.Errorf("host detail %+v is invalid, opt: %+v, rid: %s", details, opt, ctx.Kit.Rid)
		ctx.RespErrorCodeOnly(common.CCErrCommDBSelectFailed, "host detail is invalid")
		return
	}

	ctx.RespString(&details[0])
}

// SearchHostWithAgentIDInCache This function is only used to query host information based on agentID in the host
// snapshot scenario !!!
func (s *cacheService) SearchHostWithAgentIDInCache(ctx *rest.Contexts) {
	opt := new(metadata.SearchHostWithAgentID)
	if err := ctx.DecodeInto(opt); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if len(opt.AgentID) == 0 {
		blog.Errorf("host agent id %s is not set, rid: %s", opt.AgentID, ctx.Kit.Rid)
		ctx.RespErrorCodeOnly(common.CCErrCommDBSelectFailed, "host agent id is not set")
		return
	}

	listOpt := &general.ListDetailByUniqueKeyOpt{
		Resource: general.Host,
		Type:     general.AgentIDType,
		Keys:     []string{general.AgentIDKey(opt.AgentID)},
		Fields:   opt.Fields,
	}
	details, err := s.cacheSet.General.ListDetailByUniqueKey(ctx.Kit, listOpt, true)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	if len(details) != 1 {
		blog.Errorf("host detail %+v is invalid, opt: %+v, rid: %s", details, opt, ctx.Kit.Rid)
		ctx.RespErrorCodeOnly(common.CCErrCommDBSelectFailed, "host detail is invalid")
		return
	}

	ctx.RespString(&details[0])
}

// SearchHostWithHostIDInCache TODO
func (s *cacheService) SearchHostWithHostIDInCache(ctx *rest.Contexts) {
	opt := new(metadata.SearchHostWithIDOption)
	if err := ctx.DecodeInto(opt); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if opt.HostID <= 0 {
		blog.Errorf("host id %d is invalid, rid: %s", opt.HostID, ctx.Kit.Rid)
		ctx.RespErrorCodeOnly(common.CCErrCommDBSelectFailed, "host id is invalid")
		return
	}

	listOpt := &general.ListDetailByIDsOpt{
		Resource: general.Host,
		IDs:      []int64{opt.HostID},
		Fields:   opt.Fields,
	}
	details, err := s.cacheSet.General.ListDetailByIDs(ctx.Kit, listOpt, true)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	if len(details) != 1 {
		blog.Errorf("host detail %+v is invalid, opt: %+v, rid: %s", details, opt, ctx.Kit.Rid)
		ctx.RespErrorCodeOnly(common.CCErrCommDBSelectFailed, "host detail is invalid")
		return
	}

	ctx.RespString(&details[0])
}

// ListHostWithHostIDInCache list hosts info from redis with host id list.
// if a host is not exist in cache and still can not find in mongodb,
// then it will not be return. so the returned array may not equal to
// the request host ids length and the sequence is also may not same.
func (s *cacheService) ListHostWithHostIDInCache(ctx *rest.Contexts) {
	opt := new(metadata.ListWithIDOption)
	if err := ctx.DecodeInto(opt); err != nil {
		ctx.RespAutoError(err)
		return
	}

	details, err := s.listHostWithHostIDInCache(ctx.Kit, opt.IDs, opt.Fields)
	if err != nil {
		ctx.RespErrorCodeOnly(common.CCErrCommDBSelectFailed, err.Error())
		return
	}
	ctx.RespStringArray(details)
}

func (s *cacheService) listHostWithHostIDInCache(kit *rest.Kit, ids []int64, fields []string) ([]string, error) {
	if len(ids) == 0 || len(ids) > 500 {
		blog.Errorf("host id(%+v) is invalid, rid: %s", ids, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, "host ids")
	}

	listOpt := &general.ListDetailByIDsOpt{
		Resource: general.Host,
		IDs:      ids,
		Fields:   fields,
	}
	details, err := s.cacheSet.General.ListDetailByIDs(kit, listOpt, true)
	if err != nil {
		return nil, err
	}

	return details, nil
}

// ListHostWithPageInCache TODO
func (s *cacheService) ListHostWithPageInCache(ctx *rest.Contexts) {
	opt := new(metadata.ListHostWithPage)
	if err := ctx.DecodeInto(opt); err != nil {
		ctx.RespAutoError(err)
		return
	}

	ctx.Kit.Ctx = util.SetDBReadPreference(ctx.Kit.Ctx, common.SecondaryPreferredMode)

	if len(opt.HostIDs) > 0 {
		cntCond := mapstr.MapStr{common.BKHostIDField: mapstr.MapStr{common.BKDBIN: opt.HostIDs}}
		cnt, err := mongodb.Client().Table(common.BKTableNameBaseHost).Find(cntCond).Count(ctx.Kit.Ctx)
		if err != nil {
			blog.Errorf("count host failed, err: %v, cond: %+v, rid: %s", err, cntCond, ctx.Kit.Rid)
			ctx.RespErrorCodeOnly(common.CCErrCommDBSelectFailed, err.Error())
			return
		}

		details, err := s.listHostWithHostIDInCache(ctx.Kit, opt.HostIDs, opt.Fields)
		if err != nil {
			ctx.RespErrorCodeOnly(common.CCErrCommDBSelectFailed, err.Error())
			return
		}

		ctx.RespCountInfoString(int64(cnt), details)
		return
	}

	if opt.Page.Limit == 0 || opt.Page.Limit > common.BKMaxPageSize {
		ctx.RespErrorCodeOnly(common.CCErrCommDBSelectFailed, "page limit is invalid")
		return
	}

	listOpt := &general.ListDetailOpt{
		Resource: general.Host,
		Fields:   opt.Fields,
		Page: &general.PagingOption{
			StartIndex: int64(opt.Page.Start),
			Limit:      int64(opt.Page.Limit),
		},
	}

	_, details, err := s.cacheSet.General.ListData(ctx.Kit, listOpt)
	if err != nil {
		ctx.RespErrorCodeOnly(common.CCErrCommDBSelectFailed, err.Error())
		return
	}

	listOpt.Page = &general.PagingOption{
		EnableCount: true,
	}

	cnt, _, err := s.cacheSet.General.ListData(ctx.Kit, listOpt)
	if err != nil {
		ctx.RespErrorCodeOnly(common.CCErrCommDBSelectFailed, err.Error())
		return
	}

	ctx.RespCountInfoString(cnt, details)
}

// ListBusinessInCache list business with id from cache, if not exist in cache, then get from mongodb directly.
func (s *cacheService) ListBusinessInCache(ctx *rest.Contexts) {
	opt := new(metadata.ListWithIDOption)
	if err := ctx.DecodeInto(&opt); nil != err {
		ctx.RespAutoError(err)
		return
	}

	details, err := s.cacheSet.Business.ListBusiness(ctx.Kit.Ctx, opt)
	if err != nil {
		ctx.RespErrorCodeOnly(common.CCErrCommDBSelectFailed, "list business with id in cache failed, err: %v", err)
		return
	}
	ctx.RespStringArray(details)
}

// ListModulesInCache TODO
// ListModules list modules with id from cache, if not exist in cache, then get from mongodb directly.
func (s *cacheService) ListModulesInCache(ctx *rest.Contexts) {
	opt := new(metadata.ListWithIDOption)
	if err := ctx.DecodeInto(&opt); nil != err {
		ctx.RespAutoError(err)
		return
	}

	details, err := s.cacheSet.Business.ListModules(ctx.Kit.Ctx, opt)
	if err != nil {
		ctx.RespErrorCodeOnly(common.CCErrCommDBSelectFailed, "list modules with id in cache failed, err: %v", err)
		return
	}
	ctx.RespStringArray(details)
}

// ListSetsInCache TODO
// ListSets list sets with id from cache, if not exist in cache, then get from mongodb directly.
func (s *cacheService) ListSetsInCache(ctx *rest.Contexts) {
	opt := new(metadata.ListWithIDOption)
	if err := ctx.DecodeInto(&opt); nil != err {
		ctx.RespAutoError(err)
		return
	}

	details, err := s.cacheSet.Business.ListSets(ctx.Kit.Ctx, opt)
	if err != nil {
		ctx.RespErrorCodeOnly(common.CCErrCommDBSelectFailed, "list sets with id in cache failed, err: %v", err)
		return
	}
	ctx.RespStringArray(details)
}

// SearchBusinessInCache TODO
func (s *cacheService) SearchBusinessInCache(ctx *rest.Contexts) {
	bizID, err := strconv.ParseInt(ctx.Request.PathParameter(common.BKAppIDField), 10, 64)
	if err != nil {
		ctx.RespErrorCodeOnly(common.CCErrCommParamsIsInvalid, "invalid biz id")
		return
	}
	biz, err := s.cacheSet.Business.GetBusiness(ctx.Kit.Ctx, bizID)
	if err != nil {
		ctx.RespErrorCodeOnly(common.CCErrCommDBSelectFailed,
			"search biz with id in cache, but get biz failed, err: %v", err)
		return
	}
	ctx.RespString(&biz)
}

// SearchSetInCache TODO
func (s *cacheService) SearchSetInCache(ctx *rest.Contexts) {
	setID, err := strconv.ParseInt(ctx.Request.PathParameter(common.BKSetIDField), 10, 64)
	if err != nil {
		ctx.RespErrorCodeOnly(common.CCErrCommParamsIsInvalid, "invalid set id")
		return
	}

	set, err := s.cacheSet.Business.GetSet(ctx.Kit.Ctx, setID)
	if err != nil {
		ctx.RespErrorCodeOnly(common.CCErrCommDBSelectFailed, "search set with id in cache failed, err: %v", err)
		return
	}
	ctx.RespString(&set)
}

// SearchModuleInCache TODO
func (s *cacheService) SearchModuleInCache(ctx *rest.Contexts) {
	moduleID, err := strconv.ParseInt(ctx.Request.PathParameter(common.BKModuleIDField), 10, 64)
	if err != nil {
		ctx.RespErrorCodeOnly(common.CCErrCommParamsIsInvalid, "invalid module id")
		return
	}

	module, err := s.cacheSet.Business.GetModuleDetail(ctx.Kit.Ctx, moduleID)
	if err != nil {
		ctx.RespErrorCodeOnly(common.CCErrCommDBSelectFailed, "search module with id in cache failed, err: %v", err)
		return
	}
	ctx.RespString(&module)
}

// SearchCustomLayerInCache TODO
func (s *cacheService) SearchCustomLayerInCache(ctx *rest.Contexts) {
	objID := ctx.Request.PathParameter(common.BKObjIDField)

	instID, err := strconv.ParseInt(ctx.Request.PathParameter(common.BKInstIDField), 10, 64)
	if err != nil {
		ctx.RespErrorCodeOnly(common.CCErrCommParamsIsInvalid, "invalid instance id")
		return
	}

	inst, err := s.cacheSet.Business.GetCustomLevelDetail(ctx.Kit.Ctx, objID, ctx.Kit.SupplierAccount, instID)
	if err != nil {
		ctx.RespErrorCodeOnly(common.CCErrCommDBSelectFailed, "search custom layer with id in cache failed, err: %v",
			err)
		return
	}
	ctx.RespString(&inst)
}

// SearchBizTopologyNodePath is to search biz instance topology node's parent path. eg:
// from itself up to the biz instance, but not contains the node itself.
func (s *cacheService) SearchBizTopologyNodePath(ctx *rest.Contexts) {
	bizID, err := strconv.ParseInt(ctx.Request.PathParameter(common.BKAppIDField), 10, 64)
	if err != nil {
		ctx.RespErrorCodeOnly(common.CCErrCommParamsIsInvalid, "invalid biz id")
		return
	}

	if bizID <= 0 {
		ctx.RespErrorCodeOnly(common.CCErrCommParamsIsInvalid, "invalid biz id")
		return
	}

	opt := new(topotree.SearchNodePathOption)
	if err := ctx.DecodeInto(&opt); nil != err {
		ctx.RespAutoError(err)
		return
	}

	opt.Business = bizID

	paths, err := s.cacheSet.Tree.SearchNodePath(ctx.Kit.Ctx, opt, ctx.Kit.SupplierAccount)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(paths)
}

// SearchBusinessBriefTopology search a business's brief topology from biz to module
// with only required fields.
func (s *cacheService) SearchBusinessBriefTopology(ctx *rest.Contexts) {
	biz := ctx.Request.PathParameter("biz")
	bizID, err := strconv.ParseInt(biz, 10, 64)
	if err != nil {
		ctx.RespErrorCodeOnly(common.CCErrCommParamsInvalid, "search biz topology, got invalid biz id, err: %v", err)
		return
	}

	topo, err := s.cacheSet.Topology.GetBizTopology(ctx.Kit, bizID)
	if err != nil {
		ctx.RespErrorCodeOnly(common.CCErrCommDBSelectFailed, "search biz topology, select db failed, err: %v", err)
		return
	}

	ctx.RespString(topo)
}

// SearchBizTopo search business topology cache info
func (s *cacheService) SearchBizTopo(cts *rest.Contexts) {
	topoType := cts.Request.PathParameter("type")
	if len(topoType) == 0 {
		cts.RespAutoError(cts.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedSet, "type"))
		return
	}

	opt := new(types.GetBizTopoOption)
	if err := cts.DecodeInto(opt); err != nil {
		cts.RespAutoError(err)
		return
	}

	// authorize
	authRes := acmeta.ResourceAttribute{Basic: acmeta.Basic{Type: acmeta.Business, Action: acmeta.ViewBusinessResource,
		InstanceID: opt.BizID}}
	if resp, authorized := s.authManager.Authorize(cts.Kit, authRes); !authorized {
		cts.RespNoAuth(resp)
		return
	}

	res, err := s.cacheSet.Topo.GetBizTopo(cts.Kit, topoType, opt)
	if err != nil {
		cts.RespAutoError(err)
		return
	}
	cts.RespString(res)
}

// RefreshBizTopo refresh business topology cache info
func (s *cacheService) RefreshBizTopo(cts *rest.Contexts) {
	topoType := cts.Request.PathParameter("type")
	if len(topoType) == 0 {
		cts.RespAutoError(cts.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedSet, "type"))
		return
	}

	opt := new(types.RefreshBizTopoOption)
	if err := cts.DecodeInto(opt); err != nil {
		cts.RespAutoError(err)
		return
	}

	// authorize
	authRes := acmeta.ResourceAttribute{Basic: acmeta.Basic{Type: acmeta.Business, Action: acmeta.ViewBusinessResource,
		InstanceID: opt.BizID}}
	if resp, authorized := s.authManager.Authorize(cts.Kit, authRes); !authorized {
		cts.RespNoAuth(resp)
		return
	}

	err := s.cacheSet.Topo.RefreshBizTopo(cts.Kit, topoType, opt)
	if err != nil {
		cts.RespAutoError(err)
		return
	}
	cts.RespEntity(nil)
}

// ListPodLabelKey list pod label key cache info
func (s *cacheService) ListPodLabelKey(cts *rest.Contexts) {
	opt := new(customtypes.ListPodLabelKeyOption)
	if err := cts.DecodeInto(opt); err != nil {
		cts.RespAutoError(err)
		return
	}

	// authorize
	authRes := acmeta.ResourceAttribute{Basic: acmeta.Basic{Type: acmeta.KubePod, Action: acmeta.Find},
		BusinessID: opt.BizID}
	if resp, authorized := s.authManager.Authorize(cts.Kit, authRes); !authorized {
		cts.RespNoAuth(resp)
		return
	}

	res, err := s.cacheSet.Custom.ListPodLabelKey(cts.Kit, opt)
	if err != nil {
		cts.RespAutoError(err)
		return
	}
	cts.RespEntity(&customtypes.ListPodLabelKeyRes{Keys: res})
}

// ListPodLabelValue list pod label value cache info
func (s *cacheService) ListPodLabelValue(cts *rest.Contexts) {
	opt := new(customtypes.ListPodLabelValueOption)
	if err := cts.DecodeInto(opt); err != nil {
		cts.RespAutoError(err)
		return
	}

	// authorize
	authRes := acmeta.ResourceAttribute{Basic: acmeta.Basic{Type: acmeta.KubePod, Action: acmeta.Find},
		BusinessID: opt.BizID}
	if resp, authorized := s.authManager.Authorize(cts.Kit, authRes); !authorized {
		cts.RespNoAuth(resp)
		return
	}

	res, err := s.cacheSet.Custom.ListPodLabelValue(cts.Kit, opt)
	if err != nil {
		cts.RespAutoError(err)
		return
	}
	cts.RespEntity(&customtypes.ListPodLabelValueRes{Values: res})
}

// RefreshPodLabel refresh pod label key and value cache info
func (s *cacheService) RefreshPodLabel(cts *rest.Contexts) {
	opt := new(customtypes.RefreshPodLabelOption)
	if err := cts.DecodeInto(opt); err != nil {
		cts.RespAutoError(err)
		return
	}

	// authorize
	authRes := acmeta.ResourceAttribute{Basic: acmeta.Basic{Type: acmeta.KubePod, Action: acmeta.Find},
		BusinessID: opt.BizID}
	if resp, authorized := s.authManager.Authorize(cts.Kit, authRes); !authorized {
		cts.RespNoAuth(resp)
		return
	}

	err := s.cacheSet.Custom.RefreshPodLabel(cts.Kit, opt)
	if err != nil {
		cts.RespAutoError(err)
		return
	}
	cts.RespEntity(nil)
}

// WatchEvent TODO
func (s *cacheService) WatchEvent(ctx *rest.Contexts) {
	var err error
	// sleep for a while if an error occurred to avoid others using wrong input to request too frequently
	defer func() {
		if err != nil {
			time.Sleep(500 * time.Millisecond)
		}
	}()

	options := new(watch.WatchEventOptions)
	if err = ctx.DecodeInto(&options); err != nil {
		blog.Errorf("watch event, but decode request body failed, err: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	if err = options.Validate(); err != nil {
		blog.Errorf("watch event, but got invalid request options, err: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	key, err := event.GetResourceKeyWithCursorType(options.Resource)
	if err != nil {
		blog.Errorf("watch event, but get resource key with cursor type failed, err: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	// read all data from db in case secondary node's latency causes data inconsistency
	util.SetDBReadPreference(ctx.Kit.Ctx, common.PrimaryMode)

	// watch with cursor
	if len(options.Cursor) != 0 {
		events, err := s.cacheSet.Event.WatchWithCursor(ctx.Kit, key, options)
		if err != nil {
			blog.Errorf("watch event with cursor failed, cursor: %s, err: %v, rid: %s", options.Cursor, err,
				ctx.Kit.Rid)
			ctx.RespAutoError(err)
			return
		}

		// if not events is hit, then we return user's cursor, so that they can watch with this cursor again.
		ctx.RespEntity(s.generateWatchEventResp(options.Cursor, options.Resource, events))
		return
	}

	// watch with start from
	if options.StartFrom != 0 {
		events, err := s.cacheSet.Event.WatchWithStartFrom(ctx.Kit, key, options)
		if err != nil {
			blog.Errorf("watch event with start from: %s failed, err: %v, rid: %s",
				time.Unix(options.StartFrom, 0).Format(time.RFC3339), err, ctx.Kit.Rid)
			ctx.RespAutoError(err)
			return
		}

		ctx.RespEntity(s.generateWatchEventResp("", options.Resource, events))
		return
	}

	// watch from now
	events, err := s.cacheSet.Event.WatchFromNow(ctx.Kit, key, options)
	if err != nil {
		blog.Errorf("watch event from now failed, err: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(s.generateWatchEventResp("", options.Resource, []*watch.WatchEventDetail{events}))
}

func (s *cacheService) generateWatchEventResp(startCursor string, rsc watch.CursorType,
	events []*watch.WatchEventDetail) *watch.WatchResp {

	result := new(watch.WatchResp)
	if len(events) == 0 {
		result.Watched = false
		if len(startCursor) == 0 {
			result.Events = []*watch.WatchEventDetail{
				{
					Cursor:   watch.NoEventCursor,
					Resource: rsc,
				},
			}
		} else {
			// if user's watch with a start cursor, but we do not find event after this cursor,
			// then we return this start cursor directly, so that they can watch with this cursor for next round.
			result.Events = []*watch.WatchEventDetail{
				{
					Cursor:   startCursor,
					Resource: rsc,
				},
			}
		}
	} else {
		if events[0].Cursor == watch.NoEventCursor {
			result.Watched = false

			if len(startCursor) == 0 {
				// user watch with start form time, or watch from now, then return with NoEventCursor cursor.
				result.Events = []*watch.WatchEventDetail{
					{
						Cursor:   watch.NoEventCursor,
						Resource: rsc,
					},
				}
			} else {
				// if user's watch with a start cursor, but hit a NoEventCursor cursor,
				// then we return this start cursor directly, so that they can watch with this cursor for next round.
				result.Events = []*watch.WatchEventDetail{
					{
						Cursor:   startCursor,
						Resource: rsc,
					},
				}
			}
		} else if events[0].Detail == nil {
			// compatible for event happens but not hit(with different event type), last cursor is returned with no detail
			result.Watched = false
			result.Events = []*watch.WatchEventDetail{events[0]}
		} else {
			result.Watched = true
			result.Events = events
		}
	}

	return result
}

// CreateFullSyncCond create full sync cache condition
func (s *cacheService) CreateFullSyncCond(cts *rest.Contexts) {
	opt := new(fullsynccond.CreateFullSyncCondOpt)
	if err := cts.DecodeInto(opt); err != nil {
		cts.RespAutoError(err)
		return
	}

	if rawErr := opt.Validate(); rawErr.ErrCode != 0 {
		blog.Errorf("create full sync cond option is invalid, err: %v, opt: %+v, rid: %s", rawErr, opt, cts.Kit.Rid)
		cts.RespAutoError(rawErr.ToCCError(cts.Kit.CCError))
		return
	}

	// authorize
	authRes := acmeta.ResourceAttribute{Basic: acmeta.Basic{Type: acmeta.FullSyncCond, Action: acmeta.Create}}
	if resp, authorized := s.authManager.Authorize(cts.Kit, authRes); !authorized {
		cts.RespNoAuth(resp)
		return
	}

	id, err := s.cacheSet.General.FullSyncCond().CreateFullSyncCond(cts.Kit, opt)
	if err != nil {
		cts.RespAutoError(err)
		return
	}
	cts.RespEntity(metadata.RspID{ID: id})
}

// UpdateFullSyncCond update full sync cache condition
func (s *cacheService) UpdateFullSyncCond(cts *rest.Contexts) {
	opt := new(fullsynccond.UpdateFullSyncCondOpt)
	if err := cts.DecodeInto(opt); err != nil {
		cts.RespAutoError(err)
		return
	}

	if rawErr := opt.Validate(); rawErr.ErrCode != 0 {
		blog.Errorf("update full sync cond option is invalid, err: %v, opt: %+v, rid: %s", rawErr, opt, cts.Kit.Rid)
		cts.RespAutoError(rawErr.ToCCError(cts.Kit.CCError))
		return
	}

	// authorize
	authRes := acmeta.ResourceAttribute{Basic: acmeta.Basic{Type: acmeta.FullSyncCond, Action: acmeta.Update}}
	if resp, authorized := s.authManager.Authorize(cts.Kit, authRes); !authorized {
		cts.RespNoAuth(resp)
		return
	}

	err := s.cacheSet.General.FullSyncCond().UpdateFullSyncCond(cts.Kit, opt)
	if err != nil {
		cts.RespAutoError(err)
		return
	}
	cts.RespEntity(nil)
}

// DeleteFullSyncCond delete full sync cache condition
func (s *cacheService) DeleteFullSyncCond(cts *rest.Contexts) {
	opt := new(fullsynccond.DeleteFullSyncCondOpt)
	if err := cts.DecodeInto(opt); err != nil {
		cts.RespAutoError(err)
		return
	}

	if rawErr := opt.Validate(); rawErr.ErrCode != 0 {
		blog.Errorf("delete full sync cond option is invalid, err: %v, opt: %+v, rid: %s", rawErr, opt, cts.Kit.Rid)
		cts.RespAutoError(rawErr.ToCCError(cts.Kit.CCError))
		return
	}

	// authorize
	authRes := acmeta.ResourceAttribute{Basic: acmeta.Basic{Type: acmeta.FullSyncCond, Action: acmeta.Delete}}
	if resp, authorized := s.authManager.Authorize(cts.Kit, authRes); !authorized {
		cts.RespNoAuth(resp)
		return
	}

	err := s.cacheSet.General.FullSyncCond().DeleteFullSyncCond(cts.Kit, opt)
	if err != nil {
		cts.RespAutoError(err)
		return
	}
	cts.RespEntity(nil)
}

// ListFullSyncCond list full sync cache condition
func (s *cacheService) ListFullSyncCond(cts *rest.Contexts) {
	opt := new(fullsynccond.ListFullSyncCondOpt)
	if err := cts.DecodeInto(opt); err != nil {
		cts.RespAutoError(err)
		return
	}

	if rawErr := opt.Validate(); rawErr.ErrCode != 0 {
		blog.Errorf("list full sync cond option is invalid, err: %v, opt: %+v, rid: %s", rawErr, opt, cts.Kit.Rid)
		cts.RespAutoError(rawErr.ToCCError(cts.Kit.CCError))
		return
	}

	// authorize
	authRes := acmeta.ResourceAttribute{Basic: acmeta.Basic{Type: acmeta.FullSyncCond, Action: acmeta.Find}}
	if resp, authorized := s.authManager.Authorize(cts.Kit, authRes); !authorized {
		cts.RespNoAuth(resp)
		return
	}

	data, err := s.cacheSet.General.FullSyncCond().ListFullSyncCond(cts.Kit, opt)
	if err != nil {
		cts.RespAutoError(err)
		return
	}
	cts.RespEntity(data)
}

// ListCacheByFullSyncCond list resource cache by full sync condition
func (s *cacheService) ListCacheByFullSyncCond(cts *rest.Contexts) {
	opt := new(fullsynccond.ListCacheByFullSyncCondOpt)
	if err := cts.DecodeInto(opt); err != nil {
		cts.RespAutoError(err)
		return
	}

	if rawErr := opt.Validate(); rawErr.ErrCode != 0 {
		blog.Errorf("list cache by full sync cond option(%+v) is invalid, err: %v, rid: %s", opt, rawErr, cts.Kit.Rid)
		cts.RespAutoError(rawErr.ToCCError(cts.Kit.CCError))
		return
	}

	// get full sync cond
	fullSyncCond, err := s.cacheSet.General.FullSyncCond().GetFullSyncCond(cts.Kit, opt.CondID)
	if err != nil {
		cts.RespAutoError(err)
		return
	}

	// authorize
	authResp, authorized, err := s.authorizeListGeneralCache(cts.Kit, fullSyncCond.Resource)
	if err != nil {
		cts.RespAutoError(err)
		return
	}
	if !authorized {
		cts.RespNoAuth(authResp)
		return
	}

	data, err := s.cacheSet.General.ListCacheByFullSyncCond(cts.Kit, opt, fullSyncCond)
	if err != nil {
		cts.RespAutoError(err)
		return
	}
	cts.RespEntity(&general.ListGeneralCacheRes{Info: data})
}

// ListGeneralCacheByIDs list general resource cache by ids
func (s *cacheService) ListGeneralCacheByIDs(cts *rest.Contexts) {
	opt := new(general.ListDetailByIDsOpt)
	if err := cts.DecodeInto(opt); err != nil {
		cts.RespAutoError(err)
		return
	}

	if rawErr := opt.Validate(); rawErr.ErrCode != 0 {
		blog.Errorf("list general cache by ids option(%+v) is invalid, err: %v, rid: %s", opt, rawErr, cts.Kit.Rid)
		cts.RespAutoError(rawErr.ToCCError(cts.Kit.CCError))
		return
	}

	// authorize
	authResp, authorized, err := s.authorizeListGeneralCache(cts.Kit, opt.Resource)
	if err != nil {
		cts.RespAutoError(err)
		return
	}
	if !authorized {
		cts.RespNoAuth(authResp)
		return
	}

	details, err := s.cacheSet.General.ListDetailByIDs(cts.Kit, opt, false)
	if err != nil {
		cts.RespAutoError(err)
		return
	}
	cts.RespEntity(&general.ListGeneralCacheRes{Info: details})
}

// ListGeneralCacheByUniqueKey list general resource cache by unique keys
func (s *cacheService) ListGeneralCacheByUniqueKey(cts *rest.Contexts) {
	opt := new(general.ListDetailByUniqueKeyOpt)
	if err := cts.DecodeInto(opt); err != nil {
		cts.RespAutoError(err)
		return
	}

	if rawErr := opt.Validate(); rawErr.ErrCode != 0 {
		blog.Errorf("list cache by unique key option(%+v) is invalid, err: %v, rid: %s", opt, rawErr, cts.Kit.Rid)
		cts.RespAutoError(rawErr.ToCCError(cts.Kit.CCError))
		return
	}

	// authorize
	authResp, authorized, err := s.authorizeListGeneralCache(cts.Kit, opt.Resource)
	if err != nil {
		cts.RespAutoError(err)
		return
	}
	if !authorized {
		cts.RespNoAuth(authResp)
		return
	}

	details, err := s.cacheSet.General.ListDetailByUniqueKey(cts.Kit, opt, true)
	if err != nil {
		cts.RespAutoError(err)
		return
	}
	cts.RespEntity(&general.ListGeneralCacheRes{Info: details})
}

// RefreshGeneralResIDList refresh general resource id list cache
func (s *cacheService) RefreshGeneralResIDList(cts *rest.Contexts) {
	opt := new(general.RefreshIDListOpt)
	if err := cts.DecodeInto(opt); err != nil {
		cts.RespAutoError(err)
		return
	}

	if rawErr := opt.Validate(); rawErr.ErrCode != 0 {
		blog.Errorf("refresh id list cache option(%+v) is invalid, err: %v, rid: %s", opt, rawErr, cts.Kit.Rid)
		cts.RespAutoError(rawErr.ToCCError(cts.Kit.CCError))
		return
	}

	fullSyncCond := new(fullsynccond.FullSyncCond)
	// get full sync cond
	if opt.CondID > 0 {
		var err error
		fullSyncCond, err = s.cacheSet.General.FullSyncCond().GetFullSyncCond(cts.Kit, opt.CondID)
		if err != nil {
			cts.RespAutoError(err)
			return
		}

		if opt.SubRes != "" && fullSyncCond.SubResource != opt.SubRes {
			blog.Errorf("full sync cond(%+v) sub res != input(%+v) sub res, rid: %s", fullSyncCond, opt, cts.Kit.Rid)
			cts.RespAutoError(cts.Kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, fullsynccond.SubResField))
			return
		}
	}

	// authorize
	authResp, authorized, err := s.authorizeListGeneralCache(cts.Kit, opt.Resource)
	if err != nil {
		cts.RespAutoError(err)
		return
	}
	if !authorized {
		cts.RespNoAuth(authResp)
		return
	}

	err = s.cacheSet.General.RefreshIDList(cts.Kit, opt, fullSyncCond)
	if err != nil {
		cts.RespAutoError(err)
		return
	}
	cts.RespEntity(nil)
}

// RefreshGeneralResDetailByIDs refresh general resource detail cache by ids
func (s *cacheService) RefreshGeneralResDetailByIDs(cts *rest.Contexts) {
	opt := new(general.RefreshDetailByIDsOpt)
	if err := cts.DecodeInto(opt); err != nil {
		cts.RespAutoError(err)
		return
	}

	if rawErr := opt.Validate(); rawErr.ErrCode != 0 {
		blog.Errorf("refresh detail cache option(%+v) is invalid, err: %v, rid: %s", opt, rawErr, cts.Kit.Rid)
		cts.RespAutoError(rawErr.ToCCError(cts.Kit.CCError))
		return
	}

	// authorize
	authResp, authorized, err := s.authorizeListGeneralCache(cts.Kit, opt.Resource)
	if err != nil {
		cts.RespAutoError(err)
		return
	}
	if !authorized {
		cts.RespNoAuth(authResp)
		return
	}

	err = s.cacheSet.General.RefreshDetailByIDs(cts.Kit, opt)
	if err != nil {
		cts.RespAutoError(err)
		return
	}
	cts.RespEntity(nil)
}

func (s *cacheService) authorizeListGeneralCache(kit *rest.Kit, res general.ResType) (*metadata.BaseResp, bool, error) {
	// skip auth for inner request
	if httpheader.IsInnerReq(kit.Header) {
		return nil, true, nil
	}

	authRes := acmeta.ResourceAttribute{Basic: acmeta.Basic{Type: acmeta.GeneralCache, Action: acmeta.Find,
		InstanceIDEx: string(res)}}
	if resp, authorized := s.authManager.Authorize(kit, authRes); !authorized {
		return resp, authorized, nil
	}

	return nil, true, nil
}
