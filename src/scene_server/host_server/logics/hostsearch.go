/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 Tencent. All rights reserved.
 * Licensed under the MIT License (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 * We undertake not to change the open source license (MIT license) applicable
 * to the current version of the project delivered to anyone in the future.
 */

package logics

import (
	"context"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"configcenter/pkg/filter"
	"configcenter/src/ac"
	acMeta "configcenter/src/ac/meta"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	hostParse "configcenter/src/common/paraparse"
	"configcenter/src/common/util"
)

// SearchHost query the host flag: true for business access, false for no authentication in this function.
func (lgc *Logics) SearchHost(kit *rest.Kit, data *metadata.HostCommonSearch, flag bool) (
	*metadata.SearchHost, *metadata.BaseResp, error) {
	searchHostInst := NewSearchHost(kit, lgc, data)
	if err := searchHostInst.ParseCondition(); err != nil {
		return nil, nil, err
	}
	retHostInfo := &metadata.SearchHost{
		Info: make([]mapstr.MapStr, 0),
	}

	if flag {
		// for scenarios that require business access authentication,
		// the user must pass the business attribute field.
		if !searchHostInst.GetTopologyBizFlag() {
			return nil, nil, kit.CCError.CCError(common.CCErrCommParamsInvalid)
		}
		if authResp, authorized := searchHostInst.AuthorizeSearchHost(); !authorized {
			return retHostInfo, authResp, ac.NoAuthorizeError
		}
	}

	err := searchHostInst.SearchHostByConds()
	if err != nil {
		return retHostInfo, nil, err
	}
	hostInfoArr, cnt, err := searchHostInst.FillTopologyData()
	if err != nil {
		return retHostInfo, nil, err
	}

	retHostInfo.Count = cnt
	if cnt > 0 {
		retHostInfo.Info = hostInfoArr
	}
	return retHostInfo, nil, nil
}

// SearchHostForResource in the host pool view host scenario, authentication is performed through
// the host pool host view permission
func (lgc *Logics) SearchHostForResource(kit *rest.Kit, data *metadata.HostCommonSearch) (*metadata.SearchHost, error) {
	searchHostInst := NewSearchHost(kit, lgc, data)
	if err := searchHostInst.ParseCondition(); err != nil {
		return nil, err
	}
	retHostInfo := &metadata.SearchHost{
		Info: make([]mapstr.MapStr, 0),
	}

	err := searchHostInst.SearchHostByConds()
	if err != nil {
		return retHostInfo, err
	}
	hostInfoArr, cnt, err := searchHostInst.FillTopologyData()
	if err != nil {
		return retHostInfo, err
	}

	retHostInfo.Count = cnt
	if cnt > 0 {
		retHostInfo.Info = hostInfoArr
	}
	return retHostInfo, nil
}

type searchHostConds struct {
	hostCond      metadata.SearchCondition
	appCond       metadata.SearchCondition
	setCond       metadata.SearchCondition
	moduleCond    metadata.SearchCondition
	mainlineCond  metadata.SearchCondition
	platCond      metadata.SearchCondition
	objectCondMap map[string][]metadata.ConditionItem
}

type searchHostTopologyShowSection struct {
	app    bool
	set    bool
	module bool
}

type searchHostModuleHostConfig struct {
	appIDArr      []int64
	moduleIDArr   []int64
	setIDArr      []int64
	asstHostIDArr []int64
}

type searchHostIDArr struct {
	moduleHostConfig searchHostModuleHostConfig
}

type searchHostInfoMapCache struct {
	appInfoMap           map[int64]mapstr.MapStr
	setInfoMap           map[int64]mapstr.MapStr
	moduleInfoMap        map[int64]mapstr.MapStr
	cloudAsstNameInfoMap map[int64]*InstNameAsst
}

type hostInfoStruct struct {
	hostID   int64
	hostInfo mapstr.MapStr
}

type searchHost struct {
	kit             *rest.Kit
	lgc             *Logics
	pheader         http.Header
	hostSearchParam *metadata.HostCommonSearch
	//  this part need to be displayed?
	topoShowSection searchHostTopologyShowSection

	conds searchHostConds
	// search end, condition not dsetAppConfigata
	noData       bool
	idArr        searchHostIDArr
	hostInfoArr  []hostInfoStruct // int64 is hostID
	cacheInfoMap searchHostInfoMapCache
	totalHostCnt int

	paged bool

	searchedHostIDs []int64
	searchCloudIDs  []int64

	ccErr errors.DefaultCCErrorIf
	ccRid string
	ctx   context.Context
}

type setLevelInfo struct {
	setName       string
	setID         int64
	moduleInfoMap map[int64]int64
}

type appLevelInfo struct {
	appName    string
	appID      int64
	setInfoMap map[int64]*setLevelInfo
}

// searchHostInterface Too many methods, hiding private methods
type searchHostInterface interface {
	ParseCondition() error
	SearchHostByConds() errors.CCError
	AuthorizeSearchHost() (*metadata.BaseResp, bool)
	FillTopologyData() ([]mapstr.MapStr, int, errors.CCError)
	GetTopologyBizFlag() bool
}

// NewSearchHost TODO
func NewSearchHost(kit *rest.Kit, lgc *Logics, hostSearchParam *metadata.HostCommonSearch) searchHostInterface {
	sh := &searchHost{
		kit:             kit,
		lgc:             lgc,
		pheader:         kit.Header,
		hostSearchParam: hostSearchParam,
		idArr:           searchHostIDArr{},
		ccRid:           kit.Rid,
		ccErr:           kit.CCError,
		ctx:             kit.Ctx,
	}

	sh.conds.objectCondMap = make(map[string][]metadata.ConditionItem)

	return sh
}

// ParseCondition TODO
func (sh *searchHost) ParseCondition() error {

	for _, object := range sh.hostSearchParam.Condition {
		if err := sh.validateObjCond(&object); err != nil {
			return err
		}

		if object.ObjectID == common.BKInnerObjIDHost {
			sh.conds.hostCond = object
		} else if object.ObjectID == common.BKInnerObjIDSet {
			sh.conds.setCond = object
			sh.topoShowSection.set = true
		} else if object.ObjectID == common.BKInnerObjIDModule {
			sh.conds.moduleCond = object
			sh.topoShowSection.module = true
		} else if object.ObjectID == common.BKInnerObjIDApp {
			sh.conds.appCond = object
			// 只有关于biz的条件大于0才被视为通过业务进行查询主机，
			// 因为不通过业务进行查询主机也需要返回主机的业务信息，
			// 这里如果不区分场景直接置为标记后面的流程走不通。
			if len(object.Condition) > 0 {
				sh.topoShowSection.app = true
			}
		} else if object.ObjectID == common.BKInnerObjIDObject {
			sh.conds.mainlineCond = object
		} else if object.ObjectID == common.BKInnerObjIDPlat {
			sh.conds.platCond = object
		} else {
			sh.conds.objectCondMap[object.ObjectID] = object.Condition
		}
	}
	sh.hostSearchParam.Condition = nil

	sh.tryParseAppID()

	return nil
}

func (sh *searchHost) validateObjCond(objCond *metadata.SearchCondition) error {
	if len(objCond.Condition) == 0 {
		return nil
	}

	// mainline object condition only allows condition of 'bk_inst_id $eq value' form, and is validated elsewhere
	if objCond.ObjectID == common.BKInnerObjIDObject {
		return nil
	}

	attributes, err := sh.lgc.SearchObjectAttributes(sh.kit, sh.hostSearchParam.AppID, objCond.ObjectID)
	if err != nil {
		blog.Errorf("search %s obj attr for validation failed, err: %v, rid: %s", objCond.ObjectID, err, sh.kit.Rid)
		return err
	}

	attributeMap := make(map[string]string)
	for _, attribute := range attributes {
		attributeMap[attribute.PropertyID] = attribute.PropertyType
	}

	switch objCond.ObjectID {
	case common.BKInnerObjIDApp:
		attributeMap[common.BKAppIDField] = common.FieldTypeInt
		attributeMap[common.BKDefaultField] = common.FieldTypeInt
	case common.BKInnerObjIDSet:
		attributeMap[common.BKSetIDField] = common.FieldTypeInt
		attributeMap[common.BKDefaultField] = common.FieldTypeInt
		attributeMap[common.BKParentIDField] = common.FieldTypeInt
	case common.BKInnerObjIDModule:
		attributeMap[common.BKModuleIDField] = common.FieldTypeInt
		attributeMap[common.BKDefaultField] = common.FieldTypeInt
		attributeMap[common.BKParentIDField] = common.FieldTypeInt
	case common.BKInnerObjIDHost:
		attributeMap[common.BKHostIDField] = common.FieldTypeInt
		attributeMap[common.BKCloudIDField] = common.FieldTypeInt
	}

	for _, cond := range objCond.Condition {
		err = sh.validateCondItem(objCond.ObjectID, cond, attributeMap)
		if err != nil {
			return err
		}
	}
	return nil
}

func (sh *searchHost) validateCondItem(objID string, cond metadata.ConditionItem, attrMap map[string]string) error {
	attrType, exists := attrMap[cond.Field]
	if !exists {
		blog.Errorf("%s condition item field %s not exists, rid: %s", objID, cond.Field, sh.kit.Rid)
		return fmt.Errorf("condition field %s not exists", cond.Field)
	}

	supportedOpMap, exists := attrTypeSupportedOpMap[attrType]
	if !exists {
		blog.Errorf("%s condition item field %s attr type %s is invalid, rid: %s", objID, cond.Field, attrType,
			sh.kit.Rid)
		return fmt.Errorf("condition field %s is invalid", cond.Field)
	}

	_, exists = supportedOpMap[cond.Operator]
	if !exists {
		blog.Errorf("%s condition item field %s op %s is invalid, rid: %s", objID, cond.Field, cond.Operator,
			sh.kit.Rid)
		return fmt.Errorf("condition operator %s is invalid", cond.Operator)
	}

	switch cond.Operator {
	case common.BKDBIN, common.BKDBNIN:
		valueArr, ok := cond.Value.([]interface{})
		if !ok {
			blog.Errorf("%s condition item field %s op %s value(%+v) is invalid, rid: %s", objID, cond.Field,
				cond.Operator, cond.Value, sh.kit.Rid)
			return fmt.Errorf("operator %s only support array value", cond.Operator)
		}

		for _, value := range valueArr {
			if err := sh.validCondValueType(attrType, value); err != nil {
				blog.Errorf("%s condition item field %s array value(%+v) is invalid, err: %v, rid: %s", objID,
					cond.Field, value, err, sh.kit.Rid)
				return err
			}
		}
	default:
		if err := sh.validCondValueType(attrType, cond.Value); err != nil {
			blog.Errorf("%s condition item field %s value(%+v) is invalid, err: %v, rid: %s", objID, cond.Field,
				cond.Value, err, sh.kit.Rid)
			return err
		}
	}
	return nil
}

var attrTypeSupportedOpMap = make(map[string]map[string]struct{})

func init() {
	commonOps := []string{common.BKDBIN, common.BKDBNIN, common.BKDBEQ, common.BKDBNE}

	commonOpMap := make(map[string]struct{})
	for _, op := range commonOps {
		commonOpMap[op] = struct{}{}
	}

	commonAttrTypes := []string{common.FieldTypeBool, common.FieldTypeOrganization}
	for _, attrType := range commonAttrTypes {
		attrTypeSupportedOpMap[attrType] = commonOpMap
	}

	strOpMap := make(map[string]struct{})
	for _, op := range append(commonOps, common.BKDBLIKE, string(filter.Contains), string(filter.ContainsSensitive)) {
		strOpMap[op] = struct{}{}
	}

	strAttrTypes := []string{common.FieldTypeSingleChar, common.FieldTypeEnum, common.FieldTypeEnumMulti,
		common.FieldTypeDate, common.FieldTypeTime, common.FieldTypeLongChar, common.FieldTypeUser,
		common.FieldTypeTimeZone, common.FieldTypeList}
	for _, attrType := range strAttrTypes {
		attrTypeSupportedOpMap[attrType] = strOpMap
	}

	comparableOpMap := make(map[string]struct{})
	for _, op := range append(commonOps, common.BKDBGT, common.BKDBGTE, common.BKDBLT, common.BKDBLTE) {
		comparableOpMap[op] = struct{}{}
	}

	numericAttrTypes := []string{common.FieldTypeInt, common.FieldTypeFloat, common.FieldTypeDate, common.FieldTypeTime}
	for _, attrType := range numericAttrTypes {
		attrTypeSupportedOpMap[attrType] = comparableOpMap
	}
}

func (sh *searchHost) validCondValueType(attrType string, value interface{}) error {
	switch attrType {
	case common.FieldTypeInt, common.FieldTypeFloat, common.FieldTypeOrganization:
		if !util.IsNumeric(value) {
			return fmt.Errorf("%s attribute type only support numeric value", attrType)
		}
	case common.FieldTypeBool:
		if _, ok := value.(bool); !ok {
			return fmt.Errorf("%s attribute type only support bool value", attrType)
		}
	default:
		if _, ok := value.(string); !ok {
			return fmt.Errorf("%s attribute type only support string value", attrType)
		}
	}
	return nil
}

// AuthorizeSearchHost query the host according to the business access authority.
func (sh *searchHost) AuthorizeSearchHost() (*metadata.BaseResp, bool) {

	if !sh.lgc.AuthManager.Enabled() {
		return nil, true
	}

	bizMap := make(map[int64]struct{})
	if sh.hostSearchParam.AppID > 0 {
		bizMap[sh.hostSearchParam.AppID] = struct{}{}
	}

	bizCond := make(map[string]interface{})
	if len(sh.conds.appCond.Condition) > 0 {
		for _, cond := range sh.conds.appCond.Condition {
			if cond.Field == common.BKAppIDField {
				id, err := util.GetInt64ByInterface(cond.Value)
				if err != nil {
					return &metadata.BaseResp{
						Result: false,
						Code:   common.CCErrCommParamsInvalid,
						ErrMsg: err.Error(),
					}, false
				}
				bizMap[id] = struct{}{}
				continue
			}

			// 这里如果还有条件是非bizID，需要查询对应的业务列表
			bizCond[cond.Field] = map[string]interface{}{
				cond.Operator: cond.Value,
			}
		}

		if len(bizCond) > 0 {

			input := &metadata.QueryCondition{
				Condition:      bizCond,
				Fields:         []string{common.BKAppIDField},
				Page:           metadata.BasePage{Limit: common.BKNoLimit},
				DisableCounter: true,
			}

			result, err := sh.lgc.CoreAPI.CoreService().Instance().ReadInstance(sh.kit.Ctx, sh.kit.Header,
				common.BKInnerObjIDApp, input)
			if err != nil {
				blog.Errorf("get biz info failed, input: %+v, err: %v, rid: %s", input, err, sh.kit.Rid)
				return &metadata.BaseResp{
					Result: false,
					Code:   common.CCErrCommParamsValueInvalidError,
					ErrMsg: err.Error(),
				}, false
			}

			for _, biz := range result.Info {
				bizID, err := util.GetInt64ByInterface(biz[common.BKAppIDField])
				if err != nil {
					blog.Errorf("get biz id failed, biz: %v, err: %v, rid: %s", biz, err, sh.kit.Rid)
					return &metadata.BaseResp{
						Result: false,
						Code:   common.CCErrCommParamsValueInvalidError,
						ErrMsg: common.BKAppIDField,
					}, false
				}
				bizMap[bizID] = struct{}{}
			}
		}
	}

	bizIDs := make([]int64, 0)
	for id := range bizMap {
		bizIDs = append(bizIDs, id)
	}

	if err := sh.lgc.AuthManager.AuthorizeByInstanceID(sh.kit.Ctx, sh.kit.Header, acMeta.ViewBusinessResource,
		common.BKInnerObjIDApp, bizIDs...); err != nil {
		authResources := make([]acMeta.ResourceAttribute, len(bizIDs))
		for id, biz := range bizIDs {
			authResources[id] = acMeta.ResourceAttribute{Basic: acMeta.Basic{InstanceID: biz, Type: acMeta.Business,
				Action: acMeta.ViewBusinessResource}}
		}

		base, flag := sh.lgc.AuthManager.Authorize(sh.kit, authResources...)
		blog.Errorf("authorize failed, bizID: %v, err: %v, rid: %s", bizIDs, err, sh.kit.Rid)
		return base, flag
	}

	return nil, true
}

// SearchHostByConds TODO
func (sh *searchHost) SearchHostByConds() errors.CCError {

	err := sh.searchByTopo()
	if err != nil {
		return err
	}
	if sh.noData {
		return nil
	}
	err = sh.searchByHostConds()
	if err != nil {
		return err
	}
	if sh.noData {
		return nil
	}

	return nil

}
func (sh *searchHost) GetTopologyBizFlag() bool {
	return sh.topoShowSection.app
}

// FillTopologyData TODO
func (sh *searchHost) FillTopologyData() ([]mapstr.MapStr, int, errors.CCError) {

	if sh.noData {
		return nil, 0, nil
	}

	queryCond := metadata.HostModuleRelationRequest{
		HostIDArr: sh.searchedHostIDs,
		Fields:    []string{common.BKAppIDField, common.BKSetIDField, common.BKModuleIDField, common.BKHostIDField},
	}
	sh.searchedHostIDs = nil
	mhconfig, err := sh.lgc.GetHostRelations(sh.kit, queryCond)
	if err != nil {
		return nil, 0, err
	}
	appIDArr := make([]int64, 0)
	setIDArr := make([]int64, 0)
	moduleIDArr := make([]int64, 0)

	hostAppSetModuleConfig := make(map[int64]map[int64]*appLevelInfo, 0)

	blog.V(5).Infof("get modulehostconfig map:%v, rid:%s", mhconfig, sh.ccRid)
	for _, mh := range mhconfig {
		hostAppInfoLevelInst, ok := hostAppSetModuleConfig[mh.HostID]
		if !ok {
			hostAppInfoLevelInst = make(map[int64]*appLevelInfo, 0)
			hostAppSetModuleConfig[mh.HostID] = hostAppInfoLevelInst
		}

		appInfoLevelInst, ok := hostAppInfoLevelInst[mh.AppID]
		if !ok {
			appInfoLevelInst = &appLevelInfo{
				setInfoMap: make(map[int64]*setLevelInfo, 0),
			}
			hostAppInfoLevelInst[mh.AppID] = appInfoLevelInst
		}
		setInfoLevleInst, ok := appInfoLevelInst.setInfoMap[mh.SetID]
		if !ok {
			setInfoLevleInst = &setLevelInfo{
				moduleInfoMap: make(map[int64]int64, 0),
			}
			appInfoLevelInst.setInfoMap[mh.SetID] = setInfoLevleInst
		}
		setInfoLevleInst.moduleInfoMap[mh.ModuleID] = mh.ModuleID

		appIDArr = append(appIDArr, mh.AppID)
		setIDArr = append(setIDArr, mh.SetID)
		moduleIDArr = append(moduleIDArr, mh.ModuleID)
	}
	appInfoMap, err := sh.fetchTopoAppCacheInfo(appIDArr)
	if err != nil {
		return nil, 0, err
	}
	sh.cacheInfoMap.appInfoMap = appInfoMap

	setInfoMap, err := sh.fetchTopoSetCacheInfo(setIDArr)
	if err != nil {
		return nil, 0, err
	}
	sh.cacheInfoMap.setInfoMap = setInfoMap

	moduleInfoMap, err := sh.fetchTopoModuleCacheInfo(moduleIDArr)
	if err != nil {
		return nil, 0, err
	}
	sh.cacheInfoMap.moduleInfoMap = moduleInfoMap

	cloudAsstNameInfoMap, err := sh.fetchHostCloudCacheInfo()
	if err != nil {
		return nil, 0, err
	}
	sh.cacheInfoMap.cloudAsstNameInfoMap = cloudAsstNameInfoMap

	result := make([]mapstr.MapStr, 0)
	for _, hostInfoItem := range sh.hostInfoArr {
		searchHostItem := mapstr.New()
		levelInfo, ok := hostAppSetModuleConfig[hostInfoItem.hostID]
		if !ok {
			continue
		}
		searchHostItem = sh.fillHostAppInfo(levelInfo, searchHostItem)
		searchHostItem = sh.fillHostSetInfo(levelInfo, searchHostItem)
		searchHostItem = sh.fillHostModuleInfo(levelInfo, searchHostItem)
		hostInfo := sh.fillHostCloudInfo(hostInfoItem.hostInfo, searchHostItem)
		searchHostItem.Set(common.BKInnerObjIDHost, hostInfo)

		result = append(result, searchHostItem)
	}
	return result, sh.totalHostCnt, nil

}

/* ** fill host cloud info  ** */

func (sh *searchHost) fillHostCloudInfo(hostInfo, searchHostItem mapstr.MapStr) mapstr.MapStr {
	clouldID, err := hostInfo.Int64(common.BKCloudIDField)
	if err != nil {
		blog.Warnf("search host fillHostCloudInfo host get cloud id error, hostinfo:%+v, error:%s, rid:%s", hostInfo,
			err.Error(), sh.ccRid)
		hostInfo.Set(common.BKCloudIDField, make([]InstNameAsst, 0))
		return hostInfo
	}
	instAsst, ok := sh.cacheInfoMap.cloudAsstNameInfoMap[clouldID]
	if !ok {
		blog.Warnf("search host fillHostCloudInfo host  cloud id not found, cloud id:%d, hostinfo:%+v, rid:%s",
			clouldID, hostInfo, sh.ccRid)
		hostInfo.Set(common.BKCloudIDField, make([]InstNameAsst, 0))
		return hostInfo
	}
	hostInfo.Set(common.BKCloudIDField, []*InstNameAsst{instAsst})
	return hostInfo
}

func (sh *searchHost) fetchHostCloudCacheInfo() (map[int64]*InstNameAsst, errors.CCError) {

	queryInput := &metadata.QueryCondition{}
	queryInput.Condition = mapstr.MapStr{
		common.BKCloudIDField: mapstr.MapStr{
			common.BKDBIN: sh.searchCloudIDs,
		},
	}
	result, err := sh.lgc.CoreAPI.CoreService().Instance().ReadInstance(sh.ctx, sh.pheader, common.BKInnerObjIDPlat,
		queryInput)
	if err != nil {
		blog.Errorf("fetchHostCloudCacheInfo SearchObjects http do error, err:%s,input:%+v,rid:%s", err.Error(),
			queryInput, sh.ccRid)
		return nil, sh.ccErr.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	cloudInfoMap := make(map[int64]*InstNameAsst)
	for _, info := range result.Info {
		asstInst, err := sh.convInstInfoToAssociateInfo(common.BKCloudIDField, common.BKCloudNameField,
			common.BKInnerObjIDPlat, info)
		if err != nil {
			return nil, err
		}
		cloudInfoMap[asstInst.ObjectID] = asstInst
	}
	return cloudInfoMap, nil

}

func (sh *searchHost) convInstInfoToAssociateInfo(instIDKey, instNameKey, objID string,
	instInfo mapstr.MapStr) (*InstNameAsst, errors.CCError) {
	if val, exist := instInfo[instNameKey]; exist {
		asstInst := &InstNameAsst{}
		if name, can := val.(string); can {
			asstInst.Name = name
			asstInst.ObjID = objID
		}
		instID, err := instInfo.Int64(instIDKey)
		if err != nil {
			return nil, sh.ccErr.Errorf(common.CCErrCommInstFieldConvertFail, objID, instIDKey, "int", err.Error())
		}
		asstInst.ID = strconv.FormatInt(instID, 10)
		asstInst.ObjectID = instID
		return asstInst, nil
	}

	return nil, nil
}

/* ** fill host topology data ** */

func (sh *searchHost) fillHostAppInfo(appInfoLevelInst map[int64]*appLevelInfo,
	searchHostItem mapstr.MapStr) mapstr.MapStr {

	appInfoArr := make([]mapstr.MapStr, 0)
	var err error
	// appdata
	for appID, appLevelInfo := range appInfoLevelInst {

		appInfo, mapOk := sh.cacheInfoMap.appInfoMap[appID]
		if mapOk {
			appInfoArr = append(appInfoArr, appInfo)
		}
		appLevelInfo.appID = appID
		appLevelInfo.appName, err = appInfo.String(common.BKAppNameField)
		if err != nil {
			blog.Warnf("hostSearch not found app name, appInfo:%d, rid:%s", appInfo, sh.ccRid)
			continue
		}

	}
	searchHostItem.Set(common.BKInnerObjIDApp, appInfoArr)
	return searchHostItem

}

func (sh *searchHost) fillHostSetInfo(appInfoLevelInst map[int64]*appLevelInfo,
	searchHostItem mapstr.MapStr) mapstr.MapStr {

	setInfoArr := make([]mapstr.MapStr, 0)
	for _, appLevelInfo := range appInfoLevelInst {
		for setID, setLevelInfo := range appLevelInfo.setInfoMap {
			setInfo, isOk := sh.cacheInfoMap.setInfoMap[setID]
			if false == isOk {
				continue
			}

			setName, err := setInfo.String(common.BKSetNameField)
			if nil != err {
				blog.Warnf("hostSearch not found set name, setInfo:%d, rid:%s", setInfo, sh.ccRid)
				continue
			}
			setLevelInfo.setID = setID
			setLevelInfo.setName = setName
			setInfo.Set(common.TopoModuleName, appLevelInfo.appName+SplitFlag+setName)

			setInfoArr = append(setInfoArr, setInfo)
		}

	}

	searchHostItem[common.BKInnerObjIDSet] = setInfoArr
	return searchHostItem
}

func (sh *searchHost) fillHostModuleInfo(appInfoLevelInst map[int64]*appLevelInfo,
	searchHostItem mapstr.MapStr) mapstr.MapStr {

	moduleInfoArr := make([]mapstr.MapStr, 0)
	for _, appLevelInfo := range appInfoLevelInst {
		for _, setLevelInfo := range appLevelInfo.setInfoMap {
			for mdouleID := range setLevelInfo.moduleInfoMap {
				moduleInfo, ok := sh.cacheInfoMap.moduleInfoMap[mdouleID]
				if false == ok {
					blog.V(5).Infof("hostSearch not found module id, moduleID:%d, hostModuleMap:%v, rid:%s", mdouleID,
						sh.cacheInfoMap.moduleInfoMap, sh.ccRid)
					continue
				}

				moduleName, err := moduleInfo.String(common.BKModuleNameField)
				if nil != err {
					blog.V(5).Infof("hostSearch not found module name, moduleID:%d, hostModuleMap:%v, rid:%s", mdouleID,
						sh.cacheInfoMap.moduleInfoMap, sh.ccRid)
					continue
				}
				datacp := make(map[string]interface{})
				for key, val := range moduleInfo {
					datacp[key] = val
				}
				datacp[common.TopoModuleName] = appLevelInfo.appName + SplitFlag + setLevelInfo.setName + SplitFlag +
					moduleName
				moduleInfoArr = append(moduleInfoArr, datacp)
			}
		}

	}

	searchHostItem[common.BKInnerObjIDModule] = moduleInfoArr
	return searchHostItem
}

func (sh *searchHost) fetchTopoAppCacheInfo(appIDArr []int64) (map[int64]mapstr.MapStr, errors.CCError) {

	if nil != sh.conds.appCond.Fields {
		if len(sh.conds.appCond.Fields) != 0 {
			sh.conds.appCond.Fields = append(sh.conds.appCond.Fields, common.BKAppIDField)
			sh.conds.appCond.Fields = append(sh.conds.appCond.Fields, common.BKAppNameField)
		}
		cond := mapstr.New()
		celld := mapstr.New()
		celld.Set(common.BKDBIN, appIDArr)
		cond.Set(common.BKAppIDField, celld)
		return sh.lgc.GetAppMapByCond(sh.kit, sh.conds.appCond.Fields, cond)

	}
	return nil, nil
}

func (sh *searchHost) fetchTopoSetCacheInfo(setIDArr []int64) (map[int64]mapstr.MapStr, errors.CCError) {

	if nil != sh.conds.setCond.Fields {
		exist := util.InArray(common.BKSetIDField, sh.conds.setCond.Fields)
		if !exist && 0 != len(sh.conds.setCond.Fields) {
			sh.conds.setCond.Fields = append(sh.conds.setCond.Fields, common.BKSetIDField)
		}
		cond := mapstr.New()
		celld := mapstr.New()
		celld.Set(common.BKDBIN, setIDArr)
		cond.Set(common.BKSetIDField, celld)
		return sh.lgc.GetSetMapByCond(sh.kit, sh.conds.setCond.Fields, cond)
	}

	return nil, nil
}

func (sh *searchHost) fetchTopoModuleCacheInfo(moduleIDArr []int64) (map[int64]mapstr.MapStr, errors.CCError) {
	if nil != sh.conds.moduleCond.Fields {
		exist := util.InArray(common.BKModuleIDField, sh.conds.moduleCond.Fields)
		if !exist && 0 != len(sh.conds.moduleCond.Fields) {
			sh.conds.moduleCond.Fields = append(sh.conds.moduleCond.Fields, common.BKModuleIDField)
		}
		cond := mapstr.New()
		celld := mapstr.New()
		celld.Set(common.BKDBIN, moduleIDArr)
		cond.Set(common.BKModuleIDField, celld)
		return sh.lgc.GetModuleMapByCond(sh.kit, sh.conds.moduleCond.Fields, cond)

	}

	return nil, nil

}

/* ** The following is the processing of querying data according to conditions. ** */

func (sh *searchHost) searchByTopo() errors.CCError {
	err := sh.searchByApp()
	if err != nil {
		return err
	}
	err = sh.searchByMainline()
	if err != nil {
		return err
	}
	err = sh.searchByModule()
	if err != nil {
		return err
	}
	// Query host information based on associated objects, alternate code
	// sh.searchByAssocation()
	err = sh.searchByPlatCondition()
	if err != nil {
		return err
	}
	return nil
}

func (sh *searchHost) searchByPlatCondition() errors.CCError {
	if sh.noData {
		return nil
	}
	if len(sh.conds.platCond.Condition) > 0 {
		instIDArr, err := sh.lgc.GetObjectInstByCond(sh.kit, common.BKInnerObjIDPlat, sh.conds.platCond.Condition)
		if err != nil {
			return err
		}
		if len(instIDArr) == 0 {
			sh.noData = true
			return nil
		}
		sh.conds.platCond.Condition = nil
		sh.conds.hostCond.Condition = append(sh.conds.hostCond.Condition, metadata.ConditionItem{
			Field:    common.BKCloudIDField,
			Operator: common.BKDBIN,
			Value:    instIDArr,
		})
	}

	return nil
}

func (sh *searchHost) searchByApp() errors.CCError {
	if sh.noData {
		return nil
	}
	if len(sh.conds.appCond.Condition) > 0 || sh.conds.appCond.TimeCondition != nil {
		cond := metadata.ConditionWithTime{
			Condition:     sh.conds.appCond.Condition,
			TimeCondition: sh.conds.appCond.TimeCondition,
		}
		appIDArr, err := sh.lgc.GetAppIDByCond(sh.kit, cond)
		if err != nil {
			return err
		}
		if len(appIDArr) == 0 {
			sh.noData = true
			return nil
		}
		sh.conds.appCond.Condition = nil
		sh.idArr.moduleHostConfig.appIDArr = appIDArr
	}
	return nil
}

func (sh *searchHost) searchByMainline() errors.CCError {
	if sh.noData {
		return nil
	}

	var err error
	setIDArr := make([]int64, 0)
	objSetIDArr := make([]int64, 0)

	// search mainline object by cond
	if len(sh.conds.mainlineCond.Condition) > 0 {
		objSetIDArr, err = sh.lgc.GetSetIDByObjectCond(sh.kit, sh.hostSearchParam.AppID,
			sh.conds.mainlineCond.Condition)
		if err != nil {
			return err
		}
		if len(objSetIDArr) == 0 {
			sh.noData = true
			return nil
		}
		sh.conds.mainlineCond.Condition = nil
		sh.conds.setCond.Condition = append(sh.conds.setCond.Condition, metadata.ConditionItem{
			Field:    common.BKSetIDField,
			Operator: common.BKDBIN,
			Value:    objSetIDArr,
		})
	}
	// search set by appcond
	if len(sh.conds.setCond.Condition) > 0 || sh.conds.setCond.TimeCondition != nil {
		if len(sh.idArr.moduleHostConfig.appIDArr) > 0 {
			sh.conds.setCond.Condition = append(sh.conds.setCond.Condition, metadata.ConditionItem{
				Field:    common.BKAppIDField,
				Operator: common.BKDBIN,
				Value:    sh.idArr.moduleHostConfig.appIDArr,
			})
		}
		cond := metadata.ConditionWithTime{
			Condition:     sh.conds.setCond.Condition,
			TimeCondition: sh.conds.setCond.TimeCondition,
		}
		setIDArr, err = sh.lgc.GetSetIDByCond(sh.kit, cond)
		if err != nil {
			return err
		}
		if len(setIDArr) == 0 {
			sh.noData = true
			return nil
		}
		sh.conds.setCond.Condition = nil
		sh.idArr.moduleHostConfig.setIDArr = setIDArr
	}

	return nil
}

func (sh *searchHost) searchByModule() errors.CCError {
	if sh.noData {
		return nil
	}
	if len(sh.conds.moduleCond.Condition) > 0 || sh.conds.moduleCond.TimeCondition != nil {
		if len(sh.idArr.moduleHostConfig.setIDArr) > 0 {
			sh.conds.moduleCond.Condition = append(sh.conds.moduleCond.Condition, metadata.ConditionItem{
				Field:    common.BKSetIDField,
				Operator: common.BKDBIN,
				Value:    sh.idArr.moduleHostConfig.setIDArr,
			})
		}
		if len(sh.idArr.moduleHostConfig.appIDArr) > 0 {
			sh.conds.moduleCond.Condition = append(sh.conds.moduleCond.Condition, metadata.ConditionItem{
				Field:    common.BKAppIDField,
				Operator: common.BKDBIN,
				Value:    sh.idArr.moduleHostConfig.appIDArr,
			})
		}
		// search module by cond
		cond := metadata.ConditionWithTime{
			Condition:     sh.conds.moduleCond.Condition,
			TimeCondition: sh.conds.moduleCond.TimeCondition,
		}
		moduleIDArr, err := sh.lgc.GetModuleIDByCond(sh.kit, cond)
		if err != nil {
			return err
		}
		if len(moduleIDArr) == 0 {
			sh.noData = true
			return nil
		}
		sh.conds.moduleCond.Condition = nil
		sh.idArr.moduleHostConfig.moduleIDArr = moduleIDArr
	}

	return nil
}

func (sh *searchHost) searchByHostConds() errors.CCError {
	if sh.noData {
		return nil
	}

	err := sh.appendHostTopoConds()
	if err != nil {
		return err
	}
	if sh.noData {
		return nil
	}

	if 0 != len(sh.conds.hostCond.Fields) {
		sh.conds.hostCond.Fields = append(sh.conds.hostCond.Fields, common.BKHostIDField, common.BKCloudIDField)
	}

	condition, err := hostParse.ParseHostParams(sh.conds.hostCond.Condition)
	if err != nil {
		return err
	}

	condition, err = hostParse.ParseHostIPParams(sh.hostSearchParam.Ipv4Ip, sh.hostSearchParam.Ipv6Ip, condition,
		sh.kit.Rid)
	if err != nil {
		return err
	}

	if ipCond, ok := condition[common.BKDBOR].([]map[string]interface{}); ok {
		if cloudIDCond, ok := condition[common.BKCloudIDField].(map[string]interface{}); ok {
			_, inExist := cloudIDCond[common.BKDBIN]
			_, ninExist := cloudIDCond[common.BKDBNIN]
			if inExist || ninExist {
				delete(condition, common.BKCloudIDField)
			}
		}

		cloudAreaCount := len(ipCond)
		if sh.hostSearchParam.Ipv4Ip.Flag == hostParse.IOBOTH {
			cloudAreaCount = cloudAreaCount / 2
		}
		if cloudAreaCount > 50 {
			return errors.NewCCError(common.CCErrHostGetFail, "cloudArea count more than 50")
		}
	}

	return sh.searchHost(condition)
}

func (sh *searchHost) searchHost(condition map[string]interface{}) error {
	query := &metadata.QueryInput{
		Condition:     condition,
		TimeCondition: sh.conds.hostCond.TimeCondition,
		Start:         sh.hostSearchParam.Page.Start,
		Limit:         sh.hostSearchParam.Page.Limit,
		Sort:          sh.hostSearchParam.Page.Sort,
		Fields:        strings.Join(sh.conds.hostCond.Fields, ","),
	}
	sh.conds.hostCond.Fields = nil
	sh.hostSearchParam = nil

	if sh.paged {
		query.Start = 0
	}

	gResult, err := sh.lgc.CoreAPI.CoreService().Host().GetHosts(sh.ctx, sh.pheader, query)
	if err != nil {
		blog.Errorf("get hosts failed, err: %v, rid: %s", err, sh.ccRid)
		return err
	}

	if len(gResult.Info) == 0 {
		sh.noData = true
	}

	if !sh.paged {
		sh.totalHostCnt = gResult.Count
	}

	if sh.searchedHostIDs == nil {
		sh.searchedHostIDs = make([]int64, 0)
	}
	if sh.searchCloudIDs == nil {
		sh.searchCloudIDs = make([]int64, 0)
	}

	for _, host := range gResult.Info {
		hostID, err := util.GetInt64ByInterface(host[common.BKHostIDField])
		if err != nil {
			return err
		}

		cloudID, err := host.Int64(common.BKCloudIDField)
		if err != nil {
			blog.Warnf("hostSearch not found  cloud id in hsot, hostInfo:%d, rid:%s", host, sh.ccRid)
			continue
		}
		sh.searchedHostIDs = append(sh.searchedHostIDs, hostID)
		sh.searchCloudIDs = append(sh.searchCloudIDs, cloudID)

		sh.hostInfoArr = append(sh.hostInfoArr, hostInfoStruct{
			hostID:   hostID,
			hostInfo: host,
		})
	}

	sh.searchedHostIDs = util.IntArrayUnique(sh.searchedHostIDs)
	sh.searchCloudIDs = util.IntArrayUnique(sh.searchCloudIDs)

	return nil
}

func (sh *searchHost) appendHostTopoConds() errors.CCError {
	var moduleHostConfig metadata.DistinctHostIDByTopoRelationRequest
	isAddHostID := false

	if len(sh.idArr.moduleHostConfig.setIDArr) > 0 {
		moduleHostConfig.SetIDArr = sh.idArr.moduleHostConfig.setIDArr
		isAddHostID = true
	}
	if len(sh.idArr.moduleHostConfig.moduleIDArr) > 0 {
		moduleHostConfig.ModuleIDArr = sh.idArr.moduleHostConfig.moduleIDArr
		isAddHostID = true
	}
	if len(sh.conds.objectCondMap) > 0 {
		moduleHostConfig.HostIDArr = sh.idArr.moduleHostConfig.asstHostIDArr
		isAddHostID = true
	}

	if len(sh.idArr.moduleHostConfig.appIDArr) > 0 {
		// already sorted by app id.
		moduleHostConfig.ApplicationIDArr = sh.idArr.moduleHostConfig.appIDArr
		isAddHostID = true
	}

	if !isAddHostID {
		return nil
	}

	var hostIDArr []int64

	respHostIDs, err := sh.lgc.CoreAPI.CoreService().Host().GetDistinctHostIDByTopology(sh.ctx, sh.kit.Header,
		&moduleHostConfig)
	if err != nil {
		blog.Errorf("get hosts failed, err: %v, rid: %s", err, sh.ccRid)
		return err
	}

	sh.totalHostCnt = len(respHostIDs)
	// 当有根据主机实例内容查询的时候的时候，无法在程序中完成分页
	hasHostCond := false
	if len(sh.hostSearchParam.Ipv4Ip.Data) > 0 || len(sh.hostSearchParam.Ipv6Ip.Data) > 0 ||
		len(sh.conds.hostCond.Condition) > 0 || sh.conds.hostCond.TimeCondition != nil {
		hasHostCond = true
	}
	if !hasHostCond && sh.hostSearchParam.Page.Limit > 0 {
		start := sh.hostSearchParam.Page.Start
		limit := start + sh.hostSearchParam.Page.Limit

		uniqHostIDCnt := len(respHostIDs)
		// 如果用户start 设置小于0， 将start 设置为默认值
		if start < 0 {
			start = 0
		}
		if start >= uniqHostIDCnt {
			sh.noData = true
			return nil
		}
		allHostIDsArr := respHostIDs
		sort.Slice(allHostIDsArr, func(i, j int) bool { return allHostIDsArr[i] < allHostIDsArr[j] })
		if uniqHostIDCnt <= limit {
			hostIDArr = allHostIDsArr[start:]
		} else {
			hostIDArr = allHostIDsArr[start:limit]
		}
		sh.paged = true
	} else {
		if len(respHostIDs) == 0 {
			sh.noData = true
			return nil
		}
		hostIDArr = respHostIDs
	}

	cond, mergeErr := MergeHostIDToCond(sh.kit, sh.conds.hostCond.Condition, hostIDArr)
	if mergeErr != nil {
		return mergeErr
	}
	sh.conds.hostCond.Condition = cond

	return nil
}

// searchByAssociation TODO
// Query host information based on associated objects, alternate code
func (sh *searchHost) searchByAssociation() errors.CCError {
	instAsstHostIDArr := make([]int64, 0)
	// search host id by object
	firstCond := true
	if len(sh.conds.objectCondMap) > 0 {
		for objID, objCond := range sh.conds.objectCondMap {
			instIDArr, err := sh.lgc.GetObjectInstByCond(sh.kit, objID, objCond)
			if err != nil {
				return err
			}
			instHostIDArr, err := sh.lgc.GetHostIDByInstID(sh.kit, objID, instIDArr)
			if err != nil {
				return err
			}
			if firstCond {
				instAsstHostIDArr = instHostIDArr
			} else {
				instAsstHostIDArr = util.IntArrIntersection(instAsstHostIDArr, instHostIDArr)
			}
			firstCond = false
		}

	}
	instAsstHostIDArr = util.IntArrayUnique(instAsstHostIDArr)
	if len(sh.conds.objectCondMap) > 0 {
		sh.idArr.moduleHostConfig.asstHostIDArr = instAsstHostIDArr
	}

	return nil

}

func (sh *searchHost) tryParseAppID() {
	// search appID by cond
	if sh.hostSearchParam.AppID != -1 && sh.hostSearchParam.AppID != 0 {
		sh.conds.appCond.Condition = append(sh.conds.appCond.Condition, metadata.ConditionItem{
			Field:    common.BKAppIDField,
			Operator: common.BKDBEQ,
			Value:    sh.hostSearchParam.AppID,
		})
	}
	if sh.hostSearchParam.AppID != 0 && sh.hostSearchParam.AppID != -1 {
		sh.topoShowSection.app = true
	}
}

// MergeHostIDToCond merge host id to host query condition
func MergeHostIDToCond(kit *rest.Kit, conds []metadata.ConditionItem, hostIDs []int64) (
	[]metadata.ConditionItem, error) {

	// 合并两种涞源的根据 host_id 查询的 condition
	// 详情见issue: https://github.com/TencentBlueKing/bk-cmdb/issues/2461
	hostIDConditionExist := false
	for idx, cond := range conds {
		if cond.Field != common.BKHostIDField {
			continue
		}

		// merge two condition
		// {"field": "bk_host_id", "operator": "$eq", "value": 1}
		// {"field": "bk_host_id", "operator": "$eq", "value": [1, 2]}
		// ==> {"field": "bk_host_id", "operator": "", "value": {"$in": [1,2], "$eq": 1}}
		hostIDConditionExist = true
		if cond.Operator != common.BKDBIN {
			// it's somewhat trick here to use common.BKDBEQ as merge operator
			cond = metadata.ConditionItem{
				Field:    common.BKHostIDField,
				Operator: common.BKDBEQ,
				Value: map[string]interface{}{
					cond.Operator: cond.Value,
					common.BKDBIN: hostIDs,
				},
			}
			conds[idx] = cond
		} else {
			// intersection of two array
			value, ok := cond.Value.([]interface{})
			if ok == false {
				blog.Errorf("invalid query condition with $in operator, value must be []int64, but got: %+v, rid: %s",
					cond.Value, kit.Rid)
				return nil, kit.CCError.New(common.CCErrCommParamsIsInvalid, common.BKHostIDField)
			}
			hostIDMap := make(map[int64]bool)
			for _, hostID := range hostIDs {
				hostIDMap[hostID] = true
			}
			shareIDs := make([]int64, 0)
			for _, hostID := range value {
				id, err := util.GetInt64ByInterface(hostID)
				if err != nil {
					blog.Errorf("invalid query condition with $in operator, value must be []int64, but got: %+v, "+
						"rid: %s", cond.Value, kit.Rid)
					return nil, kit.CCError.New(common.CCErrCommParamsIsInvalid, common.BKHostIDField)
				}

				if hostIDMap[id] {
					shareIDs = append(shareIDs, id)
				}
			}
			conds[idx].Value = shareIDs
		}
	}
	if !hostIDConditionExist {
		conds = append(conds, metadata.ConditionItem{
			Field:    common.BKHostIDField,
			Operator: common.BKDBIN,
			Value:    hostIDs,
		})
	}

	return conds, nil
}
