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

package common

import (
	"errors"
	"fmt"
	"strings"
)

// table names
const (
	// BKTableNamePropertyGroup the table name of the property group
	BKTableNamePropertyGroup = "PropertyGroup"

	// BKTableNameAsstDes the table name of the asst des
	BKTableNameAsstDes = "AsstDes"

	// BKTableNameObjDes the table name of the object
	BKTableNameObjDes = "ObjDes"

	// BKTableNameObjUnique the table name of the object
	BKTableNameObjUnique = "ObjectUnique"

	// BKTableNameObjAttDes the table name of the object attribute
	BKTableNameObjAttDes = "ObjAttDes"

	// BKTableNameObjClassification the table name of the object classification
	BKTableNameObjClassification = "ObjClassification"

	// BKTableNameInstAsst the table name of the inst association
	BKTableNameInstAsst = "InstAsst"

	BKTableNameBaseApp    = "ApplicationBase"
	BKTableNameBaseBizSet = "BizSetBase"

	// BKTableNameModelQuoteRelation model reference relationship table name.
	BKTableNameModelQuoteRelation = "ModelQuoteRelation"

	BKTableNameBaseProject = "ProjectBase"
	BKTableNameBaseHost    = "HostBase"
	BKTableNameBaseModule  = "ModuleBase"
	BKTableNameBaseInst    = "ObjectBase"
	BKTableNameBasePlat    = "PlatBase"
	BKTableNameBaseSet     = "SetBase"
	BKTableNameBaseProcess = "Process"

	BKTableNameDelArchive     = "DelArchive"
	BKTableNameKubeDelArchive = "KubeDelArchive"

	BKTableNameModuleHostConfig = "ModuleHostConfig"
	BKTableNameSystem           = "System"
	BKTableNameHistory          = "History"
	BKTableNameHostFavorite     = "HostFavourite"
	BKTableNameAuditLog         = "AuditLog"
	BKTableNameUserAPI          = "UserAPI"
	BKTableNameDynamicGroup     = "DynamicGroup"
	BKTableNameUserCustom       = "UserCustom"
	BKTableNameObjAsst          = "ObjAsst"
	BKTableNameTopoGraphics     = "TopoGraphics"
	BKTableNameTransaction      = "Transaction"
	BKTableNameIDgenerator      = "idgenerator"
	BKTableNameGlobalConfig     = "GlobalConfig"

	BKTableNameHostLock = "HostLock"

	// process tables
	BKTableNameServiceCategory         = "ServiceCategory"
	BKTableNameServiceTemplate         = "ServiceTemplate"
	BKTableNameServiceTemplateAttr     = "ServiceTemplateAttr"
	BKTableNameServiceInstance         = "ServiceInstance"
	BKTableNameProcessTemplate         = "ProcessTemplate"
	BKTableNameProcessInstanceRelation = "ProcessInstanceRelation"

	BKTableNameSetTemplate                = "SetTemplate"
	BKTableNameSetTemplateAttr            = "SetTemplateAttr"
	BKTableNameSetServiceTemplateRelation = "SetServiceTemplateRelation"
	BKTableNameAPITask                    = "APITask"
	BKTableNameAPITaskSyncHistory         = "APITaskSyncHistory"

	// BKTableNameHostApplyRule rule for host property auto apply
	BKTableNameHostApplyRule = "HostApplyRule"

	// BKTableNameWatchToken the table to store the latest watch token for database
	BKTableNameWatchToken = "WatchToken"

	// BKTableNameLastWatchEvent is the table to store the latest watch event info for resources
	BKTableNameLastWatchEvent = "LastWatchEvent"

	// BKTableNameMainlineInstance is a virtual collection name which represent for mainline instance events
	BKTableNameMainlineInstance = "MainlineInstance"

	// BKTableNameFieldTemplate  field template table
	BKTableNameFieldTemplate = "FieldTemplate"

	// BKTableNameObjAttDesTemplate  field template  attribute description table
	BKTableNameObjAttDesTemplate = "ObjAttDesTemplate"

	// BKTableNameObjectUniqueTemplate  field template unique checklist table
	BKTableNameObjectUniqueTemplate = "ObjectUniqueTemplate"

	// BKTableNameObjFieldTemplateRelation  object and field template relationship table
	BKTableNameObjFieldTemplateRelation = "ObjFieldTemplateRelation"

	// BKTableNameTenant is the tenant table
	BKTableNameTenant = "Tenant"

	// BKTableNameTenantTemplate is the tenant template(public data that needs to be initialized for all tenants) table
	BKTableNameTenantTemplate = "TenantTemplate"

	// BKTableNameObjectBaseMapping object base mapping table
	BKTableNameObjectBaseMapping = "ObjectBaseMapping"

	// BKTableNameWatchDBRelation is the db and watch db relation table
	BKTableNameWatchDBRelation = "WatchDBRelation"

	// BKTableNameFullSyncCond is the full synchronization cache condition table
	BKTableNameFullSyncCond = "FullSyncCond"

	// BKTableNameCacheWatchToken is the cache event watch token table
	BKTableNameCacheWatchToken = "CacheWatchToken"

	// BKTableNameDefaultAreaHost is used to store the default area host and ensure that IP is not repeated.
	BKTableNameDefaultAreaHost = "DefaultAreaHost"
)

