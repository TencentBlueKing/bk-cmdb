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

package command

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/condition"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/storage/dal"
)

var (
	defaultinitUserName = "bk_init_user"
)

func backup(ctx context.Context, db dal.RDB, opt *option) error {
	dir := filepath.Dir(opt.position)
	now := time.Now().Format("2006_01_02_15_04_05")
	file := filepath.Join(dir, "backup_bk_biz_"+now+".json")
	exportOpt := *opt
	exportOpt.position = file
	exportOpt.mini = false
	exportOpt.scope = "all"
	err := export(ctx, db, &exportOpt)
	if nil != err {
		return err
	}
	fmt.Println("%s business has been backup to \033[35m"+file+"\033[0m", opt.bizName)
	return nil
}

func importBKBiz(ctx context.Context, db dal.RDB, opt *option) error {
	file, err := os.OpenFile(opt.position, os.O_RDONLY, os.ModePerm)
	if nil != err {
		return err
	}
	defer file.Close()

	importJSON := BKTopo{}
	err = json.NewDecoder(file).Decode(&importJSON)
	if nil != err {
		return err
	}

	bizID, err := getBKBizID(ctx, db)
	if err != nil {
		return err
	}

	setParentID, err := getSetParentID(ctx, bizID, db)
	if err != nil {
		return err
	}

	if err := allowInit(ctx, bizID, db, opt); err != nil {
		return err
	}
	importer := NewImporterBizTopo(importJSON, db, opt)
	if err := importer.FilterBKTopo(ctx, bizID, setParentID); err != nil {
		return err
	}

	if err := importer.ClearBKTopo(ctx, bizID); err != nil {
		return err
	}

	if err := importer.InitBKTopo(ctx, bizID, setParentID); err != nil {
		return err
	}
	if err := recordInitLog(ctx, db); err != nil {
		return err
	}

	return nil
}

// aloowInit 是否允许初始化
func allowInit(ctx context.Context, bizID int64, db dal.DB, opt *option) error {

	// 是否已经初始化过
	initFlagCond := map[string]interface{}{"cc_init_bk_biz_init": mapstr.MapStr{"$exists": true}}
	cnt, err := db.Table(common.BKTableNameSystem).Find(initFlagCond).Count(ctx)
	if err != nil {
		return fmt.Errorf("find cc_init_bk_biz_init flag from db error. err:%s ", err.Error())
	}
	if cnt != 0 {
		return fmt.Errorf("no duplicat import allowed")
	}

	// 是否已经有主机
	hostInfo := map[string]interface{}{common.BKAppIDField: bizID}
	cnt, err = db.Table(common.BKTableNameModuleHostConfig).Find(hostInfo).Count(ctx)
	if err != nil {
		return fmt.Errorf("find blueking business host error. err:%s", err.Error())
	}
	if cnt != 0 {
		return fmt.Errorf("host already exists")
	}
	return nil
}

func recordInitLog(ctx context.Context, db dal.DB) error {
	doc := map[string]interface{}{
		"cc_init_bk_biz_init": time.Now(),
	}
	err := db.Table(common.BKTableNameSystem).Insert(ctx, doc)
	if err != nil {
		return fmt.Errorf("record business blueking topology operation log error. err:%s", err.Error())
	}
	return nil
}

type importerBizTopo struct {

	// file json content
	importJSON BKTopo
	db         dal.DB
	opt        *option

	// first handle data. 检查json合法后的数据
	procFuncNameInfoMap  map[string]metadata.ProcessTemplate
	serviceTemplateMap   map[string]BKServiceTemplate
	setNameInfoMap       map[string]map[string]interface{}
	moduleSetNameInfoMap map[string]map[string]BKBizModule

	// create topo info
	newServiceTemplateMap map[string]int64
	newSetTemplate        map[string]int64

	// cache
	// 缓存业务下服务分类 map[level1 category name]catorgy id
	serviceCategoryL1CacheInfo map[string]int64
	// 缓存业务下服务分类 map[level1 category id ][level2]catorgy id
	serviceCategoryL2CacheInfo map[int64]map[string]int64
}

