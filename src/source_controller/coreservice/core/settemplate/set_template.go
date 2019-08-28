/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.,
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the ",License",); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an ",AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package settemplate

import (
	"strconv"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/source_controller/coreservice/core"
	"configcenter/src/storage/dal"
)

type setTemplateOperation struct {
	dbProxy dal.RDB
}

// New create a new model manager instance
func New(dbProxy dal.RDB) core.SetTemplateOperation {
	setTplOps := &setTemplateOperation{
		dbProxy: dbProxy,
	}
	return setTplOps
}

func (p *setTemplateOperation) ValidateBusinessID(ctx core.ContextParams, bizID int64) errors.CCErrorCoder {
	filter := map[string]interface{}{
		common.BKAppIDField: bizID,
	}
	count, err := p.dbProxy.Table(common.BKTableNameBaseApp).Find(filter).Count(ctx.Context)
	if err != nil {
		blog.Errorf("ValidateBusinessID failed, validate business id failed, db select failed, filter: %+v, err: %+v, rid: %s", filter, err, ctx.ReqID)
		return ctx.Error.CCError(common.CCErrCommDBSelectFailed)
	}
	if count == 0 {
		return ctx.Error.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField)
	}
	return nil
}
func (p *setTemplateOperation) ValidateServiceTemplateIDs(ctx core.ContextParams, bizID int64, serviceTemplateIDs ...int64) ([]int64, errors.CCErrorCoder) {
	serviceTemplateIDs = util.IntArrayUnique(serviceTemplateIDs)
	filter := map[string]interface{}{
		common.BKFieldID: map[string]interface{}{
			common.BKDBIN: serviceTemplateIDs,
		},
		metadata.MetadataBizField: strconv.FormatInt(bizID, 10),
	}
	count, err := p.dbProxy.Table(common.BKTableNameServiceTemplate).Find(filter).Count(ctx.Context)
	if err != nil {
		blog.Errorf("validate service template id field failed, filter: %+v, err: %+v, rid: %s", filter, err, ctx.ReqID)
		return nil, ctx.Error.CCError(common.CCErrCommDBSelectFailed)
	}
	if count != uint64(len(serviceTemplateIDs)) {
		blog.Errorf("validate service template failed, some id invalid, bizID: %d, serviceTemplateIDs: %+v, rid: %s", bizID, serviceTemplateIDs)
		return nil, ctx.Error.CCErrorf(common.CCErrCommParamsInvalid, "service_template_ids")
	}
	return serviceTemplateIDs, nil
}

