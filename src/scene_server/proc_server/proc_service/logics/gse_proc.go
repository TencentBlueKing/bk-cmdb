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
)

func (lgc *Logics) MatchProcessInstance(ctx context.Context, procOp *metadata.ProcessOperate, forward http.Header) (map[string]*metadata.ProcInstanceModel, error) {

	setConds := make(map[string]interface{})
	setConds[common.BKAppIDField] = procOp.ApplicationID
	setIDs, _, err := lgc.matchName(ctx, forward, procOp.SetName, common.BKInnerObjIDSet, common.BKSetIDField, common.BKSetNameField, setConds)
	if nil != err {
		blog.Errorf("MatchProcessInstance error:%s", err.Error())
		return nil, err
	}
	if 0 == len(setIDs) {
		return nil, nil
	}
	moduleConds := make(map[string]interface{}, 0)
	moduleConds[common.BKAppIDField] = procOp.ApplicationID
	moduleConds[common.BKSetIDField] = common.KvMap{common.BKDBIN: setIDs}
	moduleIDs, _, err := lgc.matchName(ctx, forward, procOp.ModuleName, common.BKInnerObjIDModule, common.BKModuleIDField, common.BKModuleNameField, moduleConds)
	if nil != err {
		blog.Errorf("MatchProcessInstance error:%s", err.Error())
		return nil, err
	}
	if 0 == len(moduleIDs) {
		return nil, nil
	}
	conds := make(map[string]interface{}, 0)
	conds[common.BKAppIDField] = procOp.ApplicationID
	conds[common.BKSetIDField] = common.KvMap{common.BKDBIN: setIDs}
	conds[common.BKModuleIDField] = common.KvMap{common.BKDBIN: moduleIDs}
	return lgc.matchID(ctx, forward, procOp.FuncID, procOp.InstanceID, conds)

}

func (lgc *Logics) matchName(ctx context.Context, forward http.Header, match, objID, instIDKey, instNameKey string, conds map[string]interface{}) (instIDs []int64, data map[int64]mapstr.MapStr, err error) {
	parseConds, notParse, err := ParseProcInstMatchCondition(match, true)
	if nil != err {
		blog.Errorf("matchName  parse set regex %s error %s", match, err.Error())
		return nil, nil, err
	}
	query := new(metadata.QueryInput)
	query.Limit = common.BKNoLimit
	if nil != parseConds {
		if nil == conds {
			conds = make(map[string]interface{})
		}
		conds[instNameKey] = parseConds

	}
	query.Condition = conds
	ret, err := lgc.CoreAPI.ObjectController().Instance().SearchObjects(ctx, objID, forward, query)
	if nil != err {
		blog.Errorf("matchName get %s instance error:%s", objID, err.Error())
		return nil, nil, fmt.Errorf("get % error:%s", objID, err.Error())
	}
	if !ret.Result {
		blog.Errorf("matchName get %s instance error:%s", objID, ret.ErrMsg)
		return nil, nil, fmt.Errorf("get % error:%s", objID, ret.ErrMsg)
	}
	var rangeRegexrole *RegexRole
	if notParse {
		rangeRegexrole, err = NewRegexRole(match, true)
		if nil != err {
			blog.Errorf("regex role %s parse error %s", match, err.Error())
			return nil, nil, fmt.Errorf("regex role %s parse error %s", match, err.Error())
		}
	}
	for _, inst := range ret.Data.Info {
		ID, err := inst.Int64(instIDKey)
		if nil != err {
			blog.Errorf("matchName %s info %v get key %s by int error", objID, inst, instIDKey)
			return nil, nil, fmt.Errorf("get %s id by int64 error", objID)
		}
		if notParse {
			name, err := inst.String(instNameKey)
			if nil != err {
				blog.Errorf("matchName %s info %v get key %s by int error", objID, inst, instNameKey)
				return nil, nil, fmt.Errorf("get %s id by int64 error", objID)
			}
			if rangeRegexrole.MatchStr(name) {
				instIDs = append(instIDs, ID)
				data[ID] = inst
			}
		} else {
			instIDs = append(instIDs, ID)
			data[ID] = inst
		}
	}

	return instIDs, data, nil
}

func (lgc *Logics) matchID(ctx context.Context, forward http.Header, funcIDMath, HostIDMatch string, conds map[string]interface{}) (data map[string]*metadata.ProcInstanceModel, err error) {

	funcIDConds, funcIDNotParse, err := ParseProcInstMatchCondition(funcIDMath, false)
	if nil != err {
		blog.Errorf("matchID  parse funcID regex %s error %s", funcIDMath, err.Error())
		return nil, err
	}
	var funcRegexRole *RegexRole
	if funcIDNotParse {
		funcRegexRole, err = NewRegexRole(funcIDMath, false)
		if nil != err {
			blog.Errorf("regex role %s parse error %s", funcIDMath, err.Error())
			return nil, fmt.Errorf("regex role %s parse error %s", funcIDMath, err.Error())
		}
	}
	hostConds, hostIDNotParse, err := ParseProcInstMatchCondition(HostIDMatch, false)
	if nil != err {
		blog.Errorf("matchID  parse host instance id regex %s error %s", HostIDMatch, err.Error())
		return nil, err
	}
	var hostRegexRole *RegexRole
	if funcIDNotParse {
		hostRegexRole, err = NewRegexRole(HostIDMatch, false)
		if nil != err {
			blog.Errorf("regex role %s parse error %s", HostIDMatch, err.Error())
			return nil, fmt.Errorf("regex role %s parse error %s", HostIDMatch, err.Error())
		}
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
		conds[common.BKFuncIDField] = hostConds

	}
	query := new(metadata.QueryInput)
	query.Limit = common.BKNoLimit
	query.Condition = conds
	ret, err := lgc.CoreAPI.ProcController().GetProcInstanceModel(ctx, forward, query)
	if nil != err {
		blog.Errorf("matchID get set instance error:%s", err.Error())
		return nil, fmt.Errorf("get set error:%s", err.Error())
	}
	if !ret.Result {
		blog.Errorf("matchID get set instance error:%s", ret.ErrMsg)
		return nil, fmt.Errorf("get set error:%s", ret.ErrMsg)
	}
	data = make(map[string]*metadata.ProcInstanceModel, 0)
	for _, item := range ret.Data.Info {
		isAdd := true
		if funcIDNotParse {
			isAdd = funcRegexRole.MatchInt64(item.FuncID)
		}
		if hostIDNotParse {
			isAdd = hostRegexRole.MatchInt64(int64(item.HostInstanID))
		}
		if isAdd {
			data[fmt.Sprintf("%d.%d.%d.%d", item.SetID, item.ModuleID, item.FuncID, item.HostID)] = &item
		}
	}

	return data, nil
}
