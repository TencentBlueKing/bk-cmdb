/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.,
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the ",License",); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an ",AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package process

import (
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	"configcenter/src/common/errors"
	"configcenter/src/common/metadata"
	"configcenter/src/source_controller/coreservice/core"
	"configcenter/src/storage/dal"
)

type processOperation struct {
	dbProxy    dal.RDB
	dependence OperationDependence
}

// OperationDependence methods definition
type OperationDependence interface {
	CreateProcessInstance(params core.ContextParams, process *metadata.Process) (*metadata.Process, errors.CCErrorCoder)
	TransferHostModuleDep(ctx core.ContextParams, input *metadata.HostsModuleRelation) ([]metadata.ExceptionResult, error)
}

// New create a new model manager instance
func New(dbProxy dal.RDB, dependence OperationDependence) core.ProcessOperation {
	processOps := &processOperation{
		dbProxy:    dbProxy,
		dependence: dependence,
	}
	return processOps
}

func (p *processOperation) validateBizID(ctx core.ContextParams, md metadata.Metadata) (int64, errors.CCErrorCoder) {
	// extract biz id from metadata
	bizID, err := metadata.BizIDFromMetadata(md)
	if err != nil {
		blog.Errorf("parse biz id from metadata failed, metadata: %+v, err: %+v, rid: %s", md, err, ctx.ReqID)
		return 0, ctx.Error.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField)
	}

	// avoid unnecessary db query
	if bizID == 0 {
		return 0, ctx.Error.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField)
	}

	// check bizID valid
	cond := condition.CreateCondition()
	cond.Field(common.BKAppIDField).Eq(bizID)
	count, err := p.dbProxy.Table(common.BKTableNameBaseApp).Find(cond.ToMapStr()).Count(ctx.Context)
	if nil != err {
		blog.Errorf("mongodb failed, table: %s, err: %+v, rid: %s", common.BKTableNameBaseApp, err, ctx.ReqID)
		return 0, ctx.Error.CCErrorf(common.CCErrCommDBSelectFailed)
	}
	if count < 1 {
		return 0, ctx.Error.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField)
	}

	return bizID, nil
}

func (p *processOperation) validateModuleID(ctx core.ContextParams, moduleID int64) (*metadata.ModuleInst, errors.CCErrorCoder) {
	// avoid unnecessary db query
	if moduleID == 0 {
		return nil, ctx.Error.CCErrorf(common.CCErrCommParamsInvalid, common.BKModuleIDField)
	}

	cond := condition.CreateCondition()
	cond.Field(common.BKModuleIDField).Eq(moduleID)

	module := &metadata.ModuleInst{}
	err := p.dbProxy.Table(common.BKTableNameBaseModule).Find(cond.ToMapStr()).One(ctx.Context, module)
	if nil != err {
		blog.Errorf("validateModuleID failed, mongodb failed, table: %s, err: %+v, rid: %s", common.BKTableNameBaseModule, err, ctx.ReqID)
		return nil, ctx.Error.CCErrorf(common.CCErrCommDBSelectFailed)
	}

	return module, nil
}

func (p *processOperation) validateHostID(ctx core.ContextParams, hostID int64) (string, errors.CCErrorCoder) {
	// avoid unnecessary db query
	if hostID == 0 {
		return "", ctx.Error.CCErrorf(common.CCErrCommParamsInvalid, common.BKHostIDField)
	}

	// check bizID valid
	host := &struct {
		InnerIP string `field:"bk_host_innerip" json:"bk_host_innerip,omitempty" bson:"bk_host_innerip"`
	}{}
	cond := condition.CreateCondition()
	cond.Field(common.BKHostIDField).Eq(hostID)
	err := p.dbProxy.Table(common.BKTableNameBaseHost).Find(cond.ToMapStr()).One(ctx.Context, host)
	if nil != err {
		blog.Errorf("validateHostID failed, mongodb failed, table: %s, err: %+v, rid: %s", common.BKTableNameBaseHost, err.Error(), ctx.ReqID)
		return "", ctx.Error.CCErrorf(common.CCErrCommDBSelectFailed)
	}

	return host.InnerIP, nil
}
