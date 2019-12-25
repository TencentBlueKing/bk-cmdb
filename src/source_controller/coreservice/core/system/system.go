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
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/metadata"
	"configcenter/src/source_controller/coreservice/core"
	"configcenter/src/storage/dal"
)

var _ core.SystemOperation = (*systemManager)(nil)

type systemManager struct {
	dbProxy dal.RDB
}

// New create a new instance manager instance
func New(dbProxy dal.RDB) core.SystemOperation {
	return &systemManager{
		dbProxy: dbProxy,
	}
}

func (sm *systemManager) GetSystemUserConfig(ctx core.ContextParams) (map[string]interface{}, errors.CCErrorCoder) {
	cond := map[string]string{"type": metadata.CCSystemUserConfigSwitch}
	result := make(map[string]interface{}, 0)
	err := sm.dbProxy.Table(common.BKTableNameSystem).Find(cond).One(ctx, &result)
	if err != nil && !sm.dbProxy.IsNotFoundError(err) {
		blog.ErrorJSON("GetSystemUserConfig find error. cond:%s, err:%s, rid:%s", cond, err.Error(), ctx.ReqID)
		return nil, ctx.Error.CCError(common.CCErrCommDBSelectFailed)
	}

	return result, nil
}
