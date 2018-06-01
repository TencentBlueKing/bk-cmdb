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
 
package logics

import (
	"configcenter/src/common"
	"configcenter/src/common/core/cc/api"
	"configcenter/src/scene_server/admin_server/migrateregister"

	_ "configcenter/src/scene_server/admin_server/migrate_service/logics/collections"
	_ "configcenter/src/scene_server/admin_server/migrate_service/logics/event"
	_ "configcenter/src/scene_server/admin_server/migrate_service/logics/host"
	_ "configcenter/src/scene_server/admin_server/migrate_service/logics/obj"
	_ "configcenter/src/scene_server/admin_server/migrate_service/logics/plat"
	_ "configcenter/src/scene_server/admin_server/migrate_service/logics/privilege"
	_ "configcenter/src/scene_server/admin_server/migrate_service/logics/proc"
	_ "configcenter/src/scene_server/admin_server/migrate_service/logics/topo"
)

func DBMigrate(ownerid string) error {
	a := api.GetAPIResource()
	if "" == ownerid {
		ownerid = common.BKDefaultOwnerID
	}
	createTable := migrateregister.GetMigrateAction(migrateregister.MigrateTypeCreateTable)
	for _, f := range createTable {
		err := f(ownerid, a.InstCli, a.InstCli)
		if nil != err {
			return err
		}
	}

	alterTable := migrateregister.GetMigrateAction(migrateregister.MigrateTypeAlterTable)
	for _, f := range alterTable {
		err := f(ownerid, a.InstCli, a.InstCli)
		if nil != err {
			return err
		}
	}

	dropTable := migrateregister.GetMigrateAction(migrateregister.MigrateTypeDropTable)
	for _, f := range dropTable {
		err := f(ownerid, a.InstCli, a.InstCli)
		if nil != err {
			return err
		}
	}

	addData := migrateregister.GetMigrateAction(migrateregister.MigrateTypeAddData)
	for _, f := range addData {
		err := f(ownerid, a.InstCli, a.InstCli)
		if nil != err {
			return err
		}
	}

	updateData := migrateregister.GetMigrateAction(migrateregister.MigrateTypeUpdateData)
	for _, f := range updateData {
		err := f(ownerid, a.InstCli, a.InstCli)
		if nil != err {
			return err
		}
	}

	delData := migrateregister.GetMigrateAction(migrateregister.MigrateTypeDelData)
	for _, f := range delData {
		err := f(ownerid, a.InstCli, a.InstCli)
		if nil != err {
			return err
		}
	}
	return nil
}
