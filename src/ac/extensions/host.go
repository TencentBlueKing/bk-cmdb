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

package extensions

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"configcenter/src/ac/iam"
	"configcenter/src/ac/meta"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

func (am *AuthManager) constructHostFromSearchResult(ctx context.Context, header http.Header, rawData []mapstr.MapStr) ([]HostSimplify, error) {
	rid := util.ExtractRequestIDFromContext(ctx)

	hostIDs := make([]int64, 0)
	hosts := make([]HostSimplify, 0)
	for _, cls := range rawData {
		host := HostSimplify{}
		if _, err := host.Parse(cls); err != nil {
			return nil, fmt.Errorf("get hosts by object failed, err: %+v", err)
		}
		hosts = append(hosts, host)
		hostIDs = append(hostIDs, host.BKHostIDField)
	}

	// inject business,set,module info to HostSimplify
	query := &metadata.HostModuleRelationRequest{
		HostIDArr: hostIDs,
		Fields:    []string{common.BKAppIDField, common.BKSetIDField, common.BKModuleIDField, common.BKHostIDField},
	}
	hostModuleResult, err := am.clientSet.CoreService().Host().GetHostModuleRelation(ctx, header, query)
	if err != nil {
		err = fmt.Errorf("get host:%+v layer failed, err: %+v", hostIDs, err)
		return nil, err
	}
	if len(rawData) == 0 {
		err = fmt.Errorf("get host:%+v layer failed, get host module config by host id not found, maybe hostID invalid", hostIDs)
		return nil, err
	}
	hostModuleMap := map[int64]HostSimplify{}
	for _, cls := range hostModuleResult.Data.Info {
		host := HostSimplify{
			BKAppIDField:    cls.AppID,
			BKModuleIDField: cls.ModuleID,
			BKSetIDField:    cls.SetID,
			BKHostIDField:   cls.HostID,
		}
		hostModuleMap[host.BKHostIDField] = host
	}
	for idx, host := range hosts {
		hostModule, exist := hostModuleMap[host.BKHostIDField]
		if !exist {
			return nil, fmt.Errorf("hostID:%+d doesn't exist in any module", host.BKHostIDField)
		}
		hosts[idx].BKAppIDField = hostModule.BKAppIDField
		hosts[idx].BKSetIDField = hostModule.BKSetIDField
		hosts[idx].BKModuleIDField = hostModule.BKModuleIDField
	}
	blog.V(9).Infof("hosts: %+v, rid: %s", hosts, rid)
	return hosts, nil
}

func (am *AuthManager) collectHostByHostIDs(ctx context.Context, header http.Header, hostIDs ...int64) ([]HostSimplify, error) {
	rid := util.ExtractRequestIDFromContext(ctx)

	// unique ids so that we can be aware of invalid id if query result length not equal ids's length
	hostIDs = util.IntArrayUnique(hostIDs)
	hosts := make([]mapstr.MapStr, 0)
	count := -1
	for offset := 0; count == -1 || offset < count; offset += common.BKMaxRecordsAtOnce {
		cond := metadata.QueryCondition{
			Fields:    []string{common.BKHostIDField, common.BKHostNameField, common.BKHostInnerIPField},
			Condition: condition.CreateCondition().Field(common.BKHostIDField).In(hostIDs).ToMapStr(),
			Page: metadata.BasePage{
				Sort:  "",
				Limit: common.BKMaxRecordsAtOnce,
				Start: offset,
			},
		}
		result, err := am.clientSet.CoreService().Instance().ReadInstance(ctx, header, common.BKInnerObjIDHost, &cond)
		if err != nil {
			blog.V(3).Infof("get hosts by id failed, err: %+v, rid: %s", err, rid)
			return nil, fmt.Errorf("get hosts by id failed, err: %+v", err)
		}
		hosts = append(hosts, result.Data.Info...)
		count = result.Data.Count
	}
	return am.constructHostFromSearchResult(ctx, header, hosts)
}

func (am *AuthManager) MakeResourcesByHosts(ctx context.Context, header http.Header, action meta.Action, hosts ...HostSimplify) ([]meta.ResourceAttribute, error) {
	rid := util.ExtractRequestIDFromContext(ctx)

	businessIDs := make([]int64, 0)
	for _, host := range hosts {
		businessIDs = append(businessIDs, host.BKAppIDField)
	}
	businessIDs = util.IntArrayUnique(businessIDs)
	bizIDCorrectMap := make(map[int64]int64)
	resPoolBizID, err := am.getResourcePoolBusinessID(ctx, header)
	if err != nil {
		return nil, fmt.Errorf("correct host related business id failed, err: %+v", err)
	}
	for _, businessID := range businessIDs {
		if businessID == resPoolBizID {
			// if this is resource pool business, then change the biz id to 0, so that it
			// represent global resources
			bizIDCorrectMap[businessID] = 0
		} else {
			bizIDCorrectMap[businessID] = businessID
		}
	}

	resources := make([]meta.ResourceAttribute, 0)
	for _, host := range hosts {
		resource := meta.ResourceAttribute{
			Basic: meta.Basic{
				Action:     action,
				Type:       meta.HostInstance,
				Name:       host.BKHostInnerIPField,
				InstanceID: host.BKHostIDField,
			},
			SupplierAccount: util.GetOwnerID(header),
			BusinessID:      bizIDCorrectMap[host.BKAppIDField],
		}
		resources = append(resources, resource)
	}

	blog.V(9).Infof("host resources for iam: %+v, rid: %s", resources, rid)
	return resources, nil
}

