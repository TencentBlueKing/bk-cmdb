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

	"configcenter/src/auth/meta"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"

	"github.com/tidwall/gjson"
)

// this package's topology filter is the latest api version
// for these resources, it also has a elder resource api version.
// TODO: if the elder api has been removed, delete their resource
// filter at the same time.

func (ps *parseStream) topologyLatest() *parseStream {
	if ps.shouldReturn() {
		return ps
	}

	ps.objectUniqueLatest().
		associationTypeLatest().
		objectAssociationLatest().
		objectInstanceAssociationLatest().
		objectInstanceLatest().
		objectLatest().
		ObjectClassificationLatest().
		objectAttributeGroupLatest().
		objectAttributeLatest().
		mainlineLatest().
		SetTemplate()

	return ps
}

var (
	createObjectUniqueLatestRegexp = regexp.MustCompile(`^/api/v3/create/objectunique/object/[^\s/]+/?$`)
	updateObjectUniqueLatestRegexp = regexp.MustCompile(`^/api/v3/update/objectunique/object/[^\s/]+/unique/[0-9]+/?$`)
	deleteObjectUniqueLatestRegexp = regexp.MustCompile(`^/api/v3/delete/objectunique/object/[^\s/]+/unique/[0-9]+/?$`)
	findObjectUniqueLatestRegexp   = regexp.MustCompile(`^/api/v3/find/objectunique/object/[^\s/]+/?$`)
)

func (ps *parseStream) objectUniqueLatest() *parseStream {
	if ps.shouldReturn() {
		return ps
	}

	// TODO: add business id for these filter rules to resources.
	// add object unique operation.
	if ps.hitRegexp(createObjectUniqueLatestRegexp, http.MethodPost) {
		bizID, err := metadata.BizIDFromMetadata(ps.RequestCtx.Metadata)
		if err != nil {
			ps.err = err
			return ps
		}
		model, err := ps.getOneModel(mapstr.MapStr{common.BKObjIDField: ps.RequestCtx.Elements[5]})
		if err != nil {
			ps.err = err
			return ps
		}
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:   meta.ModelUnique,
					Action: meta.Create,
				},
				Layers: []meta.Item{{Type: meta.Model, InstanceID: model.ID}},
			},
		}
		return ps
	}

	// update object unique operation.
	if ps.hitRegexp(updateObjectUniqueLatestRegexp, http.MethodPut) {
		bizID, err := metadata.BizIDFromMetadata(ps.RequestCtx.Metadata)
		if err != nil {
			ps.err = err
			return ps
		}
		uniqueID, err := strconv.ParseInt(ps.RequestCtx.Elements[7], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("update object unique, but got invalid unique id %s", ps.RequestCtx.Elements[7])
			return ps
		}
		model, err := ps.getOneModel(mapstr.MapStr{common.BKObjIDField: ps.RequestCtx.Elements[5]})
		if err != nil {
			ps.err = err
			return ps
		}

		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:       meta.ModelUnique,
					Action:     meta.Update,
					InstanceID: uniqueID,
				},
				Layers: []meta.Item{{Type: meta.Model, InstanceID: model.ID}},
			},
		}
		return ps
	}

	// delete object unique operation.
	if ps.hitRegexp(deleteObjectUniqueLatestRegexp, http.MethodPost) {
		bizID, err := metadata.BizIDFromMetadata(ps.RequestCtx.Metadata)
		if err != nil {
			ps.err = err
			return ps
		}
		uniqueID, err := strconv.ParseInt(ps.RequestCtx.Elements[7], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("update object unique, but got invalid unique id %s", ps.RequestCtx.Elements[7])
			return ps
		}
		model, err := ps.getOneModel(mapstr.MapStr{common.BKObjIDField: ps.RequestCtx.Elements[5]})
		if err != nil {
			ps.err = err
			return ps
		}

		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:       meta.ModelUnique,
					Action:     meta.Delete,
					InstanceID: uniqueID,
				},
				Layers:     []meta.Item{{Type: meta.Model, InstanceID: model.ID}},
				BusinessID: bizID,
			},
		}
		return ps
	}

	// find model unique operation
	if ps.hitRegexp(findObjectUniqueLatestRegexp, http.MethodPost) {
		bizID, err := metadata.BizIDFromMetadata(ps.RequestCtx.Metadata)
		if err != nil {
			ps.err = err
			return ps
		}
		model, err := ps.getOneModel(mapstr.MapStr{common.BKObjIDField: ps.RequestCtx.Elements[5]})
		if err != nil {
			ps.err = err
			return ps
		}
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:   meta.ModelUnique,
					Action: meta.FindMany,
				},
				Layers:     []meta.Item{{Type: meta.Model, InstanceID: model.ID}},
				BusinessID: bizID,
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
	updateAssociationKindLatestRegexp = regexp.MustCompile(`^/api/v3/update/associationtype/[0-9]+/?$`)
	deleteAssociationKindLatestRegexp = regexp.MustCompile(`^/api/v3/delete/associationtype/[0-9]+/?$`)
)

func (ps *parseStream) associationTypeLatest() *parseStream {
	if ps.shouldReturn() {
		return ps
	}

	// find association kind operation
	if ps.hitPattern(findManyAssociationKindLatestPattern, http.MethodPost) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:   meta.AssociationType,
					Action: meta.FindMany,
				},
			},
		}
		return ps
	}

	// create association kind operation
	if ps.hitPattern(createAssociationKindLatestPattern, http.MethodPost) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:   meta.AssociationType,
					Action: meta.Create,
				},
			},
		}
		return ps
	}

	// update association kind operation
	if ps.hitRegexp(updateAssociationKindLatestRegexp, http.MethodPut) {
		kindID, err := strconv.ParseInt(ps.RequestCtx.Elements[4], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("update association kind, but got invalid kind id %s", ps.RequestCtx.Elements[5])
			return ps
		}
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:       meta.AssociationType,
					Action:     meta.Update,
					InstanceID: kindID,
				},
			},
		}

		return ps
	}

	// delete association kind operation
	if ps.hitRegexp(deleteAssociationKindLatestRegexp, http.MethodDelete) {
		kindID, err := strconv.ParseInt(ps.RequestCtx.Elements[4], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("delete association kind, but got invalid kind id %s", ps.RequestCtx.Elements[5])
			return ps
		}
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:       meta.AssociationType,
					Action:     meta.Delete,
					InstanceID: kindID,
				},
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
	updateObjectAssociationLatestRegexp = regexp.MustCompile(`^/api/v3/update/objectassociation/[0-9]+/?$`)
	deleteObjectAssociationLatestRegexp = regexp.MustCompile(`^/api/v3/delete/objectassociation/[0-9]+/?$`)
)

