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
	"fmt"
	"strconv"

	"configcenter/src/common"
	"configcenter/src/common/auditoplog"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/mapstr"
	meta "configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/validator"
)

// helpers
func (phpapi *PHPAPI) UpdateHostMain(ctx context.Context, hostCondition, data map[string]interface{}, appID int64) (string, errors.CCError) {
	blog.V(5).Infof("hostCondition:%+v, rid:%s", hostCondition, phpapi.rid)

	_, hostIDArr, err := phpapi.GetHostMapByCond(ctx, hostCondition)

	blog.V(5).Infof("hostIDArr:%+v,rid:%s", hostIDArr, phpapi.rid)
	if nil != err {
		return "", err
	}

	lenOfHostIDArr := len(hostIDArr)
	if lenOfHostIDArr != 1 {
		blog.V(5).Infof("GetHostMapByCond condition: %+v, host:%+v,rid:%s", hostCondition, hostIDArr, phpapi.rid)
		return "", phpapi.ccErr.Errorf(common.CCErrCommFieldNotValid)
	}

	ownerID := util.GetOwnerID(phpapi.header)
	valid := validator.NewValidMapWithKeyFields(ownerID, common.BKInnerObjIDHost, []string{common.CreateTimeField, common.LastTimeField, common.BKChildStr, common.BKOwnerIDField}, phpapi.header, phpapi.logic.Engine)
	validErr := valid.ValidMap(data, common.ValidUpdate, hostIDArr[0])
	if nil != validErr {
		blog.Errorf("updateHostMain error: %s, input:%+v,rid:%s", validErr.Error(), data, phpapi.rid)
		return "", phpapi.ccErr.Errorf(common.CCErrCommFieldNotValidFail, validErr.Error())
	}

	configData, err := phpapi.logic.GetConfigByCond(ctx, meta.HostModuleRelationRequest{
		ApplicationID: appID,
		HostIDArr:     []int64{hostIDArr[0]},
	})
	if nil != err {
		return "", fmt.Errorf("GetConfigByCond error:%v, rid:%s", err, phpapi.rid)
	}

	lenOfConfigData := len(configData)
	if lenOfConfigData == 0 {
		blog.Errorf("not expected config length: appID:%d, hostID:%d, rid:%s", appID, hostIDArr[0], phpapi.rid)
		return "", fmt.Errorf("not expected config length: %d", lenOfConfigData)
	}

	hostID := configData[0].HostID

	condition := mapstr.New()
	condition.Set(common.BKHostIDField, hostID)

	param := &meta.UpdateOption{
		Condition: condition,
		Data:      mapstr.NewFromMap(data),
	}

	strHostID := strconv.FormatInt(hostID, 10)
	logContent := phpapi.logic.NewHostLog(ctx, util.GetOwnerID(phpapi.header))
	if err := logContent.WithPrevious(ctx, strHostID, nil); err != nil {
		blog.Warnf("WithPrevious failed, err: %+v, rid: %s", err, phpapi.rid)
	}
	res, err := phpapi.logic.Engine.CoreAPI.CoreService().Instance().UpdateInstance(ctx, phpapi.header, common.BKInnerObjIDHost, param)
	if nil != err {
		blog.Errorf("UpdateHostMain UpdateObject http do error, err:%s,param:%+v,rid:%s", err.Error(), param, phpapi.rid)
		return "", phpapi.ccErr.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if false == res.Result {
		blog.Errorf("UpdateHostMain UpdateObject http response error, err code:%d,err msg:%s, param:%+v,rid:%s", res.Code, res.ErrMsg, param, phpapi.rid)
		return "", phpapi.ccErr.New(res.Code, res.ErrMsg)
	}
	if nil == err && true == res.Result {
		// 操作成功，新加操作日志日志resJs, err := simplejson.NewJson([]byte(res))
		if res.Result {
			if err := logContent.WithCurrent(ctx, strHostID); err != nil {
				blog.Warnf("WithCurrent failed, err: %+v, rid: %s", err, phpapi.rid)
			}
			content := logContent.GetContent(hostID)
			// (id interface{}, Content interface{}, OpDesc string, InnerIP, ownerID, appID, user string, OpType auditoplog.AuditOpType)

			log := meta.SaveAuditLogParams{
				ID:      hostID,
				Model:   common.BKInnerObjIDHost,
				Content: content,
				OpDesc:  "enter ip host",
				OpType:  auditoplog.AuditOpTypeModify,
				BizID:   appID,
			}

			if _, err := phpapi.logic.CoreAPI.CoreService().Audit().SaveAuditLog(context.Background(), phpapi.header, log); err != nil {
				blog.Warnf("SaveAuditLog failed, err: %+v, rid: %s", err, phpapi.rid)
			}
		}
	}

	return "", nil
}

func (phpapi *PHPAPI) AddHost(ctx context.Context, data map[string]interface{}) (int64, errors.CCError) {
	hostID, err := phpapi.addObj(ctx, data, common.BKInnerObjIDHost)
	return hostID, err
}

func (phpapi *PHPAPI) AddModuleHostConfig(ctx context.Context, hostID, appID int64, moduleIDs []int64) errors.CCError {
	data := &meta.HostsModuleRelation{
		ApplicationID: appID,
		HostID:        []int64{hostID},
		ModuleID:      moduleIDs,
	}
	blog.V(5).Infof("addModuleHostConfig start, data: %+v,rid:%s", data, phpapi.rid)

	res, err := phpapi.logic.CoreAPI.CoreService().Host().TransferToNormalModule(ctx, phpapi.header, data)
	if nil != err {
		blog.Errorf("AddModuleHostConfig http do error.err:%s,param:%+v,rid:%s", err.Error(), data, phpapi.rid)
		return phpapi.ccErr.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !res.Result {
		blog.Errorf("AddModuleHostConfig http response error. reply:%#v,param:%+v,rid:%s", res, data, phpapi.rid)
		if len(res.Data) != 0 {
			return phpapi.ccErr.New(int(res.Data[0].Code), res.Data[0].Message)
		}
		return phpapi.ccErr.New(res.Code, res.ErrMsg)
	}
	return nil
}

func (phpapi *PHPAPI) addObj(ctx context.Context, data map[string]interface{}, objType string) (int64, errors.CCError) {
	resMap := make(map[string]interface{})
	input := &meta.CreateModelInstance{
		Data: mapstr.NewFromMap(data),
	}
	resp, err := phpapi.logic.CoreAPI.CoreService().Instance().CreateInstance(ctx, phpapi.header, objType, input)
	if nil != err {
		blog.Errorf("addObj http do error.err:%s,param:%+v,rid:%s", err.Error(), data, phpapi.rid)
		return 0, phpapi.ccErr.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !resp.Result {
		blog.Errorf("addObj http response error.err code:%s,err msg:%s,param:%+v,rid:%s", resp.Code, resp.ErrMsg, input, phpapi.rid)
		return 0, phpapi.ccErr.New(resp.Code, resp.ErrMsg)
	}

	blog.V(5).Infof("add object result : %+v,rid:%s", resMap, phpapi.rid)

	return int64(resp.Data.Created.ID), nil
}

// search host helpers

func (phpapi *PHPAPI) SetHostData(ctx context.Context, moduleHostConfig []meta.ModuleHost, hostMap map[int64]map[string]interface{}) ([]mapstr.MapStr, errors.CCError) {

	// total data
	hostData := make([]mapstr.MapStr, 0)

	appIDArr := make([]int64, 0)
	setIDArr := make([]int64, 0)
	moduleIDArr := make([]int64, 0)

	for _, config := range moduleHostConfig {
		setIDArr = append(setIDArr, config.SetID)
		moduleIDArr = append(moduleIDArr, config.ModuleID)
		appIDArr = append(appIDArr, config.AppID)
	}

	moduleMap, err := phpapi.logic.GetModuleMapByCond(ctx, nil, mapstr.MapStr{
		common.BKModuleIDField: map[string]interface{}{
			common.BKDBIN: moduleIDArr,
		},
	})
	if err != nil {
		return hostData, err
	}

	setMap, err := phpapi.logic.GetSetMapByCond(ctx, nil, map[string]interface{}{
		common.BKSetIDField: map[string]interface{}{
			common.BKDBIN: setIDArr,
		},
	})
	if err != nil {
		blog.Errorf("hostMap GetSetMapByCond  error, err:%s,rid:%s", err.Error(), phpapi.rid)
		return hostData, err
	}

	blog.V(5).Infof("GetAppMapByCond , appIDArr:%v, rid:%s", appIDArr, phpapi.rid)
	appMap, err := phpapi.logic.GetAppMapByCond(ctx, nil, mapstr.MapStr{
		common.BKAppIDField: map[string]interface{}{
			common.BKDBIN: appIDArr,
		},
	})

	if err != nil {
		return hostData, err
	}
	for _, config := range moduleHostConfig {
		hostItem, hasHost := hostMap[config.HostID]
		if !hasHost {
			blog.Errorf("hostMap has not hostID: %d,rid:%s", config.HostID, phpapi.rid)
			continue
		}
		host := mapstr.New()
		host.Merge(hostItem)

		module := moduleMap[config.ModuleID]
		set := setMap[config.SetID]
		app := appMap[config.AppID]

		host[common.BKModuleIDField] = module[common.BKModuleIDField]
		host[common.BKModuleNameField] = module[common.BKModuleNameField]
		host[common.BKSetIDField], _ = set.Int64(common.BKSetIDField) // [common.BKSetIDField]
		host[common.BKSetNameField] = set[common.BKSetNameField]
		host[common.BKAppIDField], _ = app.Int64(common.BKAppIDField) // [common.BKAppIDField]
		host[common.BKAppNameField] = app[common.BKAppNameField]
		host[common.BKModuleTypeField] = module[common.BKModuleTypeField]
		host[common.BKOwnerIDField] = app[common.BKOwnerIDField]
		host[common.BKOperatorField] = module[common.BKOperatorField]
		host[common.BKBakOperatorField] = module[common.BKBakOperatorField]
		host[common.BKSupplierIDField] = app[common.BKSupplierIDField]

		hostData = append(hostData, host)
	}
	return hostData, nil
}
