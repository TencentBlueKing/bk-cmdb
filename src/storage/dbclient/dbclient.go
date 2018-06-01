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
 
package dbclient

import (
	"configcenter/src/storage"
	"configcenter/src/storage/mgoclient"
	"configcenter/src/storage/redisclient"
)

// NewDB return DI instance
func NewDB(host, port, usr, pwd, mechanism, database, driverType string) (storage.DI, error) {
	if driverType == storage.DI_MONGO {
		db, err := mgoclient.NewMgoCli(host, port, usr, pwd, mechanism, database)
		if err == nil {
			return db, err
		}
		return db, err
	} else if driverType == storage.DI_REDIS {
		db, err := redisclient.NewRedis(host, port, usr, pwd, database)
		if err == nil {
			return db, err
		}
	}
	db, err := mgoclient.NewMgoCli(host, port, usr, pwd, mechanism, database)
	return db, err
}
