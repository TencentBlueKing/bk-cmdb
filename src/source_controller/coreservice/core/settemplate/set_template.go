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
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/source_controller/coreservice/core"
	"configcenter/src/storage/driver/mongodb"
)

type setTemplateOperation struct {
}

// New create a new model manager instance
func New() core.SetTemplateOperation {
	setTplOps := &setTemplateOperation{}
	return setTplOps
}

func (p *setTemplateOperation) ValidateBusinessID(kit *rest.Kit, bizID int64) errors.CCErrorCoder {
	filter := map[string]interface{}{
		common.BKAppIDField: bizID,
	}
	count, err := mongodb.Client().Table(common.BKTableNameBaseApp).Find(filter).Count(kit.Ctx)
	if err != nil {
		blog.Errorf("ValidateBusinessID failed, validate business id failed, db select failed, filter: %+v, err: %+v, rid: %s", filter, err, kit.Rid)
		return kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}
	if count == 0 {
		return kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField)
	}
	return nil
}

func (p *setTemplateOperation) ValidateServiceTemplateIDs(kit *rest.Kit, bizID int64, serviceTemplateIDs ...int64) ([]int64, errors.CCErrorCoder) {
	serviceTemplateIDs = util.IntArrayUnique(serviceTemplateIDs)
	filter := map[string]interface{}{
		common.BKFieldID: map[string]interface{}{
			common.BKDBIN: serviceTemplateIDs,
		},
		common.BKAppIDField: bizID,
	}
	count, err := mongodb.Client().Table(common.BKTableNameServiceTemplate).Find(filter).Count(kit.Ctx)
	if err != nil {
		blog.Errorf("validate service template id field failed, filter: %+v, err: %+v, rid: %s", filter, err, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}
	if count != uint64(len(serviceTemplateIDs)) {
		blog.Errorf("validate service template failed, some id invalid, bizID: %d, serviceTemplateIDs: %+v, rid: %s", bizID, serviceTemplateIDs, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, "service_template_ids")
	}
	return serviceTemplateIDs, nil
}

func (p *setTemplateOperation) CreateSetTemplate(kit *rest.Kit, bizID int64, option metadata.CreateSetTemplateOption) (metadata.SetTemplate, errors.CCErrorCoder) {
	now := time.Now()
	setTemplate := metadata.SetTemplate{
		ID:              0,
		Name:            option.Name,
		BizID:           bizID,
		Version:         0,
		Creator:         kit.User,
		Modifier:        kit.User,
		CreateTime:      now,
		LastTime:        now,
		SupplierAccount: kit.SupplierAccount,
	}
	if key, err := setTemplate.Validate(); err != nil {
		blog.Errorf("CreateSetTemplate failed, parameter invalid, key: %s, err: %+v, rid: %s", key, err, kit.Rid)
		return setTemplate, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, key)
	}

	// validate business id
	if err := p.ValidateBusinessID(kit, bizID); err != nil {
		blog.Errorf("CreateSetTemplate failed, validate business id failed, bizID: %d, err: %s, rid: %s", bizID, err.Error(), kit.Rid)
		return setTemplate, err
	}

	// validate service template id
	if option.ServiceTemplateIDs != nil && len(option.ServiceTemplateIDs) > 0 {
		serviceTemplateIDs, err := p.ValidateServiceTemplateIDs(kit, bizID, option.ServiceTemplateIDs...)
		if err != nil {
			return setTemplate, err
		}
		option.ServiceTemplateIDs = serviceTemplateIDs
	} else {
		option.ServiceTemplateIDs = make([]int64, 0)
	}

	// name unique validate
	nameFilter := map[string]interface{}{
		common.BKFieldName:  setTemplate.Name,
		common.BKAppIDField: setTemplate.BizID,
	}
	nameFilter = util.SetModOwner(nameFilter, kit.SupplierAccount)
	sameNameCount, err := mongodb.Client().Table(common.BKTableNameSetTemplate).Find(nameFilter).Count(kit.Ctx)
	if err != nil {
		blog.Errorf("create set template failed, filter same name records failed, filter: %+v, err: %+v, rid: %s", nameFilter, err, kit.Rid)
		return setTemplate, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}
	if sameNameCount > 0 {
		blog.Errorf("create set template failed, name duplicated, name: %s, rid: %s", setTemplate.Name, kit.Rid)
		return setTemplate, kit.CCError.CCErrorf(common.CCErrCommDuplicateItem, common.BKFieldName)
	}

	// generate id field
	id, err := mongodb.Client().NextSequence(kit.Ctx, common.BKTableNameSetTemplate)
	if nil != err {
		blog.Errorf("CreateSetTemplate failed, generate id failed, err: %+v, rid: %s", err, kit.Rid)
		return setTemplate, kit.CCError.CCErrorf(common.CCErrCommGenerateRecordIDFailed)
	}
	setTemplate.ID = int64(id)

	if err := mongodb.Client().Table(common.BKTableNameSetTemplate).Insert(kit.Ctx, setTemplate); err != nil {
		blog.Errorf("CreateSetTemplate failed, db insert failed, doc: %+v, err: %+v, rid: %s", setTemplate, err, kit.Rid)
		return setTemplate, kit.CCError.CCError(common.CCErrCommDBInsertFailed)
	}

	// TODO: add transaction
	if option.ServiceTemplateIDs != nil && len(option.ServiceTemplateIDs) > 0 {
		relations := make([]metadata.SetServiceTemplateRelation, 0)
		for _, serviceTemplateID := range option.ServiceTemplateIDs {
			relations = append(relations, metadata.SetServiceTemplateRelation{
				BizID:             setTemplate.BizID,
				SetTemplateID:     setTemplate.ID,
				ServiceTemplateID: serviceTemplateID,
				SupplierAccount:   kit.SupplierAccount,
			})
		}
		if err := mongodb.Client().Table(common.BKTableNameSetServiceTemplateRelation).Insert(kit.Ctx, relations); err != nil {
			blog.Errorf("CreateSetTemplate failed, db insert failed, doc: %+v, err: %+v, rid: %s", setTemplate, err, kit.Rid)
			return setTemplate, kit.CCError.CCError(common.CCErrCommDBInsertFailed)
		}
	}

	return setTemplate, nil
}

