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

package identifier

import (
	"strconv"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/source_controller/coreservice/core"
	hostutil "configcenter/src/source_controller/coreservice/core/host/util"
	"configcenter/src/storage/dal"
)

type Identifier struct {
	dbQuery *hostutil.DBExecQuery
	hosts   []metadata.HostIdentifier

	// tmp data
	setIDs           []int64
	moduleIDs        []int64
	bizIDs           []int64
	sets             map[int64]metadata.SetInst
	modules          map[int64]*metadata.ModuleInst
	bizs             map[int64]metadata.BizInst
	clouds           map[int64]metadata.CloudInst
	hostProcRelation map[int64][]metadata.HostIdentProcess
	modulehosts      map[int64][]metadata.ModuleHost
	asstMap          map[string]string
	layers           map[string]map[int64]metadata.MainlineInstInfo
}

func NewIdentifier(db dal.DB) *Identifier {
	dbQuery := hostutil.NewDBExecQuery(db)
	return &Identifier{
		dbQuery:          dbQuery,
		sets:             make(map[int64]metadata.SetInst),
		modules:          make(map[int64]*metadata.ModuleInst),
		bizs:             make(map[int64]metadata.BizInst),
		clouds:           make(map[int64]metadata.CloudInst),
		hostProcRelation: make(map[int64][]metadata.HostIdentProcess),
		modulehosts:      make(map[int64][]metadata.ModuleHost),
		asstMap:          make(map[string]string),
		layers:           make(map[string]map[int64]metadata.MainlineInstInfo),
	}
}

func (i *Identifier) Identifier(ctx core.ContextParams, hostIDs []int64) ([]metadata.HostIdentifier, error) {
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

	err = i.findHostTopoInfo(ctx)
	if err != nil {
		return nil, err
	}
	err = i.findHostCloud(ctx)
	if err != nil {
		return nil, err
	}

	err = i.findHostLayerInfo(ctx)
	if err != nil {
		return nil, err
	}

	i.build(ctx)
	return i.hosts, nil
}

// FindHost query host info
func (i *Identifier) findHost(ctx core.ContextParams, hostIDs []int64) error {
	hostCond := condition.CreateCondition().Field(common.BKHostIDField).In(hostIDs)
	condHostMap := util.SetQueryOwner(hostCond.ToMapStr(), ctx.SupplierAccount)
	// fetch all hosts
	i.hosts = make([]metadata.HostIdentifier, 0) // New func init
	err := i.dbQuery.DbProxy.Table(common.BKTableNameBaseHost).Find(condHostMap).All(ctx, &i.hosts)
	if err != nil {
		blog.ErrorJSON("findHost query host error. err:%s, conidtion:%s, rid:%s", err.Error(), condHostMap, ctx.ReqID)
		return ctx.Error.Error(common.CCErrCommDBSelectFailed)
	}

	blog.V(5).Infof("findHost query host info. host:%#v, rid;%s", i.hosts, ctx.ReqID)

	return nil
}

