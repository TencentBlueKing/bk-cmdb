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

func (s *service) preCheckExecuteTask(kit *rest.Kit, syncOption *metadata.SyncObjectTask) error {

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
		return kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKTemplateID)
	}

	if len(counts) != 1 || int(counts[0]) > 0 {
		blog.Errorf("list field template num error, opt: %+v, rid: %s", syncOption, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKTemplateID)
	}

	// 2、determine whether the model exists
	objCond := &metadata.QueryCondition{
		Fields: []string{metadata.ModelFieldIsPaused},
		Page: metadata.BasePage{
			Limit: common.BKNoLimit,
		},
		Condition: mapstr.MapStr{common.BKFieldID: syncOption.ObjectID},
	}

	result, objErr := s.clientSet.CoreService().Model().ReadModel(kit.Ctx, kit.Header, objCond)
	if objErr != nil {
		blog.Errorf("list objects failed, err: %v, opt: %+v, rid: %s", objErr, syncOption, kit.Rid)
		return kit.CCError.CCError(common.CCErrCommParamsInvalid)

	}
	if result.Count != 1 {
		blog.Errorf("objects num error,  count: %d, opt: %+v, rid: %s", result.Count, syncOption, kit.Rid)
		return kit.CCError.CCError(common.CCErrCommParamsInvalid)
	}

	// 3、determine whether the state of the model is enabled
	if result.Info[0].IsPaused {
		blog.Errorf("object status is stop, unable to execute sync task, opt: %+v, rid: %s", syncOption, kit.Rid)
		return kit.CCError.CCError(common.CCErrorTopoModelStopped)
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
		blog.Error("get illegal relation failed, cond: %+v, err: %v, rid: %s", relCond, err, kit.Rid)
		return kit.CCError.CCError(common.CCErrCommParamsInvalid)
	}

	if len(counts) != 1 || int(counts[0]) > 0 {
		blog.Error("get illegal relation count error, cond: %+v, rid: %s", relCond, kit.Rid)
		return kit.CCError.CCError(common.CCErrCommParamsInvalid)
	}

	return nil
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

	if err := s.preCheckExecuteTask(ctx.Kit, syncOption); err != nil {
		ctx.RespAutoError(err)
		return
	}

	txnErr := s.clientSet.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		if err := s.doSyncFieldTemplateTask(ctx.Kit, syncOption); err != nil {
			blog.Errorf("do sync service instance task(%#v) failed, err: %v, rid: %s", syncOption, err, ctx.Kit.Rid)
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

func (s *service) getFieldTemplateAttrForSyncByID(kit *rest.Kit, id int64, fields []string) (
	[]metadata.FieldTemplateAttr, error) {

	attrFilter := filtertools.GenAtomFilter(common.BKTemplateID, filter.Equal, id)
	listOpt := &metadata.CommonQueryOption{
		CommonFilterOption: metadata.CommonFilterOption{Filter: attrFilter},
		Page: metadata.BasePage{
			Limit: common.BKNoLimit,
		},
		Fields: fields,
	}

	// list field template attributes
	res, err := s.clientSet.CoreService().FieldTemplate().ListFieldTemplateAttr(kit.Ctx, kit.Header, listOpt)
	if err != nil {
		blog.Errorf("list template attributes failed, err: %v, template id: %+v, rid: %s", err, id, kit.Rid)
		return nil, err
	}

	if len(res.Info) == 0 {
		blog.Errorf("no template attributes founded, template id: %d, rid: %s", err, id, kit.Rid)
		return []metadata.FieldTemplateAttr{}, kit.CCError.CCError(common.CCErrCommParamsInvalid)
	}

	return res.Info, nil
}

func templateAttrConvertObjAttr(id int64, user, objID string, attr metadata.FieldTemplateAttr,
	time metadata.Time) *metadata.Attribute {

	objAttr := new(metadata.Attribute)

	objAttr.ObjectID = objID
	objAttr.TemplateID = id
	objAttr.PropertyID = attr.PropertyID
	objAttr.Placeholder = attr.Placeholder.Value
	objAttr.OwnerID = attr.OwnerID
	objAttr.IsRequired = attr.Required.Value
	objAttr.PropertyName = attr.PropertyName
	objAttr.PropertyType = attr.PropertyType
	objAttr.Default = attr.Default
	objAttr.Option = attr.Option
	objAttr.IsMultiple = &attr.IsMultiple
	objAttr.LastTime = &time
	objAttr.CreateTime = &time
	objAttr.Creator = user
	objAttr.IsEditable = attr.Editable.Value
	objAttr.Unit = attr.Unit
	return objAttr
}

func (s *service) getModelForSyncByID(kit *rest.Kit, id int64) (*metadata.Object, error) {
	objCond := &metadata.QueryCondition{
		Fields:    []string{common.BKObjIDField},
		Page:      metadata.BasePage{Limit: common.BKNoLimit},
		Condition: mapstr.MapStr{common.BKFieldID: id},
	}

	result, objErr := s.clientSet.CoreService().Model().ReadModel(kit.Ctx, kit.Header, objCond)
	if objErr != nil {
		blog.Errorf("list objects failed, err: %v, id: %+v, rid: %s", objErr, id, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrCommParamsInvalid)

	}
	if result.Count != 1 {
		blog.Errorf("objects num error,  count: %d, id: %+v, rid: %s", result.Count, id, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrCommParamsInvalid)
	}
	return &result.Info[0], nil
}

func (s *service) getFieldTemplateUniqueByID(kit *rest.Kit, op *metadata.SyncObjectTask) (
	*metadata.CompareFieldTmplUniqueOption, error) {

	// 这里得有一个attrID与propertyID的map
	fields := []string{common.BKFieldID, common.BKPropertyIDField}
	attrs, err := s.getFieldTemplateAttrForSyncByID(kit, op.TemplateID, fields)
	if err != nil {
		return nil, err
	}
	attrIDProMap := make(map[int64]string)
	for _, attr := range attrs {
		attrIDProMap[attr.ID] = attr.PropertyID
	}

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
		return nil, nil
	}
	result := new(metadata.CompareFieldTmplUniqueOption)
	result.Uniques = make([]metadata.FieldTmplUniqueForUpdate, len(res.Info))

	result.TemplateID = op.TemplateID
	result.ObjectID = op.ObjectID

	for index := range res.Info {
		propertyIDs := make([]string, 0)
		for _, key := range res.Info[index].Keys {
			propertyIDs = append(propertyIDs, attrIDProMap[key])
		}
		result.Uniques[index].Keys = propertyIDs
		result.Uniques[index].ID = res.Info[index].ID
	}

	return result, nil
}

