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
	"configcenter/pkg/filter"
	filtertools "configcenter/pkg/tools/filter"
	"configcenter/src/ac/meta"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
)

// ListObjFieldTmplRel list field template and object relations.
func (s *service) ListObjFieldTmplRel(cts *rest.Contexts) {
	opt := new(metadata.ListObjFieldTmplRelOption)
	if err := cts.DecodeInto(opt); err != nil {
		cts.RespAutoError(err)
		return
	}

	if rawErr := opt.Validate(); rawErr.ErrCode != 0 {
		cts.RespAutoError(rawErr.ToCCError(cts.Kit.CCError))
		return
	}

	authResp, authorized, err := s.hasObjOrTmplAuth(cts.Kit, opt.ObjectIDs, opt.TemplateIDs, any)
	if err != nil {
		cts.RespAutoError(err)
		return
	}
	if !authorized {
		cts.RespNoAuth(authResp)
		return
	}

	// list field templates and object relations
	var relFilter *filter.Expression
	if len(opt.TemplateIDs) > 0 {
		relFilter = filtertools.GenAtomFilter(common.BKTemplateID, filter.In, opt.TemplateIDs)
	}

	if len(opt.ObjectIDs) > 0 {
		var err error
		relFilter, err = filtertools.And(relFilter, filtertools.GenAtomFilter(common.ObjectIDField, filter.In,
			opt.ObjectIDs))
		if err != nil {
			cts.RespAutoError(cts.Kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, "filter"))
			return
		}
	}

	listOpt := &metadata.CommonQueryOption{
		CommonFilterOption: metadata.CommonFilterOption{Filter: relFilter},
		Page:               metadata.BasePage{Limit: common.BKNoLimit},
	}
	res, err := s.clientSet.CoreService().FieldTemplate().ListObjFieldTmplRel(cts.Kit.Ctx, cts.Kit.Header, listOpt)
	if err != nil {
		blog.Errorf("list field templates failed, err: %v, opt: %+v, rid: %s", err, opt, cts.Kit.Rid)
		cts.RespAutoError(err)
		return
	}

	cts.RespEntity(res)
}

// priority 当没权限时，决定优先返回是因为什么原因而没权限
type priority string

const (
	object   priority = "object"
	template priority = "template"
	any      priority = "any"
)

func (s *service) hasObjOrTmplAuth(kit *rest.Kit, tmplIDs []int64, objectIDs []int64, priorityVal priority) (
	*metadata.BaseResp, bool, error) {

	var tmplAuthResp *metadata.BaseResp
	var tmplAuthorized bool
	if len(tmplIDs) != 0 {
		// check if user has the permission of the field template
		resAttr := make([]meta.ResourceAttribute, len(tmplIDs))
		for idx, tmplID := range tmplIDs {
			resAttr[idx] = meta.ResourceAttribute{Basic: meta.Basic{Type: meta.FieldTemplate, Action: meta.Find,
				InstanceID: tmplID}}
		}

		tmplAuthResp, tmplAuthorized = s.auth.Authorize(kit, resAttr...)
	}

	var objAuthResp *metadata.BaseResp
	var objAuthorized bool
	if len(objectIDs) != 0 {
		var err error
		objAuthResp, objAuthorized, err = s.auth.HasFindModelAuthUseID(kit, objectIDs)
		if err != nil {
			return nil, false, err
		}
	}

	if tmplAuthorized || objAuthorized {
		return nil, true, nil
	}

	switch priorityVal {
	case template:
		return tmplAuthResp, false, nil
	case object:
		return objAuthResp, false, nil
	case any:
		if tmplAuthResp != nil {
			return tmplAuthResp, false, nil
		}
		if objAuthResp != nil {
			return objAuthResp, false, nil
		}
	}

	return nil, false, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, "priority")
}

