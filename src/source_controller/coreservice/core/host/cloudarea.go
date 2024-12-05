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
	"strings"
	"sync"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/storage/driver/mongodb"
	"configcenter/src/thirdparty/hooks"
)

// UpdateHostCloudAreaField update host cloud area field
func (hm *hostManager) UpdateHostCloudAreaField(kit *rest.Kit,
	input metadata.UpdateHostCloudAreaFieldOption) errors.CCErrorCoder {

	if len(input.HostIDs) == 0 {
		return kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, "bk_host_ids")
	}
	input.HostIDs = util.IntArrayUnique(input.HostIDs)

	if err := validCloudID(kit, input.CloudID); err != nil {
		return err
	}

	if err := validHost(kit, input.CloudID, input.HostIDs); err != nil {
		return err
	}

	updateFilter := map[string]interface{}{
		common.BKHostIDField: map[string]interface{}{
			common.BKDBIN: input.HostIDs,
		},
	}
	updateDoc := map[string]interface{}{
		common.BKCloudIDField: input.CloudID,
	}
	if err := mongodb.Shard(kit.ShardOpts()).Table(common.BKTableNameBaseHost).Update(kit.Ctx, updateFilter,
		updateDoc); err != nil {
		blog.Errorf("update host cloud area failed, filter: %v, doc: %v, err: %v, rid: %s", updateFilter, updateDoc,
			err, kit.Rid)
		return kit.CCError.CCError(common.CCErrCommDBUpdateFailed)
	}
	return nil
}

func validCloudID(kit *rest.Kit, cloudID int64) errors.CCErrorCoder {
	if err := hooks.ValidHostCloudIDHook(kit, cloudID); err != nil {
		return err
	}
	cloudIDFiler := map[string]interface{}{
		common.BKCloudIDField: cloudID,
	}
	count, err := mongodb.Shard(kit.ShardOpts()).Table(common.BKTableNameBasePlat).Find(cloudIDFiler).Count(kit.Ctx)
	if err != nil {
		blog.Errorf("find cloud area failed, option: %v, err: %v, rid: %s", cloudIDFiler, err, kit.Rid)
		return kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}

	if count == 0 {
		blog.Errorf("bk_cloud_id is invalid, bk_cloud_id: %d, rid: %s", cloudID, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKCloudIDField)
	}
	if count > 1 {
		blog.Errorf("get multiple cloud area, bk_cloud_id: %d, rid: %s", cloudID, kit.Rid)
		return kit.CCError.CCError(common.CCErrCommGetMultipleObject)
	}

	return nil
}

func validHost(kit *rest.Kit, cloudID int64, hostIDs []int64) errors.CCErrorCoder {
	// step1. validate bk_host_ids is exist
	hostFilter := map[string]interface{}{
		common.BKHostIDField: map[string]interface{}{
			common.BKDBIN: hostIDs,
		},
	}
	hostSimplify := make([]metadata.HostMapStr, 0)
	fields := []string{common.BKHostInnerIPField, common.BKHostInnerIPv6Field, common.BKCloudIDField,
		common.BKHostIDField, common.BKAddressingField}

	err := mongodb.Shard(kit.ShardOpts()).Table(common.BKTableNameBaseHost).Find(hostFilter).Fields(fields...).
		All(kit.Ctx, &hostSimplify)
	if err != nil {
		blog.Errorf("find host failed, option: %v, err: %v, rid: %s", hostFilter, err, kit.Rid)
		return kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}
	if len(hostIDs) != len(hostSimplify) {
		blog.Errorf("some hosts not found, hostIDs: %v, hosts: %v, rid: %s", hostIDs, hostSimplify, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKHostIDField)
	}

	// step2. validate when addressing is static,
	// unique of bk_cloud_id + bk_host_innerip and bk_cloud_id + bk_host_innerip_v6
	innerIPv4s := make([]string, 0)
	ipv4Map := make(map[string]struct{})
	innerIPv6s := make([]string, 0)
	ipv6Map := make(map[string]struct{})
	for _, item := range hostSimplify {
		addressing, ok := item[common.BKAddressingField].(string)
		if !ok {
			blog.Errorf("host field is invalid, field: %s, host: %v, rid: %s", common.BKAddressingField, item, kit.Rid)
			return kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKAddressingField)
		}

		if addressing != common.BKAddressingStatic {
			continue
		}

		ipv4, ok := item[common.BKHostInnerIPField].(string)
		if ok {
			if _, ok := ipv4Map[ipv4]; ok {
				return kit.CCError.CCErrorf(common.CCErrCommDuplicateItem, common.BKHostInnerIPField)
			}
			innerIPv4s = append(innerIPv4s, ipv4)
			ipv4Map[ipv4] = struct{}{}
		}

		ipv6, ok := item[common.BKHostInnerIPv6Field].(string)
		if ok {
			if _, ok := ipv6Map[ipv6]; ok {
				return kit.CCError.CCErrorf(common.CCErrCommDuplicateItem, common.BKHostInnerIPv6Field)
			}
			innerIPv6s = append(innerIPv6s, ipv6)
			ipv6Map[ipv6] = struct{}{}
		}
	}

	// step3. validate when addressing is static,
	// unique of bk_cloud_id + bk_inner_ip and bk_cloud_id + bk_host_innerip_v6 and in database
	if err := validDuplicatedHostInDB(kit, hostIDs, cloudID, innerIPv4s, innerIPv6s); err != nil {
		return err
	}

	return nil
}

