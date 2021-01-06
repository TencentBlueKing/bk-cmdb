package parser

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	"configcenter/src/ac/meta"
	"configcenter/src/common"
	"configcenter/src/common/json"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
)

func (ps *parseStream) hostRelated() *parseStream {
	if ps.shouldReturn() {
		return ps
	}

	ps.host().
		hostTransfer().
		dynamicGrouping().
		userCustom().
		hostFavorite().
		hostSnapshot().
		findObjectIdentifier().
		HostApply()
	return ps
}

func (ps *parseStream) parseBusinessID() (int64, error) {
	val, err := ps.RequestCtx.getValueFromBody(common.BKAppIDField)
	if err != nil {
		return 0, err
	}
	if !val.Exists() {
		return 0, nil
	}
	bizID := val.Int()
	if bizID == 0 {
		return 0, errors.New("invalid bk_biz_id value")
	}
	return bizID, nil
}

var (
	createDynamicGroupPattern = "/api/v3/dynamicgroup"
	updateDynamicGroupRegexp  = regexp.MustCompile(`^/api/v3/dynamicgroup/[0-9]+/[^\s/]+/?$`)
	deleteDynamicGroupRegexp  = regexp.MustCompile(`^/api/v3/dynamicgroup/[0-9]+/[^\s/]+/?$`)
	getDynamicGroupRegexp     = regexp.MustCompile(`^/api/v3/dynamicgroup/[0-9]+/[^\s/]+/?$`)
	searchDynamicGroupRegexp  = regexp.MustCompile(`^/api/v3/dynamicgroup/search/[0-9]+/?$`)
	executeDynamicGroupRegexp = regexp.MustCompile(`^/api/v3/dynamicgroup/data/[0-9]+/[^\s/]+/?$`)
)

func (ps *parseStream) dynamicGrouping() *parseStream {
	if ps.shouldReturn() {
		return ps
	}

	if ps.hitPattern(createDynamicGroupPattern, http.MethodPost) {
		bizID, err := ps.parseBusinessID()
		if err != nil {
			ps.err = err
			return ps
		}
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:   meta.DynamicGrouping,
					Action: meta.Create,
				},
			},
		}
		return ps
	}

	if ps.hitRegexp(updateDynamicGroupRegexp, http.MethodPut) {
		if len(ps.RequestCtx.Elements) != 5 {
			ps.err = errors.New("update dynamic group, but got invalid uri")
			return ps
		}
		bizID, err := strconv.ParseInt(ps.RequestCtx.Elements[3], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("update dynamic group failed, err: %v", err)
			return ps
		}
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
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

	if ps.hitRegexp(deleteDynamicGroupRegexp, http.MethodDelete) {
		if len(ps.RequestCtx.Elements) != 5 {
			ps.err = errors.New("delete dynamic group, but got invalid uri")
			return ps
		}
		bizID, err := strconv.ParseInt(ps.RequestCtx.Elements[3], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("update dynamic group failed, err: %v", err)
			return ps
		}
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
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

	if ps.hitRegexp(searchDynamicGroupRegexp, http.MethodPost) {
		if len(ps.RequestCtx.Elements) != 5 {
			ps.err = errors.New("search dynamic groups, but got invalid uri")
			return ps
		}

		bizID, err := strconv.ParseInt(ps.RequestCtx.Elements[4], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("search dynamic groups failed, err: %v", err)
			return ps
		}

		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:   meta.DynamicGrouping,
					Action: meta.FindMany,
				},
			},
		}
		return ps
	}

	if ps.hitRegexp(getDynamicGroupRegexp, http.MethodGet) {
		if len(ps.RequestCtx.Elements) != 5 {
			ps.err = errors.New("find dynamic group detail, but got invalid uri")
			return ps
		}

		bizID, err := strconv.ParseInt(ps.RequestCtx.Elements[3], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("find dynamic group failed, err: %v", err)
			return ps
		}
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:         meta.DynamicGrouping,
					Action:       meta.Find,
					InstanceIDEx: ps.RequestCtx.Elements[4],
				},
			},
		}
		return ps
	}

	if ps.hitRegexp(executeDynamicGroupRegexp, http.MethodPost) {
		if len(ps.RequestCtx.Elements) != 6 {
			ps.err = errors.New("execute dynamic group, but got invalid uri")
			return ps
		}

		bizID, err := strconv.ParseInt(ps.RequestCtx.Elements[4], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("execute dynamic group failed, err: %v", err)
			return ps
		}

		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
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
	saveUserCustomPattern         = `/api/v3/usercustom`
	searchUserCustomPattern       = `/api/v3/usercustom/user/search`
	getModelDefaultCustomPattern  = `/api/v3/usercustom/default/model`
	saveModelDefaultCustomPattern = regexp.MustCompile(`^/api/v3/usercustom/default/model/[^\s/]+/?$`)
)

