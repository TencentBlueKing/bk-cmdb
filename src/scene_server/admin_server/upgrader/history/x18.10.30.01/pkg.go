// Package x18_10_30_01 TODO
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
package x18_10_30_01

import (
	"context"

	"configcenter/src/common/blog"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
)

func init() {
	upgrader.RegistUpgrader("x18.10.30.01", upgrade)
}
func upgrade(ctx context.Context, db dal.RDB, conf *upgrader.Config) (err error) {
	err = createAssociationTable(ctx, db, conf)
	if err != nil {
		blog.Errorf("[upgrade x18.10.30.01] createAssociationTable error  %s", err.Error())
		return err
	}
	err = createInstanceAssociationIndex(ctx, db, conf)
	if err != nil {
		blog.Errorf("[upgrade x18.10.30.01] createInstanceAssociationIndex error  %s", err.Error())
		return err
	}
	err = addPresetAssociationType(ctx, db, conf)
	if err != nil {
		blog.Errorf("[upgrade x18.10.30.01] addPresetAssociationType error  %s", err.Error())
		return err
	}
	err = reconcilAsstData(ctx, db, conf)
	if err != nil {
		blog.Errorf("[upgrade x18.10.30.01] reconcilAsstData error  %s", err.Error())
		return err
	}
	return
}
