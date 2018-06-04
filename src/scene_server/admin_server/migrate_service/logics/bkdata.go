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
 
package logics

import (
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/core/cc/api"
	httpcli "configcenter/src/common/http/httpclient"
	"configcenter/src/common/util"
	"encoding/json"
	"errors"
	"strconv"
	"strings"

	simplejson "github.com/bitly/go-simplejson"
	restful "github.com/emicklei/go-restful"
)

// 进程:功能:port
var prc2port = []string{"job_java:java:8008,8443", "uwsgi:uwsgi:8000,8001,8003,8004",
	"paas_agent:paas_agent:4245", "common_mysql:mysqld:3306", "common_redis:redis-server:6379", "redis_cluster:redis-server:16379,6379",
	"zk_java:java:2181", "kafka_java:java:9092", "es_java:java:9300,10004", "beam.smp:beam.smp:15672,5672,25672",
	"common_nginx:nginx:80", "beanstalkd:beanstalkd:6380", "influxdb:influxdb:5620,5621", "etcd:etcd:2379,2380", "dataapi_py:python:10011", "databus_java:java:10021", "monitor_py:python:10043",
	"fta_py:python:13031,13021,13041", "gse_data:gse_data:58625", "gse_dba:gse_dba:58859",
	"gse_api:gse_api:59313,50002", "gse_task:gse_task:48668,48669,48671,48329", "gse_transit:gse_transit:58625", "gse_proc:gse_proc:52023,52025",
	"gse_btsvr:gse_btsvr:10020,58930,58925", "gse_ops:gse_ops:58725,58726", "gse_alarm:gse_alarm:53425", "gse_agent:gse_agent:60020,34334,36510",
	"license_server:license_server:443", "consul:consul:8301,8300,8302,8500,53",
	"cmdb_adminserver:adminserver:32004",
	"cmdb_apiserver:apiserver:33031",
	"cmdb_auditcontroller:auditcontroller:31004",
	"cmdb_datacollection:datacollection:32006",
	"cmdb_eventserver:eventserver:32005",
	"cmdb_hostcontroller:hostcontroller:31002",
	"cmdb_hostserver:hostserver:32001",
	"cmdb_objectcontroller:objectcontroller:31001",
	"cmdb_proccontroller:proccontroller:31003",
	"cmdb_procserver:procserver:32003",
	"cmdb_toposerver:toposerver:32002",
	"cmdb_webserver:webserver:33083"}

// 集群:模块:进程
var setModuleKv = map[string]map[string]string{"作业平台": {"job": "job_java"},
	"配置平台": {
		"adminserver":      "cmdb_adminserver",
		"apiserver":        "cmdb_apiserver",
		"auditcontroller":  "cmdb_auditcontroller",
		"datacollection":   "cmdb_datacollection",
		"eventserver":      "cmdb_eventserver",
		"hostcontroller":   "cmdb_hostcontroller",
		"hostserver":       "cmdb_hostserver",
		"objectcontroller": "cmdb_objectcontroller",
		"proccontroller":   "cmdb_proccontroller",
		"procserver":       "cmdb_procserver",
		"toposerver":       "cmdb_toposerver",
		"webserver":        "cmdb_webserver",
	},
	"管控平台":   {"gse_api": "gse_api", "gse_data": "gse_data", "gse_dba": "gse_dba", "gse_task": "gse_task", "gse_transit": "gse_transit", "gse_proc": "gse_proc", "gse_btsvr": "gse_btsvr", "gse_ops": "gse_ops", "gse_opts": "", "gse_alarm": "gse_alarm", "gse_agent": "gse_agent", "license": "license_server"},
	"故障自愈":   {"fta": "fta_py"},
	"数据服务模块": {"dataapi": "dataapi_py", "databus": "databus_java", "monitor": "monitor_py"},
	"公共组件": {"mysql": "common_mysql", "redis": "common_redis", "redis_cluster": "redis_cluster", "zookeeper": "zk_java", "kafka": "kafka_java", "elasticsearch": "es_java",
		"rabbitmq": "beam.smp", "nginx": "common_nginx", "beanstalk": "beanstalkd", "influxdb": "influxdb", "etcd": "etcd", "consul": "consul", "mongodb": "mongodb"},
	"集成平台": {"esb": "uwsgi", "login": "uwsgi", "paas": "uwsgi", "appengine": "uwsgi", "console": "uwsgi", "appo": "paas_agent", "appt": "paas_agent"},
}

