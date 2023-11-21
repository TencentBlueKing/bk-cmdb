/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2019 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package logics

import (
	"context"
	"fmt"
	"net/http"

	"configcenter/src/apimachinery"
	"configcenter/src/common"
	"configcenter/src/common/metadata"
)

// GetCustomObjects get all custom objects(without inner and mainline objects that authorize separately)
func GetCustomObjects(ctx context.Context, header http.Header, client apimachinery.ClientSetInterface) (
	[]metadata.Object, error) {
	// get mainline objects
	assoCond := &metadata.QueryCondition{
		Condition: map[string]interface{}{common.AssociationKindIDField: common.AssociationKindMainline},
		Fields:    []string{common.BKObjIDField},
	}
	assoRsp, err := client.CoreService().Association().ReadModelAssociation(ctx, header, assoCond)
	if err != nil {
		return nil, fmt.Errorf("get custom models failed, read model association cond:%#v, err: %#v", assoCond, err)
	}

	// get all excluded objectIDs
	excludedObjIDs := []string{
		common.BKInnerObjIDApp, common.BKInnerObjIDSet, common.BKInnerObjIDModule,
		common.BKInnerObjIDHost, common.BKInnerObjIDProc, common.BKInnerObjIDPlat,
	}
	for _, association := range assoRsp.Info {
		if !metadata.IsCommon(association.ObjectID) {
			excludedObjIDs = append(excludedObjIDs, association.ObjectID)
		}
	}

	objCond := &metadata.QueryCondition{
		Fields: []string{common.BKObjIDField, common.BKObjNameField, common.BKFieldID},
		Page:   metadata.BasePage{Limit: common.BKNoLimit},
		Condition: map[string]interface{}{
			common.BKIsPre: false,
			common.BKObjIDField: map[string]interface{}{
				common.BKDBNIN: excludedObjIDs,
			},
		},
	}
	resp, err := client.CoreService().Model().ReadModel(ctx, header, objCond)
	if err != nil {
		return nil, fmt.Errorf("get custom models failed, read model cond:%#v, err: %#v", objCond, err)
	}

	if len(resp.Info) == 0 {
		return []metadata.Object{}, nil
	}

	// 表格字段类型的object不注册到权限中心，这里需要将他们过滤出来
	opt := &metadata.CommonQueryOption{
		Fields: []string{common.BKDestModelField},
		Page:   metadata.BasePage{Limit: common.BKMaxPageSize},
	}
	relRes, err := client.CoreService().ModelQuote().ListModelQuoteRelation(ctx, header, opt)
	if err != nil {
		return nil, err
	}

	if len(relRes.Info) == 0 {
		return resp.Info, nil
	}

	relationObjMap := make(map[string]struct{})
	for _, res := range relRes.Info {
		relationObjMap[res.DestModel] = struct{}{}
	}

	result := make([]metadata.Object, 0)
	for _, object := range resp.Info {
		if _, ok := relationObjMap[object.ObjectID]; ok {
			continue
		}

		result = append(result, object)
	}

	return result, nil
}