// ListFieldTmplByObj list field template by related object id.
func (s *service) ListFieldTmplByObj(cts *rest.Contexts) {
	opt := new(metadata.ListFieldTmplByObjOption)
	if err := cts.DecodeInto(opt); err != nil {
		cts.RespAutoError(err)
		return
	}

	if rawErr := opt.Validate(); rawErr.ErrCode != 0 {
		cts.RespAutoError(rawErr.ToCCError(cts.Kit.CCError))
		return
	}

	authResp, authorized, err := s.hasObjOrTmplAuth(cts.Kit, nil, []int64{opt.ObjectID}, object)
	if err != nil {
		cts.RespAutoError(err)
		return
	}
	if !authorized {
		cts.RespNoAuth(authResp)
		return
	}

	// get field templates ids by object id
	relOpt := &metadata.CommonQueryOption{
		CommonFilterOption: metadata.CommonFilterOption{Filter: filtertools.GenAtomFilter(common.ObjectIDField,
			filter.Equal, opt.ObjectID)},
		Page: metadata.BasePage{Limit: common.BKNoLimit},
	}
	relRes, err := s.clientSet.CoreService().FieldTemplate().ListObjFieldTmplRel(cts.Kit.Ctx, cts.Kit.Header, relOpt)
	if err != nil {
		blog.Errorf("list field templates failed, err: %v, opt: %+v, rid: %s", err, opt, cts.Kit.Rid)
		cts.RespAutoError(err)
		return
	}

	templateIDs := make([]int64, len(relRes.Info))
	for i, relation := range relRes.Info {
		templateIDs[i] = relation.TemplateID
	}

	if len(templateIDs) == 0 {
		cts.RespEntity(&metadata.FieldTemplateInfo{Info: make([]metadata.FieldTemplate, 0)})
		return
	}

	// list filed template by ids
	listOpt := &metadata.CommonQueryOption{
		CommonFilterOption: metadata.CommonFilterOption{Filter: filtertools.GenAtomFilter(common.BKFieldID,
			filter.In, templateIDs)},
		Page: metadata.BasePage{Limit: common.BKNoLimit},
	}
	res, err := s.clientSet.CoreService().FieldTemplate().ListFieldTemplate(cts.Kit.Ctx, cts.Kit.Header, listOpt)
	if err != nil {
		blog.Errorf("list field templates failed, err: %v, opt: %+v, rid: %s", err, opt, cts.Kit.Rid)
		cts.RespAutoError(err)
		return
	}

	cts.RespEntity(res)
}

// ListObjByFieldTmpl list object by field template.
func (s *service) ListObjByFieldTmpl(cts *rest.Contexts) {
	opt := new(metadata.ListObjByFieldTmplOption)
	if err := cts.DecodeInto(opt); err != nil {
		cts.RespAutoError(err)
		return
	}

	if rawErr := opt.Validate(); rawErr.ErrCode != 0 {
		cts.RespAutoError(rawErr.ToCCError(cts.Kit.CCError))
		return
	}

	authResp, authorized, err := s.hasObjOrTmplAuth(cts.Kit, []int64{opt.TemplateID}, nil, template)
	if err != nil {
		cts.RespAutoError(err)
		return
	}
	if !authorized {
		cts.RespNoAuth(authResp)
		return
	}

	// if object detail is needed later, add object auth check

	// get object ids by field template id
	relOpt := &metadata.CommonQueryOption{
		CommonFilterOption: metadata.CommonFilterOption{Filter: filtertools.GenAtomFilter(common.BKTemplateID,
			filter.Equal, opt.TemplateID)},
		Page: metadata.BasePage{Limit: common.BKNoLimit},
	}
	relRes, err := s.clientSet.CoreService().FieldTemplate().ListObjFieldTmplRel(cts.Kit.Ctx, cts.Kit.Header, relOpt)
	if err != nil {
		blog.Errorf("list field templates failed, err: %v, opt: %+v, rid: %s", err, opt, cts.Kit.Rid)
		cts.RespAutoError(err)
		return
	}

	objectIDs := make([]int64, len(relRes.Info))
	for i, relation := range relRes.Info {
		objectIDs[i] = relation.ObjectID
	}

	if len(objectIDs) == 0 {
		cts.RespEntity(&metadata.FieldTemplateInfo{Info: make([]metadata.FieldTemplate, 0)})
		return
	}

	expr, rawErr := filtertools.And(filtertools.GenAtomFilter(common.BKFieldID, filter.In, objectIDs), opt.Filter)
	if rawErr != nil {
		blog.Errorf("merge field template filter failed, err: %v, opt: %+v, rid: %s", err, opt, cts.Kit.Rid)
		cts.RespAutoError(rawErr)
		return
	}

	// list object by ids
	listOpt := &metadata.CommonQueryOption{
		Fields:             opt.Fields,
		Page:               opt.Page,
		CommonFilterOption: metadata.CommonFilterOption{Filter: expr},
	}
	res, objErr := s.clientSet.CoreService().Model().ListModel(cts.Kit.Ctx, cts.Kit.Header, listOpt)
	if objErr != nil {
		blog.Errorf("list objects failed, err: %v, opt: %+v, rid: %s", objErr, opt, cts.Kit.Rid)
		cts.RespAutoError(objErr)
		return
	}

	cts.RespEntity(res)
}

