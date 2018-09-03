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

package v3v0v1alpha2

import (
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage"
)

func addBkStartParamRegex(db storage.DI, conf *upgrader.Config) (err error) {
	tablename := common.BKTableNameObjAttDes
	now := time.Now()

	row := &metadata.Attribute{
		ObjectID:      common.BKInnerObjIDProc,
		PropertyID:    "bk_start_param_regex",
		PropertyName:  "启动参数匹配规则",
		IsRequired:    false,
		IsOnly:        false,
		IsEditable:    true,
		PropertyGroup: "default",
		PropertyType:  common.FieldTypeLongChar,
		Option:        "",
		OwnerID:       conf.OwnerID,
		IsPre:         true,
		IsReadOnly:    false,
		CreateTime:    &now,
		Creator:       common.CCSystemOperatorUserName,
		LastTime:      &now,
		Description:   "通过进程启动参数唯一识别进程，比如kafka和zookeeper的二进制名称为java，通过启动参数包含kafka或zookeeper来区分",
	}
	_, _, err = upgrader.Upsert(db, tablename, row, "id", []string{common.BKObjIDField, common.BKPropertyIDField, common.BKOwnerIDField}, []string{})
	if nil != err {
		blog.Errorf("[upgrade v3.1.0-alpha.2] addBkStartParamRegex  %s", err)
		return err
	}

	return nil
}

func updateLanguageField(db storage.DI, conf *upgrader.Config) (err error) {
	condition := map[string]interface{}{
		common.BKObjIDField:      common.BKInnerObjIDApp,
		common.BKPropertyIDField: "language",
	}
	data := map[string]interface{}{
		"isreadonly": false,
		"editable":   true,
	}
	err = db.UpdateByCondition(common.BKTableNameObjAttDes, data, condition)
	if nil != err {
		blog.Errorf("[upgrade v3.1.0-alpha.2] updateLanguageField error  %s", err.Error())
		return err
	}
	return nil
}