func NewImporterBizTopo(importJSON BKTopo, db dal.DB, opt *option) *importerBizTopo {
	return &importerBizTopo{
		importJSON:                 importJSON,
		db:                         db,
		opt:                        opt,
		procFuncNameInfoMap:        make(map[string]metadata.ProcessTemplate, 0),
		serviceTemplateMap:         make(map[string]BKServiceTemplate, 0),
		setNameInfoMap:             make(map[string]map[string]interface{}, 0),
		moduleSetNameInfoMap:       make(map[string]map[string]BKBizModule, 0),
		newServiceTemplateMap:      make(map[string]int64),
		newSetTemplate:             make(map[string]int64, 0),
		serviceCategoryL1CacheInfo: make(map[string]int64, 0),
		serviceCategoryL2CacheInfo: make(map[int64]map[string]int64, 0),
	}
}

func (ibt *importerBizTopo) FilterBKTopo(ctx context.Context, bizID, setParentID int64) error {

	// 检查用户配置服务分类是否存在
	if err := ibt.cacheServiceCategory(ctx, bizID); err != nil {
		return err
	}

	if err := ibt.filterBKTopoProc(ctx, bizID); err != nil {
		return err
	}

	if err := ibt.filterBKTopoServiceTemplate(ctx); err != nil {
		return err
	}

	if err := ibt.filterBKTopoSet(ctx, bizID, setParentID); err != nil {
		return err
	}

	if err := ibt.filterBKTopoModule(ctx); err != nil {
		return err
	}

	return nil
}

func (ibt *importerBizTopo) filterBKTopoProc(ctx context.Context, bizID int64) error {
	for idx, proc := range ibt.importJSON.Proc {
		funcName, ok := proc[common.BKFuncName].(string)
		if !ok {
			funcName, ok = proc[common.BKProcNameField].(string)
			if !ok {
				return fmt.Errorf("process info index %d, field %s value not string", idx, common.BKFuncName)
			}
			proc[common.BKFuncName] = funcName
		} else {
			proc[common.BKProcNameField] = funcName
		}
		if _, ok := ibt.procFuncNameInfoMap[funcName]; ok {
			return fmt.Errorf("process info index %d,  %s  duplicate", idx, common.BKFuncName)
		}

		procTemp := metadata.ProcessTemplate{
			// 	set value befor insert data to db
			ID:          0,
			ProcessName: funcName,
			// set value  after create template,
			ServiceTemplateID: 0,
			BizID:             bizID,
			Creator:           defaultinitUserName,
			Modifier:          defaultinitUserName,
			CreateTime:        time.Now().UTC(),
			LastTime:          time.Now().UTC(),
			SupplierAccount:   ibt.opt.OwnerID,
			Property:          nil,
		}
		var err error
		procTemp.Property, err = convProcTemplateProperty(ctx, proc)
		if err != nil {
			return fmt.Errorf("process index %d, name:%s, %s", idx, funcName, err.Error())
		}
		ibt.procFuncNameInfoMap[funcName] = procTemp
	}
	return nil
}

func (ibt *importerBizTopo) filterBKTopoServiceTemplate(ctx context.Context) error {

	for idx, srvTemp := range ibt.importJSON.ServiceTemplateArr {
		for _, procName := range srvTemp.BindProcess {
			if _, ok := ibt.procFuncNameInfoMap[procName]; !ok {
				return fmt.Errorf("service template  index %d, name:%s, bind process name[%s] not found", idx, srvTemp.Name, procName)
			}
		}
		if len(srvTemp.ServiceCategoryName) == 0 {
			srvTemp.ServiceCategoryName = []string{common.DefaultServiceCategoryName, common.DefaultServiceCategoryName}
		}
		if len(srvTemp.ServiceCategoryName) != 2 {
			return fmt.Errorf("ervice template  index %d, name:%s, service category must be tow level. not %d", idx, srvTemp.Name, len(srvTemp.ServiceCategoryName))
		}
		srvTempL1ID, ok := ibt.serviceCategoryL1CacheInfo[srvTemp.ServiceCategoryName[0]]
		if !ok {
			return fmt.Errorf("ervice template  index %d, name:%s, service category  level1 name[%s] not found", idx, srvTemp.Name, srvTemp.ServiceCategoryName[0])
		}
		srvTempL2ID, ok := ibt.serviceCategoryL2CacheInfo[srvTempL1ID][srvTemp.ServiceCategoryName[1]]
		if !ok {
			return fmt.Errorf("ervice template  index %d, name:%s, service category level2 name[%s] not found", idx, srvTemp.Name, srvTemp.ServiceCategoryName[1])
		}

		srvTemp.ServiceCategoryID = srvTempL2ID
		ibt.serviceTemplateMap[srvTemp.Name] = srvTemp
	}
	return nil
}

