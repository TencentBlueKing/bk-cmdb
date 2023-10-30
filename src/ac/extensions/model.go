/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.,
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the ",License",); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an ",AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package extensions

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"configcenter/src/ac/iam"
	"configcenter/src/ac/meta"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/storage/dal/redis"
)

// collectObjectsByObjectIDs collect business object that belongs to business or global object,
// which both id must in objectIDs
func (am *AuthManager) collectObjectsByObjectIDs(ctx context.Context, header http.Header, businessID int64,
	objectIDs ...string) ([]metadata.Object, error) {

	// unique ids so that we can be aware of invalid id if query result length not equal ids's length
	objectIDs = util.StrArrayUnique(objectIDs)

	// get model by objID
	cond := condition.CreateCondition().Field(common.BKObjIDField).In(objectIDs)
	fCond := cond.ToMapStr()
	util.AddModelBizIDCondition(fCond, businessID)
	fCond.Remove(metadata.BKMetadata)
	queryCond := &metadata.QueryCondition{Condition: fCond}

	resp, err := am.clientSet.CoreService().Model().ReadModel(ctx, header, queryCond)
	if err != nil {
		return nil, fmt.Errorf("get model by id: %+v failed, err: %+v", objectIDs, err)
	}
	if len(resp.Info) == 0 {
		return nil, fmt.Errorf("get model by id: %+v failed, not found", objectIDs)
	}
	if len(resp.Info) != len(objectIDs) {
		return nil, fmt.Errorf("get model by id: %+v failed, get multiple model", objectIDs)
	}

	return resp.Info, nil
}

// MakeResourcesByObjects make object resource with businessID and objects
func (am *AuthManager) MakeResourcesByObjects(ctx context.Context, header http.Header, action meta.Action,
	objects ...metadata.Object) ([]meta.ResourceAttribute, error) {
	// prepare resource layers for authorization
	resources := make([]meta.ResourceAttribute, 0)
	for _, object := range objects {
		// instance
		resource := meta.ResourceAttribute{
			Basic: meta.Basic{
				Action:     action,
				Type:       meta.Model,
				Name:       object.ObjectName,
				InstanceID: object.ID,
			},
			SupplierAccount: util.GetOwnerID(header),
			BusinessID:      0,
		}
		resources = append(resources, resource)
	}

	return resources, nil
}

// AuthorizeByObjectIDs TODO
func (am *AuthManager) AuthorizeByObjectIDs(ctx context.Context, header http.Header, action meta.Action, bizID int64,
	objIDs ...string) error {

	rid := util.ExtractRequestIDFromContext(ctx)

	if !am.Enabled() {
		return nil
	}
	if am.SkipReadAuthorization && (action == meta.Find || action == meta.FindMany) {
		blog.V(4).Infof("skip authorization for reading, objIDs: %+v, rid: %s", objIDs, rid)
		return nil
	}

	if len(objIDs) == 0 {
		return nil
	}

	objects, err := am.collectObjectsByObjectIDs(ctx, header, bizID, objIDs...)
	if err != nil {
		return fmt.Errorf("get objects by objIDs(%+v) failed, err: %v, rid: %s", objIDs, err, rid)
	}

	// make auth resources
	resources, err := am.MakeResourcesByObjects(ctx, header, action, objects...)
	if err != nil {
		return fmt.Errorf("make object resources failed, err: %+v", err)
	}
	return am.batchAuthorize(ctx, header, resources...)
}

// GenObjectBatchNoPermissionResp TODO
func (am *AuthManager) GenObjectBatchNoPermissionResp(ctx context.Context, header http.Header, action meta.Action,
	bizID int64, objIDs []string) (*metadata.BaseResp, error) {

	objects, err := am.collectObjectsByObjectIDs(ctx, header, bizID, objIDs...)
	if err != nil {
		return nil, err
	}

	iamObjects := make([][]metadata.IamResourceInstance, 0)
	for _, object := range objects {
		iamObjects = append(iamObjects, []metadata.IamResourceInstance{{
			Type: string(iam.SysModel),
			ID:   strconv.FormatInt(object.ID, 10),
		}})
	}

	iamAction, err := iam.ConvertResourceAction(meta.Model, action, bizID)
	if err != nil {
		return nil, err
	}

	permission := &metadata.IamPermission{
		SystemID: iam.SystemIDCMDB,
		Actions: []metadata.IamAction{{
			ID: string(iamAction),
			RelatedResourceTypes: []metadata.IamResourceType{{
				SystemID:  iam.SystemIDCMDB,
				Type:      string(iam.SysModel),
				Instances: iamObjects,
			}},
		}},
	}
	resp := metadata.NewNoPermissionResp(permission)
	return &resp, nil
}

