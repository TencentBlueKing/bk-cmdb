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
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/emicklei/go-restful"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/storage/dal"
)

//SearchIdentifier get identifier
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
		setIDs    = []int64{}
		moduleIDs = []int64{}
		bizIDs    = []int64{}
		cloudIDs  = []int64{}
	)

	// caches
	var (
		sets    = map[int64]metadata.SetInst{}
		modules = map[int64]metadata.ModuleInst{}
		bizs    = map[int64]metadata.BizInst{}
		clouds  = map[int64]metadata.CloudInst{}
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

	// fetch hosts
	hosts := []*metadata.HostIdentifier{}
	err = cli.GetHostByCondition(nil, condition, &hosts, "", 0, 0)
	if err != nil {
		blog.Errorf("SearchIdentifier error:%s", err.Error())
		resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: defErr.New(common.CCErrObjectSelectIdentifierFailed, err.Error())})
		return
	}
	for _, host := range hosts {
		relations := []metadata.ModuleHost{}
		condiction := map[string]interface{}{
			common.BKHostIDField: host.HostID,
		}
		blog.Infof("SearchIdentifier relations condition %v ", condiction)
		err = cli.Instance.Table(common.BKTableNameModuleHostConfig).Find(condiction).All(context.Background(), &relations)
		if err != nil {
			blog.Errorf("SearchIdentifier error:%s", err.Error())
			resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: defErr.New(common.CCErrObjectSelectIdentifierFailed, err.Error())})
			return
		}

		host.HostIdentModule = map[string]*metadata.HostIdentModule{}
		for _, rela := range relations {
			host.HostIdentModule[fmt.Sprint(rela.ModuleID)] = &metadata.HostIdentModule{
				SetID:    rela.SetID,
				ModuleID: rela.ModuleID,
				BizID:    rela.AppID,
			}
			setIDs = append(setIDs, rela.SetID)
			moduleIDs = append(moduleIDs, rela.ModuleID)
			bizIDs = append(bizIDs, rela.AppID)
		}
		cloudIDs = append(cloudIDs, host.CloudID)
	}

	blog.Infof("sets: %v, modules: %v, bizs: %v, clouds: %v", setIDs, moduleIDs, bizIDs, cloudIDs)
	// fetch cache
	if len(setIDs) > 0 {
		tmps := []metadata.SetInst{}
		err = getCache(cli.Instance, common.BKTableNameBaseSet, common.BKSetIDField, setIDs, &tmps)
		if err != nil {
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
		if err != nil {
			blog.Errorf("SearchIdentifier error:%s", err.Error())
			resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: defErr.New(common.CCErrObjectSelectIdentifierFailed, err.Error())})
			return
		}
		for _, tmp := range tmps {
			modules[tmp.ModuleID] = tmp
		}
	}
	if len(bizIDs) > 0 {
		tmps := []metadata.BizInst{}
		err = getCache(cli.Instance, common.BKTableNameBaseApp, common.BKAppIDField, bizIDs, &tmps)
		if err != nil {
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
		if err != nil {
			blog.Errorf("SearchIdentifier error:%s", err.Error())
			resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: defErr.New(common.CCErrObjectSelectIdentifierFailed, err.Error())})
			return
		}
		for _, tmp := range tmps {
			clouds[tmp.CloudID] = tmp
		}
	}

	blog.V(3).Infof("sets: %v, modules: %v, bizs: %v, clouds: %v", sets, modules, bizs, clouds)

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

	resp.WriteEntity(metadata.Response{BaseResp: metadata.SuccessBaseResp, Data: info})

}

func getCache(db dal.RDB, tablename string, idfield string, ids []int64, result interface{}) error {
	condition := map[string]interface{}{
		idfield: map[string]interface{}{
			common.BKDBIN: ids,
		},
	}
	return db.Table(tablename).Find(condition).All(context.Background(), result)
}
