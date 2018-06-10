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

	"configcenter/src/apimachinery/rest"
	"configcenter/src/apimachinery/util"
	"configcenter/src/common/core/cc/api"
	metatype "configcenter/src/common/metadata"
	"configcenter/src/source_controller/api/metadata"
	metadata2 "configcenter/src/source_controller/objectcontroller/objectdata/actions/metadata"
)

type MetaInterface interface {
	SelectClassificationWithObject(ctx context.Context, h util.Headers, dat map[string]interface{}) (resp *api.BKAPIRsp, err error)
	SelectClassifications(ctx context.Context, h util.Headers, dat map[string]interface{}) (resp *metatype.QueryObjectClassificationResult, err error)
	DeleteClassification(ctx context.Context, id int, h util.Headers, dat map[string]interface{}) (resp *api.BKAPIRsp, err error)
	CreateClassification(ctx context.Context, h util.Headers, dat *metatype.Classification) (resp *metatype.CreateObjectClassificationResult, err error)
	UpdateClassification(ctx context.Context, id int, h util.Headers, dat map[string]interface{}) (resp *api.BKAPIRsp, err error)

	SearchTopoGraphics(ctx context.Context, h util.Headers, dat *metadata.TopoGraphics) (resp *api.BKAPIRsp, err error)
	UpdateTopoGraphics(ctx context.Context, h util.Headers, dat []metadata.TopoGraphics) (resp *api.BKAPIRsp, err error)

	CreatePropertyGroup(ctx context.Context, h util.Headers, dat *metatype.Group) (resp *metatype.CreateObjectGroupResult, err error)
	UpdatePropertyGroup(ctx context.Context, h util.Headers, dat *metatype.UpdateGroupCondition) (resp *api.BKAPIRsp, err error)
	DeletePropertyGroup(ctx context.Context, groupID string, h util.Headers) (resp *api.BKAPIRsp, err error)
	UpdatePropertyGroupObjectAtt(ctx context.Context, h util.Headers, dat []metadata2.PropertyGroupObjectAtt) (resp *api.BKAPIRsp, err error)
	DeletePropertyGroupObjectAtt(ctx context.Context, objID string, propertyID string, groupID string, h util.Headers) (resp *api.BKAPIRsp, err error)
	SelectPropertyGroupByObjectID(ctx context.Context, objID string, h util.Headers, dat map[string]interface{}) (resp *api.BKAPIRsp, err error)
	SelectGroup(ctx context.Context, h util.Headers, dat map[string]interface{}) (resp *metatype.QueryObjectGroupResult, err error)

	SelectObjects(ctx context.Context, h util.Headers, data interface{}) (resp *metatype.QueryObjectResult, err error)
	DeleteObject(ctx context.Context, id int, h util.Headers, dat map[string]interface{}) (resp *api.BKAPIRsp, err error)
	CreateObject(ctx context.Context, h util.Headers, dat *metatype.Object) (resp *metatype.CreateObjectResult, err error)
	UpdateObject(ctx context.Context, id int, h util.Headers, dat interface{}) (resp *api.BKAPIRsp, err error)

	SelectObjectAssociations(ctx context.Context, h util.Headers, dat map[string]interface{}) (resp *api.BKAPIRsp, err error)
	DeleteObjectAssociation(ctx context.Context, objID string, h util.Headers, dat map[string]interface{}) (resp *api.BKAPIRsp, err error)
	CreateObjectAssociation(ctx context.Context, h util.Headers, dat *metadata.ObjectAsst) (resp *api.BKAPIRsp, err error)
	UpdateObjectAssociation(ctx context.Context, objID string, h util.Headers, dat map[string]interface{}) (resp *api.BKAPIRsp, err error)
	SelectObjectAttByID(ctx context.Context, objID string, h util.Headers) (resp *api.BKAPIRsp, err error)
	SelectObjectAttWithParams(ctx context.Context, h util.Headers, dat map[string]interface{}) (resp *metatype.QueryObjectAttributeResult, err error)
	DeleteObjectAttByID(ctx context.Context, id int, h util.Headers, dat map[string]interface{}) (resp *api.BKAPIRsp, err error)
	CreateObjectAtt(ctx context.Context, h util.Headers, dat *metatype.Attribute) (resp *metatype.CreateObjectAttributeResult, err error)
	UpdateObjectAttByID(ctx context.Context, id int, h util.Headers, dat map[string]interface{}) (resp *api.BKAPIRsp, err error)
}

func NewmetaInterface(client rest.ClientInterface) MetaInterface {
	return &meta{client: client}
}

type meta struct {
	client rest.ClientInterface
}
