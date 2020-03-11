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
	"context"
	"strings"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/source_controller/coreservice/core"
	"configcenter/src/storage/dal"

	"github.com/google/go-cmp/cmp"
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
	var logRows []interface{}
	for _, log := range logs {
		if log.OperationDetail == nil || instNotChange(kit.Ctx, log.OperationDetail) {
			continue
		}
		if log.OperateFrom == "" {
			log.OperateFrom = metadata.FromUser
		}
		log.SupplierAccount = kit.SupplierAccount
		log.User = kit.User
		log.OperationTime = metadata.Now()
		logRows = append(logRows, log)
	}
	if len(logRows) == 0 {
		return nil
	}
	return m.dbProxy.Table(common.BKTableNameAuditLog).Insert(kit.Ctx, logRows)
}

func (m *auditManager) SearchAuditLog(kit *rest.Kit, param metadata.QueryInput) ([]metadata.AuditLog, uint64, error) {
	fields := param.Fields
	condition := param.Condition
	condition = util.SetQueryOwner(condition, kit.SupplierAccount)
	param.ConvTime()
	skip := param.Start
	limit := param.Limit
	fieldArr := strings.Split(fields, ",")
	rows := make([]metadata.AuditLog, 0)
	blog.V(5).Infof("Search table common.BKTableNameAuditLog with parameters: %+v, rid: %s", condition, kit.Rid)
	err := m.dbProxy.Table(common.BKTableNameAuditLog).Find(condition).Sort(param.Sort).Fields(fieldArr...).Start(uint64(skip)).Limit(uint64(limit)).All(kit.Ctx, &rows)
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

// instNotChange Determine whether the data is consistent before and after the change
func instNotChange(ctx context.Context, content metadata.DetailFactory) bool {
	rid := util.ExtractRequestIDFromContext(ctx)
	var basicContent *metadata.BasicOpDetail
	switch content.WithName() {
	case "BasicDetail":
		basicContent = content.(*metadata.BasicOpDetail)
	case "InstanceOpDetail":
		instanceContent := content.(*metadata.InstanceOpDetail)
		basicContent = &instanceContent.BasicOpDetail
	case "HostTransferOpDetail":
		hostTransferContent := content.(*metadata.HostTransferOpDetail)
		bl := cmp.Equal(hostTransferContent.PreData, hostTransferContent.CurData)
		if bl {
			blog.V(5).Infof("inst data same, %+v, rid: %s", content, rid)
		}
		return bl
	case "CloudAccountOpDetail":
		cloudAccountContent := content.(*metadata.CloudAccountOpDetail)
		bl := cmp.Equal(cloudAccountContent.PreData, cloudAccountContent.CurData)
		if bl {
			blog.V(5).Infof("inst data same, %+v, rid: %s", content, rid)
		}
		return bl
	}
	if basicContent == nil || basicContent.Details == nil || basicContent.Details.PreData == nil || basicContent.Details.CurData == nil {
		return false
	}
	preData := basicContent.Details.PreData
	curData := basicContent.Details.CurData
	delete(preData, common.LastTimeField)
	delete(curData, common.LastTimeField)
	bl := cmp.Equal(preData, curData)
	if bl {
		blog.V(5).Infof("inst data same, %+v, rid: %s", content, rid)
	}
	return bl
}
