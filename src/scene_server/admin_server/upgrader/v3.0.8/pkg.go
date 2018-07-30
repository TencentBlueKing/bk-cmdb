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

package v3v0v8

import (
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage"
)

func init() {
	upgrader.RegistUpgrader("v3.0.8", upgrade)
}

func upgrade(db storage.DI, conf *upgrader.Config) (err error) {
	err = createTable(db, conf)
	if err != nil {
		return err
	}
	err = addPresetObjects(db, conf)
	if err != nil {
		return err
	}
	err = addPlatData(db, conf)
	if err != nil {
		return err
	}
	err = addSystemData(db, conf)
	if err != nil {
		return err
	}
	err = addDefaultBiz(db, conf)
	if err != nil {
		return err
	}
	err = addBKApp(db, conf)
	if err != nil {
		return err
	}

	return
}
