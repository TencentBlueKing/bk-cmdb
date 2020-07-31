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

	"configcenter/src/apimachinery"
	"configcenter/src/auth/extensions"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
)

// Unique OperationInterface group operation methods
type UniqueOperationInterface interface {
	Create(kit *rest.Kit, objectID string, request *metadata.CreateUniqueRequest, metaData *metadata.Metadata) (uniqueID *metadata.RspID, err error)
	Update(kit *rest.Kit, objectID string, id uint64, request *metadata.UpdateUniqueRequest) (err error)
	Delete(kit *rest.Kit, objectID string, id uint64, metaData *metadata.Metadata) (err error)
	Search(kit *rest.Kit, objectID string, metaData *metadata.Metadata) (objectUniques []metadata.ObjectUnique, err error)
}

// NewUnique Operation create a new group operation instance
func NewUniqueOperation(client apimachinery.ClientSetInterface, authManager *extensions.AuthManager) UniqueOperationInterface {
	return &unique{
		clientSet:   client,
		authManager: authManager,
	}
}

type unique struct {
	clientSet   apimachinery.ClientSetInterface
	authManager *extensions.AuthManager
}

func (a *unique) Create(kit *rest.Kit, objectID string, request *metadata.CreateUniqueRequest, metaData *metadata.Metadata) (uniqueID *metadata.RspID, err error) {

	unique := metadata.ObjectUnique{
		ObjID:     request.ObjID,
		Keys:      request.Keys,
		MustCheck: request.MustCheck,
	}

	if nil != metaData {
		unique.Metadata = *metaData
	}
	resp, err := a.clientSet.CoreService().Model().CreateModelAttrUnique(context.Background(), kit.Header, objectID, metadata.CreateModelAttrUnique{Data: unique})
	if err != nil {
		blog.Errorf("[operation-unique] create for %s, %#v failed %v, rid: %s", objectID, request, err, kit.Rid)
		return nil, kit.CCError.Error(common.CCErrTopoObjectUniqueCreateFailed)
	}
	if !resp.Result {
		return nil, kit.CCError.New(resp.Code, resp.ErrMsg)
	}

	//package audit response
	err = NewObjectUniqueAudit(kit, a.clientSet, int64(resp.Data.Created.ID)).buildSnapshotForCur().SaveAuditLog(metadata.AuditCreate)
	if err != nil {
		blog.Errorf("[operation-unique] create %s unique item success, but update audit log failed: %v, rid: %s", objectID, err, kit.Rid)
	}

	return &metadata.RspID{ID: int64(resp.Data.Created.ID)}, nil
}

func (a *unique) Update(kit *rest.Kit, objectID string, id uint64, request *metadata.UpdateUniqueRequest) (err error) {
	update := metadata.UpdateModelAttrUnique{
		Data: *request,
	}
	//get PreData
	objAudit := NewObjectUniqueAudit(kit, a.clientSet, int64(id)).buildSnapshotForPre()

	resp, err := a.clientSet.CoreService().Model().UpdateModelAttrUnique(context.Background(), kit.Header, objectID, id, update)
	if err != nil {
		blog.Errorf("[operation-unique] update for %s, %d, %#v failed %v, rid: %s", objectID, id, request, err, kit.Rid)
		return kit.CCError.Error(common.CCErrTopoObjectUniqueUpdateFailed)
	}
	if !resp.Result {
		return kit.CCError.New(resp.Code, resp.ErrMsg)
	}

	// auth: update register to iam
	if err := a.authManager.UpdateRegisteredModelUniqueByID(kit.Ctx, kit.Header, int64(id)); err != nil {
		blog.V(2).Infof("update unique %d for model %s failed, authorization failed, err: %+v, rid: %s", id, objectID, err, kit.Rid)
		return err
	}

	//get CurData and saveAuditLog
	err = objAudit.buildSnapshotForCur().SaveAuditLog(metadata.AuditUpdate)
	if err != nil {
		blog.Errorf("[operation-unique] update %s unique item success, but update audit log failed: %v, rid: %s", objectID, err, kit.Rid)
	}
	return nil
}

func (a *unique) Delete(kit *rest.Kit, objectID string, id uint64, metaData *metadata.Metadata) (err error) {
	meta := metadata.Metadata{}
	if metaData != nil {
		meta = *metaData
	}
	//get PreData
	objAudit := NewObjectUniqueAudit(kit, a.clientSet, int64(id)).buildSnapshotForPre()

	resp, err := a.clientSet.CoreService().Model().DeleteModelAttrUnique(context.Background(), kit.Header, objectID, id, metadata.DeleteModelAttrUnique{Metadata: meta})
	if err != nil {
		blog.Errorf("[operation-unique] delete for %s, %d failed %v, rid: %s", objectID, id, err, kit.Rid)
		return kit.CCError.Error(common.CCErrTopoObjectUniqueDeleteFailed)
	}
	if !resp.Result {
		return kit.CCError.New(resp.Code, resp.ErrMsg)
	}
	// auth: deregister to iam
	if err := a.authManager.DeregisterModelUniqueByID(kit.Ctx, kit.Header, int64(id)); err != nil {
		blog.V(2).Infof("deregister unique %d for model %s failed, authorization failed, err: %+v, rid: %s", id, objectID, err, kit.Rid)
		return err
	}
	//saveAuditLog
	err = objAudit.SaveAuditLog(metadata.AuditDelete)
	if err != nil {
		blog.Errorf("[operation-unique] delete %s unique item success, but update audit log failed: %v, rid: %s", objectID, err, kit.Rid)
	}
	return nil
}

func (a *unique) Search(kit *rest.Kit, objectID string, metaData *metadata.Metadata) (objectUniques []metadata.ObjectUnique, err error) {
	fCond := condition.CreateCondition().Field(common.BKObjIDField).Eq(objectID).ToMapStr()
	if nil != metaData {
		fCond.Merge(metadata.PublicAndBizCondition(*metaData))
		fCond.Remove(metadata.BKMetadata)
	} else {
		fCond.Merge(metadata.BizLabelNotExist)
	}

	cond := metadata.QueryCondition{
		Condition: fCond,
	}
	resp, err := a.clientSet.CoreService().Model().ReadModelAttrUnique(context.Background(), kit.Header, cond)
	if err != nil {
		blog.Errorf("[operation-unique] search for %s, failed %v, rid: %s", objectID, err, kit.Rid)
		return nil, kit.CCError.Error(common.CCErrTopoObjectUniqueSearchFailed)
	}
	if !resp.Result {
		return nil, kit.CCError.New(resp.Code, resp.ErrMsg)
	}
	return resp.Data.Info, nil
}