var appID int = 0
var ownerID string = common.BKDefaultOwnerID
var procAPI, topoAPI string
var procName2ID map[string]int
var appModelData map[string]interface{}
var setModelData map[string]interface{}
var moduleModelData map[string]interface{}
var procModelData map[string]interface{}

//BKAppInit  init bk app
func BKAppInit(req *restful.Request, cc *api.APIResource, ownerID string) error {
	var err error
	//get api addr
	procAPI = cc.ProcAPI()
	topoAPI = cc.TopoAPI()
	//get model module
	procModelData, err = getObjectFields(cc.TopoAPI(), req, common.BKInnerObjIDProc)
	if err != nil {
		blog.Error("get procModelData err :%v ", err)
		return err
	}
	appModelData, err = getObjectFields(cc.TopoAPI(), req, common.BKInnerObjIDApp)
	if err != nil {
		blog.Error("get appModelData err :%v ", err)
		return err
	}
	setModelData, err = getObjectFields(cc.TopoAPI(), req, common.BKInnerObjIDSet)
	if err != nil {
		blog.Error("get setModelData err :%v ", err)
		return err
	}
	moduleModelData, err = getObjectFields(cc.TopoAPI(), req, common.BKInnerObjIDModule)
	if err != nil {
		blog.Error("get moduleModelData err :%v ", err)
		return err
	}

	isExist, err := BKAppIsExist(req)
	if nil != err {
		blog.Error("get app isExist err :%v ", err)
		return err
	}

	if !isExist {
		err = addBKApp(req)
		if nil != err {
			blog.Error("add bk app err :%v ", err)
			return err
		}

		err = addBKProcess(req)
		if nil != err {
			blog.Error("add bk process err :%v ", err)
			return err
		}
	}
	return nil

}

//addBKApp add bk app
func addBKApp(req *restful.Request) error {
	appModelData[common.BKAppNameField] = common.BKAppName
	appModelData[common.BKMaintainersField] = "admin"
	appModelData[common.BKTimeZoneField] = "Asia/Shanghai"
	appModelData[common.BKLanguageField] = "1" //"中文"
	appModelData[common.BKLifeCycleField] = common.DefaultAppLifeCycleNormal

	byteParams, _ := json.Marshal(appModelData)
	url := topoAPI + "/topo/v1/app/" + ownerID
	blog.Info("migrate add bk app url :%s", url)
	blog.Info("migrate add bk app content :%s", string(byteParams))
	reply, err := httpcli.ReqHttp(req, url, common.HTTPCreate, byteParams)
	blog.Info("migrate add bk app return :%s", string(reply))
	if err != nil {
		blog.Error("add bk app err :%v ", err)
		return err
	}
	js, _ := simplejson.NewJson([]byte(reply))
	output, _ := js.Map()

	code, err := util.GetIntByInterface(output["bk_error_code"])
	if err != nil || 0 != code {
		if strings.Contains(reply, "duplicate") || strings.Contains(reply, "重复") {
			return nil
		}
		blog.Error("add bk app json err :%v ", err)
		return errors.New(reply)
	}
	data, ok := output["data"].(map[string]interface{})
	if false == ok {
		blog.Error("add bk app result data err :%v ", err)
		return errors.New("get appID error")
	}
	appID, err = util.GetIntByInterface(data[common.BKAppIDField])
	if nil != err {
		blog.Error("add bk app result data app id err :%v ", err)
		return err
	}
	return nil
}

