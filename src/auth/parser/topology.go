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
	"configcenter/src/auth"
	"configcenter/src/framework/core/errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
)

func (ps *parseStream) topology() *parseStream {
	if ps.err != nil {
		return ps
	}

	ps.business().
		mainline().
		associationType().
		objectAssociation().
		objectInstanceAssociation()

	return ps
}

var (
	createBusinessRegexp       = regexp.MustCompile(`^/api/v3/biz/[\S][^/]+$`)
	updateBusinessRegexp       = regexp.MustCompile(`^/api/v3/biz/[\S][^/]+/[0-9]+$`)
	deleteBusinessRegexp       = regexp.MustCompile(`^/api/v3/biz/[\S][^/]+/[0-9]+$`)
	findBusinessRegexp         = regexp.MustCompile(`^/api/v3/biz/search/[\S][^/]+$`)
	updateBusinessStatusRegexp = regexp.MustCompile(`^/api/v3/biz/status/[\S][^/]+/[\S][^/]+/[0-9]+$`)
)

func (ps *parseStream) business() *parseStream {
	if ps.err != nil {
		return ps
	}

	// create business, this is not a normalize api.
	// TODO: update this api format.
	if createBusinessRegexp.MatchString(ps.RequestCtx.URI) && ps.RequestCtx.Method == http.MethodPost {
		ps.Attribute.Resources = []auth.Resource{
			auth.Resource{
				Name:   auth.Business,
				Action: auth.Create,
			},
		}
		return ps
	}

	// update business, this is not a normalize api.
	// TODO: update this api format.
	if updateBusinessRegexp.MatchString(ps.RequestCtx.URI) && ps.RequestCtx.Method == http.MethodPut {
		if len(ps.RequestCtx.Elements) != 5 {
			ps.err = errors.New("invalid update business request uri")
			return ps
		}

		bizID, err := strconv.ParseInt(ps.RequestCtx.Elements[4], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("udpate business, but got invalid business id %s", ps.RequestCtx.Elements[4])
			return ps
		}

		ps.Attribute.Resources = []auth.Resource{
			auth.Resource{
				Name:       auth.Business,
				Action:     auth.Update,
				InstanceID: uint64(bizID),
				BusinessID: uint64(bizID),
			},
		}
		return ps
	}

	// update business enable status, this is not a normalize api.
	// TODO: update this api format.
	if updateBusinessRegexp.MatchString(ps.RequestCtx.URI) && ps.RequestCtx.Method == http.MethodPut {
		if len(ps.RequestCtx.Elements) != 7 {
			ps.err = errors.New("invalid update business enable status request uri")
			return ps
		}

		bizID, err := strconv.ParseInt(ps.RequestCtx.Elements[6], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("udpate business enable status, but got invalid business id %s", ps.RequestCtx.Elements[4])
			return ps
		}

		ps.Attribute.Resources = []auth.Resource{
			auth.Resource{
				Name:       auth.Business,
				Action:     auth.Update,
				InstanceID: uint64(bizID),
				BusinessID: uint64(bizID),
			},
		}
		return ps
	}

	// delete business, this is not a normalize api.
	// TODO: update this api format
	if updateBusinessRegexp.MatchString(ps.RequestCtx.URI) && ps.RequestCtx.Method == http.MethodDelete {
		if len(ps.RequestCtx.Elements) != 5 {
			ps.err = errors.New("invalid delete business request uri")
			return ps
		}

		bizID, err := strconv.ParseInt(ps.RequestCtx.Elements[4], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("delete business, but got invalid business id %s", ps.RequestCtx.Elements[4])
			return ps
		}

		ps.Attribute.Resources = []auth.Resource{
			auth.Resource{
				Name:       auth.Business,
				Action:     auth.Delete,
				InstanceID: uint64(bizID),
				BusinessID: uint64(bizID),
			},
		}
		return ps
	}

	// find business, this is not a normalize api.
	// TODO: update this api format
	if findBusinessRegexp.MatchString(ps.RequestCtx.URI) && ps.RequestCtx.Method == http.MethodPost {
		ps.Attribute.Resources = []auth.Resource{
			auth.Resource{
				Name: auth.Business,
				// we don't know if one or more business is to find, so we assume it's a find many
				// business operation.
				Action: auth.FindMany,
			},
		}
		return ps
	}

	return ps
}

