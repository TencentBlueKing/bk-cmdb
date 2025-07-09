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

package y3_9_202008121631

import (
	"context"

	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
)

const currentPackageName = "y3.9.202008121631"

func init() {
	upgrader.RegistUpgrader(currentPackageName, upgrade)
}

// upgrade add field 'bk_ishidden' to object process and plat, and the value is true
func upgrade(ctx context.Context, db dal.RDB, conf *upgrader.Config) (err error) {
	blog.Infof("start execute %s", currentPackageName)

	if err := addObjectFieldIsHidden(ctx, db, conf); err != nil {
		blog.Errorf("[%s] failed to object add field %s, error: %s", currentPackageName,
			metadata.ModelFieldIsHidden, err.Error())
		return err
	}

	return nil
}