func (ps *parseStream) objectAssociationLatest() *parseStream {
	if ps.shouldReturn() {
		return ps
	}

	// search object association operation
	if ps.hitPattern(findObjectAssociationLatestPattern, http.MethodPost) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:   meta.ModelAssociation,
					Action: meta.FindMany,
				},
			},
		}
		return ps
	}

	// create object association operation
	if ps.hitPattern(createObjectAssociationLatestPattern, http.MethodPost) {
		bizID, err := metadata.BizIDFromMetadata(ps.RequestCtx.Metadata)
		if err != nil {
			blog.Warnf("get business id in metadata failed, err: %v", err)
		}

		filter := mapstr.MapStr{
			common.BKObjIDField: mapstr.MapStr{
				common.BKDBIN: []interface{}{
					gjson.GetBytes(ps.RequestCtx.Body, common.BKObjIDField).Value(),
					gjson.GetBytes(ps.RequestCtx.Body, common.BKAsstObjIDField).Value(),
				},
			},
		}
		models, err := ps.searchModels(filter)
		if err != nil {
			ps.err = err
			return ps
		}

		for _, model := range models {
			ps.Attribute.Resources = append(ps.Attribute.Resources,
				meta.ResourceAttribute{
					BusinessID: bizID,
					Basic: meta.Basic{
						Type:       meta.Model,
						Action:     meta.Update,
						InstanceID: model.ID,
					},
				},
			)
		}
		return ps
	}

	// update object association operation
	if ps.hitRegexp(updateObjectAssociationLatestRegexp, http.MethodPut) {
		bizID, err := metadata.BizIDFromMetadata(ps.RequestCtx.Metadata)
		if err != nil {
			blog.Warnf("get business id in metadata failed, err: %v", err)
		}

		if len(ps.RequestCtx.Elements) != 5 {
			ps.err = errors.New("update object association, but got invalid url")
			return ps
		}

		assoID, err := strconv.ParseInt(ps.RequestCtx.Elements[4], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("update object association, but got invalid association id %s", ps.RequestCtx.Elements[4])
			return ps
		}
		asst, err := ps.getModelAssociation(mapstr.MapStr{common.BKFieldID: assoID})
		if err != nil {
			ps.err = err
			return ps
		}

		filter := mapstr.MapStr{
			common.BKObjIDField: mapstr.MapStr{
				common.BKDBIN: []interface{}{
					asst[0].ObjectID,
					asst[0].AsstObjID,
				},
			},
		}
		models, err := ps.searchModels(filter)
		if err != nil {
			ps.err = err
			return ps
		}

		for _, model := range models {
			ps.Attribute.Resources = append(ps.Attribute.Resources,
				meta.ResourceAttribute{
					Basic: meta.Basic{
						Type:       meta.Model,
						Action:     meta.Update,
						InstanceID: model.ID,
					},
					BusinessID: bizID,
				})
		}

		return ps
	}

	// delete object association operation
	if ps.hitRegexp(deleteObjectAssociationLatestRegexp, http.MethodDelete) {
		bizID, err := metadata.BizIDFromMetadata(ps.RequestCtx.Metadata)
		if err != nil {
			blog.Warnf("get business id in metadata failed, err: %v", err)
		}
		if len(ps.RequestCtx.Elements) != 5 {
			ps.err = errors.New("delete object association, but got invalid url")
			return ps
		}

		assoID, err := strconv.ParseInt(ps.RequestCtx.Elements[4], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("delete object association, but got invalid association id %s", ps.RequestCtx.Elements[4])
			return ps
		}

		asst, err := ps.getModelAssociation(mapstr.MapStr{common.BKFieldID: assoID})
		if err != nil {
			ps.err = err
			return ps
		}

		filter := mapstr.MapStr{
			common.BKObjIDField: mapstr.MapStr{
				common.BKDBIN: []interface{}{
					asst[0].ObjectID,
					asst[0].AsstObjID,
				},
			},
		}
		models, err := ps.searchModels(filter)
		if err != nil {
			ps.err = err
			return ps
		}

		for _, model := range models {
			ps.Attribute.Resources = append(ps.Attribute.Resources,
				meta.ResourceAttribute{
					Basic: meta.Basic{
						Type:       meta.Model,
						Action:     meta.Update,
						InstanceID: model.ID,
					},
					BusinessID: bizID,
				})
		}
		return ps
	}

	// find object association with a association kind list.
	if ps.hitPattern(findObjectAssociationWithAssociationKindLatestPattern, http.MethodPost) {
		bizID, err := metadata.BizIDFromMetadata(ps.RequestCtx.Metadata)
		if err != nil {
			ps.err = fmt.Errorf("parse bizID from metadata failed, err: %s", err.Error())
			return ps
		}
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:   meta.ModelAssociation,
					Action: meta.FindMany,
				},
			},
		}
		return ps
	}

	return ps
}

const (
	findObjectInstanceAssociationLatestPattern   = "/api/v3/find/instassociation"
	createObjectInstanceAssociationLatestPattern = "/api/v3/create/instassociation"
)

var (
	deleteObjectInstanceAssociationLatestRegexp = regexp.MustCompile("^/api/v3/delete/instassociation/[0-9]+/?$")
	findObjectInstanceTopologyUILatestRegexp    = regexp.MustCompile(`^/api/v3/findmany/inst/association/object/[^\s/]+/inst_id/[0-9]+/offset/[0-9]+/limit/[0-9]+/web$`)
	findInstAssociationObjInstInfoLatestRegexp  = regexp.MustCompile(`^/api/v3/findmany/inst/association/association_object/inst_base_info$`)
)

