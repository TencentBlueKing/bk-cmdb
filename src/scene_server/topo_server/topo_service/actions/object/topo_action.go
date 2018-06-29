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
	"configcenter/src/common"
	"configcenter/src/common/bkbase"
	"configcenter/src/common/blog"
	"configcenter/src/common/core/cc/actions"
	"configcenter/src/common/errors"
	httpcli "configcenter/src/common/http/httpclient"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/topo_server/topo_service/manager"
	api "configcenter/src/source_controller/api/object"
	"encoding/json"
	"fmt"
	"github.com/bitly/go-simplejson"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"

	restful "github.com/emicklei/go-restful"
)

var topo = &topoAction{}

// MainLineObject main line object definition
type MainLineObject struct {
	api.ObjDes    `json:",inline"`
	AssociationID string `json:"bk_asst_obj_id"`
}

// topoAction
type topoAction struct {
	base.BaseAction
	mgr manager.Manager
}

func init() {

	// register action
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPCreate, Path: "/model/mainline", Params: nil, Handler: topo.CreateModel})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPDelete, Path: "/model/mainline/owners/{owner_id}/objectids/{obj_id}", Params: nil, Handler: topo.DeleteModel})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectGet, Path: "/model/{owner_id}", Params: nil, Handler: topo.SelectModel})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectGet, Path: "/model/{owner_id}/{cls_id}/{obj_id}", Params: nil, Handler: topo.SelectModelByClsID})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectGet, Path: "/inst/{owner_id}/{app_id}", Params: nil, Handler: topo.SelectInst})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectGet, Path: "/inst/child/{owner_id}/{obj_id}/{app_id}/{inst_id}", Params: nil, Handler: topo.SelectInstChild})

	// create cc object
	topo.CreateAction()
	// set manager
	manager.SetManager(topo)

}

// SetManager implement the manager's Hooker interface
func (cli *topoAction) SetManager(mgr manager.Manager) error {
	cli.mgr = mgr
	return nil
}
func (cli *topoAction) createDefaultInst(ownerID, objectID, objectName string, parentID int, req *restful.Request) (instID int, err error) {

	input := make(map[string]interface{})

	input[common.BKOwnerIDField] = ownerID
	input[common.BKInstParentStr] = parentID
	input[common.BKDefaultField] = 0
	input[common.CreateTimeField] = util.GetCurrentTimeStr()

	targetOBJ := ""
	switch objectID {
	case common.BKInnerObjIDApp:
		targetOBJ = common.BKInnerObjIDApp
		input[common.BKAppNameField] = objectName
	case common.BKInnerObjIDModule:
		targetOBJ = common.BKInnerObjIDModule
		input[common.BKModuleNameField] = objectName
	case common.BKInnerObjIDSet:
		targetOBJ = common.BKInnerObjIDSet
		input[common.BKSetNameField] = objectName
	default:
		targetOBJ = common.BKINnerObjIDObject
		input[common.BKObjIDField] = objectID
		input[common.BKInstNameField] = objectName
	}

	inputJSON, jsErr := json.Marshal(input)
	if nil != jsErr {
		blog.Error("the input json is invalid, error info is %s", jsErr.Error())
		return 0, fmt.Errorf("failed to unmarshal the json")
	}

	cURL := cli.CC.ObjCtrl() + "/object/v1/insts/" + targetOBJ
	blog.Debug("inst:%v", string(inputJSON))

	instRes, err := httpcli.ReqHttp(req, cURL, "POST", inputJSON)
	if nil != err {
		blog.Error("create inst failed, errors:%s", err.Error())
		return 0, fmt.Errorf("failed to request the inst")
	}

	rsp, rspOk := cli.IsSuccess([]byte(instRes))
	if !rspOk {
		return 0, fmt.Errorf("failed to create inst, error info is %+v", rsp.Message)
	}

	if rspMap, rspOk := rsp.Data.(map[string]interface{}); rspOk {
		if rspData, ok := rspMap[common.BKInstIDField]; ok {
			//if dataMap, dataOk := rspData.(map[string]interface{}); dataOk {
			//if objID, ok := rspData["ObjectID"]; ok {
			switch id := rspData.(type) {
			case int:
				return id, nil
			case float64:
				return int(id), nil
			case float32:
				return int(id), nil
			case json.Number:
				_id, err := id.Int64()
				if nil == err {
					return 0, err
				}
				return int(_id), nil
			}
			//}
			//return 0, fmt.Errorf("not found the 'ObjectID' filed in the response(%+v)", dataMap)
			//}
			//return 0, fmt.Errorf("can not convert to map[string]interface{}, the type is %s", reflect.TypeOf(rspMap["data"]).Kind())
		}
		return 0, fmt.Errorf("not found the 'inst id' filed in the response(%+v)", rsp.Data)
	}
	return 0, fmt.Errorf("can not convert to map[string]interface{}, the type is %s", reflect.TypeOf(rsp.Data).Kind())
}

