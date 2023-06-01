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

	// TODO add find object or template auth check after find object operation authorization is supported

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

	// TODO add find object or template auth check after find object operation authorization is supported

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

	// TODO add find object or template auth check after find object operation authorization is supported
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

	// list object by ids
	listOpt := &metadata.QueryCondition{
		Fields:    []string{common.BKFieldID, common.BKFieldName},
		Page:      opt.Page,
		Condition: mapstr.MapStr{common.BKFieldID: mapstr.MapStr{common.BKDBIN: objectIDs}},
	}
	res, objErr := s.clientSet.CoreService().Model().ReadModel(cts.Kit.Ctx, cts.Kit.Header, listOpt)
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

// SyncFieldTemplateToObjectTask synchronize field template information to model tasks.
func (s *service) SyncFieldTemplateToObjectTask(ctx *rest.Contexts) {

	syncOption := new(metadata.SyncObjectTask)
	if err := ctx.DecodeInto(syncOption); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if rawErr := syncOption.Validate(); rawErr.ErrCode != 0 {
		ctx.RespAutoError(rawErr.ToCCError(ctx.Kit.CCError))
		return
	}

	objectID, err := s.preCheckExecuteTaskAndGetObjID(ctx.Kit, syncOption)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	txnErr := s.clientSet.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		if err := s.doSyncFieldTemplateTask(ctx.Kit, syncOption, objectID); err != nil {
			blog.Errorf("do sync field template task(%#v) failed, err: %v, rid: %s", syncOption, err, ctx.Kit.Rid)
			return err
		}
		return nil
	})

	if txnErr != nil {
		ctx.RespAutoError(txnErr)
		return
	}

	ctx.RespEntity(nil)
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
		blog.Errorf("no template attributes founded, template id: %d, rid: %s", err, id, kit.Rid)
		return []metadata.FieldTemplateAttr{}, kit.CCError.CCError(common.CCErrCommParamsInvalid)
	}

	return res.Info, nil
}

// tmplAttrConvertObjAttr template attributes are converted to attributes on the model.
func tmplAttrConvertObjAttr(id int64, user, objID string, attr metadata.FieldTemplateAttr,
	time metadata.Time) *metadata.Attribute {

	return &metadata.Attribute{
		ObjectID:     objID,
		TemplateID:   id,
		PropertyID:   attr.PropertyID,
		Placeholder:  attr.Placeholder.Value,
		OwnerID:      attr.OwnerID,
		IsRequired:   attr.Required.Value,
		PropertyName: attr.PropertyName,
		PropertyType: attr.PropertyType,
		Default:      attr.Default,
		Option:       attr.Option,
		IsMultiple:   &attr.IsMultiple,
		LastTime:     &time,
		CreateTime:   &time,
		Creator:      user,
		IsEditable:   attr.Editable.Value,
		Unit:         attr.Unit,
	}

}

// getFieldTemplateUniqueByID 由于存在数据库中的模版唯一校验keys是属性的自增ID，而在新建唯一索引的场景是没有自增ID的，
// 所以为了兼容此场景，比较函数的请求统一为propertyID，此函数的作用是将数据库中的自增ID转化为propertyID。
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
	blog.ErrorJSON("8888888888888888 res: %s", res)
	// 表示本模版上没有配置唯一索引
	if len(res.Info) == 0 {
		return nil, nil
	}

	result := new(metadata.CompareFieldTmplUniqueOption)

	result.TemplateID = op.TemplateID
	result.ObjectID = op.ObjectID
	result.Uniques = make([]metadata.FieldTmplUniqueForUpdate, len(res.Info))

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
				blog.Errorf("property id not found, id: %d, rid: %s", key, kit.Rid)
				return nil, kit.CCError.CCErrorf(common.CCErrCommNotFound, common.BKPropertyIDField)
			}
			propertyIDs = append(propertyIDs, propertyID)
		}
		result.Uniques[index].Keys = propertyIDs
		result.Uniques[index].ID = res.Info[index].ID
		propertyIDs = []string{}
	}
	blog.ErrorJSON("666666666666666666 result: %s", result)
	return result, nil
}

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
	blog.ErrorJSON("1111111111111111 unique: %s", cond.Condition)

	if len(result.Info) != len(tmplUnique.Keys) {
		blog.Errorf("the number of model attrs fetched is not as expected, num: %d, keys: %+v, rid: %s",
			len(result.Info), tmplUnique.Keys, kit.Rid)
	}

	// object上的索引keys保存的是object上的属性自增ID
	attrIDProMap := make(map[string]int64)
	for _, data := range result.Info {
		attrIDProMap[data.PropertyType] = data.ID
	}
	blog.ErrorJSON("2222222222222222 attrIDProMap: %s", attrIDProMap)

	unique := metadata.ObjectUnique{}
	for _, key := range tmplUnique.Keys {
		unique.Keys = append(unique.Keys, metadata.UniqueKey{
			Kind: metadata.UniqueKeyKindProperty,
			ID:   uint64(attrIDProMap[key]),
		})
	}

	// 模型上的模版ID 指的是对应（相同property）模版属性上的自增ID
	unique.TemplateID = tmplUnique.ID
	unique.ObjID = objID
	unique.LastTime = metadata.Now()
	unique.OwnerID = kit.SupplierAccount
	blog.ErrorJSON("0000000000000 unique: %s", unique)
	return unique, nil
}

