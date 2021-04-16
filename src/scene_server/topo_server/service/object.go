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
	"time"

	"configcenter/src/ac"
	"configcenter/src/ac/iam"
	"configcenter/src/ac/meta"
	"configcenter/src/common"
	"configcenter/src/common/auth"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/topo_server/core/model"
	"configcenter/src/scene_server/topo_server/core/operation"
)

// CreateObjectBatch batch to create some objects
func (s *Service) CreateObjectBatch(ctx *rest.Contexts) {
	data := new(map[string]operation.ImportObjectData)
	if err := ctx.DecodeInto(data); err != nil {
		ctx.RespAutoError(err)
		return
	}

	// auth: check authorization
	objIDs := make([]string, 0)
	for objID := range *data {
		objIDs = append(objIDs, objID)
	}

	if err := s.AuthManager.AuthorizeByObjectIDs(ctx.Kit.Ctx, ctx.Kit.Header, meta.UpdateMany, 0, objIDs...); err != nil {
		blog.Errorf("check object authorization failed, objIDs: %+v, err: %v, rid: %s", objIDs, err, ctx.Kit.Rid)
		if err != ac.NoAuthorizeError {
			ctx.RespAutoError(err)
			return
		}

		perm, err := s.AuthManager.GenObjectBatchNoPermissionResp(ctx.Kit.Ctx, ctx.Kit.Header, meta.UpdateMany, 0, objIDs)
		if err != nil {
			ctx.RespAutoError(err)
			return
		}
		ctx.RespEntityWithError(perm, ac.NoAuthorizeError)
		return
	}

	// 创建模型前，先创建表，避免模型创建后，对模型数据查询出现下面的错误，
	// (SnapshotUnavailable) Unable to read from a snapshot due to pending collection catalog changes;
	// please retry the operation. Snapshot timestamp is Timestamp(1616747877, 51).
	// Collection minimum is Timestamp(1616747878, 5)
	if err := s.createObjectTableBatch(ctx, *data); err != nil {
		ctx.RespAutoError(err)
	}

	var ret mapstr.MapStr
	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		var err error
		ret, err = s.Core.ObjectOperation().CreateObjectBatch(ctx.Kit, *data)
		if err != nil {
			return err
		}
		return nil
	})

	if txnErr != nil {
		ctx.RespAutoError(txnErr)
		return
	}
	ctx.RespEntity(ret)
}

