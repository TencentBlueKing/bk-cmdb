/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 Tencent. All rights reserved.
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

package host

import (
	"strings"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	ccErr "configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/storage/driver/mongodb"
)

// UpdateHostInnerIP updates the inner ip fields of a static host directly. It is used by
// the dangerous host inner ip modification api, so it bypasses the generic instance
// update validation.
//
// The source ip is validated here, not only in host_server, so calling this coreservice
// api directly can not bypass the source ip check. The source ip is also used as the
// update condition, so that the ip can not be overwritten when it is modified between
// the read and the write.
func (hm *hostManager) UpdateHostInnerIP(kit *rest.Kit,
	input metadata.ModifyHostInnerIPOpt) (*metadata.UpdateHostInnerIPData, ccErr.CCErrorCoder) {

	if rawErr := input.Validate(); rawErr.ErrCode != 0 {
		return nil, rawErr.ToCCError(kit.CCError)
	}

	host, err := getHostInnerIP(kit, input.HostID)
	if err != nil {
		return nil, err
	}

	addressing := util.GetStrByInterface(host[common.BKAddressingField])
	if addressing != common.BKAddressingStatic {
		blog.Errorf("host %d addressing type is %s, only static host can modify inner ip, rid: %s", input.HostID,
			addressing, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrHostModifyInnerIPNotStatic, input.HostID)
	}

	update, err := buildHostInnerIPUpdate(kit, input, host)
	if err != nil {
		return nil, err
	}
	if len(update.updateDoc) == 0 {
		return &metadata.UpdateHostInnerIPData{Changed: false}, nil
	}

	now := time.Now()
	update.updateDoc[common.LastTimeField] = now
	update.updateDoc[common.BKUpdatedBy] = kit.User
	update.updateDoc[common.BKUpdatedAt] = now

	// the changed ip fields must still be the same as the values read before the update,
	// otherwise the request would overwrite the data modified concurrently.
	hostFilter := map[string]interface{}{common.BKHostIDField: input.HostID}
	for _, field := range update.changedFields {
		current := util.GetStrByInterface(host[field])
		if current == "" {
			// an empty ip list is stored as null
			hostFilter[field] = nil
			continue
		}
		// $all + $size compares the ip content without considering the order, which is the
		// same as the source ip confirmation.
		ips := strings.Split(current, ",")
		hostFilter[field] = map[string]interface{}{common.BKDBAll: ips, common.BKDBSize: len(ips)}
	}

	modified, dbErr := mongodb.Client().Table(common.BKTableNameBaseHost).UpdateMany(kit.Ctx, hostFilter,
		update.updateDoc)
	if dbErr != nil {
		blog.Errorf("update host %d inner ip failed, err: %v, rid: %s", input.HostID, dbErr, kit.Rid)
		if mongodb.Client().IsDuplicatedError(dbErr) {
			return nil, kit.CCError.CCErrorf(common.CCErrCommDuplicateItem, mongodb.GetDuplicateKey(dbErr))
		}
		return nil, kit.CCError.CCError(common.CCErrCommDBUpdateFailed)
	}
	if modified == 0 {
		// the update condition is no longer satisfied, the host is deleted or its inner ip
		// is modified concurrently, the request must fail instead of overwriting the data.
		blog.Errorf("update host %d inner ip failed, the host is deleted or changed concurrently, rid: %s",
			input.HostID, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrCommDBUpdateFailed)
	}

	// the pre host data is returned so that the api layer can generate the audit log in the
	// same transaction, it is not exposed to the api caller.
	return &metadata.UpdateHostInnerIPData{
		Changed:       true,
		ChangedFields: update.changedFields,
		PreHost:       mapstr.MapStr(host),
	}, nil
}