func (am *AuthManager) AuthorizeByHosts(ctx context.Context, header http.Header, action meta.Action, hosts ...HostSimplify) error {
	if !am.Enabled() {
		return nil
	}

	if len(hosts) == 0 {
		return nil
	}

	// make auth resources
	resources, err := am.MakeResourcesByHosts(ctx, header, action, hosts...)
	if err != nil {
		return fmt.Errorf("make host resources failed, err: %+v", err)
	}
	return am.batchAuthorize(ctx, header, resources...)
}

func (am *AuthManager) GenHostBatchNoPermissionResp(ctx context.Context, header http.Header, action meta.Action, hostIDs []int64) (*metadata.BaseResp, error) {
	hosts, err := am.collectHostByHostIDs(ctx, header, hostIDs...)
	if err != nil {
		return nil, err
	}
	resPoolBizID, err := am.getResourcePoolBusinessID(ctx, header)
	if err != nil {
		return nil, err
	}
	bizHosts := make([][]metadata.IamResourceInstance, 0)
	resourceHosts := make([][]metadata.IamResourceInstance, 0)
	var bizID int64
	for _, host := range hosts {
		if host.BKAppIDField == resPoolBizID {
			resourceHosts = append(resourceHosts, []metadata.IamResourceInstance{{
				Type: string(iam.SysHostRscPoolDirectory),
				ID:   strconv.FormatInt(host.BKModuleIDField, 10),
			}, {
				Type: string(iam.Host),
				ID:   strconv.FormatInt(host.BKHostIDField, 10),
			}})
		} else {
			bizID = host.BKAppIDField
			bizHosts = append(bizHosts, []metadata.IamResourceInstance{{
				Type: string(iam.Business),
				ID:   strconv.FormatInt(host.BKAppIDField, 10),
			}, {
				Type: string(iam.Host),
				ID:   strconv.FormatInt(host.BKHostIDField, 10),
			}})
		}
	}

	permission := &metadata.IamPermission{SystemID: iam.SystemIDCMDB}
	if len(bizHosts) > 0 {
		action, err := iam.ConvertResourceAction(meta.HostInstance, action, bizID)
		if err != nil {
			return nil, err
		}
		permission.Actions = append(permission.Actions, metadata.IamAction{
			ID: string(action),
			RelatedResourceTypes: []metadata.IamResourceType{{
				SystemID:  iam.SystemIDCMDB,
				Type:      string(iam.Host),
				Instances: bizHosts,
			}},
		})
	}
	if len(resourceHosts) > 0 {
		action, err := iam.ConvertResourceAction(meta.HostInstance, action, 0)
		if err != nil {
			return nil, err
		}
		permission.Actions = append(permission.Actions, metadata.IamAction{
			ID: string(action),
			RelatedResourceTypes: []metadata.IamResourceType{{
				SystemID:  iam.SystemIDCMDB,
				Type:      string(iam.Host),
				Instances: resourceHosts,
			}},
		})
	}
	resp := metadata.NewNoPermissionResp(permission)
	return &resp, nil
}

func (am *AuthManager) GenEditBizHostNoPermissionResp(ctx context.Context, header http.Header, hostIDs []int64) (*metadata.BaseResp, error) {
	hosts, err := am.collectHostByHostIDs(ctx, header, hostIDs...)
	if err != nil {
		return nil, err
	}
	instances := make([][]metadata.IamResourceInstance, len(hosts))
	for index, host := range hosts {
		instances[index] = []metadata.IamResourceInstance{{
			Type: string(iam.Business),
			ID:   strconv.FormatInt(host.BKAppIDField, 10),
		}, {
			Type: string(iam.Host),
			ID:   strconv.FormatInt(host.BKHostIDField, 10),
		}}
	}
	permission := &metadata.IamPermission{
		SystemID: iam.SystemIDCMDB,
		Actions: []metadata.IamAction{{
			ID: string(iam.EditBusinessHost),
			RelatedResourceTypes: []metadata.IamResourceType{{
				SystemID:  iam.SystemIDCMDB,
				Type:      string(iam.Host),
				Instances: instances,
			}},
		}},
	}
	resp := metadata.NewNoPermissionResp(permission)
	return &resp, nil
}

func (am *AuthManager) AuthorizeByHostsIDs(ctx context.Context, header http.Header, action meta.Action, hostIDs ...int64) error {
	rid := util.ExtractRequestIDFromContext(ctx)

	if !am.Enabled() {
		return nil
	}
	if am.SkipReadAuthorization && (action == meta.Find || action == meta.FindMany) {
		blog.V(4).Infof("skip authorization for reading, hosts: %+v, rid: %s", hostIDs, rid)
		return nil
	}

	if len(hostIDs) == 0 {
		return nil
	}
	hosts, err := am.collectHostByHostIDs(ctx, header, hostIDs...)
	if err != nil {
		return fmt.Errorf("authorize hosts failed, get hosts by id failed, err: %+v, rid: %s", err, rid)
	}
	return am.AuthorizeByHosts(ctx, header, action, hosts...)
}

func (am *AuthManager) AuthorizeCreateHost(ctx context.Context, header http.Header, bizID int64) error {
	if !am.Enabled() {
		return nil
	}

	return am.AuthorizeResourceCreate(ctx, header, bizID, meta.HostInstance)
}
