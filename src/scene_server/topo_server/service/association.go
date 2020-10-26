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
	"context"
	"io/ioutil"
	"sort"
	"strconv"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/topo_server/core/types"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

// CreateMainLineObject create a new model in the main line topo
func (s *Service) CreateMainLineObject(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (output interface{}, retErr error) {
	tx, err := s.Txn.Start(context.Background())
	if err != nil {
		blog.Errorf("create mainline model failed, start transaction failed, err: %v, rid: %s", err, params.ReqID)
		return nil, params.Err.Error(common.CCErrObjectDBOpErrno)
	}
	params.Header = tx.TxnInfo().IntoHeader(params.Header)

	mainLineAssociation := &metadata.Association{}
	_, err = mainLineAssociation.Parse(data)
	if nil != err {
		blog.Errorf("[api-asst] failed to parse the data(%#v), error info is %s, rid: %s", data, err.Error(), params.ReqID)
		return nil, params.Err.Errorf(common.CCErrCommParamsInvalid, "mainline object")
	}
	params.MetaData = &mainLineAssociation.Metadata
	ret, err := s.Core.AssociationOperation().CreateMainlineAssociation(params, mainLineAssociation)
	if err != nil {
		blog.Errorf("create mainline object: %s failed, err: %v, rid: %s", mainLineAssociation.ObjectID, err, params.ReqID)
		if txnErr := tx.Abort(context.Background()); txnErr != nil {
			blog.Errorf("create mainline object, but abort transaction[id: %s] failed; %v, rid: %s", tx.TxnInfo().TxnID, txnErr, params.ReqID)
		}
		return nil, err
	}
	if txnErr := tx.Commit(context.Background()); txnErr != nil {
		blog.Errorf("create mainline object, but commit transaction[id: %s] failed, err: %v, rid: %s", tx.TxnInfo().TxnID, txnErr, params.ReqID)
		return nil, params.Err.Error(common.CCErrTopoMainlineCreatFailed)
	}

	// auth: register mainline object
	if err := s.AuthManager.RegisterMainlineObject(params.Context, params.Header, ret.Object()); err != nil {
		blog.Errorf("create mainline object success, but register mainline model to iam failed, err: %+v, rid: %s", err, params.ReqID)
		return ret, params.Err.Error(common.CCErrCommRegistResourceToIAMFailed)
	}

	return ret, nil
}

// DeleteMainLineObject delete a object int the main line topo
func (s *Service) DeleteMainLineObject(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	tx, err := s.Txn.Start(params.Context)
	if err != nil {
		return nil, params.Err.Error(common.CCErrObjectDBOpErrno)
	}
	params.Header = tx.TxnInfo().IntoHeader(params.Header)
	objID := pathParams("bk_obj_id")

	var bizID int64
	if params.MetaData != nil {
		bizID, err = metadata.BizIDFromMetadata(*params.MetaData)
		if err != nil {
			blog.Errorf("parse business id from request failed, err: %+v, rid: %s", err, params.ReqID)
			return nil, params.Err.Error(common.CCErrCommParamsInvalid)
		}
	}

	// auth: collection iam resource before it really be deleted
	iamResources, err := s.AuthManager.MakeResourcesByObjectIDs(params.Context, params.Header, bizID, objID)
	if err != nil {
		blog.Errorf("parse business id from request failed, err: %+v, rid: %s", err, params.ReqID)
		return nil, params.Err.Error(common.CCErrTopoObjectDeleteFailed)
	}

	if err = s.Core.AssociationOperation().DeleteMainlineAssociation(params, objID); err != nil {
		if txErr := tx.Abort(context.Background()); txErr != nil {
			blog.Errorf("[api-asst] abort transaction failed; %v, rid: %s", err, params.ReqID)
			return nil, params.Err.Error(common.CCErrObjectDBOpErrno)
		}
	} else {
		if txErr := tx.Commit(context.Background()); txErr != nil {
			return nil, params.Err.Error(common.CCErrObjectDBOpErrno)
		}
	}

	// auth: do deregister
	if err := s.AuthManager.Authorize.DeregisterResource(params.Context, iamResources...); err != nil {
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

	id, err := strconv.ParseInt(pathParams("bk_biz_id"), 10, 64)
	if nil != err {
		blog.Errorf("[api-asst] failed to parse the path params id(%s), error info is %s , rid: %s", pathParams("bk_biz_id"), err.Error(), params.ReqID)
		return nil, err
	}

	withDefault := false
	if len(queryParams("with_default")) > 0 {
		withDefault = true
	}

	topoInstRst, err := s.Core.AssociationOperation().SearchMainlineAssociationInstTopo(params, common.BKInnerObjIDApp, id, withStatistics, withDefault)
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

	return s.Core.AssociationOperation().SearchMainlineAssociationInstTopo(params, objID, instID, false, false)
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
	params.Header.Add(common.ReadPreferencePolicyKey, common.SecondaryPreference)
	ret, err := s.Core.AssociationOperation().SearchInst(params, request)
	if err != nil {
		return nil, err
	}

	if ret.Code != 0 {
		return nil, params.Err.New(ret.Code, ret.ErrMsg)
	}

	return ret.Data, nil
}

//Search all associations of certain model instance,by regarding the instance as both Association source and Association target.
func (s *Service) SearchAssociationRelatedInst(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	request := &metadata.SearchAssociationRelatedInstRequest{}
	if err := data.MarshalJSONInto(request); err != nil {
		return nil, params.Err.New(common.CCErrCommParamsInvalid, err.Error())
	}
	//check condition
	if request.Condition.InstID == 0 || request.Condition.ObjectID == "" {
		return nil, params.Err.Error(common.CCErrCommHTTPInputInvalid)
	}
	//check fields,if there's any incorrect params,return err.
	if len(request.Fields) == 0 {
		return nil, params.Err.New(common.CCErrCommParamsInvalid, "there should be at least one param in 'fields'.")
	}
	//Use fixed sort parameters
	request.SortArr = []metadata.SearchSort{
		{
			IsDsc: false,
			Field: common.BKFieldID,
		},
	}
	//check Maximum limit
	if request.Limit.Limit > 500 {
		return nil, params.Err.New(common.CCErrCommParamsInvalid, "The maximum limit should be less than 500")
	}

	ret, err := s.Core.AssociationOperation().SearchInstAssociationRelated(params, request)
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
	} else if err = ret.CCError(); err != nil {
		return nil, params.Err.New(ret.Code, ret.ErrMsg)
	} else {
		return ret.Data, nil

	}
}

//Delete association batch by ID.
func (s *Service) DeleteAssociationInstBatch(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	var err error
	var count int
	request := &metadata.DeleteAssociationInstBatchRequest{}
	ret := &metadata.DeleteAssociationInstResult{}

	if err = data.MarshalJSONInto(request); err != nil {
		return nil, params.Err.New(common.CCErrCommParamsInvalid, err.Error())
	}
	if len(request.ID) == 0 {
		return nil, params.Err.Error(common.CCErrCommHTTPInputInvalid)
	}
	//check Maximum limit
	if len(request.ID) > 500 {
		return nil, params.Err.Errorf(common.CCErrCommPageLimitIsExceeded, "The number of ID should be less than 500")
	}
	for _, id := range request.ID {
		ret, err = s.Core.AssociationOperation().DeleteInst(params, id)
		if err != nil {
			return nil, err
		} else if err = ret.CCError(); err != nil {
			return nil, err
		} else {
			count++
		}
	}
	if ret.Data == "" {
		ret.Data = strconv.Itoa(count)
	}
	return ret.Data, nil
}
