/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 Tencent. All rights reserved.
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

package y3_14_202603231000

import (
	"context"

	"configcenter/src/common/blog"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
)

func init() {
	upgrader.RegistUpgrader("y3.14.202603231000", upgrade)
}

func upgrade(ctx context.Context, db dal.RDB, conf *upgrader.Config) (err error) {
	blog.Infof("start execute y3.14.202603231000, init custom resource workload collection and index")

	if err = addCustomResourceCollection(ctx, db); err != nil {
		blog.Errorf("upgrade y3.14.202603231000 add custom resource collection failed, err: %v", err)
		return err
	}

	if err = addCustomResourceCollectionIndex(ctx, db); err != nil {
		blog.Errorf("upgrade y3.14.202603231000 add custom resource collection index failed, err: %v", err)
		return err
	}

	blog.Infof("upgrade y3.14.202603231000 init custom resource workload success")
	return nil
}