func (s *service) doSyncFieldTemplateTask(kit *rest.Kit, option *metadata.SyncObjectTask) error {

	// 1、获取字段模版中的全量字段内容
	tmplAttrs, err := s.getFieldTemplateAttrForSyncByID(kit, option.TemplateID, []string{})
	if err != nil {
		return err
	}

	obj, err := s.getModelForSyncByID(kit, option.ObjectID)
	if err != nil {
		return err
	}

	// 1、校验字段和索引的合法性直接调用backend接口获取需要同步的模版字段内容
	now := metadata.Now()

	// 2、将需要创建的字段属性做一次类型转化后同步到模型属性中（templateID需要考虑赋值）
	templateAttr := new(metadata.CompareFieldTmplAttrsRes)
	attrs := make([]*metadata.Attribute, 0)
	if len(templateAttr.Create) > 0 {
		// 采用批量的方式进行创建模型的属性，这里需要封装一下
		for _, attr := range templateAttr.Create {
			objAttr := templateAttrConvertObjAttr(option.TemplateID, kit.User, obj.ObjectID, tmplAttrs[attr.Index], now)
			attrs = append(attrs, objAttr)
		}
	}
	uniques, err := s.getFieldTemplateUniqueByID(kit, option.TemplateID)

	if err != nil {
		return err
	}
	// 获取模版所有的index
	indexes := make([]metadata.ObjectUnique, 0)

	templateIndex := new(metadata.CompareFieldTmplUniquesRes)
	tpmlAttrIDPropertIDMap := make(map[int64]string)
	tpmlAttrIDs := make([]int64, 0)
	if len(templateIndex.Create) > 0 {
		tpmlAttrIDs = append(tpmlAttrIDs)
	}
	txnErr := s.clientSet.CoreService().Txn().AutoRunTxn(kit.Ctx, kit.Header, func() error {
		// 创建属性
		if err := s.logics.AttributeOperation().BatchCreateObjectAttr(kit, obj.ObjectID, attrs); err != nil {
			blog.Errorf("create model attribute failed, attr: %+v, err: %v, rid: %s", attrs, err, kit.Rid)
			return err
		}

		// 更新属性
		for _, attr := range templateAttr.Update {
			err := s.logics.AttributeOperation().UpdateObjectAttribute(kit, attr.UpdateData, attr.Data.ID, 0)
			if err != nil {
				return err
			}
		}
		// 创建索引

		// 更新索引
		return nil
	})

	if txnErr != nil {
		return txnErr
	}

	// 3、将需要更新的字段属性内容通过objectID和propertyID 进行更新
	// 4、对于字段模板的删除场景，转化为将templateID置为0
	// 5、获取需要同步的索引内容，即创建唯一索引并设置templateID
	// 6、通过对比获取到索引的创建和更新
	// 7、由于模版的索引结构和模型的索引结构不一致，这里需要做一次转换
	// 8、注意需要重新赋值一直更新的时间和更新人员
	return nil
}
