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
	"configcenter/src/common/condition"
	"configcenter/src/common/errors"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/synchronize_server/app/options"
)

// FetchInst fetch instance struct
type FetchAssociation struct {
	lgc        *Logics
	syncConfig *options.ConfigItem
	//ingoreAppID []int64
	baseConds mapstr.MapStr
	appIDArr  []int64
}

// NewFetchAssociation fetch instance struct
func (lgc *Logics) NewFetchAssociation(syncConfig *options.ConfigItem, conds mapstr.MapStr) *FetchAssociation {
	return &FetchAssociation{
		lgc:        lgc,
		syncConfig: syncConfig,
		baseConds:  conds,
	}
}

// Fetch fetch massociation
func (fa *FetchAssociation) Fetch(ctx context.Context, dataClassify string, start, limit int64) (*metadata.InstDataInfo, errors.CCError) {
	input := &metadata.SynchronizeFindInfoParameter{
		Condition: mapstr.New(),
	}
	input.Limit = uint64(limit)
	input.Start = uint64(start)
	input.DataClassify = dataClassify
	input.DataType = metadata.SynchronizeOperateDataTypeAssociation
	input.Condition.Merge(fa.baseConds)
	switch dataClassify {
	case common.SynchronizeAssociationTypeModelHost:
		input.Condition.Merge(fa.getAppCondition())
	}

	result, err := fa.lgc.synchronizeSrv.SynchronizeSrv(fa.syncConfig.Name).Find(ctx, fa.lgc.header, input)
	blog.V(5).Infof("SynchronizeSrv %s conditon:%#v, rid:%s", fa.syncConfig.Name, input, fa.lgc.rid)
	blog.V(6).Infof("SynchronizeSrv %s result:%#v, rid:%s", fa.syncConfig.Name, result, fa.lgc.rid)
	if err != nil {
		blog.Errorf("FetchModuleHostConfig http do error. err:%s,input:%#v,rid:%s", err.Error(), input, fa.lgc.rid)
		return nil, fa.lgc.ccErr.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !result.Result {
		blog.Errorf("FetchModuleHostConfig http reply error. err code:%d,err msg:%s,input:%#v,rid:%s", result.Code, result.ErrMsg, input, fa.lgc.rid)
		return nil, fa.lgc.ccErr.New(result.Code, result.ErrMsg)
	}

	return &result.Data, nil
}

// SetAppIDArr set app id
func (fa *FetchAssociation) SetAppIDArr(appIDArr []int64) {
	fa.appIDArr = appIDArr
}

func (fa *FetchAssociation) getAppCondition() mapstr.MapStr {
	conds := condition.CreateCondition()
	if len(fa.appIDArr) > 0 {
		conds.Field(common.BKAppIDField).In(fa.appIDArr)
	}
	return conds.ToMapStr()
}