func (ibt *importerBizTopo) filterBKTopoSet(ctx context.Context, bizID, setParentID int64) error {
	for idx, setInfo := range ibt.importJSON.Topo.SetArr {
		setName, ok := setInfo[common.BKSetNameField].(string)
		if !ok {
			return fmt.Errorf("set info index %d, field %s value not string", idx, common.BKSetNameField)
		}
		ibt.setNameInfoMap[setName] = setInfo
	}
	return nil
}

func (ibt *importerBizTopo) filterBKTopoModule(ctx context.Context) error {
	for idx, module := range ibt.importJSON.Topo.ModuleArr {

		if _, ok := ibt.setNameInfoMap[module.SetName]; !ok {
			return fmt.Errorf("module info index:%d, set name[%s] not found", idx, module.SetName)
		}
		if _, ok := ibt.moduleSetNameInfoMap[module.SetName]; !ok {
			ibt.moduleSetNameInfoMap[module.SetName] = make(map[string]BKBizModule, 0)
		}
		if _, ok := ibt.serviceTemplateMap[module.ServiceTemplate]; !ok {
			return fmt.Errorf("module info index:%d, service template[%s] not found", idx, module.ServiceTemplate)
		}
		ibt.moduleSetNameInfoMap[module.SetName][module.ServiceTemplate] = module
	}
	return nil
}

func (ibt *importerBizTopo) ClearBKTopo(ctx context.Context, bizID int64) error {

	// clear process template
	deleteProcTempCond := condition.CreateCondition()
	deleteProcTempCond.Field(common.BKAppIDField).Eq(bizID)
	err := ibt.db.Table(common.BKTableNameProcessTemplate).Delete(ctx, deleteProcTempCond.ToMapStr())
	if err != nil {
		return fmt.Errorf("clear business topology error. delete service template error. err:%s", err.Error())
	}
	// clear set
	deleteSetTempCond := condition.CreateCondition()
	deleteSetTempCond.Field(common.BKAppIDField).Eq(bizID)
	deleteSetTempCond.Field(common.BKDefaultField).Eq(common.NormalSetDefaultFlag)
	err = ibt.db.Table(common.BKTableNameBaseSet).Delete(ctx, deleteSetTempCond.ToMapStr())
	if err != nil {
		return fmt.Errorf("clear business topology error. delete set template error. err:%s", err.Error())
	}

	// clear service template
	deleteSrvTempCond := condition.CreateCondition()
	deleteSrvTempCond.Field(common.BKAppIDField).Eq(bizID)
	err = ibt.db.Table(common.BKTableNameServiceTemplate).Delete(ctx, deleteSrvTempCond.ToMapStr())
	if err != nil {
		return fmt.Errorf("clear business topology error. delete service template error. err:%s", err.Error())
	}
	// clear module
	deleteModuleTempCond := condition.CreateCondition()
	deleteModuleTempCond.Field(common.BKAppIDField).Eq(bizID)
	deleteModuleTempCond.Field(common.BKDefaultField).Eq(common.NormalModuleFlag)
	err = ibt.db.Table(common.BKTableNameBaseModule).Delete(ctx, deleteModuleTempCond.ToMapStr())
	if err != nil {
		return fmt.Errorf("clear business topology error. delete module template error. err:%s", err.Error())
	}

	return nil
}

