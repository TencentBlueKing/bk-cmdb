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

package y3_10_202104221702

import (
	"context"

	"configcenter/src/common/blog"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
)

func init() {
	upgrader.RegistUpgrader("y3.10.202104221702", upgrade)
}

func upgrade(ctx context.Context, db dal.RDB, conf *upgrader.Config) (err error) {
	blog.Info("y3.10.202104221702")

	if err := dropPlatVpcIDIndex(ctx, db, conf); err != nil {
		blog.Errorf("[upgrade y3.10.202104221702] migrate cloud unique index error, error:%s",
			err.Error())
		return err
	}

	if err := instanceObjectIDMapping(ctx, db, conf); err != nil {
		blog.Errorf("[upgrade y3.10.202104221702] migrate instance object id mapping table fail, error:%s",
			err.Error())
		return err
	}

	err = splitTable(ctx, db, conf)
	if err != nil {
		blog.Errorf("[upgrade y3.10.202104221702] migrate inst split table failed, error  %s", err.Error())
		return err
	}

	if err = syncInnerObjectIndex(ctx, db, conf); err != nil {
		blog.Errorf("[upgrade y3.10.202104221702] migrate inner object index failed, error  %s", err.Error())
		return err
	}

	return nil
}
