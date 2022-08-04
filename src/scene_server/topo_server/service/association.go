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
	"bytes"
	"io/ioutil"
	"sort"
	"strconv"

	"configcenter/src/ac/iam"
	"configcenter/src/ac/meta"
	"configcenter/src/common"
	"configcenter/src/common/auth"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

// CreateMainLineObject create a new model in the main line topo
func (s *Service) CreateMainLineObject(ctx *rest.Contexts) {
	data := new(metadata.MainlineAssociation)
	if err := ctx.DecodeInto(data); err != nil {
		ctx.RespAutoError(err)
		return
	}

	// 创建模型前，先创建表，避免模型创建后，对模型数据查询出现下面的错误，
	// (SnapshotUnavailable) Unable to read from a snapshot due to pending collection catalog changes;
	// please retry the operation. Snapshot timestamp is Timestamp(1616747877, 51).
	// Collection minimum is Timestamp(1616747878, 5)
	if err := s.createObjectTableByObjectID(ctx, data.ObjectID, true); err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, "mainline object"))
		return
	}

	ret := new(metadata.Object)
	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {

		var err error

		ret, err = s.Logics.AssociationOperation().CreateMainlineAssociation(ctx.Kit, data)
		if err != nil {
			blog.Errorf("create mainline object: %s failed, err: %v, rid: %s", data.ObjectID, err, ctx.Kit.Rid)
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

// DeleteMainLineObject delete a object int the main line topo
func (s *Service) DeleteMainLineObject(ctx *rest.Contexts) {
	objID := ctx.Request.PathParameter("bk_obj_id")

	// do with transaction
	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		if err := s.Logics.AssociationOperation().DeleteMainlineAssociation(ctx.Kit, objID); err != nil {
			blog.Errorf("delete mainline association failed, err: %+v, rid: %s", err, ctx.Kit.Rid)
			return ctx.Kit.CCError.CCErrorf(common.CCErrTopoObjectDeleteFailed, err.Error())
		}
		return nil
	})

	if txnErr != nil {
		ctx.RespAutoError(txnErr)
		return
	}
	ctx.RespEntity(nil)

}