func (cli *topoAction) updateOldInstParentID(ownerID string, parentInstID int, cacheInst []manager.TopoInstRst, req *restful.Request) error {

	input := make(map[string]interface{})

	for _, inst := range cacheInst {

		condition := make(map[string]interface{})

		condition[common.BKOwnerIDField] = ownerID

		input["data"] = map[string]interface{}{common.BKInstParentStr: parentInstID}
		targetOBJ := ""
		switch inst.ObjID {
		case common.BKInnerObjIDApp:
			targetOBJ = common.BKInnerObjIDApp
			condition[common.BKAppIDField] = inst.InstID
		case common.BKInnerObjIDModule:
			targetOBJ = common.BKInnerObjIDModule
			condition[common.BKModuleIDField] = inst.InstID
		case common.BKInnerObjIDSet:
			targetOBJ = common.BKInnerObjIDSet
			condition[common.BKSetIDField] = inst.InstID
		default:
			targetOBJ = common.BKINnerObjIDObject
			condition[common.BKInstIDField] = inst.InstID
			condition[common.BKObjIDField] = inst.ObjID
		}

		input["condition"] = condition
		uURL := cli.CC.ObjCtrl() + "/object/v1/insts/" + targetOBJ

		inputJSON, jsErr := json.Marshal(input)
		if nil != jsErr {
			blog.Error("failed to create json object, error info is %s", jsErr.Error())
			continue
		}

		objRes, err := httpcli.ReqHttp(req, uURL, "PUT", []byte(inputJSON))
		if nil != err {
			blog.Error("failed to update the inst, error info is %s", err.Error())
			continue
		}

		if rsp, rspOk := cli.IsSuccess([]byte(objRes)); !rspOk {
			blog.Error("failed to update the object, error info is %+v ", rsp.Message)
		}
	}

	return nil
}

func (cli *topoAction) deleteInsts(ownerID string, target []manager.TopoInstRst, req *restful.Request) error {

	for _, inst := range target {

		condition := make(map[string]interface{})

		condition[common.BKOwnerIDField] = ownerID

		targetOBJ := ""
		switch inst.ObjID {
		case common.BKInnerObjIDApp:
			targetOBJ = common.BKInnerObjIDApp
			condition[common.BKAppIDField] = inst.InstID
		case common.BKInnerObjIDModule:
			targetOBJ = common.BKInnerObjIDModule
			condition[common.BKModuleIDField] = inst.InstID
		case common.BKInnerObjIDSet:
			targetOBJ = common.BKInnerObjIDSet
			condition[common.BKSetIDField] = inst.InstID
		default:
			targetOBJ = common.BKINnerObjIDObject
			condition[common.BKObjIDField] = inst.ObjID
			condition[common.BKInstIDField] = inst.InstID
		}
		uURL := cli.CC.ObjCtrl() + "/object/v1/insts/" + targetOBJ
		inputJSON, jsErr := json.Marshal(condition)
		if nil != jsErr {
			blog.Error("failed to create json object, error info is %s", jsErr.Error())
			continue
		}

		objRes, err := httpcli.ReqHttp(req, uURL, "DELETE", []byte(inputJSON))
		if nil != err {
			blog.Error("failed to delete the inst, error info is %s", err.Error())
			continue
		}

		if rsp, rspOk := cli.IsSuccess([]byte(objRes)); !rspOk {
			blog.Error("failed to delete the object, error info is %+v ", rsp.Message)
		}
	}

	return nil
}

