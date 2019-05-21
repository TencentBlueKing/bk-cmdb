/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.,
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the ",License",); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an ",AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package host

import (
	"strconv"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/source_controller/coreservice/core"
)

type identifier struct {
	hostManager *hostManager
	hosts       []metadata.HostIdentifier

	// tmp data
	setIDs            []int64
	moduleIDs         []int64
	bizIDs            []int64
	sets              map[int64]metadata.SetInst
	modules           map[int64]*metadata.ModuleInst
	bizs              map[int64]metadata.BizInst
	clouds            map[int64]metadata.CloudInst
	procs             map[int64]metadata.HostIdentProcess
	modulehosts       map[int64][]metadata.ModuleHost
	hostProcRealtions map[int64][]metadata.ProcessInstanceRelation
}

func (h *hostManager) Identifier(ctx core.ContextParams, input *metadata.SearchHostIdentifierParam) ([]metadata.HostIdentifier, error) {
	identifier := h.NewIdentifier()

	host, err := identifier.Identifier(ctx, input.HostIDs)
	if err != nil {
		blog.ErrorJSON("Identifier get host identifier error. err:%s, input:%s, rid:%s", err.Error(), input, ctx.ReqID)
		return nil, err
	}
	return host, nil
}

func (h *hostManager) NewIdentifier() *identifier {
	return &identifier{
		hostManager:       h,
		modulehosts:       make(map[int64][]metadata.ModuleHost, 0),
		sets:              make(map[int64]metadata.SetInst, 0),
		bizs:              make(map[int64]metadata.BizInst, 0),
		clouds:            make(map[int64]metadata.CloudInst, 0),
		procs:             make(map[int64]metadata.HostIdentProcess, 0),
		hostProcRealtions: make(map[int64][]metadata.ProcessInstanceRelation, 0),
	}
}

func (i *identifier) Identifier(ctx core.ContextParams, hostIDs []int64) ([]metadata.HostIdentifier, error) {
	err := i.findHost(ctx, hostIDs)
	if err != nil {
		return nil, err
	}

	err = i.findModuleHostRelation(ctx, hostIDs)
	if err != nil {
		return nil, err
	}

	err = i.findHostServiceInst(ctx, hostIDs)
	if err != nil {
		return nil, err
	}

	err = i.identifierHostTopoInfo(ctx)
	if err != nil {
		return nil, err
	}
	err = i.identifierHostCloud(ctx)
	if err != nil {
		return nil, err
	}

	i.build(ctx)
	return i.hosts, nil
}

// FindHost query host info
func (i *identifier) findHost(ctx core.ContextParams, hostIDs []int64) error {
	hostCond := condition.CreateCondition().Field(common.BKHostIDField).In(hostIDs)
	condHostMap := util.SetQueryOwner(hostCond.ToMapStr(), ctx.SupplierAccount)
	// fetch all hosts
	i.hosts = make([]metadata.HostIdentifier, 0) // New func init
	err := i.hostManager.DbProxy.Table(common.BKTableNameBaseHost).Find(condHostMap).All(ctx, &i.hosts)
	if err != nil {
		blog.ErrorJSON("findHost query host error. err:%s, conidtion:%s, rid:%s", err.Error(), condHostMap, ctx.ReqID)
		return ctx.Error.Error(common.CCErrCommDBSelectFailed)
	}

	blog.V(5).Infof("findHost query host info. host:%#v, rid;%s", i.hosts, ctx.ReqID)

	return nil
}

// findModuleHostRelation query host and module relation
func (i *identifier) findModuleHostRelation(ctx core.ContextParams, hostIDs []int64) error {
	hostModuleCond := condition.CreateCondition().Field(common.BKHostIDField).In(hostIDs)
	condModuleHostMap := util.SetQueryOwner(hostModuleCond.ToMapStr(), ctx.SupplierAccount)
	// fetch  host and module relation
	moduleHostRelation := []metadata.ModuleHost{}
	err := i.hostManager.DbProxy.Table(common.BKTableNameModuleHostConfig).Find(condModuleHostMap).All(ctx, &moduleHostRelation)
	if err != nil {
		blog.ErrorJSON("findModuleHostRelation query host and module relation error. err:%s, conidtion:%s, rid:%s", err.Error(), condModuleHostMap, ctx.ReqID)
		return ctx.Error.Error(common.CCErrCommDBSelectFailed)
	}

	blog.V(5).Infof("findModuleHostRelation query host and module relation. relation:%#v, rid;%s", i.hosts, ctx.ReqID)

	for _, modulehost := range moduleHostRelation {
		i.modulehosts[modulehost.HostID] = append(i.modulehosts[modulehost.HostID], modulehost)
		i.setIDs = append(i.setIDs, modulehost.SetID)
		i.moduleIDs = append(i.moduleIDs, modulehost.ModuleID)
		i.bizIDs = append(i.bizIDs, modulehost.AppID)
	}

	return nil
}

