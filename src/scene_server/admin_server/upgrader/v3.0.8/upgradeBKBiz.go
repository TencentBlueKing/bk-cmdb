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

package v3v0v8

import (
	"fmt"

	"strings"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/auditoplog"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage"
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

var procName2ID map[string]int64

//addBKApp add bk app
func addBKApp(db storage.DI, conf *upgrader.Config) error {

	// add bk app
	appModelData := map[string]interface{}{}
	appModelData[common.BKAppNameField] = common.BKAppName
	appModelData[common.BKMaintainersField] = "admin"
	appModelData[common.BKTimeZoneField] = "Asia/Shanghai"
	appModelData[common.BKLanguageField] = "1" //"中文"
	appModelData[common.BKLifeCycleField] = common.DefaultAppLifeCycleNormal
	appModelData[common.BKOwnerIDField] = conf.OwnerID
	appModelData[common.BKDefaultField] = 0
	appModelData[common.BKSupplierIDField] = conf.SupplierID
	filled := fillEmptyFields(appModelData, AppRow())
	var preData map[string]interface{}
	bizID, preData, err := upgrader.Upsert(db, "cc_ApplicationBase", appModelData, common.BKAppIDField, []string{common.BKAppNameField, common.BKOwnerIDField}, append(filled, common.BKAppIDField))
	if err != nil {
		blog.Error("add addBKApp error ", err.Error())
		return err
	}

	// add audit log
	headers := []metadata.Header{}
	for _, item := range AppRow() {
		headers = append(headers, metadata.Header{
			PropertyID:   item.PropertyID,
			PropertyName: item.PropertyName,
		})
	}
	auditContent := metadata.Content{
		CurData: appModelData,
		Headers: headers,
	}
	logRow := &metadata.OperationLog{
		OwnerID:       conf.OwnerID,
		ApplicationID: bizID,
		OpType:        int(auditoplog.AuditOpTypeAdd),
		OpTarget:      "biz",
		User:          conf.User,
		ExtKey:        "",
		OpDesc:        "create app",
		Content:       auditContent,
		CreateTime:    time.Now(),
		InstID:        bizID,
	}
	if preData != nil {
		logRow.OpDesc = "update process"
		logRow.OpType = int(auditoplog.AuditOpTypeModify)
	}
	if _, err = db.Insert(logRow.TableName(), logRow); err != nil {
		blog.Error("add audit log error ", err.Error())
		return err
	}

	// add bk app default set
	inputSetInfo := make(map[string]interface{})
	inputSetInfo[common.BKAppIDField] = bizID
	inputSetInfo[common.BKInstParentStr] = bizID
	inputSetInfo[common.BKSetNameField] = common.DefaultResSetName
	inputSetInfo[common.BKDefaultField] = common.DefaultResSetFlag
	inputSetInfo[common.BKOwnerIDField] = conf.OwnerID
	filled = fillEmptyFields(inputSetInfo, SetRow())
	setID, _, err := upgrader.Upsert(db, "cc_SetBase", inputSetInfo, common.BKSetIDField, []string{common.BKOwnerIDField, common.BKAppIDField, common.BKSetNameField}, append(filled, common.BKSetIDField))
	if err != nil {
		blog.Error("add defaultSet error ", err.Error())
		return err
	}

	// add bk app default module
	inputResModuleInfo := make(map[string]interface{})
	inputResModuleInfo[common.BKSetIDField] = setID
	inputResModuleInfo[common.BKInstParentStr] = setID
	inputResModuleInfo[common.BKAppIDField] = bizID
	inputResModuleInfo[common.BKModuleNameField] = common.DefaultResModuleName
	inputResModuleInfo[common.BKDefaultField] = common.DefaultResModuleFlag
	inputResModuleInfo[common.BKOwnerIDField] = conf.OwnerID
	filled = fillEmptyFields(inputResModuleInfo, ModuleRow())
	_, _, err = upgrader.Upsert(db, "cc_ModuleBase", inputResModuleInfo, common.BKModuleIDField, []string{common.BKOwnerIDField, common.BKModuleNameField, common.BKAppIDField, common.BKSetIDField}, append(filled, common.BKModuleIDField))
	if err != nil {
		blog.Error("add defaultResModule error ", err.Error())
		return err
	}

	inputFaultModuleInfo := make(map[string]interface{})
	inputFaultModuleInfo[common.BKSetIDField] = setID
	inputFaultModuleInfo[common.BKInstParentStr] = setID
	inputFaultModuleInfo[common.BKAppIDField] = bizID
	inputFaultModuleInfo[common.BKModuleNameField] = common.DefaultFaultModuleName
	inputFaultModuleInfo[common.BKDefaultField] = common.DefaultFaultModuleFlag
	inputFaultModuleInfo[common.BKOwnerIDField] = conf.OwnerID
	filled = fillEmptyFields(inputFaultModuleInfo, ModuleRow())
	_, _, err = upgrader.Upsert(db, "cc_ModuleBase", inputFaultModuleInfo, common.BKModuleIDField, []string{common.BKOwnerIDField, common.BKModuleNameField, common.BKAppIDField, common.BKSetIDField}, append(filled, common.BKModuleIDField))
	if err != nil {
		blog.Error("add defaultFaultModule error ", err.Error())
		return err
	}

	if err := addBKProcess(db, conf, bizID); err != nil {
		blog.Error("add addBKProcess error ", err.Error())
	}
	if err := addSetInBKApp(db, conf, bizID); err != nil {
		blog.Error("add addSetInBKApp error ", err.Error())
	}

	return nil
}

//addBKProcess add bk process
func addBKProcess(db storage.DI, conf *upgrader.Config, bizID int64) error {
	procName2ID = make(map[string]int64)

	for _, procStr := range prc2port {
		procArr := strings.Split(procStr, ":")
		procName := procArr[0]
		funcName := procArr[1]
		portStr := procArr[2]
		var protocal string
		if len(procArr) > 3 {
			protocal = procArr[3]
		}
		procModelData := map[string]interface{}{}
		procModelData[common.BKProcessNameField] = procName
		procModelData[common.BKFuncName] = funcName
		procModelData[common.BKPort] = portStr
		procModelData[common.BKWorkPath] = "/data/bkee"
		procModelData[common.BKOwnerIDField] = conf.OwnerID
		procModelData[common.BKAppIDField] = bizID

		protocal = strings.ToLower(protocal)
		switch protocal {
		case "udp":
			procModelData[common.BKProtocol] = "2"
		case "tcp":
			procModelData[common.BKProtocol] = "1"
		default:
			procModelData[common.BKProtocol] = "1"
		}

		filled := fillEmptyFields(procModelData, ProcRow())
		var preData map[string]interface{}
		processID, preData, err := upgrader.Upsert(db, "cc_Process", procModelData, common.BKProcessIDField, []string{common.BKProcessNameField, common.BKAppIDField, common.BKOwnerIDField}, append(filled, common.BKProcessIDField))
		if err != nil {
			blog.Error("add addBKProcess error ", err.Error())
			return err
		}
		procName2ID[procName] = processID

		// add audit log
		headers := []metadata.Header{}
		for _, item := range ProcRow() {
			headers = append(headers, metadata.Header{
				PropertyID:   item.PropertyID,
				PropertyName: item.PropertyName,
			})
		}
		auditContent := metadata.Content{
			CurData: procModelData,
			Headers: headers,
		}
		logRow := &metadata.OperationLog{
			OwnerID:       conf.OwnerID,
			ApplicationID: bizID,
			OpType:        int(auditoplog.AuditOpTypeAdd),
			OpTarget:      "process",
			User:          conf.User,
			ExtKey:        "",
			OpDesc:        "create process",
			Content:       auditContent,
			CreateTime:    time.Now(),
			InstID:        processID,
		}
		if preData != nil {
			logRow.OpDesc = "update process"
			logRow.OpType = int(auditoplog.AuditOpTypeModify)
		}
		if _, err = db.Insert(logRow.TableName(), logRow); err != nil {
			blog.Error("add audit log error ", err.Error())
			return err
		}

	}

	return nil
}

//addSetInBKApp add set in bk app
func addSetInBKApp(db storage.DI, conf *upgrader.Config, bizID int64) error {
	for setName, moduleArr := range setModuleKv {
		setModelData := map[string]interface{}{}
		setModelData[common.BKSetNameField] = setName
		setModelData[common.BKAppIDField] = bizID
		setModelData[common.BKOwnerIDField] = conf.OwnerID
		setModelData[common.BKInstParentStr] = bizID
		setModelData[common.BKDefaultField] = 0
		setModelData[common.CreateTimeField] = time.Now()
		setModelData[common.LastTimeField] = time.Now()
		filled := fillEmptyFields(setModelData, SetRow())
		var preData map[string]interface{}
		setID, preData, err := upgrader.Upsert(db, "cc_SetBase", setModelData, common.BKSetIDField, []string{common.BKSetNameField, common.BKOwnerIDField, common.BKAppIDField}, append(filled, common.BKSetIDField))
		if err != nil {
			blog.Error("add addSetInBKApp error ", err.Error())
			return err
		}

		// add audit log
		headers := []metadata.Header{}
		for _, item := range SetRow() {
			headers = append(headers, metadata.Header{
				PropertyID:   item.PropertyID,
				PropertyName: item.PropertyName,
			})
		}
		auditContent := metadata.Content{
			CurData: setModelData,
			Headers: headers,
		}
		logRow := &metadata.OperationLog{
			OwnerID:       conf.OwnerID,
			ApplicationID: bizID,
			OpType:        int(auditoplog.AuditOpTypeAdd),
			OpTarget:      "set",
			User:          conf.User,
			ExtKey:        "",
			OpDesc:        "create set",
			Content:       auditContent,
			CreateTime:    time.Now(),
			InstID:        setID,
		}
		if preData != nil {
			logRow.OpDesc = "update set"
			logRow.OpType = int(auditoplog.AuditOpTypeModify)
		}
		if _, err = db.Insert(logRow.TableName(), logRow); err != nil {
			blog.Error("add audit log error ", err.Error())
			return err
		}

		// add module in set
		if err := addModuleInSet(db, conf, moduleArr, setID, bizID); err != nil {
			return err
		}
	}
	return nil
}

//addModuleInSet add module in set
func addModuleInSet(db storage.DI, conf *upgrader.Config, moduleArr map[string]string, setID, bizID int64) error {
	for moduleName, processNameStr := range moduleArr {
		moduleModelData := map[string]interface{}{}
		moduleModelData[common.BKModuleNameField] = moduleName
		moduleModelData[common.BKAppIDField] = bizID
		moduleModelData[common.BKSetIDField] = setID
		moduleModelData[common.BKOwnerIDField] = conf.OwnerID
		moduleModelData[common.BKInstParentStr] = setID
		moduleModelData[common.BKDefaultField] = 0
		var preData map[string]interface{}
		filled := fillEmptyFields(moduleModelData, ModuleRow())
		moduleID, preData, err := upgrader.Upsert(db, "cc_ModuleBase", moduleModelData, common.BKModuleIDField, []string{common.BKModuleNameField, common.BKOwnerIDField, common.BKAppIDField, common.BKSetIDField},
			append(filled, common.BKModuleIDField))
		if err != nil {
			blog.Error("add addModuleInSet error ", err.Error())
			return err
		}

		// add audit log
		headers := []metadata.Header{}
		for _, item := range ModuleRow() {
			headers = append(headers, metadata.Header{
				PropertyID:   item.PropertyID,
				PropertyName: item.PropertyName,
			})
		}
		auditContent := metadata.Content{
			CurData: moduleModelData,
			PreData: preData,
			Headers: headers,
		}
		logRow := &metadata.OperationLog{
			OwnerID:       conf.OwnerID,
			ApplicationID: bizID,
			OpType:        int(auditoplog.AuditOpTypeAdd),
			OpTarget:      "module",
			User:          conf.User,
			ExtKey:        "",
			OpDesc:        "create module",
			Content:       auditContent,
			CreateTime:    time.Now(),
			InstID:        moduleID,
		}
		if preData != nil {
			logRow.OpDesc = "update module"
			logRow.OpType = int(auditoplog.AuditOpTypeModify)
		}
		if _, err = db.Insert(logRow.TableName(), logRow); err != nil {
			blog.Error("add audit log error ", err.Error())
			return err
		}

		//add module process config
		if err := addModule2Process(db, conf, processNameStr, moduleName, bizID); err != nil {
			return err
		}

	}
	return nil
}

//addModule2Process add process 2 module
func addModule2Process(db storage.DI, conf *upgrader.Config, processNameStr string, moduleName string, bizID int64) (err error) {
	processNameArr := strings.Split(processNameStr, ",")
	for _, processName := range processNameArr {
		processID, ok := procName2ID[processName]
		if false == ok {
			continue
		}
		module2Process := map[string]interface{}{}
		module2Process[common.BKAppIDField] = bizID
		module2Process[common.BKModuleNameField] = moduleName
		module2Process[common.BKProcessIDField] = processID

		if _, _, err = upgrader.Upsert(db, "cc_Proc2Module", module2Process, "", []string{common.BKModuleNameField, common.BKAppIDField, common.BKProcessIDField}, nil); err != nil {
			blog.Error("add addModuleInSet error ", err.Error())
			return err
		}

		// add audit log
		headers := []metadata.Header{}
		for _, item := range ModuleRow() {
			headers = append(headers, metadata.Header{
				PropertyID:   item.PropertyID,
				PropertyName: item.PropertyName,
			})
		}
		logRow := &metadata.OperationLog{
			OwnerID:       conf.OwnerID,
			ApplicationID: bizID,
			OpType:        int(auditoplog.AuditOpTypeModify),
			OpTarget:      "module",
			User:          conf.User,
			ExtKey:        "",
			OpDesc:        fmt.Sprintf("bind module [%s]", moduleName),
			Content:       "",
			CreateTime:    time.Now(),
			InstID:        bizID,
		}
		if _, err = db.Insert(logRow.TableName(), logRow); err != nil {
			blog.Error("add audit log error ", err.Error())
			return err
		}
	}
	return nil
}
