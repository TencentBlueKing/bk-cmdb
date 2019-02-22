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

	"configcenter/src/auth"
)

// this package's topology filter is the latest api version
// for these resources, it also has a elder resource api version.
// TODO: if the elder api has been removed, delete their resource
// filter at the same time.

var (
	createObjectUniqueLatestRegexp = regexp.MustCompile(`^/api/v3/create/objectunique/object/[\S][^/]+$`)
	updateObjectUniqueLatestRegexp = regexp.MustCompile(`^/api/v3/update/objectunique/object/[\S][^/]+/unique/[0-9]+$`)
	deleteObjectUniqueLatestRegexp = regexp.MustCompile(`^/api/v3/delete/objectunique/object/[\S][^/]+/unique/[0-9]+$`)
	findObjectUniqueLatestRegexp   = regexp.MustCompile(`^/api/v3/find/objectunique/object/[\S][^/]+`)
)

func (ps *parseStream) objectUniqueLatest() *parseStream {
	if ps.err != nil {
		return ps
	}

	// TODO: add business id for these filter rules to resources.
	// add object unique operation.
	if ps.hitRegexp(createObjectUniqueLatestRegexp, http.MethodPost) {
		ps.Attribute.Resources = []auth.Resource{
			auth.Resource{
				Type:   auth.ObjectUnique,
				Name:   ps.RequestCtx.Elements[5],
				Action: auth.Create,
			},
		}
		return ps
	}

	// update object unique operation.
	if ps.hitRegexp(updateObjectUniqueLatestRegexp, http.MethodPut) {
		uniqueID, err := strconv.ParseInt(ps.RequestCtx.Elements[7], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("update object unique, but got invalid unique id %s", ps.RequestCtx.Elements[7])
			return ps
		}

		ps.Attribute.Resources = []auth.Resource{
			auth.Resource{
				Type:       auth.ObjectUnique,
				InstanceID: uniqueID,
				Action:     auth.Update,
				Affiliated: auth.Affiliated{
					Type: auth.Object,
					Name: ps.RequestCtx.Elements[5],
				},
			},
		}
		return ps
	}

	// delete object unique operation.
	if ps.hitRegexp(deleteObjectUniqueLatestRegexp, http.MethodDelete) {
		uniqueID, err := strconv.ParseInt(ps.RequestCtx.Elements[7], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("update object unique, but got invalid unique id %s", ps.RequestCtx.Elements[7])
			return ps
		}

		ps.Attribute.Resources = []auth.Resource{
			auth.Resource{
				Type:       auth.ObjectUnique,
				InstanceID: uniqueID,
				Action:     auth.Delete,
				Affiliated: auth.Affiliated{
					Type: auth.Object,
					Name: ps.RequestCtx.Elements[5],
				},
			},
		}
		return ps
	}

	// find object unique operation.
	if ps.hitRegexp(findObjectUniqueLatestRegexp, http.MethodGet) {
		ps.Attribute.Resources = []auth.Resource{
			auth.Resource{
				Type:   auth.ObjectUnique,
				Action: auth.FindMany,
				Affiliated: auth.Affiliated{
					Type: auth.Object,
					Name: ps.RequestCtx.Elements[5],
				},
			},
		}
		return ps
	}

	return ps
}

const (
	findManyAssociationKindLatestPattern = "/api/v3/find/associationtype"
	createAssociationKindLatestPattern   = "/api/v3/create/associationtype"
)

var (
	updateAssociationKindLatestRegexp = regexp.MustCompile(`^/api/v3/update/associationtype/[0-9]+$`)
	deleteAssociationKindLatestRegexp = regexp.MustCompile(`^/api/v3/delete/associationtype/[0-9]+$`)
)

func (ps *parseStream) associationTypeLatest() *parseStream {
	if ps.err != nil {
		return ps
	}

	// find association kind operation
	if ps.hitPattern(findManyAssociationKindLatestPattern, http.MethodPost) {
		ps.Attribute.Resources = []auth.Resource{
			auth.Resource{
				Type:   auth.AssociationType,
				Action: auth.FindMany,
			},
		}
		return ps
	}

	// create association kind operation
	if ps.hitPattern(createAssociationKindLatestPattern, http.MethodPost) {
		ps.Attribute.Resources = []auth.Resource{
			auth.Resource{
				Type:   auth.AssociationType,
				Action: auth.Create,
			},
		}
		return ps
	}

	// update association kind operation
	if ps.hitRegexp(updateAssociationKindLatestRegexp, http.MethodPut) {
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
				Type:       auth.AssociationType,
				Action:     auth.Update,
				InstanceID: kindID,
			},
		}

		return ps
	}

	// delete association kind operation
	if ps.hitRegexp(deleteAssociationKindLatestRegexp, http.MethodDelete) {
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
				Type:       auth.AssociationType,
				Action:     auth.Delete,
				InstanceID: kindID,
			},
		}

		return ps
	}

	return ps
}

