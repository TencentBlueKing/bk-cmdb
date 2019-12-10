package parser

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	"configcenter/src/auth/meta"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/json"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/framework/core/errors"

	"github.com/tidwall/gjson"
)

func (ps *parseStream) hostRelated() *parseStream {
	if ps.shouldReturn() {
		return ps
	}

	ps.host().
		userAPI().
		userCustom().
		hostFavorite().
		cloudResourceSync().
		hostSnapshot().
		findObjectIdentifier()

	return ps
}

var (
	createUserAPIPattern     = "/api/v3/userapi"
	updateUserAPIRegexp      = regexp.MustCompile(`^/api/v3/userapi/[0-9]+/[^\s/]+/?$`)
	deleteUserAPIRegexp      = regexp.MustCompile(`^/api/v3/userapi/[0-9]+/[^\s/]+/?$`)
	findUserAPIRegexp        = regexp.MustCompile(`^/api/v3/userapi/search/[0-9]+/?$`)
	findUserAPIDetailsRegexp = regexp.MustCompile(`^/api/v3/userapi/detail/[0-9]+/[^\s/]+/?$`)
	findWithUserAPIRegexp    = regexp.MustCompile(`^/api/v3/userapi/data/[0-9]+/[^\s/]+/[0-9]+/[0-9]+/?$`)
)

func (ps *parseStream) parseBusinessID() (int64, error) {
	if !gjson.GetBytes(ps.RequestCtx.Body, common.BKAppIDField).Exists() {
		return 0, nil
	}
	bizID := gjson.GetBytes(ps.RequestCtx.Body, common.BKAppIDField).Int()
	if bizID == 0 {
		return 0, errors.New("invalid bk_biz_id value")
	}
	return bizID, nil
}

func (ps *parseStream) userAPI() *parseStream {
	if ps.shouldReturn() {
		return ps
	}

	// create user custom query operation.
	if ps.hitPattern(createUserAPIPattern, http.MethodPost) {
		bizID, err := ps.parseBusinessID()
		if err != nil {
			ps.err = err
			return ps
		}
		ps.Attribute.Resources = []meta.ResourceAttribute{
			meta.ResourceAttribute{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:   meta.DynamicGrouping,
					Action: meta.Create,
				},
			},
		}
		return ps
	}

	// update host user custom query operation.
	if ps.hitRegexp(updateUserAPIRegexp, http.MethodPut) {
		if len(ps.RequestCtx.Elements) != 5 {
			ps.err = errors.New("update host user custom query, but got invalid uri")
			return ps
		}
		bizID, err := strconv.ParseInt(ps.RequestCtx.Elements[3], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("update host user custom query failed, err: %v", err)
			return ps
		}
		ps.Attribute.Resources = []meta.ResourceAttribute{
			meta.ResourceAttribute{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:         meta.DynamicGrouping,
					Action:       meta.Update,
					InstanceIDEx: ps.RequestCtx.Elements[4],
				},
			},
		}
		return ps

	}

	// delete host user custom query operation.
	if ps.hitRegexp(deleteUserAPIRegexp, http.MethodDelete) {
		if len(ps.RequestCtx.Elements) != 5 {
			ps.err = errors.New("delete host user custom query operation, but got invalid uri")
			return ps
		}
		bizID, err := strconv.ParseInt(ps.RequestCtx.Elements[3], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("update host user custom query failed, err: %v", err)
			return ps
		}
		ps.Attribute.Resources = []meta.ResourceAttribute{
			meta.ResourceAttribute{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:         meta.DynamicGrouping,
					Action:       meta.Delete,
					InstanceIDEx: ps.RequestCtx.Elements[4],
				},
			},
		}
		return ps

	}

	// find host user custom query operation
	if ps.hitRegexp(findUserAPIRegexp, http.MethodPost) {
		if len(ps.RequestCtx.Elements) != 5 {
			ps.err = errors.New("find host usr custom query, but got invalid uri")
			return ps
		}

		bizID, err := strconv.ParseInt(ps.RequestCtx.Elements[4], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("find host user custom query failed, err: %v", err)
			return ps
		}
		ps.Attribute.Resources = []meta.ResourceAttribute{
			meta.ResourceAttribute{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:   meta.DynamicGrouping,
					Action: meta.FindMany,
				},
			},
		}
		return ps
	}

	// find host user custom query details operation.
	if ps.hitRegexp(findUserAPIDetailsRegexp, http.MethodGet) {
		if len(ps.RequestCtx.Elements) != 6 {
			ps.err = errors.New("find host user custom details query, but got invalid uri")
			return ps
		}

		bizID, err := strconv.ParseInt(ps.RequestCtx.Elements[4], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("find host user custom query details failed, err: %v", err)
			return ps
		}
		ps.Attribute.Resources = []meta.ResourceAttribute{
			meta.ResourceAttribute{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:         meta.DynamicGrouping,
					Action:       meta.Find,
					InstanceIDEx: ps.RequestCtx.Elements[5],
				},
			},
		}
		return ps
	}

	// get data with user custom query api.
	if ps.hitRegexp(findWithUserAPIRegexp, http.MethodGet) {
		if len(ps.RequestCtx.Elements) != 8 {
			ps.err = errors.New("find host user custom details query, but got invalid uri")
			return ps
		}

		bizID, err := strconv.ParseInt(ps.RequestCtx.Elements[4], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("find host user custom query details failed, err: %v", err)
			return ps
		}

		ps.Attribute.Resources = []meta.ResourceAttribute{
			meta.ResourceAttribute{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:   meta.DynamicGrouping,
					Action: meta.Execute,
					Name:   ps.RequestCtx.Elements[5],
				},
			},
		}
		return ps
	}

	return ps
}

