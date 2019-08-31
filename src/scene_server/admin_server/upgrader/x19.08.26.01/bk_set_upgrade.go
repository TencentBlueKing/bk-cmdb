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

package x19_08_26_01

import (
	"context"
	"errors"
	"fmt"

	"configcenter/src/common"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
)

// change set name of 集成平台 to PaaS平台 in blueking business
type bizSimple struct {
	ID int64 `bson:"bk_biz_id"`
}

func changeIntegrationPlatToPaasPlat(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	cond := map[string]interface{}{
		"bk_biz_name": "蓝鲸",
	}
	biz := make([]bizSimple, 0)
	if err := db.Table(common.BKTableNameBaseApp).Find(cond).All(context.Background(), &biz); err != nil {
		return fmt.Errorf("upgrade x19_08_26_01, but get blueking business id failed, err: %v", err)
	}

	if len(biz) == 0 {
		return errors.New("upgrade x19_08_26_01, but can not find blueking business ")
	}

	if len(biz) >= 2 {
		return errors.New("upgrade x19_08_26_01, but got multiple blueking business ")
	}

	setFilter := map[string]interface{}{
		"bk_biz_id":   biz[0].ID,
		"bk_set_name": "集成平台",
	}
	newName := map[string]interface{}{
		"bk_set_name": "PaaS平台",
	}

	if err := db.Table(common.BKTableNameBaseSet).Update(context.Background(), setFilter, newName); err != nil {
		return fmt.Errorf("upgrade x19_08_26_01, but update set name failed, err: %v", err)
	}

	return nil
}