const (
	findObjectAssociationLatestPattern                    = "/api/v3/find/objectassociation"
	createObjectAssociationLatestPattern                  = "/api/v3/create/objectassociation"
	findObjectAssociationWithAssociationKindLatestPattern = "/api/v3/find/topoassociationtype"
)

var (
	updateObjectAssociationLatestRegexp = regexp.MustCompile(`^/api/v3/update/object/association/[0-9]+$`)
	deleteObjectAssociationLatestRegexp = regexp.MustCompile(`^/api/v3/delete/objectassociation/[0-9]+$`)
)

func (ps *parseStream) objectAssociationLatest() *parseStream {
	if ps.err != nil {
		return ps
	}

	// search object association operation
	if ps.hitPattern(findObjectAssociationLatestPattern, http.MethodPost) {
		ps.Attribute.Resources = []auth.Resource{
			auth.Resource{
				Type:   auth.ObjectAssociation,
				Action: auth.FindMany,
			},
		}
		return ps
	}

	// create object association operation
	if ps.hitPattern(createObjectAssociationLatestPattern, http.MethodPost) {
		ps.Attribute.Resources = []auth.Resource{
			auth.Resource{
				Type:   auth.ObjectAssociation,
				Action: auth.Create,
			},
		}
		return ps
	}

	// update object association operation
	if ps.hitRegexp(updateObjectAssociationLatestRegexp, http.MethodPut) {
		if len(ps.RequestCtx.Elements) != 5 {
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
				Type:       auth.ObjectAssociation,
				Action:     auth.Update,
				InstanceID: assoID,
			},
		}
		return ps
	}

	// delete object association operation
	if ps.hitRegexp(deleteObjectAssociationLatestRegexp, http.MethodDelete) {
		if len(ps.RequestCtx.Elements) != 5 {
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
				Type:       auth.ObjectAssociation,
				Action:     auth.Delete,
				InstanceID: assoID,
			},
		}
		return ps
	}

	// find object association with a association kind list.
	if ps.hitPattern(findObjectAssociationWithAssociationKindLatestPattern, http.MethodPost) {
		ps.Attribute.Resources = []auth.Resource{
			auth.Resource{
				Type:   auth.ObjectAssociation,
				Action: auth.FindMany,
			},
		}
		return ps
	}

	return ps
}

const (
	findObjectInstanceAssociationLatestPattern   = "/api/v3/find/instassociation"
	createObjectInstanceAssociationLatestPattern = "/api/v3/inst/association/action/create"
)

var (
	deleteObjectInstanceAssociationLatestRegexp = regexp.MustCompile("^/api/v3/delete/instassociation/[0-9]+$")
)

func (ps *parseStream) objectInstanceAssociationLatest() *parseStream {
	if ps.err != nil {
		return ps
	}

	// find object instance's association operation.
	if ps.hitPattern(findObjectInstanceAssociationLatestPattern, http.MethodPost) {
		ps.Attribute.Resources = []auth.Resource{
			auth.Resource{
				Type:   auth.ObjectInstanceAssociation,
				Action: auth.FindMany,
			},
		}
		return ps
	}

	// create object's instance association operation.
	if ps.hitPattern(createObjectInstanceAssociationLatestPattern, http.MethodPost) {
		ps.Attribute.Resources = []auth.Resource{
			auth.Resource{
				Type:   auth.ObjectInstanceAssociation,
				Action: auth.Create,
			},
		}
		return ps
	}

	// delete object's instance association operation.
	if ps.hitRegexp(deleteObjectInstanceAssociationLatestRegexp, http.MethodDelete) {
		assoID, err := strconv.ParseInt(ps.RequestCtx.Elements[4], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("delete object instance association, but got invalid association id %s", ps.RequestCtx.Elements[4])
			return ps
		}

		ps.Attribute.Resources = []auth.Resource{
			auth.Resource{
				Type:       auth.ObjectInstanceAssociation,
				Action:     auth.Delete,
				InstanceID: assoID,
			},
		}
		return ps
	}

	return ps
}

var (
	createObjectInstanceLatestRegexp          = regexp.MustCompile(`^/api/v3/create/inst/object/[\S][^/]+$`)
	findObjectInstanceLatestRegexp            = regexp.MustCompile(`^/api/v3/find/instassociation/object/[\S][^/]+$`)
	updateObjectInstanceLatestRegexp          = regexp.MustCompile(`^/api/v3/update/inst/object/[\S][^/]+/inst/[0-9]+$`)
	updateObjectInstanceBatchLatestRegexp     = regexp.MustCompile(`^/api/v3/updatemany/inst/object/[\S][^/]+$`)
	deleteObjectInstanceBatchLatestRegexp     = regexp.MustCompile(`^/api/v3/deletemany/inst/object/[\S][^/]+$`)
	deleteObjectInstanceLatestRegexp          = regexp.MustCompile(`^/api/v3/delete/inst/object/[\S][^/]+/inst/[0-9]+$`)
	findObjectInstanceSubTopologyLatestRegexp = regexp.MustCompile(`^/api/v3/find/insttopo/object/[\S][^/]+/inst/[0-9]+$`)
	findObjectInstanceTopologyLatestRegexp    = regexp.MustCompile(`^/api/v3/find/instassttopo/object/[\S][^/]+/inst/[0-9]+$`)
	findBusinessInstanceTopologyLatestRegexp  = regexp.MustCompile(`^/api/v3/find/topoinst/biz/[0-9]+$`)
	findObjectInstancesLatestRegexp           = regexp.MustCompile(`^/api/v3/find/inst/object/[\S][^/]+$`)
)

