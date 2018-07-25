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

package v3v0v9beta3

import (
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/scene_server/validator"
	"configcenter/src/storage"
)

func fixesProcessPortPattern(db storage.DI, conf *upgrader.Config) (err error) {
	condition := map[string]interface{}{
		common.BKObjIDField:      common.BKInnerObjIDProc,
		common.BKPropertyIDField: "port",
	}
	data := map[string]interface{}{
		"option": common.PatternMultiplePortRange,
	}
	err = db.UpdateByCondition(common.BKTableNameObjAttDes, data, condition)
	if nil != err {
		blog.Errorf("[upgrade v3.0.9-beta.3] fixesPortPattern error  %s", err.Error())
		return err
	}
	return nil
}

func fixesProcessPriorityPattern(db storage.DI, conf *upgrader.Config) (err error) {
	condition := map[string]interface{}{
		common.BKObjIDField:      common.BKInnerObjIDProc,
		common.BKPropertyIDField: "priority",
	}
	data := map[string]interface{}{
		"option": validator.IntOption{Min: "1", Max: "10000"},
	}
	err = db.UpdateByCondition(common.BKTableNameObjAttDes, data, condition)
	if nil != err {
		blog.Errorf("[upgrade v3.0.9-beta.3] fixesPortPattern error  %s", err.Error())
		return err
	}
	return nil
}
