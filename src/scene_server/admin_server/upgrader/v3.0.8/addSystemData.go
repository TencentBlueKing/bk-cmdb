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
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage"
)

func addSystemData(db storage.DI, conf *upgrader.Config) error {
	tablename := "cc_System"
	blog.V(3).Infof("add data for  %s table ", tablename)
	data := map[string]interface{}{
		common.HostCrossBizField: common.HostCrossBizValue}
	isExist, err := db.GetCntByCondition(tablename, data)
	if nil != err {
		blog.Errorf("add data for  %s table error  %s", tablename, err)
		return err
	}
	if isExist > 0 {
		return nil
	}
	_, err = db.Insert(tablename, data)
	if nil != err {
		blog.Errorf("add data for  %s table error  %s", tablename, err)
		return err
	}

	return nil
}