// SearchMainLineObjectTopo search the main line topo
func (s *Service) SearchMainLineObjectTopo(ctx *rest.Contexts) {
	// get biz model related mainline models (mainline relationship model)
	resp, err := s.Logics.AssociationOperation().SearchMainlineAssociationTopo(ctx.Kit, common.BKInnerObjIDApp)
	if nil != err {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(resp)
}

// SearchObjectByClassificationID search the object by classification ID
func (s *Service) SearchObjectByClassificationID(ctx *rest.Contexts) {

	objID := ctx.Request.PathParameter(common.BKObjIDField)
	exist, err := s.Logics.ObjectOperation().IsObjectExist(ctx.Kit, objID)
	if err != nil {
		blog.Errorf("find the object(%s) failed, err: %v, rid: %s", objID, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	if !exist {
		blog.Errorf("object(%s) is non-exist, rid: %s", objID, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, common.BKObjIDField))
		return
	}

	resp, err := s.Logics.AssociationOperation().SearchMainlineAssociationTopo(ctx.Kit, objID)
	if nil != err {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(resp)
}

// SearchBusinessTopoWithStatistics calculate how many service instances on each topo instance node
func (s *Service) SearchBusinessTopoWithStatistics(ctx *rest.Contexts) {
	resp, err := s.searchBusinessTopo(ctx, true, true)
	if nil != err {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(resp)
}

// SearchBusinessTopo TODO
func (s *Service) SearchBusinessTopo(ctx *rest.Contexts) {
	resp, err := s.searchBusinessTopo(ctx, false, false)
	if nil != err {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(resp)
}

// searchBusinessTopo search the business topo
// withSortName 按拼音对名字排序，ui 专用
func (s *Service) searchBusinessTopo(ctx *rest.Contexts,
	withStatistics, withSortName bool) ([]*metadata.TopoInstRst, error) {
	id, err := strconv.ParseInt(ctx.Request.PathParameter("bk_biz_id"), 10, 64)
	if nil != err {
		blog.Errorf("failed to parse the path params id(%s), err: %v , rid: %s", ctx.Request.PathParameter("app_id"),
			err, ctx.Kit.Rid)

		return nil, err
	}

	withDefault := false
	if len(ctx.Request.QueryParameter("with_default")) > 0 {
		withDefault = true
	}

	topoInstRst, err := s.Logics.InstAssociationOperation().SearchMainlineAssociationInstTopo(ctx.Kit,
		common.BKInnerObjIDApp, id, withStatistics, withDefault)
	if err != nil {
		return nil, err
	}

	if withSortName {
		// sort before response,
		SortTopoInst(topoInstRst)
	}

	return topoInstRst, nil
}

// GetTopoNodeHostAndSerInstCount calculate how many service instances amd how many hosts on toponode
func (s *Service) GetTopoNodeHostAndSerInstCount(ctx *rest.Contexts) {
	urlBizID := ctx.Request.PathParameter(common.BKAppIDField)
	bizID, err := strconv.ParseInt(urlBizID, 10, 64)
	if err != nil {
		blog.Errorf("parse biz set id: %s from url failed, err: %v , rid: %s", urlBizID, err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField))
		return
	}

	input := new(metadata.HostAndSerInstCountOption)
	if err := ctx.DecodeInto(input); err != nil {
		ctx.RespAutoError(err)
		return
	}

	const BKParamMaxLength = 1000
	if len(input.Condition) > BKParamMaxLength {
		err := ctx.Kit.CCError.Errorf(common.CCErrCommParamsInvalid, "condition length")
		ctx.RespAutoError(err)
		return
	}

	// validate if the topo nodes are all in biz
	objInstMap := make(map[string][]int64)
	for _, node := range input.Condition {
		objInstMap[node.ObjID] = append(objInstMap[node.ObjID], node.InstID)
	}

	for objID, instIDs := range objInstMap {
		filter := []map[string]interface{}{{
			common.GetInstIDField(objID): mapstr.MapStr{common.BKDBIN: instIDs},
			common.BKAppIDField:          bizID,
		}}
		counts, err := s.Engine.CoreAPI.CoreService().Count().GetCountByFilter(ctx.Kit.Ctx, ctx.Kit.Header,
			common.GetInstTableName(objID, ctx.Kit.SupplierAccount), filter)
		if err != nil {
			blog.Errorf("count topo nodes failed, err: %v, filter: %#v, rid: %s", err, filter, ctx.Kit.Rid)
			ctx.RespAutoError(err)
			return
		}

		if int(counts[0]) != len(instIDs) {
			blog.Errorf("topo nodes are not all in biz, obj: %s, inst ids: %v, rid: %s", objID, instIDs, ctx.Kit.Rid)
			ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, "condition"))
			return
		}
	}

	result, err := s.Logics.InstAssociationOperation().TopoNodeHostAndSerInstCount(ctx.Kit, input)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(result)
}

// SortTopoInst TODO
func SortTopoInst(instData []*metadata.TopoInstRst) {
	for _, data := range instData {
		instNameInGBK, _ := ioutil.ReadAll(transform.NewReader(bytes.NewReader([]byte(data.InstName)), simplifiedchinese.GBK.NewEncoder()))
		data.InstName = string(instNameInGBK)
	}

	sort.Slice(instData, func(i, j int) bool {
		return instData[i].InstName < instData[j].InstName
	})

	for _, data := range instData {
		instNameInUTF, _ := ioutil.ReadAll(transform.NewReader(bytes.NewReader([]byte(data.InstName)), simplifiedchinese.GBK.NewDecoder()))
		data.InstName = string(instNameInUTF)
	}

	for idx := range instData {
		SortTopoInst(instData[idx].Child)
	}
}

// SearchBriefBizTopo search brief topo
func (s *Service) SearchBriefBizTopo(ctx *rest.Contexts) {
	bizID, err := strconv.ParseInt(ctx.Request.PathParameter(common.BKAppIDField), 10, 64)
	if err != nil {
		blog.Errorf("SearchBriefBizTopo failed, parse bk_biz_id error, err: %s, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, "bk_biz_id"))
		return
	}

	input := new(metadata.SearchBriefBizTopoOption)
	if err := ctx.DecodeInto(input); err != nil {
		ctx.RespAutoError(err)
		return
	}

	rawErr := input.Validate()
	if rawErr.ErrCode != 0 {
		ctx.RespAutoError(rawErr.ToCCError(ctx.Kit.CCError))
		return
	}

	setDetail, err := s.getSetDetailOfTopo(ctx, bizID, input)
	if err != nil {
		blog.Errorf("SearchBriefBizTopo failed, getSetDetailOfTopo err: %v, rid:%s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	if len(setDetail) == 0 {
		ctx.RespEntity([]interface{}{})
		return
	}

	moduleDetail, setModuleMap, err := s.getModuleInfoOfTopo(ctx, bizID, input)
	if err != nil {
		blog.Errorf("SearchBriefBizTopo failed, getModuleInfoOfTopo err: %v, rid:%s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	hostDetail, moduleHostMap, err := s.getHostInfoOfTopo(ctx, bizID, input)
	if err != nil {
		blog.Errorf("SearchBriefBizTopo failed, getHostInfoOfTopo err: %v, rid:%s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	// construct the final result
	bizTopo := s.constructBizTopo(setDetail, moduleDetail, hostDetail, setModuleMap, moduleHostMap)

	ctx.RespEntity(bizTopo)
}

// getSetDetailOfTopo get set detail of topo
func (s *Service) getSetDetailOfTopo(ctx *rest.Contexts, bizID int64, input *metadata.SearchBriefBizTopoOption) (
	map[int64]map[string]interface{}, errors.CCErrorCoder) {
	setDetail := make(map[int64]map[string]interface{})
	originSetFields := make(map[string]bool)
	for _, field := range input.SetFields {
		originSetFields[field] = true
	}
	input.SetFields = append(input.SetFields, common.BKSetIDField)

	pageSize := 2000
	start := 0
	hasNext := true
	param := &metadata.QueryCondition{
		Condition: map[string]interface{}{
			common.BKAppIDField: bizID,
		},
		Fields: input.SetFields,
		Page: metadata.BasePage{
			Start: start,
			Limit: pageSize,
		},
	}

	for hasNext {
		param.Page.Start = start
		setResult, err := s.Engine.CoreAPI.CoreService().Instance().ReadInstance(ctx.Kit.Ctx, ctx.Kit.Header,
			common.BKInnerObjIDSet, param)
		if nil != err {
			blog.Errorf("getSetDetailOfTopo failed, coreservice http ReadInstance fail, param: %v, err: %v, rid:%s",
				param, err, ctx.Kit.Rid)
			return nil, ctx.Kit.CCError.CCError(common.CCErrCommHTTPDoRequestFailed)
		}

		if len(setResult.Info) == 0 {
			break
		}

		for _, info := range setResult.Info {
			setID, _ := info.Int64(common.BKSetIDField)
			if !originSetFields[common.BKDefaultField] {
				info.Remove(common.BKDefaultField)
			}
			setDetail[setID] = info
		}

		start += pageSize
		if len(setResult.Info) < pageSize {
			hasNext = false
		}
	}

	return setDetail, nil
}

// getModuleInfoOfTopo get module info of topo
func (s *Service) getModuleInfoOfTopo(ctx *rest.Contexts, bizID int64, input *metadata.SearchBriefBizTopoOption) (
	map[int64]map[string]interface{}, map[int64][]int64, errors.CCErrorCoder) {
	//  get moduleDetail, setModuleMap
	moduleDetail := make(map[int64]map[string]interface{})
	setModuleMap := make(map[int64][]int64)

	originModuleFields := make(map[string]bool)
	for _, field := range input.ModuleFields {
		originModuleFields[field] = true
	}
	input.ModuleFields = append(input.ModuleFields, common.BKModuleIDField, common.BKSetIDField)

	pageSize := 2000
	start := 0
	hasNext := true
	param := &metadata.QueryCondition{
		Condition: map[string]interface{}{
			common.BKAppIDField: bizID,
		},
		Fields: input.ModuleFields,
		Page: metadata.BasePage{
			Start: start,
			Limit: pageSize,
		},
	}

	for hasNext {
		param.Page.Start = start
		moduleResult, err := s.Engine.CoreAPI.CoreService().Instance().ReadInstance(ctx.Kit.Ctx, ctx.Kit.Header,
			common.BKInnerObjIDModule, param)
		if nil != err {
			blog.Errorf("getModuleInfoOfTopo failed, coreservice http ReadInstance fail, param: %v, err: %v, rid:%s",
				param, err, ctx.Kit.Rid)
			return nil, nil, ctx.Kit.CCError.CCError(common.CCErrCommHTTPDoRequestFailed)
		}

		if len(moduleResult.Info) == 0 {
			break
		}

		for _, info := range moduleResult.Info {
			setID, _ := info.Int64(common.BKSetIDField)
			moduleID, _ := info.Int64(common.BKModuleIDField)
			setModuleMap[setID] = append(setModuleMap[setID], moduleID)

			if !originModuleFields[common.BKDefaultField] {
				info.Remove(common.BKDefaultField)
			}
			if !originModuleFields[common.BKSetIDField] {
				info.Remove(common.BKSetIDField)
			}
			moduleDetail[moduleID] = info
		}

		start += pageSize
		if len(moduleResult.Info) < pageSize {
			hasNext = false
		}
	}

	return moduleDetail, setModuleMap, nil
}

// getHostInfoOfTopo get host info of topo
func (s *Service) getHostInfoOfTopo(ctx *rest.Contexts, bizID int64, input *metadata.SearchBriefBizTopoOption) (
	map[int64]map[string]interface{}, map[int64][]int64, errors.CCErrorCoder) {
	hostDetail := make(map[int64]map[string]interface{})
	moduleHostMap := make(map[int64][]int64)

	// get hostIDArr, moduleHostMap
	hostIDArr := make([]int64, 0)
	relationOption := &metadata.HostModuleRelationRequest{
		ApplicationID: bizID,
		Page: metadata.BasePage{
			Limit: common.BKNoLimit,
		},
		Fields: []string{common.BKModuleIDField, common.BKHostIDField},
	}
	hostModuleRelations, err := s.Engine.CoreAPI.CoreService().Host().GetHostModuleRelation(ctx.Kit.Ctx,
		ctx.Kit.Header, relationOption)
	if err != nil {
		blog.Errorf("getHostInfoOfTopo failed, option: %+v, err: %s, rid: %s", relationOption, err.Error(), ctx.Kit.Rid)
		return nil, nil, ctx.Kit.CCError.CCError(common.CCErrCommHTTPDoRequestFailed)
	}
	for _, relation := range hostModuleRelations.Info {
		hostIDArr = append(hostIDArr, relation.HostID)
		moduleHostMap[relation.ModuleID] = append(moduleHostMap[relation.ModuleID], relation.HostID)
	}

	// get hostDetail
	if len(hostIDArr) > 0 {
		pageSize := 2000
		start := 0
		hasNext := true
		input.HostFields = append(input.HostFields, common.BKHostIDField)
		param := &metadata.QueryCondition{
			Condition: map[string]interface{}{
				common.BKHostIDField: map[string]interface{}{
					common.BKDBIN: hostIDArr,
				},
			},
			Fields: input.HostFields,
			Page: metadata.BasePage{
				Start: start,
				Limit: pageSize,
			},
		}

		for hasNext {
			param.Page.Start = start
			hostResult, err := s.Engine.CoreAPI.CoreService().Instance().ReadInstance(ctx.Kit.Ctx, ctx.Kit.Header,
				common.BKInnerObjIDHost, param)
			if nil != err {
				blog.Errorf("getHostInfoOfTopo failed, coreservice http ReadInstance fail, param: %v, err: %v, "+
					"rid:%s", param, err, ctx.Kit.Rid)
				return nil, nil, ctx.Kit.CCError.CCError(common.CCErrCommHTTPDoRequestFailed)
			}

			if len(hostResult.Info) == 0 {
				break
			}

			for _, info := range hostResult.Info {
				hostID, _ := info.Int64(common.BKHostIDField)
				hostDetail[hostID] = info
			}

			start += pageSize
			if len(hostResult.Info) < pageSize {
				hasNext = false
			}
		}
	}

	return hostDetail, moduleHostMap, nil
}

// constructBizTopo construct biz topo
func (s *Service) constructBizTopo(setDetail, moduleDetail, hostDetail map[int64]map[string]interface{}, setModuleMap,
	moduleHostMap map[int64][]int64) []*metadata.SetTopo {
	bizTopo := make([]*metadata.SetTopo, 0)
	for setID, set := range setDetail {
		setTopo := new(metadata.SetTopo)
		setTopo.Set = set
		moduleTopos := make([]*metadata.ModuleTopo, 0)
		for _, moduleID := range setModuleMap[setID] {
			moduleTopo := new(metadata.ModuleTopo)
			moduleTopo.Module = moduleDetail[moduleID]
			hosts := make([]map[string]interface{}, 0)
			for _, hostID := range moduleHostMap[moduleID] {
				hosts = append(hosts, hostDetail[hostID])
			}
			moduleTopo.Hosts = hosts
			moduleTopos = append(moduleTopos, moduleTopo)
		}
		setTopo.ModuleTopos = moduleTopos
		bizTopo = append(bizTopo, setTopo)
	}
	return bizTopo
}

// SearchAssociationType search association type
func (s *Service) SearchAssociationType(ctx *rest.Contexts) {
	request := &metadata.SearchAssociationTypeRequest{}
	if err := ctx.DecodeInto(request); err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.New(common.CCErrCommParamsInvalid, err.Error()))
		return
	}
	if request.Condition == nil {
		request.Condition = make(map[string]interface{}, 0)
	}

	input := &metadata.QueryCondition{
		Condition: request.Condition,
		Page:      request.BasePage,
	}
	ret, err := s.Engine.CoreAPI.CoreService().Association().ReadAssociationType(ctx.Kit.Ctx, ctx.Kit.Header, input)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(ret)
}

// SearchObjectAssocWithAssocKindList search object association by association kind
func (s *Service) SearchObjectAssocWithAssocKindList(ctx *rest.Contexts) {

	ids := new(metadata.AssociationKindIDs)
	if err := ctx.DecodeInto(ids); err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommParamsInvalid))
		return
	}

	resp, err := s.Logics.AssociationOperation().SearchObjectAssocWithAssocKindList(ctx.Kit, ids.AsstIDs)
	if nil != err {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(resp)
}

// CreateAssociationType create association kind
func (s *Service) CreateAssociationType(ctx *rest.Contexts) {
	request := &metadata.AssociationKind{}
	if err := ctx.DecodeInto(request); err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.New(common.CCErrCommParamsInvalid, err.Error()))
		return
	}

	var ret *metadata.CreateOneDataResult
	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		var err error
		input := &metadata.CreateAssociationKind{Data: *request}
		ret, err = s.Engine.CoreAPI.CoreService().Association().CreateAssociationType(ctx.Kit.Ctx, ctx.Kit.Header,
			input)
		if err != nil {
			return err
		}

		// register association type resource creator action to iam
		if auth.EnableAuthorize() {
			iamInstance := metadata.IamInstanceWithCreator{
				Type:    string(iam.SysAssociationType),
				ID:      strconv.FormatInt(int64(ret.Created.ID), 10),
				Name:    request.AssociationKindName,
				Creator: ctx.Kit.User,
			}
			_, err = s.AuthManager.Authorizer.RegisterResourceCreatorAction(ctx.Kit.Ctx, ctx.Kit.Header, iamInstance)
			if err != nil {
				blog.Errorf("register created association type to iam failed, err: %v, rid: %s", err, ctx.Kit.Rid)
				return err
			}
		}
		return nil
	})

	if txnErr != nil {
		ctx.RespAutoError(txnErr)
		return
	}
	ctx.RespEntity(metadata.RspID{ID: int64(ret.Created.ID)})
}