func (ibt *importerBizTopo) InitBKTopo(ctx context.Context, bizID, setParentID int64) error {
	if err := ibt.initBKServiceCategory(ctx, bizID); err != nil {
		return err
	}

	if err := ibt.initBKTopoSet(ctx, bizID, setParentID); err != nil {
		return err
	}

	if err := ibt.initBKTopoModule(ctx, bizID); err != nil {
		return err
	}
	return nil
}

func (ibt *importerBizTopo) initBKServiceCategory(ctx context.Context, bizID int64) error {
	if len(ibt.serviceTemplateMap) == 0 {
		return nil
	}
	var srvTempArr []interface{}
	var procTempArr []interface{}
	for _, srvTemp := range ibt.serviceTemplateMap {
		nextSrvTempID, err := ibt.db.NextSequence(ctx, common.BKTableNameServiceTemplate)
		if err != nil {
			return fmt.Errorf("init service template, get next id error. err:%s", err.Error())
		}

		ibt.newServiceTemplateMap[srvTemp.Name] = int64(nextSrvTempID)
		srvTempArr = append(srvTempArr, metadata.ServiceTemplate{
			BizID:             bizID,
			ID:                int64(nextSrvTempID),
			Name:              srvTemp.Name,
			ServiceCategoryID: srvTemp.ServiceCategoryID,
			Creator:           defaultinitUserName,
			Modifier:          defaultinitUserName,
			CreateTime:        time.Now().UTC(),
			LastTime:          time.Now().UTC(),
			SupplierAccount:   ibt.opt.OwnerID,
		})

		for _, procName := range srvTemp.BindProcess {
			nextProcTempID, err := ibt.db.NextSequence(ctx, common.BKTableNameProcessTemplate)
			if err != nil {
				return fmt.Errorf("init service template, get next id error. err:%s", err.Error())
			}
			procTemp := ibt.procFuncNameInfoMap[procName]
			procTemp.ServiceTemplateID = int64(nextSrvTempID)
			procTemp.ID = int64(nextProcTempID)
			procTempArr = append(procTempArr, procTemp)
		}
	}

	err := ibt.db.Table(common.BKTableNameServiceTemplate).Insert(ctx, srvTempArr)
	if err != nil {
		return fmt.Errorf("init service template error. err:%s", err.Error())
	}
	err = ibt.db.Table(common.BKTableNameProcessTemplate).Insert(ctx, procTempArr)
	if err != nil {
		return fmt.Errorf("init process template error. err:%s", err.Error())
	}
	return nil
}

func (ibt *importerBizTopo) initBKTopoSet(ctx context.Context, bizID, setParentID int64) error {
	var setArr []interface{}
	for name, setInfo := range ibt.setNameInfoMap {
		nextSetID, err := ibt.db.NextSequence(ctx, common.BKTableNameBaseSet)
		if err != nil {
			return fmt.Errorf("init service template, get next id error. err:%s", err.Error())
		}
		ibt.newSetTemplate[name] = int64(nextSetID)
		setInfo[common.BKSetNameField] = name
		setInfo[common.BKSetIDField] = int64(nextSetID)
		setInfo[common.BKParentIDField] = setParentID
		setInfo[common.BKAppIDField] = bizID
		setInfo[common.BKDefaultField] = common.NormalSetDefaultFlag
		setInfo[common.CreateTimeField] = time.Now().UTC()
		setInfo[common.LastTimeField] = time.Now().UTC()
		setInfo[common.BKOwnerIDField] = ibt.opt.OwnerID
		setArr = append(setArr, setInfo)
	}

	err := ibt.db.Table(common.BKTableNameBaseSet).Insert(ctx, setArr)
	if err != nil {
		return fmt.Errorf("init set  error. err:%s", err.Error())
	}
	return nil
}

