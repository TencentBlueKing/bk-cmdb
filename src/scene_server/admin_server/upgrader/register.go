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
	"sort"
	"sync"

	"gopkg.in/mgo.v2"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	ccversion "configcenter/src/common/version"
	"configcenter/src/storage"
)

// Config config for upgrader
type Config struct {
	OwnerID    string
	SupplierID int
	User       string
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
	// blog.Infof("registed upgrader for version %s", v.version)
}

// Upgrade uprade the db datas to newest verison
func Upgrade(db storage.DI, conf *Config) (err error) {
	sort.Slice(upgraderPool, func(i, j int) bool {
		return upgraderPool[i].version < upgraderPool[j].version
	})

	cmdbVision, err := getVersion(db)
	if err != nil {
		return err
	}
	cmdbVision.Distro = ccversion.CCDistro
	cmdbVision.DistroVersion = ccversion.CCDistroVersion

	currentVision := cmdbVision.CurrentVersion
	lastVersion := ""
	for _, v := range upgraderPool {
		lastVersion = v.version
		if v.version <= currentVision {
			blog.Infof(`currentVision is "%s" skip upgrade "%s"`, currentVision, v.version)
			continue
		}
		err = v.do(db, conf)
		if err != nil {
			blog.Errorf("upgrade version %s error: %s", v.version, err.Error())
			return err
		}
		cmdbVision.CurrentVersion = v.version
		err = saveVesion(db, cmdbVision)
		if err != nil {
			blog.Errorf("save version %s error: %s", v.version, err.Error())
			return err
		}
		blog.Info("upgrade to version %s success", v.version)
	}
	if "" == cmdbVision.InitVersion {
		cmdbVision.InitVersion = lastVersion
		cmdbVision.InitDistroVersion = ccversion.CCDistroVersion
		saveVesion(db, cmdbVision)
	}
	return nil
}

func getVersion(db storage.DI) (*Version, error) {
	data := new(Version)
	condition := map[string]interface{}{
		"type": SystemTypeVersion,
	}
	err := db.GetOneByCondition(common.BKTableNameSystem, nil, condition, &data)
	if err == mgo.ErrNotFound {
		data = new(Version)
		data.Type = SystemTypeVersion
		_, err = db.Insert(common.BKTableNameSystem, data)
		if err != nil {
			return nil, err
		}
		return data, nil
	}
	if err != nil {
		blog.Error("get system version error", err.Error())
		return nil, err
	}

	return data, nil
}

func saveVesion(db storage.DI, version *Version) error {
	condition := map[string]interface{}{
		"type": SystemTypeVersion,
	}
	return db.UpdateByCondition(common.BKTableNameSystem, version, condition)
}

type System struct {
	Type string `bson:"type"`
}

type Version struct {
	System            `bson:",inline"`
	CurrentVersion    string `bson:"current_version"`
	Distro            string `bson:"distro"`
	DistroVersion     string `bson:"distro_version"`
	InitVersion       string `bson:"init_version"`
	InitDistroVersion string `bson:"init_distro_version"`
}

const SystemTypeVersion = "version"
