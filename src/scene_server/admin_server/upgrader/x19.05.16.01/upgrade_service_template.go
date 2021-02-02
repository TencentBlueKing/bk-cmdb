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

package x19_05_16_01

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/selector"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/dal/mongo/local"

	"go.mongodb.org/mongo-driver/mongo/options"
)

type ProcessInstanceRelation struct {
	Metadata metadata.Metadata `field:"metadata" json:"metadata" bson:"metadata"`

	// unique field, 1:1 mapping with ProcessInstance.
	ProcessID         int64 `field:"bk_process_id" json:"bk_process_id" bson:"bk_process_id"`
	ServiceInstanceID int64 `field:"service_instance_id" json:"service_instance_id" bson:"service_instance_id"`

	// ProcessTemplateID indicate which template are current process instantiate from.
	ProcessTemplateID int64 `field:"process_template_id" json:"process_template_id" bson:"process_template_id"`

	// redundant field for accelerating processes by HostID
	HostID          int64  `field:"bk_host_id" json:"bk_host_id" bson:"bk_host_id"`
	SupplierAccount string `field:"bk_supplier_account" json:"bk_supplier_account" bson:"bk_supplier_account"`
}

type ServiceTemplate struct {
	Metadata metadata.Metadata `field:"metadata" json:"metadata" bson:"metadata"`

	ID int64 `field:"id" json:"id,omitempty" bson:"id"`
	// name of this service, can not be empty
	Name string `field:"name" json:"name,omitempty" bson:"name"`

	// the class of this service, each field means a class label.
	// now, the class must have two labels.
	ServiceCategoryID int64 `field:"service_category_id" json:"service_category_id,omitempty" bson:"service_category_id"`

	Creator         string    `field:"creator" json:"creator,omitempty" bson:"creator"`
	Modifier        string    `field:"modifier" json:"modifier,omitempty" bson:"modifier"`
	CreateTime      time.Time `field:"create_time" json:"create_time,omitempty" bson:"create_time"`
	LastTime        time.Time `field:"last_time" json:"last_time,omitempty" bson:"last_time"`
	SupplierAccount string    `field:"bk_supplier_account" json:"bk_supplier_account,omitempty" bson:"bk_supplier_account"`
}

type ProcessTemplate struct {
	Metadata metadata.Metadata `field:"metadata" json:"metadata" bson:"metadata"`

	ID          int64  `field:"id" json:"id,omitempty" bson:"id"`
	ProcessName string `field:"bk_process_name" json:"bk_process_name" bson:"bk_process_name"`
	// the service template's, which this process template belongs to.
	ServiceTemplateID int64 `field:"service_template_id" json:"service_template_id" bson:"service_template_id"`

	// stores a process instance's data includes all the process's
	// properties's value.
	Property *ProcessProperty `field:"property" json:"property,omitempty" bson:"property"`

	Creator         string    `field:"creator" json:"creator,omitempty" bson:"creator"`
	Modifier        string    `field:"modifier" json:"modifier,omitempty" bson:"modifier"`
	CreateTime      time.Time `field:"create_time" json:"create_time,omitempty" bson:"create_time"`
	LastTime        time.Time `field:"last_time" json:"last_time,omitempty" bson:"last_time"`
	SupplierAccount string    `field:"bk_supplier_account" json:"bk_supplier_account,omitempty" bson:"bk_supplier_account"`
}