func (cli *topoAction) updateInsts(ownerID string, child []manager.TopoInstRst, parent []manager.TopoInstRst, req *restful.Request) error {

	blog.Debug("cur: %+v  parent: %+v", child, parent)
	for _, parentItem := range parent {

		for _, subParentItem := range parentItem.Child {

			for _, childItem := range child {

				// find the same inst, update the inst's child
				if subParentItem.InstID == childItem.InstID {
					cli.updateOldInstParentID(ownerID, parentItem.InstID, childItem.Child, req)
				}
			}

		}
	}

	return nil
}

func (cli *topoAction) createInsts(ownerID, objectID, objectName string, cacheInst []manager.TopoInstRst, req *restful.Request) error {

	blog.Debug("topo action, the topo inst cache, %+v", cacheInst)
	for _, inst := range cacheInst {

		// create the default inst
		instID, instIDErr := cli.createDefaultInst(ownerID, objectID, objectName, inst.InstID, req)
		if nil != instIDErr {
			blog.Error("failed to create the default inst, error info is %s ", instIDErr.Error())
			//return instIDErr
			continue
		}

		// update the old child inst parentid
		if updateInstErr := cli.updateOldInstParentID(ownerID, instID, inst.Child, req); nil != updateInstErr {
			blog.Error("failed to update the old inst parentid, error info is %s", updateInstErr.Error())
			continue
		}
	}

	return nil
}

// updateMainModule update the mainline object topo
func (cli *topoAction) updateMainModule(forward *api.ForwardParam, oldAsstItems []api.ObjAsstDes, objectID string, errProxy errors.DefaultCCErrorIf) error {

	// to update the old main line object association
	for _, objAsst := range oldAsstItems {
		newObj := map[string]interface{}{}
		newObj["bk_asst_obj_id"] = objectID
		tmpData, _ := json.Marshal(objAsst)
		js, _ := simplejson.NewJson(tmpData)
		cond, _ := js.Map()
		if err := cli.mgr.UpdateObjectAsst(forward, cond, newObj, errProxy); nil != err {
			blog.Error("failed to update the mainline association, error info is %s ", err.Error())
		}
	}

	return nil
}