// UpdateAssociationType update association kind
func (s *Service) UpdateAssociationType(ctx *rest.Contexts) {
	request := &metadata.UpdateAssociationTypeRequest{}
	if err := ctx.DecodeInto(request); err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.New(common.CCErrCommParamsInvalid, err.Error()))
		return
	}

	asstTypeID, err := strconv.ParseInt(ctx.Request.PathParameter("id"), 10, 64)
	if err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.New(common.CCErrCommParamsInvalid, err.Error()))
		return
	}

	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		var err error

		input := metadata.UpdateOption{
			Condition: mapstr.MapStr{common.BKFieldID: asstTypeID},
			Data: mapstr.MapStr{
				"bk_asst_name": request.AsstName,
				"src_des":      request.SrcDes,
				"dest_des":     request.DestDes,
				"direction":    request.Direction,
			},
		}

		_, err = s.Engine.CoreAPI.CoreService().Association().UpdateAssociationType(ctx.Kit.Ctx, ctx.Kit.Header, &input)
		if err != nil {
			blog.Errorf("update association type failed, input: %#v, err: %v, rid: %s", input, err, ctx.Kit.Rid)
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

// DeleteAssociationType delete association kind
func (s *Service) DeleteAssociationType(ctx *rest.Contexts) {
	asstTypeID, err := strconv.ParseInt(ctx.Request.PathParameter("id"), 10, 64)
	if err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.New(common.CCErrCommParamsInvalid, err.Error()))
		return
	}

	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		err := s.Logics.AssociationOperation().DeleteAssociationType(ctx.Kit, asstTypeID)
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

