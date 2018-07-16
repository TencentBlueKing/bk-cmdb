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

package operation

import (
	"context"
	"net/http"

	"configcenter/src/apimachinery"
	"configcenter/src/common"
	"configcenter/src/common/auditoplog"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/topo_server/core/model"
)

// WrapperResult the data wrapper
type WrapperResult struct {
	Header          http.Header
	User            string
	BizID           string
	InstID          int64
	SupplierAccount string
	Obj             *metadata.Object
	Data            interface{}
}

// AuditInterface audit log methods
type AuditInterface interface {
	CreateSnapshot(header http.Header, instID int64, obj model.Object, data mapstr.MapStr) *WrapperResult
	CommitCreateLog(preData, currData *WrapperResult)
	CommitDeleteLog(preData, currData *WrapperResult)
	CommitUpdateLog(preData, currData *WrapperResult)
}

type auditLog struct {
	client apimachinery.ClientSetInterface
	inst   InstOperationInterface
}

func (a *auditLog) commitSnapshot(preData, currData *WrapperResult, action auditoplog.AuditOpType) {
	if nil != currData {
		blog.Errorf("[audit] the curr data is empty")
		return
	}
	var preDataTmp interface{}
	if nil != preData {
		preDataTmp = preData.Data
	}
	desc := ""
	switch action {
	case auditoplog.AuditOpTypeAdd:
		desc = "create " + currData.Obj.GetObjectType()
	case auditoplog.AuditOpTypeDel:
		desc = "delete " + currData.Obj.GetObjectType()
	case auditoplog.AuditOpTypeModify:
		desc = "update " + currData.Obj.GetObjectType()

	}

	data := common.KvMap{
		common.BKContentField: Content{
			CurData: currData.Data,
			PreData: preDataTmp,
		},
		common.BKOpDescField:   desc,
		common.BKOpTypeField:   auditoplog.AuditOpTypeAdd,
		common.BKOpTargetField: currData.Obj.GetObjectType(),
		"inst_id":              currData.InstID,
	}

	switch currData.Obj.GetObjectType() {

	case common.BKInnerObjIDApp:
		a.client.AuditController().AddObjectLog(context.Background(), currData.SupplierAccount, currData.BizID, currData.User, currData.Header, data)
	case common.BKINnerObjIDObject:
		a.client.AuditController().AddObjectLog(context.Background(), currData.SupplierAccount, currData.BizID, currData.User, currData.Header, data)
	case common.BKInnerObjIDModule:
		a.client.AuditController().AddModuleLog(context.Background(), currData.SupplierAccount, currData.BizID, currData.User, currData.Header, data)
	case common.BKInnerObjIDSet:
		a.client.AuditController().AddSetLog(context.Background(), currData.SupplierAccount, currData.BizID, currData.User, currData.Header, data)
	}
}

func (a *auditLog) CreateSnapshot(header http.Header, instID int64, obj model.Object, data mapstr.MapStr) *WrapperResult {
	return nil
}

func (a *auditLog) CommitCreateLog(preData, currData *WrapperResult) {
	a.commitSnapshot(preData, currData, auditoplog.AuditOpTypeAdd)
}

func (a *auditLog) CommitDeleteLog(preData, currData *WrapperResult) {
	a.commitSnapshot(preData, currData, auditoplog.AuditOpTypeDel)
}

func (a *auditLog) CommitUpdateLog(preData, currData *WrapperResult) {
	a.commitSnapshot(preData, currData, auditoplog.AuditOpTypeModify)
}