// AllTables is all table names, not include the sharding tables which is created dynamically,
// such as object instance sharding table 'ObjectBase_{supplierAccount}_pub_{objectID}'.
var AllTables = []string{
	BKTableNamePropertyGroup,
	BKTableNameObjDes,
	BKTableNameObjAttDes,
	BKTableNameObjClassification,
	BKTableNameInstAsst,
	BKTableNameBaseApp,
	BKTableNameBaseHost,
	BKTableNameBaseModule,
	BKTableNameBaseInst,
	BKTableNameBasePlat,
	BKTableNameBaseSet,
	BKTableNameBaseProcess,
	BKTableNameModuleHostConfig,
	BKTableNameSystem,
	BKTableNameGlobalConfig,
	BKTableNameHistory,
	BKTableNameHostFavorite,
	BKTableNameAuditLog,
	BKTableNameUserAPI,
	BKTableNameDynamicGroup,
	BKTableNameUserCustom,
	BKTableNameObjAsst,
	BKTableNameTopoGraphics,
	BKTableNameTransaction,
	BKTableNameIDgenerator,
	BKTableNameHostLock,
	BKTableNameObjUnique,
	BKTableNameAsstDes,
	BKTableNameServiceCategory,
	BKTableNameServiceTemplate,
	BKTableNameServiceInstance,
	BKTableNameProcessTemplate,
	BKTableNameProcessInstanceRelation,
	BKTableNameSetTemplate,
	BKTableNameSetServiceTemplateRelation,
	BKTableNameHostApplyRule,
	BKTableNameAPITask,
	BKTableNameAPITaskSyncHistory,
	BKTableNameDefaultAreaHost,
}

// TableSpecifier is table specifier type which describes the metadata
// access or classification level.
type TableSpecifier string

const (
	// TableSpecifierPublic is public specifier for table.
	TableSpecifierPublic TableSpecifier = "pub"
)

const (
	// BKObjectInstShardingTablePrefix is prefix of object instance sharding table.
	BKObjectInstShardingTablePrefix = BKTableNameBaseInst + "_"

	// BKObjectInstAsstShardingTablePrefix is prefix of object instance association sharding table.
	BKObjectInstAsstShardingTablePrefix = BKTableNameInstAsst + "_"
)

// GetObjInstTableName return the object instance table name.
// example: ObjectBase_{uuid}, such as 'ObjectBase_d0hcqvfk9a52rli3fuq0'
func GetObjInstTableName(uuid string) string {
	return fmt.Sprintf("%s%s", BKObjectInstShardingTablePrefix, uuid)
}

// GetObjInstAsstTableName return the object instance association table name.
// example: InstAsst_{uuid}, such as 'InstAsst_d0hcqvfk9a52rli3fuq0'
func GetObjInstAsstTableName(uuid string) string {
	return fmt.Sprintf("%s%s", BKObjectInstAsstShardingTablePrefix, uuid)
}

// GetObjectUUIDByTableName get uuid by object sharding table name.
func GetObjectUUIDByTableName(collectionName string) (string, error) {
	if !IsObjectInstShardingTable(collectionName) {
		return "", fmt.Errorf("%s is not an object instance sharding table", collectionName)
	}
	return strings.TrimPrefix(collectionName, BKObjectInstShardingTablePrefix), nil
}

// IsObjectShardingTable returns if the target table is an object sharding table, include
// object instance and association.
func IsObjectShardingTable(tableName string) bool {
	if IsObjectInstShardingTable(tableName) {
		return true
	}
	return IsObjectInstAsstShardingTable(tableName)
}

