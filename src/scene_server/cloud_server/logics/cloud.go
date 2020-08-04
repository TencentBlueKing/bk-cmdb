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

package logics

import (
	"fmt"
	"strconv"
	"sync"

	"configcenter/src/ac/meta"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/cloud_server/cloudvendor"
	ccom "configcenter/src/scene_server/cloud_server/common"
)

// AccountVerify 验证云账户连通性
func (lgc *Logics) AccountVerify(kit *rest.Kit, conf metadata.CloudAccountConf) error {
	client, err := cloudvendor.GetVendorClient(conf)
	if err != nil {
		blog.Errorf("AccountVerify GetVendorClient failed, AccountID:%d, err:%s, rid:%s", conf.AccountID, err.Error(), kit.Rid)
		return err
	}

	_, err = client.GetRegions()
	if err != nil {
		blog.Errorf("AccountVerify GetRegions failed, AccountID:%d, err:%s, rid:%s", conf.AccountID, err.Error(), kit.Rid)
		return err
	}

	return nil
}

// SearchAccountValidity 查询云账户有效性
func (lgc *Logics) SearchAccountValidity(kit *rest.Kit, confs []metadata.CloudAccountConf) []metadata.AccountValidityInfo {
	result := make([]metadata.AccountValidityInfo, 0)
	validityInfoChan := make(chan metadata.AccountValidityInfo, 10)
	var wg, wg2 sync.WaitGroup

	for _, conf := range confs {
		wg.Add(1)
		go func(conf metadata.CloudAccountConf) {
			defer wg.Done()
			errMsg := ""
			if err := lgc.AccountVerify(kit, conf); err != nil {
				errMsg = err.Error()
			}
			validityInfoChan <- metadata.AccountValidityInfo{
				AccountID: conf.AccountID,
				ErrMsg:    errMsg,
			}
		}(conf)
	}

	wg2.Add(1)
	go func() {
		defer wg2.Done()
		for info := range validityInfoChan {
			result = append(result, info)
		}
	}()

	wg.Wait()
	close(validityInfoChan)
	wg2.Wait()

	return result
}

// GetRegionsInfo 获取地域信息
func (lgc *Logics) GetRegionsInfo(kit *rest.Kit, conf metadata.CloudAccountConf, withHostCount bool) ([]metadata.SyncRegion, error) {
	client, err := cloudvendor.GetVendorClient(conf)
	if err != nil {
		blog.Errorf("GetRegionsInfo GetVendorClient failed, AccountID:%d, err:%s, rid:%s", conf.AccountID, err.Error(), kit.Rid)
		return nil, err
	}

	regionSet, err := client.GetRegions()
	if err != nil {
		blog.Errorf("GetRegionsInfo GetRegions err:%s, rid:%s", err.Error(), kit.Rid)
		return nil, err
	}

	regionHostCnt := make(map[string]int64)
	// 需要获取地域下的主机数
	if withHostCount {
		hostCntChan := make(chan []interface{}, 10)
		var wg, wg2 sync.WaitGroup
		// 并发请求获取每个地域下的主机数
		for _, region := range regionSet {
			wg.Add(1)
			go func(region *metadata.Region) {
				defer wg.Done()
				count, err := client.GetInstancesTotalCnt(region.RegionId, nil)
				if err != nil {
					blog.Errorf("GetVpcHostCnt GetInstances failed, AccountID:%d, err:%s, rid:%s", conf.AccountID, err.Error(), kit.Rid)
					return
				}
				hostCntChan <- []interface{}{region.RegionId, count}
			}(region)
		}
		wg2.Add(1)
		go func() {
			defer wg2.Done()
			for hostCnt := range hostCntChan {
				regionHostCnt[hostCnt[0].(string)] = hostCnt[1].(int64)
			}
		}()
		wg.Wait()
		close(hostCntChan)
		wg2.Wait()
	}

	result := make([]metadata.SyncRegion, 0)
	for i := range regionSet {
		region := regionSet[i]
		result = append(result, metadata.SyncRegion{
			RegionId:    region.RegionId,
			RegionName:  region.RegionName,
			RegionState: region.RegionState,
			HostCount:   regionHostCnt[region.RegionId],
		})
	}

	return result, nil
}

