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
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
)

func (lgc *Logics) FindHostByModuleIDs(ctx context.Context, data *metadata.HostModuleFind, isDetail bool) (*metadata.SearchHost, error) {
	retHostInfo := &metadata.SearchHost{
		Info: make([]mapstr.MapStr, 0),
	}
	hostSearchParam := new(metadata.HostCommonSearch)
	hostSearchParam.Page = data.Page

	condItem := metadata.ConditionItem{Field: common.BKModuleIDField, Operator: common.BKDBIN, Value: data.ModuleIDS}
	moduleFindCond := metadata.SearchCondition{ObjectID: common.BKInnerObjIDModule, Condition: []metadata.ConditionItem{condItem}, Fields: []string{}}
	setFindCond := metadata.SearchCondition{ObjectID: common.BKInnerObjIDSet, Condition: []metadata.ConditionItem{}, Fields: []string{}}
	bizFindCond := metadata.SearchCondition{ObjectID: common.BKInnerObjIDApp, Condition: []metadata.ConditionItem{}, Fields: []string{}}
	bizID := data.AppID
	var err error
	if bizID == 0 {
		// if bk_biz_id and metadata are both not set, then return error
		if data.Metadata == nil {
			return retHostInfo, lgc.ccErr.Errorf(common.CCErrCommParamsIsInvalid, common.BKAppIDField+", "+common.MetadataField)
		}
		bizID, err = data.Metadata.Label.GetBusinessID()
		if nil != err {
			return retHostInfo, err
		}
	}
	hostSearchParam.AppID = bizID

	hostSearchParam.Condition = []metadata.SearchCondition{moduleFindCond, setFindCond, bizFindCond}

	findHostInst := NewSearchHost(ctx, lgc, hostSearchParam)
	findHostInst.ParseCondition()

	err = findHostInst.SearchHostByConds()
	if err != nil {
		return retHostInfo, err
	}
	hostInfoArr, cnt, err := findHostInst.FillTopologyData()
	if err != nil {
		return retHostInfo, err
	}

	retHostInfo.Count = cnt
	if cnt > 0 {
		retHostInfo.Info = hostInfoArr
	}
	return retHostInfo, nil
}
