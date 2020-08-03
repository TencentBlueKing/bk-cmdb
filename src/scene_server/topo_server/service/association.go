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

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/topo_server/core/model"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

// CreateMainLineObject create a new model in the main line topo
func (s *Service) CreateMainLineObject(ctx *rest.Contexts) {
	data := make(map[string]interface{})
	if err := ctx.DecodeInto(&data); err != nil {
		ctx.RespAutoError(err)
		return
	}
	mainLineAssociation := &metadata.Association{}
	_, err := mainLineAssociation.Parse(data)
	if nil != err {
		blog.Errorf("[api-asst] failed to parse the data(%#v), error info is %s, rid: %s", data, err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, "mainline object"))
		return
	}

	var ret model.Object
	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, s.EnableTxn, ctx.Kit.Header, func() error {
		var err error
		ret, err = s.Core.AssociationOperation().CreateMainlineAssociation(ctx.Kit, mainLineAssociation, s.Config.BusinessTopoLevelMax)
		if err != nil {
			blog.Errorf("create mainline object: %s failed, err: %v, rid: %s", mainLineAssociation.ObjectID, err, ctx.Kit.Rid)
			return err
		}

		// auth: register mainline object
		if err := s.AuthManager.RegisterMainlineObject(ctx.Kit.Ctx, ctx.Kit.Header, ret.Object()); err != nil {
			blog.Errorf("create mainline object success, but register mainline model to iam failed, err: %+v, rid: %s", err, ctx.Kit.Rid)
			return ctx.Kit.CCError.Error(common.CCErrCommRegistResourceToIAMFailed)
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
	var bizID int64
	md := new(MetaShell)
	if err := ctx.DecodeInto(md); err != nil {
		ctx.RespAutoError(err)
		return
	}
	if md.Metadata != nil {
		var err error
		bizID, err = metadata.BizIDFromMetadata(*md.Metadata)
		if err != nil {
			blog.Errorf("parse business id from request failed, err: %+v, rid: %s", err, ctx.Kit.Rid)
			ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommParamsInvalid))
			return
		}
	}

	// auth: collection iam resource before it really be deleted
	iamResources, err := s.AuthManager.MakeResourcesByObjectIDs(ctx.Kit.Ctx, ctx.Kit.Header, bizID, objID)
	if err != nil {
		blog.Errorf("MakeResourcesByObjectIDs failed, err: %+v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrTopoObjectDeleteFailed))
		return
	}

	// do with transaction
	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, s.EnableTxn, ctx.Kit.Header, func() error {
		if err := s.Core.AssociationOperation().DeleteMainlineAssociation(ctx.Kit, objID, md.Metadata); err != nil {
			blog.Errorf("DeleteMainlineAssociation failed, err: %+v, rid: %s", err, ctx.Kit.Rid)
			return ctx.Kit.CCError.CCError(common.CCErrTopoObjectDeleteFailed)
		}

		// auth: do deregister
		if err := s.AuthManager.Authorize.DeregisterResource(ctx.Kit.Ctx, iamResources...); err != nil {
			blog.Errorf("delete mainline association success, but deregister mainline model failed, err: %+v, rid: %s", err, ctx.Kit.Rid)
			return ctx.Kit.CCError.CCError(common.CCErrCommUnRegistResourceToIAMFailed)
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
	md := new(MetaShell)
	if err := ctx.DecodeInto(md); err != nil {
		ctx.RespAutoError(err)
		return
	}
	bizObj, err := s.Core.ObjectOperation().FindSingleObject(ctx.Kit, common.BKInnerObjIDApp, md.Metadata)
	if nil != err {
		blog.Errorf("[api-asst] failed to find the biz object, error info is %s, rid: %s", err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	// get biz model related mainline models (mainline relationship model)
	resp, err := s.Core.AssociationOperation().SearchMainlineAssociationTopo(ctx.Kit, bizObj)
	if nil != err {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(resp)
}

// SearchObjectByClassificationID search the object by classification ID
func (s *Service) SearchObjectByClassificationID(ctx *rest.Contexts) {
	md := new(MetaShell)
	if err := ctx.DecodeInto(md); err != nil {
		ctx.RespAutoError(err)
		return
	}
	bizObj, err := s.Core.ObjectOperation().FindSingleObject(ctx.Kit, ctx.Request.PathParameter("bk_obj_id"), md.Metadata)
	if nil != err {
		blog.Errorf("[api-asst] failed to find the biz object, error info is %s, rid: %s", err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	resp, err := s.Core.AssociationOperation().SearchMainlineAssociationTopo(ctx.Kit, bizObj)
	if nil != err {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(resp)
}

// SearchBusinessTopoWithStatistics calculate how many service instances on each topo instance node
func (s *Service) SearchBusinessTopoWithStatistics(ctx *rest.Contexts) {
	resp, err := s.searchBusinessTopo(ctx, true)
	if nil != err {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(resp)
}

func (s *Service) SearchBusinessTopo(ctx *rest.Contexts) {
	resp, err := s.searchBusinessTopo(ctx, false)
	if nil != err {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(resp)
}

// SearchBusinessTopo search the business topo
func (s *Service) searchBusinessTopo(ctx *rest.Contexts, withStatistics bool) ([]*metadata.TopoInstRst, error) {
	id, err := strconv.ParseInt(ctx.Request.PathParameter("bk_biz_id"), 10, 64)
	if nil != err {
		blog.Errorf("[api-asst] failed to parse the path params id(%s), error info is %s , rid: %s", ctx.Request.PathParameter("app_id"), err.Error(), ctx.Kit.Rid)

		return nil, err
	}

	md := new(MetaShell)
	if err := ctx.DecodeInto(md); err != nil {
		return nil, err
	}

	withDefault := false
	if len(ctx.Request.QueryParameter("with_default")) > 0 {
		withDefault = true
	}
	topoInstRst, err := s.Core.AssociationOperation().SearchMainlineAssociationInstTopo(ctx.Kit, common.BKInnerObjIDApp, id, withStatistics, withDefault, md.Metadata)
	if err != nil {
		return nil, err
	}

	// sort before response
	SortTopoInst(topoInstRst)

	return topoInstRst, nil
}

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

// SearchMainLineChildInstTopo search the child inst topo by a inst
func (s *Service) SearchMainLineChildInstTopo(ctx *rest.Contexts) {

	// {obj_id}/{app_id}/{inst_id}
	objID := ctx.Request.PathParameter("obj_id")
	bizID, err := strconv.ParseInt(ctx.Request.PathParameter("app_id"), 10, 64)
	if nil != err {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, "app_id"))
		return
	}

	// get the instance id of this object.
	instID, err := strconv.ParseInt(ctx.Request.PathParameter("inst_id"), 10, 64)
	if nil != err {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, "inst_id"))
		return
	}
	_ = bizID

	md := new(MetaShell)
	if err := ctx.DecodeInto(md); err != nil {
		ctx.RespAutoError(err)
		return
	}

	resp, err := s.Core.AssociationOperation().SearchMainlineAssociationInstTopo(ctx.Kit, objID, instID, false, false, md.Metadata)
	if nil != err {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(resp)
}

func (s *Service) SearchAssociationType(ctx *rest.Contexts) {
	request := &metadata.SearchAssociationTypeRequest{}
	if err := ctx.DecodeInto(request); err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.New(common.CCErrCommParamsInvalid, err.Error()))
		return
	}
	if request.Condition == nil {
		request.Condition = make(map[string]interface{}, 0)
	}

	ret, err := s.Core.AssociationOperation().SearchType(ctx.Kit, request)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	if ret.Code != 0 {
		ctx.RespAutoError(ctx.Kit.CCError.New(ret.Code, ret.ErrMsg))
		return
	}

	ctx.RespEntity(ret.Data)
}

func (s *Service) SearchObjectAssocWithAssocKindList(ctx *rest.Contexts) {

	ids := new(metadata.AssociationKindIDs)
	if err := ctx.DecodeInto(ids); err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommParamsInvalid))
		return
	}

	resp, err := s.Core.AssociationOperation().SearchObjectAssocWithAssocKindList(ctx.Kit, ids.AsstIDs)
	if nil != err {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(resp)
}

func (s *Service) CreateAssociationType(ctx *rest.Contexts) {
	request := &metadata.AssociationKind{}
	if err := ctx.DecodeInto(request); err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.New(common.CCErrCommParamsInvalid, err.Error()))
		return
	}

	var ret *metadata.CreateAssociationTypeResult
	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, s.EnableTxn, ctx.Kit.Header, func() error {
		var err error
		ret, err = s.Core.AssociationOperation().CreateType(ctx.Kit, request)
		if err != nil {
			return err
		}

		if ret.Code != 0 {
			return ctx.Kit.CCError.New(ret.Code, ret.ErrMsg)
		}
		return nil
	})

	if txnErr != nil {
		ctx.RespAutoError(txnErr)
		return
	}
	ctx.RespEntity(ret.Data)
}

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

	var ret *metadata.UpdateAssociationTypeResult
	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, s.EnableTxn, ctx.Kit.Header, func() error {
		var err error
		ret, err = s.Core.AssociationOperation().UpdateType(ctx.Kit, asstTypeID, request)
		if err != nil {
			return err
		}

		if ret.Code != 0 {
			return ctx.Kit.CCError.New(ret.Code, ret.ErrMsg)
		}
		return nil
	})

	if txnErr != nil {
		ctx.RespAutoError(txnErr)
		return
	}
	ctx.RespEntity(ret.Data)
}