var (
	saveUserCustomPattern       = `/api/v3/usercustom`
	searchUserCustomPattern     = `/api/v3/usercustom/user/search`
	getUserDefaultCustomPattern = `/api/v3/usercustom/default/search`
)

func (ps *parseStream) userCustom() *parseStream {
	if ps.shouldReturn() {
		return ps
	}

	// create user custom query operation.
	if ps.hitPattern(saveUserCustomPattern, http.MethodPost) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			meta.ResourceAttribute{
				Basic: meta.Basic{
					Type:   meta.UserCustom,
					Action: meta.Create,
				},
			},
		}
		return ps
	}

	// update host user custom query operation.
	if ps.hitPattern(searchUserCustomPattern, http.MethodPost) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			meta.ResourceAttribute{
				Basic: meta.Basic{
					Type:   meta.UserCustom,
					Action: meta.Find,
				},
			},
		}
		return ps

	}

	// delete host user custom query operation.
	if ps.hitPattern(getUserDefaultCustomPattern, http.MethodPost) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			meta.ResourceAttribute{
				Basic: meta.Basic{
					Type:   meta.UserCustom,
					Action: meta.Find,
				},
			},
		}
		return ps

	}
	return ps
}

const (
	deleteHostBatchPattern                    = "/api/v3/hosts/batch"
	addHostsToHostPoolPattern                 = "/api/v3/hosts/add"
	moveHostToBusinessModulePattern           = "/api/v3/hosts/modules"
	moveResPoolToBizIdleModulePattern         = "/api/v3/hosts/modules/resource/idle"
	moveHostsToBizFaultModulePattern          = "/api/v3/hosts/modules/fault"
	moveHostsFromModuleToResPoolPattern       = "/api/v3/hosts/modules/resource"
	moveHostsToBizIdleModulePattern           = "/api/v3/hosts/modules/idle"
	moveHostsToBizRecycleModulePattern        = "/api/v3/hosts/modules/recycle"
	moveHostsFromOneToAnotherBizModulePattern = "/api/v3/hosts/modules/biz/mutilple"
	moveHostsFromRscPoolToAppModule           = "/api/v3/hosts/host/add/module"
	cleanHostInSetOrModulePattern             = "/api/v3/hosts/modules/idle/set"
	findHostTopoRelationPattern               = "/api/v3/host/topo/relation/read"
	updateHostCloudAreaFieldPattern           = "/api/v3/updatemany/hosts/cloudarea_field"

	// used in sync framework.
	moveHostToBusinessOrModulePattern = "/api/v3/hosts/sync/new/host"
	findHostsWithConditionPattern     = "/api/v3/hosts/search"
	findBizHostsWithoutAppPattern     = "/api/v3/hosts/list_hosts_without_app"
	findHostsDetailsPattern           = "/api/v3/hosts/search/asstdetail"
	updateHostInfoBatchPattern        = "/api/v3/hosts/batch"
	updateHostPropertyBatchPattern    = "/api/v3/hosts/property/batch"
	findHostsWithModulesPattern       = "/api/v3/findmany/modulehost"

	// 特殊接口，给蓝鲸业务使用
	hostInstallPattern = "/api/v3/host/install/bk"
)

