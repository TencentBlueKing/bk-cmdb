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
	"net/http"

	"configcenter/src/apimachinery/rest"
	"configcenter/src/common/condition"
	ccErr "configcenter/src/common/errors"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/metric"
)

// Client TODO
func (a *apiServer) Client() rest.ClientInterface {
	return a.client
}

// AddDefaultApp TODO
func (a *apiServer) AddDefaultApp(ctx context.Context, h http.Header, ownerID string,
	params mapstr.MapStr) (resp *metadata.Response, err error) {

	resp = new(metadata.Response)
	subPath := "biz/default/%s"

	err = a.client.Post().
		WithContext(ctx).
		Body(params).
		SubResourcef(subPath, ownerID).
		WithHeaders(h).
		Do().
		IntoCmdbResp(resp)
	return
}

// SearchDefaultApp TODO
func (a *apiServer) SearchDefaultApp(ctx context.Context, h http.Header,
	ownerID string) (resp *metadata.QueryInstResult, err error) {
	resp = new(metadata.QueryInstResult)
	subPath := "biz/default/%s/search"

	err = a.client.Post().
		WithContext(ctx).
		Body(nil).
		SubResourcef(subPath, ownerID).
		WithHeaders(h).
		Do().
		IntoCmdbResp(resp)
	return
}

// GetObjectData get object data
func (a *apiServer) GetObjectData(ctx context.Context, h http.Header,
	cond *metadata.ExportObjectCondition) (resp *metadata.ObjectAttrBatchResult, err error) {
	resp = new(metadata.ObjectAttrBatchResult)
	subPath := "/findmany/object"

	err = a.client.Post().
		WithContext(ctx).
		Body(cond).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		IntoCmdbResp(resp)
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
		IntoCmdbResp(resp)

	if err != nil {
		return nil, err
	}

	if ccErr := resp.CCError(); ccErr != nil {
		return nil, ccErr
	}

	return resp.Data, nil
}

// GetInstDetail TODO
func (a *apiServer) GetInstDetail(ctx context.Context, h http.Header, objID string,
	params mapstr.MapStr) (resp *metadata.QueryInstResult, err error) {

	resp = new(metadata.QueryInstResult)
	subPath := "/find/instance/object/%s"
	err = a.client.Post().
		WithContext(ctx).
		Body(params).
		SubResourcef(subPath, objID).
		WithHeaders(h).
		Do().
		IntoCmdbResp(resp)

	if err != nil {
		return nil, err
	}

	if ccErr := resp.CCError(); ccErr != nil {
		return nil, ccErr
	}

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
		IntoCmdbResp(&resp)
	return
}

// CreateObjectAtt TODO
func (a *apiServer) CreateObjectAtt(ctx context.Context, h http.Header,
	obj *metadata.ObjAttDes) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "/create/objectattr"

	err = a.client.Post().
		WithContext(ctx).
		Body(obj).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		IntoCmdbResp(resp)
	return
}

// UpdateObjectAtt TODO
func (a *apiServer) UpdateObjectAtt(ctx context.Context, objID string, h http.Header,
	data map[string]interface{}) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "/update/objectattr/%s"
	err = a.client.Put().
		WithContext(ctx).
		Body(data).
		SubResourcef(subPath, objID).
		WithHeaders(h).
		Do().
		IntoCmdbResp(resp)
	return
}

// DeleteObjectAtt TODO
func (a *apiServer) DeleteObjectAtt(ctx context.Context, objID string, h http.Header) (resp *metadata.Response,
	err error) {
	resp = new(metadata.Response)
	subPath := "/delete/objectattr/%s"

	err = a.client.Delete().
		WithContext(ctx).
		Body(nil).
		SubResourcef(subPath, objID).
		WithHeaders(h).
		Do().
		IntoCmdbResp(resp)
	return
}

// GetObjectAttr TODO
func (a *apiServer) GetObjectAttr(ctx context.Context, h http.Header,
	params mapstr.MapStr) (resp *metadata.ObjectAttrResult, err error) {

	resp = new(metadata.ObjectAttrResult)
	subPath := "/find/objectattr"

	err = a.client.Post().
		WithContext(ctx).
		Body(params).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		IntoCmdbResp(resp)
	return
}

// GetHostData TODO
func (a *apiServer) GetHostData(ctx context.Context, h http.Header,
	params mapstr.MapStr) (resp *metadata.QueryInstResult, err error) {

	resp = new(metadata.QueryInstResult)
	subPath := "hosts/search/asstdetail"

	err = a.client.Post().
		WithContext(ctx).
		Body(params).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		IntoCmdbResp(resp)
	return
}

