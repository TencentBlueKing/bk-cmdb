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

package multilingual

import (
	"fmt"

	"configcenter/src/common"
	"configcenter/src/common/language"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/util"
)

var BuildInInstanceNamePkg = map[string]map[string][]string{
	common.BKInnerObjIDModule: {
		"1": {"inst_module_idle", common.BKModuleNameField},
		"2": {"inst_module_fault", common.BKModuleNameField},
		"3": {"inst_module_recycle", common.BKModuleNameField},
	},
	common.BKInnerObjIDApp: {
		"1": {"inst_biz_default", common.BKAppNameField},
	},
	common.BKInnerObjIDSet: {
		"1": {"inst_set_default", common.BKSetNameField},
	},
}

// TranslateInstanceName is used to translate build-in model(module/set/biz) instance's name to the
// corresponding language.
// Note: these instances's name is related it's default field's value, different value have different name.
// such as the module's instance, the different meaning of default value is as follows:
// 0: a common module
// 1: a idle module
// 2: a fault module
// 3: a recycle module
func TranslateInstanceName(defLang language.DefaultCCLanguageIf, objectID string, instances []mapstr.MapStr) {
	if m, ok := BuildInInstanceNamePkg[objectID]; ok {
		for idx := range instances {
			// get the default's value and it's corresponding infos from defaultNameLanguagePkg
			subResult := m[fmt.Sprint(instances[idx][common.BKDefaultField])]
			if len(subResult) >= 2 {
				instances[idx][subResult[1]] = util.FirstNotEmptyString(defLang.Language(subResult[0]),
					fmt.Sprint(instances[idx][subResult[1]]))
			}
		}
	}
}
