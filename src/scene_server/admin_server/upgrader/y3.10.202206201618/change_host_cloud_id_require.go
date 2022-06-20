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

package y3_10_202206201618

import (
	"context"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
)

func changeHostCloudIDRequire(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	cond := map[string]interface{}{
		common.BKObjIDField:      common.BKInnerObjIDHost,
		common.BKPropertyIDField: common.BKCloudIDField,
	}

	count, err := db.Table(common.BKTableNameObjAttDes).Find(cond).Count(ctx)
	if err != nil {
		blog.Errorf("find host attribute bk_cloud_id failed, filter: %v, err: %v", cond, err)
		return err
	}

	if count != 1 {
		blog.Errorf("find host attribute bk_cloud_id not only one, filter: %v, err: %v", cond, err)
		return err
	}

	filter := map[string]interface{}{
		common.BKObjIDField:      common.BKInnerObjIDHost,
		common.BKPropertyIDField: common.BKCloudIDField,
	}

	data := mapstr.MapStr{common.BKIsRequiredField: true}

	if err := db.Table(common.BKTableNameObjAttDes).Update(ctx, filter, data); err != nil {
		blog.Errorf("change bk_cloud_id's isrequired to true failed, filter: %v, data: %v, err: %v", filter, data, err)
		return err
	}
	return nil
}
