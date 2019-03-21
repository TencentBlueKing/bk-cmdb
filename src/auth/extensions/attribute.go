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

	"configcenter/src/auth/meta"
	"configcenter/src/common"
	"configcenter/src/common/condition"
	"configcenter/src/common/metadata"
)

func (am *AuthManager) collectAttributesByAttributeIDs(ctx context.Context, header http.Header, attributeIDs ...int64) ([]metadata.Attribute, error) {
	// get model by objID
	cond := condition.CreateCondition().Field(common.BKFieldID).In(attributeIDs)
	queryCond := &metadata.QueryCondition{Condition: cond.ToMapStr()}
	resp, err := am.clientSet.CoreService().Instance().ReadInstance(ctx, header, common.BKTableNameObjAttDes, queryCond)
	if err != nil {
		return nil, fmt.Errorf("get attribute by id: %+v failed, err: %+v", attributeIDs, err)
	}
	if len(resp.Data.Info) == 0 {
		return nil, fmt.Errorf("get attribute by id: %+v failed, not found", attributeIDs)
	}
	if len(resp.Data.Info) != len(attributeIDs) {
		return nil, fmt.Errorf("get attribute by id: %+v failed, get %d, expect %d", attributeIDs, len(resp.Data.Info), len(attributeIDs))
	}

	attributes := make([]metadata.Attribute, 0)
	for _, item := range resp.Data.Info {
		attribute := metadata.Attribute{}
		attribute.Parse(item)
		attributes = append(attributes, attribute)
	}
	return attributes, nil
}

func (am *AuthManager) AuthorizeByAttributeID(ctx context.Context, header http.Header, action meta.Action, attributeIDs ...int64) error {
	attributes, err := am.collectAttributesByAttributeIDs(ctx, header, attributeIDs...)
	if err != nil {
		return fmt.Errorf("get attributes by id failed, err: %+v", err)
	}
	
	objectIDs := make([]string, 0)
	for _, attribute := range attributes {
		objectIDs = append(objectIDs, attribute.ObjectID)
	}
	
	return am.AuthorizeByObjectID(ctx, header, action, objectIDs...)
}