func (ps *parseStream) objectInstanceLatest() *parseStream {
	if ps.err != nil {
		return ps
	}

	// create object instance operation.
	if ps.hitRegexp(createObjectInstanceLatestRegexp, http.MethodPost) {
		ps.Attribute.Resources = []auth.Resource{
			auth.Resource{
				Type:   auth.ObjectInstance,
				Action: auth.Create,
			},
		}
		return ps
	}

	// find object instance operation.
	if ps.hitRegexp(findObjectInstanceLatestRegexp, http.MethodPost) {
		if len(ps.RequestCtx.Elements) != 6 {
			ps.err = errors.New("search object instance, but got invalid url")
			return ps
		}
		ps.Attribute.Resources = []auth.Resource{
			auth.Resource{
				Type:   auth.ObjectInstance,
				Action: auth.Find,
				Affiliated: auth.Affiliated{
					Type: auth.Object,
					Name: ps.RequestCtx.Elements[5],
				},
			},
		}
		return ps
	}

	// update object instance operation.
	if ps.hitRegexp(updateObjectInstanceLatestRegexp, http.MethodPut) {
		if len(ps.RequestCtx.Elements) != 8 {
			ps.err = errors.New("update object instance, but got invalid url")
			return ps
		}

		instID, err := strconv.ParseInt(ps.RequestCtx.Elements[7], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("update object instance, but got invalid instance id %s", ps.RequestCtx.Elements[5])
			return ps
		}

		ps.Attribute.Resources = []auth.Resource{
			auth.Resource{
				Type:       auth.ObjectInstance,
				Action:     auth.Update,
				InstanceID: instID,
				Affiliated: auth.Affiliated{
					Type: auth.Object,
					Name: ps.RequestCtx.Elements[5],
				},
			},
		}
		return ps
	}

	// update object instance batch operation.
	if ps.hitRegexp(updateObjectInstanceBatchLatestRegexp, http.MethodPut) {
		if len(ps.RequestCtx.Elements) != 6 {
			ps.err = errors.New("update object instance batch, but got invalid url")
			return ps
		}

		ps.Attribute.Resources = []auth.Resource{
			auth.Resource{
				Type:   auth.ObjectInstance,
				Action: auth.UpdateMany,
				Affiliated: auth.Affiliated{
					Type: auth.Object,
					Name: ps.RequestCtx.Elements[5],
				},
			},
		}
		return ps
	}

	// delete object instance batch operation.
	if ps.hitRegexp(deleteObjectInstanceBatchLatestRegexp, http.MethodDelete) {
		if len(ps.RequestCtx.Elements) != 6 {
			ps.err = errors.New("delete object instance batch, but got invalid url")
			return ps
		}

		ps.Attribute.Resources = []auth.Resource{
			auth.Resource{
				Type:   auth.ObjectInstance,
				Action: auth.DeleteMany,
				Affiliated: auth.Affiliated{
					Type: auth.Object,
					Name: ps.RequestCtx.Elements[5],
				},
			},
		}
		return ps
	}

	// delete object instance operation.
	if ps.hitRegexp(deleteObjectInstanceLatestRegexp, http.MethodDelete) {
		if len(ps.RequestCtx.Elements) != 8 {
			ps.err = errors.New("delete object instance, but got invalid url")
			return ps
		}

		instID, err := strconv.ParseInt(ps.RequestCtx.Elements[7], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("delete object instance, but got invalid instance id %s", ps.RequestCtx.Elements[7])
			return ps
		}

		ps.Attribute.Resources = []auth.Resource{
			auth.Resource{
				Type:       auth.ObjectInstance,
				Action:     auth.Delete,
				InstanceID: instID,
				Affiliated: auth.Affiliated{
					Type: auth.Object,
					Name: ps.RequestCtx.Elements[5],
				},
			},
		}
		return ps
	}

	// find object instance sub topology operation
	if ps.hitRegexp(findObjectInstanceSubTopologyLatestRegexp, http.MethodPost) {
		if len(ps.RequestCtx.Elements) != 8 {
			ps.err = errors.New("find object instance topology, but got invalid url")
			return ps
		}

		instID, err := strconv.ParseInt(ps.RequestCtx.Elements[7], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("find object instance topology, but got invalid instance id %s", ps.RequestCtx.Elements[7])
			return ps
		}

		ps.Attribute.Resources = []auth.Resource{
			auth.Resource{
				Type:       auth.ObjectInstanceTopology,
				Action:     auth.Find,
				InstanceID: instID,
				Affiliated: auth.Affiliated{
					Type: auth.Object,
					Name: ps.RequestCtx.Elements[5],
				},
			},
		}
		return ps
	}

	// find object instance fully topology operation.
	if ps.hitRegexp(findObjectInstanceTopologyLatestRegexp, http.MethodPost) {
		if len(ps.RequestCtx.Elements) != 8 {
			ps.err = errors.New("find object instance topology, but got invalid url")
			return ps
		}

		bizID, err := ps.RequestCtx.Metadata.Label.GetBusinessID()
		if err != nil {
			ps.err = fmt.Errorf("find object instance, but get object id in metadata failed, err: %v", err)
			return ps
		}

		instID, err := strconv.ParseInt(ps.RequestCtx.Elements[7], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("find object instance, but get instance id %s", ps.RequestCtx.Elements[7])
			return ps
		}

		ps.Attribute.Resources = []auth.Resource{
			auth.Resource{
				Type:       auth.ObjectInstanceTopology,
				Action:     auth.Find,
				InstanceID: instID,
				BusinessID: bizID,
				Affiliated: auth.Affiliated{
					Type: auth.Object,
					Name: ps.RequestCtx.Elements[5],
				},
			},
		}

		return ps
	}

	// find business instance topology operation.
	if ps.hitRegexp(findBusinessInstanceTopologyLatestRegexp, http.MethodGet) {
		if len(ps.RequestCtx.Elements) != 6 {
			ps.err = errors.New("find business instance topology, but got invalid url")
			return ps
		}

		bizID, err := ps.RequestCtx.Metadata.Label.GetBusinessID()
		if err != nil {
			ps.err = fmt.Errorf("find business instance, but get business id in metadata failed, err: %v", err)
			return ps
		}

		ps.Attribute.Resources = []auth.Resource{
			auth.Resource{
				Type:       auth.ObjectInstanceTopology,
				Action:     auth.Find,
				BusinessID: bizID,
				Affiliated: auth.Affiliated{
					Type: auth.Object,
					Name: string(auth.Business),
				},
			},
		}
		return ps
	}

	// find object's instance list operation
	if ps.hitRegexp(findObjectInstancesLatestRegexp, http.MethodPost) {
		if len(ps.RequestCtx.Elements) != 6 {
			ps.err = errors.New("find object's instance list, but got invalid url")
			return ps
		}

		ps.Attribute.Resources = []auth.Resource{
			auth.Resource{
				Type:   auth.ObjectInstanceTopology,
				Action: auth.FindMany,
				Affiliated: auth.Affiliated{
					Type: auth.Object,
					Name: ps.RequestCtx.Elements[5],
				},
			},
		}
		return ps
	}

	return ps
}

