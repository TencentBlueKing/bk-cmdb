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

	// validate service template id
	if option.ServiceTemplateIDs != nil && len(option.ServiceTemplateIDs) > 0 {
		serviceTemplateIDs := util.IntArrayUnique(option.ServiceTemplateIDs)
		filter := map[string]interface{}{
			common.BKFieldID: map[string]interface{}{
				common.BKDBIN: serviceTemplateIDs,
			},
			metadata.MetadataBizField: strconv.FormatInt(bizID, 10),
		}
		count, err := p.dbProxy.Table(common.BKTableNameSetTemplate).Find(filter).Count(ctx.Context)
		if err != nil {
			blog.Errorf("create set template failed, validate service template id field failed, filter: %+v, err: %+v, rid: %s", filter, err, ctx.ReqID)
			return setTemplate, ctx.Error.CCError(common.CCErrCommDBSelectFailed)
		}
		if count != uint64(len(serviceTemplateIDs)) {
			blog.Errorf("create set template failed, validate service template failed, some id invalid, bizID: %d, serviceTemplateIDs: %+v, rid: %s", bizID, serviceTemplateIDs)
			return setTemplate, ctx.Error.CCErrorf(common.CCErrCommParamsInvalid, "service_template_ids")
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

	// TODO need transaction ?
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