type ServiceInstance struct {
	Metadata metadata.Metadata `field:"metadata" json:"metadata" bson:"metadata"`
	ID       int64             `field:"id" json:"id,omitempty" bson:"id"`
	Name     string            `field:"name" json:"name,omitempty" bson:"name"`
	Labels   selector.Labels   `field:"labels" json:"labels,omitempty" bson:"labels"`

	// the template id can not be updated, once the service is created.
	// it can be 0 when the service is not created with a service template.
	ServiceTemplateID int64 `field:"service_template_id" json:"service_template_id,omitempty" bson:"service_template_id"`
	HostID            int64 `field:"bk_host_id" json:"bk_host_id,omitempty" bson:"bk_host_id"`

	// the module that this service belongs to.
	ModuleID int64 `field:"bk_module_id" json:"bk_module_id,omitempty" bson:"bk_module_id"`

	Creator         string    `field:"creator" json:"creator,omitempty" bson:"creator"`
	Modifier        string    `field:"modifier" json:"modifier,omitempty" bson:"modifier"`
	CreateTime      time.Time `field:"create_time" json:"create_time,omitempty" bson:"create_time"`
	LastTime        time.Time `field:"last_time" json:"last_time,omitempty" bson:"last_time"`
	SupplierAccount string    `field:"bk_supplier_account" json:"bk_supplier_account,omitempty" bson:"bk_supplier_account"`
}

type Process struct {
	Metadata        metadata.Metadata      `field:"metadata" json:"metadata" bson:"metadata"`
	ProcNum         *int64                 `field:"proc_num" json:"proc_num,omitempty" bson:"proc_num" structs:"proc_num"`
	StopCmd         *string                `field:"stop_cmd" json:"stop_cmd,omitempty" bson:"stop_cmd" structs:"stop_cmd"`
	RestartCmd      *string                `field:"restart_cmd" json:"restart_cmd,omitempty" bson:"restart_cmd" structs:"restart_cmd"`
	ForceStopCmd    *string                `field:"face_stop_cmd" json:"face_stop_cmd,omitempty" bson:"face_stop_cmd" structs:"face_stop_cmd"`
	ProcessID       int64                  `field:"bk_process_id" json:"bk_process_id,omitempty" bson:"bk_process_id" structs:"bk_process_id"`
	FuncName        *string                `field:"bk_func_name" json:"bk_func_name,omitempty" bson:"bk_func_name" structs:"bk_func_name"`
	WorkPath        *string                `field:"work_path" json:"work_path,omitempty" bson:"work_path" structs:"work_path"`
	BindIP          *string                `field:"bind_ip" json:"bind_ip,omitempty" bson:"bind_ip" structs:"bind_ip"`
	Priority        *int64                 `field:"priority" json:"priority,omitempty" bson:"priority" structs:"priority"`
	ReloadCmd       *string                `field:"reload_cmd" json:"reload_cmd,omitempty" bson:"reload_cmd" structs:"reload_cmd"`
	ProcessName     *string                `field:"bk_process_name" json:"bk_process_name,omitempty" bson:"bk_process_name" structs:"bk_process_name"`
	Port            *string                `field:"port" json:"port,omitempty" bson:"port" structs:"port"`
	PidFile         *string                `field:"pid_file" json:"pid_file,omitempty" bson:"pid_file" structs:"pid_file"`
	AutoStart       *bool                  `field:"auto_start" json:"auto_start,omitempty" bson:"auto_start" structs:"auto_start"`
	AutoTimeGap     *int64                 `field:"auto_time_gap" json:"auto_time_gap,omitempty" bson:"auto_time_gap" structs:"auto_time_gap"`
	LastTime        time.Time              `field:"last_time" json:"last_time,omitempty" bson:"last_time" structs:"last_time"`
	CreateTime      time.Time              `field:"create_time" json:"create_time,omitempty" bson:"create_time" structs:"create_time"`
	BusinessID      int64                  `field:"bk_biz_id" json:"bk_biz_id,omitempty" bson:"bk_biz_id" structs:"bk_biz_id"`
	StartCmd        *string                `field:"start_cmd" json:"start_cmd,omitempty" bson:"start_cmd" structs:"start_cmd"`
	FuncID          *string                `field:"bk_func_id" json:"bk_func_id,omitempty" bson:"bk_func_id" structs:"bk_func_id"`
	User            *string                `field:"user" json:"user,omitempty" bson:"user" structs:"user"`
	TimeoutSeconds  *int64                 `field:"timeout" json:"timeout,omitempty" bson:"timeout" structs:"timeout"`
	Protocol        *metadata.ProtocolType `field:"protocol" json:"protocol,omitempty" bson:"protocol" structs:"protocol"`
	Description     *string                `field:"description" json:"description,omitempty" bson:"description" structs:"description"`
	SupplierAccount string                 `field:"bk_supplier_account" json:"bk_supplier_account,omitempty" bson:"bk_supplier_account" structs:"bk_supplier_account"`
	StartParamRegex *string                `field:"bk_start_param_regex" json:"bk_start_param_regex,omitempty" bson:"bk_start_param_regex,omitempty" structs:"bk_start_param_regex"`
}

