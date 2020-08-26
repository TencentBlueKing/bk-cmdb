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

package auditlog

import (
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/source_controller/coreservice/core"
	"configcenter/src/storage/dal"
)

var _ core.AuditOperation = (*auditManager)(nil)

type auditManager struct {
	dbProxy dal.RDB
}

// New create a new instance manager instance
func New(dbProxy dal.RDB) core.AuditOperation {
	return &auditManager{
		dbProxy: dbProxy,
	}
}

func (m *auditManager) CreateAuditLog(kit *rest.Kit, logs ...metadata.AuditLog) error {
	logRows := make([]metadata.AuditLog, 0)

	for _, log := range logs {
		if log.OperationDetail == nil {
			continue
		}

		if log.OperateFrom == "" {
			log.OperateFrom = metadata.FromUser
		}
		log.SupplierAccount = kit.SupplierAccount
		log.User = kit.User
		log.OperationTime = metadata.Now()

		id, err := m.dbProxy.NextSequence(kit.Ctx, common.BKTableNameAuditLog)
		if err != nil {
			blog.Errorf("get next audit log id failed, err: %s", err.Error())
			return err
		}
		log.ID = int64(id)

		logRows = append(logRows, log)
	}

	if len(logRows) == 0 {
		return nil
	}
	return m.dbProxy.Table(common.BKTableNameAuditLog).Insert(kit.Ctx, logRows)
}

func (m *auditManager) SearchAuditLog(kit *rest.Kit, param metadata.QueryCondition) ([]metadata.AuditLog, uint64, error) {
	condition := param.Condition
	condition = util.SetQueryOwner(condition, kit.SupplierAccount)
	blog.V(5).Infof("Search table common.BKTableNameAuditLog with parameters: %+v, rid: %s", condition, kit.Rid)

	rows := make([]metadata.AuditLog, 0)
	err := m.dbProxy.Table(common.BKTableNameAuditLog).Find(condition).Sort(param.Page.Sort).Fields(param.Fields...).
		Start(uint64(param.Page.Start)).Limit(uint64(param.Page.Limit)).All(kit.Ctx, &rows)
	if nil != err {
		blog.Errorf("query database error:%s, condition:%v, rid: %s", err.Error(), condition, kit.Rid)
		return nil, 0, err
	}

	cnt, err := m.dbProxy.Table(common.BKTableNameAuditLog).Find(condition).Count(kit.Ctx)
	if nil != err {
		blog.Errorf("query database error:%s, condition:%v, rid: %s", err.Error(), condition, kit.Rid)
		return nil, 0, err
	}

	return rows, cnt, nil
}