func (ps *parseStream) userCustom() *parseStream {
	if ps.shouldReturn() {
		return ps
	}

	// create user custom query operation.
	if ps.hitPattern(saveUserCustomPattern, http.MethodPost) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
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
			{
				Basic: meta.Basic{
					Type:   meta.UserCustom,
					Action: meta.Find,
				},
			},
		}
		return ps

	}

	// get default model list header
	if ps.hitPattern(getModelDefaultCustomPattern, http.MethodPost) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:   meta.UserCustom,
					Action: meta.Find,
				},
			},
		}
		return ps

	}

	// set default  model list header
	if ps.hitRegexp(saveModelDefaultCustomPattern, http.MethodPost) {
		if len(ps.RequestCtx.Elements) != 6 {
			ps.err = errors.New("search object instance, but got invalid url")
			return ps
		}
		model, err := ps.getOneModel(mapstr.MapStr{common.BKObjIDField: ps.RequestCtx.Elements[5]})
		if err != nil {
			ps.err = err
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
					Type:   meta.UserCustom,
					Action: meta.Create,
					Name:   ps.RequestCtx.Elements[5],
				},
			},
			{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:   meta.ModelAttribute,
					Action: meta.Create,
				},
				Layers: []meta.Item{{Type: meta.Model, InstanceID: model.ID}},
			},
		}
		return ps

	}
	return ps
}

const (
	deleteHostBatchPattern                = "/api/v3/hosts/batch"
	addHostsToHostPoolPattern             = "/api/v3/hosts/add"
	addHostsByExcelPattern                = "/api/v3/hosts/excel/add"
	addHostsToResourcePoolPattern         = "/api/v3/hosts/add/resource"
	moveHostToBusinessModulePattern       = "/api/v3/hosts/modules"
	moveResPoolHostToBizIdleModulePattern = "/api/v3/hosts/modules/resource/idle"
	moveHostsToBizFaultModulePattern      = "/api/v3/hosts/modules/fault"
	moveHostsFromModuleToResPoolPattern   = "/api/v3/hosts/modules/resource"
	moveHostsToBizIdleModulePattern       = "/api/v3/hosts/modules/idle"
	moveHostsToBizRecycleModulePattern    = "/api/v3/hosts/modules/recycle"
	moveHostAcrossBizPattern              = "/api/v3/hosts/modules/across/biz"
	moveRscPoolHostToRscPoolDir           = "/api/v3/host/transfer/resource/directory"
	cleanHostInSetOrModulePattern         = "/api/v3/hosts/modules/idle/set"
	findHostTopoRelationPattern           = "/api/v3/host/topo/relation/read"
	updateHostCloudAreaFieldPattern       = "/api/v3/updatemany/hosts/cloudarea_field"
	updateImportHostsPattern              = "/api/v3/hosts/update"
	getHostModuleRelationPattern          = "/api/v3/hosts/modules/read"
	lockHostPattern                       = "/api/v3/host/lock"
	unLockHostPattern                     = "/api/v3/host/lock"
	queryHostLockPattern                  = "/api/v3/host/lock/search"

	// used in sync framework.
	// moveHostToBusinessOrModulePattern = "/api/v3/hosts/sync/new/host"
	findHostsWithConditionPattern  = "/api/v3/hosts/search"
	findBizHostsWithoutAppPattern  = "/api/v3/hosts/list_hosts_without_app"
	findResourcePoolHostsPattern   = "/api/v3/hosts/list_resource_pool_hosts"
	findHostsDetailsPattern        = "/api/v3/hosts/search/asstdetail"
	updateHostInfoBatchPattern     = "/api/v3/hosts/batch"
	updateHostPropertyBatchPattern = "/api/v3/hosts/property/batch"
	cloneHostPropertyBatchPattern  = "/api/v3/hosts/property/clone"

	// 特殊接口，给蓝鲸业务使用
	hostInstallPattern = "/api/v3/host/install/bk"

	// cc system user config
	systemUserConfig = "/api/v3/system/config/user_config/blueking_modify"
)

