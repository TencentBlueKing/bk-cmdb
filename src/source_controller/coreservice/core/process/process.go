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
	"fmt"
	"strconv"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	"configcenter/src/common/metadata"
	"configcenter/src/source_controller/coreservice/core"
	"configcenter/src/storage/dal"
)

type processOperation struct {
	dbProxy dal.RDB
}

// New create a new model manager instance
func New(dbProxy dal.RDB) core.ProcessOperation {
	processOps := &processOperation{dbProxy: dbProxy}
	return processOps
}

func (p *processOperation) validateBizID(ctx core.ContextParams, md metadata.MetadataNG) (int64, error) {
	// extract biz id from metadata
	bizID, err := strconv.ParseInt(md.Label.BusinessID, 10, 64)
	if err != nil {
		blog.Errorf("parse biz id from metadata failed, bizID: %s, err: %+v", md.Label.BusinessID, err)
		return 0, err
	}

	// avoid unnecessary db query
	if bizID == 0 {
		return 0, fmt.Errorf("bizID invalid, bizID: %d", bizID)
	}

	// check bizID valid
	cond := condition.CreateCondition()
	cond.Field(common.BKAppIDField).Eq(bizID)
	count, err := p.dbProxy.Table(common.BKTableNameBaseApp).Find(cond.ToMapStr()).Count(ctx.Context)
	if nil != err {
		blog.Errorf("mongodb failed, table: %s, err: %+v, rid: %s", common.BKTableNameObjDes, err.Error(), ctx.ReqID)
		return 0, err
	}
	if count < 1 {
		return 0, fmt.Errorf("business not found, id:%d", bizID)
	}

	return bizID, nil
}

func (p *processOperation) validateModuleID(ctx core.ContextParams, moduleID int64) error {
	// avoid unnecessary db query
	if moduleID == 0 {
		return fmt.Errorf("moduleID invalid, moduleID: %d", moduleID)
	}

	// check bizID valid
	cond := condition.CreateCondition()
	cond.Field(common.BKModuleIDField).Eq(moduleID)
	count, err := p.dbProxy.Table(common.BKTableNameBaseModule).Find(cond.ToMapStr()).Count(ctx.Context)
	if nil != err {
		blog.Errorf("validateModuleID failed, mongodb failed, table: %s, err: %+v, rid: %s", common.BKTableNameBaseModule, err.Error(), ctx.ReqID)
		return err
	}
	if count < 1 {
		return fmt.Errorf("validateModuleID failed, module not found, id:%d", moduleID)
	}

	return nil
}

func (p *processOperation) validateHostID(ctx core.ContextParams, hostID int64) error {
	// avoid unnecessary db query
	if hostID == 0 {
		return fmt.Errorf("hostID invalid, bizID: %d", hostID)
	}

	// check bizID valid
	cond := condition.CreateCondition()
	cond.Field(common.BKHostIDField).Eq(hostID)
	count, err := p.dbProxy.Table(common.BKTableNameBaseHost).Find(cond.ToMapStr()).Count(ctx.Context)
	if nil != err {
		blog.Errorf("validateHostID failed, mongodb failed, table: %s, err: %+v, rid: %s", common.BKTableNameBaseHost, err.Error(), ctx.ReqID)
		return err
	}
	if count < 1 {
		return fmt.Errorf("validateHostID failed, host not found, id:%d", hostID)
	}

	return nil
}