// GetVpcHostCntInOneRegion 获取某地域下的vpc详情和主机数
func (lgc *Logics) GetVpcHostCntInOneRegion(kit *rest.Kit, conf metadata.CloudAccountConf, region string) (*metadata.VpcHostCntResult, error) {
	client, err := cloudvendor.GetVendorClient(conf)
	if err != nil {
		blog.Errorf("GetVpcHostCntInOneRegion GetVendorClient failed, AccountID:%d, err:%s, rid:%s", conf.AccountID, err.Error(), kit.Rid)
		return nil, err
	}

	// 获取该地域下的vpc详情
	vpcsInfo, err := client.GetVpcs(region, nil)
	if err != nil {
		blog.Errorf("GetVpcHostCntInOneRegion GetVpcs err:%s, rid:%s", err.Error(), kit.Rid)
		return nil, err
	}

	result := new(metadata.VpcHostCntResult)
	result.Count = vpcsInfo.Count

	if len(vpcsInfo.VpcSet) == 0 {
		return result, nil
	}

	// 获取vpc对应的主机数
	option := metadata.SearchVpcHostCntOption{}
	for _, vpc := range vpcsInfo.VpcSet {
		option.RegionVpcs = append(option.RegionVpcs, metadata.RegionVpc{
			Region: region,
			VpcID:  vpc.VpcId,
		})
	}
	vpcHostCnt, err := lgc.GetVpcHostCnt(kit, conf, option)
	if err != nil {
		blog.Errorf("GetVpcHostCntInOneRegion GetVpcHostCnt failed, AccountID:%d, err:%s, rid:%s", conf.AccountID, err.Error(), kit.Rid)
		return nil, err
	}

	for i := range vpcsInfo.VpcSet {
		vpc := vpcsInfo.VpcSet[i]
		result.Info = append(result.Info, metadata.VpcSyncInfo{
			VpcID:        vpc.VpcId,
			VpcName:      vpc.VpcName,
			Region:       region,
			VpcHostCount: vpcHostCnt[vpc.VpcId],
		})
	}

	return result, nil
}

// GetVpcHostCnt 获取多个vpc对应的主机数
func (lgc *Logics) GetVpcHostCnt(kit *rest.Kit, conf metadata.CloudAccountConf, option metadata.SearchVpcHostCntOption) (map[string]int64, error) {
	client, err := cloudvendor.GetVendorClient(conf)
	if err != nil {
		blog.Errorf("GetVpcHostCnt GetVendorClient failed, AccountID:%d, err:%s, rid:%s", conf.AccountID, err.Error(), kit.Rid)
		return nil, err
	}

	vpcHostCnt := make(map[string]int64)
	hostCntChan := make(chan []interface{}, 10)
	var wg, wg2 sync.WaitGroup
	// 并发请求获取每个vpc的实例个数
	for _, regionVpc := range option.RegionVpcs {
		wg.Add(1)
		go func(regionVpc metadata.RegionVpc) {
			defer wg.Done()
			count, err := client.GetInstancesTotalCnt(regionVpc.Region, &ccom.InstanceOpt{
				BaseOpt: ccom.BaseOpt{
					Filters: []*ccom.Filter{{ccom.StringPtr("vpc-id"), ccom.StringPtrs([]string{regionVpc.VpcID})}},
					Limit:   ccom.MaxLimit,
				},
			})
			if err != nil {
				blog.Errorf("GetVpcHostCnt GetInstances failed, AccountID:%d, err:%s, rid:%s", conf.AccountID, err.Error(), kit.Rid)
				return
			}
			hostCntChan <- []interface{}{regionVpc.VpcID, count}
		}(regionVpc)
	}
	wg2.Add(1)
	go func() {
		defer wg2.Done()
		for hostCnt := range hostCntChan {
			vpcHostCnt[hostCnt[0].(string)] = hostCnt[1].(int64)
		}
	}()
	wg.Wait()
	close(hostCntChan)
	wg2.Wait()

	return vpcHostCnt, nil
}

