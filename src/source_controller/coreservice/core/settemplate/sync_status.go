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
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/metadata"
	"configcenter/src/source_controller/coreservice/core"
)

func (p *setTemplateOperation) UpdateSetTemplateSyncStatus(ctx core.ContextParams, setID int64, option metadata.SetTemplateSyncStatus) errors.CCErrorCoder {
	if setID != option.SetID {
		return ctx.Error.CCErrorf(common.CCErrCommParamsInvalid, common.BKSetIDField)
	}

	filter := map[string]interface{}{
		common.BKSetIDField: setID,
	}
	if err := p.dbProxy.Table(common.BKTableNameSetTemplateSyncStatus).Upsert(ctx.Context, filter, option); err != nil {
		blog.Errorf("UpdateSetTemplateSyncStatus failed, db upsert sync status failed, id: %d, option: %s, err: %s, rid: %s", setID, option, err.Error(), ctx.ReqID)
		return ctx.Error.CCError(common.CCErrCommDBUpdateFailed)
	}

	if len(option.TaskID) == 0 {
		return nil
	}

	historyFilter := map[string]interface{}{
		common.BKTaskIDField: option.TaskID,
	}
	if err := p.dbProxy.Table(common.BKTableNameSetTemplateSyncHistory).Upsert(ctx.Context, historyFilter, option); err != nil {
		blog.Errorf("UpdateSetTemplateSyncStatus failed, db upsert sync history failed, id: %d, option: %s, err: %s, rid: %s", setID, option, err.Error(), ctx.ReqID)
		return ctx.Error.CCError(common.CCErrCommDBUpdateFailed)
	}

	return nil
}

func (p *setTemplateOperation) ListSetTemplateSyncStatus(ctx core.ContextParams, option metadata.ListSetTemplateSyncStatusOption) (metadata.MultipleSetTemplateSyncStatus, errors.CCErrorCoder) {
	return p.listSetTemplateSyncStatus(ctx, option, common.BKTableNameSetTemplateSyncStatus)
}

func (p *setTemplateOperation) ListSetTemplateSyncHistory(ctx core.ContextParams, option metadata.ListSetTemplateSyncStatusOption) (metadata.MultipleSetTemplateSyncStatus, errors.CCErrorCoder) {
	return p.listSetTemplateSyncStatus(ctx, option, common.BKTableNameSetTemplateSyncHistory)
}

func (p *setTemplateOperation) listSetTemplateSyncStatus(ctx core.ContextParams, option metadata.ListSetTemplateSyncStatusOption, tableName string) (metadata.MultipleSetTemplateSyncStatus, errors.CCErrorCoder) {
	result := metadata.MultipleSetTemplateSyncStatus{
		Count: 0,
		Info:  make([]metadata.SetTemplateSyncStatus, 0),
	}
	if option.BizID == 0 {
		return result, ctx.Error.CCErrorf(common.CCErrCommParamsInvalid, common.BKSetIDField)
	}
	if option.SetTemplateID == 0 {
		return result, ctx.Error.CCErrorf(common.CCErrCommParamsInvalid, common.BKSetTemplateIDField)
	}

	filter := option.ToFilter()
	querySet := p.dbProxy.Table(tableName).Find(filter)
	total, err := querySet.Count(ctx.Context)
	if err != nil {
		blog.ErrorJSON("ListSetTemplateSyncStatus failed, db count failed, filter: %s, err: %s, rid: %s", filter, err.Error(), ctx.ReqID)
		return result, ctx.Error.CCError(common.CCErrCommDBSelectFailed)
	}
	result.Count = int64(total)

	if option.Page.Start != 0 {
		querySet = querySet.Start(uint64(option.Page.Start))
	}
	if option.Page.Limit != 0 {
		querySet = querySet.Limit(uint64(option.Page.Limit))
	}
	if len(option.Page.Sort) != 0 {
		querySet = querySet.Sort(option.Page.Sort)
	}
	if err := querySet.All(ctx.Context, &result.Info); err != nil {
		blog.ErrorJSON("ListSetTemplateSyncStatus failed, db select failed, filter: %s, err: %s, rid: %s", filter, err.Error(), ctx.ReqID)
		return result, ctx.Error.CCError(common.CCErrCommDBSelectFailed)
	}

	return result, nil
}

func (p *setTemplateOperation) DeleteSetTemplateSyncStatus(ctx core.ContextParams, option metadata.DeleteSetTemplateSyncStatusOption) errors.CCErrorCoder {
	filter := map[string]interface{}{
		common.BKSetIDField: map[string]interface{}{
			common.BKDBIN: option.SetIDs,
		},
		common.BKAppIDField: option.BizID,
	}
	if err := p.dbProxy.Table(common.BKTableNameSetTemplateSyncStatus).Delete(ctx.Context, filter); err != nil {
		blog.Errorf("RemoveSetTemplateSyncStatus failed, db delete sync status failed, option: %s, err: %s, rid: %s", option, err.Error(), ctx.ReqID)
		return ctx.Error.CCError(common.CCErrCommDBUpdateFailed)
	}
	return nil
}
