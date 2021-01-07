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

package meta

import (
	"context"
	"net/http"

	"configcenter/src/apimachinery/rest"
	"configcenter/src/common/metadata"
)

type MetaInterface interface {
	SelectClassificationWithObject(ctx context.Context, ownerID string, h http.Header, dat map[string]interface{}) (resp *metadata.QueryObjectClassificationWithObjectsResult, err error)
	SelectClassifications(ctx context.Context, h http.Header, dat map[string]interface{}) (resp *metadata.QueryObjectClassificationResult, err error)
	DeleteClassification(ctx context.Context, id int64, h http.Header, dat map[string]interface{}) (resp *metadata.DeleteResult, err error)
	CreateClassification(ctx context.Context, h http.Header, dat *metadata.Classification) (resp *metadata.CreateObjectClassificationResult, err error)
	UpdateClassification(ctx context.Context, id int64, h http.Header, dat map[string]interface{}) (resp *metadata.UpdateResult, err error)
	SearchTopoGraphics(ctx context.Context, h http.Header, dat *metadata.TopoGraphics) (resp *metadata.SearchTopoGraphicsResult, err error)
	UpdateTopoGraphics(ctx context.Context, h http.Header, dat []metadata.TopoGraphics) (resp *metadata.UpdateResult, err error)
	CreatePropertyGroup(ctx context.Context, h http.Header, dat *metadata.Group) (resp *metadata.CreateObjectGroupResult, err error)
	UpdatePropertyGroup(ctx context.Context, h http.Header, dat *metadata.UpdateGroupCondition) (resp *metadata.UpdateResult, err error)
	DeletePropertyGroup(ctx context.Context, groupID string, h http.Header) (resp *metadata.DeleteResult, err error)
	UpdatePropertyGroupObjectAtt(ctx context.Context, h http.Header, dat []metadata.PropertyGroupObjectAtt) (resp *metadata.UpdateResult, err error)
	DeletePropertyGroupObjectAtt(ctx context.Context, ownerID, objID, propertyID, groupID string, h http.Header) (resp *metadata.DeleteResult, err error)
	SelectPropertyGroupByObjectID(ctx context.Context, ownerID string, objID string, h http.Header, dat map[string]interface{}) (resp *metadata.QueryObjectGroupResult, err error)
	SelectGroup(ctx context.Context, h http.Header, dat map[string]interface{}) (resp *metadata.QueryObjectGroupResult, err error)
	SelectObjects(ctx context.Context, h http.Header, data interface{}) (resp *metadata.QueryObjectResult, err error)
	DeleteObject(ctx context.Context, objID int64, h http.Header, dat map[string]interface{}) (resp *metadata.DeleteResult, err error)
	CreateObject(ctx context.Context, h http.Header, dat *metadata.Object) (resp *metadata.CreateObjectResult, err error)
	UpdateObject(ctx context.Context, objID int64, h http.Header, dat map[string]interface{}) (resp *metadata.UpdateResult, err error)
	SelectObjectAssociations(ctx context.Context, h http.Header, dat map[string]interface{}) (resp *metadata.QueryObjectAssociationResult, err error)
	DeleteObjectAssociation(ctx context.Context, objID int64, h http.Header, dat map[string]interface{}) (resp *metadata.DeleteResult, err error)
	CreateObjectAssociation(ctx context.Context, h http.Header, dat *metadata.Association) (resp *metadata.CreateResult, err error)
	UpdateObjectAssociation(ctx context.Context, objID int64, h http.Header, dat map[string]interface{}) (resp *metadata.UpdateResult, err error)
	SelectObjectAttByID(ctx context.Context, attID int64, h http.Header) (resp *metadata.QueryObjectAttributeResult, err error)
	SelectObjectAttWithParams(ctx context.Context, h http.Header, dat map[string]interface{}) (resp *metadata.QueryObjectAttributeResult, err error)
	DeleteObjectAttByID(ctx context.Context, attID int64, h http.Header, dat map[string]interface{}) (resp *metadata.DeleteResult, err error)
	CreateObjectAtt(ctx context.Context, h http.Header, dat *metadata.Attribute) (resp *metadata.CreateObjectAttributeResult, err error)
	UpdateObjectAttByID(ctx context.Context, attID int64, h http.Header, dat map[string]interface{}) (resp *metadata.UpdateResult, err error)
}

func NewmetaInterface(client rest.ClientInterface) MetaInterface {
	return &meta{client: client}
}

type meta struct {
	client rest.ClientInterface
}
