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
	"context"
	"sort"
	"strconv"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/topo_server/core/types"
)

// CreateMainLineObject create a new object in the main line topo
func (s *Service) CreateMainLineObject(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (output interface{}, retErr error) {
	tx, err := s.Txn.Start(context.Background())
	if err != nil {
		blog.Errorf("create mainline object failed, start transaction failed, err: %v", err)
		return nil, params.Err.Error(common.CCErrObjectDBOpErrno)
	}
	params.Header = tx.TxnInfo().IntoHeader(params.Header)

	mainLineAssociation := &metadata.Association{}
	_, err = mainLineAssociation.Parse(data)
	if nil != err {
		blog.Errorf("[api-asst] failed to parse the data(%#v), error info is %s", data, err.Error())
		return nil, params.Err.Errorf(common.CCErrCommParamsInvalid, "mainline object")
	}
	params.MetaData = &mainLineAssociation.Metadata
	ret, err := s.Core.AssociationOperation().CreateMainlineAssociation(params, mainLineAssociation)
	if err != nil {
		blog.Errorf("create mainline object: %s failed, err: %v", mainLineAssociation.ObjectID, err)
		if txnErr := tx.Abort(context.Background()); txnErr != nil {
			blog.Errorf("create mainline object, but abort transaction[id: %s] failed; %v", tx.TxnInfo().TxnID, txnErr)
		}
		return nil, err
	}
	if txnErr := tx.Commit(context.Background()); txnErr != nil {
		blog.Errorf("create mainline object, but commit transaction[id: %s] failed, err: %v", tx.TxnInfo().TxnID, txnErr)
		return nil, params.Err.Error(common.CCErrTopoMainlineCreatFailed)
	}

	// auth: register mainline object
	if err := s.AuthManager.RegisterMainlineObject(params.Context, params.Header, ret.Object()); err != nil {
		blog.Errorf("create mainline object success, but register mainline model to iam failed, err: %+v", err)
		return ret, params.Err.Error(common.CCErrCommRegistResourceToIAMFailed)
	}

	return ret, nil
}

// DeleteMainLineObject delete a object int the main line topo
func (s *Service) DeleteMainLineObject(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	tx, err := s.Txn.Start(context.Background())
	if err != nil {
		return nil, params.Err.Error(common.CCErrObjectDBOpErrno)
	}
	params.Header = tx.TxnInfo().IntoHeader(params.Header)
	objID := pathParams("bk_obj_id")

	// auth: deregister mainline object
	var bizID int64
	if params.MetaData != nil {
		bizID, err = metadata.BizIDFromMetadata(*params.MetaData)
		if err != nil {
			blog.Errorf("parse business id from request failed, err: %+v", err)
			return nil, params.Err.Error(common.CCErrCommParamsInvalid)
		}
	}
	if err := s.AuthManager.DeregisterMainlineModelByObjectID(params.Context, params.Header, bizID, objID); err != nil {
		blog.Errorf("delete mainline association failed, deregister mainline model failed, err: %+v", err)
		return nil, params.Err.Error(common.CCErrCommUnRegistResourceToIAMFailed)
	}

	err = s.Core.AssociationOperation().DeleteMainlineAssociaton(params, objID)

	if err != nil {
		if txerr := tx.Abort(context.Background()); txerr != nil {
			blog.Errorf("[api-asst] abort transaction failed; %v", err)
			return nil, params.Err.Error(common.CCErrObjectDBOpErrno)
		}
	} else {
		if txerr := tx.Commit(context.Background()); txerr != nil {
			return nil, params.Err.Error(common.CCErrObjectDBOpErrno)
		}
	}
	return nil, err
}

// SearchMainLineOBjectTopo search the main line topo
func (s *Service) SearchMainLineObjectTopo(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {

	bizObj, err := s.Core.ObjectOperation().FindSingleObject(params, common.BKInnerObjIDApp)
	if nil != err {
		blog.Errorf("[api-asst] failed to find the biz object, error info is %s", err.Error())
		return nil, err
	}

	// get biz model related mainline models (mainline relationship model)
	return s.Core.AssociationOperation().SearchMainlineAssociationTopo(params, bizObj)
}

// SearchObjectByClassificationID search the object by classification ID
func (s *Service) SearchObjectByClassificationID(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {

	bizObj, err := s.Core.ObjectOperation().FindSingleObject(params, pathParams("bk_obj_id"))
	if nil != err {
		blog.Errorf("[api-asst] failed to find the biz object, error info is %s", err.Error())
		return nil, err
	}

	return s.Core.AssociationOperation().SearchMainlineAssociationTopo(params, bizObj)
}

// SearchBusinessTopo search the business topo
func (s *Service) SearchBusinessTopo(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {

	paramPath := mapstr.MapStr{}
	paramPath.Set("id", pathParams("bk_biz_id"))
	id, err := paramPath.Int64("id")
	if nil != err {
		blog.Errorf("[api-asst] failed to parse the path params id(%s), error info is %s ", pathParams("app_id"), err.Error())
		return nil, err
	}

	bizObj, err := s.Core.ObjectOperation().FindSingleObject(params, common.BKInnerObjIDApp)
	if nil != err {
		return nil, err
	}

	topoInstRst, err := s.Core.AssociationOperation().SearchMainlineAssociationInstTopo(params, bizObj, id)
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

	return s.Core.AssociationOperation().SearchMainlineAssociationInstTopo(params, obj, instID)
}

func (s *Service) SearchAssociationType(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	request := &metadata.SearchAssociationTypeRequest{}
	if err := data.MarshalJSONInto(request); err != nil {
		return nil, params.Err.New(common.CCErrCommParamsInvalid, err.Error())
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

func (s *Service) SearchObjectAssoWithAssoKindList(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {

	ids := new(metadata.AssociationKindIDs)
	if err := data.MarshalJSONInto(ids); err != nil {
		return nil, params.Err.Error(common.CCErrCommParamsInvalid)
	}

	return s.Core.AssociationOperation().SearchObjectAssoWithAssoKindList(params, ids.AsstIDs)
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
	} else if ret.Code != 0 {
		return nil, params.Err.New(ret.Code, ret.ErrMsg)
	} else {
		return ret.Data, nil

	}
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