func (p *setTemplateOperation) UpdateSetTemplate(kit *rest.Kit, setTemplateID int64, option metadata.UpdateSetTemplateOption) (metadata.SetTemplate, errors.CCErrorCoder) {
	setTemplate := metadata.SetTemplate{}

	if errKey, err := option.Validate(); err != nil {
		blog.Errorf("UpdateSetTemplate failed, update option validate failed, option: %+v, key: %s, err: %+v, rid: %s", option, errKey, err, kit.Rid)
		return setTemplate, kit.CCError.CCError(common.CCErrCommHTTPBodyEmpty)
	}

	filter := map[string]interface{}{
		common.BKFieldID: setTemplateID,
	}
	filter = util.SetModOwner(filter, kit.SupplierAccount)
	if err := mongodb.Client().Table(common.BKTableNameSetTemplate).Find(filter).One(kit.Ctx, &setTemplate); err != nil {
		if mongodb.Client().IsNotFoundError(err) {
			blog.Errorf("UpdateSetTemplate failed, set template not found, id: %d, rid: %s", setTemplateID, kit.Rid)
			return setTemplate, kit.CCError.CCError(common.CCErrCommNotFound)
		}

		blog.Errorf("UpdateSetTemplate failed, select set template failed, filter: %+v, err: %+v, rid: %s", filter, err, kit.Rid)
		return setTemplate, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}

	if len(option.Name) > 0 && option.Name != setTemplate.Name {
		sameNameFilter := map[string]interface{}{
			common.BKFieldName:  option.Name,
			common.BKAppIDField: setTemplate.BizID,
			common.BKFieldID: map[string]interface{}{
				common.BKDBNE: setTemplate.ID,
			},
		}
		sameNameCount, err := mongodb.Client().Table(common.BKTableNameSetTemplate).Find(sameNameFilter).Count(kit.Ctx)
		if err != nil {
			blog.Errorf("UpdateSetTemplate failed, db select failed, filter: %+v, err: %+v, rid: %s", sameNameFilter, err, kit.Rid)
			return setTemplate, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
		}
		if sameNameCount > 0 {
			blog.Errorf("UpdateSetTemplate failed, name duplicated, filter: %+v, count: %d, rid: %s", sameNameFilter, sameNameCount, kit.Rid)
			return setTemplate, kit.CCError.CCErrorf(common.CCErrCommDuplicateItem, common.BKFieldName)
		}
		setTemplate.Name = option.Name
	}

	// TODO: add transaction
	if option.ServiceTemplateIDs != nil {
		serviceTemplateIDs, err := p.ValidateServiceTemplateIDs(kit, setTemplate.BizID, option.ServiceTemplateIDs...)
		if err != nil {
			return setTemplate, err
		}

		// list original related service template ids
		relations := make([]metadata.SetServiceTemplateRelation, 0)
		relationFilter := map[string]interface{}{
			common.BKSetTemplateIDField: setTemplateID,
		}
		if err := mongodb.Client().Table(common.BKTableNameSetServiceTemplateRelation).Find(relationFilter).All(kit.Ctx, &relations); err != nil {
			blog.Errorf("UpdateSetTemplate failed, db select relation failed, filter: %+v, err: %+v, rid: %s", relationFilter, err, kit.Rid)
		}

		// add new add service template ids
		existIDMap := make(map[int64]bool)
		for _, item := range relations {
			existIDMap[item.ServiceTemplateID] = true
		}
		addRelations := make([]metadata.SetServiceTemplateRelation, 0)
		for _, serviceTemplateID := range serviceTemplateIDs {
			if _, exist := existIDMap[serviceTemplateID]; !exist {
				addRelations = append(addRelations, metadata.SetServiceTemplateRelation{
					BizID:             setTemplate.BizID,
					SetTemplateID:     setTemplate.ID,
					ServiceTemplateID: serviceTemplateID,
					SupplierAccount:   kit.SupplierAccount,
				})
			}
		}
		if len(addRelations) > 0 {
			if err := mongodb.Client().Table(common.BKTableNameSetServiceTemplateRelation).Insert(kit.Ctx, addRelations); err != nil {
				blog.Errorf("UpdateSetTemplate failed, insert relations failed, docs: %+v, err: %+v, rid: %s", addRelations, err, kit.Rid)
				return setTemplate, kit.CCError.CCError(common.CCErrCommDBInsertFailed)
			}
		}

		// remove not need service template ids
		targetIDMap := make(map[int64]bool)
		for _, serviceTemplateID := range serviceTemplateIDs {
			targetIDMap[serviceTemplateID] = true
		}
		removeIDs := make([]int64, 0)
		for _, item := range relations {
			if _, exist := targetIDMap[item.ServiceTemplateID]; !exist {
				removeIDs = append(removeIDs, item.ServiceTemplateID)
			}
		}
		// TODO add reference check
		if len(removeIDs) > 0 {
			removeFilter := map[string]interface{}{
				common.BKSetTemplateIDField: setTemplate.ID,
				common.BKServiceTemplateIDField: map[string]interface{}{
					common.BKDBIN: removeIDs,
				},
			}
			if err := mongodb.Client().Table(common.BKTableNameSetServiceTemplateRelation).Delete(kit.Ctx, removeFilter); err != nil {
				blog.Errorf("UpdateSetTemplate failed, remove no need service template ids failed, filter: %+v, err: %+v, rid: %s", removeFilter, err, kit.Rid)
				return setTemplate, kit.CCError.CCError(common.CCErrCommDBDeleteFailed)
			}
		}
		if len(addRelations) > 0 || len(removeIDs) > 0 {
			setTemplate.Version += 1
		}
	}

	setTemplate.LastTime = time.Now()
	setTemplate.Modifier = kit.User
	if err := mongodb.Client().Table(common.BKTableNameSetTemplate).Update(kit.Ctx, filter, setTemplate); err != nil {
		blog.Errorf("UpdateSetTemplate failed, err: %+v, rid: %s", err, kit.Rid)
		return setTemplate, kit.CCError.CCError(common.CCErrCommDBUpdateFailed)
	}

	return setTemplate, nil
}

