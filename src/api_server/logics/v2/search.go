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
	"context"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/paraparse"
)

func (lgc *Logics) GetAllHostAndModuleRelation(ctx context.Context) ([]mapstr.MapStr, errors.CCError) {
	defErr := lgc.ccErr

	appConds := &metadata.SearchParams{
		Condition: map[string]interface{}{common.BKOwnerIDField: lgc.ownerID},
	}
	appList, err := lgc.CoreAPI.TopoServer().Instance().InstSearch(ctx, lgc.ownerID, common.BKInnerObjIDApp, lgc.header, appConds)
	if nil != err {
		blog.Errorf("GetAllHostAndModuleRelation http do error,err:%s,input:%#v,rid:%s", err.Error(), appConds, lgc.rid)
		return nil, defErr.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !appList.Result {
		blog.Errorf("GetAllHostAndModuleRelation http reply error, reply:%#v,input:%#v,rid:%s", appList, appConds, lgc.rid)
		return nil, defErr.New(appList.Code, appList.ErrMsg)
	}

	dat := &params.HostCommonSearch{
		//search all host
		Condition: []params.SearchCondition{
			params.SearchCondition{
				Fields:   []string{common.BKAppIDField},
				ObjectID: common.BKInnerObjIDApp,
			},
			params.SearchCondition{
				Fields:   []string{},
				ObjectID: common.BKInnerObjIDHost,
			},
			params.SearchCondition{
				Fields:   []string{common.BKAppIDField, common.BKSetIDField, common.BKModuleIDField, common.BKModuleTypeField, common.BKModuleNameField},
				ObjectID: common.BKInnerObjIDModule,
			},
		},
	}

	hostList, err := lgc.CoreAPI.HostServer().SearchHost(ctx, lgc.header, dat)
	if nil != err {
		blog.Errorf("GetAllHostAndModuleRelation error:%s, rid:%s", err.Error(), lgc.rid)
		return nil, defErr.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !hostList.Result {
		blog.Errorf("GetAllHostAndModuleRelation error:%s, rid:%s", hostList.ErrMsg, lgc.rid)
		return nil, defErr.New(hostList.Code, hostList.ErrMsg)
	}

	hostAppRelation, err := lgc.handleHostAppRelationByHostModuleInfo(hostList.Data.Info, lgc.ownerID)
	if nil != err {
		return nil, err
	}
	appInfoArr := make([]mapstr.MapStr, 0)
	for _, app := range appList.Data.Info {
		appID, err := app.Int64(common.BKAppIDField)
		if nil != err {
			blog.Errorf("GetAllHostAndModuleRelation not found app id, app info:%+v, owner:%s, rid:%s", app, lgc.ownerID, lgc.rid)
			// CCErrCommInstFieldConvFail  convert %s  field %s to %s error %s
			return nil, defErr.Errorf(common.CCErrCommInstFieldConvFail, common.BKInnerObjIDApp, common.BKAppIDField, "int", err.Error())
		}
		hostInfoArr, ok := hostAppRelation[appID]
		if ok {
			app.Set("children", hostInfoArr)
		} else {
			app.Set("children", make([]mapstr.MapStr, 0))
		}
		appInfoArr = append(appInfoArr, app)

	}

	return appInfoArr, nil

}

func (lgc *Logics) handleHostAppRelationByHostModuleInfo(hostInfoArr []mapstr.MapStr, ownerID string) (map[int64][]mapstr.MapStr, errors.CCError) {
	defErr := lgc.ccErr

	hostAppRelation := make(map[int64][]mapstr.MapStr, 0)
	for _, host := range hostInfoArr {
		moduleArr, ok := host[common.BKInnerObjIDModule].([]interface{})
		if !ok {
			blog.Warnf("GetAllHostAndModuleRelation not found module, owner:%s,hostInfo:%#v, rid:%s", ownerID, host, lgc.rid)
			continue
		}
		hostDetail, ok := host[common.BKInnerObjIDHost].(map[string]interface{})
		if !ok {
			blog.Errorf("GetAllHostAndModuleRelation not found host, owner:%s,hostInfo:%#v, rid:%s", ownerID, host, lgc.rid)
			return nil, defErr.Errorf(common.CCErrAPIServerV2DirectErr, "not foud host detail")
		}
		cloudArr, ok := hostDetail[common.BKCloudIDField].([]interface{})
		if ok {
			if 0 < len(cloudArr) {
				cloudMap, err := mapstr.NewFromInterface(cloudArr[0])
				if nil != err {
					blog.Errorf("GetAllHostAndModuleRelation not found host, owner:%s, rid:%s", ownerID, lgc.rid)
					return nil, defErr.Errorf(common.CCErrAPIServerV2DirectErr, "not foud host Source")
				}
				hostDetail[common.BKCloudIDField], err = cloudMap.Int64(common.BKInstIDField)
				if nil != err {
					blog.Errorf("GetAllHostAndModuleRelation not found host, owner:%s, rid:%s", ownerID, lgc.rid)
					return nil, defErr.Errorf(common.CCErrAPIServerV2DirectErr, "not foud host Source")
				}
			}
		}
		for _, module := range moduleArr {
			moduleMap, _ := mapstr.NewFromInterface(module)
			appID, err := moduleMap.Int64(common.BKAppIDField)
			if nil != err {
				blog.Warnf("GetAllHostAndModuleRelation not found app id from module info,, module info:%+v, owner:%s, rid:%s", host, ownerID, lgc.rid)
				continue
			}
			if _, ok := hostAppRelation[appID]; !ok {
				hostAppRelation[appID] = make([]mapstr.MapStr, 0)
			}
			moduleMap.Remove("TopModuleName")
			hostInfo := mapstr.New()
			hostInfo.Merge(hostDetail)
			hostInfo.Merge(moduleMap)
			hostAppRelation[appID] = append(hostAppRelation[appID], hostInfo)

		}

	}
	return hostAppRelation, nil
}