const (
	createObjectLatestPattern       = "/api/v3/create/object"
	findObjectsLatestPattern        = "/api/v3/find/object"
	findObjectTopologyLatestPattern = "/api/v3/find/objects/objecttopology"
)

var (
	deleteObjectLatestRegexp                = regexp.MustCompile(`^/api/v3/delete/object/[0-9]+$`)
	updateObjectLatestRegexp                = regexp.MustCompile(`^/api/v3/update/object/[0-9]+$`)
	findObjectTopologyGraphicLatestRegexp   = regexp.MustCompile(`^/api/v3/find/objecttopo/scope_type/[\S][^/]+/scope_id/[\S][^/]+$`)
	updateObjectTopologyGraphicLatestRegexp = regexp.MustCompile(`^/api/v3/update/objecttopo/scope_type/[\S][^/]+/scope_id/[\S][^/]+$`)
)

func (ps *parseStream) objectLatest() *parseStream {
	if ps.err != nil {
		return ps
	}

	// create common object operation.
	if ps.hitPattern(createObjectLatestPattern, http.MethodPost) {
		bizID, err := ps.RequestCtx.Metadata.Label.GetBusinessID()
		if err != nil {
			ps.err = fmt.Errorf("create object, but get business id in metadata failed, err: %v", err)
			return ps
		}

		ps.Attribute.Resources = []auth.Resource{
			auth.Resource{
				Type:       auth.Object,
				BusinessID: bizID,
				Action:     auth.Create,
			},
		}
		return ps
	}

	// delete object operation
	if ps.hitRegexp(deleteObjectLatestRegexp, http.MethodDelete) {
		if len(ps.RequestCtx.Elements) != 5 {
			ps.err = errors.New("delete object, but got invalid url")
			return ps
		}

		objID, err := strconv.ParseInt(ps.RequestCtx.Elements[4], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("delete object, but got invalid object's id %s", ps.RequestCtx.Elements[3])
			return ps
		}

		bizID, err := ps.RequestCtx.Metadata.Label.GetBusinessID()
		if err != nil {
			ps.err = fmt.Errorf("delete object, but get business id in metadata failed, err: %v", err)
			return ps
		}

		ps.Attribute.Resources = []auth.Resource{
			auth.Resource{
				Type:       auth.Object,
				Action:     auth.Delete,
				BusinessID: bizID,
				InstanceID: objID,
			},
		}
		return ps
	}

	// update object operation.
	if ps.hitRegexp(updateObjectLatestRegexp, http.MethodPut) {
		if len(ps.RequestCtx.Elements) != 5 {
			ps.err = errors.New("update object, but got invalid url")
			return ps
		}

		objID, err := strconv.ParseInt(ps.RequestCtx.Elements[4], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("update object, but got invalid object's id %s", ps.RequestCtx.Elements[4])
			return ps
		}

		bizID, err := ps.RequestCtx.Metadata.Label.GetBusinessID()
		if err != nil {
			ps.err = fmt.Errorf("update object, but get business id in metadata failed, err: %v", err)
			return ps
		}

		ps.Attribute.Resources = []auth.Resource{
			auth.Resource{
				Type:       auth.Object,
				Action:     auth.Update,
				BusinessID: bizID,
				InstanceID: objID,
			},
		}
		return ps
	}

	// get object operation.
	if ps.hitPattern(findObjectsLatestPattern, http.MethodPost) {
		bizID, err := ps.RequestCtx.Metadata.Label.GetBusinessID()
		if err != nil {
			ps.err = fmt.Errorf("find object, but get business id in metadata failed, err: %v", err)
			return ps
		}

		ps.Attribute.Resources = []auth.Resource{
			auth.Resource{
				Type:       auth.Object,
				BusinessID: bizID,
				Action:     auth.FindMany,
			},
		}
		return ps
	}

	// find object's topology operation.
	if ps.hitPattern(findObjectTopologyLatestPattern, http.MethodPost) {
		bizID, err := ps.RequestCtx.Metadata.Label.GetBusinessID()
		if err != nil {
			ps.err = fmt.Errorf("find object, but get business id in metadata failed, err: %v", err)
			return ps
		}
		ps.Attribute.Resources = []auth.Resource{
			auth.Resource{
				Type:       auth.ObjectTopology,
				BusinessID: bizID,
				Action:     auth.Find,
			},
		}
		return ps
	}

	// find object's topology graphic operation.
	if ps.hitRegexp(findObjectTopologyGraphicLatestRegexp, http.MethodPost) {
		bizID, err := ps.RequestCtx.Metadata.Label.GetBusinessID()
		if err != nil {
			ps.err = fmt.Errorf("find object topology graphic, but get business id in metadata failed, err: %v", err)
			return ps
		}
		ps.Attribute.Resources = []auth.Resource{
			auth.Resource{
				Type:       auth.ObjectTopology,
				BusinessID: bizID,
				Action:     auth.Find,
			},
		}
		return ps
	}

	// update object's topology graphic operation.
	// TODO: confirm if bizID is needed.
	if ps.hitRegexp(updateObjectTopologyGraphicLatestRegexp, http.MethodPost) {

		ps.Attribute.Resources = []auth.Resource{
			auth.Resource{
				Type:   auth.ObjectTopology,
				Action: auth.Update,
			},
		}
		return ps
	}

	return ps
}