func (p *setTemplateOperation) DeleteSetTemplate(kit *rest.Kit, bizID int64, option metadata.DeleteSetTemplateOption) errors.CCErrorCoder {
	// check reference
	setFilter := map[string]interface{}{
		common.BKAppIDField:   bizID,
		common.BKOwnerIDField: kit.SupplierAccount,
		common.BKSetTemplateIDField: map[string]interface{}{
			common.BKDBIN: option.SetTemplateIDs,
		},
	}
	referenceCount, err := mongodb.Client().Table(common.BKTableNameBaseSet).Find(setFilter).Count(kit.Ctx)
	if err != nil {
		blog.Error("DeleteSetTemplate failed, get template reference count failed, db select failed, filter: %+v, err: %+v, rid: %s", setFilter, err, kit.Rid)
		return kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}
	if referenceCount > 0 {
		blog.Error("DeleteSetTemplate failed, set templates referenced by %d sets, setTemplates: %+v, rid: %s", referenceCount, option.SetTemplateIDs, kit.Rid)
		return kit.CCError.CCError(common.CCErrCommRemoveReferencedRecordForbidden)
	}

	// delete set template
	setTplFilter := map[string]interface{}{
		common.BKAppIDField:   bizID,
		common.BKOwnerIDField: kit.SupplierAccount,
		common.BKFieldID: map[string]interface{}{
			common.BKDBIN: option.SetTemplateIDs,
		},
	}
	if err := mongodb.Client().Table(common.BKTableNameSetTemplate).Delete(kit.Ctx, setTplFilter); err != nil {
		blog.Errorf("DeleteSetTemplate failed, db remove failed, filter: %+v, err: %+v, rid: %s", setTplFilter, err, kit.Rid)
		return kit.CCError.CCError(common.CCErrCommDBDeleteFailed)
	}

	// delete relations
	relationFilter := map[string]interface{}{
		common.BKAppIDField:   bizID,
		common.BKOwnerIDField: kit.SupplierAccount,
		common.BKSetTemplateIDField: map[string]interface{}{
			common.BKDBIN: option.SetTemplateIDs,
		},
	}
	if err := mongodb.Client().Table(common.BKTableNameSetServiceTemplateRelation).Delete(kit.Ctx, relationFilter); err != nil {
		blog.Errorf("DeleteSetTemplate failed, db remove relations failed, filter: %+v, err: %+v, rid: %s", relationFilter, err, kit.Rid)
		return kit.CCError.CCError(common.CCErrCommDBDeleteFailed)
	}

	return nil
}

