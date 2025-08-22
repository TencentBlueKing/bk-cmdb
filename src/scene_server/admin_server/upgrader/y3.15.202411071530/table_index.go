/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 THL A29 Limited,
 * a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 * We undertake not to change the open source license (MIT license) applicable
 * to the current version of the project delivered to anyone in the future.
 */

package y3_15_202411071530

import (
	"fmt"

	fullsynccond "configcenter/pkg/cache/full-sync-cond"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/kube/types"
	"configcenter/src/scene_server/admin_server/logics"
	"configcenter/src/scene_server/admin_server/upgrader/y3.15.202411071530/data"
	"configcenter/src/storage/dal/mongo/local"
	daltypes "configcenter/src/storage/dal/types"
	"configcenter/src/storage/driver/mongodb"
)

var tableIndexMap = map[string][]daltypes.Index{
	common.BKTableNameAPITask:                    apiTaskIndexes,
	common.BKTableNameAPITaskSyncHistory:         aPITaskSyncHistoryIndexes,
	common.BKTableNameBaseApp:                    applicationBaseIndexes,
	common.BKTableNameAsstDes:                    asstDesIndexes,
	common.BKTableNameAuditLog:                   auditLogIndexes,
	common.BKTableNameBaseBizSet:                 bizSetBaseIndexes,
	common.BKTableNamePropertyGroup:              propertyGroupIndexes,
	common.BKTableNameObjDes:                     objDesIndexes,
	common.BKTableNameObjUnique:                  objectUniqueIndexes,
	common.BKTableNameObjAttDes:                  objAttDesIndexes,
	common.BKTableNameObjClassification:          objClassificationIndexes,
	common.BKTableNameInstAsst:                   instAsstIndexes,
	common.BKTableNameModelQuoteRelation:         modelQuoteRelationIndexes,
	common.BKTableNameBaseProject:                projectBaseIndexes,
	common.BKTableNameBaseHost:                   hostBaseIndexes,
	common.BKTableNameBaseModule:                 moduleBaseIndexes,
	common.BKTableNameBaseInst:                   objectBaseIndexes,
	common.BKTableNameBasePlat:                   platBaseIndexes,
	common.BKTableNameBaseSet:                    setBaseIndexes,
	common.BKTableNameBaseProcess:                processIndexes,
	common.BKTableNameModuleHostConfig:           moduleHostConfigIndexes,
	common.BKTableNameHostFavorite:               nil,
	common.BKTableNameUserAPI:                    nil,
	common.BKTableNameDynamicGroup:               dynamicGroupIndexes,
	common.BKTableNameUserCustom:                 nil,
	common.BKTableNameObjAsst:                    objAsstIndexes,
	common.BKTableNameTopoGraphics:               topoGraphicsIndexes,
	common.BKTableNameHostLock:                   hostLockIndexes,
	common.BKTableNameServiceCategory:            serviceCategoryIndexes,
	common.BKTableNameServiceTemplate:            serviceTemplateIndexes,
	common.BKTableNameServiceTemplateAttr:        serviceTemplateAttrIndexes,
	common.BKTableNameServiceInstance:            serviceInstanceIndexes,
	common.BKTableNameProcessTemplate:            processTemplateIndexes,
	common.BKTableNameProcessInstanceRelation:    processInstanceRelationIndexes,
	common.BKTableNameSetTemplate:                setTemplateIndexes,
	common.BKTableNameSetTemplateAttr:            setTemplateAttrIndexes,
	common.BKTableNameSetServiceTemplateRelation: setServiceTemplateRelationIndexes,
	common.BKTableNameHostApplyRule:              hostApplyRuleIndexes,
	common.BKTableNameFieldTemplate:              fieldTemplateIndexes,
	common.BKTableNameObjAttDesTemplate:          objAttDesTemplateIndexes,
	common.BKTableNameObjectUniqueTemplate:       objectUniqueTemplateIndexes,
	common.BKTableNameObjFieldTemplateRelation:   objFieldTemplateRelationIndexes,
	types.BKTableNameBaseCluster:                 clusterBaseIndexes,
	types.BKTableNameBaseContainer:               containerBaseIndexes,
	types.BKTableNameBaseCronJob:                 cronJobBaseIndexes,
	types.BKTableNameBaseDaemonSet:               daemonSetBaseIndexes,
	types.BKTableNameBaseDeployment:              deploymentBaseIndexes,
	fullsynccond.BKTableNameFullSyncCond:         fullSyncCondIndexes,
	types.BKTableNameGameDeployment:              gameDeploymentBaseIndexes,
	types.BKTableNameGameStatefulSet:             gameStatefulSetBaseIndexes,
	types.BKTableNameBaseJob:                     jobBaseIndexes,
	types.BKTableNameBaseStatefulSet:             statefulSetBaseIndexes,
	types.BKTableNameNsSharedClusterRel:          nsSharedClusterRelationIndexes,
	types.BKTableNameBasePodWorkload:             podWorkloadBaseIndexes,
	types.BKTableNameBaseNode:                    nodeBaseIndexes,
	types.BKTableNameBasePod:                     podBaseIndexes,
	types.BKTableNameBaseNamespace:               namespaceBaseIndexes,
	common.BKTableNameObjectBaseMapping:          objectBaseMappingIndexes,
}