// GetObjectGroup TODO
func (a *apiServer) GetObjectGroup(ctx context.Context, h http.Header, ownerID, objID string,
	params mapstr.MapStr) (resp *metadata.ObjectAttrGroupResult, err error) {

	resp = new(metadata.ObjectAttrGroupResult)
	subPath := "/find/objectattgroup/object/%s"
	err = a.client.Post().
		WithContext(ctx).
		Body(params).
		SubResourcef(subPath, objID).
		WithHeaders(h).
		Do().
		IntoCmdbResp(resp)
	return
}

// AddHost TODO
func (a *apiServer) AddHost(ctx context.Context, h http.Header,
	params mapstr.MapStr) (resp *metadata.ResponseDataMapStr, err error) {

	resp = new(metadata.ResponseDataMapStr)
	subPath := "hosts/add"

	err = a.client.Post().
		WithContext(ctx).
		Body(params).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		IntoCmdbResp(resp)
	return
}

// AddHostByExcel TODO
func (a *apiServer) AddHostByExcel(ctx context.Context, h http.Header, params mapstr.MapStr) (*metadata.ImportInstRes,
	error) {

	resp := new(metadata.ImportInstResp)
	subPath := "hosts/excel/add"

	err := a.client.Post().
		WithContext(ctx).
		Body(params).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		IntoCmdbResp(resp)

	if err != nil {
		return nil, err
	}

	if ccErr := resp.CCError(); ccErr != nil {
		return nil, ccErr
	}
	return &resp.Data, nil
}

// UpdateHost TODO
func (a *apiServer) UpdateHost(ctx context.Context, h http.Header, params mapstr.MapStr) (*metadata.ImportInstRes,
	error) {

	resp := new(metadata.ImportInstResp)
	subPath := "hosts/update"

	err := a.client.Put().
		WithContext(ctx).
		Body(params).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		IntoCmdbResp(resp)

	if err != nil {
		return nil, err
	}

	if ccErr := resp.CCError(); ccErr != nil {
		return nil, ccErr
	}
	return &resp.Data, nil
}

// GetHostModuleRelation TODO
func (a *apiServer) GetHostModuleRelation(ctx context.Context, h http.Header,
	params mapstr.MapStr) (resp *metadata.HostModuleResp, err error) {

	resp = new(metadata.HostModuleResp)
	subPath := "/hosts/modules/read"

	err = a.client.Post().
		WithContext(ctx).
		Body(params).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		IntoCmdbResp(resp)
	return
}

// AddInst TODO
func (a *apiServer) AddInst(ctx context.Context, h http.Header, ownerID, objID string,
	params mapstr.MapStr) (resp *metadata.ResponseDataMapStr, err error) {

	resp = new(metadata.ResponseDataMapStr)
	subPath := "/create/instance/object/%s"
	err = a.client.Post().
		WithContext(ctx).
		Body(params).
		SubResourcef(subPath, objID).
		WithHeaders(h).
		Do().
		IntoCmdbResp(resp)
	return
}

// AddInstByImport add instances by import excel
func (a *apiServer) AddInstByImport(ctx context.Context, h http.Header, ownerID, objID string, params mapstr.MapStr) (
	*metadata.ImportInstRes, error) {

	resp := new(metadata.ImportInstResp)
	err := a.client.Post().
		WithContext(ctx).
		Body(params).
		SubResourcef("/create/instance/object/%s/by_import", objID).
		WithHeaders(h).
		Do().
		IntoCmdbResp(resp)

	if err != nil {
		return nil, err
	}

	if ccErr := resp.CCError(); ccErr != nil {
		return nil, ccErr
	}
	return &resp.Data, nil
}

// AddObjectBatch TODO
func (a *apiServer) AddObjectBatch(ctx context.Context, h http.Header, params mapstr.MapStr) (resp *metadata.Response,
	err error) {
	resp = new(metadata.Response)
	subPath := "/createmany/object"

	err = a.client.Post().
		WithContext(ctx).
		Body(params).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		IntoCmdbResp(resp)
	return
}

// SearchAssociationInst TODO
func (a *apiServer) SearchAssociationInst(ctx context.Context, h http.Header,
	request *metadata.SearchAssociationInstRequest) (resp *metadata.SearchAssociationInstResult, err error) {
	resp = new(metadata.SearchAssociationInstResult)
	subPath := "/find/instassociation"

	err = a.client.Post().
		WithContext(ctx).
		Body(request).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		IntoCmdbResp(resp)

	return
}