// AuthorizeResourceCreate TODO
func (am *AuthManager) AuthorizeResourceCreate(ctx context.Context, header http.Header, businessID int64,
	resourceType meta.ResourceType) error {
	if !am.Enabled() {
		return nil
	}

	resource := meta.ResourceAttribute{
		Basic: meta.Basic{
			Type:   resourceType,
			Action: meta.Create,
		},
		SupplierAccount: util.GetOwnerID(header),
		BusinessID:      businessID,
	}

	return am.batchAuthorize(ctx, header, resource)
}

// CreateObjectOnIAM create object on iam including:
// 1. create iam view
// 2. register object resource creator action to iam
func (am *AuthManager) CreateObjectOnIAM(ctx context.Context, header http.Header, objects []metadata.Object,
	iamInstances []metadata.IamInstanceWithCreator, redisCli redis.Client) error {
	if !am.Enabled() {
		return nil
	}

	rid := util.ExtractRequestIDFromContext(ctx)

	// create iam view
	if err := am.Viewer.CreateView(ctx, header, objects, redisCli, rid); err != nil {
		blog.ErrorJSON("create view failed, objects:%s, err: %s, rid: %s", objects, err, rid)
		return err
	}

	// register object resource creator action to iam
	for _, iamInstance := range iamInstances {
		if _, err := am.Authorizer.RegisterResourceCreatorAction(ctx, header, iamInstance); err != nil {
			blog.ErrorJSON("register created object to iam failed, iam instance:%s, err: %s, rid: %s",
				iamInstance, err, rid)
			return err
		}
	}

	return nil
}

// HasFindModelInstAuth has find model instance auth
func (am *AuthManager) HasFindModelInstAuth(kit *rest.Kit, objIDs []string) (*metadata.BaseResp, bool, error) {
	if !am.Enabled() {
		return nil, true, nil
	}

	if len(objIDs) == 0 {
		return nil, true, nil
	}

	cond := &metadata.QueryCondition{
		Fields: []string{common.BKFieldID, common.BKObjIDField},
		Page:   metadata.BasePage{Limit: common.BKNoLimit},
		Condition: map[string]interface{}{
			common.BKObjIDField: map[string]interface{}{
				common.BKDBIN: objIDs,
			},
		},
	}
	modelResp, err := am.clientSet.CoreService().Model().ReadModel(kit.Ctx, kit.Header, cond)
	if err != nil {
		return nil, false, err
	}

	if len(modelResp.Info) == 0 {
		return nil, false, fmt.Errorf("get model by objID failed, val: %+v", objIDs)
	}

	mainlineModels, err := am.getMainlineModel(kit)
	if err != nil {
		return nil, false, err
	}

	authResources := make([]meta.ResourceAttribute, 0)
	for _, v := range modelResp.Info {
		instanceType, err := am.getInstanceTypeByObject(mainlineModels, v.ObjectID, v.ID)
		if err != nil {
			return nil, false, err
		}

		authResources = append(authResources,
			meta.ResourceAttribute{Basic: meta.Basic{Type: instanceType, Action: meta.Find}})
	}

	authResp, authorized := am.Authorize(kit, authResources...)
	return authResp, authorized, nil
}

// getSkipFindAttrAuthModel 主线模型和内置模型（不包括：交换机、路由器、防火墙、负载均衡）模型属性查看不鉴权
func (am *AuthManager) getSkipFindAttrAuthModel(kit *rest.Kit) (map[string]struct{}, error) {
	models, err := am.getMainlineModel(kit)
	if err != nil {
		return nil, err
	}
	models[common.BKInnerObjIDProc] = struct{}{}
	models[common.BKInnerObjIDPlat] = struct{}{}
	models[common.BKInnerObjIDBizSet] = struct{}{}
	models[common.BKInnerObjIDProject] = struct{}{}
	return models, nil
}