func (s *service) preCheckExecuteTaskAndGetObjID(kit *rest.Kit, syncOption *metadata.SyncObjectTask) (string, error) {

	// 1、determine whether the field template exists
	templateCond := []map[string]interface{}{
		{
			common.BKFieldID: syncOption.TemplateID,
		},
	}

	counts, err := s.clientSet.CoreService().Count().GetCountByFilter(kit.Ctx, kit.Header,
		common.BKTableNameFieldTemplate, templateCond)
	if err != nil {
		blog.Error("get field template num failed, cond: %+v, err: %v, rid: %s", templateCond, err, kit.Rid)
		return "", kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKTemplateID)
	}

	if len(counts) != 1 || int(counts[0]) != 1 {
		blog.Errorf("list field template num error, cond: %+v, rid: %s", syncOption, kit.Rid)
		return "", kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKTemplateID)
	}

	// 2、determine whether the model exists
	objCond := &metadata.QueryCondition{
		Fields: []string{metadata.ModelFieldIsPaused, common.BKObjIDField},
		Page: metadata.BasePage{
			Limit: common.BKNoLimit,
		},
		Condition: mapstr.MapStr{common.BKFieldID: syncOption.ObjectID},
	}

	result, objErr := s.clientSet.CoreService().Model().ReadModel(kit.Ctx, kit.Header, objCond)
	if objErr != nil {
		blog.Errorf("list objects failed, err: %v, cond: %+v, rid: %s", objErr, syncOption, kit.Rid)
		return "", kit.CCError.CCError(common.CCErrCommParamsInvalid)
	}

	if result.Count != 1 {
		blog.Errorf("objects num error,  count: %d, opt: %+v, rid: %s", result.Count, syncOption, kit.Rid)
		return "", kit.CCError.CCError(common.CCErrCommParamsInvalid)
	}

	// 3、determine whether the state of the model is enabled
	if result.Info[0].IsPaused {
		blog.Errorf("object status is stop, unable to execute sync task, opt: %+v, rid: %s", syncOption, kit.Rid)
		return "", kit.CCError.CCError(common.CCErrorTopoModelStopped)
	}

	// 4、determine whether the binding relationship between the model and the template exists
	relCond := []map[string]interface{}{
		{
			common.ObjectIDField: syncOption.ObjectID,
			common.BKTemplateID:  syncOption.TemplateID,
		},
	}

	counts, err = s.clientSet.CoreService().Count().GetCountByFilter(kit.Ctx, kit.Header,
		common.BKTableNameObjFieldTemplateRelation, relCond)
	if err != nil {
		blog.Error("get invalid relation failed, cond: %+v, err: %v, rid: %s", relCond, err, kit.Rid)
		return "", kit.CCError.CCError(common.CCErrCommParamsInvalid)
	}

	if len(counts) != 1 || int(counts[0]) != 1 {
		blog.Error("get invalid relation count error, cond: %+v, rid: %s", relCond, kit.Rid)
		return "", kit.CCError.CCError(common.CCErrCommParamsInvalid)
	}

	return result.Info[0].ObjectID, nil
}

func (s *service) getTemplateAttrByID(kit *rest.Kit, id int64, fields []string) ([]metadata.FieldTemplateAttr,
	errors.CCErrorCoder) {

	attrFilter := filtertools.GenAtomFilter(common.BKTemplateID, filter.Equal, id)
	listOpt := &metadata.CommonQueryOption{
		CommonFilterOption: metadata.CommonFilterOption{
			Filter: attrFilter,
		},
		Page: metadata.BasePage{
			Limit: common.BKNoLimit,
		},
		Fields: fields,
	}

	// list field template attributes
	res, err := s.clientSet.CoreService().FieldTemplate().ListFieldTemplateAttr(kit.Ctx, kit.Header, listOpt)
	if err != nil {
		blog.Errorf("list template attributes failed, template id: %+v, err: %v, rid: %s", id, err, kit.Rid)
		return nil, err
	}

	if len(res.Info) == 0 {
		blog.Errorf("no template attributes founded, template id: %d, rid: %s", id, kit.Rid)
		return []metadata.FieldTemplateAttr{}, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKTemplateID)
	}

	return res.Info, nil
}

