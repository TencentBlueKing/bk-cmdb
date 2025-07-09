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

package y3_10_202107311844

import (
	"context"

	"configcenter/src/common"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
)

// setProcBindIPNotRequired update process bind info attribute, set bind ip not required and allows empty value
func setProcBindIPNotRequired(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	bindIPFilter := map[string]interface{}{
		common.BKObjIDField:                                   common.BKInnerObjIDProc,
		common.BKPropertyIDField:                              common.BKProcBindInfo,
		common.BKOptionField + "." + common.BKPropertyIDField: common.BKIP,
	}

	doc := map[string]interface{}{
		common.LastTimeField:  metadata.Now(),
		"option.$.isrequired": false,
	}
	return db.Table(common.BKTableNameObjAttDes).Update(ctx, bindIPFilter, doc)
}
