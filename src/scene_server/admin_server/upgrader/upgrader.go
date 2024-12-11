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
	"fmt"
	"sort"
	"strings"

	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	ccversion "configcenter/src/common/version"
	"configcenter/src/scene_server/admin_server/upgrader/types"
	"configcenter/src/storage/dal"
)

// Upgrade cmdb to new version
func Upgrade(kit *rest.Kit, db dal.Dal, op *Options) (*types.MigrateInfo, error) {

	sort.Slice(upgraderPool, func(i, j int) bool {
		return strings.Compare(upgraderPool[i].version, upgraderPool[j].version) < 0
	})

	cmdbVersion, err := GetVersion(kit.Ctx, db.Shard(kit.SysShardOpts()))
	if err != nil {
		blog.Errorf("get version failed, err: %v, rid: %s", err, kit.Rid)
		return nil, err
	}
	if cmdbVersion.CurrentVersion < "y3.15" && cmdbVersion.CurrentVersion != "" {
		blog.Errorf("current version is %s, need a new environment, can not upgrade", cmdbVersion.CurrentVersion)
		return nil, fmt.Errorf("current version is %s, need a new environment, can not upgrade",
			cmdbVersion.CurrentVersion)
	}
	cmdbVersion.Distro = ccversion.CCDistro
	cmdbVersion.DistroVersion = ccversion.CCDistroVersion

	currentVersion := cmdbVersion.CurrentVersion
	preVersion := cmdbVersion.CurrentVersion
	finishedMigrations := make([]string, 0)
	for _, v := range upgraderPool {
		if strings.Compare(v.version, currentVersion) <= 0 {
			blog.Infof("currentVision is %s skip upgrade %s", currentVersion, v.version)
			continue
		}
		blog.Infof("run migration: %s", v.version)
		err = v.do(kit, db, op)
		if err != nil {
			blog.Errorf("upgrade version %s failed, error: %v", v.version, err)
			return nil, err
		}
		cmdbVersion.CurrentVersion = v.version
		err = SaveVersion(kit.Ctx, db.Shard(kit.SysShardOpts()), cmdbVersion)
		if err != nil {
			blog.Errorf("save version %s failed, error: %s", v.version, err)
			return nil, err
		}
		finishedMigrations = append(finishedMigrations, v.version)
		blog.Infof("upgrade to version %s success", v.version)
		currentVersion = v.version
	}

	if cmdbVersion.InitVersion == "" {
		cmdbVersion.InitVersion = currentVersion
		cmdbVersion.InitDistroVersion = ccversion.CCDistroVersion
		if err := SaveVersion(kit.Ctx, db.Shard(kit.SysShardOpts()), cmdbVersion); err != nil {
			return nil, err
		}
	}

	result := &types.MigrateInfo{
		PreVersion:       preVersion,
		CurrentVersion:   currentVersion,
		FinishedVersions: finishedMigrations,
	}

	return result, nil
}

// Upgrader define a version upgrader
type Upgrader struct {
	version string
	do      func(*rest.Kit, dal.Dal, *Options) error
}
