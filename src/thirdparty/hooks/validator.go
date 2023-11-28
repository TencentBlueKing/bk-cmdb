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

package hooks

import (
	"configcenter/src/apimachinery"
	ccErr "configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/storage/dal"
)

// IsSkipValidateHook is a hook to check if a insert or update option is need to validate or not.
// and check the resource's insert/operate data is valid.
func IsSkipValidateHook(kit *rest.Kit, objID string, data mapstr.MapStr) (bool, error) {
	return false, nil
}

// ValidUpdateCloudIDHook is a hook to check if an update operation on host cloud ID field is valid or not
func ValidUpdateCloudIDHook(kit *rest.Kit, objID string, originInst mapstr.MapStr, updateData mapstr.MapStr) error {
	return nil
}

// ValidateBizBsTopoHook is a hook to check if biz bk_bs_topo field is valid or not
func ValidateBizBsTopoHook(kit *rest.Kit, objID string, originData mapstr.MapStr, updateData mapstr.MapStr,
	validType string, clientSet apimachinery.ClientSetInterface) error {

	return nil
}

// ValidateHostBsInfoHook is a hook to check if host bk_bs_info field is valid or not
func ValidateHostBsInfoHook(kit *rest.Kit, objID string, data mapstr.MapStr) error {
	return nil
}

// ValidHostTransferHook is a hook to check if host transfer parameter is valid or not
func ValidHostTransferHook(kit *rest.Kit, db dal.DB, crossBizTransfer bool, srcBizIDs []int64,
	destBizID int64) ccErr.CCErrorCoder {

	return nil
}

// ValidBizSetPropertyHook is a hook to check if a specific property id is valid or not
func ValidBizSetPropertyHook(kit *rest.Kit, fieldInfo *metadata.BizSetScopeParamsInfo, info metadata.Attribute,
	propertyID interface{}) (bool, error) {
	return false, nil
}

// ValidHostCloudIDHook valid host cloud id hook
func ValidHostCloudIDHook(kit *rest.Kit, cloudID int64) ccErr.CCErrorCoder {
	return nil
}

// IsSkipValidateKeyHook is a hook to check if a insert or update option data key's value need to validate or not.
func IsSkipValidateKeyHook(kit *rest.Kit, objID string, key string, data mapstr.MapStr) (bool, error) {
	return false, nil
}

// ValidUpdateHostStatusHook is a hook to check if an update operation on host status field is valid or not
func ValidUpdateHostStatusHook(kit *rest.Kit, cs apimachinery.ClientSetInterface, objID string,
	originInst mapstr.MapStr, updateData mapstr.MapStr) error {

	return nil
}

// ValidHostApplyStatusHook is a hook to check if host apply status is valid or not
func ValidHostApplyStatusHook(kit *rest.Kit, cs apimachinery.ClientSetInterface, attrID string,
	value interface{}) ccErr.CCErrorCoder {

	return nil
}

// CanUpdateHostApplyStatusHook is a hook to check if host status can be updated by host apply
func CanUpdateHostApplyStatusHook(kit *rest.Kit, cs apimachinery.ClientSetInterface, attrID string,
	originalValue, expectValue interface{}) (bool, ccErr.CCErrorCoder) {

	return true, nil
}

// HostApplyUpdateInfo defines host apply update info
type HostApplyUpdateInfo struct {
	HostIDs    []int64
	Attributes []metadata.HostAttribute
}

// GetHostApplyUpdateInfoHook is a hook to get host apply update info
func GetHostApplyUpdateInfoHook(kit *rest.Kit, cs apimachinery.ClientSetInterface, rules []metadata.HostAttribute,
	hostIDs []int64, attrMap map[int64]string) ([]HostApplyUpdateInfo, ccErr.CCErrorCoder) {

	return []HostApplyUpdateInfo{{HostIDs: hostIDs, Attributes: rules}}, nil
}
