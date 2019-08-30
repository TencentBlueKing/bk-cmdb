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
	types "configcenter/src/common/mapstr"
	meta "configcenter/src/common/metadata"
	hutil "configcenter/src/scene_server/host_server/util"
)

func (phpapi *PHPAPI) GetDefaultModules(ctx context.Context, appID int) (types.MapStr, errors.CCError) {

	param := &meta.QueryCondition{
		Condition: mapstr.MapStr{
			common.BKAppIDField:   appID,
			common.BKDefaultField: 1,
		},
		Fields: []string{common.BKSetIDField, common.BKModuleIDField},
	}

	resMap, err := phpapi.getObjByCondition(ctx, param, common.BKInnerObjIDModule)
	if nil != err {
		return nil, err
	}

	blog.V(5).Infof("getDefaultModules complete, res: %+v.rid:%s", resMap, phpapi.rid)

	if len(resMap) == 0 {
		blog.Errorf("GetDefaultModules default module not found, appID:%d, params:%+v, rid:%s", appID, param, phpapi.rid)
		return nil, phpapi.ccErr.Errorf(common.CCErrCommNotFound)
	}

	return resMap[0], nil

}

func (phpapi *PHPAPI) GetHostByIPAndSource(ctx context.Context, innerIP string, platID int64) ([]types.MapStr, errors.CCError) {

	param := &meta.QueryCondition{
		Condition: mapstr.MapStr{
			common.BKHostInnerIPField: innerIP,
			common.BKCloudIDField:     platID,
		},
		Fields: []string{common.BKHostIDField},
	}

	resMap, err := phpapi.getObjByCondition(ctx, param, common.BKInnerObjIDHost)

	if nil != err {
		return nil, err
	}

	blog.V(5).Infof("getHostByIPAndSource res: %+v,rid:%s", resMap, phpapi.rid)

	return resMap, nil
}

func (phpapi *PHPAPI) GetHostByCond(ctx context.Context, param *meta.QueryCondition) ([]types.MapStr, errors.CCError) {
	blog.V(5).Infof("GetHostByCond param:%+v,rid:%s", param, phpapi.rid)
	resMap, err := phpapi.getObjByCondition(ctx, param, common.BKInnerObjIDHost)
	if nil != err {
		return nil, err
	}

	blog.V(5).Infof("getHostByIPArrAndSource res: %+v,rid:%s", resMap, phpapi.rid)
	return resMap, nil
}

// search host helpers
func (phpapi *PHPAPI) GetHostMapByCond(ctx context.Context, condition map[string]interface{}) (map[int64]map[string]interface{}, []int64, errors.CCError) {
	hostMap := make(map[int64]map[string]interface{})
	hostIDArr := make([]int64, 0)

	// build host controller url
	searchParams := &meta.QueryInput{
		Fields:    "",
		Condition: condition,
	}
	res, err := phpapi.logic.CoreAPI.CoreService().Host().GetHosts(ctx, phpapi.header, searchParams)

	if nil != err {
		blog.Errorf("getHostMapByCond error params:%+v, error:%s,rid:%s", condition, err.Error(), phpapi.rid)
		return hostMap, hostIDArr, err
	}

	if false == res.Result {
		blog.Errorf("getHostMapByCond GetHosts http response error, params:%+v, err code:%d, err msg:%s,rid:%s", condition, res.Code, res.ErrMsg, phpapi.rid)
		return nil, nil, phpapi.ccErr.New(res.Code, res.ErrMsg)
	}

	for _, host := range res.Data.Info {
		HostID, err := host.Int64(common.BKHostIDField)
		if nil != err {
			blog.Errorf("getHostMapByCond  hostID not integer, err:%s,input:%s,host:%+v,rid:%s", err.Error(), condition, host, phpapi.rid)
			return nil, nil, phpapi.ccErr.Errorf(common.CCErrCommInstFieldConvertFail, common.BKInnerObjIDHost, common.BKHostIDField, "int", err.Error())
		}

		hostMap[HostID] = host
		hostIDArr = append(hostIDArr, HostID)
	}
	return hostMap, hostIDArr, nil
}

// GetHostDataByConfig  get host info
func (phpapi *PHPAPI) GetHostDataByConfig(ctx context.Context, configData []meta.ModuleHost) ([]mapstr.MapStr, errors.CCError) {
	hostIDArr := make([]int64, 0)
	for _, config := range configData {
		hostIDArr = append(hostIDArr, config.HostID)
	}

	hostMapCondition := map[string]interface{}{
		common.BKHostIDField: map[string]interface{}{
			common.BKDBIN: hostIDArr,
		},
	}

	hostMap, _, err := phpapi.GetHostMapByCond(ctx, hostMapCondition)
	if nil != err {
		return nil, err
	}

	hostData, err := phpapi.SetHostData(ctx, configData, hostMap)
	if nil != err {
		return hostData, err
	}

	return hostData, nil
}

func (phpapi *PHPAPI) GetCustomerPropertyByOwner(ctx context.Context, objType string) ([]meta.Attribute, errors.CCError) {

	blog.V(5).Infof("getCustomerPropertyByOwner start,objType:%s,rid:%s", objType, phpapi.rid)
	opt := hutil.NewOperation().WithOwnerID(phpapi.logic.ownerID).WithObjID(common.BKInnerObjIDHost).WithAttrComm().MapStr()
	searchBody := &meta.QueryCondition{
		Condition: opt,
	}
	res, err := phpapi.logic.CoreAPI.CoreService().Model().ReadModelAttr(ctx, phpapi.header, common.BKInnerObjIDHost, searchBody)
	if nil != err {
		blog.Errorf("GetCustomerPropertyByOwner  http do  error, err:%s,param:%+v,objType:%s, rid:%s", err.Error(), searchBody, objType, phpapi.rid)
		return nil, phpapi.ccErr.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if false == res.Result {
		blog.Errorf("GetCustomerPropertyByOwner  http response  error, err code:%d, err msg:%s,rid:%s", res.Code, res.ErrMsg, phpapi.rid)
		return nil, phpapi.ccErr.New(res.Code, res.ErrMsg)
	}
	customAttrArr := make([]meta.Attribute, 0)
	for _, attr := range res.Data.Info { // hostAttrArr {
		if false == attr.IsPre {
			customAttrArr = append(customAttrArr, attr)
		}
	}
	return customAttrArr, nil
}

func (phpapi *PHPAPI) getObjByCondition(ctx context.Context, dat *meta.QueryCondition, objType string) ([]mapstr.MapStr, errors.CCError) {

	res, err := phpapi.logic.CoreAPI.CoreService().Instance().ReadInstance(ctx, phpapi.header, objType, dat)
	if nil != err {
		blog.Errorf("getObjByCondition  http do  error, err:%s, param:%+v, objType:%s, rid:%s", err.Error(), dat, objType, phpapi.rid)
		return nil, phpapi.ccErr.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if false == res.Result {
		blog.Errorf("getObjByCondition  http response error, err code:%d, err msg:%s,rid:%s", res.Code, res.ErrMsg, phpapi.rid)
		return nil, phpapi.ccErr.New(res.Code, res.ErrMsg)
	}

	return res.Data.Info, nil
}
