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

package hooks

import (
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/storage/dal"
)

// SetVIPInfoForProcessHook if query fields contains vip info, set vip info for processes
func SetVIPInfoForProcessHook(kit *rest.Kit, processes []mapstr.MapStr, fields []string, table string, db dal.DB) (
	[]mapstr.MapStr, error) {

	return processes, nil
}

// ParseVIPFieldsForProcessHook parse process vip fields for process
func ParseVIPFieldsForProcessHook(fields []string, table string) ([]string, []string) {

	return fields, make([]string, 0)
}

// UpdateProcessBindInfoHook if process need to update bind info, only update the specified fields
func UpdateProcessBindInfoHook(kit *rest.Kit, objID string, origin mapstr.MapStr, data mapstr.MapStr) error {
	return nil
}
