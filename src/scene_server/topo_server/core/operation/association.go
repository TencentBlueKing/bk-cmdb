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
	"configcenter/src/apimachinery"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	frtypes "configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/topo_server/core/model"
	"configcenter/src/scene_server/topo_server/core/types"
)

// AssociationOperationInterface association operation methods
type AssociationOperationInterface interface {
	CreateMainLineAssociation(params types.LogicParams, data *metadata.Association) (model.Association, error)
	CreateCommonAssociation(params types.LogicParams, data *metadata.Association) (model.Association, error)
	DeleteAssociation(params types.LogicParams, cond condition.Condition) error
	UpdateAssociation(params types.LogicParams, data frtypes.MapStr, cond condition.Condition) error
}

type association struct {
	clientSet    apimachinery.ClientSetInterface
	obj          ObjectOperationInterface
	attr         AttributeOperationInterface
	inst         InstOperationInterface
	modelFactory model.Factory
}

// NewAssociationOperation create a new association operation instance
func NewAssociationOperation(client apimachinery.ClientSetInterface, obj ObjectOperationInterface, attr AttributeOperationInterface, inst InstOperationInterface, targetModel model.Factory) AssociationOperationInterface {
	return &association{
		clientSet:    client,
		obj:          obj,
		attr:         attr,
		inst:         inst,
		modelFactory: targetModel,
	}
}

func (cli *association) CreateMainLineAssociation(params types.LogicParams, data *metadata.Association) (model.Association, error) {

	// TODO : cache the association onbject's parent model

	theAsstObjectParent, err := cli.FetchAssociationParentByChild(params, data.AsstObjID, data.ObjectAttID)
	if nil != err {
		blog.Errorf("[operation-asst] failed to fetch association parent , error info is %s", err.Error())
		return nil, err
	}
	_ = theAsstObjectParent
	// TODO: cache the old inst

	// TODO: create the main line object

	srcObj, err := cli.obj.FindSingleObject(params, data.ObjectID)
	if nil != err {
		blog.Errorf("[operation-asst] failed to find the single object(%s), error info is %s", data.ObjectID, err.Error())
		return nil, err
	}

	asstObj, err := cli.obj.FindSingleObject(params, data.AsstObjID)
	if nil != err {
		blog.Errorf("[operation-asst] failed to find the single object(%s), error info is %s", data.AsstObjID, err.Error())
		return nil, err
	}

	asst := cli.modelFactory.CreateMainLineAssociatin(params, srcObj, data.ObjectAttID, asstObj)

	if exists, err := asst.IsExists(); nil != err {
		blog.Errorf("[operation-asst] failed to create the main line association, error info is %s", err.Error())
		return nil, err
	} else if exists {
		return nil, params.Err.Errorf(common.CCErrCommDuplicateItem, "main line association")
	}

	if err := asst.Create(); nil != err {
		blog.Errorf("[operation-asst] failed to create the main line association, error info is %s", err.Error())
		return nil, err
	}

	// TODO: create the main line association

	//asst := cli.modelFactory.CreateMainLineAssociatin(params)
	//_, err := asst.Parse(data)
	//if nil != err {
	//		blog.Errorf("[operation-asst] failed to parse the data, error info is %s", err.Error())
	//		return nil, err
	//	}
	/*
		if exists, err := asst.IsExists(); nil != err {
			blog.Errorf("[operation-asst] failed to check the association(%#v), error info is %s", data, err.Error())
			return nil, err
		} else if exists {
			return nil, nil
		}

		// TODO: update the old association inst
		// TODO: create the default inst for the new created mainline object

		return asst, nil
	*/
	return nil, nil
}
func (cli *association) CreateCommonAssociation(params types.LogicParams, data *metadata.Association) (model.Association, error) {

	//  check the association
	//	cond := condition.CreateCondition()
	//	cond.Field(metadata.AssociationFieldAssociationObjectID).Eq(data.AsstObjID)
	//	cond.Field(metadata.AssociationFieldObjectAttributeID).Eq(data.ObjectAttID)

	//asst := cli.modelFactory.(params)
	//asst.Parse(data)

	//cli.clientSet.ObjectController().Meta().SelectObjectAssociations()

	return nil, nil
}
func (cli *association) DeleteAssociation(params types.LogicParams, cond condition.Condition) error {
	return nil
}
func (cli *association) UpdateAssociation(params types.LogicParams, data frtypes.MapStr, cond condition.Condition) error {
	return nil
}
