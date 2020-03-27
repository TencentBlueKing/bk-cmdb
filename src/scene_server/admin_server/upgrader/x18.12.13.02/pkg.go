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

package x18_12_13_02

import (
	"context"

	"configcenter/src/common/blog"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
)

func init() {
	upgrader.RegistUpgrader("x18_12_13_02", upgrade)
}

func upgrade(ctx context.Context, db dal.RDB, conf *upgrader.Config) (err error) {

	err = addCloudTaskTable(ctx, db, conf)
	if err != nil {
		blog.Errorf("[upgrade x18_12_13_02] addCloudTaskTable error  %s", err.Error())
		return err
	}

	err = addCloudResourceConfirmTable(ctx, db, conf)
	if err != nil {
		blog.Errorf("[upgrade x18_12_13_02] addCloudResourceConfirmTable error  %s", err.Error())
		return err
	}

	err = addCloudSyncHistoryTable(ctx, db, conf)
	if err != nil {
		blog.Errorf("[upgrade x18_12_13_02] addCloudSyncHistoryTable error  %s", err.Error())
		return err
	}

	err = addCloudConfirmHistoryTable(ctx, db, conf)
	if err != nil {
		blog.Errorf("[upgrade x18_12_13_02] addCloudConfirmHistoryTable error  %s", err.Error())
		return err
	}

	return
}
