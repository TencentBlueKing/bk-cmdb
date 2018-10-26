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

package logics

import (
	"configcenter/src/common/mapstr"
	"context"
	"errors"
	"fmt"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/language"
	"configcenter/src/common/util"
)

func (lgc *Logics) GetObjectByID(ctx context.Context, objType string, fields []string, id int64, result interface{}, sort string) error {
	tName := common.GetInstTableName(objType)
	condition := make(map[string]interface{}, 1)
	switch objType {
	case common.BKInnerObjIDApp:
		condition[common.BKAppIDField] = id
	case common.BKInnerObjIDSet:
		condition[common.BKSetIDField] = id
	case common.BKInnerObjIDModule:
		condition[common.BKModuleIDField] = id
	case common.BKINnerObjIDObject:
		condition[common.BKInstIDField] = id
	case common.BKInnerObjIDHost:
		condition[common.BKHostIDField] = id
	case common.BKInnerObjIDProc:
		condition[common.BKProcIDField] = id
	case common.BKInnerObjIDPlat:
		condition[common.BKCloudIDField] = id
	default:
		return errors.New("invalid object type")
	}
	err := lgc.Instance.Table(tName).Find(condition).Fields(fields...).One(ctx, result)
	return err
}

func (lgc *Logics) CreateObject(ctx context.Context, objType string, input interface{}, idName *string) (int64, error) {
	tName := common.GetInstTableName(objType)
	objID, err := lgc.Instance.NextSequence(ctx, tName)
	if err != nil {
		return 0, err
	}
	inputc := input.(map[string]interface{})
	*idName = common.GetInstIDField(objType)
	inputc[*idName] = objID
	err = lgc.Instance.Table(tName).Insert(ctx, inputc)
	if err != nil {
		return 0, err
	}
	return int64(objID), nil
}

var defaultNameLanguagePkg = map[string]map[string][]string{
	common.BKInnerObjIDModule: {
		"1": {"inst_module_idle", common.BKModuleNameField, common.BKModuleIDField},
		"2": {"inst_module_fault", common.BKModuleNameField, common.BKModuleIDField},
	},
	common.BKInnerObjIDApp: {
		"1": {"inst_biz_default", common.BKAppNameField, common.BKAppIDField},
	},
	common.BKInnerObjIDSet: {
		"1": {"inst_set_default", common.BKSetNameField, common.BKSetIDField},
	},
}

func (lgc *Logics) GetObjectByCondition(ctx context.Context, defLang language.DefaultCCLanguageIf, objType string, fields []string, condition interface{}, sort string, skip, limit int) ([]mapstr.MapStr, error) {
	results := make([]mapstr.MapStr, 0)
	tName := common.GetInstTableName(objType)

	dbInst := lgc.Instance.Table(tName).Find(condition).Sort(sort).Start(uint64(skip)).Limit(uint64(limit))
	if 0 < len(fields) {
		dbInst.Fields(fields...)
	}
	if err := dbInst.All(ctx, &results); err != nil {
		blog.Errorf("failed to query the inst , error info %s", err.Error())
		return nil, err
	}

	// translate language for default name
	if m, ok := defaultNameLanguagePkg[objType]; nil != defLang && ok {
		for index, info := range results {
			l := m[fmt.Sprint(info["default"])]
			if len(l) >= 3 {
				results[index][l[1]] = util.FirstNotEmptyString(defLang.Language(l[0]), fmt.Sprint(info[l[1]]), fmt.Sprint(info[l[2]]))
			}
		}
	}

	return results, nil
}