// getHostInnerIP queries the addressing and inner ip fields of the host.
//
// metadata.HostMapStr is used because it is the type allowed by the mongo driver for the
// host table query, and the special ip fields are converted to comma separated strings
// when they are loaded.
func getHostInnerIP(kit *rest.Kit, hostID int64) (metadata.HostMapStr, ccErr.CCErrorCoder) {
	hostFilter := map[string]interface{}{common.BKHostIDField: hostID}
	host := make(metadata.HostMapStr)
	if err := mongodb.Client().Table(common.BKTableNameBaseHost).Find(hostFilter).
		Fields(common.BKHostIDField, common.BKAddressingField, common.BKHostInnerIPField,
			common.BKHostInnerIPv6Field).
		One(kit.Ctx, &host); err != nil {
		if mongodb.Client().IsNotFoundError(err) {
			return nil, kit.CCError.CCError(common.CCErrHostNotFound)
		}
		blog.Errorf("query host %d inner ip failed, err: %v, rid: %s", hostID, err, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}
	return host, nil
}

// buildHostInnerIPUpdate validates the source ip and builds the fields to be updated.
// The returned doc is empty if no inner ip field is changed.
func buildHostInnerIPUpdate(kit *rest.Kit, input metadata.ModifyHostInnerIPOpt,
	host metadata.HostMapStr) (*hostInnerIPUpdate, ccErr.CCErrorCoder) {

	update := &hostInnerIPUpdate{updateDoc: make(map[string]interface{}, 2)}
	families := []struct {
		field  string
		isIPv6 bool
		raw    string
		src    []string
		dest   []string
	}{
		{
			field: common.BKHostInnerIPField,
			raw:   util.GetStrByInterface(host[common.BKHostInnerIPField]),
			src:   input.SrcIPv4,
			dest:  input.DestIPv4,
		},
		{
			field:  common.BKHostInnerIPv6Field,
			isIPv6: true,
			raw:    util.GetStrByInterface(host[common.BKHostInnerIPv6Field]),
			src:    input.SrcIPv6,
			dest:   input.DestIPv6,
		},
	}

	for _, family := range families {
		// Validate guarantees src and dest are provided together; both nil means the
		// family is not modified.
		if family.src == nil && family.dest == nil {
			continue
		}

		// the stored value is historical data, normalize and validate it before comparing.
		// the historical ipv6 may be stored in the abbreviated format.
		normalize := metadata.NormalizeHostIPv4List
		if family.isIPv6 {
			normalize = metadata.NormalizeHostIPv6List
		}
		current, err := normalize(splitHostIPList(family.raw))
		if err != nil {
			blog.Errorf("the stored %s of host %d is invalid, err: %v, rid: %s", family.field,
				input.HostID, err, kit.Rid)
			return nil, kit.CCError.CCErrorf(common.CCErrHostModifyInnerIPStoredValueInvalid, input.HostID,
				family.field, family.raw)
		}

		if !metadata.HostIPListEqual(current, family.src) {
			blog.Errorf("the %s of host %d does not match the source ip, rid: %s", family.field,
				input.HostID, kit.Rid)
			return nil, kit.CCError.CCErrorf(common.CCErrHostModifyInnerIPSourceMismatch, input.HostID,
				family.field, metadata.ConvertHostIPListToString(current),
				metadata.ConvertHostIPListToString(family.src))
		}

		// the destination ip is the same as the current ip in the same order, there is
		// nothing to update. the order is meaningful because the first ip is the primary ip.
		if metadata.HostIPListEqualInOrder(current, family.dest) {
			continue
		}

		update.changedFields = append(update.changedFields, family.field)
		if len(family.dest) == 0 {
			// use nil to clear the field, so that the host special field unique index will
			// not be hit by an empty array.
			update.updateDoc[family.field] = nil
			continue
		}
		update.updateDoc[family.field] = family.dest
	}

	return update, nil
}

// hostInnerIPUpdate is the update document, the changed fields and the update condition of
// a host inner ip update.
type hostInnerIPUpdate struct {
	updateDoc     map[string]interface{}
	changedFields []string
}

// splitHostIPList splits the comma separated ip string, the elements are trimmed and the
// duplicates are removed, it is used to compare the stored ip with the source ip.
func splitHostIPList(value string) []string {
	if value == "" {
		return nil
	}

	result := make([]string, 0)
	for _, item := range strings.Split(value, ",") {
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}
		result = append(result, item)
	}
	return util.RemoveDuplicatesAndEmptyByMap(result)
}