// tmplAttrConvertObjAttr template attributes are converted to attrs on the model.
func tmplAttrConvertObjAttr(user, objID string, attr metadata.FieldTemplateAttr) *metadata.Attribute {

	return &metadata.Attribute{
		ObjectID:     objID,
		TemplateID:   attr.ID,
		PropertyID:   attr.PropertyID,
		Placeholder:  attr.Placeholder.Value,
		IsRequired:   attr.Required.Value,
		PropertyName: attr.PropertyName,
		PropertyType: attr.PropertyType,
		Default:      attr.Default,
		Option:       attr.Option,
		IsMultiple:   &attr.IsMultiple,
		Creator:      user,
		IsEditable:   attr.Editable.Value,
		Unit:         attr.Unit,
	}

}

// getFieldTemplateUniqueByID the unique verification keys of the template in the database
// are the self-incrementing IDs of the corresponding template properties, but there is no
// self-incrementing ID in the scene of creating a unique index. In order to be compatible
// with this scenario, the request of the comparison function is unified as propertyID. The
// function of this function is to The auto-increment ID in the database is converted to propertyID.
func (s *service) getFieldTemplateUniqueByID(kit *rest.Kit, op *metadata.SyncObjectTask) (
	*metadata.CompareFieldTmplUniqueOption, error) {

	uniqueFilter := filtertools.GenAtomFilter(common.BKTemplateID, filter.Equal, op.TemplateID)
	listOpt := &metadata.CommonQueryOption{
		CommonFilterOption: metadata.CommonFilterOption{Filter: uniqueFilter},
		Page: metadata.BasePage{
			Limit: common.BKNoLimit,
		},
	}

	// list field template uniques
	res, err := s.clientSet.CoreService().FieldTemplate().ListFieldTemplateUnique(kit.Ctx, kit.Header, listOpt)
	if err != nil {
		blog.Errorf("list field template uniques failed, err: %v, id: %+v, rid: %s", err, op.TemplateID, kit.Rid)
		return nil, err
	}

	result := new(metadata.CompareFieldTmplUniqueOption)

	result.TemplateID = op.TemplateID
	result.ObjectID = op.ObjectID
	result.Uniques = make([]metadata.FieldTmplUniqueForUpdate, len(res.Info))

	// indicates that no unique index is configured on this template
	if len(res.Info) == 0 {
		return result, nil
	}

	attrs, err := s.getTemplateAttrByID(kit, op.TemplateID, []string{common.BKFieldID, common.BKPropertyIDField})
	if err != nil {
		return nil, err
	}

	attrIDProMap := make(map[int64]string)
	for _, attr := range attrs {
		attrIDProMap[attr.ID] = attr.PropertyID
	}

	propertyIDs := make([]string, 0)
	for index := range res.Info {
		for _, key := range res.Info[index].Keys {
			propertyID, ok := attrIDProMap[key]
			if !ok {
				blog.Errorf("property id not found, attr id: %s,object id: %d, rid: %s", key, op.ObjectID, kit.Rid)
				return nil, kit.CCError.CCErrorf(common.CCErrCommNotFound, common.BKPropertyIDField)
			}
			propertyIDs = append(propertyIDs, propertyID)
		}
		result.Uniques[index].Keys = propertyIDs
		result.Uniques[index].ID = res.Info[index].ID
		propertyIDs = []string{}
	}
	return result, nil
}

// tmplUniqueConvertObjUnique the unique verification of the template is converted
// into the unique verification format of the model.
func (s *service) tmplUniqueConvertObjUnique(kit *rest.Kit, objID string,
	tmplUnique metadata.FieldTmplUniqueForUpdate) (metadata.ObjectUnique, error) {

	cond := &metadata.QueryCondition{
		Condition: mapstr.MapStr{
			common.BKObjIDField: objID,
			common.BKPropertyIDField: mapstr.MapStr{
				common.BKDBIN: tmplUnique.Keys,
			},
		},
	}
	result, err := s.clientSet.CoreService().Model().ReadModelAttrByCondition(kit.Ctx, kit.Header, cond)
	if err != nil {
		blog.Errorf("failed to read model attribute, cond: %+v, err: %v, rid: %s", cond.Condition, err, kit.Rid)
		return metadata.ObjectUnique{}, err
	}

	if len(result.Info) != len(tmplUnique.Keys) {
		blog.Errorf("the number of model attrs fetched is not as expected, num: %d, keys: %+v, rid: %s",
			len(result.Info), tmplUnique.Keys, kit.Rid)
	}

	// the index keys on the model save the attribute auto-increment ID on the object
	attrIDProMap := make(map[string]int64)
	for _, data := range result.Info {
		attrIDProMap[data.PropertyID] = data.ID
	}

	unique := metadata.ObjectUnique{}
	for _, key := range tmplUnique.Keys {
		unique.Keys = append(unique.Keys, metadata.UniqueKey{
			Kind: metadata.UniqueKeyKindProperty,
			ID:   uint64(attrIDProMap[key]),
		})
	}

	// bk_template_id on the model refers to the auto-increment ID on
	// the corresponding (same propertyID) template attribute
	unique.TemplateID = tmplUnique.ID
	unique.ObjID = objID
	return unique, nil
}