// SearchInstanceAssociations searches object instance associations with the input conditions.
func (s *Service) SearchInstanceAssociations(ctx *rest.Contexts) {
	objID := ctx.Request.PathParameter("bk_obj_id")

	// decode input parameter.
	input := &metadata.CommonSearchFilter{}
	if err := ctx.DecodeInto(input); nil != err {
		ctx.RespAutoError(err)
		return
	}
	input.ObjectID = objID

	// validate input parameter.
	if invalidKey, err := input.Validate(); err != nil {
		blog.Errorf("validate search instance associations input parameters failed, err: %s, rid: %s",
			err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, invalidKey))
		return
	}

	// set read preference.
	ctx.SetReadPreference(common.SecondaryPreferredMode)

	// search instance associations.
	result, err := s.Logics.InstAssociationOperation().SearchInstanceAssociations(ctx.Kit, objID, input)
	if err != nil {
		blog.Errorf("search object[%s] instance associations failed, err: %s, rid: %s", objID, err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(mapstr.MapStr{"info": result})
}

// CountInstanceAssociations counts object instance associations with the input conditions.
func (s *Service) CountInstanceAssociations(ctx *rest.Contexts) {
	objID := ctx.Request.PathParameter("bk_obj_id")

	// decode input parameter.
	input := &metadata.CommonCountFilter{}
	if err := ctx.DecodeInto(input); nil != err {
		ctx.RespAutoError(err)
		return
	}
	input.ObjectID = objID

	// validate input parameter.
	if invalidKey, err := input.Validate(); err != nil {
		blog.Errorf("validate count instance associations input parameters failed, err: %s, rid: %s",
			err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, invalidKey))
		return
	}

	// set read preference.
	ctx.SetReadPreference(common.SecondaryPreferredMode)

	// count instance associations.
	result, err := s.Logics.InstAssociationOperation().CountInstanceAssociations(ctx.Kit, objID, input)
	if err != nil {
		blog.Errorf("count object[%s] instance associations failed, err: %s, rid: %s", objID, err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(result)
}

// SearchAssociationInst search instance association
func (s *Service) SearchAssociationInst(ctx *rest.Contexts) {
	request := &metadata.SearchAssociationInstRequest{}
	if err := ctx.DecodeInto(request); err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.New(common.CCErrCommParamsInvalid, err.Error()))
		return
	}

	if len(request.ObjID) == 0 {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedSet, common.BKObjIDField))
		return
	}

	cond := &metadata.InstAsstQueryCondition{
		Cond:  metadata.QueryCondition{Condition: request.Condition},
		ObjID: request.ObjID,
	}

	ctx.SetReadPreference(common.SecondaryPreferredMode)
	ret, err := s.Engine.CoreAPI.CoreService().Association().ReadInstAssociation(ctx.Kit.Ctx, ctx.Kit.Header, cond)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(ret.Info)
}

