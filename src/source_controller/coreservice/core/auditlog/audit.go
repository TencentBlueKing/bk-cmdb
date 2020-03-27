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
	"configcenter/src/common/util"
	"context"
	"strings"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
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

func (m *auditManager) CreateAuditLog(ctx core.ContextParams, logs ...metadata.SaveAuditLogParams) error {

	var logRows []interface{}
	for _, content := range logs {
		if instNotChange(ctx, content.Content, content.Model) {
			continue
		}
		row := &metadata.OperationLog{
			OwnerID:       ctx.SupplierAccount,
			ApplicationID: content.BizID,
			OpType:        int(content.OpType),
			OpTarget:      content.Model,
			User:          ctx.User,
			ExtKey:        content.ExtKey,
			OpDesc:        content.OpDesc,
			Content:       content.Content,
			CreateTime:    time.Now(),
			InstID:        content.ID,
		}
		logRows = append(logRows, row)

	}
	if len(logRows) == 0 {
		return nil
	}
	return m.dbProxy.Table(common.BKTableNameOperationLog).Insert(ctx, logRows)
}

func (m *auditManager) SearchAuditLog(ctx core.ContextParams, param metadata.QueryInput) ([]metadata.OperationLog, uint64, error) {
	fields := param.Fields
	condition := param.Condition
	condition = util.SetQueryOwner(condition, ctx.SupplierAccount)
	param.ConvTime()
	skip := param.Start
	limit := param.Limit
	fieldArr := strings.Split(fields, ",")
	rows := make([]metadata.OperationLog, 0)
	blog.V(5).Infof("Search table common.BKTableNameOperationLog with parameters: %+v, rid: %s", condition, ctx.ReqID)
	err := m.dbProxy.Table(common.BKTableNameOperationLog).Find(condition).Sort(param.Sort).Fields(fieldArr...).Start(uint64(skip)).Limit(uint64(limit)).All(ctx, &rows)
	if nil != err {
		blog.Errorf("query database error:%s, condition:%v, rid: %s", err.Error(), condition, ctx.ReqID)
		return nil, 0, err
	}
	cnt, err := m.dbProxy.Table(common.BKTableNameOperationLog).Find(condition).Count(ctx)
	if nil != err {
		blog.Errorf("query database error:%s, condition:%v, rid: %s", err.Error(), condition, ctx.ReqID)
		return nil, 0, err
	}

	return rows, cnt, nil
}

// instNotChange Determine whether the data is consistent before and after the change
// notice: getIgnoreOptions用来设置不参与对比变化的字段，这些字段发生变化，在instNotChange不在返回数据发生变化
func instNotChange(ctx context.Context, content interface{}, objID string) bool {
	rid := util.ExtractRequestIDFromContext(ctx)
	contentMap, ok := content.(map[string]interface{})
	if !ok {
		return false
	}
	preData, ok := contentMap["pre_data"].(map[string]interface{})
	if !ok {
		return false
	}
	curData, ok := contentMap["cur_data"].(map[string]interface{})
	if !ok {
		return false
	}

	bl := cmp.Equal(preData, curData, getIgnoreOptions(objID))
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
