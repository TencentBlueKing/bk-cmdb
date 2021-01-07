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

package x08_09_11_01

import (
	"context"

	"configcenter/src/common"
	"configcenter/src/common/mapstr"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
)

func addOperationLogIndex(ctx context.Context, db dal.RDB, conf *upgrader.Config) (err error) {

	index := dal.Index{Name: "", Keys: mapstr.MapStr{"opt_time": 1}, Background: true}

	if err = db.Table(common.BKTableNameOperationLog).CreateIndex(ctx, index); err != nil && db.IsDuplicatedError(err) {
		return err
	}

	return nil
}