func (ibt *importerBizTopo) initBKTopoModule(ctx context.Context, bizID int64) error {

	var moduleArr []interface{}
	for setName, moduleNameMap := range ibt.moduleSetNameInfoMap {

		setID := ibt.newSetTemplate[setName]
		for _, moduleTemp := range moduleNameMap {
			nextModuleID, err := ibt.db.NextSequence(ctx, common.BKTableNameBaseModule)
			if err != nil {
				return fmt.Errorf("init service template, get next id error. err:%s", err.Error())
			}
			srvTempID := ibt.newServiceTemplateMap[moduleTemp.ServiceTemplate]
			srvTempInfo := ibt.serviceTemplateMap[moduleTemp.ServiceTemplate]

			moduleInfo := make(map[string]interface{}, 0)
			for key, val := range moduleTemp.Info {
				moduleInfo[key] = val
			}
			moduleInfo[common.BKModuleNameField] = moduleTemp.ServiceTemplate
			moduleInfo[common.BKModuleIDField] = int64(nextModuleID)
			moduleInfo[common.BKSetIDField] = setID
			moduleInfo[common.BKParentIDField] = setID
			moduleInfo[common.BKAppIDField] = bizID
			moduleInfo[common.BKDefaultField] = common.NormalModuleFlag
			moduleInfo[common.CreateTimeField] = time.Now().UTC()
			moduleInfo[common.LastTimeField] = time.Now().UTC()
			moduleInfo[common.BKOwnerIDField] = ibt.opt.OwnerID
			moduleInfo[common.BKServiceCategoryIDField] = srvTempInfo.ServiceCategoryID
			moduleInfo[common.BKServiceTemplateIDField] = srvTempID
			if _, ok := moduleInfo[common.BKModuleTypeField]; !ok {
				moduleInfo[common.BKModuleTypeField] = "1"
			}
			moduleArr = append(moduleArr, moduleInfo)
		}
	}
	err := ibt.db.Table(common.BKTableNameBaseModule).Insert(ctx, moduleArr)
	if err != nil {
		return fmt.Errorf("init module error. err:%s", err)
	}

	return nil

}

func (ibt *importerBizTopo) cacheServiceCategory(ctx context.Context, bizID int64) error {

	// find build in service  category
	searchBuindInCond := condition.CreateCondition()
	searchBuindInCond.Field("is_built_in").Eq(true)
	serviceCategoryArr := make([]metadata.ServiceCategory, 0)
	err := ibt.db.Table(common.BKTableNameServiceCategory).Find(searchBuindInCond.ToMapStr()).All(ctx, &serviceCategoryArr)
	if err != nil {
		return fmt.Errorf("find build-in service category error. err:%s", err.Error())
	}
	for _, serviceCategory := range serviceCategoryArr {
		if serviceCategory.ParentID == 0 {
			ibt.serviceCategoryL1CacheInfo[serviceCategory.Name] = serviceCategory.ID
		} else {
			if _, ok := ibt.serviceCategoryL2CacheInfo[serviceCategory.ParentID]; !ok {
				ibt.serviceCategoryL2CacheInfo[serviceCategory.ParentID] = make(map[string]int64, 0)
			}
			ibt.serviceCategoryL2CacheInfo[serviceCategory.ParentID][serviceCategory.Name] = serviceCategory.ID
		}
	}

	// find build in service  category
	searchCond := condition.CreateCondition()
	searchCond.Field(common.BKAppIDField).Eq(bizID)
	serviceCategoryArr = make([]metadata.ServiceCategory, 0)
	err = ibt.db.Table(common.BKTableNameServiceCategory).Find(searchCond.ToMapStr()).All(ctx, &serviceCategoryArr)
	if err != nil {
		return fmt.Errorf("find business service category error. err:%s", err.Error())
	}
	for _, serviceCategory := range serviceCategoryArr {
		if serviceCategory.ParentID == 0 {
			ibt.serviceCategoryL1CacheInfo[serviceCategory.Name] = serviceCategory.ID
		} else {
			if _, ok := ibt.serviceCategoryL2CacheInfo[serviceCategory.ParentID]; !ok {
				ibt.serviceCategoryL2CacheInfo[serviceCategory.ParentID] = make(map[string]int64, 0)
			}
			ibt.serviceCategoryL2CacheInfo[serviceCategory.ParentID][serviceCategory.Name] = serviceCategory.ID
		}
	}

	return nil
}