// ImportAssociation TODO
func (a *apiServer) ImportAssociation(ctx context.Context, h http.Header, objID string,
	input *metadata.RequestImportAssociation) (resp *metadata.ResponeImportAssociation, err error) {
	resp = new(metadata.ResponeImportAssociation)
	subPath := "/import/instassociation/%s"

	err = a.client.Post().
		WithContext(ctx).
		Body(input).
		SubResourcef(subPath, objID).
		WithHeaders(h).
		Do().
		IntoCmdbResp(resp)

	return
}

// SearchNetCollectDevice TODO
func (a *apiServer) SearchNetCollectDevice(ctx context.Context, h http.Header,
	cond condition.Condition) (resp *metadata.ResponseInstData, err error) {
	resp = new(metadata.ResponseInstData)

	subPath := "/collector/netcollect/device/action/search"

	err = a.client.Post().
		WithContext(ctx).
		Body(cond).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		IntoCmdbResp(resp)

	return
}

// SearchNetDeviceProperty TODO
func (a *apiServer) SearchNetDeviceProperty(ctx context.Context, h http.Header,
	cond condition.Condition) (resp *metadata.ResponseInstData, err error) {
	resp = new(metadata.ResponseInstData)

	subPath := "collector/netcollect/property/action/search"

	err = a.client.Post().
		WithContext(ctx).
		Body(cond).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		IntoCmdbResp(resp)

	return
}

// SearchNetCollectDeviceBatch TODO
func (a *apiServer) SearchNetCollectDeviceBatch(ctx context.Context, h http.Header,
	cond mapstr.MapStr) (resp *metadata.ResponseInstData, err error) {
	resp = new(metadata.ResponseInstData)

	subPath := "collector/netcollect/device/action/batch"

	err = a.client.Post().
		WithContext(ctx).
		Body(cond).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		IntoCmdbResp(resp)

	return
}

// SearchNetDevicePropertyBatch TODO
func (a *apiServer) SearchNetDevicePropertyBatch(ctx context.Context, h http.Header,
	cond mapstr.MapStr) (resp *metadata.ResponseInstData, err error) {
	resp = new(metadata.ResponseInstData)

	subPath := "/collector/netcollect/property/action/batch"

	err = a.client.Post().
		WithContext(ctx).
		Body(cond).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		IntoCmdbResp(resp)

	return
}

// ListHostWithoutApp TODO
func (a *apiServer) ListHostWithoutApp(ctx context.Context, h http.Header,
	option metadata.ListHostsWithNoBizParameter) (resp *metadata.ListHostWithoutAppResponse, err error) {
	resp = new(metadata.ListHostWithoutAppResponse)

	subPath := "/hosts/list_hosts_without_app"

	err = a.client.Post().
		WithContext(ctx).
		Body(option).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		IntoCmdbResp(&resp)

	return
}

// ReadModuleAssociation get mainline topo model association
func (a *apiServer) ReadModuleAssociation(ctx context.Context, h http.Header, cond *metadata.QueryCondition) (
	*metadata.AsstResult, ccErr.CCErrorCoder) {

	resp := new(metadata.SearchAsstModelResp)
	subPath := "/find/instassociation/model"

	err := a.client.Post().
		WithContext(ctx).
		Body(cond).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		IntoCmdbResp(resp)

	if err != nil {
		return nil, ccErr.CCHttpError
	}

	if err := resp.CCError(); err != nil {
		return nil, err
	}

	return &resp.Data, nil
}

// ReadModel read object model data by obj id
func (a *apiServer) ReadModel(ctx context.Context, h http.Header, cond *metadata.QueryCondition) (
	*metadata.QueryModelDataResult, ccErr.CCErrorCoder) {

	resp := new(metadata.ReadModelResult)
	subPath := "/find/object/model"

	err := a.client.Post().
		WithContext(ctx).
		Body(cond).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		IntoCmdbResp(resp)

	if err != nil {
		return nil, ccErr.CCHttpError
	}

	if err := resp.CCError(); err != nil {
		return nil, err
	}

	return &resp.Data, nil
}