func (ps *parseStream) objectInstanceAssociationLatest() *parseStream {
	if ps.shouldReturn() {
		return ps
	}

	// find instance's association operation.
	if ps.hitPattern(findObjectInstanceAssociationLatestPattern, http.MethodPost) {
		bizID, err := metadata.BizIDFromMetadata(ps.RequestCtx.Metadata)
		if err != nil {
			ps.err = err
			return ps
		}
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:   meta.ModelInstanceAssociation,
					Action: meta.FindMany,
				},
			},
		}
		return ps
	}

	// create instance association operation.
	if ps.hitPattern(createObjectInstanceAssociationLatestPattern, http.MethodPost) {
		bizID, err := metadata.BizIDFromMetadata(ps.RequestCtx.Metadata)
		if err != nil {
			ps.err = err
			return ps
		}
		associationObjAsstID := gjson.GetBytes(ps.RequestCtx.Body, common.AssociationObjAsstIDField).String()
		filter := mapstr.MapStr{
			common.AssociationObjAsstIDField: associationObjAsstID,
		}
		asst, err := ps.getModelAssociation(filter)
		if err != nil {
			ps.err = err
			return ps
		}

		modelFilter := mapstr.MapStr{
			common.BKObjIDField: mapstr.MapStr{
				common.BKDBIN: []interface{}{
					asst[0].ObjectID,
					asst[0].AsstObjID,
				},
			},
		}
		models, err := ps.searchModels(modelFilter)
		if err != nil {
			ps.err = err
			return ps
		}

		// 处理模型自关联的情况
		if len(models) == 1 {
			ps.Attribute.Resources = []meta.ResourceAttribute{
				{
					Basic: meta.Basic{
						Type:       meta.ModelInstance,
						Action:     meta.Update,
						InstanceID: gjson.GetBytes(ps.RequestCtx.Body, common.BKInstIDField).Int(),
					},
					Layers:     []meta.Item{{Type: meta.Model, InstanceID: models[0].ID}},
					BusinessID: bizID,
				},
				{
					Basic: meta.Basic{
						Type:       meta.ModelInstance,
						Action:     meta.Update,
						InstanceID: gjson.GetBytes(ps.RequestCtx.Body, common.BKAsstInstIDField).Int(),
					},
					Layers:     []meta.Item{{Type: meta.Model, InstanceID: models[0].ID}},
					BusinessID: bizID,
				},
			}
			return ps
		}

		for _, model := range models {
			var instID int64
			if model.ObjectID == asst[0].ObjectID {
				instID = gjson.GetBytes(ps.RequestCtx.Body, common.BKInstIDField).Int()
			} else {
				instID = gjson.GetBytes(ps.RequestCtx.Body, common.BKAsstInstIDField).Int()
			}

			ps.Attribute.Resources = append(ps.Attribute.Resources,
				meta.ResourceAttribute{
					Basic: meta.Basic{
						Type:       meta.ModelInstance,
						Action:     meta.Update,
						InstanceID: instID,
					},
					Layers:     []meta.Item{{Type: meta.Model, InstanceID: model.ID}},
					BusinessID: bizID,
				})
		}
		return ps
	}

	// delete object's instance association operation. for web
	if ps.hitRegexp(deleteObjectInstanceAssociationLatestRegexp, http.MethodDelete) {
		assoID, err := strconv.ParseInt(ps.RequestCtx.Elements[4], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("delete object instance association, but got invalid association id %s", ps.RequestCtx.Elements[4])
			return ps
		}

		bizID, err := metadata.BizIDFromMetadata(ps.RequestCtx.Metadata)
		if err != nil {
			ps.err = err
			return ps
		}
		asst, err := ps.getInstAssociation(mapstr.MapStr{common.BKFieldID: assoID})
		if err != nil {
			ps.err = err
			return ps
		}
		models, err := ps.searchModels(mapstr.MapStr{common.BKObjIDField: mapstr.MapStr{common.BKDBIN: []interface{}{
			asst.ObjectID,
			asst.AsstObjectID,
		}}})
		if err != nil {
			ps.err = err
			return ps
		}

		// 处理模型自关联的情况
		if len(models) == 1 {
			ps.Attribute.Resources = []meta.ResourceAttribute{
				{
					Basic: meta.Basic{
						Type:       meta.ModelInstance,
						Action:     meta.Update,
						InstanceID: asst.InstID,
					},
					Layers:     []meta.Item{{Type: meta.Model, InstanceID: models[0].ID}},
					BusinessID: bizID,
				},
				{
					Basic: meta.Basic{
						Type:       meta.ModelInstance,
						Action:     meta.Update,
						InstanceID: asst.AsstInstID,
					},
					Layers:     []meta.Item{{Type: meta.Model, InstanceID: models[0].ID}},
					BusinessID: bizID,
				},
			}
			return ps
		}

		for _, model := range models {
			var instID int64
			if model.ObjectID == asst.ObjectID {
				instID = asst.InstID
			} else {
				instID = asst.AsstInstID
			}

			ps.Attribute.Resources = append(ps.Attribute.Resources,
				meta.ResourceAttribute{
					Basic: meta.Basic{
						Type:       meta.ModelInstance,
						Action:     meta.Update,
						InstanceID: instID,
					},
					Layers:     []meta.Item{{Type: meta.Model, InstanceID: model.ID}},
					BusinessID: bizID,
				})
		}

		return ps
	}

	// find object instance's association operation.
	if ps.hitRegexp(findObjectInstanceTopologyUILatestRegexp, http.MethodPost) {
		bizID, err := metadata.BizIDFromMetadata(ps.RequestCtx.Metadata)
		if err != nil {
			ps.err = err
			return ps
		}
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:   meta.ModelInstanceAssociation,
					Action: meta.FindMany,
				},
			},
		}
		return ps
	}

	// find object instance's association object instance info operation.
	if ps.hitRegexp(findInstAssociationObjInstInfoLatestRegexp, http.MethodPost) {
		bizID, err := metadata.BizIDFromMetadata(ps.RequestCtx.Metadata)
		if err != nil {
			ps.err = err
			return ps
		}
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:   meta.ModelInstanceAssociation,
					Action: meta.FindMany,
				},
			},
		}
		return ps
	}

	return ps
}

