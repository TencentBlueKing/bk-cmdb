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
	"configcenter/src/scene_server/synchronize_server/app/options"
)

// FetchModel fetch model description struct
type FetchModel struct {
	lgc           *Logics
	syncConfig    *options.ConfigItem
	baseCondition mapstr.MapStr
	//ingoreAppID []int64
}

// NewFetchModel fetch instance struct
func (lgc *Logics) NewFetchModel(syncConfig *options.ConfigItem, conds mapstr.MapStr) *FetchModel {
	return &FetchModel{
		lgc:           lgc,
		syncConfig:    syncConfig,
		baseCondition: conds,
	}
}

// Pretreatment pretreatment handle
func (fm *FetchModel) Pretreatment() errors.CCError {
	return nil
}

// Fetch  get model info
func (fm *FetchModel) Fetch(ctx context.Context, dataClassify string, cond mapstr.MapStr, start, limit int64) (*metadata.InstDataInfo, errors.CCError) {
	input := &metadata.SynchronizeFindInfoParameter{
		Condition: mapstr.New(),
	}
	input.DataClassify = dataClassify
	input.DataType = metadata.SynchronizeOperateDataTypeModel
	input.Limit = uint64(limit)
	input.Start = uint64(start)
	input.Condition.Merge(fm.baseCondition)
	input.Condition.Merge(cond)
	result, err := fm.lgc.synchronizeSrv.SynchronizeSrv(fm.syncConfig.Name).Find(ctx, fm.lgc.header, input)
	blog.V(5).Infof("SynchronizeSrv %s conditon:%#v, rid:%s", fm.syncConfig.Name, input, fm.lgc.rid)
	blog.V(6).Infof("SynchronizeSrv %s result:%#v, rid:%s", fm.syncConfig.Name, result, fm.lgc.rid)
	if err != nil {
		blog.Errorf("Fetch http do error. err:%s,rid:%s", err.Error(), fm.lgc.rid)
		return nil, fm.lgc.ccErr.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !result.Result {
		blog.Errorf("Fetch http reply error. err code:%d,err msg:%s,rid:%s", result.Code, result.ErrMsg, fm.lgc.rid)
		return nil, fm.lgc.ccErr.New(result.Code, result.ErrMsg)
	}
	return &result.Data, nil
}
