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

package privilege

import (
	"configcenter/src/common"
	"configcenter/src/common/bkbase"
	"configcenter/src/common/blog"
	"configcenter/src/common/core/cc/actions"
	"configcenter/src/common/core/cc/api"
	httpcli "configcenter/src/common/http/httpclient"
	"configcenter/src/common/paraparse"
	"configcenter/src/common/util"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/emicklei/go-restful"
)

var privilege *privilegeAction = &privilegeAction{}

type privilegeAction struct {
	base.BaseAction
}

type GroupListResult struct {
	Result  bool                     `json:"result"`
	Code    int                      `json:"code"`
	Message interface{}              `json:"message"`
	Data    []map[string]interface{} `json:"data"`
}

type modelConfig struct {
}

func init() {

	actions.RegisterNewAction(actions.Action{Verb: common.HTTPCreate, Path: "/privilege/group/detail/{bk_supplier_account}/{group_id}", Params: nil, Handler: group.UpdateUserGroupPrivi})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectGet, Path: "/privilege/group/detail/{bk_supplier_account}/{group_id}", Params: nil, Handler: group.GetUserGroupPrivi})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectGet, Path: "/privilege/user/detail/{bk_supplier_account}/{user_name}", Params: nil, Handler: group.GetUserPrivi})
	privilege.CreateAction()
}

//UpdateUserGroupPrivi create user group
func (cli *groupAction) UpdateUserGroupPrivi(req *restful.Request, resp *restful.Response) {
	// get the language
	language := util.GetActionLanguage(req)
	// get the error factory by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)

	cli.CallResponseEx(func() (int, interface{}, error) {

		pathParams := req.PathParameters()
		ownerID, _ := pathParams["bk_supplier_account"]
		groupID, _ := pathParams["group_id"]

		value, err := ioutil.ReadAll(req.Request.Body)
		if err != nil {
			blog.Error("read json data error :%v", err)
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommHTTPReadBodyFailed)
		}

		//get user group privilege url
		groupURL := cli.CC.ObjCtrl() + "/object/v1/privilege/group/detail/" + ownerID + "/" + groupID
		blog.Info("get user group privilege url: %s", groupURL)
		groupInfo, err := httpcli.ReqHttp(req, groupURL, common.HTTPSelectGet, value)
		if nil != err {
			blog.Error("get user group privilege error :%v", err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrTopoUserGroupPrivilegeUpdateFailed)
		}
		blog.Info("get user group return: %s", groupInfo)
		var groupPrivi params.GroupPriviResult
		err = json.Unmarshal([]byte(groupInfo), &groupPrivi)
		if nil == err && false == groupPrivi.Result {
			//create group privilege
			blog.Info("create user group privilege content: %s", string(value))
			blog.Info("create user group privilege url: %s", groupURL)
			groupInfo, err = httpcli.ReqHttp(req, groupURL, common.HTTPCreate, value)
			blog.Info("create user group privilege return: %s", groupInfo)
			if nil != err {
				blog.Error("create user group privilege error: %v", err)
			}

			return http.StatusOK, groupInfo, nil
		}
		//update group privilege
		blog.Info("update user group privilege content: %s", string(value))
		blog.Info("update user group privilege url: %s", groupURL)
		groupInfo, err = httpcli.ReqHttp(req, groupURL, common.HTTPUpdate, value)
		if nil != err {
			blog.Error("update user group privilege error :%v", err)
			blog.Info("update user group privilege url: %s", groupInfo)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrTopoUserGroupPrivilegeUpdateFailed)
		}

		return http.StatusOK, groupInfo, nil
	}, resp)
}

//GetUserGroup get user group
func (cli *groupAction) GetUserGroupPrivi(req *restful.Request, resp *restful.Response) {
	// get the language
	language := util.GetActionLanguage(req)
	// get the error factory by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)

	cli.CallResponseEx(func() (int, interface{}, error) {
		pathParams := req.PathParameters()
		ownerID, _ := pathParams["bk_supplier_account"]
		groupID, _ := pathParams["group_id"]

		//get user group privilege url
		groupURL := cli.CC.ObjCtrl() + "/object/v1/privilege/group/detail/" + ownerID + "/" + groupID
		blog.Info("get user group privilege url: %s", groupURL)
		groupInfo, err := httpcli.ReqHttp(req, groupURL, common.HTTPSelectGet, nil)
		if nil != err {
			blog.Error("get user group privilege error :%v", err)
			cli.ResponseFailed(common.CC_Err_Comm_CREATE_USER_GROUP_FAIL, common.CC_Err_Comm_CREATE_USER_GROUP_FAIL_STR, resp)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrTopoUserGroupPrivilegeSelectFailed)
		}
		blog.Info("get user group privilege return: %s", groupInfo)
		var groupPrivi params.GroupPriviResult
		err = json.Unmarshal([]byte(groupInfo), &groupPrivi)
		if err != nil || false == groupPrivi.Result {
			data := make(map[string]interface{})
			data[common.BKOwnerIDField] = ownerID
			data[common.BKUserGroupIDField] = groupID
			data[common.BKPrivilegeField] = common.KvMap{}
			return http.StatusOK, data, nil
		}
		return http.StatusOK, groupPrivi.Data, nil
	}, resp)
}

