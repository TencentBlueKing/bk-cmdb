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

package y3_9_202002131522

import (
	"context"
	"fmt"

	"configcenter/src/common"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
)

func updatePlatUniqueKey(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	attr := Attribute{}
	cond := mapstr.MapStr{
		common.BKPropertyIDField: common.BKCloudNameField,
	}
	if err := db.Table(common.BKTableNameObjAttDes).Find(cond).One(ctx, &attr); err != nil {
		return fmt.Errorf("upgrade y3.9.202002131522, but update plat unique key failed, err: %v", err)
	}

	uniques := metadata.ObjectUnique{
		ObjID:     common.BKInnerObjIDPlat,
		MustCheck: true,
		Keys: []metadata.UniqueKey{
			{
				Kind: metadata.UniqueKeyKindProperty,
				ID:   uint64(attr.ID),
			},
		},
		Ispre:    true,
		OwnerID:  conf.OwnerID,
		LastTime: metadata.Now(),
	}

	opt := mapstr.MapStr{
		common.BKObjIDField: common.BKInnerObjIDPlat,
	}
	if err := db.Table(common.BKTableNameObjUnique).Update(ctx, opt, uniques); err != nil {
		return fmt.Errorf("upgrade y3.9.202002131522, but update plat unique key failed, err: %v", err)
	}

	return nil
}
