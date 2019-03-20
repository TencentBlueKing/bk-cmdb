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
	"configcenter/src/auth/meta"
	"configcenter/src/auth/parser"
	"configcenter/src/common"
	"configcenter/src/common/condition"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"context"
	"fmt"
	"net/http"
)

// AuthorizeByObjectID authorize model by id
func (am *AuthManager) AuthorizeByObjectID(ctx context.Context, header http.Header, action meta.Action, objID string) error {
	// auth: check authorization
	// step1 get model by objID
	cond := condition.CreateCondition().Field(common.BKObjIDField).Eq(objID)
	queryCond := &metadata.QueryCondition{Condition: cond.ToMapStr()}
	resp, err := am.clientSet.CoreService().Model().ReadModel(context.Background(), header, queryCond)
	if err != nil {
		message := fmt.Sprintf("get model by id: %s failed, err: %+v", objID, err)
		return am.Err.Errorf(common.CCErrCommAuthorizeFailed, message)
	}
	if len(resp.Data.Info) == 0 {
		message := fmt.Sprintf("get model by id: %s failed, not found", objID)
		return am.Err.Errorf(common.CCErrCommAuthorizeFailed, message)
	}
	if len(resp.Data.Info) > 1 {
		message := fmt.Sprintf("get model by id: %s failed, get multiple model", objID)
		return am.Err.Errorf(common.CCErrCommAuthorizeFailed, message)
	}
	object := resp.Data.Info[0].Spec

	// step2: check authorize
	if err := am.AuthorizeByObject(ctx, header, action, object); err != nil {
		message := fmt.Sprintf("authorize failed, %s", err.Error())
		return am.Err.New(common.CCErrCommAuthorizeFailed, message)
	}
	return nil
}

// AuthorizeObject authorize by object, plz be note this method only overlay model read/update/delete, without create
func (am *AuthManager) AuthorizeByObject(ctx context.Context, header http.Header, action meta.Action, objects ...metadata.Object) error {
	
	// step1: extract business ID from object, business ID from all objects must be identical to one value
	var businessID int64
	for idx, object := range objects {
		bizID, err := object.Metadata.Label.Int64(metadata.LabelBusinessID)
		// we should ignore metadata.LabelBusinessID field not found error
		if err != nil && err != metadata.LabelKeyNotExistError{
			message := fmt.Sprintf("parse biz id from model: %+v failed, err: %+v", object, err)
			return am.Err.New(common.CCErrCommAuthorizeFailed, message)
		}
		if idx > 0 && bizID != businessID {
			message := fmt.Sprintf("authorization failed, get multiple business ID from objects")
			return am.Err.New(common.CCErrCommAuthorizeFailed, message)
		}
		businessID = bizID
	}

	resources := make([]meta.ResourceAttribute, 0)
	for _, object := range objects {
		resource := meta.ResourceAttribute{
			Basic: meta.Basic{
				Action:     action,
				Type:       meta.Model,
				Name:       object.ObjectName,
				InstanceID: object.ID,
			},
			SupplierAccount: util.GetOwnerID(header),
			BusinessID:      businessID,
		}

		resources = append(resources, resource)
	}

	return am.authorize(ctx, header, businessID, resources...)
}

func (am *AuthManager) AuthorizeResourceCreate(ctx context.Context, header http.Header, businessID int64, resourceType meta.ResourceType) error {
	resource := meta.ResourceAttribute{
		Basic: meta.Basic{
			Action:     meta.Create,
			Type:       resourceType,
		},
		SupplierAccount: util.GetOwnerID(header),
		BusinessID:      businessID,
	}

	return am.authorize(ctx, header, businessID, resource)
}

func (am *AuthManager) authorize(ctx context.Context, header http.Header, businessID int64, resources ...meta.ResourceAttribute) error {
	commonInfo, err := parser.ParseCommonInfo(&header)
	if err != nil {
		return fmt.Errorf("authentication failed, parse user info from header failed, %+v", err)
	}
	authAttribute := &meta.AuthAttribute{
		User:       commonInfo.User,
		BusinessID: businessID,
		Resources:  resources,
	}

	decision, err := am.Authorizer.Authorize(ctx, authAttribute)
	if err != nil {
		return fmt.Errorf("authorize failed, err: %+v", err)
	}
	if decision.Authorized == false {
		return fmt.Errorf("authorize failed, reason: %s", decision.Reason)
	}

	return nil
}

func (am *AuthManager) AuthorizeByAttributeID(ctx context.Context, header http.Header, action meta.Action, attID int64) error {
	// auth: check authorization
	queryCondition := &metadata.QueryCondition{
		Condition: condition.CreateCondition().Field(common.BKFieldID).Eq(attID).ToMapStr(),
	}
	response, err := am.clientSet.CoreService().Model().ReadModel(context.Background(), header, queryCondition)
	if nil != err {
		message := fmt.Sprintf("get model attribute by id:%d failed, err: %s", attID, err.Error())
		return am.Err.Errorf(common.CCErrCommAuthorizeFailed, message)
	}
	if !response.Result {
		message := fmt.Sprintf("get model attribute by id:%d failed, err: %s", attID, response.ErrMsg)
		return am.Err.Errorf(common.CCErrCommAuthorizeFailed, message)
	}
	
	// result count must be exactly one
	searchResult := response.Data.Info
	if len(searchResult) == 0 {
		message := fmt.Sprintf("target attribute:%d not found.", attID)
		return am.Err.New(common.CCErrCommAuthorizeFailed, message)
	}
	if len(searchResult) >= 1 {
		message := fmt.Sprintf("target attribute:%d found multiple.", attID)
		return am.Err.New(common.CCErrCommAuthorizeFailed, message)
	}
	
	// check authorize by low level
	objID := searchResult[0].Spec.ObjectID
	return am.AuthorizeByObjectID(ctx, header, action, objID)
}


func (am *AuthManager) AuthorizeByClassification(ctx context.Context, header http.Header, action meta.Action, classes ...*metadata.Classification) error {
	
	// extract business id
	var bizID int64
	for idx, class := range classes {
		businessID, err := class.Metadata.Label.Int64(metadata.LabelBusinessID)
		if err != nil {
			return fmt.Errorf("extract business id from classify: %+v failed, err: %+v", class, err)
		}
		if idx != 0 && bizID != businessID {
			return fmt.Errorf("classes:%+v own to multiple business", classes)
		}
		bizID = businessID
	}
	
	// make auth resources
	resources := make([]meta.ResourceAttribute, 0)
	for _, class := range classes {
		resource := meta.ResourceAttribute{
			Basic: meta.Basic{
				Type:       meta.ModelClassification,
				Name:       class.ClassificationName,
				InstanceID: class.ID,
			},
			SupplierAccount: util.GetOwnerID(header),
			BusinessID:      businessID,
		}
		resources = append(resources, resource)
	}
	return am.authorize(ctx, header, businessID, resources...)
}