func (p *setTemplateOperation) CreateSetTemplate(ctx core.ContextParams, bizID int64, option metadata.CreateSetTemplateOption) (metadata.SetTemplate, errors.CCErrorCoder) {
	now := time.Now()
	setTemplate := metadata.SetTemplate{
		Name:            option.Name,
		BizID:           bizID,
		Creator:         ctx.User,
		Modifier:        ctx.User,
		CreateTime:      now,
		LastTime:        now,
		SupplierAccount: ctx.SupplierAccount,
	}
	if key, err := setTemplate.Validate(); err != nil {
		blog.Errorf("CreateSetTemplate failed, parameter invalid, key: %s, err: %+v, rid: %s", key, err, ctx.ReqID)
		return setTemplate, ctx.Error.CCErrorf(common.CCErrCommParamsInvalid, key)
	}

	// validate business id
	if err := p.ValidateBusinessID(ctx, bizID); err != nil {
		blog.Errorf("CreateSetTemplate failed, validate business id failed, bizID: %d, err: %s, rid: %s", bizID, err.Error(), ctx.ReqID)
		return setTemplate, err
	}

	// validate service template id
	if option.ServiceTemplateIDs != nil && len(option.ServiceTemplateIDs) > 0 {
		serviceTemplateIDs, err := p.ValidateServiceTemplateIDs(ctx, bizID, option.ServiceTemplateIDs...)
		if err != nil {
			return setTemplate, err
		}
		option.ServiceTemplateIDs = serviceTemplateIDs
	}

	// name unique validate
	nameFilter := map[string]interface{}{
		common.BKFieldName:  setTemplate.Name,
		common.BKAppIDField: setTemplate.BizID,
	}
	sameNameCount, err := p.dbProxy.Table(common.BKTableNameSetTemplate).Find(nameFilter).Count(ctx.Context)
	if err != nil {
		blog.Errorf("create set template failed, filter same name records failed, filter: %+v, err: %+v, rid: %s", nameFilter, err, ctx.ReqID)
		return setTemplate, ctx.Error.CCError(common.CCErrCommDBSelectFailed)
	}
	if sameNameCount > 0 {
		blog.Errorf("create set template failed, name duplicated, name: %s, rid: %s", setTemplate.Name, ctx.ReqID)
		return setTemplate, ctx.Error.CCErrorf(common.CCErrCommDuplicateItem, common.BKFieldName)
	}

	// generate id field
	id, err := p.dbProxy.NextSequence(ctx, common.BKTableNameSetTemplate)
	if nil != err {
		blog.Errorf("CreateSetTemplate failed, generate id failed, err: %+v, rid: %s", err, ctx.ReqID)
		return setTemplate, ctx.Error.CCErrorf(common.CCErrCommGenerateRecordIDFailed)
	}
	setTemplate.ID = int64(id)

	if err := p.dbProxy.Table(common.BKTableNameSetTemplate).Insert(ctx.Context, setTemplate); err != nil {
		blog.Errorf("CreateSetTemplate failed, db insert failed, doc: %+v, err: %+v, rid: %s", setTemplate, err, ctx.ReqID)
		return setTemplate, ctx.Error.CCError(common.CCErrCommDBInsertFailed)
	}

	// TODO: add transaction
	if option.ServiceTemplateIDs != nil && len(option.ServiceTemplateIDs) > 0 {
		relations := make([]metadata.SetServiceTemplateRelation, 0)
		for _, serviceTemplateID := range option.ServiceTemplateIDs {
			relations = append(relations, metadata.SetServiceTemplateRelation{
				SetTemplateID:     setTemplate.ID,
				ServiceTemplateID: serviceTemplateID,
			})
		}
		if err := p.dbProxy.Table(common.BKTableNameSetServiceTemplateRelation).Insert(ctx.Context, relations); err != nil {
			blog.Errorf("CreateSetTemplate failed, db insert failed, doc: %+v, err: %+v, rid: %s", setTemplate, err, ctx.ReqID)
			return setTemplate, ctx.Error.CCError(common.CCErrCommDBInsertFailed)
		}
	}

	return setTemplate, nil
}

