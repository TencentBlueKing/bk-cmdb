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
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/source_controller/coreservice/core"
	"configcenter/src/storage/driver/mongodb"
)

type processOperation struct {
	dependence OperationDependence
}

// OperationDependence methods definition
type OperationDependence interface {
	CreateProcessInstance(kit *rest.Kit, process *metadata.Process) (*metadata.Process, errors.CCErrorCoder)
	CreateProcessInstances(kit *rest.Kit, processes []*metadata.Process) ([]*metadata.Process, errors.CCErrorCoder)
	CreateAuditLogDependence(kit *rest.Kit, logs ...metadata.AuditLog) error
}

// New create a new model manager instance
func New(dependence OperationDependence) core.ProcessOperation {
	processOps := &processOperation{
		dependence: dependence,
	}
	return processOps
}

func (p *processOperation) validateBizID(kit *rest.Kit, bizID int64) (int64, errors.CCErrorCoder) {
	// avoid unnecessary db query
	if bizID == 0 {
		return 0, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField)
	}

	// check bizID valid
	filter := map[string]interface{}{
		common.BKAppIDField: bizID,
	}
	count, err := mongodb.Shard(kit.ShardOpts()).Table(common.BKTableNameBaseApp).Find(filter).Count(kit.Ctx)
	if err != nil {
		blog.Errorf("count biz failed, err: %v, filter: %v, rid: %s", err, filter, kit.Rid)
		return 0, kit.CCError.CCErrorf(common.CCErrCommDBSelectFailed)
	}
	if count < 1 {
		return 0, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField)
	}

	return bizID, nil
}

func (p *processOperation) validateModuleID(kit *rest.Kit, moduleID int64) (*metadata.ModuleInst,
	errors.CCErrorCoder) {
	// avoid unnecessary db query
	if moduleID == 0 {
		return nil, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKModuleIDField)
	}

	module := &metadata.ModuleInst{}
	filter := map[string]interface{}{
		common.BKModuleIDField: moduleID,
	}
	err := mongodb.Shard(kit.ShardOpts()).Table(common.BKTableNameBaseModule).Find(filter).One(kit.Ctx, module)
	if err != nil {
		blog.Errorf("find module failed, err: %v, filter: %v, rid: %s", err, filter, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommDBSelectFailed)
	}

	return module, nil
}

func (p *processOperation) validateHostID(kit *rest.Kit, hostID int64) (string, errors.CCErrorCoder) {
	// avoid unnecessary db query
	if hostID == 0 {
		return "", kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKHostIDField)
	}

	// check bizID valid
	filter := map[string]interface{}{
		common.BKHostIDField: hostID,
	}
	host := metadata.HostMapStr{}
	err := mongodb.Shard(kit.ShardOpts()).Table(common.BKTableNameBaseHost).Find(filter).Fields(
		common.BKHostInnerIPField).One(kit.Ctx, &host)
	if err != nil {
		blog.Errorf("find module failed, err: %v, filter: %v, rid: %s", err, filter, kit.Rid)
		return "", kit.CCError.CCErrorf(common.CCErrCommDBSelectFailed)
	}

	return util.GetStrByInterface(host[common.BKHostInnerIPField]), nil
}
