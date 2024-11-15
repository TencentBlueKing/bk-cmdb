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

package y3_10_202109131607

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	iamtype "configcenter/src/ac/iam"
	"configcenter/src/common"
	"configcenter/src/common/auth"
	"configcenter/src/common/blog"
	httpheader "configcenter/src/common/http/header"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/resource/esb"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/dal/redis"
	"configcenter/src/storage/driver/mongodb/instancemapping"
)

// migrateIAMSysInstances migrate iam system instances
func migrateIAMSysInstances(ctx context.Context, db dal.RDB, cache redis.Client, iam *iamtype.IAM,
	conf *upgrader.Config) error {
	if !auth.EnableAuthorize() {
		return nil
	}

	// for the first installation, cmdb is not registered to iam,
	// skip migrate iam system instances
	isRegistered, err := iam.IsRegisteredToIAM(ctx)
	if err != nil {
		blog.Errorf("check iam system info failed, err: %v", err)
		return err
	}

	if !isRegistered {
		blog.Warnf("skip migrate iam system instances, for not registered to iam yet")
		return nil
	}

	// get all custom objects(without inner and mainline objects that authorize separately)
	associations := make([]Association, 0)
	filter := mapstr.MapStr{
		common.AssociationKindIDField: common.AssociationKindMainline,
	}

	if err := db.Table(common.BKTableNameObjAsst).Find(filter).Fields(common.BKObjIDField).
		All(ctx, &associations); err != nil {
		blog.Errorf("get mainline associations failed, err: %v", err)
		return err
	}

	excludeObjIDs := []string{
		common.BKInnerObjIDApp, common.BKInnerObjIDSet, common.BKInnerObjIDModule,
		common.BKInnerObjIDHost, common.BKInnerObjIDProc, common.BKInnerObjIDPlat,
	}
	for _, association := range associations {
		if !metadata.IsCommon(association.ObjectID) {
			excludeObjIDs = append(excludeObjIDs, association.ObjectID)
		}
	}

	objects := make([]Object, 0)
	condition := map[string]interface{}{
		common.BKIsPre: false,
		common.BKObjIDField: map[string]interface{}{
			common.BKDBNIN: excludeObjIDs,
		},
	}
	if err := db.Table(common.BKTableNameObjDes).Find(condition).All(ctx, &objects); err != nil {
		blog.Errorf("get all custom objects failed, err: %v", err)
		return err
	}

	// add new system instances
	if err := iam.SyncIAMSysInstances(ctx, cache, convertTenantObject(objects)); err != nil {
		blog.Errorf("sync iam system instances failed, err: %v", err)
		return err
	}

	fields := []iamtype.SystemQueryField{iamtype.FieldActions}
	iamResp, err := iam.Client.GetSystemInfo(ctx, fields)
	if err != nil {
		blog.Errorf("get system info failed, error: %v", err)
		return err
	}
	iamActionMap := make(map[iamtype.ActionID]struct{})
	for _, action := range iamResp.Data.Actions {
		iamActionMap[action.ID] = struct{}{}
	}
	// migrate instance auth policies
	instanceActions := []iamtype.ActionID{"create_sys_instance", "edit_sys_instance", "delete_sys_instance"}
	for _, action := range instanceActions {
		if _, ok := iamActionMap[action]; !ok {
			continue
		}
		if err := migrateModelInstancePermission(ctx, action, db, iam, objects); err != nil {
			blog.Errorf("[upgrade y3.10.202106301057] migrate model instance authorization failed, error: %v", err)
			return err
		}
	}

	// delete the old cmdb resource
	param := &iamtype.DeleteCMDBResourceParam{
		ActionIDs: []iamtype.ActionID{
			"create_sys_instance",
			"edit_sys_instance",
			"delete_sys_instance",
			"create_event_subscription",
			"edit_event_subscription",
			"delete_event_subscription",
			"find_event_subscription",
			"watch_set_template_event",
		},
		InstanceSelectionIDs: []iamtype.InstanceSelectionID{"sys_instance", "sys_event_pushing"},
		TypeIDs:              []iamtype.TypeID{"sys_instance", "sys_event_pushing"},
	}
	return iam.DeleteCMDBResource(ctx, param, convertTenantObject(objects))
}