const (
	createObjectClassificationLatestPattern         = "/api/v3/create/objectclassification"
	findObjectClassificationListLatestPattern       = "/api/v3/find/objectclassification"
	findObjectsBelongsToClassificationLatestPattern = `/api/v3/find/classificationobject`
)

var (
	deleteObjectClassificationLatestRegexp = regexp.MustCompile("^/api/v3/delete/objectclassification/[0-9]+$")
	updateObjectClassificationLatestRegexp = regexp.MustCompile("^/api/v3/update/objectclassification/[0-9]+$")
)

func (ps *parseStream) ObjectClassificationLatest() *parseStream {
	if ps.err != nil {
		return ps
	}

	// create object's classification operation.
	if ps.hitPattern(createObjectClassificationLatestPattern, http.MethodPost) {
		bizID, err := ps.RequestCtx.Metadata.Label.GetBusinessID()
		if err != nil {
			ps.err = fmt.Errorf("get business id in metadata failed, err: %v", err)
			return ps
		}
		ps.Attribute.Resources = []auth.Resource{
			auth.Resource{
				Type:       auth.ObjectClassification,
				BusinessID: bizID,
				Action:     auth.Create,
			},
		}
		return ps
	}

	// delete object's classification operation.
	if ps.hitRegexp(deleteObjectClassificationLatestRegexp, http.MethodDelete) {
		if len(ps.RequestCtx.Elements) != 5 {
			ps.err = errors.New("delete object classification, but got invalid url")
			return ps
		}

		classID, err := strconv.ParseInt(ps.RequestCtx.Elements[4], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("delete object classification, but got invalid object's id %s", ps.RequestCtx.Elements[4])
			return ps
		}

		bizID, err := ps.RequestCtx.Metadata.Label.GetBusinessID()
		if err != nil {
			ps.err = fmt.Errorf("get business id in metadata failed, err: %v", err)
			return ps
		}

		ps.Attribute.Resources = []auth.Resource{
			auth.Resource{
				Type:       auth.ObjectClassification,
				Action:     auth.Delete,
				BusinessID: bizID,
				InstanceID: classID,
			},
		}
		return ps
	}

	// update object's classification operation.
	if ps.hitRegexp(updateObjectClassificationLatestRegexp, http.MethodPut) {
		if len(ps.RequestCtx.Elements) != 5 {
			ps.err = errors.New("update object classification, but got invalid url")
			return ps
		}

		classID, err := strconv.ParseInt(ps.RequestCtx.Elements[4], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("update object classification, but got invalid object's  classification id %s", ps.RequestCtx.Elements[4])
			return ps
		}

		bizID, err := ps.RequestCtx.Metadata.Label.GetBusinessID()
		if err != nil {
			ps.err = fmt.Errorf("get business id in metadata failed, err: %v", err)
			return ps
		}

		ps.Attribute.Resources = []auth.Resource{
			auth.Resource{
				Type:       auth.ObjectClassification,
				Action:     auth.Update,
				BusinessID: bizID,
				InstanceID: classID,
			},
		}
		return ps
	}

	// find object's classification list operation.
	if ps.hitPattern(findObjectClassificationListLatestPattern, http.MethodPost) {
		bizID, err := ps.RequestCtx.Metadata.Label.GetBusinessID()
		if err != nil {
			ps.err = fmt.Errorf("get business id in metadata failed, err: %v", err)
			return ps
		}
		ps.Attribute.Resources = []auth.Resource{
			auth.Resource{
				Type:       auth.ObjectClassification,
				BusinessID: bizID,
				Action:     auth.FindMany,
			},
		}
		return ps
	}

	// find all the objects belongs to a classification
	if ps.hitPattern(findObjectsBelongsToClassificationLatestPattern, http.MethodPost) {
		bizID, err := ps.RequestCtx.Metadata.Label.GetBusinessID()
		if err != nil {
			ps.err = fmt.Errorf("get business id in metadata failed, err: %v", err)
			return ps
		}
		ps.Attribute.Resources = []auth.Resource{
			auth.Resource{
				Type:       auth.ObjectClassification,
				BusinessID: bizID,
				Action:     auth.FindMany,
			},
		}
		return ps
	}

	return ps
}