// preprocessTmplUnique 处理经过对比之后的索引，返回两组结果，一个是需要新增的索引内容，一个是需要更新的索引内容。
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
			// 将这个索引转化成object的索引
			createUnique, err := s.tmplUniqueConvertObjUnique(kit, objectID, input.Uniques[data.Index])
			if err != nil {
				return nil, nil, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, objectID)
			}
			createUniques = append(createUniques, createUnique)
		}
	}

	if len(tUnique.Update) > 0 {
		for id := range tUnique.Update {
			// -1 表示本索引从模版中删除，对于模型侧就是解除绑定关系
			if tUnique.Update[id].Index == -1 {
				input.Uniques[id].ID = 0
				unique, err := s.tmplUniqueConvertObjUnique(kit, objectID, input.Uniques[id])
				if err != nil {
					return nil, nil, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, objectID)
				}
				updateUniques = append(updateUniques, unique)
				continue
			}
			updateUniques = append(updateUniques, *tUnique.Update[id].Data)
		}
	}
	return createUniques, updateUniques, nil
}

func (s *service) getCreateAndUpdateAttr(kit *rest.Kit, option *metadata.SyncObjectTask, objectID string) (
	[]*metadata.Attribute, []metadata.CompareOneFieldTmplAttrRes, error) {

	// 1、获取字段模版中的全量字段内容
	tmplAttrs, ccErr := s.getTemplateAttrByID(kit, option.TemplateID, []string{})
	if ccErr != nil {
		return nil, nil, ccErr
	}

	// 校验字段和索引的合法性直接调用backend接口获取需要同步的模版字段内容
	opt := &metadata.CompareFieldTmplAttrOption{
		TemplateID: option.TemplateID,
		ObjectID:   option.ObjectID,
		Attrs:      tmplAttrs,
	}

	result, err := s.logics.FieldTemplateOperation().CompareFieldTemplateAttr(kit, opt, false)
	if err != nil {
		blog.Errorf("compare field template failed, cond: %+v, err: %v, rid: %s", opt, err, kit.Rid)
		return nil, nil, err
	}

	createAttr := make([]*metadata.Attribute, 0)
	if len(result.Create) > 0 {
		now := metadata.Now()
		for _, attr := range result.Create {
			objAttr := tmplAttrConvertObjAttr(option.TemplateID, kit.User, objectID, tmplAttrs[attr.Index], now)
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

	if len(createAttr) == 0 && len(updateAttr) == 0 {
		blog.Warnf("no difference between templates and models, opt: %+v, rid: %s", option, kit.Rid)
		return nil
	}

	uniqueOp, err := s.getFieldTemplateUniqueByID(kit, option)
	if err != nil {
		return err
	}

	createUniques, updateUniques := make([]metadata.ObjectUnique, 0), make([]metadata.ObjectUnique, 0)

	if uniqueOp != nil {
		res, err := s.logics.FieldTemplateOperation().CompareFieldTemplateUnique(kit, uniqueOp, false)
		if err != nil {
			blog.Errorf("get field template unique failed, cond: %+v, err: %v, rid: %s", uniqueOp, err, kit.Rid)
			return err
		}
		blog.ErrorJSON("qqqqqqqqqqqqqqqqqqq res: %s", res)
		create, update, err := s.preprocessTmplUnique(kit, objectID, uniqueOp, res)
		if err != nil {
			blog.Errorf("get object unique failed, object: %+v, unique: %+v, err: %v, rid: %s", objectID,
				uniqueOp.Uniques, err, kit.Rid)
			return err
		}
		createUniques, updateUniques = create, update
	}

	txnErr := s.clientSet.CoreService().Txn().AutoRunTxn(kit.Ctx, kit.Header, func() error {
		// 创建属性
		if len(createAttr) > 0 {
			if err := s.logics.AttributeOperation().BatchCreateObjectAttr(kit, objectID, createAttr); err != nil {
				blog.Errorf("create model attribute failed, attr: %+v, err: %v, rid: %s", createAttr, err, kit.Rid)
				return err
			}

		}
		// 更新属性
		if len(updateAttr) > 0 {
			for _, attr := range updateAttr {
				err := s.logics.AttributeOperation().UpdateObjectAttribute(kit, attr.UpdateData, attr.Data.ID, 0)
				if err != nil {
					return err
				}
			}
		}

		if len(createUniques) > 0 {
			for _, unique := range createUniques {
				op := metadata.CreateModelAttrUnique{Data: unique}
				_, err := s.clientSet.CoreService().Model().CreateModelAttrUnique(kit.Ctx, kit.Header, objectID, op)
				if err != nil {
					blog.Errorf("create unique failed for failed: raw: %#v, err: %v, rid: %s", unique, err, kit.Rid)
					return err
				}
			}
			return nil
		}

		if len(updateUniques) > 0 {
			for _, unique := range updateUniques {
				op := metadata.CreateModelAttrUnique{Data: unique}
				_, err := s.clientSet.CoreService().Model().CreateModelAttrUnique(kit.Ctx, kit.Header, objectID, op)
				if err != nil {
					blog.Errorf("create unique failed, raw: %#v, err: %v, rid: %s", unique, err, kit.Rid)
					return err
				}
			}
		}
		return nil
	})

	if txnErr != nil {
		return txnErr
	}
	return nil
}

// 3、将需要更新的字段属性内容通过objectID和propertyID 进行更新
// 4、对于字段模板的删除场景，转化为将templateID置为0
// 5、获取需要同步的索引内容，即创建唯一索引并设置templateID
// 6、通过对比获取到索引的创建和更新
// 8、注意需要重新赋值一直更新的时间和更新人员