const (
	createMainlineObjectPattern = "/api/v3/topo/model/mainline"
)

var (
	deleteMainlineObjectRegexp        = regexp.MustCompile(`^/api/v3/topo/model/mainline/owners/[\S][^/]+/objectids/[\S][^/]+$`)
	findMainlineObjectTopoRegexp      = regexp.MustCompile(`^/api/v3/topo/model/[\S][^/]+$`)
	findMainlineInstanceTopoRegexp    = regexp.MustCompile(`^/api/v3/topo/inst/[\S][^/]+/[0-9]+$`)
	findMainineSubInstanceTopoRegexp  = regexp.MustCompile(`^/api/v3/topo/inst/child/[\S][^/]+/[\S][^/]+/[0-9]+/[0-9]+$`)
	findMainlineIdleFaultModuleRegexp = regexp.MustCompile(`^/api/v3/topo/internal/[\S][^/]+/[0-9]+$`)
)

func (ps *parseStream) mainline() *parseStream {
	if ps.err != nil {
		return ps
	}

	// create mainline object operation.
	if ps.RequestCtx.URI == createMainlineObjectPattern && ps.RequestCtx.Method == http.MethodPost {
		ps.Attribute.Resources = []auth.Resource{
			auth.Resource{
				Name:   auth.MainlineObject,
				Action: auth.Create,
			},
		}
		return ps
	}

	// delete mainline object operation
	if deleteMainlineObjectRegexp.MatchString(ps.RequestCtx.URI) && ps.RequestCtx.Method == http.MethodDelete {
		ps.Attribute.Resources = []auth.Resource{
			auth.Resource{
				Name:   auth.MainlineObject,
				Action: auth.Delete,
			},
		}

		return ps
	}

	// find internal idle and fault machine module operation.

	// get mainline object operation
	if findMainlineObjectTopoRegexp.MatchString(ps.RequestCtx.URI) && ps.RequestCtx.Method == http.MethodGet {
		ps.Attribute.Resources = []auth.Resource{
			auth.Resource{
				Name:   auth.MainlineObjectTopology,
				Action: auth.Find,
			},
		}

		return ps
	}

	// find mainline instance topology operation.
	if findMainlineInstanceTopoRegexp.MatchString(ps.RequestCtx.URI) && ps.RequestCtx.Method == http.MethodGet {
		if len(ps.RequestCtx.Elements) != 6 {
			ps.err = errors.New("find mainline instance topology, but got invalid url")
			return ps
		}

		bizID, err := strconv.ParseInt(ps.RequestCtx.Elements[5], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("find mainline instance topology, but got invalid business id %s", ps.RequestCtx.Elements[5])
			return ps
		}

		ps.Attribute.Resources = []auth.Resource{
			auth.Resource{
				Name:       auth.MainlineInstanceTopology,
				Action:     auth.Find,
				BusinessID: uint64(bizID),
			},
		}

		return ps
	}

	// find mainline object instance's sub-instance topology operation.
	if findMainineSubInstanceTopoRegexp.MatchString(ps.RequestCtx.URI) && ps.RequestCtx.Method == http.MethodGet {
		if len(ps.RequestCtx.Elements) != 9 {
			ps.err = errors.New("find mainline object's sub instance topology, but got invalid url")
			return ps
		}

		bizID, err := strconv.ParseInt(ps.RequestCtx.Elements[7], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("find mainline object's sub instance topology, but got invalid business id %s", ps.RequestCtx.Elements[7])
			return ps
		}

		ps.Attribute.Resources = []auth.Resource{
			auth.Resource{
				Name:       auth.MainlineInstanceTopology,
				Action:     auth.Find,
				BusinessID: uint64(bizID),
			},
		}

		return ps
	}

	// find mainline internal idle and fault module operation.
	if findMainlineIdleFaultModuleRegexp.MatchString(ps.RequestCtx.URI) && ps.RequestCtx.Method == http.MethodGet {
		if len(ps.RequestCtx.Elements) != 6 {
			ps.err = errors.New("find mainline object's sub instance topology, but got invalid url")
			return ps
		}

		bizID, err := strconv.ParseInt(ps.RequestCtx.Elements[5], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("find mainline object's sub instance topology, but got invalid business id %s", ps.RequestCtx.Elements[5])
			return ps
		}

		ps.Attribute.Resources = []auth.Resource{
			auth.Resource{
				Name:       auth.MainlineObject,
				Action:     auth.Find,
				BusinessID: uint64(bizID),
			},
		}

		return ps
	}

	return ps
}

