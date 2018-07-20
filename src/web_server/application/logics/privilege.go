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
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/web_server/application/middleware/privilege"
)

//GetUserAppPri get user privilege
func GetUserAppPri(apiAddr string, userName string, ownerID, lang string) (userPriveApp map[int64][]string, rolePrivi map[string][]string, modelConfigPrivi map[string][]string, sysPrivi []string, mainLineObjIDArr []string) {
	p, _ := privilege.NewPrivilege(userName, apiAddr, ownerID, lang)
	appRole := p.GetAppRole()
	rolePrivi = make(map[string][]string)
	blog.Info("get app role result:%v", appRole)
	if nil != appRole {
		userPriveApp = p.GetUserPrivilegeApp(appRole)
		blog.Info("get user privi app:%v", userPriveApp)
		if nil != userPriveApp {
			for _, role := range appRole {
				rolePrivilege := p.GetRolePrivilege(common.BKInnerObjIDApp, role)
				rolePrivi[role] = rolePrivilege
				blog.Info("get role privilege:%v", rolePrivilege)
			}
		}

	}

	modelConfigPrivi, sysPrivi = p.GetUserPrivilegeConfig()
	mainLineObjIDArr = p.GetAllMainLineObject()
	return

}