func migrateModelInstancePermission(ctx context.Context, action iamtype.ActionID, db dal.DB, iam *iamtype.IAM,
	objects []Object) error {

	var timestamp, pageSize int64 = 0, 500

	for page := 1; ; page++ {
		listPoliciesParam := &iamtype.ListPoliciesParams{
			ActionID:  action,
			Page:      int64(page),
			PageSize:  pageSize,
			Timestamp: timestamp,
		}

		listPoliciesResp, err := iam.Client.ListPolicies(ctx, listPoliciesParam)
		if err != nil {
			blog.Errorf("list %s policies failed, page: %d, timestamp: %d, error: %v", action, page, timestamp, err)
			return err
		}

		if listPoliciesResp.Metadata.System != iamtype.SystemIDCMDB {
			blog.Errorf("list %s policies, but system id %s does not match", action, listPoliciesResp.Metadata.System)
			return errors.New("system id does not match")
		}

		timestamp = listPoliciesResp.Metadata.Timestamp

		policyIDs := make([]int64, len(listPoliciesResp.Results))
		for idx, policyRes := range listPoliciesResp.Results {
			if err := migrateModelInstancePolicy(ctx, action, db, policyRes, objects); err != nil {
				blog.ErrorJSON("migrate %s policies %s failed, error: %s", action, policyRes, err)
				return err
			}
			policyIDs[idx] = policyRes.ID
		}

		blog.Infof("successfully migrate policies: %+v", policyIDs)

		if len(listPoliciesResp.Results) < page {
			return nil
		}
	}
}

func migrateModelInstancePolicy(ctx context.Context, action iamtype.ActionID, db dal.DB,
	policyRes iamtype.PolicyResult, objects []Object) error {

	objectIDs := make([]int64, len(objects))
	objMap := make(map[int64]Object)
	for index, object := range objects {
		objectIDs[index] = object.ID
		objMap[object.ID] = object
	}

	parsedPolicy, err := parseInstancePolicy(policyRes.Expression)
	if err != nil {
		blog.ErrorJSON("parse %s policies %s failed, error: %s", action, policyRes.Expression, err)
		return err
	}

	// if user has permission to any instances, then migrate it to all current objects' instance action
	if parsedPolicy.isAny {
		for _, id := range objectIDs {
			err := batchGrantInstanceAuth(ctx, policyRes.Subject, policyRes.ExpiredAt, action, nil, id)
			if err != nil {
				blog.ErrorJSON("batch grant instance auth failed, err: %s, policy: %s, id: %s", err, policyRes, id)
				return err
			}
		}
		return nil
	}

	// migrate objects permissions to the 'any' permission of this object's instance action
	for _, id := range parsedPolicy.objectIDs {
		if _, exists := objMap[id]; !exists {
			blog.Errorf("iam policy has an object(%d) that is not in cc, **skip this object**", id)
			continue
		}
		err := batchGrantInstanceAuth(ctx, policyRes.Subject, policyRes.ExpiredAt, action, nil, id)
		if err != nil {
			blog.ErrorJSON("batch grant instance auth failed, err: %s, policy: %s, id: %s", err, policyRes, id)
			return err
		}
	}

	for objID, instanceIDs := range parsedPolicy.objInstIDMap {
		if _, exists := objMap[objID]; !exists {
			blog.Errorf("iam policy has an object(%d) that is not in cc, **skip this object**", objID)
			continue
		}

		// get objectID by instance IDs to judge if the instances belongs to the object specified
		instIDObjMappings, err := instancemapping.GetInstanceObjectMapping(instanceIDs)
		if err != nil {
			blog.Errorf("get instance object mapping from instance ids(%+v) failed, err: %v", instanceIDs, err)
			return err
		}

		object := objMap[objID]
		objInstIDs := make([]int64, 0)
		for _, objectInfo := range instIDObjMappings {
			// only use the instances that belongs to the object
			if object.ObjectID == objectInfo.ObjectID {
				objInstIDs = append(objInstIDs, objectInfo.ID)
			}
		}

		if len(objInstIDs) == 0 {
			continue
		}

		if err := grantAuthForInstances(ctx, objInstIDs, db, policyRes, action, object); err != nil {
			return err
		}
	}

	return nil
}