// preprocessTmplUnique process the unique verification after comparison, and return two sets of results,
// one is the unique verification content that needs to be added, and the other is the unique verification
// content that needs to be updated.
func (s *service) preprocessTmplUnique(kit *rest.Kit, objectID string, input *metadata.CompareFieldTmplUniqueOption,
	tUnique *metadata.CompareFieldTmplUniquesRes) ([]metadata.ObjectUnique, []metadata.ObjectUnique, error) {

	if len(tUnique.Update) == 0 && len(tUnique.Create) == 0 {
		return nil, nil, nil
	}
	// mainline object's unique can not be changed.
	yes, err := s.logics.AssociationOperation().IsMainlineObject(kit, objectID)
	if err != nil {
		return nil, nil, err
	}
	if yes && objectID != common.BKInnerObjIDHost {
		return nil, nil, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, objectID)
	}

	createUniques := make([]metadata.ObjectUnique, 0)
	updateUniques := make([]metadata.ObjectUnique, 0)
	if len(tUnique.Create) > 0 {
		for _, data := range tUnique.Create {
			createUnique, err := s.tmplUniqueConvertObjUnique(kit, objectID, input.Uniques[data.Index])
			if err != nil {
				return nil, nil, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, objectID)
			}
			createUniques = append(createUniques, createUnique)
		}
	}

	if len(tUnique.Update) > 0 {
		for id := range tUnique.Update {
			// -1 means that this index is deleted from the template,
			// and for the model side, it means unbinding
			if tUnique.Update[id].Index == -1 {
				tUnique.Update[id].Data.TemplateID = 0
				updateUniques = append(updateUniques, *tUnique.Update[id].Data)
				continue
			}
			updateUniques = append(updateUniques, *tUnique.Update[id].Data)
		}
	}
	return createUniques, updateUniques, nil
}

func (s *service) getCreateAndUpdateAttr(kit *rest.Kit, option *metadata.SyncObjectTask, objectID string) (
	[]*metadata.Attribute, []metadata.CompareOneFieldTmplAttrRes, error) {

	// 1、get the full content of the field template
	tmplAttrs, ccErr := s.getTemplateAttrByID(kit, option.TemplateID, []string{})
	if ccErr != nil {
		return nil, nil, ccErr
	}

	// verify the validity of field attributes and obtain
	// the contents of template fields that need to be synchronized
	opt := &metadata.CompareFieldTmplAttrOption{
		TemplateID: option.TemplateID,
		ObjectID:   option.ObjectID,
		Attrs:      tmplAttrs,
	}

	result, _, err := s.logics.FieldTemplateOperation().CompareFieldTemplateAttr(kit, opt, false)
	if err != nil {
		blog.Errorf("compare field template failed, cond: %+v, err: %v, rid: %s", opt, err, kit.Rid)
		return nil, nil, err
	}

	createAttr := make([]*metadata.Attribute, 0)
	if len(result.Create) > 0 {
		for _, attr := range result.Create {
			objAttr := tmplAttrConvertObjAttr(kit.User, objectID, tmplAttrs[attr.Index])
			createAttr = append(createAttr, objAttr)
		}
	}
	return createAttr, result.Update, nil
}