func (p *setTemplateOperation) GetSetTemplate(kit *rest.Kit, bizID int64, setTemplateID int64) (metadata.SetTemplate, errors.CCErrorCoder) {
	setTemplate := metadata.SetTemplate{}
	filter := map[string]interface{}{
		common.BKFieldID:    setTemplateID,
		common.BKAppIDField: bizID,
	}
	filter = util.SetQueryOwner(filter, kit.SupplierAccount)
	if err := mongodb.Client().Table(common.BKTableNameSetTemplate).Find(filter).One(kit.Ctx, &setTemplate); err != nil {
		if mongodb.Client().IsNotFoundError(err) {
			blog.Errorf("GetSetTemplate failed, db select failed, not found, filter: %+v, err: %+v, rid: %s", filter, err, kit.Rid)
			return setTemplate, kit.CCError.CCError(common.CCErrCommNotFound)
		}
		blog.Errorf("GetSetTemplate failed, db select failed, filter: %+v, err: %+v, rid: %s", filter, err, kit.Rid)
		return setTemplate, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}
	return setTemplate, nil
}

func (p *setTemplateOperation) ListSetTemplate(kit *rest.Kit, bizID int64, option metadata.ListSetTemplateOption) (metadata.MultipleSetTemplateResult, errors.CCErrorCoder) {
	result := metadata.MultipleSetTemplateResult{}
	if option.Page.Limit > common.BKMaxPageSize && option.Page.Limit != common.BKNoLimit {
		return result, kit.CCError.CCError(common.CCErrCommPageLimitIsExceeded)
	}

	filter := map[string]interface{}{
		common.BKAppIDField: bizID,
	}
	filter = util.SetQueryOwner(filter, kit.SupplierAccount)
	if option.SetTemplateIDs != nil {
		filter[common.BKFieldID] = map[string]interface{}{
			common.BKDBIN: option.SetTemplateIDs,
		}
	}
	query := mongodb.Client().Table(common.BKTableNameSetTemplate).Find(filter)
	total, err := query.Count(kit.Ctx)
	if err != nil {
		blog.ErrorJSON("ListSetTemplate failed, db count failed, filter: %s, err: %s, rid: %s", filter, err.Error(), kit.Rid)
		return result, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}
	result.Count = int64(total)

	if len(option.Page.Sort) > 0 {
		query = query.Sort(option.Page.Sort)
	}
	if option.Page.Limit > 0 && option.Page.Limit != common.BKNoLimit {
		query = query.Limit(uint64(option.Page.Limit))
	}
	if option.Page.Start > 0 {
		query = query.Start(uint64(option.Page.Start))
	}

	setTemplates := make([]metadata.SetTemplate, 0)
	if err := query.All(kit.Ctx, &setTemplates); err != nil {
		blog.ErrorJSON("ListSetTemplate failed, db select failed, filter: %s, err: %s, rid: %s", filter, err.Error(), kit.Rid)
		return result, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}

	result.Info = setTemplates
	return result, nil
}