func (s *Service) DeleteAssociationType(ctx *rest.Contexts) {
	asstTypeID, err := strconv.ParseInt(ctx.Request.PathParameter("id"), 10, 64)
	if err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.New(common.CCErrCommParamsInvalid, err.Error()))
		return
	}

	var ret *metadata.DeleteAssociationTypeResult
	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, s.EnableTxn, ctx.Kit.Header, func() error {
		var err error
		ret, err = s.Core.AssociationOperation().DeleteType(ctx.Kit, asstTypeID)
		if err != nil {
			return err
		}

		if ret.Code != 0 {
			return ctx.Kit.CCError.New(ret.Code, ret.ErrMsg)
		}
		return nil
	})

	if txnErr != nil {
		ctx.RespAutoError(txnErr)
		return
	}
	ctx.RespEntity(ret.Data)
}

func (s *Service) SearchAssociationInst(ctx *rest.Contexts) {
	request := &metadata.SearchAssociationInstRequest{}
	if err := ctx.DecodeInto(request); err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.New(common.CCErrCommParamsInvalid, err.Error()))
		return
	}

	ret, err := s.Core.AssociationOperation().SearchInst(ctx.Kit, request)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	if ret.Code != 0 {
		ctx.RespAutoError(ctx.Kit.CCError.New(ret.Code, ret.ErrMsg))
		return
	}

	ctx.RespEntity(ret.Data)
}

