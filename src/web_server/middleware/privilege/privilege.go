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

package privilege

import (
	"context"
	"net/http"
	"strings"

	"configcenter/src/common"
	"configcenter/src/common/backbone"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/util"

	"github.com/gin-gonic/gin"
)

type Privilege struct {
	Engine   *backbone.Engine
	UserName string
	OwnerID  string
	language string
	header   http.Header
}

func NewPrivilege(userName, ownerID, language string) *Privilege {
	privi := new(Privilege)
	privi.UserName = userName
	privi.OwnerID = ownerID
	header := make(http.Header)
	header.Add(common.BKHTTPHeaderUser, userName)
	header.Add(common.BKHTTPLanguage, language)
	header.Add(common.BKHTTPOwnerID, ownerID)
	privi.header = header
	return privi
}

func ValidPrivilege() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(200, gin.H{
			"bk_error_msg": "pong",
		})
	}
}

// GetRolePrivilege get role privilege
func (p *Privilege) GetRolePrivilege(objID, role string) []string {
	result, err := p.Engine.CoreAPI.ApiServer().GetRolePrivilege(context.Background(), p.header, p.OwnerID, objID, role)
	if nil != err || !result.Result {
		blog.Warnf("get role privilege json error: %v", err)
		return []string{}
	}
	return result.Data
}

// GetAppRole get app role
func (p *Privilege) GetAppRole() []string {
	data := make([]string, 0)
	params := mapstr.MapStr{common.BKPropertyTypeField: "objuser", common.BKObjIDField: common.BKInnerObjIDApp}
	result, err := p.Engine.CoreAPI.ApiServer().GetAppRole(context.Background(), p.header, params)
	if nil != err || !result.Result {
		blog.Warnf("get role privilege json error: %v", err)
		return data
	}
	for _, i := range result.Data {
		propertyID, ok := i[common.BKPropertyIDField].(string)
		if false == ok {
			continue
		}
		data = append(data, propertyID)
	}
	return data
}

// GetUserPrivilegeApp get user privilege app
func (p *Privilege) GetUserPrivilegeApp(appRole []string) map[int64][]string {
	orCond := make([]mapstr.MapStr, 0)
	allCond := mapstr.MapStr{}
	condition := mapstr.MapStr{}
	for _, role := range appRole {
		cell := mapstr.MapStr{}
		d := mapstr.MapStr{}
		cell[common.BKDBLIKE] = p.UserName
		d[role] = cell
		orCond = append(orCond, d)
	}
	allCond[common.BKDBOR] = orCond
	condition["condition"] = allCond
	condition["native"] = 1
	userRole := make(map[int64][]string)

	result, err := p.Engine.CoreAPI.ApiServer().GetUserPrivilegeApp(context.Background(), p.header, p.OwnerID, p.UserName, condition)
	if nil != err || !result.Result {
		blog.Errorf("get role privilege app error: %v", err)
		return userRole
	}

	for _, i := range result.Data.Info {
		appID, err := util.GetInt64ByInterface(i[common.BKAppIDField])
		if nil != err {
			continue
		}
		userRoleArr := make([]string, 0)
		for _, j := range appRole {
			roleData, ok := i[j]
			if false == ok {
				continue
			}
			roleStr, ok := roleData.(string)
			if false == ok {
				continue
			}
			roleArr := strings.Split(roleStr, ",")
			for _, k := range roleArr {
				if k == p.UserName {
					userRoleArr = append(userRoleArr, j)
				}
			}
			userRole[appID] = userRoleArr
		}
	}
	return userRole
}

// GetUserPrivilegeConfig get user privilege config
func (p *Privilege) GetUserPrivilegeConfig() (map[string][]string, []string) {
	result, err := p.Engine.CoreAPI.ApiServer().GetUserPrivilegeConfig(context.Background(), p.header, p.OwnerID, p.UserName)

	if nil != err || false == result.Result {
		blog.Warnf("get user privilege json error: %v", err)
		return nil, nil
	}
	sysConfig := make([]string, 0)
	modelConfig := make(map[string][]string, 0)
	for _, i := range result.Data.SysConfig.BackConfig {
		sysConfig = append(sysConfig, i)
	}

	for _, i := range result.Data.SysConfig.Globalbusi {
		sysConfig = append(sysConfig, i)
	}
	for _, j := range result.Data.ModelConfig {
		for m, n := range j {
			modelConfig[m] = n
		}
	}
	return modelConfig, sysConfig
}

// GetAllMainLineObject get all main line object
func (p *Privilege) GetAllMainLineObject() []string {
	mainLineObjName := make([]string, 0)
	result, err := p.Engine.CoreAPI.ApiServer().GetAllMainLineObject(context.Background(), p.header, p.OwnerID, p.UserName)
	if nil != err || false == result.Result {
		blog.Warnf("get all main line object error: %v", err)
		return mainLineObjName
	}
	for _, data := range result.Data {
		objID, ok := data[common.BKObjIDField].(string)
		if ok {
			mainLineObjName = append(mainLineObjName, objID)
		}
	}
	return mainLineObjName
}