var (
	createObjectInstanceLatestRegexp          = regexp.MustCompile(`^/api/v3/create/instance/object/[^\s/]+/?$`)
	findObjectInstanceAssociationLatestRegexp = regexp.MustCompile(`^/api/v3/find/instassociation/object/[^\s/]+/?$`)
	updateObjectInstanceLatestRegexp          = regexp.MustCompile(`^/api/v3/update/instance/object/[^\s/]+/inst/[0-9]+/?$`)
	updateObjectInstanceBatchLatestRegexp     = regexp.MustCompile(`^/api/v3/updatemany/instance/object/[^\s/]+/?$`)
	deleteObjectInstanceBatchLatestRegexp     = regexp.MustCompile(`^/api/v3/deletemany/instance/object/[^\s/]+/?$`)
	deleteObjectInstanceLatestRegexp          = regexp.MustCompile(`^/api/v3/delete/instance/object/[^\s/]+/inst/[0-9]+/?$`)
	// TODO remove it
	findObjectInstanceSubTopologyLatestRegexp = regexp.MustCompile(`^/api/v3/find/insttopo/object/[^\s/]+/inst/[0-9]+/?$`)
	findObjectInstanceTopologyLatestRegexp    = regexp.MustCompile(`^/api/v3/find/instassttopo/object/[^\s/]+/inst/[0-9]+/?$`)
	findObjectInstancesLatestRegexp           = regexp.MustCompile(`^/api/v3/find/instance/object/[^\s/]+/?$`)
)

func (ps *parseStream) objectInstanceLatest() *parseStream {
	if ps.shouldReturn() {
		return ps
	}

	// create instance operation
	if ps.hitRegexp(createObjectInstanceLatestRegexp, http.MethodPost) {
		bizID, err := metadata.BizIDFromMetadata(ps.RequestCtx.Metadata)
		if err != nil {
			ps.err = err
			return ps
		}

		filter := mapstr.MapStr{
			common.BKObjIDField: ps.RequestCtx.Elements[5],
		}
		model, err := ps.getOneModel(filter)
		if err != nil {
			ps.err = err
			return ps
		}

		var modelType = meta.ModelInstance
		isMainline, err := ps.isMainlineModel(model.ObjectID)
		if err != nil {
			ps.err = err
			return ps
		}
		if isMainline {
			if bizID == 0 {
				ps.err = errors.New("create mainline instance must have metadata with biz id")
				return ps
			}
			modelType = meta.MainlineInstance
		}

		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:   modelType,
					Action: meta.Create,
				},
				Layers: []meta.Item{{Type: meta.Model, InstanceID: model.ID}},
			},
		}
		return ps
	}

	// search instance association
	if ps.hitRegexp(findObjectInstanceAssociationLatestRegexp, http.MethodPost) {
		if len(ps.RequestCtx.Elements) != 6 {
			ps.err = errors.New("search instance association, but got invalid url")
			return ps
		}
		objectID := ps.RequestCtx.Elements[5]
		filter := mapstr.MapStr{
			common.BKObjIDField: objectID,
		}
		model, err := ps.getOneModel(filter)
		if err != nil {
			ps.err = err
			return ps
		}

		var modelType = meta.ModelInstance
		isMainline, err := ps.isMainlineModel(objectID)
		if err != nil {
			ps.err = err
			return ps
		}
		if isMainline {
			modelType = meta.MainlineInstance
		}

		bizID, err := metadata.BizIDFromMetadata(ps.RequestCtx.Metadata)
		if err != nil {
			ps.err = err
			return ps
		}

		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:   modelType,
					Action: meta.Find,
				},
				Layers: []meta.Item{{Type: meta.Model, InstanceID: model.ID}},
			},
		}

		return ps
	}

	// update instance operation
	if ps.hitRegexp(updateObjectInstanceLatestRegexp, http.MethodPut) {
		if len(ps.RequestCtx.Elements) != 8 {
			ps.err = errors.New("update object instance, but got invalid url")
			return ps
		}

		bizID, err := metadata.BizIDFromMetadata(ps.RequestCtx.Metadata)
		if err != nil {
			ps.err = err
			return ps
		}

		instID, err := strconv.ParseInt(ps.RequestCtx.Elements[7], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("update object instance, but got invalid instance id %s", ps.RequestCtx.Elements[5])
			return ps
		}

		objectID := ps.RequestCtx.Elements[5]
		filter := mapstr.MapStr{
			common.BKObjIDField: objectID,
		}
		model, err := ps.getOneModel(filter)
		if err != nil {
			ps.err = err
			return ps
		}

		var modelType = meta.ModelInstance
		isMainline, err := ps.isMainlineModel(objectID)
		if err != nil {
			ps.err = err
			return ps
		}
		if isMainline {
			modelType = meta.MainlineInstance
		}

		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:       modelType,
					Action:     meta.Update,
					InstanceID: instID,
				},
				Layers: []meta.Item{{Type: meta.Model, InstanceID: model.ID}},
			},
		}
		return ps
	}

	// batch update instance operation
	if ps.hitRegexp(updateObjectInstanceBatchLatestRegexp, http.MethodPut) {
		if len(ps.RequestCtx.Elements) != 6 {
			ps.err = errors.New("update object instance batch, but got invalid url")
			return ps
		}

		bizID, err := metadata.BizIDFromMetadata(ps.RequestCtx.Metadata)
		if err != nil {
			ps.err = err
			return ps
		}

		objectID := ps.RequestCtx.Elements[5]
		filter := mapstr.MapStr{
			common.BKObjIDField: objectID,
		}
		model, err := ps.getOneModel(filter)
		if err != nil {
			ps.err = err
			return ps
		}

		ids := make([]int64, 0)
		gjson.GetBytes(ps.RequestCtx.Body, "update.#.inst_id").ForEach(
			func(key, value gjson.Result) bool {
				ids = append(ids, value.Int())
				return true
			})

		for _, id := range ids {
			ps.Attribute.Resources = append(ps.Attribute.Resources, meta.ResourceAttribute{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:       meta.ModelInstance,
					Action:     meta.UpdateMany,
					InstanceID: id,
				},
				Layers: []meta.Item{{Type: meta.Model, InstanceID: model.ID}},
			})
		}

		return ps
	}

	// batch delete instance operation
	if ps.hitRegexp(deleteObjectInstanceBatchLatestRegexp, http.MethodDelete) {
		if len(ps.RequestCtx.Elements) != 6 {
			ps.err = errors.New("delete object instance batch, but got invalid url")
			return ps
		}

		bizID, err := metadata.BizIDFromMetadata(ps.RequestCtx.Metadata)
		if err != nil {
			ps.err = err
			return ps
		}

		objectID := ps.RequestCtx.Elements[5]
		filter := mapstr.MapStr{
			common.BKObjIDField: objectID,
		}
		model, err := ps.getOneModel(filter)
		if err != nil {
			ps.err = err
			return ps
		}

		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:   meta.ModelInstance,
					Action: meta.DeleteMany,
				},
				Layers: []meta.Item{{Type: meta.Model, InstanceID: model.ID}},
			},
		}

		return ps
	}

	// delete instance operation.
	if ps.hitRegexp(deleteObjectInstanceLatestRegexp, http.MethodDelete) {
		if len(ps.RequestCtx.Elements) != 8 {
			ps.err = errors.New("delete object instance, but got invalid url")
			return ps
		}

		bizID, err := metadata.BizIDFromMetadata(ps.RequestCtx.Metadata)
		if err != nil {
			ps.err = err
			return ps
		}
		e7 := ps.RequestCtx.Elements[7]
		instID, err := strconv.ParseInt(e7, 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("delete object instance, but got invalid instance id %s", e7)
			return ps
		}

		filter := mapstr.MapStr{
			common.BKObjIDField: ps.RequestCtx.Elements[5],
		}
		model, err := ps.getOneModel(filter)
		if err != nil {
			ps.err = err
			return ps
		}

		var modelType = meta.ModelInstance
		isMainline, err := ps.isMainlineModel(model.ObjectID)
		if err != nil {
			ps.err = err
			return ps
		}
		if isMainline {
			// special logic for mainline object's instance authorization.
			if bizID == 0 {
				ps.err = errors.New("delete mainline instance must have metadata with biz id")
				return ps
			}
			modelType = meta.MainlineInstance
		}

		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:       modelType,
					Action:     meta.Delete,
					InstanceID: instID,
				},
				Layers: []meta.Item{{Type: meta.Model, InstanceID: model.ID}},
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

		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:       meta.ModelInstanceTopology,
					Action:     meta.Find,
					InstanceID: instID,
				},
				Layers: []meta.Item{
					{
						Type: meta.Model,
						Name: ps.RequestCtx.Elements[5],
					},
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

		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:   meta.ModelInstanceTopology,
					Action: meta.Find,
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

		bizID, err := metadata.BizIDFromMetadata(ps.RequestCtx.Metadata)
		if err != nil {
			ps.err = fmt.Errorf("parse bizID from metadata failed, err: %s", err.Error())
			return ps
		}
		filter := mapstr.MapStr{
			common.BKObjIDField: ps.RequestCtx.Elements[5],
		}
		model, err := ps.getOneModel(filter)
		if err != nil {
			ps.err = err
			return ps
		}
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:   meta.ModelInstance,
					Action: meta.FindMany,
				},
				Layers: []meta.Item{{Type: meta.Model, InstanceID: model.ID}},
			},
		}
		return ps
	}

	return ps
}