func upgradeServiceTemplate(ctx context.Context, db dal.RDB, conf *upgrader.Config) (err error) {
	mongo, ok := db.(*local.Mongo)
	if !ok {
		return fmt.Errorf("db is not *local.Mongo type")
	}
	dbc := mongo.GetDBClient()

	categoryID, err := addDefaultCategory(ctx, db, conf)
	if err != nil {
		return fmt.Errorf("addDefaultCategory failed: %v", err)
	}
	if err = backupProcessBase(ctx, db, conf); err != nil {
		return fmt.Errorf("backupProcessBase failed: %v", err)
	}

	allmodules := make([]metadata.ModuleInst, 0)
	if err = db.Table(common.BKTableNameBaseModule).Find(condition.CreateCondition().
		Field(common.BKDefaultField).Eq(common.DefaultFlagDefaultValue).ToMapStr()).
		All(ctx, &allmodules); err != nil {
		return err
	}

	// bizID:moduleName:modules
	biz2Module := map[int64]map[string][]metadata.ModuleInst{}
	for _, module := range allmodules {
		_, ok := biz2Module[module.BizID]
		if !ok {
			biz2Module[module.BizID] = map[string][]metadata.ModuleInst{}
		}
		biz2Module[module.BizID][module.ModuleName] = append(biz2Module[module.BizID][module.ModuleName], module)
	}

	hostMap := make(map[int64]map[string]interface{}, 0)

	for bizID, bizModules := range biz2Module {
		ownerID := common.BKDefaultOwnerID

		svcTemplateIDs, err := db.NextSequences(ctx, common.BKTableNameServiceTemplate, len(bizModules))
		if err != nil {
			return err
		}
		svcTemplateIndex := 0

		for moduleName, modules := range bizModules {
			if len(modules) == 0 {
				continue
			}
			// modules would always more than 0, so would never panic here
			if modules[0].SupplierAccount != "" {
				ownerID = modules[0].SupplierAccount
			}

			processMappingInModuleCond := mapstr.MapStr{common.BKAppIDField: bizID, common.BKModuleNameField: moduleName}
			processMappingInModule := make([]metadata.ProcessModule, 0)
			if err = db.Table("cc_Proc2Module").Find(processMappingInModuleCond).All(ctx, &processMappingInModule); err != nil {
				return err
			}
			if len(processMappingInModule) <= 0 {
				// this module does not bounded with a process, do not need to create service instance related info.
				continue
			}

			svcTemplateID := svcTemplateIDs[svcTemplateIndex]
			svcTemplateIndex++

			// build service template
			serviceTemplate := ServiceTemplate{
				Metadata:          metadata.NewMetadata(bizID),
				ID:                int64(svcTemplateID),
				Name:              moduleName,
				ServiceCategoryID: categoryID,
				Creator:           conf.User,
				Modifier:          conf.User,
				CreateTime:        time.Now(),
				LastTime:          time.Now(),
				SupplierAccount:   ownerID,
			}
			if err = db.Table(common.BKTableNameServiceTemplate).Insert(ctx, serviceTemplate); err != nil {
				return err
			}

			// set module template
			moduleFilter := map[string]interface{}{
				common.BKModuleIDField: modules[0].ModuleID,
			}
			moduleUpdateData := map[string]interface{}{
				common.BKServiceCategoryIDField: categoryID,
				common.BKServiceTemplateIDField: svcTemplateID,
			}
			if err = db.Table(common.BKTableNameBaseModule).Update(ctx, moduleFilter, moduleUpdateData); err != nil {
				return err
			}

			processIDInModule := make([]int64, 0)
			for _, mapping := range processMappingInModule {
				processIDInModule = append(processIDInModule, mapping.ProcessID)
			}
			oldProcess := make([]Process, 0)
			processBaseCond := condition.CreateCondition().Field(common.BKProcessIDField).In(processIDInModule).
				Field(common.BKAppIDField).Eq(bizID).ToMapStr()
			if err = db.Table(common.BKTableNameBaseProcess).Find(processBaseCond).All(ctx, &oldProcess); err != nil {
				blog.Errorf("find process failed: %v %v", processBaseCond, err)
				return err
			}
			if len(oldProcess) <= 0 {
				// no process in this bounded module,
				// normally, this can not be happen.
				continue
			}

			inst2ProcessInstTemplate := map[int64]ProcessTemplate{}
			procTemplateIDs, err := db.NextSequences(ctx, common.BKTableNameProcessTemplate, len(oldProcess))
			if err != nil {
				return err
			}

			for index, oldInst := range oldProcess {

				procName := ""
				if oldInst.ProcessName != nil {
					procName = *oldInst.ProcessName
				}
				if procName == "" && oldInst.FuncName != nil {
					procName = *oldInst.FuncName
				}

				procTemplate := ProcessTemplate{
					Metadata:          metadata.NewMetadata(bizID),
					ID:                int64(procTemplateIDs[index]),
					ProcessName:       procName,
					ServiceTemplateID: serviceTemplate.ID,
					Property:          procInstToProcTemplate(oldInst),
					Creator:           conf.User,
					Modifier:          conf.User,
					CreateTime:        time.Now(),
					LastTime:          time.Now(),
				}
				inst2ProcessInstTemplate[oldInst.ProcessID] = procTemplate
				blog.InfoJSON("procTemplate: %s", procTemplate)
				if err = db.Table(common.BKTableNameProcessTemplate).Insert(ctx, procTemplate); err != nil {
					return err
				}
			}

			// build service instance
			for _, module := range modules {
				moduleHosts := make([]metadata.ModuleHost, 0)
				if err = db.Table(common.BKTableNameModuleHostConfig).Find(
					condition.CreateCondition().Field(common.BKModuleIDField).Eq(module.ModuleID).ToMapStr(),
				).All(ctx, &moduleHosts); err != nil {
					return err
				}

				srvInstIDs, err := db.NextSequences(ctx, common.BKTableNameServiceInstance, len(moduleHosts))
				if err != nil {
					return err
				}

				for index, moduleHost := range moduleHosts {
					srvInst := ServiceInstance{
						Metadata:          metadata.NewMetadata(bizID),
						ID:                int64(srvInstIDs[index]),
						Name:              moduleName,
						ServiceTemplateID: serviceTemplate.ID,
						HostID:            moduleHost.HostID,
						ModuleID:          module.ModuleID,
						SupplierAccount:   ownerID,
						Creator:           conf.User,
						Modifier:          conf.User,
						CreateTime:        time.Now(),
						LastTime:          time.Now(),
					}
					blog.InfoJSON("srvInst: %s", srvInst)
					if err = db.Table(common.BKTableNameServiceInstance).Insert(ctx, srvInst); err != nil {
						return err
					}

					// build process instance
					procInstIDs, err := db.NextSequences(ctx, common.BKTableNameBaseProcess, len(oldProcess))
					if err != nil {
						return err
					}

					for index, inst := range oldProcess {
						processTemplateID := inst2ProcessInstTemplate[inst.ProcessID].ID
						inst.ProcessID = int64(procInstIDs[index])
						inst.Metadata = metadata.NewMetaDataFromBusinessID(strconv.FormatInt(bizID, 10))
						inst.BusinessID = bizID
						inst.CreateTime = time.Now()
						inst.LastTime = time.Now()
						if inst.BindIP != nil {
							tplBindIP := metadata.SocketBindType(*inst.BindIP)

							if tplBindIP == metadata.BindInnerIP || tplBindIP == metadata.BindOuterIP {
								if hostMap[moduleHost.HostID] == nil {
									findOpts := &options.FindOptions{}
									findOpts.SetLimit(1)
									findOpts.SetProjection(map[string]int{common.BKHostInnerIPField: 1, common.BKHostOuterIPField: 1})
									filter := map[string]interface{}{common.BKHostIDField: moduleHost.HostID}
									host := make(map[string]interface{})

									cursor, err := dbc.Database(mongo.GetDBName()).Collection(common.BKTableNameBaseHost).Find(ctx, filter, findOpts)
									if err != nil {
										blog.Errorf("find host %d failed, err: %s", moduleHost.HostID, err.Error())
										return err
									}

									if !cursor.Next(ctx) {
										return fmt.Errorf("host %d not exist", moduleHost.HostID)
									}

									if err := cursor.Decode(&host); err != nil {
										blog.Errorf("decode host %d failed, err: %s", moduleHost.HostID, err.Error())
										return err
									}
									hostMap[moduleHost.HostID] = host
								}
							}

							bindIP, err := tplBindIP.IP(hostMap[moduleHost.HostID])
							if err != nil {
								return err
							}

							*inst.BindIP = bindIP
						} else {
							inst.BindIP = new(string)
						}
						inst.SupplierAccount = ownerID
						blog.InfoJSON("procInst: %s", inst)
						if err = db.Table(common.BKTableNameBaseProcess).Insert(ctx, inst); err != nil {
							return err
						}

						// build service instance relation
						relation := ProcessInstanceRelation{
							Metadata:          metadata.NewMetaDataFromBusinessID(strconv.FormatInt(bizID, 10)),
							ProcessID:         inst.ProcessID,
							ServiceInstanceID: srvInst.ID,
							ProcessTemplateID: processTemplateID,
							HostID:            moduleHost.HostID,
							SupplierAccount:   ownerID,
						}
						blog.InfoJSON("relation: %s", relation)
						if err = db.Table(common.BKTableNameProcessInstanceRelation).Insert(ctx, relation); err != nil {
							return err
						}
					}
				}
			}
			blog.Info("done \n")
		}
	}

	// 填充默认值：service_template_id, service_category_id
	notSetFilter := map[string]interface{}{
		common.BKServiceCategoryIDField: map[string]interface{}{
			common.BKDBExists: false,
		},
	}
	defaultData := map[string]interface{}{
		common.BKServiceCategoryIDField: categoryID,
		common.BKServiceTemplateIDField: 0,
	}
	if err = db.Table(common.BKTableNameBaseModule).Update(ctx, notSetFilter, defaultData); err != nil {
		return err
	}

	return db.Table(common.BKTableNameBaseProcess).Delete(ctx, mapstr.MapStr{"old_flag": true})
}

