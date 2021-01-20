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
)

func (hm *hostManager) UpdateHostCloudAreaField(kit *rest.Kit, input metadata.UpdateHostCloudAreaFieldOption) errors.CCErrorCoder {
	rid := kit.Rid
	context := kit.Ctx

	if len(input.HostIDs) == 0 {
		return kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, "bk_host_ids")
	}
	input.HostIDs = util.IntArrayUnique(input.HostIDs)

	// step1. validate bk_cloud_id
	cloudIDFiler := map[string]interface{}{
		common.BKCloudIDField: input.CloudID,
	}
	count, err := mongodb.Client().Table(common.BKTableNameBasePlat).Find(cloudIDFiler).Count(context)
	if err != nil {
		blog.ErrorJSON("UpdateHostCloudAreaField failed, db select failed, table: %s, option: %s, err: %s, rid: %s", common.BKTableNameBasePlat, cloudIDFiler, err.Error(), rid)
		return kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}
	if count == 0 {
		blog.Errorf("UpdateHostCloudAreaField failed, bk_cloud_id invalid, bk_cloud_id: %d, rid: %s", input.CloudID, rid)
		return kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKCloudIDField)
	}
	if count > 1 {
		blog.Errorf("UpdateHostCloudAreaField failed, get multiple cloud area, bk_cloud_id: %d, rid: %s", input.CloudID, rid)
		return kit.CCError.CCError(common.CCErrCommGetMultipleObject)
	}

	// step2. validate bk_host_ids
	hostFilter := map[string]interface{}{
		common.BKHostIDField: map[string]interface{}{
			common.BKDBIN: input.HostIDs,
		},
	}
	hostSimplify := make([]metadata.HostMapStr, 0)
	fields := []string{common.BKHostInnerIPField, common.BKCloudIDField, common.BKHostIDField}
	if err := mongodb.Client().Table(common.BKTableNameBaseHost).Find(hostFilter).Fields(fields...).All(context, &hostSimplify); err != nil {
		blog.ErrorJSON("UpdateHostCloudAreaField failed, db select failed, table: %s, option: %s, err: %s, rid: %s", common.BKTableNameBaseHost, hostFilter, err.Error(), rid)
		return kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}
	if len(input.HostIDs) != len(hostSimplify) {
		blog.Errorf("UpdateHostCloudAreaField failed, maybe some hosts not found, hostIDs:%s, hosts:%s, rid:%s", input.HostIDs, hostSimplify, rid)
		return kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKHostIDField)
	}

	// step3. validate unique of bk_cloud_id + bk_host_innerip in input parameters
	innerIPs := make([]string, 0)
	for _, item := range hostSimplify {
		innerIPs = append(innerIPs, item[common.BKHostInnerIPField].(string))
	}
	innerIPs = util.StrArrayUnique(innerIPs)
	if len(input.HostIDs) != len(innerIPs) {
		return kit.CCError.CCErrorf(common.CCErrCommDuplicateItem, common.BKHostInnerIPField)
	}

	// step4. validate unique of bk_cloud_id + bk_inner_ip in database
	ipCond := make([]map[string]interface{}, len(innerIPs))
	for index, innerIP := range innerIPs {
		innerIPArr := strings.Split(innerIP, ",")
		ipCond[index] = map[string]interface{}{
			common.BKHostInnerIPField: map[string]interface{}{
				common.BKDBIN: innerIPArr,
			},
		}
	}
	dbHostFilter := map[string]interface{}{
		common.BKHostIDField: map[string]interface{}{
			common.BKDBNIN: input.HostIDs,
		},
		common.BKCloudIDField: input.CloudID,
		common.BKDBOR:         ipCond,
	}
	duplicatedHosts := make([]metadata.HostMapStr, 0)
	if err := mongodb.Client().Table(common.BKTableNameBaseHost).Find(dbHostFilter).Fields(fields...).All(context, &duplicatedHosts); err != nil {
		blog.ErrorJSON("UpdateHostCloudAreaField failed, db select failed, table: %s, option: %s, err: %s, rid: %s", common.BKTableNameBaseHost, dbHostFilter, err.Error(), rid)
		return kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}
	if len(duplicatedHosts) > 0 {
		blog.ErrorJSON("UpdateHostCloudAreaField failed, bk_cloud_id + bk_host_innerip duplicated, input: %s, duplicated hosts: %s, rid: %s", input, duplicatedHosts, rid)
		return kit.CCError.CCErrorf(common.CCErrCommDuplicateItem, common.BKHostInnerIPField)
	}

	// step5. update hosts bk_cloud_id field
	updateFilter := map[string]interface{}{
		common.BKHostIDField: map[string]interface{}{
			common.BKDBIN: input.HostIDs,
		},
	}
	updateDoc := map[string]interface{}{
		common.BKCloudIDField: input.CloudID,
	}
	if err := mongodb.Client().Table(common.BKTableNameBaseHost).Update(context, updateFilter, updateDoc); err != nil {
		blog.ErrorJSON("UpdateHostCloudAreaField failed, db update failed, table: %s, filter: %s, doc: %s, err: %s, rid: %s", common.BKTableNameBaseHost, updateFilter, updateDoc, err.Error(), rid)
		return kit.CCError.CCError(common.CCErrCommDBUpdateFailed)
	}
	return nil
}

func (hm *hostManager) FindCloudAreaHostCount(kit *rest.Kit, input metadata.CloudAreaHostCount) ([]metadata.CloudAreaHostCountElem, error) {
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
			hostCnt, err := mongodb.Client().Table(common.BKTableNameBaseHost).Find(filter).Count(kit.Ctx)
			if err != nil {
				blog.ErrorJSON("UpdateHostCloudAreaField failed, db selected failed, table: %s, filter: %s, err: %s, rid: %s", common.BKTableNameBaseHost, filter, err.Error(), kit.Rid)
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