const (
	createObjectLatestPattern       = "/api/v3/create/object"
	findObjectsLatestPattern        = "/api/v3/find/object"
	findObjectTopologyLatestPattern = "/api/v3/find/objecttopology"
)

var (
	deleteObjectLatestRegexp = regexp.MustCompile(`^/api/v3/delete/object/[0-9]+/?$`)
	updateObjectLatestRegexp = regexp.MustCompile(`^/api/v3/update/object/[0-9]+/?$`)

	// TODO remove it
	// 获取模型拓扑图及位置信息-Web
	findObjectTopologyGraphicLatestRegexp = regexp.MustCompile(`^/api/v3/find/objecttopo/scope_type/[^\s/]+/scope_id/[^\s/]+/?$`)
	// 设置模型拓扑图及位置信息-Web
	updateObjectTopologyGraphicLatestRegexp = regexp.MustCompile(`^/api/v3/update/objecttopo/scope_type/[^\s/]+/scope_id/[^\s/]+/?$`)
)

func (ps *parseStream) objectLatest() *parseStream {
	if ps.shouldReturn() {
		return ps
	}

	// create common object operation.
	if ps.hitPattern(createObjectLatestPattern, http.MethodPost) {
		bizID, err := metadata.BizIDFromMetadata(ps.RequestCtx.Metadata)
		if err != nil {
			blog.Warnf("create object, but get business id in metadata failed, err: %v", err)
		}
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:   meta.Model,
					Action: meta.Create,
				},
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

		id, err := strconv.ParseInt(ps.RequestCtx.Elements[4], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("delete object, but got invalid object's id %s", ps.RequestCtx.Elements[3])
			return ps
		}

		filter := map[string]interface{}{
			common.BKFieldID: id,
		}
		model, err := ps.getOneModel(filter)
		if err != nil {
			ps.err = fmt.Errorf("delete object, but model by objectID failed, err: %v", err)
			return ps
		}
		bizID, err := metadata.BizIDFromMetadata(model.Metadata)
		if err != nil {
			blog.Warnf("delete object, but get business id in metadata failed, err: %v", err)
		}
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:       meta.Model,
					Action:     meta.Delete,
					InstanceID: id,
				},
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

		id, err := strconv.ParseInt(ps.RequestCtx.Elements[4], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("update object, but got invalid object's id %s", ps.RequestCtx.Elements[4])
			return ps
		}
		filter := map[string]interface{}{
			common.BKFieldID: id,
		}
		model, err := ps.getOneModel(filter)
		if err != nil {
			ps.err = fmt.Errorf("delete object, but model by objectID failed, err: %v", err)
			return ps
		}
		bizID, err := metadata.BizIDFromMetadata(model.Metadata)
		if err != nil {
			ps.err = fmt.Errorf("update object, but get business id in metadata failed, err: %v", err)
			return ps
		}
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:       meta.Model,
					Action:     meta.Update,
					InstanceID: id,
				},
			},
		}
		return ps
	}

	// get object operation.
	if ps.hitPattern(findObjectsLatestPattern, http.MethodPost) {
		bizID, err := metadata.BizIDFromMetadata(ps.RequestCtx.Metadata)
		if err != nil {
			blog.Warnf("find object, but get business id in metadata failed, err: %v", err)
		}
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:   meta.Model,
					Action: meta.FindMany,
				},
			},
		}
		return ps
	}

	// find object's topology operation.
	if ps.hitPattern(findObjectTopologyLatestPattern, http.MethodPut) {
		bizID, err := metadata.BizIDFromMetadata(ps.RequestCtx.Metadata)
		if err != nil {
			blog.Warnf("find object, but get business id in metadata failed, err: %v", err)
		}
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:   meta.ModelTopology,
					Action: meta.Find,
				},
			},
		}
		return ps
	}

	// find object's topology graphic operation.
	if ps.hitRegexp(findObjectTopologyGraphicLatestRegexp, http.MethodPost) {
		bizID, err := metadata.BizIDFromMetadata(ps.RequestCtx.Metadata)
		if err != nil {
			blog.Warnf("find object topology graphic, but get business id in metadata failed, err: %v", err)
		}
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type: meta.ModelTopology,
					// Action: meta.Find,
					Action: meta.SkipAction,
				},
			},
		}
		return ps
	}

	// update object's topology graphic operation.
	// TODO: confirm if bizID is needed.
	if ps.hitRegexp(updateObjectTopologyGraphicLatestRegexp, http.MethodPost) {

		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type: meta.ModelTopology,
					// Action: meta.Update,
					Action: meta.SkipAction,
				},
			},
		}
		return ps
	}

	return ps
}

