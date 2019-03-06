package parser

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	"configcenter/src/auth/meta"
	"configcenter/src/framework/core/errors"
)

func (ps *parseStream) hostRelated() *parseStream {
	if ps.err != nil {
		return ps
	}

	ps.host().
		userCustom().
		hostFavorite()

	return ps
}

const (
	createUserCustomPattern = "/api/v3/userapi"
)

var (
	updateUserCustomRegexp      = regexp.MustCompile(`^/api/v3/userapi/[\S][^/]+/[0-9]+$`)
	deleteUserCustomRegexp      = regexp.MustCompile(`^/api/v3/userapi/[\S][^/]+/[0-9]+$`)
	findUserCustomRegexp        = regexp.MustCompile(`^/api/v3/userapi/search/[0-9]+$`)
	findUserCustomDetailsRegexp = regexp.MustCompile(`^/api/v3/userapi/detail/[0-9]+/[\S][^/]+$`)
	findWithUserCustomRegexp    = regexp.MustCompile(`^/api/v3/userapi/data/[0-9]+/[\S][^/]+/[0-9]+/[0-9]+$`)
)

func (ps *parseStream) userCustom() *parseStream {
	if ps.err != nil {
		return ps
	}

	// create user custom query operation.
	if ps.hitPattern(createUserCustomPattern, http.MethodPost) {
		type Business struct {
			BusinessID int64
		}
		biz := new(Business)
		if err := json.Unmarshal(ps.RequestCtx.Body, biz); err != nil {
			ps.err = fmt.Errorf("create host user custom query, but get business id failed, err: %v", err)
			return ps
		}

		ps.Attribute.Resources = []meta.ResourceAttribute{
			meta.ResourceAttribute{
				Basic: meta.Basic{
					Type:   meta.HostUserCustom,
					Action: meta.Create,
				},
			},
		}
		return ps
	}

	// update host user custom query operation.
	if ps.hitRegexp(updateUserCustomRegexp, http.MethodPut) {
		if len(ps.RequestCtx.Elements) != 5 {
			ps.err = errors.New("update host user custom query, but got invalid uri")
			return ps
		}
		bizID, err := strconv.ParseInt(ps.RequestCtx.Elements[3], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("update host user custom query failed, err: %v", err)
			return ps
		}
		ps.Attribute.BusinessID = bizID
		ps.Attribute.Resources = []meta.ResourceAttribute{
			meta.ResourceAttribute{
				Basic: meta.Basic{
					Type:   meta.HostUserCustom,
					Action: meta.Update,
					Name:   ps.RequestCtx.Elements[4],
				},
			},
		}
		return ps

	}

	// delete host user custom query operation.
	if ps.hitRegexp(deleteUserCustomRegexp, http.MethodDelete) {
		if len(ps.RequestCtx.Elements) != 5 {
			ps.err = errors.New("delete host user custom query operation, but got invalid uri")
			return ps
		}
		bizID, err := strconv.ParseInt(ps.RequestCtx.Elements[3], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("update host user custom query failed, err: %v", err)
			return ps
		}
		ps.Attribute.BusinessID = bizID
		ps.Attribute.Resources = []meta.ResourceAttribute{
			meta.ResourceAttribute{
				Basic: meta.Basic{
					Type:   meta.HostUserCustom,
					Action: meta.Delete,
					Name:   ps.RequestCtx.Elements[4],
				},
			},
		}
		return ps

	}

	// find host user custom query operation
	if ps.hitRegexp(findUserCustomRegexp, http.MethodPost) {
		if len(ps.RequestCtx.Elements) != 5 {
			ps.err = errors.New("find host usr custom query, but got invalid uri")
			return ps
		}

		bizID, err := strconv.ParseInt(ps.RequestCtx.Elements[4], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("find host user custom query failed, err: %v", err)
			return ps
		}
		ps.Attribute.BusinessID = bizID
		ps.Attribute.Resources = []meta.ResourceAttribute{
			meta.ResourceAttribute{
				Basic: meta.Basic{
					Type:   meta.HostUserCustom,
					Action: meta.FindMany,
				},
			},
		}
		return ps
	}

	// find host user custom query details operation.
	if ps.hitRegexp(findUserCustomDetailsRegexp, http.MethodGet) {
		if len(ps.RequestCtx.Elements) != 6 {
			ps.err = errors.New("find host user custom details query, but got invalid uri")
			return ps
		}

		bizID, err := strconv.ParseInt(ps.RequestCtx.Elements[4], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("find host user custom query details failed, err: %v", err)
			return ps
		}
		ps.Attribute.BusinessID = bizID
		ps.Attribute.Resources = []meta.ResourceAttribute{
			meta.ResourceAttribute{
				Basic: meta.Basic{
					Type:   meta.HostUserCustom,
					Action: meta.Find,
					Name:   ps.RequestCtx.Elements[5],
				},
			},
		}
		return ps
	}

	// get data with user custom query api.
	if ps.hitRegexp(findWithUserCustomRegexp, http.MethodGet) {
		if len(ps.RequestCtx.Elements) != 8 {
			ps.err = errors.New("find host user custom details query, but got invalid uri")
			return ps
		}

		bizID, err := strconv.ParseInt(ps.RequestCtx.Elements[4], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("find host user custom query details failed, err: %v", err)
			return ps
		}

		ps.Attribute.BusinessID = bizID
		ps.Attribute.Resources = []meta.ResourceAttribute{
			meta.ResourceAttribute{
				Basic: meta.Basic{
					Type:   meta.HostUserCustom,
					Action: meta.Find,
					Name:   ps.RequestCtx.Elements[5],
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
	moveHostsFromOneToAnotherBizModulePattern = "/api/v3/hosts/modules/biz/mutilple"
	cleanHostInSetOrModulePattern             = "/api/v3/hosts/modules/idle/set"
	// used in sync framework.
	moveHostToBusinessOrModulePattern = "/api/v3/hosts/sync/new/host"
	findHostsWithConditionPattern     = "/api/v3/hosts/search"
	updateHostInfoBatchPattern        = "/api/v3/hosts/batch"
)

func (ps *parseStream) host() *parseStream {
	if ps.err != nil {
		return ps
	}

	// TODO: add host lock authorize filter if needed.

	// delete hosts batch operation.
	if ps.hitPattern(deleteHostBatchPattern, http.MethodDelete) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			meta.ResourceAttribute{
				Basic: meta.Basic{
					Type:   meta.ModelInstance,
					Action: meta.DeleteMany,
					Name:   "host",
				},
			},
		}
		return ps
	}

	// TODO: add host clone authorize filter if needed.

	// add new hosts to resource pool
	if ps.hitPattern(addHostsToHostPoolPattern, http.MethodPost) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			meta.ResourceAttribute{
				Basic: meta.Basic{
					Type:   meta.ModelInstance,
					Action: meta.AddHostToResourcePool,
					Name:   meta.Host,
				},
			},
		}

		return ps
	}

	// move hosts from a module to resource pool.
	if ps.hitPattern(moveHostsFromModuleToResPoolPattern, http.MethodPost) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			meta.ResourceAttribute{
				Basic: meta.Basic{
					Type:   meta.ModelInstance,
					Action: meta.MoveHostFromModuleToResPool,
					Name:   meta.Host,
				},
			},
		}

		return ps
	}

	// move hosts to business module operation.
	if ps.hitPattern(moveHostToBusinessModulePattern, http.MethodPost) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			meta.ResourceAttribute{
				Basic: meta.Basic{
					Type:   meta.ModelInstance,
					Action: meta.MoveHostToModule,
					Name:   meta.Host,
				},
			},
		}

		return ps
	}

	// move resource pool hosts to a business idle module operation.
	if ps.hitPattern(moveResPoolToBizIdleModulePattern, http.MethodPost) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			meta.ResourceAttribute{
				Basic: meta.Basic{
					Type:   meta.ModelInstance,
					Action: meta.MoveResPoolHostToBizIdleModule,
					Name:   meta.Host,
				},
			},
		}

		return ps
	}

	// move host to a business fault module.
	if ps.hitPattern(moveHostsToBizFaultModulePattern, http.MethodPost) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			meta.ResourceAttribute{
				Basic: meta.Basic{
					Type:   meta.ModelInstance,
					Action: meta.MoveHostToBizFaultModule,
					Name:   meta.Host,
				},
			},
		}

		return ps
	}

	// move hosts to a business idle module.
	if ps.hitPattern(moveHostsToBizIdleModulePattern, http.MethodPost) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			meta.ResourceAttribute{
				Basic: meta.Basic{
					Type:   meta.ModelInstance,
					Action: meta.MoveHostToBizIdleModule,
					Name:   meta.Host,
				},
			},
		}

		return ps
	}

	// move hosts from one business module to another business module.
	if ps.hitPattern(moveHostsFromOneToAnotherBizModulePattern, http.MethodPost) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			meta.ResourceAttribute{
				Basic: meta.Basic{
					Type:   meta.ModelInstance,
					Action: meta.MoveHostToAnotherBizModule,
					Name:   meta.Host,
				},
			},
		}

		return ps
	}

	// clean the hosts in a set or module, and move these hosts to the business idle module
	// when these hosts only exist in this set or module. otherwise these hosts will only be
	// removed from this set or module.
	if ps.hitPattern(cleanHostInSetOrModulePattern, http.MethodPost) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			meta.ResourceAttribute{
				Basic: meta.Basic{
					Type:   meta.ModelInstance,
					Action: meta.CleanHostInSetOrModule,
					Name:   meta.Host,
				},
			},
		}

		return ps
	}

	// synchronize hosts directly to a module in a business if this host does not exist.
	// otherwise, this operation will only change host's attribute.
	if ps.hitPattern(moveHostToBusinessOrModulePattern, http.MethodPost) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			meta.ResourceAttribute{
				Basic: meta.Basic{
					Type:   meta.ModelInstance,
					Action: meta.MoveHostsToBusinessOrModule,
					Name:   meta.Host,
				},
			},
		}

		return ps
	}

	// find hosts with condition operation.
	if ps.hitPattern(findHostsWithConditionPattern, http.MethodPost) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			meta.ResourceAttribute{
				Basic: meta.Basic{
					Type:   meta.ModelInstance,
					Action: meta.FindMany,
					Name:   meta.Host,
				},
			},
		}

		return ps
	}

	// update hosts batch.
	if ps.hitPattern(updateHostInfoBatchPattern, http.MethodPut) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			meta.ResourceAttribute{
				Basic: meta.Basic{
					Type:   meta.ModelInstance,
					Action: meta.UpdateMany,
					Name:   meta.Host,
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
	updateHostFavoriteRegexp   = regexp.MustCompile(`^/api/v3/hosts/favorite/[\S][^/]+$`)
	deleteHostFavoriteRegexp   = regexp.MustCompile(`^/api/v3/hosts/favorite/[\S][^/]+$`)
	increaseHostFavoriteRegexp = regexp.MustCompile(`^/api/v3/hosts/favorite/[\S][^/]+/incr$`)
)

func (ps *parseStream) hostFavorite() *parseStream {
	if ps.err != nil {
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

	//

	return ps
}
