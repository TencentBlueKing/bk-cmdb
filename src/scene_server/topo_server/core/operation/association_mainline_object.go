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
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/topo_server/core/model"
	"configcenter/src/scene_server/topo_server/core/types"
)

func (cli *association) CreateMainlineAssociation(params types.LogicParams, data *metadata.Association) (model.Association, error) {

	// check and fetch the association object's classification
	objCls, err := cli.cls.FindSingleClassification(params, data.ClassificationID)
	if nil != err {
		blog.Errorf("[opration-asst] failed to find the single classification, error info is %s", err.Error())
		return nil, params.Err.Errorf(common.CCErrCommParamsInvalid, data.ClassificationID)
	}

	// find the mainline parent object
	parentObj, err := cli.obj.FindSingleObject(params, data.AsstObjID)
	switch t := err.(type) {
	case nil:
	default:
		blog.Errorf("[operation-asst] failed to find the single object(%s), error info is %s", data.ObjectID, t.Error())
		return nil, t
	case errors.CCErrorCoder:
		if t.GetCode() == common.CCErrTopoObjectSelectFailed {
			blog.Errorf("[operation-asst] failed to find the single object(%s), error info is %s", data.ObjectID, t.Error())
			return nil, t
		}
	}

	// find the mainline child object for the parent
	childObj, err := parentObj.GetMainlineChildObject()
	if nil != err {
		blog.Errorf("[operation-asst] failed to find the child object for the object(%s), error info is %s", parentObj.GetID(), err.Error())
		return nil, err
	}

	// check and create the association mainline object
	currentObj, err := cli.obj.FindSingleObject(params, data.ObjectID)
	switch t := err.(type) {
	case nil:
	default:
		blog.Errorf("[operation-asst] failed to find the single object(%s), error info is %s", data.AsstObjID, err.Error())
		return nil, t
	case errors.CCErrorCoder:
		if t.GetCode() == common.CCErrTopoObjectSelectFailed {

			currentObj = cli.modelFactory.CreaetObject(params)
			currentObj.SetID(data.ObjectID)
			currentObj.SetName(data.ObjectName)
			currentObj.SetIcon(data.ObjectIcon)
			currentObj.SetClassification(objCls)

			if err = currentObj.Save(); nil != err {
				blog.Errorf("[operation-asst] failed to create the object(%s), error info is %s", data.AsstObjID, err.Error())
				return nil, err
			}

			attr := currentObj.CreateAttribute()
			attr.SetID(common.BKChildStr)
			if err = attr.Save(); nil != err {
				blog.Errorf("[operation-asst] failed to create the object(%s) attribute(%s), error info is %s", data.AsstObjID, common.BKChildStr, err.Error())
				return nil, err
			}
		}
	}

	// update the mainline topo inst association
	if err = cli.SetMainlineInstAssociation(params, parentObj, currentObj, childObj); nil != err {
		blog.Errorf("[operation-asst] failed set the mainline inst association, error info is %s", err.Error())
		return nil, err
	}

	// reset the parent's child object
	if err = parentObj.SetMainlineChildObject(currentObj.GetID()); nil != err {
		blog.Errorf("[operation-asst] failed to set the mainline object, error info is %s", err.Error())
		return nil, err
	}

	// reset the current's child object
	if err = currentObj.SetMainlineChildObject(childObj.GetID()); nil != err {
		blog.Errorf("[operation-asst] failed to set the mainline object, error info is %s ", err.Error())
		return nil, err
	}

	return nil, nil
}