func (p *setTemplateOperation) ListSetServiceTemplateRelations(kit *rest.Kit, bizID int64, setTemplateID int64) ([]metadata.SetServiceTemplateRelation, errors.CCErrorCoder) {
	filter := map[string]interface{}{
		common.BKAppIDField:         bizID,
		common.BKSetTemplateIDField: setTemplateID,
	}
	filter = util.SetQueryOwner(filter, kit.SupplierAccount)

	setServiceTemplateRelations := make([]metadata.SetServiceTemplateRelation, 0)
	if err := mongodb.Client().Table(common.BKTableNameSetServiceTemplateRelation).Find(filter).All(kit.Ctx, &setServiceTemplateRelations); err != nil {
		blog.Errorf("ListSetServiceTemplateRelations failed, db select failed, filter: %+v, err: %+v, rid: %s", filter, err, kit.Rid)
		return setServiceTemplateRelations, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}

	return setServiceTemplateRelations, nil
}

func (p *setTemplateOperation) ListSetTplRelatedSvcTpl(kit *rest.Kit, bizID, setTemplateID int64) ([]metadata.ServiceTemplate, errors.CCErrorCoder) {
	relations, err := p.ListSetServiceTemplateRelations(kit, bizID, setTemplateID)
	if err != nil {
		blog.Errorf("ListSetTplRelatedSvcTpl failed, do core service list failed, bizID: %d, setTemplateID: %+v, err: %+v, rid: %s", bizID, setTemplateID, err, kit.Rid)
		return nil, err
	}
	serviceTemplateIDs := make([]int64, 0)
	for _, relation := range relations {
		serviceTemplateIDs = append(serviceTemplateIDs, relation.ServiceTemplateID)
	}
	filter := map[string]interface{}{
		common.BKFieldID: map[string]interface{}{
			common.BKDBIN: serviceTemplateIDs,
		},
	}

	// is it appropriate to visit service template table here?
	serviceTemplates := make([]metadata.ServiceTemplate, 0)
	if err := mongodb.Client().Table(common.BKTableNameServiceTemplate).Find(filter).All(kit.Ctx, &serviceTemplates); err != nil {
		blog.Errorf("ListSetTplRelatedSvcTpl failed, db select failed, serviceTemplateIDs: %+v, err: %s, rid: %s", serviceTemplateIDs, err.Error(), kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}

	return serviceTemplates, nil
}
