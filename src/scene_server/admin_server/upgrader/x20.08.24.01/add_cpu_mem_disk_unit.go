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

package x20_08_24_01

import (
	"context"
	"fmt"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
)

//add unit to bk_property_name.
func addCpuMemDiskUnit(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	var err error

	//cpu
	filter := map[string]interface{}{
		common.BKObjIDField:      common.BKInnerObjIDHost,
		common.BKPropertyIDField: "bk_cpu_mhz",
	}
	doc := map[string]interface{}{
		common.BKPropertyNameField: "CPU频率(MHz)",
	}
	if err = db.Table(common.BKTableNameObjAttDes).Update(ctx, filter, doc); err != nil {
		blog.Errorf("addCpuMemDiskUnit(bk_cpu_mhz) failed, err: %+v", err)
		return fmt.Errorf("addCpuMemDiskUnit(bk_cpu_mhz) failed, err: %v", err)
	}

	//mem
	filter[common.BKPropertyIDField] = "bk_mem"
	doc[common.BKPropertyNameField] = "内存容量(MB)"
	if err = db.Table(common.BKTableNameObjAttDes).Update(ctx, filter, doc); err != nil {
		blog.Errorf("addCpuMemDiskUnit(bk_mem) failed, err: %+v", err)
		return fmt.Errorf("addCpuMemDiskUnit(bk_mem) failed, err: %v", err)
	}

	//disk
	filter[common.BKPropertyIDField] = "bk_disk"
	doc[common.BKPropertyNameField] = "磁盘容量(GB)"
	if err = db.Table(common.BKTableNameObjAttDes).Update(ctx, filter, doc); err != nil {
		blog.Errorf("addCpuMemDiskUnit(bk_disk) failed, err: %+v", err)
		return fmt.Errorf("addCpuMemDiskUnit(bk_disk) failed, err: %v", err)
	}
	return nil
}