const (
	findManyAssociationKindPattern = "/api/v3/topo/association/type/action/search"
	createAssociationKindPattern   = "/api/v3/topo/association/type/action/search"
)

var (
	updateAssociationKindRegexp = regexp.MustCompile(`^/api/v3/topo/association/type/[0-9]+/action/update$`)
	deleteAssociationKindRegexp = regexp.MustCompile(`^/api/v3/topo/association/type/[0-9]+/action/delete$`)
)

func (ps *parseStream) associationType() *parseStream {
	if ps.err != nil {
		return ps
	}

	// find association kind operation
	if ps.RequestCtx.URI == findManyAssociationKindPattern && ps.RequestCtx.Method == http.MethodPost {
		ps.Attribute.Resources = []auth.Resource{
			auth.Resource{
				Name:   auth.AssociationKind,
				Action: auth.FindMany,
			},
		}
		return ps
	}

	// create association kind operation
	if ps.RequestCtx.URI == createAssociationKindPattern && ps.RequestCtx.Method == http.MethodPost {
		ps.Attribute.Resources = []auth.Resource{
			auth.Resource{
				Name:   auth.AssociationKind,
				Action: auth.Create,
			},
		}
		return ps
	}

	// update association kind operation
	if updateAssociationKindRegexp.MatchString(ps.RequestCtx.URI) && ps.RequestCtx.Method == http.MethodPut {
		if len(ps.RequestCtx.Elements) != 8 {
			ps.err = errors.New("update association kind, but got invalid url")
			return ps
		}

		kindID, err := strconv.ParseInt(ps.RequestCtx.Elements[5], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("update association kind, but got invalid kind id %s", ps.RequestCtx.Elements[5])
			return ps
		}
		ps.Attribute.Resources = []auth.Resource{
			auth.Resource{
				Name:       auth.AssociationKind,
				Action:     auth.Update,
				InstanceID: uint64(kindID),
			},
		}

		return ps
	}

	// delete association kind operation
	if deleteAssociationKindRegexp.MatchString(ps.RequestCtx.URI) && ps.RequestCtx.Method == http.MethodDelete {
		if len(ps.RequestCtx.Elements) != 8 {
			ps.err = errors.New("delete association kind, but got invalid url")
			return ps
		}

		kindID, err := strconv.ParseInt(ps.RequestCtx.Elements[5], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("delete association kind, but got invalid kind id %s", ps.RequestCtx.Elements[5])
			return ps
		}
		ps.Attribute.Resources = []auth.Resource{
			auth.Resource{
				Name:       auth.AssociationKind,
				Action:     auth.Delete,
				InstanceID: uint64(kindID),
			},
		}

		return ps
	}

	return ps
}

const (
	findObjectAssociationPattern                    = "/api/v3/object/association/action/search"
	createObjectAssociationPattern                  = "/api/v3/object/association/action/create"
	findObjectAssociationWithAssociationKindPattern = "/api/v3/topo/association/type/action/search/batch"
)

var (
	updateObjectAssociationRegexp = regexp.MustCompile(`^/api/v3/object/association/[0-9]+/action/update$`)
	deleteObjectAssociationRegexp = regexp.MustCompile(`^/api/v3/object/association/[0-9]+/action/delete$`)
)