// GetCloudHostResource 获取需要同步的云主机资源信息
func (lgc *Logics) GetCloudHostResource(kit *rest.Kit, conf metadata.CloudAccountConf, syncVpcs []metadata.VpcSyncInfo) (*metadata.CloudHostResource, error) {
	client, err := cloudvendor.GetVendorClient(conf)
	if err != nil {
		blog.Errorf("GetCloudHostResource GetVendorClient failed, AccountID:%d, err:%s, rid:%s", conf.AccountID, err.Error(), kit.Rid)
		return nil, err
	}

	blog.V(4).Infof("GetCloudHostResource syncVpcs %#v, rid:%s", syncVpcs, kit.Rid)

	// 不再同步已经被销毁的vpc
	allVpcs := make([]metadata.VpcSyncInfo, 0)
	for _, vpc := range syncVpcs {
		if vpc.Destroyed {
			continue
		}
		allVpcs = append(allVpcs, vpc)
	}

	vpcHostDetail := make(map[string][]*metadata.Instance)
	hostDetailChan := make(chan []*metadata.Instance, 10)
	destroyedVpcs := make(map[string]bool)
	destroyedVpcsChan := make(chan string, 10)
	errs := make([]error, 0)
	errChan := make(chan error, 10)
	var wg, wg2, wg3, wg4 sync.WaitGroup
	// 并发请求获取被销毁的vpc数据和没被销毁的vpc下主机实例详情
	for _, vpc := range allVpcs {
		wg.Add(1)
		go func(vpc metadata.VpcSyncInfo) {
			defer wg.Done()

			vpcInfo, err := client.GetVpcs(vpc.Region, &ccom.VpcOpt{
				BaseOpt: ccom.BaseOpt{
					Filters: []*ccom.Filter{{ccom.StringPtr("vpc-id"), ccom.StringPtrs([]string{vpc.VpcID})}},
					Limit:   ccom.MaxLimit,
				},
			})
			if err != nil {
				blog.Errorf("GetCloudHostResource GetVpcs failed, AccountID:%d, err:%s, rid:%s", conf.AccountID, err.Error(), kit.Rid)
				errChan <- err
				return
			}
			if len(vpcInfo.VpcSet) == 0 {
				blog.Errorf("GetCloudHostResource add destroyed vpcID:%s, AccountID:%d, vpcInfo.VpcSet:%#v, param vpc:%#v, conf:%#v, rid:%s",
					vpc.VpcID, conf.AccountID, vpcInfo.VpcSet, vpc, conf, kit.Rid)
				destroyedVpcsChan <- vpc.VpcID
				return
			}

			instancesInfo, err := client.GetInstances(vpc.Region, &ccom.InstanceOpt{
				BaseOpt: ccom.BaseOpt{
					Filters: []*ccom.Filter{{ccom.StringPtr("vpc-id"), ccom.StringPtrs([]string{vpc.VpcID})}},
					Limit:   ccom.MaxLimit,
				},
			})
			if err != nil {
				blog.Errorf("GetCloudHostResource GetInstances failed, AccountID:%d, err:%s, rid:%s", conf.AccountID, err.Error(), kit.Rid)
				errChan <- err
				return
			}
			blog.V(4).Infof("GetCloudHostResource vpc-id:%s, instances count %#v, rid:%s", vpc.VpcID, instancesInfo.Count, kit.Rid)
			hostDetailChan <- instancesInfo.InstanceSet
		}(vpc)
	}
	// 收集vpc实例详情
	wg2.Add(1)
	go func() {
		defer wg2.Done()
		for hostDetail := range hostDetailChan {
			if len(hostDetail) > 0 {
				vpcHostDetail[hostDetail[0].VpcId] = hostDetail
			}

		}
	}()
	// 收集被销毁的vpc数据
	wg3.Add(1)
	go func() {
		defer wg3.Done()
		for vpc := range destroyedVpcsChan {
			destroyedVpcs[vpc] = true
		}
	}()
	// 收集错误
	wg4.Add(1)
	go func() {
		defer wg4.Done()
		for err := range errChan {
			errs = append(errs, err)
		}
	}()
	wg.Wait()
	close(hostDetailChan)
	close(destroyedVpcsChan)
	close(errChan)
	wg2.Wait()
	wg3.Wait()
	wg4.Wait()

	// 调用云厂商接口出现过错误则直接返回
	if len(errs) > 0 {
		return nil, errs[0]
	}

	result := new(metadata.CloudHostResource)

	for i := range allVpcs {
		vpc := allVpcs[i]
		// 被销毁的vpc数据
		if destroyedVpcs[vpc.VpcID] {
			result.DestroyedVpcs = append(result.DestroyedVpcs, &vpc)
			continue
		}
		// 没被销毁的vpc下的云主机资源数据
		result.HostResource = append(result.HostResource, &metadata.VpcInstances{
			Vpc:       &vpc,
			Instances: vpcHostDetail[vpc.VpcID],
		})
	}

	return result, nil
}