// IsObjectInstShardingTable returns if the target table is an object instance sharding table.
func IsObjectInstShardingTable(tableName string) bool {
	// check object instance table, ObjectBase_{uuid}
	return strings.HasPrefix(tableName, BKObjectInstShardingTablePrefix)
}

// IsObjectInstAsstShardingTable returns if the target table is an object instance association sharding table.
func IsObjectInstAsstShardingTable(tableName string) bool {
	// check object instance association table, InstAsst_{uuid}
	return strings.HasPrefix(tableName, BKObjectInstAsstShardingTablePrefix)
}

// GetInstTableName returns inst data table name
func GetInstTableName(objID string, uuid string) string {
	if IsInnerModel(objID) {
		return GetInnerInstTableName(objID)
	}

	return GetObjInstTableName(uuid)
}

// GetInnerInstTableName returns inner object instance table name
func GetInnerInstTableName(objID string) string {
	switch objID {
	case BKInnerObjIDApp:
		return BKTableNameBaseApp
	case BKInnerObjIDBizSet:
		return BKTableNameBaseBizSet
	case BKInnerObjIDProject:
		return BKTableNameBaseProject
	case BKInnerObjIDSet:
		return BKTableNameBaseSet
	case BKInnerObjIDModule:
		return BKTableNameBaseModule
	case BKInnerObjIDHost:
		return BKTableNameBaseHost
	case BKInnerObjIDProc:
		return BKTableNameBaseProcess
	case BKInnerObjIDPlat:
		return BKTableNameBasePlat
	}
	return ""
}

var platformTableMap = map[string]struct{}{
	BKTableNameSystem:             {},
	BKTableNameIDgenerator:        {},
	BKTableNameTenant:             {},
	BKTableNameTenantTemplate:     {},
	BKTableNameWatchToken:         {},
	BKTableNameAPITask:            {},
	BKTableNameAPITaskSyncHistory: {},
	BKTableNameWatchDBRelation:    {},
	BKTableNameFullSyncCond:       {},
	BKTableNameCacheWatchToken:    {},
	"SrcSyncDataToken":            {},
	"SrcSyncDataCursor":           {},
	BKTableNameGlobalConfig:       {},
	BKTableNameDefaultAreaHost:    {},
}

// IsPlatformTable returns if the target table is a platform table
func IsPlatformTable(tableName string) bool {
	_, exists := platformTableMap[tableName]
	return exists
}

var platformTableWithTenantMap = map[string]struct{}{
	BKTableNameAPITask:            {},
	BKTableNameAPITaskSyncHistory: {},
	BKTableNameFullSyncCond:       {},
	BKTableNameCacheWatchToken:    {},
	"SrcSyncDataToken":            {},
	"SrcSyncDataCursor":           {},
	BKTableNameGlobalConfig:       {},
}

// IsPlatformTableWithTenant returns if the target table is a platform table with tenant id field
func IsPlatformTableWithTenant(tableName string) bool {
	_, exists := platformTableWithTenantMap[tableName]
	return exists
}

// PlatformTables returns platform tables
func PlatformTables() []string {
	tables := make([]string, 0)
	for tableName := range platformTableMap {
		tables = append(tables, tableName)
	}
	return tables
}

// GenTenantTableName generate tenant table name by table name and tenant id
func GenTenantTableName(tenantID, tableName string) string {
	return fmt.Sprintf("%s_%s", tenantID, tableName)
}

// SplitTenantTableName split tenant table name to tenant id and table name
func SplitTenantTableName(tenantTableName string) (string, string, error) {
	if IsPlatformTable(tenantTableName) {
		return "", tenantTableName, nil
	}

	if strings.Contains(tenantTableName, "_"+BKObjectInstShardingTablePrefix) {
		sepIdx := strings.LastIndex(tenantTableName, "_"+BKObjectInstShardingTablePrefix)
		return tenantTableName[:sepIdx], tenantTableName[sepIdx+1:], nil
	}

	if strings.Contains(tenantTableName, "_"+BKObjectInstAsstShardingTablePrefix) {
		sepIdx := strings.LastIndex(tenantTableName, "_"+BKObjectInstAsstShardingTablePrefix)
		return tenantTableName[:sepIdx], tenantTableName[sepIdx+1:], nil
	}

	sepIdx := strings.LastIndex(tenantTableName, "_")
	if sepIdx == -1 {
		return "", "", errors.New("tenant table name is invalid")
	}
	return tenantTableName[:sepIdx], tenantTableName[sepIdx+1:], nil
}
