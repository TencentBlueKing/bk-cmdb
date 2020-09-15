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
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
	"configcenter/src/storage/driver/mongodb"
)

func (p *setTemplateOperation) UpdateSetTemplateSyncStatus(kit *rest.Kit, setID int64, option metadata.SetTemplateSyncStatus) errors.CCErrorCoder {
	if setID != option.SetID {
		return kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKSetIDField)
	}

	filter := map[string]interface{}{
		common.BKSetIDField: setID,
	}
	if err := mongodb.Client().Table(common.BKTableNameSetTemplateSyncStatus).Upsert(kit.Ctx, filter, option); err != nil {
		blog.Errorf("UpdateSetTemplateSyncStatus failed, db upsert sync status failed, id: %d, option: %s, err: %s, rid: %s", setID, option, err.Error(), kit.Rid)
		return kit.CCError.CCError(common.CCErrCommDBUpdateFailed)
	}

	if len(option.TaskID) == 0 {
		return nil
	}

	historyFilter := map[string]interface{}{
		common.BKTaskIDField: option.TaskID,
	}
	if err := mongodb.Client().Table(common.BKTableNameSetTemplateSyncHistory).Upsert(kit.Ctx, historyFilter, option); err != nil {
		blog.Errorf("UpdateSetTemplateSyncStatus failed, db upsert sync history failed, id: %d, option: %s, err: %s, rid: %s", setID, option, err.Error(), kit.Rid)
		return kit.CCError.CCError(common.CCErrCommDBUpdateFailed)
	}

	return nil
}

func (p *setTemplateOperation) ListSetTemplateSyncStatus(kit *rest.Kit, option metadata.ListSetTemplateSyncStatusOption) (metadata.MultipleSetTemplateSyncStatus, errors.CCErrorCoder) {
	return p.listSetTemplateSyncStatus(kit, option, common.BKTableNameSetTemplateSyncStatus)
}

func (p *setTemplateOperation) ListSetTemplateSyncHistory(kit *rest.Kit, option metadata.ListSetTemplateSyncStatusOption) (metadata.MultipleSetTemplateSyncStatus, errors.CCErrorCoder) {
	return p.listSetTemplateSyncStatus(kit, option, common.BKTableNameSetTemplateSyncHistory)
}

func (p *setTemplateOperation) listSetTemplateSyncStatus(kit *rest.Kit, option metadata.ListSetTemplateSyncStatusOption, tableName string) (metadata.MultipleSetTemplateSyncStatus, errors.CCErrorCoder) {
	result := metadata.MultipleSetTemplateSyncStatus{
		Count: 0,
		Info:  make([]metadata.SetTemplateSyncStatus, 0),
	}
	if option.BizID == 0 {
		return result, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKSetIDField)
	}
	if option.SetTemplateID == 0 {
		return result, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKSetTemplateIDField)
	}

	filter := option.ToFilter()
	querySet := mongodb.Client().Table(tableName).Find(filter)
	total, err := querySet.Count(kit.Ctx)
	if err != nil {
		blog.ErrorJSON("ListSetTemplateSyncStatus failed, db count failed, filter: %s, err: %s, rid: %s", filter, err.Error(), kit.Rid)
		return result, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
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
	if err := querySet.All(kit.Ctx, &result.Info); err != nil {
		blog.ErrorJSON("ListSetTemplateSyncStatus failed, db select failed, filter: %s, err: %s, rid: %s", filter, err.Error(), kit.Rid)
		return result, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}

	return result, nil
}

func (p *setTemplateOperation) DeleteSetTemplateSyncStatus(kit *rest.Kit, option metadata.DeleteSetTemplateSyncStatusOption) errors.CCErrorCoder {
	filter := map[string]interface{}{
		common.BKSetIDField: map[string]interface{}{
			common.BKDBIN: option.SetIDs,
		},
		common.BKAppIDField: option.BizID,
	}
	if err := mongodb.Client().Table(common.BKTableNameSetTemplateSyncStatus).Delete(kit.Ctx, filter); err != nil {
		blog.Errorf("RemoveSetTemplateSyncStatus failed, db delete sync status failed, option: %s, err: %s, rid: %s", option, err.Error(), kit.Rid)
		return kit.CCError.CCError(common.CCErrCommDBUpdateFailed)
	}
	return nil
}

func (sto *setTemplateOperation) ModifySetTemplateSyncStatus(kit *rest.Kit, setID int64, sysncStatus metadata.SyncStatus) errors.CCErrorCoder {

	// 最好有前后状态对比，避免跨状态转移
	filter := map[string]interface{}{
		common.BKSetIDField: setID,
	}
	doc := map[string]interface{}{
		common.BKStatusField: sysncStatus,
	}
	// check 数据是否存在
	cnt, err := mongodb.Client().Table(common.BKTableNameSetTemplateSyncStatus).Find(filter).Count(kit.Ctx)
	if err != nil {
		blog.Errorf("ModifyStatus failed, find set template sync info error, id: %d, status: %s, err: %s, rid: %s", setID, sysncStatus, err.Error(), kit.Rid)
		return kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}
	if cnt <= 0 {
		blog.Errorf("ModifyStatus failed, not find set template sync info, id: %d, status: %s, rid: %s", setID, sysncStatus, kit.Rid)
		return kit.CCError.CCError(common.CCErrCommNotFound)
	}

	if err := mongodb.Client().Table(common.BKTableNameSetTemplateSyncStatus).Upsert(kit.Ctx, filter, doc); err != nil {
		blog.Errorf("ModifyStatus failed, db upsert sync status failed, id: %d, status: %s, err: %s, rid: %s", setID, sysncStatus, err.Error(), kit.Rid)
		return kit.CCError.CCError(common.CCErrCommDBUpdateFailed)
	}

	return nil
}