// ReadModelForUI read object model data by obj id, NOTE: this api is only used for UI, do not authorize
func (a *apiServer) ReadModelForUI(ctx context.Context, h http.Header, cond *metadata.QueryCondition) (
	*metadata.QueryModelDataResult, ccErr.CCErrorCoder) {

	resp := new(metadata.ReadModelResult)
	subPath := "/find/object/model/web"

	err := a.client.Post().
		WithContext(ctx).
		Body(cond).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		IntoCmdbResp(resp)

	if err != nil {
		return nil, ccErr.CCHttpError
	}

	if err := resp.CCError(); err != nil {
		return nil, err
	}

	return &resp.Data, nil
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
		IntoCmdbResp(resp)

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
		IntoCmdbResp(resp)
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
		IntoCmdbResp(resp)

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
		IntoCmdbResp(resp)

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
		IntoCmdbResp(&resp)

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
		IntoCmdbResp(resp)

	if err != nil {
		return nil, err
	}

	if ccErr := resp.CCError(); ccErr != nil {
		return nil, ccErr
	}

	return &resp.Data, nil
}

// SearchPlatformSetting find platform config.
func (a *apiServer) SearchPlatformSetting(ctx context.Context, h http.Header, status string) (
	resp *metadata.PlatformSettingResult, err error) {

	resp = new(metadata.PlatformSettingResult)
	subPath := "/admin/find/system_config/platform_setting/%s"

	err = a.client.Get().
		WithContext(ctx).
		SubResourcef(subPath, status).
		WithHeaders(h).
		Do().
		IntoCmdbResp(resp)

	if ccErr := resp.CCError(); ccErr != nil {
		return nil, ccErr
	}

	return
}

// CountObjectInstances counts object instances num with the input conditions.
func (a *apiServer) CountObjectInstances(ctx context.Context, h http.Header, objID string,
	input *metadata.CommonCountFilter) (*metadata.CommonCountResult, ccErr.CCErrorCoder) {

	resp := new(metadata.CommonCountResp)
	subPath := "/count/instances/object/%s"

	err := a.client.Post().
		WithContext(ctx).
		Body(input).
		SubResourcef(subPath, objID).
		WithHeaders(h).
		Do().
		IntoCmdbResp(resp)

	if err != nil {
		return nil, ccErr.CCHttpError
	}

	if err := resp.CCError(); err != nil {
		return nil, err
	}

	return resp.Data, nil
}

// CountObjInstByFilters counts object instances num by filters.
func (a *apiServer) CountObjInstByFilters(ctx context.Context, h http.Header, objID string,
	filters []map[string]interface{}) ([]int64, ccErr.CCErrorCoder) {

	resp := new(metadata.BatchCountResp)
	subPath := "/count/%s/instances"

	err := a.client.Post().
		WithContext(ctx).
		Body(filters).
		SubResourcef(subPath, objID).
		WithHeaders(h).
		Do().
		IntoCmdbResp(resp)

	if err != nil {
		return nil, ccErr.CCHttpError
	}

	if err := resp.CCError(); err != nil {
		return nil, err
	}

	return resp.Data, nil
}

// GroupRelResByIDs group related resource by ids
func (a *apiServer) GroupRelResByIDs(ctx context.Context, h http.Header, kind metadata.GroupByResKind,
	opt *metadata.GroupRelResByIDsOption) (map[int64][]interface{}, ccErr.CCErrorCoder) {

	resp := new(metadata.GroupRelResByIDsResp)
	subPath := "/group/related/%s/resource/by_ids"

	err := a.client.Post().
		WithContext(ctx).
		Body(opt).
		SubResourcef(subPath, kind).
		WithHeaders(h).
		Do().
		IntoCmdbResp(resp)

	if err != nil {
		return nil, ccErr.CCHttpError
	}

	if err := resp.CCError(); err != nil {
		return nil, err
	}

	return resp.Data, nil
}

// HealthCheck check if api-server is healthy.
func (a *apiServer) HealthCheck() (bool, error) {
	resp := new(metric.HealthResponse)
	err := a.client.Get().
		WithContext(context.Background()).
		WithBaseURL("/").
		SubResourcef("/healthz").
		Body(nil).
		Do().
		IntoCmdbResp(resp)

	if err != nil {
		return false, err
	}

	if !resp.Result {
		return false, ccErr.New(resp.Code, resp.Message)
	}

	return true, nil
}

// SearchProject search project
func (a *apiServer) SearchProject(ctx context.Context, h http.Header, params *metadata.SearchProjectOption) (
	*metadata.InstResult, error) {

	resp := new(metadata.SearchInstResult)
	subPath := "/findmany/project"

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