func (s *service) doSyncFieldTemplateTask(kit *rest.Kit, option *metadata.SyncObjectTask, objectID string) error {

	createAttr, updateAttr, err := s.getCreateAndUpdateAttr(kit, option, objectID)
	if err != nil {
		return err
	}

	txnErr := s.clientSet.CoreService().Txn().AutoRunTxn(kit.Ctx, kit.Header, func() error {
		// 1、create model properties
		if len(createAttr) > 0 {
			if err := s.logics.AttributeOperation().BatchCreateObjectAttr(kit, objectID, createAttr, true); err != nil {
				blog.Errorf("create model attribute failed, attr: %+v, err: %v, rid: %s", createAttr, err, kit.Rid)
				return err
			}
		}

		// 2、update model properties
		if len(updateAttr) > 0 {
			for _, attr := range updateAttr {
				if len(attr.UpdateData) == 0 {
					continue
				}
				err := s.logics.AttributeOperation().UpdateObjectAttribute(kit, attr.UpdateData, attr.Data.ID, 0, true)
				if err != nil {
					return err
				}
			}
		}

		// note: be sure to process the unique check after the attrs are processed, and the order cannot be reversed

		// 3、unique validation for preprocessed template synchronization

		uniqueOp, err := s.getFieldTemplateUniqueByID(kit, option)
		if err != nil {
			return err
		}

		res, _, err := s.logics.FieldTemplateOperation().CompareFieldTemplateUnique(kit, uniqueOp, false)
		if err != nil {
			blog.Errorf("get field template unique failed, cond: %+v, err: %v, rid: %s", uniqueOp, err, kit.Rid)
			return err
		}

		create, update, err := s.preprocessTmplUnique(kit, objectID, uniqueOp, res)
		if err != nil {
			blog.Errorf("get object unique failed, object: %+v, unique: %+v, err: %v, rid: %s", objectID,
				uniqueOp.Uniques, err, kit.Rid)
			return err
		}

		// 4、update the unique validation of the model
		if len(update) > 0 {
			for _, unique := range update {
				input := metadata.UpdateUniqueRequest{
					TemplateID: unique.TemplateID,
					Keys:       unique.Keys,
					LastTime:   metadata.Now()}

				op := metadata.UpdateModelAttrUnique{Data: input, FromTemplate: true}
				_, err := s.clientSet.CoreService().Model().UpdateModelAttrUnique(kit.Ctx, kit.Header,
					objectID, unique.ID, op)
				if err != nil {
					blog.Errorf("create unique failed, raw: %#v, err: %v, rid: %s", unique, err, kit.Rid)
					return err
				}
			}
		}

		// 5、create a unique check for the model
		if len(create) > 0 {
			for _, unique := range create {
				op := metadata.CreateModelAttrUnique{Data: unique, FromTemplate: true}
				_, err := s.clientSet.CoreService().Model().CreateModelAttrUnique(kit.Ctx, kit.Header, objectID, op)
				if err != nil {
					blog.Errorf("create unique failed for failed: raw: %#v, err: %v, rid: %s", unique, err, kit.Rid)
					return err
				}
			}
			return nil
		}

		return nil
	})

	if txnErr != nil {
		return txnErr
	}
	return nil
}

// CountFieldTemplateObj count field templates related objects
func (s *service) CountFieldTemplateObj(cts *rest.Contexts) {
	opt := new(metadata.CountFieldTmplResOption)
	if err := cts.DecodeInto(opt); err != nil {
		cts.RespAutoError(err)
		return
	}

	if rawErr := opt.Validate(); rawErr.ErrCode != 0 {
		cts.RespAutoError(rawErr.ToCCError(cts.Kit.CCError))
		return
	}

	relFilter := filtertools.GenAtomFilter(common.BKTemplateID, filter.In, opt.TemplateIDs)
	listOpt := &metadata.CommonQueryOption{
		CommonFilterOption: metadata.CommonFilterOption{Filter: relFilter},
		Page:               metadata.BasePage{Limit: common.BKNoLimit},
	}

	res, err := s.clientSet.CoreService().FieldTemplate().ListObjFieldTmplRel(cts.Kit.Ctx, cts.Kit.Header, listOpt)
	if err != nil {
		blog.Errorf("list field templates failed, err: %v, opt: %+v, rid: %s", err, opt, cts.Kit.Rid)
		cts.RespAutoError(err)
		return
	}

	countMap := make(map[int64]int, 0)
	for _, relation := range res.Info {
		countMap[relation.TemplateID]++
	}

	countInfos := make([]metadata.FieldTmplResCount, len(opt.TemplateIDs))
	for i, templateID := range opt.TemplateIDs {
		countInfos[i] = metadata.FieldTmplResCount{
			TemplateID: templateID,
			Count:      countMap[templateID],
		}
	}

	cts.RespEntity(countInfos)
}