// identifierHostTopoInfo handle host biz,set, module info
func (i *identifier) identifierHostTopoInfo(ctx core.ContextParams) error {

	// fetch set info
	if len(i.setIDs) > 0 {
		setInfoArr := []metadata.SetInst{}
		cond := condition.CreateCondition().Field(common.BKSetIDField).In(i.setIDs)
		err := i.hostManager.dbExecQuery(ctx, common.BKTableNameBaseSet, nil, cond.ToMapStr(), &setInfoArr)
		if err != nil {
			blog.Errorf("identifierHostTopoInfo query set info error. condition:%#v, rid:%s", cond.ToMapStr(), ctx.ReqID)
			return err
		}
		for _, setInfo := range setInfoArr {
			i.sets[setInfo.SetID] = setInfo
		}
	}
	if len(i.moduleIDs) > 0 {
		moduleInfoArr := []*metadata.ModuleInst{}
		cond := condition.CreateCondition().Field(common.BKModuleIDField).In(i.moduleIDs)
		err := i.hostManager.dbExecQuery(ctx, common.BKTableNameBaseModule, nil, cond.ToMapStr(), &moduleInfoArr)
		if err != nil {
			blog.Errorf("identifierHostTopoInfo query module info error. condition:%#v, rid:%s", cond.ToMapStr(), ctx.ReqID)
			return err
		}
		for _, moduleInfo := range moduleInfoArr {
			i.modules[moduleInfo.ModuleID] = moduleInfo
		}
	}
	if len(i.bizIDs) > 0 {
		bizInfoArr := []metadata.BizInst{}
		cond := condition.CreateCondition().Field(common.BKAppIDField).In(i.bizIDs)
		err := i.hostManager.dbExecQuery(ctx, common.BKTableNameBaseApp, nil, cond.ToMapStr(), &bizInfoArr)
		if err != nil {
			blog.Errorf("identifierHostTopoInfo query biz info error. rid:%s", ctx.ReqID)
			return err
		}
		for _, bizInfo := range bizInfoArr {
			i.bizs[bizInfo.BizID] = bizInfo
		}
	}
	blog.V(5).Infof("identifierHostTopoInfo query host topo info. bizIDs:%#v, setIDs:%#v, moduleIDs:%#v, biz:%#v, set:%#v, module:%#v, rid;%s", i.bizIDs, i.setIDs, i.moduleIDs, i.bizs, i.sets, i.modules, ctx.ReqID)

	return nil
}

// identifierHostCloud find host cloud area info
func (i *identifier) identifierHostCloud(ctx core.ContextParams) error {
	var cloudIDs []int64
	// parse  host  cloud id
	for _, host := range i.hosts {
		cloudIDs = append(cloudIDs, host.CloudID)
	}

	if len(cloudIDs) > 0 {
		cloudInfoArr := []metadata.CloudInst{}
		cond := condition.CreateCondition().Field(common.BKCloudIDField).In(cloudIDs)
		err := i.hostManager.dbExecQuery(ctx, common.BKTableNameBasePlat, nil, cond.ToMapStr(), &cloudInfoArr)
		if err != nil {
			blog.Errorf("identifierHostCloud query cloud info error. condition:%#v, rid:%s", cond.ToMapStr(), ctx.ReqID)
			return err
		}

		blog.V(5).Infof("identifierHostCloud query cloud info. cloud id:%#v, cloud info:%#v, rid;%s", cloudIDs, cloudInfoArr, ctx.ReqID)

		for _, cloudInfo := range cloudInfoArr {
			i.clouds[cloudInfo.CloudID] = cloudInfo
		}
	}

	return nil
}

