/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 THL A29 Limited,
 * a Tencent company. All rights reserved.
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

package upgrader

import (
	"context"
	"fmt"
	"regexp"
	"sort"
	"sync"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/dal/mongo/sharding"
)

// RegisterUpgrade register upgrade handler
func RegisterUpgrade(version string, handlerFunc func(*rest.Kit, dal.Dal) error) {

	if err := validateVersionFormat(version); err != nil {
		blog.Errorf("invalid version format, err: %v", err)
		return
	}
	registerLock.Lock()
	v := Upgrader{
		version: version,
		do: func(kit *rest.Kit, rdb dal.Dal, op *Options) error {
			return handlerFunc(kit, rdb)
		},
	}
	upgraderPool = append(upgraderPool, v)
	registerLock.Unlock()
}

func validateVersionFormat(version string) error {

	if !validVersionRegx.MatchString(version) {
		blog.Errorf("invalid migration version format: %s", version)
		return fmt.Errorf("invalid migration version format: %s", version)
	}

	return validateDateTimeString(version)
}

func validateDateTimeString(version string) error {
	const layoutISO = "200601021504"
	versionInfo := version[0 : len(version)-len(layoutISO)-1]
	versionPattern := `[\d\.]+`
	re := regexp.MustCompile(versionPattern)
	match := re.FindString(versionInfo)
	if match < "3.15" {
		blog.Errorf("version %s should be greater than v3.15", version)
		return fmt.Errorf("version %s should be greater than v3.15", version)
	}
	datetime := version[len(version)-len(layoutISO):]
	t, err := time.Parse(layoutISO, datetime)
	if err != nil {
		return fmt.Errorf("invalid datetime format: %s", datetime)
	}

	year, month, day, hour, minute := t.Year(), int(t.Month()), t.Day(), t.Hour(), t.Minute()

	if year == 0 || month <= 0 || month > 12 || day <= 0 || day > 31 || hour < 0 || hour > 23 || minute < 0 || minute > 59 {
		return fmt.Errorf("invalid datetime format: %s", datetime)
	}

	maxMigrationTime := time.Now().AddDate(0, 0, 1)
	maxVersionCurrently := maxMigrationTime.Format(layoutISO)
	if datetime >= maxVersionCurrently {
		return fmt.Errorf("invalid datetime format, time is longer than now: %s", datetime)
	}
	return nil
}

// GetVersion get version info
func GetVersion(ctx context.Context, db dal.RDB) (*Version, error) {
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

// SaveVersion save version info
func SaveVersion(ctx context.Context, db dal.RDB, version *Version) error {
	condition := map[string]interface{}{
		"type": SystemTypeVersion,
	}
	return db.Table(common.BKTableNameSystem).Update(ctx, condition, version)
}

// DBReady 已经执行过init_db. 数据库初始化成功
func DBReady(ctx context.Context, db dal.Dal) (bool, error) {

	sort.Slice(upgraderPool, func(i, j int) bool {
		return VersionCmp(upgraderPool[i].version, upgraderPool[j].version) < 0
	})

	cmdbVersion, err := GetVersion(ctx, db.Shard(sharding.NewShardOpts().WithIgnoreTenant()))
	if err != nil {
		return false, fmt.Errorf("getVersion failed, err: %s", err)
	}

	currentVersion := ""
	for _, v := range upgraderPool {
		if VersionCmp(v.version, currentVersion) <= 0 {
			blog.Infof(`currentVision is "%s" skip upgrade "%s"`, currentVersion, v.version)
			continue
		}
		currentVersion = v.version
	}
	if currentVersion == cmdbVersion.CurrentVersion {
		return true, nil
	}
	return false, nil
}

// UpgradeSpecifyVersion 强制执行version版本的migrate, 不会修改数据库cc_System表中migrate 版本
func UpgradeSpecifyVersion(kit *rest.Kit, db dal.Dal, version string) (err error) {

	sort.Slice(upgraderPool, func(i, j int) bool {
		return VersionCmp(upgraderPool[i].version, upgraderPool[j].version) < 0
	})

	hasCurrent := false
	for _, v := range upgraderPool {
		if v.version != version {
			continue
		}
		blog.Infof(`run specify migration: %s`, v.version)
		err = v.do(kit, db, &Options{})
		if err != nil {
			blog.Errorf("upgrade specify version %s error: %s", v.version, err.Error())
			return fmt.Errorf("run specify migration %s failed, err: %s", v.version, err.Error())
		}
		hasCurrent = true
	}
	if !hasCurrent {
		return fmt.Errorf("run specify migration %s failed, err: not found", version)
	}

	return nil
}

// System system info
type System struct {
	Type string `bson:"type"`
}

// Version version info
type Version struct {
	System            `bson:",inline"`
	CurrentVersion    string `bson:"current_version"`
	Distro            string `bson:"distro"`
	DistroVersion     string `bson:"distro_version"`
	InitVersion       string `bson:"init_version"`
	InitDistroVersion string `bson:"init_distro_version"`
}

// SystemTypeVersion system type version
const SystemTypeVersion = "version"

var (
	registerLock     sync.Mutex
	upgraderPool     = make([]Upgrader, 0)
	validVersionRegx = regexp.MustCompile(`^y(\d+\.){2}\d{12}$`) // yx.xx.xxxxxxxxxxxx
)

// Options upgrader options
type Options struct{}
