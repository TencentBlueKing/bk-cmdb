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
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
)

func upgradeServiceTemplate(ctx context.Context, db dal.RDB, conf *upgrader.Config) (err error) {
	categoryID, err := addDefaultCategory(ctx, db, conf)
	if err != nil {
		return fmt.Errorf("addDefaultCategory failed: %v", err)
	}
	if err = backupProcessBase(ctx, db, conf); err != nil {
		return fmt.Errorf("backupProcessBase failed: %v", err)
	}

	allmodules := []metadata.ModuleInst{}
	if err = db.Table(common.BKTableNameBaseModule).Find(condition.CreateCondition().
		Field(common.BKDefaultField).NotGt(0).ToMapStr()).
		All(ctx, &allmodules); err != nil {
		return err
	}

	// bizID:modulename:modules
	biz2Module := map[int64]map[string][]metadata.ModuleInst{}
	for _, module := range allmodules {
		_, ok := biz2Module[module.BizID]
		if !ok {
			biz2Module[module.BizID] = map[string][]metadata.ModuleInst{}
		}
		biz2Module[module.BizID][module.ModuleName] = append(biz2Module[module.BizID][module.ModuleName], module)
	}

	for bizID, bizModules := range biz2Module {
		ownerID := common.BKDefaultOwnerID
		for modulename, modules := range bizModules {
		    if len(modules) == 0 {
		        continue
            }
			// modules would always more than 0, so would never panic here
			if modules[0].SupplierAccount != "" {
				ownerID = modules[0].SupplierAccount
			}

            processMappingInModuleCond := mapstr.MapStr{common.BKAppIDField: bizID, common.BKModuleNameField: modulename}
            processMappingInModule := []metadata.ProcessModule{}
            if err = db.Table(common.BKTableNameProcModule).Find(processMappingInModuleCond).All(ctx, &processMappingInModule); err != nil {
                return err
            }
            if len(processMappingInModule) <= 0 {
                // this module does not bounded with a process, do not need to create service instance related info.
                continue
            }
            
			// build service template
			svcTemplateID, err := db.NextSequence(ctx, common.BKTableNameServiceTemplate)
			if err != nil {
				return err
			}
			serviceTemplate := metadata.ServiceTemplate{
				Metadata:          metadata.NewMetaDataFromBusinessID(strconv.FormatInt(bizID, 10)),
				ID:                int64(svcTemplateID),
				Name:              modulename,
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
			

			processIDInModule := []int64{}
			for _, mapping := range processMappingInModule {
				processIDInModule = append(processIDInModule, mapping.ProcessID)
			}
			oldProcess := []metadata.Process{}
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

			inst2ProcessInstTemplate := map[int64]metadata.ProcessTemplate{}
			for _, oldInst := range oldProcess {
				procTemplateID, err := db.NextSequence(ctx, common.BKTableNameProcessTemplate)
				if err != nil {
					return err
				}

				procTemplate := metadata.ProcessTemplate{
					Metadata:          metadata.NewMetaDataFromBusinessID(strconv.FormatInt(bizID, 10)),
					ID:                int64(procTemplateID),
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
				moduleHosts := []metadata.ModuleHost{}
				if err = db.Table(common.BKTableNameModuleHostConfig).Find(
					condition.CreateCondition().Field(common.BKModuleIDField).Eq(module.ModuleID).ToMapStr(),
				).All(ctx, &moduleHosts); err != nil {
					return err
				}

				for _, moduleHost := range moduleHosts {
					srvInstID, err := db.NextSequence(ctx, common.BKTableNameServiceInstance)
					if err != nil {
						return err
					}
					srvInst := metadata.ServiceInstance{
						Metadata:          metadata.NewMetaDataFromBusinessID(strconv.FormatInt(bizID, 10)),
						ID:                int64(srvInstID),
						Name:              modulename,
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
					for _, inst := range oldProcess {
						processTemplateID := inst2ProcessInstTemplate[inst.ProcessID].ID
						procInstID, err := db.NextSequence(ctx, common.BKTableNameBaseProcess)
						if err != nil {
							return err
						}
						inst.ProcessID = int64(procInstID)
						inst.Metadata = metadata.NewMetaDataFromBusinessID(strconv.FormatInt(bizID, 10))
						inst.CreateTime = time.Now()
						inst.LastTime = time.Now()
						if inst.BindIP != nil {
							tplBindIP := metadata.SocketBindType(*inst.BindIP)
							*inst.BindIP = tplBindIP.IP()
						} else {
							inst.BindIP = new(string)
						}
						blog.InfoJSON("procInst: %s", inst)
						if err = db.Table(common.BKTableNameBaseProcess).Insert(ctx, inst); err != nil {
							return err
						}

						// build service instance relation
						relateion := metadata.ProcessInstanceRelation{
							Metadata:          metadata.NewMetaDataFromBusinessID(strconv.FormatInt(bizID, 10)),
							ProcessID:         inst.ProcessID,
							ServiceInstanceID: srvInst.ID,
							ProcessTemplateID: processTemplateID,
							HostID:            moduleHost.HostID,
							SupplierAccount:   ownerID,
						}
						blog.InfoJSON("relation: %s", relateion)
						if err = db.Table(common.BKTableNameProcessInstanceRelation).Insert(ctx, relateion); err != nil {
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
		process := []mapstr.MapStr{}
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

func procInstToProcTemplate(inst metadata.Process) *metadata.ProcessProperty {
	var True = true
	template := metadata.ProcessProperty{}
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
