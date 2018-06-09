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
	"configcenter/src/common/bkbase"
	"configcenter/src/common/blog"
	"configcenter/src/common/core/cc/actions"
	httpcli "configcenter/src/common/http/httpclient"
	"encoding/json"
	"io/ioutil"

	simplejson "github.com/bitly/go-simplejson"
	"github.com/emicklei/go-restful"
)

var app = &appAction{}

type appAction struct {
	base.BaseAction
}

func init() {

	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "/app/searchAll", Params: nil, Handler: app.SearchAllApp})
	//actions.RegisterNewAction(actions.Action{common.HTTPCreate, "/openapi/app/addApp", nil, app.AddApp,nil})
	//actions.RegisterNewAction(actions.Action{common.HTTPUpdate, "/openapi/app/editApp", nil, app.EditApp,nil})
	//actions.RegisterNewAction(actions.Action{common.HTTPDelete, "/openapi/app/deleteApp", nil, app.DeleteApp,nil})

	// create CC object
	app.CreateAction()
}

//SearchAllApp: search all application
func (cli *appAction) SearchAllApp(req *restful.Request, resp *restful.Response) {
	blog.Debug("SearchAllApp start")
	value, _ := ioutil.ReadAll(req.Request.Body)
	js, err := simplejson.NewJson([]byte(value))
	if err != nil {
		blog.Error("get all app failed, err msg : %v", err)
		cli.ResponseFailed(common.CC_Err_Comm_http_Input_Params, common.CC_Err_Comm_http_Input_Params_STR, resp)
		return
	}

	rq_para, err := js.Map()
	if err != nil {
		blog.Error("get all app failed, err msg : %v", err)
		cli.ResponseFailed(common.CC_Err_Comm_http_Input_Params, common.CC_Err_Comm_http_Input_Params_STR, resp)
		return
	}

	searchParams := make(map[string]interface{})
	searchParams["condition"] = rq_para["condition"]
	searchParams["fields"] = rq_para["fields"].(string)
	inputJson, _ := json.Marshal(searchParams)

	sAppURL := cli.CC.ObjCtrl() + "/object/v1/insts/" + common.BKInnerObjIDApp + "/search"
	appInfo, err := httpcli.ReqHttp(req, sAppURL, common.HTTPSelectPost, []byte(inputJson))
	if nil != err {
		blog.Error("search all app failed, err msg : %v", err)
		cli.ResponseFailed(common.CC_Err_Comm_APP_QUERY_FAIL, common.CC_Err_Comm_APP_QUERY_FAIL_STR, resp)
		return
	}

	json, err := simplejson.NewJson([]byte(appInfo))
	appResData, _ := json.Map()
	cli.ResponseSuccess(appResData["data"], resp)

}