func backupProcessBase(ctx context.Context, db dal.RDB, conf *upgrader.Config) (err error) {
	start := uint64(0)
	limit := uint64(100)
	for {
		process := make([]mapstr.MapStr, 0)
		if err := db.Table(common.BKTableNameBaseProcess).Find(nil).Start(start).Limit(limit).
			All(ctx, &process); err != nil {
			return err
		}
		if len(process) <= 0 {
			break
		}
		if err := db.Table("cc_Process_backup").Insert(ctx, process); err != nil {
			return err
		}
		start += limit
	}
	return db.Table(common.BKTableNameBaseProcess).Update(ctx, nil, mapstr.MapStr{"old_flag": true})
}

func procInstToProcTemplate(inst Process) *ProcessProperty {
	var True = true
	template := ProcessProperty{}
	if inst.ProcNum != nil && *inst.ProcNum > 0 {
		template.ProcNum.Value = inst.ProcNum
		template.ProcNum.AsDefaultValue = &True
	}
	if inst.StopCmd != nil && len(*inst.StopCmd) > 0 {
		template.StopCmd.Value = inst.StopCmd
		template.StopCmd.AsDefaultValue = &True
	}
	if inst.RestartCmd != nil && len(*inst.RestartCmd) > 0 {
		template.RestartCmd.Value = inst.RestartCmd
		template.RestartCmd.AsDefaultValue = &True
	}
	if inst.ForceStopCmd != nil && len(*inst.ForceStopCmd) > 0 {
		template.ForceStopCmd.Value = inst.ForceStopCmd
		template.ForceStopCmd.AsDefaultValue = &True
	}
	if inst.FuncName != nil && len(*inst.FuncName) > 0 {
		template.FuncName.Value = inst.FuncName
		template.FuncName.AsDefaultValue = &True
	} else if inst.ProcessName != nil && len(*inst.ProcessName) > 0 {
		// FuncName empty, try use ProcessName
		template.FuncName.Value = inst.ProcessName
		template.FuncName.AsDefaultValue = &True
	}
	if inst.WorkPath != nil && len(*inst.WorkPath) > 0 {
		template.WorkPath.Value = inst.WorkPath
		template.WorkPath.AsDefaultValue = &True
	}
	if inst.BindIP != nil {
		template.BindIP.Value = new(metadata.SocketBindType)
		*template.BindIP.Value = metadata.SocketBindType(*inst.BindIP)
		template.BindIP.AsDefaultValue = &True
	}
	if inst.Priority != nil && *inst.Priority > 0 {
		template.Priority.Value = inst.Priority
		template.Priority.AsDefaultValue = &True
	}
	if inst.ReloadCmd != nil && len(*inst.ReloadCmd) > 0 {
		template.ReloadCmd.Value = inst.ReloadCmd
		template.ReloadCmd.AsDefaultValue = &True
	}
	if inst.ProcessName != nil && len(*inst.ProcessName) > 0 {
		template.ProcessName.Value = inst.ProcessName
		template.ProcessName.AsDefaultValue = &True
	}
	if inst.Port != nil && len(*inst.Port) > 0 {
		template.Port.Value = inst.Port
		template.Port.AsDefaultValue = &True
	}
	if inst.PidFile != nil && len(*inst.PidFile) > 0 {
		template.PidFile.Value = inst.PidFile
		template.PidFile.AsDefaultValue = &True
	}
	if inst.AutoStart != nil {
		template.AutoStart.Value = inst.AutoStart
		template.AutoStart.AsDefaultValue = &True
	}
	if inst.AutoTimeGap != nil && *inst.AutoTimeGap > 0 {
		template.AutoTimeGapSeconds.Value = inst.AutoTimeGap
		template.AutoTimeGapSeconds.AsDefaultValue = &True
	}
	if inst.StartCmd != nil && len(*inst.StartCmd) > 0 {
		template.StartCmd.Value = inst.StartCmd
		template.StartCmd.AsDefaultValue = &True
	}
	if inst.FuncID != nil && len(*inst.FuncID) > 0 {
		template.FuncID.Value = inst.FuncID
		template.FuncID.AsDefaultValue = &True
	}
	if inst.User != nil && len(*inst.User) > 0 {
		template.User.Value = inst.User
		template.User.AsDefaultValue = &True
	}
	if inst.TimeoutSeconds != nil && *inst.TimeoutSeconds > 0 {
		template.TimeoutSeconds.Value = inst.TimeoutSeconds
		template.TimeoutSeconds.AsDefaultValue = &True
	}
	if inst.Protocol != nil && inst.Protocol.String() != "" {
		template.Protocol.Value = inst.Protocol
		template.Protocol.AsDefaultValue = &True
	}
	if inst.Description != nil && len(*inst.Description) > 0 {
		template.Description.Value = inst.Description
		template.Description.AsDefaultValue = &True
	}
	if inst.StartParamRegex != nil && len(*inst.StartParamRegex) > 0 {
		template.StartParamRegex.Value = inst.StartParamRegex
		template.StartParamRegex.AsDefaultValue = &True
	}

	return &template
}

