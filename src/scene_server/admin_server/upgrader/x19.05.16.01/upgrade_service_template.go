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

	allmodules := []metadata.ModuleInst{}
	if err = db.Table(common.BKTableNameBaseModule).Find(nil).All(ctx, &allmodules); err != nil {
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
		ownerID := ""
		for modulename, modules := range bizModules {
			// modules would always more than 0, so would never panic here
			ownerID = modules[0].SupplierAccount
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
			blog.InfoJSON("serviceTemplate: %s", serviceTemplate)
			if err = db.Table(common.BKTableNameServiceTemplate).Insert(ctx, serviceTemplate); err != nil {
				return err
			}

			// build process template
			processMappingInModuleCond := mapstr.MapStr{common.BKAppIDField: bizID, common.BKModuleNameField: modulename}
			processMappingInModule := []metadata.ProcessModule{}
			if err = db.Table(common.BKTableNameProcModule).Find(processMappingInModuleCond).All(ctx, &processMappingInModule); err != nil {
				return err
			}
			if len(processMappingInModule) <= 0 {
				continue
			}

			processIDInModule := []int64{}
			for _, mapping := range processMappingInModule {
				processIDInModule = append(processIDInModule, mapping.ProcessID)
			}
			oldProcess := []metadata.Process{}
			if err = db.Table(common.BKTableNameBaseProcess).Find(
				condition.CreateCondition().Field(common.BKProcessIDField).In(processIDInModule).
					Field(common.BKAppIDField).Eq(bizID).ToMapStr(),
			).All(ctx, &oldProcess); err != nil {
				return err
			}
			if len(oldProcess) <= 0 {
				continue
			}

			inst2ProcessInstTemplate := map[int64]metadata.ProcessTemplate{}
			for _, oldInst := range oldProcess {
				procTemplateID, err := db.NextSequence(ctx, common.BKTableNameProcessTemplate)
				if err != nil {
					return err
				}

				procTemplate := metadata.ProcessTemplate{
					ID:                int64(procTemplateID),
					ServiceTemplateID: serviceTemplate.ID,
					Template:          procInstToProcTemplate(oldInst),
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
						blog.InfoJSON("procInst: %s", inst)
						if err = db.Table(common.BKTableNameBaseProcess).Insert(ctx, inst); err != nil {
							return err
						}

						// build service instance relation
						relateion := metadata.ServiceInstanceRelations{
							ProcessID:         inst.ProcessID,
							ServiceInstanceID: srvInst.ID,
							ProcessTemplateID: processTemplateID,
							HostID:            moduleHost.HostID,
							SupplierAccount:   ownerID,
						}
						blog.InfoJSON("relateion: %s", relateion)
						if err = db.Table(common.BKTableNameServiceInstanceRelations).Insert(ctx, relateion); err != nil {
							return err
						}
					}
				}
			}
			blog.Info("done \n")
		}
	}

	return nil
}

func procInstToProcTemplate(inst metadata.Process) metadata.ProcessProperty {
	template := metadata.ProcessProperty{}
	if inst.ProcNum > 0 {
		template.ProcNum.Value = inst.ProcNum
		template.ProcNum.AsDefaultValue = true
	}
	if inst.StopCmd != "" {
		template.StopCmd.Value = inst.StopCmd
		template.StopCmd.AsDefaultValue = true
	}
	if inst.RestartCmd != "" {
		template.RestartCmd.Value = inst.RestartCmd
		template.RestartCmd.AsDefaultValue = true
	}
	if inst.ForceStopCmd != "" {
		template.ForceStopCmd.Value = inst.ForceStopCmd
		template.ForceStopCmd.AsDefaultValue = true
	}
	if inst.FuncName != "" {
		template.FuncName.Value = inst.FuncName
		template.FuncName.AsDefaultValue = true
	}
	if inst.WorkPath != "" {
		template.WorkPath.Value = inst.WorkPath
		template.WorkPath.AsDefaultValue = true
	}
	if inst.BindIP != "" {
		template.BindIP.Value = inst.BindIP
		template.BindIP.AsDefaultValue = true
	}
	if inst.Priority > 0 {
		template.Priority.Value = inst.Priority
		template.Priority.AsDefaultValue = true
	}
	if inst.ReloadCmd != "" {
		template.ReloadCmd.Value = inst.ReloadCmd
		template.ReloadCmd.AsDefaultValue = true
	}
	if inst.ProcessName != "" {
		template.ProcessName.Value = inst.ProcessName
		template.ProcessName.AsDefaultValue = true
	}
	if inst.Port != "" {
		template.Port.Value = inst.Port
		template.Port.AsDefaultValue = true
	}
	if inst.PidFile != "" {
		template.PidFile.Value = inst.PidFile
		template.PidFile.AsDefaultValue = true
	}
	if inst.AutoStart == true {
		template.AutoStart.Value = inst.AutoStart
		template.AutoStart.AsDefaultValue = true
	}
	if inst.AutoTimeGap > 0 {
		template.AutoTimeGap.Value = inst.AutoTimeGap
		template.AutoTimeGap.AsDefaultValue = true
	}
	if inst.StartCmd != "" {
		template.StartCmd.Value = inst.StartCmd
		template.StartCmd.AsDefaultValue = true
	}
	if inst.FuncID != "" {
		template.FuncID.Value = inst.FuncID
		template.FuncID.AsDefaultValue = true
	}
	if inst.User != "" {
		template.User.Value = inst.User
		template.User.AsDefaultValue = true
	}
	if inst.TimeoutSeconds > 0 {
		template.TimeoutSeconds.Value = inst.TimeoutSeconds
		template.TimeoutSeconds.AsDefaultValue = true
	}
	if inst.Protocol != "" {
		template.Protocol.Value = inst.Protocol
		template.Protocol.AsDefaultValue = true
	}
	if inst.Description != "" {
		template.Description.Value = inst.Description
		template.Description.AsDefaultValue = true
	}

	return template
}
