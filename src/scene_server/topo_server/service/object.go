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

package service

import (
	"fmt"
	"strconv"

	"configcenter/src/ac"
	"configcenter/src/ac/iam"
	"configcenter/src/ac/meta"
	"configcenter/src/common"
	"configcenter/src/common/auth"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
)

// CreateObjectBatch batch to create some objects
func (s *Service) CreateObjectBatch(ctx *rest.Contexts) {
	data := new(map[string]metadata.ImportObjectData)
	if err := ctx.DecodeInto(data); err != nil {
		ctx.RespAutoError(err)
		return
	}

	// auth: check authorization
	objIDs := make([]string, 0)
	for objID := range *data {
		objIDs = append(objIDs, objID)
	}

	err := s.AuthManager.AuthorizeByObjectIDs(ctx.Kit.Ctx, ctx.Kit.Header, meta.UpdateMany, 0, objIDs...)
	if err != nil {
		blog.Errorf("check object authorization failed, objIDs: %+v, err: %v, rid: %s", objIDs, err, ctx.Kit.Rid)
		if err != ac.NoAuthorizeError {
			ctx.RespAutoError(err)
			return
		}

		perm, err := s.AuthManager.GenObjectBatchNoPermissionResp(ctx.Kit.Ctx, ctx.Kit.Header, meta.UpdateMany, 0,
			objIDs)
		if err != nil {
			ctx.RespAutoError(err)
			return
		}
		ctx.RespEntityWithError(perm, ac.NoAuthorizeError)
		return
	}

	var ret mapstr.MapStr
	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		var err error
		ret, err = s.Logics.AttributeOperation().CreateObjectBatch(ctx.Kit, *data)
		if err != nil {
			return err
		}
		return nil
	})

	if txnErr != nil {
		ctx.RespEntityWithError(ret, txnErr)
		return
	}
	ctx.RespEntity(ret)
}

// SearchObjectBatch batch to search some objects
func (s *Service) SearchObjectBatch(ctx *rest.Contexts) {
	data := metadata.ExportObjectCondition{}
	if err := ctx.DecodeInto(&data); nil != err {
		ctx.RespAutoError(err)
		return
	}
	resp, err := s.Logics.AttributeOperation().FindObjectBatch(ctx.Kit, data.ObjIDs)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(resp)
}

// CreateObject create a new object
func (s *Service) CreateObject(ctx *rest.Contexts) {
	data := new(mapstr.MapStr)
	if err := ctx.DecodeInto(&data); err != nil {
		ctx.RespAutoError(err)
		return
	}

	// 创建模型前，先创建表，避免模型创建后，对模型数据查询出现下面的错误，
	// (SnapshotUnavailable) Unable to read from a snapshot due to pending collection catalog changes;
	// please retry the operation. Snapshot timestamp is Timestamp(1616747877, 51).
	// Collection minimum is Timestamp(1616747878, 5)
	if err := s.createObjectTable(ctx, *data); err != nil {
		ctx.RespAutoError(err)
		return
	}

	var rsp *metadata.Object
	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		var err error
		rsp, err = s.Logics.ObjectOperation().CreateObject(ctx.Kit, false, *data)
		if nil != err {
			return err
		}

		if auth.EnableAuthorize() {
			objects := []metadata.Object{*rsp}
			iamInstances := []metadata.IamInstanceWithCreator{{
				Type:    string(iam.SysModel),
				ID:      strconv.FormatInt(rsp.ID, 10),
				Name:    rsp.ObjectName,
				Creator: ctx.Kit.User,
			}}
			if err := s.AuthManager.CreateObjectOnIAM(ctx.Kit.Ctx, ctx.Kit.Header, objects, iamInstances); err != nil {
				blog.ErrorJSON("create object on iam failed, objects: %s, iam instances: %s, err: %s, rid: %s",
					objects, iamInstances, err, ctx.Kit.Rid)
				return err
			}
		}
		return nil
	})

	if txnErr != nil {
		ctx.RespAutoError(txnErr)
		return
	}
	ctx.RespEntity(rsp.ToMapStr())
}

// SearchObject search some objects by condition
func (s *Service) SearchObject(ctx *rest.Contexts) {
	data := mapstr.New()
	if err := ctx.DecodeInto(&data); err != nil {
		ctx.RespAutoError(err)
		return
	}

	query := &metadata.QueryCondition{Condition: data, DisableCounter: true}
	resp, err := s.Engine.CoreAPI.CoreService().Model().ReadModel(ctx.Kit.Ctx, ctx.Kit.Header, query)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(resp.Info)
}

