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

package instdata

import (
	"configcenter/src/common"
	"configcenter/src/common/base"
	"configcenter/src/common/blog"
	"configcenter/src/common/core/cc/actions"
	"configcenter/src/common/util"
	eventtypes "configcenter/src/scene_server/event_server/types"
	"configcenter/src/source_controller/common/commondata"
	"configcenter/src/source_controller/common/eventdata"
	"configcenter/src/source_controller/common/instdata"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	simplejson "github.com/bitly/go-simplejson"
	"github.com/emicklei/go-restful"
)

var obj = &objectAction{}

// ObjectAction
type objectAction struct {
	base.BaseAction
}

func init() {
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "/insts/{obj_type}/search", Params: nil, Handler: obj.SearchObjects})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPCreate, Path: "/insts/{obj_type}", Params: nil, Handler: obj.CreateObject})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPDelete, Path: "/insts/{obj_type}", Params: nil, Handler: obj.DelObject})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPUpdate, Path: "/insts/{obj_type}", Params: nil, Handler: obj.UpdateObject})
	// set cc api interface
	obj.CreateAction()
}

//delete object
func (cli *objectAction) DelObject(req *restful.Request, resp *restful.Response) {
	blog.Info("delete insts")
	// get the language
	language := util.GetActionLanguage(req)
	ownerID := util.GetActionOnwerID(req)
	// get the error factory by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)
	defLang := cli.CC.Lang.CreateDefaultCCLanguageIf(language)
	cli.CallResponseEx(func() (int, interface{}, error) {

		pathParams := req.PathParameters()
		objType := pathParams["obj_type"]
		instdata.DataH = cli.CC.InstCli
		value, _ := ioutil.ReadAll(req.Request.Body)
		js, _ := simplejson.NewJson([]byte(value))
		input, _ := js.Map()
		util.SetModOwner(input, ownerID)

		// retrieve original datas
		originDatas := make([]map[string]interface{}, 0)
		getErr := instdata.GetObjectByCondition(defLang, objType, nil, input, &originDatas, "", 0, 0)
		if getErr != nil {
			blog.Error("retrieve original data error:%v", getErr)
		}

		blog.Info("delete object type:%s,input:%v ", objType, input)
		err := instdata.DelObjByCondition(objType, input)
		if err != nil {
			blog.Error("delete object type:%s,input:%v error:%v", objType, input, err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrObjectDeleteInstFailed)
		}

		// send events
		if len(originDatas) > 0 {
			ec := eventdata.NewEventContextByReq(req)
			for _, originData := range originDatas {
				err := ec.InsertEvent(eventtypes.EventTypeInstData, objType, eventtypes.EventActionDelete, nil, originData, ownerID)
				if err != nil {
					blog.Error("create event error:%v", err)
				}
			}
		}
		return http.StatusOK, nil, nil
	}, resp)

}

//update object
func (cli *objectAction) UpdateObject(req *restful.Request, resp *restful.Response) {
	blog.Info("update insts")
	// get the language
	language := util.GetActionLanguage(req)
	ownerID := util.GetActionOnwerID(req)
	// get the error factory by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)
	defLang := cli.CC.Lang.CreateDefaultCCLanguageIf(language)
	cli.CallResponseEx(func() (int, interface{}, error) {

		pathParams := req.PathParameters()
		objType := pathParams["obj_type"]
		instdata.DataH = cli.CC.InstCli
		value, _ := ioutil.ReadAll(req.Request.Body)
		js, _ := simplejson.NewJson([]byte(value))
		input, _ := js.Map()
		data := input["data"].(map[string]interface{})
		data[common.LastTimeField] = time.Now()
		condition := input["condition"]
		condition = util.SetModOwner(condition, ownerID)
		data = util.SetModOwner(data, ownerID)

		// retrieve original datas
		originDatas := make([]map[string]interface{}, 0)
		getErr := instdata.GetObjectByCondition(defLang, objType, nil, condition, &originDatas, "", 0, 0)
		if getErr != nil {
			blog.Error("retrieve original datas error:%v", getErr)
		}

		blog.Info("update object type:%s,data:%v,condition:%v", objType, data, condition)
		err := instdata.UpdateObjByCondition(objType, data, condition)
		if err != nil {
			blog.Error("update object type:%s,data:%v,condition:%v,error:%v", objType, data, condition, err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrObjectUpdateInstFailed)
		}

		// record event
		if len(originDatas) > 0 {
			newdatas := []map[string]interface{}{}
			if err := instdata.GetObjectByCondition(defLang, objType, nil, condition, &newdatas, "", 0, 0); err != nil {
				blog.Error("create event error:%v", err)
			} else {
				ec := eventdata.NewEventContextByReq(req)
				idname := instdata.GetIDNameByType(objType)
				for _, originData := range originDatas {
					newData := map[string]interface{}{}
					id, err := strconv.Atoi(fmt.Sprintf("%v", originData[idname]))
					if err != nil {
						blog.Errorf("create event error:%v", err)
						continue
					}
					if err := instdata.GetObjectByID(objType, nil, id, &newData, ""); err != nil {
						blog.Error("create event error:%v", err)
					} else {
						err := ec.InsertEvent(eventtypes.EventTypeInstData, objType, eventtypes.EventActionUpdate, newData, originData, ownerID)
						if err != nil {
							blog.Error("create event error:%v", err)
						}
					}
				}
			}
		}
		return http.StatusOK, nil, nil
	}, resp)
}

