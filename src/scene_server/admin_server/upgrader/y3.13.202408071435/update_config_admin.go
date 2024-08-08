/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 THL A29 Limited,
 * a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 * We undertake not to change the open source license (MIT license) applicable
 * to the current version of the project delivered to anyone in the future.
 */

package y3_13_202408071435

import (
	"context"
	"fmt"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"

	"github.com/tidwall/gjson"
)

func updateConfigAdmin(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	cond := map[string]interface{}{
		"_id": common.ConfigAdminID,
	}

	dbData := make(map[string]string)
	err := db.Table(common.BKTableNameSystem).Find(cond).Fields(common.ConfigAdminValueField).One(ctx, &dbData)
	if err != nil {
		blog.Errorf("get db config admin config failed, err: %v", err)
		return err
	}

	dbConfStr := dbData[common.ConfigAdminValueField]

	// id generator config already exists, do not need to add default value
	if gjson.Get(dbConfStr, "id_generator").Exists() {
		return nil
	}

	// update id generator step to default value 1
	updateConf := gjson.Get(fmt.Sprintf("[%s,%s]", dbConfStr, `{"id_generator":{"step":1}}`), "@join").String()
	updateData := map[string]interface{}{
		common.ConfigAdminValueField: updateConf,
		common.LastTimeField:         time.Now(),
	}

	err = db.Table(common.BKTableNameSystem).Update(ctx, cond, updateData)
	if err != nil {
		blog.Errorf("update db config admin config failed, err: %v, update config: %s", err, updateConf)
		return err
	}

	return nil
}