const (
	createObjectAttributeGroupLatestPattern = "/api/v3/create/objectattgroup"
	updateObjectAttributeGroupLatestPattern = "/api/v3/update/objectattgroup"
)

var (
	findObjectAttributeGroupLatestRegexp     = regexp.MustCompile(`^/api/v3/find/objectattgroup/object/[\S][^/]+$`)
	deleteObjectAttributeGroupLatestRegexp   = regexp.MustCompile(`^/api/v3/delete/objectattgroup/[0-9]+$`)
	removeAttributeAwayFromGroupLatestRegexp = regexp.MustCompile(`^/api/v3/delete/objectattgroupasst/object/[\S][^/]+/property/[\S][^/]+/group/[\S][^/]+$`)
)

func (ps *parseStream) objectAttributeGroupLatest() *parseStream {
	if ps.err != nil {
		return ps
	}

	// create object's attribute group operation.
	if ps.hitPattern(createObjectAttributeGroupLatestPattern, http.MethodPost) {
		bizID, err := ps.RequestCtx.Metadata.Label.GetBusinessID()
		if err != nil {
			ps.err = fmt.Errorf("get business id in metadata failed, err: %v", err)
			return ps
		}
		ps.Attribute.Resources = []auth.Resource{
			auth.Resource{
				Type:       auth.ObjectAttributeGroup,
				BusinessID: bizID,
				Action:     auth.Create,
			},
		}
		return ps
	}

	// find object's attribute group operation.
	if ps.hitRegexp(findObjectAttributeGroupLatestRegexp, http.MethodPost) {
		if len(ps.RequestCtx.Elements) != 6 {
			ps.err = errors.New("find object's attribute group, but got invalid uri")
			return ps
		}

		bizID, err := ps.RequestCtx.Metadata.Label.GetBusinessID()
		if err != nil {
			ps.err = fmt.Errorf("get business id in metadata failed, err: %v", err)
			return ps
		}

		ps.Attribute.Resources = []auth.Resource{
			auth.Resource{
				Type:       auth.ObjectAttributeGroup,
				Action:     auth.Find,
				BusinessID: bizID,
				Affiliated: auth.Affiliated{
					Type: auth.Object,
					Name: ps.RequestCtx.Elements[5],
				},
			},
		}
		return ps
	}

	// update object's attribute group operation.
	if ps.hitPattern(updateObjectAttributeGroupLatestPattern, http.MethodPut) {
		bizID, err := ps.RequestCtx.Metadata.Label.GetBusinessID()
		if err != nil {
			ps.err = fmt.Errorf("get business id in metadata failed, err: %v", err)
			return ps
		}
		ps.Attribute.Resources = []auth.Resource{
			auth.Resource{
				Type:       auth.ObjectClassification,
				BusinessID: bizID,
				Action:     auth.Update,
			},
		}
		return ps
	}

	// delete object's attribute group operation.
	if ps.hitRegexp(deleteObjectAttributeGroupLatestRegexp, http.MethodDelete) {
		if len(ps.RequestCtx.Elements) != 5 {
			ps.err = errors.New("delete object's attribute group, but got invalid url")
			return ps
		}

		groupID, err := strconv.ParseInt(ps.RequestCtx.Elements[4], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("delete object's attribute group, but got invalid group's id %s", ps.RequestCtx.Elements[4])
			return ps
		}

		bizID, err := ps.RequestCtx.Metadata.Label.GetBusinessID()
		if err != nil {
			ps.err = fmt.Errorf("get business id in metadata failed, err: %v", err)
			return ps
		}

		ps.Attribute.Resources = []auth.Resource{
			auth.Resource{
				Type:       auth.ObjectAttributeGroup,
				Action:     auth.Delete,
				BusinessID: bizID,
				InstanceID: groupID,
			},
		}
		return ps
	}

	// remove a object's attribute away from a group.
	if ps.hitRegexp(removeAttributeAwayFromGroupLatestRegexp, http.MethodDelete) {
		if len(ps.RequestCtx.Elements) != 12 {
			ps.err = errors.New("remove a object attribute away from a group, but got invalid uri")
			return ps
		}

		bizID, err := ps.RequestCtx.Metadata.Label.GetBusinessID()
		if err != nil {
			ps.err = fmt.Errorf("get business id in metadata failed, err: %v", err)
			return ps
		}

		ps.Attribute.Resources = []auth.Resource{
			auth.Resource{
				Type:       auth.ObjectAttributeGroup,
				Name:       ps.RequestCtx.Elements[11],
				BusinessID: bizID,
				Action:     auth.Delete,
			},
		}
		return ps
	}

	return ps
}

