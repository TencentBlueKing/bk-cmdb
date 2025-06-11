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

package logics

import (
	"fmt"

	"configcenter/src/apimachinery"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/storage/dal/mongo/local"
)

// GetObjUUIDFromCache get object uuid by object id from cache
func GetObjUUIDFromCache(kit *rest.Kit, clientSet apimachinery.ClientSetInterface, objID string) (string, error) {
	if clientSet == nil {
		blog.Errorf("client set is nil, rid: %s", kit.Rid)
		return "", fmt.Errorf("client set is nil")
	}

	uuid, err := clientSet.CacheService().Cache().Object().GetUUIDByObj(kit.Ctx, kit.Header, objID)
	if err != nil {
		blog.Errorf("get object %s uuid failed from core service, err: %v, rid: %s", objID, err, kit.Rid)
		return "", err
	}

	return uuid, nil
}

// GetObjUUIDFromDB get object uuid by object id form db
func GetObjUUIDFromDB(kit *rest.Kit, db local.DB, objID string) (string, error) {
	result := new(metadata.Object)
	err := db.Table(common.BKTableNameObjDes).Find(mapstr.MapStr{common.BKObjIDField: objID}).Fields(
		metadata.ModelFieldObjUUID).One(kit.Ctx, result)
	if err != nil {
		blog.Errorf("get object %s uuid failed from db, err: %v", objID, err)
		return "", err
	}
	return result.UUID, nil
}

// GetObjInstTableFromDB get object instance table name from db
func GetObjInstTableFromDB(kit *rest.Kit, db local.DB, objID string) (string, error) {
	if common.IsInnerModel(objID) {
		return common.GetInnerInstTableName(objID), nil
	}

	objUUID, err := GetObjUUIDFromDB(kit, db, objID)
	if err != nil {
		blog.Errorf("get object uuid failed from db, err: %v", err)
		return "", err
	}

	return common.GetObjInstTableName(objUUID), nil
}

// GetObjInstTableFromCache get object instance table name from cache
func GetObjInstTableFromCache(kit *rest.Kit, clientSet apimachinery.ClientSetInterface, objID string) (string, error) {
	if common.IsInnerModel(objID) {
		return common.GetInnerInstTableName(objID), nil
	}

	objUUID, err := GetObjUUIDFromCache(kit, clientSet, objID)
	if err != nil {
		blog.Errorf("get object uuid failed from db, err: %v", err)
		return "", err
	}

	return common.GetObjInstTableName(objUUID), nil
}

// GetObjInstAsstTableFromDB get object instance association table name from db
func GetObjInstAsstTableFromDB(kit *rest.Kit, db local.DB, objID string) (string, error) {
	objUUID, err := GetObjUUIDFromDB(kit, db, objID)
	if err != nil {
		blog.Errorf("get object uuid %s failed, object: %s, err: %v, rid: %s", objID, err, kit.Rid)
		return "", err
	}

	return common.GetObjInstAsstTableName(objUUID), nil
}

// GetObjInstAsstTableFromCache get object instance association table name from cache
func GetObjInstAsstTableFromCache(kit *rest.Kit, clientSet apimachinery.ClientSetInterface, objID string) (
	string, error) {

	objUUID, err := GetObjUUIDFromCache(kit, clientSet, objID)
	if err != nil {
		blog.Errorf("get object uuid %s failed, object: %s, err: %v, rid: %s", objID, err, kit.Rid)
		return "", err
	}

	return common.GetObjInstAsstTableName(objUUID), nil
}