// findModuleHostRelation query host and module relation
func (i *Identifier) findModuleHostRelation(ctx core.ContextParams, hostIDs []int64) error {
	hostModuleCond := condition.CreateCondition().Field(common.BKHostIDField).In(hostIDs)
	condModuleHostMap := util.SetQueryOwner(hostModuleCond.ToMapStr(), ctx.SupplierAccount)
	// fetch  host and module relation
	moduleHostRelation := make([]metadata.ModuleHost, 0)
	err := i.dbQuery.DbProxy.Table(common.BKTableNameModuleHostConfig).Find(condModuleHostMap).All(ctx, &moduleHostRelation)
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

// findHostTopoInfo handle host biz,set, module info
func (i *Identifier) findHostTopoInfo(ctx core.ContextParams) error {

	// fetch set info
	if len(i.setIDs) > 0 {
		setInfoArr := make([]metadata.SetInst, 0)
		cond := condition.CreateCondition().Field(common.BKSetIDField).In(i.setIDs)
		err := i.dbQuery.ExecQuery(ctx, common.BKTableNameBaseSet, nil, cond.ToMapStr(), &setInfoArr)
		if err != nil {
			blog.Errorf("findHostTopoInfo query set info error. condition:%#v, rid:%s", cond.ToMapStr(), ctx.ReqID)
			return err
		}
		for _, setInfo := range setInfoArr {
			i.sets[setInfo.SetID] = setInfo
		}
	}
	if len(i.moduleIDs) > 0 {
		moduleInfoArr := make([]*metadata.ModuleInst, 0)
		cond := condition.CreateCondition().Field(common.BKModuleIDField).In(i.moduleIDs)
		err := i.dbQuery.ExecQuery(ctx, common.BKTableNameBaseModule, nil, cond.ToMapStr(), &moduleInfoArr)
		if err != nil {
			blog.Errorf("findHostTopoInfo query module info error. condition:%#v, rid:%s", cond.ToMapStr(), ctx.ReqID)
			return err
		}
		for _, moduleInfo := range moduleInfoArr {
			i.modules[moduleInfo.ModuleID] = moduleInfo
		}
	}
	if len(i.bizIDs) > 0 {
		bizInfoArr := make([]metadata.BizInst, 0)
		cond := condition.CreateCondition().Field(common.BKAppIDField).In(i.bizIDs)
		err := i.dbQuery.ExecQuery(ctx, common.BKTableNameBaseApp, nil, cond.ToMapStr(), &bizInfoArr)
		if err != nil {
			blog.Errorf("findHostTopoInfo query biz info error. rid:%s", ctx.ReqID)
			return err
		}
		for _, bizInfo := range bizInfoArr {
			i.bizs[bizInfo.BizID] = bizInfo
		}
	}
	blog.V(5).Infof("findHostTopoInfo query host topo info. bizIDs:%#v, setIDs:%#v, moduleIDs:%#v, biz:%#v, set:%#v, module:%#v, rid;%s", i.bizIDs, i.setIDs, i.moduleIDs, i.bizs, i.sets, i.modules, ctx.ReqID)

	return nil
}

// findHostCloud find host cloud area info
func (i *Identifier) findHostCloud(ctx core.ContextParams) error {
	var cloudIDs []int64
	// parse  host  cloud id
	for _, host := range i.hosts {
		cloudIDs = append(cloudIDs, host.CloudID)
	}

	if len(cloudIDs) > 0 {
		cloudInfoArr := make([]metadata.CloudInst, 0)
		cond := condition.CreateCondition().Field(common.BKCloudIDField).In(cloudIDs)
		err := i.dbQuery.ExecQuery(ctx, common.BKTableNameBasePlat, nil, cond.ToMapStr(), &cloudInfoArr)
		if err != nil {
			blog.Errorf("findHostCloud query cloud info error. condition:%#v, rid:%s", cond.ToMapStr(), ctx.ReqID)
			return err
		}

		blog.V(5).Infof("findHostCloud query cloud info. cloud id:%#v, cloud info:%#v, rid;%s", cloudIDs, cloudInfoArr, ctx.ReqID)

		for _, cloudInfo := range cloudInfoArr {
			i.clouds[cloudInfo.CloudID] = cloudInfo
		}
	}

	return nil
}

// findHostServiceInst handle host service instance
func (i *Identifier) findHostServiceInst(ctx core.ContextParams, hostIDs []int64) error {
	relations := make([]metadata.ProcessInstanceRelation, 0)

	// query process id with host id
	relationFilter := map[string]interface{}{
		common.BKHostIDField: map[string]interface{}{
			common.BKDBIN: hostIDs,
		},
	}
	err := i.dbQuery.ExecQuery(ctx, common.BKTableNameProcessInstanceRelation, nil, relationFilter, &relations)
	if err != nil {
		blog.ErrorJSON("findHostServiceInst query table %s err. cond:%s, rid:%s", common.BKTableNameProcessInstanceRelation, relationFilter, ctx.ReqID)
		return err
	}

	blog.V(5).Infof("findHostServiceInst query host and process relation. hostID:%#v, relation:%#v, rid;%s", hostIDs, relations, ctx.ReqID)

	var procIDs []int64
	var serviceInstIDs []int64
	// 进程与服务实例的关系
	procServiceInstMap := make(map[int64][]int64, 0)
	for _, relation := range relations {
		procIDs = append(procIDs, relation.ProcessID)
		serviceInstIDs = append(serviceInstIDs, relation.ServiceInstanceID)
		procServiceInstMap[relation.ProcessID] = append(procServiceInstMap[relation.ProcessID], relation.ServiceInstanceID)
	}

	serviceInstInfos := make([]metadata.ServiceInstance, 0)
	serviceInstFilter := map[string]interface{}{
		common.BKFieldID: map[string]interface{}{
			common.BKDBIN: serviceInstIDs,
		},
	}
	err = i.dbQuery.ExecQuery(ctx, common.BKTableNameServiceInstance, nil, serviceInstFilter, &serviceInstInfos)
	if err != nil {
		blog.ErrorJSON("findHostServiceInst query table %s err. cond:%s, rid:%s", common.BKTableNameServiceInstance, serviceInstFilter, ctx.ReqID)
		return err
	}
	blog.V(5).Infof("findHostServiceInst query service instance info. service instance id:%#v, info:%#v, rid;%s", serviceInstIDs, serviceInstInfos, ctx.ReqID)
	// 服务实例与模块的关系
	serviceInstModuleRelation := make(map[int64][]int64, 0)
	for _, serviceInstInfo := range serviceInstInfos {
		serviceInstModuleRelation[serviceInstInfo.ID] = append(serviceInstModuleRelation[serviceInstInfo.ID], serviceInstInfo.ModuleID)
	}

	procModuleRelation := make(map[int64][]int64, 0)
	for procID, serviceInstIDs := range procServiceInstMap {
		for _, serviceInstID := range serviceInstIDs {
			for _, moduleID := range serviceInstModuleRelation[serviceInstID] {
				procModuleRelation[procID] = append(procModuleRelation[procID], moduleID)
			}
		}
	}

	procInfos := make([]metadata.HostIdentProcess, 0)
	// query process info with process id
	cond := condition.CreateCondition().Field(common.BKProcIDField).In(procIDs)
	err = i.dbQuery.ExecQuery(ctx, common.BKTableNameBaseProcess, nil, cond.ToMapStr(), &procInfos)
	if err != nil {
		blog.ErrorJSON("findHostServiceInst query table %s err. cond:%s, rid:%s", common.BKTableNameBaseProcess, cond.ToMapStr(), ctx.ReqID)
		return err
	}

	blog.V(5).Infof("findHostServiceInst query process info. procIDs:%#v, info:%#v, rid;%s", procIDs, procInfos, ctx.ReqID)

	procs := make(map[int64]metadata.HostIdentProcess, 0)
	for _, procInfo := range procInfos {
		if moduleIDs, ok := procModuleRelation[procInfo.ProcessID]; ok {
			procInfo.BindModules = moduleIDs
		}
		procs[procInfo.ProcessID] = procInfo
	}

	// 主机和进程之间的关系,生成主机与进程的关系
	for _, relation := range relations {
		if procInfo, ok := procs[relation.ProcessID]; ok {
			i.hostProcRelation[relation.HostID] = append(i.hostProcRelation[relation.HostID], procInfo)
		}
	}

	return nil
}

// findHostLayerInfo handle host layer info
func (i *Identifier) findHostLayerInfo(ctx core.ContextParams) error {
	// find mainline association
	asstArr := make([]metadata.Association, 0)
	cond := condition.CreateCondition().Field(common.AssociationKindIDField).Eq(common.AssociationKindMainline)
	err := i.dbQuery.ExecQuery(ctx, common.BKTableNameObjAsst, nil, cond.ToMapStr(), &asstArr)
	if err != nil {
		blog.ErrorJSON("findHostLayerInfo query mainline association info error. condition:%s, rid:%s", cond.ToMapStr(), ctx.ReqID)
		return err
	}

	for _, asst := range asstArr {
		i.asstMap[asst.ObjectID] = asst.AsstObjID
	}

	// initialize parent inst search param
	parentIDs := make([]int64, 0)
	for _, set := range i.sets {
		parentIDs = append(parentIDs, set.ParentID)
	}
	curObj, ok := i.asstMap[common.BKInnerObjIDSet]
	if !ok {
		return nil
	}

	// find layer info
	for curObj != "" && curObj != common.BKInnerObjIDApp {
		layers := make([]metadata.MainlineInstInfo, 0)
		cond := condition.CreateCondition().Field(common.BKInstIDField).In(parentIDs)
		cond.Field(common.BKObjIDField).Eq(curObj)
		err := i.dbQuery.ExecQuery(ctx, common.BKTableNameBaseInst, nil, cond.ToMapStr(), &layers)
		if err != nil {
			blog.Errorf("findHostLayerInfo query layer info error. condition:%#v, rid:%s", cond.ToMapStr(), ctx.ReqID)
			return err
		}

		parentIDs = make([]int64, 0)
		curObj = i.asstMap[curObj]

		for _, layer := range layers {
			if i.layers[layer.ObjID] == nil {
				i.layers[layer.ObjID] = make(map[int64]metadata.MainlineInstInfo)
			}
			i.layers[layer.ObjID][layer.InstID] = layer
			parentIDs = append(parentIDs, layer.ParentID)
		}
	}

	blog.V(5).Infof("findHostLayerInfo query host layer info. layers: %#v, rid;%s", i.layers, ctx.ReqID)
	return nil
}

func (i *Identifier) build(ctx core.ContextParams) {
	for idx, host := range i.hosts {
		if cloudInfo, ok := i.clouds[host.CloudID]; ok {
			host.CloudName = cloudInfo.CloudName
		}
		// 填充主机身份中的 业务，模块，集群，自定义层级信息
		for _, relation := range i.modulehosts[host.HostID] {
			mod := &metadata.HostIdentModule{
				SetID:    relation.SetID,
				ModuleID: relation.ModuleID,
				BizID:    relation.AppID,
			}

			if host.HostIdentModule == nil {
				host.HostIdentModule = make(map[string]*metadata.HostIdentModule, 0)
			}
			host.HostIdentModule[strconv.FormatInt(mod.ModuleID, 10)] = mod

			if biz, ok := i.bizs[mod.BizID]; ok {
				mod.BizName = biz.BizName
				host.SupplierID = biz.SupplierID
			}

			var parentID int64
			if set, ok := i.sets[mod.SetID]; ok {
				mod.SetName = set.SetName
				mod.SetEnv = set.SetEnv
				mod.SetStatus = set.SetStatus
				parentID = set.ParentID
			}

			if module, ok := i.modules[mod.ModuleID]; ok {
				mod.ModuleName = module.ModuleName
			}

			curObj, ok := i.asstMap[common.BKInnerObjIDSet]
			if !ok {
				continue
			}
			var layer *metadata.Layer
			for curObj != "" && curObj != common.BKInnerObjIDApp {
				objLayers, ok := i.layers[curObj]
				if !ok {
					curObj = i.asstMap[curObj]
					continue
				}
				objLayer, ok := objLayers[parentID]
				if !ok {
					curObj = i.asstMap[curObj]
					continue
				}

				layer = &metadata.Layer{
					InstID:   objLayer.InstID,
					InstName: objLayer.InstName,
					ObjID:    objLayer.ObjID,
					Child:    layer,
				}

				curObj = i.asstMap[curObj]
				parentID = objLayer.ParentID
			}
			mod.Layer = layer
		}
		host.Process = i.hostProcRelation[host.HostID]
		i.hosts[idx] = host
	}
}