//addBKProcess add bk process
func addBKProcess(req *restful.Request) error {
	procName2ID = make(map[string]int)
	appIDStr := strconv.Itoa(appID)

	for _, procStr := range prc2port {
		procArr := strings.Split(procStr, ":")
		procName := procArr[0]
		funcName := procArr[1]
		portStr := procArr[2]
		var protocal string
		if len(procArr) > 3 {
			protocal = procArr[3]
		}
		procModelData[common.BKProcNameField] = procName
		procModelData[common.BKFuncName] = funcName
		procModelData[common.BKPort] = portStr
		procModelData[common.BKWorkPath] = "/data/bkee"

		protocal = strings.ToLower(protocal)
		switch protocal {
		case "udp":
			procModelData[common.BKProtocol] = "2"
		case "tcp":
			procModelData[common.BKProtocol] = "1"
		default:
			procModelData[common.BKProtocol] = "1"
		}

		byteParams, _ := json.Marshal(procModelData)
		url := procAPI + "/process/v1/" + ownerID + "/" + appIDStr
		blog.Info("migrate add process url :%s", url)
		blog.Info("migrate add process content :%s", string(byteParams))
		reply, err := httpcli.ReqHttp(req, url, common.HTTPCreate, byteParams)
		blog.Info("migrate add process return :%s", string(reply))
		if err != nil {
			blog.Error("add process err :%v ", err)
			procName2ID[procName] = 0
			continue
		}
		js, err := simplejson.NewJson([]byte(reply))
		if nil != err {
			blog.Error("add bk process data return not json err :%v ", err)
			return err
		}
		output, err := js.Map()
		if nil != err {
			blog.Error("add bk process data return not json err :%v ", err)
			return err
		}
		code, err := util.GetIntByInterface(output["bk_error_code"])
		if err != nil || 0 != code {
			blog.Error("add process code err :%v ", err)
			continue
		}
		data, ok := output["data"].(map[string]interface{})
		if false == ok {
			blog.Error("add process data err :%v ", err)
			continue
		}
		procIDi, ok := data[common.BKProcIDField]
		if false == ok {
			blog.Error("add process data process ID err :%v ", err)
			continue
		}
		procID, err := util.GetIntByInterface(procIDi)
		if nil != err {
			continue
		}
		procName2ID[procName] = procID
	}
	addSetInBKApp(req)
	return nil
}

//addSetInBKApp add set in bk app
func addSetInBKApp(req *restful.Request) {
	appIDStr := strconv.Itoa(appID)
	for setName, moduleArr := range setModuleKv {
		setModelData[common.BKSetNameField] = setName
		setModelData[common.BKAppIDField] = appID
		setModelData[common.BKOwnerIDField] = common.BKDefaultOwnerID
		setModelData[common.BKInstParentStr] = appID
		byteParams, _ := json.Marshal(setModelData)
		url := topoAPI + "/topo/v1/set" + "/" + appIDStr
		blog.Info("migrate add set url :%s", url)
		blog.Info("migrate add set content :%s", string(byteParams))
		reply, err := httpcli.ReqHttp(req, url, common.HTTPCreate, byteParams)
		blog.Info("migrate add set return :%s", string(reply))
		if err != nil {
			blog.Error("add set data err :%v ", err)
			continue
		}
		js, _ := simplejson.NewJson([]byte(reply))
		output, _ := js.Map()

		code, err := util.GetIntByInterface(output["bk_error_code"])
		if err != nil || 0 != code {
			blog.Error("add set data code err :%v ", err)
			continue
		}
		data, ok := output["data"].(map[string]interface{})
		if false == ok {
			blog.Error("add set data result err :%v ", err)
			continue
		}
		setIDi, ok := data[common.BKSetIDField]
		if false == ok {
			continue
		}
		setID, err := util.GetIntByInterface(setIDi)
		if nil != err {
			continue
		}
		// add module in set
		addModuleInSet(req, moduleArr, setID)
	}
}

