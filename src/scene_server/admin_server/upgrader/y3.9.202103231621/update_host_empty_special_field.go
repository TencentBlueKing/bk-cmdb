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

package y3_9_202103231621

import (
	"context"
	"fmt"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
)

// updateHostEmptySpecialField update host special field's empty value to nil so that unique index will ignore them
func updateHostEmptySpecialField(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	specialFields := []string{common.BKHostInnerIPField, common.BKHostOuterIPField, common.BKOperatorField,
		common.BKBakOperatorField}

	filters := make([]map[string]interface{}, len(specialFields))
	for index, field := range specialFields {
		filters[index] = map[string]interface{}{
			field: map[string]interface{}{common.BKDBSize: 0},
		}
	}
	filter := map[string]interface{}{
		common.BKDBOR: filters,
	}

	for {
		hosts := make([]metadata.HostMapStr, 0)
		fields := append(specialFields, common.BKHostIDField)

		if err := db.Table(common.BKTableNameBaseHost).Find(filter).Fields(fields...).Limit(common.BKMaxPageSize).
			All(ctx, &hosts); err != nil {

			blog.Errorf("find hosts failed, filter: %#v, err: %v", filter, err)
			return err
		}

		for _, host := range hosts {
			updateHost := make(map[string]interface{})
			for _, field := range specialFields {
				if host[field] == nil {
					continue
				}

				value, ok := host[field].(string)
				if !ok {
					return fmt.Errorf("host(%#v) special field %s is not valid", host, field)
				}

				if len(value) == 0 {
					updateHost[field] = nil
				}
			}

			if len(updateHost) == 0 {
				continue
			}

			filter := map[string]interface{}{
				common.BKHostIDField: host[common.BKHostIDField],
			}
			if err := db.Table(common.BKTableNameBaseHost).Update(ctx, filter, updateHost); err != nil {
				blog.ErrorJSON("update host failed, filter: %s, data: %s, err: %s", filter, updateHost, err)
				return err
			}
		}

		if len(hosts) < common.BKMaxPageSize {
			return nil
		}
	}
	return nil
}