const (
	createObjectClassificationLatestPattern   = "/api/v3/create/objectclassification"
	findObjectClassificationListLatestPattern = "/api/v3/find/objectclassification"
	// 查找模型分组及分组下的模型列表
	findObjectsBelongsToClassificationLatestPattern = `/api/v3/find/classificationobject`
)

var (
	deleteObjectClassificationLatestRegexp = regexp.MustCompile("^/api/v3/delete/objectclassification/[0-9]+/?$")
	updateObjectClassificationLatestRegexp = regexp.MustCompile("^/api/v3/update/objectclassification/[0-9]+/?$")
)

func (ps *parseStream) ObjectClassificationLatest() *parseStream {
	if ps.shouldReturn() {
		return ps
	}

	// create object's classification operation.
	if ps.hitPattern(createObjectClassificationLatestPattern, http.MethodPost) {
		bizID, err := metadata.BizIDFromMetadata(ps.RequestCtx.Metadata)
		if err != nil {
			blog.Warnf("get business id in metadata failed, err: %v", err)
		}
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:   meta.ModelClassification,
					Action: meta.Create,
				},
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

		filter := map[string]interface{}{
			common.BKFieldID: classID,
		}
		classification, err := ps.getOneClassification(filter)
		if err != nil {
			ps.err = fmt.Errorf("delete object classification, but get by id failed")
			return ps
		}
		bizID, err := metadata.BizIDFromMetadata(classification.Metadata)
		if err != nil {
			ps.err = fmt.Errorf("delete object classification, get business id in metadata failed, err: %v", err)
			return ps
		}
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:       meta.ModelClassification,
					Action:     meta.Delete,
					InstanceID: classID,
				},
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
		filter := map[string]interface{}{
			common.BKFieldID: classID,
		}
		classification, err := ps.getOneClassification(filter)
		if err != nil {
			ps.err = fmt.Errorf("delete object classification, but get by id failed")
			return ps
		}
		bizID, err := metadata.BizIDFromMetadata(classification.Metadata)
		if err != nil {
			ps.err = fmt.Errorf("get business id in metadata failed, err: %v", err)
			return ps
		}
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:       meta.ModelClassification,
					Action:     meta.Update,
					InstanceID: classID,
				},
			},
		}
		return ps
	}

	// find object's classification list operation.
	if ps.hitPattern(findObjectClassificationListLatestPattern, http.MethodPost) {
		bizID, err := metadata.BizIDFromMetadata(ps.RequestCtx.Metadata)
		if err != nil {
			blog.Warnf("get business id in metadata failed, err: %v", err)
		}
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:   meta.ModelClassification,
					Action: meta.FindMany,
				},
			},
		}
		return ps
	}
	// find all the objects belongs to a classification
	if ps.hitPattern(findObjectsBelongsToClassificationLatestPattern, http.MethodPost) {
		bizID, err := metadata.BizIDFromMetadata(ps.RequestCtx.Metadata)
		if err != nil {
			blog.Warnf("get business id in metadata failed, err: %v", err)
		}
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:   meta.Model,
					Action: meta.FindMany,
				},
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
	findObjectAttributeGroupLatestRegexp   = regexp.MustCompile(`^/api/v3/find/objectattgroup/object/[^\s/]+/?$`)
	deleteObjectAttributeGroupLatestRegexp = regexp.MustCompile(`^/api/v3/delete/objectattgroup/[0-9]+/?$`)
	// TODO remove it, interface implementation not found
	removeAttributeAwayFromGroupLatestRegexp = regexp.MustCompile(`^/api/v3/delete/objectattgroupasst/object/[^\s/]+/property/[^\s/]+/group/[^\s/]+/?$`)
)

