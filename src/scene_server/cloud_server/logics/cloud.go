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
	"context"
	"fmt"
	"sync"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/cloud_server/cloudvendor"
	ccom "configcenter/src/scene_server/cloud_server/common"
)

func (lgc *Logics) AccountVerify(conf metadata.CloudAccountConf) error {
	client, err := cloudvendor.GetVendorClient(conf)
	if err != nil {
		blog.Errorf("AccountVerify GetVendorClient failed, AccountID:%d, err:%s", conf.AccountID, err.Error())
		return err
	}

	_, err = client.GetRegions(nil)
	if err != nil {
		blog.Errorf("AccountVerify GetRegions failed, AccountID:%d, err:%s", conf.AccountID, err.Error())
		return err
	}

	return nil
}

// 获取地域信息
func (lgc *Logics) GetRegionsInfo(conf metadata.CloudAccountConf, withHostCount bool) ([]metadata.SyncRegion, error) {
	client, err := cloudvendor.GetVendorClient(conf)
	if err != nil {
		blog.Errorf("GetRegionsInfo GetVendorClient failed, AccountID:%d, err:%s", conf.AccountID, err.Error())
		return nil, err
	}

	regionSet, err := client.GetRegions(nil)
	if err != nil {
		blog.Errorf("GetRegionsInfo GetRegions err:%s", err.Error())
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
					blog.Errorf("GetVpcHostCnt GetInstances failed, AccountID:%d, err:%s", conf.AccountID, err.Error())
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
	for i, _ := range regionSet {
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

// 获取某地域下的vpc详情和主机数
func (lgc *Logics) GetVpcHostCntInOneRegion(conf metadata.CloudAccountConf, region string) (*metadata.VpcHostCntResult, error) {
	client, err := cloudvendor.GetVendorClient(conf)
	if err != nil {
		blog.Errorf("GetVpcHostCntInOneRegion GetVendorClient failed, AccountID:%d, err:%s", conf.AccountID, err.Error())
		return nil, err
	}

	// 获取该地域下的vpc详情
	vpcsInfo, err := client.GetVpcs(region, nil)
	if err != nil {
		blog.Errorf("GetVpcHostCntInOneRegion GetVpcs err:%s", err.Error())
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
	vpcHostCnt, err := lgc.GetVpcHostCnt(conf, option)
	if err != nil {
		blog.Errorf("GetVpcHostCntInOneRegion GetVpcHostCnt failed, AccountID:%d, err:%s", conf.AccountID, err.Error())
		return nil, err
	}

	for i, _ := range vpcsInfo.VpcSet {
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

// 获取多个vpc对应的主机数
func (lgc *Logics) GetVpcHostCnt(conf metadata.CloudAccountConf, option metadata.SearchVpcHostCntOption) (map[string]int64, error) {
	client, err := cloudvendor.GetVendorClient(conf)
	if err != nil {
		blog.Errorf("GetVpcHostCnt GetVendorClient failed, AccountID:%d, err:%s", conf.AccountID, err.Error())
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
			count, err := client.GetInstancesTotalCnt(regionVpc.Region, &ccom.RequestOpt{
				Filters: []*ccom.Filter{&ccom.Filter{ccom.StringPtr("vpc-id"), ccom.StringPtrs([]string{regionVpc.VpcID})}},
			})
			if err != nil {
				blog.Errorf("GetVpcHostCnt GetInstances failed, AccountID:%d, err:%s", conf.AccountID, err.Error())
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

// 获取需要同步的云主机资源信息
func (lgc *Logics) GetCloudHostResource(conf metadata.CloudAccountConf, syncVpcs []metadata.VpcSyncInfo) (*metadata.CloudHostResource, error) {
	client, err := cloudvendor.GetVendorClient(conf)
	if err != nil {
		blog.Errorf("GetCloudHostResource GetVendorClient failed, AccountID:%d, err:%s", conf.AccountID, err.Error())
		return nil, err
	}

	blog.V(4).Infof("GetCloudHostResource syncVpcs %#v", syncVpcs)

	// 不再同步已经被销毁的vpc
	for i, vpc := range syncVpcs {
		if vpc.Destroyed {
			syncVpcs = append(syncVpcs[:i], syncVpcs[i:]...)
		}
	}

	vpcHostDetail := make(map[string][]*metadata.Instance)
	hostDetailChan := make(chan []*metadata.Instance, 10)
	destroyedVpcs := make(map[string]bool)
	destroyedVpcsChan := make(chan string, 10)
	errs := make([]error, 0)
	errChan := make(chan error, 10)
	var wg, wg2, wg3, wg4 sync.WaitGroup
	// 并发请求获取被销毁的vpc数据和没被销毁的vpc下主机实例详情
	for _, vpc := range syncVpcs {
		wg.Add(1)
		go func(vpc metadata.VpcSyncInfo) {
			defer wg.Done()
			vpcInfo, err := client.GetVpcs(vpc.Region, &ccom.RequestOpt{
				Filters: []*ccom.Filter{&ccom.Filter{ccom.StringPtr("vpc-id"), ccom.StringPtrs([]string{vpc.VpcID})}},
				Limit:   ccom.Int64Ptr(ccom.MaxLimit),
			})
			if err != nil {
				blog.Errorf("GetCloudHostResource GetVpcs failed, AccountID:%d, err:%s", conf.AccountID, err.Error())
				errChan <- err
				return
			}
			if len(vpcInfo.VpcSet) == 0 {
				destroyedVpcsChan <- vpc.VpcID
				return
			}
			instancesInfo, err := client.GetInstances(vpc.Region, &ccom.RequestOpt{
				Filters: []*ccom.Filter{&ccom.Filter{ccom.StringPtr("vpc-id"), ccom.StringPtrs([]string{vpc.VpcID})}},
				Limit:   ccom.Int64Ptr(ccom.MaxLimit),
			})
			if err != nil {
				blog.Errorf("GetCloudHostResource GetInstances failed, AccountID:%d, err:%s", conf.AccountID, err.Error())
				errChan <- err
				return
			}
			blog.V(4).Infof("GetCloudHostResource vpc-id:%s, instances count %#v", vpc.VpcID, instancesInfo.Count)
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
		return nil , errs[0]
	}

	result := new(metadata.CloudHostResource)

	for i, _ := range syncVpcs {
		vpc := syncVpcs[i]
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

// 获取云厂商账户配置
func (lgc *Logics) GetCloudAccountConf(accountID int64) (*metadata.CloudAccountConf, error) {
	option := &metadata.SearchCloudOption{Condition: mapstr.MapStr{common.BKCloudAccountID: accountID}}
	result, err := lgc.CoreAPI.CoreService().Cloud().SearchAccountConf(context.Background(), ccom.GetHeader(), option)
	if err != nil {
		blog.Errorf("SearchAccountConf failed, accountID: %v, err: %s", accountID, err.Error())
		return nil, err
	}
	if len(result.Info) == 0 {
		blog.Errorf("GetCloudAccountConf failed, accountID: %v is not exist", accountID)
		return nil, fmt.Errorf("GetAccountConf failed, accountID: %v is not exist", accountID)
	}

	return &result.Info[0], nil
}
