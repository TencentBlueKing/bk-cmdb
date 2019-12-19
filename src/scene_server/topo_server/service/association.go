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
	"sort"
	"strconv"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/mapstruct"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/topo_server/core/types"
)

// CreateMainLineObject create a new model in the main line topo
func (s *Service) CreateMainLineObject(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (output interface{}, retErr error) {

	var txnErr error
	// 判断是否使用事务
	if s.EnableTxn {
		sess, err := s.DB.StartSession()
		if err != nil {
			txnErr = err
			blog.Errorf("StartSession err: %s, rid: %s", err.Error(), params.ReqID)
			return nil, err
		}
		// 获取事务信息，将其存入context中
		txnInfo, err := sess.TxnInfo()
		if err != nil {
			txnErr = err
			blog.Errorf("TxnInfo err: %+v", err)
			return nil, params.Err.Error(common.CCErrObjectDBOpErrno)
		}
		params.Header = txnInfo.IntoHeader(params.Header)
		params.Context = util.TnxIntoContext(params.Context, txnInfo)
		err = sess.StartTransaction(params.Context)
		if err != nil {
			txnErr = err
			blog.Errorf("StartTransaction err: %+v", err)
			return nil, params.Err.Error(common.CCErrObjectDBOpErrno)
		}
		defer func() {
			if txnErr == nil {
				err = sess.CommitTransaction(params.Context)
				if err != nil {
					blog.Errorf("CommitTransaction err: %+v", err)
				}
			} else {
				blog.Errorf("Occur err:%v, begin AbortTransaction", txnErr)
				err = sess.AbortTransaction(params.Context)
				if err != nil {
					blog.Errorf("AbortTransaction err: %+v", err)
				}
			}
			sess.EndSession(params.Context)
		}()
	}

	mainLineAssociation := &metadata.Association{}
	_, err := mainLineAssociation.Parse(data)
	if nil != err {
		txnErr = err
		blog.Errorf("[api-asst] failed to parse the data(%#v), error info is %s, rid: %s", data, err.Error(), params.ReqID)
		return nil, params.Err.Errorf(common.CCErrCommParamsInvalid, "mainline object")
	}
	params.MetaData = &mainLineAssociation.Metadata
	ret, err := s.Core.AssociationOperation().CreateMainlineAssociation(params, mainLineAssociation)
	if err != nil {
		txnErr = err
		blog.Errorf("create mainline object: %s failed, err: %v, rid: %s", mainLineAssociation.ObjectID, err, params.ReqID)
		return nil, err
	}

	// auth: register mainline object
	if err := s.AuthManager.RegisterMainlineObject(params.Context, params.Header, ret.Object()); err != nil {
		txnErr = err
		blog.Errorf("create mainline object success, but register mainline model to iam failed, err: %+v, rid: %s", err, params.ReqID)
		return ret, params.Err.Error(common.CCErrCommRegistResourceToIAMFailed)
	}

	return ret, nil
}

// DeleteMainLineObject delete a object int the main line topo
func (s *Service) DeleteMainLineObject(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {

	var txnErr error
	// 判断是否使用事务
	if s.EnableTxn {
		sess, err := s.DB.StartSession()
		if err != nil {
			txnErr = err
			blog.Errorf("StartSession err: %s, rid: %s", err.Error(), params.ReqID)
			return nil, err
		}
		// 获取事务信息，将其存入context中
		txnInfo, err := sess.TxnInfo()
		if err != nil {
			txnErr = err
			blog.Errorf("TxnInfo err: %+v", err)
			return nil, params.Err.Error(common.CCErrObjectDBOpErrno)
		}
		params.Header = txnInfo.IntoHeader(params.Header)
		params.Context = util.TnxIntoContext(params.Context, txnInfo)
		err = sess.StartTransaction(params.Context)
		if err != nil {
			txnErr = err
			blog.Errorf("StartTransaction err: %+v", err)
			return nil, params.Err.Error(common.CCErrObjectDBOpErrno)
		}
		defer func() {
			if txnErr == nil {
				err = sess.CommitTransaction(params.Context)
				if err != nil {
					blog.Errorf("CommitTransaction err: %+v", err)
				}
			} else {
				blog.Errorf("Occur err:%v, begin AbortTransaction", txnErr)
				err = sess.AbortTransaction(params.Context)
				if err != nil {
					blog.Errorf("AbortTransaction err: %+v", err)
				}
			}
			sess.EndSession(params.Context)
		}()
	}

	objID := pathParams("bk_obj_id")

	var bizID int64
	var err error
	if params.MetaData != nil {
		bizID, err = metadata.BizIDFromMetadata(*params.MetaData)
		if err != nil {
			txnErr = err
			blog.Errorf("parse business id from request failed, err: %+v, rid: %s", err, params.ReqID)
			return nil, params.Err.Error(common.CCErrCommParamsInvalid)
		}
	}

	// auth: collection iam resource before it really be deleted
	iamResources, err := s.AuthManager.MakeResourcesByObjectIDs(params.Context, params.Header, bizID, objID)
	if err != nil {
		txnErr = err
		blog.Errorf("MakeResourcesByObjectIDs failed, err: %+v, rid: %s", err, params.ReqID)
		return nil, params.Err.Error(common.CCErrTopoObjectDeleteFailed)
	}

	if err = s.Core.AssociationOperation().DeleteMainlineAssociation(params, objID); err != nil {
		txnErr = err
		blog.Errorf("DeleteMainlineAssociation failed, err: %+v, rid: %s", err, params.ReqID)
		return nil, params.Err.Error(common.CCErrTopoObjectDeleteFailed)
	}

	// auth: do deregister
	if err := s.AuthManager.Authorize.DeregisterResource(params.Context, iamResources...); err != nil {
		txnErr = err
		blog.Errorf("delete mainline association success, but deregister mainline model failed, err: %+v, rid: %s", err, params.ReqID)
		return nil, params.Err.Error(common.CCErrCommUnRegistResourceToIAMFailed)
	}

	return nil, err
}

// SearchMainLineObjectTopo search the main line topo
func (s *Service) SearchMainLineObjectTopo(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {

	bizObj, err := s.Core.ObjectOperation().FindSingleObject(params, common.BKInnerObjIDApp)
	if nil != err {
		blog.Errorf("[api-asst] failed to find the biz object, error info is %s, rid: %s", err.Error(), params.ReqID)
		return nil, err
	}

	// get biz model related mainline models (mainline relationship model)
	return s.Core.AssociationOperation().SearchMainlineAssociationTopo(params, bizObj)
}

// SearchObjectByClassificationID search the object by classification ID
func (s *Service) SearchObjectByClassificationID(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {

	bizObj, err := s.Core.ObjectOperation().FindSingleObject(params, pathParams("bk_obj_id"))
	if nil != err {
		blog.Errorf("[api-asst] failed to find the biz object, error info is %s, rid: %s", err.Error(), params.ReqID)
		return nil, err
	}

	return s.Core.AssociationOperation().SearchMainlineAssociationTopo(params, bizObj)
}

// SearchBusinessTopoWithStatistics calculate how many service instances on each topo instance node
func (s *Service) SearchBusinessTopoWithStatistics(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	return s.searchBusinessTopo(params, pathParams, queryParams, data, true)
}

func (s *Service) SearchBusinessTopo(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	return s.searchBusinessTopo(params, pathParams, queryParams, data, false)
}

// SearchBusinessTopo search the business topo
func (s *Service) searchBusinessTopo(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr, withStatistics bool) ([]*metadata.TopoInstRst, error) {

	paramPath := mapstr.MapStr{}
	paramPath.Set("id", pathParams("bk_biz_id"))
	id, err := paramPath.Int64("id")
	if nil != err {
		blog.Errorf("[api-asst] failed to parse the path params id(%s), error info is %s , rid: %s", pathParams("app_id"), err.Error(), params.ReqID)
		return nil, err
	}

	bizObj, err := s.Core.ObjectOperation().FindSingleObject(params, common.BKInnerObjIDApp)
	if nil != err {
		return nil, err
	}

	topoInstRst, err := s.Core.AssociationOperation().SearchMainlineAssociationInstTopo(params, bizObj, id, withStatistics)
	if err != nil {
		return nil, err
	}

	// sort before response
	SortTopoInst(topoInstRst)

	return topoInstRst, nil
}

func SortTopoInst(instData []*metadata.TopoInstRst) {
	sort.Slice(instData, func(i, j int) bool {
		return instData[i].InstName < instData[j].InstName
	})
	for idx := range instData {
		SortTopoInst(instData[idx].Child)
	}
}

// SearchMainLineChildInstTopo search the child inst topo by a inst
func (s *Service) SearchMainLineChildInstTopo(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {

	// {obj_id}/{app_id}/{inst_id}
	objID := pathParams("obj_id")
	bizID, err := strconv.ParseInt(pathParams("app_id"), 10, 64)
	if nil != err {
		return nil, params.Err.Errorf(common.CCErrCommParamsIsInvalid, "app_id")
	}

	// get the instance id of this object.
	instID, err := strconv.ParseInt(pathParams("inst_id"), 10, 64)
	if nil != err {
		return nil, params.Err.Errorf(common.CCErrCommParamsIsInvalid, "inst_id")
	}
	_ = bizID

	obj, err := s.Core.ObjectOperation().FindSingleObject(params, objID)
	if nil != err {
		return nil, err
	}

	return s.Core.AssociationOperation().SearchMainlineAssociationInstTopo(params, obj, instID, false)
}

func (s *Service) SearchAssociationType(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	request := &metadata.SearchAssociationTypeRequest{}
	if err := data.MarshalJSONInto(request); err != nil {
		return nil, params.Err.New(common.CCErrCommParamsInvalid, err.Error())
	}
	if request.Condition == nil {
		request.Condition = make(map[string]interface{}, 0)
	}

	ret, err := s.Core.AssociationOperation().SearchType(params, request)
	if err != nil {
		return nil, err
	}

	if ret.Code != 0 {
		return nil, params.Err.New(ret.Code, ret.ErrMsg)
	}

	return ret.Data, nil
}

func (s *Service) SearchObjectAssocWithAssocKindList(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {

	ids := new(metadata.AssociationKindIDs)
	if err := data.MarshalJSONInto(ids); err != nil {
		return nil, params.Err.Error(common.CCErrCommParamsInvalid)
	}

	return s.Core.AssociationOperation().SearchObjectAssocWithAssocKindList(params, ids.AsstIDs)
}

func (s *Service) CreateAssociationType(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	request := &metadata.AssociationKind{}
	if err := data.MarshalJSONInto(request); err != nil {
		return nil, params.Err.New(common.CCErrCommParamsInvalid, err.Error())
	}
	ret, err := s.Core.AssociationOperation().CreateType(params, request)
	if err != nil {
		return nil, err
	}

	if ret.Code != 0 {
		return nil, params.Err.New(ret.Code, ret.ErrMsg)
	}

	return ret.Data, nil
}

func (s *Service) UpdateAssociationType(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	request := &metadata.UpdateAssociationTypeRequest{}
	if err := data.MarshalJSONInto(request); err != nil {
		return nil, params.Err.New(common.CCErrCommParamsInvalid, err.Error())
	}

	asstTypeID, err := strconv.ParseInt(pathParams("id"), 10, 64)
	if err != nil {
		return nil, params.Err.New(common.CCErrCommParamsInvalid, err.Error())
	}

	ret, err := s.Core.AssociationOperation().UpdateType(params, asstTypeID, request)
	if err != nil {
		return nil, err
	}

	if ret.Code != 0 {
		return nil, params.Err.New(ret.Code, ret.ErrMsg)
	}

	return ret.Data, nil
}

func (s *Service) DeleteAssociationType(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	asstTypeID, err := strconv.ParseInt(pathParams("id"), 10, 64)
	if err != nil {
		return nil, params.Err.New(common.CCErrCommParamsInvalid, err.Error())
	}

	ret, err := s.Core.AssociationOperation().DeleteType(params, asstTypeID)
	if err != nil {
		return nil, err
	}

	if ret.Code != 0 {
		return nil, params.Err.New(ret.Code, ret.ErrMsg)
	}

	return ret.Data, nil
}

func (s *Service) SearchAssociationInst(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	request := &metadata.SearchAssociationInstRequest{}
	if err := data.MarshalJSONInto(request); err != nil {
		return nil, params.Err.New(common.CCErrCommParamsInvalid, err.Error())
	}

	ret, err := s.Core.AssociationOperation().SearchInst(params, request)
	if err != nil {
		return nil, err
	}

	if ret.Code != 0 {
		return nil, params.Err.New(ret.Code, ret.ErrMsg)
	}

	return ret.Data, nil
}

func (s *Service) CreateAssociationInst(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	request := &metadata.CreateAssociationInstRequest{}
	if err := data.MarshalJSONInto(request); err != nil {
		return nil, params.Err.New(common.CCErrCommParamsInvalid, err.Error())
	}

	ret, err := s.Core.AssociationOperation().CreateInst(params, request)
	if err != nil {
		return nil, err
	} else if ret.Code != 0 {
		return nil, params.Err.New(ret.Code, ret.ErrMsg)
	} else {
		return ret.Data, nil
	}
}

func (s *Service) DeleteAssociationInst(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {

	id, err := strconv.ParseInt(pathParams("association_id"), 10, 64)
	if err != nil {
		return nil, params.Err.Error(common.CCErrCommParamsIsInvalid)
	}

	ret, err := s.Core.AssociationOperation().DeleteInst(params, id)
	if err != nil {
		return nil, err
	} else if ret.Code != 0 {
		return nil, params.Err.New(ret.Code, ret.ErrMsg)
	} else {
		return ret.Data, nil

	}
}

func (s *Service) SearchTopoPath(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	rid := params.ReqID

	bizIDStr := pathParams(common.BKAppIDField)
	bizID, err := strconv.ParseInt(bizIDStr, 10, 64)
	if nil != err {
		blog.Errorf("SearchTopoPath failed, bizIDStr: %s, err: %s, rid: %s", bizIDStr, err.Error(), rid)
		return nil, params.Err.Errorf(common.CCErrCommParamsInvalid, common.BKAppIDField)
	}

	input := metadata.FindTopoPathRequest{}
	if err := mapstruct.Decode2Struct(data, &input); err != nil {
		blog.ErrorJSON("SearchTopoPath failed, parse request body failed, data: %s, err: %s, rid: %s", data, err.Error(), rid)
		return nil, params.Err.Errorf(common.CCErrCommPostInputParseError)
	}
	if len(input.Nodes) == 0 {
		return nil, params.Err.Errorf(common.CCErrCommHTTPBodyEmpty)
	}

	topoRoot, err := s.Engine.CoreAPI.CoreService().Mainline().SearchMainlineInstanceTopo(params.Context, params.Header, bizID, false)
	if err != nil {
		blog.Errorf("SearchTopoPath failed, SearchMainlineInstanceTopo failed, bizID:%d, err:%s, rid:%s", bizID, err.Error(), rid)
		return nil, err
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

	return result, nil
}