func (ps *parseStream) objectAssociation() *parseStream {
	if ps.err != nil {
		return ps
	}

	// search object association operation
	if ps.RequestCtx.URI == findObjectAssociationPattern && ps.RequestCtx.Method == http.MethodPost {
		ps.Attribute.Resources = []auth.Resource{
			auth.Resource{
				Name:   auth.ObjectAssociation,
				Action: auth.FindMany,
			},
		}
		return ps
	}

	// create object association operation
	if ps.RequestCtx.URI == createObjectAssociationPattern && ps.RequestCtx.Method == http.MethodPost {
		ps.Attribute.Resources = []auth.Resource{
			auth.Resource{
				Name:   auth.ObjectAssociation,
				Action: auth.Create,
			},
		}
		return ps
	}

	// update object association operation
	if updateObjectAssociationRegexp.MatchString(ps.RequestCtx.URI) && ps.RequestCtx.Method == http.MethodPut {
		if len(ps.RequestCtx.Elements) != 7 {
			ps.err = errors.New("update object association, but got invalid url")
			return ps
		}

		assoID, err := strconv.ParseInt(ps.RequestCtx.Elements[4], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("update object association, but got invalid association id %s", ps.RequestCtx.Elements[4])
			return ps
		}

		ps.Attribute.Resources = []auth.Resource{
			auth.Resource{
				Name:       auth.ObjectAssociation,
				Action:     auth.Update,
				InstanceID: uint64(assoID),
			},
		}
		return ps
	}

	// delete object association operation
	if deleteObjectAssociationRegexp.MatchString(ps.RequestCtx.URI) && ps.RequestCtx.Method == http.MethodDelete {
		if len(ps.RequestCtx.Elements) != 7 {
			ps.err = errors.New("delete object association, but got invalid url")
			return ps
		}

		assoID, err := strconv.ParseInt(ps.RequestCtx.Elements[4], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("delete object association, but got invalid association id %s", ps.RequestCtx.Elements[4])
			return ps
		}

		ps.Attribute.Resources = []auth.Resource{
			auth.Resource{
				Name:       auth.ObjectAssociation,
				Action:     auth.Delete,
				InstanceID: uint64(assoID),
			},
		}
		return ps
	}

	// find object association with a association kind list.
	if ps.RequestCtx.URI == findObjectAssociationWithAssociationKindPattern && ps.RequestCtx.Method == http.MethodPost {
		ps.Attribute.Resources = []auth.Resource{
			auth.Resource{
				Name:   auth.ObjectAssociation,
				Action: auth.FindMany,
			},
		}
		return ps
	}

	return ps
}

const (
	findObjectInstanceAssociationPattern   = "/api/v3/inst/association/action/search"
	createObjectInstanceAssociationPattern = "/api/v3/inst/association/action/create"
)

var (
	deleteObjectInstanceAssociationRegexp = regexp.MustCompile("/api/v3/inst/association/[0-9]+/action/create")
)

func (ps *parseStream) objectInstanceAssociation() *parseStream {
	if ps.err != nil {
		return ps
	}

	// find object instance's association operation.
	if ps.RequestCtx.URI == findObjectInstanceAssociationPattern && ps.RequestCtx.Method == http.MethodPost {
		ps.Attribute.Resources = []auth.Resource{
			auth.Resource{
				Name:   auth.ObjectInstanceAssociation,
				Action: auth.FindMany,
			},
		}
		return ps
	}

	// create object's instance association operation.
	if ps.RequestCtx.URI == createObjectInstanceAssociationPattern && ps.RequestCtx.Method == http.MethodPost {
		ps.Attribute.Resources = []auth.Resource{
			auth.Resource{
				Name:   auth.ObjectInstanceAssociation,
				Action: auth.Create,
			},
		}
		return ps
	}

	// delete object's instance association operation.
	if ps.hitRegexp(deleteObjectInstanceAssociationRegexp, http.MethodDelete) {
		if len(ps.RequestCtx.Elements) != 7 {
			ps.err = errors.New("delete object instance association, but got invalid url")
			return ps
		}

		assoID, err := strconv.ParseInt(ps.RequestCtx.Elements[4], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("delete object instance association, but got invalid association id %s", ps.RequestCtx.Elements[4])
			return ps
		}

		ps.Attribute.Resources = []auth.Resource{
			auth.Resource{
				Name:       auth.ObjectInstanceAssociation,
				Action:     auth.Delete,
				InstanceID: uint64(assoID),
			},
		}
		return ps
	}

	return ps
}

var (
	createObjectInstanceRegexp = regexp.MustCompile(`^/api/v3/inst/[\S][^/]+/[\S][^/]+`)
	findObjectInstanceRegexp   = regexp.MustCompile(`^/api/v3/inst/association/search/owner/[\S][^/]+/object/[\S][^/]+$`)
)

func (ps *parseStream) objectInstance() *parseStream {
	if ps.err != nil {
		return ps
	}

	// create object instance operation.
	if ps.hitRegexp(createObjectInstanceRegexp, http.MethodPost) {
		ps.Attribute.Resources = []auth.Resource{
			auth.Resource{
				Name:   auth.ObjectInstance,
				Action: auth.Create,
			},
		}
		return ps
	}

	// find object instance operation.
	if ps.hitRegexp(findObjectInstanceRegexp, http.MethodPost) {
	    
    }

	return ps
}
