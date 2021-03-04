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

package y3_9_202103031533

import (
	"context"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
)

// removeHostCPUMhzField remove host bk_cpu_mhz field in object attribute table and from host instances
func removeHostCPUMhzField(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	// remove bk_cpu_mhz attributes
	attrFilter := map[string]interface{}{
		common.BKObjIDField:      common.BKInnerObjIDHost,
		common.BKPropertyIDField: "bk_cpu_mhz",
		common.BKAppIDField:      0,
	}
	if err := db.Table(common.BKTableNameObjAttDes).Delete(ctx, attrFilter); err != nil {
		blog.Errorf("remove host object bk_cpu_mhz attribute failed, err: %v", err)
		return err
	}

	// remove all host instances' bk_cpu_mhz fields
	hostFilter := map[string]interface{}{
		"bk_cpu_mhz": map[string]interface{}{common.BKDBExists: true},
	}

	for {
		hosts := make([]map[string]interface{}, 0)
		err := db.Table(common.BKTableNameBaseHost).Find(hostFilter).Start(0).Limit(common.BKMaxPageSize).
			Fields(common.BKHostIDField).All(ctx, &hosts)
		if err != nil {
			blog.Errorf("get host ids to remove fields failed, err: %v", err)
			return err
		}

		if len(hosts) == 0 {
			break
		}

		hostIDs := make([]int64, len(hosts))
		for index, host := range hosts {
			hostID, err := util.GetInt64ByInterface(host[common.BKHostIDField])
			if err != nil {
				blog.Errorf("get host id failed, host: %+v, err: %v", host, err)
				return err
			}
			hostIDs[index] = hostID
		}

		hostFilter := map[string]interface{}{
			common.BKHostIDField: map[string]interface{}{
				common.BKDBIN: hostIDs,
			},
		}

		if err := db.Table(common.BKTableNameBaseHost).DropColumns(ctx, hostFilter, []string{"bk_cpu_mhz"}); err != nil {
			blog.Errorf("remove host instances(%+v) bk_cpu_mhz field failed, err: %v", hostIDs, err)
			return err
		}

		if len(hosts) < common.BKMaxPageSize {
			break
		}
	}
	return nil
}
