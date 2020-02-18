/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package y3_9_202002181444

import (
	"configcenter/src/common"
	"configcenter/src/common/mapstr"
	"context"
	"fmt"

	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
)

func removeOldPlatAttrs(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	cond := mapstr.MapStr{
		common.BKObjIDField: common.BKInnerObjIDPlat,
	}

	err := db.Table(common.BKTableNameObjAttDes).Delete(ctx, cond)
	if err != nil {
		return fmt.Errorf("upgrade y3.9.202002181444, remove old plat attrs failed, err: %v", err)
	}
	return nil
}
