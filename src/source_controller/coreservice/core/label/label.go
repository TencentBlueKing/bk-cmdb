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
	"configcenter/src/common/http/rest"
	"configcenter/src/common/selector"
	"configcenter/src/common/util"
	"configcenter/src/source_controller/coreservice/core"
	"configcenter/src/storage/driver/mongodb"
)

type labelOperation struct {
}

// New create a new model manager instance
func New() core.LabelOperation {
	labelOps := &labelOperation{}
	return labelOps
}

func (p *labelOperation) AddLabel(kit *rest.Kit, tableName string, option selector.LabelAddOption) errors.CCErrorCoder {
	if field, err := option.Labels.Validate(); err != nil {
		blog.Infof("addLabel failed, validate failed, field:%s, err: %+v, rid: %s", field, err, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, "label."+field)
	}

	idField := common.GetInstIDField(tableName)

	// check all instance validate
	option.InstanceIDs = util.IntArrayUnique(option.InstanceIDs)
	countFilter := map[string]interface{}{
		idField: map[string]interface{}{
			common.BKDBIN: option.InstanceIDs,
		},
	}
	if count, err := mongodb.Client().Table(tableName).Find(countFilter).Count(kit.Ctx); err != nil {
		blog.ErrorJSON("AddLabel failed, db count instances failed, filter: %s, err: %s, rid: %s", countFilter, err.Error(), kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCommDBSelectFailed)
	} else if count != uint64(len(option.InstanceIDs)) {
		blog.ErrorJSON("AddLabel failed, some instance not valid, filter: %s, result count: %d, rid: %s", countFilter, count, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, "instance_ids")
	}

	for _, instanceID := range option.InstanceIDs {
		filter := map[string]interface{}{
			idField: instanceID,
		}
		data := &selector.LabelInstance{}
		if err := mongodb.Client().Table(tableName).Find(filter).One(kit.Ctx, data); err != nil {
			blog.Errorf("AddLabel failed, get instance failed, instanceID: %+v, err: %+v, rid: %s", instanceID, err, kit.Rid)
			return kit.CCError.CCErrorf(common.CCErrCommDBSelectFailed)
		}
		if data.Labels != nil {
			data.Labels.AddLabel(option.Labels)
		} else {
			data.Labels = option.Labels
		}
		if err := mongodb.Client().Table(tableName).Update(kit.Ctx, filter, data); err != nil {
			blog.Errorf("AddLabel failed, update instance failed, instanceID: %+v, err: %+v, rid: %s", instanceID, err, kit.Rid)
			return kit.CCError.CCErrorf(common.CCErrCommDBUpdateFailed)
		}
	}
	return nil
}

func (p *labelOperation) RemoveLabel(kit *rest.Kit, tableName string, option selector.LabelRemoveOption) errors.CCErrorCoder {
	idField := common.GetInstIDField(tableName)

	// check all instance validate
	option.InstanceIDs = util.IntArrayUnique(option.InstanceIDs)
	countFilter := map[string]interface{}{
		idField: map[string]interface{}{
			common.BKDBIN: option.InstanceIDs,
		},
	}
	if count, err := mongodb.Client().Table(tableName).Find(countFilter).Count(kit.Ctx); err != nil {
		blog.ErrorJSON("RemoveLabel failed, db count instances failed, filter: %s, err: %s, rid: %s", countFilter, err.Error(), kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCommDBSelectFailed)
	} else if count != uint64(len(option.InstanceIDs)) {
		blog.ErrorJSON("RemoveLabel failed, some instance not valid, filter: %s, result count: %d, rid: %s", countFilter, count, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, "instance_ids")
	}

	for _, instanceID := range option.InstanceIDs {
		filter := map[string]interface{}{
			idField: instanceID,
		}
		data := &selector.LabelInstance{}
		if err := mongodb.Client().Table(tableName).Find(filter).One(kit.Ctx, data); err != nil {
			blog.Errorf("RemoveLabel failed, get instance failed, instanceID: %+v, err: %+v, rid: %s", instanceID, err, kit.Rid)
			return kit.CCError.CCErrorf(common.CCErrCommDBSelectFailed)
		}
		if data.Labels != nil {
			data.Labels.RemoveLabel(option.Keys)
		} else {
			data.Labels = make(map[string]string)
		}
		if err := mongodb.Client().Table(tableName).Update(kit.Ctx, filter, data); err != nil {
			blog.Errorf("RemoveLabel failed, update instance failed, instanceID: %+v, err: %+v, rid: %s", instanceID, err, kit.Rid)
			return kit.CCError.CCErrorf(common.CCErrCommDBUpdateFailed)
		}
	}
	return nil
}