var (
	findBizHostsRegex     = regexp.MustCompile(`/api/v3/hosts/app/\d+/list_hosts`)
	findBizHostsTopoRegex = regexp.MustCompile(`/api/v3/hosts/app/\d+/list_hosts_topo`)
	// find host instance's object properties info
	findHostInstanceObjectPropertiesRegexp = regexp.MustCompile(`^/api/v3/hosts/[^\s/]+/[0-9]+/?$`)

	transferHostWithAutoClearServiceInstanceRegex        = regexp.MustCompile("^/api/v3/host/transfer_with_auto_clear_service_instance/bk_biz_id/[0-9]+/?$")
	transferHostWithAutoClearServiceInstancePreviewRegex = regexp.MustCompile("^/api/v3/host/transfer_with_auto_clear_service_instance/bk_biz_id/[0-9]+/preview/?$")

	countHostByTopoNodeRegexp = regexp.MustCompile(`^/api/v3/host/count_by_topo_node/bk_biz_id/[0-9]+$`)
)

func (ps *parseStream) host() *parseStream {
	if ps.shouldReturn() {
		return ps
	}

	if ps.hitRegexp(findHostInstanceObjectPropertiesRegexp, http.MethodGet) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			meta.ResourceAttribute{
				Basic: meta.Basic{
					Type:   meta.HostInstance,
					Action: meta.Find,
				},
			},
		}
		return ps
	}

	if ps.hitPattern(findHostsWithModulesPattern, http.MethodPost) {
		bizID := gjson.GetBytes(ps.RequestCtx.Body, common.BKAppIDField).Int()
		if bizID == 0 {
			var err error
			bizID, err = metadata.BizIDFromMetadata(ps.RequestCtx.Metadata)
			if err != nil {
				ps.err = fmt.Errorf("find hosts with modules, but parse business id failed, err: %v", err)
				return ps
			}
		}
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:   meta.HostInstance,
					Action: meta.FindMany,
				},
			},
		}
		return ps
	}

	// TODO: add host lock authorize filter if needed.

	// delete hosts batch operation.
	if ps.hitPattern(deleteHostBatchPattern, http.MethodDelete) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			meta.ResourceAttribute{
				Basic: meta.Basic{
					Type: meta.HostInstance,
					// Action: meta.DeleteMany,
					Action: meta.SkipAction,
				},
			},
		}
		return ps
	}

	if ps.hitPattern(updateHostCloudAreaFieldPattern, http.MethodPut) {
		input := metadata.UpdateHostCloudAreaFieldOption{}
		if err := json.Unmarshal(ps.RequestCtx.Body, &input); err != nil {
			ps.err = fmt.Errorf("unmarshal request body failed, err: %+v", err)
			return ps
		}

		ps.Attribute.Resources = make([]meta.ResourceAttribute, 0)
		for _, hostID := range input.HostIDs {
			iamResource := meta.ResourceAttribute{
				Basic: meta.Basic{
					Type:       meta.HostInstance,
					Action:     meta.UpdateMany,
					InstanceID: hostID,
				},
				BusinessID: input.BizID,
			}
			ps.Attribute.Resources = append(ps.Attribute.Resources, iamResource)
		}
		return ps
	}

	// TODO: add host clone authorize filter if needed.

	// add new hosts to resource pool
	if ps.hitPattern(addHostsToHostPoolPattern, http.MethodPost) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			meta.ResourceAttribute{
				Basic: meta.Basic{
					Type:   meta.HostInstance,
					Action: meta.AddHostToResourcePool,
				},
			},
		}

		return ps
	}

	// move hosts from a module to resource pool.
	if ps.hitPattern(moveHostsFromModuleToResPoolPattern, http.MethodPost) {
		bizID, err := ps.parseBusinessID()
		if err != nil {
			ps.err = err
			return ps
		}
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:   meta.HostInstance,
					Action: meta.SkipAction,
				},
			},
		}

		return ps
	}

	// move hosts to business module operation.
	// skip this auth operation, it has already been done
	// in host server.
	if ps.hitPattern(moveHostToBusinessModulePattern, http.MethodPost) {
		bizID, err := ps.parseBusinessID()
		if err != nil {
			ps.err = err
			return ps
		}
		ps.Attribute.Resources = []meta.ResourceAttribute{
			meta.ResourceAttribute{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:   meta.HostInstance,
					Action: meta.SkipAction,
				},
			},
		}

		return ps
	}

	// move resource pool hosts to a business idle module operation.
	// authcenter: system->host/resource_pool->edit
	if ps.hitPattern(moveResPoolToBizIdleModulePattern, http.MethodPost) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:   meta.HostInstance,
					Action: meta.Update,
				},
			},
		}
		return ps
	}

	// move host to a business fault module.
	if ps.hitPattern(moveHostsToBizFaultModulePattern, http.MethodPost) {
		bizID, err := ps.parseBusinessID()
		if err != nil {
			ps.err = err
			return ps
		}
		ps.Attribute.Resources = []meta.ResourceAttribute{
			meta.ResourceAttribute{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type: meta.HostInstance,
					// auth this resource in scene layer, as is host server
					Action: meta.SkipAction,
				},
			},
		}

		return ps
	}

	// move host to a business recycle module.
	if ps.hitPattern(moveHostsToBizRecycleModulePattern, http.MethodPost) {
		bizID, err := ps.parseBusinessID()
		if err != nil {
			ps.err = err
			return ps
		}
		ps.Attribute.Resources = []meta.ResourceAttribute{
			meta.ResourceAttribute{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type: meta.HostInstance,
					// auth this resource in scene layer, as is host server
					Action: meta.SkipAction,
				},
			},
		}

		return ps
	}

	// move hosts to a business idle module.
	if ps.hitPattern(moveHostsToBizIdleModulePattern, http.MethodPost) {
		bizID, err := ps.parseBusinessID()
		if err != nil {
			ps.err = err
			return ps
		}
		ps.Attribute.Resources = []meta.ResourceAttribute{
			meta.ResourceAttribute{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type: meta.HostInstance,
					// auth this resource in host server.
					Action: meta.SkipAction,
				},
			},
		}

		return ps
	}

	// move hosts from one business module to another business module.
	if ps.hitPattern(moveHostsFromOneToAnotherBizModulePattern, http.MethodPost) {
		bizID, err := ps.parseBusinessID()
		if err != nil {
			ps.err = err
			return ps
		}
		ps.Attribute.Resources = []meta.ResourceAttribute{
			meta.ResourceAttribute{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:   meta.HostInstance,
					Action: meta.MoveHostToAnotherBizModule,
				},
			},
		}

		return ps
	}

	if ps.hitPattern(moveHostsFromRscPoolToAppModule, http.MethodPost) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			meta.ResourceAttribute{
				Basic: meta.Basic{
					Type: meta.HostInstance,
					// Action: meta.MoveResPoolHostToBizIdleModule,
					// auth this resource in scene layer
					Action: meta.SkipAction,
				},
			},
		}

		return ps
	}

	// clean the hosts in a set or module, and move these hosts to the business idle module
	// when these hosts only exist in this set or module. otherwise these hosts will only be
	// removed from this set or module.
	if ps.hitPattern(cleanHostInSetOrModulePattern, http.MethodPost) {
		bizID, err := ps.parseBusinessID()
		if err != nil {
			ps.err = err
			return ps
		}
		ps.Attribute.Resources = []meta.ResourceAttribute{
			meta.ResourceAttribute{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:   meta.HostInstance,
					Action: meta.CleanHostInSetOrModule,
				},
			},
		}

		return ps
	}

	if ps.hitPattern(findHostTopoRelationPattern, http.MethodPost) {
		bizID, err := ps.parseBusinessID()
		if err != nil {
			ps.err = err
			return ps
		}
		ps.Attribute.Resources = []meta.ResourceAttribute{
			meta.ResourceAttribute{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:   meta.HostInstance,
					Action: meta.FindMany,
				},
			},
		}

		return ps
	}

	// synchronize hosts directly to a module in a business if this host does not exist.
	// otherwise, this operation will only change host's attribute.
	if ps.hitPattern(moveHostToBusinessOrModulePattern, http.MethodPost) {
		bizID, err := ps.parseBusinessID()
		if err != nil {
			ps.err = err
			return ps
		}
		ps.Attribute.Resources = []meta.ResourceAttribute{
			meta.ResourceAttribute{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:   meta.HostInstance,
					Action: meta.MoveHostsToBusinessOrModule,
				},
			},
		}

		return ps
	}

	if ps.hitRegexp(transferHostWithAutoClearServiceInstanceRegex, http.MethodPost) ||
		ps.hitRegexp(transferHostWithAutoClearServiceInstancePreviewRegex, http.MethodPost) {
		match := BizIDRegex.FindStringSubmatch(ps.RequestCtx.URI)
		bizID, err := util.GetInt64ByInterface(match[1])
		if err != nil {
			blog.Errorf("get business id from request path failed, name: %s, err: %v, rid: %s", transferHostWithAutoClearServiceInstanceRegex, err, ps.RequestCtx.Rid)
			ps.err = fmt.Errorf("parse biz id from url failed, err: %s", err.Error())
			return ps
		}
		ps.Attribute.Resources = []meta.ResourceAttribute{
			meta.ResourceAttribute{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:   meta.HostInstance,
					Action: meta.MoveHostsToBusinessOrModule,
				},
			},
		}

		return ps
	}

	if ps.hitRegexp(countHostByTopoNodeRegexp, http.MethodPost) {
		match := BizIDRegex.FindStringSubmatch(ps.RequestCtx.URI)
		bizID, err := util.GetInt64ByInterface(match[1])
		if err != nil {
			blog.Errorf("get business id from request path failed, name: %s, err: %v, rid: %s", countHostByTopoNodeRegexp, err, ps.RequestCtx.Rid)
			ps.err = fmt.Errorf("parse biz id from url failed, err: %s", err.Error())
			return ps
		}
		ps.Attribute.Resources = []meta.ResourceAttribute{
			meta.ResourceAttribute{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:   meta.HostInstance,
					Action: meta.FindMany,
				},
			},
		}

		return ps
	}

	// find hosts with condition operation.
	if ps.hitPattern(findHostsWithConditionPattern, http.MethodPost) {
		bizID, err := ps.parseBusinessID()
		if err != nil {
			ps.err = err
			return ps
		}
		ps.Attribute.Resources = []meta.ResourceAttribute{
			meta.ResourceAttribute{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:   meta.HostInstance,
					Action: meta.FindMany,
				},
			},
		}

		return ps
	}
	// find hosts without app id
	if ps.hitPattern(findBizHostsWithoutAppPattern, http.MethodPost) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			meta.ResourceAttribute{
				Basic: meta.Basic{
					Type:   meta.HostInstance,
					Action: meta.FindMany,
				},
			},
		}

		return ps
	}

	// find hosts under business specified by path parameter
	if ps.hitRegexp(findBizHostsRegex, http.MethodPost) {
		bizID, err := strconv.ParseInt(ps.RequestCtx.Elements[4], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("list business's hosts, but got invalid business id: %s", ps.RequestCtx.Elements[4])
			return ps
		}
		ps.Attribute.Resources = []meta.ResourceAttribute{
			meta.ResourceAttribute{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:   meta.HostInstance,
					Action: meta.FindMany,
				},
			},
		}

		return ps
	}

	// find hosts under business specified by path parameter with their topology information
	if ps.hitRegexp(findBizHostsTopoRegex, http.MethodPost) {
		bizID, err := strconv.ParseInt(ps.RequestCtx.Elements[4], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("list business's hosts with topo, but got invalid business id: %s", ps.RequestCtx.Elements[4])
			return ps
		}
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:   meta.HostInstance,
					Action: meta.FindMany,
				},
			},
		}

		return ps
	}

	if ps.hitPattern(findHostsDetailsPattern, http.MethodPost) {
		bizID, err := ps.parseBusinessID()
		if err != nil {
			ps.err = err
			return ps
		}
		ps.Attribute.Resources = []meta.ResourceAttribute{
			meta.ResourceAttribute{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:   meta.HostInstance,
					Action: meta.FindMany,
				},
			},
		}

		return ps
	}

	// update hosts batch. but can not get the exactly host id.
	if ps.hitPattern(updateHostInfoBatchPattern, http.MethodPut) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			meta.ResourceAttribute{
				Basic: meta.Basic{
					Type: meta.HostInstance,
					// Action: meta.UpdateMany,
					Action: meta.SkipAction,
				},
			},
		}

		return ps
	}

	// update hosts property batch. but can not get the exactly host id.
	if ps.hitPattern(updateHostPropertyBatchPattern, http.MethodPut) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type: meta.HostInstance,
					// Action: meta.UpdateMany,
					Action: meta.SkipAction,
				},
			},
		}

		return ps
	}

	if ps.hitPattern(hostInstallPattern, http.MethodPost) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:   meta.InstallBK,
					Action: meta.Update,
				},
			},
		}
		return ps
	}

	return ps
}