func convProcTemplateProperty(ctx context.Context, proc map[string]interface{}) (*metadata.ProcessProperty, error) {
	processProperty := &metadata.ProcessProperty{}
	blTrue := true
	for key, val := range proc {
		switch key {
		case "proc_num":
			procNum, err := util.GetInt64ByInterface(val)
			if err != nil {
				return nil, fmt.Errorf("%s not integer. val:%s", key, val)
			}
			processProperty.ProcNum.Value = &procNum
			processProperty.ProcNum.AsDefaultValue = &blTrue
			if err := processProperty.ProcNum.Validate(); err != nil {
				return nil, fmt.Errorf("%s illegal. val:%s. err:%s", key, val, err.Error())
			}
		case "stop_cmd":
			stopCmd, ok := val.(string)
			if !ok {
				return nil, fmt.Errorf("%s not string. val:%s", key, val)
			}
			processProperty.StopCmd.Value = &stopCmd
			processProperty.StopCmd.AsDefaultValue = &blTrue
			if err := processProperty.StopCmd.Validate(); err != nil {
				return nil, fmt.Errorf("%s illegal. val:%s. err:%s", key, val, err.Error())
			}
		case "restart_cmd":
			restartCmd, ok := val.(string)
			if !ok {
				return nil, fmt.Errorf("%s not string. val:%s", key, val)
			}
			processProperty.RestartCmd.Value = &restartCmd
			processProperty.RestartCmd.AsDefaultValue = &blTrue
			if err := processProperty.RestartCmd.Validate(); err != nil {
				return nil, fmt.Errorf("%s illegal. val:%s. err:%s", key, val, err.Error())
			}
		case "face_stop_cmd":
			restartCmd, ok := val.(string)
			if !ok {
				return nil, fmt.Errorf("%s not string. val:%s", key, val)
			}
			processProperty.RestartCmd.Value = &restartCmd
			processProperty.RestartCmd.AsDefaultValue = &blTrue
			if err := processProperty.RestartCmd.Validate(); err != nil {
				return nil, fmt.Errorf("%s illegal. val:%s. err:%s", key, val, err.Error())
			}
		case "bk_func_name":
			funcName, ok := val.(string)
			if !ok {
				return nil, fmt.Errorf("%s not string. val:%s", key, val)
			}
			processProperty.FuncName.Value = &funcName
			processProperty.FuncName.AsDefaultValue = &blTrue
			if err := processProperty.FuncName.Validate(); err != nil {
				return nil, fmt.Errorf("%s illegal. val:%s. err:%s", key, val, err.Error())
			}
		case "work_path":
			workPath, ok := val.(string)
			if !ok {
				return nil, fmt.Errorf("%s not string. val:%s", key, val)
			}
			processProperty.WorkPath.Value = &workPath
			processProperty.WorkPath.AsDefaultValue = &blTrue
			if err := processProperty.WorkPath.Validate(); err != nil {
				return nil, fmt.Errorf("%s illegal. val:%s. err:%s", key, val, err.Error())
			}
		case "bind_ip":
			bindIP, ok := val.(string)
			if !ok {
				return nil, fmt.Errorf("%s not string. val:%s", key, val)
			}
			bindIPAlias := metadata.SocketBindType(bindIP)
			processProperty.BindIP.Value = &bindIPAlias
			processProperty.BindIP.AsDefaultValue = &blTrue
			if err := processProperty.BindIP.Validate(); err != nil {
				return nil, fmt.Errorf("%s illegal. val:%s. err:%s", key, val, err.Error())
			}
		case "priority":
			priority, err := util.GetInt64ByInterface(val)
			if err != nil {
				return nil, fmt.Errorf("%s not integer. val:%s", key, val)
			}
			processProperty.Priority.Value = &priority
			processProperty.Priority.AsDefaultValue = &blTrue
			if err := processProperty.Priority.Validate(); err != nil {
				return nil, fmt.Errorf("%s illegal. val:%s. err:%s", key, val, err.Error())
			}
		case "reload_cmd":
			reloadCmd, ok := val.(string)
			if !ok {
				return nil, fmt.Errorf("%s not string. val:%s", key, val)
			}
			processProperty.ReloadCmd.Value = &reloadCmd
			processProperty.ReloadCmd.AsDefaultValue = &blTrue
			if err := processProperty.ReloadCmd.Validate(); err != nil {
				return nil, fmt.Errorf("%s illegal. val:%s. err:%s", key, val, err.Error())
			}
		case "bk_process_name":
			procName, ok := val.(string)
			if !ok {
				return nil, fmt.Errorf("%s not string. val:%s", key, val)
			}
			processProperty.ProcessName.Value = &procName
			processProperty.ProcessName.AsDefaultValue = &blTrue
			if err := processProperty.ProcessName.Validate(); err != nil {
				return nil, fmt.Errorf("%s illegal. val:%s. err:%s", key, val, err.Error())
			}
		case "port":
			port, ok := val.(string)
			if !ok {
				return nil, fmt.Errorf("%s not string. val:%s", key, val)
			}
			processProperty.Port.Value = &port
			processProperty.Port.AsDefaultValue = &blTrue
			if err := processProperty.Port.Validate(); err != nil {
				return nil, fmt.Errorf("%s illegal. val:%s. err:%s", key, val, err.Error())
			}
		case "pid_file":
			pidFile, ok := val.(string)
			if !ok {
				return nil, fmt.Errorf("%s not string. val:%s", key, val)
			}
			processProperty.PidFile.Value = &pidFile
			processProperty.PidFile.AsDefaultValue = &blTrue
			if err := processProperty.PidFile.Validate(); err != nil {
				return nil, fmt.Errorf("%s illegal. val:%s. err:%s", key, val, err.Error())
			}
		case "auto_start":
			autoStart, ok := val.(bool)
			if !ok {
				return nil, fmt.Errorf("%s not boolean. val:%s", key, val)
			}
			processProperty.AutoStart.Value = &autoStart
			processProperty.AutoStart.AsDefaultValue = &blTrue
			if err := processProperty.AutoStart.Validate(); err != nil {
				return nil, fmt.Errorf("%s illegal. val:%s. err:%s", key, val, err.Error())
			}
		case "auto_time_gap":
			autoTimeGap, err := util.GetInt64ByInterface(val)
			if err != nil {
				return nil, fmt.Errorf("%s not integer. val:%s", key, val)
			}
			processProperty.AutoTimeGapSeconds.Value = &autoTimeGap
			processProperty.AutoTimeGapSeconds.AsDefaultValue = &blTrue
			if err := processProperty.AutoTimeGapSeconds.Validate(); err != nil {
				return nil, fmt.Errorf("%s illegal. val:%s. err:%s", key, val, err.Error())
			}
		case "start_cmd":
			startCmd, ok := val.(string)
			if !ok {
				return nil, fmt.Errorf("%s not string. val:%s", key, val)
			}
			processProperty.StartCmd.Value = &startCmd
			processProperty.StartCmd.AsDefaultValue = &blTrue
			if err := processProperty.StartCmd.Validate(); err != nil {
				return nil, fmt.Errorf("%s illegal. val:%s. err:%s", key, val, err.Error())
			}
		case "bk_func_id":
			funcID, ok := val.(string)
			if !ok {
				return nil, fmt.Errorf("%s not string. val:%s", key, val)
			}
			processProperty.FuncID.Value = &funcID
			processProperty.FuncID.AsDefaultValue = &blTrue
			if err := processProperty.FuncID.Validate(); err != nil {
				return nil, fmt.Errorf("%s illegal. val:%s. err:%s", key, val, err.Error())
			}
		case "user":
			user, ok := val.(string)
			if !ok {
				return nil, fmt.Errorf("%s not string. val:%s", key, val)
			}
			processProperty.User.Value = &user
			processProperty.User.AsDefaultValue = &blTrue
			if err := processProperty.User.Validate(); err != nil {
				return nil, fmt.Errorf("%s illegal. val:%s. err:%s", key, val, err.Error())
			}
		case "timeout":
			timeout, err := util.GetInt64ByInterface(val)
			if err != nil {
				return nil, fmt.Errorf("%s not integer. val:%s", key, val)
			}
			processProperty.TimeoutSeconds.Value = &timeout
			processProperty.TimeoutSeconds.AsDefaultValue = &blTrue
			if err := processProperty.TimeoutSeconds.Validate(); err != nil {
				return nil, fmt.Errorf("%s illegal. val:%s. err:%s", key, val, err.Error())
			}
		case "protocol":
			protocol, ok := val.(string)
			if !ok {
				return nil, fmt.Errorf("%s not string. val:%s", key, val)
			}
			protocalAlias := metadata.ProtocolType(protocol)
			processProperty.Protocol.Value = &protocalAlias
			processProperty.Protocol.AsDefaultValue = &blTrue
			if err := processProperty.Protocol.Validate(); err != nil {
				return nil, fmt.Errorf("%s illegal. val:%s. err:%s", key, val, err.Error())
			}
		case "description":
			desc, ok := val.(string)
			if !ok {
				return nil, fmt.Errorf("%s not string. val:%s", key, val)
			}
			processProperty.Description.Value = &desc
			processProperty.Description.AsDefaultValue = &blTrue
			if err := processProperty.Description.Validate(); err != nil {
				return nil, fmt.Errorf("%s illegal. val:%s. err:%s", key, val, err.Error())
			}
		case "bk_start_param_regex":
			regex, ok := val.(string)
			if !ok {
				return nil, fmt.Errorf("%s not string. val:%s", key, val)
			}
			processProperty.StartParamRegex.Value = &regex
			processProperty.StartParamRegex.AsDefaultValue = &blTrue
			if err := processProperty.StartParamRegex.Validate(); err != nil {
				return nil, fmt.Errorf("%s illegal. val:%s. err:%s", key, val, err.Error())
			}
		default:
			return nil, fmt.Errorf("%s illegal. val:%s", key, val)
		}

	}

	if field, err := processProperty.Validate(); err != nil {
		return nil, fmt.Errorf("process illegal. field:%s, err:%s", field, err.Error())
	}
	return processProperty, nil

}

