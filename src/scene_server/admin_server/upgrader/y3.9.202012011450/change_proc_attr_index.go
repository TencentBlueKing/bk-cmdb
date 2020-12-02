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

package y3_9_202012011450

import (
	"context"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
)

// changeProcAttrIndex 调整进程属性index
func changeProcAttrIndex(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	// 切片里的元素下标顺序代表调整后的前端展示顺序，其中auto_start属性的bk_isapi为true，在前端不会展示，放在最后
	attrs := []string{
		"work_path",
		"pid_file",
		"user",
		"proc_num",
		"priority",
		"start_cmd",
		"stop_cmd",
		"restart_cmd",
		"face_stop_cmd",
		"reload_cmd",
		"bk_start_check_secs",
		"timeout",
		"auto_start",
	}

	filter := map[string]interface{}{
		common.BKObjIDField: common.BKInnerObjIDProc,
	}

	doc := map[string]interface{}{
		common.BKPropertyIndexField: 0,
	}

	for idx, attr := range attrs {
		filter[common.BKPropertyIDField] = attr
		doc[common.BKPropertyIndexField] = idx + 1
		if err := db.Table(common.BKTableNameObjAttDes).Update(ctx, filter, doc); err != nil {
			blog.Errorf("update failed, filter:%#v, doc:%#v, err:%v", filter, doc, err)
			return err
		}
	}

	return nil
}

// changeProcAttrOption 调整进程属性bk_start_check_secs的option
func changeProcAttrOption(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	filter := map[string]interface{}{
		common.BKObjIDField:      common.BKInnerObjIDProc,
		common.BKPropertyIDField: "bk_start_check_secs",
	}

	doc := map[string]interface{}{
		common.BKOptionField: metadata.IntOption{Min: "0", Max: "600"},
	}

	if err := db.Table(common.BKTableNameObjAttDes).Update(ctx, filter, doc); err != nil {
		blog.Errorf("update failed, filter:%#v, doc:%#v, err:%v", filter, doc, err)
		return err
	}

	return nil
}