const (
	createHostFavoritePattern   = "/api/v3/hosts/favorites"
	findManyHostFavoritePattern = "/api/v3/hosts/favorites/search"
)

var (
	updateHostFavoriteRegexp   = regexp.MustCompile(`^/api/v3/hosts/favorites/[^\s/]+/?$`)
	deleteHostFavoriteRegexp   = regexp.MustCompile(`^/api/v3/hosts/favorites/[^\s/]+/?$`)
	increaseHostFavoriteRegexp = regexp.MustCompile(`^/api/v3/hosts/favorites/[^\s/]+/incr$`)
)

func (ps *parseStream) hostFavorite() *parseStream {
	if ps.shouldReturn() {
		return ps
	}

	// create host favorite operation.
	if ps.hitPattern(createHostFavoritePattern, http.MethodPost) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			meta.ResourceAttribute{
				Basic: meta.Basic{
					Type:   meta.HostFavorite,
					Action: meta.Create,
				},
			},
		}

		return ps
	}

	// update host favorite operation.
	if ps.hitRegexp(updateHostFavoriteRegexp, http.MethodPut) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			meta.ResourceAttribute{
				Basic: meta.Basic{
					Type:   meta.HostFavorite,
					Action: meta.Update,
					Name:   ps.RequestCtx.Elements[4],
				},
			},
		}

		return ps
	}

	// delete host favorite operation.
	if ps.hitRegexp(deleteHostFavoriteRegexp, http.MethodDelete) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			meta.ResourceAttribute{
				Basic: meta.Basic{
					Type:   meta.HostFavorite,
					Action: meta.DeleteMany,
					Name:   ps.RequestCtx.Elements[4],
				},
			},
		}

		return ps
	}

	// find many host favorite operation.
	if ps.hitPattern(findManyHostFavoritePattern, http.MethodPost) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			meta.ResourceAttribute{
				Basic: meta.Basic{
					Type:   meta.HostFavorite,
					Action: meta.FindMany,
				},
			},
		}

		return ps
	}

	// increase host favorite count by one.
	if ps.hitRegexp(increaseHostFavoriteRegexp, http.MethodPut) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			meta.ResourceAttribute{
				Basic: meta.Basic{
					Type:   meta.HostFavorite,
					Action: meta.Update,
					Name:   ps.RequestCtx.Elements[4],
				},
			},
		}

		return ps
	}
	return ps
}

