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
	BKTableNamePlatformAuditLog = "PlatformAuditLog"
	BKTableNameUserAPI          = "UserAPI"
	BKTableNameDynamicGroup     = "DynamicGroup"
	BKTableNameUserCustom       = "UserCustom"
	BKTableNameObjAsst          = "ObjAsst"
	BKTableNameTopoGraphics     = "TopoGraphics"
	BKTableNameTransaction      = "Transaction"
	BKTableNameIDgenerator      = "idgenerator"

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

// GetObjectInstTableName return the object instance table name in sharding mode base on
// the object ID. Format: ObjectBase_{supplierAccount}_{Specifier}_{ObjectID}, such as 'ObjectBase_0_pub_switch'.
func GetObjectInstTableName(objID, tenantID string) string {
	return fmt.Sprintf("%s%s_%s_%s", BKObjectInstShardingTablePrefix, tenantID, TableSpecifierPublic, objID)
}

// GetObjectInstObjIDByTableName return the object id
// example: ObjectBase_{supplierAccount}_{Specifier}_{ObjectID}, such as 'ObjectBase_0_pub_switch',return switch.
func GetObjectInstObjIDByTableName(collectionName, tenantID string) (string, error) {
	prefix := fmt.Sprintf("%s%s_", BKObjectInstShardingTablePrefix, tenantID)
	suffix := strings.TrimPrefix(collectionName, prefix)
	suffixSlice := strings.Split(suffix, "_")
	if len(suffixSlice) <= 1 {
		return "", fmt.Errorf("collection name is error, collection name: %s", collectionName)
	}
	return strings.Join(suffixSlice[1:], "_"), nil
}

// GetObjectInstAsstTableName return the object instance association table name in sharding mode base on
// the object ID. Format: InstAsst_{supplierAccount}_{Specifier}_{ObjectID}, such as 'InstAsst_0_pub_switch'.
func GetObjectInstAsstTableName(objID, tenantID string) string {
	return fmt.Sprintf("%s%s_%s_%s", BKObjectInstAsstShardingTablePrefix, tenantID, TableSpecifierPublic, objID)
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
	// check object instance table, ObjectBase_{Specifier}_{ObjectID}
	return strings.HasPrefix(tableName, BKObjectInstShardingTablePrefix)
}

// IsObjectInstAsstShardingTable returns if the target table is an object instance association sharding table.
func IsObjectInstAsstShardingTable(tableName string) bool {
	// check object instance association table, InstAsst_{Specifier}_{ObjectID}
	return strings.HasPrefix(tableName, BKObjectInstAsstShardingTablePrefix)
}

// GetInstTableName returns inst data table name
func GetInstTableName(objID, tenantID string) string {
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
	default:
		return GetObjectInstTableName(objID, tenantID)
	}
}

// GetInstObjIDByTableName get objID by table name
func GetInstObjIDByTableName(collectionName, tenantID string) (string, error) {
	switch collectionName {
	case BKTableNameBaseApp:
		return BKInnerObjIDApp, nil
	case BKTableNameBaseBizSet:
		return BKInnerObjIDBizSet, nil
	case BKTableNameBaseProject:
		return BKInnerObjIDProject, nil
	case BKTableNameBaseSet:
		return BKInnerObjIDSet, nil
	case BKTableNameBaseModule:
		return BKInnerObjIDModule, nil
	case BKTableNameBaseHost:
		return BKInnerObjIDHost, nil
	case BKTableNameBaseProcess:
		return BKInnerObjIDProc, nil
	case BKTableNameBasePlat:
		return BKInnerObjIDPlat, nil
	default:
		return GetObjectInstObjIDByTableName(collectionName, tenantID)
	}
}

var platformTableMap = map[string]struct{}{
	BKTableNameSystem:             {},
	BKTableNameIDgenerator:        {},
	BKTableNameTenant:             {},
	BKTableNameTenantTemplate:     {},
	BKTableNamePlatformAuditLog:   {},
	BKTableNameWatchToken:         {},
	BKTableNameAPITask:            {},
	BKTableNameAPITaskSyncHistory: {},
	BKTableNameWatchDBRelation:    {},
}

// IsPlatformTable returns if the target table is a platform table
func IsPlatformTable(tableName string) bool {
	_, exists := platformTableMap[tableName]
	return exists
}

var platformTableWithTenantMap = map[string]struct{}{
	BKTableNameAPITask:            {},
	BKTableNameAPITaskSyncHistory: {},
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

	// TODO compatible for old table name with cc_ prefix, change this after this prefix is removed
	sepIdx := strings.LastIndex(tenantTableName, "_cc_")
	if sepIdx == -1 {
		return "", "", errors.New("tenant table name is invalid")
	}
	return tenantTableName[:sepIdx], tenantTableName[sepIdx+1:], nil
}