const (
	createObjectAttributeLatestPattern = "/api/v3/create/objectattr"
	findObjectAttributeLatestPattern   = "/api/v3/find/objectattr"
)

var (
	deleteObjectAttributeLatestRegexp = regexp.MustCompile(`^/api/v3/delete/objectattr/[0-9]+$`)
	updateObjectAttributeLatestRegexp = regexp.MustCompile(`^/api/v3/update/objectattr/[0-9]+$`)
)

func (ps *parseStream) objectAttributeLatest() *parseStream {
	if ps.err != nil {
		return ps
	}

	// create object's attribute operation.
	if ps.hitPattern(createObjectAttributeLatestPattern, http.MethodPost) {
		bizID, err := ps.RequestCtx.Metadata.Label.GetBusinessID()
		if err != nil {
			ps.err = fmt.Errorf("get business id in metadata failed, err: %v", err)
			return ps
		}
		ps.Attribute.Resources = []auth.Resource{
			auth.Resource{
				Type:       auth.ObjectAttribute,
				BusinessID: bizID,
				Action:     auth.Create,
			},
		}
		return ps
	}

	// delete object's attribute operation.
	if ps.hitRegexp(deleteObjectAttributeLatestRegexp, http.MethodDelete) {
		if len(ps.RequestCtx.Elements) != 5 {
			ps.err = errors.New("delete object attribute, but got invalid url")
			return ps
		}

		attrID, err := strconv.ParseInt(ps.RequestCtx.Elements[4], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("delete object attribute, but got invalid attribute id %s", ps.RequestCtx.Elements[4])
			return ps
		}

		bizID, err := ps.RequestCtx.Metadata.Label.GetBusinessID()
		if err != nil {
			ps.err = fmt.Errorf("get business id in metadata failed, err: %v", err)
			return ps
		}

		ps.Attribute.Resources = []auth.Resource{
			auth.Resource{
				Type:       auth.ObjectAttribute,
				Action:     auth.Delete,
				BusinessID: bizID,
				InstanceID: attrID,
			},
		}
		return ps
	}

	// update object attribute operation
	if ps.hitRegexp(updateObjectAttributeLatestRegexp, http.MethodPut) {
		if len(ps.RequestCtx.Elements) != 5 {
			ps.err = errors.New("update object attribute, but got invalid url")
			return ps
		}

		attrID, err := strconv.ParseInt(ps.RequestCtx.Elements[4], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("update object attribute, but got invalid attribute id %s", ps.RequestCtx.Elements[4])
			return ps
		}

		bizID, err := ps.RequestCtx.Metadata.Label.GetBusinessID()
		if err != nil {
			ps.err = fmt.Errorf("get business id in metadata failed, err: %v", err)
			return ps
		}

		ps.Attribute.Resources = []auth.Resource{
			auth.Resource{
				Type:       auth.ObjectAttribute,
				Action:     auth.Update,
				BusinessID: bizID,
				InstanceID: attrID,
			},
		}
		return ps
	}

	// get object's attribute operation.
	if ps.hitPattern(findObjectAttributeLatestPattern, http.MethodPost) {
		bizID, err := ps.RequestCtx.Metadata.Label.GetBusinessID()
		if err != nil {
			ps.err = fmt.Errorf("get business id in metadata failed, err: %v", err)
			return ps
		}

		ps.Attribute.Resources = []auth.Resource{
			auth.Resource{
				Type:       auth.ObjectAttribute,
				BusinessID: bizID,
				Action:     auth.Find,
			},
		}
		return ps
	}

	return ps
}

