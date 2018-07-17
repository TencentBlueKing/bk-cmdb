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
	"strings"
	"time"

	"configcenter/src/common/auditoplog"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
)

// AddLogMulti insert multiple row
func (lgc *Logics) AddLogMulti(appID int64, opType auditoplog.AuditOpType, opTarget string, contents []auditoplog.AuditLogContext, opDesc, ownerID, user string) error {
	var logRows []interface{}

	for _, content := range contents {
		row := &metadata.OperationLog{
			OwnerID:       ownerID,
			ApplicationID: appID,
			OpType:        int(opType),
			OpTarget:      opTarget,
			User:          user,
			ExtKey:        "",
			OpDesc:        opDesc,
			Content:       content.Content,
			CreateTime:    time.Now(),
			InstID:        content.ID,
		}
		logRows = append(logRows, row)

	}
	if len(logRows) == 0 {
		return nil
	}
	logRow := metadata.OperationLog{}
	err := lgc.Instance.InsertMuti(logRow.TableName(), logRows...)
	return err
}

// AddLogMultiWithExtKey insert multiple row with  extension key
func (lgc *Logics) AddLogMultiWithExtKey(appID int64, opType auditoplog.AuditOpType, opTarget string, contents []auditoplog.AuditLogExt, opDesc, ownerID, user string) error {
	var logRows []interface{}

	for _, content := range contents {
		row := &metadata.OperationLog{
			OwnerID:       ownerID,
			ApplicationID: appID,
			OpType:        int(opType),
			OpTarget:      opTarget,
			User:          user,
			ExtKey:        content.ExtKey,
			OpDesc:        opDesc,
			Content:       content.Content,
			CreateTime:    time.Now(),
			InstID:        content.ID,
		}
		logRows = append(logRows, row)

	}
	if len(logRows) == 0 {
		return nil
	}
	logRow := metadata.OperationLog{}
	err := lgc.Instance.InsertMuti(logRow.TableName(), logRows...)
	return err
}

// AddLogWithStr insert row
func (lgc *Logics) AddLogWithStr(appID, instID int64, opType auditoplog.AuditOpType, opTarget string, content interface{}, extKey, opDesc, ownerID, user string) error {
	logRow := &metadata.OperationLog{
		OwnerID:       ownerID,
		ApplicationID: appID,
		OpType:        int(opType),
		OpTarget:      opTarget,
		User:          user,
		ExtKey:        extKey,
		OpDesc:        opDesc,
		Content:       content,
		CreateTime:    time.Now(),
		InstID:        instID,
	}
	_, err := lgc.Instance.Insert(logRow.TableName(), logRow)
	return err
}

// Search query operation log
func (lgc *Logics) Search(dat *metadata.ObjQueryInput) ([]metadata.OperationLog, int, error) {
	fields := dat.Fields
	condition := dat.Condition
	dat.ConvTime()
	skip := dat.Start
	limit := dat.Limit
	sort := dat.Sort
	fieldArr := strings.Split(fields, ",")
	rows := make([]metadata.OperationLog, 0)
	logRow := metadata.OperationLog{}
	err := lgc.Instance.GetMutilByCondition(logRow.TableName(), fieldArr, condition, &rows, sort, skip, limit)
	if nil != err {
		blog.Errorf("query database error:%s, condition:%v", err.Error, condition)
		return nil, 0, err
	}
	cnt, err := lgc.Instance.GetCntByCondition(logRow.TableName(), condition)
	if nil != err {
		blog.Errorf("query database error:%s, condition:%v", err.Error, condition)
		return nil, 0, err
	}

	return rows, cnt, nil

}