//addModuleInSet add module in set
func addModuleInSet(req *restful.Request, moduleArr map[string]string, setID int) {
	appIDStr := strconv.Itoa(appID)
	for moduleName, processNameStr := range moduleArr {
		moduleModelData[common.BKModuleNameField] = moduleName
		moduleModelData[common.BKAppIDField] = appID
		moduleModelData[common.BKSetIDField] = setID
		moduleModelData[common.BKOwnerIDField] = common.BKDefaultOwnerID
		moduleModelData[common.BKInstParentStr] = setID
		setIDStr := strconv.Itoa(setID)
		byteParams, _ := json.Marshal(moduleModelData)
		url := topoAPI + "/topo/v1/module" + "/" + appIDStr + "/" + setIDStr
		blog.Info("migrate add module url :%s", url)
		blog.Info("migrate add module content :%s", string(byteParams))
		reply, err := httpcli.ReqHttp(req, url, common.HTTPCreate, byteParams)
		blog.Info("migrate add module return :%s", string(reply))
		if err != nil {
			continue
		}
		js, _ := simplejson.NewJson([]byte(reply))
		output, _ := js.Map()

		code, err := util.GetIntByInterface(output["bk_error_code"])
		if err != nil || 0 != code {
			continue
		}
		//add module process config
		addModule2Process(req, processNameStr, moduleName)
	}
}

//addModule2Process add process 2 module
func addModule2Process(req *restful.Request, processNameStr string, moduleName string) {
	appIDStr := strconv.Itoa(appID)
	processNameArr := strings.Split(processNameStr, ",")
	for _, processName := range processNameArr {
		processID, ok := procName2ID[processName]
		if false == ok {
			continue
		}
		processIDStr := strconv.Itoa(processID)
		url := procAPI + "/process/v1/module" + "/" + common.BKDefaultOwnerID + "/" + appIDStr + "/" + processIDStr + "/" + moduleName
		blog.Info("migrate add module process config url :%s", url)
		reply, err := httpcli.ReqHttp(req, url, common.HTTPUpdate, nil)
		if err != nil {
			blog.Error("migrate add module process config %v", err)
			continue
		}
		blog.Info("migrate add module process return :%s", string(reply))
		js, err := simplejson.NewJson([]byte(reply))
		if err != nil {
			blog.Error("migrate add module process config json err %v", err)
			continue
		}
		output, err := js.Map()
		if err != nil {
			blog.Error("migrate add module process config data not map err %v", err)
			continue
		}

		code, err := util.GetIntByInterface(output["bk_error_code"])
		if err != nil || 0 != code {
			blog.Error("migrate add module process config code err %v", err)
			continue
		}
	}
}

//BKAppIsExist is bk app exist
func BKAppIsExist(req *restful.Request) (bool, error) {

	params := make(map[string]interface{})
	conditon := make(map[string]interface{})
	conditon[common.BKAppNameField] = common.BKAppName
	conditon[common.BKOwnerIDField] = ownerID
	params["condition"] = conditon
	params["fields"] = []string{common.BKAppIDField}
	params["start"] = 0
	params["limit"] = 20

	byteParams, _ := json.Marshal(params)
	url := topoAPI + "/topo/v1/app/search/" + ownerID
	blog.Info("Get bk app url :%s", url)
	blog.Info("Get bk app content :%s", string(byteParams))
	reply, err := httpcli.ReqHttp(req, url, common.HTTPSelectPost, byteParams)
	if err != nil {
		blog.Error("Get bk app error :%v", err)
		return false, err
	}
	blog.Info("Get bk app return :%s", string(reply))
	js, err := simplejson.NewJson([]byte(reply))
	if nil != err {
		blog.Error("Get bk app data not json error :%v", err)
	}
	output, err := js.Map()
	if nil != err {
		blog.Error("Get bk app data not map error :%v", err)
	}
	code, err := util.GetIntByInterface(output["bk_error_code"])
	if err != nil {
		blog.Error("Get bk app data not map error :%v", err)
		return false, errors.New(reply)
	}
	if 0 != code {
		blog.Error("Get bk app data not map error :%v", err)
		return false, errors.New(output["message"].(string))
	}
	cnt, err := js.Get("data").Get("count").Int()
	if err != nil {
		blog.Error("Get bk app data not count error :%v", err)
		return false, errors.New(reply)
	}
	if 0 == cnt {
		return false, nil
	}
	return true, nil
}
