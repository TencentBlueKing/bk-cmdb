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
	"fmt"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/util"
	ccversion "configcenter/src/common/version"
	"configcenter/src/storage/dal"
)

// Config config for upgrader
type Config struct {
	OwnerID      string
	SupplierID   int
	User         string
	CCApiSrvAddr string // cmdb nginx address
}

// Upgrader define a version upgrader
type Upgrader struct {
	version string // v3.0.8-beta.11
	do      func(context.Context, dal.RDB, *Config) error
}

var upgraderPool = []Upgrader{}

var registLock sync.Mutex

/*
	v3.5.x format:
		x19_09_03_02
	v3.6.x format:
		y3.6.201911081042
	legacy format:
		v3.0.8
		v3.0.9-beta.1
		x18_12_13_02
		x08.09.04.01
		x18.09.30.01
		x19.01.18.01
		x19_09_03_02
*/

var LegacyMigrationVersion = []string{
	"v3.0.8",
	"v3.0.9-beta.1",
	"v3.0.9-beta.3",
	"v3.1.0-alpha.2",
	"x08.09.04.01",
	"x08.09.11.01",
	"x08.09.13.01",
	"x08.09.17.01",
	"x08.09.18.01",
	"x08.09.26.01",
	"x18.09.30.01",
	"x18.10.10.01",
	"x18.10.30.01",
	"x18.10.30.02",
	"x18.11.19.01",
	"x18.12.12.01",
	"x18.12.12.02",
	"x18.12.12.03",
	"x18.12.12.04",
	"x18.12.12.05",
	"x18.12.12.06",
	"x18.12.13.01",
	"x18_12_13_02",
	"x19.01.18.01",
	"x19.02.15.10",
	"x19.04.16.01",
	"x19.04.16.02",
	"x19.04.16.03",
	"x19.05.16.01",
	"x19_05_22_01",
	"x19_08_19_01",
	"x19_08_20_01",
	"x19_08_26_01",
	"x19_08_26_02",
	"x19_09_03_01",
	"x19_09_03_02",
	"x19_09_03_03",
	"x19_09_03_04",
	"x19_09_03_05",
	"x19_09_03_06",
	"x19_09_03_07",
	"x19_09_03_08",
	"x19_09_06_01",
	"x19_09_27_01",
	"x19_10_09_01",
	"x19_10_22_01",
}

var ValidMigrationVersionFormat = []*regexp.Regexp{
	// regexp.MustCompile(`^v(\d+\.){2}\d+$`),
	// regexp.MustCompile(`^v(\d+\.){2}\d+\-beta\.\d+$`),
	// regexp.MustCompile(`^v(\d+\.){2}\d+\-alpha\.\d+$`),
	// regexp.MustCompile(`^x(\d+\.){3}\d+$`),

	// v3.5.x version format
	regexp.MustCompile(`^x(\d+_){3}\d+$`),
	// v3.6.x version format
	regexp.MustCompile(`^y(\d+\.){2}\d{12}$`),
}

func ValidateMigrationVersionFormat(version string) error {
	// only check newer add migration's version
	if util.InStrArr(LegacyMigrationVersion, version) == true {
		return nil
	}
	match := false
	for _, re := range ValidMigrationVersionFormat {
		if re.MatchString(version) == true {
			match = true
			break
		}
	}
	if match == false {
		err := fmt.Errorf(`
	invalid migration version: %s,
    please use a valid format:
      x19_09_03_02(v3.5.x)
      y3.6.201911081042(>=v3.6.x)
	`, version)
		return err
	}

	// since v3.6.x migration version must
	if strings.HasPrefix(version, VersionNgPrefix) {
		ngVersion, err := ParseNgVersion(version)
		if err != nil {
			return err
		}

		// third field in version split by `.` shouldn't greater than tomorrow
		timeFormat := "200601021504"
		maxMigrationTime := time.Now().AddDate(0, 0, 1)
		maxVersionCurrently := maxMigrationTime.Format(timeFormat)
		if ngVersion.Patch >= maxVersionCurrently {
			err := fmt.Errorf(`
	invalid time field of migration version: %s,
    please use current time as part of migration version:
      ex: y3.6.%s
	`, version, time.Now().Format(timeFormat))
			return err
		}
	}
	return nil
}

// RegistUpgrader register upgrader
func RegistUpgrader(version string, handlerFunc func(context.Context, dal.RDB, *Config) error) {
	if err := ValidateMigrationVersionFormat(version); err != nil {
		blog.Fatalf("ValidateMigrationVersionFormat failed, err: %s", err.Error())
	}
	registLock.Lock()
	defer registLock.Unlock()
	v := Upgrader{version: version, do: handlerFunc}
	upgraderPool = append(upgraderPool, v)
}

// Upgrade upgrade the db data to newest version
// we use date instead of version later since 2018.09.04, because the version wasn't manage by the developer
// ps: when use date instead of version, the date should add x prefix cause x > v
func Upgrade(ctx context.Context, db dal.RDB, conf *Config) (currentVersion string, finishedMigrations []string, err error) {

	sort.Slice(upgraderPool, func(i, j int) bool {
		return VersionCmp(upgraderPool[i].version, upgraderPool[j].version) < 0
	})

	cmdbVersion, err := getVersion(ctx, db)
	if err != nil {
		return "", nil, fmt.Errorf("getVersion failed, err: %s", err.Error())
	}
	cmdbVersion.Distro = ccversion.CCDistro
	cmdbVersion.DistroVersion = ccversion.CCDistroVersion

	currentVersion = remapVersion(cmdbVersion.CurrentVersion)
	lastVersion := ""
	finishedMigrations = make([]string, 0)
	for _, v := range upgraderPool {
		lastVersion = remapVersion(v.version)
		if VersionCmp(v.version, currentVersion) <= 0 {
			blog.Infof(`currentVision is "%s" skip upgrade "%s"`, currentVersion, v.version)
			continue
		}
		blog.Infof(`run migration: %s`, v.version)
		err = v.do(ctx, db, conf)
		if err != nil {
			blog.Errorf("upgrade version %s error: %s", v.version, err.Error())
			return currentVersion, finishedMigrations, fmt.Errorf("run migration %s failed, err: %s", v.version, err.Error())
		}
		cmdbVersion.CurrentVersion = v.version
		err = saveVersion(ctx, db, cmdbVersion)
		if err != nil {
			blog.Errorf("save version %s error: %s", v.version, err.Error())
			return currentVersion, finishedMigrations, fmt.Errorf("saveVersion failed, err: %s", err.Error())
		}
		finishedMigrations = append(finishedMigrations, v.version)
		blog.Infof("upgrade to version %s success", v.version)
	}

	if "" == cmdbVersion.InitVersion {
		cmdbVersion.InitVersion = lastVersion
		cmdbVersion.InitDistroVersion = ccversion.CCDistroVersion
		if err := saveVersion(ctx, db, cmdbVersion); err != nil {
			return currentVersion, finishedMigrations, fmt.Errorf("saveVersion failed, err: %s", err.Error())
		}
	}
	return currentVersion, finishedMigrations, nil
}

func remapVersion(v string) string {
	if correct, ok := wrongVersion[v]; ok {
		return correct
	}
	return v
}

var wrongVersion = map[string]string{
	"x18_10_10_01": "x18.10.10.01",
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
		blog.Errorf("get system version error,err:%s", err.Error())
		return nil, err
	}

	return data, nil
}

func saveVersion(ctx context.Context, db dal.RDB, version *Version) error {
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
