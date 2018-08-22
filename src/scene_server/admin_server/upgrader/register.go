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
	"context"
	"sort"
	"sync"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	ccversion "configcenter/src/common/version"
	"configcenter/src/storage/dal"
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
	do      func(context.Context, dal.RDB, *Config) error
}

var upgraderPool = []Upgrader{}

var registlock sync.Mutex

// RegistUpgrader register upgrader
func RegistUpgrader(version string, handlerFunc func(context.Context, dal.RDB, *Config) error) {
	registlock.Lock()
	defer registlock.Unlock()
	v := Upgrader{version: version, do: handlerFunc}
	upgraderPool = append(upgraderPool, v)
	// blog.Infof("registed upgrader for version %s", v.version)
}

// Upgrade uprade the db datas to newest verison
func Upgrade(ctx context.Context, db dal.RDB, conf *Config) (err error) {
	sort.Slice(upgraderPool, func(i, j int) bool {
		return upgraderPool[i].version < upgraderPool[j].version
	})

	cmdbVision, err := getVersion(ctx, db)
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
		err = v.do(ctx, db, conf)
		if err != nil {
			blog.Errorf("upgrade version %s error: %s", v.version, err.Error())
			return err
		}
		cmdbVision.CurrentVersion = v.version
		err = saveVesion(ctx, db, cmdbVision)
		if err != nil {
			blog.Errorf("save version %s error: %s", v.version, err.Error())
			return err
		}
		blog.Info("upgrade to version %s success", v.version)
	}
	if "" == cmdbVision.InitVersion {
		cmdbVision.InitVersion = lastVersion
		cmdbVision.InitDistroVersion = ccversion.CCDistroVersion
		saveVesion(ctx, db, cmdbVision)
	}
	return nil
}

func getVersion(ctx context.Context, db dal.RDB) (*Version, error) {
	data := new(Version)
	condition := map[string]interface{}{
		"type": SystemTypeVersion,
	}
	err := db.Table(common.BKTableNameSystem).Find(condition).One(ctx, &data)
	if db.IsNotFoundError(err) {
		data = new(Version)
		data.Type = SystemTypeVersion

		err = db.Table(common.BKTableNameSystem).Insert(ctx, data)
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

func saveVesion(ctx context.Context, db dal.RDB, version *Version) error {
	condition := map[string]interface{}{
		"type": SystemTypeVersion,
	}
	return db.Table(common.BKTableNameSystem).Update(ctx, condition, version)
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