// getSetParentID 获取set的parent id， 多个层级这个值不是业务id
func getSetParentID(ctx context.Context, bizID int64, db dal.DB) (int64, error) {

	searchCond := condition.CreateCondition()
	searchCond.Field(common.BKAppIDField).Eq(bizID)
	searchCond.Field(common.BKDefaultField).Eq(common.NormalModuleFlag)

	result := make(map[string]int64, 0)
	err := db.Table(common.BKTableNameBaseSet).Find(searchCond.ToMapStr()).Fields(common.BKParentIDField).One(ctx, result)
	if err != nil {
		return 0, fmt.Errorf("find set parent id error. err:%s", err.Error())
	}

	if result[common.BKParentIDField] == 0 {
		return 0, fmt.Errorf("set parent id = 0. illegal")
	}

	return result[common.BKParentIDField], nil

}

// getBkBizID 获取蓝鲸业务的business id
func getBKBizID(ctx context.Context, db dal.DB) (int64, error) {
	searchCond := condition.CreateCondition()
	searchCond.Field(common.BKAppNameField).Eq(common.BKAppName)
	result := make(map[string]int64, 0)
	err := db.Table(common.BKTableNameBaseApp).Find(searchCond.ToMapStr()).Fields(common.BKAppIDField).One(ctx, result)
	if err != nil {
		return 0, fmt.Errorf("find 蓝鲸 business id error. err:%s", err.Error())
	}

	return result[common.BKAppIDField], nil

}
