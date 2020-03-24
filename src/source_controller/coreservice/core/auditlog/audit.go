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
	"github.com/google/go-cmp/cmp/cmpopts"
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
// notice: getIgnoreOptions用来设置不参与对比变化的字段，这些字段发生变化，在instNotChange不在返回数据发生变化
func instNotChange(ctx context.Context, content metadata.DetailFactory) bool {
	rid := util.ExtractRequestIDFromContext(ctx)
	modelID := ""
	var basicContent *metadata.BasicOpDetail
	switch content.WithName() {
	case "BasicDetail":
		basicContent = content.(*metadata.BasicOpDetail)
	case "InstanceOpDetail":
		instanceContent := content.(*metadata.InstanceOpDetail)
		modelID = instanceContent.ModelID
		basicContent = &instanceContent.BasicOpDetail
	case "HostTransferOpDetail":
		hostTransferContent := content.(*metadata.HostTransferOpDetail)
		// ignore default field
		bl := cmp.Equal(hostTransferContent.PreData, hostTransferContent.CurData, getIgnoreOptions(""))
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
	bl := cmp.Equal(preData, curData, getIgnoreOptions(modelID))
	if bl {
		blog.V(5).Infof("inst data same, %+v, rid: %s", content, rid)
	}
	return bl
}

// getIgnoreOptions ignore fields options,不参与对比变化的字段，这些字段发生变化，在instNotChange不在返回数据发生变化
// params objID 模型id，预留字段，为根据不同模型实现不同忽略字段,
func getIgnoreOptions(objID string) cmp.Option {
	field := make(map[string]interface{}, 0)
	switch objID {
	default:
		field = ignoreCmpFields["default"]
	}
	if len(field) == 0 {
		return nil
	}
	ignoreCmpFunc := func(key string, val interface{}) bool {
		if _, ok := field[key]; ok {
			return ok
		}
		return false
	}

	option := cmpopts.IgnoreMapEntries(ignoreCmpFunc)
	return option
}

var (
	ignoreCmpFields = map[string]map[string]interface{}{
		// default 默认情况下忽略的字段
		"default": map[string]interface{}{
			"_id":                nil,
			common.LastTimeField: nil,
		},
	}
)
