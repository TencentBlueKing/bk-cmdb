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

package label

import (
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/selector"
	"configcenter/src/common/util"
	"configcenter/src/source_controller/coreservice/core"
	"configcenter/src/storage/dal"
)

type labelOperation struct {
	dbProxy dal.RDB
}

// New create a new model manager instance
func New(dbProxy dal.RDB) core.LabelOperation {
	labelOps := &labelOperation{
		dbProxy: dbProxy,
	}
	return labelOps
}

func (p *labelOperation) AddLabel(ctx core.ContextParams, tableName string, option selector.LabelAddOption) errors.CCErrorCoder {
	if field, err := option.Labels.Validate(); err != nil {
		blog.Infof("addLabel failed, validate failed, field:%s, err: %+v, rid: %s", field, err, ctx.ReqID)
		return ctx.Error.CCErrorf(common.CCErrCommParamsInvalid, "label."+field)
	}

	idField := common.GetInstIDField(tableName)

	// check all instance validate
	option.InstanceIDs = util.IntArrayUnique(option.InstanceIDs)
	countFilter := map[string]interface{}{
		idField: map[string]interface{}{
			common.BKDBIN: option.InstanceIDs,
		},
	}
	if count, err := p.dbProxy.Table(tableName).Find(countFilter).Count(ctx.Context); err != nil {
		blog.ErrorJSON("AddLabel failed, db count instances failed, filter: %s, err: %s, rid: %s", countFilter, err.Error(), ctx.ReqID)
		return ctx.Error.CCErrorf(common.CCErrCommDBSelectFailed)
	} else if count != uint64(len(option.InstanceIDs)) {
		blog.ErrorJSON("AddLabel failed, some instance not valid, filter: %s, result count: %d, rid: %s", countFilter, count, ctx.ReqID)
		return ctx.Error.CCErrorf(common.CCErrCommParamsInvalid, "instance_ids")
	}

	for _, instanceID := range option.InstanceIDs {
		filter := map[string]interface{}{
			idField: instanceID,
		}
		data := &selector.LabelInstance{}
		if err := p.dbProxy.Table(tableName).Find(filter).One(ctx.Context, data); err != nil {
			blog.Errorf("AddLabel failed, get instance failed, instanceID: %+v, err: %+v, rid: %s", instanceID, err, ctx.ReqID)
			return ctx.Error.CCErrorf(common.CCErrCommDBSelectFailed)
		}
		if data.Labels != nil {
			data.Labels.AddLabel(option.Labels)
		} else {
			data.Labels = option.Labels
		}
		if err := p.dbProxy.Table(tableName).Update(ctx.Context, filter, data); err != nil {
			blog.Errorf("AddLabel failed, update instance failed, instanceID: %+v, err: %+v, rid: %s", instanceID, err, ctx.ReqID)
			return ctx.Error.CCErrorf(common.CCErrCommDBUpdateFailed)
		}
	}
	return nil
}

func (p *labelOperation) RemoveLabel(ctx core.ContextParams, tableName string, option selector.LabelRemoveOption) errors.CCErrorCoder {
	idField := common.GetInstIDField(tableName)

	// check all instance validate
	option.InstanceIDs = util.IntArrayUnique(option.InstanceIDs)
	countFilter := map[string]interface{}{
		idField: map[string]interface{}{
			common.BKDBIN: option.InstanceIDs,
		},
	}
	if count, err := p.dbProxy.Table(tableName).Find(countFilter).Count(ctx.Context); err != nil {
		blog.ErrorJSON("RemoveLabel failed, db count instances failed, filter: %s, err: %s, rid: %s", countFilter, err.Error(), ctx.ReqID)
		return ctx.Error.CCErrorf(common.CCErrCommDBSelectFailed)
	} else if count != uint64(len(option.InstanceIDs)) {
		blog.ErrorJSON("RemoveLabel failed, some instance not valid, filter: %s, result count: %d, rid: %s", countFilter, count, ctx.ReqID)
		return ctx.Error.CCErrorf(common.CCErrCommParamsInvalid, "instance_ids")
	}

	for _, instanceID := range option.InstanceIDs {
		filter := map[string]interface{}{
			idField: instanceID,
		}
		data := &selector.LabelInstance{}
		if err := p.dbProxy.Table(tableName).Find(filter).One(ctx.Context, data); err != nil {
			blog.Errorf("RemoveLabel failed, get instance failed, instanceID: %+v, err: %+v, rid: %s", instanceID, err, ctx.ReqID)
			return ctx.Error.CCErrorf(common.CCErrCommDBSelectFailed)
		}
		if data.Labels != nil {
			data.Labels.RemoveLabel(option.Keys)
		} else {
			data.Labels = make(map[string]string)
		}
		if err := p.dbProxy.Table(tableName).Update(ctx.Context, filter, data); err != nil {
			blog.Errorf("RemoveLabel failed, update instance failed, instanceID: %+v, err: %+v, rid: %s", instanceID, err, ctx.ReqID)
			return ctx.Error.CCErrorf(common.CCErrCommDBUpdateFailed)
		}
	}
	return nil
}