func (s *Service) CreateAssociationInst(ctx *rest.Contexts) {
	request := &metadata.CreateAssociationInstRequest{}
	if err := ctx.DecodeInto(request); err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.New(common.CCErrCommParamsInvalid, err.Error()))
		return
	}

	var ret *metadata.CreateAssociationInstResult
	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, s.EnableTxn, ctx.Kit.Header, func() error {
		var err error
		ret, err = s.Core.AssociationOperation().CreateInst(ctx.Kit, request)
		if err != nil {
			return err
		}

		if ret.Code != 0 {
			return ctx.Kit.CCError.New(ret.Code, ret.ErrMsg)
		}
		return nil
	})

	if txnErr != nil {
		ctx.RespAutoError(txnErr)
		return
	}
	ctx.RespEntity(ret.Data)
}

func (s *Service) DeleteAssociationInst(ctx *rest.Contexts) {
	id, err := strconv.ParseInt(ctx.Request.PathParameter("association_id"), 10, 64)
	if err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommParamsIsInvalid))
		return
	}
	md := new(MetaShell)
	if err := ctx.DecodeInto(md); err != nil {
		ctx.RespAutoError(err)
		return
	}

	var ret *metadata.DeleteAssociationInstResult
	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, s.EnableTxn, ctx.Kit.Header, func() error {
		var err error
		ret, err = s.Core.AssociationOperation().DeleteInst(ctx.Kit, id, md.Metadata)
		if err != nil {
			return err
		}

		if ret.Code != 0 {
			return ctx.Kit.CCError.New(ret.Code, ret.ErrMsg)
		}

		return nil
	})

	if txnErr != nil {
		ctx.RespAutoError(txnErr)
		return
	}
	ctx.RespEntity(ret.Data)
}

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

