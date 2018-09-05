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

package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"

	"github.com/emicklei/go-restful"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	cccondition "configcenter/src/common/condition"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/storage"
)

//search object
func (cli *Service) SearchIdentifier(req *restful.Request, resp *restful.Response) {
	// get the language
	language := util.GetActionLanguage(req)
	ownerID := util.GetOwnerID(req.Request.Header)
	// get the error factory by the language
	defErr := cli.Core.CCErr.CreateDefaultCCErrorIf(language)

	param := new(metadata.SearchIdentifierParam)
	err := json.NewDecoder(req.Request.Body).Decode(param)
	if err != nil {
		blog.Errorf("SearchIdentifier error:%s", err.Error())
		resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: defErr.New(common.CCErrCommJSONUnmarshalFailed, err.Error())})
		return

	}

	var (
		hostIDs        = []int64{}
		setIDs         = []int64{}
		moduleIDs      = []int64{}
		bizIDs         = []int64{}
		cloudIDs       = []int64{}
		procIDs        = []int64{}
		appmodulenames = map[int64][]string{}
	)

	// caches
	var (
		sets        = map[int64]metadata.SetInst{}
		modules     = map[int64]metadata.ModuleInst{}
		bizs        = map[int64]metadata.BizInst{}
		clouds      = map[int64]metadata.CloudInst{}
		procs       = map[int64]metadata.HostIdentProcess{}
		modulehosts = map[int64][]metadata.ModuleHost{}
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
	condition = util.SetQueryOwner(condition, ownerID)

	// fetch all hosts
	hosts := []*metadata.HostIdentifier{}
	err = cli.GetHostByCondition(nil, condition, &hosts, "", 0, 0)
	if err != nil && !cli.Instance.IsNotFoundErr(err) {
		blog.Errorf("SearchIdentifier error:%s", err.Error())
		resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: defErr.New(common.CCErrObjectSelectIdentifierFailed, err.Error())})
		return
	}

	// fetch all host module relations
	for _, host := range hosts {
		hostIDs = append(hostIDs, host.HostID)
		cloudIDs = append(cloudIDs, host.CloudID)
	}
	relations := []metadata.ModuleHost{}
	cond := cccondition.CreateCondition().Field(common.BKHostIDField).In(hostIDs)
	err = cli.Instance.GetMutilByCondition(common.BKTableNameModuleHostConfig, nil, cond.ToMapStr(), &relations, "", -1, -1)
	if err != nil && !cli.Instance.IsNotFoundErr(err) {
		blog.Errorf("SearchIdentifier error:%s", err.Error())
		resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: defErr.New(common.CCErrObjectSelectIdentifierFailed, err.Error())})
		return
	}
	for _, modulehost := range relations {
		modulehosts[modulehost.HostID] = append(modulehosts[modulehost.HostID], modulehost)
		setIDs = append(setIDs, modulehost.SetID)
		moduleIDs = append(moduleIDs, modulehost.ModuleID)
		bizIDs = append(bizIDs, modulehost.AppID)
	}

	blog.Infof("sets: %v, modules: %v, bizs: %v, clouds: %v", setIDs, moduleIDs, bizIDs, cloudIDs)
	// fetch cache
	if len(setIDs) > 0 {
		tmps := []metadata.SetInst{}
		err = getCache(cli.Instance, common.BKTableNameBaseSet, common.BKSetIDField, setIDs, &tmps)
		if err != nil && !cli.Instance.IsNotFoundErr(err) {
			blog.Errorf("SearchIdentifier error:%s", err.Error())
			resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: defErr.New(common.CCErrObjectSelectIdentifierFailed, err.Error())})
			return
		}
		for _, tmp := range tmps {
			sets[tmp.SetID] = tmp
		}
	}
	if len(moduleIDs) > 0 {
		tmps := []metadata.ModuleInst{}
		err = getCache(cli.Instance, common.BKTableNameBaseModule, common.BKModuleIDField, moduleIDs, &tmps)
		if err != nil && !cli.Instance.IsNotFoundErr(err) {
			blog.Errorf("SearchIdentifier error:%s", err.Error())
			resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: defErr.New(common.CCErrObjectSelectIdentifierFailed, err.Error())})
			return
		}
		for _, tmp := range tmps {
			modules[tmp.ModuleID] = tmp
			appmodulenames[tmp.BizID] = append(appmodulenames[tmp.BizID], tmp.ModuleName)
		}
	}
	if len(bizIDs) > 0 {
		tmps := []metadata.BizInst{}
		err = getCache(cli.Instance, common.BKTableNameBaseApp, common.BKAppIDField, bizIDs, &tmps)
		if err != nil && !cli.Instance.IsNotFoundErr(err) {
			blog.Errorf("SearchIdentifier error:%s", err.Error())
			resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: defErr.New(common.CCErrObjectSelectIdentifierFailed, err.Error())})
			return
		}
		for _, tmp := range tmps {
			bizs[tmp.BizID] = tmp
		}
	}
	if len(cloudIDs) > 0 {
		tmps := []metadata.CloudInst{}
		err = getCache(cli.Instance, common.BKTableNameBasePlat, common.BKCloudIDField, cloudIDs, &tmps)
		if err != nil && !cli.Instance.IsNotFoundErr(err) {
			blog.Errorf("SearchIdentifier error:%s", err.Error())
			resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: defErr.New(common.CCErrObjectSelectIdentifierFailed, err.Error())})
			return
		}
		for _, tmp := range tmps {
			clouds[tmp.CloudID] = tmp
		}
	}

	appmodulename2ProcIDs := map[string][]int64{}
	if len(appmodulenames) > 0 {
		for appID, modulenames := range appmodulenames {
			proc2modules := []metadata.ProcessModule{}
			cond := cccondition.CreateCondition().Field(common.BKAppIDField).Eq(appID).Field(common.BKModuleNameField).In(modulenames)
			err = cli.Instance.GetMutilByCondition(common.BKTableNameProcModule, nil, cond.ToMapStr(), &proc2modules, "", -1, -1)
			if err != nil && !cli.Instance.IsNotFoundErr(err) {
				blog.Errorf("SearchIdentifier error:%s", err.Error())
				resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: defErr.New(common.CCErrObjectSelectIdentifierFailed, err.Error())})
				return
			}
			for _, proc2module := range proc2modules {
				key := fmt.Sprintf("%d-%s", proc2module.AppID, proc2module.ModuleName)
				appmodulename2ProcIDs[key] = append(appmodulename2ProcIDs[key], proc2module.ProcessID)
				procIDs = append(procIDs, proc2module.ProcessID)
			}
		}
	}

	if len(procIDs) > 0 {
		tmps := []metadata.HostIdentProcess{}
		err = getCache(cli.Instance, common.BKTableNameBaseProcess, common.BKProcIDField, procIDs, &tmps)
		if err != nil && !cli.Instance.IsNotFoundErr(err) {
			blog.Errorf("SearchIdentifier error:%s", err.Error())
			resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: defErr.New(common.CCErrObjectSelectIdentifierFailed, err.Error())})
			return
		}
		for _, tmp := range tmps {
			procs[tmp.ProcessID] = tmp
		}
	}

	// fill hostidentifier
	for _, inst := range hosts {
		inst.HostIdentModule = map[string]*metadata.HostIdentModule{}
		// fill cloud
		if _, ok := clouds[inst.CloudID]; ok {
			cloud := clouds[inst.CloudID]
			inst.CloudName = cloud.CloudName
		}
		// fill module
		appmodulename2moduleIDs := map[string][]int64{}
		for _, rela := range modulehosts[inst.HostID] {
			mod := &metadata.HostIdentModule{
				SetID:    rela.SetID,
				ModuleID: rela.ModuleID,
				BizID:    rela.AppID,
			}
			inst.HostIdentModule[fmt.Sprint(rela.ModuleID)] = mod

			if biz, ok := bizs[mod.BizID]; ok {
				mod.BizName = biz.BizName
				inst.SupplierAccount = biz.SupplierAccount
				inst.SupplierID = biz.SupplierID
			}

			if set, ok := sets[mod.SetID]; ok {
				mod.SetName = set.SetName
				mod.SetEnv = set.SetEnv
				mod.SetStatus = set.SetStatus
			}

			if module, ok := modules[mod.ModuleID]; ok {
				mod.ModuleName = module.ModuleName
				key := fmt.Sprintf("%d-%s", mod.BizID, mod.ModuleName)
				appmodulename2moduleIDs[key] = append(appmodulename2moduleIDs[key], mod.ModuleID)
			}
		}
		// fill process
		appmoduleProcID2moduleIDs := map[int64][]int64{} // ProcID->moduleIDs
		for key, moduleIDs := range appmodulename2moduleIDs {
			for _, procID := range appmodulename2ProcIDs[key] {
				appmoduleProcID2moduleIDs[procID] = append(appmoduleProcID2moduleIDs[procID], moduleIDs...)
			}
		}

		for procID, moduleIDs := range appmoduleProcID2moduleIDs {
			if proc, ok := procs[procID]; ok {
				proc.BindModules = append(proc.BindModules, moduleIDs...)
				inst.Process = append(inst.Process, proc)
			}
		}

	}

	for _, host := range hosts {
		sort.Sort(metadata.HostIdentProcessSorter(host.Process))
	}

	// returns
	info := make(map[string]interface{})
	info["count"] = len(hosts)
	info["info"] = hosts

	resp.WriteEntity(metadata.Response{BaseResp: metadata.SuccessBaseResp, Data: info})

}

func getCache(db storage.DI, tablename string, idfield string, ids []int64, result interface{}) error {
	condition := map[string]interface{}{
		idfield: map[string]interface{}{
			common.BKDBIN: ids,
		},
	}
	return db.GetMutilByCondition(tablename, nil, condition, result, "", 0, 0)
}