// SearchAssociationRelatedInst search all associations of certain model instance, including it as source & target
func (s *Service) SearchAssociationRelatedInst(ctx *rest.Contexts) {
	request := &metadata.SearchAssociationRelatedInstRequest{}
	if err := ctx.DecodeInto(request); err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrCommParamsInvalid, err.Error()))
		return
	}

	// check condition
	if request.Condition.InstID == 0 || request.Condition.ObjectID == "" {
		ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrCommParamsNeedSet, "'bk_inst_id' & 'bk_obj_id'"))
		return
	}

	// check fields,if there's none param,return err.
	if len(request.Fields) == 0 {
		ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrCommParamsNeedSet, "fields"))
		return
	}

	// Use id as sort parameters
	request.Page.Sort = common.BKFieldID

	// check Maximum limit
	if request.Page.Limit > common.BKMaxInstanceLimit {
		ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrCommPageLimitIsExceeded))
		return
	}

	queryCond := &metadata.InstAsstQueryCondition{
		ObjID: request.Condition.ObjectID,
		Cond: metadata.QueryCondition{
			Fields: request.Fields,
			Page:   request.Page,
			Condition: mapstr.MapStr{
				common.BKDBOR: []mapstr.MapStr{
					{
						common.BKObjIDField:  request.Condition.ObjectID,
						common.BKInstIDField: request.Condition.InstID,
					},
					{
						common.BKAsstObjIDField:  request.Condition.ObjectID,
						common.BKAsstInstIDField: request.Condition.InstID,
					},
				},
			},
		},
	}

	res, err := s.Engine.CoreAPI.CoreService().Association().ReadInstAssociation(ctx.Kit.Ctx, ctx.Kit.Header, queryCond)
	if err != nil {
		blog.Errorf("search instance association failed, err: %v, query: %#v, rid: %s", err, queryCond, ctx.Kit.Rid)
		ctx.RespAutoError(err)
	}

	ctx.RespEntity(res.Info)
}