var (
	findBizHostsRegex     = regexp.MustCompile(`/api/v3/hosts/app/\d+/list_hosts`)
	findBizHostsTopoRegex = regexp.MustCompile(`/api/v3/hosts/app/\d+/list_hosts_topo`)
	// find host instance's object properties info
	findHostInstanceObjectPropertiesRegexp = regexp.MustCompile(`^/api/v3/hosts/[^\s/]+/[0-9]+/?$`)

	transferHostWithAutoClearServiceInstanceRegex        = regexp.MustCompile("^/api/v3/host/transfer_with_auto_clear_service_instance/bk_biz_id/[0-9]+/?$")
	transferHostWithAutoClearServiceInstancePreviewRegex = regexp.MustCompile("^/api/v3/host/transfer_with_auto_clear_service_instance/bk_biz_id/[0-9]+/preview/?$")

	countHostByTopoNodeRegexp = regexp.MustCompile(`^/api/v3/host/count_by_topo_node/bk_biz_id/[0-9]+$`)

	findHostsByServiceTemplatesRegex = regexp.MustCompile(`^/api/v3/findmany/hosts/by_service_templates/biz/\d+$`)
	findHostsBySetTemplatesRegex     = regexp.MustCompile(`^/api/v3/findmany/hosts/by_set_templates/biz/\d+$`)
	findHostModuleRelationsRegex     = regexp.MustCompile(`^/api/v3/findmany/module_relation/bk_biz_id/[0-9]+/?$`)
	findHostsByTopoRegex             = regexp.MustCompile(`^/api/v3/findmany/hosts/by_topo/biz/\d+$`)
)

