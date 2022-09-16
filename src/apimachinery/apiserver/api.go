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

package apiserver

import (
	"context"
	"errors"
	"net/http"

	"configcenter/src/apimachinery/rest"
	"configcenter/src/common"
	"configcenter/src/common/condition"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
)

// Client TODO
func (a *apiServer) Client() rest.ClientInterface {
	return a.client
}

// AddDefaultApp TODO
func (a *apiServer) AddDefaultApp(ctx context.Context, h http.Header, ownerID string, params mapstr.MapStr) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "biz/default/%s"

	err = a.client.Post().
		WithContext(ctx).
		Body(params).
		SubResourcef(subPath, ownerID).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

// SearchDefaultApp TODO
func (a *apiServer) SearchDefaultApp(ctx context.Context, h http.Header, ownerID string) (resp *metadata.QueryInstResult, err error) {
	resp = new(metadata.QueryInstResult)
	subPath := "biz/default/%s/search"

	err = a.client.Post().
		WithContext(ctx).
		Body(nil).
		SubResourcef(subPath, ownerID).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

// GetObjectData TODO
func (a *apiServer) GetObjectData(ctx context.Context, h http.Header, params mapstr.MapStr) (resp *metadata.ObjectAttrBatchResult, err error) {
	resp = new(metadata.ObjectAttrBatchResult)
	subPath := "/findmany/object"

	err = a.client.Post().
		WithContext(ctx).
		Body(params).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

// SearchObjectWithTotalInfo search object with it's attribute and association
func (a *apiServer) SearchObjectWithTotalInfo(ctx context.Context, h http.Header, params *metadata.BatchExportObject) (
	*metadata.TotalObjectInfo, error) {

	resp := new(metadata.ListObjectTopoResponse)
	subPath := "/findmany/object/total/info"

	err := a.client.Post().
		WithContext(ctx).
		Body(params).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)

	if err != nil {
		return nil, err
	}

	if ccErr := resp.CCError(); ccErr != nil {
		return nil, ccErr
	}

	return resp.Data, nil
}

// GetInstDetail TODO
func (a *apiServer) GetInstDetail(ctx context.Context, h http.Header, objID string, params mapstr.MapStr) (resp *metadata.QueryInstResult, err error) {

	resp = new(metadata.QueryInstResult)
	subPath := "/find/instance/object/%s"
	err = a.client.Post().
		WithContext(ctx).
		Body(params).
		SubResourcef(subPath, objID).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

// GetInstUniqueFields TODO
func (a *apiServer) GetInstUniqueFields(ctx context.Context, h http.Header, objID string, uniqueID int64,
	params mapstr.MapStr) (resp metadata.QueryUniqueFieldsResult, err error) {

	subPath := "/find/instance/object/%s/unique_fields/by/unique/%d"
	err = a.client.Post().
		WithContext(ctx).
		Body(params).
		SubResourcef(subPath, objID, uniqueID).
		WithHeaders(h).
		Do().
		Into(&resp)
	return
}

// CreateObjectAtt TODO
func (a *apiServer) CreateObjectAtt(ctx context.Context, h http.Header, obj *metadata.ObjAttDes) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "/create/objectattr"

	err = a.client.Post().
		WithContext(ctx).
		Body(obj).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

// UpdateObjectAtt TODO
func (a *apiServer) UpdateObjectAtt(ctx context.Context, objID string, h http.Header, data map[string]interface{}) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "/update/objectattr/%s"
	err = a.client.Put().
		WithContext(ctx).
		Body(data).
		SubResourcef(subPath, objID).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

// DeleteObjectAtt TODO
func (a *apiServer) DeleteObjectAtt(ctx context.Context, objID string, h http.Header) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "/delete/objectattr/%s"

	err = a.client.Delete().
		WithContext(ctx).
		Body(nil).
		SubResourcef(subPath, objID).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

// GetObjectAttr TODO
func (a *apiServer) GetObjectAttr(ctx context.Context, h http.Header, params mapstr.MapStr) (resp *metadata.ObjectAttrResult, err error) {

	resp = new(metadata.ObjectAttrResult)
	subPath := "/find/objectattr"

	err = a.client.Post().
		WithContext(ctx).
		Body(params).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

// GetHostData TODO
func (a *apiServer) GetHostData(ctx context.Context, h http.Header, params mapstr.MapStr) (resp *metadata.QueryInstResult, err error) {

	resp = new(metadata.QueryInstResult)
	subPath := "hosts/search/asstdetail"

	err = a.client.Post().
		WithContext(ctx).
		Body(params).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

// GetObjectGroup TODO
func (a *apiServer) GetObjectGroup(ctx context.Context, h http.Header, ownerID, objID string, params mapstr.MapStr) (resp *metadata.ObjectAttrGroupResult, err error) {

	resp = new(metadata.ObjectAttrGroupResult)
	subPath := "/find/objectattgroup/object/%s"
	err = a.client.Post().
		WithContext(ctx).
		Body(params).
		SubResourcef(subPath, objID).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

// AddHost TODO
func (a *apiServer) AddHost(ctx context.Context, h http.Header, params mapstr.MapStr) (resp *metadata.ResponseDataMapStr, err error) {

	resp = new(metadata.ResponseDataMapStr)
	subPath := "hosts/add"

	err = a.client.Post().
		WithContext(ctx).
		Body(params).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

// AddHostByExcel TODO
func (a *apiServer) AddHostByExcel(ctx context.Context, h http.Header, params mapstr.MapStr) (resp *metadata.ResponseDataMapStr, err error) {

	resp = new(metadata.ResponseDataMapStr)
	subPath := "hosts/excel/add"

	err = a.client.Post().
		WithContext(ctx).
		Body(params).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

// UpdateHost TODO
func (a *apiServer) UpdateHost(ctx context.Context, h http.Header, params mapstr.MapStr) (resp *metadata.ResponseDataMapStr, err error) {

	resp = new(metadata.ResponseDataMapStr)
	subPath := "hosts/update"

	err = a.client.Put().
		WithContext(ctx).
		Body(params).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

// GetHostModuleRelation TODO
func (a *apiServer) GetHostModuleRelation(ctx context.Context, h http.Header, params mapstr.MapStr) (resp *metadata.HostModuleResp, err error) {

	resp = new(metadata.HostModuleResp)
	subPath := "/hosts/modules/read"

	err = a.client.Post().
		WithContext(ctx).
		Body(params).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

// AddInst TODO
func (a *apiServer) AddInst(ctx context.Context, h http.Header, ownerID, objID string, params mapstr.MapStr) (resp *metadata.ResponseDataMapStr, err error) {

	resp = new(metadata.ResponseDataMapStr)
	subPath := "/create/instance/object/%s"
	err = a.client.Post().
		WithContext(ctx).
		Body(params).
		SubResourcef(subPath, objID).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

// AddInstByImport add instances by import excel
func (a *apiServer) AddInstByImport(ctx context.Context, h http.Header, ownerID, objID string, params mapstr.MapStr) (
	*metadata.ResponseDataMapStr, error) {

	resp := new(metadata.ResponseDataMapStr)
	err := a.client.Post().
		WithContext(ctx).
		Body(params).
		SubResourcef("/create/instance/object/%s/by_import", objID).
		WithHeaders(h).
		Do().
		Into(resp)

	if err != nil {
		return nil, err
	}

	if ccErr := resp.CCError(); ccErr != nil {
		return nil, ccErr
	}
	return resp, nil
}

// AddObjectBatch TODO
func (a *apiServer) AddObjectBatch(ctx context.Context, h http.Header, params mapstr.MapStr) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "/createmany/object"

	err = a.client.Post().
		WithContext(ctx).
		Body(params).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

// SearchAssociationInst TODO
func (a *apiServer) SearchAssociationInst(ctx context.Context, h http.Header, request *metadata.SearchAssociationInstRequest) (resp *metadata.SearchAssociationInstResult, err error) {
	resp = new(metadata.SearchAssociationInstResult)
	subPath := "/find/instassociation"

	err = a.client.Post().
		WithContext(ctx).
		Body(request).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)

	return
}

// ImportAssociation TODO
func (a *apiServer) ImportAssociation(ctx context.Context, h http.Header, objID string, input *metadata.RequestImportAssociation) (resp *metadata.ResponeImportAssociation, err error) {
	resp = new(metadata.ResponeImportAssociation)
	subPath := "/import/instassociation/%s"

	err = a.client.Post().
		WithContext(ctx).
		Body(input).
		SubResourcef(subPath, objID).
		WithHeaders(h).
		Do().
		Into(resp)

	return
}

// GetUserAuthorizedBusinessList TODO
func (a *apiServer) GetUserAuthorizedBusinessList(ctx context.Context, h http.Header, user string) (*metadata.InstDataInfo, error) {
	h.Add(common.BKHTTPHeaderUser, user)
	subPath := "/auth/business-list"
	resp := new(metadata.ResponseInstData)

	err := a.client.Get().
		WithContext(ctx).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)

	if err != nil {
		return nil, err
	}

	if !resp.Result {
		return nil, errors.New(resp.ErrMsg)
	}

	return &resp.Data, nil
}

// SearchNetCollectDevice TODO
func (a *apiServer) SearchNetCollectDevice(ctx context.Context, h http.Header, cond condition.Condition) (resp *metadata.ResponseInstData, err error) {
	resp = new(metadata.ResponseInstData)

	subPath := "/collector/netcollect/device/action/search"

	err = a.client.Post().
		WithContext(ctx).
		Body(cond).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)

	return
}

// SearchNetDeviceProperty TODO
func (a *apiServer) SearchNetDeviceProperty(ctx context.Context, h http.Header, cond condition.Condition) (resp *metadata.ResponseInstData, err error) {
	resp = new(metadata.ResponseInstData)

	subPath := "collector/netcollect/property/action/search"

	err = a.client.Post().
		WithContext(ctx).
		Body(cond).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)

	return
}

// SearchNetCollectDeviceBatch TODO
func (a *apiServer) SearchNetCollectDeviceBatch(ctx context.Context, h http.Header, cond mapstr.MapStr) (resp *metadata.ResponseInstData, err error) {
	resp = new(metadata.ResponseInstData)

	subPath := "collector/netcollect/device/action/batch"

	err = a.client.Post().
		WithContext(ctx).
		Body(cond).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)

	return
}

// SearchNetDevicePropertyBatch TODO
func (a *apiServer) SearchNetDevicePropertyBatch(ctx context.Context, h http.Header, cond mapstr.MapStr) (resp *metadata.ResponseInstData, err error) {
	resp = new(metadata.ResponseInstData)

	subPath := "/collector/netcollect/property/action/batch"

	err = a.client.Post().
		WithContext(ctx).
		Body(cond).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)

	return
}

// ListHostWithoutApp TODO
func (a *apiServer) ListHostWithoutApp(ctx context.Context, h http.Header, option metadata.ListHostsWithNoBizParameter) (resp *metadata.ListHostWithoutAppResponse, err error) {
	resp = new(metadata.ListHostWithoutAppResponse)

	subPath := "/hosts/list_hosts_without_app"

	err = a.client.Post().
		WithContext(ctx).
		Body(option).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(&resp)

	return
}

// ReadModuleAssociation get mainline topo model association
func (a *apiServer) ReadModuleAssociation(ctx context.Context, h http.Header,
	cond *metadata.QueryCondition) (resp *metadata.SearchAsstModelResp, err error) {
	resp = new(metadata.SearchAsstModelResp)
	subPath := "/find/instassociation/model"

	err = a.client.Post().
		WithContext(ctx).
		Body(cond).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)

	return
}

// ReadModel read object model data by obj id
func (a *apiServer) ReadModel(ctx context.Context, h http.Header, cond *metadata.QueryCondition) (resp *metadata.
	ReadModelResult, err error) {
	resp = new(metadata.ReadModelResult)
	subPath := "/find/object/model"

	err = a.client.Post().
		WithContext(ctx).
		Body(cond).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)

	return
}

// ReadInstance read instance by obj id and condition
func (a *apiServer) ReadInstance(ctx context.Context, h http.Header, objID string,
	cond *metadata.QueryCondition) (resp *metadata.QueryConditionResult, err error) {
	resp = new(metadata.QueryConditionResult)
	subPath := "/find/instance/%s"

	err = a.client.Post().
		WithContext(ctx).
		Body(cond).
		SubResourcef(subPath, objID).
		WithHeaders(h).
		Do().
		Into(resp)

	return
}

// SearchObjectUnique TODO
func (t *apiServer) SearchObjectUnique(ctx context.Context, objID string, h http.Header) (
	resp *metadata.SearchUniqueResult, err error) {
	resp = new(metadata.SearchUniqueResult)
	subPath := "/find/objectunique/object/%s"

	err = t.client.Post().
		WithContext(ctx).
		Body(nil).
		SubResourcef(subPath, objID).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

// FindAssociationByObjectAssociationID TODO
func (a *apiServer) FindAssociationByObjectAssociationID(ctx context.Context, h http.Header, objID string,
	input metadata.FindAssociationByObjectAssociationIDRequest) (
	resp *metadata.FindAssociationByObjectAssociationIDResponse, err error) {
	resp = new(metadata.FindAssociationByObjectAssociationIDResponse)
	subPath := "/topo/find/object/%s/association/by/bk_obj_asst_id"

	err = a.client.Post().
		WithContext(ctx).
		Body(input).
		SubResourcef(subPath, objID).
		WithHeaders(h).
		Do().
		Into(resp)

	return
}

// SearchObjectAssociation TODO
func (asst *apiServer) SearchObjectAssociation(ctx context.Context, h http.Header,
	request *metadata.SearchAssociationObjectRequest) (resp *metadata.SearchAssociationObjectResult, err error) {
	resp = new(metadata.SearchAssociationObjectResult)
	subPath := "/find/objectassociation"

	err = asst.client.Post().
		WithContext(ctx).
		Body(request).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)

	return
}

// CreateManyObject create many object
func (a *apiServer) CreateManyObject(ctx context.Context, h http.Header, params metadata.ImportObjects) (
	[]metadata.Object, error) {

	resp := struct {
		metadata.BaseResp `json:",inline"`
		Data              []metadata.Object `json:"data"`
	}{}
	subPath := "/createmany/object/by_import"

	err := a.client.Post().
		WithContext(ctx).
		Body(params).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(&resp)

	if err != nil {
		return nil, err
	}

	if ccErr := resp.CCError(); ccErr != nil {
		return nil, ccErr
	}

	return resp.Data, nil
}

// SearchCloudArea TODO
func (a *apiServer) SearchCloudArea(ctx context.Context, h http.Header, params metadata.CloudAreaSearchParam) (
	*metadata.SearchDataResult, error) {

	resp := new(metadata.SearchResp)
	subPath := "/findmany/cloudarea"

	err := a.client.Post().
		WithContext(ctx).
		Body(params).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)

	if err != nil {
		return nil, err
	}

	if ccErr := resp.CCError(); ccErr != nil {
		return nil, ccErr
	}

	return &resp.Data, nil
}