// SearchInstAssociationAndInstDetail search association, source object inst and destination object inst
// related issue: https://github.com/Tencent/bk-cmdb/issues/5807
func (s *Service) SearchInstAssociationAndInstDetail(ctx *rest.Contexts) {
	request := metadata.InstAndAssocRequest{}
	if err := ctx.DecodeInto(&request); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if ccErr := request.Validate(); ccErr.ErrCode != 0 {
		blog.Errorf("validate request failed, err: %v, rid: %s", ccErr, ctx.Kit.Rid)
		ctx.RespAutoError(ccErr.ToCCError(ctx.Kit.CCError))
		return
	}

	assocFilter, errKey, err := request.Condition.AsstFilter.ToMgo()
	if err != nil {
		blog.Errorf("parse asst_filter to mongo searcher failed, errKey: %s, err: %v", errKey, err)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, errKey))
		return
	}

	objID := ctx.Request.PathParameter(common.BKObjIDField)
	cond := &metadata.InstAsstQueryCondition{
		Cond: metadata.QueryCondition{
			Condition:      assocFilter,
			Fields:         request.Condition.AsstFields,
			Page:           request.Page,
			DisableCounter: true,
		},
		ObjID: objID,
	}

	ctx.SetReadPreference(common.SecondaryPreferredMode)
	assocRsp, err := s.Engine.CoreAPI.CoreService().Association().ReadInstAssociation(ctx.Kit.Ctx, ctx.Kit.Header, cond)
	if err != nil {
		blog.Errorf("search inst association failed, err: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	resp := metadata.InstAndAssocDetailData{Asst: assocRsp.Info}
	if !request.Condition.SrcDetail && !request.Condition.DstDetail {
		ctx.RespEntity(resp)
		return
	}

	srcObjInst := make(map[string][]int64)
	dstObjInst := make(map[string][]int64)
	for _, item := range assocRsp.Info {
		srcObjInst[item.ObjectID] = append(srcObjInst[item.ObjectID], item.InstID)
		dstObjInst[item.AsstObjectID] = append(dstObjInst[item.AsstObjectID], item.AsstInstID)
	}

	if request.Condition.SrcDetail {
		for obj, insts := range srcObjInst {
			srcRsp, err := s.searchInstForAssocDetail(ctx, obj, request.Condition.SrcFields, insts)
			if err != nil {
				blog.Errorf("search src inst failed, err: %v, rid: %s", err, ctx.Kit.Rid)
				ctx.RespAutoError(err)
				return
			}

			resp.Src = append(resp.Src, srcRsp...)
		}
	}

	if request.Condition.DstDetail {
		for obj, insts := range dstObjInst {
			dstRsp, err := s.searchInstForAssocDetail(ctx, obj, request.Condition.DstFields, insts)
			if err != nil {
				blog.Errorf("search dst inst failed, err: %v, rid: %s", err, ctx.Kit.Rid)
				ctx.RespAutoError(err)
				return
			}

			resp.Dst = append(resp.Dst, dstRsp...)
		}
	}

	ctx.RespEntity(resp)
}

func (s *Service) searchInstForAssocDetail(ctx *rest.Contexts, objID string, fields []string, insts []int64) (
	[]mapstr.MapStr, error) {

	isMainline, err := s.Logics.AssociationOperation().IsMainlineObject(ctx.Kit, objID)
	if err != nil {
		blog.Errorf("check object is mainline failed, err: %v, rid: %s", err, ctx.Kit.Rid)
		return nil, err
	}

	// common object's auth is different from mainline object, mainline need to check business id
	// mainline object will verify whether user has the viewBusinessResource of the business which the object belongs
	action := meta.FindMany
	authObject := objID
	authInst := insts
	if isMainline {
		action = meta.ViewBusinessResource
		authObject = common.BKInnerObjIDApp
		authInst = make([]int64, 0)
		if objID == common.BKInnerObjIDHost {
			input := &metadata.HostModuleRelationRequest{HostIDArr: insts}
			relation, err := s.Engine.CoreAPI.CoreService().Host().GetHostModuleRelation(ctx.Kit.Ctx, ctx.Kit.Header,
				input)
			if err != nil {
				blog.Errorf("search host business relation failed, err: %v, rid: %s", err, ctx.Kit.Rid)
				return nil, err
			}
			for _, item := range relation.Info {
				authInst = append(authInst, item.AppID)
			}
		} else {
			cond := &metadata.QueryCondition{
				Condition:      mapstr.MapStr{common.GetInstIDField(objID): mapstr.MapStr{common.BKDBIN: insts}},
				Fields:         []string{common.BKAppIDField},
				DisableCounter: true,
			}
			rsp, err := s.Logics.InstOperation().FindInst(ctx.Kit, objID, cond)
			if err != nil {
				blog.Errorf("search object[%s] instance failed, err: %v, rid: %s", objID, err, ctx.Kit.Rid)
				return nil, err
			}

			for _, item := range rsp.Info {
				bizID, err := item.Int64(common.BKAppIDField)
				if err != nil {
					blog.Errorf("get object[%s] instance parentID failed, err: %v, rid: %s", objID, err, ctx.Kit.Rid)
					return nil, err
				}
				authInst = append(authInst, bizID)
			}
		}
	}

	err = s.AuthManager.AuthorizeByInstanceID(ctx.Kit.Ctx, ctx.Kit.Header, action, authObject, authInst...)
	if err != nil {
		blog.Errorf("authorize object[%s] failed, action: %s, insts: %v, err: %v, rid: %s", objID, action, insts, err,
			ctx.Kit.Rid)
		return nil, err
	}

	query := &metadata.QueryCondition{
		Condition:      mapstr.MapStr{common.GetInstIDField(objID): mapstr.MapStr{common.BKDBIN: insts}},
		Fields:         fields,
		DisableCounter: true,
	}
	rsp, err := s.Logics.InstOperation().FindInst(ctx.Kit, objID, query)
	if err != nil {
		blog.Errorf("search object[%s] insts failed, err: %v, rid: %s", objID, err, ctx.Kit.Rid)
		return nil, err
	}
	return rsp.Info, nil
}

// CreateAssociationInst create instance associaiton
func (s *Service) CreateAssociationInst(ctx *rest.Contexts) {
	request := &metadata.CreateAssociationInstRequest{}
	if err := ctx.DecodeInto(request); err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.New(common.CCErrCommParamsInvalid, err.Error()))
		return
	}

	var ret *metadata.RspID
	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		var err error
		ret, err = s.Logics.InstAssociationOperation().CreateInstanceAssociation(ctx.Kit, request)
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