func (ps *parseStream) objectAttributeGroupLatest() *parseStream {
	if ps.shouldReturn() {
		return ps
	}
	// create object's attribute group operation.
	if ps.hitPattern(createObjectAttributeGroupLatestPattern, http.MethodPost) {
		// 业务ID的解释
		// case  0: 创建公共的属性分组
		// case ~0: 创建业务私有的属性分组
		bizID, err := metadata.BizIDFromMetadata(ps.RequestCtx.Metadata)
		if err != nil {
			blog.Warnf("get business id in metadata failed, err: %v", err)
		}
		filter := mapstr.MapStr{
			common.BKObjIDField: gjson.GetBytes(ps.RequestCtx.Body, common.BKObjIDField).Value(),
		}
		model, err := ps.getOneModel(filter)
		if err != nil {
			ps.err = err
			return ps
		}

		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:   meta.ModelAttributeGroup,
					Action: meta.Create,
				},
				Layers: []meta.Item{{Type: meta.Model, InstanceID: model.ID}},
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

		model, err := ps.getOneModel(mapstr.MapStr{common.BKObjIDField: ps.RequestCtx.Elements[5]})
		if err != nil {
			ps.err = err
			return ps
		}

		// 业务ID的解释
		// case  0: 仅查询公共的属性分组
		// case ~0: 查询业务私有的属性分组 + 公用属性分组
		bizID, err := metadata.BizIDFromMetadata(ps.RequestCtx.Metadata)
		if err != nil {
			blog.Warnf("get business id in metadata failed, err: %v", err)
		}
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:   meta.ModelAttributeGroup,
					Action: meta.FindMany,
				},
				Layers: []meta.Item{{Type: meta.Model, InstanceID: model.ID}},
			},
		}
		return ps
	}

	// update object's attribute group operation.
	if ps.hitPattern(updateObjectAttributeGroupLatestPattern, http.MethodPut) {
		groups, err := ps.getAttributeGroup(gjson.GetBytes(ps.RequestCtx.Body, "condition").Value())
		if err != nil {
			ps.err = err
			return ps
		}

		for _, group := range groups {
			bizID, err := metadata.BizIDFromMetadata(group.Metadata)
			if err != nil {
				blog.Warnf("get business id in metadata failed, err: %v", err)
			}

			filter := mapstr.MapStr{
				common.BKObjIDField: group.ObjectID,
			}
			model, err := ps.getOneModel(filter)
			if err != nil {
				ps.err = err
				return ps
			}
			ps.Attribute.Resources = append(ps.Attribute.Resources,
				meta.ResourceAttribute{
					BusinessID: bizID,
					Basic: meta.Basic{
						Type:       meta.ModelAttributeGroup,
						Action:     meta.Update,
						InstanceID: group.ID,
					},
					Layers: []meta.Item{{Type: meta.Model, InstanceID: model.ID}},
				})
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

		groups, err := ps.getAttributeGroup(mapstr.MapStr{"id": groupID})
		if err != nil {
			ps.err = err
			return ps
		}

		bizID, err := metadata.BizIDFromMetadata(groups[0].Metadata)
		if err != nil {
			blog.Warnf("get business id in metadata failed, err: %v", err)
		}

		model, err := ps.getOneModel(mapstr.MapStr{common.BKObjIDField: groups[0].ObjectID})
		if err != nil {
			ps.err = err
			return ps
		}
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:       meta.ModelAttributeGroup,
					Action:     meta.Delete,
					InstanceID: groupID,
				},
				Layers: []meta.Item{{Type: meta.Model, InstanceID: model.ID}},
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

		bizID, err := metadata.BizIDFromMetadata(ps.RequestCtx.Metadata)
		if err != nil {
			blog.Warnf("get business id in metadata failed, err: %v", err)
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

const (
	createObjectAttributeLatestPattern = "/api/v3/create/objectattr"
	findObjectAttributeLatestPattern   = "/api/v3/find/objectattr"
)

var (
	deleteObjectAttributeLatestRegexp = regexp.MustCompile(`^/api/v3/delete/objectattr/[0-9]+/?$`)
	updateObjectAttributeLatestRegexp = regexp.MustCompile(`^/api/v3/update/objectattr/[0-9]+/?$`)
)

func (ps *parseStream) objectAttributeLatest() *parseStream {
	if ps.shouldReturn() {
		return ps
	}

	// create object's attribute operation.
	if ps.hitPattern(createObjectAttributeLatestPattern, http.MethodPost) {
		// 注意业务ID是否为0表示创建两种不同的属性
		// case 0: 创建公共属性，这种属性相比业务私有属性，所有业务都可见
		// case ~0: 创建业务私有属性，业务私有属性，其它业务不可见
		bizID, err := metadata.BizIDFromMetadata(ps.RequestCtx.Metadata)
		if err != nil {
			blog.Warnf("get business id in metadata failed, err: %v", err)
		}
		modelEn := gjson.GetBytes(ps.RequestCtx.Body, common.BKObjIDField).String()
		model, err := ps.getOneModel(mapstr.MapStr{common.BKObjIDField: modelEn})
		if err != nil {
			ps.err = err
			return ps
		}
		ps.Attribute.Resources = []meta.ResourceAttribute{
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

		attr, err := ps.getModelAttribute(mapstr.MapStr{common.BKFieldID: attrID})
		if err != nil {
			ps.err = fmt.Errorf("delete object attribute, but fetch attribute by %v failed %v", mapstr.MapStr{common.BKFieldID: attrID}, err)
			return ps
		}

		model, err := ps.getOneModel(mapstr.MapStr{common.BKObjIDField: attr[0].ObjectID})
		if err != nil {
			ps.err = err
			return ps
		}

		// 对属性操作的鉴权，依赖于属性是公有属性，还是业务私有属性
		bizID, err := metadata.BizIDFromMetadata(attr[0].Metadata)
		if err != nil {
			blog.Warnf("get business id in metadata failed, err: %v", err)
		}
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:       meta.ModelAttribute,
					Action:     meta.Delete,
					InstanceID: attrID,
				},
				Layers: []meta.Item{{Type: meta.Model, InstanceID: model.ID}},
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

		attr, err := ps.getModelAttribute(mapstr.MapStr{common.BKFieldID: attrID})
		if err != nil {
			ps.err = fmt.Errorf("delete object attribute, but fetch attribute by %v failed %v", mapstr.MapStr{common.BKFieldID: attrID}, err)
			return ps
		}
		model, err := ps.getOneModel(mapstr.MapStr{common.BKObjIDField: attr[0].ObjectID})
		if err != nil {
			ps.err = err
			return ps
		}

		// 对属性操作的鉴权，依赖于属性是公有属性，还是业务私有属性
		bizID, err := metadata.BizIDFromMetadata(attr[0].Metadata)
		if err != nil {
			blog.Warnf("get business id in metadata failed, err: %v", err)
		}
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:       meta.ModelAttribute,
					Action:     meta.Update,
					InstanceID: attrID,
				},
				Layers: []meta.Item{{Type: meta.Model, InstanceID: model.ID}},
			},
		}
		return ps
	}

	// get object's attribute operation.
	if ps.hitPattern(findObjectAttributeLatestPattern, http.MethodPost) {
		// 注意：业务ID是否为0表示两种不同的操作
		// case 0: 读取模型的公有属性
		// case ~0: 读取业务私有属性+公有属性
		bizID, err := metadata.BizIDFromMetadata(ps.RequestCtx.Metadata)
		if err != nil {
			blog.V(5).Infof("get business id in metadata failed, err: %v", err)
		}

		modelCond := gjson.GetBytes(ps.RequestCtx.Body, common.BKObjIDField).Value()
		models, err := ps.searchModels(mapstr.MapStr{common.BKObjIDField: modelCond})
		if err != nil {
			ps.err = err
			return ps
		}
		for _, model := range models {

			ps.Attribute.Resources = append(ps.Attribute.Resources,
				meta.ResourceAttribute{
					BusinessID: bizID,
					Basic: meta.Basic{
						Type:   meta.ModelAttribute,
						Action: meta.FindMany,
					},
					Layers: []meta.Item{{Type: meta.Model, InstanceID: model.ID}},
				})
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
	deleteMainlineObjectLatestRegexp                       = regexp.MustCompile(`^/api/v3/delete/topomodelmainline/object/[^\s/]+/?$`)
	findBusinessInstanceTopologyLatestRegexp               = regexp.MustCompile(`^/api/v3/find/topoinst/biz/[0-9]+/?$`)
	findBusinessInstanceTopologyPathRegexp                 = regexp.MustCompile(`^/api/v3/find/topopath/biz/[0-9]+/?$`)
	findBusinessInstanceTopologyWithStatisticsLatestRegexp = regexp.MustCompile(`^/api/v3/find/topoinst_with_statistics/biz/[0-9]+/?$`)
	// TODO remove it, interface implementation not found
	findMainlineSubInstanceTopoLatestRegexp = regexp.MustCompile(`^/api/v3/topoinstchild/object/[^\s/]+/biz/[0-9]+/inst/[0-9]+/?$`)
	// TODO remove it, interface implementation not found
	findMainlineIdleFaultModuleLatestRegexp = regexp.MustCompile(`^/api/v3/find/topointernal/biz/[0-9]+/?$`)
)

func (ps *parseStream) mainlineLatest() *parseStream {
	if ps.shouldReturn() {
		return ps
	}

	// create mainline object operation.
	if ps.hitPattern(createMainlineObjectLatestPattern, http.MethodPost) {
		bizID, err := metadata.BizIDFromMetadata(ps.RequestCtx.Metadata)
		if err != nil {
			blog.Warnf("get business id in metadata failed, err: %v", err)
		}
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:   meta.MainlineModel,
					Action: meta.Create,
				},
			},
		}
		return ps
	}

	// delete mainline object operation
	if ps.hitRegexp(deleteMainlineObjectLatestRegexp, http.MethodDelete) {
		bizID, err := metadata.BizIDFromMetadata(ps.RequestCtx.Metadata)
		if err != nil {
			blog.Warnf("get business id in metadata failed, err: %v", err)
		}
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:   meta.MainlineModel,
					Action: meta.Delete,
				},
			},
		}

		return ps
	}

	// get mainline object operation
	if ps.hitPattern(findMainlineObjectTopoLatestPattern, http.MethodPost) {
		bizID, err := metadata.BizIDFromMetadata(ps.RequestCtx.Metadata)
		if err != nil {
			blog.Warnf("get business id in metadata failed, err: %v", err)
		}
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type: meta.MainlineModelTopology,
					// Action: meta.Find,
					Action: meta.SkipAction,
				},
			},
		}

		return ps
	}

	// find mainline object instance's sub-instance topology operation.
	if ps.hitRegexp(findMainlineSubInstanceTopoLatestRegexp, http.MethodGet) {
		if len(ps.RequestCtx.Elements) != 9 {
			ps.err = errors.New("find mainline object's sub instance topology, but got invalid url")
			return ps
		}

		bizID, err := strconv.ParseInt(ps.RequestCtx.Elements[6], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("find mainline object's sub instance topology, but got invalid business id %s", ps.RequestCtx.Elements[6])
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
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:   meta.MainlineModel,
					Action: meta.Find,
				},
			},
		}

		return ps
	}

	// find business instance topology operation.
	// also is find mainline instance topology operation.
	if ps.hitRegexp(findBusinessInstanceTopologyLatestRegexp, http.MethodPost) ||
		ps.hitRegexp(findBusinessInstanceTopologyPathRegexp, http.MethodPost) ||
		ps.hitRegexp(findBusinessInstanceTopologyWithStatisticsLatestRegexp, http.MethodPost) {
		if len(ps.RequestCtx.Elements) != 6 {
			ps.err = errors.New("find business instance topology, but got invalid url")
			return ps
		}

		bizID, err := metadata.BizIDFromMetadata(ps.RequestCtx.Metadata)
		if err != nil {
			blog.Warnf("find business instance, but get business id in metadata failed, err: %v", err)
		}
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:   meta.ModelInstanceTopology,
					Action: meta.Find,
				},
			},
		}
		return ps
	}

	return ps
}