// CreateModel create main line object
func (cli *topoAction) CreateModel(req *restful.Request, resp *restful.Response) {

	blog.Info("create main line object")
	// get the language
	language := util.GetActionLanguage(req)

	// get the default error by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)
	forward := &api.ForwardParam{Header: req.Request.Header}

	cli.CallResponseEx(func() (int, interface{}, error) {
		// read data
		val, err := ioutil.ReadAll(req.Request.Body)
		if err != nil {
			blog.Error("read request body failed, error information is %s", err.Error())
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommHTTPReadBodyFailed)
		}

		// deal data
		var obj MainLineObject
		if jsErr := json.Unmarshal(val, &obj); nil != jsErr {
			blog.Error("unmarshal json failed, error information is %v", jsErr)
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommJSONUnmarshalFailed)
		}

		// check the object level limit
		level := common.BKTopoBusinessLevelDefault
		if config, err := cli.CC.ParseConfig(); nil != err {
			blog.Errorf("failed to get the parse the conigure, error info is %s", err.Error())
		} else if cfg, ok := config[common.BKTopoBusinessLevelLimit]; ok {
			level, err := strconv.Atoi(cfg)
			if nil != err {
				blog.Errorf("can not convert level(%s) to int, error info is %s", cfg, err.Error())
			}
			if level <= 0 { // the min level limit is 3
				level = common.BKTopoBusinessLevelDefault
			}
		}

		rstItems, ctrErr := cli.mgr.SelectTopoModel(forward, nil, obj.OwnerID, common.BKInnerObjIDApp, "", "", "", defErr)
		if nil != ctrErr {
			blog.Error("select topo model failed, error information is %v, res:%v", ctrErr, rstItems)
			return http.StatusBadRequest, nil, ctrErr
		}

		//blog.Debug("select module for insts:%v", rstItems)

		if level <= len(rstItems) {
			blog.Errorf("business topology level exceeds the limit, %d", level)
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrTopoBizTopoLevelOverLimit)
		}

		// to cache the old main association by the parent
		asstSearch := map[string]interface{}{}
		asstSearch[common.BKOwnerIDField] = obj.OwnerID
		asstSearch["bk_asst_obj_id"] = obj.AssociationID
		asstSearch["bk_object_att_id"] = common.BKChildStr
		asstDesItems, asstDesItemsErr := cli.mgr.SelectObjectAsst(forward, asstSearch, defErr)
		if nil != asstDesItemsErr {
			blog.Error("failed to cache the old asst, error info is %s", asstDesItemsErr.Error())
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommJSONUnmarshalFailed)
		}

		// to cache the old main inst data, read the parent insts
		asstInstItems, rsterr := cli.SelectInstTopo(forward, obj.OwnerID, obj.AssociationID, 0, 0, 2, req)
		if nil != rsterr {
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrTopoMainlineCreatFailed)
		}

		// create a new main line object
		if _, err := cli.mgr.CreateObject(forward, val, defErr); nil != err {
			blog.Error("failed to create the main line object, error info is %s", err.Error())
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrTopoMainlineCreatFailed)
		}

		// create a new main line topo
		objAtt := api.ObjAttDes{}
		objAtt.ObjectID = obj.ObjectID
		objAtt.OwnerID = obj.OwnerID
		objAtt.AssociationID = obj.AssociationID
		objAtt.AssoType = common.BKChild
		if _, ctrErr := cli.mgr.CreateTopoModel(forward, objAtt, defErr); nil != ctrErr {
			blog.Error("create objectatt failed, error information is %s", ctrErr.Error())
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrTopoMainlineCreatFailed)
		}
		// to update the old asst, must be first
		cli.updateMainModule(forward, asstDesItems, obj.ObjectID, defErr)

		// to create insts, must be second
		if err := cli.createInsts(obj.OwnerID, obj.ObjectID, obj.ObjectName, asstInstItems, req); nil != err {
			blog.Error("failed to create the default inst , error info is %s", err.Error())
		}

		return http.StatusOK, nil, nil
	}, resp)
}

// hasChildInstNameRepeat check the deleted inst name wether it is repeated
func (cli *topoAction) hasChildInstNameRepeat(current []manager.TopoInstRst, parent []manager.TopoInstRst) bool {

	blog.Debug("check name cur: %+v ", current)
	tmpItems := map[string]*struct{}{}

	for _, parentItem := range parent {

		for _, subParentItem := range parentItem.Child {

			for _, curItem := range current {

				// find the same inst, update the inst's child
				if subParentItem.InstID == curItem.InstID {

					for _, subCurItem := range curItem.Child {
						key := fmt.Sprintf("%d_%s", parentItem.InstID, subCurItem.InstName)
						if _, ok := tmpItems[key]; ok {
							blog.Debug("already exists %v, target %s", tmpItems, key)
							return true
						}

						tmpItems[key] = nil
					}
				}
			}
		}
	}

	return false
}

