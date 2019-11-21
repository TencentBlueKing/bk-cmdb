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
	"configcenter/src/common/errors"
	"configcenter/src/common/eventclient"
	"configcenter/src/common/metadata"
	"configcenter/src/source_controller/coreservice/core"
	"configcenter/src/storage/dal"

	"gopkg.in/redis.v5"
)

type processOperation struct {
	dbProxy    dal.RDB
	dependence OperationDependence
	eventCli   eventclient.Client
}

// OperationDependence methods definition
type OperationDependence interface {
	CreateProcessInstance(params core.ContextParams, process *metadata.Process) (*metadata.Process, errors.CCErrorCoder)
	TransferHostModuleDep(ctx core.ContextParams, input *metadata.HostsModuleRelation) ([]metadata.ExceptionResult, error)
}

// New create a new model manager instance
func New(dbProxy dal.RDB, dependence OperationDependence, cache *redis.Client) core.ProcessOperation {
	processOps := &processOperation{
		dbProxy:    dbProxy,
		dependence: dependence,
		eventCli:   eventclient.NewClientViaRedis(cache, dbProxy),
	}
	return processOps
}

func (p *processOperation) validateBizID(ctx core.ContextParams, bizID int64) (int64, errors.CCErrorCoder) {
	// avoid unnecessary db query
	if bizID == 0 {
		return 0, ctx.Error.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField)
	}

	// check bizID valid
	filter := map[string]interface{}{
		common.BKAppIDField: bizID,
	}
	count, err := p.dbProxy.Table(common.BKTableNameBaseApp).Find(filter).Count(ctx.Context)
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

	module := &metadata.ModuleInst{}
	filter := map[string]interface{}{
		common.BKModuleIDField: moduleID,
	}
	err := p.dbProxy.Table(common.BKTableNameBaseModule).Find(filter).One(ctx.Context, module)
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
	filter := map[string]interface{}{
		common.BKHostIDField: hostID,
	}
	host := &struct {
		InnerIP string `field:"bk_host_innerip" json:"bk_host_innerip,omitempty" bson:"bk_host_innerip"`
	}{}
	err := p.dbProxy.Table(common.BKTableNameBaseHost).Find(filter).One(ctx.Context, host)
	if nil != err {
		blog.Errorf("validateHostID failed, mongodb failed, table: %s, err: %+v, rid: %s", common.BKTableNameBaseHost, err.Error(), ctx.ReqID)
		return "", ctx.Error.CCErrorf(common.CCErrCommDBSelectFailed)
	}

	return host.InnerIP, nil
}
