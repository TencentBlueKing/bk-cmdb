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

package x18_09_30_01

import (
	"fmt"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/mapstr"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage"
)

func cleanBKCloud(db storage.DI, conf *upgrader.Config) (err error) {

	clouds := []map[string]interface{}{}

	err = db.GetMutilByCondition(common.BKTableNameBasePlat, nil, mapstr.MapStr{}, &clouds, "create_time", 0, 0)
	if nil != err && !db.IsNotFoundErr(err) {
		return err
	}

	flag := "updateflag"
	existDefault := false
	expects := map[string]map[string]interface{}{}
	for _, cloud := range clouds {
		if cloud[common.BKCloudNameField] == common.DefaultCloudName {
			cloud[common.BKCloudIDField] = 0
			existDefault = true
		}
		cloud[flag] = true
		expects[fmt.Sprintf("%v:%v", cloud[common.BKOwnerIDField], cloud[common.BKCloudNameField])] = cloud
	}

	if !existDefault {
		expects["0:"+common.DefaultCloudName] = map[string]interface{}{
			common.BKCloudNameField: common.DefaultCloudName,
			common.BKOwnerIDField:   common.BKDefaultOwnerID,
			common.BKCloudIDField:   common.BKDefaultDirSubArea,
			common.CreateTimeField:  time.Now(),
			common.LastTimeField:    time.Now(),
			flag:                    true,
		}
	}

	for _, expect := range expects {
		if _, err = db.Insert(common.BKTableNameBasePlat, expect); err != nil {
			return err
		}
	}

	if err = db.DelByCondition(common.BKTableNameBasePlat, map[string]interface{}{
		flag: map[string]interface{}{
			common.BKDBNE: true,
		},
	}); err != nil {
		return err
	}

	if err = db.DropColumn(common.BKTableNameBasePlat, flag); err != nil {
		return err
	}

	return nil
}
