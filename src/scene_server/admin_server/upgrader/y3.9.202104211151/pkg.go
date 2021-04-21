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

package y3_9_202104211151

import (
	"context"

	"configcenter/src/common/blog"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
)

func init() {
	upgrader.RegistUpgrader("y3.9.202104211151", upgrade)
}

/*
	This upgrader is allowed to create sets with the same name in different custom levels
	Delete unique index 'idx_unique_bizID_setName'
	在不同自定义层级下允许创建同名集群,删除唯一索引'idx_unique_bizID_setName'
*/
func upgrade(ctx context.Context, db dal.RDB, conf *upgrader.Config) (err error) {
	blog.Infof("start execute y3.9.202104211151，remove set name unique with biz index, and add set name unique with parent id index")

	err = changeSetUniqueIndex(ctx, db, conf)
	if err != nil {
		blog.Errorf("[upgrade y3.9.202104211151] change unique index failed, err: %v", err)
		return err
	}
	return nil
}
