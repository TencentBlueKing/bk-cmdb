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
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/topo_server/core/types"
	"strconv"
	"time"

	"gopkg.in/mgo.v2"
)

// CreateObjectBatch batch to create some objects
func (s *Service) CreateObjectBatch(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	data.Remove(metadata.BKMetadata)
	return s.Core.ObjectOperation().CreateObjectBatch(params, data)
}

// SearchObjectBatch batch to search some objects
func (s *Service) SearchObjectBatch(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {

	data.Remove(metadata.BKMetadata)
	return s.Core.ObjectOperation().FindObjectBatch(params, data)
}

// CreateObject create a new object
func (s *Service) CreateObject(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	rsp, err := s.Core.ObjectOperation().CreateObject(params, false, data)
	if nil != err {
		return nil, err
	}
	//////////
	dail_info := &mgo.DialInfo{
		Addrs:     []string{"127.0.0.1"},
		Direct:    false,
		Timeout:   time.Second * 1,
		PoolLimit: 1024,
	}

	session, err := mgo.DialWithInfo(dail_info)
	if err != nil {
		return nil, err
	}

	defer session.Close()

	session.SetMode(mgo.Monotonic, true)
	c := session.DB("cmdb").C("cc_OperationLog")

	//get meta info
	_meta_data, _ := data.Get("metadata")
	_bk_bi_id := _meta_data.(common.KvMap)["lable"]
	log := common.KvMap{
		"bk_supplier_account": common.BKDefaultSupplierID,
		"bk_biz_id":           _bk_bi_id,
		"ext_key":             "",
		"op_desc":             "create module",
		"op_type":             1,
		"op_target":           "module",
		"content":             data,
		"operator":            params.User,
		"op_from":             "",
		"ext_info":            "",
		"op_time":             time.Now(),
		"inst_id":             nil,
	}
	err = c.Insert(log)

	if err != nil {
		return nil, err
	}

	return rsp.ToMapStr()
}

// SearchObject search some objects by condition
func (s *Service) SearchObject(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	cond := condition.CreateCondition()
	if err := cond.Parse(data); nil != err {
		return nil, err
	}

	return s.Core.ObjectOperation().FindObject(params, cond)
}

// SearchObjectTopo search the object topo
func (s *Service) SearchObjectTopo(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	cond := condition.CreateCondition()
	err := cond.Parse(data)
	if nil != err {
		return nil, params.Err.New(common.CCErrTopoObjectSelectFailed, err.Error())
	}

	return s.Core.ObjectOperation().FindObjectTopo(params, cond)
}

// UpdateObject update the object
func (s *Service) UpdateObject(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	id, err := strconv.ParseInt(pathParams("id"), 10, 64)
	if nil != err {
		blog.Errorf("[api-obj] failed to parse the path params id(%s), error info is %s ", pathParams("id"), err.Error())
		return nil, params.Err.Errorf(common.CCErrCommParamsNeedInt, "object id")
	}
	err = s.Core.ObjectOperation().UpdateObject(params, data, id)
	return nil, err
}

// DeleteObject delete the object
func (s *Service) DeleteObject(params types.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	paramPath := mapstr.MapStr{}
	paramPath.Set("id", pathParams("id"))
	id, err := paramPath.Int64("id")
	if nil != err {
		blog.Errorf("[api-obj] failed to parse the path params id(%s), error info is %s ", pathParams("id"), err.Error())
		return nil, err
	}

	cond := condition.CreateCondition()
	err = s.Core.ObjectOperation().DeleteObject(params, id, cond, true)
	return nil, err
}