// DeleteModule 删除模型
func (cli *topoAction) DeleteModel(req *restful.Request, resp *restful.Response) {

	blog.Info("delete model")
	// get the language
	language := util.GetActionLanguage(req)

	// get the default error by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)
	forward := &api.ForwardParam{Header: req.Request.Header}
	cli.CallResponseEx(func() (int, interface{}, error) {

		// read parameters
		ownerID := req.PathParameter("owner_id")
		objID := req.PathParameter("obj_id")
		switch objID {
		case common.BKInnerObjIDApp, common.BKInnerObjIDHost, common.BKInnerObjIDPlat,
			common.BKInnerObjIDModule, common.BKInnerObjIDSet, common.BKInnerObjIDProc:
			blog.Error("inner object, forbiden to delete %s", objID)
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrTopoForbiddenToDeleteModelFailed)
		}

		// to cache the old main association by the parent
		asstSearch := map[string]interface{}{}
		asstSearch[common.BKOwnerIDField] = ownerID
		asstSearch["bk_asst_obj_id"] = objID
		asstSearch["bk_object_att_id"] = common.BKChildStr
		asstChildDesItems, asstDesItemsErr := cli.mgr.SelectObjectAsst(forward, asstSearch, defErr)
		if nil != asstDesItemsErr {
			blog.Error("failed to cache the old asst, error info is %s", asstDesItemsErr.Error())
			//cli.ResponseFailed(common.CC_Err_Comm_http_DO, asstDesItemsErr.Error(), resp)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrTopoMainlineDeleteFailed)
		}
		delete(asstSearch, "bk_asst_obj_id")
		asstSearch[common.BKObjIDField] = objID
		asstParentDesItems, asstDesItemsErr := cli.mgr.SelectObjectAsst(forward, asstSearch, defErr)
		if nil != asstDesItemsErr {
			blog.Error("failed to cache the old asst, error info is %s", asstDesItemsErr.Error())
			//cli.ResponseFailed(common.CC_Err_Comm_http_DO, asstDesItemsErr.Error(), resp)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrTopoMainlineDeleteFailed)
		}

		if 0 == len(asstParentDesItems) {
			blog.Error("not found the parent object,  %s", objID)
			//cli.ResponseFailed(common.CC_Err_Comm_http_DO, "not found the parent object", resp)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrTopoMainlineDeleteFailed)
		}

		// to cache the old main inst data, read the parent insts
		asstChildInstItems, rsterr := cli.SelectInstTopo(forward, ownerID, objID, 0, 0, 2, req)
		if nil != rsterr {
			//cli.ResponseFailed(common.CC_Err_Comm_http_DO, rsterr.Error(), resp)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrTopoMainlineDeleteFailed)
		}

		asstParentInstItems, rstErr := cli.SelectInstTopo(forward, ownerID, asstParentDesItems[0].AsstObjID, 0, 0, 2, req)
		if nil != rstErr {
			cli.ResponseFailed(common.CC_Err_Comm_http_DO, rsterr.Error(), resp)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrTopoMainlineDeleteFailed)
		}

		if cli.hasChildInstNameRepeat(asstChildInstItems, asstParentInstItems) {
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrTopoDeleteMainLineObjectAndInstNameRepeat)
		}

		// deal data
		if ctrErr := cli.mgr.DeleteTopoModel(forward, ownerID, objID, common.BKChild, defErr); nil == ctrErr {

			// to delete the old main line object
			objDes := map[string]interface{}{}
			objDes[common.BKOwnerIDField] = ownerID
			objDes[common.BKObjIDField] = objID
			objDesJSON, _ := json.Marshal(objDes)
			if err := cli.mgr.DeleteObject(forward, 0, objDesJSON, defErr); nil != err {
				blog.Error("failed to delete the object module(%s), data(%s) error info is %s", objID, string(objDesJSON), err.Error())
			}

			// update the main line module association
			cli.updateMainModule(forward, asstChildDesItems, asstParentDesItems[0].AsstObjID, defErr)

			// update the main inst association
			cli.updateInsts(ownerID, asstChildInstItems, asstParentInstItems, req)

			// deleete the old main inst
			cli.deleteInsts(ownerID, asstChildInstItems, req)

		} else {
			blog.Error("create objectatt failed, error information is %s", ctrErr.Error())
			cli.ResponseFailed(common.CC_Err_Comm_http_DO, ctrErr.Error(), resp)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrTopoMainlineDeleteFailed)
		}
		return http.StatusOK, nil, nil

	}, resp)

}