var platTableIndexesArr = map[string][]daltypes.Index{
	common.BKTableNameSystem:          nil,
	common.BKTableNameIDgenerator:     nil,
	common.BKTableNameTenantTemplate:  templateIndexes,
	common.BKTableNameTenant:          tenantIndexes,
	common.BKTableNameGlobalConfig:    globalConfigIndexes,
	common.BKTableNameDefaultAreaHost: defaultAreaHostIndexes,
}

var tableInstAsstArr = []string{
	common.BKInnerObjIDApp,
	common.BKInnerObjIDModule,
	common.BKProcessObjectName,
	common.BKInnerObjIDHost,
	common.BKInnerObjIDProject,
	common.BKInnerObjIDBizSet,
	common.BKInnerObjIDPlat,
	common.BKInnerObjIDSet,
}

func initTableIndex(kit *rest.Kit, db local.DB, tableIndexMap map[string][]daltypes.Index) error {
	for table, index := range tableIndexMap {
		if err := logics.CreateTable(kit, db, table); err != nil {
			blog.Errorf("create table %s failed, err: %v", table, err)
			return err
		}
		if err := logics.CreateIndexes(kit, db, table, index); err != nil {
			blog.Errorf("create table %s failed, err: %v", table, err)
			return err
		}
	}

	if err := createInstAsstTable(kit, db); err != nil {
		fmt.Errorf("create instance assosiation table failed, err: %v", err)
		return err
	}

	// create plat table and indexes
	for table, index := range platTableIndexesArr {
		if err := logics.CreateTable(kit, mongodb.Dal().Shard(kit.SysShardOpts()), table); err != nil {
			fmt.Errorf("create platfrom table failed, err: %v", err)
			return err
		}

		if len(index) == 0 {
			continue
		}
		if err := logics.CreateIndexes(kit, mongodb.Dal().Shard(kit.SysShardOpts()), table, index); err != nil {
			blog.Errorf("create table %s failed, err: %v", table, err)
			return err
		}
	}

	return nil
}

// createInstAsstTable add object data and create instance association table indexes
func createInstAsstTable(kit *rest.Kit, db local.DB) error {
	objUUIDMap, err := data.AddObjectData(kit, db)
	if err != nil {
		blog.Errorf("add object data failed, err: %v", err)
		return err
	}

	for _, obj := range tableInstAsstArr {
		uuid, exist := objUUIDMap[obj]
		if !exist {
			blog.Errorf("object uuid not exist, obj: %s", obj)
			return fmt.Errorf("object uuid not exist")
		}
		tableName := fmt.Sprintf("InstAsst_%s", uuid)
		if err = logics.CreateTable(kit, db, tableName); err != nil {
			blog.Errorf("create table %s failed, err: %v", tableName, err)
			return err
		}

		if err = logics.CreateIndexes(kit, db, tableName, instAsstCommonIndexes); err != nil {
			return err
		}
	}

	return nil
}
