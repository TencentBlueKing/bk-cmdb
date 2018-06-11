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

	"configcenter/src/apimachinery/rest"
	"configcenter/src/apimachinery/util"
	"configcenter/src/common/core/cc/api"
	sencapi "configcenter/src/scene_server/api"
	actobj "configcenter/src/scene_server/topo_server/actions/object"
	obj "configcenter/src/source_controller/api/object"
	"configcenter/src/source_controller/objectcontroller/objectdata/actions/metadata"
)

type ObjectInterface interface {
	// graphic
	SelectObjectTopoGraphics(ctx context.Context, scopeType string, scopeID string, h util.Headers) (resp *api.BKAPIRsp, err error)
	UpdateObjectTopoGraphics(ctx context.Context, scopeType string, scopeID string, h util.Headers) (resp *api.BKAPIRsp, err error)

	// object
	CreateObjectBatch(ctx context.Context, h util.Headers, data map[string]interface{}) (resp *api.BKAPIRsp, err error)
	SearchObjectBatch(ctx context.Context, h util.Headers, data map[string]interface{}) (resp *api.BKAPIRsp, err error)
	CreateObject(ctx context.Context, h util.Headers, obj sencapi.ObjectDes) (resp *api.BKAPIRsp, err error)
	SelectObjectWithParams(ctx context.Context, h util.Headers, data map[string]interface{}) (resp *api.BKAPIRsp, err error)
	SelectObjectTopo(ctx context.Context, h util.Headers, data map[string]interface{}) (resp *api.BKAPIRsp, err error)
	UpdateObject(ctx context.Context, objID string, h util.Headers, data map[string]interface{}) (resp *api.BKAPIRsp, err error)
	DeleteObject(ctx context.Context, objID string, h util.Headers, data map[string]interface{}) (resp *api.BKAPIRsp, err error)

	// attribute
	CreateObjectAtt(ctx context.Context, h util.Headers, obj *obj.ObjAttDes) (resp *api.BKAPIRsp, err error)
	SelectObjectAttWithParams(ctx context.Context, h util.Headers, data map[string]interface{}) (resp *api.BKAPIRsp, err error)
	UpdateObjectAtt(ctx context.Context, objID string, h util.Headers, data map[string]interface{}) (resp *api.BKAPIRsp, err error)
	DeleteObjectAtt(ctx context.Context, objID string, h util.Headers) (resp *api.BKAPIRsp, err error)

	// class
	CreateClassification(ctx context.Context, h util.Headers, obj *sencapi.ObjectClsDes) (resp *api.BKAPIRsp, err error)
	SelectClassificationWithObjects(ctx context.Context, h util.Headers, data map[string]interface{}) (resp *api.BKAPIRsp, err error)
	SelectClassificationWithParams(ctx context.Context, h util.Headers, data map[string]interface{}) (resp *api.BKAPIRsp, err error)
	UpdateClassification(ctx context.Context, classID string, h util.Headers, data map[string]interface{}) (resp *api.BKAPIRsp, err error)
	DeleteClassification(ctx context.Context, classID string, h util.Headers, data map[string]interface{}) (resp *api.BKAPIRsp, err error)

	// action
	CreateModel(ctx context.Context, h util.Headers, model *actobj.MainLineObject) (resp *api.BKAPIRsp, err error)
	DeleteModel(ctx context.Context, objID string, h util.Headers) (resp *api.BKAPIRsp, err error)
	SelectModel(ctx context.Context, h util.Headers) (resp *api.BKAPIRsp, err error)
	SelectModelByClsID(ctx context.Context, clsID string, objID string, h util.Headers) (resp *api.BKAPIRsp, err error)
	SelectInst(ctx context.Context, appID string, h util.Headers) (resp *api.BKAPIRsp, err error)
	SelectInstChild(ctx context.Context, objID string, appID string, instID string, h util.Headers) (resp *api.BKAPIRsp, err error)

	// group
	CreatePropertyGroup(ctx context.Context, h util.Headers, dat obj.ObjAttGroupDes) (resp *api.BKAPIRsp, err error)
	UpdatePropertyGroup(ctx context.Context, h util.Headers, cond *metadata.PropertyGroupCondition) (resp *api.BKAPIRsp, err error)
	DeletePropertyGroup(ctx context.Context, groupID string, h util.Headers) (resp *api.BKAPIRsp, err error)
	UpdatePropertyGroupObjectAtt(ctx context.Context, h util.Headers, data metadata.PropertyGroupObjectAtt) (resp *api.BKAPIRsp, err error)
	DeletePropertyGroupObjectAtt(ctx context.Context, objID string, propertyID string, groupID string, h util.Headers) (resp *api.BKAPIRsp, err error)
	SelectPropertyGroupByObjectID(ctx context.Context, objID string, h util.Headers, data map[string]interface{}) (resp *api.BKAPIRsp, err error)
}

func NewObjectInterface(client rest.ClientInterface) ObjectInterface {
	return &object{
		client: client,
	}
}

type object struct {
	client rest.ClientInterface
}
