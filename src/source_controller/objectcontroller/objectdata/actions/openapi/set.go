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

package openapi

import (
	"configcenter/src/common"
	"configcenter/src/common/base"
	"configcenter/src/common/blog"
	"configcenter/src/common/core/cc/actions"
	"configcenter/src/common/util"
	eventtypes "configcenter/src/scene_server/event_server/types"
	"configcenter/src/source_controller/common/eventdata"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"gopkg.in/mgo.v2/bson"

	"github.com/emicklei/go-restful"
)

var set *setAction = &setAction{}

type setAction struct {
	base.BaseAction
}

func init() {

	actions.RegisterNewAction(actions.Action{Verb: common.HTTPDelete, Path: "/openapi/set/delhost", Params: nil, Handler: set.DeleteSetHost})

	// create CC object
	set.CreateAction()
}

func (cli *setAction) DeleteSetHost(req *restful.Request, resp *restful.Response) {
	// get the language
	language := util.GetActionLanguage(req)
	ownerID := util.GetActionOnwerID(req)
	// get the error factory by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)

	cli.CallResponseEx(func() (int, interface{}, error) {

		blog.Debug("DeleteSetHost start !")
		value, err := ioutil.ReadAll(req.Request.Body)
		if nil != err {
			blog.Error("read request body failed, error:%v", err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrCommHTTPReadBodyFailed)
		}

		blog.Debug("DeleteSetHost http body data: %s", value)
		input := make(map[string]interface{})
		err = json.Unmarshal(value, &input)
		if nil != err {
			blog.Error("unmarshal json error:%v", err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrCommJSONUnmarshalFailed)
		}
		input = util.SetModOwner(input, ownerID)

		err = delModuleConfigSet(input, ownerID, req)
		if err != nil {
			blog.Error("fail to delSetConfigHost: %v", err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrCommParamsInvalid)

		}

		return http.StatusOK, nil, nil
	}, resp)
}

// TODO
func getModuleConfigCount(con map[string]interface{}) (int, error) {
	count, err := set.CC.InstCli.GetCntByCondition("cc_ModuleHostConfig", con)
	if err != nil {
		blog.Error("fail getModuleConfigCount error:%v", err)
		return 0, err
	}
	return count, err
}

func delModuleConfigSet(input map[string]interface{}, ownerID string, req *restful.Request) error {
	tableName := "cc_ModuleHostConfig"

	appID, ok := input[common.BKAppIDField]
	if false == ok {
		blog.Errorf("params ApplicationID is required")
		return errors.New("params ApplicationID is required")
	}
	var oldContents []interface{}
	getErr := set.CC.InstCli.GetMutilByCondition(tableName, nil, input, &oldContents, "", 0, common.BKNoLimit)
	if getErr != nil {
		blog.Errorf("fail to delSetConfigHost: %v", getErr)
		return getErr
	}

	setID, moduleID, defErr := GetIdleModule(appID, ownerID)
	if nil != defErr {
		blog.Errorf("get idle module error:%v", defErr)
		return defErr
	}

	err := set.CC.InstCli.DelByCondition(tableName, input)
	if err != nil {
		blog.Error("fail to delSetConfigHost: %v", err)
		return err
	}
	//发送删除主机关系事件
	ec := eventdata.NewEventContextByReq(req)
	for oldContent := range oldContents {
		err = ec.InsertEvent(eventtypes.EventTypeRelation, common.BKInnerObjIDHost, eventtypes.EventActionDelete, oldContent, nil, ownerID)
		if err != nil {
			blog.Error("create event error:%v", err)
		}
	}

	var hostIDs []interface{}    //all hostid
	mapHostIDs := common.KvMap{} //distinct hostid
	for _, item := range oldContents {
		mapItem, _ := item.(bson.M)

		hostIDs = append(hostIDs, mapItem[common.BKHostIDField])
		mapHostIDs[fmt.Sprintf("%v", mapItem[common.BKHostIDField])] = mapItem[common.BKHostIDField]
	}
	//del host from set, get host module relation
	params := common.KvMap{common.BKAppIDField: appID, common.BKHostIDField: common.KvMap{"$in": hostIDs}}
	params = util.SetModOwner(params, ownerID)
	var hostRelations []interface{}
	getErr = set.CC.InstCli.GetMutilByCondition(tableName, nil, params, &hostRelations, "", 0, common.BKNoLimit)
	if getErr != nil {
		blog.Error("fail to exist relation host error: %v", getErr)
		return getErr
	}

	existRelationHostID := common.KvMap{}
	for _, item := range hostRelations {
		mapItem, _ := item.(bson.M)
		existRelationHostID[fmt.Sprintf("%v", mapItem[common.BKHostIDField])] = 1
	}

	//get host not moodule
	var addIdleModuleDatas []interface{}
	setID, _ = util.GetIntByInterface(setID)
	moduleID, _ = util.GetIntByInterface(moduleID)
	for strHostID, rawHostID := range mapHostIDs {
		_, ok := existRelationHostID[strHostID]
		if !ok {
			param := map[string]interface{}{common.BKAppIDField: appID, common.BKSetIDField: setID, common.BKModuleIDField: moduleID, common.BKHostIDField: rawHostID}
			//set.CC.InstCli.InsertMuti
			param = util.SetModOwner(param, ownerID)
			addIdleModuleDatas = append(addIdleModuleDatas, param)
		}

	}
	if 0 < len(addIdleModuleDatas) {
		err := set.CC.InstCli.InsertMuti(tableName, addIdleModuleDatas...)
		if getErr != nil {
			blog.Error("fail to exist relation host error: %v", err)
			return err
		}
		//推送新加到空闲机器的关系
		for _, row := range addIdleModuleDatas {
			err = ec.InsertEvent(eventtypes.EventTypeRelation, common.BKInnerObjIDHost, eventtypes.EventActionCreate, nil, row, ownerID)
			if err != nil {
				blog.Error("create event error:%v", err)
			}
		}

	}

	return nil
}

func GetIdleModule(appID interface{}, ownerID string) (interface{}, interface{}, error) {
	params := common.KvMap{common.BKAppIDField: appID, common.BKDefaultField: common.DefaultResModuleFlag, common.BKModuleNameField: common.DefaultResModuleName}
	params = util.SetModOwner(params, ownerID)
	var result bson.M
	err := set.CC.InstCli.GetOneByCondition("cc_ModuleBase", []string{common.BKModuleIDField, common.BKSetIDField}, params, &result)

	if nil != err {
		return nil, nil, err
	}
	return result[common.BKSetIDField], result[common.BKModuleIDField], nil
}