type ProcessProperty struct {
	ProcNum            metadata.PropertyInt64    `field:"proc_num" json:"proc_num" bson:"proc_num" validate:"max=10000,min=1"`
	StopCmd            metadata.PropertyString   `field:"stop_cmd" json:"stop_cmd" bson:"stop_cmd"`
	RestartCmd         metadata.PropertyString   `field:"restart_cmd" json:"restart_cmd" bson:"restart_cmd"`
	ForceStopCmd       metadata.PropertyString   `field:"face_stop_cmd" json:"face_stop_cmd" bson:"face_stop_cmd"`
	FuncName           metadata.PropertyString   `field:"bk_func_name" json:"bk_func_name" bson:"bk_func_name" validate:"required"`
	WorkPath           metadata.PropertyString   `field:"work_path" json:"work_path" bson:"work_path"`
	BindIP             metadata.PropertyBindIP   `field:"bind_ip" json:"bind_ip" bson:"bind_ip"`
	Priority           metadata.PropertyInt64    `field:"priority" json:"priority" bson:"priority" validate:"max=10000,min=1"`
	ReloadCmd          metadata.PropertyString   `field:"reload_cmd" json:"reload_cmd" bson:"reload_cmd"`
	ProcessName        metadata.PropertyString   `field:"bk_process_name" json:"bk_process_name" bson:"bk_process_name" validate:"required"`
	Port               metadata.PropertyPort     `field:"port" json:"port" bson:"port"`
	PidFile            metadata.PropertyString   `field:"pid_file" json:"pid_file" bson:"pid_file"`
	AutoStart          metadata.PropertyBool     `field:"auto_start" json:"auto_start" bson:"auto_start"`
	AutoTimeGapSeconds metadata.PropertyInt64    `field:"auto_time_gap" json:"auto_time_gap" bson:"auto_time_gap" validate:"max=10000,min=1"`
	StartCmd           metadata.PropertyString   `field:"start_cmd" json:"start_cmd" bson:"start_cmd"`
	FuncID             metadata.PropertyString   `field:"bk_func_id" json:"bk_func_id" bson:"bk_func_id"`
	User               metadata.PropertyString   `field:"user" json:"user" bson:"user"`
	TimeoutSeconds     metadata.PropertyInt64    `field:"timeout" json:"timeout" bson:"timeout" validate:"max=10000,min=1"`
	Protocol           metadata.PropertyProtocol `field:"protocol" json:"protocol" bson:"protocol"`
	Description        metadata.PropertyString   `field:"description" json:"description" bson:"description"`
	StartParamRegex    metadata.PropertyString   `field:"bk_start_param_regex" json:"bk_start_param_regex" bson:"bk_start_param_regex"`
	PortEnable         metadata.PropertyBool     `field:"bk_enable_port" json:"bk_enable_port" bson:"bk_enable_port"`
	GatewayIP          metadata.PropertyString   `field:"bk_gateway_ip" json:"bk_gateway_ip" bson:"bk_gateway_ip"`
	GatewayPort        metadata.PropertyString   `field:"bk_gateway_port" json:"bk_gateway_port" bson:"bk_gateway_port"`
	GatewayProtocol    metadata.PropertyProtocol `field:"bk_gateway_protocol" json:"bk_gateway_protocol" bson:"bk_gateway_protocol"`
	GatewayCity        metadata.PropertyString   `field:"bk_gateway_city" json:"bk_gateway_city" bson:"bk_gateway_city"`
}