// SearchObjectTopo search the object topo
func (s *Service) SearchObjectTopo(ctx *rest.Contexts) {
	data := mapstr.New()
	if err := ctx.DecodeInto(&data); err != nil {
		ctx.RespAutoError(err)
		return
	}

	resp, err := s.Logics.ObjectOperation().FindObjectTopo(ctx.Kit, data)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(resp)
}

// UpdateObject update the object
func (s *Service) UpdateObject(ctx *rest.Contexts) {
	idStr := ctx.Request.PathParameter(common.BKFieldID)
	id, err := strconv.ParseInt(idStr, 10, 64)
	if nil != err {
		blog.Errorf("failed to parse the path params id(%s), err: %v , rid: %s", idStr, err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedInt, common.BKFieldID))
		return
	}
	// update model
	data := make(map[string]interface{})
	if err := ctx.DecodeInto(&data); err != nil {
		ctx.RespAutoError(err)
		return
	}

	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		err = s.Logics.ObjectOperation().UpdateObject(ctx.Kit, data, id)
		if err != nil {
			return err
		}

		// sync the object name to IAM only if the object name is updated
		if auth.EnableAuthorize() {
			// judge whether the data contains object name
			if _, ok := data[common.BKObjNameField]; !ok {
				return nil
			}

			cond := &metadata.QueryCondition{
				Page:           metadata.BasePage{Limit: common.BKNoLimit},
				Condition:      mapstr.MapStr{common.BKFieldID: id},
				DisableCounter: true,
			}

			resp, err := s.Engine.CoreAPI.CoreService().Model().ReadModel(ctx.Kit.Ctx, ctx.Kit.Header, cond)
			if err != nil {
				blog.ErrorJSON("find object failed, cond: %s, err: %s, rid: %s", cond, err, ctx.Kit.Rid)
				return err
			}
			if len(resp.Info) != 1 {
				blog.ErrorJSON("object count is wrong, id: %d, resp: %#v, rid: %s", id, resp, ctx.Kit.Rid)
				return ctx.Kit.CCError.New(common.CCErrCommDBSelectFailed, "object count is wrong")
			}

			objects := []metadata.Object{resp.Info[0]}
			if err := s.AuthManager.Viewer.UpdateView(ctx.Kit.Ctx, ctx.Kit.Header, objects); err != nil {
				blog.Errorf("update view failed, err: %s, rid: %s", err, ctx.Kit.Rid)
				return err
			}
		}

		return nil
	})

	if txnErr != nil {
		ctx.RespAutoError(txnErr)
		return
	}
	ctx.RespEntity(nil)
}

// DeleteObject delete the object
func (s *Service) DeleteObject(ctx *rest.Contexts) {
	idStr := ctx.Request.PathParameter(common.BKFieldID)
	id, err := strconv.ParseInt(idStr, 10, 64)
	if nil != err {
		blog.Errorf("failed to parse the path params id(%s), err: %v , rid: %s", idStr, err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKFieldID))
		return
	}

	var obj *metadata.Object
	cond := mapstr.MapStr{common.BKFieldID: id}
	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		obj, err = s.Logics.ObjectOperation().DeleteObject(ctx.Kit, cond, true)
		if err != nil {
			return err
		}
		return nil
	})
	if txnErr != nil {
		ctx.RespAutoError(txnErr)
		return
	}

	if auth.EnableAuthorize() {
		// delete iam view
		// if some errors occur, the sync iam task will delete the iam view in the end
		objects := []metadata.Object{*obj}
		// use new transaction, need a new header
		ctx.Kit.Header = ctx.Kit.NewHeader()
		if err := s.AuthManager.Viewer.DeleteView(ctx.Kit.Ctx, ctx.Kit.Header, objects); err != nil {
			blog.Errorf("delete view failed, err: %s, rid: %s", err, ctx.Kit.Rid)
		}
	}

	ctx.RespEntity(nil)
}

