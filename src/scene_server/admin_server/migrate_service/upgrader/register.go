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

package upgrader

import (
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/storage"
	"sort"
	"sync"
)

type Ctx struct{}

type Version struct {
	version string // v3.0.8-beta.11
	do      func(storage.DI, *Ctx) error
}

var versionPool = []Version{}

var registlock sync.Mutex

func RegistVersion(version string, handlerFunc func(storage.DI, *Ctx) error) {
	registlock.Lock()
	defer registlock.Unlock()
	v := Version{version: version, do: handlerFunc}
	versionPool = append(versionPool, v)
}

// Upgrade uprade the db datas to newest verison
func Upgrade(db storage.DI) (err error) {
	sort.Slice(versionPool, func(i, j int) bool {
		return versionPool[i].version > versionPool[j].version
	})

	for _, v := range versionPool {
		err = v.do(db, nil)
		if err != nil {
			blog.Errorf("upgrade version %s error: %s", v.version, err.Error())
			return err
		}
		err = saveVesion(db, v.version)
		if err != nil {
			blog.Errorf("save version %s error: %s", v.version, err.Error())
			return err
		}
		blog.Info("upgrade version success")
	}
	return nil
}

func saveVesion(db storage.DI, version string) error {
	condition := map[string]interface{}{
		"type": "version",
	}
	data := map[string]interface{}{
		"type":            "version",
		"current_version": version,
	}
	count, err := db.GetCntByCondition(common.SystemTableName, condition)
	if err != nil {
		return err
	}
	if count > 0 {
		_, err = db.Insert(common.SystemTableName, data)
		if err != nil {
			return err
		}
	}

	return db.UpdateByCondition(common.SystemTableName, data, condition)
}