// CreateManyInstAssociation batch create instance association
func (s *Service) CreateManyInstAssociation(ctx *rest.Contexts) {
	request := &metadata.CreateManyInstAsstRequest{}
	if err := ctx.DecodeInto(request); err != nil {
		blog.Errorf("deserialization failed, err: %s, rid: %s", err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	ret, err := s.Logics.InstAssociationOperation().CreateManyInstAssociation(ctx.Kit, request)
	if err != nil {
		blog.Errorf("create many instance association failed, err: %s, rid: %s", err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(ret)
}

// DeleteAssociationInst delete instance association
func (s *Service) DeleteAssociationInst(ctx *rest.Contexts) {
	objID := ctx.Request.PathParameter(common.BKObjIDField)
	if len(objID) == 0 {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedSet, common.BKObjIDField))
		return
	}

	id, err := strconv.ParseInt(ctx.Request.PathParameter("association_id"), 10, 64)
	if err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedInt, "association_id"))
		return
	}

	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		_, err := s.Logics.InstAssociationOperation().DeleteInstAssociation(ctx.Kit, objID, []int64{id})
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

// DeleteAssociationInstBatch batch delete instance association
func (s *Service) DeleteAssociationInstBatch(ctx *rest.Contexts) {
	request := &metadata.DeleteAssociationInstBatchRequest{}
	if err := ctx.DecodeInto(request); err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.New(common.CCErrCommParamsInvalid, err.Error()))
		return
	}
	if len(request.ID) == 0 {
		ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrCommHTTPInputInvalid))
		return
	}
	if len(request.ID) > common.BKMaxInstanceLimit {
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommPageLimitIsExceeded))
		return
	}
	if len(request.ObjectID) == 0 {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedSet, common.BKObjIDField))
		return
	}

	result := &metadata.DeleteAssociationInstBatchResult{}
	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {

		rsp, err := s.Logics.InstAssociationOperation().DeleteInstAssociation(ctx.Kit, request.ObjectID, request.ID)
		if err != nil {
			return err
		}

		result.Data = int(rsp)
		return nil
	})
	if txnErr != nil {
		ctx.RespAutoError(txnErr)
		return
	}

	ctx.RespEntity(result.Data)
}

