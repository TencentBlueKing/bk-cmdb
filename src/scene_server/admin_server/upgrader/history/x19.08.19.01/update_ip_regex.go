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

package x19_08_19_01

import (
	"context"

	"configcenter/src/common"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
)

var IPRegex = `^((1?\d{1,2}|2[0-4]\d|25[0-5])[.]){3}(1?\d{1,2}|2[0-4]\d|25[0-5])(,((1?\d{1,2}|2[0-4]\d|25[0-5])[.]){3}(1?\d{1,2}|2[0-4]\d|25[0-5]))*$`

func updateIPRegex(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	// host
	filter := map[string]interface{}{
		common.BKPropertyIDField: map[string]interface{}{
			common.BKDBIN: []string{common.BKHostInnerIPField, common.BKHostOuterIPField},
		},
		common.BKObjIDField: common.BKInnerObjIDHost,
	}
	doc := map[string]interface{}{
		common.BKOptionField: IPRegex,
	}
	err := db.Table(common.BKTableNameObjAttDes).Update(ctx, filter, doc)
	if err != nil {
		return err
	}

	// bk-switch
	switchFilter := map[string]interface{}{
		common.BKPropertyIDField: map[string]interface{}{
			common.BKDBIN: []string{"bk_admin_ip"},
		},
		common.BKObjIDField: common.BKInnerObjIDSwitch,
	}
	switchDoc := map[string]interface{}{
		common.BKOptionField: IPRegex,
	}
	err = db.Table(common.BKTableNameObjAttDes).Update(ctx, switchFilter, switchDoc)
	if err != nil {
		return err
	}
	return nil
}