func (p *setTemplateOperation) UpdateSetTemplate(ctx core.ContextParams, setTemplateID int64, option metadata.UpdateSetTemplateOption) (metadata.SetTemplate, errors.CCErrorCoder) {
	setTemplate := metadata.SetTemplate{}

	if errKey, err := option.Validate(); err != nil {
		blog.Errorf("UpdateSetTemplate failed, update option validate failed, option: %+v, key: %s, err: %+v, rid: %s", option, errKey, err, ctx.ReqID)
		return setTemplate, ctx.Error.CCError(common.CCErrCommHTTPBodyEmpty)
	}

	filter := map[string]interface{}{
		common.BKFieldID:      setTemplateID,
		common.BKOwnerIDField: ctx.SupplierAccount,
	}
	if err := p.dbProxy.Table(common.BKTableNameSetTemplate).Find(filter).One(ctx.Context, &setTemplate); err != nil {
		if p.dbProxy.IsNotFoundError(err) {
			blog.Errorf("UpdateSetTemplate failed, set template not found, id: %d, rid: %s", setTemplateID, ctx.ReqID)

			return setTemplate, ctx.Error.CCError(common.CCErrCommNotFound)
		}

		blog.Errorf("UpdateSetTemplate failed, select set template failed, filter: %+v, err: %+v, rid: %s", filter, err, ctx.ReqID)
		return setTemplate, ctx.Error.CCError(common.CCErrCommDBSelectFailed)
	}

	if len(option.Name) > 0 && option.Name != setTemplate.Name {
		sameNameFilter := map[string]interface{}{
			common.BKFieldName:  option.Name,
			common.BKAppIDField: setTemplate.BizID,
			common.BKFieldID: map[string]interface{}{
				common.BKDBNE: setTemplate.ID,
			},
		}
		sameNameCount, err := p.dbProxy.Table(common.BKTableNameSetTemplate).Find(sameNameFilter).Count(ctx.Context)
		if err != nil {
			blog.Errorf("UpdateSetTemplate failed, db select failed, filter: %+v, err: %+v, rid: %s", sameNameFilter, err, ctx.ReqID)
			return setTemplate, ctx.Error.CCError(common.CCErrCommDBSelectFailed)
		}
		if sameNameCount > 0 {
			blog.Errorf("UpdateSetTemplate failed, name duplicated, filter: %+v, count: %d, rid: %s", sameNameFilter, sameNameCount, ctx.ReqID)
			return setTemplate, ctx.Error.CCErrorf(common.CCErrCommDuplicateItem, common.BKFieldName)
		}
		setTemplate.Name = option.Name
	}

	// TODO: add transaction
	if option.ServiceTemplateIDs != nil {
		serviceTemplateIDs, err := p.ValidateServiceTemplateIDs(ctx, setTemplate.BizID, option.ServiceTemplateIDs...)
		if err != nil {
			return setTemplate, err
		}

		// list original related service template ids
		relations := make([]metadata.SetServiceTemplateRelation, 0)
		relationFilter := map[string]interface{}{
			common.BKSetTemplateIDField: setTemplateID,
		}
		if err := p.dbProxy.Table(common.BKTableNameSetServiceTemplateRelation).Find(relationFilter).All(ctx.Context, &relations); err != nil {
			blog.Errorf("UpdateSetTemplate failed, db select relation failed, filter: %+v, err: %+v, rid: %s", relationFilter, err, ctx.ReqID)
		}

		// add new add service template ids
		existIDMap := make(map[int64]bool)
		for _, item := range relations {
			existIDMap[item.ServiceTemplateID] = true
		}
		addRelations := make([]metadata.SetServiceTemplateRelation, 0)
		for _, serviceTemplateID := range serviceTemplateIDs {
			if _, exist := existIDMap[serviceTemplateID]; exist == false {
				addRelations = append(addRelations, metadata.SetServiceTemplateRelation{
					SetTemplateID:     setTemplate.ID,
					ServiceTemplateID: serviceTemplateID,
				})
			}
		}
		if len(addRelations) > 0 {
			if err := p.dbProxy.Table(common.BKTableNameSetServiceTemplateRelation).Insert(ctx.Context, addRelations); err != nil {
				blog.Errorf("UpdateSetTemplate failed, insert relations failed, docs: %+v, err: %+v, rid: %s", addRelations, err, ctx.ReqID)
				return setTemplate, ctx.Error.CCError(common.CCErrCommDBInsertFailed)
			}
		}

		// remove not need service template ids
		targetIDMap := make(map[int64]bool)
		for _, serviceTemplateID := range serviceTemplateIDs {
			targetIDMap[serviceTemplateID] = true
		}
		removeIDs := make([]int64, 0)
		for _, item := range relations {
			if _, exist := targetIDMap[item.ServiceTemplateID]; exist == false {
				removeIDs = append(removeIDs, item.ServiceTemplateID)
			}
		}
		if len(removeIDs) > 0 {
			removeFilter := map[string]interface{}{
				common.BKSetTemplateIDField: setTemplate.ID,
				common.BKServiceTemplateIDField: map[string]interface{}{
					common.BKDBIN: removeIDs,
				},
			}
			if err := p.dbProxy.Table(common.BKTableNameSetServiceTemplateRelation).Delete(ctx.Context, removeFilter); err != nil {
				blog.Errorf("UpdateSetTemplate failed, remove no need service template ids failed, filter: %+v, err: %+v, rid: %s", removeFilter, err, ctx.ReqID)
				return setTemplate, ctx.Error.CCError(common.CCErrCommDBDeleteFailed)
			}
		}
	}

	setTemplate.LastTime = time.Now()
	setTemplate.Modifier = ctx.User
	if err := p.dbProxy.Table(common.BKTableNameSetTemplate).Update(ctx.Context, filter, setTemplate); err != nil {
		blog.Errorf("UpdateSetTemplate failed, err: %+v, rid: %s", err, ctx.ReqID)
		return setTemplate, ctx.Error.CCError(common.CCErrCommDBUpdateFailed)
	}

	return setTemplate, nil
}