// SearchCustomMainlineRelation search set module relation below one custom mainline instance
func (s *Service) SearchCustomMainlineRelation(ctx *rest.Contexts) {
	rid := ctx.Kit.Rid

	bizIDStr := ctx.Request.PathParameter(common.BKAppIDField)
	bizID, err := strconv.ParseInt(bizIDStr, 10, 64)
	if nil != err {
		blog.Errorf("SearchCustomMainlineRelation failed, bizIDStr: %s, err: %s, rid: %s", bizIDStr, err.Error(), rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField))
		return
	}
	if bizID == 0 {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField))
	}

	input := metadata.SearchCustomMainlineRelationParam{}
	if err := ctx.DecodeInto(&input); err != nil {
		ctx.RespAutoError(err)
		return
	}

	rawErr := input.Validate()
	if rawErr.ErrCode != 0 {
		ctx.RespAutoError(rawErr.ToCCError(ctx.Kit.CCError))
		return
	}

	// find the mainline topo before set
	topoRoot, err := s.Engine.CoreAPI.CoreService().Mainline().SearchMainlineInstanceTopoBeforeSet(ctx.Kit.Ctx, ctx.Kit.Header, bizID, false)
	if err != nil {
		blog.Errorf("SearchCustomMainlineRelation failed, bizID:%d, SearchMainlineInstanceTopoBeforeSet err:%s, rid:%s", bizID, err.Error(), rid)
		ctx.RespAutoError(err)
		return
	}
	if topoRoot == nil {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrorMainlineInstTopoNotFound, bizID))
		return
	}

	// check objID,instID and get set parent instance IDs
	isCustomMainlineObj := false
	var targetNode *metadata.TopoInstanceNode

	topoRoot.DeepFirstTraversal(func(node *metadata.TopoInstanceNode) {
		if node.ObjectID == input.ObjectID {
			if node.ObjectID != common.BKInnerObjIDApp && node.ObjectID != common.BKInnerObjIDSet &&
				node.ObjectID != common.BKInnerObjIDModule && node.ObjectID != common.BKInnerObjIDHost {
				// input.ObjectID may not be any topo node.ObjectID, program never arrive here, this is why isCustomMainlineObj is set false at beginning
				isCustomMainlineObj = true
			}
			if node.InstanceID == input.InstID {
				targetNode = node
			}
		}
	})

	if !isCustomMainlineObj {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrorNotCustomMainlineModel, input.ObjectID))
		return
	}

	if targetNode == nil {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrorMainlineInstIDNotFound, input.InstID))
		return
	}

	setParentInstIDs := make([]int64, 0)
	targetNode.DeepFirstTraversal(func(node *metadata.TopoInstanceNode) {
		if len(node.Children) == 0 {
			setParentInstIDs = append(setParentInstIDs, node.InstanceID)
		}
	})

	if len(setParentInstIDs) == 0 {
		ctx.RespEntityWithCount(0, []map[string]interface{}{})
		return
	}

	setQuery := &metadata.QueryCondition{
		Fields: []string{common.BKSetIDField},
		Condition: map[string]interface{}{
			common.BKAppIDField: bizID,
			common.BKParentIDField: map[string]interface{}{
				common.BKDBIN: setParentInstIDs,
			},
		},
		Page: metadata.BasePage{
			Sort:  common.BKSetIDField,
			Limit: input.Page.Limit,
			Start: input.Page.Start,
		},
	}
	setResult, err := s.Engine.CoreAPI.CoreService().Instance().ReadInstance(ctx.Kit.Ctx, ctx.Kit.Header, common.BKInnerObjIDSet, setQuery)
	if err != nil {
		blog.Errorf("failed to get set, setQuery: %#v, err: %s, rid: %s", setQuery, err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	if len(setResult.Data.Info) == 0 {
		ctx.RespEntityWithCount(0, []map[string]interface{}{})
		return
	}

	setIDs := make([]int64, len(setResult.Data.Info))
	for i, set := range setResult.Data.Info {
		setID, err := set.Int64(common.BKSetIDField)
		if err != nil {
			blog.Errorf("failed to parse setID into int64, set: %#v, err: %s, rid: %s", set, err.Error(), ctx.Kit.Rid)
			ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommInstFieldConvertFail, common.BKInnerObjIDSet, common.BKSetIDField, "int64", err.Error()))
			return
		}
		setIDs[i] = setID
	}

	moduleQuery := &metadata.QueryCondition{
		Fields: []string{common.BKSetIDField, common.BKModuleIDField},
		Condition: map[string]interface{}{
			common.BKSetIDField: map[string]interface{}{
				common.BKDBIN: setIDs,
			},
		},
		Page: metadata.BasePage{
			Limit: common.BKNoLimit,
		},
	}
	moduleResult, err := s.Engine.CoreAPI.CoreService().Instance().ReadInstance(ctx.Kit.Ctx, ctx.Kit.Header, common.BKInnerObjIDModule, moduleQuery)
	if err != nil {
		blog.Errorf("failed to get module, moduleQuery: %#v, err: %s, rid: %s", moduleQuery, err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	if len(moduleResult.Data.Info) == 0 {
		ctx.RespEntityWithCount(0, []map[string]interface{}{})
		return
	}

	setModuleMap := make(map[int64][]int64)
	for _, module := range moduleResult.Data.Info {
		setID, err := module.Int64(common.BKSetIDField)
		if err != nil {
			blog.Errorf("failed to parse setID into int64, set: %#v, err: %s, rid: %s", module, err.Error(), ctx.Kit.Rid)
			ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommInstFieldConvertFail, common.BKInnerObjIDSet, common.BKSetIDField, "int64", err.Error()))
			return
		}
		moduleID, err := module.Int64(common.BKModuleIDField)
		if err != nil {
			blog.Errorf("failed to parse moduleID into int64, set: %#v, err: %s, rid: %s", module, err.Error(), ctx.Kit.Rid)
			ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommInstFieldConvertFail, common.BKInnerObjIDModule, common.BKModuleIDField, "int64", err.Error()))
			return
		}
		if _, ok := setModuleMap[setID]; !ok {
			setModuleMap[setID] = []int64{}
		}
		setModuleMap[setID] = append(setModuleMap[setID], moduleID)
	}

	result := make([]map[string]interface{}, 0)
	for setID, moduleIDs := range setModuleMap {
		relation := map[string]interface{}{}
		relation["bk_set_id"] = setID
		relation["bk_module_ids"] = moduleIDs
		result = append(result, relation)
	}

	ctx.RespEntityWithCount(int64(len(setIDs)), result)
}