// SelectModule 查询模型拓扑
func (cli *topoAction) SelectModel(req *restful.Request, resp *restful.Response) {
	language := util.GetActionLanguage(req)
	// get the error factory by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)
	forward := &api.ForwardParam{Header: req.Request.Header}

	cli.CallResponseEx(func() (int, interface{}, error) {
		blog.Info("select model topo")

		ownerID := req.PathParameter("owner_id")

		// deal data
		result, ctrErr := cli.mgr.SelectTopoModel(forward, nil, ownerID, common.BKInnerObjIDApp, "", "", "", defErr)
		if nil != ctrErr {
			blog.Error("select topo model failed, error information is %s", ctrErr.Error())
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrTopoMainlineSelectFailed)
		}
		return http.StatusOK, result, nil
	}, resp)

}

// SelectModelByClsID
func (cli *topoAction) SelectModelByClsID(req *restful.Request, resp *restful.Response) {

	blog.Info("select model topo by objcls")
	language := util.GetActionLanguage(req)
	// get the error factory by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)
	forward := &api.ForwardParam{Header: req.Request.Header}

	cli.CallResponseEx(func() (int, interface{}, error) {

		ownerID := req.PathParameter("owner_id")
		clsID := req.PathParameter("cls_id")
		objID := req.PathParameter("obj_id")

		result, ctrErr := cli.mgr.SelectTopoModel(forward, nil, ownerID, objID, clsID, "", "", defErr)
		if nil != ctrErr {
			blog.Error("select topo model failed, error information is %s", ctrErr.Error())
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrTopoMainlineSelectFailed)
		}
		return http.StatusOK, result, nil
	}, resp)

}

// selectInstDetail select instDetail
func (cli *topoAction) selectInstDetail(req *restful.Request, ownerID, objID string, appID, instID, parentID int) (*ObjInstRsp, error) {

	blog.Debug("select inst detail ownerid: %s objid: %s appid: %d parentid: %d", ownerID, objID, appID, parentID)
	condition := make(map[string]interface{})

	condition[common.BKOwnerIDField] = ownerID

	searchParams := make(map[string]interface{})

	searchParams["fields"] = ""
	searchParams["start"] = 0
	searchParams["limit"] = common.BKNoLimit

	var targetobj string

	switch objID {
	case common.BKInnerObjIDApp:
		targetobj = common.BKInnerObjIDApp
		searchParams["sort"] = common.BKAppNameField
		if 0 != appID {
			condition[common.BKAppIDField] = appID
		}
	case common.BKInnerObjIDSet:
		targetobj = common.BKInnerObjIDSet
		searchParams["sort"] = common.BKSetNameField
		if 0 != parentID {
			condition[common.BKInstParentStr] = parentID
		}
		if 0 != instID {
			condition[common.BKSetIDField] = instID
		}
		condition[common.BKDefaultField] = map[string]interface{}{
			common.BKDBNE: common.DefaultResSetFlag,
		}

	case common.BKInnerObjIDModule:
		targetobj = common.BKInnerObjIDModule
		searchParams["sort"] = common.BKModuleNameField
		if 0 != parentID {
			condition[common.BKInstParentStr] = parentID
		}
		if 0 != instID {
			condition[common.BKModuleIDField] = instID
		}
	default:
		targetobj = common.BKINnerObjIDObject
		condition[common.BKObjIDField] = objID
		searchParams["sort"] = common.BKInstNameField
		if 0 != parentID {
			condition[common.BKInstParentStr] = parentID
		}
		if 0 != instID {
			condition[common.BKInstIDField] = instID
		}
	}

	searchParams["condition"] = condition

	//search
	sURL := cli.CC.ObjCtrl() + "/object/v1/insts/" + targetobj + "/search"
	inputJSON, _ := json.Marshal(searchParams)
	objRes, err := httpcli.ReqHttp(req, sURL, common.HTTPSelectPost, []byte(inputJSON))
	blog.Debug("search inst detail params: %s", string(inputJSON))
	if nil != err {
		blog.Error("search inst defail failed, error: %v", err)
		return nil, err
	}

	var instRsp ObjInstRsp
	jsErr := json.Unmarshal([]byte(objRes), &instRsp)
	if nil != jsErr {
		blog.Error("unmarshal failed, data:%s error:%v", objRes, jsErr)
		return nil, jsErr
	}

	blog.Debug("inst result:%v", instRsp)
	return &instRsp, nil
}

