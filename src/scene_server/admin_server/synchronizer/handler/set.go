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

package handler

import (
	"context"
	"fmt"
	"strings"

	"configcenter/src/auth/authcenter"
	authmeta "configcenter/src/auth/meta"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/admin_server/synchronizer/meta"
	"configcenter/src/scene_server/admin_server/synchronizer/utils"
)

// HandleSetSync do sync set of one business
func (ih *IAMHandler) HandleSetSync(task *meta.WorkRequest) error {
	businessSimplify := task.Data.(meta.BusinessSimplify)
	header := utils.NewAPIHeaderByBusiness(&businessSimplify)
	coreService := ih.CoreAPI.CoreService()

	// step1 get instances by business from core service
	cond := condition.CreateCondition()
	cond.Field(common.BKAppIDField).Eq(businessSimplify.BKAppIDField)
	query := &metadata.QueryCondition{
		Fields:    []string{common.BKAppIDField, common.BKSetIDField, common.BKSetNameField},
		Condition: cond.ToMapStr(),
		Limit:     metadata.SearchLimit{Limit: common.BKNoLimit},
	}
	instances, err := coreService.Instance().ReadInstance(context.Background(), *header, common.BKInnerObjIDSet, query)
	if err != nil {
		blog.Errorf("get set:%+v by businessID:%d failed, err: %+v", businessSimplify.BKAppIDField, err)
		return fmt.Errorf("get set by businessID:%d failed, err: %+v", businessSimplify.BKAppIDField, err)
	}

	if len(instances.Data.Info) == 0 {
		blog.V(2).Infof("business: %d has no instances, skip synchronize sets.", businessSimplify.BKAppIDField)
		return nil
	}

	// extract sets
	setArr := make([]meta.SetSimplify, 0)
	for _, instance := range instances.Data.Info {
		val, exist := instance[common.BKSetIDField]
		if exist == false {
			continue
		}
		idVal, err := util.GetInt64ByInterface(val)
		if err != nil {
			blog.V(2).Infof("synchronize task skip set:%+v, as parse setID field failed, err: %+v", instance, err)
			continue
		}

		var nameVal interface{}
		var name string
		nameVal, exist = instance[common.BKSetNameField]
		if exist == false {
			name = fmt.Sprintf("set:%d", idVal)
		} else {
			name = util.GetStrByInterface(nameVal)
		}
		setArr = append(setArr, meta.SetSimplify{
			BKAppIDField:   businessSimplify.BKAppIDField,
			BKSetIDField:   idVal,
			BKSetNameField: name,
		})
	}

	blog.V(4).Infof("list sets by business:%d result: %+v", businessSimplify.BKAppIDField, setArr)

	// step2 get set by business from iam
	rs := &authmeta.ResourceAttribute{
		Basic: authmeta.Basic{
			Type: authmeta.ModelSet,
		},
		SupplierAccount: "",
		BusinessID:      businessSimplify.BKAppIDField,
		Layers: make([]authmeta.Item, 0),
		// iam don't support set layers yet.
		// Layers:          []authmeta.Item{{Type: authmeta.Business, InstanceID: businessID,}},
	}
	resultResources, err := ih.Authorizer.ListResources(context.Background(), rs)
	if err != nil {
		blog.Errorf("synchronize set instance failed, ListResources failed, err: %+v", err)
		return err
	}
	// filter set from topo model instances
	realResources := make([]authmeta.BackendResource, 0)
	for _, iamResources := range resultResources {
		if strings.Contains(iamResources[len(iamResources)-1].ResourceID, "set") {
			realResources = append(realResources, iamResources)
		}
	}
	blog.InfoJSON("realResources is: %s", realResources)

	// init key:hit map for
	iamResourceKeyMap := map[string]int{}
	for _, iamResource := range realResources {
		key := generateIAMResourceKey(iamResource)
		// init hit count 0
		iamResourceKeyMap[key] = 0
	}

	// step6 register set not exist in iam
	// step5 diff step2 and step4 result
	scope := authcenter.ScopeInfo{}
	needRegister := make([]authmeta.ResourceAttribute, 0)
	for _, set := range setArr {
		resource := authmeta.ResourceAttribute{
			Basic: authmeta.Basic{
				Type:       authmeta.ModelSet,
				Name:       set.BKSetNameField,
				InstanceID: set.BKSetIDField,
			},
			SupplierAccount: "",
			BusinessID:      businessSimplify.BKAppIDField,
			// Layers:          layer[0:1],
		}
		targetResource, err := ih.Authorizer.DryRunRegisterResource(context.Background(), resource)
		if err != nil {
			blog.Errorf("synchronize set instance failed, dry run register resource failed, err: %+v", err)
			return err
		}
		if len(targetResource.Resources) != 1 {
			blog.Errorf("synchronize instance:%+v failed, dry run register result is: %+v", resource, targetResource)
			continue
		}
		scope.ScopeID = targetResource.Resources[0].ScopeID
		scope.ScopeType = targetResource.Resources[0].ScopeType
		resourceKey := generateCMDBResourceKey(&targetResource.Resources[0])
		_, exist := iamResourceKeyMap[resourceKey]
		if exist {
			iamResourceKeyMap[resourceKey]++
		} else {
			needRegister = append(needRegister, resource)
		}
	}
	blog.V(5).Infof("iamResourceKeyMap: %+v", iamResourceKeyMap)
	blog.V(5).Infof("needRegister: %+v", needRegister)
	if len(needRegister) > 0 {
		blog.V(2).Infof("sychronizer register resource that only in cmdb, resources: %+v", needRegister)
		err = ih.Authorizer.RegisterResource(context.Background(), needRegister...)
		if err != nil {
			blog.ErrorJSON("sychronizer register resource that only in cmdb failed, resources: %s, err: %+v", needRegister, err)
		}
	}

	// step7 deregister resource id that hasn't been hit
	needDeregister := make([]authmeta.BackendResource, 0)
	for _, iamResource := range realResources {
		resourceKey := generateIAMResourceKey(iamResource)
		if iamResourceKeyMap[resourceKey] == 0 {
			needDeregister = append(needDeregister, iamResource)
		}
	}
	blog.V(5).Infof("needDeregister: %+v", needDeregister)
	if len(needDeregister) != 0 {
		blog.V(2).Infof("sychronizer deregister resource that only in iam, resources: %+v", needDeregister)
		err = ih.Authorizer.RawDeregisterResource(context.Background(), scope, needDeregister...)
		if err != nil {
			blog.ErrorJSON("sychronizer deregister resource that only in iam failed, resources: %s, err: %+v", needDeregister, err)
		}
	}

	return nil
}