/*
// AddApp: 新增业务
func (cli *appAction) AddApp(req *restful.Request, resp *restful.Response) {
	value, _ := ioutil.ReadAll(req.Request.Body)

	js, err := simplejson.NewJson([]byte(value))
	if err != nil {
		blog.Error("add app failed, err msg : %v", err)
		cli.ResponseFailed(common.CC_Err_Comm_http_Input_Params, common.CC_Err_Comm_http_Input_Params_STR, resp)
		return
	}
	rq_para, err := js.Map()
	blog.Debug("add app rq_para:%v",rq_para)
	if err != nil {
		blog.Error("add app failed, err msg : %v", err)
		cli.ResponseFailed(common.CC_Err_Comm_http_Input_Params, common.CC_Err_Comm_http_Input_Params_STR, resp)
		return
	}

	const BK_DEFAULT_LEVER  = 3
	const BK_DEFAULT = 0
	const BK_LANGUAGE = "中文"
	Params := make(map[string]interface{})
	Params["OwnerID"] = common.BKDefaultOwnerID
	Params["Language"] = BK_LANGUAGE
	Params["Level"] = BK_DEFAULT_LEVER
	Params["Default"] = BK_DEFAULT
	Params["ApplicationName"] = rq_para["ApplicationName"]
	Params["Maintainers"] = rq_para["Maintainers"]
	Params["Creator"] = rq_para["Creator"]
	Params["LifeCycle"] = rq_para["LifeCycle"]
	Params["ProductPm"] = rq_para["ProductPm"]
	Params["Developer"] = rq_para["Developer"]
	Params["Tester"] = rq_para["Tester"]
	Params["Operator"] = rq_para["Operator"]
	inputJson, _ := json.Marshal(Params)

	sAppUrl :=  cli.CC.ObjCtrl() + "/object/v1/insts/"+common.BKInnerObjIDApp
	res, err := httpcli.ReqHttp(req, sAppUrl, common.HTTPCreate, []byte(inputJson))
	blog.Debug("add app res:%v",res)
	if nil != err {
		blog.Error("search all app failed, err msg : %v", err)
		cli.ResponseFailed(common.CC_Err_Comm_APP_QUERY_FAIL, common.CC_Err_Comm_APP_QUERY_FAIL_STR, resp)
		return
	}
	// deal result
	var rst api.APIRsp
	if jserr := json.Unmarshal([]byte(res), &rst); nil == jserr {
		cli.Response(&rst, resp)
		return
	} else {
		blog.Error("unmarshal the json failed, error information is %v", jserr)
		cli.ResponseFailed(common.CC_Err_Comm_CREATE_PLAT_FAIL, common.CC_Err_Comm_CREATE_PLAT_FAIL_STR, resp)
	}
}

//EditApp 修改业务
func (cli *appAction) EditApp(req *restful.Request, resp *restful.Response) {
	value, _ := ioutil.ReadAll(req.Request.Body)
	js, err := simplejson.NewJson([]byte(value))
	if err != nil {
		blog.Error("add app failed, err msg : %v", err)
		cli.ResponseFailed(common.CC_Err_Comm_http_Input_Params, common.CC_Err_Comm_http_Input_Params_STR, resp)
		return
	}

	rq_para, err := js.Map()
	if err != nil {
		blog.Error("edit failed, err msg : %v", err)
		cli.ResponseFailed(common.CC_Err_Comm_http_Input_Params, common.CC_Err_Comm_http_Input_Params_STR, resp)
		return
	}
	appId := rq_para["ApplicationID"]

	condition := map[string]interface{}{
		"ApplicationID": appId,
	}
	appMap,err:= GetOneApp(req,cli.CC.ObjCtrl(),condition)
	if len(appMap)==0{
		msg := fmt.Sprintf("not find app in ApplicationID:%v",appId)
		blog.Error("edit app failed, err msg : %v", msg)
		cli.ResponseFailed(common.CC_Err_Comm_APP_ID_ERR, msg, resp)
		return
	}
	blog.Debug("appMap:%v",appMap)
	for key,value := range rq_para{
		appMap[key] = value
	}
	blog.Debug("appMap:%v",appMap)
	params := map[string]interface{}{}
	params["condition"] = condition
	params["data"] = appMap
	inputJson, _ := json.Marshal(params)
	sAppUrl := cli.CC.ObjCtrl() + "/object/v1/insts/"+common.BKInnerObjIDApp
	res, err := httpcli.ReqHttp(req, sAppUrl, common.HTTPUpdate, []byte(inputJson))
	if nil != err {
		blog.Error("edit app failed, err msg : %v", err)
		cli.ResponseFailed(common.CC_Err_Comm_APP_QUERY_FAIL, common.CC_Err_Comm_APP_QUERY_FAIL_STR, resp)
		return
	}
	blog.Debug("edit app res:%v",res)

	cli.ResponseSuccess("", resp)

}

//DeleteApp: 删除业务
func (cli *appAction) DeleteApp(req *restful.Request, resp *restful.Response) {
	value, _ := ioutil.ReadAll(req.Request.Body)
	js, err := simplejson.NewJson([]byte(value))
	if err != nil {
		blog.Error("delete app failed, err msg : %v", err)
		cli.ResponseFailed(common.CC_Err_Comm_http_Input_Params, common.CC_Err_Comm_http_Input_Params_STR, resp)
		return
	}
	input, err := js.Map()
	if err != nil {
		blog.Error("delete app failed, err msg : %v", err)
		cli.ResponseFailed(common.CC_Err_Comm_http_Input_Params, common.CC_Err_Comm_http_Input_Params_STR, resp)
		return
	}

	appIdInt,_:=strconv.Atoi(input["appId"].(string))
	condition := make(map[string]interface{})
	condition["ApplicationID"] = appIdInt
	fileds := "ApplicationID,ApplicationName,Default"
	appMap,err := GetAppMapByCond(req,fileds,cli.CC.ObjCtrl(),condition)
	if nil != err {
		blog.Debug("GetAppMapByCond:%v",err)
		cli.ResponseFailed(common.CC_Err_Comm_APP_DEL_FAIL,common.CC_Err_Comm_APP_DEL_FAIL_STR , resp)
		return
	}
	blog.Debug("appMap:%v",appMap)
	if len(appMap) == 0 {
		cli.ResponseFailed(common.CC_Err_Comm_APP_DEL_FAIL, "ApplicationID not found" , resp)
		return
	}

	Default := appMap[appIdInt].(map[string]interface{})["Default"]
	if Default == "1" {
		cli.ResponseFailed(common.CC_Err_Comm_APP_DEL_FAIL, "该业务不能删除" , resp)
		return
	}

	const Not_DETELE_APPLICATIONBASE = "资源池"
	notDelArr := strings.Split(Not_DETELE_APPLICATIONBASE,",")
	appName := appMap[appIdInt].(map[string]interface{})["ApplicationName"]
	for _,value :=range notDelArr{
		if appName == value{
			cli.ResponseFailed(common.CC_Err_Comm_APP_DEL_FAIL, "该业务不能删除" , resp)
			return
		}
	}
	appIdArr := []int{appIdInt}
	hostIdArr,err:= GetHostIdByCond(req,cli.CC.HostCtrl(),map[string]interface{}{
		"ApplicationID": appIdArr,
	})
	if nil != err {
		blog.Error("GetAppMapByCond:%v",err)
		cli.ResponseFailed(common.CC_Err_Comm_APP_DEL_FAIL,common.CC_Err_Comm_APP_DEL_FAIL_STR , resp)
		return
	}
	if len(hostIdArr) >0 {
		cli.ResponseFailed(common.CC_Err_Comm_APP_DEL_FAIL, "业务下存在主机，不能删除" , resp)
		return
	}

	inputJson, _ := json.Marshal(condition)
	url := "http://" + cli.CC.ObjCtrl() + "/object/v1/insts/"+common.BKInnerObjIDApp+"/"
	res, err := httpcli.ReqHttp(req, url, common.HTTPDelete, []byte(inputJson))
	if nil != err {
		blog.Error("delete app failed, err msg : %v", err)
		cli.ResponseFailed(common.CC_Err_Comm_APP_DEL_FAIL, common.CC_Err_Comm_APP_DEL_FAIL_STR, resp)
		return
	}
	//deal result
	var rst api.APIRsp
	if jserr := json.Unmarshal([]byte(res), &rst); nil == jserr {
		cli.Response(&rst, resp)
		return
	} else {
		blog.Error("unmarshal the json failed, error information is %v", jserr)
	}
}


//helper
//get appmap by cond
func GetAppMapByCond(req *restful.Request, fields string, objUrl string, cond interface{}) (map[int]interface{}, error) {
	appMap := make(map[int]interface{})
	condition := make(map[string]interface{})
	condition["fields"] = fields
	condition["sort"] = "ApplicationID"
	condition["start"] = 0
	condition["limit"] = 0
	condition["condition"] = cond
	bodyContent, _ := json.Marshal(condition)
	url := "http://" + objUrl + "/object/v1/insts/"+common.BKInnerObjIDApp+"/search"
	blog.Info("GetAppMapByCond url :%s", url)
	blog.Info("GetAppMapByCond content :%s", string(bodyContent))
	reply, err := httpcli.ReqHttp(req, url, common.HTTPSelectPost, []byte(bodyContent))
	blog.Info("GetAppMapByCond return :%s", string(reply))
	if err != nil {
		return appMap, err
	}
	js, _ := simplejson.NewJson([]byte(reply))
	output, _ := js.Map()
	appData := output["data"]
	appResult := appData.(map[string]interface{})
	appInfo := appResult["info"].([]interface{})
	for _, i := range appInfo {
		app := i.(map[string]interface{})
		appId, _ := app["ApplicationID"].(json.Number).Int64()
		appMap[int(appId)] = i
	}
	return appMap, nil
}

//get module host config
func GetHostIdByCond(req *restful.Request, hostUrl string, cond interface{}) ([]int, error) {
	hostIdArr := make([]int, 0)
	bodyContent, _ := json.Marshal(cond)
	url := "http://" + hostUrl + "/host/v1/meta/hosts/module/config/search"
	blog.Info("GetUserConfig ModuleHostConfig url :%s", url)
	blog.Info("GetUserConfig ModuleHostConfig content :%s", string(bodyContent))
	reply, err := httpcli.ReqHttp(req, url, common.HTTPSelectPost, []byte(bodyContent))
	blog.Info("GetUserConfig ModuleHostConfig return :%s", string(reply))
	if err != nil {
		return hostIdArr, err
	}
	js, _ := simplejson.NewJson([]byte(reply))
	output, _ := js.Map()
	configData := output["data"]
	configInfo, ok := configData.([]interface{})
	if !ok {
		return hostIdArr, nil
	}
	for _, i := range configInfo {
		host := i.(map[string]interface{})
		hostId, _ := host["HostID"].(json.Number).Int64()
		hostIdArr = append(hostIdArr, int(hostId))
	}
	return hostIdArr, err
}


func GetOneApp(req *restful.Request, objUrl string, cond interface{}) (map[string]interface{}, error) {
	condition := make(map[string]interface{})
	condition["sort"] = "ApplicationID"
	condition["start"] = 0
	condition["limit"] = 1
	condition["condition"] = cond
	bodyContent, _ := json.Marshal(condition)
	url := "http://" + objUrl + "/object/v1/insts/"+common.BKInnerObjIDApp+"/search"
	blog.Info("GetOneApp url :%s", url)
	blog.Info("GetOneApp content :%s", string(bodyContent))

	reply, err := httpcli.ReqHttp(req, url, common.HTTPSelectPost, []byte(bodyContent))
	blog.Info("GetOneApp return :%s", string(reply))
	if err != nil {
		blog.Info("GetOneApp return http request error:%s", string(reply))
		return nil, err
	}
	js, _ := simplejson.NewJson([]byte(reply))
	output, _ := js.Map()
	appData := output["data"]
	appResult := appData.(map[string]interface{})
	appInfo := appResult["info"].([]interface{})
	for _, i := range appInfo {
		app, _ := i.(map[string]interface{})
		return app, nil
	}
	return nil, nil
}
*/