// GetCloudAccountConf 获取云账户配置
func (lgc *Logics) GetCloudAccountConf(kit *rest.Kit, accountID int64) (*metadata.CloudAccountConf, error) {
	option := &metadata.SearchCloudOption{Condition: mapstr.MapStr{common.BKCloudAccountID: accountID}}
	result, err := lgc.CoreAPI.CoreService().Cloud().SearchAccountConf(kit.Ctx, kit.Header, option)
	if err != nil {
		blog.Errorf("SearchAccountConf failed, accountID: %v, err: %s，rid:%s", accountID, err.Error(), kit.Rid)
		return nil, err
	}
	if len(result.Info) == 0 {
		blog.Errorf("GetCloudAccountConf failed, accountID: %v is not exist", accountID)
		return nil, fmt.Errorf("GetAccountConf failed, accountID: %v is not exist, rid:%s", accountID, kit.Rid)
	}

	accountConf := result.Info[0]
	// 解密云账户密钥
	if lgc.cryptor != nil {
		secretKey, err := lgc.cryptor.Decrypt(accountConf.SecretKey)
		if err != nil {
			blog.Errorf("GetCloudAccountConf failed, Encrypt err: %st, rid:%s", err.Error(), kit.Rid)
			return nil, err
		}
		accountConf.SecretKey = secretKey
	}

	return &accountConf, nil
}

// GetCloudAccountConfBatch 批量获取云账户配置
func (lgc *Logics) GetCloudAccountConfBatch(kit *rest.Kit, accountIDs []int64) ([]metadata.CloudAccountConf, error) {
	if len(accountIDs) == 0 {
		return nil, nil
	}

	option := &metadata.SearchCloudOption{
		Condition: mapstr.MapStr{
			common.BKCloudAccountID: mapstr.MapStr{
				common.BKDBIN: accountIDs,
			},
		},
	}
	result, err := lgc.CoreAPI.CoreService().Cloud().SearchAccountConf(kit.Ctx, kit.Header, option)
	if err != nil {
		blog.Errorf("GetCloudAccountConfBatch failed, accountIDs: %v, SearchAccountConf err: %s, rid:%s", accountIDs, err.Error(), kit.Rid)
		return nil, err
	}
	if len(result.Info) == 0 {
		blog.Errorf("GetCloudAccountConfBatch failed, accountIDs: %v are not exist, rid:%s", accountIDs, kit.Rid)
		return nil, fmt.Errorf("GetAccountConf failed, accountIDs: %v are not exist", accountIDs)
	}

	accountConfs := result.Info
	// 解密云账户密钥
	if lgc.cryptor != nil {
		for i, _ := range accountConfs {
			secretKey, err := lgc.cryptor.Decrypt(accountConfs[i].SecretKey)
			if err != nil {
				blog.Errorf("GetCloudAccountConfBatch failed, Encrypt err: %s, rid:%s", err.Error(), kit.Rid)
				return nil, err
			}
			accountConfs[i].SecretKey = secretKey
		}

	}

	return accountConfs, nil
}

func (lgc *Logics) ListAuthorizedResources(kit *rest.Kit, typ meta.ResourceType, act meta.Action) ([]int64, error) {

	authInput := meta.ListAuthorizedResourcesParam{
		UserName:     kit.User,
		ResourceType: typ,
		Action:       act,
	}

	authList, err := lgc.authorizer.ListAuthorizedResources(kit.Ctx, kit.Header, authInput)
	if err != nil {
		blog.ErrorJSON("list authorized %s failed, options: %s, err: %v, rid: %s", typ, authInput, err, kit.Rid)
		return nil, err
	}

	accountIDList := make([]int64, 0)
	for _, id := range authList {
		subscriptionID, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			blog.Errorf("parse account id(%s) failed, err: %v, rid: %s", id, err, kit.Rid)
			return nil, err
		}
		accountIDList = append(accountIDList, subscriptionID)
	}
	return accountIDList, nil
}
