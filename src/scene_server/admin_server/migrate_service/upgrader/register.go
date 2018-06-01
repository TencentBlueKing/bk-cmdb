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
	ccversion "configcenter/src/common/version"
	"configcenter/src/storage"
	"sort"
	"sync"
)

// Config config for upgrader
type Config struct {
	OwnerID string
}

// Upgrader define a version upgrader
type Upgrader struct {
	version string // v3.0.8-beta.11
	do      func(storage.DI, *Config) error
}

var upgraderPool = []Upgrader{}

var registlock sync.Mutex

// RegistUpgrader register upgrader
func RegistUpgrader(version string, handlerFunc func(storage.DI, *Config) error) {
	registlock.Lock()
	defer registlock.Unlock()
	v := Upgrader{version: version, do: handlerFunc}
	upgraderPool = append(upgraderPool, v)
}

// Upgrade uprade the db datas to newest verison
func Upgrade(db storage.DI) (err error) {
	sort.Slice(upgraderPool, func(i, j int) bool {
		return upgraderPool[i].version > upgraderPool[j].version
	})

	for _, v := range upgraderPool {
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
		"distro":          ccversion.CCDistro,
		"distro_version":  ccversion.CCDistroVersion,
	}
	count, err := db.GetCntByCondition(common.SystemTableName, condition)
	if err != nil {
		return err
	}
	if count > 0 {
		data["init_version"] = version
		data["init_distro_version"] = ccversion.CCDistroVersion
		_, err = db.Insert(common.SystemTableName, data)
		if err != nil {
			return err
		}
	}

	return db.UpdateByCondition(common.SystemTableName, data, condition)
}