func validDuplicatedHostInDB(kit *rest.Kit, hostIDs []int64, cloudID int64, innerIPv4s []string,
	innerIPv6s []string) errors.CCErrorCoder {

	ipCond := make([]map[string]interface{}, len(innerIPv4s)+len(innerIPv6s))
	for idx, ip := range innerIPv4s {
		ipArr := strings.Split(ip, ",")
		ipCond[idx] = map[string]interface{}{
			common.BKHostInnerIPField: map[string]interface{}{
				common.BKDBIN: ipArr,
			},
		}
	}

	for idx, ip := range innerIPv6s {
		ipArr := strings.Split(ip, ",")
		ipCond[idx+len(innerIPv4s)] = map[string]interface{}{
			common.BKHostInnerIPv6Field: map[string]interface{}{
				common.BKDBIN: ipArr,
			},
		}
	}

	dbHostFilter := map[string]interface{}{
		common.BKAddressingField: common.BKAddressingStatic,
		common.BKHostIDField: map[string]interface{}{
			common.BKDBNIN: hostIDs,
		},
		common.BKCloudIDField: cloudID,
		common.BKDBOR:         ipCond,
	}
	duplicatedHosts := make([]metadata.HostMapStr, 0)
	fields := []string{common.BKHostInnerIPField, common.BKHostInnerIPv6Field, common.BKCloudIDField,
		common.BKHostIDField, common.BKAddressingField}

	err := mongodb.Shard(kit.ShardOpts()).Table(common.BKTableNameBaseHost).Find(dbHostFilter).Fields(fields...).
		All(kit.Ctx, &duplicatedHosts)
	if err != nil {
		blog.Errorf("find host failed, option: %v, err: %v, rid: %s", dbHostFilter, err, kit.Rid)
		return kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}

	if len(duplicatedHosts) > 0 {
		blog.Errorf("duplicated hosts exits, option: %v, host: %v, rid: %s", dbHostFilter, duplicatedHosts, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCommDuplicateItem,
			common.BKHostInnerIPField+" or "+common.BKHostInnerIPv6Field)
	}

	return nil
}

// FindCloudAreaHostCount TODO
func (hm *hostManager) FindCloudAreaHostCount(kit *rest.Kit, input metadata.CloudAreaHostCount) ([]metadata.
	CloudAreaHostCountElem, error) {
	if len(input.CloudIDs) == 0 {
		return nil, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, "bk_cloud_ids")
	}

	cloudIDs := util.IntArrayUnique(input.CloudIDs)

	// to speed up, multi goroutine to query host count for multi cloudarea
	var wg sync.WaitGroup
	var lock sync.RWMutex
	var firstErr errors.CCErrorCoder
	pipeline := make(chan bool, 10)
	cloudCountMap := make(map[int64]int64)

	for _, cloudID := range cloudIDs {
		pipeline <- true
		wg.Add(1)

		go func(cloudID int64) {
			defer func() {
				wg.Done()
				<-pipeline
			}()

			filter := map[string]interface{}{common.BKCloudIDField: cloudID}
			hostCnt, err := mongodb.Shard(kit.ShardOpts()).Table(common.BKTableNameBaseHost).Find(filter).
				Count(kit.Ctx)
			if err != nil {
				blog.Errorf("count host cloud area failed, table: %s, filter: %v, err: %v, rid: %s",
					common.BKTableNameBaseHost, filter, err, kit.Rid)
				if firstErr == nil {
					firstErr = kit.CCError.CCError(common.CCErrCommDBSelectFailed)
				}
				return
			}

			lock.Lock()
			cloudCountMap[cloudID] = int64(hostCnt)
			lock.Unlock()

		}(cloudID)
	}

	wg.Wait()

	if firstErr != nil {
		return nil, firstErr
	}

	ret := make([]metadata.CloudAreaHostCountElem, len(input.CloudIDs))
	for idx, cloudID := range input.CloudIDs {
		ret[idx] = metadata.CloudAreaHostCountElem{
			CloudID:   cloudID,
			HostCount: cloudCountMap[cloudID],
		}
	}

	return ret, nil
}
