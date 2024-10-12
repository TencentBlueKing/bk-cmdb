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

package synchronize

import (
	"fmt"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

type hostSyncer struct {
}

// ParseDataArr parse data array to actual type
func (h *hostSyncer) ParseDataArr(kit *rest.Kit, data any) (any, error) {
	hosts, err := parseDataArr[metadata.HostMapStr](kit, data)
	if err != nil {
		return nil, err
	}

	for idx := range hosts {
		hosts[idx], err = metadata.ConvertHostSpecialStringToArray(hosts[idx])
		if err != nil {
			blog.Errorf("convert host(%+v) special string to array failed, err: %v, rid: %s", hosts[idx], err, kit.Rid)
			return nil, err
		}
	}
	return hosts, nil
}

// Validate host sync data
func (h *hostSyncer) Validate(kit *rest.Kit, subRes string, data any) error {
	hosts, ok := data.([]metadata.HostMapStr)
	if !ok {
		return kit.CCError.New(common.CCErrCommParamsIsInvalid, fmt.Sprintf("data type %T is invalid", data))
	}

	cloudIDs := make([]int64, 0)
	for _, host := range hosts {
		hostID, err := util.GetInt64ByInterface(host[common.BKHostIDField])
		if err != nil {
			blog.Errorf("parse host id(%v) failed, err: %v, rid: %s", host[common.BKHostIDField], err, kit.Rid)
			return kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, common.BKHostIDField)
		}
		if hostID <= 0 {
			blog.Errorf("host id(%d) is invalid,  rid: %s", hostID, kit.Rid)
			return kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, common.BKHostIDField)
		}

		cloudID, err := util.GetInt64ByInterface(host[common.BKCloudIDField])
		if err != nil {
			blog.Errorf("parse cloud id(%v) failed, err: %v, rid: %s", host[common.BKCloudIDField], err, kit.Rid)
			return kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, common.BKCloudIDField)
		}
		cloudIDs = append(cloudIDs, cloudID)
	}

	if err := validateDependency(kit, common.BKTableNameBasePlat, common.BKCloudIDField, cloudIDs); err != nil {
		return err
	}

	return nil
}

// TableName returns table name for host syncer
func (h *hostSyncer) TableName(subRes, supplierAccount string) string {
	return common.BKTableNameBaseHost
}

type hostRelSyncer struct {
}

// ParseDataArr parse data array to actual type
func (h *hostRelSyncer) ParseDataArr(kit *rest.Kit, data any) (any, error) {
	return parseDataArr[metadata.ModuleHost](kit, data)
}

// Validate host relation sync data
func (h *hostRelSyncer) Validate(kit *rest.Kit, subRes string, data any) error {
	relations, ok := data.([]metadata.ModuleHost)
	if !ok {
		return kit.CCError.New(common.CCErrCommParamsIsInvalid, fmt.Sprintf("data type %T is invalid", data))
	}

	hostIDs, moduleIDs := make([]int64, 0), make([]int64, 0)
	for _, rel := range relations {
		hostIDs = append(hostIDs, rel.HostID)
		moduleIDs = append(moduleIDs, rel.ModuleID)
	}

	if err := validateDependency(kit, common.BKTableNameBaseModule, common.BKModuleIDField, moduleIDs); err != nil {
		return err
	}

	if err := validateDependency(kit, common.BKTableNameBaseHost, common.BKHostIDField, hostIDs); err != nil {
		return err
	}

	return nil
}

// TableName returns table name for host relation syncer
func (h *hostRelSyncer) TableName(subRes, supplierAccount string) string {
	return common.BKTableNameModuleHostConfig
}