const (
	createMainlineObjectLatestPattern   = "/api/v3/create/topomodelmainline"
	findMainlineObjectTopoLatestPattern = "/api/v3/find/topomodelmainline"
)

var (
	deleteMainlineObjectLatestRegexp        = regexp.MustCompile(`^/api/v3/delete/topomodelmainline/object/[\S][^/]+$`)
	findMainlineInstanceTopoLatestRegexp    = regexp.MustCompile(`^/api/v3/find/topoinst/biz/[0-9]+$`)
	findMainineSubInstanceTopoLatestRegexp  = regexp.MustCompile(`^/api/v3/topoinstchild/object/[\S][^/]+/biz/[0-9]+/inst/[0-9]+$`)
	findMainlineIdleFaultModuleLatestRegexp = regexp.MustCompile(`^/api/v3/find/topointernal/biz/[0-9]+$`)
)

func (ps *parseStream) mainlineLatest() *parseStream {
	if ps.err != nil {
		return ps
	}

	// create mainline object operation.
	if ps.hitPattern(createMainlineObjectLatestPattern, http.MethodPost) {
		bizID, err := ps.RequestCtx.Metadata.Label.GetBusinessID()
		if err != nil {
			ps.err = fmt.Errorf("get business id in metadata failed, err: %v", err)
			return ps
		}

		ps.Attribute.Resources = []auth.Resource{
			auth.Resource{
				Type:       auth.MainlineObject,
				BusinessID: bizID,
				Action:     auth.Create,
			},
		}
		return ps
	}

	// delete mainline object operation
	if ps.hitRegexp(deleteMainlineObjectLatestRegexp, http.MethodDelete) {
		bizID, err := ps.RequestCtx.Metadata.Label.GetBusinessID()
		if err != nil {
			ps.err = fmt.Errorf("get business id in metadata failed, err: %v", err)
			return ps
		}

		ps.Attribute.Resources = []auth.Resource{
			auth.Resource{
				Type:       auth.MainlineObject,
				BusinessID: bizID,
				Action:     auth.Delete,
			},
		}

		return ps
	}

	// get mainline object operation
	if ps.hitPattern(findMainlineObjectTopoLatestPattern, http.MethodGet) {
		bizID, err := ps.RequestCtx.Metadata.Label.GetBusinessID()
		if err != nil {
			ps.err = fmt.Errorf("get business id in metadata failed, err: %v", err)
			return ps
		}
		ps.Attribute.Resources = []auth.Resource{
			auth.Resource{
				Type:       auth.MainlineObjectTopology,
				BusinessID: bizID,
				Action:     auth.Find,
			},
		}

		return ps
	}

	// find mainline instance topology operation.
	// TODO: confirm this api about multiple biz filed in url and metadata.
	if ps.hitRegexp(findMainlineInstanceTopoLatestRegexp, http.MethodGet) {
		bizID, err := ps.RequestCtx.Metadata.Label.GetBusinessID()
		if err != nil {
			ps.err = fmt.Errorf("get business id in metadata failed, err: %v", err)
			return ps
		}

		ps.Attribute.Resources = []auth.Resource{
			auth.Resource{
				Type:       auth.MainlineInstanceTopology,
				Action:     auth.Find,
				BusinessID: bizID,
			},
		}

		return ps
	}

	// find mainline object instance's sub-instance topology operation.
	if ps.hitRegexp(findMainineSubInstanceTopoLatestRegexp, http.MethodGet) {
		if len(ps.RequestCtx.Elements) != 9 {
			ps.err = errors.New("find mainline object's sub instance topology, but got invalid url")
			return ps
		}

		bizID, err := strconv.ParseInt(ps.RequestCtx.Elements[6], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("find mainline object's sub instance topology, but got invalid business id %s", ps.RequestCtx.Elements[6])
			return ps
		}

		ps.Attribute.Resources = []auth.Resource{
			auth.Resource{
				Type:       auth.MainlineInstanceTopology,
				Action:     auth.Find,
				BusinessID: bizID,
			},
		}

		return ps
	}

	// find mainline internal idle and fault module operation.
	if ps.hitRegexp(findMainlineIdleFaultModuleLatestRegexp, http.MethodGet) {
		if len(ps.RequestCtx.Elements) != 6 {
			ps.err = errors.New("find mainline idle and fault module, but got invalid url")
			return ps
		}

		bizID, err := strconv.ParseInt(ps.RequestCtx.Elements[5], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("find mainline idle and fault module, but got invalid business id %s", ps.RequestCtx.Elements[5])
			return ps
		}

		ps.Attribute.Resources = []auth.Resource{
			auth.Resource{
				Type:       auth.MainlineObject,
				Action:     auth.Find,
				BusinessID: bizID,
			},
		}

		return ps
	}

	return ps
}