func (cli *topoAction) selectChildInstTopo(req *restful.Request, topoItem []manager.TopoModelRsp, level, appID, parentID int, childs *[]manager.TopoInstRst) error {

	blog.Debug("level[%d]%v", level, topoItem)

	if 0 == len(topoItem) {
		return nil
	}

	// 取当前父实例节点的子实例节点
	childInst, instErr := cli.selectInstDetail(req, topoItem[0].OwnerID, topoItem[0].ObjID, appID, 0, parentID)

	if nil != instErr {
		blog.Error("can not select insts %v", instErr)
		return instErr
	}

	for _, childInstItem := range childInst.Data.Info {

		blog.Debug("child inst item:%v", childInstItem)

		childInstRst := manager.TopoInstRst{}
		childInstRst.ObjID = topoItem[0].ObjID
		childInstRst.ObjName = topoItem[0].ObjName

		if dftval, ok := childInstItem[common.BKDefaultField]; ok {
			switch dftval.(type) {
			case string:
				childInstRst.Default, _ = strconv.Atoi(dftval.(string))
			case float64:
				childInstRst.Default = int(dftval.(float64))
			}
		}

		switch topoItem[0].ObjID {
		case common.BKInnerObjIDApp:
			if id, ok := childInstItem[common.BKAppIDField]; ok {
				childInstRst.InstID = int(id.(float64))
			}
			if name, ok := childInstItem[common.BKAppNameField]; ok {
				childInstRst.InstName = name.(string)
			}
		case common.BKInnerObjIDModule:
			if id, ok := childInstItem[common.BKModuleIDField]; ok {
				childInstRst.InstID = int(id.(float64))
			}
			if name, ok := childInstItem[common.BKModuleNameField]; ok {
				childInstRst.InstName = name.(string)
			}
		case common.BKInnerObjIDSet:
			if id, ok := childInstItem[common.BKSetIDField]; ok {
				childInstRst.InstID = int(id.(float64))
			}
			if name, ok := childInstItem[common.BKSetNameField]; ok {
				childInstRst.InstName = name.(string)
			}
		default:
			if id, ok := childInstItem[common.BKInstIDField]; ok {
				childInstRst.InstID = int(id.(float64))
			}
			if name, ok := childInstItem[common.BKInstNameField]; ok {
				childInstRst.InstName = name.(string)
			}
		}

		if level >= 0 {
			childInstRst.Child = make([]manager.TopoInstRst, 0)
			cli.selectChildInstTopo(req, topoItem[1:], level-1, appID, childInstRst.InstID, &childInstRst.Child)
		}
		*childs = append(*childs, childInstRst)
		blog.Debug("childs:%v", childs)
	}

	return nil
}