func (cli *groupAction) GetUserPrivi(req *restful.Request, resp *restful.Response) {
	// get the language
	language := util.GetActionLanguage(req)
	// get the error factory by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)

	cli.CallResponseEx(func() (int, interface{}, error) {
		//get group by user
		pathParams := req.PathParameters()
		ownerID, _ := pathParams["bk_supplier_account"]
		userName, _ := pathParams["user_name"]
		cond := make(map[string]interface{})
		userNameMap := make(map[string]interface{})
		userNameMap[common.BKDBLIKE] = userName
		cond[common.BKUserListField] = userNameMap
		var gPrivilege params.Gprivilege

		//get cross biz privilege
		isHostCrossBiz := false
		url := cli.CC.ObjCtrl() + "/object/v1/system/" + common.HostCrossBizField + "/" + ownerID
		blog.Info("get system config url :%s", url)
		reply, err := httpcli.ReqHttp(req, url, common.HTTPSelectGet, nil)
		if err != nil {
			blog.Error("get system config : %v", err)
			isHostCrossBiz = false
		}
		blog.Info("get system config return :%s", string(reply))

		var result api.APIRsp
		err = json.Unmarshal([]byte(reply), &result)
		if nil != err {
			blog.Error("get system config error : %v", err)
			isHostCrossBiz = false
		}
		if result.Result {
			isHostCrossBiz = true
		}
		gPrivilege.IsHostCrossBiz = isHostCrossBiz

		gPrivilege.ModelConfig = make(map[string]map[string][]string)
		jsonStr, _ := json.Marshal(cond)
		groupURL := cli.CC.ObjCtrl() + "/object/v1/privilege/group/" + ownerID + "/" + "search"

		blog.Info("search user group url: %s", groupURL)
		blog.Info("search user group info: %s", jsonStr)
		searchInfo, err := httpcli.ReqHttp(req, groupURL, common.HTTPSelectPost, jsonStr)
		blog.Info("search user group result: %s", searchInfo)
		if nil != err {
			blog.Error("search user group error :%v", err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrTopoUserPrivilegeSelectFailed)
		}
		var data GroupListResult
		err = json.Unmarshal([]byte(searchInfo), &data)
		if nil != err || !data.Result {
			return http.StatusOK, gPrivilege, nil
		}
		var groupIDArr []string
		for _, i := range data.Data {
			id, ok := i[common.BKUserGroupIDField]
			if false == ok {
				continue
			}
			idStr, ok := id.(string)
			if false == ok {
				continue
			}
			userNames, ok := i[common.BKUserListField]
			if false == ok {
				continue
			}
			userNameStr, ok := userNames.(string)
			if false == ok {
				continue
			}
			userNameList := strings.Split(userNameStr, ";")
			if util.InArray(userName, userNameList) {
				groupIDArr = append(groupIDArr, idStr)
			}

		}
		//get group pri
		ugroupIDArr := util.ArrayUnique(groupIDArr)
		var gglconfig []string
		var gbkconfig []string
		var modelCls []string
		modelPrivi := make(map[string][]string)
		modelClsConfig := make(map[string]string)
		for _, i := range ugroupIDArr {
			//get user group privilege url
			groupID := i.(string)
			groupURL := cli.CC.ObjCtrl() + "/object/v1/privilege/group/detail/" + ownerID + "/" + groupID
			blog.Info("get user group privilege url: %s", groupURL)
			groupInfo, err := httpcli.ReqHttp(req, groupURL, common.HTTPSelectGet, nil)
			blog.Info("get user group privilege result: %v", groupInfo)
			if nil != err {
				continue
			}
			var groupPrivi params.GroupPriviResult
			err = json.Unmarshal([]byte(groupInfo), &groupPrivi)
			if nil != err || !groupPrivi.Result {
				continue
			}

			if nil != groupPrivi.Data.Privilege.SysConfig {
				sysConfig := *groupPrivi.Data.Privilege.SysConfig
				for _, i := range sysConfig.Globalbusi {
					gglconfig = append(gglconfig, i)
				}
				for _, j := range sysConfig.BackConfig {
					gbkconfig = append(gbkconfig, j)
				}
			}

			for m, n := range groupPrivi.Data.Privilege.ModelConfig {
				for i, j := range n {
					for _, k := range j {
						modelPrivi[i] = append(modelPrivi[i], k)
					}
					modelClsConfig[i] = m
				}
				modelCls = append(modelCls, m)

			}
		}
		umodelCls := util.RemoveDuplicatesAndEmpty(modelCls)
		cls := make(map[string]map[string][]string)
		for _, i := range umodelCls {
			modelCls := make(map[string][]string)
			for j, k := range modelPrivi {
				if modelClsConfig[j] == i {
					modelCls[j] = k
				}
			}
			cls[i] = modelCls
		}
		gPrivilege.SysConfig.BackConfig = util.RemoveDuplicatesAndEmpty(gbkconfig)
		gPrivilege.SysConfig.Globalbusi = util.RemoveDuplicatesAndEmpty(gglconfig)
		gPrivilege.ModelConfig = cls
		return http.StatusOK, gPrivilege, nil
	}, resp)
}
