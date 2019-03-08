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

package object

import (
	"context"
	"net/http"

	"configcenter/src/apimachinery/rest"
	"configcenter/src/common/metadata"
)

type ObjectInterface interface {
	CreateModel(ctx context.Context, h http.Header, model *metadata.MainLineObject) (resp *metadata.Response, err error)
	DeleteModel(ctx context.Context, ownerID string, objID string, h http.Header) (resp *metadata.Response, err error)
	SelectModel(ctx context.Context, ownerID string, h http.Header) (resp *metadata.MainlineObjectTopoResult, err error)
	SelectModelByClsID(ctx context.Context, ownerID string, clsID string, objID string, h http.Header) (resp *metadata.Response, err error)
	SelectInst(ctx context.Context, ownerID string, appID string, h http.Header) (resp *metadata.Response, err error)
	SelectInstChild(ctx context.Context, ownerID string, objID string, appID string, instID string, h http.Header) (resp *metadata.Response, err error)
	CreateObjectAtt(ctx context.Context, h http.Header, obj *metadata.ObjAttDes) (resp *metadata.Response, err error)
	SelectObjectAttWithParams(ctx context.Context, h http.Header, data map[string]interface{}) (resp *metadata.Response, err error)
	UpdateObjectAtt(ctx context.Context, objID string, h http.Header, data map[string]interface{}) (resp *metadata.Response, err error)
	DeleteObjectAtt(ctx context.Context, objID string, h http.Header) (resp *metadata.Response, err error)
	CreateClassification(ctx context.Context, h http.Header, obj *metadata.Classification) (resp *metadata.Response, err error)
	SelectClassificationWithObjects(ctx context.Context, ownerID string, h http.Header, data map[string]interface{}) (resp *metadata.Response, err error)
	SelectClassificationWithParams(ctx context.Context, h http.Header, data map[string]interface{}) (resp *metadata.Response, err error)
	UpdateClassification(ctx context.Context, classID string, h http.Header, data map[string]interface{}) (resp *metadata.Response, err error)
	DeleteClassification(ctx context.Context, classID string, h http.Header, data map[string]interface{}) (resp *metadata.Response, err error)
	SelectObjectTopoGraphics(ctx context.Context, scopeType string, scopeID string, h http.Header) (resp *metadata.Response, err error)
	UpdateObjectTopoGraphics(ctx context.Context, scopeType string, scopeID string, h http.Header) (resp *metadata.Response, err error)
	CreatePropertyGroup(ctx context.Context, h http.Header, dat metadata.Group) (resp *metadata.Response, err error)
	UpdatePropertyGroup(ctx context.Context, h http.Header, cond *metadata.PropertyGroupCondition) (resp *metadata.Response, err error)
	DeletePropertyGroup(ctx context.Context, groupID string, h http.Header) (resp *metadata.Response, err error)
	UpdatePropertyGroupObjectAtt(ctx context.Context, h http.Header, data metadata.PropertyGroupObjectAtt) (resp *metadata.Response, err error)
	DeletePropertyGroupObjectAtt(ctx context.Context, ownerID string, objID string, propertyID string, groupID string, h http.Header) (resp *metadata.Response, err error)
	SelectPropertyGroupByObjectID(ctx context.Context, ownerID string, objID string, h http.Header, data map[string]interface{}) (resp *metadata.Response, err error)
	CreateObjectBatch(ctx context.Context, h http.Header, data map[string]interface{}) (resp *metadata.Response, err error)
	SearchObjectBatch(ctx context.Context, h http.Header, data map[string]interface{}) (resp *metadata.Response, err error)
	CreateObject(ctx context.Context, h http.Header, obj metadata.Object) (resp *metadata.Response, err error)
	SelectObjectWithParams(ctx context.Context, h http.Header, data map[string]interface{}) (resp *metadata.Response, err error)
	SelectObjectTopo(ctx context.Context, h http.Header, data map[string]interface{}) (resp *metadata.Response, err error)
	UpdateObject(ctx context.Context, objID string, h http.Header, data map[string]interface{}) (resp *metadata.Response, err error)
	DeleteObject(ctx context.Context, objID string, h http.Header, data map[string]interface{}) (resp *metadata.Response, err error)
}

func NewObjectInterface(client rest.ClientInterface) ObjectInterface {
	return &object{
		client: client,
	}
}

type object struct {
	client rest.ClientInterface
}
