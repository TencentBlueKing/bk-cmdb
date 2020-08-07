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
	"context"
	"strings"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
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

var procName2ID map[string]uint64

//addBKApp add bk app
func addBKApp(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	if count, err := db.Table(common.BKTableNameBaseApp).Find(mapstr.MapStr{common.BKAppNameField: common.BKAppName}).Count(ctx); err != nil {
		return err
	} else if count >= 1 {
		return nil
	}

	// add bk app
	appModelData := map[string]interface{}{}
	appModelData[common.BKAppNameField] = common.BKAppName
	appModelData[common.BKMaintainersField] = admin
	appModelData[common.BKTimeZoneField] = "Asia/Shanghai"
	appModelData[common.BKLanguageField] = "1" // "中文"
	appModelData[common.BKLifeCycleField] = common.DefaultAppLifeCycleNormal
	appModelData[common.BKOwnerIDField] = conf.OwnerID
	appModelData[common.BKDefaultField] = common.DefaultFlagDefaultValue
	appModelData[common.BKSupplierIDField] = conf.SupplierID
	filled := fillEmptyFields(appModelData, AppRow())
	var preData map[string]interface{}
	bizID, preData, err := upgrader.Upsert(ctx, db, common.BKTableNameBaseApp, appModelData, common.BKAppIDField, []string{common.BKAppNameField, common.BKOwnerIDField}, append(filled, common.BKAppIDField))
	if err != nil {
		blog.Error("add addBKApp error ", err.Error())
		return err
	}

	// add audit log
	properties := make([]metadata.Property, 0)
	for _, item := range AppRow() {
		properties = append(properties, metadata.Property{
			PropertyID:   item.PropertyID,
			PropertyName: item.PropertyName,
		})
	}
	log := metadata.AuditLog{
		AuditType:       metadata.BusinessType,
		SupplierAccount: conf.OwnerID,
		User:            conf.User,
		ResourceType:    metadata.BusinessRes,
		Action:          metadata.AuditCreate,
		OperateFrom:     metadata.FromCCSystem,
		OperationDetail: &metadata.InstanceOpDetail{
			BasicOpDetail: metadata.BasicOpDetail{
				BusinessID:   int64(bizID),
				BusinessName: common.BKAppName,
				ResourceID:   int64(bizID),
				ResourceName: common.BKAppName,
				Details: &metadata.BasicContent{
					PreData:    preData,
					CurData:    appModelData,
					Properties: properties,
				},
			},
			ModelID: common.BKInnerObjIDApp,
		},
		OperationTime: metadata.Now(),
		Label:         nil,
	}
	if preData != nil {
		log.Action = metadata.AuditUpdate
	}
	if err = db.Table(common.BKTableNameAuditLog).Insert(ctx, log); err != nil {
		blog.ErrorJSON("add audit log %s error %s", log, err.Error())
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
	setID, _, err := upgrader.Upsert(ctx, db, common.BKTableNameBaseSet, inputSetInfo, common.BKSetIDField, []string{common.BKOwnerIDField, common.BKAppIDField, common.BKSetNameField}, append(filled, common.BKSetIDField))
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
	_, _, err = upgrader.Upsert(ctx, db, common.BKTableNameBaseModule, inputResModuleInfo, common.BKModuleIDField, []string{common.BKOwnerIDField, common.BKModuleNameField, common.BKAppIDField, common.BKSetIDField}, append(filled, common.BKModuleIDField))
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
	_, _, err = upgrader.Upsert(ctx, db, common.BKTableNameBaseModule, inputFaultModuleInfo, common.BKModuleIDField, []string{common.BKOwnerIDField, common.BKModuleNameField, common.BKAppIDField, common.BKSetIDField}, append(filled, common.BKModuleIDField))
	if err != nil {
		blog.Error("add defaultFaultModule error ", err.Error())
		return err
	}

	if err := addBKProcess(ctx, db, conf, bizID); err != nil {
		blog.Error("add addBKProcess error ", err.Error())
	}
	if err := addSetInBKApp(ctx, db, conf, bizID); err != nil {
		blog.Error("add addSetInBKApp error ", err.Error())
	}

	return nil
}

//addBKProcess add bk process
func addBKProcess(ctx context.Context, db dal.RDB, conf *upgrader.Config, bizID uint64) error {
	procName2ID = make(map[string]uint64)

	for _, procStr := range prc2port {
		procArr := strings.Split(procStr, ":")
		procName := procArr[0]
		funcName := procArr[1]
		portStr := procArr[2]
		var protocol string
		if len(procArr) > 3 {
			protocol = procArr[3]
		}
		procModelData := map[string]interface{}{}
		procModelData[common.BKProcessNameField] = procName
		procModelData[common.BKFuncName] = funcName
		procModelData[common.BKPort] = portStr
		procModelData[common.BKWorkPath] = "/data/bkee"
		procModelData[common.BKOwnerIDField] = conf.OwnerID
		procModelData[common.BKAppIDField] = bizID

		protocol = strings.ToLower(protocol)
		switch protocol {
		case "udp":
			procModelData[common.BKProtocol] = "2"
		case "tcp":
			procModelData[common.BKProtocol] = "1"
		default:
			procModelData[common.BKProtocol] = "1"
		}

		filled := fillEmptyFields(procModelData, ProcRow())
		processID, _, err := upgrader.Upsert(ctx, db, common.BKTableNameBaseProcess, procModelData, common.BKProcessIDField, []string{common.BKProcessNameField, common.BKAppIDField, common.BKOwnerIDField}, append(filled, common.BKProcessIDField))
		if err != nil {
			blog.Error("add addBKProcess error ", err.Error())
			return err
		}
		procName2ID[procName] = processID
	}

	return nil
}

//addSetInBKApp add set in bk app
func addSetInBKApp(ctx context.Context, db dal.RDB, conf *upgrader.Config, bizID uint64) error {
	for setName, moduleArr := range setModuleKv {
		setModelData := map[string]interface{}{}
		setModelData[common.BKSetNameField] = setName
		setModelData[common.BKAppIDField] = bizID
		setModelData[common.BKOwnerIDField] = conf.OwnerID
		setModelData[common.BKInstParentStr] = bizID
		setModelData[common.BKDefaultField] = common.DefaultFlagDefaultValue
		setModelData[common.CreateTimeField] = time.Now()
		setModelData[common.LastTimeField] = time.Now()
		filled := fillEmptyFields(setModelData, SetRow())
		var preData map[string]interface{}
		setID, preData, err := upgrader.Upsert(ctx, db, common.BKTableNameBaseSet, setModelData, common.BKSetIDField, []string{common.BKSetNameField, common.BKOwnerIDField, common.BKAppIDField}, append(filled, common.BKSetIDField))
		if err != nil {
			blog.Error("add addSetInBKApp error ", err.Error())
			return err
		}

		// add audit log
		properties := make([]metadata.Property, 0)
		for _, item := range SetRow() {
			properties = append(properties, metadata.Property{
				PropertyID:   item.PropertyID,
				PropertyName: item.PropertyName,
			})
		}
		log := metadata.AuditLog{
			AuditType:       metadata.BusinessResourceType,
			SupplierAccount: conf.OwnerID,
			User:            conf.User,
			ResourceType:    metadata.SetRes,
			Action:          metadata.AuditCreate,
			OperateFrom:     metadata.FromCCSystem,
			OperationDetail: &metadata.InstanceOpDetail{
				BasicOpDetail: metadata.BasicOpDetail{
					BusinessID:   int64(bizID),
					BusinessName: common.BKAppName,
					ResourceID:   int64(setID),
					ResourceName: setName,
					Details: &metadata.BasicContent{
						PreData:    preData,
						CurData:    setModelData,
						Properties: properties,
					},
				},
				ModelID: common.BKInnerObjIDSet,
			},
			OperationTime: metadata.Now(),
			Label:         map[string]string{metadata.LabelBizTopology: ""},
		}
		if preData != nil {
			log.Action = metadata.AuditUpdate
		}
		if err = db.Table(common.BKTableNameAuditLog).Insert(ctx, log); err != nil {
			blog.ErrorJSON("add audit log %s error %s", log, err.Error())
			return err
		}

		// add module in set
		if err := addModuleInSet(ctx, db, conf, moduleArr, setID, bizID); err != nil {
			return err
		}
	}
	return nil
}

// addModuleInSet add module in set
func addModuleInSet(ctx context.Context, db dal.RDB, conf *upgrader.Config, moduleArr map[string]string, setID, bizID uint64) error {
	for moduleName, processNameStr := range moduleArr {
		moduleModelData := map[string]interface{}{}
		moduleModelData[common.BKModuleNameField] = moduleName
		moduleModelData[common.BKAppIDField] = bizID
		moduleModelData[common.BKSetIDField] = setID
		moduleModelData[common.BKOwnerIDField] = conf.OwnerID
		moduleModelData[common.BKInstParentStr] = setID
		moduleModelData[common.BKDefaultField] = common.DefaultFlagDefaultValue
		var preData map[string]interface{}
		filled := fillEmptyFields(moduleModelData, ModuleRow())
		moduleID, preData, err := upgrader.Upsert(ctx, db, common.BKTableNameBaseModule, moduleModelData, common.BKModuleIDField, []string{common.BKModuleNameField, common.BKOwnerIDField, common.BKAppIDField, common.BKSetIDField},
			append(filled, common.BKModuleIDField))
		if err != nil {
			blog.Error("add addModuleInSet error ", err.Error())
			return err
		}

		// add audit log
		properties := make([]metadata.Property, 0)
		for _, item := range ModuleRow() {
			properties = append(properties, metadata.Property{
				PropertyID:   item.PropertyID,
				PropertyName: item.PropertyName,
			})
		}
		log := metadata.AuditLog{
			AuditType:       metadata.BusinessResourceType,
			SupplierAccount: conf.OwnerID,
			User:            conf.User,
			ResourceType:    metadata.ModuleRes,
			Action:          metadata.AuditCreate,
			OperateFrom:     metadata.FromCCSystem,
			OperationDetail: &metadata.InstanceOpDetail{
				BasicOpDetail: metadata.BasicOpDetail{
					BusinessID:   int64(bizID),
					BusinessName: common.BKAppName,
					ResourceID:   int64(moduleID),
					ResourceName: moduleName,
					Details: &metadata.BasicContent{
						PreData:    preData,
						CurData:    moduleModelData,
						Properties: properties,
					},
				},
				ModelID: common.BKInnerObjIDModule,
			},
			OperationTime: metadata.Now(),
			Label:         map[string]string{metadata.LabelBizTopology: ""},
		}
		if preData != nil {
			log.Action = metadata.AuditUpdate
		}
		if err = db.Table(common.BKTableNameAuditLog).Insert(ctx, log); err != nil {
			blog.ErrorJSON("add audit log %s error %s", log, err.Error())
			return err
		}

		//add module process config
		if err := addModule2Process(ctx, db, conf, processNameStr, moduleName, bizID); err != nil {
			return err
		}

	}
	return nil
}

//addModule2Process add process 2 module
func addModule2Process(ctx context.Context, db dal.RDB, conf *upgrader.Config, processNameStr string, moduleName string, bizID uint64) (err error) {
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

		if _, _, err = upgrader.Upsert(ctx, db, common.BKTableNameProcModule, module2Process, "", []string{common.BKModuleNameField, common.BKAppIDField, common.BKProcessIDField}, nil); err != nil {
			blog.Error("add addModuleInSet error ", err.Error())
			return err
		}
	}
	return nil
}