func (cli *topoAction) SelectInstTopo(forward *api.ForwardParam, ownerID, objID string, appID, instID, level int, req *restful.Request) ([]manager.TopoInstRst, error) {

	// get the language
	language := util.GetActionLanguage(req)

	// get the default error by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)

	rstItems, ctrErr := cli.mgr.SelectTopoModel(forward, nil, ownerID, objID, "", "", "", defErr)
	if nil != ctrErr {
		blog.Error("select topo model failed, error information is %v, res:%v", ctrErr, rstItems)
		return nil, ctrErr
	}

	blog.Debug("select module for insts:%v", rstItems)

	if 0 == len(rstItems) {
		blog.Warn("module items is empty")
		return []manager.TopoInstRst{}, nil
	}

	if level > len(rstItems) || level <= 0 {
		level = len(rstItems)
	}

	instRstItems := make([]manager.TopoInstRst, 0)

	// the default sequence：a=>b=>c=>d
	item := rstItems[0]

	// get the current inst
	inst, instErr := cli.selectInstDetail(req, ownerID, item.ObjID, appID, instID, 0)

	if nil != instErr {
		blog.Error("can not select insts %v", instErr)
		return nil, instErr
	}

	// foreach all insts of the object
	for _, instItem := range inst.Data.Info {
		blog.Debug("child inst item:%v", instItem)
		instRst := manager.TopoInstRst{}
		instRst.ObjID = item.ObjID
		instRst.ObjName = item.ObjName

		if dftval, ok := instItem[common.BKDefaultField]; ok {
			instRst.Default = int(dftval.(float64))
		}

		switch item.ObjID {

		case common.BKInnerObjIDApp:
			if id, ok := instItem[common.BKAppIDField]; ok {
				instRst.InstID = int(id.(float64))
			}
			if name, ok := instItem[common.BKAppNameField]; ok {
				instRst.InstName = name.(string)
			}
		case common.BKInnerObjIDModule:
			if id, ok := instItem[common.BKModuleIDField]; ok {
				instRst.InstID = int(id.(float64))
			}
			if name, ok := instItem[common.BKModuleNameField]; ok {
				instRst.InstName = name.(string)
			}
		case common.BKInnerObjIDSet:
			if id, ok := instItem[common.BKSetIDField]; ok {
				instRst.InstID = int(id.(float64))
			}
			if name, ok := instItem[common.BKSetNameField]; ok {
				instRst.InstName = name.(string)
			}
		default:
			if id, ok := instItem[common.BKInstIDField]; ok {
				instRst.InstID = int(id.(float64))
			}
			if name, ok := instItem[common.BKInstNameField]; ok {
				instRst.InstName = name.(string)
			}
		}

		// select child inst
		instRst.Child = make([]manager.TopoInstRst, 0)
		if instErr := cli.selectChildInstTopo(req, rstItems[1:], level-1, appID, instRst.InstID, &instRst.Child); nil != instErr {
			blog.Error("read child failed, error:%s", instErr.Error())
			return nil, instErr
		}

		blog.Debug("instRsp:%v", instRst)
		instRstItems = append(instRstItems, instRst)
	}

	return instRstItems, nil
}

// SelectInst select inst
func (cli *topoAction) SelectInst(req *restful.Request, resp *restful.Response) {
	// get the language
	language := util.GetActionLanguage(req)
	// get the error factory by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)
	forward := &api.ForwardParam{Header: req.Request.Header}

	cli.CallResponseEx(func() (int, interface{}, error) {

		blog.Info("select model inst ")
		ownerID := req.PathParameter("owner_id")
		appID, _ := strconv.Atoi(req.PathParameter("app_id"))
		level := req.QueryParameter("level")
		iLevel := 2
		if 0 != len(level) {
			if tmp, err := strconv.Atoi(level); nil != err {
				blog.Errorf("the level(%s) is error,  error info is %s", level, err.Error())
			} else {
				iLevel = tmp
			}
		}
		rst, rstErr := cli.SelectInstTopo(forward, ownerID, common.BKInnerObjIDApp, appID, 0, iLevel, req)
		if nil != rstErr {
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrTopoTopoSelectFailed)
		}
		return http.StatusOK, rst, nil
	}, resp)

}

// SelectInstChild 查询子实例
func (cli *topoAction) SelectInstChild(req *restful.Request, resp *restful.Response) {
	// get the language
	language := util.GetActionLanguage(req)
	// get the error factory by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)
	forward := &api.ForwardParam{Header: req.Request.Header}

	cli.CallResponseEx(func() (int, interface{}, error) {
		blog.Info("select model inst child")

		ownerID := req.PathParameter("owner_id")
		appID, _ := strconv.Atoi(req.PathParameter("app_id"))
		instID, _ := strconv.Atoi(req.PathParameter("inst_id"))
		objID := req.PathParameter("obj_id")

		rst, rstErr := cli.SelectInstTopo(forward, ownerID, objID, appID, instID, 2, req)
		if nil != rstErr {
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrTopoTopoSelectFailed)
		}
		return http.StatusOK, rst, nil
	}, resp)

}
