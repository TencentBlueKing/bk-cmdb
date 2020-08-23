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

package system

import (
	"encoding/json"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
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

func (sm *systemManager) GetSystemUserConfig(kit *rest.Kit) (map[string]interface{}, errors.CCErrorCoder) {
	cond := map[string]string{"type": metadata.CCSystemUserConfigSwitch}
	result := make(map[string]interface{}, 0)
	err := mongodb.Client().Table(common.BKTableNameSystem).Find(cond).One(kit.Ctx, &result)
	if err != nil && !mongodb.Client().IsNotFoundError(err) {
		blog.ErrorJSON("GetSystemUserConfig find error. cond:%s, err:%s, rid:%s", cond, err.Error(), kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}

	return result, nil
}

func (sm *systemManager) SearchConfigAdmin(kit *rest.Kit) (*metadata.ConfigAdmin, errors.CCErrorCoder) {
	cond := map[string]interface{}{
		"_id": common.ConfigAdminID,
	}

	ret := struct {
		Config string `json:"config"`
	}{}
	err := mongodb.Client().Table(common.BKTableNameSystem).Find(cond).Fields(common.ConfigAdminValueField).One(kit.Ctx, &ret)
	if err != nil {
		blog.Errorf("SearchConfigAdmin failed, err: %+v, rid: %s", err, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}
	conf := new(metadata.ConfigAdmin)
	if err := json.Unmarshal([]byte(ret.Config), conf); err != nil {
		blog.Errorf("SearchConfigAdmin failed, Unmarshal err: %v, config:%+v,rid:%s", err, ret.Config, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrCommJSONUnmarshalFailed)
	}

	return conf, nil
}