// findHostServiceInst handle host service instance
func (i *identifier) findHostServiceInst(ctx core.ContextParams, hostIDs []int64) error {
	relationCond := condition.CreateCondition().Field(common.BKHostIDField).In(hostIDs)
	relations := []metadata.ProcessInstanceRelation{}

	// query process id with host id
	err := i.hostManager.dbExecQuery(ctx, common.BKTableNameProcessInstanceRelation, nil, relationCond.ToMapStr(), &relations)
	if err != nil {
		blog.ErrorJSON("findHostServiceInst query table %s err. cond:%s, rid:%s", common.BKTableNameProcessInstanceRelation, relationCond.ToMapStr(), ctx.ReqID)
		return err
	}

	blog.V(5).Infof("findHostServiceInst query host and process relation. hostID:%#v, relation:%#v, rid;%s", hostIDs, relations, ctx.ReqID)

	var procIDs []int64
	for _, relation := range relations {
		procIDs = append(procIDs, relation.ProcessID)
		i.hostProcRealtions[relation.HostID] = append(i.hostProcRealtions[relation.HostID], relation)
	}

	procInfos := make([]metadata.HostIdentProcess, 0)
	// query process info with process id
	cond := condition.CreateCondition().Field(common.BKProcIDField).In(procIDs)
	err = i.hostManager.dbExecQuery(ctx, common.BKTableNameBaseProcess, nil, cond.ToMapStr(), &procInfos)
	if err != nil {
		blog.ErrorJSON("findHostServiceInst query table %s err. cond:%s, rid:%s", common.BKTableNameBaseProcess, cond.ToMapStr(), ctx.ReqID)
		return err
	}

	blog.V(5).Infof("findHostServiceInst query process info. procIDs:%#v, info:%#v, rid;%s", procIDs, procInfos, ctx.ReqID)

	for _, procInfo := range procInfos {
		i.procs[procInfo.ProcessID] = procInfo
	}

	return nil
}

func (i *identifier) build(ctx core.ContextParams) {
	for idx, host := range i.hosts {
		if cloudInfo, ok := i.clouds[host.CloudID]; ok {
			host.CloudName = cloudInfo.CloudName
		}
		// 填充主机身份中的 业务，模块，集群信息
		for _, relation := range i.modulehosts[host.HostID] {
			mod := &metadata.HostIdentModule{
				SetID:    relation.SetID,
				ModuleID: relation.ModuleID,
				BizID:    relation.AppID,
			}

			host.HostIdentModule[strconv.FormatInt(mod.ModuleID, 10)] = mod

			if biz, ok := i.bizs[mod.BizID]; ok {
				mod.BizName = biz.BizName
				host.SupplierID = biz.SupplierID
			}

			if set, ok := i.sets[mod.SetID]; ok {
				mod.SetName = set.SetName
				mod.SetEnv = set.SetEnv
				mod.SetStatus = set.SetStatus
			}

			// 根据主机ID及服务实例。获取主机上的进程信息
			for _, hostProcRealtions := range i.hostProcRealtions[host.HostID] {
				if proc, ok := i.procs[hostProcRealtions.ProcessID]; ok {
					proc.BindModules = append(proc.BindModules, hostProcRealtions.ServiceInstanceID)
					host.Process = append(host.Process, proc)
				}

			}

		}
		i.hosts[idx] = host

	}
}

// dbExecQuery get info from table with condition
func (h *hostManager) dbExecQuery(ctx core.ContextParams, tableName string, fields []string, condMap mapstr.MapStr, result interface{}) error {
	newCondMap := util.SetQueryOwner(condMap, ctx.SupplierAccount)
	dbFind := h.DbProxy.Table(tableName).Find(newCondMap)
	if len(fields) > 0 {
		dbFind = dbFind.Fields(fields...)
	}
	err := dbFind.All(ctx, result)
	if err != nil {
		blog.ErrorJSON("findAll query table[%s] error. err:%s, conidtion:%s, rid:%s", tableName, err.Error(), newCondMap, ctx.ReqID)
		return ctx.Error.Error(common.CCErrCommDBSelectFailed)
	}
	return nil
}
