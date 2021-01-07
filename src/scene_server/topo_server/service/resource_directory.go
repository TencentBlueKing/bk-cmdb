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

	"configcenter/src/ac/iam"
	"configcenter/src/common"
	"configcenter/src/common/auditlog"
	"configcenter/src/common/auth"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/paraparse"
	"configcenter/src/common/util"
)

func (s *Service) CreateResourceDirectory(ctx *rest.Contexts) {
	data := mapstr.MapStr{}
	if err := ctx.DecodeInto(&data); nil != err {
		blog.Errorf("CreateResourceDirectory failed to parse the params, error info is %s, rid: %s", err.Error(), ctx.Kit.Rid)
		ctx.RespErrorCodeOnly(common.CCErrCommJSONUnmarshalFailed, "")
		return
	}

	// 给资源池目录加上资源池(业务id)和空闲机池（集群id）, service_category_id, service_template_id
	_, bizID, setID, err := s.getResourcePoolIDAndSetID(ctx)
	if err != nil {
		blog.ErrorJSON("CreateResourceDirectory fail with getResourcePoolIDAndSetID failed, err: %s, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	data[common.BKAppIDField] = bizID
	data[common.BKSetIDField] = setID
	data[common.BKServiceCategoryIDField] = 0
	data[common.BKServiceTemplateIDField] = 0
	data[common.BKInstParentStr] = setID
	data[common.BKSetTemplateIDField] = 0
	data.Set(common.BKOperatorField, "")
	data.Set(common.BKBakOperatorField, "")
	data[common.BKChildStr] = nil

	// 设置资源池自定义目录的default值
	data[common.BKDefaultField] = common.DefaultResSelfDefinedModuleFlag
	input := &metadata.CreateModelInstance{Data: data}
	rsp, err := s.Engine.CoreAPI.CoreService().Instance().CreateInstance(ctx.Kit.Ctx, ctx.Kit.Header, common.BKInnerObjIDModule, input)
	if err != nil {
		blog.ErrorJSON("CreateResourceDirectory, failed to CreateInstance, err: %s, rid: %s", err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	if !rsp.Result {
		blog.ErrorJSON("CreateResourceDirectory, failed to CreateInstance, errMsg: %s, rid: %s", rsp.ErrMsg, ctx.Kit.Rid)
		ctx.RespAutoError(errors.New(rsp.Code, rsp.ErrMsg))
		return
	}

	query := &metadata.QueryCondition{Condition: mapstr.MapStr{common.BKModuleIDField: rsp.Data.Created.ID}}
	readInstanceResult, err := s.Engine.CoreAPI.CoreService().Instance().ReadInstance(ctx.Kit.Ctx, ctx.Kit.Header, common.BKInnerObjIDModule, query)
	if err != nil {
		blog.ErrorJSON("CreateResourceDirectory success, but add host audit log failed, err: %s,rid:%s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	if !readInstanceResult.Result {
		blog.ErrorJSON("CreateResourceDirectory success, but add host audit log failed, err: %s,rid:%s", err, ctx.Kit.Rid)
		ctx.RespAutoError(errors.New(rsp.Code, rsp.ErrMsg))
		return
	}
	if len(readInstanceResult.Data.Info) <= 0 {
		err := fmt.Errorf("not find resource directory")
		blog.Errorf("create resource directory success, but add host audit log failed, err: %v, rid: %s",
			err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	// generate audit log.
	audit := auditlog.NewResourceDirAuditLog(s.Engine.CoreAPI.CoreService())
	generateAuditParameter := auditlog.NewGenerateAuditCommonParameter(ctx.Kit, metadata.AuditCreate)
	auditLog, err := audit.GenerateAuditLog(generateAuditParameter, int64(rsp.Data.Created.ID), bizID, readInstanceResult.Data.Info[0])
	if err != nil {
		blog.Errorf("generate audit log failed after create resource directory, err: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	// save audit log.
	if err := audit.SaveAuditLog(ctx.Kit, *auditLog); err != nil {
		blog.Errorf("save audit log failed after create resource directory, err: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	// register resource directory resource creator action to iam
	if auth.EnableAuthorize() {
		iamInstance := metadata.IamInstanceWithCreator{
			Type:    string(iam.SysResourcePoolDirectory),
			ID:      strconv.FormatUint(rsp.Data.Created.ID, 10),
			Name:    util.GetStrByInterface(readInstanceResult.Data.Info[0][common.BKModuleNameField]),
			Creator: ctx.Kit.User,
		}
		_, err = s.AuthManager.Authorizer.RegisterResourceCreatorAction(ctx.Kit.Ctx, ctx.Kit.Header, iamInstance)
		if err != nil {
			blog.Errorf("register created resource directory to iam failed, err: %s, rid: %s", err, ctx.Kit.Rid)
			ctx.RespAutoError(err)
			return
		}
	}

	ctx.RespEntity(rsp.Data)
}

func (s *Service) getResourcePoolIDAndSetID(ctx *rest.Contexts) (string, int64, int64, error) {
	query := &metadata.QueryCondition{
		Condition: mapstr.MapStr{common.BKDefaultField: 1},
	}
	bizRsp, err := s.Engine.CoreAPI.CoreService().Instance().ReadInstance(ctx.Kit.Ctx, ctx.Kit.Header, common.BKInnerObjIDApp, query)
	if err != nil {
		blog.ErrorJSON("getResourcePoolIDAndSetID, failed to find business by query condition: %s, err: %s, rid: %s", query, err.Error(), ctx.Kit.Rid)
		return "", 0, 0, err
	}
	if !bizRsp.Result {
		return "", 0, 0, ctx.Kit.CCError.New(bizRsp.Code, bizRsp.ErrMsg)
	}
	if bizRsp.Data.Count <= 0 {
		return "", 0, 0, fmt.Errorf("get resource pool info success, but count < 0")
	}

	intBizID, err := bizRsp.Data.Info[0].Int64(common.BKAppIDField)
	if err != nil {
		blog.Errorf("getResourcePoolIDAndSetID, bizID convert to float64 failed, err:%v, rid: %v", err, ctx.Kit.Rid)
		return "", 0, 0, err
	}
	bizName, err := bizRsp.Data.Info[0].String(common.BKAppNameField)
	if err != nil {
		blog.Errorf("getResourcePoolIDAndSetID, bizName convert to string failed, err:%v, rid: %v", err, ctx.Kit.Rid)
		return "", 0, 0, err
	}

	query.Condition = mapstr.MapStr{common.BKAppIDField: intBizID}
	setRsp, err := s.Engine.CoreAPI.CoreService().Instance().ReadInstance(ctx.Kit.Ctx, ctx.Kit.Header, common.BKInnerObjIDSet, query)
	if err != nil {
		blog.ErrorJSON("getResourcePoolIDAndSetID, failed to find business by query condition: %s, err: %s, rid: %s", query, err.Error(), ctx.Kit.Rid)
		return "", 0, 0, err
	}
	if !setRsp.Result {
		return "", 0, 0, ctx.Kit.CCError.New(setRsp.Code, setRsp.ErrMsg)
	}
	if setRsp.Data.Count <= 0 {
		return "", 0, 0, fmt.Errorf("get set info success, but count < 0")
	}

	intSetID, err := setRsp.Data.Info[0].Int64(common.BKSetIDField)
	if err != nil {
		blog.Errorf("getResourcePoolIDAndSetID, setID convert to float64 failed, err:%v, rid: %v", err, ctx.Kit.Rid)
		return "", 0, 0, err
	}

	return bizName, intBizID, intSetID, nil
}

func (s *Service) UpdateResourceDirectory(ctx *rest.Contexts) {
	input := mapstr.MapStr{}
	if err := ctx.DecodeInto(&input); nil != err {
		blog.Errorf("UpdateResourceDirectory failed to parse the params, error info is %s, rid: %s", err.Error(), ctx.Kit.Rid)
		ctx.RespErrorCodeOnly(common.CCErrCommJSONUnmarshalFailed, "")
		return
	}

	if !input.Exists(common.BKModuleNameField) {
		ctx.RespErrorCodeF(common.CCErrorTopoOnlyResourceDirNameCanBeUpdated, "UpdateResourceDirectory, field bk_module_name not exist, rid: %s", ctx.Kit.Rid)
		return
	}
	if len(input) > 1 {
		delete(input, common.BKModuleNameField)
		ctx.RespErrorCodeF(common.CCErrorTopoOnlyResourceDirNameCanBeUpdated, "UpdateResourceDirectory invalid params %s, rid: %s", input, ctx.Kit.Rid)
		return
	}

	moduleID := ctx.Request.PathParameter(common.BKModuleIDField)
	intModuleID, err := strconv.ParseInt(moduleID, 10, 64)
	if err != nil {
		blog.Errorf("DeleteResourceDirectory, moduleID convert to int64 failed, err:%v, rid: %v", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	_, bizID, _, err := s.getResourcePoolIDAndSetID(ctx)
	if err != nil {
		blog.ErrorJSON("failed to get resource pollID and setID in before update resource directory, err: %s, rid: %s",
			err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	// generate audit log.
	audit := auditlog.NewResourceDirAuditLog(s.Engine.CoreAPI.CoreService())
	generateAuditParameter := auditlog.NewGenerateAuditCommonParameter(ctx.Kit, metadata.AuditUpdate).WithUpdateFields(input)
	auditLog, err := audit.GenerateAuditLog(generateAuditParameter, intModuleID, bizID, nil)
	if err != nil {
		blog.Errorf("generate audit log failed before update resource directory, err: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	// to update.
	option := &metadata.UpdateOption{
		Data:      input,
		Condition: mapstr.MapStr{common.BKModuleIDField: intModuleID},
	}
	rsp, err := s.Engine.CoreAPI.CoreService().Instance().UpdateInstance(ctx.Kit.Ctx, ctx.Kit.Header, common.BKInnerObjIDModule, option)
	if err != nil {
		blog.Errorf("UpdateResourceDirectory failed, coreservice UpdateInstance http call fail, option: %v, err: %v, rid:%s", option.Data, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	if !rsp.Result {
		blog.ErrorJSON("UpdateResourceDirectory, failed to UpdateInstance, errMsg: %s, rid: %s", rsp.ErrMsg, ctx.Kit.Rid)
		ctx.RespAutoError(errors.New(rsp.Code, rsp.ErrMsg))
		return
	}

	// save audit log.
	if err := audit.SaveAuditLog(ctx.Kit, *auditLog); err != nil {
		blog.Errorf("save audit log failed after update resource directory, err: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(rsp.Data)
}

func (s *Service) SearchResourceDirectory(ctx *rest.Contexts) {
	input := new(metadata.SearchResourceDirParams)
	if err := ctx.DecodeInto(input); nil != err {
		blog.Errorf("SearchResourceDirectory failed to parse the params, error info is %s, rid: %s", err.Error(), ctx.Kit.Rid)
		ctx.RespErrorCodeOnly(common.CCErrCommJSONUnmarshalFailed, "")
		return
	}

	// if fuzzy search, change the string query to regexp
	if input.IsFuzzy == true {
		for k, v := range input.Condition {
			field, ok := v.(string)
			if ok {
				input.Condition[k] = mapstr.MapStr{
					common.BKDBLIKE: params.SpecialCharChange(field),
					"$options":      "i",
				}
			}
		}
	}

	_, bizID, setID, err := s.getResourcePoolIDAndSetID(ctx)
	if err != nil {
		blog.ErrorJSON("SearchResourceDirectory fail with getResourcePoolIDAndSetID failed, err: %s, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	if len(input.Condition) == 0 {
		input.Condition = mapstr.MapStr{}
	}
	input.Condition[common.BKAppIDField] = bizID
	input.Condition[common.BKSetIDField] = setID
	fields := input.Fields
	// must have fields bk_moudle_id and bk_module_name
	if !util.InArray(fields, common.BKModuleIDField) {
		fields = append(fields, common.BKModuleIDField)
	}
	if !util.InArray(fields, common.BKModuleNameField) {
		fields = append(fields, common.BKModuleNameField)
	}
	query := &metadata.QueryCondition{
		Fields: fields,
		Page: metadata.BasePage{
			Limit: common.BKNoLimit,
			Sort:  input.Page.Sort,
		},
		Condition: input.Condition,
	}
	rsp, err := s.Engine.CoreAPI.CoreService().Instance().ReadInstance(ctx.Kit.Ctx, ctx.Kit.Header, common.BKInnerObjIDModule, query)
	if err != nil {
		blog.Errorf("SearchResourceDirectory failed, coreservice http ReadInstance fail, input: %v, err: %v, %s", input, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	if !rsp.Result {
		blog.ErrorJSON("SearchResourceDirectory, failed to SearchResourceDirectory, errMsg: %s, rid: %s", rsp.ErrMsg, ctx.Kit.Rid)
		ctx.RespAutoError(errors.New(rsp.Code, rsp.ErrMsg))
		return
	}

	moduleIDArr := make([]int64, 0)
	mapModuleIdInfo := make(map[int64]mapstr.MapStr)
	IdleMoudleID := int64(0)
	for _, item := range rsp.Data.Info {
		moduleID, err := item.Int64(common.BKModuleIDField)
		if err != nil {
			blog.ErrorJSON("SearchResourceDirectory fail with moduleID convert from interface to int64 failed, err: %s, rid: %s", err, ctx.Kit.Rid)
			ctx.RespAutoError(err)
			return
		}
		defaultVal, err := item.Int64(common.BKDefaultField)
		if err != nil {
			blog.ErrorJSON("SearchResourceDirectory fail with defaultVal convert from interface to int64 failed, err: %s, rid: %s", err, ctx.Kit.Rid)
			ctx.RespAutoError(err)
			return
		}
		if int(defaultVal) == common.DefaultAppFlag {
			IdleMoudleID = moduleID
		} else {
			moduleIDArr = append(moduleIDArr, moduleID)

		}
		mapModuleIdInfo[moduleID] = item
	}
	// always put idle module in the first position
	if IdleMoudleID != int64(0) {
		moduleIDArr = append([]int64{IdleMoudleID}, moduleIDArr...)
	}

	// count hosts
	listHostOption := &metadata.HostModuleRelationRequest{
		ApplicationID: bizID,
		SetIDArr:      []int64{setID},
		ModuleIDArr:   moduleIDArr,
		Page: metadata.BasePage{
			Limit: common.BKNoLimit,
		},
	}
	hostModuleRelations, e := s.Engine.CoreAPI.CoreService().Host().GetHostModuleRelation(ctx.Kit.Ctx, ctx.Kit.Header, listHostOption)
	if e != nil {
		blog.Errorf("GetInternalModuleWithStatistics failed, list host modules failed, option: %+v, err: %s, rid: %s", listHostOption, e.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(e)
		return
	}
	moduleHostsCount := make(map[int64]int64)
	for _, item := range hostModuleRelations.Data.Info {
		if _, exist := moduleHostsCount[item.ModuleID]; exist == false {
			moduleHostsCount[item.ModuleID] = 0
		}
		moduleHostsCount[item.ModuleID] += 1
	}
	retInfo := make([]mapstr.MapStr, 0)
	for _, moduleID := range moduleIDArr {
		moduleInfo := mapModuleIdInfo[moduleID]
		moduleInfo["host_count"] = 0
		if count, exist := moduleHostsCount[moduleID]; exist == true {
			moduleInfo["host_count"] = count
		}
		retInfo = append(retInfo, moduleInfo)
	}

	ret := make(map[string]interface{}, 0)
	ret["count"] = rsp.Data.Count
	ret["info"] = retInfo
	ctx.RespEntity(ret)
}

func (s *Service) DeleteResourceDirectory(ctx *rest.Contexts) {
	moduleID := ctx.Request.PathParameter(common.BKModuleIDField)
	intModuleID, err := strconv.ParseInt(moduleID, 10, 64)
	if err != nil {
		blog.Errorf("DeleteResourceDirectory, moduleID convert to int64 failed, err:%v, rid: %v", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	_, bizID, setID, err := s.getResourcePoolIDAndSetID(ctx)
	if err != nil {
		blog.ErrorJSON("DeleteResourceDirectory fail with getResourcePoolIDAndSetID fail, err: %s, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	// 资源池目录下是否有主机，有主机的不能删除
	hasHost, err := s.hasHost(ctx, bizID, []int64{setID}, []int64{intModuleID})
	if err != nil {
		blog.Errorf("DeleteResourceDirectory, check if resource directory has host failed, err: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	if hasHost {
		blog.ErrorJSON("DeleteResourceDirectory fail, resource directory has host, rid: %s", ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.Error(common.CCErrTopoHasHostCheckFailed))
		return
	}

	// 资源池目录是否在云同步任务中被使用，使用了的不能删除
	syncDirs, err := s.GetResourceDirsInCloudSync(ctx)
	if err != nil {
		blog.Errorf("DeleteResourceDirectory failed, err: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	if _, ok := syncDirs[intModuleID]; ok {
		blog.Errorf("DeleteResourceDirectory failed, Resource dir is being used in cloud sync task, rid: %s", ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.Error(common.CCErrorTopoResourceDirUsedInCloudSync))
		return
	}

	language := util.GetLanguage(ctx.Kit.Header)
	query := &metadata.QueryCondition{Condition: mapstr.MapStr{common.BKModuleIDField: intModuleID}}
	curData, err := s.Engine.CoreAPI.CoreService().Instance().ReadInstance(ctx.Kit.Ctx, ctx.Kit.Header, common.BKInnerObjIDModule, query)
	if err != nil {
		blog.Errorf("DeleteResourceDirectory success but fail to create audiLog, coreservice http ReadInstance fail, err: %v, %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	if !curData.Result {
		blog.ErrorJSON("DeleteResourceDirectory success but fail to create audiLog, errMsg: %s, rid: %s", curData.ErrMsg, ctx.Kit.Rid)
		ctx.RespAutoError(errors.New(curData.Code, curData.ErrMsg))
		return
	}
	if len(curData.Data.Info) <= 0 {
		blog.Errorf("DeleteResourceDirectory fail, resource pool directory not exist, bk_module_id: %d, rid: %s", intModuleID, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrorTopoOperateReourceDirFailNotExist, s.Language.Language(language, "delete")))
		return
	}

	// 空闲机目录不能被删除
	moduleDefault, err := curData.Data.Info[0].Int64(common.BKDefaultField)
	if err != nil {
		blog.ErrorJSON("DeleteResourceDirectory fail, idle module can not delete, err: %s, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	if moduleDefault == 1 {
		blog.ErrorJSON("DeleteResourceDirectory fail, idle module can not delete, err: %s, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.Error(common.CCErrorTopoResourceDirIdleModuleCanNotRemove))
		return
	}

	// generate audit log.
	audit := auditlog.NewResourceDirAuditLog(s.Engine.CoreAPI.CoreService())
	generateAuditParameter := auditlog.NewGenerateAuditCommonParameter(ctx.Kit, metadata.AuditDelete)
	auditLog, err := audit.GenerateAuditLog(generateAuditParameter, intModuleID, bizID, curData.Data.Info[0])
	if err != nil {
		blog.Errorf("generate audit log failed before delete resource directory, err: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	// to delete.
	cond := &metadata.DeleteOption{Condition: mapstr.MapStr{common.BKModuleIDField: intModuleID}}
	rsp, err := s.Engine.CoreAPI.CoreService().Instance().DeleteInstance(ctx.Kit.Ctx, ctx.Kit.Header, common.BKInnerObjIDModule, cond)
	if err != nil {
		blog.Errorf("DeleteResourceDirectory, coreservice http DeleteInstance fail, bk_module_id: %d, err: %v, rid: %s")
		ctx.RespAutoError(err)
		return
	}
	if !rsp.Result {
		blog.ErrorJSON("DeleteResourceDirectory, failed to DeleteInstance, errMsg: %s, rid: %s", rsp.ErrMsg, ctx.Kit.Rid)
		ctx.RespAutoError(errors.New(rsp.Code, rsp.ErrMsg))
		return
	}

	// save audit log.
	if err := audit.SaveAuditLog(ctx.Kit, *auditLog); err != nil {
		blog.Errorf("save audit log failed after delete resource directory, err: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(rsp.Data)
}

func (s *Service) hasHost(ctx *rest.Contexts, bizID int64, setIDs, moduleIDS []int64) (bool, error) {
	option := &metadata.HostModuleRelationRequest{
		ApplicationID: bizID,
		ModuleIDArr:   moduleIDS,
	}
	if len(setIDs) > 0 {
		option.SetIDArr = setIDs
	}
	if len(moduleIDS) > 0 {
		option.ModuleIDArr = moduleIDS
	}
	rsp, err := s.Engine.CoreAPI.CoreService().Host().GetHostModuleRelation(ctx.Kit.Ctx, ctx.Kit.Header, option)
	if nil != err {
		blog.Errorf("[resource-directory] failed to request the object controller, err: %s, rid: %s", err.Error(), ctx.Kit.Rid)
		return false, ctx.Kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if !rsp.Result {
		blog.Errorf("[resource-directory]  failed to search the host module configures, err: %s, rid: %s", rsp.ErrMsg, ctx.Kit.Rid)
		return false, ctx.Kit.CCError.New(rsp.Code, rsp.ErrMsg)
	}

	return 0 != len(rsp.Data.Info), nil
}

// 获取云同步任务有关联的所有资源池目录
func (s *Service) GetResourceDirsInCloudSync(ctx *rest.Contexts) (map[int64]bool, error) {
	option := &metadata.SearchCloudOption{
		Page: metadata.BasePage{
			Limit: common.BKNoLimit,
		},
	}
	rsp, err := s.Engine.CoreAPI.CoreService().Cloud().SearchSyncTask(ctx.Kit.Ctx, ctx.Kit.Header, option)
	if nil != err {
		blog.Errorf("GetResourceDirsInCloudSync failed, err: %s, rid: %s", err.Error(), ctx.Kit.Rid)
		return nil, err
	}

	result := make(map[int64]bool)
	for _, task := range rsp.Info {
		for _, syncInfo := range task.SyncVpcs {
			result[syncInfo.SyncDir] = true
		}
	}

	return result, nil

}