func grantAuthForInstances(ctx context.Context, objInstIDs []int64, db dal.DB,
	policyRes iamtype.PolicyResult, action iamtype.ActionID, object Object) error {

	instanceFilter := map[string]interface{}{
		common.BKInstIDField: map[string]interface{}{common.BKDBIN: objInstIDs},
	}

	// get instance names by ids
	instances := make([]iamtype.SimplifiedInstance, 0)
	if err := db.Table(common.GetObjectInstTableName(object.ObjectID, object.OwnerID)).Find(instanceFilter).
		Fields(common.BKInstIDField, common.BKInstNameField).All(ctx, &instances); err != nil {
		blog.Errorf("get instances failed, error: %v, instIDs: %+v", err, objInstIDs)
		return err
	}

	instanceMap := make(map[int64]string)
	for _, instance := range instances {
		instanceMap[instance.InstanceID] = instance.InstanceName
	}

	instIDs := make([]int64, 0)
	for index, instID := range objInstIDs {
		if _, exists := instanceMap[instID]; !exists {
			blog.Errorf("iam policy has an instance(%d) that is not in cc, **skip this object**", instID)
			continue
		}

		instIDs = append(instIDs, instID)

		// iam only allows granting 20 permissions at a time
		if len(instIDs) < 20 && index != len(objInstIDs)-1 {
			continue
		}

		iamInsts := make([]metadata.IamInstance, 0)
		for _, instanceID := range instIDs {
			iamInsts = append(iamInsts, metadata.IamInstance{
				ID:   strconv.FormatInt(instanceID, 10),
				Name: instanceMap[instanceID],
			})
		}

		err := batchGrantInstanceAuth(ctx, policyRes.Subject, policyRes.ExpiredAt, action, iamInsts, object.ID)
		if err != nil {
			blog.ErrorJSON("batch grant instance auth failed, err: %s, policy: %s, instIDs: %s, objID: %s",
				err, policyRes, instIDs, object.ID)
			return err
		}
		instIDs = make([]int64, 0)
	}
	return nil
}

func batchGrantInstanceAuth(ctx context.Context, subject iamtype.PolicySubject, expiredAt int64,
	actionID iamtype.ActionID, instances []metadata.IamInstance, objectID int64) error {

	header := http.Header{}
	httpheader.SetUser(header, common.CCSystemOperatorUserName)
	req := &metadata.IamBatchOperateInstanceAuthReq{
		Asynchronous: false,
		Operate:      metadata.IamGrantOperation,
		System:       iamtype.SystemIDCMDB,
		Actions:      []metadata.ActionWithID{{ID: convertOldInstanceAction(actionID, objectID)}},
		Subject: metadata.IamSubject{
			Type: subject.Type,
			Id:   subject.ID,
		},
		ExpiredAt: expiredAt,
	}

	// create model instance is related to no resources
	if actionID == "create_sys_instance" {
		req.Resources = make([]metadata.IamInstances, 0)
		_, err := esb.EsbClient().IamSrv().BatchOperateInstanceAuth(ctx, header, req)
		if err != nil {
			blog.ErrorJSON("batch register instance auth failed, err: %s, input: %s", err, req)
			return err
		}
		return nil
	}

	req.Resources = []metadata.IamInstances{{
		System:    iamtype.SystemIDCMDB,
		Type:      string(iamtype.GenIAMDynamicResTypeID(objectID)),
		Instances: make([]metadata.IamInstance, 0),
	}}

	// if specified no instances, then grant permissions to all instances in the object
	if len(instances) > 0 {
		req.Resources[0].Instances = instances
	}

	_, err := esb.EsbClient().IamSrv().BatchOperateInstanceAuth(ctx, header, req)
	if err != nil {
		blog.ErrorJSON("batch register instance auth failed, err: %s, input: %s", err, req)
		return err
	}
	return nil
}

// convertOldInstanceAction convert old form of iam instance action to the new form
func convertOldInstanceAction(actionID iamtype.ActionID, id int64) string {
	switch actionID {
	case "create_sys_instance":
		return string(iamtype.GenDynamicActionID(iamtype.Create, id))
	case "edit_sys_instance":
		return string(iamtype.GenDynamicActionID(iamtype.Edit, id))
	case "delete_sys_instance":
		return string(iamtype.GenDynamicActionID(iamtype.Delete, id))
	}
	return ""
}