func (ps *parseStream) host() *parseStream {
	if ps.shouldReturn() {
		return ps
	}

	if ps.hitRegexp(findHostInstanceObjectPropertiesRegexp, http.MethodGet) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:   meta.HostInstance,
					Action: meta.Find,
				},
			},
		}
		return ps
	}

	if ps.hitRegexp(findHostsByServiceTemplatesRegex, http.MethodPost) {
		bizID, err := strconv.ParseInt(ps.RequestCtx.Elements[6], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("find hosts by service templates, but got invalid business id: %s", ps.RequestCtx.Elements[6])
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

	if ps.hitRegexp(findHostModuleRelationsRegex, http.MethodPost) {
		bizID, err := strconv.ParseInt(ps.RequestCtx.Elements[5], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("find host module relations, but got invalid business id: %s", ps.RequestCtx.Elements[5])
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

	if ps.hitRegexp(findHostsBySetTemplatesRegex, http.MethodPost) {
		bizID, err := strconv.ParseInt(ps.RequestCtx.Elements[6], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("find hosts by set templates, but got invalid business id: %s", ps.RequestCtx.Elements[6])
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

	if ps.hitRegexp(findHostsByTopoRegex, http.MethodPost) {
		bizID, err := strconv.ParseInt(ps.RequestCtx.Elements[6], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("find hosts by set templates, but got invalid business id: %s", ps.RequestCtx.Elements[6])
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

	// host lock authorize filter
	if ps.hitPattern(lockHostPattern, http.MethodPost) {
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

	if ps.hitPattern(unLockHostPattern, http.MethodDelete) {
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

	if ps.hitPattern(queryHostLockPattern, http.MethodPost) {
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

	// delete hosts batch operation.
	if ps.hitPattern(deleteHostBatchPattern, http.MethodDelete) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
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
			{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:   meta.ProcessServiceInstance,
					Action: meta.Update,
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

	if ps.hitRegexp(countHostByTopoNodeRegexp, http.MethodPost) {
		bizID, err := strconv.ParseInt(ps.RequestCtx.Elements[5], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("count host by topo node, but got invalid business id: %s", ps.RequestCtx.Elements[5])
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

	// find hosts with condition operation.
	if ps.hitPattern(findHostsWithConditionPattern, http.MethodPost) {
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
					Action: meta.FindMany,
				},
			},
		}

		return ps
	}
	// find hosts without app id
	if ps.hitPattern(findBizHostsWithoutAppPattern, http.MethodPost) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:   meta.HostInstance,
					Action: meta.FindMany,
				},
			},
		}

		return ps
	}

	// find resource pool hosts
	if ps.hitPattern(findResourcePoolHostsPattern, http.MethodPost) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
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

	// update hosts batch. but can not get the exactly host id.
	if ps.hitPattern(updateHostInfoBatchPattern, http.MethodPut) {
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

	// update import hosts batch. but can not get the exactly host id.
	if ps.hitPattern(updateImportHostsPattern, http.MethodPut) {
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

	// clone hosts property, but can not get the exactly host id.
	if ps.hitPattern(cloneHostPropertyBatchPattern, http.MethodPut) {
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

	if ps.hitPattern(systemUserConfig, http.MethodPost) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:   meta.SystemConfig,
					Action: meta.FindMany,
				},
			},
		}
		return ps
	}

	if ps.hitPattern(getHostModuleRelationPattern, http.MethodPost) {
		bizID, err := ps.parseBusinessID()
		if err != nil {
			ps.err = err
			return ps
		}
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				BusinessID: bizID,
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

func (ps *parseStream) hostTransfer() *parseStream {
	if ps.shouldReturn() {
		return ps
	}

	// add new hosts to resource pool
	if ps.hitPattern(addHostsToHostPoolPattern, http.MethodPost) {
		dirID, err := ps.getResourcePoolDefaultDirID()
		if err != nil {
			ps.err = fmt.Errorf("invalid directory id value, %s", err.Error())
			return ps
		}
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:   meta.HostInstance,
					Action: meta.AddHostToResourcePool,
				},
				Layers: []meta.Item{
					{
						Type:       meta.ResourcePoolDirectory,
						InstanceID: dirID,
					},
				},
			},
		}

		return ps
	}

	// add new hosts come from excel to resource pool directory
	if ps.hitPattern(addHostsByExcelPattern, http.MethodPost) {
		val, err := ps.RequestCtx.getValueFromBody("bk_module_id")
		if err != nil {
			ps.err = err
			return ps
		}
		dirID := val.Int()
		if dirID == 0 {
			var err error
			dirID, err = ps.getResourcePoolDefaultDirID()
			if err != nil {
				ps.err = fmt.Errorf("invalid directory id value, %s", err.Error())
				return ps
			}
		}
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:   meta.HostInstance,
					Action: meta.AddHostToResourcePool,
				},
				Layers: []meta.Item{
					{
						Type:       meta.ResourcePoolDirectory,
						InstanceID: dirID,
					},
				},
			},
		}

		return ps
	}

	// add hosts to resource pool directory
	if ps.hitPattern(addHostsToResourcePoolPattern, http.MethodPost) {
		val, err := ps.RequestCtx.getValueFromBody("directory")
		if err != nil {
			ps.err = err
			return ps
		}
		dirID := val.Int()
		if dirID == 0 {
			var err error
			dirID, err = ps.getResourcePoolDefaultDirID()
			if err != nil {
				ps.err = fmt.Errorf("invalid directory id value, %s", err.Error())
				return ps
			}
		}
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:   meta.HostInstance,
					Action: meta.AddHostToResourcePool,
				},
				Layers: []meta.Item{
					{
						Type:       meta.ResourcePoolDirectory,
						InstanceID: dirID,
					},
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

		rawDirID, err := ps.RequestCtx.getValueFromBody(common.BKModuleIDField)
		if err != nil {
			ps.err = err
			return ps
		}
		dirID := rawDirID.Int()

		// if directory id is not specified, transfer host to the default directory, use it to authorize
		if dirID == 0 {
			dirID, err = ps.getResourcePoolDefaultDirID()
			if err != nil {
				ps.err = fmt.Errorf("invalid directory id value, %s", err.Error())
				return ps
			}
		}

		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:   meta.HostInstance,
					Action: meta.MoveBizHostFromModuleToResPool,
				},
				Layers: []meta.Item{
					{
						Type:       meta.Business,
						InstanceID: bizID,
					},
					{
						Type:       meta.ResourcePoolDirectory,
						InstanceID: dirID,
					},
				},
			},
		}

		return ps
	}

	// move hosts to business module operation, transfer host in the same business.
	if ps.hitPattern(moveHostToBusinessModulePattern, http.MethodPost) {
		bizID, err := ps.parseBusinessID()
		if err != nil {
			ps.err = err
			return ps
		}
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:   meta.ProcessServiceInstance,
					Action: meta.Update,
				},
			},
		}

		return ps
	}

	// move resource pool hosts to a business idle module operation.
	if ps.hitPattern(moveResPoolHostToBizIdleModulePattern, http.MethodPost) {
		opt := new(hostPool)
		body, err := ps.RequestCtx.getRequestBody()
		if err != nil {
			ps.err = err
			return ps
		}
		if err := json.Unmarshal(body, opt); err != nil {
			ps.err = err
			return ps
		}

		relation, err := ps.getRscPoolHostModuleRelation(opt.HostID)
		if err != nil {
			ps.err = err
			return ps
		}

		for _, id := range opt.HostID {
			srcModuleID, exist := relation[id]
			if !exist {
				ps.err = errors.New("host not exist in resource pool")
				return ps
			}

			ps.Attribute.Resources = []meta.ResourceAttribute{
				{
					BusinessID: opt.Business,
					Basic: meta.Basic{
						Type:       meta.HostInstance,
						Action:     meta.MoveResPoolHostToBizIdleModule,
						InstanceID: id,
					},
					Layers: []meta.Item{{Type: meta.ModelModule, InstanceID: srcModuleID},
						{Type: meta.Business, InstanceID: opt.Business}},
				},
			}
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
			{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:   meta.ProcessServiceInstance,
					Action: meta.Update,
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
			{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:   meta.ProcessServiceInstance,
					Action: meta.Update,
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
			{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:   meta.ProcessServiceInstance,
					Action: meta.Update,
				},
			},
		}

		return ps
	}

	// transfer host to another business
	if ps.hitPattern(moveHostAcrossBizPattern, http.MethodPost) {
		val, err := ps.RequestCtx.getValueFromBody("src_bk_biz_id")
		if err != nil {
			ps.err = err
			return ps
		}
		srcBizID := val.Int()
		if srcBizID == 0 {
			ps.err = errors.New("src_bk_biz_id invalid")
			return ps
		}
		val, err = ps.RequestCtx.getValueFromBody("dst_bk_biz_id")
		if err != nil {
			ps.err = err
			return ps
		}
		dstBizID := val.Int()
		if dstBizID == 0 {
			ps.err = errors.New("dst_bk_biz_id invalid")
			return ps
		}
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				BusinessID: srcBizID,
				Basic: meta.Basic{
					Type:   meta.HostInstance,
					Action: meta.MoveHostToAnotherBizModule,
				},
				Layers: []meta.Item{
					{
						Type:       meta.Business,
						InstanceID: srcBizID,
					},
					{
						Type:       meta.Business,
						InstanceID: dstBizID,
					},
				},
			},
		}
		return ps
	}

	// synchronize hosts directly to a module in a business if this host does not exist.
	// otherwise, this operation will only change host's attribute.
	//if ps.hitPattern(moveHostToBusinessOrModulePattern, http.MethodPost) {
	//	bizID, err := ps.parseBusinessID()
	//	if err != nil {
	//		ps.err = err
	//		return ps
	//	}
	//	ps.Attribute.Resources = []meta.ResourceAttribute{
	//		{
	//			BusinessID: bizID,
	//			Basic: meta.Basic{
	//				Type:   meta.HostInstance,
	//				Action: meta.MoveHostsToBusinessOrModule,
	//			},
	//		},
	//	}
	//
	//	return ps
	//}

	if ps.hitRegexp(transferHostWithAutoClearServiceInstanceRegex, http.MethodPost) ||
		ps.hitRegexp(transferHostWithAutoClearServiceInstancePreviewRegex, http.MethodPost) {

		bizID, err := strconv.ParseInt(ps.RequestCtx.Elements[5], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("transfer host with auto clear service instance, but got invalid business id: %s", ps.RequestCtx.Elements[5])
			return ps
		}

		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:   meta.ProcessServiceInstance,
					Action: meta.Create,
				},
			},
			{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:   meta.ProcessServiceInstance,
					Action: meta.Update,
				},
			},
			{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:   meta.ProcessServiceInstance,
					Action: meta.Delete,
				},
			},
		}

		return ps
	}

	// move the resource pool host to another dir in resource pool
	if ps.hitPattern(moveRscPoolHostToRscPoolDir, http.MethodPost) {
		opt := new(metadata.TransferHostResourceDirectory)
		body, err := ps.RequestCtx.getRequestBody()
		if err != nil {
			ps.err = err
			return ps
		}
		if err := json.Unmarshal(body, opt); err != nil {
			ps.err = err
			return ps
		}

		relation, err := ps.getRscPoolHostModuleRelation(opt.HostID)
		if err != nil {
			ps.err = err
			return ps
		}

		for _, id := range opt.HostID {
			srcModuleID, exist := relation[id]
			if !exist {
				ps.err = errors.New("host not exist in resource pool")
				return ps
			}

			ps.Attribute.Resources = []meta.ResourceAttribute{
				{
					Basic: meta.Basic{
						Type:       meta.HostInstance,
						Action:     meta.MoveResPoolHostToDirectory,
						InstanceID: id,
					},
					Layers: []meta.Item{{Type: meta.ModelModule, InstanceID: srcModuleID},
						{Type: meta.ModelModule, InstanceID: opt.ModuleID}},
				},
			}
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
			{
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
			{
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
			{
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
			{
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
			{
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
	findHostSnapshotAPIRegexp = regexp.MustCompile(`^/api/v3/hosts/snapshot/[0-9]+/?$`)

	findHostSnapshotBatchPattern = "/api/v3/hosts/snapshot/batch"
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

	if ps.hitPattern(findHostSnapshotBatchPattern, http.MethodPost) {
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