//search object
func (cli *objectAction) SearchObjects(req *restful.Request, resp *restful.Response) {
	blog.Info("select insts")
	// get the language
	language := util.GetActionLanguage(req)
	ownerID := util.GetActionOnwerID(req)
	// get the error factory by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)
	defLang := cli.CC.Lang.CreateDefaultCCLanguageIf(language)
	cli.CallResponseEx(func() (int, interface{}, error) {
		pathParams := req.PathParameters()
		objType := pathParams["obj_type"]
		instdata.DataH = cli.CC.InstCli

		value, err := ioutil.ReadAll(req.Request.Body)
		var dat commondata.ObjQueryInput
		err = json.Unmarshal([]byte(value), &dat)
		if err != nil {
			blog.Error("get object type:%s,input:%v error:%v", string(objType), value, err)
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommJSONUnmarshalFailed)

		}
		//dat.ConvTime()
		fields := dat.Fields
		condition := util.SetModOwner(dat.Condition, ownerID)

		skip := dat.Start
		limit := dat.Limit
		sort := dat.Sort
		fieldArr := strings.Split(fields, ",")
		result := make([]map[string]interface{}, 0)
		count, err := instdata.GetCntByCondition(objType, condition)
		if err != nil {
			blog.Error("get object type:%s,input:%v error:%v", objType, string(value), err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrObjectSelectInstFailed)

		}
		err = instdata.GetObjectByCondition(defLang, objType, fieldArr, condition, &result, sort, skip, limit)
		if err != nil {
			blog.Error("get object type:%s,input:%v error:%v", string(objType), string(value), err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrObjectSelectInstFailed)
		}
		info := make(map[string]interface{})
		info["count"] = count
		info["info"] = result
		return http.StatusOK, info, nil
	}, resp)
}

//create object
func (cli *objectAction) CreateObject(req *restful.Request, resp *restful.Response) {
	blog.Info("create insts")
	// get the language
	language := util.GetActionLanguage(req)
	ownerID := util.GetActionOnwerID(req)
	// get the error factory by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)

	cli.CallResponseEx(func() (int, interface{}, error) {
		pathParams := req.PathParameters()
		objType := pathParams["obj_type"]
		instdata.DataH = cli.CC.InstCli
		value, _ := ioutil.ReadAll(req.Request.Body)
		js, _ := simplejson.NewJson([]byte(value))
		input, _ := js.Map()
		input[common.CreateTimeField] = time.Now()
		input[common.LastTimeField] = time.Now()
		util.SetModOwner(input, ownerID)
		blog.Info("create object type:%s,data:%v", objType, input)
		var idName string
		id, err := instdata.CreateObject(objType, input, &idName)
		if err != nil {
			blog.Error("create object type:%s,data:%v error:%v", objType, input, err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrObjectCreateInstFailed)
		}

		// record event
		origindata := map[string]interface{}{}
		if err := instdata.GetObjectByID(objType, nil, id, origindata, ""); err != nil {
			blog.Error("create event error:%v", err)
		} else {
			ec := eventdata.NewEventContextByReq(req)
			err := ec.InsertEvent(eventtypes.EventTypeInstData, objType, eventtypes.EventActionCreate, origindata, nil, ownerID)
			if err != nil {
				blog.Error("create event error:%v", err)
			}
		}

		info := make(map[string]int)
		info[idName] = id
		return http.StatusOK, info, nil
	}, resp)
}