// SearchObjectBatch batch to search some objects
func (s *Service) SearchObjectBatch(ctx *rest.Contexts) {
	data := struct {
		operation.ExportObjectCondition `json:",inline"`
	}{}
	if err := ctx.DecodeInto(&data); nil != err {
		ctx.RespAutoError(err)
		return
	}
	resp, err := s.Core.ObjectOperation().FindObjectBatch(ctx.Kit, data.ObjIDS)
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

	var rsp model.Object
	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		var err error
		rsp, err = s.Core.ObjectOperation().CreateObject(ctx.Kit, false, *data)
		if nil != err {
			return err
		}

		// 注册: 1.创建者权限 2.新资源类型 3.新实例视图 4.新鉴权动作 5.新鉴权动作分组
		if auth.EnableAuthorize() {
			iamInstance := metadata.IamInstanceWithCreator{
				Type:    string(iam.SysModel),
				ID:      strconv.FormatInt(rsp.Object().ID, 10),
				Name:    rsp.Object().ObjectName,
				Creator: ctx.Kit.User,
			}
			// 1.register object resource creator action to iam
			_, err = s.AuthManager.Authorizer.RegisterResourceCreatorAction(ctx.Kit.Ctx, ctx.Kit.Header, iamInstance)
			if err != nil {
				blog.Errorf("register created object to iam failed, err: %s, rid: %s", err, ctx.Kit.Rid)
				return err
			}

			modelList := []metadata.Object{rsp.Object()}
			// 2.register new resourceType when new model is created
			err = s.AuthManager.Authorizer.RegisterModelResourceTypes(ctx.Kit.Ctx, ctx.Kit.Header, modelList)
			if err != nil {
				blog.Errorf("register created object to iam failed, err: %s, rid: %s", err, ctx.Kit.Rid)
				return err
			}

			// 3.register new instance_selection when new model is created
			err = s.AuthManager.Authorizer.RegisterModelInstanceSelections(ctx.Kit.Ctx, ctx.Kit.Header, modelList)
			if err != nil {
				blog.Errorf("register created object to iam failed, err: %s, rid: %s", err, ctx.Kit.Rid)
				return err
			}

			// 4.register IAM actions when new model is created
			err = s.AuthManager.Authorizer.CreateModelInstanceActions(ctx.Kit.Ctx, ctx.Kit.Header, modelList)
			if err != nil {
				blog.Errorf("register created object to iam failed, err: %s, rid: %s", err, ctx.Kit.Rid)
				return err
			}
			time.Sleep(time.Duration(1) * time.Second)

			// 5.register IAM action groups when new model is created
			// 注意: 此接口目前为全量更新, IAM暂时没有提供增量更新接口
			err = s.AuthManager.Authorizer.UpdateModelInstanceActionGroups(ctx.Kit.Ctx, ctx.Kit.Header)
			if err != nil {
				blog.Errorf("register created object to iam failed, err: %s, rid: %s", err, ctx.Kit.Rid)
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
	data := new(mapstr.MapStr)
	if err := ctx.DecodeInto(data); err != nil {
		ctx.RespAutoError(err)
		return
	}
	cond := condition.CreateCondition()
	if err := cond.Parse(*data); nil != err {
		ctx.RespAutoError(err)
		return
	}

	ctx.SetReadPreference(common.SecondaryPreferredMode)
	resp, err := s.Core.ObjectOperation().FindObject(ctx.Kit, cond)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(resp)
}

// SearchObjectTopo search the object topo
func (s *Service) SearchObjectTopo(ctx *rest.Contexts) {
	data := new(mapstr.MapStr)
	if err := ctx.DecodeInto(data); err != nil {
		ctx.RespAutoError(err)
		return
	}
	cond := condition.CreateCondition()
	err := cond.Parse(*data)
	if nil != err {
		ctx.RespAutoError(ctx.Kit.CCError.New(common.CCErrTopoObjectSelectFailed, err.Error()))
		return
	}

	resp, err := s.Core.ObjectOperation().FindObjectTopo(ctx.Kit, cond)
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
		blog.Errorf("[api-obj] failed to parse the path params id(%s), error info is %s , rid: %s", idStr, err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedInt, common.BKFieldID))
		return
	}
	//update model
	data := make(map[string]interface{})
	if err := ctx.DecodeInto(&data); err != nil {
		ctx.RespAutoError(err)
		return
	}

	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		err = s.Core.ObjectOperation().UpdateObject(ctx.Kit, data, id)
		if err != nil {
			return err
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
		blog.Errorf("[api-obj] failed to parse the path params id(%s), error info is %s , rid: %s", idStr, err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKFieldID))
		return
	}

	obj := &metadata.Object{}
	//delete model
	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		obj, err = s.Core.ObjectOperation().DeleteObject(ctx.Kit, id, true)
		if err != nil {
			return err
		}
		// 事务内部可以允许不多于一个的其它系统调用函数, 排在事务尾部, 仍可保证原子性
		// 由于ActionGroups会直观的让用户在IAM内看到访问界面的变化, 应尽量优先保证与cmdb操作的事务性
		// 1.update IAM action groups
		if auth.EnableAuthorize() {
			err = s.AuthManager.Authorizer.UpdateModelInstanceActionGroups(ctx.Kit.Ctx, ctx.Kit.Header)
			if err != nil {
				blog.Errorf("[api-obj] delete model, but update instance action group to iam failed, err: %s, rid: %s", err, ctx.Kit.Rid)
				return err
			}
		}
		return nil
	})
	if txnErr != nil {
		ctx.RespAutoError(txnErr)
		return
	}

	// 如果注销操作有部分成功, 然后失败; 此时cmdb数据如果回滚会导致IAM的数据出现不可预知的情况; 故此处操作不能实现事务
	// 此处的注销操作如果失败, 仍可依赖周期性任务[此任务挂在authServer进程上]保证IAM数据同步
	if auth.EnableAuthorize() {
		objList := []metadata.Object{*obj}
		// 2.delete action
		err = s.AuthManager.Authorizer.DeleteModelInstanceActions(ctx.Kit.Ctx, ctx.Kit.Header, objList)
		if err != nil {
			blog.Errorf("[api-obj] delete model, but unregister actions failed, err: %s, rid: %s", err, ctx.Kit.Rid)

			// 这里的失败, 用户是无感知的, 可依赖周期性任务实现最终一致
			// ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommUnRegistResourceToIAMFailed))
			ctx.RespEntity(nil)
			return
		}

		// 3.instance selection
		err = s.AuthManager.Authorizer.UnregisterModelInstanceSelections(ctx.Kit.Ctx, ctx.Kit.Header, objList)
		if err != nil {
			blog.Errorf("[api-obj] delete model, but unregister instance selections failed, err: %s, rid: %s", err, ctx.Kit.Rid)

			// 这里的失败, 用户是无感知的, 可依赖周期性任务实现最终一致
			// ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommUnRegistResourceToIAMFailed))
			ctx.RespEntity(nil)
			return
		}

		// 4.unregister resourceType
		err = s.AuthManager.Authorizer.UnregisterModelResourceTypes(ctx.Kit.Ctx, ctx.Kit.Header, objList)
		if err != nil {
			blog.Errorf("[api-obj] delete model, but unregister model resource types failed, err: %s, rid: %s", err, ctx.Kit.Rid)

			// 这里的失败, 用户是无感知的, 可依赖周期性任务实现最终一致
			// ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommUnRegistResourceToIAMFailed))
			ctx.RespEntity(nil)
			return
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

// 创建模型前，先创建表，避免模型创建后，对模型数据查询出现下面的错误，
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

// 创建模型前，先创建表，避免模型创建后，对模型数据查询出现下面的错误，
// (SnapshotUnavailable) Unable to read from a snapshot due to pending collection catalog changes;
// please retry the operation. Snapshot timestamp is Timestamp(1616747877, 51).
// Collection minimum is Timestamp(1616747878, 5)
func (s *Service) createObjectTableBatch(ctx *rest.Contexts, objectMap map[string]operation.ImportObjectData) error {

	input := &metadata.CreateModelTable{
		IsMainLine: false,
	}
	for objID := range objectMap {
		if objID != "" {
			input.ObjectIDs = append(input.ObjectIDs, objID)

		}
	}

	// 表创建成功后，需要sleep 要不然有查询操作会报错
	time.Sleep(time.Millisecond * 400)
	return s.Engine.CoreAPI.CoreService().Model().CreateModelTables(ctx.Kit.Ctx, ctx.Kit.Header, input)
}

// 创建模型前，先创建表，避免模型创建后，对模型数据查询出现下面的错误，
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
