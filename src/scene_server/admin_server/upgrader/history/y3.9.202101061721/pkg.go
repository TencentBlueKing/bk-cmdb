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

package y3_9_202101061721

import (
	"context"

	"configcenter/src/common/blog"
	"configcenter/src/scene_server/admin_server/upgrader/history"
	"configcenter/src/storage/dal"
)

func init() {
	history.RegistUpgrader("y3.9.202101061721", upgrade)
}

func upgrade(ctx context.Context, db dal.RDB, conf *history.Config) (err error) {
	blog.Infof("start execute y3.9.202101061721")

	err = delPreviousDelArchiveData(ctx, db, conf)
	if err != nil {
		blog.Errorf("[upgrade y3.9.202101061721] delete previous del archive data failed, error %v", err)
		return err
	}

	err = addDelArchiveIndex(ctx, db, conf)
	if err != nil {
		blog.Errorf("[upgrade y3.9.202101061721] add del archive table index failed, error %v", err)
		return err
	}

	return nil
}