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

package operation

import (
	"configcenter/src/ac/extensions"
	"configcenter/src/apimachinery"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/topo_server/core/inst"
	"configcenter/src/scene_server/topo_server/core/model"
)

// ClassificationOperationInterface classification operation methods
type TTOperationInterface interface {
	SetProxy(modelFactory model.Factory, instFactory inst.Factory, asst AssociationOperationInterface, obj ObjectOperationInterface)

	// FindSingleClassification(kit *rest.Kit, classificationID string) (model.Classification, error)
	//CreateClassification(kit *rest.Kit, data mapstr.MapStr) (model.Classification, error)
	//DeleteClassification(kit *rest.Kit, id int64, cond condition.Condition) error
	//FindClassification(kit *rest.Kit, cond condition.Condition) ([]model.Classification, error)
	FindTTWithObjects(kit *rest.Kit, cond condition.Condition) ([]metadata.ClassificationWithObject, error)
	//UpdateClassification(kit *rest.Kit, data mapstr.MapStr, id int64, cond condition.Condition) error
}

// NewClassificationOperation create a new classification operation instance
func NewTTOperation(client apimachinery.ClientSetInterface, authManager *extensions.AuthManager) TTOperationInterface {
	//return &tt{
	//	clientSet:   client,
	//	authManager: authManager,
	//}
	return &tt{
		clientSet:   client,
		authManager: authManager,
	}
}

type tt struct {
	clientSet    apimachinery.ClientSetInterface
	authManager  *extensions.AuthManager
	asst         AssociationOperationInterface
	obj          ObjectOperationInterface
	modelFactory model.Factory
	instFactory  inst.Factory
}

func (c *tt) SetProxy(modelFactory model.Factory, instFactory inst.Factory, asst AssociationOperationInterface, obj ObjectOperationInterface) {
	c.modelFactory = modelFactory
	c.instFactory = instFactory
	c.asst = asst
	c.obj = obj
}

func (c *tt) FindTTWithObjects(kit *rest.Kit, cond condition.Condition) ([]metadata.ClassificationWithObject, error) {
	fCond := cond.ToMapStr()
	//xxx add tt     src/apimachinery/coreservice/model/model.go
	// xxx add tt   src/apimachinery/coreservice/model/api.go
	rsp, err := c.clientSet.CoreService().Model().ReadTT(kit.Ctx, kit.Header, &metadata.QueryCondition{Condition: fCond})
	if nil != err {
		blog.Errorf("[operation-cls]failed to request the object controller, error info is %s, rid: %s", err.Error(), kit.Rid)
		return nil, err
	}

	//if !rsp.Result {
	//	blog.Errorf("[operation-cls] failed to search the classification by the condition(%#v), error info is %s, rid: %s", fCond, rsp.ErrMsg, kit.Rid)
	//	return nil, kit.CCError.New(rsp.Code, rsp.ErrMsg)
	//}
	//
	//clsIDs := make([]string, 0)
	//for _, cls := range rsp.Data.Info {
	//	clsIDs = append(clsIDs, cls.ClassificationID)
	//}
	//clsIDs = util.StrArrayUnique(clsIDs)
	//queryObjectCond := condition.CreateCondition().Field(common.BKClassificationIDField).In(clsIDs)
	//queryObjectResp, err := c.clientSet.CoreService().Model().ReadModel(kit.Ctx, kit.Header, &metadata.QueryCondition{Condition: queryObjectCond.ToMapStr()})
	//if nil != err {
	//	blog.Errorf("[operation-cls]failed to request the object controller, error info is %s, rid: %s", err.Error(), kit.Rid)
	//	return nil, err
	//}
	//if !queryObjectResp.Result {
	//	blog.Errorf("[operation-cls] failed to search the classification by the condition(%#v), error info is %s, rid: %s", fCond, queryObjectResp.ErrMsg, kit.Rid)
	//	return nil, kit.CCError.New(queryObjectResp.Code, queryObjectResp.ErrMsg)
	//}
	objMap := make(map[string][]metadata.Object)
	//objIDs := make([]string, 0)
	//for _, info := range queryObjectResp.Data.Info {
	//	objIDs = append(objIDs, info.Spec.ObjectID)
	//	objMap[info.Spec.ObjCls] = append(objMap[info.Spec.ObjCls], info.Spec)
	//}
	//
	datas := make([]metadata.ClassificationWithObject, 0)
	for _, cls := range rsp.Data.Info {
		clsItem := metadata.ClassificationWithObject{
			Classification: cls,
			Objects:        []metadata.Object{},
		}
		if obj, ok := objMap[cls.ClassificationID]; ok {
			clsItem.Objects = obj
		}
		datas = append(datas, clsItem)
	}

	return datas, nil
}
