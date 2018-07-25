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
	"errors"
	"fmt"
	"reflect"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/language"
	"configcenter/src/common/util"
)

func (lgc *Logics) GetObjectByID(objType string, fields []string, id int64, result interface{}, sort string) error {
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
	err := lgc.Instance.GetOneByCondition(tName, fields, condition, result)
	return err
}

func (lgc *Logics) CreateObject(objType string, input interface{}, idName *string) (int64, error) {
	tName := common.GetInstTableName(objType)
	objID, err := lgc.Instance.GetIncID(tName)
	if err != nil {
		return 0, err
	}
	inputc := input.(map[string]interface{})
	*idName = common.GetInstIDField(objType)
	inputc[*idName] = objID
	_, err = lgc.Instance.Insert(tName, inputc)
	if err != nil {
		return 0, err
	}
	return objID, nil
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

func (lgc *Logics) GetObjectByCondition(defLang language.DefaultCCLanguageIf, objType string, fields []string, condition, result interface{}, sort string, skip, limit int) error {
	tName := common.GetInstTableName(objType)
	if err := lgc.Instance.GetMutilByCondition(tName, fields, condition, result, sort, skip, limit); err != nil {
		blog.Errorf("failed to query the inst , error info %s", err.Error())
		return err
	}

	// translate language for default name
	if m, ok := defaultNameLanguagePkg[objType]; nil != defLang && ok {
		switch result.(type) {
		case *[]map[string]interface{}:
			results := *result.(*[]map[string]interface{})
			for index := range results {
				l := m[fmt.Sprint(results[index]["default"])]
				if len(l) >= 3 {
					results[index][l[1]] = util.FirstNotEmptyString(defLang.Language(l[0]), fmt.Sprint(results[index][l[1]]), fmt.Sprint(results[index][l[2]]))
				}
			}
		default:
			blog.Infof("get object by condition translate error: %v", reflect.TypeOf(result))
		}
	}

	return nil
}