// SearchTopoPath TODO
func (s *Service) SearchTopoPath(ctx *rest.Contexts) {
	rid := ctx.Kit.Rid

	bizIDStr := ctx.Request.PathParameter(common.BKAppIDField)
	bizID, err := strconv.ParseInt(bizIDStr, 10, 64)
	if nil != err {
		blog.Errorf("SearchTopoPath failed, bizIDStr: %s, err: %s, rid: %s", bizIDStr, err.Error(), rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField))
		return
	}

	input := metadata.FindTopoPathRequest{}
	if err := ctx.DecodeInto(&input); err != nil {
		ctx.RespAutoError(err)
		return
	}
	if len(input.Nodes) == 0 {
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommHTTPBodyEmpty))
		return
	}

	topoRoot, err := s.Engine.CoreAPI.CoreService().Mainline().SearchMainlineInstanceTopo(ctx.Kit.Ctx, ctx.Kit.Header, bizID, false)
	if err != nil {
		blog.Errorf("SearchTopoPath failed, SearchMainlineInstanceTopo failed, bizID:%d, err:%s, rid:%s", bizID, err.Error(), rid)
		ctx.RespAutoError(err)
		return
	}
	result := metadata.TopoPathResult{}
	for _, node := range input.Nodes {
		topoPath := topoRoot.TraversalFindNode(node.ObjectID, node.InstanceID)
		path := make([]*metadata.TopoInstanceNodeSimplify, 0)
		for _, item := range topoPath {
			simplify := item.ToSimplify()
			path = append(path, simplify)
		}
		nodeTopoPath := metadata.NodeTopoPath{
			BizID: bizID,
			Node:  node,
			Path:  path,
		}
		result.Nodes = append(result.Nodes, nodeTopoPath)
	}

	ctx.RespEntity(result)
}