// GetModelStatistics 用于统计各个模型的实例数(Web页面展示需要)
func (s *Service) GetModelStatistics(ctx *rest.Contexts) {
	result, err := s.Engine.CoreAPI.CoreService().Model().GetModelStatistics(ctx.Kit.Ctx, ctx.Kit.Header)
	if err != nil {
		blog.Errorf("GetModelStatistics failed, err: %s, rid: %s", err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(result.Data)
}

// SearchModel search some model by condition
func (s *Service) SearchModel(ctx *rest.Contexts) {
	data := new(metadata.QueryCondition)
	if err := ctx.DecodeInto(data); err != nil {
		ctx.RespAutoError(err)
		return
	}

	resp, err := s.Engine.CoreAPI.CoreService().Model().ReadModel(ctx.Kit.Ctx, ctx.Kit.Header, data)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(resp)
}

// createObjectTable 创建模型前，先创建表，避免模型创建后，对模型数据查询出现下面的错误，
// (SnapshotUnavailable) Unable to read from a snapshot due to pending collection catalog changes;
// please retry the operation. Snapshot timestamp is Timestamp(1616747877, 51).
// Collection minimum is Timestamp(1616747878, 5)
func (s *Service) createObjectTable(ctx *rest.Contexts, object map[string]interface{}) error {

	input := &metadata.CreateModelTable{
		IsMainLine: false,
	}
	if objID := object[common.BKObjIDField]; objID != nil {
		strObjID := fmt.Sprintf("%v", objID)
		input.ObjectIDs = []string{strObjID}
		return s.Engine.CoreAPI.CoreService().Model().CreateModelTables(ctx.Kit.Ctx, ctx.Kit.Header, input)

	}
	return nil
}

// createObjectTableByObjectID 创建模型前，先创建表，避免模型创建后，对模型数据查询出现下面的错误，
// (SnapshotUnavailable) Unable to read from a snapshot due to pending collection catalog changes;
// please retry the operation. Snapshot timestamp is Timestamp(1616747877, 51).
// Collection minimum is Timestamp(1616747878, 5)
func (s *Service) createObjectTableByObjectID(ctx *rest.Contexts, objectID string, isMainline bool) error {
	input := &metadata.CreateModelTable{
		IsMainLine: isMainline,
	}

	if objectID != "" {
		input.ObjectIDs = []string{objectID}
		return s.Engine.CoreAPI.CoreService().Model().CreateModelTables(ctx.Kit.Ctx, ctx.Kit.Header, input)

	}
	return nil
}

// CreateManyObject batch create object with it's attr and asst
func (s *Service) CreateManyObject(ctx *rest.Contexts) {
	data := new(metadata.ImportObjects)
	if err := ctx.DecodeInto(data); err != nil {
		ctx.RespAutoError(err)
		return
	}

	for _, item := range data.Objects {
		// 创建模型前，先创建表，避免模型创建后，对模型数据查询出现下面的错误，
		// (SnapshotUnavailable) Unable to read from a snapshot due to pending collection catalog changes;
		// please retry the operation. Snapshot timestamp is Timestamp(1616747877, 51).
		// Collection minimum is Timestamp(1616747878, 5)
		if err := s.createObjectTable(ctx, mapstr.MapStr{common.BKObjIDField: item.ObjectID}); err != nil {
			ctx.RespAutoError(err)
			return
		}
	}

	var rsp []metadata.Object
	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		var err error
		rsp, err = s.Logics.ObjectOperation().CreateObjectByImport(ctx.Kit, data.Objects)
		if err != nil {
			return err
		}

		if auth.EnableAuthorize() {
			iamInstances := make([]metadata.IamInstanceWithCreator, 0)
			for _, obj := range rsp {
				iamInstances = append(iamInstances, metadata.IamInstanceWithCreator{
					Type:    string(iam.SysModel),
					ID:      strconv.FormatInt(obj.ID, 10),
					Name:    obj.ObjectName,
					Creator: ctx.Kit.User,
				})
			}
			err := s.AuthManager.CreateObjectOnIAM(ctx.Kit.Ctx, ctx.Kit.Header, rsp, iamInstances)
			if err != nil {
				blog.Errorf("create object on iam failed, objects: %v, iam instances: %v, err: %v, rid: %s",
					rsp, iamInstances, err, ctx.Kit.Rid)
				return err
			}
		}

		if len(data.Asst) != 0 {
			if err = s.Logics.AssociationOperation().CreateOrUpdateAssociationType(ctx.Kit, data.Asst); err != nil {
				blog.Errorf("create or update association kind failed, err: %v, rid: %s", err, ctx.Kit.Rid)
				return err
			}
		}

		return nil
	})

	if txnErr != nil {
		ctx.RespAutoError(txnErr)
		return
	}

	ctx.RespEntity(rsp)
}

// SearchObjectWithTotalInfo search object with it's attribute and association
func (s *Service) SearchObjectWithTotalInfo(ctx *rest.Contexts) {
	data := new(metadata.BatchExportObject)
	if err := ctx.DecodeInto(data); nil != err {
		ctx.RespAutoError(err)
		return
	}
	resp, err := s.Logics.ObjectOperation().SearchObjectsWithTotalInfo(ctx.Kit, data.ObjectID, data.ExcludedAsstID)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(resp)
}
