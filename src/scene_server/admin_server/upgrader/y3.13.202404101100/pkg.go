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

package y3_13_202404101100

import (
	"context"

	"configcenter/src/common/blog"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
)

func init() {
	upgrader.RegistUpgrader("y3.13.202404101100", upgrade)
}

func upgrade(ctx context.Context, db dal.RDB, conf *upgrader.Config) (err error) {
	blog.Infof("start execute y3.13.202404101100")

	if err = addUnassignedCloudArea(ctx, db, conf); err != nil {
		blog.Errorf("upgrade y3.13.202404101100 add unassigned cloud area failed, err: %v", err)
		return err
	}

	if err = updateCloudAreaPlaceHolder(ctx, db, conf); err != nil {
		blog.Errorf("upgrade y3.13.202404101100 update host cloud area place holder failed, err: %v", err)
		return err
	}

	if err = changeDefaultAreaDefaultField(ctx, db, conf); err != nil {
		blog.Errorf("upgrade y3.13.202404101100 change default cloud area default field failed, err: %v", err)
		return err
	}

	blog.Infof("upgrade y3.13.202404101100 add unassigned cloud area, update host cloud area place holder, change " +
		"default cloud area default field success")
	return nil
}