// Object object metadata definition
type Object struct {
	ID         int64  `field:"id" json:"id" bson:"id" mapstructure:"id"`
	ObjCls     string `field:"bk_classification_id" json:"bk_classification_id" bson:"bk_classification_id" mapstructure:"bk_classification_id"`
	ObjIcon    string `field:"bk_obj_icon" json:"bk_obj_icon" bson:"bk_obj_icon" mapstructure:"bk_obj_icon"`
	ObjectID   string `field:"bk_obj_id" json:"bk_obj_id" bson:"bk_obj_id" mapstructure:"bk_obj_id"`
	ObjectName string `field:"bk_obj_name" json:"bk_obj_name" bson:"bk_obj_name" mapstructure:"bk_obj_name"`

	// IsHidden front-end don't display the object if IsHidden is true
	IsHidden bool `field:"bk_ishidden" json:"bk_ishidden" bson:"bk_ishidden" mapstructure:"bk_ishidden"`

	IsPre         bool           `field:"ispre" json:"ispre" bson:"ispre" mapstructure:"ispre"`
	IsPaused      bool           `field:"bk_ispaused" json:"bk_ispaused" bson:"bk_ispaused" mapstructure:"bk_ispaused"`
	Position      string         `field:"position" json:"position" bson:"position" mapstructure:"position"`
	OwnerID       string         `field:"bk_supplier_account" json:"bk_supplier_account" bson:"bk_supplier_account" mapstructure:"bk_supplier_account"`
	Description   string         `field:"description" json:"description" bson:"description" mapstructure:"description"`
	Creator       string         `field:"creator" json:"creator" bson:"creator" mapstructure:"creator"`
	Modifier      string         `field:"modifier" json:"modifier" bson:"modifier" mapstructure:"modifier"`
	CreateTime    *metadata.Time `field:"create_time" json:"create_time" bson:"create_time" mapstructure:"create_time"`
	LastTime      *metadata.Time `field:"last_time" json:"last_time" bson:"last_time" mapstructure:"last_time"`
	ObjSortNumber int64          `field:"obj_sort_number" json:"obj_sort_number" bson:"obj_sort_number" mapstructure:"obj_sort_number"`
}

// Association defines the association between two objects.
type Association struct {
	ID      int64  `field:"id" json:"id" bson:"id"`
	OwnerID string `field:"bk_supplier_account" json:"bk_supplier_account" bson:"bk_supplier_account"`

	// the unique id belongs to  this association, should be generated with rules as follows:
	// "$ObjectID"_"$AsstID"_"$AsstObjID"
	AssociationName string `field:"bk_obj_asst_id" json:"bk_obj_asst_id" bson:"bk_obj_asst_id"`
	// the alias name of this association, which is a substitute name in the association kind $AsstKindID
	AssociationAliasName string `field:"bk_obj_asst_name" json:"bk_obj_asst_name" bson:"bk_obj_asst_name"`

	// describe which object this association is defined for.
	ObjectID string `field:"bk_obj_id" json:"bk_obj_id" bson:"bk_obj_id"`
	// describe where the Object associate with.
	AsstObjID string `field:"bk_asst_obj_id" json:"bk_asst_obj_id" bson:"bk_asst_obj_id"`
	// the association kind used by this association.
	AsstKindID string `field:"bk_asst_id" json:"bk_asst_id" bson:"bk_asst_id"`

	// defined which kind of association can be used between the source object and destination object.
	Mapping AssociationMapping `field:"mapping" json:"mapping" bson:"mapping"`
	// describe the action when this association is deleted.
	OnDelete AssociationOnDeleteAction `field:"on_delete" json:"on_delete" bson:"on_delete"`
	// describe whether this association is a pre-defined association or not,
	// if true, it means this association is used by cmdb itself.
	IsPre *bool `field:"ispre" json:"ispre" bson:"ispre"`
}

// AssociationOnDeleteAction TODO
type AssociationOnDeleteAction string

// AssociationMapping TODO
type AssociationMapping string

func convertTenantObject(objs []Object) []metadata.Object {
	objects := make([]metadata.Object, 0)
	for _, obj := range objs {
		metaObj := metadata.Object{
			ID:            obj.ID,
			ObjCls:        obj.ObjCls,
			ObjIcon:       obj.ObjIcon,
			ObjectID:      obj.ObjectID,
			ObjectName:    obj.ObjectName,
			IsHidden:      obj.IsHidden,
			IsPre:         obj.IsPre,
			IsPaused:      obj.IsPaused,
			Position:      obj.Position,
			TenantID:      obj.OwnerID,
			Description:   obj.Description,
			Creator:       obj.Creator,
			Modifier:      obj.Modifier,
			CreateTime:    obj.CreateTime,
			LastTime:      obj.LastTime,
			ObjSortNumber: obj.ObjSortNumber,
		}
		objects = append(objects, metaObj)
	}
	return objects
}
