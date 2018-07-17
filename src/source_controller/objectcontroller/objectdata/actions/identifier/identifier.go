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

package instdata

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/emicklei/go-restful"

	"configcenter/src/common"
	"configcenter/src/common/base"
	"configcenter/src/common/blog"
	"configcenter/src/common/core/cc/actions"
	"configcenter/src/common/util"
	"configcenter/src/source_controller/api/metadata"
	"configcenter/src/source_controller/common/instdata"
	"configcenter/src/storage"
)

var obj = &identifierAction{}

// identifierAction
type identifierAction struct {
	base.BaseAction
}

func init() {
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "/identifier/{obj_type}/search", Params: nil, Handler: obj.SearchIdentifier})
	// set cc api interface
	obj.CreateAction()
}

//search object
func (cli *identifierAction) SearchIdentifier(req *restful.Request, resp *restful.Response) {
	// get the language
	language := util.GetActionLanguage(req)
	// get the error factory by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)
	cli.CallResponseEx(func() (int, interface{}, error) {
		instdata.DataH = cli.CC.InstCli
		param := new(SearchIdentifierParam)
		err := json.NewDecoder(req.Request.Body).Decode(param)
		if err != nil {
			blog.Errorf("SearchIdentifier error:%s", err.Error())
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommJSONUnmarshalFailed)
		}
		blog.Infof("SearchIdentifier %v", param)

		var (
			setIDs    = []int{}
			moduleIDs = []int{}
			bizIDs    = []int{}
			cloudIDs  = []int{}
		)

		// caches
		var (
			sets    = map[int]metadata.SetInst{}
			modules = map[int]metadata.ModuleInst{}
			bizs    = map[int]metadata.BizInst{}
			clouds  = map[int]metadata.CloudInst{}
		)

		condition := map[string]interface{}{
			common.BKDBOR: []map[string]interface{}{
				{
					common.BKHostInnerIPField: map[string]interface{}{
						common.BKDBIN: param.IP.Data,
					},
				}, {
					common.BKHostOuterIPField: map[string]interface{}{
						common.BKDBIN: param.IP.Data,
					},
				},
			},
		}

		if param.IP.CloudID != nil {
			condition[common.BKCloudIDField] = *param.IP.CloudID
		}

		// fetch hosts
		hosts := []*metadata.HostIdentifier{}
		err = instdata.GetHostByCondition(nil, condition, &hosts, "", 0, 0)
		if err != nil {
			blog.Errorf("SearchIdentifier error:%s", err.Error())
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrObjectSelectIdentifierFailed)
		}
		for _, host := range hosts {
			relations := []metadata.ModuleHostConfig{}
			condiction := map[string]interface{}{
				common.BKHostIDField: host.HostID,
			}
			blog.Infof("SearchIdentifier relations condition %v ", condiction)
			err = cli.CC.InstCli.GetMutilByCondition(common.BKTableNameModuleHostConfig, nil, condiction, &relations, "", -1, -1)
			if err != nil {
				blog.Errorf("SearchIdentifier error:%s", err.Error())
				return http.StatusInternalServerError, nil, defErr.Error(common.CCErrObjectSelectIdentifierFailed)
			}
			host.HostIdentModule = map[string]*metadata.HostIdentModule{}
			for _, rela := range relations {
				host.HostIdentModule[fmt.Sprint(rela.ModuleID)] = &metadata.HostIdentModule{
					SetID:    rela.SetID,
					ModuleID: rela.ModuleID,
					BizID:    rela.ApplicationID,
				}
				setIDs = append(setIDs, rela.SetID)
				moduleIDs = append(moduleIDs, rela.ModuleID)
				bizIDs = append(bizIDs, rela.ApplicationID)
			}
			cloudIDs = append(cloudIDs, host.CloudID)
		}

		blog.Infof("sets: %v, modules: %v, bizs: %v, clouds: %v", setIDs, moduleIDs, bizIDs, cloudIDs)
		// fetch cache
		if len(setIDs) > 0 {
			tmps := []metadata.SetInst{}
			err = getCache(cli.CC.InstCli, common.BKTableNameBaseSet, common.BKSetIDField, setIDs, &tmps)
			if err != nil {
				blog.Errorf("SearchIdentifier error:%s", err.Error())
				return http.StatusInternalServerError, nil, defErr.Error(common.CCErrObjectSelectIdentifierFailed)
			}
			for _, tmp := range tmps {
				sets[tmp.SetID] = tmp
			}
		}
		if len(moduleIDs) > 0 {
			tmps := []metadata.ModuleInst{}
			err = getCache(cli.CC.InstCli, common.BKTableNameBaseModule, common.BKModuleIDField, moduleIDs, &tmps)
			if err != nil {
			}
			for _, tmp := range tmps {
				modules[tmp.ModuleID] = tmp
			}
		}
		if len(bizIDs) > 0 {
			tmps := []metadata.BizInst{}
			err = getCache(cli.CC.InstCli, common.BKTableNameBaseApp, common.BKAppIDField, bizIDs, &tmps)
			if err != nil {
				blog.Errorf("SearchIdentifier error:%s", err.Error())
				return http.StatusInternalServerError, nil, defErr.Error(common.CCErrObjectSelectIdentifierFailed)
			}
			for _, tmp := range tmps {
				bizs[tmp.BizID] = tmp
			}
		}
		if len(cloudIDs) > 0 {
			tmps := []metadata.CloudInst{}
			err = getCache(cli.CC.InstCli, common.BKTableNameBasePlat, common.BKCloudIDField, cloudIDs, &tmps)
			if err != nil {
				blog.Errorf("SearchIdentifier error:%s", err.Error())
				return http.StatusInternalServerError, nil, defErr.Error(common.CCErrObjectSelectIdentifierFailed)
			}
			for _, tmp := range tmps {
				clouds[tmp.CloudID] = tmp
			}
		}

		blog.Infof("sets: %v, modules: %v, bizs: %v, clouds: %v", sets, modules, bizs, clouds)

		// fill hostidentifier
		for _, inst := range hosts {
			for _, mod := range inst.HostIdentModule {
				if _, ok := bizs[mod.BizID]; ok {
					biz := bizs[mod.BizID]
					mod.BizName = biz.BizName
					inst.SupplierAccount = biz.SupplierAccount
					inst.SupplierID = biz.SupplierID
				}

				if _, ok := sets[mod.SetID]; ok {
					set := sets[mod.SetID]
					mod.SetName = set.SetName
					mod.SetEnv = set.SetEnv
					mod.SetStatus = set.SetStatus
				}

				if _, ok := modules[mod.ModuleID]; ok {
					module := modules[mod.ModuleID]
					mod.ModuleName = module.ModuleName
				}
			}
			if _, ok := clouds[inst.CloudID]; ok {
				cloud := clouds[inst.CloudID]
				inst.CloudName = cloud.CloudName
			}
		}

		// returns
		info := make(map[string]interface{})
		info["count"] = len(hosts)
		info["info"] = hosts
		return http.StatusOK, info, nil
	}, resp)
}

func getCache(db storage.DI, tablename string, idfield string, ids []int, result interface{}) error {
	condition := map[string]interface{}{
		idfield: map[string]interface{}{
			common.BKDBIN: ids,
		},
	}
	return db.GetMutilByCondition(tablename, nil, condition, result, "", 0, 0)
}

// SearchIdentifierParam defines the param
type SearchIdentifierParam struct {
	IP   IPParam `json:"ip"`
	Page struct {
		Start int    `json:"start"`
		Limit int    `json:"limit"`
		Sort  string `json:"sort"`
	} `json:"page"`
}

type IPParam struct {
	Data    []string `json:"data"`
	CloudID *int64   `json:"bk_cloud_id"`
}
