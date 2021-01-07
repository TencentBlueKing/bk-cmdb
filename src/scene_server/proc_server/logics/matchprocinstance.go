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
	"net/http"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

func (lgc *Logics) MatchProcessInstance(ctx context.Context, params *metadata.MatchProcInstParam, header http.Header) (map[string]*metadata.ProcInstanceModel, error) {
	setConds := mapstr.MapStr{}
	setConds.Set(common.BKAppIDField, params.ApplicationID)
	setIDs, _, err := lgc.matchName(ctx, header, params.SetName, common.BKInnerObjIDSet, common.BKSetIDField, common.BKSetNameField, setConds)
	if nil != err {
		blog.Errorf("MatchProcessInstance error:%s", err.Error())
		return nil, err
	}
	if 0 == len(setIDs) {
		return nil, nil
	}
	moduleConds := make(map[string]interface{}, 0)
	moduleConds[common.BKAppIDField] = params.ApplicationID
	moduleConds[common.BKSetIDField] = mapstr.MapStr{common.BKDBIN: setIDs}
	moduleIDs, _, err := lgc.matchName(ctx, header, params.ModuleName, common.BKInnerObjIDModule, common.BKModuleIDField, common.BKModuleNameField, moduleConds)
	if nil != err {
		blog.Errorf("MatchProcessInstance error:%s", err.Error())
		return nil, err
	}
	if 0 == len(moduleIDs) {
		return nil, nil
	}
	conds := make(map[string]interface{}, 0)
	conds[common.BKAppIDField] = params.ApplicationID
	conds[common.BKSetIDField] = mapstr.MapStr{common.BKDBIN: setIDs}
	conds[common.BKModuleIDField] = mapstr.MapStr{common.BKDBIN: moduleIDs}
	return lgc.matchFuncIDInstID(ctx, header, params.FuncID, params.HostInstanceID, conds)

}

// matchName match module or set by match role
func (lgc *Logics) matchName(ctx context.Context, header http.Header, match, objID, instIDKey, instNameKey string, conds mapstr.MapStr) (instIDs []int64, data map[int64]mapstr.MapStr, err error) {
	defErr := lgc.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(header))
	// paseConds mongodb query condition,
	scopeMatch := NewScopeMatch(match, true)

	parseConds, err := scopeMatch.ParseConds()
	if nil != err {
		blog.Errorf("matchName  parse set regex %s error %s", match, err.Error())
		return nil, nil, defErr.Errorf(common.CCErrCommUtilHandleFail, fmt.Sprintf("parse math %s", match), err.Error())
	}
	query := new(metadata.QueryInput)
	query.Limit = common.BKNoLimit
	if nil != parseConds {
		if nil == conds {
			conds = make(mapstr.MapStr, 0)
		}
		conds[instNameKey] = parseConds
	}
	query.Condition = conds
	ret, err := lgc.CoreAPI.ObjectController().Instance().SearchObjects(ctx, objID, header, query)
	if nil != err {
		blog.Errorf("matchName get %s instance error:%s", objID, err.Error())
		return nil, nil, defErr.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !ret.Result {
		blog.Errorf("matchName get %s instance error:%s", objID, ret.ErrMsg)
		return nil, nil, defErr.New(ret.Code, ret.ErrMsg)
	}
	instIDs = make([]int64, 0)
	data = make(map[int64]mapstr.MapStr, 0)
	for _, inst := range ret.Data.Info {
		ID, err := inst.Int64(instIDKey)
		if nil != err {
			blog.Errorf("matchName %s info %v get key %s by int error", objID, inst, instIDKey)
			return nil, nil, defErr.Errorf(common.CCErrCommInstFieldConvFail, objID, "inst id", "int", err.Error())
		}
		name, err := inst.String(instNameKey)
		if nil != err {
			blog.Errorf("matchName %s info %v get key %s by int error", objID, inst, instNameKey)
			return nil, nil, defErr.Errorf(common.CCErrCommInstFieldConvFail, objID, "inst name", "string", err.Error())
		}
		if scopeMatch.MatchStr(name) {
			instIDs = append(instIDs, ID)
			data[ID] = inst
		}
	}

	return instIDs, data, nil
}

// matchID get the matching rules list by funcID and Instance ID
func (lgc *Logics) matchFuncIDInstID(ctx context.Context, header http.Header, funcIDMath, HostIDMatch string, conds map[string]interface{}) (data map[string]*metadata.ProcInstanceModel, err error) {
	defErr := lgc.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(header))
	funcScopeMatch := NewScopeMatch(funcIDMath, false)
	// funcIDConds mongodb query condition,
	funcIDConds, err := funcScopeMatch.ParseConds()
	if nil != err {
		blog.Errorf("matchID  parse funcID regex %s error %s", funcIDMath, err.Error())
		return nil, defErr.Errorf(common.CCErrCommUtilHandleFail, fmt.Sprintf("parse math %s", funcIDMath), err.Error())
	}
	hostScopeMatch := NewScopeMatch(HostIDMatch, false)
	// hostConds mongodb query condition,
	hostConds, err := hostScopeMatch.ParseConds()
	if nil != err {
		blog.Errorf("matchID  parse host instance id regex %s error %s", HostIDMatch, err.Error())
		return nil, defErr.Errorf(common.CCErrCommUtilHandleFail, fmt.Sprintf("parse math %s", HostIDMatch), err.Error())
	}

	if nil != funcIDConds {
		if nil == conds {
			conds = make(map[string]interface{})
		}
		conds[common.BKFuncIDField] = funcIDConds

	}
	if nil != hostConds {
		if nil == conds {
			conds = make(map[string]interface{})
		}
		conds[common.BKProcinstanceID] = hostConds

	}
	query := new(metadata.QueryInput)
	query.Limit = common.BKNoLimit
	query.Condition = conds
	ret, err := lgc.CoreAPI.ProcController().GetProcInstanceModel(ctx, header, query)
	if nil != err {
		blog.Errorf("matchID get set instance error:%s", err.Error())
		return nil, defErr.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !ret.Result {
		blog.Errorf("matchID get set instance error:%s", ret.ErrMsg)
		return nil, defErr.New(ret.Code, ret.ErrMsg)
	}
	data = make(map[string]*metadata.ProcInstanceModel, 0)
	for idx, item := range ret.Data.Info {
		if funcScopeMatch.MatchInt64(item.FuncID) && hostScopeMatch.MatchInt64(int64(item.HostInstanID)) {
			data[fmt.Sprintf("%d.%d.%d.%d", item.SetID, item.ModuleID, item.FuncID, item.HostInstanID)] = &ret.Data.Info[idx]
		}
	}

	return data, nil
}