var objTypeMap = map[string]meta.ResourceType{
	common.BKInnerObjIDPlat:    meta.CloudAreaInstance,
	common.BKInnerObjIDHost:    meta.HostInstance,
	common.BKInnerObjIDModule:  meta.ModelModule,
	common.BKInnerObjIDSet:     meta.ModelSet,
	common.BKInnerObjIDApp:     meta.Business,
	common.BKInnerObjIDProc:    meta.Process,
	common.BKInnerObjIDBizSet:  meta.BizSet,
	common.BKInnerObjIDProject: meta.Project,
}

func (am *AuthManager) getInstanceTypeByObject(mainlineModel map[string]struct{}, objID string, id int64) (
	meta.ResourceType, error) {

	if objType, ok := objTypeMap[objID]; ok {
		return objType, nil
	}

	if _, ok := mainlineModel[objID]; ok {
		return meta.MainlineInstance, nil
	}

	return iam.GenCMDBDynamicResType(id), nil
}

func (am *AuthManager) getMainlineModel(kit *rest.Kit) (map[string]struct{}, error) {
	cond := mapstr.MapStr{common.AssociationKindIDField: common.AssociationKindMainline}
	asst, err := am.clientSet.CoreService().Association().ReadModelAssociation(kit.Ctx, kit.Header,
		&metadata.QueryCondition{Condition: cond})
	if err != nil {
		return nil, err
	}

	if len(asst.Info) <= 0 {
		return nil, fmt.Errorf("model association [%+v] not found", cond)
	}

	mainlineModel := make(map[string]struct{})
	for _, mainline := range asst.Info {
		mainlineModel[mainline.AsstObjID] = struct{}{}
		mainlineModel[mainline.ObjectID] = struct{}{}
	}

	return mainlineModel, nil
}

// HasFindModelAuthUseObjID use the objID parameter to determine whether you have permission to find the model
func (am *AuthManager) HasFindModelAuthUseObjID(kit *rest.Kit, objIDs []string) (*metadata.BaseResp, bool, error) {
	if !am.Enabled() {
		return nil, true, nil
	}

	if len(objIDs) == 0 {
		return nil, true, nil
	}

	models, err := am.getSkipFindAttrAuthModel(kit)
	if err != nil {
		return nil, false, err
	}
	finalObjIDs := make([]string, 0)
	for _, obj := range objIDs {
		if _, ok := models[obj]; !ok {
			finalObjIDs = append(finalObjIDs, obj)
		}
	}

	if len(finalObjIDs) == 0 {
		return nil, true, nil
	}

	cond := &metadata.QueryCondition{
		Fields: []string{common.BKFieldID},
		Page:   metadata.BasePage{Limit: common.BKNoLimit},
		Condition: map[string]interface{}{
			common.BKObjIDField: map[string]interface{}{
				common.BKDBIN: finalObjIDs,
			},
		},
	}
	modelResp, err := am.clientSet.CoreService().Model().ReadModel(kit.Ctx, kit.Header, cond)
	if err != nil {
		return nil, false, err
	}
	if len(modelResp.Info) == 0 {
		return nil, true, nil
	}

	ids := make([]int64, len(modelResp.Info))
	for idx, val := range modelResp.Info {
		ids[idx] = val.ID
	}

	return am.HasFindModelAuthUseID(kit, ids)
}

// HasFindModelAuthUseID use the id parameter to determine whether you have permission to find the model
func (am *AuthManager) HasFindModelAuthUseID(kit *rest.Kit, ids []int64) (*metadata.BaseResp, bool, error) {
	if !am.Enabled() {
		return nil, true, nil
	}

	if len(ids) == 0 {
		return nil, true, nil
	}

	authResources := make([]meta.ResourceAttribute, len(ids))
	for idx, val := range ids {
		authResources[idx] = meta.ResourceAttribute{Basic: meta.Basic{InstanceID: val, Type: meta.Model,
			Action: meta.Find}}
	}

	authResp, authorized := am.Authorize(kit, authResources...)
	return authResp, authorized, nil
}