var (
	searchSyncTask       = `/api/v3/hosts/cloud/search`
	confirmSyncTResource = `/api/v3/hosts/cloud/searchConfirm`
)

func (ps *parseStream) cloudResourceSync() *parseStream {
	if ps.shouldReturn() {
		return ps
	}

	if ps.hitPattern(searchSyncTask, http.MethodPost) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			meta.ResourceAttribute{
				Basic: meta.Basic{
					Type:   meta.ResourceSync,
					Action: meta.FindMany,
				},
			},
		}
		return ps
	}

	if ps.hitPattern(confirmSyncTResource, http.MethodPost) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			meta.ResourceAttribute{
				Basic: meta.Basic{
					Type:   meta.ResourceSync,
					Action: meta.FindMany,
				},
			},
		}
		return ps
	}

	return ps
}

var (
	findHostSnapshotAPIRegexp = regexp.MustCompile(`^/api/v3/hosts/snapshot/[0-9]+/?$`)
)

func (ps *parseStream) hostSnapshot() *parseStream {
	if ps.shouldReturn() {
		return ps
	}

	if ps.hitRegexp(findHostSnapshotAPIRegexp, http.MethodGet) {
		if len(ps.RequestCtx.Elements) != 5 {
			ps.err = errors.New("find host snapshot details query, but got invalid uri")
			return ps
		}

		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:   meta.HostInstance,
					Action: meta.SkipAction,
				},
			},
		}
		return ps
	}
	return ps
}

var (
	findIdentifierAPIRegexp = regexp.MustCompile(`^/api/v3/identifier/[^\s/]+/search/?$`)
)

func (ps *parseStream) findObjectIdentifier() *parseStream {
	if ps.shouldReturn() {
		return ps
	}

	if ps.hitRegexp(findIdentifierAPIRegexp, http.MethodPost) {
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
