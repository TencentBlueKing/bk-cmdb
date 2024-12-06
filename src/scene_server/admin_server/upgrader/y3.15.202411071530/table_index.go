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
	"configcenter/src/scene_server/admin_server/upgrader/tools"
	upgradertypes "configcenter/src/scene_server/admin_server/upgrader/types"
	"configcenter/src/storage/dal"
	daltypes "configcenter/src/storage/dal/types"
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
	common.BKTableNameDelArchive:                 delArchiveIndexes,
	common.BKTableNameKubeDelArchive:             kubeDelArchiveIndexes,
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
	"cc_ObjectBaseMapping":                       objectBaseMappingIndexes,
}

var platTableIndexesArr = []string{common.BKTableNameSystem, common.BKTableNameIDgenerator}

var tableInstAsstIndexesMap = map[string][]daltypes.Index{
	common.BKInnerObjIDApp:     bizInstanceIndexes,
	common.BKInnerObjIDModule:  moduleInstanceIndexes,
	common.BKProcessObjectName: processInstanceIndexes,
	common.BKInnerObjIDHost:    hostInstanceIndexes,
	common.BKInnerObjIDProject: bkProjectInstanceIndexes,
	common.BKInnerObjIDBizSet:  bkBizSetObjInstIndexes,
	common.BKInnerObjIDPlat:    platInstanceIndexes,
	common.BKInnerObjIDSet:     setInstanceIndexes,
}

func initTableIndex(kit *rest.Kit, db dal.Dal, tableIndexMap map[string][]daltypes.Index) error {
	for table, index := range tableIndexMap {
		if err := tools.CreateTable(kit, db.Shard(kit.ShardOpts()), table); err != nil {
			blog.Errorf("create table %s failed, err: %v", table, err)
			return err
		}

		if err := tools.CreateIndexes(kit, db, table, index); err != nil {
			blog.Errorf("create table %s failed, err: %v", table, err)
			return err
		}
	}

	if err := createInstAsstTable(kit, db); err != nil {
		fmt.Errorf("create instance assosiation table failed, err: %v", err)
		return err
	}

	if kit.TenantID != upgradertypes.GetBlueKing() {
		return nil
	}
	for _, table := range platTableIndexesArr {
		if err := tools.CreateTable(kit, db.Shard(kit.SysShardOpts()), table); err != nil {
			fmt.Errorf("create plat table failed, err: %v", err)
			return err
		}
	}

	return nil
}

func buildInstAsstTableName(objID, tenantID string) string {
	return fmt.Sprintf("cc_InstAsst_%s_pub_%s", tenantID, objID)
}

func createInstAsstTable(kit *rest.Kit, db dal.Dal) error {
	for obj, index := range tableInstAsstIndexesMap {
		tableName := buildInstAsstTableName(obj, kit.TenantID)
		if err := tools.CreateTable(kit, db.Shard(kit.ShardOpts()), tableName); err != nil {
			blog.Errorf("create table %s failed, err: %v", tableName, err)
			return err
		}

		if err := tools.CreateIndexes(kit, db, tableName, index); err != nil {
			return err
		}
	}

	return nil
}
