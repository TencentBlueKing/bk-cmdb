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

// Package system TODO
package system

import (
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/source_controller/coreservice/core"
	"configcenter/src/storage/driver/mongodb"
)

var _ core.SystemOperation = (*systemManager)(nil)

type systemManager struct {
}

// New create a new instance manager instance
func New() core.SystemOperation {
	return &systemManager{}
}

// GetSystemUserConfig TODO
func (sm *systemManager) GetSystemUserConfig(kit *rest.Kit) (map[string]interface{}, errors.CCErrorCoder) {
	cond := map[string]string{"type": metadata.CCSystemUserConfigSwitch}
	result := make(map[string]interface{}, 0)
	err := mongodb.Shard(kit.SysShardOpts()).Table(common.BKTableNameSystem).Find(cond).One(kit.Ctx, &result)
	if err != nil && !mongodb.IsNotFoundError(err) {
		blog.Errorf("find system user config failed. cond: %v, err: %v, rid: %s", cond, err, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}

	return result, nil
}

// SearchGlobalSettingConfig search global setting.
func (sm *systemManager) SearchGlobalSettingConfig(kit *rest.Kit, options *metadata.GlobalConfOptions) (
	*metadata.GlobalSettingConfig, errors.CCErrorCoder) {

	ret := new(metadata.GlobalSettingConfig)
	err := mongodb.Shard(kit.ShardOpts()).Table(common.BKTableNameGlobalConfig).Find(mapstr.MapStr{}).Fields(
		options.Fields...).One(kit.Ctx, ret)
	if err != nil {
		blog.Errorf("search platform setting failed, err: %v, rid: %s", err, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}

	return ret, nil
}

// UpdatePlatformSettingConfig update platform setting.
func (sm *systemManager) UpdatePlatformSettingConfig(kit *rest.Kit, input mapstr.MapStr,
	typeId string) errors.CCErrorCoder {

	if _, ok := input[typeId]; !ok {
		blog.Errorf("type %s is not exist, typeId: %s, rid: %s", typeId, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, "type %s is not exist", typeId)
	}
	data := map[string]interface{}{
		typeId:               input[typeId],
		common.LastTimeField: time.Now(),
	}
	err := mongodb.Shard(kit.ShardOpts()).Table(common.BKTableNameGlobalConfig).Update(kit.Ctx, mapstr.MapStr{},
		data)
	if err != nil {
		blog.Errorf("update global config %s failed, update: %v, err: %v, rid: %s", typeId, data, err, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCommDBUpdateFailed, err)
	}

	return nil
}
