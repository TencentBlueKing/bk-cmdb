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

package x19_09_03_05

import (
	"context"
	"fmt"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/dal/types"
)

func CreateObjectIDIndex(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	idx := types.Index{
		Keys: map[string]int32{
			common.BKObjIDField: 1,
		},
		Name:       common.BKObjIDField,
		Unique:     false,
		Background: false,
	}

	if err := db.Table(common.BKTableNameObjUnique).CreateIndex(ctx, idx); err != nil {
		blog.Errorf("CreateObjectIDIndex failed, err: %+v", err)
		return fmt.Errorf("CreateObjectIDIndex failed, err: %v", err)
	}
	return nil
}
