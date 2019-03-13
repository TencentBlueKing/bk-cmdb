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
	"strings"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	"configcenter/src/common/errors"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/synchronize_server/app/options"
)

// FetchInst fetch instance struct
type FetchInst struct {
	lgc        *Logics
	syncConfig *options.ConfigItem
	//ingoreAppID []int64
	baseConds mapstr.MapStr
	appIDArr  []int64
}

// NewFetchInst fetch instance struct
func (lgc *Logics) NewFetchInst(syncConfig *options.ConfigItem, baseConds mapstr.MapStr) *FetchInst {
	return &FetchInst{
		lgc:        lgc,
		baseConds:  baseConds,
		syncConfig: syncConfig,
		appIDArr:   make([]int64, 0),
	}
}

// Pretreatment pretreatment handle
func (fi *FetchInst) Pretreatment() errors.CCError {
	conds := condition.CreateCondition()
	if fi.syncConfig.SynchronizeFlag != "" {
		conds.Field(combineMongoDBKey(metadata.BKMetadata, common.MetaDataSynchronizeFlagField)).Eq(fi.syncConfig.SynchronizeFlag)
	}
	if len(fi.syncConfig.SupplerAccount) > 0 {
		conds.Field(common.BKOwnerIDField).In(fi.syncConfig.SupplerAccount)
	}
	fi.baseConds = conds.ToMapStr()
	return nil
}

func combineMongoDBKey(keys ...string) string {
	return strings.Join(keys, ".")
}

// Fetch fetch instance data
func (fi *FetchInst) Fetch(ctx context.Context, objID string, start, limit int64) (*metadata.InstDataInfo, errors.CCError) {
	input := &metadata.SynchronizeFindInfoParameter{
		Condition: mapstr.New(),
	}
	input.Limit = uint64(limit)
	input.Start = uint64(start)
	switch objID {
	case common.BKInnerObjIDApp:
		conds := condition.CreateCondition()
		// Unsynchronized resource pool
		if !fi.syncConfig.SyncResource {
			conds.Field(common.BKDefaultField).Eq(0)
		}
		if len(fi.syncConfig.AppNames) > 0 {
			if fi.syncConfig.WiteList {
				conds.Field(common.BKAppNameField).In(fi.syncConfig.AppNames)
			} else {
				conds.Field(common.BKAppNameField).NotIn(fi.syncConfig.AppNames)
			}
		}
		input.Condition.Merge(conds.ToMapStr())
	case common.BKInnerObjIDSet:
		input.Condition.Merge(fi.getAppCondition())
	case common.BKInnerObjIDModule:
		input.Condition.Merge(fi.getAppCondition())
	case common.BKInnerObjIDHost:
	case common.BKInnerObjIDProc:
		input.Condition.Merge(fi.getAppCondition())
	// not synchronized
	case common.BKInnerObjIDConfigTemp:
		return nil, nil
	// not synchronized
	case common.BKInnerObjIDTempVersion:
		return nil, nil
	case common.BKInnerObjIDPlat:
		// object get all model
	default:

	}
	input.Condition.Merge(fi.baseConds)
	input.DataClassify = objID
	input.DataType = metadata.SynchronizeOperateDataTypeInstance

	result, err := fi.lgc.synchronizeSrv.SynchronizeSrv(fi.syncConfig.Name).Find(ctx, fi.lgc.header, input)
	if err != nil {
		blog.Errorf("fetchInst http do error. err:%s,objID:%s,input:%#v,rid:%s", err.Error(), objID, input, fi.lgc.rid)
		return nil, fi.lgc.ccErr.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !result.Result {
		blog.Errorf("fetchInst http reply error. err code:%d,err msg:%s,objID:%s,input:%#v,rid:%s", result.Code, result.ErrMsg, objID, input, fi.lgc.rid)
		return nil, fi.lgc.ccErr.New(result.Code, result.ErrMsg)
	}
	return &result.Data, nil

}

// SetAppIDArr set app id
func (fi *FetchInst) SetAppIDArr(appIDArr []int64) {
	fi.appIDArr = appIDArr
}

func (fi *FetchInst) getAppCondition() mapstr.MapStr {
	conds := condition.CreateCondition()
	if len(fi.appIDArr) > 0 {
		conds.Field(common.BKAppIDField).In(fi.appIDArr)
	}
	return conds.ToMapStr()
}
